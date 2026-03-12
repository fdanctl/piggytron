package views

import (
	expensecategory "github.com/fdanctl/piggytron/internal/domain/expense_category"
	incomecategory "github.com/fdanctl/piggytron/internal/domain/income_category"
)

type IncomeCategories struct {
	Id   incomecategory.ID
	Name string
}

type ExpenseCategories struct {
	Id          expensecategory.ID
	Name        string
	ExpenseType expensecategory.ExpenseType
}

type CategoriesView struct {
	IncomeCategories  []IncomeCategories
	ExpenseCategories []ExpenseCategories
}
