package incomecategory

import "context"

// Repository persists IncomeCategory aggregates; implemented by the postgres
// package.
type Repository interface {
	Create(ctx context.Context, category *IncomeCategory) error
	Update(ctx context.Context, category *IncomeCategory) error
	Delete(ctx context.Context, id ID) error
	Archive(ctx context.Context, id ID) error
	FindByID(ctx context.Context, id ID) (*IncomeCategory, error)
	FindAllByUser(ctx context.Context, userID ID) ([]*IncomeCategory, error)
}
