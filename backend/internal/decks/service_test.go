package decks

import (
	"context"
	"errors"
	"net/http"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/base-go/backend/internal/shared/factory"
	"github.com/base-go/backend/internal/shared/models"
)

// --- fake repository ---

type fakeRepo struct {
	decks   map[uuid.UUID]*models.Deck
	cards   map[uuid.UUID]*models.Card
	reviews []models.Review
	// lastQueueLimit records the limit passed to GetDueCards
	lastQueueLimit int
	// forceErr, when set, is returned by every method to drive 500 paths
	forceErr error
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		decks: make(map[uuid.UUID]*models.Deck),
		cards: make(map[uuid.UUID]*models.Card),
	}
}

func (f *fakeRepo) CreateDeck(_ context.Context, deck *models.Deck) error {
	if f.forceErr != nil {
		return f.forceErr
	}
	deck.ID = uuid.New()
	now := time.Now()
	deck.CreatedAt, deck.UpdatedAt = now, now
	f.decks[deck.ID] = deck
	return nil
}

func (f *fakeRepo) GetDecksByUser(_ context.Context, userID uuid.UUID) ([]DeckWithCount, error) {
	if f.forceErr != nil {
		return nil, f.forceErr
	}
	var result []DeckWithCount
	for _, d := range f.decks {
		if d.UserID == userID && d.DeletedAt == nil {
			var count int64
			for _, c := range f.cards {
				if c.DeckID == d.ID && c.DeletedAt == nil {
					count++
				}
			}
			result = append(result, DeckWithCount{Deck: *d, CardCount: count})
		}
	}
	return result, nil
}

func (f *fakeRepo) GetDeckByID(_ context.Context, id uuid.UUID) (*models.Deck, error) {
	if f.forceErr != nil {
		return nil, f.forceErr
	}
	d, ok := f.decks[id]
	if !ok || d.DeletedAt != nil {
		return nil, gorm.ErrRecordNotFound
	}
	return d, nil
}

func (f *fakeRepo) CountCards(_ context.Context, deckID uuid.UUID) (int64, error) {
	if f.forceErr != nil {
		return 0, f.forceErr
	}
	var count int64
	for _, c := range f.cards {
		if c.DeckID == deckID && c.DeletedAt == nil {
			count++
		}
	}
	return count, nil
}

func (f *fakeRepo) UpdateDeck(_ context.Context, deck *models.Deck) error {
	if f.forceErr != nil {
		return f.forceErr
	}
	f.decks[deck.ID] = deck
	return nil
}

func (f *fakeRepo) SoftDeleteDeck(_ context.Context, id uuid.UUID) error {
	if f.forceErr != nil {
		return f.forceErr
	}
	now := time.Now()
	f.decks[id].DeletedAt = &now
	return nil
}

func (f *fakeRepo) CreateCard(_ context.Context, card *models.Card) error {
	if f.forceErr != nil {
		return f.forceErr
	}
	card.ID = uuid.New()
	now := time.Now()
	card.CreatedAt, card.UpdatedAt = now, now
	f.cards[card.ID] = card
	return nil
}

func (f *fakeRepo) GetCardsByDeck(_ context.Context, deckID uuid.UUID) ([]models.Card, error) {
	if f.forceErr != nil {
		return nil, f.forceErr
	}
	var cards []models.Card
	for _, c := range f.cards {
		if c.DeckID == deckID && c.DeletedAt == nil {
			cards = append(cards, *c)
		}
	}
	return cards, nil
}

func (f *fakeRepo) GetCardByID(_ context.Context, id uuid.UUID) (*models.Card, error) {
	if f.forceErr != nil {
		return nil, f.forceErr
	}
	c, ok := f.cards[id]
	if !ok || c.DeletedAt != nil {
		return nil, gorm.ErrRecordNotFound
	}
	return c, nil
}

func (f *fakeRepo) UpdateCard(_ context.Context, card *models.Card) error {
	if f.forceErr != nil {
		return f.forceErr
	}
	f.cards[card.ID] = card
	return nil
}

func (f *fakeRepo) SoftDeleteCard(_ context.Context, id uuid.UUID) error {
	if f.forceErr != nil {
		return f.forceErr
	}
	now := time.Now()
	f.cards[id].DeletedAt = &now
	return nil
}

func (f *fakeRepo) GetDueCards(_ context.Context, deckID uuid.UUID, before time.Time, limit int) ([]models.Card, error) {
	if f.forceErr != nil {
		return nil, f.forceErr
	}
	f.lastQueueLimit = limit
	var due []models.Card
	for _, c := range f.cards {
		if c.DeckID == deckID && c.DeletedAt == nil && !c.DueAt.After(before) {
			due = append(due, *c)
		}
	}
	sort.Slice(due, func(i, j int) bool { return due[i].DueAt.Before(due[j].DueAt) })
	if len(due) > limit {
		due = due[:limit]
	}
	return due, nil
}

func (f *fakeRepo) CreateReview(_ context.Context, review *models.Review) error {
	if f.forceErr != nil {
		return f.forceErr
	}
	review.ID = uuid.New()
	f.reviews = append(f.reviews, *review)
	return nil
}

// --- helpers ---

func seedDeck(repo *fakeRepo, userID uuid.UUID) *models.Deck {
	deck := factory.Deck().WithUserID(userID).Build()
	repo.decks[deck.ID] = &deck
	return &deck
}

func seedCard(repo *fakeRepo, deckID uuid.UUID) *models.Card {
	card := factory.Card().WithDeckID(deckID).Build()
	repo.cards[card.ID] = &card
	return &card
}

// --- deck tests ---

func TestCreateDeck(t *testing.T) {
	userID := uuid.New()

	t.Run("creates a deck for the user", func(t *testing.T) {
		repo := newFakeRepo()
		svc := NewService(repo)

		resp, status, err := svc.CreateDeck(context.Background(), userID, CreateDeckRequest{
			Name: "Japanese Basics", SourceLang: "Japanese", TargetLang: "English",
		})

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, status)
		assert.Equal(t, "Japanese Basics", resp.Name)
		assert.Equal(t, "Japanese", resp.SourceLang)
		assert.Equal(t, "English", resp.TargetLang)
		assert.Zero(t, resp.CardCount)
		require.Len(t, repo.decks, 1)
		for _, d := range repo.decks {
			assert.Equal(t, userID, d.UserID, "deck must belong to the creator")
		}
	})

	t.Run("rejects missing fields", func(t *testing.T) {
		svc := NewService(newFakeRepo())

		_, status, err := svc.CreateDeck(context.Background(), userID, CreateDeckRequest{})

		require.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("maps repository failure to 500", func(t *testing.T) {
		repo := newFakeRepo()
		repo.forceErr = errors.New("db down")
		svc := NewService(repo)

		_, status, err := svc.CreateDeck(context.Background(), userID, CreateDeckRequest{
			Name: "n", SourceLang: "a", TargetLang: "b",
		})

		require.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, status)
	})
}

func TestListDecks(t *testing.T) {
	userID := uuid.New()

	t.Run("returns only the user's decks with card counts", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		seedCard(repo, deck.ID)
		seedCard(repo, deck.ID)
		seedDeck(repo, uuid.New()) // someone else's deck
		svc := NewService(repo)

		resp, status, err := svc.ListDecks(context.Background(), userID)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		require.Len(t, resp, 1)
		assert.Equal(t, deck.ID, resp[0].ID)
		assert.EqualValues(t, 2, resp[0].CardCount)
	})

	t.Run("returns empty list when the user has no decks", func(t *testing.T) {
		svc := NewService(newFakeRepo())

		resp, status, err := svc.ListDecks(context.Background(), userID)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		assert.NotNil(t, resp, "must serialize as [] not null")
		assert.Empty(t, resp)
	})
}

func TestGetDeck(t *testing.T) {
	userID := uuid.New()

	t.Run("returns the deck with card count", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		seedCard(repo, deck.ID)
		svc := NewService(repo)

		resp, status, err := svc.GetDeck(context.Background(), userID, deck.ID)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, deck.ID, resp.ID)
		assert.EqualValues(t, 1, resp.CardCount)
	})

	t.Run("404 on unknown deck", func(t *testing.T) {
		svc := NewService(newFakeRepo())

		_, status, err := svc.GetDeck(context.Background(), userID, uuid.New())

		require.Error(t, err)
		assert.Equal(t, http.StatusNotFound, status)
	})

	t.Run("404 on another user's deck (no existence leak)", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, uuid.New())
		svc := NewService(repo)

		_, status, err := svc.GetDeck(context.Background(), userID, deck.ID)

		require.Error(t, err)
		assert.Equal(t, http.StatusNotFound, status)
	})
}

func TestUpdateDeck(t *testing.T) {
	userID := uuid.New()

	t.Run("updates name and languages", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		svc := NewService(repo)

		resp, status, err := svc.UpdateDeck(context.Background(), userID, deck.ID, UpdateDeckRequest{
			Name: "JLPT N5", SourceLang: "Japanese", TargetLang: "English",
		})

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, "JLPT N5", resp.Name)
		assert.Equal(t, "JLPT N5", repo.decks[deck.ID].Name)
	})

	t.Run("rejects invalid payload", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		svc := NewService(repo)

		_, status, err := svc.UpdateDeck(context.Background(), userID, deck.ID, UpdateDeckRequest{})

		require.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("404 on another user's deck", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, uuid.New())
		svc := NewService(repo)

		_, status, err := svc.UpdateDeck(context.Background(), userID, deck.ID, UpdateDeckRequest{
			Name: "x", SourceLang: "a", TargetLang: "b",
		})

		require.Error(t, err)
		assert.Equal(t, http.StatusNotFound, status)
	})
}

func TestDeleteDeck(t *testing.T) {
	userID := uuid.New()

	t.Run("soft deletes an owned deck", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		svc := NewService(repo)

		status, err := svc.DeleteDeck(context.Background(), userID, deck.ID)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		assert.NotNil(t, repo.decks[deck.ID].DeletedAt)
	})

	t.Run("404 on another user's deck", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, uuid.New())
		svc := NewService(repo)

		status, err := svc.DeleteDeck(context.Background(), userID, deck.ID)

		require.Error(t, err)
		assert.Equal(t, http.StatusNotFound, status)
		assert.Nil(t, repo.decks[deck.ID].DeletedAt)
	})
}

// --- card tests ---

func TestCreateCard(t *testing.T) {
	userID := uuid.New()

	t.Run("creates a card due immediately with default SRS state", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		svc := NewService(repo)

		before := time.Now()
		resp, status, err := svc.CreateCard(context.Background(), userID, deck.ID, CreateCardRequest{
			Front: "ねこ", Back: "cat",
		})

		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, status)
		assert.Equal(t, "ねこ", resp.Front)
		assert.Equal(t, "cat", resp.Back)
		assert.InDelta(t, 2.5, resp.EaseFactor, 0.0001)
		assert.Zero(t, resp.IntervalDays)
		assert.Zero(t, resp.Repetitions)
		assert.False(t, resp.DueAt.Before(before.Add(-time.Second)), "new cards are due now")
	})

	t.Run("404 when the deck belongs to someone else", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, uuid.New())
		svc := NewService(repo)

		_, status, err := svc.CreateCard(context.Background(), userID, deck.ID, CreateCardRequest{
			Front: "a", Back: "b",
		})

		require.Error(t, err)
		assert.Equal(t, http.StatusNotFound, status)
		assert.Empty(t, repo.cards)
	})

	t.Run("rejects empty front or back", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		svc := NewService(repo)

		_, status, err := svc.CreateCard(context.Background(), userID, deck.ID, CreateCardRequest{})

		require.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, status)
	})
}

func TestListCards(t *testing.T) {
	userID := uuid.New()

	t.Run("lists the deck's cards", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		seedCard(repo, deck.ID)
		seedCard(repo, deck.ID)
		svc := NewService(repo)

		resp, status, err := svc.ListCards(context.Background(), userID, deck.ID)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		assert.Len(t, resp, 2)
	})

	t.Run("404 when the deck belongs to someone else", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, uuid.New())
		svc := NewService(repo)

		_, status, err := svc.ListCards(context.Background(), userID, deck.ID)

		require.Error(t, err)
		assert.Equal(t, http.StatusNotFound, status)
	})

	t.Run("empty list serializes as [] not null", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		svc := NewService(repo)

		resp, status, err := svc.ListCards(context.Background(), userID, deck.ID)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		assert.NotNil(t, resp)
		assert.Empty(t, resp)
	})
}

func TestUpdateCard(t *testing.T) {
	userID := uuid.New()

	t.Run("updates front and back but never SRS state", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		card := seedCard(repo, deck.ID)
		card.EaseFactor = 1.9
		card.IntervalDays = 12
		card.Repetitions = 4
		svc := NewService(repo)

		resp, status, err := svc.UpdateCard(context.Background(), userID, card.ID, UpdateCardRequest{
			Front: "いぬ", Back: "dog",
		})

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, "いぬ", resp.Front)
		assert.Equal(t, "dog", resp.Back)
		assert.InDelta(t, 1.9, resp.EaseFactor, 0.0001, "editing text must not reset scheduling")
		assert.Equal(t, 12, resp.IntervalDays)
		assert.Equal(t, 4, resp.Repetitions)
	})

	t.Run("404 on unknown card", func(t *testing.T) {
		svc := NewService(newFakeRepo())

		_, status, err := svc.UpdateCard(context.Background(), userID, uuid.New(), UpdateCardRequest{
			Front: "a", Back: "b",
		})

		require.Error(t, err)
		assert.Equal(t, http.StatusNotFound, status)
	})

	t.Run("404 when the card's deck belongs to someone else", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, uuid.New())
		card := seedCard(repo, deck.ID)
		svc := NewService(repo)

		_, status, err := svc.UpdateCard(context.Background(), userID, card.ID, UpdateCardRequest{
			Front: "a", Back: "b",
		})

		require.Error(t, err)
		assert.Equal(t, http.StatusNotFound, status)
	})
}

func TestDeleteCard(t *testing.T) {
	userID := uuid.New()

	t.Run("soft deletes an owned card", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, userID)
		card := seedCard(repo, deck.ID)
		svc := NewService(repo)

		status, err := svc.DeleteCard(context.Background(), userID, card.ID)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		assert.NotNil(t, repo.cards[card.ID].DeletedAt)
	})

	t.Run("404 when the card's deck belongs to someone else", func(t *testing.T) {
		repo := newFakeRepo()
		deck := seedDeck(repo, uuid.New())
		card := seedCard(repo, deck.ID)
		svc := NewService(repo)

		status, err := svc.DeleteCard(context.Background(), userID, card.ID)

		require.Error(t, err)
		assert.Equal(t, http.StatusNotFound, status)
		assert.Nil(t, repo.cards[card.ID].DeletedAt)
	})
}
