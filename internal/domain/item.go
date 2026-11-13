package domain

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ItemType defines the supported media categories.
type ItemType string

const (
	ItemTypeMovie   ItemType = "movie"
	ItemTypeYouTube ItemType = "youtube"
	ItemTypeBook    ItemType = "book"
	ItemTypeAnime   ItemType = "anime"
	ItemTypePodcast ItemType = "podcast"
)

// IsValid validates whether an ItemType belongs to the supported set.
func (t ItemType) IsValid() bool {
	switch t {
	case ItemTypeMovie, ItemTypeYouTube, ItemTypeBook, ItemTypeAnime, ItemTypePodcast:
		return true
	default:
		return false
	}
}

// Item represents a universal content entity across different platforms.
type Item struct {
	ID          uuid.UUID              `json:"id"`
	Type        ItemType               `json:"type"`
	Title       string                 `json:"title"`
	Creator     string                 `json:"creator"`
	CoverURL    string                 `json:"cover_url"`
	ExternalURL string                 `json:"external_url"`
	SourceID    string                 `json:"source_id"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// Validate verifies the integrity of an Item entity.
func (i *Item) Validate() error {
	if !i.Type.IsValid() {
		return fmt.Errorf("%w: unsupported item type '%s'", ErrInvalidInput, i.Type)
	}
	if strings.TrimSpace(i.Title) == "" {
		return fmt.Errorf("%w: title cannot be empty", ErrInvalidInput)
	}
	if strings.TrimSpace(i.ExternalURL) == "" {
		return fmt.Errorf("%w: external_url cannot be empty", ErrInvalidInput)
	}
	if i.Metadata == nil {
		i.Metadata = make(map[string]interface{})
	}
	return nil
}

// ItemFilter provides pagination and query parameters for catalog lookups.
type ItemFilter struct {
	Type   *ItemType
	Limit  int
	Offset int
}

// ItemRepository defines the persistence port for Items.
type ItemRepository interface {
	Create(ctx context.Context, item *Item) error
	GetByID(ctx context.Context, id uuid.UUID) (*Item, error)
	GetBySourceID(ctx context.Context, sourceID string) (*Item, error)
	GetByExternalURL(ctx context.Context, url string) (*Item, error)
	List(ctx context.Context, filter ItemFilter) ([]*Item, error)
}
