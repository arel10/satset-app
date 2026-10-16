package repository

import (
	"context"

	"gorm.io/gorm"

	"satset-app/internal/model"
)

// AIRepository handles database operations for AI conversations and messages.
type AIRepository struct {
	db *gorm.DB
}

// NewAIRepository creates a new AIRepository.
func NewAIRepository(db *gorm.DB) *AIRepository {
	return &AIRepository{db: db}
}

// CreateConversation inserts a new AI conversation.
func (r *AIRepository) CreateConversation(ctx context.Context, conv *model.AIConversation) error {
	return r.db.WithContext(ctx).Create(conv).Error
}

// FindConversationByID retrieves a conversation by ID.
func (r *AIRepository) FindConversationByID(ctx context.Context, id uint) (*model.AIConversation, error) {
	var conv model.AIConversation
	err := r.db.WithContext(ctx).
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		First(&conv, id).Error
	if err != nil {
		return nil, err
	}
	return &conv, nil
}

// FindConversationsByUserID retrieves all conversations for a user (without messages, for list view).
func (r *AIRepository) FindConversationsByUserID(ctx context.Context, userID uint, offset, limit int) ([]model.AIConversation, int64, error) {
	var conversations []model.AIConversation
	var total int64

	err := r.db.WithContext(ctx).Model(&model.AIConversation{}).
		Where("user_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Offset(offset).Limit(limit).
		Find(&conversations).Error
	if err != nil {
		return nil, 0, err
	}

	return conversations, total, nil
}

// CreateMessage inserts a new AI message.
func (r *AIRepository) CreateMessage(ctx context.Context, msg *model.AIMessage) error {
	return r.db.WithContext(ctx).Create(msg).Error
}

// UpdateMessage updates an AI message (e.g., after processing completes).
func (r *AIRepository) UpdateMessage(ctx context.Context, msg *model.AIMessage) error {
	return r.db.WithContext(ctx).Save(msg).Error
}

// UpdateConversation updates a conversation (e.g., title).
func (r *AIRepository) UpdateConversation(ctx context.Context, conv *model.AIConversation) error {
	return r.db.WithContext(ctx).Save(conv).Error
}

// DeleteConversation soft-deletes a conversation and its messages.
func (r *AIRepository) DeleteConversation(ctx context.Context, id uint, userID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete messages first
		if err := tx.Where("conversation_id = ?", id).Delete(&model.AIMessage{}).Error; err != nil {
			return err
		}
		// Delete conversation (only if owned by user)
		return tx.Where("id = ? AND user_id = ?", id, userID).Delete(&model.AIConversation{}).Error
	})
}

// CountTodayMessages counts how many AI questions a user has asked today (for daily limit).
func (r *AIRepository) CountTodayMessages(ctx context.Context, userID uint, todayStart string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.AIMessage{}).
		Joins("JOIN ai_conversations ON ai_conversations.id = ai_messages.conversation_id").
		Where("ai_conversations.user_id = ? AND ai_messages.role = ? AND ai_messages.created_at >= ?",
			userID, model.AIRoleUser, todayStart).
		Count(&count).Error
	return count, err
}
