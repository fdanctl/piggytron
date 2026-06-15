package account

import "context"

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
