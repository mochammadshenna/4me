-- Seed demo user
INSERT INTO users (username, email, password_hash, created_at, updated_at) 
VALUES ('demo', 'demo@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', NOW(), NOW())
ON CONFLICT (username) DO NOTHING;

-- Seed demo project
INSERT INTO projects (name, description, user_id, created_at, updated_at)
SELECT 'Demo Project', 'A sample project to get started', u.id, NOW(), NOW()
FROM users u WHERE u.username = 'demo'
ON CONFLICT DO NOTHING;
