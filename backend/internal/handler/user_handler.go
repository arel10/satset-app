package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"satset-app/internal/service"
)

// UserHandler handles user management endpoints.
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Create handles POST /api/users
func (h *UserHandler) Create(c *gin.Context) {
	var req service.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondBadRequest(c, "Data pengguna tidak valid: "+err.Error())
		return
	}

	creatorID := c.GetUint("user_id")
	ipAddress := c.ClientIP()

	user, err := h.userService.CreateUser(c.Request.Context(), req, creatorID, ipAddress)
	if err != nil {
		RespondBadRequest(c, err.Error())
		return
	}

	RespondCreated(c, "Pengguna berhasil dibuat", user)
}

// GetByID handles GET /api/users/:id
func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		RespondBadRequest(c, "ID tidak valid")
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), uint(id))
	if err != nil {
		RespondNotFound(c, err.Error())
		return
	}

	RespondOK(c, "", user)
}

// List handles GET /api/users
func (h *UserHandler) List(c *gin.Context) {
	page, perPage := parsePagination(c)

	users, total, err := h.userService.ListUsers(c.Request.Context(), page, perPage)
	if err != nil {
		RespondInternalError(c, "Gagal mengambil data pengguna")
		return
	}

	meta := NewMeta(page, perPage, total)
	RespondSuccessWithMeta(c, 200, "", users, meta)
}

// Update handles PUT /api/users/:id
func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		RespondBadRequest(c, "ID tidak valid")
		return
	}

	var req service.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondBadRequest(c, "Data tidak valid: "+err.Error())
		return
	}

	updaterID := c.GetUint("user_id")
	ipAddress := c.ClientIP()

	user, err := h.userService.UpdateUser(c.Request.Context(), uint(id), req, updaterID, ipAddress)
	if err != nil {
		RespondBadRequest(c, err.Error())
		return
	}

	RespondOK(c, "Pengguna berhasil diperbarui", user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		RespondBadRequest(c, "ID tidak valid")
		return
	}

	actorID := c.GetUint("user_id")
	ipAddress := c.ClientIP()

	err = h.userService.DeleteUser(c.Request.Context(), uint(id), actorID, ipAddress)
	if err != nil {
		RespondBadRequest(c, err.Error())
		return
	}

	RespondOK(c, "Pengguna berhasil dihapus", nil)
}
