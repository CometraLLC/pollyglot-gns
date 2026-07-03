package stats

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/base-go/backend/pkg/database"
)

type Repository interface {
	// GetDayCounts returns reviews-per-day since the given time,
	// keyed by YYYY-MM-DD.
	GetDayCounts(ctx context.Context, userID uuid.UUID, since time.Time) (map[string]int64, error)
	CountReviews(ctx context.Context, userID uuid.UUID) (int64, error)
	CountDistinctCards(ctx context.Context, userID uuid.UUID) (int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db database.Database) Repository {
	return &repository{db: db.GetDB()}
}

func (r *repository) GetDayCounts(ctx context.Context, userID uuid.UUID, since time.Time) (map[string]int64, error) {
	var rows []struct {
		Day   string
		Count int64
	}
	err := r.db.WithContext(ctx).
		Table("reviews").
		Select("to_char(reviewed_at, 'YYYY-MM-DD') AS day, COUNT(*) AS count").
		Where("user_id = ? AND reviewed_at >= ?", userID, since).
		Group("day").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64, len(rows))
	for _, row := range rows {
		counts[row.Day] = row.Count
	}
	return counts, nil
}

func (r *repository) CountReviews(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("reviews").
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (r *repository) CountDistinctCards(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("reviews").
		Where("user_id = ?", userID).
		Distinct("card_id").
		Count(&count).Error
	return count, err
}
