package query

import (
	"context"
	"time"
)

type AccountIDName struct {
	ID   string
	Name string
}

type AccountWithCategory struct {
	ID       string
	UserID   string
	Type     string
	Name     string
	Currency string
	// goal-specific
	TargetAmount *int
	TargetDate   *time.Time
	Category     *CategoryNameDTO

	CreatedAt time.Time
	UpdatedAt time.Time
}

type AccountWithSum struct {
	AccountWithCategory
	Sum int
}

type AccountQueryService interface {
	FindIDNamesIncludes(ctx context.Context, ids []string) ([]AccountIDName, error)
	FindGoalsIDNames(ctx context.Context, uid string) ([]AccountIDName, error)
	FindAllGoalsWithSum(ctx context.Context, uid string) ([]AccountWithSum, error)
	FindOneWithSum(ctx context.Context, id string) (AccountWithSum, error)
}
