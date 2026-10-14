package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/satset/api/internal/pkg/response"
	"github.com/satset/api/internal/platform/logger"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware_RequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestID())
	r.GET("/test", func(c *gin.Context) {
		reqID, exists := c.Get(RequestIDKey)
		assert.True(t, exists)
		assert.NotEmpty(t, reqID)
		c.String(http.StatusOK, "ok")
	})

	// Case 1: Without incoming header
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get(RequestIDHeader))

	// Case 2: With incoming header
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req2.Header.Set(RequestIDHeader, "custom-req-id-12345")
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, "custom-req-id-12345", w2.Header().Get(RequestIDHeader))
}

func TestMiddleware_SecurityHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(SecurityHeaders())
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, "nosniff", w.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "SAMEORIGIN", w.Header().Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", w.Header().Get("X-XSS-Protection"))
	assert.Equal(t, "strict-origin-when-cross-origin", w.Header().Get("Referrer-Policy"))
}

func TestMiddleware_Recovery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	log := logger.New("development")
	r := gin.New()
	r.Use(RequestID())
	r.Use(Recovery(log))
	r.GET("/panic", func(c *gin.Context) {
		panic("fatal unexpected exception")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp response.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "INTERNAL_ERROR", resp.Error.Code)
}

func TestMiddleware_StrictCORS(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(StrictCORS("http://localhost:3000", false))
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// Allowed origin
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	r.ServeHTTP(w, req)

	assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))

	// Disallowed origin
	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/test", nil)
	req2.Header.Set("Origin", "http://evil-hacker.com")
	r.ServeHTTP(w2, req2)

	assert.Empty(t, w2.Header().Get("Access-Control-Allow-Origin"))
}
