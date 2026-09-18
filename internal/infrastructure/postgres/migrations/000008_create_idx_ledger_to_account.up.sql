-- Account balance computation: covers to_account JOINs
CREATE INDEX IF NOT EXISTS idx_ledger_to_account
  ON ledger (to_account_id);
