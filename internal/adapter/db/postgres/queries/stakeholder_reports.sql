-- ============================================================================
-- UC28: Get total reports count
-- ============================================================================
-- name: GetStakeholderTotalReportsCount :one
SELECT COUNT(*) AS count
FROM reports;

-- ============================================================================
-- UC29: Get under review reports count
-- ============================================================================
-- name: GetStakeholderUnderReviewReportsCount :one
SELECT COUNT(*) AS count
FROM reports
WHERE status = 'under_review';

-- ============================================================================
-- UC30: Get remediated reports count
-- ============================================================================
-- name: GetStakeholderRemediatedReportsCount :one
SELECT COUNT(*) AS count
FROM reports
WHERE status = 'remediated';

-- ============================================================================
-- UC32: List reports with status filter and pagination
-- ============================================================================
-- name: ListStakeholderReports :many
SELECT 
    r.id,
    r.title,
    u.full_name AS submitted_by,
    r.created_at AS submitted_date,
    p.name AS project_name,
    r.status,
    (SELECT COUNT(*) FROM report_vulnerabilities rv WHERE rv.report_id = r.id) AS vulnerabilities_count,
    (SELECT COUNT(*) FROM evidence_files ef WHERE ef.report_id = r.id) AS evidence_count
FROM reports r
LEFT JOIN users u ON r.submitted_by = u.id
LEFT JOIN projects p ON r.project_id = p.id
WHERE 
    (sqlc.narg('status')::TEXT IS NULL OR sqlc.narg('status')::TEXT = 'all' OR r.status = sqlc.narg('status')::TEXT)
ORDER BY r.created_at DESC
LIMIT $1 OFFSET $2;

-- ============================================================================
-- UC32: Get total count for pagination
-- ============================================================================
-- name: GetStakeholderReportsCount :one
SELECT COUNT(*) AS count
FROM reports r
WHERE 
    (sqlc.narg('status')::TEXT IS NULL OR sqlc.narg('status')::TEXT = 'all' OR r.status = sqlc.narg('status')::TEXT);

-- ============================================================================
-- UC37: Get report details by ID
-- ============================================================================
-- name: GetStakeholderReportByID :one
SELECT 
    r.id,
    r.title,
    u.full_name AS submitted_by,
    r.created_at AS submission_date,
    p.name AS project_name,
    r.status,
    r.executive_summary
FROM reports r
LEFT JOIN users u ON r.submitted_by = u.id
LEFT JOIN projects p ON r.project_id = p.id
WHERE r.id = $1;

-- ============================================================================
-- UC38: Get vulnerabilities for a report
-- ============================================================================
-- name: GetStakeholderReportVulnerabilities :many
SELECT 
    rv.id,
    rv.vulnerability_title AS title,
    rv.severity,
    rv.asset_target AS domain,
    'open' AS status,
    rv.vulnerability_description AS description,
    rv.remediation_recommendation AS remediation
FROM report_vulnerabilities rv
WHERE rv.report_id = $1
ORDER BY 
    CASE rv.severity
        WHEN 'critical' THEN 1
        WHEN 'high' THEN 2
        WHEN 'medium' THEN 3
        WHEN 'low' THEN 4
    END,
    rv.created_at DESC;

-- ============================================================================
-- UC39: Get evidence files for a report with pagination
-- ============================================================================
-- name: GetStakeholderReportEvidenceFiles :many
SELECT 
    id,
    file_name,
    file_size,
    created_at AS upload_date,
    file_path
FROM evidence_files
WHERE report_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- ============================================================================
-- UC39: Get evidence files count for pagination
-- ============================================================================
-- name: GetStakeholderReportEvidenceFilesCount :one
SELECT COUNT(*) AS count
FROM evidence_files
WHERE report_id = $1;

-- ============================================================================
-- UC41: Get evidence file by ID for download
-- ============================================================================
-- name: GetStakeholderEvidenceFileByID :one
SELECT 
    id,
    file_name,
    file_size,
    created_at AS upload_date,
    file_path
FROM evidence_files
WHERE id = $1;

-- ============================================================================
-- UC42: Get report file path for download
-- ============================================================================
-- name: GetStakeholderReportFilePath :one
SELECT file_path
FROM reports
WHERE id = $1;

