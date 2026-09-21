package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/budgethold"
)

// BudgetHoldRepository persists monthly budgets via raw SQL.
type BudgetHoldRepository struct {
	db DBTX
}

// NewBudgetHoldRepository builds a budget repository over a *sql.DB.
func NewBudgetHoldRepository(db DBTX) *BudgetHoldRepository {
	return &BudgetHoldRepository{
		db: db,
	}
}

// budgetDto is the database row shape of the monthly_budgets table.
type budgetHoldDto struct {
	UserID    budgethold.ID
	Month     time.Time
	Amount    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Save upserts a budget on (category_id, month).
func (r *BudgetHoldRepository) Save(
	ctx context.Context,
	b *budgethold.BudgetHold,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`
		INSERT INTO budget_holds (user_id, month, amount, created_at, updated_at)
		VALUES($1,$2,$3,$4,$5)
		ON CONFLICT(user_id, month)
		DO UPDATE SET
			amount = EXCLUDED.amount,
			updated_at = EXCLUDED.updated_at`,
		b.UserID(),
		b.Month().Time(),
		b.Amount(),
		b.CreatedAt(),
		b.UpdatedAt(),
	)
	return err
}

// FindByUserAndMonth loads a budget, mapping missing rows to
// budgethold.ErrNotFound.
func (r *BudgetHoldRepository) FindByUserAndMonth(
	ctx context.Context,
	uid budgethold.ID,
	month budgethold.Month,
) (*budgethold.BudgetHold, error) {
	row := r.db.QueryRowContext(
		ctx,
		`
		SELECT user_id, month, amount, created_at, updated_at
		FROM budget_holds
		WHERE user_id = $1 AND month = $2
		`,
		uid,
		month.Time(),
	)

	var c budgetHoldDto
	err := row.Scan(
		&c.UserID,
		&c.Month,
		&c.Amount,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, budgethold.ErrNotFound
		}
		return nil, err
	}
	hold := budgethold.Rehydrate(
		c.UserID,
		c.Month,
		c.Amount,
		c.CreatedAt,
		c.UpdatedAt,
	)
	return hold, err
}
