package memory

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/recco-app/recco/internal/domain"
)

// ItemRepo provides a thread-safe in-memory implementation of domain.ItemRepository.
type ItemRepo struct {
	mu    sync.RWMutex
	items map[uuid.UUID]*domain.Item
}

// NewItemRepo creates a new in-memory ItemRepo.
func NewItemRepo() *ItemRepo {
	return &ItemRepo{
		items: make(map[uuid.UUID]*domain.Item),
	}
}

// Create inserts or replaces an item in memory.
func (r *ItemRepo) Create(ctx context.Context, item *domain.Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	now := time.Now().UTC()
	item.CreatedAt = now
	item.UpdatedAt = now

	cp := *item
	r.items[item.ID] = &cp
	return nil
}

// GetByID finds an item by ID.
func (r *ItemRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.items[id]
	if !exists {
		return nil, domain.ErrNotFound
	}
	cp := *item
	return &cp, nil
}

// GetBySourceID finds an item by source provider id.
func (r *ItemRepo) GetBySourceID(ctx context.Context, sourceID string) (*domain.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, item := range r.items {
		if item.SourceID == sourceID {
			cp := *item
			return &cp, nil
		}
	}
	return nil, domain.ErrNotFound
}

// GetByExternalURL finds an item by its external URL.
func (r *ItemRepo) GetByExternalURL(ctx context.Context, url string) (*domain.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, item := range r.items {
		if item.ExternalURL == url {
			cp := *item
			return &cp, nil
		}
	}
	return nil, domain.ErrNotFound
}

// List returns items matching the filter.
func (r *ItemRepo) List(ctx context.Context, filter domain.ItemFilter) ([]*domain.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.Item
	for _, item := range r.items {
		if filter.Type != nil && item.Type != *filter.Type {
			continue
		}
		cp := *item
		result = append(result, &cp)
	}

	start := filter.Offset
	if start > len(result) {
		return []*domain.Item{}, nil
	}

	end := start + filter.Limit
	if filter.Limit <= 0 || end > len(result) {
		end = len(result)
	}

	return result[start:end], nil
}
