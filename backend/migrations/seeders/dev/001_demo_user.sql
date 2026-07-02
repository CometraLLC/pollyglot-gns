-- Development-only: seeded account for manual testing.
-- Sign in with demo@pollyglot.dev / Password123!
INSERT INTO users (id, email, password_hash, name, is_active, email_verified)
VALUES (
    'a0000000-0000-4000-8000-000000000001',
    'demo@pollyglot.dev',
    '$2a$12$OuB.l181ohywyXICmH7SuuxxaKsqtC05Za2E9nVpOTg58k7IpmZ0m',
    'Demo User',
    TRUE,
    TRUE
)
ON CONFLICT (email) DO NOTHING;

-- Give the demo user the standard role (role ids are serial, so look it up)
INSERT INTO user_roles (user_id, role_id)
SELECT 'a0000000-0000-4000-8000-000000000001', r.id
FROM roles r
WHERE r.name = 'User'
  AND NOT EXISTS (
    SELECT 1 FROM user_roles ur
    WHERE ur.user_id = 'a0000000-0000-4000-8000-000000000001' AND ur.role_id = r.id
  );
