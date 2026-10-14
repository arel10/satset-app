package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockDB struct {
	err error
}

func (m *mockDB) PingContext(ctx context.Context) error {
	return m.err
}

type mockRedis struct {
	err error
}

func (m *mockRedis) Ping(ctx context.Context) error {
	return m.err
}

func TestHealthHandler_Healthz(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHealthHandler(&mockDB{}, &mockRedis{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/healthz", nil)

	handler.Healthz(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"ok"`)
}

func TestHealthHandler_Readyz_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHealthHandler(&mockDB{err: nil}, &mockRedis{err: nil})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/readyz", nil)

	handler.Readyz(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"ok"`)
	assert.Contains(t, w.Body.String(), `"database":"ok"`)
	assert.Contains(t, w.Body.String(), `"redis":"ok"`)
}

func TestHealthHandler_Readyz_Failure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHealthHandler(&mockDB{err: errors.New("db connection timeout")}, &mockRedis{err: nil})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/readyz", nil)

	handler.Readyz(c)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Contains(t, w.Body.String(), `"status":"degraded"`)
	assert.Contains(t, w.Body.String(), `"db connection timeout"`)
}
