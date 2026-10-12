package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-management-api/internal/config"
	"user-management-api/internal/db"
	"user-management-api/internal/db/sqlc"
	"user-management-api/internal/routes"
	"user-management-api/internal/utils"
	"user-management-api/internal/validation"
	"user-management-api/pkg/auth"
	"user-management-api/pkg/cache"
	"user-management-api/pkg/logger"
	"user-management-api/pkg/mail"
	"user-management-api/pkg/rabbitmq"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Module interface {
	Routes() routes.Route
}

type Application struct {
	config  *config.Config
	router  *gin.Engine
	modules []Module
}

type ModuleContext struct {
	DB              sqlc.Querier
	Redis           *redis.Client
	JWTService      auth.TokenService
	MailService     mail.EmailProviderService
	RabbitMQService rabbitmq.RabbitMQService
}

func NewApplication(cfg *config.Config) (*Application, error) {
	if err := validation.InitValidator(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Validator init failed")
		return nil, err
	}

	r := gin.Default()

	// init DB
	if err := db.InitDB(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Database init failed")
		return nil, err
	}

	// init redis
	redisClient := config.NewRedisClient()
	cacheRedisService := cache.NewRedisCacheService(redisClient)

	// init jwt
	tokenService := auth.NewJWTService(cacheRedisService)

	// init mail
	mailLogger := utils.NewLoggerWithPath("mail.log", "info")
	factory, err := mail.NewProviderFactory(mail.ProviderMailtrap)
	if err != nil {
		mailLogger.Error().Err(err).Msg("Failed to create mail provider factory")
		return nil, err
	}

	mailService, err := mail.NewMailService(cfg, mailLogger, factory)
	if err != nil {
		mailLogger.Error().Err(err).Msg("Failed to a initialize mail service")
		return nil, err
	}

	// init rabbitmq
	rabbitMQLogger := utils.NewLoggerWithPath("worker.log", "info")
	rabbitMQService, _ := rabbitmq.NewRabbitMQService(
		utils.GetEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"),
		rabbitMQLogger,
	)

	ctx := &ModuleContext{
		DB:              db.DB,
		Redis:           redisClient,
		JWTService:      tokenService,
		MailService:     mailService,
		RabbitMQService: rabbitMQService,
	}

	modules := []Module{
		NewUserModule(ctx),
		NewAuthModule(ctx),
	}

	routes.RegisterRoutes(r, tokenService, cacheRedisService, getModuleRoutes(modules)...)

	return &Application{
		config:  cfg,
		router:  r,
		modules: modules,
	}, nil
}

func (a *Application) Run() error {
	// return a.router.Run(a.config.ServerAddress)
	srv := &http.Server{
		Addr:    a.config.ServerAddress,
		Handler: a.router,
	}

	quit := make(chan os.Signal, 1)

	// syscall.SIGINT --> Ctrl + C
	// syscall.SIGTERM --> Kill service
	// syscall.SIGHUP --> Reload service
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		logger.Log.Info().Msgf("Server is running: %s", a.config.ServerAddress)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	<-quit
	logger.Log.Warn().Msg("Shutdown signal received...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	logger.Log.Info().Msg("Server exited gracefully")

	return nil
}

func getModuleRoutes(modules []Module) []routes.Route {
	routeList := make([]routes.Route, len(modules))

	for i, module := range modules {
		routeList[i] = module.Routes()
	}

	return routeList
}
