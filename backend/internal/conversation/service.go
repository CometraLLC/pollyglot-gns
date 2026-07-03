package conversation

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/base-go/backend/internal/shared/models"
	"github.com/base-go/backend/pkg/validator"
)

var ErrConversationNotFound = errors.New("conversation not found")

type Service interface {
	CreateConversation(ctx context.Context, userID uuid.UUID, req CreateConversationRequest) (*ConversationResponse, int, error)
	ListConversations(ctx context.Context, userID uuid.UUID) ([]ConversationResponse, int, error)
	GetMessages(ctx context.Context, userID, conversationID uuid.UUID) ([]MessageResponse, int, error)
	SendMessage(ctx context.Context, userID, conversationID uuid.UUID, req SendMessageRequest) (*ExchangeResponse, int, error)
}

type service struct {
	repo  Repository
	tutor TutorProvider
}

func NewService(repo Repository, tutor TutorProvider) Service {
	return &service{repo: repo, tutor: tutor}
}

func conversationResponse(c *models.Conversation) *ConversationResponse {
	return &ConversationResponse{
		ID:        c.ID,
		Title:     c.Title,
		Language:  c.Language,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func messageResponse(m *models.ConversationMessage) MessageResponse {
	return MessageResponse{
		ID:        m.ID,
		Role:      m.Role,
		Content:   m.Content,
		CreatedAt: m.CreatedAt,
	}
}

// getOwnedConversation enforces ownership; non-existent and non-owned are
// both 404 (no existence leak), matching the decks module.
func (s *service) getOwnedConversation(ctx context.Context, userID, conversationID uuid.UUID) (*models.Conversation, int, error) {
	conversation, err := s.repo.GetConversationByID(ctx, conversationID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, ErrConversationNotFound
		}
		return nil, http.StatusInternalServerError, err
	}
	if conversation.UserID != userID {
		return nil, http.StatusNotFound, ErrConversationNotFound
	}
	return conversation, http.StatusOK, nil
}

func (s *service) CreateConversation(ctx context.Context, userID uuid.UUID, req CreateConversationRequest) (*ConversationResponse, int, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	title := req.Title
	if title == "" {
		title = fmt.Sprintf("Practice %s", req.Language)
	}

	conversation := &models.Conversation{
		UserID:   userID,
		Title:    title,
		Language: req.Language,
	}
	if err := s.repo.CreateConversation(ctx, conversation); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	greeting := &models.ConversationMessage{
		ConversationID: conversation.ID,
		Role:           models.RoleTutor,
		Content:        s.tutor.Greeting(req.Language),
	}
	if err := s.repo.CreateMessage(ctx, greeting); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return conversationResponse(conversation), http.StatusCreated, nil
}

func (s *service) ListConversations(ctx context.Context, userID uuid.UUID) ([]ConversationResponse, int, error) {
	conversations, err := s.repo.GetConversationsByUser(ctx, userID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	result := make([]ConversationResponse, 0, len(conversations))
	for i := range conversations {
		result = append(result, *conversationResponse(&conversations[i]))
	}
	return result, http.StatusOK, nil
}

func (s *service) GetMessages(ctx context.Context, userID, conversationID uuid.UUID) ([]MessageResponse, int, error) {
	if _, status, err := s.getOwnedConversation(ctx, userID, conversationID); err != nil {
		return nil, status, err
	}

	messages, err := s.repo.GetMessagesByConversation(ctx, conversationID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	result := make([]MessageResponse, 0, len(messages))
	for i := range messages {
		result = append(result, messageResponse(&messages[i]))
	}
	return result, http.StatusOK, nil
}

func (s *service) SendMessage(ctx context.Context, userID, conversationID uuid.UUID, req SendMessageRequest) (*ExchangeResponse, int, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	conversation, status, err := s.getOwnedConversation(ctx, userID, conversationID)
	if err != nil {
		return nil, status, err
	}

	history, err := s.repo.GetMessagesByConversation(ctx, conversationID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Ask the provider before persisting anything: a failed exchange
	// must not leave a user message without a reply.
	reply, err := s.tutor.Reply(ctx, conversation.Language, history, req.Content)
	if err != nil {
		return nil, http.StatusBadGateway, err
	}

	userMessage := &models.ConversationMessage{
		ConversationID: conversationID,
		Role:           models.RoleUser,
		Content:        req.Content,
	}
	if err := s.repo.CreateMessage(ctx, userMessage); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	tutorMessage := &models.ConversationMessage{
		ConversationID: conversationID,
		Role:           models.RoleTutor,
		Content:        reply,
	}
	if err := s.repo.CreateMessage(ctx, tutorMessage); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return &ExchangeResponse{
		UserMessage:  messageResponse(userMessage),
		TutorMessage: messageResponse(tutorMessage),
	}, http.StatusOK, nil
}
