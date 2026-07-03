package speech

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestElevenLabsSynthesize(t *testing.T) {
	var gotPath, gotKey string
	var gotBody map[string]any

	stub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotKey = r.Header.Get("xi-api-key")
		require.NoError(t, json.NewDecoder(r.Body).Decode(&gotBody))
		w.Header().Set("Content-Type", "audio/mpeg")
		_, _ = w.Write([]byte("mp3-bytes"))
	}))
	defer stub.Close()

	provider := NewElevenLabs("test-key", stub.URL)

	audio, contentType, err := provider.Synthesize(context.Background(), "こんにちは", "Japanese")

	require.NoError(t, err)
	assert.Equal(t, []byte("mp3-bytes"), audio)
	assert.Equal(t, "audio/mpeg", contentType)
	assert.Contains(t, gotPath, "/v1/text-to-speech/", "path carries the voice endpoint")
	assert.Equal(t, "test-key", gotKey, "API key travels in the xi-api-key header")
	assert.Equal(t, "こんにちは", gotBody["text"])
	assert.Equal(t, "eleven_multilingual_v2", gotBody["model_id"], "multilingual model handles any language")
}

func TestElevenLabsSurfacesAPIErrors(t *testing.T) {
	stub := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"detail":"invalid key"}`))
	}))
	defer stub.Close()

	provider := NewElevenLabs("bad-key", stub.URL)

	_, _, err := provider.Synthesize(context.Background(), "hi", "English")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "401")
}

func TestNewProviderSelection(t *testing.T) {
	t.Run("unconfigured yields no provider", func(t *testing.T) {
		t.Setenv("SPEECH_PROVIDER", "")
		t.Setenv("ELEVENLABS_API_KEY", "")

		assert.Nil(t, NewProvider())
	})

	t.Run("elevenlabs without a key yields no provider", func(t *testing.T) {
		t.Setenv("SPEECH_PROVIDER", "elevenlabs")
		t.Setenv("ELEVENLABS_API_KEY", "")

		assert.Nil(t, NewProvider())
	})

	t.Run("elevenlabs with a key is selected", func(t *testing.T) {
		t.Setenv("SPEECH_PROVIDER", "elevenlabs")
		t.Setenv("ELEVENLABS_API_KEY", "k")

		provider := NewProvider()

		require.NotNil(t, provider)
		assert.IsType(t, &ElevenLabs{}, provider)
	})
}
