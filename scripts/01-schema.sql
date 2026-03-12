CREATE TABLE users (
  id UUID PRIMARY KEY,
  name VARCHAR(50) NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE TABLE income_categories (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  name VARCHAR(30) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP
);

-- expense_type needs = 1
-- expense_type wants = 2
-- expense_type savings = 3
CREATE TABLE expense_categories (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  name VARCHAR(30) NOT NULL,
  expense_type INT2,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP
);

CREATE TABLE banks (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  title VARCHAR(30) NOT NULL,
  currency VARCHAR(10) NOT NULL,
  initial_balance BIGINT NOT NULL, -- in cents
  is_savings BOOLEAN NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP
);

CREATE TABLE goals (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users (id) on DELETE CASCADE,
  category_id UUID NOT NULL REFERENCES expense_categories (id),
  title VARCHAR(40) NOT NULL,
  currency VARCHAR(10) NOT NULL,
  target_amount BIGINT NOT NULL, -- in cents
  target_date TIMESTAMP,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP
);

CREATE TABLE income (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users (id) on DELETE CASCADE,
  category_id UUID NOT NULL REFERENCES income_categories (id),
  description TEXT NOT NULL,
  amount BIGINT NOT NULL,
  distination_id UUID NOT NULL REFERENCES banks (id),
  date TIMESTAMP NOT NULL
);

CREATE TABLE expenses (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users (id) on DELETE CASCADE,
  category_id UUID REFERENCES expense_categories (id),
  description TEXT NOT NULL,
  amount BIGINT NOT NULL,
  source_id UUID REFERENCES banks (id),
  date TIMESTAMP NOT NULL
);

CREATE TABLE transfers (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users (id) on DELETE CASCADE,
  description TEXT NOT NULL,
  amount BIGINT NOT NULL,
  source_id UUID REFERENCES banks (id),
  distination_id UUID NOT NULL REFERENCES banks (id),
  date TIMESTAMP NOT NULL
);

-- TODO CREATE TABLE monthly_budget
