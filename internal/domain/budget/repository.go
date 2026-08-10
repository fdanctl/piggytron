package budget

import "context"

// Repository persists Budget aggregates; implemented by the postgres package.
type Repository interface {
	Save(ctx context.Context, budget *Budget) error
	FindByCategoryAndMonth(ctx context.Context, cid ID, month Month) (*Budget, error)
	// CopyLastMonthBudget carries last month's budget over to this month's
	// empty or zero budgets, and returns the number of categories updated.
	CopyLastMonthBudget(ctx context.Context, uid ID, month Month) (int, error)
}
