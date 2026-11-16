package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"golang.org/x/sync/singleflight"
	"github.com/recco-app/recco/internal/domain"
	"github.com/recco-app/recco/internal/platform/cache"
	"github.com/recco-app/recco/internal/platform/enricher"
)

// EnrichmentService orchestrates concurrent metadata extraction with singleflight deduplication and caching.
type EnrichmentService struct {
	providers   []enricher.Provider
	cache       *cache.TTLCache[string, *domain.EnrichedMetadata]
	flightGroup singleflight.Group
	timeout     time.Duration
}

// NewEnrichmentService initializes the metadata enrichment pipeline.
func NewEnrichmentService(
	cacheTTL time.Duration,
	timeout time.Duration,
	providers ...enricher.Provider,
) *EnrichmentService {
	if timeout <= 0 {
		timeout = 7 * time.Second
	}
	if cacheTTL <= 0 {
		cacheTTL = 24 * time.Hour
	}

	return &EnrichmentService{
		providers: providers,
		cache:     cache.NewTTLCache[string, *domain.EnrichedMetadata](cacheTTL),
		timeout:   timeout,
	}
}

// Enrich extracts metadata for a given URL with singleflight deduplication and caching.
func (s *EnrichmentService) Enrich(ctx context.Context, rawURL string) (*domain.EnrichedMetadata, error) {
	cleanURL := strings.TrimSpace(rawURL)
	if cleanURL == "" {
		return nil, fmt.Errorf("%w: url cannot be empty", domain.ErrInvalidInput)
	}

	// 1. Check TTL cache
	if cached, ok := s.cache.Get(cleanURL); ok {
		return cached, nil
	}

	// 2. Singleflight deduplication to mitigate cache stampedes
	result, err, _ := s.flightGroup.Do(cleanURL, func() (interface{}, error) {
		// Double check cache inside flight in case a concurrent call just populated it
		if cached, ok := s.cache.Get(cleanURL); ok {
			return cached, nil
		}

		provider := s.selectProvider(cleanURL)
		if provider == nil {
			return nil, fmt.Errorf("no provider available to handle URL: %s", cleanURL)
		}

		ctxTimeout, cancel := context.WithTimeout(ctx, s.timeout)
		defer cancel()

		data, err := provider.Extract(ctxTimeout, cleanURL)
		if err != nil {
			return nil, fmt.Errorf("provider %s failed: %w", provider.Name(), err)
		}

		// Store in cache
		s.cache.Set(cleanURL, data)
		return data, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*domain.EnrichedMetadata), nil
}

func (s *EnrichmentService) selectProvider(rawURL string) enricher.Provider {
	for _, p := range s.providers {
		if p.CanHandle(rawURL) {
			return p
		}
	}
	return nil
}
