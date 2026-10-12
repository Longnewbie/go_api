package app

import (
	"user-management-api/internal/handler"
	"user-management-api/internal/repository"
	"user-management-api/internal/routes"
	"user-management-api/internal/service"
)

type UserModule struct {
	routes routes.Route
}

func NewUserModule() *UserModule {
	// init repository
	userRepo := repository.NewInMemoryUserRepository()

	// init service
	userService := service.NewUserService(userRepo)

	// init handler
	userHandler := handler.NewUserHandler(userService)

	// init routes
	userRoutes := routes.NewUserRoutes(userHandler)

	return &UserModule{routes: userRoutes}
}

func (m *UserModule) Routes() routes.Route {
	return m.routes
}
