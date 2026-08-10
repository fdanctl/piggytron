// Package expensecategory defines the expense category aggregate with its
// needs / wants / savings type.
package expensecategory

import "time"

// ID is an expense category identifier.
type ID string

// ExpenseType classifies a category as needs, wants or savings.
type ExpenseType string

const (
	Needs   ExpenseType = "needs"
	Wants   ExpenseType = "wants"
	Savings ExpenseType = "savings"
)

// NewExpenseType parses a string into an ExpenseType, returning
// ErrInvalidType otherwise.
func NewExpenseType(str string) (ExpenseType, error) {
	switch str {
	case "needs":
		return Needs, nil

	case "wants":
		return Wants, nil

	case "savings":
		return Savings, nil

	default:
		return "", ErrInvalidType
	}
}

// ExpenseCategory groups expenses of a single type, owned by one user.
type ExpenseCategory struct {
	id          ID
	userID      ID
	name        string
	expenseType ExpenseType
	createdAt   time.Time
	updatedAt   time.Time
}

// New builds a validated expense category.
func New(id ID, userID ID, name string, expenseType ExpenseType) (*ExpenseCategory, error) {
	if name == "" || len(name) > 30 {
		return nil, ErrInvalidName
	}

	now := time.Now()

	return &ExpenseCategory{
		id:          id,
		userID:      userID,
		name:        name,
		expenseType: expenseType,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// Rehydrate rebuilds an ExpenseCategory from persistence without re-running
// validation (the database constraints already guard it).
func Rehydrate(
	id ID,
	userID ID,
	name string,
	expenseType ExpenseType,
	createdAt, updatedAt time.Time,
) *ExpenseCategory {
	return &ExpenseCategory{
		id:          id,
		userID:      userID,
		name:        name,
		expenseType: expenseType,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

// ID returns the category id.
func (ec *ExpenseCategory) ID() ID {
	return ec.id
}

// UserID returns the id of the user who owns the category.
func (ec *ExpenseCategory) UserID() ID {
	return ec.userID
}

// Name returns the category name.
func (ec *ExpenseCategory) Name() string {
	return ec.name
}

// ExpenseType returns the category type (needs, wants or savings).
func (ec *ExpenseCategory) ExpenseType() ExpenseType {
	return ec.expenseType
}

// CreatedAt returns when the category was created.
func (ec *ExpenseCategory) CreatedAt() time.Time {
	return ec.createdAt
}

// UpdatedAt returns when the category was last updated.
func (ec *ExpenseCategory) UpdatedAt() time.Time {
	return ec.updatedAt
}

// ChangeName renames the category.
func (ec *ExpenseCategory) ChangeName(name string) error {
	if name == "" || len(name) > 30 {
		return ErrInvalidName
	}
	ec.name = name
	ec.updatedAt = time.Now()
	return nil
}

// ChangeType reclassifies the category, rejecting a no-op change.
func (ec *ExpenseCategory) ChangeType(t ExpenseType) error {
	if t == ec.expenseType {
		return ErrSameType
	}
	ec.expenseType = t
	ec.updatedAt = time.Now()
	return nil
}
