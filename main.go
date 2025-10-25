package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"gobackend/app"
	"gobackend/infra/db"
)

const (
	defaultHTTPAddr   = ":8080"
	readHeaderTimeout = 5 * time.Second
	appHTTPAddrEnv    = "APP_HTTP_ADDR"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not loaded, falling back to environment variables")
	}

	database, err := db.OpenConnection()
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}
	defer database.Close()

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(cors.New(corsConfig()))

	if err := app.RegisterAuthFeature(router, database); err != nil {
		return fmt.Errorf("register auth feature: %w", err)
	}
	if err := app.RegisterUserFeature(router, database); err != nil {
		return fmt.Errorf("register user feature: %w", err)
	}

	server := &http.Server{
		Addr:              httpAddr(),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	log.Printf("HTTP server listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server error: %w", err)
	}

	return nil
}

func httpAddr() string {
	if addr := os.Getenv(appHTTPAddrEnv); addr != "" {
		return addr
	}

	return defaultHTTPAddr
}

func corsConfig() cors.Config {
	return cors.Config{
		AllowOrigins:     allowedOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}

func allowedOrigins() []string {
	raw := os.Getenv("APP_ALLOWED_ORIGINS")
	if strings.TrimSpace(raw) == "" {
		return []string{"http://localhost:5173"}
	}

	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			origins = append(origins, trimmed)
		}
	}

	if len(origins) == 0 {
		origins = append(origins, "http://localhost:5173")
	}

	return origins
}
