CREATE TABLE IF NOT EXISTS budget_holds (
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  month DATE NOT NULL, -- e.g. '2026-03-01'
  amount BIGINT NOT NULL CHECK (amount >= 0), -- in cents
  --
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  --
  PRIMARY KEY (user_id, month)
);
