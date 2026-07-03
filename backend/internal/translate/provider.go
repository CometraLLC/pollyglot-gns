package translate

import (
	"context"
	"errors"
	"strings"
)

// ErrNoTranslation means the provider has no translation for the input.
var ErrNoTranslation = errors.New("no translation available")

// Translator is the pluggable translation provider (see decision D-007).
// Implementations must be deterministic in tests; network-backed providers
// (ML service, LLM API) plug in behind this same interface.
type Translator interface {
	Translate(ctx context.Context, text, from, to string) (string, error)
}

// DictionaryTranslator is the deterministic default provider: a built-in
// bidirectional dictionary keyed by lowercase language names and words.
type DictionaryTranslator struct {
	// entries[from][to][word] = translation
	entries map[string]map[string]map[string]string
}

// pairs are stored one-directionally here and mirrored in the constructor.
var dictionarySeed = []struct {
	langA, langB string
	pairs        [][2]string
}{
	{"japanese", "english", [][2]string{
		{"こんにちは", "hello"},
		{"ありがとう", "thank you"},
		{"ねこ", "cat"},
		{"みず", "water"},
		{"いぬ", "dog"},
		{"生き甲斐", "reason for being"},
		{"さようなら", "goodbye"},
		{"はい", "yes"},
		{"いいえ", "no"},
		{"りんご", "apple"},
		{"ほん", "book"},
		{"ともだち", "friend"},
	}},
	{"spanish", "english", [][2]string{
		{"gato", "cat"},
		{"perro", "dog"},
		{"agua", "water"},
		{"hola", "hello"},
		{"gracias", "thank you"},
		{"adiós", "goodbye"},
		{"sí", "yes"},
		{"no", "no"},
		{"manzana", "apple"},
		{"libro", "book"},
		{"amigo", "friend"},
	}},
	{"french", "english", [][2]string{
		{"merci", "thank you"},
		{"bonjour", "hello"},
		{"chat", "cat"},
		{"chien", "dog"},
		{"eau", "water"},
		{"au revoir", "goodbye"},
		{"oui", "yes"},
		{"non", "no"},
		{"pomme", "apple"},
		{"livre", "book"},
		{"ami", "friend"},
	}},
	{"german", "english", [][2]string{
		{"katze", "cat"},
		{"hund", "dog"},
		{"wasser", "water"},
		{"hallo", "hello"},
		{"danke", "thank you"},
		{"tschüss", "goodbye"},
		{"ja", "yes"},
		{"nein", "no"},
		{"apfel", "apple"},
		{"buch", "book"},
		{"freund", "friend"},
	}},
	{"spanish", "japanese", [][2]string{
		{"hola", "こんにちは"},
		{"gato", "ねこ"},
		{"agua", "みず"},
		{"perro", "いぬ"},
		{"gracias", "ありがとう"},
	}},
}

func NewDictionaryTranslator() *DictionaryTranslator {
	t := &DictionaryTranslator{entries: make(map[string]map[string]map[string]string)}
	for _, seed := range dictionarySeed {
		for _, pair := range seed.pairs {
			t.add(seed.langA, seed.langB, pair[0], pair[1])
			t.add(seed.langB, seed.langA, pair[1], pair[0])
		}
	}
	return t
}

func (t *DictionaryTranslator) add(from, to, word, translation string) {
	if t.entries[from] == nil {
		t.entries[from] = make(map[string]map[string]string)
	}
	if t.entries[from][to] == nil {
		t.entries[from][to] = make(map[string]string)
	}
	t.entries[from][to][strings.ToLower(word)] = translation
}

func (t *DictionaryTranslator) Translate(_ context.Context, text, from, to string) (string, error) {
	pairs, ok := t.entries[strings.ToLower(strings.TrimSpace(from))][strings.ToLower(strings.TrimSpace(to))]
	if !ok {
		return "", ErrNoTranslation
	}

	translation, ok := pairs[strings.ToLower(strings.TrimSpace(text))]
	if !ok {
		return "", ErrNoTranslation
	}
	return translation, nil
}
