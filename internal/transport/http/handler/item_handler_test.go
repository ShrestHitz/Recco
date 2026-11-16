package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/recco-app/recco/internal/domain"
	"github.com/recco-app/recco/internal/platform/enricher"
	"github.com/recco-app/recco/internal/repository/memory"
	"github.com/recco-app/recco/internal/service"
	transporthttp "github.com/recco-app/recco/internal/transport/http"
	"github.com/recco-app/recco/internal/transport/http/handler"
)

func TestItemHandler_HTTP(t *testing.T) {
	itemRepo := memory.NewItemRepo()
	shelfRepo := memory.NewShelfRepo(itemRepo)

	ytProvider := enricher.NewYouTubeProvider(nil)
	enrichmentSvc := service.NewEnrichmentService(0, 0, ytProvider)

	itemService := service.NewItemService(itemRepo, enrichmentSvc)
	shelfService := service.NewShelfService(shelfRepo, itemRepo)

	itemHandler := handler.NewItemHandler(itemService)
	shelfHandler := handler.NewShelfHandler(shelfService)
	router := transporthttp.NewRouter(itemHandler, shelfHandler)

	server := httptest.NewServer(router)
	defer server.Close()

	// 1. Health check
	resp, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatalf("failed calling health check: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// 2. Create an item via POST /api/v1/items
	itemPayload := map[string]interface{}{
		"type":         "movie",
		"title":        "Inception",
		"creator":      "Christopher Nolan",
		"external_url": "https://www.themoviedb.org/movie/27205",
	}
	body, _ := json.Marshal(itemPayload)

	postResp, err := http.Post(server.URL+"/api/v1/items", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed calling POST /api/v1/items: %v", err)
	}
	if postResp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", postResp.StatusCode)
	}

	var createAPIResp struct {
		Success bool        `json:"success"`
		Data    domain.Item `json:"data"`
	}
	if err := json.NewDecoder(postResp.Body).Decode(&createAPIResp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if createAPIResp.Data.Title != "Inception" {
		t.Errorf("expected title Inception, got %s", createAPIResp.Data.Title)
	}

	// 3. Retrieve item via GET /api/v1/items/{id}
	getResp, err := http.Get(server.URL + "/api/v1/items/" + createAPIResp.Data.ID.String())
	if err != nil {
		t.Fatalf("failed calling GET item: %v", err)
	}
	if getResp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", getResp.StatusCode)
	}

	// 4. Test Metadata Enrichment endpoint: GET /api/v1/items/enrich?url=...
	enrichResp, err := http.Get(server.URL + "/api/v1/items/enrich?url=https://youtu.be/dQw4w9WgXcQ")
	if err != nil {
		t.Fatalf("failed calling GET enrich: %v", err)
	}
	if enrichResp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 for enrich, got %d", enrichResp.StatusCode)
	}

	var enrichAPIResp struct {
		Success bool                    `json:"success"`
		Data    domain.EnrichedMetadata `json:"data"`
	}
	if err := json.NewDecoder(enrichResp.Body).Decode(&enrichAPIResp); err != nil {
		t.Fatalf("failed decoding enrich response: %v", err)
	}
	if enrichAPIResp.Data.Type != domain.ItemTypeYouTube {
		t.Errorf("expected YouTube type, got %s", enrichAPIResp.Data.Type)
	}
}
