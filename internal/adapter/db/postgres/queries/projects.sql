-- name: CreateProject :one
INSERT INTO projects (
    id,
    name,
    type,
    assigned_to,
    deadline,
    scope,
    status,
    current_phase,
    created_by,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetProjectByID :one
SELECT * FROM projects
WHERE id = $1;

-- name: GetProjectByName :one
SELECT * FROM projects
WHERE name = $1;

-- name: ListProjects :many
SELECT 
    p.*,
    u.full_name as assigned_to_name,
    COALESCE(v.vuln_count, 0) as vulnerability_count
FROM projects p
LEFT JOIN users u ON p.assigned_to = u.id
LEFT JOIN (
    SELECT project_id, COUNT(*) as vuln_count
    FROM vulnerabilities
    GROUP BY project_id
) v ON p.id = v.project_id
ORDER BY p.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountProjects :one
SELECT COUNT(*) FROM projects;

-- name: GetProjectStats :one
SELECT 
    COUNT(*) FILTER (WHERE status = 'open') as open_count,
    COUNT(*) FILTER (WHERE status = 'in_progress') as in_progress_count,
    COUNT(*) FILTER (WHERE status = 'completed') as completed_count
FROM projects;

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = $1;

-- name: UpdateProjectStatus :exec
UPDATE projects 
SET status = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: GetProjectsByPentester :many
SELECT 
    p.*,
    u.full_name as assigned_to_name,
    COALESCE(v.vuln_count, 0) as vulnerability_count
FROM projects p
LEFT JOIN users u ON p.assigned_to = u.id
LEFT JOIN (
    SELECT project_id, COUNT(*) as vuln_count
    FROM vulnerabilities
    GROUP BY project_id
) v ON p.id = v.project_id
WHERE p.assigned_to = $1
ORDER BY p.created_at DESC;

-- name: GetPentesters :many
SELECT id, full_name, email
FROM users
WHERE role = 'pentester' AND is_active = true
ORDER BY full_name ASC;

