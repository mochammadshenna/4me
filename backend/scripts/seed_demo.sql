-- Seed demo user for 4me Todos
-- Username: demo
-- Password: password (hashed with bcrypt cost 10)

-- Delete existing demo user if exists
DELETE FROM users WHERE username = 'demo' OR email = 'demo@example.com';

-- Create demo user
-- Password hash for 'password' with bcrypt cost 10: $2a$10$rKZB8ybH3dGYCHYvqxKPkO9F7QpqVZQ3xW8bB1z5o6xH9pY5.8Cne
INSERT INTO users (username, email, password_hash, created_at, updated_at)
VALUES (
    'demo',
    'demo@example.com',
    '$2a$10$rKZB8ybH3dGYCHYvqxKPkO9F7QpqVZQ3xW8bB1z5o6xH9pY5.8Cne',
    NOW(),
    NOW()
);

-- Verify creation
SELECT id, username, email, created_at FROM users WHERE username = 'demo';
