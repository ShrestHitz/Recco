package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ShelfStatus represents the current lifecycle state of a saved item.
type ShelfStatus string

const (
	ShelfStatusQueued   ShelfStatus = "queued"   // In personal queue ("want to consume")
	ShelfStatusFinished ShelfStatus = "finished" // Completed, rated, and reviewed
	ShelfStatusDropped  ShelfStatus = "dropped"  // Discontinued
)

// IsValid checks if the status is supported.
func (s ShelfStatus) IsValid() bool {
	switch s {
	case ShelfStatusQueued, ShelfStatusFinished, ShelfStatusDropped:
		return true
	default:
		return false
	}
}

// ShelfEntry represents an individual's engagement with an item in their personal taste journal.
type ShelfEntry struct {
	ID            uuid.UUID   `json:"id"`
	UserID        uuid.UUID   `json:"user_id"`
	ItemID        uuid.UUID   `json:"item_id"`
	Item          *Item       `json:"item,omitempty"`
	RecommendedBy *uuid.UUID  `json:"recommended_by,omitempty"`
	Status        ShelfStatus `json:"status"`
	Rating        *int        `json:"rating,omitempty"` // 1 - 10
	Take          string      `json:"take,omitempty"`   // Max 280 characters one-line take
	FinishedAt    *time.Time  `json:"finished_at,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

// Validate verifies rules for shelf entries.
func (e *ShelfEntry) Validate() error {
	if !e.Status.IsValid() {
		return fmt.Errorf("%w: invalid shelf status '%s'", ErrInvalidInput, e.Status)
	}

	if len(e.Take) > 280 {
		return fmt.Errorf("%w: take reflection exceeds maximum 280 characters", ErrInvalidInput)
	}

	if e.Status == ShelfStatusFinished {
		if e.Rating == nil || *e.Rating < 1 || *e.Rating > 10 {
			return fmt.Errorf("%w: finished entries must have a rating between 1 and 10", ErrInvalidInput)
		}
		if e.FinishedAt == nil {
			now := time.Now().UTC()
			e.FinishedAt = &now
		}
	}

	return nil
}

// ShelfFilter provides criteria to query user queues and shelf collections.
type ShelfFilter struct {
	UserID uuid.UUID
	Status *ShelfStatus
	Limit  int
	Offset int
}

// ShelfRepository defines the persistence port for ShelfEntry entities.
type ShelfRepository interface {
	Upsert(ctx context.Context, entry *ShelfEntry) error
	GetByID(ctx context.Context, id uuid.UUID) (*ShelfEntry, error)
	GetByUserAndItem(ctx context.Context, userID, itemID uuid.UUID) (*ShelfEntry, error)
	List(ctx context.Context, filter ShelfFilter) ([]*ShelfEntry, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
