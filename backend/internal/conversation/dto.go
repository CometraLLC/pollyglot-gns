package conversation

import (
	"time"

	"github.com/google/uuid"
)

type CreateConversationRequest struct {
	Title    string `json:"title" validate:"omitempty,max=100"`
	Language string `json:"language" validate:"required,max=50"`
}

type SendMessageRequest struct {
	Content string `json:"content" validate:"required,max=2000"`
}

type ConversationResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Language  string    `json:"language"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MessageResponse struct {
	ID        uuid.UUID `json:"id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// ExchangeResponse is one full turn: the learner's message and the
// tutor's reply, in the order they were persisted.
type ExchangeResponse struct {
	UserMessage  MessageResponse `json:"user_message"`
	TutorMessage MessageResponse `json:"tutor_message"`
}
