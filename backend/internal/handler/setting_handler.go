package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"satset-app/internal/service"
)

// SettingHandler handles business settings endpoints.
type SettingHandler struct {
	settingService *service.SettingService
}

// NewSettingHandler creates a new SettingHandler.
func NewSettingHandler(settingService *service.SettingService) *SettingHandler {
	return &SettingHandler{settingService: settingService}
}

// GetAll handles GET /api/settings
func (h *SettingHandler) GetAll(c *gin.Context) {
	settings, err := h.settingService.GetAll(c.Request.Context())
	if err != nil {
		RespondInternalError(c, "Gagal mengambil pengaturan")
		return
	}

	RespondOK(c, "", settings)
}

// Update handles PUT /api/settings
func (h *SettingHandler) Update(c *gin.Context) {
	var req service.UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondBadRequest(c, "Data tidak valid: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")
	ipAddress := c.ClientIP()

	if err := h.settingService.Update(c.Request.Context(), req, userID, ipAddress); err != nil {
		RespondInternalError(c, err.Error())
		return
	}

	RespondOK(c, "Pengaturan berhasil diperbarui", nil)
}

// UploadLogo handles POST /api/settings/logo (multipart/form-data)
func (h *SettingHandler) UploadLogo(c *gin.Context) {
	file, err := c.FormFile("logo")
	if err != nil {
		RespondBadRequest(c, "File logo tidak ditemukan: "+err.Error())
		return
	}

	// Validate file size (max 2MB)
	if file.Size > 2*1024*1024 {
		RespondBadRequest(c, "Ukuran file logo maksimal 2MB")
		return
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		RespondBadRequest(c, "Format file logo harus berupa JPG, PNG, atau WebP")
		return
	}

	// Create uploads directory if not exists
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		RespondInternalError(c, "Gagal membuat direktori upload")
		return
	}

	// Save file with unique timestamp
	filename := fmt.Sprintf("logo_%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		RespondInternalError(c, "Gagal menyimpan file logo: "+err.Error())
		return
	}

	logoURL := "/uploads/" + filename
	userID := c.GetUint("user_id")
	ipAddress := c.ClientIP()

	if err := h.settingService.UpdateLogo(c.Request.Context(), logoURL, userID, ipAddress); err != nil {
		RespondInternalError(c, err.Error())
		return
	}

	RespondOK(c, "Logo berhasil diperbarui", gin.H{"logo_url": logoURL})
}
