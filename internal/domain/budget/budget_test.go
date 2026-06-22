package budget

import (
	"errors"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		id         ID
		userID     ID
		categoryID ID
		month      time.Time
		amount     int
		wantErr    error
	}{
		{
			name:       "valid budget",
			id:         ID("420"),
			userID:     ID("420"),
			categoryID: ID("420"),
			month:      time.Now(),
			amount:     100,
			wantErr:    nil,
		},
		{
			name:       "valid budget zero amount",
			id:         ID("420"),
			userID:     ID("420"),
			categoryID: ID("420"),
			month:      time.Now(),
			amount:     0,
			wantErr:    nil,
		},
		{
			name:       "negative amount",
			id:         ID("420"),
			userID:     ID("420"),
			categoryID: ID("420"),
			month:      time.Now(),
			amount:     -100,
			wantErr:    ErrInvalidAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.id, tt.userID, tt.categoryID, tt.month, tt.amount)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
