package views

import (
	expensecategory "github.com/fdanctl/piggytron/internal/domain/expense_category"
	incomecategory "github.com/fdanctl/piggytron/internal/domain/income_category"
)

type Category interface {
	GetId() string
	GetName() string
	GetExpenseType() string
}

type IncomeCategory struct {
	Id   incomecategory.ID
	Name string
}

func NewIncomeCategory(c *incomecategory.IncomeCategory) IncomeCategory {
	return IncomeCategory{
		Id:   c.ID(),
		Name: c.Name(),
	}
}

func (c IncomeCategory) GetId() string {
	return string(c.Id)
}

func (c IncomeCategory) GetName() string {
	return c.Name
}

func (c IncomeCategory) GetExpenseType() string {
	return ""
}

type ExpenseCategory struct {
	Id          expensecategory.ID
	Name        string
	ExpenseType expensecategory.ExpenseType
}

func NewExpenseCategory(c *expensecategory.ExpenseCategory) ExpenseCategory {
	return ExpenseCategory{
		Id:          c.ID(),
		Name:        c.Name(),
		ExpenseType: c.ExpenseType(),
	}
}

func (c ExpenseCategory) GetId() string {
	return string(c.Id)
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
