package query

import (
	"context"
	"time"
)

type CategoryNameDTO struct {
	ID   string
	Name string
}

type CategoryDTO struct {
	ID   string
	Name string
	Type string // this is the "category type" — "income", "needs", "wants", "savings"
}

type ExpenseCategoryBudgetSpent struct {
	CategoryID string
	BudgetID   string
	Type       string
	Name       string
	Budgeted   int
	Spent      int
}

type CategoryBudget struct {
	Name  string
	Type  string
	Value int // if it's income amount will be the money in
}

type CategoryMonthlyValue struct {
	ID    string
	Name  string
	Month int
	Value int
}

type CategoryQueryService interface {
	FindByID(ctx context.Context, id string) (*CategoryDTO, error)
	FindAllCategories(ctx context.Context, uid string) ([]CategoryDTO, error)
	FindCategoriesIDIncludes(ctx context.Context, ids []string) ([]CategoryNameDTO, error)
	GetExpenseCategoriesBudgetSpent(
		ctx context.Context,
		uid string,
		minDate time.Time,
		maxDate time.Time,
	) ([]ExpenseCategoryBudgetSpent, error)
	GetCategoriesBudgetSpent(
		ctx context.Context,
		uid string,
		minDate time.Time,
		maxDate time.Time,
	) ([]CategoryBudget, error)
	GetYearMonthlyValue(
		ctx context.Context,
		year int,
		id string,
	) ([]CategoryMonthlyValue, error)
}
