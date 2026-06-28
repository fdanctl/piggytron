package query

import (
	"context"
	"time"
)

type TransactionDTO struct {
	ID     string
	UserID string

	Type string

	FromAccount *string
	ToAccount   *string

	IncomeCategory  *string
	ExpenseCategory *string

	Amount      int
	Description string
	Date        time.Time
	CreatedAt   time.Time
}

type TransactionsWithTotalCount struct {
	Data  []TransactionDTO
	Total int
}

type CategoryExpense struct {
	ID     string
	Amount int
}

type CategoryExpenseWithTotal struct {
	Data  []CategoryExpense
	Total int
}

type TransactionQueryService interface {
	FindByID(ctx context.Context, id string) (*TransactionDTO, error)
	FindFiltered(
		ctx context.Context,
		uid string,
		filters *TransactionFilters,
		limit, offset uint,
	) ([]TransactionDTO, error)
	FindFilteredWithCount(
		ctx context.Context,
		uid string,
		filters *TransactionFilters,
		limit, offset uint,
	) (*TransactionsWithTotalCount, error)
	CountFilteredResults(
		ctx context.Context, uid string, filters *TransactionFilters,
	) (int, error)
	GetExpensesByCategoryBetweenDates(
		ctx context.Context, uid string, minDate time.Time, maxDate time.Time,
	) (*CategoryExpenseWithTotal, error)
	GetRecentTransactions(
		ctx context.Context, uid string, limit uint,
	) ([]TransactionDTO, error)
	GetMinMax(ctx context.Context, uid string) (int, int, error)
}
