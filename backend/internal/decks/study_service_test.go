package decks

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func intPtr(v int) *int { return &v }

func TestReviewCard(t *testing.T) {
	userID := uuid.New()

	t.Run("applies SM-2 and records the review", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		card := seedCard(repo, deck.ID) // new card: ease 2.5, reps 0
		svc := NewService(repo)

		before := time.Now()
		resp, status, err := svc.ReviewCard(context.Background(), userID, card.ID, ReviewCardRequest{
			Rating: intPtr(4), // Got it!
		})

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, 1, resp.Repetitions)
		assert.Equal(t, 1, resp.IntervalDays)
		assert.InDelta(t, 2.6, resp.EaseFactor, 0.0001)
		assert.True(t, resp.DueAt.After(before.Add(23*time.Hour)), "next due ~1 day out")

		// persisted card state matches the response
		assert.Equal(t, 1, repo.cards[card.ID].Repetitions)

		// a review row was recorded for stats
		require.Len(t, repo.reviews, 1)
		assert.Equal(t, card.ID, repo.reviews[0].CardID)
		assert.Equal(t, userID, repo.reviews[0].UserID)
		assert.Equal(t, 4, repo.reviews[0].Rating)
	})

	t.Run("rating zero (Forgot) is a valid rating", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		card := seedCard(repo, deck.ID)
		svc := NewService(repo)

		resp, status, err := svc.ReviewCard(context.Background(), userID, card.ID, ReviewCardRequest{
			Rating: intPtr(0),
		})

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		assert.Zero(t, resp.Repetitions, "Forgot resets repetitions")
		require.Len(t, repo.reviews, 1)
		assert.Zero(t, repo.reviews[0].Rating)
	})

	t.Run("rejects a missing rating", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		card := seedCard(repo, deck.ID)
		svc := NewService(repo)

		_, status, err := svc.ReviewCard(context.Background(), userID, card.ID, ReviewCardRequest{})

		require.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, status)
		assert.Empty(t, repo.reviews)
	})

	t.Run("rejects out-of-range ratings", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		card := seedCard(repo, deck.ID)
		svc := NewService(repo)

		for _, rating := range []int{-1, 5, 42} {
			_, status, err := svc.ReviewCard(context.Background(), userID, card.ID, ReviewCardRequest{
				Rating: intPtr(rating),
			})
			require.Error(t, err, "rating %d", rating)
			assert.Equal(t, http.StatusBadRequest, status)
		}
		assert.Empty(t, repo.reviews)
	})

	t.Run("404 on another user's card, nothing recorded", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, uuid.New())
		card := seedCard(repo, deck.ID)
		svc := NewService(repo)

		_, status, err := svc.ReviewCard(context.Background(), userID, card.ID, ReviewCardRequest{
			Rating: intPtr(3),
		})

		require.Error(t, err)
		assert.Equal(t, http.StatusNotFound, status)
		assert.Empty(t, repo.reviews)
		assert.Zero(t, repo.cards[card.ID].Repetitions, "card state untouched")
	})
}

func TestGetStudyQueue(t *testing.T) {
	userID := uuid.New()

	t.Run("returns due cards oldest first, future cards excluded", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)

		overdue := seedCard(repo, deck.ID)
		overdue.DueAt = time.Now().Add(-48 * time.Hour)
		dueNow := seedCard(repo, deck.ID)
		dueNow.DueAt = time.Now().Add(-time.Minute)
		future := seedCard(repo, deck.ID)
		future.DueAt = time.Now().Add(72 * time.Hour)

		svc := NewService(repo)

		resp, status, err := svc.GetStudyQueue(context.Background(), userID, deck.ID, 20)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		require.Len(t, resp, 2)
		assert.Equal(t, overdue.ID, resp[0].ID, "most overdue first")
		assert.Equal(t, dueNow.ID, resp[1].ID)
		for _, c := range resp {
			assert.NotEqual(t, future.ID, c.ID)
		}
	})

	t.Run("respects the limit", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		for range 5 {
			c := seedCard(repo, deck.ID)
			c.DueAt = time.Now().Add(-time.Hour)
		}
		svc := NewService(repo)

		resp, _, err := svc.GetStudyQueue(context.Background(), userID, deck.ID, 3)

		require.NoError(t, err)
		assert.Len(t, resp, 3)
	})

	t.Run("zero or negative limit falls back to the default", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		c := seedCard(repo, deck.ID)
		c.DueAt = time.Now().Add(-time.Hour)
		svc := NewService(repo)

		resp, _, err := svc.GetStudyQueue(context.Background(), userID, deck.ID, 0)

		require.NoError(t, err)
		assert.Len(t, resp, 1)
	})

	t.Run("limit is capped at 100", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		svc := NewService(repo)

		_, _, err := svc.GetStudyQueue(context.Background(), userID, deck.ID, 5000)

		require.NoError(t, err)
		assert.LessOrEqual(t, repo.lastQueueLimit, 100, "repo must never be asked for more than the cap")
	})

	t.Run("404 on another user's deck", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, uuid.New())
		svc := NewService(repo)

		_, status, err := svc.GetStudyQueue(context.Background(), userID, deck.ID, 20)

		require.Error(t, err)
		assert.Equal(t, http.StatusNotFound, status)
	})

	t.Run("empty queue serializes as [] not null", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		svc := NewService(repo)

		resp, status, err := svc.GetStudyQueue(context.Background(), userID, deck.ID, 20)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		assert.NotNil(t, resp)
		assert.Empty(t, resp)
	})
}
