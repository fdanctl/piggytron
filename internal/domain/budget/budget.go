// Package budget defines the monthly budget aggregate: an amount (in cents)
// per expense category and month, together with the Month type used across
// the application.
package budget

import "time"

// ID is a budget identifier.
type ID string

// Budget is the planned spending limit for one expense category in one month.
type Budget struct {
	categoryID ID
	month      Month
	amount     int

	createdAt time.Time
	updatedAt time.Time
}

// New builds a validated budget; a zero amount means "no budget set".
func New(
	categoryID ID,
	month Month,
	amount int,
) (*Budget, error) {
	if amount < 0 {
		return nil, ErrInvalidAmount
	}

	now := time.Now()

	return &Budget{
		categoryID: categoryID,
		month:      month,
		amount:     amount,

		createdAt: now,
		updatedAt: now,
	}, nil
}

// Rehydrate rebuilds a Budget from persistence; the month is normalized to
// the first of the month.
func Rehydrate(
	categoryID ID,
	month time.Time,
	amount int,
	createdAt, updatedAt time.Time,
) *Budget {
	return &Budget{
		categoryID: categoryID,
		month:      NewMonth(month),
		amount:     amount,

		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// ChangeAmount replaces the budget limit, requiring a non-negative value.
func (b *Budget) ChangeAmount(nAmount int) error {
	if nAmount < 0 {
		return ErrInvalidAmount
	}

	b.amount = nAmount
	b.updatedAt = time.Now()
	return nil
}

// CategoryID returns the budgeted expense category.
func (b *Budget) CategoryID() ID {
	return b.categoryID
}

// Month returns the budgeted month.
func (b *Budget) Month() Month {
	return b.month
}

// Amount returns the budget limit in cents.
func (b *Budget) Amount() int {
	return b.amount
}

// CreatedAt returns when the budget was created.
func (b *Budget) CreatedAt() time.Time {
	return b.createdAt
}

// UpdatedAt returns when the budget was last updated.
func (b *Budget) UpdatedAt() time.Time {
	return b.updatedAt
}
