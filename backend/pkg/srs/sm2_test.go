package srs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var now = time.Date(2026, 7, 2, 12, 0, 0, 0, time.UTC)

func newCard() State {
	return State{EaseFactor: 2.5, IntervalDays: 0, Repetitions: 0}
}

func TestReviewRatingsFromNewCard(t *testing.T) {
	tests := []struct {
		name         string
		rating       Rating
		wantReps     int
		wantInterval int
		wantEase     float64
	}{
		// Forgot (q=1) and Difficult (q=2) are lapses: repetitions reset,
		// card comes back tomorrow, ease untouched (classic SM-2).
		{"Forgot is a lapse", Forgot, 0, 1, 2.5},
		{"Difficult is a lapse", Difficult, 0, 1, 2.5},
		// Okay (q=3): first successful repetition, ease drops by 0.14.
		{"Okay passes with ease penalty", Okay, 1, 1, 2.36},
		// Almost (q=4): ease unchanged.
		{"Almost passes with ease unchanged", Almost, 1, 1, 2.5},
		// Got it! (q=5): ease grows by 0.1.
		{"Got it passes with ease bonus", GotIt, 1, 1, 2.6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state, due := Review(newCard(), tt.rating, now)

			assert.Equal(t, tt.wantReps, state.Repetitions)
			assert.Equal(t, tt.wantInterval, state.IntervalDays)
			assert.InDelta(t, tt.wantEase, state.EaseFactor, 0.0001)
			assert.Equal(t, now.AddDate(0, 0, tt.wantInterval), due)
		})
	}
}

func TestReviewProgression(t *testing.T) {
	// Three perfect reviews: 1 day, 6 days, then round(6 * ease).
	state := newCard()

	state, due := Review(state, GotIt, now)
	assert.Equal(t, 1, state.IntervalDays)
	assert.Equal(t, 1, state.Repetitions)
	assert.InDelta(t, 2.6, state.EaseFactor, 0.0001)
	assert.Equal(t, now.AddDate(0, 0, 1), due)

	state, due = Review(state, GotIt, now)
	assert.Equal(t, 6, state.IntervalDays)
	assert.Equal(t, 2, state.Repetitions)
	assert.InDelta(t, 2.7, state.EaseFactor, 0.0001)
	assert.Equal(t, now.AddDate(0, 0, 6), due)

	state, due = Review(state, GotIt, now)
	// round(6 * 2.7) = 16
	assert.Equal(t, 16, state.IntervalDays)
	assert.Equal(t, 3, state.Repetitions)
	assert.InDelta(t, 2.8, state.EaseFactor, 0.0001)
	assert.Equal(t, now.AddDate(0, 0, 16), due)
}

func TestReviewLapseResetsProgressButKeepsEase(t *testing.T) {
	// Build up a mature card, then forget it.
	state := State{EaseFactor: 2.7, IntervalDays: 16, Repetitions: 3}

	state, due := Review(state, Forgot, now)

	assert.Equal(t, 0, state.Repetitions, "lapse restarts the repetition sequence")
	assert.Equal(t, 1, state.IntervalDays, "lapsed card comes back tomorrow")
	assert.InDelta(t, 2.7, state.EaseFactor, 0.0001, "SM-2 leaves ease unchanged on lapse")
	assert.Equal(t, now.AddDate(0, 0, 1), due)

	// The repetition after a lapse starts over at interval 1 → 6 → ...
	state, _ = Review(state, GotIt, now)
	assert.Equal(t, 1, state.IntervalDays)
	assert.Equal(t, 1, state.Repetitions)
}

func TestReviewEaseNeverDropsBelowFloor(t *testing.T) {
	state := newCard()

	// Repeated barely-passing reviews grind ease down; it must stop at 1.3.
	for range 20 {
		state, _ = Review(state, Okay, now)
	}

	assert.InDelta(t, 1.3, state.EaseFactor, 0.0001)
	assert.GreaterOrEqual(t, state.EaseFactor, 1.3)
}

func TestReviewMatureCardIntervalUsesPreReviewEase(t *testing.T) {
	// Canonical SM-2: interval multiplies by the ease factor as it stood
	// before this review; the ease penalty applies afterwards.
	state := State{EaseFactor: 2.0, IntervalDays: 10, Repetitions: 5}

	state, _ = Review(state, Okay, now)

	// interval = round(10 * 2.0) = 20, then ease drops to 1.86
	assert.Equal(t, 20, state.IntervalDays)
	assert.InDelta(t, 1.86, state.EaseFactor, 0.0001)
}

func TestReviewClampsOutOfRangeRatings(t *testing.T) {
	below, _ := Review(newCard(), Rating(-3), now)
	forgot, _ := Review(newCard(), Forgot, now)
	assert.Equal(t, forgot, below, "ratings below range behave like Forgot")

	above, _ := Review(newCard(), Rating(99), now)
	gotIt, _ := Review(newCard(), GotIt, now)
	assert.Equal(t, gotIt, above, "ratings above range behave like Got it!")
}

func TestReviewNormalizesCorruptState(t *testing.T) {
	// Zero-value or corrupt ease (e.g. bad import) must not collapse
	// scheduling to sub-1.3 multipliers.
	state := State{EaseFactor: 0, IntervalDays: 4, Repetitions: 3}

	state, _ = Review(state, GotIt, now)

	assert.GreaterOrEqual(t, state.EaseFactor, 1.3)
	assert.GreaterOrEqual(t, state.IntervalDays, 1)
}
