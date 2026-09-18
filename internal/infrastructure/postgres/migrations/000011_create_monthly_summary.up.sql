CREATE TABLE IF NOT EXISTS monthly_summary (
  account_id UUID NOT NULL REFERENCES accounts (id) ON DELETE CASCADE,
  month DATE NOT NULL CHECK (month = DATE_TRUNC('month', month)), -- first day: '2026-06-01'
  money_in BIGINT NOT NULL CHECK (money_in >= 0), -- income + transfers_into_this_account
  money_out BIGINT NOT NULL CHECK (money_out >= 0), -- expense + transfers_from_this_account
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  PRIMARY KEY (account_id, month)
);
