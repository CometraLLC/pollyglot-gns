package cloze

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeletions(t *testing.T) {
	tests := []struct {
		name  string
		front string
		want  []string
	}{
		{"single deletion", "水を{{c1::飲みます}}", []string{"飲みます"}},
		{"multiple deletions", "{{c1::猫}}が{{c2::魚}}を食べる", []string{"猫", "魚"}},
		{"no deletions", "plain text", nil},
		{"empty string", "", nil},
		{"unclosed marker is not a deletion", "水を{{c1::飲みます", nil},
		{"empty deletion is ignored", "水を{{c1::}}", nil},
		{"marker without cN prefix is not a deletion", "{{飲みます}}", nil},
		{"latin text", "the {{c1::cat}} sat on the {{c2::mat}}", []string{"cat", "mat"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Deletions(tt.front))
		})
	}
}

func TestBlank(t *testing.T) {
	tests := []struct {
		name  string
		front string
		want  string
	}{
		{"single deletion becomes a blank", "水を{{c1::飲みます}}", "水を[…]"},
		{"multiple blanks", "{{c1::猫}}が{{c2::魚}}を食べる", "[…]が[…]を食べる"},
		{"plain text unchanged", "plain text", "plain text"},
		{"unclosed marker left as-is", "水を{{c1::飲みます", "水を{{c1::飲みます"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Blank(tt.front))
		})
	}
}

func TestReveal(t *testing.T) {
	tests := []struct {
		name  string
		front string
		want  string
	}{
		{"markers stripped", "水を{{c1::飲みます}}", "水を飲みます"},
		{"multiple markers stripped", "{{c1::猫}}が{{c2::魚}}を食べる", "猫が魚を食べる"},
		{"plain text unchanged", "plain text", "plain text"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Reveal(tt.front))
		})
	}
}
