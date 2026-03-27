-- name: GetTotalNotificationsSent :one
SELECT COUNT(*) as total
FROM notifications
WHERE created_by = $1;

-- name: ListNotifications :many
SELECT 
    n.id,
    n.title,
    n.message,
    n.sent_to,
    n.recipient_id,
    n.created_by,
    n.status,
    n.type,
    n.is_read,
    n.created_at,
    u.full_name as created_by_name,
    CASE 
        WHEN n.sent_to = 'specific_user' THEN ru.full_name
        ELSE NULL
    END as recipient_name
FROM notifications n
LEFT JOIN users u ON n.created_by = u.id
LEFT JOIN users ru ON n.recipient_id = ru.id
WHERE n.created_by = $1
ORDER BY n.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountNotifications :one
SELECT COUNT(*) as total
FROM notifications
WHERE created_by = $1;

-- name: CreateNotification :one
INSERT INTO notifications (
    id,
    title,
    message,
    sent_to,
    recipient_id,
    created_by,
    status,
    type,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetNotificationByID :one
SELECT 
    n.id,
    n.title,
    n.message,
    n.sent_to,
    n.recipient_id,
    n.created_by,
    n.status,
    n.type,
    n.is_read,
    n.created_at,
    u.full_name as created_by_name,
    CASE 
        WHEN n.sent_to = 'specific_user' THEN ru.full_name
        ELSE NULL
    END as recipient_name
FROM notifications n
LEFT JOIN users u ON n.created_by = u.id
LEFT JOIN users ru ON n.recipient_id = ru.id
WHERE n.id = $1;

-- name: ListImportantAlerts :many
SELECT
    a.id,
    a.title,
    a.message,
    a.alert_type,
    a.priority,
    a.source,
    a.sender_id,
    a.is_resolved,
    a.resolved_by,
    a.resolved_at,
    a.created_at,
    a.updated_at,
    u.full_name as pentester_name
FROM alerts a
LEFT JOIN users u ON a.sender_id = u.id
WHERE a.alert_type = 'important_update' AND a.source = 'pentester'
ORDER BY a.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountImportantAlerts :one
SELECT COUNT(*) as total
FROM alerts
WHERE alert_type = 'important_update' AND source = 'pentester';

-- name: CreateAlert :one
INSERT INTO alerts (
    id,
    title,
    message,
    alert_type,
    priority,
    source,
    sender_id,
    recipient_id,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetAlertByID :one
SELECT
    a.id,
    a.title,
    a.message,
    a.alert_type,
    a.priority,
    a.source,
    a.sender_id,
    a.is_resolved,
    a.resolved_by,
    a.resolved_at,
    a.created_at,
    a.updated_at,
    u.full_name as pentester_name
FROM alerts a
LEFT JOIN users u ON a.sender_id = u.id
WHERE a.id = $1;

-- name: GetUsersByRole :many
SELECT id, email, full_name, role
FROM users
WHERE role = $1 AND is_active = true;

