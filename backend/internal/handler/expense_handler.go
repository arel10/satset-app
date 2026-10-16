package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"satset-app/internal/repository"
	"satset-app/internal/service"
)

// ExpenseHandler handles expense management endpoints.
type ExpenseHandler struct {
	expenseService *service.ExpenseService
}

// NewExpenseHandler creates a new ExpenseHandler.
func NewExpenseHandler(expenseService *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{expenseService: expenseService}
}

// --- Categories ---

// CreateCategory handles POST /api/expense-categories
func (h *ExpenseHandler) CreateCategory(c *gin.Context) {
	var req service.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondBadRequest(c, "Data kategori tidak valid: "+err.Error())
		return
	}

	cat, err := h.expenseService.CreateCategory(c.Request.Context(), req)
	if err != nil {
		RespondBadRequest(c, err.Error())
		return
	}

	RespondCreated(c, "Kategori berhasil dibuat", cat)
}

// ListCategories handles GET /api/expense-categories
func (h *ExpenseHandler) ListCategories(c *gin.Context) {
	activeOnly := c.Query("active") == "true"

	categories, err := h.expenseService.GetCategories(c.Request.Context(), activeOnly)
	if err != nil {
		RespondInternalError(c, "Gagal mengambil kategori")
		return
	}

	RespondOK(c, "", categories)
}

// UpdateCategory handles PUT /api/expense-categories/:id
func (h *ExpenseHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		RespondBadRequest(c, "ID tidak valid")
		return
	}

	var req service.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondBadRequest(c, "Data tidak valid: "+err.Error())
		return
	}

	cat, err := h.expenseService.UpdateCategory(c.Request.Context(), uint(id), req)
	if err != nil {
		RespondBadRequest(c, err.Error())
		return
	}

	RespondOK(c, "Kategori berhasil diperbarui", cat)
}

// --- Expenses ---

// Create handles POST /api/expenses
func (h *ExpenseHandler) Create(c *gin.Context) {
	var req service.CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondBadRequest(c, "Data pengeluaran tidak valid: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")
	ipAddress := c.ClientIP()

	expense, err := h.expenseService.CreateExpense(c.Request.Context(), req, userID, ipAddress)
	if err != nil {
		RespondBadRequest(c, err.Error())
		return
	}

	RespondCreated(c, "Pengeluaran berhasil dicatat", expense)
}

// GetByID handles GET /api/expenses/:id
func (h *ExpenseHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		RespondBadRequest(c, "ID tidak valid")
		return
	}

	expense, err := h.expenseService.GetExpense(c.Request.Context(), uint(id))
	if err != nil {
		RespondNotFound(c, err.Error())
		return
	}

	RespondOK(c, "", expense)
}

// List handles GET /api/expenses
func (h *ExpenseHandler) List(c *gin.Context) {
	page, perPage := parsePagination(c)

	catID := uint(0)
	if cid := c.Query("category_id"); cid != "" {
		if parsed, err := strconv.ParseUint(cid, 10, 32); err == nil {
			catID = uint(parsed)
		}
	}

	filter := repository.ExpenseFilter{
		CategoryID:       catID,
		StartDate:        c.Query("start_date"),
		EndDate:          c.Query("end_date"),
		IncludeCancelled: c.Query("include_cancelled") == "true",
	}

	expenses, total, err := h.expenseService.ListExpenses(c.Request.Context(), filter, page, perPage)
	if err != nil {
		RespondInternalError(c, "Gagal mengambil data pengeluaran")
		return
	}

	meta := NewMeta(page, perPage, total)
	RespondSuccessWithMeta(c, 200, "", expenses, meta)
}

// Update handles PUT /api/expenses/:id
func (h *ExpenseHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		RespondBadRequest(c, "ID tidak valid")
		return
	}

	var req service.UpdateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondBadRequest(c, "Data tidak valid: "+err.Error())
		return
	}

	userID := c.GetUint("user_id")
	ipAddress := c.ClientIP()

	expense, err := h.expenseService.UpdateExpense(c.Request.Context(), uint(id), req, userID, ipAddress)
	if err != nil {
		RespondBadRequest(c, err.Error())
		return
	}

	RespondOK(c, "Pengeluaran berhasil diperbarui", expense)
}

// Cancel handles PUT /api/expenses/:id/cancel
func (h *ExpenseHandler) Cancel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		RespondBadRequest(c, "ID tidak valid")
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondBadRequest(c, "Alasan pembatalan diperlukan")
		return
	}

	userID := c.GetUint("user_id")
	ipAddress := c.ClientIP()

	expense, err := h.expenseService.CancelExpense(c.Request.Context(), uint(id), req.Reason, userID, ipAddress)
	if err != nil {
		RespondBadRequest(c, err.Error())
		return
	}

	RespondOK(c, "Pengeluaran berhasil dibatalkan", expense)
}
