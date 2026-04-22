-- Migration: Add Lead Status Relation to Leads Table
-- Date: 2025-12-22
-- Description: Adds foreign key relationship between leads and lead_statuses tables

-- =====================================================
-- 1. ADD LEAD_STATUS_ID COLUMN TO LEADS TABLE
-- =====================================================
ALTER TABLE leads ADD COLUMN IF NOT EXISTS lead_status_id UUID;

-- =====================================================
-- 2. CREATE FOREIGN KEY CONSTRAINT
-- =====================================================
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'fk_leads_lead_status_id'
    ) THEN
        ALTER TABLE leads
        ADD CONSTRAINT fk_leads_lead_status_id
        FOREIGN KEY (lead_status_id)
        REFERENCES lead_statuses(id)
        ON DELETE SET NULL
        ON UPDATE CASCADE;
    END IF;
END $$;

-- =====================================================
-- 3. CREATE INDEX FOR PERFORMANCE
-- =====================================================
CREATE INDEX IF NOT EXISTS idx_leads_lead_status_id 
    ON leads(lead_status_id) WHERE deleted_at IS NULL;

-- =====================================================
-- 4. MIGRATE EXISTING DATA
-- =====================================================
-- Map existing lead_status string values to lead_status_id based on code
-- This assumes lead_status column contains codes like 'NEW', 'CONTACTED', etc.
UPDATE leads l
SET lead_status_id = ls.id
FROM lead_statuses ls
WHERE l.lead_status = ls.code
  AND l.lead_status_id IS NULL
  AND ls.deleted_at IS NULL;

-- For leads with status that doesn't match any code, set to default status
UPDATE leads l
SET lead_status_id = (
    SELECT id FROM lead_statuses 
    WHERE is_default = true AND deleted_at IS NULL 
    LIMIT 1
)
WHERE l.lead_status_id IS NULL;

-- =====================================================
-- 5. COMMENTS FOR DOCUMENTATION
-- =====================================================
COMMENT ON COLUMN leads.lead_status_id IS 'Foreign key reference to lead_statuses table. Replaces string-based lead_status for better data integrity.';
COMMENT ON COLUMN leads.lead_status IS 'Legacy string field. Will be deprecated in favor of lead_status_id.';






