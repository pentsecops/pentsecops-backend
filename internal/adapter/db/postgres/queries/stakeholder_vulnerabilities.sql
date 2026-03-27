-- ============================================================================
-- UC11: Get critical vulnerabilities count
-- ============================================================================
-- name: GetStakeholderCriticalVulnerabilitiesCount :one
SELECT COUNT(*) AS count
FROM vulnerabilities
WHERE severity = 'critical';

-- ============================================================================
-- UC12: Get high severity vulnerabilities count
-- ============================================================================
-- name: GetStakeholderHighSeverityVulnerabilitiesCount :one
SELECT COUNT(*) AS count
FROM vulnerabilities
WHERE severity = 'high';

-- ============================================================================
-- UC13: Get open issues count (open + in_progress)
-- ============================================================================
-- name: GetStakeholderOpenIssuesCount :one
SELECT COUNT(*) AS count
FROM vulnerabilities
WHERE status IN ('open', 'in_progress');

-- ============================================================================
-- UC14: Get remediation count (remediated + verified)
-- ============================================================================
-- name: GetStakeholderRemediationCount :one
SELECT COUNT(*) AS count
FROM vulnerabilities
WHERE status IN ('remediated', 'verified');

-- ============================================================================
-- UC19: List vulnerabilities with search, filters, and pagination
-- ============================================================================
-- name: ListStakeholderVulnerabilities :many
SELECT 
    id,
    title,
    severity,
    domain,
    status,
    created_at AS discovered_date,
    due_date,
    assigned_to
FROM vulnerabilities
WHERE 
    (sqlc.narg('search')::TEXT IS NULL OR 
     title ILIKE '%' || sqlc.narg('search')::TEXT || '%' OR 
     domain ILIKE '%' || sqlc.narg('search')::TEXT || '%')
    AND (sqlc.narg('severity')::TEXT IS NULL OR sqlc.narg('severity')::TEXT = 'all' OR severity = sqlc.narg('severity')::TEXT)
    AND (sqlc.narg('status')::TEXT IS NULL OR sqlc.narg('status')::TEXT = 'all' OR status = sqlc.narg('status')::TEXT)
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- ============================================================================
-- UC19: Get total count for pagination
-- ============================================================================
-- name: GetStakeholderVulnerabilitiesCount :one
SELECT COUNT(*) AS count
FROM vulnerabilities
WHERE 
    (sqlc.narg('search')::TEXT IS NULL OR 
     title ILIKE '%' || sqlc.narg('search')::TEXT || '%' OR 
     domain ILIKE '%' || sqlc.narg('search')::TEXT || '%')
    AND (sqlc.narg('severity')::TEXT IS NULL OR sqlc.narg('severity')::TEXT = 'all' OR severity = sqlc.narg('severity')::TEXT)
    AND (sqlc.narg('status')::TEXT IS NULL OR sqlc.narg('status')::TEXT = 'all' OR status = sqlc.narg('status')::TEXT);

-- ============================================================================
-- UC24: Export vulnerabilities to CSV (all matching filters)
-- ============================================================================
-- name: ExportStakeholderVulnerabilities :many
SELECT 
    id,
    title,
    severity,
    domain,
    status,
    created_at AS discovered_date,
    due_date,
    assigned_to
FROM vulnerabilities
WHERE 
    (sqlc.narg('search')::TEXT IS NULL OR 
     title ILIKE '%' || sqlc.narg('search')::TEXT || '%' OR 
     domain ILIKE '%' || sqlc.narg('search')::TEXT || '%')
    AND (sqlc.narg('severity')::TEXT IS NULL OR sqlc.narg('severity')::TEXT = 'all' OR severity = sqlc.narg('severity')::TEXT)
    AND (sqlc.narg('status')::TEXT IS NULL OR sqlc.narg('status')::TEXT = 'all' OR status = sqlc.narg('status')::TEXT)
ORDER BY created_at DESC;

-- ============================================================================
-- UC25: Get critical vulnerabilities overdue count
-- ============================================================================
-- name: GetStakeholderCriticalOverdueCount :one
SELECT COUNT(*) AS count
FROM vulnerabilities
WHERE severity = 'critical'
  AND status IN ('open', 'in_progress')
  AND due_date IS NOT NULL
  AND due_date < CURRENT_DATE;

-- ============================================================================
-- UC26: Get high severity approaching deadline count (within 3 days)
-- ============================================================================
-- name: GetStakeholderHighApproachingDeadlineCount :one
SELECT COUNT(*) AS count
FROM vulnerabilities
WHERE severity = 'high'
  AND status IN ('open', 'in_progress')
  AND due_date IS NOT NULL
  AND due_date BETWEEN CURRENT_DATE AND CURRENT_DATE + INTERVAL '3 days';

-- ============================================================================
-- UC27: Get overall SLA compliance data
-- ============================================================================
-- name: GetStakeholderSLAComplianceData :one
SELECT 
    COUNT(*) FILTER (WHERE due_date IS NOT NULL) AS total_with_due_date,
    COUNT(*) FILTER (WHERE due_date IS NOT NULL AND status IN ('remediated', 'verified') AND remediated_date <= due_date) AS remediated_on_time
FROM vulnerabilities;

