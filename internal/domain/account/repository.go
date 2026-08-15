package account

import "context"

// Repository persists Account aggregates; implemented by the postgres package.
type Repository interface {
	Create(ctx context.Context, acc *Account) error
	Update(ctx context.Context, acc *Account) error
	UpdateStatus(ctx context.Context, acc *Account) error
	Delete(ctx context.Context, id ID) error
	FindByID(ctx context.Context, id ID) (*Account, error)
	FindAllByUser(ctx context.Context, uid ID) ([]*Account, error)
	FindAllBanksByUser(ctx context.Context, uid ID) ([]*Account, error)
	FindAllGoalsByUser(ctx context.Context, uid ID) ([]*Account, error)
}
