package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/recco-app/recco/internal/domain"
)

// ItemService handles business operations on cultural items.
type ItemService struct {
	repo domain.ItemRepository
}

// NewItemService creates a new ItemService.
func NewItemService(repo domain.ItemRepository) *ItemService {
	return &ItemService{repo: repo}
}

// CreateItemInput contains parameters to register a new item.
type CreateItemInput struct {
	Type        domain.ItemType        `json:"type"`
	Title       string                 `json:"title"`
	Creator     string                 `json:"creator"`
	CoverURL    string                 `json:"cover_url"`
	ExternalURL string                 `json:"external_url"`
	SourceID    string                 `json:"source_id"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// CreateItem validates and registers a new item, deduplicating existing external URLs.
func (s *ItemService) CreateItem(ctx context.Context, input CreateItemInput) (*domain.Item, error) {
	// 1. Check if external URL already exists to deduplicate catalog entries
	if input.ExternalURL != "" {
		existing, err := s.repo.GetByExternalURL(ctx, input.ExternalURL)
		if err == nil && existing != nil {
			return existing, nil
		}
		if err != nil && !errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("failed checking duplicate URL: %w", err)
		}
	}

	// 2. Check if source ID already exists
	if input.SourceID != "" {
		existing, err := s.repo.GetBySourceID(ctx, input.SourceID)
		if err == nil && existing != nil {
			return existing, nil
		}
		if err != nil && !errors.Is(err, domain.ErrNotFound) {
			return nil, fmt.Errorf("failed checking duplicate source ID: %w", err)
		}
	}

	item := &domain.Item{
		ID:          uuid.New(),
		Type:        input.Type,
		Title:       input.Title,
		Creator:     input.Creator,
		CoverURL:    input.CoverURL,
		ExternalURL: input.ExternalURL,
		SourceID:    input.SourceID,
		Metadata:    input.Metadata,
	}

	if err := item.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to persist item: %w", err)
	}

	return item, nil
}

// GetItem retrieves an item by its unique ID.
func (s *ItemService) GetItem(ctx context.Context, id uuid.UUID) (*domain.Item, error) {
	return s.repo.GetByID(ctx, id)
}

// ListItems returns a list of items based on filter criteria.
func (s *ItemService) ListItems(ctx context.Context, filter domain.ItemFilter) ([]*domain.Item, error) {
	return s.repo.List(ctx, filter)
}
