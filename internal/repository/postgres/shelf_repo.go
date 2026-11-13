package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/recco-app/recco/internal/domain"
)

// ShelfRepo implements domain.ShelfRepository using PostgreSQL.
type ShelfRepo struct {
	pool *pgxpool.Pool
}

// NewShelfRepo creates a new ShelfRepo instance.
func NewShelfRepo(pool *pgxpool.Pool) *ShelfRepo {
	return &ShelfRepo{pool: pool}
}

// Upsert inserts a new shelf entry or updates the existing one for the same user and item.
func (r *ShelfRepo) Upsert(ctx context.Context, entry *domain.ShelfEntry) error {
	if entry.ID == uuid.Nil {
		entry.ID = uuid.New()
	}
	now := time.Now().UTC()
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = now
	}
	entry.UpdatedAt = now

	query := `
		INSERT INTO shelf_entries (
			id, user_id, item_id, recommended_by, status, rating, take, finished_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (user_id, item_id) DO UPDATE SET
			status = EXCLUDED.status,
			rating = EXCLUDED.rating,
			take = EXCLUDED.take,
			finished_at = EXCLUDED.finished_at,
			updated_at = EXCLUDED.updated_at
		RETURNING id, created_at, updated_at
	`
	err := r.pool.QueryRow(ctx, query,
		entry.ID,
		entry.UserID,
		entry.ItemID,
		entry.RecommendedBy,
		string(entry.Status),
		entry.Rating,
		entry.Take,
		entry.FinishedAt,
		entry.CreatedAt,
		entry.UpdatedAt,
	).Scan(&entry.ID, &entry.CreatedAt, &entry.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to upsert shelf entry: %w", err)
	}

	return nil
}

// GetByID retrieves a shelf entry by its ID.
func (r *ShelfRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.ShelfEntry, error) {
	query := `
		SELECT s.id, s.user_id, s.item_id, s.recommended_by, s.status, s.rating, s.take, s.finished_at, s.created_at, s.updated_at,
		       i.id, i.type, i.title, i.creator, i.cover_url, i.external_url, i.source_id, i.metadata, i.created_at, i.updated_at
		FROM shelf_entries s
		INNER JOIN items i ON s.item_id = i.id
		WHERE s.id = $1
	`
	var (
		entry        domain.ShelfEntry
		item         domain.Item
		statusStr    string
		itemTypeStr  string
		metadataJSON []byte
	)

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&entry.ID,
		&entry.UserID,
		&entry.ItemID,
		&entry.RecommendedBy,
		&statusStr,
		&entry.Rating,
		&entry.Take,
		&entry.FinishedAt,
		&entry.CreatedAt,
		&entry.UpdatedAt,
		&item.ID,
		&itemTypeStr,
		&item.Title,
		&item.Creator,
		&item.CoverURL,
		&item.ExternalURL,
		&item.SourceID,
		&metadataJSON,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to scan shelf entry by id: %w", err)
	}

	entry.Status = domain.ShelfStatus(statusStr)
	item.Type = domain.ItemType(itemTypeStr)
	if len(metadataJSON) > 0 {
		_ = json.Unmarshal(metadataJSON, &item.Metadata)
	}
	entry.Item = &item

	return &entry, nil
}

// GetByUserAndItem retrieves an entry for a specific user and item.
func (r *ShelfRepo) GetByUserAndItem(ctx context.Context, userID, itemID uuid.UUID) (*domain.ShelfEntry, error) {
	query := `
		SELECT id, user_id, item_id, recommended_by, status, rating, take, finished_at, created_at, updated_at
		FROM shelf_entries
		WHERE user_id = $1 AND item_id = $2
	`
	var (
		entry     domain.ShelfEntry
		statusStr string
	)

	err := r.pool.QueryRow(ctx, query, userID, itemID).Scan(
		&entry.ID,
		&entry.UserID,
		&entry.ItemID,
		&entry.RecommendedBy,
		&statusStr,
		&entry.Rating,
		&entry.Take,
		&entry.FinishedAt,
		&entry.CreatedAt,
		&entry.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to scan shelf entry by user and item: %w", err)
	}

	entry.Status = domain.ShelfStatus(statusStr)
	return &entry, nil
}

// List retrieves a paginated collection of entries for a user, optionally filtered by status.
func (r *ShelfRepo) List(ctx context.Context, filter domain.ShelfFilter) ([]*domain.ShelfEntry, error) {
	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	query := `
		SELECT s.id, s.user_id, s.item_id, s.recommended_by, s.status, s.rating, s.take, s.finished_at, s.created_at, s.updated_at,
		       i.id, i.type, i.title, i.creator, i.cover_url, i.external_url, i.source_id, i.metadata, i.created_at, i.updated_at
		FROM shelf_entries s
		INNER JOIN items i ON s.item_id = i.id
		WHERE s.user_id = $1
		  AND ($2::text IS NULL OR s.status = $2)
		ORDER BY 
			CASE WHEN s.status = 'finished' THEN s.finished_at ELSE s.created_at END DESC
		LIMIT $3 OFFSET $4
	`
	var statusParam *string
	if filter.Status != nil {
		s := string(*filter.Status)
		statusParam = &s
	}

	rows, err := r.pool.Query(ctx, query, filter.UserID, statusParam, limit, filter.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list shelf entries: %w", err)
	}
	defer rows.Close()

	entries := make([]*domain.ShelfEntry, 0)
	for rows.Next() {
		var (
			entry        domain.ShelfEntry
			item         domain.Item
			statusStr    string
			itemTypeStr  string
			metadataJSON []byte
		)

		if err := rows.Scan(
			&entry.ID,
			&entry.UserID,
			&entry.ItemID,
			&entry.RecommendedBy,
			&statusStr,
			&entry.Rating,
			&entry.Take,
			&entry.FinishedAt,
			&entry.CreatedAt,
			&entry.UpdatedAt,
			&item.ID,
			&itemTypeStr,
			&item.Title,
			&item.Creator,
			&item.CoverURL,
			&item.ExternalURL,
			&item.SourceID,
			&metadataJSON,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan shelf entry row: %w", err)
		}

		entry.Status = domain.ShelfStatus(statusStr)
		item.Type = domain.ItemType(itemTypeStr)
		if len(metadataJSON) > 0 {
			_ = json.Unmarshal(metadataJSON, &item.Metadata)
		}
		entry.Item = &item
		entries = append(entries, &entry)
	}

	return entries, nil
}

// Delete removes a shelf entry.
func (r *ShelfRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM shelf_entries WHERE id = $1`
	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete shelf entry: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
