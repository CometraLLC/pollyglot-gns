package decks

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/base-go/backend/pkg/middleware"
	"github.com/base-go/backend/pkg/response"
)

// fakeService records calls and returns canned values
type fakeService struct {
	deckResp  *DeckResponse
	decksResp []DeckResponse
	cardResp  *CardResponse
	cardsResp []CardResponse
	status    int
	err       error

	gotUserID uuid.UUID
	gotDeckID uuid.UUID
	gotCardID uuid.UUID
}

func (f *fakeService) CreateDeck(_ context.Context, userID uuid.UUID, _ CreateDeckRequest) (*DeckResponse, int, error) {
	f.gotUserID = userID
	return f.deckResp, f.status, f.err
}

func (f *fakeService) ListDecks(_ context.Context, userID uuid.UUID) ([]DeckResponse, int, error) {
	f.gotUserID = userID
	return f.decksResp, f.status, f.err
}

func (f *fakeService) GetDeck(_ context.Context, userID, deckID uuid.UUID) (*DeckResponse, int, error) {
	f.gotUserID, f.gotDeckID = userID, deckID
	return f.deckResp, f.status, f.err
}

func (f *fakeService) UpdateDeck(_ context.Context, userID, deckID uuid.UUID, _ UpdateDeckRequest) (*DeckResponse, int, error) {
	f.gotUserID, f.gotDeckID = userID, deckID
	return f.deckResp, f.status, f.err
}

func (f *fakeService) DeleteDeck(_ context.Context, userID, deckID uuid.UUID) (int, error) {
	f.gotUserID, f.gotDeckID = userID, deckID
	return f.status, f.err
}

func (f *fakeService) CreateCard(_ context.Context, userID, deckID uuid.UUID, _ CreateCardRequest) (*CardResponse, int, error) {
	f.gotUserID, f.gotDeckID = userID, deckID
	return f.cardResp, f.status, f.err
}

func (f *fakeService) ListCards(_ context.Context, userID, deckID uuid.UUID) ([]CardResponse, int, error) {
	f.gotUserID, f.gotDeckID = userID, deckID
	return f.cardsResp, f.status, f.err
}

func (f *fakeService) UpdateCard(_ context.Context, userID, cardID uuid.UUID, _ UpdateCardRequest) (*CardResponse, int, error) {
	f.gotUserID, f.gotCardID = userID, cardID
	return f.cardResp, f.status, f.err
}

func (f *fakeService) DeleteCard(_ context.Context, userID, cardID uuid.UUID) (int, error) {
	f.gotUserID, f.gotCardID = userID, cardID
	return f.status, f.err
}

// testRouter mounts the handler exactly as pkg/router does, with a stub
// auth middleware injecting the given user (or none).
func testRouter(svc Service, user *response.UserContext) *chi.Mux {
	h := NewHandler(svc)
	mux := chi.NewRouter()
	mux.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if user != nil {
				r = r.WithContext(middleware.SetUserContext(r.Context(), *user))
			}
			next.ServeHTTP(w, r)
		})
	})
	mux.Route("/v1", func(r chi.Router) {
		RegisterRoutes(r, h)
	})
	return mux
}

func doJSON(t *testing.T, mux *chi.Mux, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		require.NoError(t, json.NewEncoder(&buf).Encode(body))
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func authedUser() (*response.UserContext, uuid.UUID) {
	id := uuid.New()
	return &response.UserContext{UserID: id.String(), Email: "test@pollyglot.dev"}, id
}

func TestHandlerRequiresAuthContext(t *testing.T) {
	svc := &fakeService{status: http.StatusOK}
	mux := testRouter(svc, nil)

	for _, tc := range []struct{ method, path string }{
		{http.MethodGet, "/v1/decks"},
		{http.MethodPost, "/v1/decks"},
		{http.MethodGet, "/v1/decks/" + uuid.NewString()},
		{http.MethodPost, "/v1/decks/" + uuid.NewString() + "/cards"},
		{http.MethodPut, "/v1/cards/" + uuid.NewString()},
		{http.MethodDelete, "/v1/cards/" + uuid.NewString()},
	} {
		rec := doJSON(t, mux, tc.method, tc.path, nil)
		assert.Equal(t, http.StatusUnauthorized, rec.Code, "%s %s", tc.method, tc.path)
	}
}

func TestHandlerRejectsMalformedIDs(t *testing.T) {
	user, _ := authedUser()
	svc := &fakeService{status: http.StatusOK}
	mux := testRouter(svc, user)

	rec := doJSON(t, mux, http.MethodGet, "/v1/decks/not-a-uuid", nil)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	rec = doJSON(t, mux, http.MethodDelete, "/v1/cards/42", nil)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandlerCreateDeck(t *testing.T) {
	user, userID := authedUser()
	svc := &fakeService{
		deckResp: &DeckResponse{ID: uuid.New(), Name: "Japanese Basics"},
		status:   http.StatusCreated,
	}
	mux := testRouter(svc, user)

	rec := doJSON(t, mux, http.MethodPost, "/v1/decks", CreateDeckRequest{
		Name: "Japanese Basics", SourceLang: "Japanese", TargetLang: "English",
	})

	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, userID, svc.gotUserID, "handler must pass the authenticated user's ID")

	var body DeckResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "Japanese Basics", body.Name)
}

func TestHandlerCreateDeckRejectsBadJSON(t *testing.T) {
	user, _ := authedUser()
	svc := &fakeService{status: http.StatusCreated}
	mux := testRouter(svc, user)

	req := httptest.NewRequest(http.MethodPost, "/v1/decks", bytes.NewBufferString("{not json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandlerListDecks(t *testing.T) {
	user, _ := authedUser()
	svc := &fakeService{
		decksResp: []DeckResponse{{Name: "A"}, {Name: "B"}},
		status:    http.StatusOK,
	}
	mux := testRouter(svc, user)

	rec := doJSON(t, mux, http.MethodGet, "/v1/decks", nil)

	assert.Equal(t, http.StatusOK, rec.Code)
	var body []DeckResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Len(t, body, 2)
}

func TestHandlerServiceErrorPassthrough(t *testing.T) {
	user, _ := authedUser()
	svc := &fakeService{status: http.StatusNotFound, err: ErrDeckNotFound}
	mux := testRouter(svc, user)

	rec := doJSON(t, mux, http.MethodGet, "/v1/decks/"+uuid.NewString(), nil)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	var body map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "deck not found", body["error"])
}

func TestHandlerCardRoutes(t *testing.T) {
	user, _ := authedUser()
	deckID := uuid.New()
	cardID := uuid.New()
	svc := &fakeService{
		cardResp:  &CardResponse{ID: cardID, Front: "ねこ", Back: "cat"},
		cardsResp: []CardResponse{{ID: cardID}},
		status:    http.StatusOK,
	}
	mux := testRouter(svc, user)

	rec := doJSON(t, mux, http.MethodGet, "/v1/decks/"+deckID.String()+"/cards", nil)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, deckID, svc.gotDeckID)

	svc.status = http.StatusCreated
	rec = doJSON(t, mux, http.MethodPost, "/v1/decks/"+deckID.String()+"/cards", CreateCardRequest{Front: "ねこ", Back: "cat"})
	assert.Equal(t, http.StatusCreated, rec.Code)

	svc.status = http.StatusOK
	rec = doJSON(t, mux, http.MethodPut, "/v1/cards/"+cardID.String(), UpdateCardRequest{Front: "いぬ", Back: "dog"})
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, cardID, svc.gotCardID)

	rec = doJSON(t, mux, http.MethodDelete, "/v1/cards/"+cardID.String(), nil)
	assert.Equal(t, http.StatusOK, rec.Code)
}
