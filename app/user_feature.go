package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"

	userdelivery "gobackend/src/users/delivery"
	userrepository "gobackend/src/users/repository"
	userroutes "gobackend/src/users/routes"
	userservice "gobackend/src/users/service"

	logdelivery "gobackend/src/users/logs/delivery"
	logrepository "gobackend/src/users/logs/repository"
	logroutes "gobackend/src/users/logs/routes"
	logservice "gobackend/src/users/logs/service"
)

// RegisterUserFeature wires the user endpoints into the router.
func RegisterUserFeature(router gin.IRouter, database *sql.DB) error {
	if router == nil {
		return fmt.Errorf("register user feature: router is nil")
	}

	if database == nil {
		return fmt.Errorf("register user feature: database is nil")
	}

	repo := userrepository.NewPostgresUserRepository(database)
	service := userservice.NewUserService(repo)
	handler := userdelivery.NewHandler(service)

	userroutes.Register(router, handler)

	logRepo := logrepository.NewPostgresRepository(database)
	if err := logRepo.EnsureSchema(context.Background()); err != nil {
		return fmt.Errorf("ensure user logs schema: %w", err)
	}
	logService := logservice.NewLogService(logRepo)
	logHandler := logdelivery.NewHandler(logService)
	logroutes.Register(router, logHandler)

	return nil
}
