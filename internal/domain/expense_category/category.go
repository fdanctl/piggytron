package expensecategory

import "time"

type ID string

func NewId(str string) (ID, error) {
	if str == "" {
		return "", ErrInvalidID
	}
	return ID(str), nil
}

type ExpenseType uint8

const (
	needs ExpenseType = iota + 1
	wants
	savings
)

func NewExpenseType(num uint8) (ExpenseType, error) {
	if num <= 0 || num > 3 {
		return 0, ErrInvalidType
	}
	return ExpenseType(num), nil
}

type ExpenseCategory struct {
	id          ID
	userId      ID
	name        string
	expenseType ExpenseType
	createdAt   time.Time
	updatedAt   time.Time
}

func New(id ID, userId ID, name string, expenseType ExpenseType) (*ExpenseCategory, error) {
	if name == "" {
		return nil, ErrInvalidName
	}

	now := time.Now()

	return &ExpenseCategory{
		id:          id,
		userId:      userId,
		name:        name,
		expenseType: expenseType,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

func Rehydrate(
	id ID,
	userId ID,
	name string,
	expenseType ExpenseType,
	createdAt, updatedAt time.Time,
) *ExpenseCategory {
	return &ExpenseCategory{
		id:          id,
		userId:      userId,
		name:        name,
		expenseType: expenseType,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func (ec *ExpenseCategory) ID() ID {
	return ec.id
}

func (ec *ExpenseCategory) UserId() ID {
	return ec.userId
}

func (ec *ExpenseCategory) Name() string {
	return ec.name
}

func (ec *ExpenseCategory) ExpenseType() ExpenseType {
	return ec.expenseType
}

func (ec *ExpenseCategory) CreatedAt() time.Time {
	return ec.createdAt
}

func (ec *ExpenseCategory) UpdatedAt() time.Time {
	return ec.updatedAt
}

func (ec *ExpenseCategory) ChangeName(name string) error {
	if name == "" {
		return ErrInvalidName
	}
	ec.name = name
	ec.updatedAt = time.Now()
	return nil
}
