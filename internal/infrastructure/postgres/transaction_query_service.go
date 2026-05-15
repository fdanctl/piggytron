package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/fdanctl/piggytron/internal/query"
)

type TransactionQueryService struct {
	db *sql.DB
}

func NewTransactionQueryService(db *sql.DB) *TransactionQueryService {
	return &TransactionQueryService{
		db: db,
	}
}

func (r *TransactionQueryService) FindFiltered(
	ctx context.Context,
	uid string,
	filters *query.TransactionFilters,
	limit, offset uint,
) ([]query.TransactionDTO, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT
			t.id,
			t.user_id,
			t.type,
			fa.name,
			ta.name,
			ic.name,
			ec.name,
			t.amount,
			t.description,
			t.date,
			t.created_at
		 FROM transactions t
		 LEFT JOIN accounts fa
			ON t.from_account_id = fa.id
		 LEFT JOIN accounts ta
			ON t.to_account_id = ta.id
		 LEFT JOIN expense_categories ec
			ON t.expense_category_id = ec.id
		 LEFT JOIN income_categories ic
			ON t.income_category_id = ic.id
		 WHERE t.user_id = $1
		 AND ($2::TEXT[] IS NULL OR t.type = ANY($2))
		 AND ($3::UUID[] IS NULL OR t.from_account_id = ANY($3) OR t.to_account_id = ANY($3))
		 AND ($4::UUID[] IS NULL OR t.income_category_id = ANY($4) OR t.expense_category_id = ANY($4))
		 AND ($5::BIGINT IS NULL OR t.amount >= $5)
		 AND ($6::BIGINT IS NULL OR t.amount <= $6)
	     AND ($7::TIMESTAMP IS NULL OR t.date >= $7)
	     AND ($8::TIMESTAMP IS NULL OR t.date < $8)
		 ORDER BY date DESC
		 LIMIT NULLIF($9, 0)
		 OFFSET $10`,
		uid,
		filters.Types,
		filters.AccountIDs,
		filters.CategoryIDs,
		filters.MinAmount,
		filters.MaxAmount,
		filters.MinDate,
		filters.MaxDate,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []query.TransactionDTO

	for rows.Next() {
		var dto query.TransactionDTO
		err := rows.Scan(
			&dto.ID,
			&dto.UserID,
			&dto.Type,
			&dto.FromAccount,
			&dto.ToAccount,
			&dto.IncomeCategory,
			&dto.ExpenseCategory,
			&dto.Amount,
			&dto.Description,
			&dto.Date,
			&dto.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, dto)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *TransactionQueryService) FindFilteredWithCount(
	ctx context.Context,
	uid string,
	filters *query.TransactionFilters,
	limit, offset uint,
) (*query.TransactionsWithTotalCount, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT
			t.id,
			t.user_id,
			t.type,
			fa.name,
			ta.name,
			ic.name,
			ec.name,
			t.amount,
			t.description,
			t.date,
			t.created_at,
			COUNT(*) OVER() AS total_count
		 FROM transactions t
		 LEFT JOIN accounts fa
			ON t.from_account_id = fa.id
		 LEFT JOIN accounts ta
			ON t.to_account_id = ta.id
		 LEFT JOIN expense_categories ec
			ON t.expense_category_id = ec.id
		 LEFT JOIN income_categories ic
			ON t.income_category_id = ic.id
		 WHERE t.user_id = $1
	     AND ($2::TEXT[] IS NULL OR t.type = ANY($2))
	     AND ($3::UUID[] IS NULL OR t.from_account_id = ANY($3) OR t.to_account_id = ANY($3))
	     AND ($4::UUID[] IS NULL OR t.income_category_id = ANY($4) OR t.expense_category_id = ANY($4))
	     AND ($5::BIGINT IS NULL OR t.amount >= $5)
	     AND ($6::BIGINT IS NULL OR t.amount <= $6)
	     AND ($7::TIMESTAMP IS NULL OR t.date >= $7)
	     AND ($8::TIMESTAMP IS NULL OR t.date < $8)
		 ORDER BY date DESC
		 LIMIT NULLIF($9, 0)
		 OFFSET $10`,
		uid,
		filters.Types,
		filters.AccountIDs,
		filters.CategoryIDs,
		filters.MinAmount,
		filters.MaxAmount,
		filters.MinDate,
		filters.MaxDate,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []query.TransactionDTO
	var totalCount int

	for rows.Next() {
		var dto query.TransactionDTO
		var count int
		err := rows.Scan(
			&dto.ID,
			&dto.UserID,
			&dto.Type,
			&dto.FromAccount,
			&dto.ToAccount,
			&dto.IncomeCategory,
			&dto.ExpenseCategory,
			&dto.Amount,
			&dto.Description,
			&dto.Date,
			&dto.CreatedAt,
			&count,
		)
		if err != nil {
			return nil, err
		}
		totalCount = count
		transactions = append(transactions, dto)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &query.TransactionsWithTotalCount{
		Data:  transactions,
		Total: totalCount,
	}, nil
}

func (r *TransactionQueryService) CountFilteredResults(
	ctx context.Context,
	uid string,
	filters *query.TransactionFilters,
) (int, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*)
		 FROM transactions
		 WHERE user_id = $1
	     AND ($2::TEXT[] IS NULL OR type = ANY($2))
	     AND ($3::UUID[] IS NULL OR from_account_id = ANY($3) OR to_account_id = ANY($3))
	     AND ($4::UUID[] IS NULL OR income_category_id = ANY($4) OR expense_category_id = ANY($4))
	     AND ($5::BIGINT IS NULL OR amount >= $5)
	     AND ($6::BIGINT IS NULL OR amount <= $6)`,
		// TODO add min max date
		uid,
		filters.Types,
		filters.AccountIDs,
		filters.CategoryIDs,
		filters.MinAmount,
		filters.MaxAmount,
	)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *TransactionQueryService) GetExpensesByCategoryBetweenDates(
	ctx context.Context,
	uid string,
	minDate, maxDate time.Time,
) (*query.CategoryExpenseWithTotal, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT
			COALESCE(expense_category_id, '00000000-0000-0000-0000-000000000000') AS expense_category_id,
			COALESCE(SUM(amount), 0) AS total
		 FROM
  			transactions
		 WHERE
  			type = 'expense'
  			AND date >= $1
  			AND date < $2
		 GROUP BY
  		 ROLLUP (expense_category_id)`,
		minDate,
		maxDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var catExpenses []query.CategoryExpense
	var totalExpense int

	for rows.Next() {
		var dto query.CategoryExpense
		err := rows.Scan(
			&dto.ID,
			&dto.Amount,
		)
		if err != nil {
			return nil, err
		}
		if dto.ID == "00000000-0000-0000-0000-000000000000" {
			totalExpense = dto.Amount
		} else {
			catExpenses = append(catExpenses, dto)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &query.CategoryExpenseWithTotal{
		Data:  catExpenses,
		Total: totalExpense,
	}, nil
}
