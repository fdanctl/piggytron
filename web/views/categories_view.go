package views

import (
	"github.com/fdanctl/piggytron/internal/domain/expensecategory"
	"github.com/fdanctl/piggytron/internal/domain/incomecategory"
)

// Category is the common view-model interface for income and expense
// categories.
type Category interface {
	GetID() string
	GetName() string
	GetExpenseType() string
}

// IncomeCategory is the view-model for an income category.
type IncomeCategory struct {
	ID   incomecategory.ID
	Name string
}

// NewIncomeCategory builds the view-model from the domain category.
func NewIncomeCategory(c *incomecategory.IncomeCategory) IncomeCategory {
	return IncomeCategory{
		ID:   c.ID(),
		Name: c.Name(),
	}
}

// GetID returns the category id as a string.
func (c IncomeCategory) GetID() string {
	return string(c.ID)
}

// GetName returns the category name.
func (c IncomeCategory) GetName() string {
	return c.Name
}

// GetExpenseType returns "" — income categories have no expense type.
func (c IncomeCategory) GetExpenseType() string {
	return ""
}

// ExpenseCategory is the view-model for an expense category.
type ExpenseCategory struct {
	ID          expensecategory.ID
	Name        string
	ExpenseType expensecategory.ExpenseType
}

// NewExpenseCategory builds the view-model from the domain category.
func NewExpenseCategory(c *expensecategory.ExpenseCategory) ExpenseCategory {
	return ExpenseCategory{
		ID:          c.ID(),
		Name:        c.Name(),
		ExpenseType: c.ExpenseType(),
	}
}

// GetID returns the category id as a string.
func (c ExpenseCategory) GetID() string {
	return string(c.ID)
}

// GetName returns the category name.
func (c ExpenseCategory) GetName() string {
	return c.Name
}

// GetExpenseType returns the category's expense type ("needs", "wants" or
// "savings").
func (c ExpenseCategory) GetExpenseType() string {
	return string(c.ExpenseType)
}

// CategoriesView groups the two category lists for the categories page.
type CategoriesView struct {
	IncomeCategories  []IncomeCategory
	ExpenseCategories []ExpenseCategory
}
