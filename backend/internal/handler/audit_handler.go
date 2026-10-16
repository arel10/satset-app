package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"satset-app/internal/repository"
	"satset-app/internal/service"
)

// AuditHandler handles audit log endpoints.
type AuditHandler struct {
	auditService *service.AuditService
}

// NewAuditHandler creates a new AuditHandler.
func NewAuditHandler(auditService *service.AuditService) *AuditHandler {
	return &AuditHandler{auditService: auditService}
}

// List handles GET /api/audit-logs
func (h *AuditHandler) List(c *gin.Context) {
	page, perPage := parsePagination(c)

	var filter repository.AuditFilter
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if uid, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			filter.UserID = uint(uid)
		}
	}
	if entityIDStr := c.Query("entity_id"); entityIDStr != "" {
		if eid, err := strconv.ParseUint(entityIDStr, 10, 32); err == nil {
			filter.EntityID = uint(eid)
		}
	}
	filter.Action = c.Query("action")
	filter.EntityType = c.Query("entity_type")
	filter.StartDate = c.Query("start_date")
	filter.EndDate = c.Query("end_date")

	logs, total, err := h.auditService.ListAuditLogs(c.Request.Context(), filter, page, perPage)
	if err != nil {
		RespondInternalError(c, "Gagal mengambil audit log")
		return
	}

	meta := NewMeta(page, perPage, total)
	RespondSuccessWithMeta(c, 200, "", logs, meta)
}
