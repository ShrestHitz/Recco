package memory

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/recco-app/recco/internal/domain"
)

// ShelfRepo provides a thread-safe in-memory implementation of domain.ShelfRepository.
type ShelfRepo struct {
	mu       sync.RWMutex
	entries  map[uuid.UUID]*domain.ShelfEntry
	itemRepo domain.ItemRepository
}

// NewShelfRepo creates an in-memory ShelfRepo.
func NewShelfRepo(itemRepo domain.ItemRepository) *ShelfRepo {
	return &ShelfRepo{
		entries:  make(map[uuid.UUID]*domain.ShelfEntry),
		itemRepo: itemRepo,
	}
}

// Upsert adds or updates a shelf entry.
func (r *ShelfRepo) Upsert(ctx context.Context, entry *domain.ShelfEntry) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	// Check if user and item already exist
	for id, existing := range r.entries {
		if existing.UserID == entry.UserID && existing.ItemID == entry.ItemID {
			entry.ID = id
			entry.CreatedAt = existing.CreatedAt
			entry.UpdatedAt = now
			cp := *entry
			r.entries[id] = &cp
			return nil
		}
	}

	if entry.ID == uuid.Nil {
		entry.ID = uuid.New()
	}
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = now
	}
	entry.UpdatedAt = now

	cp := *entry
	r.entries[entry.ID] = &cp
	return nil
}

// GetByID retrieves a shelf entry by ID.
func (r *ShelfRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.ShelfEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entry, exists := r.entries[id]
	if !exists {
		return nil, domain.ErrNotFound
	}
	cp := *entry
	if r.itemRepo != nil {
		if item, err := r.itemRepo.GetByID(ctx, entry.ItemID); err == nil {
			cp.Item = item
		}
	}
	return &cp, nil
}

// GetByUserAndItem retrieves an entry by user ID and item ID.
func (r *ShelfRepo) GetByUserAndItem(ctx context.Context, userID, itemID uuid.UUID) (*domain.ShelfEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, entry := range r.entries {
		if entry.UserID == userID && entry.ItemID == itemID {
			cp := *entry
			if r.itemRepo != nil {
				if item, err := r.itemRepo.GetByID(ctx, entry.ItemID); err == nil {
					cp.Item = item
				}
			}
			return &cp, nil
		}
	}
	return nil, domain.ErrNotFound
}

// List returns filtered entries for a user.
func (r *ShelfRepo) List(ctx context.Context, filter domain.ShelfFilter) ([]*domain.ShelfEntry, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*domain.ShelfEntry
	for _, entry := range r.entries {
		if entry.UserID != filter.UserID {
			continue
		}
		if filter.Status != nil && entry.Status != *filter.Status {
			continue
		}
		cp := *entry
		if r.itemRepo != nil {
			if item, err := r.itemRepo.GetByID(ctx, entry.ItemID); err == nil {
				cp.Item = item
			}
		}
		result = append(result, &cp)
	}

	// Sort: finished entries by finished_at DESC, otherwise created_at DESC
	sort.Slice(result, func(i, j int) bool {
		if result[i].Status == domain.ShelfStatusFinished && result[i].FinishedAt != nil &&
			result[j].Status == domain.ShelfStatusFinished && result[j].FinishedAt != nil {
			return result[i].FinishedAt.After(*result[j].FinishedAt)
		}
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})

	start := filter.Offset
	if start > len(result) {
		return []*domain.ShelfEntry{}, nil
	}

	end := start + filter.Limit
	if filter.Limit <= 0 || end > len(result) {
		end = len(result)
	}

	return result[start:end], nil
}

// Delete removes a shelf entry.
func (r *ShelfRepo) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.entries[id]; !exists {
		return domain.ErrNotFound
	}
	delete(r.entries, id)
	return nil
}
