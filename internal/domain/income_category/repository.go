package incomecategory

import "context"

type Repository interface {
	Save(ctx context.Context, category *IncomeCategory) error
	FindByID(ctx context.Context, id ID) (*IncomeCategory, error)
	FindByNameAndUser(ctx context.Context, userID ID, name string) (*IncomeCategory, error)
	FindAllByUser(ctx context.Context, userID ID) ([]*IncomeCategory, error)
}
