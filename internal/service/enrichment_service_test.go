package service_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/recco-app/recco/internal/domain"
	"github.com/recco-app/recco/internal/service"
)

// mockSlowProvider simulates a high-latency external HTTP API call.
type mockSlowProvider struct {
	callCount int32
}

func (m *mockSlowProvider) Name() string {
	return "mock-slow"
}

func (m *mockSlowProvider) CanHandle(rawURL string) bool {
	return true
}

func (m *mockSlowProvider) Extract(ctx context.Context, rawURL string) (*domain.EnrichedMetadata, error) {
	atomic.AddInt32(&m.callCount, 1)
	time.Sleep(60 * time.Millisecond) // Simulate network delay
	return &domain.EnrichedMetadata{
		Type:        domain.ItemTypeMovie,
		Title:       "Concurrently Enriched Masterpiece",
		Creator:     "Director",
		ExternalURL: rawURL,
		CoverURL:    "https://example.com/cover.jpg",
	}, nil
}

func TestEnrichmentService_SingleflightDeduplication(t *testing.T) {
	mockProv := &mockSlowProvider{}
	svc := service.NewEnrichmentService(1*time.Hour, 5*time.Second, mockProv)

	targetURL := "https://example.com/item/12345"
	concurrentCallers := 30
	var wg sync.WaitGroup

	results := make([]*domain.EnrichedMetadata, concurrentCallers)
	errors := make([]error, concurrentCallers)

	// Launch 30 concurrent goroutines querying the exact same URL at the same time
	for i := 0; i < concurrentCallers; i++ {
		wg.Add(1)
		idx := i
		go func() {
			defer wg.Done()
			meta, err := svc.Enrich(context.Background(), targetURL)
			results[idx] = meta
			errors[idx] = err
		}()
	}

	wg.Wait()

	// Verify all concurrent callers succeeded
	for i := 0; i < concurrentCallers; i++ {
		if errors[i] != nil {
			t.Fatalf("caller %d failed: %v", i, errors[i])
		}
		if results[i].Title != "Concurrently Enriched Masterpiece" {
			t.Errorf("caller %d got unexpected title: %s", i, results[i].Title)
		}
	}

	// Singleflight must guarantee only ONE actual call was made to the slow provider!
	actualCalls := atomic.LoadInt32(&mockProv.callCount)
	if actualCalls != 1 {
		t.Errorf("expected singleflight to coalesce requests to 1 call, but got %d calls", actualCalls)
	}

	// Subsequent sequential call should hit the TTL cache and not invoke the provider again
	_, err := svc.Enrich(context.Background(), targetURL)
	if err != nil {
		t.Fatalf("cached call failed: %v", err)
	}
	if atomic.LoadInt32(&mockProv.callCount) != 1 {
		t.Errorf("expected cache hit with 1 total provider call, got %d calls", mockProv.callCount)
	}
}
