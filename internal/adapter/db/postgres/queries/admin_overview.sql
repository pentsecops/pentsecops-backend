-- name: GetTotalProjectsCount :one
SELECT COUNT(*) as total FROM projects;

-- name: GetOpenProjectsCount :one
SELECT COUNT(*) as count FROM projects WHERE status = 'Open';

-- name: GetInProgressProjectsCount :one
SELECT COUNT(*) as count FROM projects WHERE status = 'In Progress';

-- name: GetCompletedProjectsCount :one
SELECT COUNT(*) as count FROM projects WHERE status = 'Completed';

-- name: GetTotalVulnerabilitiesCount :one
SELECT COUNT(*) as total FROM vulnerabilities;

-- name: GetActiveUsersCount :one
SELECT COUNT(*) as total FROM users WHERE status = 'Active';

-- name: GetActivePentestersCount :one
SELECT COUNT(*) as count FROM users WHERE status = 'Active' AND role = 'Pentester';

-- name: GetActiveStakeholdersCount :one
SELECT COUNT(*) as count FROM users WHERE status = 'Active' AND role = 'Stakeholder';

-- name: GetOpenIssuesCount :one
SELECT COUNT(*) as count FROM vulnerabilities WHERE status = 'Open';

-- name: GetCriticalIssuesCount :one
SELECT COUNT(*) as count FROM vulnerabilities WHERE severity = 'Critical';

-- name: GetVulnerabilitiesBySeverity :many
SELECT 
    severity,
    COUNT(*) as count
FROM vulnerabilities
GROUP BY severity
ORDER BY 
    CASE severity
        WHEN 'Critical' THEN 1
        WHEN 'High' THEN 2
        WHEN 'Medium' THEN 3
        WHEN 'Low' THEN 4
    END;

-- name: GetTop5DomainsByVulnerabilities :many
SELECT
    COALESCE(d.domain_name, 'Unknown') as domain,
    COUNT(*) as vulnerability_count
FROM vulnerabilities v
LEFT JOIN domains d ON v.domain_id = d.id
WHERE v.domain_id IS NOT NULL
GROUP BY d.domain_name
ORDER BY vulnerability_count DESC
LIMIT 5;

-- name: GetProjectStatusDistribution :many
SELECT 
    status,
    COUNT(*) as count
FROM projects
GROUP BY status
ORDER BY 
    CASE status
        WHEN 'Open' THEN 1
        WHEN 'In Progress' THEN 2
        WHEN 'Completed' THEN 3
    END;

-- name: GetRecentActivityLogs :many
SELECT 
    id,
    user_id,
    user_name,
    action,
    entity_type,
    entity_id,
    ip_address,
    user_agent,
    created_at
FROM activity_logs
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetActivityLogsCount :one
SELECT COUNT(*) as total FROM activity_logs;

-- name: CreateActivityLog :one
INSERT INTO activity_logs (
    id,
    user_id,
    user_name,
    action,
    entity_type,
    entity_id,
    ip_address,
    user_agent,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

