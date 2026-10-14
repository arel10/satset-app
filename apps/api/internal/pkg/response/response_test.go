package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestResponse_SendSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SendSuccess(c, http.StatusOK, map[string]string{"name": "SatSet"})

	assert.Equal(t, http.StatusOK, w.Code)
	var resp SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotNil(t, resp.Data)
}

func TestResponse_SendError_AppError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SendError(c, ErrNotFound)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "RESOURCE_NOT_FOUND", resp.Error.Code)
	assert.Equal(t, "Data yang dicari tidak ditemukan", resp.Error.Message)
}

func TestResponse_SendError_GenericError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SendError(c, errors.New("unhandled db error"))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "INTERNAL_ERROR", resp.Error.Code)
}

func TestResponse_SendValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	details := map[string]interface{}{"phone": "nomor telepon wajib diisi"}
	SendValidationError(c, details)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var resp ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "VALIDATION_ERROR", resp.Error.Code)
	assert.Equal(t, "nomor telepon wajib diisi", resp.Error.Details["phone"])
}
