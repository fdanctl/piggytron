package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/budget"
)

// BudgetRepository persists monthly budgets via raw SQL.
type BudgetRepository struct {
	db DBTX
}

// NewBudgetRepository builds a budget repository over a *sql.DB.
func NewBudgetRepository(db DBTX) *BudgetRepository {
	return &BudgetRepository{
		db: db,
	}
}

// budgetDto is the database row shape of the monthly_budgets table.
type budgetDto struct {
	CategoryID budget.ID
	Month      time.Time
	Amount     int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Save upserts a budget on (category_id, month).
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
	return err
}

// FindByCategoryAndMonth loads a budget, mapping missing rows to
// budget.ErrNotFound.
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
		month.Time(),
	)

	var c budgetDto
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

func (r *BudgetRepository) FindAllByCategory(
	ctx context.Context,
	cid budget.ID,
) ([]*budget.Budget, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT category_id, month, amount, created_at, updated_at
		 FROM monthly_budgets
		 WHERE category_id = $1`,
		cid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var budgets []*budget.Budget

	for rows.Next() {
		var dto budgetDto
		if err := rows.Scan(
			&dto.CategoryID,
			&dto.Month,
			&dto.Amount,
			&dto.CreatedAt,
			&dto.UpdatedAt,
		); err != nil {
			return nil, err
		}
		b := budget.Rehydrate(
			dto.CategoryID,
			dto.Month,
			dto.Amount,
			dto.CreatedAt,
			dto.UpdatedAt,
		)
		budgets = append(budgets, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return budgets, nil
}

// CopyLastMonthBudget copies last month's positive budgets into month for
// the user's categories, overwriting only budgets currently <= 0, and
// returns the number of rows affected.
func (r *BudgetRepository) CopyLastMonthBudget(
	ctx context.Context,
	uid budget.ID,
	month budget.Month,
) (int, error) {
	row := r.db.QueryRowContext(
		ctx,
		`
        WITH
        copy AS (
          INSERT INTO
            monthly_budgets (
              category_id,
              month,
              amount,
              created_at,
              updated_at
            )
          SELECT
            ms.category_id,
            $2,
            ms.amount,
            NOW(),
            NOW()
          FROM
            monthly_budgets ms
            JOIN expense_categories c ON ms.category_id = c.id
            AND ms.month = $3 AND ms.amount > 0
          WHERE
            c.user_id = $1
          ON CONFLICT (category_id, month) DO UPDATE
          SET
            amount = EXCLUDED.amount,
            updated_at = EXCLUDED.updated_at
          WHERE
            monthly_budgets.amount <= 0
          RETURNING
            1 as inserted
        )
        SELECT
          COALESCE(SUM(inserted), 0) as total_inserted
        FROM
        copy`,
		uid,
		month.Time(),
		time.Date(
			month.Time().Year(),
			month.Time().Month()-1,
			1,
			0,
			0,
			0,
			0,
			month.Time().Location(),
		),
	)

	var amountUpdated int
	err := row.Scan(
		&amountUpdated,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, budget.ErrNotFound
		}
		return 0, err
	}
	return amountUpdated, nil
}
