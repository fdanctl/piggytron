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
	UserId    incomecategory.ID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *IncomeCategoryRepository) FindAllByUser(
	ctx context.Context,
	userId incomecategory.ID,
) ([]*incomecategory.IncomeCategory, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, name, created_at, updated_at
		 FROM income_categories
		 WHERE user_id = $1`,
		userId,
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
			&c.UserId,
			&c.Name,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}

		ec := incomecategory.Rehydrate(
			c.ID,
			c.UserId,
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
