package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apphttp "github.com/satset/api/internal/http"
	"github.com/satset/api/internal/platform/config"
	"github.com/satset/api/internal/platform/db"
	"github.com/satset/api/internal/platform/logger"
	"github.com/satset/api/internal/platform/redis"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.AppEnv)

	log.Info().
		Str("env", cfg.AppEnv).
		Str("port", cfg.Port).
		Msg("Starting SatSet API server...")

	// Initialize MySQL
	database, err := db.Connect(cfg.DSN())
	if err != nil {
		log.Warn().Err(err).Msg("Database connection initialization failed, continuing with degraded status")
	} else {
		defer database.Close()
	}

	// Initialize Redis
	redisClient, err := redis.Connect(cfg.RedisAddr, cfg.RedisPass, cfg.RedisDB)
	if err != nil {
		log.Warn().Err(err).Msg("Redis client initialization failed, continuing with degraded status")
	} else {
		defer redisClient.Close()
	}

	healthHandler := apphttp.NewHealthHandler(database, redisClient)
	router := apphttp.NewRouter(log, healthHandler)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info().Msgf("Server listening on :%s", cfg.Port)
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
