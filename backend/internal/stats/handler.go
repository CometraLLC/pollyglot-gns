package stats

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

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
	r.Get("/stats", h.GetStats)
}

func (h Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := middleware.GetUserContext(r.Context())
	if !ok {
		response.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	userID, err := uuid.Parse(userCtx.UserID)
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	resp, status, err := h.service.GetStats(r.Context(), userID)
	if err != nil {
		response.ResponseError(w, status, err.Error())
		return
	}
	response.ResponseJSON(w, status, resp)
}
