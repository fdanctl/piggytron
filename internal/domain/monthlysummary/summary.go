// Package monthlysummary defines the per-account monthly rollup
// (money_in / money_out) used to compute balances without scanning the whole
// ledger. Its month is always the first day of the month.
package monthlysummary

import (
	"time"
)

// ID is a monthly summary identifier (account id + month).
type ID string

// MonthlySummary rolls up money in and money out for one account in one
// month, avoiding full ledger scans when computing balances.
type MonthlySummary struct {
	accountID ID
	month     Month
	moneyIn   int
	moneyOut  int
	createdAt time.Time
	updatedAt time.Time
}

// New builds a validated summary with non-negative in/out totals.
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

// Rehydrate rebuilds a MonthlySummary from persistence; the month is
// normalized to the first of the month.
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

// AccountID returns the summarized account.
func (s *MonthlySummary) AccountID() ID {
	return s.accountID
}

// Month returns the summarized month.
func (s *MonthlySummary) Month() Month {
	return s.month
}

// MoneyIn returns the total money received in cents.
func (s *MonthlySummary) MoneyIn() int {
	return s.moneyIn
}

// MoneyOut returns the total money spent in cents.
func (s *MonthlySummary) MoneyOut() int {
	return s.moneyOut
}

// CreatedAt returns when the summary row was created.
func (s *MonthlySummary) CreatedAt() time.Time {
	return s.createdAt
}

// UpdatedAt returns when the summary row was last updated.
func (s *MonthlySummary) UpdatedAt() time.Time {
	return s.updatedAt
}

// AddMoneyIn adjusts the money-in total by v (v may be negative to undo an
// entry), rejecting a negative result.
func (s *MonthlySummary) AddMoneyIn(v int) error {
	if s.moneyIn+v < 0 {
		return ErrInvalidAmount
	}
	s.moneyIn += v
	s.updatedAt = time.Now()
	return nil
}

// AddMoneyOut adjusts the money-out total by v (v may be negative to undo an
// entry), rejecting a negative result.
func (s *MonthlySummary) AddMoneyOut(v int) error {
	if s.moneyOut+v < 0 {
		return ErrInvalidAmount
	}
	s.moneyOut += v
	s.updatedAt = time.Now()
	return nil
}
