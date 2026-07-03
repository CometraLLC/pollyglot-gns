package decks

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/base-go/backend/internal/shared/factory"
	"github.com/base-go/backend/internal/shared/models"
)

func TestExportCards(t *testing.T) {
	cards := []models.Card{
		factory.Card().WithFront("こんにちは").WithBack("hello").Build(),
		factory.Card().WithFront("with, comma").WithBack(`with "quotes"`).Build(),
		factory.Card().Cloze("水を{{c1::飲みます}}").WithBack("drink water").Build(),
	}

	t.Run("csv with header, quoting, and card types", func(t *testing.T) {
		out, err := ExportCards(cards, "csv")

		require.NoError(t, err)
		lines := strings.Split(strings.TrimSpace(out), "\n")
		require.Len(t, lines, 4)
		assert.Equal(t, "front,back,card_type", lines[0])
		assert.Equal(t, "こんにちは,hello,basic", lines[1])
		assert.Equal(t, `"with, comma","with ""quotes""",basic`, lines[2])
		assert.Equal(t, "水を{{c1::飲みます}},drink water,cloze", lines[3])
	})

	t.Run("tsv is Anki-import compatible", func(t *testing.T) {
		out, err := ExportCards(cards, "tsv")

		require.NoError(t, err)
		lines := strings.Split(strings.TrimSpace(out), "\n")
		assert.Equal(t, "front\tback\tcard_type", lines[0])
		assert.Equal(t, "こんにちは\thello\tbasic", lines[1])
	})

	t.Run("unknown format rejected", func(t *testing.T) {
		_, err := ExportCards(cards, "xlsx")
		require.Error(t, err)
	})
}

func TestParseImport(t *testing.T) {
	t.Run("parses rows and skips the header", func(t *testing.T) {
		input := "front,back,card_type\nこんにちは,hello,basic\n水を{{c1::飲みます}},drink water,cloze\n"

		rows, rowErrors, err := ParseImport(strings.NewReader(input), "csv")

		require.NoError(t, err)
		assert.Empty(t, rowErrors)
		require.Len(t, rows, 2)
		assert.Equal(t, "こんにちは", rows[0].Front)
		assert.Equal(t, "basic", rows[0].CardType)
		assert.Equal(t, "cloze", rows[1].CardType)
	})

	t.Run("card_type column is optional and defaults to basic", func(t *testing.T) {
		input := "ねこ,cat\nいぬ,dog\n"

		rows, rowErrors, err := ParseImport(strings.NewReader(input), "csv")

		require.NoError(t, err)
		assert.Empty(t, rowErrors)
		require.Len(t, rows, 2)
		assert.Equal(t, "basic", rows[0].CardType)
	})

	t.Run("bad rows are reported with line numbers, good rows survive", func(t *testing.T) {
		input := strings.Join([]string{
			"front,back,card_type",
			",missing front,basic",
			"missing back,,basic",
			"ok,fine,basic",
			"not a cloze,really,cloze",
			"weird,type,audio",
		}, "\n")

		rows, rowErrors, err := ParseImport(strings.NewReader(input), "csv")

		require.NoError(t, err)
		require.Len(t, rows, 1)
		assert.Equal(t, "ok", rows[0].Front)
		require.Len(t, rowErrors, 4)
		assert.Equal(t, 2, rowErrors[0].Line)
		assert.Contains(t, rowErrors[0].Error, "front")
		assert.Equal(t, 5, rowErrors[2].Line)
		assert.Contains(t, rowErrors[2].Error, "deletion")
		assert.Equal(t, 6, rowErrors[3].Line)
		assert.Contains(t, rowErrors[3].Error, "card_type")
	})

	t.Run("blank lines and CRLF are tolerated", func(t *testing.T) {
		input := "front,back\r\nねこ,cat\r\n\r\nいぬ,dog\r\n"

		rows, rowErrors, err := ParseImport(strings.NewReader(input), "csv")

		require.NoError(t, err)
		assert.Empty(t, rowErrors)
		assert.Len(t, rows, 2)
	})

	t.Run("tsv parses tab-separated values", func(t *testing.T) {
		input := "こんにちは\thello\nwith, comma\tstill one field\n"

		rows, rowErrors, err := ParseImport(strings.NewReader(input), "tsv")

		require.NoError(t, err)
		assert.Empty(t, rowErrors)
		require.Len(t, rows, 2)
		assert.Equal(t, "with, comma", rows[1].Front)
	})

	t.Run("row cap is enforced", func(t *testing.T) {
		var sb strings.Builder
		for range 1001 {
			sb.WriteString("a,b\n")
		}

		_, _, err := ParseImport(strings.NewReader(sb.String()), "csv")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "1000")
	})

	t.Run("round-trips its own export", func(t *testing.T) {
		cards := []models.Card{
			factory.Card().WithFront("with, comma").WithBack(`with "quotes"`).Build(),
			factory.Card().Cloze("{{c1::猫}}が好き").WithBack("I like cats").Build(),
		}
		out, err := ExportCards(cards, "csv")
		require.NoError(t, err)

		rows, rowErrors, err := ParseImport(strings.NewReader(out), "csv")

		require.NoError(t, err)
		assert.Empty(t, rowErrors)
		require.Len(t, rows, 2)
		assert.Equal(t, cards[0].Front, rows[0].Front)
		assert.Equal(t, cards[0].Back, rows[0].Back)
		assert.Equal(t, cards[1].Front, rows[1].Front)
		assert.Equal(t, "cloze", rows[1].CardType)
	})
}
