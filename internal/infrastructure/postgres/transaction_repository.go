package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/transaction"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

type TransactionDto struct {
	id     transaction.ID
	userId transaction.ID

	ttype transaction.Ttype

	fromAccountId *transaction.ID
	toAccountId   *transaction.ID

	incomeCategoryId  *transaction.ID
	expenseCategoryId *transaction.ID

	amount      uint
	description string
	date        time.Time
	createdAt   time.Time
}

// save

func (r *TransactionRepository) FindById(
	ctx context.Context,
	id transaction.ID,
) (*transaction.Transaction, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, type, from_account_id, to_account_id, income_category_id, expense_category_id, amount, description, date, created_at
		 FROM transactions
		 WHERE id = $1`,
		id,
	)

	var dto TransactionDto
	err := row.Scan(
		&dto.id,
		&dto.userId,
		&dto.ttype,
		&dto.fromAccountId,
		&dto.toAccountId,
		&dto.incomeCategoryId,
		&dto.expenseCategoryId,
		&dto.amount,
		&dto.description,
		&dto.date,
		&dto.createdAt,
	)
	if err != nil {
		return nil, err
	}

	transaction := transaction.Rehydrate(
		dto.id,
		dto.userId,
		dto.ttype,
		dto.fromAccountId,
		dto.toAccountId,
		dto.incomeCategoryId,
		dto.expenseCategoryId,
		dto.amount,
		dto.description,
		dto.date,
		dto.createdAt,
	)
	return transaction, nil
}

func (r *TransactionRepository) FindAllByUser(
	ctx context.Context,
	uid transaction.ID,
) ([]*transaction.Transaction, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, type, from_account_id, to_account_id, income_category_id, expense_category_id, amount, description, date, created_at
		 FROM transactions
		 WHERE user_id = $1`,
		uid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*transaction.Transaction

	for rows.Next() {
		var dto TransactionDto
		err := rows.Scan(
			&dto.id,
			&dto.userId,
			&dto.ttype,
			&dto.fromAccountId,
			&dto.toAccountId,
			&dto.incomeCategoryId,
			&dto.expenseCategoryId,
			&dto.amount,
			&dto.description,
			&dto.date,
			&dto.createdAt,
		)
		if err != nil {
			return nil, err
		}
		transaction := transaction.Rehydrate(
			dto.id,
			dto.userId,
			dto.ttype,
			dto.fromAccountId,
			dto.toAccountId,
			dto.incomeCategoryId,
			dto.expenseCategoryId,
			dto.amount,
			dto.description,
			dto.date,
			dto.createdAt,
		)
		transactions = append(transactions, transaction)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
