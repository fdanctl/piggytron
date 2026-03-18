package bank

import "time"

type ID string

func NewId(str string) (ID, error) {
	if str == "" {
		return "", ErrInvalidID
	}
	return ID(str), nil
}

type BankType string

type Bank struct {
	id        ID
	userId    ID
	name      string
	currency  string
	createdAt time.Time
	updatedAt time.Time
}

func New(
	id, userId ID,
	name string,
	currency string,
) (*Bank, error) {
	if name == "" || len(name) > 50 {
		return nil, ErrInvalidName
	}
	if currency == "" || len(currency) > 10 {
		return nil, ErrInvalidCurrency
	}

	now := time.Now()

	return &Bank{
		id:        id,
		userId:    userId,
		name:      name,
		currency:  currency,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func Rehydrate(
	id, userId ID,
	name string,
	currency string,
	createdAt, updatedAt time.Time,
) *Bank {
	return &Bank{
		id:        id,
		userId:    userId,
		name:      name,
		currency:  currency,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (b *Bank) ID() ID {
	return b.id
}

func (b *Bank) UserId() ID {
	return b.userId
}

func (b *Bank) Name() string {
	return b.name
}

func (b *Bank) Currency() string {
	return b.currency
}

func (b *Bank) CreatedAt() time.Time {
	return b.createdAt
}

func (b *Bank) UpdatedAt() time.Time {
	return b.updatedAt
}

// TODO change name and currency
