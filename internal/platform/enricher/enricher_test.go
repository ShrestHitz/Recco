package enricher_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/recco-app/recco/internal/domain"
	"github.com/recco-app/recco/internal/platform/enricher"
)

func TestYouTubeProvider_CanHandleAndExtract(t *testing.T) {
	yt := enricher.NewYouTubeProvider(nil)

	if !yt.CanHandle("https://www.youtube.com/watch?v=dQw4w9WgXcQ") {
		t.Errorf("expected YouTubeProvider to handle standard watch URL")
	}
	if !yt.CanHandle("https://youtu.be/dQw4w9WgXcQ") {
		t.Errorf("expected YouTubeProvider to handle short URL")
	}
	if yt.CanHandle("https://themoviedb.org/movie/123") {
		t.Errorf("expected YouTubeProvider not to handle TMDB URL")
	}

	// Offline fallback test (without making live outbound calls)
	ctx := context.Background()
	meta, err := yt.Extract(ctx, "https://youtu.be/dQw4w9WgXcQ")
	if err != nil {
		t.Fatalf("expected fallback extraction without error, got %v", err)
	}
	if meta.Type != domain.ItemTypeYouTube {
		t.Errorf("expected YouTube item type, got %s", meta.Type)
	}
	if meta.SourceID != "yt:dQw4w9WgXcQ" {
		t.Errorf("expected source_id yt:dQw4w9WgXcQ, got %s", meta.SourceID)
	}
}

func TestOpenGraphProvider_Extraction(t *testing.T) {
	// Create a mock server that returns HTML with OpenGraph tags
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		html := `
			<!DOCTYPE html>
			<html>
			<head>
				<title>Dune: Part Two (2024)</title>
				<meta property="og:title" content="Dune: Part Two">
				<meta property="og:image" content="https://image.tmdb.org/t/p/w500/dune.jpg">
				<meta property="og:description" content="Paul Atreides unites with Chani and the Fremen.">
				<meta property="og:site_name" content="The Movie Database (TMDB)">
				<meta name="author" content="Denis Villeneuve">
			</head>
			<body><h1>Dune: Part Two</h1></body>
			</html>
		`
		fmt.Fprint(w, html)
	}))
	defer mockServer.Close()

	og := enricher.NewOpenGraphProvider(mockServer.Client())
	ctx := context.Background()

	meta, err := og.Extract(ctx, mockServer.URL+"/movie/693134-dune-part-two")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if meta.Title != "Dune: Part Two" {
		t.Errorf("expected title 'Dune: Part Two', got '%s'", meta.Title)
	}
	if meta.Creator != "Denis Villeneuve" {
		t.Errorf("expected creator 'Denis Villeneuve', got '%s'", meta.Creator)
	}
	if meta.CoverURL != "https://image.tmdb.org/t/p/w500/dune.jpg" {
		t.Errorf("expected cover 'https://image.tmdb.org/t/p/w500/dune.jpg', got '%s'", meta.CoverURL)
	}
}
