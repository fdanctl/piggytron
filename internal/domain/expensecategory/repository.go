package expensecategory

import "context"

// Repository persists ExpenseCategory aggregates; implemented by the
// postgres package.
type Repository interface {
	Create(ctx context.Context, category *ExpenseCategory) error
	Update(ctx context.Context, category *ExpenseCategory) error
	Delete(ctx context.Context, id ID) error
	Archive(ctx context.Context, id ID) error
	FindByID(ctx context.Context, id ID) (*ExpenseCategory, error)
	FindAllByUser(ctx context.Context, uid ID) ([]*ExpenseCategory, error)
}
