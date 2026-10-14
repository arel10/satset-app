package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/satset/api/internal/pkg/response"
	"github.com/satset/api/internal/platform/logger"
)

const RequestIDHeader = "X-Request-ID"
const RequestIDKey = "request_id"

func generateRequestID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader(RequestIDHeader)
		if reqID == "" {
			reqID = generateRequestID()
		}
		c.Set(RequestIDKey, reqID)
		c.Header(RequestIDHeader, reqID)
		c.Next()
	}
}

func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		if raw != "" {
			path = path + "?" + raw
		}

		reqID, _ := c.Get(RequestIDKey)
		reqIDStr, _ := reqID.(string)
		statusCode := c.Writer.Status()

		event := log.Info()
		if statusCode >= 500 {
			event = log.Error()
		} else if statusCode >= 400 {
			event = log.Warn()
		}

		event.
			Str("request_id", reqIDStr).
			Str("method", c.Request.Method).
			Str("path", path).
			Int("status", statusCode).
			Dur("latency", time.Since(start)).
			Str("client_ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Msg("HTTP request completed")
	}
}

func Recovery(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}

				reqID, _ := c.Get(RequestIDKey)
				reqIDStr, _ := reqID.(string)

				log.Error().
					Str("request_id", reqIDStr).
					Err(err).
					Str("stack", string(debug.Stack())).
					Msg("Unhandled panic recovered")

				c.Abort()
				response.SendError(c, response.ErrInternal)
			}
		}()
		c.Next()
	}
}

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "SAMEORIGIN")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}

func StrictCORS(allowedBaseURL string, isProduction bool) gin.HandlerFunc {
	origins := []string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
	}

	if allowedBaseURL != "" && !contains(origins, allowedBaseURL) {
		origins = append(origins, allowedBaseURL)
	}

	corsConfig := cors.Config{
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			RequestIDHeader,
			"Idempotency-Key",
		},
		ExposeHeaders: []string{
			RequestIDHeader,
			"Content-Length",
			"Retry-After",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	if isProduction {
		// Strict matching in production
		corsConfig.AllowOrigins = []string{allowedBaseURL}
	} else {
		corsConfig.AllowOriginFunc = func(origin string) bool {
			for _, o := range origins {
				if strings.EqualFold(o, origin) {
					return true
				}
			}
			return false
		}
	}

	return cors.New(corsConfig)
}

func Gzip() gin.HandlerFunc {
	return gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPaths([]string{"/healthz", "/readyz"}))
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if strings.EqualFold(item, val) {
			return true
		}
	}
	return false
}
