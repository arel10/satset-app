package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DatabasePinger interface {
	PingContext(ctx context.Context) error
}

type RedisPinger interface {
	Ping(ctx context.Context) error
}

type HealthHandler struct {
	db    DatabasePinger
	redis RedisPinger
}

func NewHealthHandler(db DatabasePinger, r RedisPinger) *HealthHandler {
	return &HealthHandler{
		db:    db,
		redis: r,
	}
}

func (h *HealthHandler) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *HealthHandler) Readyz(c *gin.Context) {
	parentCtx := context.Background()
	if c.Request != nil && c.Request.Context() != nil {
		parentCtx = c.Request.Context()
	}

	ctx, cancel := context.WithTimeout(parentCtx, 2*time.Second)
	defer cancel()

	checks := make(map[string]string)
	isReady := true

	if h.db != nil {
		if err := h.db.PingContext(ctx); err != nil {
			checks["database"] = err.Error()
			isReady = false
		} else {
			checks["database"] = "ok"
		}
	} else {
		checks["database"] = "not configured"
		isReady = false
	}

	if h.redis != nil {
		if err := h.redis.Ping(ctx); err != nil {
			checks["redis"] = err.Error()
			isReady = false
		} else {
			checks["redis"] = "ok"
		}
	} else {
		checks["redis"] = "not configured"
		isReady = false
	}

	status := "ok"
	httpStatus := http.StatusOK
	if !isReady {
		status = "degraded"
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, gin.H{
		"status": status,
		"checks": checks,
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}
