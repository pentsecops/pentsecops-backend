-- name: GetVulnerabilitiesForSecurityScore :many
SELECT severity, status
FROM vulnerabilities;

-- name: GetActiveProjectsCounts :one
SELECT 
    COUNT(*) FILTER (WHERE status IN ('open', 'in_progress')) AS total_active,
    COUNT(*) FILTER (WHERE status = 'in_progress') AS in_progress,
    COUNT(*) FILTER (WHERE status = 'completed') AS completed
FROM projects;

-- name: GetStakeholderCriticalIssuesCount :one
SELECT COUNT(*) AS count
FROM vulnerabilities
WHERE severity = 'critical';

-- name: GetOpenVulnerabilitiesCurrentMonth :one
SELECT COUNT(*) AS count
FROM vulnerabilities
WHERE status = 'open';

-- name: GetOpenVulnerabilitiesLastMonth :one
SELECT COUNT(*) AS count
FROM vulnerabilities
WHERE status = 'open'
  AND created_at < DATE_TRUNC('month', CURRENT_DATE);

-- name: GetRemediationRateCounts :one
SELECT 
    COUNT(*) AS total,
    COUNT(*) FILTER (WHERE status IN ('remediated', 'verified')) AS remediated
FROM vulnerabilities;

-- name: GetSLAComplianceCounts :one
SELECT 
    COUNT(*) FILTER (WHERE due_date IS NOT NULL) AS total_with_due_date,
    COUNT(*) FILTER (WHERE due_date IS NOT NULL AND status IN ('remediated', 'verified') AND remediated_date <= due_date) AS remediated_on_time
FROM vulnerabilities;

-- name: GetVulnerabilityTrendByMonth :many
SELECT
    EXTRACT(EPOCH FROM DATE_TRUNC('month', created_at))::BIGINT AS month,
    COUNT(*) FILTER (WHERE status = 'open') AS open_count,
    COUNT(*) FILTER (WHERE status IN ('remediated', 'verified')) AS resolved_count
FROM vulnerabilities
WHERE created_at >= DATE_TRUNC('month', CURRENT_DATE) - INTERVAL '5 months'
GROUP BY DATE_TRUNC('month', created_at)
ORDER BY DATE_TRUNC('month', created_at) ASC;

-- name: GetAssetVulnerabilityCounts :many
SELECT 
    domain AS asset_name,
    COUNT(*) FILTER (WHERE severity = 'critical') AS critical,
    COUNT(*) FILTER (WHERE severity = 'high') AS high,
    COUNT(*) FILTER (WHERE severity = 'medium') AS medium,
    COUNT(*) FILTER (WHERE severity = 'low') AS low
FROM vulnerabilities
GROUP BY domain
ORDER BY critical DESC, high DESC, medium DESC, low DESC
LIMIT 10;

-- name: GetRecentSecurityEvents :many
SELECT 
    id,
    action AS description,
    entity_type AS event_type,
    created_at
FROM activity_logs
ORDER BY created_at DESC
LIMIT $1;

-- name: GetRemediationUpdates :many
SELECT 
    v.id,
    v.title AS vuln_title,
    'open' AS previous_status,
    v.status AS new_status,
    v.remediated_date,
    COALESCE(u.full_name, v.assigned_to, 'Unassigned') AS assigned_team
FROM vulnerabilities v
LEFT JOIN users u ON v.discovered_by = u.id
WHERE v.status IN ('remediated', 'verified')
  AND v.remediated_date IS NOT NULL
ORDER BY v.remediated_date DESC
LIMIT $1;

