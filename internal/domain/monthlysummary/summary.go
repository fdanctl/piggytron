package monthlysummary

import "time"

type ID string

type MonthlySummary struct {
	accountID ID
	month     Month
	moneyIn   int
	moneyOut  int
	createdAt time.Time
	updatedAt time.Time
}

func New(accountID ID, month Month, moneyIn, moneyOut int) (*MonthlySummary, error) {
	if moneyIn < 0 || moneyOut < 0 {
		return nil, ErrInvalidAmount
	}

	now := time.Now()
	return &MonthlySummary{
		accountID: accountID,
		month:     month,
		moneyIn:   moneyIn,
		moneyOut:  moneyOut,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func Rehydrate(
	accountID ID,
	month time.Time,
	moneyIn, moneyOut int,
	createdAt, updatedAt time.Time,
) *MonthlySummary {
	return &MonthlySummary{
		accountID: accountID,
		month:     NewMonth(month),
		moneyIn:   moneyIn,
		moneyOut:  moneyOut,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (s *MonthlySummary) AccountID() ID {
	return s.accountID
}

func (s *MonthlySummary) Month() Month {
	return s.month
}

func (s *MonthlySummary) MoneyIn() int {
	return s.moneyIn
}

func (s *MonthlySummary) MoneyOut() int {
	return s.moneyOut
}

func (s *MonthlySummary) CreatedAt() time.Time {
	return s.createdAt
}

func (s *MonthlySummary) UpdatedAt() time.Time {
	return s.updatedAt
}

func (s *MonthlySummary) SubMoneyIn(v int) error {
	if v < 0 {
		return ErrInvalidAmount
	}
	if s.moneyIn-v < 0 {
		return ErrInvalidAmount
	}
	s.moneyIn -= v
	s.updatedAt = time.Now()
	return nil
}

func (s *MonthlySummary) SubMoneyOut(v int) error {
	if v < 0 {
		return ErrInvalidAmount
	}
	if s.moneyOut-v < 0 {
		return ErrInvalidAmount
	}
	s.moneyOut -= v
	s.updatedAt = time.Now()
	return nil
}
