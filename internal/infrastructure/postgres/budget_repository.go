package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/budget"
)

type BudgetRepository struct {
	db *sql.DB
}

func NewBudgetRepository(db *sql.DB) *BudgetRepository {
	return &BudgetRepository{
		db: db,
	}
}

type BudgetDto struct {
	CategoryID budget.ID
	Month      time.Time
	Amount     int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (r *BudgetRepository) Save(
	ctx context.Context,
	b *budget.Budget,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`
		INSERT INTO monthly_budgets (category_id, month, amount, created_at, updated_at)
		VALUES($1,$2,$3,$4,$5)
		ON CONFLICT(category_id, month)
		DO UPDATE SET
			amount = EXCLUDED.amount,
			updated_at = EXCLUDED.updated_at`,
		b.CategoryID(),
		b.Month().Time(),
		b.Amount(),
		b.CreatedAt(),
		b.UpdatedAt(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *BudgetRepository) FindByCategoryAndMonth(
	ctx context.Context,
	cid budget.ID,
	month budget.Month,
) (*budget.Budget, error) {
	row := r.db.QueryRowContext(
		ctx,
		`
		SELECT category_id, month, amount, created_at, updated_at
		FROM monthly_budgets
		WHERE category_id = $1 AND month = $2
		`,
		cid,
		month,
	)

	var c BudgetDto
	err := row.Scan(
		&c.CategoryID,
		&c.Month,
		&c.Amount,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, budget.ErrNotFound
		}
		return nil, err
	}
	category := budget.Rehydrate(
		c.CategoryID,
		c.Month,
		c.Amount,
		c.CreatedAt,
		c.UpdatedAt,
	)
	return category, err
}
