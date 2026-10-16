package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"satset-app/internal/service"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login handles POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondBadRequest(c, "Data login tidak valid")
		return
	}

	ipAddress := c.ClientIP()
	result, err := h.authService.Login(c.Request.Context(), req, ipAddress)
	if err != nil {
		RespondUnauthorized(c, err.Error())
		return
	}

	RespondOK(c, "Login berhasil", result)
}

// RefreshToken handles POST /api/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondBadRequest(c, "Refresh token diperlukan")
		return
	}

	result, err := h.authService.RefreshAccessToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		RespondUnauthorized(c, err.Error())
		return
	}

	RespondOK(c, "Token diperbarui", result)
}

// Logout handles POST /api/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondBadRequest(c, "Refresh token diperlukan")
		return
	}

	userID := c.GetUint("user_id")
	ipAddress := c.ClientIP()

	if err := h.authService.Logout(c.Request.Context(), req.RefreshToken, userID, ipAddress); err != nil {
		RespondInternalError(c, err.Error())
		return
	}

	RespondOK(c, "Logout berhasil", nil)
}

// Me handles GET /api/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userInfo := service.UserInfo{
		ID:       c.GetUint("user_id"),
		Username: c.GetString("username"),
		Role:     c.GetString("role"),
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    userInfo,
	})
}
