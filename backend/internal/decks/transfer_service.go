package decks

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/base-go/backend/internal/shared/models"
)

// ImportResult summarizes an import: how many cards landed and which
// rows were rejected (with line numbers matching the uploaded file).
type ImportResult struct {
	Imported int        `json:"imported"`
	Skipped  []RowError `json:"skipped"`
}

var filenameUnsafe = regexp.MustCompile(`[^a-z0-9]+`)

// ExportDeck renders an owned deck's cards as CSV/TSV plus a
// download-friendly filename derived from the deck name.
func (s *service) ExportDeck(ctx context.Context, userID, deckID uuid.UUID, format string) (string, string, int, error) {
	deck, status, err := s.getOwnedDeck(ctx, userID, deckID)
	if err != nil {
		return "", "", status, err
	}

	cards, err := s.repo.GetCardsByDeck(ctx, deckID)
	if err != nil {
		return "", "", http.StatusInternalServerError, err
	}

	content, err := ExportCards(cards, format)
	if err != nil {
		return "", "", http.StatusBadRequest, err
	}

	slug := strings.Trim(filenameUnsafe.ReplaceAllString(strings.ToLower(deck.Name), "-"), "-")
	if slug == "" {
		slug = "deck"
	}
	return fmt.Sprintf("%s.%s", slug, format), content, http.StatusOK, nil
}

// ImportDeck parses an upload and creates the valid rows as new cards.
func (s *service) ImportDeck(ctx context.Context, userID, deckID uuid.UUID, file io.Reader, format string) (*ImportResult, int, error) {
	if _, status, err := s.getOwnedDeck(ctx, userID, deckID); err != nil {
		return nil, status, err
	}

	rows, rowErrors, err := ParseImport(file, format)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	for _, row := range rows {
		card := &models.Card{
			DeckID:       deckID,
			Front:        row.Front,
			Back:         row.Back,
			CardType:     row.CardType,
			EaseFactor:   2.5,
			IntervalDays: 0,
			Repetitions:  0,
			DueAt:        time.Now(),
		}
		if err := s.repo.CreateCard(ctx, card); err != nil {
			return nil, http.StatusInternalServerError, err
		}
	}

	if rowErrors == nil {
		rowErrors = []RowError{}
	}
	return &ImportResult{Imported: len(rows), Skipped: rowErrors}, http.StatusOK, nil
}
