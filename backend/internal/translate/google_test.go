package translate

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLanguageCode(t *testing.T) {
	tests := []struct {
		in     string
		want   string
		wantOK bool
	}{
		{"Japanese", "ja", true},
		{"japanese", "ja", true},
		{" Spanish ", "es", true},
		{"English", "en", true},
		{"German", "de", true},
		{"French", "fr", true},
		{"Indonesian", "id", true},
		{"ja", "ja", true}, // already a code passes through
		{"EN", "en", true},
		{"Klingon", "", false},
		{"", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, ok := LanguageCode(tt.in)
			assert.Equal(t, tt.wantOK, ok)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGoogleTranslate(t *testing.T) {
	var gotKey string
	var gotBody map[string]any

	stub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotKey = r.URL.Query().Get("key")
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"translations": []map[string]string{{"translatedText": "hello"}},
			},
		})
	}))
	defer stub.Close()

	tr := NewGoogleTranslator("g-key", stub.URL)

	got, err := tr.Translate(context.Background(), "こんにちは", "Japanese", "English")

	require.NoError(t, err)
	assert.Equal(t, "hello", got)
	assert.Equal(t, "g-key", gotKey, "API key travels as a query param")
	assert.Equal(t, "こんにちは", gotBody["q"])
	assert.Equal(t, "ja", gotBody["source"], "language names map to ISO codes")
	assert.Equal(t, "en", gotBody["target"])
	assert.Equal(t, "text", gotBody["format"], "text format avoids HTML entities")
}

func TestGoogleTranslateUnknownLanguageIsNoTranslation(t *testing.T) {
	called := false
	stub := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		called = true
	}))
	defer stub.Close()

	tr := NewGoogleTranslator("k", stub.URL)

	_, err := tr.Translate(context.Background(), "hello", "English", "Klingon")

	require.ErrorIs(t, err, ErrNoTranslation, "unknown language is a 422, not a provider failure")
	assert.False(t, called, "no request leaves the process for an unmappable language")
}

func TestGoogleTranslateSurfacesAPIErrors(t *testing.T) {
	stub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error":{"message":"quota"}}`))
	}))
	defer stub.Close()

	tr := NewGoogleTranslator("k", stub.URL)

	_, err := tr.Translate(context.Background(), "hello", "English", "Spanish")

	require.Error(t, err)
	assert.NotErrorIs(t, err, ErrNoTranslation)
	assert.Contains(t, err.Error(), "403")
}

func TestNewTranslatorSelection(t *testing.T) {
	t.Run("google with a key is selected", func(t *testing.T) {
		t.Setenv("TRANSLATOR_PROVIDER", "google")
		t.Setenv("GOOGLE_TRANSLATE_API_KEY", "k")

		_, ok := NewTranslator().(*GoogleTranslator)
		assert.True(t, ok)
	})

	t.Run("google without a key falls back to the dictionary", func(t *testing.T) {
		t.Setenv("TRANSLATOR_PROVIDER", "google")
		t.Setenv("GOOGLE_TRANSLATE_API_KEY", "")

		_, ok := NewTranslator().(*DictionaryTranslator)
		assert.True(t, ok)
	})

	t.Run("default remains the dictionary", func(t *testing.T) {
		t.Setenv("TRANSLATOR_PROVIDER", "")

		_, ok := NewTranslator().(*DictionaryTranslator)
		assert.True(t, ok)
	})
}
