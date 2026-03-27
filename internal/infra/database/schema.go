package database

import (
	"database/sql"
	"fmt"
	"log"
)

// InitializeSchema creates all database tables if they don't exist
func InitializeSchema(db *sql.DB) error {
	log.Println("Initializing database schema...")

	// Execute schema creation
	schema := `
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('admin', 'pentester', 'stakeholder')),
    is_active BOOLEAN DEFAULT true,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add new columns to users table if they don't exist (for existing databases)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='force_password_change') THEN
        ALTER TABLE users ADD COLUMN force_password_change BOOLEAN DEFAULT false;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='failed_login_attempts') THEN
        ALTER TABLE users ADD COLUMN failed_login_attempts INTEGER DEFAULT 0;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='last_failed_login') THEN
        ALTER TABLE users ADD COLUMN last_failed_login TIMESTAMP;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='users' AND column_name='account_locked_until') THEN
        ALTER TABLE users ADD COLUMN account_locked_until TIMESTAMP;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);
CREATE INDEX IF NOT EXISTS idx_users_account_locked_until ON users(account_locked_until);

-- Projects table
CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
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

-- Add current_phase column if it doesn't exist (for existing databases)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='projects' AND column_name='current_phase') THEN
        ALTER TABLE projects ADD COLUMN current_phase VARCHAR(100) DEFAULT 'pending' CHECK (current_phase IN ('pending', 'reconnaissance', 'scanning', 'exploitation', 'post_exploitation', 'reporting', 'completed'));
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_projects_name ON projects(name);
CREATE INDEX IF NOT EXISTS idx_projects_type ON projects(type);
CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);
CREATE INDEX IF NOT EXISTS idx_projects_current_phase ON projects(current_phase);
CREATE INDEX IF NOT EXISTS idx_projects_assigned_to ON projects(assigned_to);
CREATE INDEX IF NOT EXISTS idx_projects_deadline ON projects(deadline);
CREATE INDEX IF NOT EXISTS idx_projects_created_at ON projects(created_at);

-- Domains table
CREATE TABLE IF NOT EXISTS domains (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    domain_name VARCHAR(255) UNIQUE NOT NULL,
    ip_address VARCHAR(45),
    description TEXT,
    risk_score DECIMAL(3,1) CHECK (risk_score >= 0 AND risk_score <= 10),
    is_active BOOLEAN DEFAULT true,
    last_scanned TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_domains_domain_name ON domains(domain_name);
CREATE INDEX IF NOT EXISTS idx_domains_risk_score ON domains(risk_score);
CREATE INDEX IF NOT EXISTS idx_domains_is_active ON domains(is_active);

-- Vulnerabilities table
CREATE TABLE IF NOT EXISTS vulnerabilities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
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

CREATE INDEX IF NOT EXISTS idx_vulnerabilities_severity ON vulnerabilities(severity);
CREATE INDEX IF NOT EXISTS idx_vulnerabilities_status ON vulnerabilities(status);
CREATE INDEX IF NOT EXISTS idx_vulnerabilities_domain ON vulnerabilities(domain);
CREATE INDEX IF NOT EXISTS idx_vulnerabilities_domain_id ON vulnerabilities(domain_id);
CREATE INDEX IF NOT EXISTS idx_vulnerabilities_project ON vulnerabilities(project_id);
CREATE INDEX IF NOT EXISTS idx_vulnerabilities_discovered_by ON vulnerabilities(discovered_by);
CREATE INDEX IF NOT EXISTS idx_vulnerabilities_due_date ON vulnerabilities(due_date);
CREATE INDEX IF NOT EXISTS idx_vulnerabilities_remediated_date ON vulnerabilities(remediated_date);
CREATE INDEX IF NOT EXISTS idx_vulnerabilities_created_at ON vulnerabilities(created_at);

-- Tasks table
CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
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

CREATE INDEX IF NOT EXISTS idx_tasks_project ON tasks(project_id);
CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
CREATE INDEX IF NOT EXISTS idx_tasks_assigned_to ON tasks(assigned_to);
CREATE INDEX IF NOT EXISTS idx_tasks_deadline ON tasks(deadline);
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);

-- Activity Logs table (Enhanced for comprehensive logging)
-- Drop and recreate if structure is incompatible
DO $$
BEGIN
    -- Check if table exists and has incompatible structure
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name='activity_logs') THEN
        -- Check if user_email column exists, if not, drop and recreate
        IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='activity_logs' AND column_name='user_email') THEN
            DROP TABLE IF EXISTS activity_logs CASCADE;
        END IF;
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS activity_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    user_email VARCHAR(255) NOT NULL DEFAULT 'unknown@system.local',
    user_role VARCHAR(50) NOT NULL DEFAULT 'system',
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50),
    entity_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address VARCHAR(45) NOT NULL DEFAULT 'unknown',
    user_agent TEXT NOT NULL DEFAULT 'unknown',
    endpoint VARCHAR(255) NOT NULL DEFAULT '/unknown',
    method VARCHAR(10) NOT NULL DEFAULT 'UNKNOWN',
    status_code INTEGER NOT NULL DEFAULT 200,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Enhanced indexes for activity logs
CREATE INDEX IF NOT EXISTS idx_activity_logs_user_id ON activity_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_activity_logs_action ON activity_logs(action);
CREATE INDEX IF NOT EXISTS idx_activity_logs_entity_type ON activity_logs(entity_type);
CREATE INDEX IF NOT EXISTS idx_activity_logs_entity_id ON activity_logs(entity_id);
CREATE INDEX IF NOT EXISTS idx_activity_logs_created_at ON activity_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_activity_logs_user_email ON activity_logs(user_email);
CREATE INDEX IF NOT EXISTS idx_activity_logs_user_action ON activity_logs(user_id, action);
CREATE INDEX IF NOT EXISTS idx_activity_logs_entity_type_id ON activity_logs(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_activity_logs_date_range ON activity_logs(created_at, user_id);

-- Activity logs table is now created above with proper structure

-- Reports table
CREATE TABLE IF NOT EXISTS reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
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

-- Add submitted_by, feedback, and executive_summary columns if they don't exist (for existing databases)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='reports' AND column_name='submitted_by') THEN
        ALTER TABLE reports ADD COLUMN submitted_by UUID REFERENCES users(id) ON DELETE SET NULL;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='reports' AND column_name='feedback') THEN
        ALTER TABLE reports ADD COLUMN feedback TEXT;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='reports' AND column_name='executive_summary') THEN
        ALTER TABLE reports ADD COLUMN executive_summary TEXT;
    END IF;
    -- Drop old column if exists
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='reports' AND column_name='generated_by') THEN
        ALTER TABLE reports DROP COLUMN generated_by;
    END IF;
    -- Update status constraint to include stakeholder statuses
    ALTER TABLE reports DROP CONSTRAINT IF EXISTS reports_status_check;
    ALTER TABLE reports ADD CONSTRAINT reports_status_check CHECK (status IN ('draft', 'received', 'under_review', 'shared', 'remediated', 'pending_review', 'approved', 'rejected', 'final', 'archived'));
END $$;

CREATE INDEX IF NOT EXISTS idx_reports_project ON reports(project_id);
CREATE INDEX IF NOT EXISTS idx_reports_type ON reports(report_type);
CREATE INDEX IF NOT EXISTS idx_reports_status ON reports(status);
CREATE INDEX IF NOT EXISTS idx_reports_submitted_by ON reports(submitted_by);
CREATE INDEX IF NOT EXISTS idx_reports_created_at ON reports(created_at);

-- Report Vulnerabilities table (links reports to vulnerability details)
CREATE TABLE IF NOT EXISTS report_vulnerabilities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
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

CREATE INDEX IF NOT EXISTS idx_report_vulnerabilities_report ON report_vulnerabilities(report_id);
CREATE INDEX IF NOT EXISTS idx_report_vulnerabilities_severity ON report_vulnerabilities(severity);
CREATE INDEX IF NOT EXISTS idx_report_vulnerabilities_created_at ON report_vulnerabilities(created_at);

-- Evidence Files table
CREATE TABLE IF NOT EXISTS evidence_files (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    vulnerability_id UUID REFERENCES vulnerabilities(id) ON DELETE CASCADE,
    report_id UUID REFERENCES reports(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_type VARCHAR(100),
    file_size BIGINT,
    uploaded_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add report_id column if it doesn't exist (for existing databases)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='evidence_files' AND column_name='report_id') THEN
        ALTER TABLE evidence_files ADD COLUMN report_id UUID REFERENCES reports(id) ON DELETE CASCADE;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_evidence_files_vulnerability ON evidence_files(vulnerability_id);
CREATE INDEX IF NOT EXISTS idx_evidence_files_report ON evidence_files(report_id);
CREATE INDEX IF NOT EXISTS idx_evidence_files_uploaded_by ON evidence_files(uploaded_by);

-- Notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
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

CREATE INDEX IF NOT EXISTS idx_notifications_sent_to ON notifications(sent_to);
CREATE INDEX IF NOT EXISTS idx_notifications_recipient ON notifications(recipient_id);
CREATE INDEX IF NOT EXISTS idx_notifications_created_by ON notifications(created_by);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);

-- Alerts table (for pentester alerts and pentester-to-admin communication)
CREATE TABLE IF NOT EXISTS alerts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    alert_type VARCHAR(100),
    priority VARCHAR(50) DEFAULT 'medium',
    source VARCHAR(50) DEFAULT 'system',
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

-- Add columns to existing alerts table if they don't exist
DO $$
BEGIN
    -- Drop old severity column if exists (we use priority instead)
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='alerts' AND column_name='severity') THEN
        ALTER TABLE alerts DROP COLUMN severity;
    END IF;
    -- Drop old created_by column if exists (we use sender_id instead)
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='alerts' AND column_name='created_by') THEN
        ALTER TABLE alerts DROP COLUMN created_by;
    END IF;
    -- Add new columns
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='alerts' AND column_name='recipient_id') THEN
        ALTER TABLE alerts ADD COLUMN recipient_id UUID REFERENCES users(id) ON DELETE CASCADE;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='alerts' AND column_name='sender_id') THEN
        ALTER TABLE alerts ADD COLUMN sender_id UUID REFERENCES users(id) ON DELETE SET NULL;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='alerts' AND column_name='priority') THEN
        ALTER TABLE alerts ADD COLUMN priority VARCHAR(50) DEFAULT 'medium';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='alerts' AND column_name='is_read') THEN
        ALTER TABLE alerts ADD COLUMN is_read BOOLEAN DEFAULT false;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='alerts' AND column_name='is_dismissed') THEN
        ALTER TABLE alerts ADD COLUMN is_dismissed BOOLEAN DEFAULT false;
    END IF;
    -- Update existing alert_type values to match new constraint
    UPDATE alerts SET alert_type = 'system_notification' WHERE alert_type NOT IN ('report_review', 'project_assignment', 'system_notification', 'important_update') OR alert_type IS NULL;

    -- Update alert_type constraint
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='alerts' AND column_name='alert_type') THEN
        ALTER TABLE alerts DROP CONSTRAINT IF EXISTS alerts_alert_type_check;
        ALTER TABLE alerts ADD CONSTRAINT alerts_alert_type_check CHECK (alert_type IN ('report_review', 'project_assignment', 'system_notification', 'important_update'));
    END IF;
    -- Add priority constraint if not exists
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'alerts_priority_check') THEN
        ALTER TABLE alerts ADD CONSTRAINT alerts_priority_check CHECK (priority IN ('low', 'medium', 'high'));
    END IF;
    -- Update source constraint
    ALTER TABLE alerts DROP CONSTRAINT IF EXISTS alerts_source_check;
    ALTER TABLE alerts ADD CONSTRAINT alerts_source_check CHECK (source IN ('pentester', 'system', 'admin'));
END $$;

-- Create indexes for alerts table
CREATE INDEX IF NOT EXISTS idx_alerts_alert_type ON alerts(alert_type);
CREATE INDEX IF NOT EXISTS idx_alerts_priority ON alerts(priority);
CREATE INDEX IF NOT EXISTS idx_alerts_source ON alerts(source);
CREATE INDEX IF NOT EXISTS idx_alerts_recipient_id ON alerts(recipient_id);
CREATE INDEX IF NOT EXISTS idx_alerts_sender_id ON alerts(sender_id);
CREATE INDEX IF NOT EXISTS idx_alerts_is_read ON alerts(is_read);
CREATE INDEX IF NOT EXISTS idx_alerts_is_dismissed ON alerts(is_dismissed);
CREATE INDEX IF NOT EXISTS idx_alerts_is_resolved ON alerts(is_resolved);
CREATE INDEX IF NOT EXISTS idx_alerts_created_at ON alerts(created_at);

-- Refresh Tokens table
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

-- Security Metrics table
CREATE TABLE IF NOT EXISTS security_metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    domain_id UUID REFERENCES domains(id) ON DELETE CASCADE,
    metric_name VARCHAR(100) NOT NULL,
    metric_value DECIMAL(5,2),
    measured_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_security_metrics_domain ON security_metrics(domain_id);
CREATE INDEX IF NOT EXISTS idx_security_metrics_name ON security_metrics(metric_name);
CREATE INDEX IF NOT EXISTS idx_security_metrics_measured_at ON security_metrics(measured_at);

-- Refresh Tokens table
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

-- Project Assets table
CREATE TABLE IF NOT EXISTS project_assets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    domain_id UUID REFERENCES domains(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_project_assets_project ON project_assets(project_id);
CREATE INDEX IF NOT EXISTS idx_project_assets_domain ON project_assets(domain_id);

-- Project Requirements table
CREATE TABLE IF NOT EXISTS project_requirements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    requirement_text TEXT NOT NULL,
    is_completed BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_project_requirements_project ON project_requirements(project_id);
CREATE INDEX IF NOT EXISTS idx_project_requirements_is_completed ON project_requirements(is_completed);
`

	// Execute the schema
	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}

	// Create triggers for updated_at
	triggers := `
-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for all tables with updated_at
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_users_updated_at') THEN
        CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_projects_updated_at') THEN
        CREATE TRIGGER update_projects_updated_at BEFORE UPDATE ON projects FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_domains_updated_at') THEN
        CREATE TRIGGER update_domains_updated_at BEFORE UPDATE ON domains FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_vulnerabilities_updated_at') THEN
        CREATE TRIGGER update_vulnerabilities_updated_at BEFORE UPDATE ON vulnerabilities FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_tasks_updated_at') THEN
        CREATE TRIGGER update_tasks_updated_at BEFORE UPDATE ON tasks FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_reports_updated_at') THEN
        CREATE TRIGGER update_reports_updated_at BEFORE UPDATE ON reports FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'update_alerts_updated_at') THEN
        CREATE TRIGGER update_alerts_updated_at BEFORE UPDATE ON alerts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;
`

	_, err = db.Exec(triggers)
	if err != nil {
		return fmt.Errorf("failed to create triggers: %w", err)
	}

	// Insert test users with proper Argon2 hashes
	// Delete existing test users first
	_, err = db.Exec("DELETE FROM users WHERE email IN ('admin@test.com', 'pentester@test.com', 'stakeholder@test.com')")
	if err != nil {
		log.Printf("Warning: Failed to delete existing test users: %v", err)
	}

	// Insert test users
	testUsers := []struct {
		email    string
		hash     string
		name     string
		role     string
	}{
		{"admin@test.com", "6e011471524e6706b3b68178cbf4930eb048a46b7491370add2557115f0e4034", "Test Admin", "admin"},
		{"pentester@test.com", "b1894abd4b06d0410903ecec7a2bdbbc611c76ea7becd781aa2aba23fa3fd7e0", "Test Pentester", "pentester"},
		{"stakeholder@test.com", "2861c866377e97a32a0507da96838f179e1115cb75f8057ad2a86700fbfa655e", "Test Stakeholder", "stakeholder"},
	}

	for _, user := range testUsers {
		_, err = db.Exec(`
			INSERT INTO users (id, email, password_hash, full_name, role, is_active, failed_login_attempts, created_at, updated_at)
			VALUES (uuid_generate_v4(), $1, $2, $3, $4, true, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`,
			user.email, user.hash, user.name, user.role)
		if err != nil {
			log.Printf("Failed to insert user %s: %v", user.email, err)
		} else {
			log.Printf("✓ Inserted test user: %s (%s)", user.email, user.role)
		}
	}

	log.Println("✓ Database schema initialized successfully!")
	return nil
}
