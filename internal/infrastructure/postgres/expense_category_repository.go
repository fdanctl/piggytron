package postgres

import (
	"context"
	"database/sql"
	"time"

	expensecategory "github.com/fdanctl/piggytron/internal/domain/expense_category"
)

type ExpenseCategoryRepository struct {
	db *sql.DB
}

func NewExpenseCategoryRepository(db *sql.DB) *ExpenseCategoryRepository {
	return &ExpenseCategoryRepository{
		db: db,
	}
}

type ExpenseCategoryDto struct {
	ID          expensecategory.ID
	UserId      expensecategory.ID
	Name        string
	ExpenseType expensecategory.ExpenseType
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r *ExpenseCategoryRepository) FindAllByUser(
	ctx context.Context,
	userId expensecategory.ID,
) ([]*expensecategory.ExpenseCategory, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, name, expense_type, created_at, updated_at
		 FROM expense_categories
		 WHERE user_id = $1`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*expensecategory.ExpenseCategory

	for rows.Next() {
		var c ExpenseCategoryDto
		if err := rows.Scan(
			&c.ID,
			&c.UserId,
			&c.Name,
			&c.ExpenseType,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}

		ec := expensecategory.Rehydrate(
			c.ID,
			c.UserId,
			c.Name,
			c.ExpenseType,
			c.CreatedAt,
			c.UpdatedAt,
		)
		categories = append(categories, ec)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, err
}
