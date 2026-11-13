package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/recco-app/recco/internal/domain"
	"github.com/recco-app/recco/internal/repository/memory"
	"github.com/recco-app/recco/internal/service"
	transporthttp "github.com/recco-app/recco/internal/transport/http"
	"github.com/recco-app/recco/internal/transport/http/handler"
)

func TestShelfHandler_HTTP(t *testing.T) {
	itemRepo := memory.NewItemRepo()
	shelfRepo := memory.NewShelfRepo(itemRepo)

	itemService := service.NewItemService(itemRepo)
	shelfService := service.NewShelfService(shelfRepo, itemRepo)

	itemHandler := handler.NewItemHandler(itemService)
	shelfHandler := handler.NewShelfHandler(shelfService)
	router := transporthttp.NewRouter(itemHandler, shelfHandler)

	server := httptest.NewServer(router)
	defer server.Close()

	// Seed item
	item, err := itemService.CreateItem(t.Context(), service.CreateItemInput{
		Type:        domain.ItemTypeAnime,
		Title:       "Frieren: Beyond Journey's End",
		Creator:     "Madhouse",
		ExternalURL: "https://anilist.co/anime/154587",
	})
	if err != nil {
		t.Fatalf("failed creating seed item: %v", err)
	}

	userID := uuid.New()

	// 1. Add to Queue: POST /api/v1/shelf/queue
	queuePayload, _ := json.Marshal(map[string]interface{}{
		"user_id": userID,
		"item_id": item.ID,
	})
	res, err := http.Post(server.URL+"/api/v1/shelf/queue", "application/json", bytes.NewReader(queuePayload))
	if err != nil || res.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created on queue, got %v (status %d)", err, res.StatusCode)
	}

	// 2. Fetch Queue: GET /api/v1/shelf/queue?user_id=...
	res, err = http.Get(server.URL + "/api/v1/shelf/queue?user_id=" + userID.String())
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK on list queue, got %v (status %d)", err, res.StatusCode)
	}

	var queueResp struct {
		Success bool                `json:"success"`
		Data    []domain.ShelfEntry `json:"data"`
		Meta    map[string]any      `json:"meta"`
	}
	if err := json.NewDecoder(res.Body).Decode(&queueResp); err != nil {
		t.Fatalf("failed decoding queue response: %v", err)
	}
	if len(queueResp.Data) != 1 {
		t.Fatalf("expected 1 queued entry, got %d", len(queueResp.Data))
	}

	// 3. Finish Item: POST /api/v1/shelf/finish
	finishPayload, _ := json.Marshal(map[string]interface{}{
		"user_id": userID,
		"item_id": item.ID,
		"rating":  10,
		"take":    "A quiet masterpiece reflecting on mortality and memory.",
	})
	res, err = http.Post(server.URL+"/api/v1/shelf/finish", "application/json", bytes.NewReader(finishPayload))
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK on finish item, got %v (status %d)", err, res.StatusCode)
	}

	// 4. Verify Queue is now empty (Queue Zero)
	res, err = http.Get(server.URL + "/api/v1/shelf/queue?user_id=" + userID.String())
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK on list queue, got %v", err)
	}
	_ = json.NewDecoder(res.Body).Decode(&queueResp)
	if len(queueResp.Data) != 0 {
		t.Errorf("expected 0 entries in queue, got %d", len(queueResp.Data))
	}

	// 5. Fetch Shelf (Taste Journal): GET /api/v1/shelf?user_id=...
	res, err = http.Get(server.URL + "/api/v1/shelf?user_id=" + userID.String())
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK on list shelf, got %v", err)
	}
	var shelfResp struct {
		Success bool                `json:"success"`
		Data    []domain.ShelfEntry `json:"data"`
	}
	_ = json.NewDecoder(res.Body).Decode(&shelfResp)
	if len(shelfResp.Data) != 1 {
		t.Fatalf("expected 1 shelf entry, got %d", len(shelfResp.Data))
	}
	if shelfResp.Data[0].Rating == nil || *shelfResp.Data[0].Rating != 10 {
		t.Errorf("expected shelf rating 10, got %v", shelfResp.Data[0].Rating)
	}
}
