package conversation

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/base-go/backend/internal/shared/models"
	"github.com/base-go/backend/pkg/database"
)

type Repository interface {
	CreateConversation(ctx context.Context, conversation *models.Conversation) error
	GetConversationsByUser(ctx context.Context, userID uuid.UUID) ([]models.Conversation, error)
	GetConversationByID(ctx context.Context, id uuid.UUID) (*models.Conversation, error)

	CreateMessage(ctx context.Context, message *models.ConversationMessage) error
	GetMessagesByConversation(ctx context.Context, conversationID uuid.UUID) ([]models.ConversationMessage, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db database.Database) Repository {
	return &repository{db: db.GetDB()}
}

func (r *repository) CreateConversation(ctx context.Context, conversation *models.Conversation) error {
	return r.db.WithContext(ctx).Create(conversation).Error
}

func (r *repository) GetConversationsByUser(ctx context.Context, userID uuid.UUID) ([]models.Conversation, error) {
	var conversations []models.Conversation
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("updated_at DESC").
		Find(&conversations).Error
	return conversations, err
}

func (r *repository) GetConversationByID(ctx context.Context, id uuid.UUID) (*models.Conversation, error) {
	var conversation models.Conversation
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&conversation).Error
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *repository) CreateMessage(ctx context.Context, message *models.ConversationMessage) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *repository) GetMessagesByConversation(ctx context.Context, conversationID uuid.UUID) ([]models.ConversationMessage, error) {
	var messages []models.ConversationMessage
	err := r.db.WithContext(ctx).
		Where("conversation_id = ?", conversationID).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}
