package budget

import "context"

type Repository interface {
	Save(ctx context.Context, budget *Budget) error
	FindByCategoryAndMonth(ctx context.Context, cid ID, month Month) (*Budget, error)
	// FindAllByUser(ctx context.Context, uid ID) ([]*Budget, error)
}
