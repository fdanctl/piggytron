CREATE TABLE IF NOT EXISTS ledger (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  --
  type VARCHAR(20) NOT NULL CHECK (
    type IN (
      'income',
      'expense',
      'transfer',
      'interest',
      'goal-fulfillment',
      'initial-balance'
    )
  ),
  --
  from_account_id UUID REFERENCES accounts (id),
  to_account_id UUID REFERENCES accounts (id),
  --
  income_category_id UUID REFERENCES income_categories (id),
  expense_category_id UUID REFERENCES expense_categories (id),
  --
  amount BIGINT NOT NULL CHECK (amount > 0),
  description TEXT NOT NULL,
  date TIMESTAMP NOT NULL,
  note TEXT,
  --
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  -- rules
  CHECK (
    from_account_id IS NOT NULL
    OR to_account_id IS NOT NULL
  ),
  --
  CHECK (
    from_account_id IS NULL
    OR to_account_id IS NULL
    OR from_account_id <> to_account_id
  ),
  --
  CHECK (
    (
      type = 'income'
      AND to_account_id IS NOT NULL
      AND from_account_id IS NULL
      AND income_category_id IS NOT NULL
      AND expense_category_id IS NULL
    )
    OR (
      type = 'expense'
      AND from_account_id IS NOT NULL
      AND to_account_id IS NULL
      AND expense_category_id IS NOT NULL
      AND income_category_id IS NULL
    )
    OR (
      type = 'transfer'
      AND from_account_id IS NOT NULL
      AND to_account_id IS NOT NULL
      AND income_category_id IS NULL
      -- expense_category_id optional
    )
    OR (
      type = 'interest'
      AND from_account_id IS NULL
      AND to_account_id IS NOT NULL
      AND income_category_id IS NOT NULL
      AND expense_category_id IS NULL
    )
    OR (
      type = 'goal-fulfillment'
      AND from_account_id IS NOT NULL
      AND to_account_id IS NULL
      AND income_category_id IS NULL
      AND expense_category_id IS NOT NULL
    )
    OR (
      type = 'initial-balance'
      AND from_account_id IS NULL
      AND to_account_id IS NOT NULL
      AND income_category_id IS NULL
      AND expense_category_id IS NULL
    )
  )
);
