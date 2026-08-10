package account

import "context"

// Repository persists Account aggregates; implemented by the postgres package.
type Repository interface {
	Create(ctx context.Context, acc *Account) error
	Update(ctx context.Context, acc *Account) error
	FindByID(ctx context.Context, id ID) (*Account, error)
	FindBankByNameAndUser(ctx context.Context, uid ID, name string) (*Account, error)
	FindGoalByNameAndUser(ctx context.Context, uid ID, name string) (*Account, error)
	FindAllByUser(ctx context.Context, uid ID) ([]*Account, error)
	FindAllBanksByUser(ctx context.Context, uid ID) ([]*Account, error)
	FindAllGoalsByUser(ctx context.Context, uid ID) ([]*Account, error)
}
