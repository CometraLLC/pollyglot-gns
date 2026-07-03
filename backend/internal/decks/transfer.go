package decks

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/base-go/backend/internal/shared/models"
	"github.com/base-go/backend/pkg/cloze"
)

// maxImportRows caps one import request; larger files should be split.
const maxImportRows = 1000

// CardImport is one validated row from an import file.
type CardImport struct {
	Front    string
	Back     string
	CardType string
}

// RowError reports one rejected row (1-indexed line numbers, header
// included in the count so numbers match what an editor shows).
type RowError struct {
	Line  int    `json:"line"`
	Error string `json:"error"`
}

func separator(format string) (rune, error) {
	switch format {
	case "csv":
		return ',', nil
	case "tsv":
		return '\t', nil
	default:
		return 0, fmt.Errorf("unsupported format %q (use csv or tsv)", format)
	}
}

// ExportCards serializes cards as CSV/TSV with a header row. TSV output
// imports cleanly into Anki.
func ExportCards(cards []models.Card, format string) (string, error) {
	sep, err := separator(format)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	w := csv.NewWriter(&sb)
	w.Comma = sep

	if err := w.Write([]string{"front", "back", "card_type"}); err != nil {
		return "", err
	}
	for _, card := range cards {
		if err := w.Write([]string{card.Front, card.Back, card.CardType}); err != nil {
			return "", err
		}
	}
	w.Flush()
	return sb.String(), w.Error()
}

// ParseImport reads CSV/TSV rows into card imports. A header row is
// detected and skipped; bad rows land in rowErrors while good rows
// survive; more than maxImportRows rows is an error for the whole file.
func ParseImport(r io.Reader, format string) ([]CardImport, []RowError, error) {
	sep, err := separator(format)
	if err != nil {
		return nil, nil, err
	}

	reader := csv.NewReader(r)
	reader.Comma = sep
	reader.FieldsPerRecord = -1 // validated per row below

	var rows []CardImport
	var rowErrors []RowError
	line := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		line++
		if err != nil {
			rowErrors = append(rowErrors, RowError{Line: line, Error: "unparseable row"})
			continue
		}
		if line > maxImportRows {
			return nil, nil, fmt.Errorf("too many rows: imports are capped at %d", maxImportRows)
		}

		// header detection: exact column names in the first row
		if line == 1 && len(record) >= 2 &&
			strings.EqualFold(strings.TrimSpace(record[0]), "front") &&
			strings.EqualFold(strings.TrimSpace(record[1]), "back") {
			continue
		}

		row, rowErr := parseRow(record)
		if rowErr != "" {
			rowErrors = append(rowErrors, RowError{Line: line, Error: rowErr})
			continue
		}
		rows = append(rows, row)
	}

	return rows, rowErrors, nil
}

func parseRow(record []string) (CardImport, string) {
	if len(record) < 2 {
		return CardImport{}, "expected at least front and back columns"
	}

	front := strings.TrimSpace(record[0])
	back := strings.TrimSpace(record[1])
	cardType := models.CardTypeBasic
	if len(record) >= 3 && strings.TrimSpace(record[2]) != "" {
		cardType = strings.TrimSpace(record[2])
	}

	if front == "" {
		return CardImport{}, "front is required"
	}
	if back == "" {
		return CardImport{}, "back is required"
	}
	if cardType != models.CardTypeBasic && cardType != models.CardTypeCloze {
		return CardImport{}, fmt.Sprintf("card_type %q is not basic or cloze", cardType)
	}
	if cardType == models.CardTypeCloze && len(cloze.Deletions(front)) == 0 {
		return CardImport{}, "cloze rows need at least one {{c1::...}} deletion"
	}

	return CardImport{Front: front, Back: back, CardType: cardType}, ""
}
