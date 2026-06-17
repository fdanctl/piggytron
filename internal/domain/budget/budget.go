package budget

import "time"

type ID string

type Budget struct {
	id         ID
	userID     ID
	categoryID ID
	month      time.Time
	amount     int

	createdAt time.Time
	updatedAt time.Time
}

func New(
	id, userID, categoryID ID,
	month time.Time,
	amount int,
) (*Budget, error) {
	if amount < 0 {
		return nil, ErrInvalidAmount
	}

	now := time.Now()

	return &Budget{
		id:         id,
		userID:     userID,
		categoryID: categoryID,
		month:      month,
		amount:     amount,

		createdAt: now,
		updatedAt: now,
	}, nil
}

func Rehydrate(
	id, userID, categoryID ID,
	month time.Time,
	amount int,
	createdAt, updatedAt time.Time,
) *Budget {
	return &Budget{
		id:         id,
		userID:     userID,
		categoryID: categoryID,
		month:      month,
		amount:     amount,

		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (b *Budget) ChangeAmount(nAmount int) error {
	if nAmount <= 0 {
		return ErrInvalidAmount
	}

	b.amount = nAmount
	b.updatedAt = time.Now()
	return nil
}

func (b *Budget) ID() ID {
	return b.id
}

func (b *Budget) UserID() ID {
	return b.userID
}

func (b *Budget) CategoryID() ID {
	return b.categoryID
}

func (b *Budget) Month() time.Time {
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
