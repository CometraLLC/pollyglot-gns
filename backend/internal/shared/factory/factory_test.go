package factory

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserFactoryDefaultsAndOverrides(t *testing.T) {
	user := User().Build()

	assert.NotEqual(t, uuid.Nil, user.ID)
	assert.NotEmpty(t, user.Email)
	assert.NotEmpty(t, user.Name)
	assert.True(t, user.IsActive)

	id := uuid.New()
	custom := User().WithID(id).WithEmail("marc@pollyglot.dev").WithName("Marc").Build()
	assert.Equal(t, id, custom.ID)
	assert.Equal(t, "marc@pollyglot.dev", custom.Email)
	assert.Equal(t, "Marc", custom.Name)
}

func TestUserFactoryBuildsAreIndependent(t *testing.T) {
	f := User()
	a := f.Build()
	b := f.WithEmail("other@pollyglot.dev").Build()

	assert.NotEqual(t, a.Email, b.Email, "later WithX must not mutate earlier builds")
	assert.Equal(t, a.ID, b.ID, "same factory keeps its identity unless overridden")

	assert.NotEqual(t, User().Build().ID, User().Build().ID, "separate factories get fresh IDs")
}

func TestDeckFactoryDefaultsAndOverrides(t *testing.T) {
	deck := Deck().Build()

	assert.NotEqual(t, uuid.Nil, deck.ID)
	assert.NotEqual(t, uuid.Nil, deck.UserID)
	assert.NotEmpty(t, deck.Name)
	assert.NotEmpty(t, deck.SourceLang)
	assert.NotEmpty(t, deck.TargetLang)
	assert.Nil(t, deck.DeletedAt)

	owner := uuid.New()
	custom := Deck().WithUserID(owner).WithName("JLPT N5").Build()
	assert.Equal(t, owner, custom.UserID)
	assert.Equal(t, "JLPT N5", custom.Name)
}

func TestCardFactoryDefaultsMatchNewCardSemantics(t *testing.T) {
	card := Card().Build()

	assert.NotEqual(t, uuid.Nil, card.ID)
	assert.NotEqual(t, uuid.Nil, card.DeckID)
	assert.NotEmpty(t, card.Front)
	assert.NotEmpty(t, card.Back)
	assert.InDelta(t, 2.5, card.EaseFactor, 0.0001, "new cards start at default ease")
	assert.Zero(t, card.IntervalDays)
	assert.Zero(t, card.Repetitions)
	assert.False(t, card.DueAt.After(time.Now()), "new cards are due immediately")
}

func TestCardFactoryOverrides(t *testing.T) {
	deckID := uuid.New()
	due := time.Now().Add(72 * time.Hour)

	card := Card().
		WithDeckID(deckID).
		WithFront("ねこ").
		WithBack("cat").
		WithSRS(1.9, 12, 4).
		DueAt(due).
		Build()

	assert.Equal(t, deckID, card.DeckID)
	assert.Equal(t, "ねこ", card.Front)
	assert.Equal(t, "cat", card.Back)
	assert.InDelta(t, 1.9, card.EaseFactor, 0.0001)
	assert.Equal(t, 12, card.IntervalDays)
	assert.Equal(t, 4, card.Repetitions)
	assert.Equal(t, due, card.DueAt)
}

func TestReviewFactory(t *testing.T) {
	review := Review().Build()
	assert.NotEqual(t, uuid.Nil, review.CardID)
	assert.NotEqual(t, uuid.Nil, review.UserID)
	assert.GreaterOrEqual(t, review.Rating, 0)
	assert.LessOrEqual(t, review.Rating, 4)

	cardID, userID := uuid.New(), uuid.New()
	at := time.Now().Add(-24 * time.Hour)
	custom := Review().WithCardID(cardID).WithUserID(userID).WithRating(0).ReviewedAt(at).Build()
	assert.Equal(t, cardID, custom.CardID)
	assert.Equal(t, userID, custom.UserID)
	assert.Zero(t, custom.Rating)
	assert.Equal(t, at, custom.ReviewedAt)
}

func TestSeededMirrorsDevSeeder(t *testing.T) {
	// These constants must stay in sync with backend/migrations/seeders/dev/.
	assert.Equal(t, "demo@pollyglot.dev", Seeded.Email)
	assert.Equal(t, "Password123!", Seeded.Password)
	assert.Equal(t, uuid.MustParse("a0000000-0000-4000-8000-000000000001"), Seeded.UserID)
	assert.Equal(t, uuid.MustParse("b0000000-0000-4000-8000-000000000001"), Seeded.DeckID)
}
