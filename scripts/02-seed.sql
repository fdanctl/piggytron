INSERT INTO
  users (id, name, password_hash, created_at, updated_at)
VALUES
  (
    gen_random_uuid(),
    'gopher',
    -- 123
    '$argon2id$v=19$m=65536,t=1,p=4$gmqRhs4/Neoo708g05eW8A$fpWmdn+np0pi31xsvBqwEmcXaVnGGxwnD3NPnAHrm9k',
    '2026-03-01',
    '2026-03-01'
  );

-- INCOME CATEGORIES
INSERT INTO
  income_categories (id, user_id, name, created_at, updated_at)
SELECT
  gen_random_uuid(),
  u.id,
  'Job',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  income_categories (id, user_id, name, created_at, updated_at)
SELECT
  gen_random_uuid(),
  u.id,
  'Side Hustle',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  income_categories (id, user_id, name, created_at, updated_at)
SELECT
  gen_random_uuid(),
  u.id,
  'Gifts',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  income_categories (id, user_id, name, created_at, updated_at)
SELECT
  gen_random_uuid(),
  u.id,
  'Others',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

-- EXPENSES CATEGORIES
INSERT INTO
  expense_categories (id, user_id, name, type, created_at, updated_at)
SELECT
  gen_random_uuid(),
  u.id,
  'House',
  'needs',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (id, user_id, name, type, created_at, updated_at)
SELECT
  gen_random_uuid(),
  u.id,
  'Utilities',
  'needs',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (id, user_id, name, type, created_at, updated_at)
SELECT
  gen_random_uuid(),
  u.id,
  'Transportation',
  'needs',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (id, user_id, name, type, created_at, updated_at)
SELECT
  gen_random_uuid(),
  u.id,
  'Groceries',
  'needs',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (id, user_id, name, type, created_at, updated_at)
SELECT
  gen_random_uuid(),
  u.id,
  'Entertainment',
  'wants',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (id, user_id, name, type, created_at, updated_at)
SELECT
  gen_random_uuid(),
  u.id,
  'Clothing',
  'wants',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (id, user_id, name, type, created_at, updated_at)
SELECT
  gen_random_uuid(),
  u.id,
  'Shopping',
  'wants',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (id, user_id, name, type, created_at, updated_at)
SELECT
  gen_random_uuid(),
  u.id,
  'Eating Out',
  'wants',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (id, user_id, name, type, created_at, updated_at)
SELECT
  gen_random_uuid(),
  u.id,
  'Health',
  'needs',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (id, user_id, name, type, created_at, updated_at)
SELECT
  gen_random_uuid(),
  u.id,
  'Household Items',
  'needs',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (id, user_id, name, type, created_at, updated_at)
SELECT
  gen_random_uuid(),
  u.id,
  'Personal',
  'wants',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (id, user_id, name, type, created_at, updated_at)
SELECT
  gen_random_uuid(),
  u.id,
  'Savings',
  'savings',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (id, user_id, name, type, created_at, updated_at)
SELECT
  gen_random_uuid(),
  u.id,
  'Investing',
  'savings',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (id, user_id, name, type, created_at, updated_at)
SELECT
  gen_random_uuid(),
  u.id,
  'Emergency Fund',
  'savings',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (id, user_id, name, type, created_at, updated_at)
SELECT
  gen_random_uuid(),
  u.id,
  'Gifts/Donations',
  'wants',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (id, user_id, name, type, created_at, updated_at)
SELECT
  gen_random_uuid(),
  u.id,
  'Miscellaneous',
  'needs',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';
