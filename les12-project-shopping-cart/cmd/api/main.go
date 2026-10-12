package main

import (
	"path/filepath"
	"user-management-api/internal/app"
	"user-management-api/internal/config"
	"user-management-api/internal/utils"
	"user-management-api/pkg/logger"

	"github.com/joho/godotenv"
)

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
		logger.Log.Info().Msg("Load successfully .env in api process")
	}

	// init configuration
	cfg := config.NewConfig()

	// init application
	application, err := app.NewApplication(cfg)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to initialize application")
	}

	// start server
	if err := application.Run(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Application run failed")
	}
}
