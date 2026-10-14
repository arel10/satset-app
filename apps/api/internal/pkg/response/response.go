package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorPayload struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details"`
}

type ErrorResponse struct {
	Error ErrorPayload `json:"error"`
}

type SuccessResponse struct {
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta,omitempty"`
}

type AppError struct {
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	HTTPStatus int                    `json:"-"`
	Details    map[string]interface{} `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Details:    make(map[string]interface{}),
	}
}

func (e *AppError) WithDetails(key string, val interface{}) *AppError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = val
	return e
}

// Predefined Domain Errors
var (
	ErrBadRequest = &AppError{
		Code:       "BAD_REQUEST",
		Message:    "Permintaan tidak valid",
		HTTPStatus: http.StatusBadRequest,
		Details:    make(map[string]interface{}),
	}

	ErrValidation = &AppError{
		Code:       "VALIDATION_ERROR",
		Message:    "Terdapat kesalahan validasi pada data yang dikirimkan",
		HTTPStatus: http.StatusBadRequest,
		Details:    make(map[string]interface{}),
	}

	ErrUnauthorized = &AppError{
		Code:       "UNAUTHORIZED",
		Message:    "Autentikasi diperlukan untuk mengakses resource ini",
		HTTPStatus: http.StatusUnauthorized,
		Details:    make(map[string]interface{}),
	}

	ErrForbidden = &AppError{
		Code:       "FORBIDDEN",
		Message:    "Anda tidak memiliki hak akses untuk tindakan ini",
		HTTPStatus: http.StatusForbidden,
		Details:    make(map[string]interface{}),
	}

	ErrNotFound = &AppError{
		Code:       "RESOURCE_NOT_FOUND",
		Message:    "Data yang dicari tidak ditemukan",
		HTTPStatus: http.StatusNotFound,
		Details:    make(map[string]interface{}),
	}

	ErrConflict = &AppError{
		Code:       "CONFLICT",
		Message:    "Terjadi konflik pada data transaksi",
		HTTPStatus: http.StatusConflict,
		Details:    make(map[string]interface{}),
	}

	ErrOrderInvalidTransition = &AppError{
		Code:       "ORDER_INVALID_TRANSITION",
		Message:    "Perubahan status pesanan tidak diizinkan",
		HTTPStatus: http.StatusBadRequest,
		Details:    make(map[string]interface{}),
	}

	ErrRateLimited = &AppError{
		Code:       "RATE_LIMITED",
		Message:    "Batas permintaan tercapai, silakan coba beberapa saat lagi",
		HTTPStatus: http.StatusTooManyRequests,
		Details:    make(map[string]interface{}),
	}

	ErrInternal = &AppError{
		Code:       "INTERNAL_ERROR",
		Message:    "Terjadi kesalahan internal pada server",
		HTTPStatus: http.StatusInternalServerError,
		Details:    make(map[string]interface{}),
	}

	ErrServiceUnavailable = &AppError{
		Code:       "SERVICE_UNAVAILABLE",
		Message:    "Layanan sementara tidak tersedia",
		HTTPStatus: http.StatusServiceUnavailable,
		Details:    make(map[string]interface{}),
	}
)

func SendError(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		details := appErr.Details
		if details == nil {
			details = make(map[string]interface{})
		}
		c.JSON(appErr.HTTPStatus, ErrorResponse{
			Error: ErrorPayload{
				Code:    appErr.Code,
				Message: appErr.Message,
				Details: details,
			},
		})
		return
	}

	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error: ErrorPayload{
			Code:    ErrInternal.Code,
			Message: ErrInternal.Message,
			Details: make(map[string]interface{}),
		},
	})
}

func SendValidationError(c *gin.Context, details map[string]interface{}) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error: ErrorPayload{
			Code:    ErrValidation.Code,
			Message: ErrValidation.Message,
			Details: details,
		},
	})
}

func SendSuccess(c *gin.Context, httpStatus int, data interface{}) {
	c.JSON(httpStatus, SuccessResponse{
		Data: data,
	})
}

func SendSuccessWithMeta(c *gin.Context, httpStatus int, data interface{}, meta interface{}) {
	c.JSON(httpStatus, SuccessResponse{
		Data: data,
		Meta: meta,
	})
}
