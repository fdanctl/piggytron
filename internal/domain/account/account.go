package account

import (
	"errors"
	"time"
)

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
	isSaving *bool // bank-specific
	currency string
	// goal-specific
	targetAmount *int
	startDate    *time.Time
	targetDate   *time.Time
	categoryID   *ID

	createdAt time.Time
	updatedAt time.Time
}

func NewBank(
	id, userID ID,
	name string,
	currency string,
	isSaving bool,
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
		isSaving:  &isSaving,
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
	startDate time.Time,
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
		startDate:    &startDate,
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
	isSaving *bool,
	targetAmount *int,
	startDate *time.Time,
	targetDate *time.Time,
	categoryID *ID,
	currency string,
	createdAt, updatedAt time.Time,
) *Account {
	return &Account{
		id:           id,
		userID:       userID,
		aType:        aType,
		name:         name,
		isSaving:     isSaving,
		targetAmount: targetAmount,
		startDate:    startDate,
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

func (b *Account) IsSaving() *bool {
	return b.isSaving
}

func (b *Account) Currency() string {
	return b.currency
}

func (b *Account) TargetAmount() *int {
	return b.targetAmount
}

func (b *Account) StartDate() *time.Time {
	return b.startDate
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

func (b *Account) CanReceiveIncome() error {
	if b.aType == goal {
		return errors.New("goals can't receive from outside")
	}
	if b.isSaving != nil && *b.isSaving {
		return errors.New("savings accounts can't receive from outside")
	}
	return nil
}

func (b *Account) CanMakeExpense() error {
	if b.aType == goal {
		return errors.New("goals can't make expenses")
	}
	if b.isSaving != nil && *b.isSaving {
		return errors.New("savings accounts can't make expenses")
	}
	return nil
}

// updates

func (b *Account) ChangeName(name string) error {
	if name == "" || len(name) > 50 {
		return ErrInvalidName
	}
	b.name = name
	b.updatedAt = time.Now()
	return nil
}

func (b *Account) ChangeTargetAmount(amount int) error {
	if b.aType != goal {
		return ErrAccountWrongType
	}
	if amount <= 0 {
		return ErrNegativeNumber
	}
	b.targetAmount = &amount
	b.updatedAt = time.Now()
	return nil
}

func (b *Account) ChangeStartDate(date time.Time, minPossible *time.Time) error {
	if b.aType != goal {
		return ErrAccountWrongType
	}
	if minPossible != nil && date.Compare(*minPossible) == 1 {
		return ErrContributionBeforeStartDate
	}
	b.startDate = &date
	b.updatedAt = time.Now()
	return nil
}

func (b *Account) ChangeTargetDate(date *time.Time) error {
	if b.aType != goal {
		return ErrAccountWrongType
	}
	b.targetDate = date
	b.updatedAt = time.Now()
	return nil
}

func (b *Account) ChangeCategory(cid ID) error {
	if b.aType != goal {
		return ErrAccountWrongType
	}
	b.categoryID = &cid
	b.updatedAt = time.Now()
	return nil
}
