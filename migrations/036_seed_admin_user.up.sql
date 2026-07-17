INSERT INTO users (user_id, username, password_hash, name, email, role, is_active, created_by, created_at, updated_by, updated_at)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'admin',
    '$2a$10$rmKmO2zXrigg/ImKKjjQ1uAjkAi3216rfuuM4c0xsurOiRO.mIo9u',
    'Administrador',
    'admin@crmstickers.com',
    'admin',
    TRUE,
    'system',
    NOW(),
    'system',
    NOW()
)
ON CONFLICT (user_id) DO NOTHING;
