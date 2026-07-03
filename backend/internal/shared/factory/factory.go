// Package factory provides chainable test-data factories for the domain
// models, plus the fixed identifiers of the development seed data.
// Use in tests only.
package factory

import (
	"time"

	"github.com/google/uuid"

	"github.com/base-go/backend/internal/shared/models"
)

// Seeded mirrors backend/migrations/seeders/dev — the account and deck
// available in every development environment (see README).
var Seeded = struct {
	UserID   uuid.UUID
	DeckID   uuid.UUID
	Email    string
	Password string
}{
	UserID:   uuid.MustParse("a0000000-0000-4000-8000-000000000001"),
	DeckID:   uuid.MustParse("b0000000-0000-4000-8000-000000000001"),
	Email:    "demo@pollyglot.dev",
	Password: "Password123!",
}

// --- UserFactory ---

type UserFactory struct {
	user models.User
}

func User() *UserFactory {
	now := time.Now()
	return &UserFactory{user: models.User{
		ID:            uuid.New(),
		Email:         "test@pollyglot.dev",
		Name:          "Test User",
		IsActive:      true,
		EmailVerified: true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}}
}

func (f *UserFactory) WithID(id uuid.UUID) *UserFactory    { f.user.ID = id; return f }
func (f *UserFactory) WithEmail(email string) *UserFactory { f.user.Email = email; return f }
func (f *UserFactory) WithName(name string) *UserFactory   { f.user.Name = name; return f }
func (f *UserFactory) Inactive() *UserFactory              { f.user.IsActive = false; return f }
func (f *UserFactory) Build() models.User                  { return f.user }

// --- DeckFactory ---

type DeckFactory struct {
	deck models.Deck
}

func Deck() *DeckFactory {
	now := time.Now()
	return &DeckFactory{deck: models.Deck{
		ID:         uuid.New(),
		UserID:     uuid.New(),
		Name:       "Japanese Basics",
		SourceLang: "Japanese",
		TargetLang: "English",
		CreatedAt:  now,
		UpdatedAt:  now,
	}}
}

func (f *DeckFactory) WithID(id uuid.UUID) *DeckFactory     { f.deck.ID = id; return f }
func (f *DeckFactory) WithUserID(id uuid.UUID) *DeckFactory { f.deck.UserID = id; return f }
func (f *DeckFactory) WithName(name string) *DeckFactory    { f.deck.Name = name; return f }
func (f *DeckFactory) WithLanguages(src, tgt string) *DeckFactory {
	f.deck.SourceLang, f.deck.TargetLang = src, tgt
	return f
}
func (f *DeckFactory) Build() models.Deck { return f.deck }

// --- CardFactory ---

type CardFactory struct {
	card models.Card
}

func Card() *CardFactory {
	now := time.Now()
	return &CardFactory{card: models.Card{
		ID:           uuid.New(),
		DeckID:       uuid.New(),
		Front:        "こんにちは",
		Back:         "hello",
		EaseFactor:   2.5,
		IntervalDays: 0,
		Repetitions:  0,
		DueAt:        now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}}
}

func (f *CardFactory) WithID(id uuid.UUID) *CardFactory     { f.card.ID = id; return f }
func (f *CardFactory) WithDeckID(id uuid.UUID) *CardFactory { f.card.DeckID = id; return f }
func (f *CardFactory) WithFront(front string) *CardFactory  { f.card.Front = front; return f }
func (f *CardFactory) WithBack(back string) *CardFactory    { f.card.Back = back; return f }

// WithSRS sets the scheduling state (ease factor, interval days, repetitions).
func (f *CardFactory) WithSRS(ease float64, intervalDays, repetitions int) *CardFactory {
	f.card.EaseFactor = ease
	f.card.IntervalDays = intervalDays
	f.card.Repetitions = repetitions
	return f
}

func (f *CardFactory) DueAt(t time.Time) *CardFactory { f.card.DueAt = t; return f }
func (f *CardFactory) Build() models.Card             { return f.card }

// --- ReviewFactory ---

type ReviewFactory struct {
	review models.Review
}

func Review() *ReviewFactory {
	return &ReviewFactory{review: models.Review{
		ID:         uuid.New(),
		CardID:     uuid.New(),
		UserID:     uuid.New(),
		Rating:     3,
		ReviewedAt: time.Now(),
	}}
}

func (f *ReviewFactory) WithCardID(id uuid.UUID) *ReviewFactory { f.review.CardID = id; return f }
func (f *ReviewFactory) WithUserID(id uuid.UUID) *ReviewFactory { f.review.UserID = id; return f }
func (f *ReviewFactory) WithRating(rating int) *ReviewFactory   { f.review.Rating = rating; return f }
func (f *ReviewFactory) ReviewedAt(t time.Time) *ReviewFactory  { f.review.ReviewedAt = t; return f }
func (f *ReviewFactory) Build() models.Review                   { return f.review }

// --- ConversationFactory ---

type ConversationFactory struct {
	conversation models.Conversation
}

func Conversation() *ConversationFactory {
	now := time.Now()
	return &ConversationFactory{conversation: models.Conversation{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Title:     "Practice Japanese",
		Language:  "Japanese",
		CreatedAt: now,
		UpdatedAt: now,
	}}
}

func (f *ConversationFactory) WithID(id uuid.UUID) *ConversationFactory {
	f.conversation.ID = id
	return f
}
func (f *ConversationFactory) WithUserID(id uuid.UUID) *ConversationFactory {
	f.conversation.UserID = id
	return f
}
func (f *ConversationFactory) WithLanguage(language string) *ConversationFactory {
	f.conversation.Language = language
	return f
}
func (f *ConversationFactory) Build() models.Conversation { return f.conversation }

// --- MessageFactory ---

type MessageFactory struct {
	message models.ConversationMessage
}

func Message() *MessageFactory {
	return &MessageFactory{message: models.ConversationMessage{
		ID:             uuid.New(),
		ConversationID: uuid.New(),
		Role:           models.RoleUser,
		Content:        "こんにちは!",
		CreatedAt:      time.Now(),
	}}
}

func (f *MessageFactory) WithConversationID(id uuid.UUID) *MessageFactory {
	f.message.ConversationID = id
	return f
}
func (f *MessageFactory) FromTutor() *MessageFactory {
	f.message.Role = models.RoleTutor
	return f
}
func (f *MessageFactory) WithContent(content string) *MessageFactory {
	f.message.Content = content
	return f
}
func (f *MessageFactory) Build() models.ConversationMessage { return f.message }
