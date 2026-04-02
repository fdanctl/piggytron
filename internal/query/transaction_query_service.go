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

type TransactionQueryService interface {
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
}
