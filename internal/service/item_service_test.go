package service_test

import (
	"context"
	"testing"

	"github.com/recco-app/recco/internal/domain"
	"github.com/recco-app/recco/internal/repository/memory"
	"github.com/recco-app/recco/internal/service"
)

func TestItemService_CreateItem(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewItemRepo()
	svc := service.NewItemService(repo)

	// 1. Create first item
	input := service.CreateItemInput{
		Type:        domain.ItemTypeYouTube,
		Title:       "How to Design Scalable Systems",
		Creator:     "TechChannel",
		ExternalURL: "https://youtube.com/watch?v=unique123",
		SourceID:    "yt:unique123",
	}

	created, err := svc.CreateItem(ctx, input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if created.Title != input.Title {
		t.Errorf("expected title %s, got %s", input.Title, created.Title)
	}

	// 2. Attempt duplicate creation with same URL (should deduplicate and return existing)
	duplicateInput := service.CreateItemInput{
		Type:        domain.ItemTypeYouTube,
		Title:       "Duplicate Name",
		ExternalURL: "https://youtube.com/watch?v=unique123",
	}

	deduped, err := svc.CreateItem(ctx, duplicateInput)
	if err != nil {
		t.Fatalf("expected no error on duplicate, got %v", err)
	}
	if deduped.ID != created.ID {
		t.Errorf("expected deduped ID %s, got %s", created.ID, deduped.ID)
	}
}
