package speech

import (
	"context"
	"errors"
	"net/http"

	"github.com/base-go/backend/pkg/validator"
)

var ErrNotConfigured = errors.New("speech provider not configured")

type SpeakRequest struct {
	Text     string `json:"text" validate:"required,max=500"`
	Language string `json:"language" validate:"omitempty,max=50"`
}

type Service interface {
	Speak(ctx context.Context, req SpeakRequest) (audio []byte, contentType string, status int, err error)
}

type service struct {
	provider Provider
}

func NewService(provider Provider) Service {
	return &service{provider: provider}
}

func (s *service) Speak(ctx context.Context, req SpeakRequest) ([]byte, string, int, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, "", http.StatusBadRequest, err
	}

	if s.provider == nil {
		return nil, "", http.StatusServiceUnavailable, ErrNotConfigured
	}

	audio, contentType, err := s.provider.Synthesize(ctx, req.Text, req.Language)
	if err != nil {
		return nil, "", http.StatusBadGateway, err
	}
	return audio, contentType, http.StatusOK, nil
}
