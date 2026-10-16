package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"satset-app/internal/service"
)

// PaymentHandler handles payment endpoints.
type PaymentHandler struct {
	paymentService *service.PaymentService
}

// NewPaymentHandler creates a new PaymentHandler.
func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

// RecordPayment handles POST /api/payments
func (h *PaymentHandler) RecordPayment(c *gin.Context) {
	var req service.RecordPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondBadRequest(c, "Data pembayaran tidak valid: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")
	ipAddress := c.ClientIP()

	payment, err := h.paymentService.RecordPayment(c.Request.Context(), req, userID, ipAddress)
	if err != nil {
		RespondBadRequest(c, err.Error())
		return
	}

	RespondCreated(c, "Pembayaran berhasil dicatat", payment)
}

// CancelPayment handles PUT /api/payments/:id/cancel
func (h *PaymentHandler) CancelPayment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		RespondBadRequest(c, "ID tidak valid")
		return
	}

	var req service.CancelPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondBadRequest(c, "Alasan pembatalan diperlukan: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")
	ipAddress := c.ClientIP()

	payment, err := h.paymentService.CancelPayment(c.Request.Context(), uint(id), req, userID, ipAddress)
	if err != nil {
		RespondBadRequest(c, err.Error())
		return
	}

	RespondOK(c, "Pembayaran berhasil dibatalkan", payment)
}

// GetByOrder handles GET /api/orders/:id/payments
func (h *PaymentHandler) GetByOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		RespondBadRequest(c, "ID tidak valid")
		return
	}

	payments, err := h.paymentService.GetPaymentsByOrder(c.Request.Context(), uint(id))
	if err != nil {
		RespondInternalError(c, "Gagal mengambil data pembayaran")
		return
	}

	RespondOK(c, "", payments)
}
