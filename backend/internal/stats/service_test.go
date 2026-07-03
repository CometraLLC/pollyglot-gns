package stats

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeRepo struct {
	dayCounts     map[string]int64
	totalReviews  int64
	distinctCards int64
	forceErr      error

	gotSince time.Time
}

func (f *fakeRepo) GetDayCounts(_ context.Context, _ uuid.UUID, since time.Time) (map[string]int64, error) {
	if f.forceErr != nil {
		return nil, f.forceErr
	}
	f.gotSince = since
	return f.dayCounts, nil
}

func (f *fakeRepo) CountReviews(_ context.Context, _ uuid.UUID) (int64, error) {
	if f.forceErr != nil {
		return 0, f.forceErr
	}
	return f.totalReviews, nil
}

func (f *fakeRepo) CountDistinctCards(_ context.Context, _ uuid.UUID) (int64, error) {
	if f.forceErr != nil {
		return 0, f.forceErr
	}
	return f.distinctCards, nil
}

func TestGetStats(t *testing.T) {
	userID := uuid.New()
	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	t.Run("aggregates counts, streak, and a 30-day chart series", func(t *testing.T) {
		repo := &fakeRepo{
			dayCounts:     map[string]int64{today: 5, yesterday: 3},
			totalReviews:  42,
			distinctCards: 17,
		}
		svc := NewService(repo)

		resp, status, err := svc.GetStats(context.Background(), userID)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		assert.EqualValues(t, 42, resp.TotalReviews)
		assert.EqualValues(t, 17, resp.UniqueCards)
		assert.EqualValues(t, 5, resp.ReviewsToday)
		assert.Equal(t, 2, resp.StreakDays)

		require.Len(t, resp.ReviewsPerDay, 30, "chart series is exactly 30 days")
		assert.Equal(t, today, resp.ReviewsPerDay[29].Date, "series ends today")
		assert.EqualValues(t, 5, resp.ReviewsPerDay[29].Count)
		assert.EqualValues(t, 3, resp.ReviewsPerDay[28].Count)
		assert.Zero(t, resp.ReviewsPerDay[0].Count, "old days are zero-filled")
	})

	t.Run("brand-new user gets zeros, not errors", func(t *testing.T) {
		svc := NewService(&fakeRepo{})

		resp, status, err := svc.GetStats(context.Background(), userID)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		assert.Zero(t, resp.TotalReviews)
		assert.Zero(t, resp.StreakDays)
		assert.Zero(t, resp.ReviewsToday)
		assert.Len(t, resp.ReviewsPerDay, 30)
	})

	t.Run("streak window looks a year back, not just 30 days", func(t *testing.T) {
		repo := &fakeRepo{dayCounts: map[string]int64{}}
		svc := NewService(repo)

		_, _, err := svc.GetStats(context.Background(), userID)

		require.NoError(t, err)
		assert.True(t, repo.gotSince.Before(time.Now().AddDate(0, 0, -300)),
			"day counts must cover enough history for long streaks")
	})

	t.Run("maps repo failure to 500", func(t *testing.T) {
		svc := NewService(&fakeRepo{forceErr: errors.New("db down")})

		_, status, err := svc.GetStats(context.Background(), userID)

		require.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, status)
	})
}
