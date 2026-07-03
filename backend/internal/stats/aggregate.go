package stats

import "time"

const dayFormat = "2006-01-02"

// DayCount is one day's review count, Date formatted YYYY-MM-DD.
type DayCount struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

// Streak counts consecutive reviewed days ending today — or ending
// yesterday when today has no review yet, so a live streak isn't shown
// as broken before the learner studies. Future dates are ignored.
func Streak(days map[string]bool, today time.Time) int {
	start := today
	if !days[today.Format(dayFormat)] {
		start = today.AddDate(0, 0, -1)
		if !days[start.Format(dayFormat)] {
			return 0
		}
	}

	streak := 0
	for d := start; days[d.Format(dayFormat)]; d = d.AddDate(0, 0, -1) {
		streak++
	}
	return streak
}

// FillDays expands sparse day counts into exactly n consecutive entries
// ending today, zero-filled, oldest first — chart-ready.
func FillDays(counts map[string]int64, today time.Time, n int) []DayCount {
	filled := make([]DayCount, 0, n)
	for i := n - 1; i >= 0; i-- {
		date := today.AddDate(0, 0, -i).Format(dayFormat)
		filled = append(filled, DayCount{Date: date, Count: counts[date]})
	}
	return filled
}
