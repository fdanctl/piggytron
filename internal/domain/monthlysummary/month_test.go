package monthlysummary

import (
	"errors"
	"testing"
)

func TestParseMonth(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "valid month",
			input:   "2026-05",
			wantErr: nil,
		},
		{
			name:    "short year",
			input:   "26-05",
			wantErr: ErrInvalidMonth,
		},
		{
			name:    "no separator",
			input:   "202605",
			wantErr: ErrInvalidMonth,
		},
		{
			name:    "full date not allowed",
			input:   "2026-05-22",
			wantErr: ErrInvalidMonth,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseMonth(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ParseMonth(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}
