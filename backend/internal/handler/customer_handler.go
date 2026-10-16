package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"satset-app/internal/service"
)

// CustomerHandler handles customer management endpoints.
type CustomerHandler struct {
	customerService *service.CustomerService
	orderService    *service.OrderService
}

// NewCustomerHandler creates a new CustomerHandler.
func NewCustomerHandler(customerService *service.CustomerService, orderService *service.OrderService) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
		orderService:    orderService,
	}
}

// Create handles POST /api/customers
func (h *CustomerHandler) Create(c *gin.Context) {
	var req service.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondBadRequest(c, "Data pelanggan tidak valid: "+err.Error())
		return
	}

	customer, err := h.customerService.CreateCustomer(c.Request.Context(), req)
	if err != nil {
		RespondBadRequest(c, err.Error())
		return
	}

	RespondCreated(c, "Pelanggan berhasil ditambahkan", customer)
}

// GetByID handles GET /api/customers/:id
func (h *CustomerHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		RespondBadRequest(c, "ID tidak valid")
		return
	}

	customer, err := h.customerService.GetCustomer(c.Request.Context(), uint(id))
	if err != nil {
		RespondNotFound(c, err.Error())
		return
	}

	RespondOK(c, "", customer)
}

// Search handles GET /api/customers
func (h *CustomerHandler) Search(c *gin.Context) {
	query := c.Query("q")
	page, perPage := parsePagination(c)

	customers, total, err := h.customerService.SearchCustomers(c.Request.Context(), query, page, perPage)
	if err != nil {
		RespondInternalError(c, "Gagal mencari pelanggan")
		return
	}

	meta := NewMeta(page, perPage, total)
	RespondSuccessWithMeta(c, 200, "", customers, meta)
}

// Update handles PUT /api/customers/:id
func (h *CustomerHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		RespondBadRequest(c, "ID tidak valid")
		return
	}

	var req service.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondBadRequest(c, "Data tidak valid: "+err.Error())
		return
	}

	customer, err := h.customerService.UpdateCustomer(c.Request.Context(), uint(id), req)
	if err != nil {
		RespondBadRequest(c, err.Error())
		return
	}

	RespondOK(c, "Pelanggan berhasil diperbarui", customer)
}

// GetOrderHistory handles GET /api/customers/:id/orders
func (h *CustomerHandler) GetOrderHistory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		RespondBadRequest(c, "ID tidak valid")
		return
	}

	limit := 50
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	orders, err := h.orderService.GetOrdersByCustomer(c.Request.Context(), uint(id), limit)
	if err != nil {
		RespondInternalError(c, "Gagal mengambil riwayat pesanan")
		return
	}

	RespondOK(c, "", orders)
}
