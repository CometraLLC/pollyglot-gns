package translate

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeTranslator drives provider behavior from tests
type fakeTranslator struct {
	result string
	err    error

	gotText, gotFrom, gotTo string
}

func (f *fakeTranslator) Translate(_ context.Context, text, from, to string) (string, error) {
	f.gotText, f.gotFrom, f.gotTo = text, from, to
	return f.result, f.err
}

func TestServiceTranslate(t *testing.T) {
	t.Run("returns the translation with echo of the request", func(t *testing.T) {
		provider := &fakeTranslator{result: "hello"}
		svc := NewService(provider)

		resp, status, err := svc.Translate(context.Background(), TranslateRequest{
			Text: "こんにちは", From: "Japanese", To: "English",
		})

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, "こんにちは", resp.Text)
		assert.Equal(t, "Japanese", resp.From)
		assert.Equal(t, "English", resp.To)
		assert.Equal(t, "hello", resp.Translation)
		assert.Equal(t, "こんにちは", provider.gotText)
	})

	t.Run("validates required fields", func(t *testing.T) {
		svc := NewService(&fakeTranslator{})

		for _, req := range []TranslateRequest{
			{},
			{Text: "hello"},
			{Text: "hello", From: "English"},
			{From: "English", To: "Spanish"},
		} {
			_, status, err := svc.Translate(context.Background(), req)
			require.Error(t, err)
			assert.Equal(t, http.StatusBadRequest, status)
		}
	})

	t.Run("rejects text over 500 characters", func(t *testing.T) {
		svc := NewService(&fakeTranslator{})
		long := make([]byte, 501)
		for i := range long {
			long[i] = 'a'
		}

		_, status, err := svc.Translate(context.Background(), TranslateRequest{
			Text: string(long), From: "English", To: "Spanish",
		})

		require.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, status)
	})

	t.Run("maps a missing translation to 422", func(t *testing.T) {
		svc := NewService(&fakeTranslator{err: ErrNoTranslation})

		_, status, err := svc.Translate(context.Background(), TranslateRequest{
			Text: "xyzzy", From: "English", To: "Spanish",
		})

		require.Error(t, err)
		assert.Equal(t, http.StatusUnprocessableEntity, status)
	})

	t.Run("maps provider failures to 502", func(t *testing.T) {
		svc := NewService(&fakeTranslator{err: errors.New("provider down")})

		_, status, err := svc.Translate(context.Background(), TranslateRequest{
			Text: "hello", From: "English", To: "Spanish",
		})

		require.Error(t, err)
		assert.Equal(t, http.StatusBadGateway, status)
	})
}
