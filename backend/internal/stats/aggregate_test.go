package stats

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func day(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestStreak(t *testing.T) {
	today := day("2026-07-02")

	tests := []struct {
		name string
		days []string
		want int
	}{
		{"no reviews ever", nil, 0},
		{"only today", []string{"2026-07-02"}, 1},
		{"today and yesterday", []string{"2026-07-01", "2026-07-02"}, 2},
		{
			"today not yet reviewed keeps yesterday's streak alive",
			[]string{"2026-06-30", "2026-07-01"},
			2,
		},
		{"only yesterday", []string{"2026-07-01"}, 1},
		{"a gap breaks the streak", []string{"2026-06-28", "2026-06-29", "2026-07-01", "2026-07-02"}, 2},
		{"last review two days ago is a dead streak", []string{"2026-06-30"}, 0},
		{
			"long unbroken run counts fully",
			[]string{"2026-06-28", "2026-06-29", "2026-06-30", "2026-07-01", "2026-07-02"},
			5,
		},
		{"duplicate days count once", []string{"2026-07-02", "2026-07-02"}, 1},
		{"future noise is ignored", []string{"2026-07-09", "2026-07-02"}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			days := make(map[string]bool, len(tt.days))
			for _, d := range tt.days {
				days[d] = true
			}
			assert.Equal(t, tt.want, Streak(days, today))
		})
	}
}

func TestFillDays(t *testing.T) {
	today := day("2026-07-02")

	t.Run("fills missing days with zero, oldest first", func(t *testing.T) {
		counts := map[string]int64{
			"2026-07-02": 3,
			"2026-06-30": 1,
		}

		filled := FillDays(counts, today, 4)

		assert.Equal(t, []DayCount{
			{Date: "2026-06-29", Count: 0},
			{Date: "2026-06-30", Count: 1},
			{Date: "2026-07-01", Count: 0},
			{Date: "2026-07-02", Count: 3},
		}, filled)
	})

	t.Run("empty input still yields n zeroed days", func(t *testing.T) {
		filled := FillDays(nil, today, 3)

		assert.Len(t, filled, 3)
		for _, dc := range filled {
			assert.Zero(t, dc.Count)
		}
		assert.Equal(t, "2026-07-02", filled[2].Date)
	})
}
