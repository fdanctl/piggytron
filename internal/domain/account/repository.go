package account

import "context"

type AccountIdName struct {
	Id   ID
	Name string
}

type Repository interface {
	Save(ctx context.Context, category *Account) error
	FindById(ctx context.Context, id ID) (*Account, error)
	FindBankByNameAndUser(ctx context.Context, uid ID, name string) (*Account, error)
	FindGoalByNameAndUser(ctx context.Context, uid ID, name string) (*Account, error)
	FindAllByUser(ctx context.Context, uid ID) ([]*Account, error)
	FindAllBanksByUser(ctx context.Context, uid ID) ([]*Account, error)
	FindAllGoalsByUser(ctx context.Context, uid ID) ([]*Account, error)
	FindIdNamesIncludes(ctx context.Context, ids []string) ([]*AccountIdName, error)
}
