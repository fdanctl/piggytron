package expensecategory

import "context"

// Repository persists ExpenseCategory aggregates; implemented by the
// postgres package.
// TODO: update
type Repository interface {
	Create(ctx context.Context, category *ExpenseCategory) error
	FindByID(ctx context.Context, id ID) (*ExpenseCategory, error)
	FindByNameAndUser(ctx context.Context, uid ID, name string) (*ExpenseCategory, error)
	FindAllByUser(ctx context.Context, uid ID) ([]*ExpenseCategory, error)
}
