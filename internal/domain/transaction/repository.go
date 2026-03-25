package transaction

import "context"

type Repository interface {
	// Save(ctx context.Context, category *Transaction) error
	FindById(ctx context.Context, id ID) (*Transaction, error)
	FindAllByUser(ctx context.Context, uid ID) ([]*Transaction, error)
	FindAllByAccount(ctx context.Context, aid ID) ([]*Transaction, error)
	FindAllByCategory(ctx context.Context, cid ID) ([]*Transaction, error)
	FindFiltered(
		ctx context.Context,
		uid ID,
		filters *Filters,
		limit, offset uint,
	) ([]*Transaction, error)
	FindFilteredWithCount(
		ctx context.Context,
		uid ID,
		filters *Filters,
		limit, offset uint,
	) ([]*Transaction, int, error)
	CountFilteredResults(ctx context.Context, uid ID, filters *Filters) (int, error)
}
