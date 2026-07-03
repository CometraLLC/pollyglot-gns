package decks

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateShareCode(t *testing.T) {
	seen := map[string]bool{}
	for range 50 {
		code, err := GenerateShareCode()
		require.NoError(t, err)
		assert.Len(t, code, 10)
		for _, r := range code {
			assert.Contains(t, shareAlphabet, string(r), "codes use only unambiguous characters")
		}
		assert.False(t, seen[code], "codes must not repeat")
		seen[code] = true
	}
}

func TestShareDeck(t *testing.T) {
	userID := uuid.New()

	t.Run("assigns a code and is idempotent", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		svc := NewService(repo)

		first, status, err := svc.ShareDeck(context.Background(), userID, deck.ID)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		assert.Len(t, first.ShareCode, 10)

		second, _, err := svc.ShareDeck(context.Background(), userID, deck.ID)
		require.NoError(t, err)
		assert.Equal(t, first.ShareCode, second.ShareCode, "sharing twice keeps the same code")
	})

	t.Run("404 on another user's deck", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, uuid.New())
		svc := NewService(repo)

		_, status, err := svc.ShareDeck(context.Background(), userID, deck.ID)

		require.Error(t, err)
		assert.Equal(t, http.StatusNotFound, status)
	})

	t.Run("unshare clears the code", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		svc := NewService(repo)
		_, _, err := svc.ShareDeck(context.Background(), userID, deck.ID)
		require.NoError(t, err)

		status, err := svc.UnshareDeck(context.Background(), userID, deck.ID)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		assert.Nil(t, repo.decks[deck.ID].ShareCode)
	})
}

func TestGetSharedDeck(t *testing.T) {
	userID := uuid.New()

	t.Run("previews a shared deck with sample cards", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		for range 7 {
			seedCard(repo, deck.ID)
		}
		svc := NewService(repo)
		share, _, err := svc.ShareDeck(context.Background(), userID, deck.ID)
		require.NoError(t, err)

		preview, status, err := svc.GetSharedDeck(context.Background(), share.ShareCode)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, deck.Name, preview.Name)
		assert.Equal(t, deck.SourceLang, preview.SourceLang)
		assert.EqualValues(t, 7, preview.CardCount)
		assert.LessOrEqual(t, len(preview.SampleCards), 5, "preview shows at most five cards")
		assert.NotEmpty(t, preview.SampleCards)
	})

	t.Run("unknown code is a 404", func(t *testing.T) {
		svc := NewService(newFakeRepo())

		_, status, err := svc.GetSharedDeck(context.Background(), "NOPE123456")

		require.Error(t, err)
		assert.Equal(t, http.StatusNotFound, status)
	})
}

func TestCloneSharedDeck(t *testing.T) {
	owner := uuid.New()
	cloner := uuid.New()

	t.Run("copies the deck and cards with fresh SRS state, no share code", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, owner)
		mature := seedCard(repo, deck.ID)
		mature.EaseFactor = 1.7
		mature.Repetitions = 9
		svc := NewService(repo)
		share, _, err := svc.ShareDeck(context.Background(), owner, deck.ID)
		require.NoError(t, err)

		clone, status, err := svc.CloneSharedDeck(context.Background(), cloner, share.ShareCode)

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, status)
		assert.Equal(t, deck.Name, clone.Name)
		assert.NotEqual(t, deck.ID, clone.ID)

		cloned := repo.decks[clone.ID]
		require.NotNil(t, cloned)
		assert.Equal(t, cloner, cloned.UserID)
		assert.Nil(t, cloned.ShareCode, "clones start unshared")

		var clonedCards int
		for _, c := range repo.cards {
			if c.DeckID == clone.ID {
				clonedCards++
				assert.InDelta(t, 2.5, c.EaseFactor, 0.0001, "clone SRS state is fresh")
				assert.Zero(t, c.Repetitions)
			}
		}
		assert.Equal(t, 1, clonedCards)
	})

	t.Run("unknown or unshared codes are 404", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, owner) // never shared
		_ = deck
		svc := NewService(repo)

		_, status, err := svc.CloneSharedDeck(context.Background(), cloner, "NOPE123456")

		require.Error(t, err)
		assert.Equal(t, http.StatusNotFound, status)
	})
}
