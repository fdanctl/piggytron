package transaction

import "context"

type Repository interface {
	Create(ctx context.Context, transaction *Transaction) error
	UpdateMany(ctx context.Context, transactions []*Transaction) error
	Delete(ctx context.Context, id ID) error
	FindByID(ctx context.Context, id ID) (*Transaction, error)
	FindAllByUser(ctx context.Context, uid ID) ([]*Transaction, error)
	FindAllByAccount(ctx context.Context, aid ID) ([]*Transaction, error)
	FindAllByCategory(ctx context.Context, cid ID) ([]*Transaction, error)
}
