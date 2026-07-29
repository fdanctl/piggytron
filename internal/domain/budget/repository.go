package budget

import "context"

type Repository interface {
	Save(ctx context.Context, budget *Budget) error
	FindByCategoryAndMonth(ctx context.Context, cid ID, month Month) (*Budget, error)
	// CopyLastMonthBudget copies last month budget to this month's empty or 0 budgets,
	// and returns the number of categories updated
	CopyLastMonthBudget(ctx context.Context, uid ID, month Month) (int, error)
}
