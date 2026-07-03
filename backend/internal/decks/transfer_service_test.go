package decks

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExportDeck(t *testing.T) {
	userID := uuid.New()

	t.Run("exports the deck's cards with a filename", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		seedCard(repo, deck.ID)
		svc := NewService(repo)

		filename, content, status, err := svc.ExportDeck(context.Background(), userID, deck.ID, "csv")

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, "japanese-basics.csv", filename)
		assert.Contains(t, content, "front,back,card_type")
		assert.Contains(t, content, "こんにちは,hello,basic")
	})

	t.Run("404 on another user's deck", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, uuid.New())
		svc := NewService(repo)

		_, _, status, err := svc.ExportDeck(context.Background(), userID, deck.ID, "csv")

		require.Error(t, err)
		assert.Equal(t, http.StatusNotFound, status)
	})

	t.Run("bad format is a 400", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		svc := NewService(repo)

		_, _, status, err := svc.ExportDeck(context.Background(), userID, deck.ID, "xlsx")

		require.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, status)
	})
}

func TestImportDeck(t *testing.T) {
	userID := uuid.New()

	t.Run("imports good rows and reports the bad ones", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		svc := NewService(repo)
		file := strings.NewReader("front,back\nねこ,cat\n,missing\nいぬ,dog\n")

		result, status, err := svc.ImportDeck(context.Background(), userID, deck.ID, file, "csv")

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, 2, result.Imported)
		require.Len(t, result.Skipped, 1)
		assert.Equal(t, 3, result.Skipped[0].Line)

		assert.Len(t, repo.cards, 2)
		for _, c := range repo.cards {
			assert.Equal(t, deck.ID, c.DeckID)
			assert.InDelta(t, 2.5, c.EaseFactor, 0.0001, "imported cards start as new")
			assert.Zero(t, c.Repetitions)
		}
	})

	t.Run("404 on another user's deck, nothing imported", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, uuid.New())
		svc := NewService(repo)

		_, status, err := svc.ImportDeck(context.Background(), userID, deck.ID, strings.NewReader("a,b\n"), "csv")

		require.Error(t, err)
		assert.Equal(t, http.StatusNotFound, status)
		assert.Empty(t, repo.cards)
	})

	t.Run("oversized files are a 400", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		svc := NewService(repo)
		var sb strings.Builder
		for range 1001 {
			sb.WriteString("a,b\n")
		}

		_, status, err := svc.ImportDeck(context.Background(), userID, deck.ID, strings.NewReader(sb.String()), "csv")

		require.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, status)
		assert.Empty(t, repo.cards)
	})
}
