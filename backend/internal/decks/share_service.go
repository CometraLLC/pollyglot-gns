package decks

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"net/http"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/base-go/backend/internal/shared/models"
)

var ErrSharedDeckNotFound = errors.New("shared deck not found")

// shareAlphabet omits ambiguous characters (0/O, 1/I/L).
const shareAlphabet = "ABCDEFGHJKMNPQRSTUVWXYZ23456789"
const shareCodeLength = 10

// GenerateShareCode returns a cryptographically random, human-friendly code.
func GenerateShareCode() (string, error) {
	code := make([]byte, shareCodeLength)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(shareAlphabet))))
		if err != nil {
			return "", err
		}
		code[i] = shareAlphabet[n.Int64()]
	}
	return string(code), nil
}

type ShareResponse struct {
	ShareCode string `json:"share_code"`
}

type SampleCard struct {
	Front string `json:"front"`
	Back  string `json:"back"`
}

type SharedDeckPreview struct {
	Name        string       `json:"name"`
	SourceLang  string       `json:"source_lang"`
	TargetLang  string       `json:"target_lang"`
	CardCount   int64        `json:"card_count"`
	SampleCards []SampleCard `json:"sample_cards"`
}

// ShareDeck enables sharing, keeping an existing code (idempotent).
func (s *service) ShareDeck(ctx context.Context, userID, deckID uuid.UUID) (*ShareResponse, int, error) {
	deck, status, err := s.getOwnedDeck(ctx, userID, deckID)
	if err != nil {
		return nil, status, err
	}

	if deck.ShareCode != nil {
		return &ShareResponse{ShareCode: *deck.ShareCode}, http.StatusOK, nil
	}

	// retry a few times in case the unique index catches a collision
	for range 3 {
		code, err := GenerateShareCode()
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		if err := s.repo.SetShareCode(ctx, deckID, &code); err != nil {
			continue
		}
		return &ShareResponse{ShareCode: code}, http.StatusOK, nil
	}
	return nil, http.StatusInternalServerError, errors.New("could not allocate a share code")
}

// UnshareDeck disables sharing.
func (s *service) UnshareDeck(ctx context.Context, userID, deckID uuid.UUID) (int, error) {
	if _, status, err := s.getOwnedDeck(ctx, userID, deckID); err != nil {
		return status, err
	}

	if err := s.repo.SetShareCode(ctx, deckID, nil); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

// GetSharedDeck previews a shared deck for any authenticated user.
func (s *service) GetSharedDeck(ctx context.Context, code string) (*SharedDeckPreview, int, error) {
	deck, err := s.repo.GetDeckByShareCode(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, ErrSharedDeckNotFound
		}
		return nil, http.StatusInternalServerError, err
	}

	cards, err := s.repo.GetCardsByDeck(ctx, deck.ID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	samples := make([]SampleCard, 0, 5)
	for i, card := range cards {
		if i == 5 {
			break
		}
		samples = append(samples, SampleCard{Front: card.Front, Back: card.Back})
	}

	return &SharedDeckPreview{
		Name:        deck.Name,
		SourceLang:  deck.SourceLang,
		TargetLang:  deck.TargetLang,
		CardCount:   int64(len(cards)),
		SampleCards: samples,
	}, http.StatusOK, nil
}

// CloneSharedDeck copies a shared deck and its cards to the caller with
// fresh SRS state; the clone starts unshared.
func (s *service) CloneSharedDeck(ctx context.Context, userID uuid.UUID, code string) (*DeckResponse, int, error) {
	source, err := s.repo.GetDeckByShareCode(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, http.StatusNotFound, ErrSharedDeckNotFound
		}
		return nil, http.StatusInternalServerError, err
	}

	cards, err := s.repo.GetCardsByDeck(ctx, source.ID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	clone := &models.Deck{
		UserID:     userID,
		Name:       source.Name,
		SourceLang: source.SourceLang,
		TargetLang: source.TargetLang,
	}
	if err := s.repo.CreateDeck(ctx, clone); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	for _, card := range cards {
		copyCard := &models.Card{
			DeckID:       clone.ID,
			Front:        card.Front,
			Back:         card.Back,
			CardType:     card.CardType,
			EaseFactor:   2.5,
			IntervalDays: 0,
			Repetitions:  0,
			DueAt:        time.Now(),
		}
		if err := s.repo.CreateCard(ctx, copyCard); err != nil {
			return nil, http.StatusInternalServerError, err
		}
	}

	return deckResponse(clone, int64(len(cards)), int64(len(cards))), http.StatusCreated, nil
}
