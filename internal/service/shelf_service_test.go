package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/recco-app/recco/internal/domain"
	"github.com/recco-app/recco/internal/repository/memory"
	"github.com/recco-app/recco/internal/service"
)

func TestShelfService_QueueAndFinishWorkflow(t *testing.T) {
	ctx := context.Background()
	itemRepo := memory.NewItemRepo()
	shelfRepo := memory.NewShelfRepo(itemRepo)

	itemService := service.NewItemService(itemRepo)
	shelfService := service.NewShelfService(shelfRepo, itemRepo)

	// Create test item
	item, err := itemService.CreateItem(ctx, service.CreateItemInput{
		Type:        domain.ItemTypeBook,
		Title:       "Designing Data-Intensive Applications",
		Creator:     "Martin Kleppmann",
		ExternalURL: "https://openlibrary.org/books/OL123",
	})
	if err != nil {
		t.Fatalf("failed to create item: %v", err)
	}

	userID := uuid.New()

	// 1. Add to Queue
	entry, err := shelfService.AddToQueue(ctx, service.AddToQueueInput{
		UserID: userID,
		ItemID: item.ID,
	})
	if err != nil {
		t.Fatalf("failed adding to queue: %v", err)
	}
	if entry.Status != domain.ShelfStatusQueued {
		t.Errorf("expected status queued, got %s", entry.Status)
	}

	// 2. Verify queue listing
	queue, err := shelfService.ListQueue(ctx, userID, 10, 0)
	if err != nil {
		t.Fatalf("failed listing queue: %v", err)
	}
	if len(queue) != 1 {
		t.Fatalf("expected queue length 1, got %d", len(queue))
	}

	// 3. Finish item (Move to Shelf)
	finished, err := shelfService.FinishItem(ctx, service.FinishItemInput{
		UserID: userID,
		ItemID: item.ID,
		Rating: 10,
		Take:   "The definitive guide to distributed systems storage and consensus.",
	})
	if err != nil {
		t.Fatalf("failed finishing item: %v", err)
	}
	if finished.Status != domain.ShelfStatusFinished {
		t.Errorf("expected status finished, got %s", finished.Status)
	}
	if *finished.Rating != 10 {
		t.Errorf("expected rating 10, got %d", *finished.Rating)
	}

	// 4. Verify queue is now empty (Inbox/Queue Zero reached)
	emptyQueue, err := shelfService.ListQueue(ctx, userID, 10, 0)
	if err != nil {
		t.Fatalf("failed listing queue: %v", err)
	}
	if len(emptyQueue) != 0 {
		t.Errorf("expected empty queue, got %d items", len(emptyQueue))
	}

	// 5. Verify shelf now has 1 rated entry
	shelf, err := shelfService.ListShelf(ctx, userID, 10, 0)
	if err != nil {
		t.Fatalf("failed listing shelf: %v", err)
	}
	if len(shelf) != 1 {
		t.Fatalf("expected shelf length 1, got %d", len(shelf))
	}
	if shelf[0].Item == nil || shelf[0].Item.Title != item.Title {
		t.Errorf("expected hydrated item title %s, got %v", item.Title, shelf[0].Item)
	}
}
