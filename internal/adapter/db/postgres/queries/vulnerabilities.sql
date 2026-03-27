-- name: CreateVulnerability :one
INSERT INTO vulnerabilities (
    id,
    title,
    description,
    severity,
    domain,
    status,
    discovered_date,
    due_date,
    assigned_to,
    cvss_score,
    cwe_id,
    domain_id,
    project_id,
    discovered_by,
    remediation_notes,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
) RETURNING *;

-- name: GetVulnerabilityByID :one
SELECT * FROM vulnerabilities
WHERE id = $1;

-- name: UpdateVulnerability :one
UPDATE vulnerabilities
SET
    title = $2,
    description = $3,
    severity = $4,
    domain = $5,
    status = $6,
    discovered_date = $7,
    due_date = $8,
    assigned_to = $9,
    cvss_score = $10,
    cwe_id = $11,
    remediation_notes = $12,
    updated_at = $13
WHERE id = $1
RETURNING *;

-- name: DeleteVulnerability :exec
DELETE FROM vulnerabilities
WHERE id = $1;

-- name: ListVulnerabilities :many
SELECT * FROM vulnerabilities
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountVulnerabilities :one
SELECT COUNT(*) FROM vulnerabilities;

-- name: SearchAndFilterVulnerabilities :many
SELECT * FROM vulnerabilities
WHERE
    (LOWER(title) LIKE LOWER($1) OR LOWER(domain) LIKE LOWER($1))
    AND ($2 = '' OR severity = $2)
    AND ($3 = '' OR status = $3)
ORDER BY created_at DESC
LIMIT $4 OFFSET $5;

-- name: CountSearchAndFilterVulnerabilities :one
SELECT COUNT(*) FROM vulnerabilities
WHERE
    (LOWER(title) LIKE LOWER($1) OR LOWER(domain) LIKE LOWER($1))
    AND ($2 = '' OR severity = $2)
    AND ($3 = '' OR status = $3);

-- name: GetVulnerabilityStats :one
SELECT
    COUNT(*) as total,
    COUNT(*) FILTER (WHERE severity = 'critical') as critical,
    COUNT(*) FILTER (WHERE severity = 'high') as high,
    COUNT(*) FILTER (WHERE severity = 'medium') as medium,
    COUNT(*) FILTER (WHERE severity = 'low') as low,
    COUNT(*) FILTER (WHERE status = 'open') as open,
    COUNT(*) FILTER (WHERE status = 'in_progress') as in_progress,
    COUNT(*) FILTER (WHERE status = 'remediated') as remediated,
    COUNT(*) FILTER (WHERE status = 'verified') as verified
FROM vulnerabilities;

-- name: GetSLACompliance :one
SELECT
    COUNT(*) FILTER (WHERE severity = 'critical' AND due_date < NOW() AND status NOT IN ('remediated', 'verified')) as critical_overdue,
    COUNT(*) FILTER (WHERE severity = 'high' AND due_date < NOW() + INTERVAL '7 days' AND due_date > NOW() AND status NOT IN ('remediated', 'verified')) as high_approaching,
    COUNT(*) FILTER (WHERE due_date IS NOT NULL) as total_with_due_date,
    COUNT(*) FILTER (WHERE status IN ('remediated', 'verified') AND due_date >= created_at) as remediated_on_time
FROM vulnerabilities;

-- name: ExportVulnerabilities :many
SELECT
    id,
    title,
    severity,
    domain,
    status,
    discovered_date,
    due_date,
    assigned_to
FROM vulnerabilities
WHERE
    (LOWER(title) LIKE LOWER($1) OR LOWER(domain) LIKE LOWER($1))
    AND ($2 = '' OR severity = $2)
    AND ($3 = '' OR status = $3)
ORDER BY created_at DESC;

