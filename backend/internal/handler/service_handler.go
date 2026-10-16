package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"satset-app/internal/service"
)

// ServiceHandler handles laundry service management endpoints.
type ServiceHandler struct {
	serviceService *service.ServiceService
}

// NewServiceHandler creates a new ServiceHandler.
func NewServiceHandler(serviceService *service.ServiceService) *ServiceHandler {
	return &ServiceHandler{serviceService: serviceService}
}

// Create handles POST /api/services
func (h *ServiceHandler) Create(c *gin.Context) {
	var req service.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondBadRequest(c, "Data layanan tidak valid: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")
	ipAddress := c.ClientIP()

	svc, err := h.serviceService.CreateService(c.Request.Context(), req, userID, ipAddress)
	if err != nil {
		RespondBadRequest(c, err.Error())
		return
	}

	RespondCreated(c, "Layanan berhasil dibuat", svc)
}

// GetByID handles GET /api/services/:id
func (h *ServiceHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		RespondBadRequest(c, "ID tidak valid")
		return
	}

	svc, err := h.serviceService.GetService(c.Request.Context(), uint(id))
	if err != nil {
		RespondNotFound(c, err.Error())
		return
	}

	RespondOK(c, "", svc)
}

// List handles GET /api/services
func (h *ServiceHandler) List(c *gin.Context) {
	activeOnly := c.Query("active") == "true"
	page, perPage := parsePagination(c)

	services, total, err := h.serviceService.ListServices(c.Request.Context(), activeOnly, page, perPage)
	if err != nil {
		RespondInternalError(c, "Gagal mengambil data layanan")
		return
	}

	meta := NewMeta(page, perPage, total)
	RespondSuccessWithMeta(c, 200, "", services, meta)
}

// GetActive handles GET /api/services/active (for POS)
func (h *ServiceHandler) GetActive(c *gin.Context) {
	services, err := h.serviceService.GetActiveServices(c.Request.Context())
	if err != nil {
		RespondInternalError(c, "Gagal mengambil data layanan aktif")
		return
	}

	RespondOK(c, "", services)
}

// Update handles PUT /api/services/:id
func (h *ServiceHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		RespondBadRequest(c, "ID tidak valid")
		return
	}

	var req service.UpdateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondBadRequest(c, "Data tidak valid: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")
	ipAddress := c.ClientIP()

	svc, err := h.serviceService.UpdateService(c.Request.Context(), uint(id), req, userID, ipAddress)
	if err != nil {
		RespondBadRequest(c, err.Error())
		return
	}

	RespondOK(c, "Layanan berhasil diperbarui", svc)
}
