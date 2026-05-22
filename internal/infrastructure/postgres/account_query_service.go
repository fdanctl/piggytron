package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/fdanctl/piggytron/internal/query"
)

type AccountQueryService struct {
	db *sql.DB
}

func NewAccountQueryService(db *sql.DB) *AccountQueryService {
	return &AccountQueryService{
		db: db,
	}
}

func (r *AccountQueryService) FindIDNamesIncludes(
	ctx context.Context,
	ids []string,
) ([]query.AccountIDName, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))

	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	qquery := fmt.Sprintf(
		`SELECT id, name
		 FROM accounts
		 WHERE id IN (%s)`,
		strings.Join(placeholders, ","),
	)

	rows, err := r.db.QueryContext(
		ctx,
		qquery,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []query.AccountIDName
	for rows.Next() {
		var a query.AccountIDName
		if err := rows.Scan(
			&a.ID,
			&a.Name,
		); err != nil {
			return nil, err
		}
		results = append(results, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *AccountQueryService) FindGoalsIDNames(
	ctx context.Context,
	uid string,
) ([]query.AccountIDName, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, name
		 FROM accounts
		 WHERE user_id = $1 and type = 'goal'`,
		uid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []query.AccountIDName

	for rows.Next() {
		var g query.AccountIDName
		if err := rows.Scan(
			&g.ID,
			&g.Name,
		); err != nil {
			return nil, err
		}
		results = append(results, g)
	}
	return results, nil
}

func (r *AccountQueryService) FindAllWithSum(
	ctx context.Context,
	uid string,
) ([]query.AccountWithSum, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT 
			a.id, 
			a.user_id, 
			a.type, 
			a.name, 
			a.is_saving, 
			a.currency, 
			a.target_amount, 
			a.target_date, 
			COALESCE(c.id, '00000000-0000-0000-0000-000000000000'),
			COALESCE(c.name,''),
			a.created_at, 
			a.updated_at, 
			COALESCE(
				SUM(
					CASE
					  WHEN t.from_account_id = a.id THEN t.amount * -1
					  ELSE t.amount
					END
				),
				0
			) AS sum
		 FROM accounts a
		 LEFT JOIN expense_categories c
			ON a.category_id = c.id
		 LEFT JOIN transactions t 
			ON a.id = t.to_account_id OR a.id = t.from_account_id
		 WHERE
			a.user_id = $1
		 GROUP BY
			a.id, c.id`,
		uid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []query.AccountWithSum

	for rows.Next() {
		var g query.AccountWithSum
		var c query.CategoryNameDTO
		if err := rows.Scan(
			&g.ID,
			&g.UserID,
			&g.Type,
			&g.Name,
			&g.IsSaving,
			&g.Currency,
			&g.TargetAmount,
			&g.TargetDate,
			&c.ID,
			&c.Name,
			&g.CreatedAt,
			&g.UpdatedAt,
			&g.Sum,
		); err != nil {
			return nil, err
		}
		g.Category = &c
		results = append(results, g)
	}
	return results, nil
}

func (r *AccountQueryService) FindAllGoalsWithSum(
	ctx context.Context,
	uid string,
) ([]query.AccountWithSum, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT 
			a.id, 
			a.user_id, 
			a.type, 
			a.name, 
			a.is_saving, 
			a.currency, 
			a.target_amount, 
			a.target_date, 
			c.id,
			c.name,
			a.created_at, 
			a.updated_at, 
			COALESCE(
				SUM(
					CASE
					  WHEN t.from_account_id = a.id THEN t.amount * -1
					  ELSE t.amount
					END
				),
				0
			) AS sum
		 FROM accounts a
		 LEFT JOIN expense_categories c
			ON a.category_id = c.id
		 LEFT JOIN transactions t 
			ON a.id = t.to_account_id OR a.id = t.from_account_id
		 WHERE
			a.user_id = $1 AND a.type = 'goal'
		 GROUP BY
			a.id, c.id`,
		uid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []query.AccountWithSum

	for rows.Next() {
		var g query.AccountWithSum
		var c query.CategoryNameDTO
		if err := rows.Scan(
			&g.ID,
			&g.UserID,
			&g.Type,
			&g.Name,
			&g.IsSaving,
			&g.Currency,
			&g.TargetAmount,
			&g.TargetDate,
			&c.ID,
			&c.Name,
			&g.CreatedAt,
			&g.UpdatedAt,
			&g.Sum,
		); err != nil {
			return nil, err
		}
		g.Category = &c
		results = append(results, g)
	}
	return results, nil
}

func (r *AccountQueryService) FindOneWithSum(
	ctx context.Context,
	id string,
) (query.AccountWithSum, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT
			a.id,
			a.user_id,
			a.type,
			a.name,
			a.is_saving, 
			a.currency,
			a.target_amount,
			a.target_date,
			c.id,
			c.name,
			a.created_at,
			a.updated_at,
			COALESCE(
				SUM(
					CASE
					  WHEN t.from_account_id = a.id THEN t.amount * -1
					  ELSE t.amount
					END
				),
				0
			) AS sum
		 FROM accounts a
		 LEFT JOIN expense_categories c
			ON a.category_id = c.id
		 LEFT JOIN transactions t
			ON a.id = t.to_account_id OR a.id = t.from_account_id
		 WHERE
			a.id = $1	
		 GROUP BY
			a.id, c.id`,
		id,
	)

	var a query.AccountWithSum
	var c query.CategoryNameDTO
	err := row.Scan(
		&a.ID,
		&a.UserID,
		&a.Type,
		&a.Name,
		&a.IsSaving,
		&a.Currency,
		&a.TargetAmount,
		&a.TargetDate,
		&c.ID,
		&c.Name,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.Sum,
	)
	a.Category = &c
	return a, err
}

func (r *AccountQueryService) GetBanksDailyChange(
	ctx context.Context,
	uid string,
) ([]query.AccountDailyChange, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT
			a.id,
			a.name,
			DATE(date) AS day,
			SUM(
				CASE
				  WHEN t.from_account_id = a.id THEN t.amount * -1
				  ELSE t.amount
				END
			) AS change
		 FROM accounts a
		 LEFT JOIN transactions t
			ON a.id = t.to_account_id OR a.id = t.from_account_id
		 WHERE
			a.user_id = $1 AND a.type = 'bank'
		 GROUP BY DATE(date), a.id
		 ORDER BY day`,

		uid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []query.AccountDailyChange

	for rows.Next() {
		var r query.AccountDailyChange
		if err := rows.Scan(
			&r.ID,
			&r.Name,
			&r.Date,
			&r.Change,
		); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}

func (r *AccountQueryService) GetAccountDailyChange(
	ctx context.Context,
	id string,
) ([]query.AccountDailyChange, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT
			a.id,
			a.name,
			DATE(date) AS day,
			SUM(
				CASE
				  WHEN t.from_account_id = a.id THEN t.amount * -1
				  ELSE t.amount
				END
			) AS change
		 FROM accounts a
		 LEFT JOIN transactions t
			ON a.id = t.to_account_id OR a.id = t.from_account_id
		 WHERE
			a.id = $1
		 GROUP BY DATE(date), a.id
		 ORDER BY day`,

		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []query.AccountDailyChange

	for rows.Next() {
		var r query.AccountDailyChange
		if err := rows.Scan(
			&r.ID,
			&r.Name,
			&r.Date,
			&r.Change,
		); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}
