package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/bank"
)

type BankRepository struct {
	db *sql.DB
}

func NewBankRepository(db *sql.DB) *BankRepository {
	return &BankRepository{
		db: db,
	}
}

type BankDto struct {
	ID        bank.ID
	UserId    bank.ID
	Name      string
	Currency  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *BankRepository) Save(ctx context.Context, bank *bank.Bank) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO accounts (id, user_id, type, name, currency, created_at, updated_at)
	 	 VALUES($1,$2,$3,$4,$5,$6,$7)`,
		bank.ID(),
		bank.UserId(),
		"bank",
		bank.Name(),
		bank.Currency(),
		bank.CreatedAt(),
		bank.UpdatedAt(),
	)
	return err
}

func (r *BankRepository) FindById(ctx context.Context, id bank.ID) (*bank.Bank, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, name, currency, created_at, updated_at
		 FROM accounts
		 WHERE id = $1`,
		id,
	)

	var b BankDto
	err := row.Scan(
		&b.ID,
		&b.UserId,
		&b.Name,
		&b.Currency,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	bank := bank.Rehydrate(
		b.ID,
		b.UserId,
		b.Name,
		b.Currency,
		b.CreatedAt,
		b.UpdatedAt,
	)
	return bank, err
}

func (r *BankRepository) FindByNameAndUser(
	ctx context.Context,
	uid bank.ID,
	name string,
) (*bank.Bank, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, name, currency, created_at, updated_at
		 FROM accounts
		 WHERE user_id = $1 AND name = $2`,
		uid,
		name,
	)

	var b BankDto
	err := row.Scan(
		&b.ID,
		&b.UserId,
		&b.Name,
		&b.Currency,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	bank := bank.Rehydrate(
		b.ID,
		b.UserId,
		b.Name,
		b.Currency,
		b.CreatedAt,
		b.UpdatedAt,
	)
	return bank, err
}

func (r *BankRepository) FindAllByUser(ctx context.Context, uid bank.ID) ([]*bank.Bank, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, name, currency, created_at, updated_at
		 FROM accounts
		 WHERE user_id = $1 AND type = $2`,
		uid,
		"bank",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var banks []*bank.Bank

	for rows.Next() {
		var dto BankDto
		if err := rows.Scan(
			&dto.ID,
			&dto.UserId,
			&dto.Name,
			&dto.Currency,
			&dto.CreatedAt,
			&dto.UpdatedAt,
		); err != nil {
			return nil, err
		}
		b := bank.Rehydrate(
			dto.ID,
			dto.UserId,
			dto.Name,
			dto.Currency,
			dto.CreatedAt,
			dto.UpdatedAt,
		)
		banks = append(banks, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return banks, nil
}
