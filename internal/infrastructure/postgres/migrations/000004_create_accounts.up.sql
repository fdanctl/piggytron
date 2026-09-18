CREATE TABLE IF NOT EXISTS accounts (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  --
  name VARCHAR(50) NOT NULL,
  type VARCHAR(10) NOT NULL CHECK (type IN ('checking', 'savings', 'goal')),
  status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'closed')),
  --
  currency VARCHAR(10) NOT NULL,
  -- Goal-specific fields (NULL if type = 'checking' or type = 'savings')
  target_amount BIGINT,
  start_date TIMESTAMP,
  target_date TIMESTAMP,
  category_id UUID REFERENCES expense_categories (id),
  note TEXT,
  completed_at TIMESTAMP,
  cancelled_at TIMESTAMP,
  finalized_amount BIGINT,
  --
  closed_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  -- rules
  CHECK (
    (
      type != 'goal'
      AND target_amount IS NULL
      AND start_date IS NULL
      AND target_date IS NULL
      AND category_id IS NULL
      AND note IS NULL
      AND completed_at IS NULL
      AND cancelled_at IS NULL
      AND finalized_amount IS NULL
    )
    OR (
      type = 'goal'
      AND target_amount IS NOT NULL
      AND start_date IS NOT NULL
      -- AND target_date IS NOT NULL
      AND category_id IS NOT NULL
    )
  ),
  UNIQUE (user_id, type, name)
);
