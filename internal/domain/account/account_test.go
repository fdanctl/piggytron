package account

import (
	"errors"
	"testing"
	"time"
)

func TestNewBank(t *testing.T) {
	tests := []struct {
		name     string
		id       ID
		userID   ID
		bankname string
		currency string
		isSaving bool
		wantErr  error
	}{
		{
			name:     "valid savings bank",
			id:       ID("420"),
			userID:   ID("420"),
			bankname: "savings",
			currency: "USD",
			isSaving: true,
			wantErr:  nil,
		},
		{
			name:     "valid normal bank",
			id:       ID("420"),
			userID:   ID("420"),
			bankname: "main",
			currency: "USD",
			isSaving: false,
			wantErr:  nil,
		},
		{
			name:     "empty name",
			id:       ID("420"),
			userID:   ID("420"),
			bankname: "",
			currency: "USD",
			isSaving: false,
			wantErr:  ErrInvalidName,
		},
		{
			name:     "name to long",
			id:       ID("420"),
			userID:   ID("420"),
			bankname: "H,g?)FTm&LHu0EyqQH*sAc_9h<tNrjJd?jucL$49ec!w)w2FU,)", // 51 chars
			currency: "USD",
			isSaving: false,
			wantErr:  ErrInvalidName,
		},
		{
			name:     "empty currency",
			id:       ID("420"),
			userID:   ID("420"),
			bankname: "main",
			currency: "",
			isSaving: false,
			wantErr:  ErrInvalidCurrency,
		},
		{
			name:     "long currency",
			id:       ID("420"),
			userID:   ID("420"),
			bankname: "main",
			currency: "UUUUUUUUUUU", // 11 chars
			isSaving: false,
			wantErr:  ErrInvalidCurrency,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewBank(tt.id, tt.userID, tt.bankname, tt.currency, tt.isSaving)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewBank() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewGoal(t *testing.T) {
	now := time.Now()
	yesterday := time.Date(
		now.Year(),
		now.Month(),
		now.Day()-1,
		now.Hour(),
		now.Minute(),
		now.Second(),
		now.Nanosecond(),
		now.Location(),
	)
	tests := []struct {
		name         string
		id           ID
		userID       ID
		goalname     string
		currency     string
		targetAmount int
		startDate    time.Time
		targetDate   *time.Time
		categoryID   ID
		wantErr      error
	}{
		{
			name:         "valid goal",
			id:           ID("420"),
			userID:       ID("420"),
			goalname:     "goal",
			currency:     "USD",
			targetAmount: 10000,
			startDate:    now,
			targetDate:   &now,
			categoryID:   ID("420"),
			wantErr:      nil,
		},
		{
			name:         "valid goal without target date",
			id:           ID("420"),
			userID:       ID("420"),
			goalname:     "goal",
			currency:     "USD",
			targetAmount: 10000,
			startDate:    now,
			targetDate:   nil,
			wantErr:      nil,
		},
		{
			name:         "empty name",
			id:           ID("420"),
			userID:       ID("420"),
			goalname:     "",
			currency:     "USD",
			targetAmount: 10000,
			startDate:    now,
			targetDate:   &now,
			categoryID:   ID("420"),
			wantErr:      ErrInvalidName,
		},
		{
			name:         "name to long",
			id:           ID("420"),
			userID:       ID("420"),
			goalname:     "H,g?)FTm&LHu0EyqQH*sAc_9h<tNrjJd?jucL$49ec!w)w2FU,)", // 51 chars
			currency:     "USD",
			targetAmount: 10000,
			startDate:    now,
			targetDate:   &now,
			categoryID:   ID("420"),
			wantErr:      ErrInvalidName,
		},
		{
			name:         "empty currency",
			id:           ID("420"),
			userID:       ID("420"),
			goalname:     "goal",
			currency:     "",
			targetAmount: 10000,
			startDate:    now,
			targetDate:   &now,
			categoryID:   ID("420"),
			wantErr:      ErrInvalidCurrency,
		},
		{
			name:         "long currency",
			id:           ID("420"),
			userID:       ID("420"),
			goalname:     "goal",
			currency:     "UUUUUUUUUUU", // 11 chars
			targetAmount: 10000,
			startDate:    now,
			targetDate:   &now,
			categoryID:   ID("420"),
			wantErr:      ErrInvalidCurrency,
		},
		{
			name:         "zero amount",
			id:           ID("420"),
			userID:       ID("420"),
			goalname:     "goal",
			currency:     "USD",
			targetAmount: 0,
			startDate:    now,
			targetDate:   &now,
			categoryID:   ID("420"),
			wantErr:      ErrNegativeNumber,
		},
		{
			name:         "negative amount",
			id:           ID("420"),
			userID:       ID("420"),
			goalname:     "goal",
			currency:     "USD",
			targetAmount: -10000,
			startDate:    now,
			targetDate:   &now,
			categoryID:   ID("420"),
			wantErr:      ErrNegativeNumber,
		},
		{
			name:         "start date after target",
			id:           ID("420"),
			userID:       ID("420"),
			goalname:     "goal",
			currency:     "USD",
			targetAmount: 10000,
			startDate:    now,
			targetDate:   &yesterday,
			categoryID:   ID("420"),
			wantErr:      ErrStartDateAfterTarget,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewGoal(
				tt.id,
				tt.userID,
				tt.goalname,
				tt.currency,
				tt.targetAmount,
				tt.startDate,
				tt.targetDate,
				tt.categoryID,
			)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewGoal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCanReceiveIncome(t *testing.T) {
	// Goal type
	goal, _ := NewGoal(ID("420"), ID("420"), "Goal", "USD", 5000, time.Now(), nil, ID("420"))
	if err := goal.CanReceiveIncome(); err == nil {
		t.Error("goals should not receive income directly")
	}

	// Savings bank
	saving, _ := NewBank(ID("420"), ID("420"), "Savings", "USD", true)
	if err := saving.CanReceiveIncome(); err == nil {
		t.Error("savings accounts should not receive income directly")
	}

	// Regular bank (should succeed)
	main, _ := NewBank(ID("420"), ID("420"), "Main", "USD", false)
	if err := main.CanReceiveIncome(); err != nil {
		t.Errorf("checking account should receive income, got: %v", err)
	}
}

func TestChangeName(t *testing.T) {
	goal, _ := NewGoal(ID("420"), ID("420"), "Goal", "USD", 5000, time.Now(), nil, ID("420"))
	upAt := goal.UpdatedAt()

	// valid
	newName := "New Name"
	if err := goal.ChangeName(newName); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// the name or the updated at did not change
	if goal.Name() != newName || goal.UpdatedAt().Compare(upAt) < 1 {
		t.Errorf("expected '%s', got '%s'", newName, goal.Name())
	}

	upAt = goal.UpdatedAt()
	// empty name
	if err := goal.ChangeName(""); err == nil {
		t.Error("expected error for empty name")
	}
	if goal.Name() != newName || goal.UpdatedAt().Compare(upAt) != 0 {
		t.Errorf("name should remain '%s', got '%s'", newName, goal.Name())
	}

	// name to long
	if err := goal.ChangeName("H,g?)FTm&LHu0EyqQH*sAc_9h<tNrjJd?jucL$49ec!w)w2FU,)"); err == nil {
		t.Error("expected error for name to long")
	}
	if goal.Name() != newName || goal.UpdatedAt().Compare(upAt) != 0 {
		t.Errorf("name should remain '%s', got '%s'", newName, goal.Name())
	}
}

func TestChangeTargetAmount(t *testing.T) {
	amount := 1000

	// Bank type
	bank, _ := NewBank(ID("420"), ID("420"), "Savings", "USD", true)
	upAt := bank.UpdatedAt()
	if err := bank.ChangeTargetAmount(amount); err == nil {
		t.Error("banks can't change target amount")
	}
	if bank.TargetAmount() != nil || bank.UpdatedAt().Compare(upAt) != 0 {
		t.Error("bank target amount should be 'nil'")
	}

	// Goal type
	goal, _ := NewGoal(ID("420"), ID("420"), "Goal", "USD", 5000, time.Now(), nil, ID("420"))
	upAt = goal.UpdatedAt()
	if err := goal.ChangeTargetAmount(amount); err != nil {
		t.Error("expected error for valid amount")
	}
	if *goal.TargetAmount() != amount || goal.UpdatedAt().Compare(upAt) < 1 {
		t.Errorf("expected target amount '%d', got '%d'", amount, *goal.TargetAmount())
	}

	upAt = goal.UpdatedAt()
	// zero amount
	if err := goal.ChangeTargetAmount(0); err == nil {
		t.Error("expected error for empty name")
	}
	if *goal.TargetAmount() != amount || goal.UpdatedAt().Compare(upAt) != 0 {
		t.Errorf("target amount should remain '%d', got '%d'", amount, *goal.TargetAmount())
	}

	// negative amount
	if err := goal.ChangeTargetAmount(-100); err == nil {
		t.Error("target amount should not receive income directly")
	}
	if *goal.TargetAmount() != amount || goal.UpdatedAt().Compare(upAt) != 0 {
		t.Errorf("target amount should remain '%d', got '%d'", amount, *goal.TargetAmount())
	}

	// start date after target
}

func TestChangeStartDate(t *testing.T) {
	today := time.Now()
	lastMonth := time.Date(
		today.Year(),
		today.Month()-1,
		today.Day(),
		today.Hour(),
		today.Minute(),
		today.Second(),
		today.Nanosecond(),
		today.Location(),
	)
	nextMonth := time.Date(
		today.Year(),
		today.Month()+1,
		today.Day(),
		today.Hour(),
		today.Minute(),
		today.Second(),
		today.Nanosecond(),
		today.Location(),
	)
	yesterday := time.Date(
		today.Year(),
		today.Month(),
		today.Day()-1,
		today.Hour(),
		today.Minute(),
		today.Second(),
		today.Nanosecond(),
		today.Location(),
	)
	nextTwoMonth := time.Date(
		today.Year(),
		today.Month()+2,
		today.Day(),
		today.Hour(),
		today.Minute(),
		today.Second(),
		today.Nanosecond(),
		today.Location(),
	)
	nextYear := time.Date(
		today.Year()+1,
		today.Month(),
		today.Day(),
		today.Hour(),
		today.Minute(),
		today.Second(),
		today.Nanosecond(),
		today.Location(),
	)

	// Bank type
	bank, _ := NewBank(ID("420"), ID("420"), "Savings", "USD", true)
	upAt := bank.UpdatedAt()
	if err := bank.ChangeStartDate(yesterday, nil); err == nil {
		t.Error("banks can't change start date")
	}
	if bank.TargetAmount() != nil || bank.UpdatedAt().Compare(upAt) != 0 {
		t.Error("bank start date should be 'nil'")
	}

	// normal change
	goal, _ := NewGoal(ID("420"), ID("420"), "Goal", "USD", 5000, yesterday, &nextMonth, ID("420"))
	upAt = goal.UpdatedAt()
	if err := goal.ChangeStartDate(today, nil); err != nil {
		t.Error("unexpected error for normal start date change")
	}
	if goal.StartDate().Compare(today) != 0 || goal.UpdatedAt().Compare(upAt) < 1 {
		t.Errorf(
			"expected start date '%s', got '%s'",
			today.String(),
			goal.StartDate().String(),
		)
	}
	upAt = goal.UpdatedAt()

	// valid change with min possible
	if err := goal.ChangeStartDate(yesterday, &nextMonth); err != nil {
		t.Error("unexpected error for start date change with min possible")
	}
	if goal.StartDate().Compare(yesterday) != 0 || goal.UpdatedAt().Compare(upAt) < 1 {
		t.Errorf(
			"expected start date '%s', got '%s'",
			yesterday.String(),
			goal.StartDate().String(),
		)
	}
	upAt = goal.UpdatedAt()

	// invalid change with min possible
	if err := goal.ChangeStartDate(today, &lastMonth); err == nil {
		t.Error("expected error for min possible are greater than start date")
	}
	if goal.StartDate().Compare(yesterday) != 0 || goal.UpdatedAt().Compare(upAt) != 0 {
		t.Errorf(
			"expected start date '%s', got '%s'",
			yesterday.String(),
			goal.StartDate().String(),
		)
	}

	// start date after target date
	if err := goal.ChangeStartDate(nextTwoMonth, &nextYear); err == nil {
		t.Error("expected error for start date after target date change")
	}
	if goal.StartDate().Compare(yesterday) != 0 || goal.UpdatedAt().Compare(upAt) != 0 {
		t.Errorf(
			"expected start date '%s', got '%s'",
			yesterday.String(),
			goal.StartDate().String(),
		)
	}
}

func TestChangeTargetDate(t *testing.T) {
	today := time.Now()
	nextMonth := time.Date(
		today.Year(),
		today.Month()+1,
		today.Day(),
		today.Hour(),
		today.Minute(),
		today.Second(),
		today.Nanosecond(),
		today.Location(),
	)
	yesterday := time.Date(
		today.Year(),
		today.Month(),
		today.Day()-1,
		today.Hour(),
		today.Minute(),
		today.Second(),
		today.Nanosecond(),
		today.Location(),
	)
	nextTwoMonth := time.Date(
		today.Year(),
		today.Month()+2,
		today.Day(),
		today.Hour(),
		today.Minute(),
		today.Second(),
		today.Nanosecond(),
		today.Location(),
	)

	// Bank type
	bank, _ := NewBank(ID("420"), ID("420"), "Savings", "USD", true)
	upAt := bank.UpdatedAt()
	if err := bank.ChangeTargetDate(&nextMonth); err == nil {
		t.Error("banks can't change target date")
	}
	if bank.TargetAmount() != nil || bank.UpdatedAt().Compare(upAt) != 0 {
		t.Error("bank start date should be 'nil'")
	}

	// normal change
	goal, _ := NewGoal(ID("420"), ID("420"), "Goal", "USD", 5000, today, &nextMonth, ID("420"))
	upAt = goal.UpdatedAt()
	if err := goal.ChangeTargetDate(&nextTwoMonth); err != nil {
		t.Error("unexpected error for normal start date change")
	}
	if goal.TargetDate().Compare(nextTwoMonth) != 0 || goal.UpdatedAt().Compare(upAt) < 1 {
		t.Errorf(
			"expected target date '%s', got '%s'",
			yesterday.String(),
			goal.StartDate().String(),
		)
	}
	upAt = goal.UpdatedAt()

	// start date after target date
	if err := goal.ChangeTargetDate(&yesterday); err == nil {
		t.Error("expected error for start date after target date change")
	}
	if goal.TargetDate().Compare(nextTwoMonth) != 0 || goal.UpdatedAt().Compare(upAt) != 0 {
		t.Errorf("expected target date '%s', got '%s'", nextTwoMonth, goal.StartDate().String())
	}
}

func TestChangeCategory(t *testing.T) {
	id := ID("2")
	bank, _ := NewBank(ID("420"), ID("420"), "Savings", "USD", true)
	upAt := bank.UpdatedAt()
	if err := bank.ChangeCategory(id); err == nil {
		t.Error("banks can't change category")
	}
	if bank.CategoryID() != nil || bank.UpdatedAt().Compare(upAt) != 0 {
		t.Error("bank category should be 'nil'")
	}

	goal, _ := NewGoal(ID("420"), ID("420"), "Goal", "USD", 5000, time.Now(), nil, ID("420"))
	upAt = goal.UpdatedAt()
	if err := goal.ChangeCategory(id); err != nil {
		t.Error("unexpected error for normal start date change")
	}
	if *goal.CategoryID() != id || goal.UpdatedAt().Compare(upAt) < 1 {
		t.Errorf("expected category '%s', got '%s'", id, goal.StartDate().String())
	}
}
