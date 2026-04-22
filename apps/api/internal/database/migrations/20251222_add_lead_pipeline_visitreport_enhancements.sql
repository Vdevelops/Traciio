-- Migration: Add Lead, Pipeline, and Visit Report enhancements
-- Date: 2025-12-22
-- Description: Adds BANT qualification fields, deal history tracking, and visit report outcomes

-- =====================================================
-- 1. LEADS TABLE ENHANCEMENTS
-- =====================================================
-- Add probability and estimated value fields
ALTER TABLE leads ADD COLUMN IF NOT EXISTS probability INTEGER DEFAULT 10;
ALTER TABLE leads ADD COLUMN IF NOT EXISTS estimated_value BIGINT DEFAULT 0;

-- Add BANT qualification fields
ALTER TABLE leads ADD COLUMN IF NOT EXISTS budget_confirmed BOOLEAN DEFAULT FALSE;
ALTER TABLE leads ADD COLUMN IF NOT EXISTS budget_amount BIGINT;
ALTER TABLE leads ADD COLUMN IF NOT EXISTS authority_confirmed BOOLEAN DEFAULT FALSE;
ALTER TABLE leads ADD COLUMN IF NOT EXISTS authority_person VARCHAR(255);
ALTER TABLE leads ADD COLUMN IF NOT EXISTS need_confirmed BOOLEAN DEFAULT FALSE;
ALTER TABLE leads ADD COLUMN IF NOT EXISTS need_description TEXT;
ALTER TABLE leads ADD COLUMN IF NOT EXISTS timeline_confirmed BOOLEAN DEFAULT FALSE;
ALTER TABLE leads ADD COLUMN IF NOT EXISTS expected_close_date DATE;

-- Add indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_leads_probability ON leads(probability) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_leads_estimated_value ON leads(estimated_value) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_leads_expected_close_date ON leads(expected_close_date) WHERE deleted_at IS NULL;

-- =====================================================
-- 2. DEALS TABLE ENHANCEMENTS
-- =====================================================
-- Add BANT fields carried from Lead
ALTER TABLE deals ADD COLUMN IF NOT EXISTS budget_confirmed BOOLEAN DEFAULT FALSE;
ALTER TABLE deals ADD COLUMN IF NOT EXISTS authority_confirmed BOOLEAN DEFAULT FALSE;
ALTER TABLE deals ADD COLUMN IF NOT EXISTS need_confirmed BOOLEAN DEFAULT FALSE;
ALTER TABLE deals ADD COLUMN IF NOT EXISTS timeline_confirmed BOOLEAN DEFAULT FALSE;

-- Add close reason field
ALTER TABLE deals ADD COLUMN IF NOT EXISTS close_reason TEXT;

-- Add indexes
CREATE INDEX IF NOT EXISTS idx_deals_budget_confirmed ON deals(budget_confirmed) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_deals_authority_confirmed ON deals(authority_confirmed) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_deals_need_confirmed ON deals(need_confirmed) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_deals_timeline_confirmed ON deals(timeline_confirmed) WHERE deleted_at IS NULL;

-- =====================================================
-- 3. VISIT REPORTS TABLE ENHANCEMENTS
-- =====================================================
-- Add outcome and next steps fields
ALTER TABLE visit_reports ADD COLUMN IF NOT EXISTS outcome VARCHAR(50);
ALTER TABLE visit_reports ADD COLUMN IF NOT EXISTS next_steps TEXT;

-- Add index for outcome filtering
CREATE INDEX IF NOT EXISTS idx_visit_reports_outcome ON visit_reports(outcome) WHERE deleted_at IS NULL;

-- Add check constraint for outcome values
ALTER TABLE visit_reports DROP CONSTRAINT IF EXISTS check_visit_report_outcome;
ALTER TABLE visit_reports ADD CONSTRAINT check_visit_report_outcome 
    CHECK (outcome IS NULL OR outcome IN ('positive', 'neutral', 'negative', 'very_positive'));

-- =====================================================
-- 4. DEAL HISTORIES TABLE (NEW)
-- =====================================================
-- Create deal_histories table for audit trail
CREATE TABLE IF NOT EXISTS deal_histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    deal_id UUID NOT NULL REFERENCES deals(id) ON DELETE CASCADE,
    from_stage_id UUID REFERENCES pipeline_stages(id),
    from_stage_name VARCHAR(255),
    to_stage_id UUID NOT NULL REFERENCES pipeline_stages(id),
    to_stage_name VARCHAR(255) NOT NULL,
    from_probability INTEGER DEFAULT 0,
    to_probability INTEGER DEFAULT 0,
    days_in_prev_stage INTEGER,
    changed_by UUID NOT NULL REFERENCES users(id),
    changed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reason TEXT,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Add indexes for deal_histories
CREATE INDEX IF NOT EXISTS idx_deal_histories_deal_id ON deal_histories(deal_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_deal_histories_changed_by ON deal_histories(changed_by) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_deal_histories_changed_at ON deal_histories(changed_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_deal_histories_from_stage_id ON deal_histories(from_stage_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_deal_histories_to_stage_id ON deal_histories(to_stage_id) WHERE deleted_at IS NULL;

-- Add composite index for common queries (deal timeline)
CREATE INDEX IF NOT EXISTS idx_deal_histories_deal_changed_at ON deal_histories(deal_id, changed_at DESC) WHERE deleted_at IS NULL;

-- =====================================================
-- 5. COMMENTS & DOCUMENTATION
-- =====================================================
COMMENT ON COLUMN leads.probability IS 'Win probability percentage (0-100), default 10';
COMMENT ON COLUMN leads.estimated_value IS 'Estimated deal value in smallest currency unit (sen)';
COMMENT ON COLUMN leads.budget_confirmed IS 'BANT: Whether budget is confirmed';
COMMENT ON COLUMN leads.budget_amount IS 'BANT: Confirmed budget amount';
COMMENT ON COLUMN leads.authority_confirmed IS 'BANT: Whether decision maker authority is confirmed';
COMMENT ON COLUMN leads.authority_person IS 'BANT: Name of decision maker';
COMMENT ON COLUMN leads.need_confirmed IS 'BANT: Whether customer need is confirmed';
COMMENT ON COLUMN leads.need_description IS 'BANT: Description of customer need';
COMMENT ON COLUMN leads.timeline_confirmed IS 'BANT: Whether buying timeline is confirmed';
COMMENT ON COLUMN leads.expected_close_date IS 'Expected close date for the deal';

COMMENT ON COLUMN deals.budget_confirmed IS 'BANT: Carried from lead conversion';
COMMENT ON COLUMN deals.authority_confirmed IS 'BANT: Carried from lead conversion';
COMMENT ON COLUMN deals.need_confirmed IS 'BANT: Carried from lead conversion';
COMMENT ON COLUMN deals.timeline_confirmed IS 'BANT: Carried from lead conversion';
COMMENT ON COLUMN deals.close_reason IS 'Reason for won/lost status';

COMMENT ON COLUMN visit_reports.outcome IS 'Visit outcome: positive, neutral, negative, very_positive';
COMMENT ON COLUMN visit_reports.next_steps IS 'Next action items after visit';

COMMENT ON TABLE deal_histories IS 'Audit trail for deal stage transitions and changes';
COMMENT ON COLUMN deal_histories.from_stage_id IS 'Previous stage ID (NULL for initial creation)';
COMMENT ON COLUMN deal_histories.to_stage_id IS 'New stage ID';
COMMENT ON COLUMN deal_histories.days_in_prev_stage IS 'Calculated duration in previous stage';
COMMENT ON COLUMN deal_histories.changed_by IS 'User who made the change';
COMMENT ON COLUMN deal_histories.changed_at IS 'Timestamp of the change';
COMMENT ON COLUMN deal_histories.reason IS 'Reason for stage change';
COMMENT ON COLUMN deal_histories.notes IS 'Additional notes about the change';

-- =====================================================
-- 6. DATA VALIDATION
-- =====================================================
-- Ensure existing leads have default probability
UPDATE leads SET probability = 10 WHERE probability IS NULL;

-- Ensure existing leads have default estimated_value
UPDATE leads SET estimated_value = 0 WHERE estimated_value IS NULL;

-- Ensure existing deals have default BANT values
UPDATE deals SET budget_confirmed = FALSE WHERE budget_confirmed IS NULL;
UPDATE deals SET authority_confirmed = FALSE WHERE authority_confirmed IS NULL;
UPDATE deals SET need_confirmed = FALSE WHERE need_confirmed IS NULL;
UPDATE deals SET timeline_confirmed = FALSE WHERE timeline_confirmed IS NULL;
