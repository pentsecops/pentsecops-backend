-- name: CreateTask :one
INSERT INTO tasks (
    id,
    project_id,
    title,
    description,
    status,
    priority,
    assigned_to,
    deadline,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: GetTaskByID :one
SELECT * FROM tasks
WHERE id = $1;

-- name: ListTasksByProject :many
SELECT 
    t.*,
    u.full_name as assigned_to_name
FROM tasks t
LEFT JOIN users u ON t.assigned_to = u.id
WHERE t.project_id = $1
ORDER BY t.created_at DESC;

-- name: ListAllTasks :many
SELECT 
    t.*,
    u.full_name as assigned_to_name,
    p.name as project_name
FROM tasks t
LEFT JOIN users u ON t.assigned_to = u.id
LEFT JOIN projects p ON t.project_id = p.id
ORDER BY t.created_at DESC;

-- name: UpdateTaskStatus :exec
UPDATE tasks 
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: UpdateTask :exec
UPDATE tasks 
SET 
    title = $2,
    description = $3,
    priority = $4,
    assigned_to = $5,
    deadline = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = $1;

-- name: GetTasksByStatus :many
SELECT 
    t.*,
    u.full_name as assigned_to_name,
    p.name as project_name
FROM tasks t
LEFT JOIN users u ON t.assigned_to = u.id
LEFT JOIN projects p ON t.project_id = p.id
WHERE t.status = $1
ORDER BY t.created_at DESC;

-- name: CountTasksByProject :one
SELECT COUNT(*) FROM tasks
WHERE project_id = $1;

