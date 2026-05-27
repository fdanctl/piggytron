package query

import (
	"context"
	"errors"
	"time"
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
	Category     *CategoryNameDTO

	CreatedAt time.Time
	UpdatedAt time.Time
}

type AccountWithSum struct {
	AccountWithCategory
	Sum int
}

type AccountDailyChange struct {
	ID     string
	Name   string
	Date   time.Time
	Change int
}

type AccountQueryService interface {
	FindIDNamesIncludes(ctx context.Context, ids []string) ([]AccountIDName, error)
	FindBanksIDNames(ctx context.Context, uid string) ([]AccountIDName, error)
	FindGoalsIDNames(ctx context.Context, uid string) ([]AccountIDName, error)
	FindWithSum(ctx context.Context, id string) (*AccountWithSum, error)
	FindAllWithSum(ctx context.Context, uid string) ([]AccountWithSum, error)
	FindAllGoalsWithSum(ctx context.Context, uid string) ([]AccountWithSum, error)
	FindOneWithSum(ctx context.Context, id string) (AccountWithSum, error)
	GetBanksDailyChange(ctx context.Context, uid string) ([]AccountDailyChange, error)
	GetAccountDailyChange(ctx context.Context, id string) ([]AccountDailyChange, error)
}
