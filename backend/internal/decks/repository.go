package decks

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/base-go/backend/internal/shared/models"
	"github.com/base-go/backend/pkg/database"
)

// DeckWithCount pairs a deck with its (non-deleted) card count for list views
type DeckWithCount struct {
	models.Deck
	CardCount int64
	DueCount  int64
}

type Repository interface {
	CreateDeck(ctx context.Context, deck *models.Deck) error
	GetDecksByUser(ctx context.Context, userID uuid.UUID) ([]DeckWithCount, error)
	GetDeckByID(ctx context.Context, id uuid.UUID) (*models.Deck, error)
	CountCards(ctx context.Context, deckID uuid.UUID) (int64, error)
	CountDueCards(ctx context.Context, deckID uuid.UUID, before time.Time) (int64, error)
	UpdateDeck(ctx context.Context, deck *models.Deck) error
	SoftDeleteDeck(ctx context.Context, id uuid.UUID) error

	CreateCard(ctx context.Context, card *models.Card) error
	GetCardsByDeck(ctx context.Context, deckID uuid.UUID) ([]models.Card, error)
	GetCardByID(ctx context.Context, id uuid.UUID) (*models.Card, error)
	UpdateCard(ctx context.Context, card *models.Card) error
	SoftDeleteCard(ctx context.Context, id uuid.UUID) error

	GetDueCards(ctx context.Context, deckID uuid.UUID, before time.Time, limit int) ([]models.Card, error)
	CreateReview(ctx context.Context, review *models.Review) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db database.Database) Repository {
	return &repository{db: db.GetDB()}
}

func (r *repository) CreateDeck(ctx context.Context, deck *models.Deck) error {
	return r.db.WithContext(ctx).Create(deck).Error
}

func (r *repository) GetDecksByUser(ctx context.Context, userID uuid.UUID) ([]DeckWithCount, error) {
	var decks []models.Deck
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Order("created_at DESC").
		Find(&decks).Error
	if err != nil {
		return nil, err
	}

	now := time.Now()
	result := make([]DeckWithCount, 0, len(decks))
	for _, deck := range decks {
		count, err := r.CountCards(ctx, deck.ID)
		if err != nil {
			return nil, err
		}
		due, err := r.CountDueCards(ctx, deck.ID, now)
		if err != nil {
			return nil, err
		}
		result = append(result, DeckWithCount{Deck: deck, CardCount: count, DueCount: due})
	}
	return result, nil
}

func (r *repository) GetDeckByID(ctx context.Context, id uuid.UUID) (*models.Deck, error) {
	var deck models.Deck
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&deck).Error
	if err != nil {
		return nil, err
	}
	return &deck, nil
}

func (r *repository) CountCards(ctx context.Context, deckID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Card{}).
		Where("deck_id = ? AND deleted_at IS NULL", deckID).
		Count(&count).Error
	return count, err
}

func (r *repository) CountDueCards(ctx context.Context, deckID uuid.UUID, before time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Card{}).
		Where("deck_id = ? AND deleted_at IS NULL AND due_at <= ?", deckID, before).
		Count(&count).Error
	return count, err
}

func (r *repository) UpdateDeck(ctx context.Context, deck *models.Deck) error {
	return r.db.WithContext(ctx).Save(deck).Error
}

func (r *repository) SoftDeleteDeck(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Deck{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *repository) CreateCard(ctx context.Context, card *models.Card) error {
	return r.db.WithContext(ctx).Create(card).Error
}

func (r *repository) GetCardsByDeck(ctx context.Context, deckID uuid.UUID) ([]models.Card, error) {
	var cards []models.Card
	err := r.db.WithContext(ctx).
		Where("deck_id = ? AND deleted_at IS NULL", deckID).
		Order("created_at ASC").
		Find(&cards).Error
	return cards, err
}

func (r *repository) GetCardByID(ctx context.Context, id uuid.UUID) (*models.Card, error) {
	var card models.Card
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&card).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (r *repository) UpdateCard(ctx context.Context, card *models.Card) error {
	return r.db.WithContext(ctx).Save(card).Error
}

func (r *repository) SoftDeleteCard(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Card{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *repository) GetDueCards(ctx context.Context, deckID uuid.UUID, before time.Time, limit int) ([]models.Card, error) {
	var cards []models.Card
	err := r.db.WithContext(ctx).
		Where("deck_id = ? AND deleted_at IS NULL AND due_at <= ?", deckID, before).
		Order("due_at ASC").
		Limit(limit).
		Find(&cards).Error
	return cards, err
}

func (r *repository) CreateReview(ctx context.Context, review *models.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}
