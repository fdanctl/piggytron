package transaction

import "context"

type Repository interface {
	// Save(ctx context.Context, category *Transaction) error
	FindById(ctx context.Context, id ID) (*Transaction, error)
	FindAllByUser(ctx context.Context, uid ID) ([]*Transaction, error)
	FindAllByCategory(ctx context.Context, cid ID) ([]*Transaction, error)
}
