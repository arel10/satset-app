package handler

import (
	"github.com/gin-gonic/gin"

	"satset-app/internal/service"
)

// DashboardHandler handles dashboard endpoints.
type DashboardHandler struct {
	dashboardService *service.DashboardService
}

// NewDashboardHandler creates a new DashboardHandler.
func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

// GetOwnerDashboard handles GET /api/dashboard/owner
func (h *DashboardHandler) GetOwnerDashboard(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	data, err := h.dashboardService.GetOwnerDashboard(c.Request.Context(), startDate, endDate)
	if err != nil {
		RespondInternalError(c, "Gagal mengambil data dashboard")
		return
	}

	RespondOK(c, "", data)
}

// GetKasirDashboard handles GET /api/dashboard/kasir
func (h *DashboardHandler) GetKasirDashboard(c *gin.Context) {
	data, err := h.dashboardService.GetKasirDashboard(c.Request.Context())
	if err != nil {
		RespondInternalError(c, "Gagal mengambil data dashboard")
		return
	}

	RespondOK(c, "", data)
}
