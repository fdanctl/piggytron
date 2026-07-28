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

func TestAddMoneyIn(t *testing.T) {
	tests := []struct {
		name        string
		initialIn   int
		delta       int
		wantErr     error
		wantMoneyIn int
	}{
		{
			name:        "add positive amount",
			initialIn:   1000,
			delta:       500,
			wantErr:     nil,
			wantMoneyIn: 1500,
		},
		{
			name:        "add zero",
			initialIn:   1000,
			delta:       0,
			wantErr:     nil,
			wantMoneyIn: 1000,
		},
		{
			name:        "subtract within bounds",
			initialIn:   1000,
			delta:       -500,
			wantErr:     nil,
			wantMoneyIn: 500,
		},
		{
			name:        "subtract to zero",
			initialIn:   1000,
			delta:       -1000,
			wantErr:     nil,
			wantMoneyIn: 0,
		},
		{
			name:        "subtract below zero rejected",
			initialIn:   1000,
			delta:       -1001,
			wantErr:     ErrInvalidAmount,
			wantMoneyIn: 1000,
		},
		{
			name:        "add to zero",
			initialIn:   0,
			delta:       500,
			wantErr:     nil,
			wantMoneyIn: 500,
		},
		{
			name:        "subtract from zero rejected",
			initialIn:   0,
			delta:       -1,
			wantErr:     ErrInvalidAmount,
			wantMoneyIn: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms, err := New(ID("acc1"), NewMonth(time.Now()), tt.initialIn, 0)
			if err != nil {
				t.Fatalf("New() setup error: %v", err)
			}

			err = ms.AddMoneyIn(tt.delta)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("AddMoneyIn(%d) error = %v, wantErr %v", tt.delta, err, tt.wantErr)
			}
			if ms.MoneyIn() != tt.wantMoneyIn {
				t.Errorf("AddMoneyIn(%d) MoneyIn = %d, want %d", tt.delta, ms.MoneyIn(), tt.wantMoneyIn)
			}
		})
	}
}

func TestAddMoneyOut(t *testing.T) {
	tests := []struct {
		name         string
		initialOut   int
		delta        int
		wantErr      error
		wantMoneyOut int
	}{
		{
			name:         "add positive amount",
			initialOut:   1000,
			delta:        500,
			wantErr:      nil,
			wantMoneyOut: 1500,
		},
		{
			name:         "add zero",
			initialOut:   1000,
			delta:        0,
			wantErr:      nil,
			wantMoneyOut: 1000,
		},
		{
			name:         "subtract within bounds",
			initialOut:   1000,
			delta:        -500,
			wantErr:      nil,
			wantMoneyOut: 500,
		},
		{
			name:         "subtract to zero",
			initialOut:   1000,
			delta:        -1000,
			wantErr:      nil,
			wantMoneyOut: 0,
		},
		{
			name:         "subtract below zero rejected",
			initialOut:   1000,
			delta:        -1001,
			wantErr:      ErrInvalidAmount,
			wantMoneyOut: 1000,
		},
		{
			name:         "add to zero",
			initialOut:   0,
			delta:        500,
			wantErr:      nil,
			wantMoneyOut: 500,
		},
		{
			name:         "subtract from zero rejected",
			initialOut:   0,
			delta:        -1,
			wantErr:      ErrInvalidAmount,
			wantMoneyOut: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms, err := New(ID("acc1"), NewMonth(time.Now()), 0, tt.initialOut)
			if err != nil {
				t.Fatalf("New() setup error: %v", err)
			}

			err = ms.AddMoneyOut(tt.delta)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("AddMoneyOut(%d) error = %v, wantErr %d", tt.delta, err, tt.wantErr)
			}
			if ms.MoneyOut() != tt.wantMoneyOut {
				t.Errorf("AddMoneyOut(%d) MoneyOut = %d, want %d", tt.delta, ms.MoneyOut(), tt.wantMoneyOut)
			}
		})
	}
}
