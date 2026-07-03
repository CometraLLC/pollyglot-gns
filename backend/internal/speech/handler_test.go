package speech

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

type fakeProvider struct {
	audio       []byte
	contentType string
	err         error

	gotText string
}

func (f *fakeProvider) Synthesize(_ context.Context, text, _ string) ([]byte, string, error) {
	f.gotText = text
	return f.audio, f.contentType, f.err
}

func testRouter(provider Provider, authed bool) *chi.Mux {
	h := NewHandler(NewService(provider))
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

func post(t *testing.T, mux *chi.Mux, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, json.NewEncoder(&buf).Encode(body))
	req := httptest.NewRequest(http.MethodPost, "/v1/speech", &buf)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func TestHandlerSpeakStreamsAudio(t *testing.T) {
	provider := &fakeProvider{audio: []byte("mp3"), contentType: "audio/mpeg"}
	mux := testRouter(provider, true)

	rec := post(t, mux, SpeakRequest{Text: "こんにちは", Language: "Japanese"})

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "audio/mpeg", rec.Header().Get("Content-Type"))
	assert.Equal(t, "mp3", rec.Body.String())
	assert.Equal(t, "こんにちは", provider.gotText)
}

func TestHandlerSpeakRequiresAuth(t *testing.T) {
	mux := testRouter(&fakeProvider{}, false)

	rec := post(t, mux, SpeakRequest{Text: "hi"})

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestHandlerSpeak503WhenUnconfigured(t *testing.T) {
	mux := testRouter(nil, true)

	rec := post(t, mux, SpeakRequest{Text: "hi"})

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	var body map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "speech provider not configured", body["error"])
}

func TestHandlerSpeakValidation(t *testing.T) {
	mux := testRouter(&fakeProvider{}, true)

	rec := post(t, mux, SpeakRequest{})

	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestHandlerSpeak502OnProviderFailure(t *testing.T) {
	mux := testRouter(&fakeProvider{err: errors.New("quota exceeded")}, true)

	rec := post(t, mux, SpeakRequest{Text: "hi"})

	assert.Equal(t, http.StatusBadGateway, rec.Code)
}
