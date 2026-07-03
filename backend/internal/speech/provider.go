// Package speech synthesizes audio for tutor messages and pronunciation
// behind a pluggable provider (D-007). The frontend falls back to browser
// SpeechSynthesis when no provider is configured (503).
package speech

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Provider turns text into audio.
type Provider interface {
	Synthesize(ctx context.Context, text, language string) (audio []byte, contentType string, err error)
}

// NewProvider selects the provider from SPEECH_PROVIDER. Returning nil
// means "not configured" — the service answers 503 and clients fall back.
func NewProvider() Provider {
	if os.Getenv("SPEECH_PROVIDER") != "elevenlabs" {
		return nil
	}
	key := os.Getenv("ELEVENLABS_API_KEY")
	if key == "" {
		return nil
	}
	baseURL := os.Getenv("ELEVENLABS_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.elevenlabs.io"
	}
	return NewElevenLabs(key, baseURL)
}

// ElevenLabs calls the ElevenLabs text-to-speech API. The multilingual
// model infers the language from the text itself.
type ElevenLabs struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// defaultVoice is "Rachel", ElevenLabs' multilingual default voice.
const defaultVoice = "21m00Tcm4TlvDq8ikWAM"

func NewElevenLabs(apiKey, baseURL string) *ElevenLabs {
	return &ElevenLabs{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (e *ElevenLabs) Synthesize(ctx context.Context, text, _ string) ([]byte, string, error) {
	payload, err := json.Marshal(map[string]any{
		"text":     text,
		"model_id": "eleven_multilingual_v2",
	})
	if err != nil {
		return nil, "", err
	}

	url := fmt.Sprintf("%s/v1/text-to-speech/%s?output_format=mp3_44100_128", e.baseURL, defaultVoice)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("xi-api-key", e.apiKey)

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, "", fmt.Errorf("elevenlabs returned %d: %s", resp.StatusCode, body)
	}

	audio, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "audio/mpeg"
	}
	return audio, contentType, nil
}
