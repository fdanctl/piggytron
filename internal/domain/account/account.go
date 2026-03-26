package account

import "time"

type ID string

func NewID(str string) (ID, error) {
	if str == "" {
		return "", ErrInvalidID
	}
	return ID(str), nil
}

type AccountType string

const (
	bank AccountType = "bank"
	goal AccountType = "goal"
)

type Account struct {
	id       ID
	userID   ID
	aType    AccountType
	name     string
	currency string
	// goal-specific
	targetAmount *int
	targetDate   *time.Time
	categoryID   *ID

	createdAt time.Time
	updatedAt time.Time
}

func NewBank(
	id, userID ID,
	name string,
	currency string,
) (*Account, error) {
	if name == "" || len(name) > 50 {
		return nil, ErrInvalidName
	}
	if currency == "" || len(currency) > 10 {
		return nil, ErrInvalidCurrency
	}

	now := time.Now()

	return &Account{
		id:        id,
		userID:    userID,
		name:      name,
		aType:     bank,
		currency:  currency,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func NewGoal(
	id, userID ID,
	name string,
	currency string,
	targetAmount int,
	targetDate *time.Time,
	categoryID ID,
) (*Account, error) {
	if name == "" || len(name) > 50 {
		return nil, ErrInvalidName
	}
	if currency == "" || len(currency) > 10 {
		return nil, ErrInvalidCurrency
	}
	if targetAmount <= 0 {
		return nil, ErrNegativeNumber
	}

	now := time.Now()

	return &Account{
		id:           id,
		userID:       userID,
		name:         name,
		aType:        goal,
		targetAmount: &targetAmount,
		targetDate:   targetDate,
		categoryID:   &categoryID,
		currency:     currency,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

func Rehydrate(
	id, userID ID,
	aType AccountType,
	name string,
	targetAmount *int,
	targetDate *time.Time,
	categoryID *ID,
	currency string,
	createdAt, updatedAt time.Time,
) *Account {
	return &Account{
		id:           id,
		userID:       userID,
		name:         name,
		targetAmount: targetAmount,
		targetDate:   targetDate,
		categoryID:   categoryID,
		currency:     currency,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

func (b *Account) ID() ID {
	return b.id
}

func (b *Account) UserID() ID {
	return b.userID
}

func (b *Account) Name() string {
	return b.name
}

func (b *Account) Type() AccountType {
	return b.aType
}

func (b *Account) Currency() string {
	return b.currency
}

func (b *Account) TargetAmount() *int {
	return b.targetAmount
}

func (b *Account) TargetDate() *time.Time {
	return b.targetDate
}

func (b *Account) CategoryID() *ID {
	return b.categoryID
}

func (b *Account) CreatedAt() time.Time {
	return b.createdAt
}

func (b *Account) UpdatedAt() time.Time {
	return b.updatedAt
}

// TODO change name and currency
