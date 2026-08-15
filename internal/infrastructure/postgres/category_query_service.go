package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/budget"
	"github.com/fdanctl/piggytron/internal/query"
)

// CategoryQueryService implements the category read-model queries declared
// in internal/query using raw SQL.
type CategoryQueryService struct {
	db DBTX
}

// NewCategoryQueryService builds the service over a DBTX (sql.DB or sql.Tx).
func NewCategoryQueryService(db DBTX) *CategoryQueryService {
	return &CategoryQueryService{
		db: db,
	}
}

// FindByID looks a category up across both income and expense categories.
func (s *CategoryQueryService) FindByID(
	ctx context.Context,
	id string,
) (*query.CategoryDTO, error) {
	row := s.db.QueryRowContext(
		ctx,
		`SELECT id, name, 'income' AS type
		 FROM income_categories
		 WHERE id = $1
		 UNION
		 SELECT id, name, type
		 FROM expense_categories
		 WHERE id = $1`,
		id,
	)
	var c query.CategoryDTO
	err := row.Scan(
		&c.ID,
		&c.Name,
		&c.Type,
	)
	if err != nil {
		return nil, err
	}
	return &c, err
}

// FindAllCategories returns all categories (income and expense) of the user.
func (s *CategoryQueryService) FindAllCategories(
	ctx context.Context,
	uid string,
) ([]query.CategoryDTO, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, name, 'income' AS type
		 FROM income_categories
		 WHERE user_id = $1
		 UNION
		 SELECT id, name, type
		 FROM expense_categories
		 WHERE user_id = $1`,
		uid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []query.CategoryDTO

	for rows.Next() {
		var c query.CategoryDTO
		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Type,
		); err != nil {
			return nil, err
		}
		results = append(results, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// FindCategoriesIDIncludes returns id/name pairs for the given category ids,
// across both income and expense categories.
func (s *CategoryQueryService) FindCategoriesIDIncludes(
	ctx context.Context,
	ids []string,
) ([]query.CategoryNameDTO, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))

	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	queryStr := fmt.Sprintf(
		`SELECT id, name
		 FROM income_categories
		 WHERE id IN (%s)
		 UNION
		 SELECT id, name
		 FROM expense_categories
		 WHERE id IN (%s)`,
		strings.Join(placeholders, ","),
		strings.Join(placeholders, ","),
	)

	rows, err := s.db.QueryContext(
		ctx,
		queryStr,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []query.CategoryNameDTO

	for rows.Next() {
		var c query.CategoryNameDTO
		if err := rows.Scan(
			&c.ID,
			&c.Name,
		); err != nil {
			return nil, err
		}
		results = append(results, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// GetCategoriesBudgetSpentValue returns, for month, every category with its
// budget, spent value and previous totals, plus the month's net cash flow
// and the running balance of the non-savings bank accounts.
func (s *CategoryQueryService) GetCategoriesBudgetSpentValue(
	ctx context.Context,
	uid string,
	month budget.Month,
) (*query.MonthExpenseCategoryBudgetSpentWithBalance, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`
		WITH
		  month_net AS (
		    SELECT
		      COALESCE(ms.month, $2) as month,
		      COALESCE(SUM(ms.money_in - ms.money_out), 0) AS net,
		      SUM(SUM(ms.money_in - ms.money_out)) OVER (
		        ORDER BY
		          month
		      ) AS running_total
		    FROM
		      accounts a
		      LEFT JOIN monthly_summary ms ON a.id = ms.account_id
		      AND ms.month <= $2
		    WHERE
		      user_id = $1
		      AND type = 'checking'
		    GROUP BY
		      ms.month
		  ),
		  categories as (
		    SELECT
		      c.id as cid,
		      COALESCE(b.month, $2) AS month,
		      c.type,
		      c.name,
		      COALESCE(b.amount, 0) as budgeted,
		      COALESCE(SUM(t.amount), 0) as spent_income,
		      (
		        SELECT
		          COALESCE(SUM(mb.amount), 0)
		        FROM
		          monthly_budgets mb
		        WHERE
		          mb.category_id = c.id
		          AND mb.month < COALESCE(b.month, $2)
		      ) as total_budgeted_prev,
		      (
		        SELECT
		          COALESCE(SUM(l.amount), 0)
		        FROM
		          ledger l
		        WHERE
		          l.expense_category_id = c.id
		          AND l.date < COALESCE(b.month, $2)
		      ) as total_spent_prev
		    FROM
		      expense_categories c
		      LEFT JOIN ledger t ON c.id = t.expense_category_id
		      AND t.date >= $2
		      AND t.date < date_trunc('month', $2::TIMESTAMP) + INTERVAL '1 month'
		      LEFT JOIN monthly_budgets b ON c.id = b.category_id
		      AND b.month >= $2
		      AND b.month < date_trunc('month', $2::TIMESTAMP) + INTERVAL '1 month'
		    WHERE
		      c.user_id = $1
		    GROUP BY
		      c.id,
		      b.category_id,
		      b.month
		    UNION ALL
		    SELECT
		      c.id as cid,
		      $2 AS month,
		      'income' AS type,
		      c.name as name,
		      0 as budgeted,
		      COALESCE(SUM(t.amount), 0) as spent_income,
		      0 as total_budgeted_prev,
		      0 as total_spent_prev
		    FROM
		      income_categories c
		      LEFT JOIN ledger t ON c.id = t.income_category_id
		      AND t.date >= $2
		      AND t.date < date_trunc('month', $2::TIMESTAMP) + INTERVAL '1 month'
		    WHERE
		      c.user_id = $1
		    GROUP BY
		      c.id
		  )
		SELECT
		  c.*,
		  COALESCE(n.net, 0) as net,
		  COALESCE(n.running_total, 0) as balance
		FROM
		  categories c
		  LEFT JOIN month_net n ON n.month = $2
		ORDER BY c.type
		`,
		uid,
		month.Time(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []query.CategoryBudgetValue
	var monthNet int
	var balance int

	for rows.Next() {
		var r query.CategoryBudgetValue
		err := rows.Scan(
			&r.CategoryID,
			&r.Month,
			&r.Type,
			&r.Name,
			&r.Budgeted,
			&r.Value,
			&r.PrevTotalBudget,
			&r.PrevTotalSpent,
			&monthNet,
			&balance,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &query.MonthExpenseCategoryBudgetSpentWithBalance{
		Data:     data,
		MonthNet: monthNet,
		Balance:  balance,
	}, nil
}

// GetCategoriesBudgetSpent returns, for a date range, the budgeted value per
// expense category and the received income per income category.
func (s *CategoryQueryService) GetCategoriesBudgetSpent(
	ctx context.Context,
	uid string,
	minDate, maxDate time.Time,
) ([]query.CategoryBudget, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`
        SELECT
          c.name AS name,
          c.type AS category_type,
          COALESCE(b.amount, 0) AS value
        FROM
          expense_categories c
          LEFT JOIN monthly_budgets b ON c.id = b.category_id
          AND b.month >= $1
          AND b.month < $2
        WHERE
          c.user_id = $3
        GROUP BY
          c.id,
		  b.category_id,
          b.month

		UNION ALL

		SELECT
          c.name AS name,
    	  'income' AS category_type,
    	  COALESCE(SUM(t.amount), 0) AS value
		FROM
    	  income_categories c
    	  LEFT JOIN ledger t 
          ON c.id = t.income_category_id
          AND t.date >= $1
          AND t.date < $2
		WHERE
		  c.user_id = $3
		GROUP BY
    	  c.id`,
		minDate,
		maxDate,
		uid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []query.CategoryBudget

	for rows.Next() {
		var dto query.CategoryBudget
		if err := rows.Scan(
			&dto.Name,
			&dto.Type,
			&dto.Value,
		); err != nil {
			return nil, err
		}
		results = append(results, dto)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// GetYearMonthlyValue returns, for one category and year, the monthly
// totals of its ledger movements.
func (s *CategoryQueryService) GetYearMonthlyValue(
	ctx context.Context,
	year int,
	id string,
) ([]query.CategoryMonthlyValue, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT
		  ic.id,
		  ic.name,
		  EXTRACT(
			MONTH
			FROM
			  t.date
		  ) as month,
		  SUM(t.amount) AS value
		FROM
		  income_categories ic
		  LEFT JOIN ledger t ON ic.id = t.income_category_id
		WHERE
		  ic.id = $1
		  AND EXTRACT(
			YEAR
			FROM
			  t.date
		  ) = $2
		GROUP BY
		  month,
		  ic.id

		UNION ALL

		SELECT
		  ec.id,
		  ec.name,
		  EXTRACT(
			MONTH
			FROM
			  t.date
		  ) as month,
		  SUM(t.amount) AS value
		FROM
		  expense_categories ec
		  LEFT JOIN ledger t ON ec.id = t.expense_category_id
		WHERE
		  ec.id = $1
		  AND EXTRACT(
			YEAR
			FROM
			  t.date
		  ) = $2
		GROUP BY
		  month,
		  ec.id
		ORDER BY
		  month`,
		id,
		year,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, query.ErrNoHistory
		}
		return nil, err
	}
	defer rows.Close()

	var results []query.CategoryMonthlyValue

	for rows.Next() {
		var r query.CategoryMonthlyValue
		if err := rows.Scan(
			&r.ID,
			&r.Name,
			&r.Month,
			&r.Value,
		); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, nil
}
