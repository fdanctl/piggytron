package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/monthlysummary"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/internal/util"
)

type AccountQueryService struct {
	db DBTX
}

func NewAccountQueryService(db DBTX) *AccountQueryService {
	return &AccountQueryService{
		db: db,
	}
}

func (s *AccountQueryService) FindIDNamesIncludes(
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

	rows, err := s.db.QueryContext(
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

func (s *AccountQueryService) FindBanksIDNames(
	ctx context.Context,
	uid string,
) ([]query.AccountIDName, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, name
		 FROM accounts
		 WHERE user_id = $1 and type = 'bank'`,
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

func (s *AccountQueryService) FindGoalsIDNames(
	ctx context.Context,
	uid string,
) ([]query.AccountIDName, error) {
	rows, err := s.db.QueryContext(
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

func (s *AccountQueryService) FindWithSum(
	ctx context.Context,
	id string,
) (*query.AccountWithSum, error) {
	row := s.db.QueryRowContext(
		ctx,
		`SELECT 
			a.id, 
			a.user_id, 
			a.type, 
			a.name, 
			a.is_saving, 
			a.currency, 
			a.target_amount, 
			a.start_date, 
			a.target_date, 
			COALESCE(c.id, $1),
			COALESCE(c.name,''),
			COALESCE(c.type,'income'),
			a.created_at, 
			a.updated_at, 
			COALESCE(SUM(ms.money_in - ms.money_out), 0) AS sum
		 FROM accounts a
		 LEFT JOIN expense_categories c
			ON a.category_id = c.id
		 LEFT JOIN monthly_summary ms
			ON a.id = ms.account_id
		 WHERE
			a.id = $2
		 GROUP BY
			a.id, c.id`,
		util.ZeroUUID,
		id,
	)
	var g query.AccountWithSum
	var c query.CategoryDTO
	if err := row.Scan(
		&g.ID,
		&g.UserID,
		&g.Type,
		&g.Name,
		&g.IsSaving,
		&g.Currency,
		&g.TargetAmount,
		&g.StartDate,
		&g.TargetDate,
		&c.ID,
		&c.Name,
		&c.Type,
		&g.CreatedAt,
		&g.UpdatedAt,
		&g.Sum,
	); err != nil {
		return nil, err
	}
	g.Category = &c
	return &g, nil
}

func (s *AccountQueryService) FindAllWithSum(
	ctx context.Context,
	uid string,
) ([]query.AccountWithSum, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT 
			a.id, 
			a.user_id, 
			a.type, 
			a.name, 
			a.is_saving, 
			a.currency, 
			a.target_amount, 
			a.start_date, 
			a.target_date, 
			COALESCE(c.id, $1),
			COALESCE(c.name,''),
			COALESCE(c.type,'income'),
			a.created_at, 
			a.updated_at, 
			COALESCE(SUM(ms.money_in - ms.money_out), 0) AS sum
		 FROM accounts a
		 LEFT JOIN expense_categories c
			ON a.category_id = c.id
		 LEFT JOIN monthly_summary ms
			ON a.id = ms.account_id
		 WHERE
			a.user_id = $2
		 GROUP BY
			a.id, c.id`,
		util.ZeroUUID,
		uid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []query.AccountWithSum

	for rows.Next() {
		var g query.AccountWithSum
		var c query.CategoryDTO
		if err := rows.Scan(
			&g.ID,
			&g.UserID,
			&g.Type,
			&g.Name,
			&g.IsSaving,
			&g.Currency,
			&g.TargetAmount,
			&g.StartDate,
			&g.TargetDate,
			&c.ID,
			&c.Name,
			&c.Type,
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

func (s *AccountQueryService) FindAllWithSumAndMonthChange(
	ctx context.Context,
	uid string,
	month monthlysummary.Month,
) ([]query.AccountWithSumAndMonthChange, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT 
			a.id, 
			a.user_id, 
			a.type, 
			a.name, 
			a.is_saving, 
			a.currency, 
			a.target_amount, 
			a.start_date, 
			a.target_date, 
			COALESCE(c.id, $1),
			COALESCE(c.name,''),
			COALESCE(c.type,'income'),
			a.created_at, 
			a.updated_at, 
			COALESCE(SUM(ms.money_in - ms.money_out), 0) AS sum,
			COALESCE(SUM(ms.money_in)  FILTER (WHERE ms.month = $3), 0) AS month_money_in,
			COALESCE(SUM(ms.money_out) FILTER (WHERE ms.month = $3), 0) AS month_money_out
		 FROM accounts a
		 LEFT JOIN expense_categories c
			ON a.category_id = c.id
		 LEFT JOIN monthly_summary ms
			ON a.id = ms.account_id and ms.month <= $3
		 WHERE
			a.user_id = $2
		 GROUP BY
			a.id, c.id`,
		util.ZeroUUID,
		uid,
		month.Time(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []query.AccountWithSumAndMonthChange

	for rows.Next() {
		var g query.AccountWithSumAndMonthChange
		var c query.CategoryDTO
		if err := rows.Scan(
			&g.ID,
			&g.UserID,
			&g.Type,
			&g.Name,
			&g.IsSaving,
			&g.Currency,
			&g.TargetAmount,
			&g.StartDate,
			&g.TargetDate,
			&c.ID,
			&c.Name,
			&c.Type,
			&g.CreatedAt,
			&g.UpdatedAt,
			&g.Sum,
			&g.MoneyIn,
			&g.MoneyOut,
		); err != nil {
			return nil, err
		}
		g.Category = &c
		results = append(results, g)
	}
	return results, nil
}

func (s *AccountQueryService) FindAllGoalsWithSum(
	ctx context.Context,
	uid string,
) ([]query.AccountWithSum, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT 
			a.id, 
			a.user_id, 
			a.type, 
			a.name, 
			a.is_saving, 
			a.currency, 
			a.target_amount, 
			a.start_date, 
			a.target_date, 
			c.id,
			c.name,
			COALESCE(c.type,'income'),
			a.created_at, 
			a.updated_at, 
			COALESCE(SUM(ms.money_in - ms.money_out), 0) AS sum
		 FROM accounts a
		 LEFT JOIN expense_categories c
			ON a.category_id = c.id
		 LEFT JOIN monthly_summary ms
			ON a.id = ms.account_id
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
		var c query.CategoryDTO
		if err := rows.Scan(
			&g.ID,
			&g.UserID,
			&g.Type,
			&g.Name,
			&g.IsSaving,
			&g.Currency,
			&g.TargetAmount,
			&g.StartDate,
			&g.TargetDate,
			&c.ID,
			&c.Name,
			&c.Type,
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

func (s *AccountQueryService) GetBanksDailyBalanceSince(
	ctx context.Context,
	uid string,
	since time.Time,
) ([]query.AccountDailyBalance, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`
		WITH
		  -- since is normalized to the 1st of its month because monthly_summary
		  -- is month-granular: comparing ms.month < $2 directly would pull the
		  -- current month's summary into the baseline and double-count it with
		  -- the ledger deltas. The days series also starts at start_day so the
		  -- month-to-date deltas before $2 still accumulate into the running
		  -- balance; the chart simply ignores the days before $2.
		  params AS (
		    SELECT
		      DATE_TRUNC('month', $2::TIMESTAMP)::date AS start_day
		  ),
		  -- 1. One row per calendar day, from the 1st of the month
		  days AS (
		    SELECT
			  generate_series((SELECT start_day FROM params), CURRENT_DATE, INTERVAL '1 day')::date AS day
		  ),
		  -- 2. Balance carried into the month (all history before it)
		  baseline AS (
		    SELECT
		      a.id,
		      a.name,
		      COALESCE(SUM(ms.money_in - ms.money_out), 0) AS bal
		    FROM
		      accounts a
		      LEFT JOIN monthly_summary ms ON ms.account_id = a.id
		      AND ms.month < (SELECT start_day FROM params)
		    WHERE
		      a.user_id = $1
		      AND a.type = 'bank'
		      AND a.deleted_at IS NULL
		    GROUP BY
		      a.id,
		      a.name
		  ),
		  -- 3. Net movement per account per day, from the 1st of the month
		  deltas AS (
		    SELECT
		      a.id,
		      t.date::date AS day,
		      SUM(
		        CASE
		          WHEN t.from_account_id = a.id THEN - t.amount
		          ELSE t.amount
		        END
		      ) AS delta
		    FROM
		      accounts a
		      JOIN ledger t ON t.date >= (SELECT start_day FROM params)
		      AND (
		        t.to_account_id = a.id
		        OR t.from_account_id = a.id
		      )
		    WHERE
		      a.user_id = $1
		      AND a.type = 'bank'
		      AND a.deleted_at IS NULL
		    GROUP BY
		      a.id,
		      t.date::date
		  )
		  -- 4. Cross join days x accounts, then run the cumulative sum
		SELECT
		  d.day,
		  a.id,
		  a.name,
		  b.bal + COALESCE(
		    SUM(dt.delta) OVER (
		      PARTITION BY
		        a.id
		      ORDER BY
		        d.day ROWS BETWEEN UNBOUNDED PRECEDING
		        AND CURRENT ROW
		    ),
		    0
		  ) AS balance
		FROM
		  days d
		  CROSS JOIN accounts a
		  JOIN baseline b ON b.id = a.id
		  LEFT JOIN deltas dt ON dt.id = a.id
		  AND dt.day = d.day
		WHERE
		  a.user_id = $1
		  AND a.type = 'bank'
		  AND a.deleted_at IS NULL
		ORDER BY
		  d.day,
		  a.id;`,
		uid,
		since,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []query.AccountDailyBalance

	for rows.Next() {
		var r query.AccountDailyBalance
		if err := rows.Scan(
			&r.Day,
			&r.ID,
			&r.Name,
			&r.Balance,
		); err != nil {
			return nil, err
		}

		results = append(results, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, query.ErrNoHistory
	}
	return results, nil
}

func (s *AccountQueryService) GetAccountDailyBalanceSince(
	ctx context.Context,
	id string,
	since time.Time,
) ([]query.AccountDailyBalance, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`
		WITH
		  -- since is normalized to the 1st of its month because monthly_summary
		  -- is month-granular: comparing ms.month < $2 directly would pull the
		  -- current month's summary into the baseline and double-count it with
		  -- the ledger deltas. The days series also starts at start_day so the
		  -- month-to-date deltas before $2 still accumulate into the running
		  -- balance; the chart simply ignores the days before $2.
		  params AS (
		    SELECT
		      DATE_TRUNC('month', $2::TIMESTAMP)::date AS start_day
		  ),
		  -- 1. One row per calendar day, from the 1st of the month
		  days AS (
		    SELECT
				generate_series((SELECT start_day FROM params), CURRENT_DATE, INTERVAL '1 day')::date AS day
		  ),
		  -- 2. Balance carried into the month (all history before it)
		  baseline AS (
		    SELECT
		      a.id,
		      a.name,
		      COALESCE(SUM(ms.money_in - ms.money_out), 0) AS bal
		    FROM
		      accounts a
		      LEFT JOIN monthly_summary ms ON ms.account_id = a.id
		      AND ms.month < (SELECT start_day FROM params)
		    WHERE
		      a.id = $1
		      AND a.deleted_at IS NULL
		    GROUP BY
		      a.id,
		      a.name
		  ),
		  -- 3. Net movement per account per day, from the 1st of the month
		  deltas AS (
		    SELECT
		      a.id,
		      t.date::date AS day,
		      SUM(
		        CASE
		          WHEN t.from_account_id = a.id THEN - t.amount
		          ELSE t.amount
		        END
		      ) AS delta
		    FROM
		      accounts a
		      JOIN ledger t ON t.date >= (SELECT start_day FROM params)
		      AND (
		        t.to_account_id = a.id
		        OR t.from_account_id = a.id
		      )
		    WHERE
		      a.id = $1
		      AND a.deleted_at IS NULL
		    GROUP BY
		      a.id,
		      t.date::date
		  )
		  -- 4. Cross join days x accounts, then run the cumulative sum
		SELECT
		  d.day,
		  a.id,
		  a.name,
		  b.bal + COALESCE(
		    SUM(dt.delta) OVER (
		      PARTITION BY
		        a.id
		      ORDER BY
		        d.day ROWS BETWEEN UNBOUNDED PRECEDING
		        AND CURRENT ROW
		    ),
		    0
		  ) AS balance
		FROM
		  days d
		  CROSS JOIN accounts a
		  JOIN baseline b ON b.id = a.id
		  LEFT JOIN deltas dt ON dt.id = a.id
		  AND dt.day = d.day
		WHERE
		  a.id = $1
		  AND a.deleted_at IS NULL
		ORDER BY
		  d.day,
		  a.id;
		`,
		id,
		since,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []query.AccountDailyBalance

	for rows.Next() {
		var r query.AccountDailyBalance
		if err := rows.Scan(
			&r.Day,
			&r.ID,
			&r.Name,
			&r.Balance,
		); err != nil {
			return nil, err
		}

		results = append(results, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, query.ErrNoHistory
	}
	return results, nil
}

func (s *AccountQueryService) GetAccountDailyBalanceAndStatsSince(
	ctx context.Context,
	id string,
	since time.Time,
) (*query.AccountDailyBalanceWithStatsSince, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`
		WITH
		  -- since is normalized to the 1st of its month because monthly_summary
		  -- is month-granular: comparing ms.month < $2 directly would pull the
		  -- current month's summary into the baseline and double-count it with
		  -- the ledger deltas. The days series also starts at start_day so the
		  -- month-to-date deltas before $2 still accumulate into the running
		  -- balance; the chart simply ignores the days before $2.
		  params AS (
		    SELECT
		      DATE_TRUNC('month', $2::TIMESTAMP)::date AS start_day
		  ),
		  -- 1. One row per calendar day, from the 1st of the month
		  days AS (
		    SELECT
				generate_series((SELECT start_day FROM params), CURRENT_DATE, INTERVAL '1 day')::date AS day
		  ),
		  -- 2. Balance carried into the month (all history before it)
		  baseline AS (
		    SELECT
		      a.id,
		      a.name,
		      COALESCE(SUM(ms.money_in - ms.money_out), 0) AS bal
		    FROM
		      accounts a
		      LEFT JOIN monthly_summary ms ON ms.account_id = a.id
		      AND ms.month < (SELECT start_day FROM params)
		    WHERE
		      a.id = $1
		      AND a.deleted_at IS NULL
		    GROUP BY
		      a.id,
		      a.name
		  ),
		  -- 3. Net movement per account per day, from the 1st of the month
		  deltas AS (
		    SELECT
		      a.id,
		      t.date::date AS day,
		      SUM(
		        CASE
		          WHEN t.from_account_id = a.id THEN - t.amount
		          ELSE t.amount
		        END
		      ) AS delta,
		      SUM(
		        CASE
		          WHEN t.to_account_id = a.id THEN t.amount
		          ELSE 0
		        END
		      ) AS money_in,
		      SUM(
		        CASE
		          WHEN t.from_account_id = a.id THEN t.amount
		          ELSE 0
		        END
		      ) AS money_out,
		      SUM(1) AS transaction_count
		    FROM
		      accounts a
		      JOIN ledger t ON t.date >= (SELECT start_day FROM params)
		      AND (
		        t.to_account_id = a.id
		        OR t.from_account_id = a.id
		      )
		    WHERE
		      a.id = $1
		      AND a.deleted_at IS NULL
		    GROUP BY
		      a.id,
		      t.date::date
		  )
		  -- 4. Cross join days x accounts, then run the cumulative sum
		SELECT
		  d.day,
		  a.id,
		  a.name,
		  b.bal + COALESCE(
		    SUM(dt.delta) OVER (
		      PARTITION BY
		        a.id
		      ORDER BY
		        d.day ROWS BETWEEN UNBOUNDED PRECEDING
		        AND CURRENT ROW
		    ),
		    0
		  ) AS balance,
		  COALESCE(
		    SUM(dt.money_in) OVER (
		      PARTITION BY
		        a.id
		      ORDER BY
		        d.day ROWS BETWEEN UNBOUNDED PRECEDING
		        AND CURRENT ROW
		    ),
		    0
		  ) AS money_in,
		  COALESCE(
		    SUM(dt.money_out) OVER (
		      PARTITION BY
		        a.id
		      ORDER BY
		        d.day ROWS BETWEEN UNBOUNDED PRECEDING
		        AND CURRENT ROW
		    ),
		    0
		  ) AS money_out,
		  COALESCE(
		    SUM(dt.transaction_count) OVER (
		      PARTITION BY
		        a.id
		      ORDER BY
		        d.day ROWS BETWEEN UNBOUNDED PRECEDING
		        AND CURRENT ROW
		    ),
		    0
		  ) AS transaction_count
		FROM
		  days d
		  CROSS JOIN accounts a
		  JOIN baseline b ON b.id = a.id
		  LEFT JOIN deltas dt ON dt.id = a.id
		  AND dt.day = d.day
		WHERE
		  a.id = $1
		  AND a.deleted_at IS NULL
		ORDER BY
		  d.day,
		  a.id;
		`,
		id,
		since,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []query.AccountDailyBalance
	var moneyInSince, moneyOutSince, transactionsSince int

	for rows.Next() {
		var r query.AccountDailyBalance
		var moneyIn, moneyOut, transactionCount int
		if err := rows.Scan(
			&r.Day,
			&r.ID,
			&r.Name,
			&r.Balance,
			&moneyIn,
			&moneyOut,
			&transactionCount,
		); err != nil {
			return nil, err
		}
		results = append(results, r)
		// money_in/money_out/transaction_count are cumulative window columns and
		// rows arrive ordered by day, so the last row scanned holds the range
		// totals for MoneyIn/MoneyOut/Transactions.
		moneyInSince, moneyOutSince, transactionsSince = moneyIn, moneyOut, transactionCount
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, query.ErrNoHistory
	}
	return &query.AccountDailyBalanceWithStatsSince{
		Data:         results,
		MoneyIn:      moneyInSince,
		MoneyOut:     moneyOutSince,
		Transactions: transactionsSince,
	}, nil
}

func (s *AccountQueryService) GetAccountWithMinRunningBalance(
	ctx context.Context,
	id string,
	fromDate time.Time,
	untilDate *time.Time,
	excludeEntryID *string,
) (*query.AccountWithMinRunningBalance, error) {
	row := s.db.QueryRowContext(
		ctx,
		`
		WITH
		  baseline AS (
		    SELECT
		      COALESCE(SUM(ms.money_in - ms.money_out), 0) AS bal
		    FROM
		      monthly_summary ms
		    WHERE
		      ms.account_id = $1
		      AND ms.month < $2
		  ),
		  transactions AS (
		    SELECT
		      day,
		      SUM(net) AS net
		    FROM
		      (
		        SELECT
		          t.date::date AS day,
		          CASE
		            WHEN t.from_account_id = $1 THEN t.amount * -1
		            ELSE t.amount
		          END AS net
		        FROM
		          ledger t
		        WHERE
		          (
		            t.from_account_id = $1
		            OR t.to_account_id = $1
		          )
		          AND t.date >= $2
		          AND ($4::TIMESTAMP IS NULL OR t.date <= $4)
		          AND ($5::UUID IS NULL OR t.id != $5)
		        UNION ALL
		        -- Guarantee fromDate always has a row so its balance is never
		        -- excluded from the minimum search.
		        SELECT
		          $3::DATE AS day,
		          0 AS net
		        WHERE
		          $4::TIMESTAMP IS NULL OR $3 <= $4
		      ) combined
		    GROUP BY
		      day
		    ORDER BY
		      day
		  ),
		  running AS (
		    SELECT
		      day,
		      net,
		      (
		        SELECT
		          bal
		        FROM
		          baseline
		      ) + SUM(net) OVER (
		        ORDER BY
		          day
		      ) AS running_balance
		    FROM
		      transactions
		  ),
          result AS (
            SELECT
              running_balance AS min_running_balance,
              day AS min_date
            FROM
              running
            WHERE
              running_balance = (
                SELECT
                  MIN(running_balance)
                FROM
                  running
                WHERE
                  day >= $3
              )
		      AND day >= $3
            ORDER BY
              day
            LIMIT
              1
          )
		SELECT
			a.id,
			a.user_id,
			a.type,
			a.name,
			a.is_saving,
			a.currency,
			a.target_amount,
			a.start_date,
			a.target_date,
			COALESCE(c.id, $1),
			COALESCE(c.name,''),
			COALESCE(c.type,'income'),
			a.created_at,
			a.updated_at,
			COALESCE(r.min_running_balance, (SELECT bal FROM baseline)),
			COALESCE(r.min_date, $3)
		FROM accounts a
		LEFT JOIN expense_categories c
			ON a.category_id = c.id
		LEFT JOIN result r ON TRUE
		WHERE
			a.id = $1`,
		id,
		time.Date(fromDate.Year(), fromDate.Month(), 1, 0, 0, 0, 0, fromDate.Location()),
		fromDate,
		untilDate,
		excludeEntryID,
	)

	var g query.AccountWithMinRunningBalance
	var c query.CategoryDTO

	if err := row.Scan(
		&g.ID,
		&g.UserID,
		&g.Type,
		&g.Name,
		&g.IsSaving,
		&g.Currency,
		&g.TargetAmount,
		&g.StartDate,
		&g.TargetDate,
		&c.ID,
		&c.Name,
		&c.Type,
		&g.CreatedAt,
		&g.UpdatedAt,
		&g.MinRunningBalance,
		&g.MinDate,
	); err != nil {
		return nil, err
	}
	g.Category = &c
	return &g, nil
}
