package expensecategory

import "time"

type ID string

func NewId(str string) (ID, error) {
	if str == "" {
		return "", ErrInvalidID
	}
	return ID(str), nil
}

type ExpenseType string

const (
	needs   ExpenseType = "needs"
	wants   ExpenseType = "wants"
	savings ExpenseType = "savings"
)

func NewExpenseType(str string) (ExpenseType, error) {
	switch str {
	case "needs":
		return needs, nil

	case "wants":
		return wants, nil

	case "savings":
		return savings, nil

	default:
		return "", ErrInvalidType
	}
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
	if name == "" || len(name) > 30 {
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

// TODO
func (ec *ExpenseCategory) ChangeName(name string) error {
	if name == "" {
		return ErrInvalidName
	}
	ec.name = name
	ec.updatedAt = time.Now()
	return nil
}
