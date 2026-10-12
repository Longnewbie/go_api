package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"time"
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/utils"
	"user-management-api/pkg/cache"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	cache cache.RedisCacheService
}

type Claims struct {
	jwt.RegisteredClaims
}

type EncryptPayLoad struct {
	UserUUID string `json:"user_uuid"`
	Email    string `json:"email"`
	Role     int32  `json:"role"`
}

type RefreshToken struct {
	Token     string    `json:"token"`
	UserUUID  string    `json:"user_uuid"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `json:"revoked"`
}

var (
	jwtSecretKey  = []byte(utils.GetEnv("JWT_SECRET", "abcxyz"))
	jwtEncryptKey = []byte(utils.GetEnv("JWT_ENCRYPT_KEY", "12345678901234567890123456789004"))
)

const (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 1 * 24 * time.Hour
)

func NewJWTService(cache cache.RedisCacheService) TokenService {
	return &JWTService{
		cache: cache,
	}
}

/** access token */
func (js *JWTService) GenerateAccessToken(user sqlc.User) (string, error) {
	// 1. create payload from user data
	payload := &EncryptPayLoad{
		UserUUID: user.UserUuid.String(),
		Email:    user.UserEmail,
		Role:     user.UserLevel,
	}

	// 2. marshal payload to JSON
	rawData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	// 3. encrypt data
	encrypted, err := utils.EncryptAES(rawData, jwtEncryptKey)
	if err != nil {
		return "", err
	}

	// 4. create JWT claims
	claims := jwt.MapClaims{
		"data": encrypted,
		"jti":  uuid.NewString(),
		"exp":  time.Now().Add(AccessTokenTTL).Unix(),
		"iat":  time.Now().Unix(),
		"iss":  "golang-gialong04",
	}

	// 5. sign JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecretKey)
}

func (js *JWTService) ParseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	// 1. parse token
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return jwtSecretKey, nil
	})

	// 2. validate token
	if err != nil || !token.Valid {
		return nil, nil, utils.NewError("Invalid token", utils.ErrCodeUnauthorized)
	}

	// 3. get claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, utils.NewError("Invalid token", utils.ErrCodeUnauthorized)
	}

	return token, claims, nil
}

func (js *JWTService) DecryptAccessTokenPayload(tokenString string) (*EncryptPayLoad, error) {
	// 1. parse token and get claims
	_, claims, err := js.ParseToken(tokenString)
	if err != nil {
		return nil, utils.WrapError(err, "Cannot parse token", utils.ErrCodeInternal)
	}

	// 2. decrypt data from claims
	encryptedData, ok := claims["data"].(string)
	if !ok {
		return nil, utils.NewError("Encoded data not found", utils.ErrCodeUnauthorized)
	}

	decryptedBytes, err := utils.DecryptAES(encryptedData, jwtEncryptKey)
	if err != nil {
		return nil, utils.WrapError(err, "Cannot decode data", utils.ErrCodeInternal)
	}

	// 3. unmarshal payload to struct
	var payload EncryptPayLoad
	if err := json.Unmarshal(decryptedBytes, &payload); err != nil {
		return nil, utils.WrapError(err, "Invalid data format", utils.ErrCodeInternal)
	}

	return &payload, nil
}

/** refresh token */
func (js *JWTService) GenerateRefreshToken(user sqlc.User) (RefreshToken, error) {
	// 1. generate random token
	tokenBytes := make([]byte, 32)

	if _, err := rand.Read(tokenBytes); err != nil {
		return RefreshToken{}, err
	}

	// 2. encode token
	token := base64.URLEncoding.EncodeToString(tokenBytes)

	// 3. create refresh token
	return RefreshToken{
		Token:     token,
		UserUUID:  user.UserUuid.String(),
		ExpiresAt: time.Now().Add(RefreshTokenTTL),
		Revoked:   false,
	}, nil
}

func (js *JWTService) StoreRefreshToken(token RefreshToken) error {
	// create cache key and store the refresh token in Redis
	cacheKey := "refresh_token:" + token.Token
	return js.cache.Set(cacheKey, token, RefreshTokenTTL)
}

func (js *JWTService) ValidateRefreshToken(token string) (RefreshToken, error) {
	// Get refresh token from redis
	cacheKey := "refresh_token:" + token

	var refreshToken RefreshToken

	err := js.cache.Get(cacheKey, &refreshToken)
	if err != nil || refreshToken.Revoked || refreshToken.ExpiresAt.Before(time.Now()) {
		return RefreshToken{}, utils.WrapError(err, "Cannot get refresh token", utils.ErrCodeInternal)
	}

	return refreshToken, nil
}

func (js *JWTService) RevokeRefreshToken(token string) error {
	// Get refresh token from redis
	cacheKey := "refresh_token:" + token

	var refreshToken RefreshToken

	err := js.cache.Get(cacheKey, &refreshToken)
	if err != nil {
		return utils.WrapError(err, "Cannot get refresh token", utils.ErrCodeInternal)
	}

	// revoked the refresh token
	refreshToken.Revoked = true

	return js.cache.Set(cacheKey, refreshToken, time.Until(refreshToken.ExpiresAt))
}
