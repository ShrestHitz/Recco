package http

import (
	"net/http"

	"github.com/recco-app/recco/internal/transport/http/handler"
	"github.com/recco-app/recco/internal/transport/http/middleware"
)

// NewRouter registers application routes with middleware.
func NewRouter(
	itemHandler *handler.ItemHandler,
	shelfHandler *handler.ShelfHandler,
) http.Handler {
	mux := http.NewServeMux()

	// Health Check & Service Metadata
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		handler.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"status":  "healthy",
			"service": "recco-api",
			"phase":   "phase-1-core-foundation-shelf",
		})
	})

	// Items Endpoints
	mux.HandleFunc("POST /api/v1/items", itemHandler.CreateItem)
	mux.HandleFunc("GET /api/v1/items", itemHandler.ListItems)
	mux.HandleFunc("GET /api/v1/items/{id}", itemHandler.GetItem)

	// Personal Shelf & Queue Endpoints
	mux.HandleFunc("POST /api/v1/shelf/queue", shelfHandler.AddToQueue)
	mux.HandleFunc("POST /api/v1/shelf/finish", shelfHandler.FinishItem)
	mux.HandleFunc("POST /api/v1/shelf/drop", shelfHandler.DropItem)
	mux.HandleFunc("GET /api/v1/shelf/queue", shelfHandler.ListQueue)
	mux.HandleFunc("GET /api/v1/shelf", shelfHandler.ListShelf)
	mux.HandleFunc("DELETE /api/v1/shelf/{id}", shelfHandler.DeleteEntry)

	// Apply Middlewares (Outer to Inner: CORS -> Recovery -> RequestLogger)
	var h http.Handler = mux
	h = middleware.RequestLogger(h)
	h = middleware.Recovery(h)
	h = middleware.CORS(h)

	return h
}
