package views

import (
	"time"

	"github.com/fdanctl/piggytron/internal/domain/expensecategory"
	"github.com/fdanctl/piggytron/internal/domain/incomecategory"
)

// Category is the common view-model interface for income and expense
// categories.
type Category interface {
	GetID() string
	GetStatus() string
	GetName() string
	GetType() string
	GetExpenseType() string
	GetArchivedAt() *time.Time
}

// IncomeCategory is the view-model for an income category.
type IncomeCategory struct {
	ID         incomecategory.ID
	Status     string
	Name       string
	Type       string
	ArchivedAt *time.Time
}

// NewIncomeCategory builds the view-model from the domain category.
func NewIncomeCategory(c *incomecategory.IncomeCategory) IncomeCategory {
	return IncomeCategory{
		ID:         c.ID(),
		Status:     string(c.Status()),
		Type:       "income",
		Name:       c.Name(),
		ArchivedAt: c.ArchivedAt(),
	}
}

// GetID returns the category id as a string.
func (c IncomeCategory) GetID() string {
	return string(c.ID)
}

// GetStatus returns the category status (active or archived).
func (c IncomeCategory) GetStatus() string {
	return c.Status
}

// GetName returns the category name.
func (c IncomeCategory) GetName() string {
	return c.Name
}

// GetType returns the category type (income or exprense).
func (c IncomeCategory) GetType() string {
	return c.Type
}

// GetExpenseType returns "" — income categories have no expense type.
func (c IncomeCategory) GetExpenseType() string {
	return ""
}

// GetArchivedAt returns the date the category was archived or nil it's active.
func (c IncomeCategory) GetArchivedAt() *time.Time {
	return c.ArchivedAt
}

// ExpenseCategory is the view-model for an expense category.
type ExpenseCategory struct {
	ID          expensecategory.ID
	Status      string
	Name        string
	Type        string
	ExpenseType expensecategory.ExpenseType
	ArchivedAt  *time.Time
}

// NewExpenseCategory builds the view-model from the domain category.
func NewExpenseCategory(c *expensecategory.ExpenseCategory) ExpenseCategory {
	return ExpenseCategory{
		ID:          c.ID(),
		Status:      string(c.Status()),
		Name:        c.Name(),
		Type:        "expense",
		ExpenseType: c.ExpenseType(),
		ArchivedAt:  c.ArchivedAt(),
	}
}

// GetID returns the category id as a string.
func (c ExpenseCategory) GetID() string {
	return string(c.ID)
}

// GetStatus returns the category status (active or archived).
func (c ExpenseCategory) GetStatus() string {
	return c.Status
}

// GetName returns the category name.
func (c ExpenseCategory) GetName() string {
	return c.Name
}

// GetType returns the category type (income or exprense).
func (c ExpenseCategory) GetType() string {
	return c.Type
}

// GetExpenseType returns the category's expense type ("needs", "wants" or
// "savings").
func (c ExpenseCategory) GetExpenseType() string {
	return string(c.ExpenseType)
}

// GetArchivedAt returns the date the category was archived or nil it's active.
func (c ExpenseCategory) GetArchivedAt() *time.Time {
	return c.ArchivedAt
}

// CategoriesView groups the two category lists for the categories page.
type CategoriesView struct {
	IncomeCategories  []IncomeCategory
	ExpenseCategories []ExpenseCategory
}
