package transaction

import "context"

type Repository interface {
	// Save(ctx context.Context, category *Transaction) error
	FindById(ctx context.Context, id ID) (*Transaction, error)
	FindAllByUser(ctx context.Context, uid ID, limit, offset uint) ([]*Transaction, error)
	FindAllByCategory(ctx context.Context, cid ID, limit, offset uint) ([]*Transaction, error)
	FindWithFilters(
		ctx context.Context,
		uid ID,
		filters *Filters,
		limit, offset uint,
	) ([]*Transaction, error)
}
