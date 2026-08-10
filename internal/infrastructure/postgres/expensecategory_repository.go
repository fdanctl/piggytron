package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/expensecategory"
	"github.com/lib/pq"
)

// ExpenseCategoryRepository persists expense category aggregates via raw SQL.
type ExpenseCategoryRepository struct {
	db *sql.DB
}

// NewExpenseCategoryRepository builds the repository over a *sql.DB.
func NewExpenseCategoryRepository(db *sql.DB) *ExpenseCategoryRepository {
	return &ExpenseCategoryRepository{
		db: db,
	}
}

// expenseCategoryDto is the database row shape of the expense_categories
// table.
type expenseCategoryDto struct {
	ID          expensecategory.ID
	UserID      expensecategory.ID
	Name        string
	ExpenseType expensecategory.ExpenseType
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Create inserts the category, mapping unique-name violations to
// expensecategory.ErrDuplicate.
func (r *ExpenseCategoryRepository) Create(
	ctx context.Context,
	category *expensecategory.ExpenseCategory,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO expense_categories (id, user_id, name, type, created_at, updated_at)
		 VALUES($1,$2,$3,$4,$5,$6)`,
		category.ID(),
		category.UserID(),
		category.Name(),
		category.ExpenseType(),
		category.CreatedAt(),
		category.UpdatedAt(),
	)
	if err != nil {
		var pqErr *pq.Error

		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				return expensecategory.ErrDuplicate
			}
		}

		return err
	}

	return nil
}

// FindByID loads one category, mapping missing rows to
// expensecategory.ErrNotFound.
func (r *ExpenseCategoryRepository) FindByID(
	ctx context.Context,
	id expensecategory.ID,
) (*expensecategory.ExpenseCategory, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, name, type, created_at, updated_at
		 FROM expense_categories
		 WHERE id = $1`,
		id,
	)

	var c expenseCategoryDto
	err := row.Scan(
		&c.ID,
		&c.UserID,
		&c.Name,
		&c.ExpenseType,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, expensecategory.ErrNotFound
		}
		return nil, err
	}
	category := expensecategory.Rehydrate(
		c.ID,
		c.UserID,
		c.Name,
		c.ExpenseType,
		c.CreatedAt,
		c.UpdatedAt,
	)
	return category, err
}

// FindByNameAndUser looks up a category by its (user, name) uniqueness key.
// TODO: remove, not used
func (r *ExpenseCategoryRepository) FindByNameAndUser(
	ctx context.Context,
	userID expensecategory.ID,
	name string,
) (*expensecategory.ExpenseCategory, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, name, type, created_at, updated_at
		 FROM expense_categories
		 WHERE user_id = $1 AND name = $2`,
		userID,
		name,
	)

	var c expenseCategoryDto
	err := row.Scan(
		&c.ID,
		&c.UserID,
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
		c.UserID,
		c.Name,
		c.ExpenseType,
		c.CreatedAt,
		c.UpdatedAt,
	)
	return category, err
}

// FindAllByUser returns every expense category of the user.
func (r *ExpenseCategoryRepository) FindAllByUser(
	ctx context.Context,
	userID expensecategory.ID,
) ([]*expensecategory.ExpenseCategory, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, name, type, created_at, updated_at
		 FROM expense_categories
		 WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*expensecategory.ExpenseCategory

	for rows.Next() {
		var c expenseCategoryDto
		if err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.Name,
			&c.ExpenseType,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}

		ec := expensecategory.Rehydrate(
			c.ID,
			c.UserID,
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
