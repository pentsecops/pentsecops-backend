-- name: GetUserForAuth :one
SELECT id, email, password_hash, full_name, role, is_active, force_password_change,
       failed_login_attempts, last_failed_login, account_locked_until, last_login,
       created_at, updated_at
FROM users
WHERE email = $1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateLastLoginTime :exec
UPDATE users
SET last_login = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserForcePasswordChange :exec
UPDATE users
SET force_password_change = $2, updated_at = NOW()
WHERE id = $1;

-- name: IncrementFailedLoginAttempts :exec
UPDATE users
SET failed_login_attempts = failed_login_attempts + 1,
    last_failed_login = NOW(),
    updated_at = NOW()
WHERE id = $1;

-- name: ResetFailedLoginAttempts :exec
UPDATE users
SET failed_login_attempts = 0,
    last_failed_login = NULL,
    updated_at = NOW()
WHERE id = $1;

-- name: LockAccount :exec
UPDATE users
SET account_locked_until = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: UnlockAccount :exec
UPDATE users
SET account_locked_until = NULL,
    failed_login_attempts = 0,
    last_failed_login = NULL,
    updated_at = NOW()
WHERE id = $1;

-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at)
VALUES ($1, $2, $3, $4, NOW());

-- name: GetRefreshToken :one
SELECT id, user_id, token_hash, expires_at, created_at
FROM refresh_tokens
WHERE token_hash = $1 AND expires_at > NOW();

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE token_hash = $1;

-- name: DeleteAllUserRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE user_id = $1;

-- name: DeleteExpiredRefreshTokens :exec
DELETE FROM refresh_tokens
WHERE expires_at <= NOW();

