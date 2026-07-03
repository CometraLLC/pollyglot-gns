package decks

import (
	"time"

	"github.com/google/uuid"
)

// --- Requests ---

type CreateDeckRequest struct {
	Name       string `json:"name" validate:"required,max=100"`
	SourceLang string `json:"source_lang" validate:"required,max=50"`
	TargetLang string `json:"target_lang" validate:"required,max=50"`
}

type UpdateDeckRequest struct {
	Name       string `json:"name" validate:"required,max=100"`
	SourceLang string `json:"source_lang" validate:"required,max=50"`
	TargetLang string `json:"target_lang" validate:"required,max=50"`
}

type CreateCardRequest struct {
	Front    string `json:"front" validate:"required,max=2000"`
	Back     string `json:"back" validate:"required,max=2000"`
	CardType string `json:"card_type" validate:"omitempty,oneof=basic cloze"`
	// Reverse also creates a mirrored (back→front) card. Basic cards only.
	Reverse bool `json:"reverse"`
}

type UpdateCardRequest struct {
	Front string `json:"front" validate:"required,max=2000"`
	Back  string `json:"back" validate:"required,max=2000"`
}

// ReviewCardRequest carries the study rating. Rating is a pointer so
// that 0 (Forgot) survives the required check.
type ReviewCardRequest struct {
	Rating *int `json:"rating" validate:"required,gte=0,lte=4"`
}

// --- Responses ---

type DeckResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	SourceLang string    `json:"source_lang"`
	TargetLang string    `json:"target_lang"`
	CardCount  int64     `json:"card_count"`
	DueCount   int64     `json:"due_count"`
	ShareCode  *string   `json:"share_code,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CardResponse struct {
	ID           uuid.UUID `json:"id"`
	DeckID       uuid.UUID `json:"deck_id"`
	Front        string    `json:"front"`
	Back         string    `json:"back"`
	CardType     string    `json:"card_type"`
	EaseFactor   float64   `json:"ease_factor"`
	IntervalDays int       `json:"interval_days"`
	Repetitions  int       `json:"repetitions"`
	DueAt        time.Time `json:"due_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
