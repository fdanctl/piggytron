package views

import (
	expensecategory "github.com/fdanctl/piggytron/internal/domain/expense_category"
	incomecategory "github.com/fdanctl/piggytron/internal/domain/income_category"
)

type Category interface {
	GetId() string
	GetName() string
	GetExpenseType() uint8
}

type IncomeCategory struct {
	Id   incomecategory.ID
	Name string
}

func (c IncomeCategory) GetId() string {
	return string(c.Id)
}

func (c IncomeCategory) GetName() string {
	return c.Name
}

func (c IncomeCategory) GetExpenseType() uint8 {
	return 0
}

type ExpenseCategory struct {
	Id          expensecategory.ID
	Name        string
	ExpenseType expensecategory.ExpenseType
}

func (c ExpenseCategory) GetId() string {
	return string(c.Id)
}

func (c ExpenseCategory) GetName() string {
	return c.Name
}

func (c ExpenseCategory) GetExpenseType() uint8 {
	return uint8(c.ExpenseType)
}

type CategoriesView struct {
	IncomeCategories  []IncomeCategory
	ExpenseCategories []ExpenseCategory
}
