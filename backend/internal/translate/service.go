package translate

import (
	"context"
	"errors"
	"net/http"

	"github.com/base-go/backend/pkg/validator"
)

type TranslateRequest struct {
	Text string `json:"text" validate:"required,max=500"`
	From string `json:"from" validate:"required,max=50"`
	To   string `json:"to" validate:"required,max=50"`
}

type TranslateResponse struct {
	Text        string `json:"text"`
	From        string `json:"from"`
	To          string `json:"to"`
	Translation string `json:"translation"`
}

type Service interface {
	Translate(ctx context.Context, req TranslateRequest) (*TranslateResponse, int, error)
}

type service struct {
	provider Translator
}

func NewService(provider Translator) Service {
	return &service{provider: provider}
}

func (s *service) Translate(ctx context.Context, req TranslateRequest) (*TranslateResponse, int, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, http.StatusBadRequest, err
	}

	translation, err := s.provider.Translate(ctx, req.Text, req.From, req.To)
	if err != nil {
		if errors.Is(err, ErrNoTranslation) {
			return nil, http.StatusUnprocessableEntity, ErrNoTranslation
		}
		return nil, http.StatusBadGateway, err
	}

	return &TranslateResponse{
		Text:        req.Text,
		From:        req.From,
		To:          req.To,
		Translation: translation,
	}, http.StatusOK, nil
}
