package bank

import "context"

type Repository interface {
	Save(ctx context.Context, category *Bank) error
	FindById(ctx context.Context, id ID) (*Bank, error)
	FindByNameAndUser(ctx context.Context, uid ID, name string) (*Bank, error)
	FindAllByUser(ctx context.Context, uid ID) ([]*Bank, error)
}
