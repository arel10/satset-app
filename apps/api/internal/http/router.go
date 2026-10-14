package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/satset/api/internal/http/middleware"
	"github.com/satset/api/internal/pkg/response"
	"github.com/satset/api/internal/platform/config"
	"github.com/satset/api/internal/platform/logger"
)

func NewRouter(cfg *config.Config, log *logger.Logger, healthHandler *HealthHandler) *gin.Engine {
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()

	// Global middlewares
	router.Use(middleware.RequestID())
	router.Use(middleware.Logger(log))
	router.Use(middleware.Recovery(log))
	router.Use(middleware.SecurityHeaders())
	router.Use(middleware.StrictCORS(cfg.AppBaseURL, cfg.IsProduction()))
	router.Use(middleware.Gzip())

	// Health check probes
	router.GET("/healthz", healthHandler.Healthz)
	router.GET("/readyz", healthHandler.Readyz)

	// Swagger UI (non-production only)
	if !cfg.IsProduction() {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Ping godoc
		// @Summary Ping endpoint
		// @Description Menguji koneksi API v1
		// @Tags System
		// @Produce json
		// @Success 200 {object} response.SuccessResponse
		// @Router /api/v1/ping [get]
		v1.GET("/ping", func(c *gin.Context) {
			response.SendSuccess(c, http.StatusOK, gin.H{
				"message": "pong",
				"version": "v1",
			})
		})
	}

	// 404 Route Handler
	router.NoRoute(func(c *gin.Context) {
		response.SendError(c, response.ErrNotFound)
	})

	return router
}
