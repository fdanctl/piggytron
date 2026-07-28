package budget

import (
	"errors"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		categoryID ID
		month      Month
		amount     int
		wantErr    error
	}{
		{
			name:       "valid budget",
			categoryID: ID("420"),
			month:      NewMonth(time.Now()),
			amount:     100,
			wantErr:    nil,
		},
		{
			name:       "valid budget zero amount",
			categoryID: ID("420"),
			month:      NewMonth(time.Now()),
			amount:     0,
			wantErr:    nil,
		},
		{
			name:       "negative amount",
			categoryID: ID("420"),
			month:      NewMonth(time.Now()),
			amount:     -100,
			wantErr:    ErrInvalidAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.categoryID, tt.month, tt.amount)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChangeAmount(t *testing.T) {
	amount := 1000
	b, _ := New(ID("420"), NewMonth(time.Now()), 100)
	upAt := b.UpdatedAt()

	if err := b.ChangeAmount(amount); err != nil {
		t.Error("unexpected error for valid amount")
	}
	if b.Amount() != amount || b.UpdatedAt().Compare(upAt) < 1 {
		t.Errorf("expected target amount '%d', got '%d'", amount, b.Amount())
	}

	upAt = b.UpdatedAt()
	// zero amount
	if err := b.ChangeAmount(0); err != nil {
		t.Error("unexpected error for zero amount")
	}
	if b.Amount() != 0 || b.UpdatedAt().Compare(upAt) < 0 {
		t.Errorf("expected '%d', got '%d'", amount, b.Amount())
	}
	upAt = b.UpdatedAt()

	// negative amount
	if err := b.ChangeAmount(-100); err == nil {
		t.Error("expected error for negative amount")
	}
	if b.Amount() != 0 || b.UpdatedAt().Compare(upAt) != 0 {
		t.Errorf("target amount should remain '%d', got '%d'", 0, b.Amount())
	}
}
