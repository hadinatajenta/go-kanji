package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	authdelivery "gobackend/src/auth/delivery"
	authrepository "gobackend/src/auth/repository"
	authroutes "gobackend/src/auth/routes"
	authservice "gobackend/src/auth/service"
	logrepository "gobackend/src/logs/repository"
	logservice "gobackend/src/logs/service"
)

const (
	googleClientIDEnv      = "GOOGLE_CLIENT_ID"
	googleClientSecretEnv  = "GOOGLE_CLIENT_SECRET"
	googleRedirectURIEnv   = "GOOGLE_REDIRECT_URI"
	jwtSecretEnv           = "JWT_SECRET"
	jwtTTLMinutesEnv       = "JWT_TOKEN_TTL_MINUTES"
	authSuccessRedirectEnv = "AUTH_SUCCESS_REDIRECT_URL"
	authFailureRedirectEnv = "AUTH_FAILURE_REDIRECT_URL"

	defaultJWTTokenTTL = time.Hour
)

// RegisterAuthFeature wires the auth feature (repository, service, handlers, routes) into the provided router.
func RegisterAuthFeature(router gin.IRouter, database *sql.DB) error {
	if router == nil {
		return fmt.Errorf("register auth feature: router is nil")
	}

	if database == nil {
		return fmt.Errorf("register auth feature: database is nil")
	}

	userRepository, err := authrepository.NewPostgresUserRepository(database)
	if err != nil {
		return fmt.Errorf("initialise auth repository: %w", err)
	}

	logRepo := logrepository.NewPostgresRepository(database)
	if err := logRepo.EnsureSchema(context.Background()); err != nil {
		return fmt.Errorf("ensure user logs schema: %w", err)
	}
	activityLogService := logservice.NewLogService(logRepo)

	authConfig := authservice.GoogleAuthConfig{
		ClientID:     os.Getenv(googleClientIDEnv),
		ClientSecret: os.Getenv(googleClientSecretEnv),
		RedirectURL:  os.Getenv(googleRedirectURIEnv),
		JWTSecret:    os.Getenv(jwtSecretEnv),
		TokenTTL:     readJWTTTL(),
		LogService:   activityLogService,
	}

	authService, err := authservice.NewGoogleAuthService(userRepository, authConfig)
	if err != nil {
		return fmt.Errorf("initialise google auth service: %w", err)
	}

	successRedirectURL := os.Getenv(authSuccessRedirectEnv)
	failureRedirectURL := os.Getenv(authFailureRedirectEnv)
	handler := authdelivery.NewHandler(authService, successRedirectURL, failureRedirectURL, activityLogService)
	authroutes.Register(router, handler)

	return nil
}

func readJWTTTL() time.Duration {
	value := os.Getenv(jwtTTLMinutesEnv)
	if value == "" {
		return defaultJWTTokenTTL
	}

	minutes, err := strconv.Atoi(value)
	if err != nil || minutes <= 0 {
		log.Printf("invalid %s value %q, defaulting to %s", jwtTTLMinutesEnv, value, defaultJWTTokenTTL)
		return defaultJWTTokenTTL
	}

	return time.Duration(minutes) * time.Minute
}
