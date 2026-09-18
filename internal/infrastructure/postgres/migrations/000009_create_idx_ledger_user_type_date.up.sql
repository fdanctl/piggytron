-- Budget page: covers type + date range filtering per user
CREATE INDEX IF NOT EXISTS idx_ledger_user_type_date
  ON ledger (user_id, type, date DESC);
