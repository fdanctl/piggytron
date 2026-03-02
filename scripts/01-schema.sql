CREATE TABLE users (
  id UUID PRIMARY KEY,
  name varchar(50) NOT NULL UNIQUE,
  password_hash text NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);
