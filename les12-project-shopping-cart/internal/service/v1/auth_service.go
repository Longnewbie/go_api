package v1service

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/repository"
	"user-management-api/internal/utils"
	"user-management-api/pkg/auth"
	"user-management-api/pkg/cache"
	"user-management-api/pkg/logger"
	"user-management-api/pkg/mail"
	"user-management-api/pkg/rabbitmq"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/time/rate"
)

type authService struct {
	userRepo     repository.UserRepository
	tokenService auth.TokenService
	cacheService cache.RedisCacheService
	mailService  mail.EmailProviderService
	rabbitmq     rabbitmq.RabbitMQService
}

type LoginAttempt struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	mu              sync.Mutex
	clients         = make(map[string]*LoginAttempt)
	LoginAttemptTTL = 5 * time.Minute
	MaxLoginAttempt = 5
)

func NewAuthService(repo repository.UserRepository, tokenService auth.TokenService, cacheService *redis.Client, mailService mail.EmailProviderService, rabbitmqService rabbitmq.RabbitMQService) AuthService {
	return &authService{
		userRepo:     repo,
		tokenService: tokenService,
		cacheService: cache.NewRedisCacheService(cacheService),
		mailService:  mailService,
		rabbitmq:     rabbitmqService,
	}
}

// get client IP address from request context
func (as *authService) getClientIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()

	if ip == "" {
		ip = ctx.Request.RemoteAddr
	}

	return ip
}

// getLoginAttempt returns the rate limiter for the given IP address, creating a new one if necessary
func (as *authService) getLoginAttempt(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	client, exists := clients[ip]

	if !exists {
		limiter := rate.NewLimiter(rate.Limit(float32(MaxLoginAttempt)/float32(LoginAttemptTTL.Seconds())), MaxLoginAttempt) // 1 req/ 1m, bursts = 5
		newClient := &LoginAttempt{limiter, time.Now()}
		clients[ip] = newClient
		return limiter
	}

	client.lastSeen = time.Now()
	return client.limiter
}

// checkLoginAttempt checks if the given IP address has exceeded the maximum number of login attempts
func (as *authService) checkLoginAttempt(ip string) error {
	limiter := as.getLoginAttempt(ip)
	if !limiter.Allow() {
		return utils.NewError("Too many login attempt. Please retry again later", utils.ErrCodeTooManyRequests)
	}

	return nil
}

// cleanupClients removes clients that have not been seen for a certain amount of time
func (as *authService) CleanupClients(ip string) {
	mu.Lock()
	defer mu.Unlock()
	delete(clients, ip)
}

func (as *authService) Login(ctx *gin.Context, email, password string) (string, string, int, error) {
	context := ctx.Request.Context()

	ip := as.getClientIP(ctx)
	if err := as.checkLoginAttempt(ip); err != nil {
		return "", "", 0, err
	}

	email = utils.NormalizeString(email)

	// get user by email
	user, err := as.userRepo.GetByEmail(context, email)
	if err != nil {
		as.getLoginAttempt(ip)
		return "", "", 0, utils.NewError("Invalid email or password", utils.ErrCodeUnauthorized)
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(password)); err != nil {
		as.getLoginAttempt(ip)
		return "", "", 0, utils.NewError("Invalid email or password", utils.ErrCodeUnauthorized)
	}

	// generate access_token
	accessToken, err := as.tokenService.GenerateAccessToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "Unable to create access token", utils.ErrCodeInternal)
	}

	// generate refresh_token
	refreshToken, err := as.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "Unable to create refresh token", utils.ErrCodeInternal)
	}

	// Save the refresh token to Redis
	if err := as.tokenService.StoreRefreshToken(refreshToken); err != nil {
		return "", "", 0, utils.WrapError(err, "Cannot save refresh token", utils.ErrCodeInternal)
	}

	as.CleanupClients(ip)

	return accessToken, refreshToken.Token, int(auth.AccessTokenTTL.Seconds()), nil
}

func (as *authService) Logout(ctx *gin.Context, refreshToken string) error {
	// 1. Get the access token from the Authorization header
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return utils.NewError("Missing authorization header", utils.ErrCodeUnauthorized)
	}

	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	_, claims, err := as.tokenService.ParseToken(accessToken)
	if err != nil {
		return utils.NewError("Invalid access token", utils.ErrCodeUnauthorized)
	}

	// 2. Add the access token's jti to the blacklist in Redis with the same expiration time as the token
	if jti, ok := claims["jti"].(string); ok {
		expUnix, _ := claims["exp"].(float64)
		exp := time.Unix(int64(expUnix), 0)
		key := "blacklist:" + jti
		ttl := time.Until(exp)
		as.cacheService.Set(key, "revoked", ttl)
	}

	// 3. Validate the refresh token and check if it exists in Redis
	_, err = as.tokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return utils.NewError("Refresh token is invalid or revoked", utils.ErrCodeUnauthorized)
	}

	// 4. Revoke the refresh token
	if err := as.tokenService.RevokeRefreshToken(refreshToken); err != nil {
		return utils.WrapError(err, "Unable to revoke token", utils.ErrCodeInternal)
	}

	return nil
}

func (as *authService) RefreshToken(ctx *gin.Context, refreshTokenString string) (string, string, int, error) {
	context := ctx.Request.Context()

	// 1. Check the refresh token. Return the user_uuid
	token, err := as.tokenService.ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return "", "", 0, utils.NewError("Refresh token is invalid or revoked", utils.ErrCodeUnauthorized)
	}

	// 2. Retrieving user information from user_uuid
	userUuid, _ := uuid.Parse(token.UserUUID)
	user, err := as.userRepo.GetByUUID(context, userUuid)
	if err != nil {
		return "", "", 0, utils.NewError("User not found", utils.ErrCodeNotFound)
	}

	// 3. Create a new access token
	accessToken, err := as.tokenService.GenerateAccessToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "Unable to create access token", utils.ErrCodeInternal)
	}

	// 4. Create a new refresh token
	refreshToken, err := as.tokenService.GenerateRefreshToken(user)
	if err != nil {
		return "", "", 0, utils.WrapError(err, "Unable to create refresh token", utils.ErrCodeInternal)
	}

	// 5. Disable old refresh token
	if err := as.tokenService.RevokeRefreshToken(refreshTokenString); err != nil {
		return "", "", 0, utils.WrapError(err, "Unable to revoke token", utils.ErrCodeInternal)
	}

	// 6. Save the new refresh token
	if err := as.tokenService.StoreRefreshToken(refreshToken); err != nil {
		return "", "", 0, utils.WrapError(err, "Cannot save refresh token", utils.ErrCodeInternal)
	}

	return accessToken, refreshToken.Token, int(auth.AccessTokenTTL.Seconds()), nil
}

func (as *authService) ForgotPassword(ctx *gin.Context, email string) error {
	context := ctx.Request.Context()

	email = utils.NormalizeString(email)

	// Check if the user has requested a password reset recently
	rateLimitKey := fmt.Sprintf("reset:rateLimit:%s", email)
	if exists, err := as.cacheService.Exists(rateLimitKey); err == nil && exists {
		return utils.NewError("Please wait before requesting another password reset", utils.ErrCodeTooManyRequests)
	}

	// get user by email
	user, err := as.userRepo.GetByEmail(context, email)
	if err != nil {
		return utils.NewError("Email not found", utils.ErrCodeNotFound)
	}

	// generate reset token
	token, err := utils.GenerateRandomString(16)
	if err != nil {
		return utils.NewError("Failed to generate reset token", utils.ErrCodeInternal)
	}

	// store reset token in Redis with a 1-hour expiration
	err = as.cacheService.Set("reset:"+token, user.UserUuid, time.Hour)
	if err != nil {
		return utils.NewError("Failed to store reset token", utils.ErrCodeInternal)
	}

	// store rate limit in Redis with a 5-minute expiration
	err = as.cacheService.Set(rateLimitKey, "1", 5*time.Minute)
	if err != nil {
		return utils.NewError("Failed to store rate limit reset password", utils.ErrCodeInternal)
	}

	resetLink := fmt.Sprintf("http://abc.com/view-to-reset-password?token=%s", token)

	logger.Log.Info().Msg(resetLink)

	// send reset link to email
	mailContent := &mail.Email{
		To: []mail.Address{
			{Email: email},
		},
		Subject: "Password Reset Request",
		Text:    fmt.Sprintf("Hi %s, \n\n You requested to reset your password. Please click link below to reset it:\n%s\n\n The link will expire in 1 hour. \n\n Best regard, \nGiaLong04", user.UserEmail, resetLink),
	}

	// send email to queue
	if err := as.rabbitmq.Publish(context, "auth_email_queue", mailContent); err != nil {
		return utils.NewError("Failed to send password reset email", utils.ErrCodeInternal)
	}

	// send email
	// if err := as.mailService.SendMail(context, mailContent); err != nil {
	// 	return utils.NewError("Failed to send password reset email", utils.ErrCodeInternal)
	// }

	return nil
}

func (as *authService) ResetPassword(ctx *gin.Context, token, password string) error {
	context := ctx.Request.Context()

	var userUUIDStr string

	// get user_uuid from token
	err := as.cacheService.Get("reset:"+token, &userUUIDStr)
	if err == redis.Nil || userUUIDStr == "" {
		return utils.NewError("Invalid or expired token", utils.ErrCodeNotFound)
	}

	if err != nil {
		return utils.NewError("Failed to get reset token", utils.ErrCodeInternal)
	}

	// user_uuid is valid, update the password
	userUuid, err := uuid.Parse(userUUIDStr)
	if err != nil {
		return utils.WrapError(err, "Uuid is invalid", utils.ErrCodeInternal)
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return utils.WrapError(err, "Failed to hash password", utils.ErrCodeInternal)
	}

	var input = sqlc.UpdatePasswordParams{
		UserPassword: string(hashPassword),
		UserUuid:     userUuid,
	}

	_, err = as.userRepo.UpdatePassword(context, input)
	if err != nil {
		return utils.NewError("Failed to update new password", utils.ErrCodeInternal)
	}

	// clear cache redis
	as.cacheService.Clear("reset:" + token)

	return nil
}
