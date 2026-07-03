package stats

import (
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

func testRouter(repo Repository, authed bool) *chi.Mux {
	h := NewHandler(NewService(repo))
	mux := chi.NewRouter()
	mux.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if authed {
				user := response.UserContext{UserID: uuid.NewString()}
				r = r.WithContext(middleware.SetUserContext(r.Context(), user))
			}
			next.ServeHTTP(w, r)
		})
	})
	mux.Route("/v1", func(r chi.Router) {
		RegisterRoutes(r, h)
	})
	return mux
}

func TestHandlerGetStats(t *testing.T) {
	mux := testRouter(&fakeRepo{totalReviews: 7, distinctCards: 4}, true)

	req := httptest.NewRequest(http.MethodGet, "/v1/stats", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var body StatsResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.EqualValues(t, 7, body.TotalReviews)
	assert.EqualValues(t, 4, body.UniqueCards)
	assert.Len(t, body.ReviewsPerDay, 30)
}

func TestHandlerGetStatsRequiresAuth(t *testing.T) {
	mux := testRouter(&fakeRepo{}, false)

	req := httptest.NewRequest(http.MethodGet, "/v1/stats", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
