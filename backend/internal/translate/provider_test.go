package translate

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDictionaryTranslatorKnownWords(t *testing.T) {
	tr := NewDictionaryTranslator()

	tests := []struct {
		name string
		text string
		from string
		to   string
		want string
	}{
		{"japanese to english", "こんにちは", "Japanese", "English", "hello"},
		{"english to japanese", "hello", "English", "Japanese", "こんにちは"},
		{"spanish to english", "gato", "Spanish", "English", "cat"},
		{"english to spanish", "cat", "English", "Spanish", "gato"},
		{"french to english", "merci", "French", "English", "thank you"},
		{"case insensitive lookup", "HELLO", "English", "Japanese", "こんにちは"},
		{"language names case insensitive", "hello", "english", "japanese", "こんにちは"},
		{"surrounding whitespace ignored", "  hello  ", "English", "Spanish", "hola"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tr.Translate(context.Background(), tt.text, tt.from, tt.to)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDictionaryTranslatorUnknownWord(t *testing.T) {
	tr := NewDictionaryTranslator()

	_, err := tr.Translate(context.Background(), "pneumonoultramicroscopic", "English", "Japanese")

	require.ErrorIs(t, err, ErrNoTranslation)
}

func TestDictionaryTranslatorUnknownLanguagePair(t *testing.T) {
	tr := NewDictionaryTranslator()

	_, err := tr.Translate(context.Background(), "hello", "English", "Klingon")

	require.ErrorIs(t, err, ErrNoTranslation)
}

func TestDictionaryTranslatorSeededDeckWords(t *testing.T) {
	// The dev-seeded starter deck words must all translate, so the demo
	// account's translate page works out of the box.
	tr := NewDictionaryTranslator()

	seeded := map[string]string{
		"こんにちは": "hello",
		"ありがとう": "thank you",
		"ねこ":    "cat",
		"みず":    "water",
		"いぬ":    "dog",
	}
	for jp, en := range seeded {
		got, err := tr.Translate(context.Background(), jp, "Japanese", "English")
		require.NoError(t, err, jp)
		assert.Equal(t, en, got)
	}
}
