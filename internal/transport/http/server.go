package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/recco-app/recco/internal/platform/config"
)

// Server wraps the standard http.Server with graceful lifecycle control.
type Server struct {
	httpServer *http.Server
}

// NewServer initializes the HTTP server with configured timeouts.
func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%s", cfg.Port),
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
	}
}

// Start listens for incoming requests.
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
