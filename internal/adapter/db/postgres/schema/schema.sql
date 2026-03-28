-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'pentester', 'stakeholder')),
    is_active BOOLEAN DEFAULT true,
    force_password_change BOOLEAN DEFAULT false,
    failed_login_attempts INTEGER DEFAULT 0,
    last_failed_login TIMESTAMP,
    account_locked_until TIMESTAMP,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_account_locked_until ON users(account_locked_until);

-- Projects table
CREATE TABLE projects (
    id UUID PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('web', 'network', 'api', 'mobile')),
    assigned_to UUID REFERENCES users(id) ON DELETE SET NULL,
    deadline TIMESTAMP NOT NULL,
    scope TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'in_progress', 'completed')),
    current_phase VARCHAR(100) DEFAULT 'pending' CHECK (current_phase IN ('pending', 'reconnaissance', 'scanning', 'exploitation', 'post_exploitation', 'reporting', 'completed')),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_projects_name ON projects(name);
CREATE INDEX idx_projects_type ON projects(type);
CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_current_phase ON projects(current_phase);
CREATE INDEX idx_projects_assigned_to ON projects(assigned_to);
CREATE INDEX idx_projects_deadline ON projects(deadline);
CREATE INDEX idx_projects_created_at ON projects(created_at);

-- Tasks table
CREATE TABLE tasks (
    id UUID PRIMARY KEY,
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'to_do' CHECK (status IN ('to_do', 'in_progress', 'done')),
    priority VARCHAR(50) CHECK (priority IN ('low', 'medium', 'high', 'critical')),
    assigned_to UUID REFERENCES users(id) ON DELETE SET NULL,
    deadline TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tasks_project ON tasks(project_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_assigned_to ON tasks(assigned_to);
CREATE INDEX idx_tasks_deadline ON tasks(deadline);
CREATE INDEX idx_tasks_created_at ON tasks(created_at);

-- Domains table
CREATE TABLE domains (
    id UUID PRIMARY KEY,
    domain_name VARCHAR(255) UNIQUE NOT NULL,
    ip_address VARCHAR(45),
    description TEXT,
    risk_score DECIMAL(3,1) CHECK (risk_score >= 0 AND risk_score <= 10),
    is_active BOOLEAN DEFAULT true,
    last_scanned TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_domains_domain_name ON domains(domain_name);
CREATE INDEX idx_domains_is_active ON domains(is_active);

-- Vulnerabilities table
CREATE TABLE vulnerabilities (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    severity VARCHAR(50) NOT NULL CHECK (severity IN ('critical', 'high', 'medium', 'low')),
    domain VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'in_progress', 'remediated', 'verified')),
    discovered_date TIMESTAMP,
    due_date TIMESTAMP NOT NULL,
    assigned_to VARCHAR(255),
    cvss_score DECIMAL(3,1) CHECK (cvss_score >= 0 AND cvss_score <= 10),
    cwe_id VARCHAR(50),
    domain_id UUID REFERENCES domains(id) ON DELETE SET NULL,
    project_id UUID REFERENCES projects(id) ON DELETE SET NULL,
    discovered_by UUID REFERENCES users(id) ON DELETE SET NULL,
    remediation_notes TEXT,
    remediated_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_vulnerabilities_severity ON vulnerabilities(severity);
CREATE INDEX idx_vulnerabilities_status ON vulnerabilities(status);
CREATE INDEX idx_vulnerabilities_domain ON vulnerabilities(domain);
CREATE INDEX idx_vulnerabilities_domain_id ON vulnerabilities(domain_id);
CREATE INDEX idx_vulnerabilities_project ON vulnerabilities(project_id);
CREATE INDEX idx_vulnerabilities_due_date ON vulnerabilities(due_date);
CREATE INDEX idx_vulnerabilities_remediated_date ON vulnerabilities(remediated_date);

-- Security Metrics table
CREATE TABLE security_metrics (
    id UUID PRIMARY KEY,
    domain_id UUID REFERENCES domains(id) ON DELETE CASCADE,
    metric_name VARCHAR(100) NOT NULL,
    metric_value DECIMAL(5,2),
    measured_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_security_metrics_domain ON security_metrics(domain_id);
CREATE INDEX idx_security_metrics_name ON security_metrics(metric_name);
CREATE INDEX idx_security_metrics_measured_at ON security_metrics(measured_at);

-- Activity Logs table
CREATE TABLE activity_logs (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    user_name VARCHAR(255),
    action VARCHAR(255) NOT NULL,
    entity_type VARCHAR(100),
    entity_id UUID,
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_activity_logs_user ON activity_logs(user_id);
CREATE INDEX idx_activity_logs_created_at ON activity_logs(created_at);

-- Notifications table
CREATE TABLE notifications (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    sent_to VARCHAR(100) NOT NULL CHECK (sent_to IN ('all_pentesters', 'all_stakeholders', 'specific_user')),
    recipient_id UUID REFERENCES users(id) ON DELETE CASCADE,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'sent' CHECK (status IN ('sent', 'failed', 'pending')),
    type VARCHAR(50) CHECK (type IN ('info', 'warning', 'error', 'success')),
    is_read BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_sent_to ON notifications(sent_to);
CREATE INDEX idx_notifications_recipient ON notifications(recipient_id);
CREATE INDEX idx_notifications_created_by ON notifications(created_by);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);

-- Alerts table (for pentester alerts and pentester-to-admin communication)
CREATE TABLE alerts (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    alert_type VARCHAR(100) CHECK (alert_type IN ('report_review', 'project_assignment', 'system_notification', 'important_update')),
    priority VARCHAR(50) DEFAULT 'medium' CHECK (priority IN ('low', 'medium', 'high')),
    source VARCHAR(50) DEFAULT 'system' CHECK (source IN ('pentester', 'system', 'admin')),
    recipient_id UUID REFERENCES users(id) ON DELETE CASCADE,
    sender_id UUID REFERENCES users(id) ON DELETE SET NULL,
    is_read BOOLEAN DEFAULT false,
    is_dismissed BOOLEAN DEFAULT false,
    is_resolved BOOLEAN DEFAULT false,
    resolved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    resolved_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_alerts_alert_type ON alerts(alert_type);
CREATE INDEX idx_alerts_priority ON alerts(priority);
CREATE INDEX idx_alerts_source ON alerts(source);
CREATE INDEX idx_alerts_recipient_id ON alerts(recipient_id);
CREATE INDEX idx_alerts_sender_id ON alerts(sender_id);
CREATE INDEX idx_alerts_is_read ON alerts(is_read);
CREATE INDEX idx_alerts_is_dismissed ON alerts(is_dismissed);
CREATE INDEX idx_alerts_is_resolved ON alerts(is_resolved);
CREATE INDEX idx_alerts_created_at ON alerts(created_at);

-- Refresh Tokens table
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_refresh_tokens_user ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

-- Reports table
CREATE TABLE reports (
    id UUID PRIMARY KEY,
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    report_type VARCHAR(50) CHECK (report_type IN ('executive', 'technical', 'compliance', 'summary')),
    file_path VARCHAR(500),
    file_size BIGINT,
    submitted_by UUID REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(50) DEFAULT 'received' CHECK (status IN ('draft', 'received', 'under_review', 'shared', 'remediated', 'pending_review', 'approved', 'rejected', 'final', 'archived')),
    executive_summary TEXT,
    feedback TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_reports_project ON reports(project_id);
CREATE INDEX idx_reports_status ON reports(status);
CREATE INDEX idx_reports_submitted_by ON reports(submitted_by);
CREATE INDEX idx_reports_created_at ON reports(created_at);

-- Insert test users with Argon2 hashed passwords
-- Credentials: admin@test.com/admin123, pentester@test.com/pentester123, stakeholder@test.com/stakeholder123
INSERT INTO users (id, email, password_hash, full_name, role, is_active, created_at, updated_at)
VALUES 
    (gen_random_uuid(), 'admin@test.com', '6e011471524e6706b3b68178cbf4930eb048a46b7491370add2557115f0e4034', 'Admin', 'admin', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), 'pentester@test.com', 'b1894abd4b06d0410903ecec7a2bdbbc611c76ea7becd781aa2aba23fa3fd7e0', 'Pentester', 'pentester', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    (gen_random_uuid(), 'stakeholder@test.com', '2861c866377e97a32a0507da96838f179e1115cb75f8057ad2a86700fbfa655e', 'Stakeholder', 'stakeholder', true, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (email) DO UPDATE
SET password_hash = EXCLUDED.password_hash,
    failed_login_attempts = 0,
    account_locked_until = NULL,
    last_failed_login = NULL;



-- Report Vulnerabilities table (links reports to vulnerability details)
CREATE TABLE report_vulnerabilities (
    id UUID PRIMARY KEY,
    report_id UUID REFERENCES reports(id) ON DELETE CASCADE,
    asset_target VARCHAR(500) NOT NULL,
    vulnerability_title VARCHAR(500) NOT NULL,
    severity VARCHAR(50) NOT NULL CHECK (severity IN ('critical', 'high', 'medium', 'low')),
    attack_vector TEXT,
    vulnerability_description TEXT NOT NULL,
    evidence_poc TEXT,
    remediation_recommendation TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_report_vulnerabilities_report ON report_vulnerabilities(report_id);
CREATE INDEX idx_report_vulnerabilities_severity ON report_vulnerabilities(severity);
CREATE INDEX idx_report_vulnerabilities_created_at ON report_vulnerabilities(created_at);

-- Evidence Files table
CREATE TABLE evidence_files (
    id UUID PRIMARY KEY,
    report_id UUID REFERENCES reports(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    file_type VARCHAR(100),
    uploaded_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_evidence_files_report ON evidence_files(report_id);
CREATE INDEX idx_evidence_files_uploaded_by ON evidence_files(uploaded_by);
CREATE INDEX idx_evidence_files_created_at ON evidence_files(created_at);

-- Project Assets table
CREATE TABLE project_assets (
    id UUID PRIMARY KEY,
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    domain_id UUID REFERENCES domains(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_project_assets_project ON project_assets(project_id);
CREATE INDEX idx_project_assets_domain ON project_assets(domain_id);

-- Project Requirements table
CREATE TABLE project_requirements (
    id UUID PRIMARY KEY,
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    requirement_text TEXT NOT NULL,
    is_completed BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_project_requirements_project ON project_requirements(project_id);
CREATE INDEX idx_project_requirements_is_completed ON project_requirements(is_completed);
