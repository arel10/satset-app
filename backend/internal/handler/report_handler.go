package handler

import (
	"github.com/gin-gonic/gin"

	"satset-app/internal/service"
)

// ReportHandler handles report endpoints.
type ReportHandler struct {
	reportService *service.ReportService
}

// NewReportHandler creates a new ReportHandler.
func NewReportHandler(reportService *service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

// GetBusinessReport handles GET /api/reports
func (h *ReportHandler) GetBusinessReport(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	report, err := h.reportService.GetBusinessReport(c.Request.Context(), startDate, endDate)
	if err != nil {
		RespondBadRequest(c, err.Error())
		return
	}

	RespondOK(c, "", report)
}
