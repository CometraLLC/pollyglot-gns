// Package srs implements the SM-2 spaced-repetition scheduler
// (SuperMemo-2, the algorithm behind Anki's ancestor) as pure functions.
// See docs/DECISIONS.md D-006 for the rationale and rating mapping.
package srs

import (
	"math"
	"time"
)

// Rating is the learner's self-assessment, one per study button.
type Rating int

const (
	Forgot    Rating = 0
	Difficult Rating = 1
	Okay      Rating = 2
	Almost    Rating = 3
	GotIt     Rating = 4
)

// State is a card's scheduling state between reviews.
type State struct {
	EaseFactor   float64
	IntervalDays int
	Repetitions  int
}

const (
	easeFloor   = 1.3
	defaultEase = 2.5
)

// Review applies one review to a card's state and returns the new state
// plus the next due time. The 0–4 rating maps to SM-2 quality 1–5;
// quality < 3 (Forgot, Difficult) is a lapse: the repetition sequence
// restarts at a one-day interval and, per classic SM-2, ease is unchanged.
// Out-of-range ratings clamp to the nearest valid one, and corrupt ease
// values are normalized so intervals can never shrink below the floor.
func Review(state State, rating Rating, now time.Time) (State, time.Time) {
	if rating < Forgot {
		rating = Forgot
	}
	if rating > GotIt {
		rating = GotIt
	}
	if state.EaseFactor < easeFloor {
		if state.EaseFactor == 0 {
			state.EaseFactor = defaultEase
		} else {
			state.EaseFactor = easeFloor
		}
	}

	quality := float64(rating) + 1

	if quality < 3 {
		state.Repetitions = 0
		state.IntervalDays = 1
	} else {
		// Canonical SM-2: the new interval uses the ease factor as it
		// stood before this review; the ease update applies afterwards.
		state.Repetitions++
		switch state.Repetitions {
		case 1:
			state.IntervalDays = 1
		case 2:
			state.IntervalDays = 6
		default:
			state.IntervalDays = int(math.Round(float64(state.IntervalDays) * state.EaseFactor))
		}

		state.EaseFactor += 0.1 - (5-quality)*(0.08+(5-quality)*0.02)
		if state.EaseFactor < easeFloor {
			state.EaseFactor = easeFloor
		}
	}

	if state.IntervalDays < 1 {
		state.IntervalDays = 1
	}

	return state, now.AddDate(0, 0, state.IntervalDays)
}
