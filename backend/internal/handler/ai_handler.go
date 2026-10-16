package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"satset-app/internal/ai/tools"
	"satset-app/internal/service"
)

// AIHandler handles AI Assistant endpoints (Owner only).
type AIHandler struct {
	aiService *service.AIService
}

// NewAIHandler creates a new AIHandler instance.
func NewAIHandler(aiService *service.AIService) *AIHandler {
	return &AIHandler{aiService: aiService}
}

// Chat handles POST /api/ai/chat
func (h *AIHandler) Chat(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req service.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		RespondBadRequest(c, "Permintaan tidak valid: "+err.Error())
		return
	}

	ip := c.ClientIP()
	resp, err := h.aiService.Chat(c.Request.Context(), userID, req, ip)
	if err != nil {
		RespondInternalError(c, err.Error())
		return
	}

	RespondOK(c, "Jawaban AI berhasil diperoleh", resp)
}

// GetConversations handles GET /api/ai/conversations
func (h *AIHandler) GetConversations(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, perPage := parsePagination(c)

	conversations, total, err := h.aiService.GetConversations(c.Request.Context(), userID, page, perPage)
	if err != nil {
		RespondInternalError(c, "Gagal mengambil riwayat percakapan")
		return
	}

	meta := NewMeta(page, perPage, total)
	RespondSuccessWithMeta(c, 200, "", conversations, meta)
}

// GetConversationDetail handles GET /api/ai/conversations/:id
func (h *AIHandler) GetConversationDetail(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		RespondBadRequest(c, "ID percakapan tidak valid")
		return
	}

	conv, err := h.aiService.GetConversationDetail(c.Request.Context(), userID, uint(id))
	if err != nil {
		RespondNotFound(c, err.Error())
		return
	}

	RespondOK(c, "", conv)
}

// DeleteConversation handles DELETE /api/ai/conversations/:id
func (h *AIHandler) DeleteConversation(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		RespondBadRequest(c, "ID percakapan tidak valid")
		return
	}

	if err := h.aiService.DeleteConversation(c.Request.Context(), userID, uint(id)); err != nil {
		RespondInternalError(c, "Gagal menghapus percakapan: "+err.Error())
		return
	}

	RespondOK(c, "Percakapan berhasil dihapus", nil)
}

// ListTools handles GET /api/ai/tools
func (h *AIHandler) ListTools(c *gin.Context) {
	toolDefs := tools.GetToolDefinitions()
	RespondOK(c, "", toolDefs)
}
