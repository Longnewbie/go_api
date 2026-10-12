package main

import (
	"user-management-api/internal/app"
	"user-management-api/internal/config"
)

func main() {
	// init configuration
	cfg := config.NewConfig()

	// init application
	application := app.NewApplication(cfg)

	// start server
	if err := application.Run(); err != nil {
		panic(err)
	}
}
