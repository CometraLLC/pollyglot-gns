package conversation

import (
	"encoding/json"
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
	r.Route("/conversations", func(r chi.Router) {
		r.Get("/", h.ListConversations)
		r.Post("/", h.CreateConversation)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/messages", h.GetMessages)
			r.Post("/messages", h.SendMessage)
		})
	})
}

func userID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	userCtx, ok := middleware.GetUserContext(r.Context())
	if !ok {
		response.ResponseError(w, http.StatusUnauthorized, "Unauthorized")
		return uuid.Nil, false
	}
	id, err := uuid.Parse(userCtx.UserID)
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid user ID")
		return uuid.Nil, false
	}
	return id, true
}

func pathID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid ID")
		return uuid.Nil, false
	}
	return id, true
}

func decode(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Invalid request body")
		return false
	}
	return true
}

func (h Handler) CreateConversation(w http.ResponseWriter, r *http.Request) {
	uid, ok := userID(w, r)
	if !ok {
		return
	}
	var req CreateConversationRequest
	if !decode(w, r, &req) {
		return
	}

	resp, status, err := h.service.CreateConversation(r.Context(), uid, req)
	if err != nil {
		response.ResponseError(w, status, err.Error())
		return
	}
	response.ResponseJSON(w, status, resp)
}

func (h Handler) ListConversations(w http.ResponseWriter, r *http.Request) {
	uid, ok := userID(w, r)
	if !ok {
		return
	}

	resp, status, err := h.service.ListConversations(r.Context(), uid)
	if err != nil {
		response.ResponseError(w, status, err.Error())
		return
	}
	response.ResponseJSON(w, status, resp)
}

func (h Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	uid, ok := userID(w, r)
	if !ok {
		return
	}
	conversationID, ok := pathID(w, r)
	if !ok {
		return
	}

	resp, status, err := h.service.GetMessages(r.Context(), uid, conversationID)
	if err != nil {
		response.ResponseError(w, status, err.Error())
		return
	}
	response.ResponseJSON(w, status, resp)
}

func (h Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	uid, ok := userID(w, r)
	if !ok {
		return
	}
	conversationID, ok := pathID(w, r)
	if !ok {
		return
	}
	var req SendMessageRequest
	if !decode(w, r, &req) {
		return
	}

	resp, status, err := h.service.SendMessage(r.Context(), uid, conversationID, req)
	if err != nil {
		response.ResponseError(w, status, err.Error())
		return
	}
	response.ResponseJSON(w, status, resp)
}
