package translate

import (
	"os"

	"github.com/sirupsen/logrus"
)

// NewTranslator selects the provider from TRANSLATOR_PROVIDER. The
// dictionary provider is the deterministic default; network-backed
// providers (ml, llm) register here as they are implemented (D-007).
func NewTranslator() Translator {
	provider := os.Getenv("TRANSLATOR_PROVIDER")
	switch provider {
	case "", "dictionary":
		return NewDictionaryTranslator()
	default:
		logrus.Warnf("unknown TRANSLATOR_PROVIDER %q, falling back to dictionary", provider)
		return NewDictionaryTranslator()
	}
}
