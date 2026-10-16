package handler

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Response is the standard API response structure.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta contains pagination information.
type Meta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// NewMeta creates pagination metadata from the given parameters.
func NewMeta(page, perPage int, total int64) *Meta {
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	return &Meta{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}
}

// RespondSuccess sends a success JSON response.
func RespondSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// RespondSuccessWithMeta sends a success JSON response with pagination metadata.
func RespondSuccessWithMeta(c *gin.Context, statusCode int, message string, data interface{}, meta *Meta) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

// RespondError sends an error JSON response.
func RespondError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   message,
	})
}

// RespondCreated sends a 201 Created response.
func RespondCreated(c *gin.Context, message string, data interface{}) {
	RespondSuccess(c, http.StatusCreated, message, data)
}

// RespondOK sends a 200 OK response.
func RespondOK(c *gin.Context, message string, data interface{}) {
	RespondSuccess(c, http.StatusOK, message, data)
}

// RespondBadRequest sends a 400 Bad Request response.
func RespondBadRequest(c *gin.Context, message string) {
	RespondError(c, http.StatusBadRequest, message)
}

// RespondUnauthorized sends a 401 Unauthorized response.
func RespondUnauthorized(c *gin.Context, message string) {
	RespondError(c, http.StatusUnauthorized, message)
}

// RespondForbidden sends a 403 Forbidden response.
func RespondForbidden(c *gin.Context, message string) {
	RespondError(c, http.StatusForbidden, message)
}

// RespondNotFound sends a 404 Not Found response.
func RespondNotFound(c *gin.Context, message string) {
	RespondError(c, http.StatusNotFound, message)
}

// RespondConflict sends a 409 Conflict response.
func RespondConflict(c *gin.Context, message string) {
	RespondError(c, http.StatusConflict, message)
}

// RespondTooManyRequests sends a 429 Too Many Requests response.
func RespondTooManyRequests(c *gin.Context, message string) {
	RespondError(c, http.StatusTooManyRequests, message)
}

// RespondInternalError sends a 500 Internal Server Error response.
func RespondInternalError(c *gin.Context, message string) {
	RespondError(c, http.StatusInternalServerError, message)
}

// parsePagination extracts page and per_page from query params with sensible defaults and limits.
func parsePagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	if perPage < 1 {
		perPage = 20
	} else if perPage > 100 {
		perPage = 100
	}

	return page, perPage
}
