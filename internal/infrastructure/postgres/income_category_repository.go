package postgres

import (
	"context"
	"database/sql"
	"time"

	incomecategory "github.com/fdanctl/piggytron/internal/domain/income_category"
)

type IncomeCategoryRepository struct {
	db *sql.DB
}

func NewIncomeCategoryRepository(db *sql.DB) *IncomeCategoryRepository {
	return &IncomeCategoryRepository{
		db: db,
	}
}

type IncomeCategoryDto struct {
	ID        incomecategory.ID
	UserID    incomecategory.ID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *IncomeCategoryRepository) Save(
	ctx context.Context,
	category *incomecategory.IncomeCategory,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO income_categories (id, user_id, name, created_at, updated_at)
		 VALUES($1,$2,$3,$4,$5)`,
		category.ID(),
		category.UserID(),
		category.Name(),
		category.CreatedAt(),
		category.UpdatedAt(),
	)
	return err
}

func (r *IncomeCategoryRepository) FindByID(
	ctx context.Context,
	id incomecategory.ID,
) (*incomecategory.IncomeCategory, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, name, created_at, updated_at
		 FROM income_categories
		 WHERE id = $1`,
		id,
	)

	var c IncomeCategoryDto
	err := row.Scan(
		&c.ID,
		&c.UserID,
		&c.Name,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	category := incomecategory.Rehydrate(
		c.ID,
		c.UserID,
		c.Name,
		c.CreatedAt,
		c.CreatedAt,
	)
	return category, err
}

func (r *IncomeCategoryRepository) FindByNameAndUser(
	ctx context.Context,
	userID incomecategory.ID,
	name string,
) (*incomecategory.IncomeCategory, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, name, created_at, updated_at
		 FROM income_categories
		 WHERE user_id = $1 AND name = $2`,
		userID,
		name,
	)

	var c IncomeCategoryDto
	err := row.Scan(
		&c.ID,
		&c.UserID,
		&c.Name,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	category := incomecategory.Rehydrate(
		c.ID,
		c.UserID,
		c.Name,
		c.CreatedAt,
		c.CreatedAt,
	)
	return category, err
}

func (r *IncomeCategoryRepository) FindAllByUser(
	ctx context.Context,
	userID incomecategory.ID,
) ([]*incomecategory.IncomeCategory, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, name, created_at, updated_at
		 FROM income_categories
		 WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*incomecategory.IncomeCategory

	for rows.Next() {
		var c IncomeCategoryDto
		if err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.Name,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}

		ec := incomecategory.Rehydrate(
			c.ID,
			c.UserID,
			c.Name,
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
