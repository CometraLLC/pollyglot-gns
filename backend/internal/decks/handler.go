package decks

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
	r.Route("/decks", func(r chi.Router) {
		r.Get("/", h.ListDecks)
		r.Post("/", h.CreateDeck)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetDeck)
			r.Put("/", h.UpdateDeck)
			r.Delete("/", h.DeleteDeck)

			r.Get("/cards", h.ListCards)
			r.Post("/cards", h.CreateCard)
		})
	})

	r.Route("/cards", func(r chi.Router) {
		r.Put("/{id}", h.UpdateCard)
		r.Delete("/{id}", h.DeleteCard)
	})
}

// userID extracts the authenticated user or writes 401/400 and reports false.
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

// pathID parses the {id} URL parameter or writes 400 and reports false.
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

// --- deck handlers ---

func (h Handler) CreateDeck(w http.ResponseWriter, r *http.Request) {
	uid, ok := userID(w, r)
	if !ok {
		return
	}
	var req CreateDeckRequest
	if !decode(w, r, &req) {
		return
	}

	resp, status, err := h.service.CreateDeck(r.Context(), uid, req)
	if err != nil {
		response.ResponseError(w, status, err.Error())
		return
	}
	response.ResponseJSON(w, status, resp)
}

func (h Handler) ListDecks(w http.ResponseWriter, r *http.Request) {
	uid, ok := userID(w, r)
	if !ok {
		return
	}

	resp, status, err := h.service.ListDecks(r.Context(), uid)
	if err != nil {
		response.ResponseError(w, status, err.Error())
		return
	}
	response.ResponseJSON(w, status, resp)
}

func (h Handler) GetDeck(w http.ResponseWriter, r *http.Request) {
	uid, ok := userID(w, r)
	if !ok {
		return
	}
	deckID, ok := pathID(w, r)
	if !ok {
		return
	}

	resp, status, err := h.service.GetDeck(r.Context(), uid, deckID)
	if err != nil {
		response.ResponseError(w, status, err.Error())
		return
	}
	response.ResponseJSON(w, status, resp)
}

func (h Handler) UpdateDeck(w http.ResponseWriter, r *http.Request) {
	uid, ok := userID(w, r)
	if !ok {
		return
	}
	deckID, ok := pathID(w, r)
	if !ok {
		return
	}
	var req UpdateDeckRequest
	if !decode(w, r, &req) {
		return
	}

	resp, status, err := h.service.UpdateDeck(r.Context(), uid, deckID, req)
	if err != nil {
		response.ResponseError(w, status, err.Error())
		return
	}
	response.ResponseJSON(w, status, resp)
}

func (h Handler) DeleteDeck(w http.ResponseWriter, r *http.Request) {
	uid, ok := userID(w, r)
	if !ok {
		return
	}
	deckID, ok := pathID(w, r)
	if !ok {
		return
	}

	status, err := h.service.DeleteDeck(r.Context(), uid, deckID)
	if err != nil {
		response.ResponseError(w, status, err.Error())
		return
	}
	response.ResponseJSON(w, status, map[string]string{"message": "Deck deleted"})
}

// --- card handlers ---

func (h Handler) CreateCard(w http.ResponseWriter, r *http.Request) {
	uid, ok := userID(w, r)
	if !ok {
		return
	}
	deckID, ok := pathID(w, r)
	if !ok {
		return
	}
	var req CreateCardRequest
	if !decode(w, r, &req) {
		return
	}

	resp, status, err := h.service.CreateCard(r.Context(), uid, deckID, req)
	if err != nil {
		response.ResponseError(w, status, err.Error())
		return
	}
	response.ResponseJSON(w, status, resp)
}

func (h Handler) ListCards(w http.ResponseWriter, r *http.Request) {
	uid, ok := userID(w, r)
	if !ok {
		return
	}
	deckID, ok := pathID(w, r)
	if !ok {
		return
	}

	resp, status, err := h.service.ListCards(r.Context(), uid, deckID)
	if err != nil {
		response.ResponseError(w, status, err.Error())
		return
	}
	response.ResponseJSON(w, status, resp)
}

func (h Handler) UpdateCard(w http.ResponseWriter, r *http.Request) {
	uid, ok := userID(w, r)
	if !ok {
		return
	}
	cardID, ok := pathID(w, r)
	if !ok {
		return
	}
	var req UpdateCardRequest
	if !decode(w, r, &req) {
		return
	}

	resp, status, err := h.service.UpdateCard(r.Context(), uid, cardID, req)
	if err != nil {
		response.ResponseError(w, status, err.Error())
		return
	}
	response.ResponseJSON(w, status, resp)
}

func (h Handler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	uid, ok := userID(w, r)
	if !ok {
		return
	}
	cardID, ok := pathID(w, r)
	if !ok {
		return
	}

	status, err := h.service.DeleteCard(r.Context(), uid, cardID)
	if err != nil {
		response.ResponseError(w, status, err.Error())
		return
	}
	response.ResponseJSON(w, status, map[string]string{"message": "Card deleted"})
}
