package translate

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/base-go/backend/pkg/middleware"
	"github.com/base-go/backend/pkg/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return Handler{service: service}
}

// RegisterRoutes mounts the module's routes on r. Callers are expected to
// have JWT auth middleware installed; handlers still verify the context.
func RegisterRoutes(r chi.Router, h Handler) {
	r.Post("/translate", h.Translate)
}

func (h Handler) Translate(w http.ResponseWriter, r *http.Request) {
	if _, ok := middleware.GetUserContext(r.Context()); !ok {
		response.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req TranslateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	resp, status, err := h.service.Translate(r.Context(), req)
	if err != nil {
		response.ResponseError(w, status, err.Error())
		return
	}
	response.ResponseJSON(w, status, resp)
}
