-- Migration: Create Lead Statuses table with scoring system
-- Date: 2025-01-28
-- Description: Adds dynamic lead status management with scoring (0-100%)

-- =====================================================
-- 1. CREATE LEAD_STATUSES TABLE
-- =====================================================
CREATE TABLE IF NOT EXISTS lead_statuses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Status Information
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL,
    description TEXT,
    
    -- Scoring (0-100%)
    score INTEGER NOT NULL DEFAULT 0 CHECK (score >= 0 AND score <= 100),
    
    -- Display Settings
    color VARCHAR(20) DEFAULT '#3B82F6',
    "order" INTEGER DEFAULT 0,
    
    -- Status Flags
    is_active BOOLEAN DEFAULT TRUE,
    is_default BOOLEAN DEFAULT FALSE,
    is_converted BOOLEAN DEFAULT FALSE,
    
    -- Audit Fields
    created_by VARCHAR(255),
    updated_by VARCHAR(255),
    deleted_by VARCHAR(255),
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- =====================================================
-- 2. CREATE INDEXES
-- =====================================================
CREATE UNIQUE INDEX IF NOT EXISTS idx_lead_statuses_name 
    ON lead_statuses(name) WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_lead_statuses_code 
    ON lead_statuses(code) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_lead_statuses_order 
    ON lead_statuses("order") WHERE deleted_at IS NULL AND is_active = TRUE;

CREATE INDEX IF NOT EXISTS idx_lead_statuses_is_default 
    ON lead_statuses(is_default) WHERE deleted_at IS NULL AND is_default = TRUE;

CREATE INDEX IF NOT EXISTS idx_lead_statuses_is_converted 
    ON lead_statuses(is_converted) WHERE deleted_at IS NULL AND is_converted = TRUE;

CREATE INDEX IF NOT EXISTS idx_lead_statuses_is_active 
    ON lead_statuses(is_active) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_lead_statuses_score 
    ON lead_statuses(score) WHERE deleted_at IS NULL;

-- =====================================================
-- 3. ADD LEAD_SCORE TO LEADS TABLE
-- =====================================================
-- Add lead_score field to track current score based on status
ALTER TABLE leads ADD COLUMN IF NOT EXISTS lead_score INTEGER DEFAULT 0 CHECK (lead_score >= 0 AND lead_score <= 100);

-- Create index for score-based queries and sorting
CREATE INDEX IF NOT EXISTS idx_leads_lead_score 
    ON leads(lead_score) WHERE deleted_at IS NULL;

-- =====================================================
-- 4. UPDATE TRIGGER FOR UPDATED_AT
-- =====================================================
-- Create or replace trigger function
CREATE OR REPLACE FUNCTION update_lead_statuses_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger
DROP TRIGGER IF EXISTS trigger_lead_statuses_updated_at ON lead_statuses;
CREATE TRIGGER trigger_lead_statuses_updated_at
    BEFORE UPDATE ON lead_statuses
    FOR EACH ROW
    EXECUTE FUNCTION update_lead_statuses_updated_at();

-- =====================================================
-- 5. COMMENTS FOR DOCUMENTATION
-- =====================================================
COMMENT ON TABLE lead_statuses IS 'Dynamic lead status management with scoring system';
COMMENT ON COLUMN lead_statuses.score IS 'Score percentage (0-100) that automatically updates lead.lead_score when status changes';
COMMENT ON COLUMN lead_statuses.is_default IS 'Marks the default status for new leads (only one should be true)';
COMMENT ON COLUMN lead_statuses.is_converted IS 'Marks the status used when lead is converted to opportunity';
COMMENT ON COLUMN lead_statuses."order" IS 'Display order in UI (ascending)';
COMMENT ON COLUMN lead_statuses.color IS 'Hex color code for UI display (#RRGGBB format)';
COMMENT ON COLUMN leads.lead_score IS 'Current lead score (0-100%) based on lead_status score';
