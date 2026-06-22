package transaction

import (
	"errors"
	"testing"
	"time"
)

func TestNewIncome(t *testing.T) {
	tests := []struct {
		name        string
		id          ID
		userID      ID
		toAccID     ID
		iCategoryID ID
		amount      int
		description string
		date        time.Time
		wantErr     error
	}{
		{
			name:        "valid income",
			id:          ID("420"),
			userID:      ID("420"),
			toAccID:     ID("420"),
			iCategoryID: ID("420"),
			amount:      10000,
			description: "a good description",
			date:        time.Now(),
			wantErr:     nil,
		},
		{
			name:        "empty description",
			id:          ID("420"),
			userID:      ID("420"),
			toAccID:     ID("420"),
			iCategoryID: ID("420"),
			amount:      10000,
			description: "",
			date:        time.Now(),
			wantErr:     ErrInvalidDescription,
		},
		{
			name:        "zero amount",
			id:          ID("420"),
			userID:      ID("420"),
			toAccID:     ID("420"),
			iCategoryID: ID("420"),
			amount:      0,
			description: "a good description",
			date:        time.Now(),
			wantErr:     ErrInvalidAmount,
		},
		{
			name:        "negative amount",
			id:          ID("420"),
			userID:      ID("420"),
			toAccID:     ID("420"),
			iCategoryID: ID("420"),
			amount:      -10000,
			description: "a good description",
			date:        time.Now(),
			wantErr:     ErrInvalidAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewIncome(
				tt.id,
				tt.userID,
				tt.toAccID,
				tt.iCategoryID,
				tt.amount,
				tt.description,
				tt.date,
			)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewIncome() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewExpense(t *testing.T) {
	tests := []struct {
		name          string
		id            ID
		userID        ID
		fromAccID     ID
		eCategoryID   ID
		amount        int
		description   string
		date          time.Time
		sourceBalance int
		wantErr       error
	}{
		{
			name:          "valid expense",
			id:            ID("420"),
			userID:        ID("420"),
			fromAccID:     ID("420"),
			eCategoryID:   ID("420"),
			amount:        10000,
			description:   "a good description",
			date:          time.Now(),
			sourceBalance: 20000,
			wantErr:       nil,
		},
		{
			name:          "empty description",
			id:            ID("420"),
			userID:        ID("420"),
			fromAccID:     ID("420"),
			eCategoryID:   ID("420"),
			amount:        10000,
			description:   "",
			date:          time.Now(),
			sourceBalance: 20000,
			wantErr:       ErrInvalidDescription,
		},
		{
			name:          "zero amount",
			id:            ID("420"),
			userID:        ID("420"),
			fromAccID:     ID("420"),
			eCategoryID:   ID("420"),
			amount:        0,
			description:   "a good description",
			date:          time.Now(),
			sourceBalance: 20000,
			wantErr:       ErrInvalidAmount,
		},
		{
			name:          "negative amount",
			id:            ID("420"),
			userID:        ID("420"),
			fromAccID:     ID("420"),
			eCategoryID:   ID("420"),
			amount:        -10000,
			description:   "a good description",
			date:          time.Now(),
			sourceBalance: 20000,
			wantErr:       ErrInvalidAmount,
		},
		{
			name:          "zero balance",
			id:            ID("420"),
			userID:        ID("420"),
			fromAccID:     ID("420"),
			eCategoryID:   ID("420"),
			amount:        10000,
			description:   "a good description",
			date:          time.Now(),
			sourceBalance: 10000,
			wantErr:       nil,
		},
		{
			name:          "negative balance",
			id:            ID("420"),
			userID:        ID("420"),
			fromAccID:     ID("420"),
			eCategoryID:   ID("420"),
			amount:        20000,
			description:   "a good description",
			date:          time.Now(),
			sourceBalance: 10000,
			wantErr:       ErrNegativeBalance,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewExpense(
				tt.id,
				tt.userID,
				tt.fromAccID,
				tt.eCategoryID,
				tt.amount,
				tt.description,
				tt.date,
				tt.sourceBalance,
			)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewExpense() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewTransfer(t *testing.T) {
	expenseCatID := ID("000")
	expenseCatID1 := ID("111")
	toCatID := ID("111")
	tests := []struct {
		name        string
		id          ID
		userID      ID
		toAccID     ID
		fromAccID   ID
		eCategoryID *ID
		amount      int
		description string
		date        time.Time

		sourceBalance    int
		toAccCatID       *ID
		toAccountCatType string
		isToAccSavings   bool
		wantErr          error
	}{
		{
			name:             "valid transfer",
			id:               ID("420"),
			userID:           ID("420"),
			toAccID:          ID("444"),
			fromAccID:        ID("420"),
			eCategoryID:      &expenseCatID,
			amount:           10000,
			description:      "a good description",
			date:             time.Now(),
			sourceBalance:    20000,
			toAccCatID:       nil,
			toAccountCatType: "needs",
			isToAccSavings:   false,
			wantErr:          nil,
		},
		{
			name:             "valid transfer without category",
			id:               ID("420"),
			userID:           ID("420"),
			toAccID:          ID("444"),
			fromAccID:        ID("420"),
			eCategoryID:      nil,
			amount:           10000,
			description:      "a good description",
			date:             time.Now(),
			sourceBalance:    20000,
			toAccCatID:       nil,
			toAccountCatType: "needs",
			isToAccSavings:   false,
			wantErr:          nil,
		},
		{
			name:             "empty description",
			id:               ID("420"),
			userID:           ID("420"),
			toAccID:          ID("444"),
			fromAccID:        ID("420"),
			eCategoryID:      nil,
			amount:           10000,
			description:      "",
			date:             time.Now(),
			sourceBalance:    20000,
			toAccCatID:       nil,
			toAccountCatType: "needs",
			isToAccSavings:   false,
			wantErr:          ErrInvalidDescription,
		},
		{
			name:             "zero amount",
			id:               ID("420"),
			userID:           ID("420"),
			toAccID:          ID("444"),
			fromAccID:        ID("420"),
			eCategoryID:      nil,
			amount:           0,
			description:      "a good description",
			date:             time.Now(),
			sourceBalance:    20000,
			toAccCatID:       nil,
			toAccountCatType: "needs",
			isToAccSavings:   false,
			wantErr:          ErrInvalidAmount,
		},
		{
			name:             "negative amount",
			id:               ID("420"),
			userID:           ID("420"),
			toAccID:          ID("444"),
			fromAccID:        ID("420"),
			eCategoryID:      nil,
			amount:           -10000,
			description:      "a good description",
			date:             time.Now(),
			sourceBalance:    20000,
			toAccCatID:       nil,
			toAccountCatType: "needs",
			isToAccSavings:   false,
			wantErr:          ErrInvalidAmount,
		},
		{
			name:             "zero balance",
			id:               ID("420"),
			userID:           ID("420"),
			toAccID:          ID("444"),
			fromAccID:        ID("420"),
			eCategoryID:      nil,
			amount:           10000,
			description:      "a good description",
			date:             time.Now(),
			sourceBalance:    10000,
			toAccCatID:       nil,
			toAccountCatType: "needs",
			isToAccSavings:   false,
			wantErr:          nil,
		},
		{
			name:             "negative balance",
			id:               ID("420"),
			userID:           ID("420"),
			toAccID:          ID("444"),
			fromAccID:        ID("420"),
			eCategoryID:      nil,
			amount:           20000,
			description:      "a good description",
			date:             time.Now(),
			sourceBalance:    10000,
			toAccCatID:       nil,
			toAccountCatType: "needs",
			isToAccSavings:   false,
			wantErr:          ErrNegativeBalance,
		},
		{
			name:             "same account transfer",
			id:               ID("420"),
			userID:           ID("420"),
			toAccID:          ID("420"),
			fromAccID:        ID("420"),
			eCategoryID:      nil,
			amount:           10000,
			description:      "a good description",
			date:             time.Now(),
			sourceBalance:    20000,
			toAccCatID:       nil,
			toAccountCatType: "needs",
			isToAccSavings:   false,
			wantErr:          ErrSameAccountTransfer,
		},
		{
			name:             "goal transfer without category",
			id:               ID("420"),
			userID:           ID("420"),
			toAccID:          ID("444"),
			fromAccID:        ID("420"),
			eCategoryID:      nil,
			amount:           10000,
			description:      "a good description",
			date:             time.Now(),
			sourceBalance:    20000,
			toAccCatID:       &toCatID,
			toAccountCatType: "needs",
			isToAccSavings:   false,
			wantErr:          ErrGoalCategory,
		},
		{
			name:             "goal transfer with wrong category",
			id:               ID("420"),
			userID:           ID("420"),
			toAccID:          ID("444"),
			fromAccID:        ID("420"),
			eCategoryID:      &expenseCatID,
			amount:           10000,
			description:      "a good description",
			date:             time.Now(),
			sourceBalance:    20000,
			toAccCatID:       &toCatID,
			toAccountCatType: "needs",
			isToAccSavings:   false,
			wantErr:          ErrGoalCategory,
		},
		{
			name:             "goal transfer with correct category",
			id:               ID("420"),
			userID:           ID("420"),
			toAccID:          ID("444"),
			fromAccID:        ID("420"),
			eCategoryID:      &expenseCatID1,
			amount:           10000,
			description:      "a good description",
			date:             time.Now(),
			sourceBalance:    20000,
			toAccCatID:       &toCatID,
			toAccountCatType: "needs",
			isToAccSavings:   false,
			wantErr:          nil,
		},
		{
			name:             "to savings account without category",
			id:               ID("420"),
			userID:           ID("420"),
			toAccID:          ID("444"),
			fromAccID:        ID("420"),
			eCategoryID:      nil,
			amount:           10000,
			description:      "a good description",
			date:             time.Now(),
			sourceBalance:    20000,
			toAccCatID:       nil,
			toAccountCatType: "",
			isToAccSavings:   true,
			wantErr:          ErrNotSavingsCategory,
		},
		{
			name:             "to savings account not savings category",
			id:               ID("420"),
			userID:           ID("420"),
			toAccID:          ID("444"),
			fromAccID:        ID("420"),
			eCategoryID:      &expenseCatID,
			amount:           10000,
			description:      "a good description",
			date:             time.Now(),
			sourceBalance:    20000,
			toAccCatID:       nil,
			toAccountCatType: "needs",
			isToAccSavings:   true,
			wantErr:          ErrNotSavingsCategory,
		},
		{
			name:             "valid transfer to savings account",
			id:               ID("420"),
			userID:           ID("420"),
			toAccID:          ID("444"),
			fromAccID:        ID("420"),
			eCategoryID:      &expenseCatID,
			amount:           10000,
			description:      "a good description",
			date:             time.Now(),
			sourceBalance:    20000,
			toAccCatID:       nil,
			toAccountCatType: "savings",
			isToAccSavings:   true,
			wantErr:          nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTransfer(
				tt.id,
				tt.userID,
				tt.fromAccID,
				tt.toAccID,
				tt.eCategoryID,
				tt.amount,
				tt.description,
				tt.date,
				tt.sourceBalance,
				tt.toAccCatID,
				tt.toAccountCatType,
				tt.isToAccSavings,
			)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewTransfer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
