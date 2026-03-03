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
