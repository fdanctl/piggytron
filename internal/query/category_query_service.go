package query

import (
	"context"
	"time"
)

type CategoryNameDTO struct {
	ID   string
	Name string
}

type ExpenseCategoryBudgetSpent struct {
	CID      string
	BID      string
	Type     string
	Name     string
	Budgeted int
	Spent    int
}

type CategoryQueryService interface {
	FindAllCategories(ctx context.Context, uid string) ([]CategoryNameDTO, error)
	FindCategoriesIDIncludes(ctx context.Context, ids []string) ([]CategoryNameDTO, error)
	GetExpenseCategoriesBudgetSpent(
		ctx context.Context,
		uid string,
		minDate time.Time,
		maxDate time.Time,
	) ([]ExpenseCategoryBudgetSpent, error)
}
