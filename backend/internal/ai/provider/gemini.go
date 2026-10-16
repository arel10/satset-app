package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"satset-app/internal/ai/prompts"
	"satset-app/internal/ai/tools"
)

// Message represents a conversation turn.
type Message struct {
	Role             string                 `json:"role"` // "user", "model", "function"
	Content          string                 `json:"content,omitempty"`
	FunctionName     string                 `json:"function_name,omitempty"`
	FunctionArgs     map[string]interface{} `json:"function_args,omitempty"`
	FunctionID       string                 `json:"function_id,omitempty"`
	ThoughtSignature string                 `json:"thought_signature,omitempty"`
	FunctionResponse map[string]interface{} `json:"function_response,omitempty"`
}

// AIResponse contains the result of an AI generation.
type AIResponse struct {
	Text             string
	ToolCalls        []ToolCallRequest
	PromptTokens     int
	CandidatesTokens int
	TotalTokens      int
}

// ToolCallRequest represents a function call made by the model.
type ToolCallRequest struct {
	Name             string
	Args             map[string]interface{}
	ID               string
	ThoughtSignature string
}

// AIProvider abstracts the AI model interactions.
type AIProvider interface {
	Generate(ctx context.Context, history []Message, toolDefs []tools.ToolDefinition) (*AIResponse, error)
}

// GeminiProvider implements AIProvider using Google's Gemini REST API.
type GeminiProvider struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewGeminiProvider creates a new GeminiProvider instance.
func NewGeminiProvider(apiKey, modelName string, timeoutSecs int) *GeminiProvider {
	if modelName == "" {
		modelName = "gemini-1.5-flash"
	}
	if timeoutSecs <= 0 {
		timeoutSecs = 30
	}

	return &GeminiProvider{
		apiKey: apiKey,
		model:  modelName,
		httpClient: &http.Client{
			Timeout: time.Duration(timeoutSecs) * time.Second,
		},
	}
}

// Internal Gemini REST API payloads
type geminiRequest struct {
	SystemInstruction *geminiContent  `json:"system_instruction,omitempty"`
	Contents          []geminiContent `json:"contents"`
	Tools             []geminiTool    `json:"tools,omitempty"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text             string                  `json:"text,omitempty"`
	ThoughtSignature string                  `json:"thoughtSignature,omitempty"`
	FunctionCall     *geminiFunctionCall     `json:"functionCall,omitempty"`
	FunctionResponse *geminiFunctionResponse `json:"functionResponse,omitempty"`
}

type geminiFunctionCall struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
	ID   string                 `json:"id,omitempty"`
}

type geminiFunctionResponse struct {
	Name     string                 `json:"name"`
	Response map[string]interface{} `json:"response"`
}

type geminiTool struct {
	FunctionDeclarations []tools.ToolDefinition `json:"function_declarations,omitempty"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []geminiPart `json:"parts"`
			Role  string       `json:"role"`
		} `json:"content"`
		FinishReason string `json:"finishReason"`
	} `json:"candidates"`
	UsageMetadata struct {
		PromptTokenCount     int `json:"promptTokenCount"`
		CandidatesTokenCount int `json:"candidatesTokenCount"`
		TotalTokenCount      int `json:"totalTokenCount"`
	} `json:"usageMetadata"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error,omitempty"`
}

// Generate sends messages and tools to Gemini and returns the model output.
func (p *GeminiProvider) Generate(ctx context.Context, history []Message, toolDefs []tools.ToolDefinition) (*AIResponse, error) {
	if p.apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY belum dikonfigurasi pada server")
	}

	reqBody := geminiRequest{
		SystemInstruction: &geminiContent{
			Parts: []geminiPart{{Text: prompts.SystemPrompt}},
		},
		Contents: make([]geminiContent, 0, len(history)),
	}

	for _, m := range history {
		var part geminiPart
		role := m.Role

		if m.FunctionName != "" && m.FunctionResponse != nil {
			role = "user" // In Gemini generateContent API, functionResponse MUST have role "user"
			part = geminiPart{
				FunctionResponse: &geminiFunctionResponse{
					Name:     m.FunctionName,
					Response: m.FunctionResponse,
				},
			}
		} else if m.FunctionName != "" && m.FunctionArgs != nil {
			role = "model"
			part = geminiPart{
				FunctionCall: &geminiFunctionCall{
					Name: m.FunctionName,
					Args: m.FunctionArgs,
					ID:   m.FunctionID,
				},
				ThoughtSignature: m.ThoughtSignature,
			}
		} else {
			if role == "assistant" {
				role = "model"
			}
			part = geminiPart{Text: m.Content}
		}

		reqBody.Contents = append(reqBody.Contents, geminiContent{
			Role:  role,
			Parts: []geminiPart{part},
		})
	}

	if len(toolDefs) > 0 {
		reqBody.Tools = []geminiTool{
			{FunctionDeclarations: toolDefs},
		}
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("gagal menyusun request Gemini: %w", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", p.model, p.apiKey)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("gagal membuat http request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("gagal menghubungi Gemini API: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca response Gemini: %w", err)
	}

	var geminiResp geminiResponse
	if err := json.Unmarshal(bodyBytes, &geminiResp); err != nil {
		return nil, fmt.Errorf("gagal parsing response Gemini: %w", err)
	}

	if geminiResp.Error != nil {
		return nil, fmt.Errorf("gemini API error (%d): %s", geminiResp.Error.Code, geminiResp.Error.Message)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return &AIResponse{Text: "Maaf, tidak ada tanggapan yang dihasilkan."}, nil
	}

	candidate := geminiResp.Candidates[0]
	result := &AIResponse{
		PromptTokens:     geminiResp.UsageMetadata.PromptTokenCount,
		CandidatesTokens: geminiResp.UsageMetadata.CandidatesTokenCount,
		TotalTokens:      geminiResp.UsageMetadata.TotalTokenCount,
	}

	for _, part := range candidate.Content.Parts {
		if part.Text != "" {
			if result.Text != "" {
				result.Text += "\n"
			}
			result.Text += part.Text
		}
		if part.FunctionCall != nil {
			result.ToolCalls = append(result.ToolCalls, ToolCallRequest{
				Name:             part.FunctionCall.Name,
				Args:             part.FunctionCall.Args,
				ID:               part.FunctionCall.ID,
				ThoughtSignature: part.ThoughtSignature,
			})
		}
	}

	return result, nil
}
