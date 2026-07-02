package models

import (
	"time"

	"github.com/google/uuid"
)

// Review records one study review of a card (feeds progress stats)
type Review struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CardID     uuid.UUID `gorm:"type:uuid;not null;index"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
	Rating     int       `gorm:"not null"`
	ReviewedAt time.Time `gorm:"not null;default:now()"`
}

func (Review) TableName() string {
	return "reviews"
}
