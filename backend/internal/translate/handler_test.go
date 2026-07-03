package translate

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/base-go/backend/pkg/middleware"
	"github.com/base-go/backend/pkg/response"
)

func testRouter(translator Translator, authed bool) *chi.Mux {
	h := NewHandler(NewService(translator))
	mux := chi.NewRouter()
	mux.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if authed {
				user := response.UserContext{UserID: "d0793289-7d71-48eb-826b-d5ea9648c1c6"}
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

func post(t *testing.T, mux *chi.Mux, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, json.NewEncoder(&buf).Encode(body))
	req := httptest.NewRequest(http.MethodPost, "/v1/translate", &buf)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func TestHandlerTranslate(t *testing.T) {
	mux := testRouter(&fakeTranslator{result: "hola"}, true)

	rec := post(t, mux, TranslateRequest{Text: "hello", From: "English", To: "Spanish"})

	assert.Equal(t, http.StatusOK, rec.Code)
	var body TranslateResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "hola", body.Translation)
}

func TestHandlerTranslateRequiresAuth(t *testing.T) {
	mux := testRouter(&fakeTranslator{result: "hola"}, false)

	rec := post(t, mux, TranslateRequest{Text: "hello", From: "English", To: "Spanish"})

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestHandlerTranslateBadJSON(t *testing.T) {
	mux := testRouter(&fakeTranslator{}, true)

	req := httptest.NewRequest(http.MethodPost, "/v1/translate", bytes.NewBufferString("{nope"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandlerTranslateNoTranslation(t *testing.T) {
	mux := testRouter(&fakeTranslator{err: ErrNoTranslation}, true)

	rec := post(t, mux, TranslateRequest{Text: "xyzzy", From: "English", To: "Klingon"})

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	var body map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "no translation available", body["error"])
}
