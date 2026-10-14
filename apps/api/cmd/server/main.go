package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/satset/api/docs"
	apphttp "github.com/satset/api/internal/http"
	"github.com/satset/api/internal/platform/config"
	"github.com/satset/api/internal/platform/db"
	"github.com/satset/api/internal/platform/logger"
	"github.com/satset/api/internal/platform/redis"
)

// @title SatSet Laundry Management API
// @version 1.0
// @description Backend API modular monolith platform SatSet dengan AI Business Assistant.
// @contact.name SatSet Engineering Team
// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token dalam format: Bearer <token>
func main() {
	cfg, err := config.Load()
	if err != nil {
		// Logger may not be initialized yet, print to stderr and fail-fast
		os.Stderr.WriteString("Configuration fatal error: " + err.Error() + "\n")
		os.Exit(1)
	}

	log := logger.New(cfg.AppEnv)

	log.Info().
		Str("env", cfg.AppEnv).
		Str("port", cfg.AppPort).
		Msg("Starting SatSet API server...")

	// Initialize MySQL with connection pooling from config
	database, err := db.Connect(
		cfg.DSN(),
		cfg.DBMaxOpenConns,
		cfg.DBMaxIdleConns,
		cfg.DBConnMaxLifetime(),
	)
	if err != nil {
		log.Warn().Err(err).Msg("Database connection initialization failed, continuing with degraded status")
	} else {
		defer database.Close()
	}

	// Initialize Redis
	redisClient, err := redis.Connect(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Warn().Err(err).Msg("Redis client initialization failed, continuing with degraded status")
	} else {
		defer redisClient.Close()
	}

	healthHandler := apphttp.NewHealthHandler(database, redisClient)
	router := apphttp.NewRouter(cfg, log, healthHandler)

	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info().Msgf("Server listening on :%s", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Server forced to shutdown")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown after timeout")
	}

	log.Info().Msg("Server exited successfully")
}
