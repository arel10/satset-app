package model

import (
	"time"
)

// AIConversation represents a chat session between Owner and AI Assistant.
type AIConversation struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	UserID    uint        `gorm:"index;not null" json:"user_id"`
	User      User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Title     string      `gorm:"size:200" json:"title,omitempty"`
	Messages  []AIMessage `gorm:"foreignKey:ConversationID" json:"messages,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// TableName specifies the database table name for AIConversation.
func (AIConversation) TableName() string {
	return "ai_conversations"
}

// AIMessage represents a single message in an AI conversation.
type AIMessage struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	ConversationID   uint      `gorm:"index;not null" json:"conversation_id"`
	Role             string    `gorm:"size:20;not null" json:"role"` // user, assistant
	Content          string    `gorm:"type:longtext;not null" json:"content"`
	ToolCallsMeta    *string   `gorm:"type:json" json:"tool_calls_meta,omitempty"`
	ProcessingStatus string    `gorm:"size:20;default:completed" json:"processing_status"`
	CreatedAt        time.Time `json:"created_at"`
}

// TableName specifies the database table name for AIMessage.
func (AIMessage) TableName() string {
	return "ai_messages"
}

// AI message roles
const (
	AIRoleUser      = "user"
	AIRoleAssistant = "assistant"
)

// AI processing statuses
const (
	AIStatusProcessing = "processing"
	AIStatusCompleted  = "completed"
	AIStatusFailed     = "failed"
)
