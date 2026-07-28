package budget

import "time"

type ID string

type Budget struct {
	categoryID ID
	month      Month
	amount     int

	createdAt time.Time
	updatedAt time.Time
}

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

func (b *Budget) ChangeAmount(nAmount int) error {
	if nAmount < 0 {
		return ErrInvalidAmount
	}

	b.amount = nAmount
	b.updatedAt = time.Now()
	return nil
}

func (b *Budget) CategoryID() ID {
	return b.categoryID
}

func (b *Budget) Month() Month {
	return b.month
}

func (b *Budget) Amount() int {
	return b.amount
}

func (b *Budget) CreatedAt() time.Time {
	return b.createdAt
}

func (b *Budget) UpdatedAt() time.Time {
	return b.updatedAt
}
