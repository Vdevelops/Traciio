-- Migration: Create Monthly Targets table
-- Date: 2025-12-23
-- Description: Adds monthly target management for divisions and users (user targets override division targets)

-- =====================================================
-- 1. CREATE MONTHLY_TARGETS TABLE
-- =====================================================
CREATE TABLE IF NOT EXISTS monthly_targets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Target Scope (either division_id OR user_id, not both)
    division_id UUID REFERENCES divisions(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    
    -- Target Period
    year INTEGER NOT NULL CHECK (year >= 2000 AND year <= 2100),
    month INTEGER NOT NULL CHECK (month >= 1 AND month <= 12),
    
    -- Target Amount (in smallest currency unit - sen)
    target_amount BIGINT NOT NULL DEFAULT 0 CHECK (target_amount >= 0),
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    -- Constraint: Either division_id OR user_id must be set, but not both
    CONSTRAINT chk_monthly_targets_scope CHECK (
        (division_id IS NOT NULL AND user_id IS NULL) OR
        (division_id IS NULL AND user_id IS NOT NULL)
    )
);

-- =====================================================
-- 2. CREATE INDEXES FOR MONTHLY_TARGETS
-- =====================================================
-- Index for division targets
CREATE INDEX IF NOT EXISTS idx_monthly_targets_division 
    ON monthly_targets(division_id, year, month) WHERE deleted_at IS NULL;

-- Index for user targets
CREATE INDEX IF NOT EXISTS idx_monthly_targets_user 
    ON monthly_targets(user_id, year, month) WHERE deleted_at IS NULL;

-- Unique constraint: One target per division per month/year
CREATE UNIQUE INDEX IF NOT EXISTS idx_monthly_targets_division_unique 
    ON monthly_targets(division_id, year, month) 
    WHERE deleted_at IS NULL AND division_id IS NOT NULL;

-- Unique constraint: One target per user per month/year
CREATE UNIQUE INDEX IF NOT EXISTS idx_monthly_targets_user_unique 
    ON monthly_targets(user_id, year, month) 
    WHERE deleted_at IS NULL AND user_id IS NOT NULL;

-- Index for querying by year/month
CREATE INDEX IF NOT EXISTS idx_monthly_targets_period 
    ON monthly_targets(year, month) WHERE deleted_at IS NULL;

-- =====================================================
-- 3. UPDATE TRIGGER FOR UPDATED_AT
-- =====================================================
-- Create or replace trigger function
CREATE OR REPLACE FUNCTION update_monthly_targets_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger
DROP TRIGGER IF EXISTS trigger_monthly_targets_updated_at ON monthly_targets;
CREATE TRIGGER trigger_monthly_targets_updated_at
    BEFORE UPDATE ON monthly_targets
    FOR EACH ROW
    EXECUTE FUNCTION update_monthly_targets_updated_at();

-- =====================================================
-- 4. COMMENTS FOR DOCUMENTATION
-- =====================================================
COMMENT ON TABLE monthly_targets IS 'Monthly targets for divisions (default) and users (override). User targets take precedence over division targets.';
COMMENT ON COLUMN monthly_targets.division_id IS 'Division ID for division-level targets (default for all users in division)';
COMMENT ON COLUMN monthly_targets.user_id IS 'User ID for user-level targets (overrides division target)';
COMMENT ON COLUMN monthly_targets.year IS 'Target year (2000-2100)';
COMMENT ON COLUMN monthly_targets.month IS 'Target month (1-12)';
COMMENT ON COLUMN monthly_targets.target_amount IS 'Target amount in smallest currency unit (sen)';

