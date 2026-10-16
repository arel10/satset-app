package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"satset-app/internal/repository"
	"satset-app/internal/service"
)

// OrderHandler handles order and POS endpoints.
type OrderHandler struct {
	orderService *service.OrderService
}

// NewOrderHandler creates a new OrderHandler.
func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// Create handles POST /api/orders
func (h *OrderHandler) Create(c *gin.Context) {
	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondBadRequest(c, "Data pesanan tidak valid: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")
	ipAddress := c.ClientIP()
	idempotencyKey := c.GetHeader("Idempotency-Key")

	order, err := h.orderService.CreateOrder(c.Request.Context(), req, userID, ipAddress, idempotencyKey)
	if err != nil {
		RespondBadRequest(c, err.Error())
		return
	}

	RespondCreated(c, "Pesanan berhasil dibuat", order)
}

// GetByID handles GET /api/orders/:id
func (h *OrderHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		RespondBadRequest(c, "ID tidak valid")
		return
	}

	order, err := h.orderService.GetOrder(c.Request.Context(), uint(id))
	if err != nil {
		RespondNotFound(c, err.Error())
		return
	}

	RespondOK(c, "", order)
}

// GetByNumber handles GET /api/orders/number/:orderNumber
func (h *OrderHandler) GetByNumber(c *gin.Context) {
	orderNumber := c.Param("orderNumber")
	if orderNumber == "" {
		RespondBadRequest(c, "Nomor pesanan diperlukan")
		return
	}

	order, err := h.orderService.GetOrderByNumber(c.Request.Context(), orderNumber)
	if err != nil {
		RespondNotFound(c, err.Error())
		return
	}

	RespondOK(c, "", order)
}

// List handles GET /api/orders
func (h *OrderHandler) List(c *gin.Context) {
	page, perPage := parsePagination(c)

	customerID := uint(0)
	if cid := c.Query("customer_id"); cid != "" {
		if parsed, err := strconv.ParseUint(cid, 10, 32); err == nil {
			customerID = uint(parsed)
		}
	}

	filter := repository.OrderFilter{
		Query:         c.Query("q"),
		WorkStatus:    c.Query("work_status"),
		PaymentStatus: c.Query("payment_status"),
		StartDate:     c.Query("start_date"),
		EndDate:       c.Query("end_date"),
		Overdue:       c.Query("overdue") == "true",
		CustomerID:    customerID,
	}

	orders, total, err := h.orderService.ListOrders(c.Request.Context(), filter, page, perPage)
	if err != nil {
		RespondInternalError(c, "Gagal mengambil data pesanan")
		return
	}

	meta := NewMeta(page, perPage, total)
	RespondSuccessWithMeta(c, 200, "", orders, meta)
}

// UpdateStatus handles PUT /api/orders/:id/status
func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		RespondBadRequest(c, "ID tidak valid")
		return
	}

	var req service.UpdateWorkStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondBadRequest(c, "Data tidak valid: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")
	userRole := c.GetString("role")
	ipAddress := c.ClientIP()

	order, err := h.orderService.UpdateWorkStatus(c.Request.Context(), uint(id), req, userID, userRole, ipAddress)
	if err != nil {
		RespondBadRequest(c, err.Error())
		return
	}

	RespondOK(c, "Status pesanan berhasil diperbarui", order)
}

// Cancel handles PUT /api/orders/:id/cancel
func (h *OrderHandler) Cancel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		RespondBadRequest(c, "ID tidak valid")
		return
	}

	var req service.CancelOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondBadRequest(c, "Alasan pembatalan diperlukan: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")
	userRole := c.GetString("role")
	ipAddress := c.ClientIP()

	order, err := h.orderService.CancelOrder(c.Request.Context(), uint(id), req, userID, userRole, ipAddress)
	if err != nil {
		RespondBadRequest(c, err.Error())
		return
	}

	RespondOK(c, "Pesanan berhasil dibatalkan", order)
}
