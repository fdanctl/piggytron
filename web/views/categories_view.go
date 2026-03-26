package views

import (
	expensecategory "github.com/fdanctl/piggytron/internal/domain/expense_category"
	incomecategory "github.com/fdanctl/piggytron/internal/domain/income_category"
)

type Category interface {
	GetID() string
	GetName() string
	GetExpenseType() string
}

type IncomeCategory struct {
	ID   incomecategory.ID
	Name string
}

func NewIncomeCategory(c *incomecategory.IncomeCategory) IncomeCategory {
	return IncomeCategory{
		ID:   c.ID(),
		Name: c.Name(),
	}
}

func (c IncomeCategory) GetID() string {
	return string(c.ID)
}

func (c IncomeCategory) GetName() string {
	return c.Name
}

func (c IncomeCategory) GetExpenseType() string {
	return ""
}

type ExpenseCategory struct {
	ID          expensecategory.ID
	Name        string
	ExpenseType expensecategory.ExpenseType
}

func NewExpenseCategory(c *expensecategory.ExpenseCategory) ExpenseCategory {
	return ExpenseCategory{
		ID:          c.ID(),
		Name:        c.Name(),
		ExpenseType: c.ExpenseType(),
	}
}

func (c ExpenseCategory) GetID() string {
	return string(c.ID)
}

func (c ExpenseCategory) GetName() string {
	return c.Name
}

func (c ExpenseCategory) GetExpenseType() string {
	return string(c.ExpenseType)
}

type CategoriesView struct {
	IncomeCategories  []IncomeCategory
	ExpenseCategories []ExpenseCategory
}
