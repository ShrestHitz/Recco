package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/recco-app/recco/internal/domain"
	"github.com/recco-app/recco/internal/service"
)

// ItemHandler exposes HTTP endpoints for managing items.
type ItemHandler struct {
	svc *service.ItemService
}

// NewItemHandler creates a new ItemHandler.
func NewItemHandler(svc *service.ItemService) *ItemHandler {
	return &ItemHandler{svc: svc}
}

// CreateItem handles POST /api/v1/items.
func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var input service.CreateItemInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		WriteError(w, domain.ErrInvalidInput)
		return
	}

	item, err := h.svc.CreateItem(r.Context(), input)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusCreated, item)
}

// GetItem handles GET /api/v1/items/{id}.
func (h *ItemHandler) GetItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		WriteError(w, domain.ErrInvalidInput)
		return
	}

	item, err := h.svc.GetItem(r.Context(), id)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, item)
}

// ListItems handles GET /api/v1/items with optional filtering.
func (h *ItemHandler) ListItems(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))

	filter := domain.ItemFilter{
		Limit:  limit,
		Offset: offset,
	}

	if typeStr := query.Get("type"); typeStr != "" {
		t := domain.ItemType(typeStr)
		if t.IsValid() {
			filter.Type = &t
		}
	}

	items, err := h.svc.ListItems(r.Context(), filter)
	if err != nil {
		WriteError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, items, map[string]interface{}{
		"count":  len(items),
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}
