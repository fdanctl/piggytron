CREATE TABLE users (
  id UUID PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE TABLE income_categories (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  --
  name VARCHAR(30) NOT NULL,
  --
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP
);

CREATE TABLE expense_categories (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  --
  name VARCHAR(30) NOT NULL,
  type VARCHAR(10) NOT NULL CHECK (type IN ('needs', 'wants', 'savings')),
  --
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP
);

CREATE TABLE accounts (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  --
  name VARCHAR(50) NOT NULL,
  type VARCHAR(10) NOT NULL CHECK (type IN ('bank', 'goal')),
  --
  currency VARCHAR(10) NOT NULL,
  -- Goal-specific fields (NULL if type = 'bank')
  target_amount BIGINT,
  target_date TIMESTAMP,
  category_id UUID REFERENCES expense_categories (id),
  --
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP
  -- rules
  CHECK (
    (
      type = 'bank'
      AND target_amount IS NULL
      AND target_date IS NULL
      AND category_id IS NULL
    )
    OR (
      type = 'goal'
      AND target_amount IS NOT NULL
      -- AND target_date IS NOT NULL
      AND category_id IS NOT NULL
    )
  )
);

CREATE TABLE transactions (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  --
  type VARCHAR(10) NOT NULL CHECK (type IN ('income', 'expense', 'transfer')),
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
  --
  created_at TIMESTAMP NOT NULL,
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
  )
);

-- one per category and month
CREATE TABLE monthly_budgets (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  category_id UUID NOT NULL REFERENCES expense_categories (id),
  --
  month DATE NOT NULL, -- e.g. '2026-03-01'
  amount BIGINT NOT NULL CHECK (amount > 0), -- in cents
  --
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  --
  UNIQUE (category_id, month)
);
