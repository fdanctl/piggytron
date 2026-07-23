package monthlysummary

import "context"

type Repository interface {
	Save(ctx context.Context, summary *MonthlySummary) error
	Update(ctx context.Context, summary *MonthlySummary) error
	FindByAccountAndMonth(
		ctx context.Context,
		accountID string,
		month Month,
	) (*MonthlySummary, error)
	FindAllByAccount(ctx context.Context, accountID string) ([]*MonthlySummary, error)
	FindAllByUser(ctx context.Context, userID string) ([]*MonthlySummary, error)
}
