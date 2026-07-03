package translate

import (
	"os"

	"github.com/sirupsen/logrus"
)

// NewTranslator selects the provider from TRANSLATOR_PROVIDER. The
// dictionary provider is the deterministic default (D-007/D-014);
// "google" uses the Google Cloud Translation API when a key is set.
func NewTranslator() Translator {
	provider := os.Getenv("TRANSLATOR_PROVIDER")
	switch provider {
	case "", "dictionary":
		return NewDictionaryTranslator()
	case "google":
		key := os.Getenv("GOOGLE_TRANSLATE_API_KEY")
		if key == "" {
			logrus.Warn("TRANSLATOR_PROVIDER=google but GOOGLE_TRANSLATE_API_KEY is empty, falling back to dictionary")
			return NewDictionaryTranslator()
		}
		return NewGoogleTranslator(key, os.Getenv("GOOGLE_TRANSLATE_BASE_URL"))
	default:
		logrus.Warnf("unknown TRANSLATOR_PROVIDER %q, falling back to dictionary", provider)
		return NewDictionaryTranslator()
	}
}
