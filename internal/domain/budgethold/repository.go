package budgethold

import "context"

// Repository persists Budget aggregates; implemented by the postgres package.
type Repository interface {
	Save(ctx context.Context, budget *BudgetHold) error
	FindByUserAndMonth(ctx context.Context, uid ID, month Month) (*BudgetHold, error)
}
