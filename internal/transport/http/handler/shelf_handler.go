package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/recco-app/recco/internal/domain"
	"github.com/recco-app/recco/internal/service"
)

// ShelfHandler handles personal queue and taste journal endpoints.
type ShelfHandler struct {
	svc *service.ShelfService
}

// NewShelfHandler creates a new ShelfHandler.
func NewShelfHandler(svc *service.ShelfService) *ShelfHandler {
	return &ShelfHandler{svc: svc}
}

// AddToQueue handles POST /api/v1/shelf/queue.
func (h *ShelfHandler) AddToQueue(w http.ResponseWriter, r *http.Request) {
	var input service.AddToQueueInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteError(w, domain.ErrInvalidInput)
		return
	}

	entry, err := h.svc.AddToQueue(r.Context(), input)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusCreated, entry)
}

// FinishItem handles POST /api/v1/shelf/finish.
func (h *ShelfHandler) FinishItem(w http.ResponseWriter, r *http.Request) {
	var input service.FinishItemInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteError(w, domain.ErrInvalidInput)
		return
	}

	entry, err := h.svc.FinishItem(r.Context(), input)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, entry)
}

// DropItem handles POST /api/v1/shelf/drop.
func (h *ShelfHandler) DropItem(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID uuid.UUID `json:"user_id"`
		ItemID uuid.UUID `json:"item_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteError(w, domain.ErrInvalidInput)
		return
	}

	entry, err := h.svc.DropItem(r.Context(), input.UserID, input.ItemID)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, entry)
}

// ListQueue handles GET /api/v1/shelf/queue?user_id=...
func (h *ShelfHandler) ListQueue(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userIDStr := query.Get("user_id")
	if userIDStr == "" {
		WriteError(w, domain.ErrInvalidInput)
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		WriteError(w, domain.ErrInvalidInput)
		return
	}

	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))

	entries, err := h.svc.ListQueue(r.Context(), userID, limit, offset)
	if err != nil {
		WriteError(w, err)
		return
	}

	// Finite inbox / queue zero indicator
	WriteJSON(w, http.StatusOK, entries, map[string]interface{}{
		"count":      len(entries),
		"queue_zero": len(entries) == 0,
		"limit":      limit,
		"offset":     offset,
	})
}

// ListShelf handles GET /api/v1/shelf?user_id=...
func (h *ShelfHandler) ListShelf(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userIDStr := query.Get("user_id")
	if userIDStr == "" {
		WriteError(w, domain.ErrInvalidInput)
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		WriteError(w, domain.ErrInvalidInput)
		return
	}

	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))

	entries, err := h.svc.ListShelf(r.Context(), userID, limit, offset)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, entries, map[string]interface{}{
		"count":  len(entries),
		"limit":  limit,
		"offset": offset,
	})
}

// DeleteEntry handles DELETE /api/v1/shelf/{id}.
func (h *ShelfHandler) DeleteEntry(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, domain.ErrInvalidInput)
		return
	}

	if err := h.svc.RemoveEntry(r.Context(), id); err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"message": "entry successfully removed"})
}
