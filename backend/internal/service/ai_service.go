package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"satset-app/internal/ai/provider"
	"satset-app/internal/ai/tools"
	"satset-app/internal/audit"
	"satset-app/internal/config"
	"satset-app/internal/model"
	"satset-app/internal/repository"
)

// AIService orchestrates AI assistant conversations and tool execution.
type AIService struct {
	aiRepo       *repository.AIRepository
	provider     provider.AIProvider
	toolExecutor *tools.ToolExecutor
	cfg          *config.Config
	auditLogger  *audit.Logger
}

// NewAIService creates a new AIService instance.
func NewAIService(
	aiRepo *repository.AIRepository,
	prov provider.AIProvider,
	toolExec *tools.ToolExecutor,
	cfg *config.Config,
	auditLogger *audit.Logger,
) *AIService {
	return &AIService{
		aiRepo:       aiRepo,
		provider:     prov,
		toolExecutor: toolExec,
		cfg:          cfg,
		auditLogger:  auditLogger,
	}
}

// ChatRequest contains input for AI conversation.
type ChatRequest struct {
	ConversationID uint   `json:"conversation_id"` // 0 to create new conversation
	Question       string `json:"question"`
	Message        string `json:"message"`
}

// ChatResponse contains the AI output and conversation details.
type ChatResponse struct {
	ConversationID uint      `json:"conversation_id"`
	Title          string    `json:"title"`
	Answer         string    `json:"answer"`
	Message        string    `json:"message"` // frontend compatibility alias
	ToolCallsMade  int       `json:"tool_calls_made"`
	TokensUsed     int       `json:"tokens_used"`
	CreatedAt      time.Time `json:"created_at"`
}

// Chat handles a user's question, manages context, executes tools, and saves history.
func (s *AIService) Chat(ctx context.Context, userID uint, req ChatRequest, ipAddress string) (*ChatResponse, error) {
	// 1. Resolve and validate question length
	targetQuestion := strings.TrimSpace(req.Question)
	if targetQuestion == "" {
		targetQuestion = strings.TrimSpace(req.Message)
	}
	if len(targetQuestion) == 0 {
		return nil, fmt.Errorf("pertanyaan tidak boleh kosong")
	}
	req.Question = targetQuestion
	if s.cfg.AI.MaxQuestionLength > 0 && len(req.Question) > s.cfg.AI.MaxQuestionLength {
		return nil, fmt.Errorf("pertanyaan terlalu panjang (maksimal %d karakter)", s.cfg.AI.MaxQuestionLength)
	}

	// 2. Daily usage limit check
	today := time.Now().Format("2006-01-02")
	count, err := s.aiRepo.CountTodayMessages(ctx, userID, today)
	if err != nil {
		return nil, fmt.Errorf("gagal memeriksa batas penggunaan AI: %w", err)
	}
	if s.cfg.AI.DailyLimitPerUser > 0 && int(count) >= s.cfg.AI.DailyLimitPerUser {
		return nil, fmt.Errorf("batas pertanyaan harian (%d) telah tercapai. Coba lagi besok.", s.cfg.AI.DailyLimitPerUser)
	}

	// 3. Get or create conversation
	var conv *model.AIConversation
	if req.ConversationID > 0 {
		existingConv, err := s.aiRepo.FindConversationByID(ctx, req.ConversationID)
		if err != nil {
			return nil, fmt.Errorf("percakapan tidak ditemukan")
		}
		if existingConv.UserID != userID {
			return nil, fmt.Errorf("akses ditolak untuk percakapan ini")
		}
		conv = existingConv
	} else {
		title := req.Question
		if len(title) > 40 {
			title = title[:40] + "..."
		}
		conv = &model.AIConversation{
			UserID: userID,
			Title:  title,
		}
		if err := s.aiRepo.CreateConversation(ctx, conv); err != nil {
			return nil, fmt.Errorf("gagal membuat percakapan baru: %w", err)
		}
	}

	// 4. Save User Message
	userMsg := &model.AIMessage{
		ConversationID: conv.ID,
		Role:           model.AIRoleUser,
		Content:        req.Question,
	}
	if err := s.aiRepo.CreateMessage(ctx, userMsg); err != nil {
		return nil, fmt.Errorf("gagal menyimpan pesan: %w", err)
	}

	// 5. Build conversation history
	history := make([]provider.Message, 0, len(conv.Messages)+1)
	for _, m := range conv.Messages {
		history = append(history, provider.Message{
			Role:    m.Role,
			Content: m.Content,
		})
	}
	history = append(history, provider.Message{
		Role:    "user",
		Content: req.Question,
	})

	// 6. Tool execution loop (max iterations guarded)
	toolDefs := tools.GetToolDefinitions()
	maxIterations := s.cfg.AI.MaxToolCalls
	if maxIterations <= 0 {
		maxIterations = 5
	}

	totalTokens := 0
	toolCallsCount := 0
	var finalAnswer string

	for i := 0; i < maxIterations; i++ {
		resp, err := s.provider.Generate(ctx, history, toolDefs)
		if err != nil {
			return nil, fmt.Errorf("asisten AI mengalami kendala: %w", err)
		}

		totalTokens += resp.TotalTokens

		if len(resp.ToolCalls) == 0 {
			finalAnswer = resp.Text
			break
		}

		// Handle tool calls
		for _, tc := range resp.ToolCalls {
			toolCallsCount++
			argsBytes, _ := json.Marshal(tc.Args)
			toolResult, execErr := s.toolExecutor.Execute(ctx, tc.Name, string(argsBytes))
			var resultPayload map[string]interface{}
			if execErr != nil {
				resultPayload = map[string]interface{}{"error": execErr.Error()}
			} else {
				if m, ok := toolResult.(map[string]interface{}); ok {
					resultPayload = m
				} else {
					resultPayload = map[string]interface{}{"result": toolResult}
				}
			}

			// Add assistant call and function response to history
			history = append(history, provider.Message{
				Role:             "model",
				FunctionName:     tc.Name,
				FunctionArgs:     tc.Args,
				FunctionID:       tc.ID,
				ThoughtSignature: tc.ThoughtSignature,
			})
			history = append(history, provider.Message{
				Role:             "user",
				FunctionName:     tc.Name,
				FunctionResponse: resultPayload,
			})
		}
	}

	if finalAnswer == "" {
		finalAnswer = "Maaf, Asisten AI membutuhkan terlalu banyak langkah untuk menjawab. Silakan coba tanyakan lebih spesifik."
	}

	// 7. Save Assistant Message
	metaBytes, _ := json.Marshal(map[string]interface{}{
		"tools_called": toolCallsCount,
		"tokens_used":  totalTokens,
	})
	metaStr := string(metaBytes)
	assistantMsg := &model.AIMessage{
		ConversationID:   conv.ID,
		Role:             model.AIRoleAssistant,
		Content:          finalAnswer,
		ToolCallsMeta:    &metaStr,
		ProcessingStatus: model.AIStatusCompleted,
	}
	_ = s.aiRepo.CreateMessage(ctx, assistantMsg)

	// Update conversation timestamp
	_ = s.aiRepo.UpdateConversation(ctx, conv)

	// 8. Audit log
	s.auditLogger.Log(userID, model.AuditActionAIChat, model.AuditEntityAIConversation, conv.ID, nil, map[string]interface{}{
		"tools_called": toolCallsCount,
		"tokens_used":  totalTokens,
	}, ipAddress)

	return &ChatResponse{
		ConversationID: conv.ID,
		Title:          conv.Title,
		Answer:         finalAnswer,
		Message:        finalAnswer,
		ToolCallsMade:  toolCallsCount,
		TokensUsed:     totalTokens,
		CreatedAt:      assistantMsg.CreatedAt,
	}, nil
}

// GetConversations retrieves paginated conversation history for the user.
func (s *AIService) GetConversations(ctx context.Context, userID uint, page, perPage int) ([]model.AIConversation, int64, error) {
	offset := (page - 1) * perPage
	return s.aiRepo.FindConversationsByUserID(ctx, userID, offset, perPage)
}

// GetConversationDetail retrieves a single conversation with all its messages.
func (s *AIService) GetConversationDetail(ctx context.Context, userID, convID uint) (*model.AIConversation, error) {
	conv, err := s.aiRepo.FindConversationByID(ctx, convID)
	if err != nil {
		return nil, fmt.Errorf("percakapan tidak ditemukan")
	}
	if conv.UserID != userID {
		return nil, fmt.Errorf("akses ditolak")
	}
	return conv, nil
}

// DeleteConversation deletes a conversation.
func (s *AIService) DeleteConversation(ctx context.Context, userID, convID uint) error {
	return s.aiRepo.DeleteConversation(ctx, convID, userID)
}
