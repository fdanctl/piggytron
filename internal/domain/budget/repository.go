package budget

import "context"

type Repository interface {
	Create(ctx context.Context, budget *Budget) error
	// Update(ctx context.Context, budget *Budget) error
	UpdateAmount(ctx context.Context, id ID, amount int) error
	FindByID(ctx context.Context, id ID) (*Budget, error)
	// FindAllByUser(ctx context.Context, uid ID) ([]*Budget, error)
}
