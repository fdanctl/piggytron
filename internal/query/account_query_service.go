package query

import (
	"context"
	"errors"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/monthlysummary"
)

var ErrNoHistory = errors.New("no history found")

type AccountIDName struct {
	ID   string
	Name string
}

type AccountWithCategory struct {
	ID       string
	UserID   string
	Type     string
	Name     string
	IsSaving *bool
	Currency string
	// goal-specific
	TargetAmount *int
	StartDate    *time.Time
	TargetDate   *time.Time
	Category     *CategoryDTO

	CreatedAt time.Time
	UpdatedAt time.Time
}

type AccountWithSum struct {
	AccountWithCategory
	Sum int
}

type AccountWithSumAndMonthChange struct {
	AccountWithSum
	MoneyIn  int
	MoneyOut int
}

type AccountWithMinRunningBalance struct {
	AccountWithCategory
	MinRunningBalance int
	MinDate           time.Time
}

type AccountDailyBalance struct {
	Day     time.Time
	ID      string
	Name    string
	Balance int
}

type AccountDailyBalanceWithStatsSince struct {
	Data         []AccountDailyBalance
	MoneyIn      int
	MoneyOut     int
	Transactions int
}

type AccountQueryService interface {
	FindIDNamesIncludes(ctx context.Context, ids []string) ([]AccountIDName, error)
	FindBanksIDNames(ctx context.Context, uid string) ([]AccountIDName, error)
	FindGoalsIDNames(ctx context.Context, uid string) ([]AccountIDName, error)
	FindWithSum(ctx context.Context, id string) (*AccountWithSum, error)
	FindAllWithSum(ctx context.Context, uid string) ([]AccountWithSum, error)
	FindAllWithSumAndMonthChange(
		ctx context.Context,
		uid string,
		month monthlysummary.Month,
	) ([]AccountWithSumAndMonthChange, error)
	FindAllGoalsWithSum(ctx context.Context, uid string) ([]AccountWithSum, error)
	GetBanksDailyBalanceSince(
		ctx context.Context,
		uid string,
		since time.Time,
	) ([]AccountDailyBalance, error)
	GetAccountDailyBalanceSince(
		ctx context.Context,
		id string,
		since time.Time,
	) ([]AccountDailyBalance, error)
	GetAccountDailyBalanceAndStatsSince(
		ctx context.Context,
		id string,
		since time.Time,
	) (*AccountDailyBalanceWithStatsSince, error)
	GetAccountWithMinRunningBalance(
		ctx context.Context,
		id string,
		fromDate time.Time,
		untilDate *time.Time,
		excludeEntryID *string,
	) (*AccountWithMinRunningBalance, error)
}
