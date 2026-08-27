// Package query defines the read-model contracts: the DTOs returned to the
// HTTP layer and the query service interfaces implemented by the Postgres
// infrastructure. It is the seam that keeps SQL out of the interface layer.
package query

import (
	"context"
	"errors"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/monthlysummary"
)

// ErrNoHistory is returned by queries that produce an empty result, such as a
// daily-balance range with no chartable days.
var ErrNoHistory = errors.New("no history found")

// AccountIDName is a lightweight account reference (id + name) used for
// dropdowns and filters.
type AccountIDName struct {
	ID   string
	Name string
}

// AccountWithCategory is an account row joined with its category (goals only).
type AccountWithCategory struct {
	ID       string
	UserID   string
	Type     string
	Name     string
	Status   string
	Currency string
	// goal-specific
	TargetAmount    *int
	StartDate       *time.Time
	TargetDate      *time.Time
	Category        *CategoryDTO
	CompletedAt     *time.Time
	CancelledAt     *time.Time
	FinalizedAmount *int

	ClosedAt  *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// AccountWithSum is an account row with its current balance (Sum), in cents.
type AccountWithSum struct {
	AccountWithCategory
	Sum int
}

// AccountWithSumAndMonthChange is an account row with its balance and the
// money in/out for the given month.
type AccountWithSumAndMonthChange struct {
	AccountWithSum
	MoneyIn  int
	MoneyOut int
}

// AccountWithMinRunningBalance is an account row plus the minimum running
// balance observed between two dates and the day it occurred.
type AccountWithMinRunningBalance struct {
	AccountWithCategory
	MinRunningBalance int
	MinDate           time.Time
}

// AccountDailyBalance is one chart point: the account balance at the end of
// a day (in cents).
type AccountDailyBalance struct {
	Day     time.Time
	ID      string
	Name    string
	Type    string
	Balance int
}

// AccountDailyBalanceWithStatsSince is a daily balance series plus the
// cumulative money in/out and transaction count over the whole period.
type AccountDailyBalanceWithStatsSince struct {
	Data         []AccountDailyBalance
	MoneyIn      int
	MoneyOut     int
	Transactions int
}

// AccountQueryService is the read-model contract for account.
// It is implemented by the Postgres infrastructure.
type AccountQueryService interface {
	FindIDNamesIncludes(ctx context.Context, ids []string) ([]AccountIDName, error)
	FindBanksIDNames(ctx context.Context, uid string) ([]AccountIDName, error)
	FindGoalsIDNames(ctx context.Context, uid string) ([]AccountIDName, error)
	FindGoalsByCategory(ctx context.Context, cid string) ([]AccountWithCategory, error)
	FindWithSum(ctx context.Context, id string) (*AccountWithSum, error)
	FindAllWithSum(ctx context.Context, uid string) ([]AccountWithSum, error)
	FindAllWithSumAndMonthChange(
		ctx context.Context,
		uid string,
		month monthlysummary.Month,
	) ([]AccountWithSumAndMonthChange, error)
	FindAllGoalsWithSum(ctx context.Context, uid string) ([]AccountWithSum, error)
	// GetAllDailyBalanceSince returns the daily balance of every account
	// from the 1st of the month of since to today. Since is truncated
	// to the month because monthly_summary is month-granular.
	GetAllDailyBalanceSince(
		ctx context.Context,
		uid string,
		since time.Time,
	) ([]AccountDailyBalance, error)
	// GetAccountDailyBalanceSince returns the daily balance of one account
	// from the 1st of the month of since to today.
	GetAccountDailyBalanceSince(
		ctx context.Context,
		id string,
		since time.Time,
	) ([]AccountDailyBalance, error)
	// GetAccountDailyBalanceAndStatsSince is GetAccountDailyBalanceSince plus
	// the cumulative money in/out and transaction count over the period.
	GetAccountDailyBalanceAndStatsSince(
		ctx context.Context,
		id string,
		since time.Time,
	) (*AccountDailyBalanceWithStatsSince, error)
	// GetAccountWithMinRunningBalance returns the minimum running balance of
	// the account between fromDate and untilDate, optionally excluding one
	// entry (used when previewing an edit).
	GetAccountWithMinRunningBalance(
		ctx context.Context,
		id string,
		fromDate time.Time,
		untilDate *time.Time,
		excludeEntryID *string,
	) (*AccountWithMinRunningBalance, error)
	// GetAccountFirstEntryDate returns the date of the account's oldest ledger entry.
	GetAccountFirstEntryDate(ctx context.Context, id string) (time.Time, error)
	GetAccountLastEntryDate(ctx context.Context, id string) (time.Time, error)
}
