package query

import (
	"context"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/budget"
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

type CategoryBudgetValue struct {
	CategoryID      string
	Month           time.Time
	Type            string
	Name            string
	Budgeted        int
	Value           int // spent or income
	PrevTotalBudget int
	PrevTotalSpent  int
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

type MonthExpenseCategoryBudgetSpentWithBalance struct {
	Data     []CategoryBudgetValue
	MonthNet int
	Balance  int
}

type CategoryQueryService interface {
	FindByID(ctx context.Context, id string) (*CategoryDTO, error)
	FindAllCategories(ctx context.Context, uid string) ([]CategoryDTO, error)
	FindCategoriesIDIncludes(ctx context.Context, ids []string) ([]CategoryNameDTO, error)
	GetCategoriesBudgetSpentValue(
		ctx context.Context,
		uid string,
		month budget.Month,
	) (*MonthExpenseCategoryBudgetSpentWithBalance, error)
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
