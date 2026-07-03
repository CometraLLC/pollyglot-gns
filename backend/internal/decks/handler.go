package decks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
			r.Get("/queue", h.GetStudyQueue)
			r.Get("/export", h.ExportDeck)
			r.Post("/import", h.ImportDeck)
			r.Post("/share", h.ShareDeck)
			r.Delete("/share", h.UnshareDeck)
		})
	})

	r.Route("/shared/{code}", func(r chi.Router) {
		r.Get("/", h.GetSharedDeck)
		r.Post("/clone", h.CloneSharedDeck)
	})

	r.Route("/cards", func(r chi.Router) {
		r.Put("/{id}", h.UpdateCard)
		r.Delete("/{id}", h.DeleteCard)
		r.Post("/{id}/review", h.ReviewCard)
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

func (h Handler) ReviewCard(w http.ResponseWriter, r *http.Request) {
	uid, ok := userID(w, r)
	if !ok {
		return
	}
	cardID, ok := pathID(w, r)
	if !ok {
		return
	}
	var req ReviewCardRequest
	if !decode(w, r, &req) {
		return
	}

	resp, status, err := h.service.ReviewCard(r.Context(), uid, cardID, req)
	if err != nil {
		response.ResponseError(w, status, err.Error())
		return
	}
	response.ResponseJSON(w, status, resp)
}

func (h Handler) GetStudyQueue(w http.ResponseWriter, r *http.Request) {
	uid, ok := userID(w, r)
	if !ok {
		return
	}
	deckID, ok := pathID(w, r)
	if !ok {
		return
	}
	// invalid or absent limits become 0; the service applies the default
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	resp, status, err := h.service.GetStudyQueue(r.Context(), uid, deckID, limit)
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

func (h Handler) ExportDeck(w http.ResponseWriter, r *http.Request) {
	uid, ok := userID(w, r)
	if !ok {
		return
	}
	deckID, ok := pathID(w, r)
	if !ok {
		return
	}
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "csv"
	}

	filename, content, status, err := h.service.ExportDeck(r.Context(), uid, deckID, format)
	if err != nil {
		response.ResponseError(w, status, err.Error())
		return
	}

	contentType := "text/csv; charset=utf-8"
	if format == "tsv" {
		contentType = "text/tab-separated-values; charset=utf-8"
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	w.WriteHeader(status)
	_, _ = w.Write([]byte(content))
}

func (h Handler) ImportDeck(w http.ResponseWriter, r *http.Request) {
	uid, ok := userID(w, r)
	if !ok {
		return
	}
	deckID, ok := pathID(w, r)
	if !ok {
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.ResponseError(w, http.StatusBadRequest, "Missing file upload field \"file\"")
		return
	}
	defer func() { _ = file.Close() }()

	format := r.URL.Query().Get("format")
	if format == "" {
		if strings.HasSuffix(strings.ToLower(header.Filename), ".tsv") {
			format = "tsv"
		} else {
			format = "csv"
		}
	}

	result, status, err := h.service.ImportDeck(r.Context(), uid, deckID, file, format)
	if err != nil {
		response.ResponseError(w, status, err.Error())
		return
	}
	response.ResponseJSON(w, status, result)
}

func (h Handler) ShareDeck(w http.ResponseWriter, r *http.Request) {
	uid, ok := userID(w, r)
	if !ok {
		return
	}
	deckID, ok := pathID(w, r)
	if !ok {
		return
	}

	resp, status, err := h.service.ShareDeck(r.Context(), uid, deckID)
	if err != nil {
		response.ResponseError(w, status, err.Error())
		return
	}
	response.ResponseJSON(w, status, resp)
}

func (h Handler) UnshareDeck(w http.ResponseWriter, r *http.Request) {
	uid, ok := userID(w, r)
	if !ok {
		return
	}
	deckID, ok := pathID(w, r)
	if !ok {
		return
	}

	status, err := h.service.UnshareDeck(r.Context(), uid, deckID)
	if err != nil {
		response.ResponseError(w, status, err.Error())
		return
	}
	response.ResponseJSON(w, status, map[string]string{"message": "Sharing disabled"})
}

func (h Handler) GetSharedDeck(w http.ResponseWriter, r *http.Request) {
	if _, ok := userID(w, r); !ok {
		return
	}
	code := chi.URLParam(r, "code")

	resp, status, err := h.service.GetSharedDeck(r.Context(), code)
	if err != nil {
		response.ResponseError(w, status, err.Error())
		return
	}
	response.ResponseJSON(w, status, resp)
}

func (h Handler) CloneSharedDeck(w http.ResponseWriter, r *http.Request) {
	uid, ok := userID(w, r)
	if !ok {
		return
	}
	code := chi.URLParam(r, "code")

	resp, status, err := h.service.CloneSharedDeck(r.Context(), uid, code)
	if err != nil {
		response.ResponseError(w, status, err.Error())
		return
	}
	response.ResponseJSON(w, status, resp)
}
