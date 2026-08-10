package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/monthlysummary"
)

// MonthlySummaryRepository persists monthly_summary rows via raw SQL.
type MonthlySummaryRepository struct {
	db DBTX
}

// NewMonthlySummaryRepository builds the repository over a DBTX (sql.DB or
// sql.Tx).
func NewMonthlySummaryRepository(db DBTX) *MonthlySummaryRepository {
	return &MonthlySummaryRepository{
		db: db,
	}
}

// monthlySummaryDTO is the raw monthly_summary row projection used when
// scanning query results.
type monthlySummaryDTO struct {
	accountID monthlysummary.ID
	month     time.Time
	moneyIn   int
	moneyOut  int
	createdAt time.Time
	updatedAt time.Time
}

// Save upserts a summary, adding the incoming money_in/money_out to the
// existing row on conflict (account_id, month).
func (r *MonthlySummaryRepository) Save(
	ctx context.Context,
	summary *monthlysummary.MonthlySummary,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO monthly_summary (
			account_id,
			month,
			money_in,
			money_out,
			created_at,
			updated_at
		 )
		 VALUES($1,$2,$3,$4,$5,$6)
		 ON CONFLICT(account_id, month)
		 DO UPDATE SET
			money_in = monthly_summary.money_in + EXCLUDED.money_in,
			money_out = monthly_summary.money_out + EXCLUDED.money_out,
			updated_at = EXCLUDED.updated_at`,
		summary.AccountID(),
		summary.Month().Time(),
		summary.MoneyIn(),
		summary.MoneyOut(),
		summary.CreatedAt(),
		summary.UpdatedAt(),
	)
	if err != nil {
		return err
	}

	return nil
}

// Update overwrites the money in/out totals of an existing row.
func (r *MonthlySummaryRepository) Update(
	ctx context.Context,
	summary *monthlysummary.MonthlySummary,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE monthly_summary 
		SET
			money_in = $3,
			money_out = $4,
			updated_at = $5
		WHERE account_id = $1 AND month = $2`,
		summary.AccountID(),
		summary.Month().Time(),
		summary.MoneyIn(),
		summary.MoneyOut(),
		summary.UpdatedAt(),
	)
	if err != nil {
		return err
	}

	return nil
}

// FindByAccountAndMonth loads one summary, mapping missing rows to
// monthlysummary.ErrNotFound.
func (r *MonthlySummaryRepository) FindByAccountAndMonth(
	ctx context.Context,
	accountID string,
	month monthlysummary.Month,
) (*monthlysummary.MonthlySummary, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT 
			account_id,
			month,
			money_in,
			money_out,
			created_at,
			updated_at
		 FROM monthly_summary
		 WHERE account_id = $1 AND month = $2`,
		accountID,
		month.Time(),
	)

	var (
		aid       monthlysummary.ID
		date      time.Time
		moneyIn   int
		moneyOut  int
		createdAt time.Time
		updatedAt time.Time
	)

	err := row.Scan(
		&aid,
		&date,
		&moneyIn,
		&moneyOut,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, monthlysummary.ErrNotFound
		}
		return nil, err
	}

	u := monthlysummary.Rehydrate(aid, date, moneyIn, moneyOut, createdAt, updatedAt)
	return u, err
}

// FindAllByAccount returns every summary row of an account.
func (r *MonthlySummaryRepository) FindAllByAccount(
	ctx context.Context,
	accountID string,
) ([]*monthlysummary.MonthlySummary, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT
			account_id,
			month,
			money_in,
			money_out,
			created_at,
			updated_at
		 FROM monthly_summary
		 WHERE account_id = $1`,
		accountID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []*monthlysummary.MonthlySummary

	for rows.Next() {
		var dto monthlySummaryDTO
		err := rows.Scan(
			&dto.accountID,
			&dto.month,
			&dto.moneyIn,
			&dto.moneyOut,
			&dto.createdAt,
			&dto.updatedAt,
		)
		if err != nil {
			return nil, err
		}
		ms := monthlysummary.Rehydrate(
			dto.accountID,
			dto.month,
			dto.moneyIn,
			dto.moneyOut,
			dto.createdAt,
			dto.updatedAt,
		)
		summaries = append(summaries, ms)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return summaries, nil
}

// FindAllByUser returns the summary rows of every account owned by the user.
func (r *MonthlySummaryRepository) FindAllByUser(
	ctx context.Context,
	userID string,
) ([]*monthlysummary.MonthlySummary, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT
			s.account_id,
			s.month,
			s.money_in,
			s.money_out,
			s.created_at,
			s.updated_at
		 FROM monthly_summary s
		 JOIN accounts a ON s.account_id = a.id
		 WHERE a.user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []*monthlysummary.MonthlySummary

	for rows.Next() {
		var dto monthlySummaryDTO
		err := rows.Scan(
			&dto.accountID,
			&dto.month,
			&dto.moneyIn,
			&dto.moneyOut,
			&dto.createdAt,
			&dto.updatedAt,
		)
		if err != nil {
			return nil, err
		}
		ms := monthlysummary.Rehydrate(
			dto.accountID,
			dto.month,
			dto.moneyIn,
			dto.moneyOut,
			dto.createdAt,
			dto.updatedAt,
		)
		summaries = append(summaries, ms)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return summaries, nil
}
