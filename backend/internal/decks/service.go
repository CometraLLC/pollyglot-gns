package decks

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/base-go/backend/internal/shared/models"
	"github.com/base-go/backend/pkg/validator"
)

var (
	ErrDeckNotFound = errors.New("deck not found")
	ErrCardNotFound = errors.New("card not found")
)

type Service interface {
	CreateDeck(ctx context.Context, userID uuid.UUID, req CreateDeckRequest) (*DeckResponse, int, error)
	ListDecks(ctx context.Context, userID uuid.UUID) ([]DeckResponse, int, error)
	GetDeck(ctx context.Context, userID, deckID uuid.UUID) (*DeckResponse, int, error)
	UpdateDeck(ctx context.Context, userID, deckID uuid.UUID, req UpdateDeckRequest) (*DeckResponse, int, error)
	DeleteDeck(ctx context.Context, userID, deckID uuid.UUID) (int, error)

	CreateCard(ctx context.Context, userID, deckID uuid.UUID, req CreateCardRequest) (*CardResponse, int, error)
	ListCards(ctx context.Context, userID, deckID uuid.UUID) ([]CardResponse, int, error)
	UpdateCard(ctx context.Context, userID, cardID uuid.UUID, req UpdateCardRequest) (*CardResponse, int, error)
	DeleteCard(ctx context.Context, userID, cardID uuid.UUID) (int, error)

	ReviewCard(ctx context.Context, userID, cardID uuid.UUID, req ReviewCardRequest) (*CardResponse, int, error)
	GetStudyQueue(ctx context.Context, userID, deckID uuid.UUID, limit int) ([]CardResponse, int, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// getOwnedDeck loads a deck and enforces ownership. Non-existent and
// non-owned decks are indistinguishable to the caller (404, no leak).
func (s *service) getOwnedDeck(ctx context.Context, userID, deckID uuid.UUID) (*models.Deck, int, error) {
	deck, err := s.repo.GetDeckByID(ctx, deckID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, ErrDeckNotFound
		}
		return nil, http.StatusInternalServerError, err
	}
	if deck.UserID != userID {
		return nil, http.StatusNotFound, ErrDeckNotFound
	}
	return deck, http.StatusOK, nil
}

// getOwnedCard loads a card and enforces ownership through its deck.
func (s *service) getOwnedCard(ctx context.Context, userID, cardID uuid.UUID) (*models.Card, int, error) {
	card, err := s.repo.GetCardByID(ctx, cardID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, ErrCardNotFound
		}
		return nil, http.StatusInternalServerError, err
	}
	if _, status, err := s.getOwnedDeck(ctx, userID, card.DeckID); err != nil {
		if status == http.StatusNotFound {
			return nil, http.StatusNotFound, ErrCardNotFound
		}
		return nil, status, err
	}
	return card, http.StatusOK, nil
}

func deckResponse(deck *models.Deck, cardCount int64) *DeckResponse {
	return &DeckResponse{
		ID:         deck.ID,
		Name:       deck.Name,
		SourceLang: deck.SourceLang,
		TargetLang: deck.TargetLang,
		CardCount:  cardCount,
		CreatedAt:  deck.CreatedAt,
		UpdatedAt:  deck.UpdatedAt,
	}
}

func cardResponse(card *models.Card) *CardResponse {
	return &CardResponse{
		ID:           card.ID,
		DeckID:       card.DeckID,
		Front:        card.Front,
		Back:         card.Back,
		EaseFactor:   card.EaseFactor,
		IntervalDays: card.IntervalDays,
		Repetitions:  card.Repetitions,
		DueAt:        card.DueAt,
		CreatedAt:    card.CreatedAt,
		UpdatedAt:    card.UpdatedAt,
	}
}

// --- deck operations ---

func (s *service) CreateDeck(ctx context.Context, userID uuid.UUID, req CreateDeckRequest) (*DeckResponse, int, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	deck := &models.Deck{
		UserID:     userID,
		Name:       req.Name,
		SourceLang: req.SourceLang,
		TargetLang: req.TargetLang,
	}
	if err := s.repo.CreateDeck(ctx, deck); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return deckResponse(deck, 0), http.StatusCreated, nil
}

func (s *service) ListDecks(ctx context.Context, userID uuid.UUID) ([]DeckResponse, int, error) {
	decks, err := s.repo.GetDecksByUser(ctx, userID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	result := make([]DeckResponse, 0, len(decks))
	for _, d := range decks {
		deck := d.Deck
		result = append(result, *deckResponse(&deck, d.CardCount))
	}
	return result, http.StatusOK, nil
}

func (s *service) GetDeck(ctx context.Context, userID, deckID uuid.UUID) (*DeckResponse, int, error) {
	deck, status, err := s.getOwnedDeck(ctx, userID, deckID)
	if err != nil {
		return nil, status, err
	}

	count, err := s.repo.CountCards(ctx, deckID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return deckResponse(deck, count), http.StatusOK, nil
}

func (s *service) UpdateDeck(ctx context.Context, userID, deckID uuid.UUID, req UpdateDeckRequest) (*DeckResponse, int, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	deck, status, err := s.getOwnedDeck(ctx, userID, deckID)
	if err != nil {
		return nil, status, err
	}

	deck.Name = req.Name
	deck.SourceLang = req.SourceLang
	deck.TargetLang = req.TargetLang
	deck.UpdatedAt = time.Now()
	if err := s.repo.UpdateDeck(ctx, deck); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	count, err := s.repo.CountCards(ctx, deckID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return deckResponse(deck, count), http.StatusOK, nil
}

func (s *service) DeleteDeck(ctx context.Context, userID, deckID uuid.UUID) (int, error) {
	if _, status, err := s.getOwnedDeck(ctx, userID, deckID); err != nil {
		return status, err
	}

	if err := s.repo.SoftDeleteDeck(ctx, deckID); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

// --- card operations ---

func (s *service) CreateCard(ctx context.Context, userID, deckID uuid.UUID, req CreateCardRequest) (*CardResponse, int, error) {
	if _, status, err := s.getOwnedDeck(ctx, userID, deckID); err != nil {
		return nil, status, err
	}

	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	card := &models.Card{
		DeckID:       deckID,
		Front:        req.Front,
		Back:         req.Back,
		EaseFactor:   2.5,
		IntervalDays: 0,
		Repetitions:  0,
		DueAt:        time.Now(),
	}
	if err := s.repo.CreateCard(ctx, card); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return cardResponse(card), http.StatusCreated, nil
}

func (s *service) ListCards(ctx context.Context, userID, deckID uuid.UUID) ([]CardResponse, int, error) {
	if _, status, err := s.getOwnedDeck(ctx, userID, deckID); err != nil {
		return nil, status, err
	}

	cards, err := s.repo.GetCardsByDeck(ctx, deckID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	result := make([]CardResponse, 0, len(cards))
	for i := range cards {
		result = append(result, *cardResponse(&cards[i]))
	}
	return result, http.StatusOK, nil
}

func (s *service) UpdateCard(ctx context.Context, userID, cardID uuid.UUID, req UpdateCardRequest) (*CardResponse, int, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	card, status, err := s.getOwnedCard(ctx, userID, cardID)
	if err != nil {
		return nil, status, err
	}

	card.Front = req.Front
	card.Back = req.Back
	card.UpdatedAt = time.Now()
	if err := s.repo.UpdateCard(ctx, card); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return cardResponse(card), http.StatusOK, nil
}

func (s *service) DeleteCard(ctx context.Context, userID, cardID uuid.UUID) (int, error) {
	if _, status, err := s.getOwnedCard(ctx, userID, cardID); err != nil {
		return status, err
	}

	if err := s.repo.SoftDeleteCard(ctx, cardID); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
