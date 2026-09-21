// Package budgethold stores the amount on hold for the next month
package budgethold

import "time"

// ID is a budget identifier.
type ID string

// BudgetHold is the hold amount for the next month
type BudgetHold struct {
	userID ID
	month  Month
	amount int

	createdAt time.Time
	updatedAt time.Time
}

// New builds a validated budgethold;
func New(
	uid ID,
	month Month,
	amount int,
) (*BudgetHold, error) {
	if amount < 0 {
		return nil, ErrInvalidAmount
	}

	now := time.Now()

	return &BudgetHold{
		userID: uid,
		month:  month,
		amount: amount,

		createdAt: now,
		updatedAt: now,
	}, nil
}

// Rehydrate rebuilds a BudgetHold from persistence; the month is normalized to
// the first of the month.
func Rehydrate(
	uid ID,
	month time.Time,
	amount int,
	createdAt, updatedAt time.Time,
) *BudgetHold {
	return &BudgetHold{
		userID: uid,
		month:  NewMonth(month),
		amount: amount,

		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// UserID returns the user id.
func (b *BudgetHold) UserID() ID {
	return b.userID
}

// Month returns the budgeted month.
func (b *BudgetHold) Month() Month {
	return b.month
}

// Amount returns the budget limit in cents.
func (b *BudgetHold) Amount() int {
	return b.amount
}

// CreatedAt returns when the budget was created.
func (b *BudgetHold) CreatedAt() time.Time {
	return b.createdAt
}

// UpdatedAt returns when the budget was last updated.
func (b *BudgetHold) UpdatedAt() time.Time {
	return b.updatedAt
}
