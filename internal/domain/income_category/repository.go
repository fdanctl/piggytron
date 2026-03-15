package incomecategory

import "context"

type Repository interface {
	Save(ctx context.Context, category *IncomeCategory) error
	FindById(ctx context.Context, id ID) (*IncomeCategory, error)
	FindByNameAndUser(ctx context.Context, userId ID, name string) (*IncomeCategory, error)
	FindAllByUser(ctx context.Context, userId ID) ([]*IncomeCategory, error)
}
