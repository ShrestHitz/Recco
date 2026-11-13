package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/recco-app/recco/internal/domain"
)

// ShelfService coordinates the personal queue and taste journal workflow.
type ShelfService struct {
	shelfRepo domain.ShelfRepository
	itemRepo  domain.ItemRepository
}

// NewShelfService creates a new ShelfService.
func NewShelfService(shelfRepo domain.ShelfRepository, itemRepo domain.ItemRepository) *ShelfService {
	return &ShelfService{
		shelfRepo: shelfRepo,
		itemRepo:  itemRepo,
	}
}

// AddToQueueInput holds data needed to save an item to the user's queue.
type AddToQueueInput struct {
	UserID        uuid.UUID  `json:"user_id"`
	ItemID        uuid.UUID  `json:"item_id"`
	RecommendedBy *uuid.UUID `json:"recommended_by,omitempty"`
}

// AddToQueue places an item into the user's queue ("want to consume").
func (s *ShelfService) AddToQueue(ctx context.Context, input AddToQueueInput) (*domain.ShelfEntry, error) {
	if input.UserID == uuid.Nil {
		return nil, fmt.Errorf("%w: user_id is required", domain.ErrInvalidInput)
	}
	if input.ItemID == uuid.Nil {
		return nil, fmt.Errorf("%w: item_id is required", domain.ErrInvalidInput)
	}

	// Verify item exists
	item, err := s.itemRepo.GetByID(ctx, input.ItemID)
	if err != nil {
		return nil, fmt.Errorf("item not found: %w", err)
	}

	entry := &domain.ShelfEntry{
		ID:            uuid.New(),
		UserID:        input.UserID,
		ItemID:        input.ItemID,
		RecommendedBy: input.RecommendedBy,
		Status:        domain.ShelfStatusQueued,
		Item:          item,
	}

	if err := entry.Validate(); err != nil {
		return nil, err
	}

	if err := s.shelfRepo.Upsert(ctx, entry); err != nil {
		return nil, fmt.Errorf("failed to save entry to queue: %w", err)
	}

	return entry, nil
}

// FinishItemInput captures rating and personal take when marking an item completed.
type FinishItemInput struct {
	UserID uuid.UUID `json:"user_id"`
	ItemID uuid.UUID `json:"item_id"`
	Rating int       `json:"rating"` // 1 to 10
	Take   string    `json:"take"`   // Optional one-line reflection
}

// FinishItem transitions an item from the personal queue to the rated shelf.
func (s *ShelfService) FinishItem(ctx context.Context, input FinishItemInput) (*domain.ShelfEntry, error) {
	if input.UserID == uuid.Nil || input.ItemID == uuid.Nil {
		return nil, fmt.Errorf("%w: user_id and item_id are required", domain.ErrInvalidInput)
	}

	entry, err := s.shelfRepo.GetByUserAndItem(ctx, input.UserID, input.ItemID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			// If not in queue yet, allow directly finishing an item seen elsewhere
			entry = &domain.ShelfEntry{
				ID:     uuid.New(),
				UserID: input.UserID,
				ItemID: input.ItemID,
			}
		} else {
			return nil, fmt.Errorf("failed to fetch shelf entry: %w", err)
		}
	}

	now := time.Now().UTC()
	entry.Status = domain.ShelfStatusFinished
	entry.Rating = &input.Rating
	entry.Take = input.Take
	entry.FinishedAt = &now

	if err := entry.Validate(); err != nil {
		return nil, err
	}

	if err := s.shelfRepo.Upsert(ctx, entry); err != nil {
		return nil, fmt.Errorf("failed to update shelf entry: %w", err)
	}

	// Hydrate item details if available
	if item, err := s.itemRepo.GetByID(ctx, entry.ItemID); err == nil {
		entry.Item = item
	}

	return entry, nil
}

// DropItem marks an item as discontinued/dropped.
func (s *ShelfService) DropItem(ctx context.Context, userID, itemID uuid.UUID) (*domain.ShelfEntry, error) {
	entry, err := s.shelfRepo.GetByUserAndItem(ctx, userID, itemID)
	if err != nil {
		return nil, err
	}

	entry.Status = domain.ShelfStatusDropped
	if err := entry.Validate(); err != nil {
		return nil, err
	}

	if err := s.shelfRepo.Upsert(ctx, entry); err != nil {
		return nil, fmt.Errorf("failed to update shelf entry to dropped: %w", err)
	}

	return entry, nil
}

// ListQueue returns the user's active, finite queue of items waiting to be consumed.
func (s *ShelfService) ListQueue(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.ShelfEntry, error) {
	status := domain.ShelfStatusQueued
	return s.shelfRepo.List(ctx, domain.ShelfFilter{
		UserID: userID,
		Status: &status,
		Limit:  limit,
		Offset: offset,
	})
}

// ListShelf returns the user's finished taste journal entries with ratings and reflections.
func (s *ShelfService) ListShelf(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.ShelfEntry, error) {
	status := domain.ShelfStatusFinished
	return s.shelfRepo.List(ctx, domain.ShelfFilter{
		UserID: userID,
		Status: &status,
		Limit:  limit,
		Offset: offset,
	})
}

// RemoveEntry completely removes an entry from the queue or shelf.
func (s *ShelfService) RemoveEntry(ctx context.Context, id uuid.UUID) error {
	return s.shelfRepo.Delete(ctx, id)
}
