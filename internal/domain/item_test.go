package domain_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/recco-app/recco/internal/domain"
)

func TestItem_Validate(t *testing.T) {
	tests := []struct {
		name    string
		item    domain.Item
		wantErr bool
	}{
		{
			name: "valid youtube item",
			item: domain.Item{
				ID:          uuid.New(),
				Type:        domain.ItemTypeYouTube,
				Title:       "Clean Code in Go",
				ExternalURL: "https://www.youtube.com/watch?v=123",
			},
			wantErr: false,
		},
		{
			name: "valid movie item",
			item: domain.Item{
				ID:          uuid.New(),
				Type:        domain.ItemTypeMovie,
				Title:       "Interstellar",
				ExternalURL: "https://www.themoviedb.org/movie/157336",
			},
			wantErr: false,
		},
		{
			name: "invalid item type",
			item: domain.Item{
				Type:        domain.ItemType("unsupported_media"),
				Title:       "Random Video",
				ExternalURL: "https://example.com/video",
			},
			wantErr: true,
		},
		{
			name: "missing title",
			item: domain.Item{
				Type:        domain.ItemTypeBook,
				Title:       "   ",
				ExternalURL: "https://openlibrary.org/works/OL123",
			},
			wantErr: true,
		},
		{
			name: "missing external url",
			item: domain.Item{
				Type:        domain.ItemTypeAnime,
				Title:       "Steins;Gate",
				ExternalURL: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.item.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Item.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
