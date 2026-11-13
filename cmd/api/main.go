package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/recco-app/recco/internal/domain"
	"github.com/recco-app/recco/internal/platform/config"
	"github.com/recco-app/recco/internal/platform/database"
	"github.com/recco-app/recco/internal/platform/logger"
	"github.com/recco-app/recco/internal/repository/memory"
	"github.com/recco-app/recco/internal/repository/postgres"
	"github.com/recco-app/recco/internal/service"
	transporthttp "github.com/recco-app/recco/internal/transport/http"
	"github.com/recco-app/recco/internal/transport/http/handler"
)

func main() {
	// 1. Load 12-factor configuration
	cfg := config.Load()

	// 2. Setup structured logging
	log := logger.Setup(cfg.Env, "info")
	log.Info("initializing Recco backend platform",
		"env", cfg.Env,
		"port", cfg.Port,
	)

	// 3. Connect to Database (with graceful dev fallback if offline)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		itemRepo  domain.ItemRepository
		shelfRepo domain.ShelfRepository
	)

	pool, err := database.NewPostgresPool(ctx, cfg)
	if err != nil {
		if cfg.Env == "production" {
			log.Error("failed to connect to PostgreSQL primary in production", "error", err)
			os.Exit(1)
		}

		log.Warn("PostgreSQL unavailable in local dev mode, falling back to thread-safe in-memory persistence",
			"error", err,
		)
		memItemRepo := memory.NewItemRepo()
		memShelfRepo := memory.NewShelfRepo(memItemRepo)
		itemRepo = memItemRepo
		shelfRepo = memShelfRepo
	} else {
		defer pool.Close()
		log.Info("connected to PostgreSQL connection pool successfully",
			"max_conns", cfg.DBMaxConns,
			"min_conns", cfg.DBMinConns,
		)
		itemRepo = postgres.NewItemRepo(pool)
		shelfRepo = postgres.NewShelfRepo(pool)
	}

	// 4. Initialize Domain Services
	itemService := service.NewItemService(itemRepo)
	shelfService := service.NewShelfService(shelfRepo, itemRepo)

	// 5. Initialize HTTP Handlers & Router
	itemHandler := handler.NewItemHandler(itemService)
	shelfHandler := handler.NewShelfHandler(shelfService)
	router := transporthttp.NewRouter(itemHandler, shelfHandler)

	// 6. Start HTTP Server
	server := transporthttp.NewServer(cfg, router)

	serverErr := make(chan error, 1)
	go func() {
		log.Info("HTTP server listening", "addr", ":"+cfg.Port)
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	// 7. Graceful Shutdown on OS Signal
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		log.Error("server startup error", "error", err)
		os.Exit(1)
	case sig := <-shutdownSignal:
		log.Info("shutdown signal received, draining active connections", "signal", sig.String())
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("graceful shutdown failed, forcing termination", "error", err)
		os.Exit(1)
	}

	log.Info("server stopped gracefully")
}
