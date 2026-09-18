-- Account balance computation: covers from_account JOINs
CREATE INDEX IF NOT EXISTS idx_ledger_from_account
  ON ledger (from_account_id);
