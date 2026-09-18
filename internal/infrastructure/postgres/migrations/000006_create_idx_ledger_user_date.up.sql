-- Main list & filtering: covers user_id filter + date ordering
CREATE INDEX IF NOT EXISTS idx_ledger_user_date
  ON ledger (user_id, date DESC, created_at DESC);

