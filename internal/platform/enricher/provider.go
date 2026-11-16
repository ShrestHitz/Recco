package enricher

import (
	"context"

	"github.com/recco-app/recco/internal/domain"
)

// Provider abstracts a platform-specific metadata extractor.
type Provider interface {
	Name() string
	CanHandle(rawURL string) bool
	Extract(ctx context.Context, rawURL string) (*domain.EnrichedMetadata, error)
}
