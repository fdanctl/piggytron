package incomecategory

import "context"

type Repository interface {
	// Save(ctx context.Context, category *IncomeCategory) error
	// FindById(ctx context.Context, id ID) (*IncomeCategory, error)
	// FindByNameAndUser(ctx context.Context, uid ID, name string) (*IncomeCategory, error)
	FindAllByUser(ctx context.Context, uid ID) ([]*IncomeCategory, error)
}
