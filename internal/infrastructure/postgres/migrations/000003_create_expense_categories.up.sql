CREATE TABLE IF NOT EXISTS expense_categories (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  --
  name VARCHAR(30) NOT NULL,
  type VARCHAR(10) NOT NULL CHECK (type IN ('needs', 'wants', 'savings')),
  status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'archived')),
  --
  archived_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  --
  UNIQUE (user_id, name)
);
