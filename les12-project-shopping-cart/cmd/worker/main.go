package main

import (
	"context"
	"encoding/json"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
	"user-management-api/internal/config"
	"user-management-api/internal/utils"
	"user-management-api/pkg/logger"
	"user-management-api/pkg/mail"
	"user-management-api/pkg/rabbitmq"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type Worker struct {
	rabbitMQ    rabbitmq.RabbitMQService
	mailService mail.EmailProviderService
	cfg         *config.Config
	logger      *zerolog.Logger
}

func NewWorker(cfg *config.Config) *Worker {
	log := utils.NewLoggerWithPath("worker.log", "info")

	// Connect RabbitMQ
	rabbitMQ, err := rabbitmq.NewRabbitMQService(
		utils.GetEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/"),
		log,
	)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to init RabbitMQ service")
	}

	// Init Mail Service
	mailLogger := utils.NewLoggerWithPath("mail.log", "info")
	factory, err := mail.NewProviderFactory(mail.ProviderMailtrap)
	if err != nil {
		mailLogger.Error().Err(err).Msg("Failed to create mail provider factory")
		return nil
	}

	mailService, err := mail.NewMailService(cfg, mailLogger, factory)
	if err != nil {
		mailLogger.Error().Err(err).Msg("Failed to a initialize mail service")
		return nil
	}

	return &Worker{
		rabbitMQ:    rabbitMQ,
		mailService: mailService,
		cfg:         cfg,
		logger:      log,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	const emailQueueName = "auth_email_queue"

	// Define the handler function for processing messages
	handler := func(body []byte) error {
		w.logger.Debug().Msgf("Received message: %s", string(body))

		var email mail.Email
		if err := json.Unmarshal(body, &email); err != nil {
			w.logger.Error().Err(err).Msg("Failed to unmarshal message")
			return err
		}

		// Send email using the mail service
		if err := w.mailService.SendMail(ctx, &email); err != nil {
			return utils.NewError("Failed to send password reset email", utils.ErrCodeInternal)
		}

		w.logger.Info().Msgf("Email sent successfully to: %s", email.To)

		return nil
	}

	// Start consuming messages from the queue
	if err := w.rabbitMQ.Consume(ctx, emailQueueName, handler); err != nil {
		w.logger.Error().Err(err).Msg("Failed to start consumer")
		return err
	}

	w.logger.Info().Msgf("Worker started, consuming from queue: %s", emailQueueName)

	// Wait for context cancellation
	<-ctx.Done()

	w.logger.Info().Msg("Worker stopped consuming due to context cancellation")

	return ctx.Err()
}

func (w *Worker) Shutdown(ctx context.Context) error {
	w.logger.Info().Msg("Shutting down worker...")

	// Close the RabbitMQ connection
	if err := w.rabbitMQ.Close(); err != nil {
		w.logger.Error().Err(err).Msg("Failed to close RabbitMQ")
		return err
	}

	w.logger.Info().Msg("RabbitMQ connection close successfully")

	// Check if the context has a deadline and log a warning if it does
	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			w.logger.Warn().Msg("Shutdown timeout exceeded")
			return ctx.Err()
		}
	default:
	}

	w.logger.Info().Msg("Worker shutdown completed")

	return nil
}

func main() {
	rootDir := utils.MustGetWorkingDir()

	logFile := filepath.Join(rootDir, "internal/logs/app.log")

	logger.InitLogger(logger.LoggerConfig{
		Levels:     "info",
		Filename:   logFile,
		MaxSize:    1, // megabytes
		MaxBackups: 5,
		MaxAge:     5,    // days
		Compress:   true, // disabled by default
		IsDev:      utils.GetEnv("APP_EVN", "development"),
	})

	// load .env
	if err := godotenv.Load(filepath.Join(rootDir, ".env")); err != nil {
		logger.Log.Fatal().Err(err).Msg("No .env file found")
	} else {
		logger.Log.Info().Msg("Load successfully .env in worker process")
	}

	// init configuration
	cfg := config.NewConfig()

	// init worker
	worker := NewWorker(cfg)
	if worker == nil {
		logger.Log.Fatal().Msg("Failed to create worker")
	}

	// Start the worker in a separate goroutine
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer stop()

	// Use a WaitGroup to wait for the worker to finish
	var wg sync.WaitGroup
	wg.Add(1)

	// Start the worker in a separate goroutine
	go func() {
		defer wg.Done()
		if err := worker.Start(ctx); err != nil && err != context.Canceled {
			logger.Log.Error().Err(err).Msg("Worker failed to start")
		}
	}()

	<-ctx.Done()
	logger.Log.Info().Msg("Received shutdown signal")

	// Create a new context with a timeout for the shutdown process
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the worker
	if err := worker.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error().Err(err).Msg("Shutdown failed")
	}

	// Wait for the worker to finish
	wg.Wait()
	logger.Log.Info().Msg("Main process terminated")
}
