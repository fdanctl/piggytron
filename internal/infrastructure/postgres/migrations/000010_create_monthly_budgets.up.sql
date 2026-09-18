-- one per category and month
CREATE TABLE IF NOT EXISTS monthly_budgets (
  category_id UUID NOT NULL REFERENCES expense_categories (id) ON DELETE CASCADE,
  month DATE NOT NULL, -- e.g. '2026-03-01'
  amount BIGINT NOT NULL CHECK (amount >= 0), -- in cents
  --
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  --
  PRIMARY KEY (category_id, month)
);
