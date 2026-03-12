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

func (r *ExpenseCategoryRepository) Save(
	ctx context.Context,
	category *expensecategory.ExpenseCategory,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO expense_categories (id, user_id, name, expense_type, created_at, updated_at)
		 VALUES($1,$2,$3,$4,$5,$6)`,
		category.ID(),
		category.UserId(),
		category.Name(),
		category.ExpenseType(),
		category.CreatedAt(),
		category.UpdatedAt(),
	)
	return err
}

func (r *ExpenseCategoryRepository) FindByNameAndUser(
	ctx context.Context,
	userId expensecategory.ID,
	name string,
) (*expensecategory.ExpenseCategory, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, name, expense_type, created_at, updated_at
		 FROM expense_categories
		 WHERE user_id = $1 AND name = $2`,
		userId,
		name,
	)

	var c ExpenseCategoryDto
	err := row.Scan(
		&c.ID,
		&c.UserId,
		&c.Name,
		&c.ExpenseType,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	category := expensecategory.Rehydrate(
		c.ID,
		c.UserId,
		c.Name,
		c.ExpenseType,
		c.CreatedAt,
		c.CreatedAt,
	)
	return category, err
}

func (r *ExpenseCategoryRepository) FindAllByUser(
	ctx context.Context,
	userId expensecategory.ID,
) ([]*expensecategory.ExpenseCategory, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, name, expense_type, created_at, updated_at
		 FROM expense_categories
		 WHERE user_id = $1
		 ORDER BY expense_type`,
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
