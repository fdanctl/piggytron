package query

import (
	"context"
	"time"
)

// LedgerEntryDTO is a ledger entry.
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
	Note        *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// EntriesWithTotalCount is a page of ledger entries with the total number of
// matching rows.
type EntriesWithTotalCount struct {
	Data  []LedgerEntryDTO
	Total int
}

// CategoryExpense is the amount spent in one expense category over a range.
type CategoryExpense struct {
	ID     string
	Amount int
}

// CategoryExpenseWithTotal is the per-category expense report plus the total
// spent.
type CategoryExpenseWithTotal struct {
	Data  []CategoryExpense
	Total int
}

// LedgerQueryService is the read-model contract for ledger.
// It is implemented by the Postgres infrastructure.
type LedgerQueryService interface {
	FindByID(ctx context.Context, id string) (*LedgerEntryDTO, error)
	FindFiltered(
		ctx context.Context,
		uid string,
		filters *LedgerFilters,
		limit, offset uint,
	) ([]LedgerEntryDTO, error)
	FindAllWithExpenseCategoryWithCount(
		ctx context.Context,
		uid string,
		minDate, maxDate time.Time,
		limit, offset uint,
	) (*EntriesWithTotalCount, error)
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
	// GetMinMaxAmountAndDate returns minAmount, maxAmount, minDate, maxDate, error
	GetMinMaxAmountAndDate(ctx context.Context, uid string) (int, int, time.Time, time.Time, error)
	// GetFirstEntryDate returns the date of the user's oldest ledger entry.
	GetFirstEntryDate(ctx context.Context, uid string) (time.Time, error)
}
