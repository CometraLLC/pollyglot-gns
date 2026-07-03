package models

import (
	"time"

	"github.com/google/uuid"
)

// Message roles in a conversation
const (
	RoleUser  = "user"
	RoleTutor = "tutor"
)

// Conversation is a practice session between a user and the tutor
type Conversation struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index"`
	Title     string     `gorm:"type:varchar(100);not null"`
	Language  string     `gorm:"type:varchar(50);not null"`
	CreatedAt time.Time  `gorm:"not null;default:now()"`
	UpdatedAt time.Time  `gorm:"not null;default:now()"`
	DeletedAt *time.Time `gorm:"index"`
}

func (Conversation) TableName() string {
	return "conversations"
}

// ConversationMessage is one turn in a conversation
type ConversationMessage struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ConversationID uuid.UUID `gorm:"type:uuid;not null;index"`
	Role           string    `gorm:"type:varchar(10);not null"`
	Content        string    `gorm:"type:text;not null"`
	CreatedAt      time.Time `gorm:"not null;default:now()"`
}

func (ConversationMessage) TableName() string {
	return "conversation_messages"
}
