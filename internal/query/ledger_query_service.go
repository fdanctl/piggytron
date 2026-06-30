package query

import (
	"context"
	"time"
)

type LedgerEntryDTO struct {
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

type EntriesWithTotalCount struct {
	Data  []LedgerEntryDTO
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

type LedgerQueryService interface {
	FindByID(ctx context.Context, id string) (*LedgerEntryDTO, error)
	FindFiltered(
		ctx context.Context,
		uid string,
		filters *LedgerFilters,
		limit, offset uint,
	) ([]LedgerEntryDTO, error)
	FindFilteredWithCount(
		ctx context.Context,
		uid string,
		filters *LedgerFilters,
		limit, offset uint,
	) (*EntriesWithTotalCount, error)
	CountFilteredResults(
		ctx context.Context, uid string, filters *LedgerFilters,
	) (int, error)
	GetExpensesByCategoryBetweenDates(
		ctx context.Context, uid string, minDate time.Time, maxDate time.Time,
	) (*CategoryExpenseWithTotal, error)
	GetRecentEntries(
		ctx context.Context, uid string, limit uint,
	) ([]LedgerEntryDTO, error)
	// GetMinMaxAmountAndDate return minAmount, maxAmount, minDate, maxDate, error
	GetMinMaxAmountAndDate(ctx context.Context, uid string) (int, int, time.Time, time.Time, error)
}
