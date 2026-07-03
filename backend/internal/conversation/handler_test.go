package conversation

import (
	"bytes"
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

func testRouter(repo Repository, tutor TutorProvider, user *response.UserContext) *chi.Mux {
	h := NewHandler(NewService(repo, tutor))
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
	return &response.UserContext{UserID: id.String()}, id
}

func TestHandlerRequiresAuth(t *testing.T) {
	mux := testRouter(newFakeRepo(), &fakeTutor{}, nil)

	for _, tc := range []struct{ method, path string }{
		{http.MethodGet, "/v1/conversations"},
		{http.MethodPost, "/v1/conversations"},
		{http.MethodGet, "/v1/conversations/" + uuid.NewString() + "/messages"},
		{http.MethodPost, "/v1/conversations/" + uuid.NewString() + "/messages"},
	} {
		rec := doJSON(t, mux, tc.method, tc.path, nil)
		assert.Equal(t, http.StatusUnauthorized, rec.Code, "%s %s", tc.method, tc.path)
	}
}

func TestHandlerCreateAndExchange(t *testing.T) {
	user, userID := authedUser()
	repo := newFakeRepo()
	mux := testRouter(repo, &fakeTutor{greeting: "Welcome?", reply: "What does it mean?"}, user)

	// create
	rec := doJSON(t, mux, http.MethodPost, "/v1/conversations", CreateConversationRequest{Language: "Japanese"})
	require.Equal(t, http.StatusCreated, rec.Code)
	var conv ConversationResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &conv))
	assert.Equal(t, "Practice Japanese", conv.Title)
	assert.Equal(t, userID, repo.conversations[conv.ID].UserID)

	// messages include the greeting
	rec = doJSON(t, mux, http.MethodGet, "/v1/conversations/"+conv.ID.String()+"/messages", nil)
	require.Equal(t, http.StatusOK, rec.Code)
	var msgs []MessageResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &msgs))
	require.Len(t, msgs, 1)
	assert.Equal(t, "Welcome?", msgs[0].Content)

	// send a message, get the exchange back
	rec = doJSON(t, mux, http.MethodPost, "/v1/conversations/"+conv.ID.String()+"/messages", SendMessageRequest{Content: "こんにちは"})
	require.Equal(t, http.StatusOK, rec.Code)
	var exchange ExchangeResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &exchange))
	assert.Equal(t, "こんにちは", exchange.UserMessage.Content)
	assert.Equal(t, "What does it mean?", exchange.TutorMessage.Content)
}

func TestHandlerListConversations(t *testing.T) {
	user, userID := authedUser()
	repo := newFakeRepo()
	seedConversationFor(repo, userID)
	mux := testRouter(repo, &fakeTutor{}, user)

	rec := doJSON(t, mux, http.MethodGet, "/v1/conversations", nil)

	assert.Equal(t, http.StatusOK, rec.Code)
	var body []ConversationResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Len(t, body, 1)
}

func TestHandlerRejectsMalformedIDsAndBodies(t *testing.T) {
	user, _ := authedUser()
	mux := testRouter(newFakeRepo(), &fakeTutor{}, user)

	rec := doJSON(t, mux, http.MethodGet, "/v1/conversations/not-a-uuid/messages", nil)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	req := httptest.NewRequest(http.MethodPost, "/v1/conversations", bytes.NewBufferString("{nope"))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(middleware.SetUserContext(req.Context(), *user))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

// seedConversationFor mirrors seedConversation but for a specific user id
func seedConversationFor(repo *fakeRepo, userID uuid.UUID) {
	seedConversation(repo, userID)
}
