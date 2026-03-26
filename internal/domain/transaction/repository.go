package transaction

import "context"

type Repository interface {
	// Save(ctx context.Context, category *Transaction) error
	FindByID(ctx context.Context, id ID) (*Transaction, error)
	FindAllByUser(ctx context.Context, uid ID) ([]*Transaction, error)
	FindAllByAccount(ctx context.Context, aid ID) ([]*Transaction, error)
	FindAllByCategory(ctx context.Context, cid ID) ([]*Transaction, error)
}
