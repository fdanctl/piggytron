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

-- acounts
INSERT INTO
  accounts (
    id,
    user_id,
    name,
    type,
    currency,
    created_at,
    updated_at
  )
SELECT
  gen_random_uuid(),
  u.id,
  'Cash',
  'bank',
  'eur',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  accounts (
    id,
    user_id,
    name,
    type,
    currency,
    created_at,
    updated_at
  )
SELECT
  gen_random_uuid(),
  u.id,
  'Main',
  'bank',
  'eur',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  accounts (
    id,
    user_id,
    name,
    type,
    currency,
    created_at,
    updated_at
  )
SELECT
  gen_random_uuid(),
  u.id,
  'Savings',
  'bank',
  'eur',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
WHERE
  u.name = 'gopher';

INSERT INTO
  accounts (
    id,
    user_id,
    name,
    type,
    currency,
    target_amount,
    target_date,
    category_id,
    created_at,
    updated_at
  )
SELECT
  gen_random_uuid(),
  u.id,
  'Ferrari',
  'goal',
  'eur',
  2000000,
  '2026-06-06',
  c.id,
  '2026-03-11',
  '2026-03-11'
FROM
  users u
  JOIN expense_categories c ON u.id = c.user_id
WHERE
  u.name = 'gopher'
  AND c.name = 'Shopping';

-- ai
-- 1️⃣ Income transactions
INSERT INTO
  transactions (
    id,
    user_id,
    type,
    to_account_id,
    income_category_id,
    expense_category_id,
    amount,
    description,
    date,
    created_at
  )
SELECT
  gen_random_uuid(),
  u.id,
  'income',
  a.id,
  ic.id,
  NULL,
  ROUND(random() * 3000 + 500)::BIGINT,
  'Salary / Income',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
  JOIN accounts a ON u.id = a.user_id
  JOIN income_categories ic ON u.id = ic.user_id
WHERE
  u.name = 'gopher'
  AND ic.name = 'Job';

-- 2️⃣ Expense transactions
INSERT INTO
  transactions (
    id,
    user_id,
    type,
    from_account_id,
    income_category_id,
    expense_category_id,
    amount,
    description,
    date,
    created_at
  )
SELECT
  gen_random_uuid(),
  u.id,
  'expense',
  a.id,
  NULL,
  ec.id,
  ROUND(random() * 200 + 10)::BIGINT,
  'Monthly expense',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
  JOIN accounts a ON u.id = a.user_id
  JOIN expense_categories ec ON u.id = ec.user_id
WHERE
  u.name = 'gopher'
  AND ec.type IN ('needs', 'wants');

-- 3️⃣ Transfer transactions
INSERT INTO
  transactions (
    id,
    user_id,
    type,
    from_account_id,
    to_account_id,
    income_category_id,
    expense_category_id,
    amount,
    description,
    date,
    created_at
  )
SELECT
  gen_random_uuid(),
  u.id,
  'transfer',
  a_from.id,
  a_to.id,
  NULL,
  NULL,
  ROUND(random() * 1000 + 100)::BIGINT,
  'Transfer between accounts',
  '2026-03-11',
  '2026-03-11'
FROM
  users u
  JOIN accounts a_from ON u.id = a_from.user_id
  JOIN accounts a_to ON u.id = a_to.user_id
WHERE
  u.name = 'gopher'
  AND a_from.name = 'Cash'
  AND a_to.name = 'Main';
