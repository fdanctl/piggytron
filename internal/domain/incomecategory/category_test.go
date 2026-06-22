package incomecategory

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name         string
		id           ID
		userID       ID
		categoryName string
		wantErr      error
	}{
		{
			name:         "valid category",
			id:           ID("420"),
			userID:       ID("420"),
			categoryName: "Rent",
			wantErr:      nil,
		},
		{
			name:         "empty name",
			id:           ID("420"),
			userID:       ID("420"),
			categoryName: "",
			wantErr:      ErrInvalidName,
		},
		{
			name:         "long name",
			id:           ID("420"),
			userID:       ID("420"),
			categoryName: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", // 31 chars
			wantErr:      ErrInvalidName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.id, tt.userID, tt.categoryName)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChangeName(t *testing.T) {
	ec, _ := New(ID("420"), ID("420"), "Rent")
	upAt := ec.UpdatedAt()

	// valid name
	name := "Transportation"
	if err := ec.ChangeName(name); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if ec.Name() != name || ec.UpdatedAt().Compare(upAt) < 1 {
		t.Errorf("expected '%s', got '%s'", name, ec.Name())
	}
	upAt = ec.UpdatedAt()

	// empty name
	if err := ec.ChangeName(""); err == nil {
		t.Error("expected error for empty name")
	}
	if ec.Name() != name || ec.UpdatedAt().Compare(upAt) != 0 {
		t.Errorf("name should remain '%s', got '%s'", name, ec.Name())
	}

	// name to long
	if err := ec.ChangeName("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"); err == nil {
		t.Error("expected error for name to long")
	}
	if ec.Name() != name || ec.UpdatedAt().Compare(upAt) != 0 {
		t.Errorf("name should remain '%s', got '%s'", name, ec.Name())
	}
}
