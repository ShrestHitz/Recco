package domain_test

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/recco-app/recco/internal/domain"
)

func TestShelfEntry_Validate(t *testing.T) {
	ratingValid := 8
	ratingZero := 0
	ratingTooHigh := 11
	now := time.Now().UTC()

	tests := []struct {
		name    string
		entry   domain.ShelfEntry
		wantErr bool
	}{
		{
			name: "valid queued entry",
			entry: domain.ShelfEntry{
				ID:     uuid.New(),
				UserID: uuid.New(),
				ItemID: uuid.New(),
				Status: domain.ShelfStatusQueued,
			},
			wantErr: false,
		},
		{
			name: "valid finished entry with rating and take",
			entry: domain.ShelfEntry{
				ID:         uuid.New(),
				UserID:     uuid.New(),
				ItemID:     uuid.New(),
				Status:     domain.ShelfStatusFinished,
				Rating:     &ratingValid,
				Take:       "Brilliant narrative pacing and profound ending.",
				FinishedAt: &now,
			},
			wantErr: false,
		},
		{
			name: "finished entry missing rating fails",
			entry: domain.ShelfEntry{
				ID:     uuid.New(),
				UserID: uuid.New(),
				ItemID: uuid.New(),
				Status: domain.ShelfStatusFinished,
				Rating: nil,
			},
			wantErr: true,
		},
		{
			name: "finished entry rating zero fails",
			entry: domain.ShelfEntry{
				ID:     uuid.New(),
				UserID: uuid.New(),
				ItemID: uuid.New(),
				Status: domain.ShelfStatusFinished,
				Rating: &ratingZero,
			},
			wantErr: true,
		},
		{
			name: "finished entry rating 11 fails",
			entry: domain.ShelfEntry{
				ID:     uuid.New(),
				UserID: uuid.New(),
				ItemID: uuid.New(),
				Status: domain.ShelfStatusFinished,
				Rating: &ratingTooHigh,
			},
			wantErr: true,
		},
		{
			name: "take exceeding 280 characters fails",
			entry: domain.ShelfEntry{
				ID:     uuid.New(),
				UserID: uuid.New(),
				ItemID: uuid.New(),
				Status: domain.ShelfStatusQueued,
				Take:   strings.Repeat("a", 281),
			},
			wantErr: true,
		},
		{
			name: "invalid status fails",
			entry: domain.ShelfEntry{
				ID:     uuid.New(),
				UserID: uuid.New(),
				ItemID: uuid.New(),
				Status: domain.ShelfStatus("unknown_status"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.entry.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ShelfEntry.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
