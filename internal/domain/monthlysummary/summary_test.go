package monthlysummary

import (
	"errors"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		accountID ID
		month     Month
		moneyIn   int
		moneyOut  int
		wantErr   error
	}{
		{
			name:      "valid summary",
			accountID: ID("420"),
			month:     NewMonth(time.Now()),
			moneyIn:   420,
			moneyOut:  69,
			wantErr:   nil,
		},
		{
			name:      "negative moneyIn",
			accountID: ID("420"),
			month:     NewMonth(time.Now()),
			moneyIn:   -420,
			moneyOut:  69,
			wantErr:   ErrInvalidAmount,
		},
		{
			name:      "negative moneyOut",
			accountID: ID("420"),
			month:     NewMonth(time.Now()),
			moneyIn:   420,
			moneyOut:  -69,
			wantErr:   ErrInvalidAmount,
		},
		{
			name:      "zero moneyIn",
			accountID: ID("420"),
			month:     NewMonth(time.Now()),
			moneyIn:   0,
			moneyOut:  69,
			wantErr:   nil,
		},
		{
			name:      "zero moneyOut",
			accountID: ID("420"),
			month:     NewMonth(time.Now()),
			moneyIn:   420,
			moneyOut:  0,
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.accountID, tt.month, tt.moneyIn, tt.moneyOut)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubMoneyIn(t *testing.T) {
	tests := []struct {
		name    string
		initial int
		sub     int
		want    int
		wantErr error
	}{
		{
			name:    "sub positive",
			initial: 100,
			sub:     30,
			want:    70,
			wantErr: nil,
		},
		{
			name:    "sub zero",
			initial: 100,
			sub:     0,
			want:    100,
			wantErr: nil,
		},
		{
			name:    "sub negative",
			initial: 100,
			sub:     -10,
			want:    100,
			wantErr: ErrInvalidAmount,
		},
		{
			name:    "sub more than available",
			initial: 50,
			sub:     100,
			want:    50,
			wantErr: ErrInvalidAmount,
		},
		{
			name:    "sub exact amount",
			initial: 100,
			sub:     100,
			want:    0,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := New("acc1", NewMonth(time.Now()), tt.initial, 0)

			err := s.SubMoneyIn(tt.sub)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SubMoneyIn() error = %v, wantErr %v", err, tt.wantErr)
			}
			if s.MoneyIn() != tt.want {
				t.Errorf("MoneyIn() = %d, want %d", s.MoneyIn(), tt.want)
			}
		})
	}
}

func TestSubMoneyOut(t *testing.T) {
	tests := []struct {
		name    string
		initial int
		sub     int
		want    int
		wantErr error
	}{
		{
			name:    "sub positive",
			initial: 200,
			sub:     60,
			want:    140,
			wantErr: nil,
		},
		{
			name:    "sub zero",
			initial: 200,
			sub:     0,
			want:    200,
			wantErr: nil,
		},
		{
			name:    "sub negative",
			initial: 200,
			sub:     -30,
			want:    200,
			wantErr: ErrInvalidAmount,
		},
		{
			name:    "sub more than available",
			initial: 40,
			sub:     90,
			want:    40,
			wantErr: ErrInvalidAmount,
		},
		{
			name:    "sub exact amount",
			initial: 200,
			sub:     200,
			want:    0,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := New("acc1", NewMonth(time.Now()), 0, tt.initial)
			err := s.SubMoneyOut(tt.sub)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SubMoneyOut() error = %v, wantErr %v", err, tt.wantErr)
			}
			if s.MoneyOut() != tt.want {
				t.Errorf("MoneyOut() = %d, want %d", s.MoneyOut(), tt.want)
			}
		})
	}
}
