package models

import (
	"time"

	"github.com/google/uuid"
)

// Deck is a user's collection of flashcards for one language pair
type Deck struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null;index"`
	Name       string     `gorm:"type:varchar(100);not null"`
	SourceLang string     `gorm:"type:varchar(50);not null"`
	TargetLang string     `gorm:"type:varchar(50);not null"`
	ShareCode  *string    `gorm:"type:varchar(12);uniqueIndex"`
	CreatedAt  time.Time  `gorm:"not null;default:now()"`
	UpdatedAt  time.Time  `gorm:"not null;default:now()"`
	DeletedAt  *time.Time `gorm:"index"`
}

func (Deck) TableName() string {
	return "decks"
}

// Card types
const (
	CardTypeBasic = "basic"
	CardTypeCloze = "cloze"
)

// Card is a flashcard with its SM-2 spaced-repetition state
type Card struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	DeckID       uuid.UUID  `gorm:"type:uuid;not null;index"`
	Front        string     `gorm:"type:text;not null"`
	Back         string     `gorm:"type:text;not null"`
	CardType     string     `gorm:"type:varchar(10);not null;default:basic"`
	EaseFactor   float64    `gorm:"not null;default:2.5"`
	IntervalDays int        `gorm:"not null;default:0"`
	Repetitions  int        `gorm:"not null;default:0"`
	DueAt        time.Time  `gorm:"not null;default:now()"`
	CreatedAt    time.Time  `gorm:"not null;default:now()"`
	UpdatedAt    time.Time  `gorm:"not null;default:now()"`
	DeletedAt    *time.Time `gorm:"index"`
}

func (Card) TableName() string {
	return "cards"
}
