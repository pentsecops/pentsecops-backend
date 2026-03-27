-- name: GetDomainsStats :one
SELECT 
    COUNT(DISTINCT d.id) as total_domains,
    COALESCE(AVG(d.risk_score), 0) as avg_risk_score,
    COALESCE(SUM(CASE WHEN v.severity = 'critical' THEN 1 ELSE 0 END), 0) as critical_issues,
    COALESCE(
        CAST(
            SUM(CASE 
                WHEN v.due_date IS NOT NULL 
                AND v.status IN ('remediated', 'verified') 
                AND v.remediated_date IS NOT NULL 
                AND v.remediated_date <= v.due_date 
                THEN 1 
                ELSE 0 
            END) * 100.0 / NULLIF(SUM(CASE WHEN v.due_date IS NOT NULL THEN 1 ELSE 0 END), 0)
        AS DECIMAL(5,2)), 
        0
    ) as sla_compliance_percent
FROM domains d
LEFT JOIN vulnerabilities v ON d.domain_name = v.domain
WHERE d.is_active = true;

-- name: ListDomains :many
SELECT 
    d.id,
    d.domain_name,
    d.ip_address,
    d.description,
    d.risk_score,
    d.is_active,
    d.last_scanned,
    d.created_at,
    d.updated_at,
    COALESCE(COUNT(v.id), 0) as total_vulnerabilities,
    COALESCE(SUM(CASE WHEN v.severity = 'critical' THEN 1 ELSE 0 END), 0) as critical_count,
    COALESCE(SUM(CASE WHEN v.severity = 'high' THEN 1 ELSE 0 END), 0) as high_count,
    COALESCE(SUM(CASE WHEN v.severity = 'medium' THEN 1 ELSE 0 END), 0) as medium_count,
    COALESCE(SUM(CASE WHEN v.severity = 'low' THEN 1 ELSE 0 END), 0) as low_count,
    COALESCE(SUM(CASE WHEN v.status = 'open' THEN 1 ELSE 0 END), 0) as open_issues,
    COALESCE(
        CAST(
            SUM(CASE 
                WHEN v.due_date IS NOT NULL 
                AND v.status IN ('remediated', 'verified') 
                AND v.remediated_date IS NOT NULL 
                AND v.remediated_date <= v.due_date 
                THEN 1 
                ELSE 0 
            END) * 100.0 / NULLIF(SUM(CASE WHEN v.due_date IS NOT NULL THEN 1 ELSE 0 END), 0)
        AS DECIMAL(5,2)), 
        0
    ) as sla_compliance
FROM domains d
LEFT JOIN vulnerabilities v ON d.domain_name = v.domain
WHERE d.is_active = true
GROUP BY d.id, d.domain_name, d.ip_address, d.description, d.risk_score, d.is_active, d.last_scanned, d.created_at, d.updated_at
ORDER BY d.risk_score DESC NULLS LAST, d.domain_name ASC
LIMIT $1 OFFSET $2;

-- name: CountDomains :one
SELECT COUNT(DISTINCT d.id) as total
FROM domains d
WHERE d.is_active = true;

-- name: GetDomainByID :one
SELECT 
    d.id,
    d.domain_name,
    d.ip_address,
    d.description,
    d.risk_score,
    d.is_active,
    d.last_scanned,
    d.created_at,
    d.updated_at
FROM domains d
WHERE d.id = $1 AND d.is_active = true;

-- name: GetSecurityMetrics :many
SELECT 
    metric_name,
    COALESCE(AVG(metric_value), 0) as avg_value
FROM security_metrics
WHERE domain_id IN (
    SELECT id FROM domains WHERE is_active = true
)
AND metric_name IN ('authentication', 'authorization', 'input_validation', 'encryption', 'configuration', 'network_security')
GROUP BY metric_name;

-- name: GetSLABreachAnalysis :many
SELECT 
    d.domain_name,
    COALESCE(
        CAST(
            SUM(CASE 
                WHEN v.due_date IS NOT NULL 
                AND v.status IN ('remediated', 'verified') 
                AND v.remediated_date IS NOT NULL 
                AND v.remediated_date <= v.due_date 
                THEN 1 
                ELSE 0 
            END) * 100.0 / NULLIF(SUM(CASE WHEN v.due_date IS NOT NULL THEN 1 ELSE 0 END), 0)
        AS DECIMAL(5,2)), 
        0
    ) as sla_compliance_percent
FROM domains d
LEFT JOIN vulnerabilities v ON d.domain_name = v.domain
WHERE d.is_active = true
GROUP BY d.domain_name
HAVING SUM(CASE WHEN v.due_date IS NOT NULL THEN 1 ELSE 0 END) > 0
ORDER BY sla_compliance_percent ASC, d.domain_name ASC;

-- name: CreateDomain :one
INSERT INTO domains (
    id,
    domain_name,
    ip_address,
    description,
    risk_score,
    is_active,
    last_scanned,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING id, domain_name, ip_address, description, risk_score, is_active, last_scanned, created_at, updated_at;

-- name: UpdateDomain :one
UPDATE domains
SET 
    domain_name = $2,
    ip_address = $3,
    description = $4,
    risk_score = $5,
    last_scanned = $6,
    updated_at = $7
WHERE id = $1 AND is_active = true
RETURNING id, domain_name, ip_address, description, risk_score, is_active, last_scanned, created_at, updated_at;

-- name: DeleteDomain :exec
UPDATE domains
SET is_active = false, updated_at = $2
WHERE id = $1;

-- name: CreateSecurityMetric :one
INSERT INTO security_metrics (
    id,
    domain_id,
    metric_name,
    metric_value,
    measured_at,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING id, domain_id, metric_name, metric_value, measured_at, created_at;

