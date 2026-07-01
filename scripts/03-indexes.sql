-- Main list & filtering: covers user_id filter + date ordering
CREATE INDEX IF NOT EXISTS idx_ledger_user_date
  ON ledger (user_id, date DESC, created_at DESC);

-- Account balance computation: covers from_account/to_account JOINs
CREATE INDEX IF NOT EXISTS idx_ledger_from_account
  ON ledger (from_account_id);
CREATE INDEX IF NOT EXISTS idx_ledger_to_account
  ON ledger (to_account_id);

-- Budget page: covers type + date range filtering per user
CREATE INDEX IF NOT EXISTS idx_ledger_user_type_date
  ON ledger (user_id, type, date DESC);
