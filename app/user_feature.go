package app

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"gobackend/shared/identity"
	logdelivery "gobackend/src/logs/delivery"
	logrepository "gobackend/src/logs/repository"
	logroutes "gobackend/src/logs/routes"
	logservice "gobackend/src/logs/service"
	userdelivery "gobackend/src/users/delivery"
	userrepository "gobackend/src/users/repository"
	userroutes "gobackend/src/users/routes"
	userservice "gobackend/src/users/service"
)

// RegisterUserFeature wires the user endpoints into the router.
func RegisterUserFeature(router gin.IRouter, database *sql.DB) error {
	if router == nil {
		return fmt.Errorf("register user feature: router is nil")
	}

	if database == nil {
		return fmt.Errorf("register user feature: database is nil")
	}

	referenceSalt := os.Getenv("USER_REFERENCE_SALT")
	if referenceSalt == "" {
		referenceSalt = os.Getenv("JWT_SECRET")
	}
	if referenceSalt == "" {
		referenceSalt = "default-user-reference-salt"
	}

	refEncoder, err := identity.NewUserReferenceEncoder(referenceSalt)
	if err != nil {
		return fmt.Errorf("initialise user reference encoder: %w", err)
	}

	repo := userrepository.NewPostgresUserRepository(database)
	service := userservice.NewUserService(repo, refEncoder)
	handler := userdelivery.NewHandler(service)

	userroutes.Register(router, handler)

	logRepo := logrepository.NewPostgresRepository(database)
	if err := logRepo.EnsureSchema(context.Background()); err != nil {
		return fmt.Errorf("ensure user logs schema: %w", err)
	}
	logService := logservice.NewLogService(logRepo)
	logHandler := logdelivery.NewHandler(logService, refEncoder)
	logroutes.Register(router, logHandler)

	return nil
}
