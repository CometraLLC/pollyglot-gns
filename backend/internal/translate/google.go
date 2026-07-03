package translate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// languageCodes maps human language names to ISO-639-1 for the Google API.
var languageCodes = map[string]string{
	"japanese":   "ja",
	"english":    "en",
	"spanish":    "es",
	"french":     "fr",
	"german":     "de",
	"italian":    "it",
	"portuguese": "pt",
	"korean":     "ko",
	"chinese":    "zh",
	"indonesian": "id",
	"dutch":      "nl",
	"russian":    "ru",
	"arabic":     "ar",
	"hindi":      "hi",
}

// LanguageCode resolves a language name (or an ISO code, passed through,
// lowercased) to ISO-639-1. ok is false for unknown languages.
func LanguageCode(name string) (string, bool) {
	normalized := strings.ToLower(strings.TrimSpace(name))
	if code, ok := languageCodes[normalized]; ok {
		return code, true
	}
	for _, code := range languageCodes {
		if normalized == code {
			return code, true
		}
	}
	return "", false
}

// GoogleTranslator implements Translator with the Google Cloud
// Translation API (v2). The key stays server-side.
type GoogleTranslator struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewGoogleTranslator(apiKey, baseURL string) *GoogleTranslator {
	if baseURL == "" {
		baseURL = "https://translation.googleapis.com"
	}
	return &GoogleTranslator{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{Timeout: 15 * time.Second},
	}
}

func (g *GoogleTranslator) Translate(ctx context.Context, text, from, to string) (string, error) {
	source, ok := LanguageCode(from)
	if !ok {
		return "", ErrNoTranslation
	}
	target, ok := LanguageCode(to)
	if !ok {
		return "", ErrNoTranslation
	}

	payload, err := json.Marshal(map[string]string{
		"q":      text,
		"source": source,
		"target": target,
		"format": "text",
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/language/translate/v2?key=%s", g.baseURL, g.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return "", fmt.Errorf("google translate returned %d: %s", resp.StatusCode, body)
	}

	var parsed struct {
		Data struct {
			Translations []struct {
				TranslatedText string `json:"translatedText"`
			} `json:"translations"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", err
	}
	if len(parsed.Data.Translations) == 0 {
		return "", ErrNoTranslation
	}
	return parsed.Data.Translations[0].TranslatedText, nil
}
