package domain

import "context"

// EnrichedMetadata represents cultural metadata extracted from an external source.
type EnrichedMetadata struct {
	Type        ItemType               `json:"type"`
	Title       string                 `json:"title"`
	Creator     string                 `json:"creator"`
	CoverURL    string                 `json:"cover_url"`
	ExternalURL string                 `json:"external_url"`
	SourceID    string                 `json:"source_id"`
	Description string                 `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// MetadataEnricher defines the port for extracting metadata from external URLs.
type MetadataEnricher interface {
	Enrich(ctx context.Context, rawURL string) (*EnrichedMetadata, error)
}
