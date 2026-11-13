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

// ItemRepo implements domain.ItemRepository using PostgreSQL.
type ItemRepo struct {
	pool *pgxpool.Pool
}

// NewItemRepo creates a new ItemRepo instance.
func NewItemRepo(pool *pgxpool.Pool) *ItemRepo {
	return &ItemRepo{pool: pool}
}

// Create inserts a new Item into the database.
func (r *ItemRepo) Create(ctx context.Context, item *domain.Item) error {
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	now := time.Now().UTC()
	item.CreatedAt = now
	item.UpdatedAt = now

	metadataJSON, err := json.Marshal(item.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal item metadata: %w", err)
	}

	query := `
		INSERT INTO items (id, type, title, creator, cover_url, external_url, source_id, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err = r.pool.Exec(ctx, query,
		item.ID,
		string(item.Type),
		item.Title,
		item.Creator,
		item.CoverURL,
		item.ExternalURL,
		item.SourceID,
		metadataJSON,
		item.CreatedAt,
		item.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert item: %w", err)
	}

	return nil
}

// GetByID retrieves an Item by its UUID.
func (r *ItemRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Item, error) {
	query := `
		SELECT id, type, title, creator, cover_url, external_url, source_id, metadata, created_at, updated_at
		FROM items
		WHERE id = $1
	`
	var (
		item         domain.Item
		typeStr      string
		metadataJSON []byte
	)

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&item.ID,
		&typeStr,
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
		return nil, fmt.Errorf("failed to scan item by id: %w", err)
	}

	item.Type = domain.ItemType(typeStr)
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &item.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return &item, nil
}

// GetBySourceID retrieves an Item by external source provider id (e.g. yt:123).
func (r *ItemRepo) GetBySourceID(ctx context.Context, sourceID string) (*domain.Item, error) {
	query := `
		SELECT id, type, title, creator, cover_url, external_url, source_id, metadata, created_at, updated_at
		FROM items
		WHERE source_id = $1
	`
	var (
		item         domain.Item
		typeStr      string
		metadataJSON []byte
	)

	err := r.pool.QueryRow(ctx, query, sourceID).Scan(
		&item.ID,
		&typeStr,
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
		return nil, fmt.Errorf("failed to scan item by source_id: %w", err)
	}

	item.Type = domain.ItemType(typeStr)
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &item.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return &item, nil
}

// GetByExternalURL retrieves an Item by its original external URL.
func (r *ItemRepo) GetByExternalURL(ctx context.Context, url string) (*domain.Item, error) {
	query := `
		SELECT id, type, title, creator, cover_url, external_url, source_id, metadata, created_at, updated_at
		FROM items
		WHERE external_url = $1
	`
	var (
		item         domain.Item
		typeStr      string
		metadataJSON []byte
	)

	err := r.pool.QueryRow(ctx, query, url).Scan(
		&item.ID,
		&typeStr,
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
		return nil, fmt.Errorf("failed to scan item by external_url: %w", err)
	}

	item.Type = domain.ItemType(typeStr)
	if len(metadataJSON) > 0 {
		if err := json.Unmarshal(metadataJSON, &item.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
	}

	return &item, nil
}

// List returns items based on filtering parameters.
func (r *ItemRepo) List(ctx context.Context, filter domain.ItemFilter) ([]*domain.Item, error) {
	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	query := `
		SELECT id, type, title, creator, cover_url, external_url, source_id, metadata, created_at, updated_at
		FROM items
		WHERE ($1::text IS NULL OR type = $1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	var typeParam *string
	if filter.Type != nil {
		t := string(*filter.Type)
		typeParam = &t
	}

	rows, err := r.pool.Query(ctx, query, typeParam, limit, filter.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list items: %w", err)
	}
	defer rows.Close()

	items := make([]*domain.Item, 0)
	for rows.Next() {
		var (
			item         domain.Item
			typeStr      string
			metadataJSON []byte
		)

		if err := rows.Scan(
			&item.ID,
			&typeStr,
			&item.Title,
			&item.Creator,
			&item.CoverURL,
			&item.ExternalURL,
			&item.SourceID,
			&metadataJSON,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan item row: %w", err)
		}

		item.Type = domain.ItemType(typeStr)
		if len(metadataJSON) > 0 {
			_ = json.Unmarshal(metadataJSON, &item.Metadata)
		}
		items = append(items, &item)
	}

	return items, nil
}
