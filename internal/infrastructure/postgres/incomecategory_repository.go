package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/incomecategory"
	"github.com/lib/pq"
)

// IncomeCategoryRepository persists income category aggregates via raw SQL.
type IncomeCategoryRepository struct {
	db DBTX
}

// NewIncomeCategoryRepository builds the repository over a *sql.DB.
func NewIncomeCategoryRepository(db DBTX) *IncomeCategoryRepository {
	return &IncomeCategoryRepository{
		db: db,
	}
}

// incomeCategoryDto is the database row shape of the income_categories
// table.
type incomeCategoryDto struct {
	ID         incomecategory.ID
	UserID     incomecategory.ID
	Name       string
	Status     incomecategory.CategoryStatus
	ArchivedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Create inserts the category, mapping unique-name violations to
// incomecategory.ErrDuplicate.
func (r *IncomeCategoryRepository) Create(
	ctx context.Context,
	category *incomecategory.IncomeCategory,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO income_categories (id, user_id, name, archived_at, created_at, updated_at)
		 VALUES($1,$2,$3,$4,$5)`,
		category.ID(),
		category.UserID(),
		category.Name(),
		category.ArchivedAt(),
		category.CreatedAt(),
		category.UpdatedAt(),
	)
	if err != nil {
		var pqErr *pq.Error

		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				return incomecategory.ErrDuplicate
			}
		}

		return err
	}

	return nil
}

// Update persists the mutable fields of an income category.
func (r *IncomeCategoryRepository) Update(
	ctx context.Context,
	category *incomecategory.IncomeCategory,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`
		UPDATE income_categories
		SET
			name = $2,
			updated_at = $3
		WHERE id = $1`,
		category.ID(),
		category.Name(),
		category.UpdatedAt(),
	)
	if err != nil {
		var pqErr *pq.Error

		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				return incomecategory.ErrDuplicate
			}
		}

		return err
	}

	return nil
}

func (r *IncomeCategoryRepository) Delete(
	ctx context.Context,
	id incomecategory.ID,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM income_categories WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *IncomeCategoryRepository) Archive(
	ctx context.Context,
	id incomecategory.ID,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE income_categories
		SET
			status = 'archived',
			archived_at = $2,
			updated_at = $2
		WHERE id = $1`,
		id,
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

// FindByID loads one category, mapping missing rows to
// incomecategory.ErrNotFound.
func (r *IncomeCategoryRepository) FindByID(
	ctx context.Context,
	id incomecategory.ID,
) (*incomecategory.IncomeCategory, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, name, status, archived_at, created_at, updated_at
		 FROM income_categories
		 WHERE id = $1`,
		id,
	)

	var c incomeCategoryDto
	err := row.Scan(
		&c.ID,
		&c.UserID,
		&c.Name,
		&c.Status,
		&c.ArchivedAt,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, incomecategory.ErrNotFound
		}
		return nil, err
	}
	category := incomecategory.Rehydrate(
		c.ID,
		c.UserID,
		c.Name,
		c.Status,
		c.ArchivedAt,
		c.CreatedAt,
		c.UpdatedAt,
	)
	return category, err
}

// FindAllByUser returns every income category of the user.
func (r *IncomeCategoryRepository) FindAllByUser(
	ctx context.Context,
	userID incomecategory.ID,
) ([]*incomecategory.IncomeCategory, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, name, status, archived_at, created_at, updated_at
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
		var c incomeCategoryDto
		if err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.Name,
			&c.Status,
			&c.ArchivedAt,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}

		ec := incomecategory.Rehydrate(
			c.ID,
			c.UserID,
			c.Name,
			c.Status,
			c.ArchivedAt,
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
