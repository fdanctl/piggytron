package incomecategory

import "context"

type Repository interface {
	Create(ctx context.Context, category *IncomeCategory) error
	// Update(ctx context.Context, budget *IncomeCategory) error
	FindByID(ctx context.Context, id ID) (*IncomeCategory, error)
	FindByNameAndUser(ctx context.Context, userID ID, name string) (*IncomeCategory, error)
	FindAllByUser(ctx context.Context, userID ID) ([]*IncomeCategory, error)
}
