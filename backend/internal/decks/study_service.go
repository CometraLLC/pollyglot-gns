package decks

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/base-go/backend/internal/shared/models"
	"github.com/base-go/backend/pkg/srs"
	"github.com/base-go/backend/pkg/validator"
)

const (
	defaultQueueLimit = 20
	maxQueueLimit     = 100
)

// ReviewCard applies an SM-2 review to an owned card, persists the new
// scheduling state, and records the review for progress stats.
func (s *service) ReviewCard(ctx context.Context, userID, cardID uuid.UUID, req ReviewCardRequest) (*CardResponse, int, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	card, status, err := s.getOwnedCard(ctx, userID, cardID)
	if err != nil {
		return nil, status, err
	}

	now := time.Now()
	state, due := srs.Review(srs.State{
		EaseFactor:   card.EaseFactor,
		IntervalDays: card.IntervalDays,
		Repetitions:  card.Repetitions,
	}, srs.Rating(*req.Rating), now)

	card.EaseFactor = state.EaseFactor
	card.IntervalDays = state.IntervalDays
	card.Repetitions = state.Repetitions
	card.DueAt = due
	card.UpdatedAt = now
	if err := s.repo.UpdateCard(ctx, card); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	review := &models.Review{
		CardID:     card.ID,
		UserID:     userID,
		Rating:     *req.Rating,
		ReviewedAt: now,
	}
	if err := s.repo.CreateReview(ctx, review); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return cardResponse(card), http.StatusOK, nil
}

// GetStudyQueue returns the deck's due cards, most overdue first.
func (s *service) GetStudyQueue(ctx context.Context, userID, deckID uuid.UUID, limit int) ([]CardResponse, int, error) {
	if _, status, err := s.getOwnedDeck(ctx, userID, deckID); err != nil {
		return nil, status, err
	}

	if limit <= 0 {
		limit = defaultQueueLimit
	}
	if limit > maxQueueLimit {
		limit = maxQueueLimit
	}

	cards, err := s.repo.GetDueCards(ctx, deckID, time.Now(), limit)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	result := make([]CardResponse, 0, len(cards))
	for i := range cards {
		result = append(result, *cardResponse(&cards[i]))
	}
	return result, http.StatusOK, nil
}
