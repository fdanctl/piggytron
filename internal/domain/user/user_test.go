package user

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		id       ID
		username string
		password string
		wantErr  error
	}{
		{
			name:     "valid user",
			id:       ID("420"),
			username: "gopher",
			password: "password",
			wantErr:  nil,
		},
		{
			name:     "empty name",
			id:       ID("420"),
			username: "",
			password: "password",
			wantErr:  ErrInvalidName,
		},
		{
			name:     "name to long",
			id:       ID("420"),
			username: "H,g?)FTm&LHu0EyqQH*sAc_9h<tNrjJd?jucL$49ec!w)w2FU,)", // 51 chars
			password: "password",
			wantErr:  ErrInvalidName,
		},
		{
			name:     "empty password",
			id:       ID("420"),
			username: "gopher",
			password: "",
			wantErr:  ErrInvalidPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.id, tt.username, tt.password)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChangeName(t *testing.T) {
	u, _ := New(ID("420"), "gopher", "password")
	upAt := u.UpdatedAt()

	newName := "New Name"
	if err := u.ChangeName(newName); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if u.Name() != newName || u.UpdatedAt().Compare(upAt) < 1 {
		t.Errorf("expected '%s', got '%s'", newName, u.Name())
	}
	upAt = u.UpdatedAt()

	// empty name
	if err := u.ChangeName(""); err == nil {
		t.Error("expected error for empty name")
	}
	if u.Name() != newName || u.UpdatedAt().Compare(upAt) != 0 {
		t.Errorf("name should remain '%s', got '%s'", newName, u.Name())
	}

	// name to long
	if err := u.ChangeName("H,g?)FTm&LHu0EyqQH*sAc_9h<tNrjJd?jucL$49ec!w)w2FU,)"); err == nil {
		t.Error("expected error for name to long")
	}
	if u.Name() != newName || u.UpdatedAt().Compare(upAt) != 0 {
		t.Errorf("name should remain '%s', got '%s'", newName, u.Name())
	}
}
