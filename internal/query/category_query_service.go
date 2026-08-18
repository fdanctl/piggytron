package query

import (
	"context"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/budget"
)

// CategoryNameDTO is a lightweight category reference (id + name).
type CategoryNameDTO struct {
	ID   string
	Name string
}

// CategoryDTO is a category row with its type: "income" for income
// categories, or one of "needs", "wants", "savings" for expense categories.
type CategoryDTO struct {
	ID         string
	Name       string
	Type       string // this is the "category type" — "income", "needs", "wants", "savings"
	Status     string // "active" or "archived"
	ArchivedAt *time.Time
}

// CategoryBudgetValue is one row of a category budget report: how much was
// budgeted and spent (or earned) for a category in a month.
type CategoryBudgetValue struct {
	CategoryID      string
	Month           time.Time
	Type            string
	Name            string
	Budgeted        int
	Value           int // spent or income
	ArchivedAt      *time.Time
	PrevTotalBudget int
	PrevTotalSpent  int
}

// CategoryBudget aggregates the budget spent (or income earned) per category
// over a date range.
type CategoryBudget struct {
	Name       string
	Type       string
	Value      int // if it's income amount will be the money in
	ArchivedAt *time.Time
}

// CategoryMonthlyValue is the monthly value (spent or earned) of one category
// across the months of a year.
type CategoryMonthlyValue struct {
	Month time.Time
	ID    string
	Name  string
	Value int
}

// MonthExpenseCategoryBudgetSpentWithBalance is the per-category budget spent
// report for a month, plus the month net and the non-savings bank accounts balance.
type MonthExpenseCategoryBudgetSpentWithBalance struct {
	Data     []CategoryBudgetValue
	MonthNet int
	Balance  int
}

// CategoryQueryService is the read-model contract for category and budget
// reporting data. It is implemented by the Postgres infrastructure.
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
	GetMonthlyValueSince(
		ctx context.Context,
		id string,
		since time.Time,
	) ([]CategoryMonthlyValue, error)
}
