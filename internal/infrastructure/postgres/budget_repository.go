package postgres

import (
	"context"
	"database/sql"
	"time"

	budget "github.com/fdanctl/piggytron/internal/domain/monthly_budget"
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
	ID         budget.ID
	UserID     budget.ID
	CategoryID budget.ID
	Month      time.Time
	Amount     int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (r *BudgetRepository) Create(
	ctx context.Context,
	budget *budget.Budget,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`
		INSERT INTO monthly_budgets (id, user_id, category_id, month, amount, created_at, updated_at)
		VALUES($1,$2,$3,$4,$5,$6,$7)
		`,
		budget.ID(),
		budget.UserID(),
		budget.CategoryID(),
		budget.Month(),
		budget.Amount(),
		budget.CreatedAt(),
		budget.UpdatedAt(),
	)
	return err
}

// func (r *BudgetRepository) Update(
// 	ctx context.Context,
// 	budget *budget.Budget,
// ) error {
// 	_, err := r.db.ExecContext(
// 		ctx,
// 		`
//         UPDATE monthly_budgets
//         SET
//             month = $1,
//             amount = $2,
//             updated_at = $3
//         WHERE id = $4
//         `,
// 		budget.Month(),
// 		budget.Amount(),
// 		budget.UpdatedAt(),
// 		budget.ID(),
// 	)
//
// 	return err
// }

func (r *BudgetRepository) UpdateAmount(
	ctx context.Context,
	id budget.ID,
	amount int,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`
        UPDATE monthly_budgets
        SET
            amount = $1,
            updated_at = $2
        WHERE id = $3
        `,
		amount,
		time.Now(),
		id,
	)

	return err
}

func (r *BudgetRepository) FindByID(
	ctx context.Context,
	id budget.ID,
) (*budget.Budget, error) {
	row := r.db.QueryRowContext(
		ctx,
		`
		SELECT id, user_id, category_id, month, amount, created_at, updated_at
		FROM monthly_budgets
		WHERE id = $1
		`,
		id,
	)

	var c BudgetDto
	err := row.Scan(
		&c.ID,
		&c.UserID,
		&c.CategoryID,
		&c.Month,
		&c.Amount,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	category := budget.Rehydrate(
		c.ID,
		c.UserID,
		c.CategoryID,
		c.Month,
		c.Amount,
		c.CreatedAt,
		c.CreatedAt,
	)
	return category, err
}
