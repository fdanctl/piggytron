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
  expense_categories (
    id,
    user_id,
    name,
    expense_type,
    created_at,
    updated_at
  )
SELECT
  gen_random_uuid(),
  u.id,
  'House',
  1,
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (
    id,
    user_id,
    name,
    expense_type,
    created_at,
    updated_at
  )
SELECT
  gen_random_uuid(),
  u.id,
  'Utilities',
  1,
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (
    id,
    user_id,
    name,
    expense_type,
    created_at,
    updated_at
  )
SELECT
  gen_random_uuid(),
  u.id,
  'Transportation',
  1,
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (
    id,
    user_id,
    name,
    expense_type,
    created_at,
    updated_at
  )
SELECT
  gen_random_uuid(),
  u.id,
  'Groceries',
  1,
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (
    id,
    user_id,
    name,
    expense_type,
    created_at,
    updated_at
  )
SELECT
  gen_random_uuid(),
  u.id,
  'Entertainment',
  2,
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (
    id,
    user_id,
    name,
    expense_type,
    created_at,
    updated_at
  )
SELECT
  gen_random_uuid(),
  u.id,
  'Clothing',
  2,
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (
    id,
    user_id,
    name,
    expense_type,
    created_at,
    updated_at
  )
SELECT
  gen_random_uuid(),
  u.id,
  'Shopping',
  2,
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (
    id,
    user_id,
    name,
    expense_type,
    created_at,
    updated_at
  )
SELECT
  gen_random_uuid(),
  u.id,
  'Eating Out',
  2,
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (
    id,
    user_id,
    name,
    expense_type,
    created_at,
    updated_at
  )
SELECT
  gen_random_uuid(),
  u.id,
  'Health',
  1,
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (
    id,
    user_id,
    name,
    expense_type,
    created_at,
    updated_at
  )
SELECT
  gen_random_uuid(),
  u.id,
  'Household Items',
  1,
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (
    id,
    user_id,
    name,
    expense_type,
    created_at,
    updated_at
  )
SELECT
  gen_random_uuid(),
  u.id,
  'Personal',
  2,
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (
    id,
    user_id,
    name,
    expense_type,
    created_at,
    updated_at
  )
SELECT
  gen_random_uuid(),
  u.id,
  'Savings',
  3,
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (
    id,
    user_id,
    name,
    expense_type,
    created_at,
    updated_at
  )
SELECT
  gen_random_uuid(),
  u.id,
  'Investing',
  3,
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (
    id,
    user_id,
    name,
    expense_type,
    created_at,
    updated_at
  )
SELECT
  gen_random_uuid(),
  u.id,
  'Emergency Fund',
  3,
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (
    id,
    user_id,
    name,
    expense_type,
    created_at,
    updated_at
  )
SELECT
  gen_random_uuid(),
  u.id,
  'Gifts/Donations',
  2,
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  expense_categories (
    id,
    user_id,
    name,
    expense_type,
    created_at,
    updated_at
  )
SELECT
  gen_random_uuid(),
  u.id,
  'Miscellaneous',
  3,
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';
