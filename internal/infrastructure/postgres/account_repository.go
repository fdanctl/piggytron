package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/account"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

type AccountDto struct {
	ID     account.ID
	UserID account.ID
	Type   account.AccountType
	Name   string

	TargetAmount *int
	TargetDate   *time.Time
	CategoryID   *account.ID

	Currency  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *AccountRepository) Save(ctx context.Context, account *account.Account) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO accounts (id, user_id, type, name, currency, target_amount, target_date, category_id, created_at, updated_at)
	 	 VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		account.ID(),
		account.UserID(),
		account.Type(),
		account.Name(),
		account.Currency(),
		account.TargetAmount(),
		account.TargetDate(),
		account.CategoryID(),
		account.CreatedAt(),
		account.UpdatedAt(),
	)
	return err
}

func (r *AccountRepository) FindByID(ctx context.Context, id account.ID) (*account.Account, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, type, name, currency, target_amount, target_date, category_id, created_at, updated_at
		 FROM accounts
		 WHERE id = $1`,
		id,
	)

	var b AccountDto
	err := row.Scan(
		&b.ID,
		&b.UserID,
		&b.Type,
		&b.Name,
		&b.Currency,
		&b.TargetAmount,
		&b.TargetDate,
		&b.CategoryID,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	account := account.Rehydrate(
		b.ID,
		b.UserID,
		b.Type,
		b.Name,
		b.TargetAmount,
		b.TargetDate,
		b.CategoryID,
		b.Currency,
		b.CreatedAt,
		b.UpdatedAt,
	)
	return account, err
}

func (r *AccountRepository) FindBankByNameAndUser(
	ctx context.Context,
	uid account.ID,
	name string,
) (*account.Account, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, name, currency, created_at, updated_at
		 FROM accounts
		 WHERE user_id = $1 AND name = $2 and type = $3`,
		uid,
		name,
		"bank",
	)

	var b AccountDto
	err := row.Scan(
		&b.ID,
		&b.UserID,
		&b.Name,
		&b.Currency,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	account := account.Rehydrate(
		b.ID,
		b.UserID,
		"bank",
		b.Name,
		nil,
		nil,
		nil,
		b.Currency,
		b.CreatedAt,
		b.UpdatedAt,
	)
	return account, err
}

func (r *AccountRepository) FindGoalByNameAndUser(
	ctx context.Context,
	uid account.ID,
	name string,
) (*account.Account, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, type, name, currency, target_amount, target_date, category_id, created_at, updated_at
		 FROM accounts
		 WHERE user_id = $1 AND name = $2 and type = $3`,
		uid,
		name,
		"goal",
	)

	var b AccountDto
	err := row.Scan(
		&b.ID,
		&b.UserID,
		&b.Type,
		&b.Name,
		&b.Currency,
		&b.TargetAmount,
		&b.TargetDate,
		&b.CategoryID,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	account := account.Rehydrate(
		b.ID,
		b.UserID,
		b.Type,
		b.Name,
		b.TargetAmount,
		b.TargetDate,
		b.CategoryID,
		b.Currency,
		b.CreatedAt,
		b.UpdatedAt,
	)
	return account, err
}

func (r *AccountRepository) FindAllByUser(
	ctx context.Context,
	uid account.ID,
) ([]*account.Account, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, type, name, currency, target_amount, target_date, category_id, created_at, updated_at
		 FROM accounts
		 WHERE user_id = $1`,
		uid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*account.Account

	for rows.Next() {
		var dto AccountDto
		if err := rows.Scan(
			&dto.ID,
			&dto.UserID,
			&dto.Type,
			&dto.Name,
			&dto.Currency,
			&dto.TargetAmount,
			&dto.TargetDate,
			&dto.CategoryID,
			&dto.CreatedAt,
			&dto.UpdatedAt,
		); err != nil {
			return nil, err
		}
		b := account.Rehydrate(
			dto.ID,
			dto.UserID,
			dto.Type,
			dto.Name,
			dto.TargetAmount,
			dto.TargetDate,
			dto.CategoryID,
			dto.Currency,
			dto.CreatedAt,
			dto.UpdatedAt,
		)
		accounts = append(accounts, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *AccountRepository) FindAllBanksByUser(
	ctx context.Context,
	uid account.ID,
) ([]*account.Account, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, name, currency, created_at, updated_at
		 FROM accounts
		 WHERE user_id = $1 AND type = 'bank'`,
		uid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*account.Account

	for rows.Next() {
		var dto AccountDto
		if err := rows.Scan(
			&dto.ID,
			&dto.UserID,
			&dto.Name,
			&dto.Currency,
			&dto.CreatedAt,
			&dto.UpdatedAt,
		); err != nil {
			return nil, err
		}
		b := account.Rehydrate(
			dto.ID,
			dto.UserID,
			"bank",
			dto.Name,
			nil,
			nil,
			nil,
			dto.Currency,
			dto.CreatedAt,
			dto.UpdatedAt,
		)
		accounts = append(accounts, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *AccountRepository) FindAllGoalsByUser(
	ctx context.Context,
	uid account.ID,
) ([]*account.Account, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, type, name, currency, target_amount, target_date, category_id, created_at, updated_at
		 FROM accounts
		 WHERE user_id = $1 and type = 'goal'`,
		uid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*account.Account

	for rows.Next() {
		var dto AccountDto
		if err := rows.Scan(
			&dto.ID,
			&dto.UserID,
			&dto.Type,
			&dto.Name,
			&dto.Currency,
			&dto.TargetAmount,
			&dto.TargetDate,
			&dto.CategoryID,
			&dto.CreatedAt,
			&dto.UpdatedAt,
		); err != nil {
			return nil, err
		}
		b := account.Rehydrate(
			dto.ID,
			dto.UserID,
			dto.Type,
			dto.Name,
			dto.TargetAmount,
			dto.TargetDate,
			dto.CategoryID,
			dto.Currency,
			dto.CreatedAt,
			dto.UpdatedAt,
		)
		accounts = append(accounts, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}
