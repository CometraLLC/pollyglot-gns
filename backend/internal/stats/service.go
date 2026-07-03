package stats

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type StatsResponse struct {
	ReviewsToday  int64      `json:"reviews_today"`
	TotalReviews  int64      `json:"total_reviews"`
	UniqueCards   int64      `json:"unique_cards"`
	StreakDays    int        `json:"streak_days"`
	ReviewsPerDay []DayCount `json:"reviews_per_day"`
}

type Service interface {
	GetStats(ctx context.Context, userID uuid.UUID) (*StatsResponse, int, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

const chartDays = 30

func (s *service) GetStats(ctx context.Context, userID uuid.UUID) (*StatsResponse, int, error) {
	now := time.Now()

	// A year of day counts covers streaks; the chart uses the last 30.
	dayCounts, err := s.repo.GetDayCounts(ctx, userID, now.AddDate(-1, 0, 0))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	totalReviews, err := s.repo.CountReviews(ctx, userID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	uniqueCards, err := s.repo.CountDistinctCards(ctx, userID)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	days := make(map[string]bool, len(dayCounts))
	for date, count := range dayCounts {
		if count > 0 {
			days[date] = true
		}
	}

	return &StatsResponse{
		ReviewsToday:  dayCounts[now.Format(dayFormat)],
		TotalReviews:  totalReviews,
		UniqueCards:   uniqueCards,
		StreakDays:    Streak(days, now),
		ReviewsPerDay: FillDays(dayCounts, now, chartDays),
	}, http.StatusOK, nil
}
