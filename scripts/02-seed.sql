INSERT INTO
  users (id, name, password_hash, created_at, updated_at)
VALUES
  (
    gen_random_uuid(),
    'user',
    'password',
    '2026-03-01',
    '2026-03-01'
  );
