-- name: CreateUser :one
INSERT INTO users (
    id,
    email,
    password_hash,
    full_name,
    role,
    is_active,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT 
    u.id,
    u.email,
    u.full_name,
    u.role,
    u.is_active,
    u.last_login,
    u.created_at,
    u.updated_at,
    COUNT(DISTINCT CASE 
        WHEN u.role = 'pentester' THEN p.id 
        ELSE NULL 
    END) as project_count
FROM users u
LEFT JOIN projects p ON u.id = p.assigned_to AND u.role = 'pentester'
GROUP BY u.id, u.email, u.full_name, u.role, u.is_active, u.last_login, u.created_at, u.updated_at
ORDER BY u.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: CountUsersByRole :one
SELECT COUNT(*) FROM users
WHERE role = $1 AND is_active = $2;

-- name: CountInactiveUsers :one
SELECT COUNT(*) FROM users
WHERE is_active = false;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: UpdateUserLastLogin :exec
UPDATE users
SET last_login = $1, updated_at = $2
WHERE id = $3;

-- name: GetUserStats :one
SELECT 
    COUNT(*) FILTER (WHERE role = 'pentester' AND is_active = true) as active_pentesters,
    COUNT(*) FILTER (WHERE role = 'stakeholder' AND is_active = true) as active_stakeholders,
    COUNT(*) FILTER (WHERE is_active = false) as inactive_users,
    COUNT(*) as total_users
FROM users;

-- name: ListAllUsersForExport :many
SELECT 
    u.id,
    u.email,
    u.full_name,
    u.role,
    u.is_active,
    u.last_login,
    u.created_at,
    COUNT(DISTINCT CASE 
        WHEN u.role = 'pentester' THEN p.id 
        ELSE NULL 
    END) as project_count
FROM users u
LEFT JOIN projects p ON u.id = p.assigned_to AND u.role = 'pentester'
GROUP BY u.id, u.email, u.full_name, u.role, u.is_active, u.last_login, u.created_at
ORDER BY u.created_at DESC;

-- name: CheckEmailExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);

