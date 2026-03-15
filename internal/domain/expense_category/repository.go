package expensecategory

import "context"

type Repository interface {
	Save(ctx context.Context, category *ExpenseCategory) error
	FindById(ctx context.Context, id ID) (*ExpenseCategory, error)
	FindByNameAndUser(ctx context.Context, uid ID, name string) (*ExpenseCategory, error)
	FindAllByUser(ctx context.Context, uid ID) ([]*ExpenseCategory, error)
}
