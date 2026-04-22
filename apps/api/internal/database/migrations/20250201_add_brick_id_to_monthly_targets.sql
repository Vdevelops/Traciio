-- Migration: Add brick_id to monthly_targets table
-- Date: 2025-02-01
-- Description: Adds brick support to monthly targets (group_id OR user_id OR brick_id, but only one)

-- =====================================================
-- 1. ADD BRICK_ID COLUMN TO MONTHLY_TARGETS
-- =====================================================
ALTER TABLE monthly_targets
ADD COLUMN IF NOT EXISTS brick_id UUID,
ADD CONSTRAINT fk_monthly_targets_brick FOREIGN KEY (brick_id)
    REFERENCES bricks(id) ON DELETE CASCADE;

-- =====================================================
-- 2. UPDATE CONSTRAINT: group_id OR user_id OR brick_id (but only one)
-- =====================================================
-- Drop existing constraint if it exists (may have different names or structure)
DO $$
BEGIN
    -- Try to drop constraint if it exists (handle both group_id and division_id cases)
    IF EXISTS (
        SELECT 1 FROM pg_constraint 
        WHERE conname = 'chk_monthly_targets_scope'
    ) THEN
        ALTER TABLE monthly_targets DROP CONSTRAINT chk_monthly_targets_scope;
    END IF;
END $$;

-- Add new constraint that supports group_id, user_id, or brick_id (only one)
ALTER TABLE monthly_targets
ADD CONSTRAINT chk_monthly_targets_scope CHECK (
    (group_id IS NOT NULL AND user_id IS NULL AND brick_id IS NULL) OR
    (group_id IS NULL AND user_id IS NOT NULL AND brick_id IS NULL) OR
    (group_id IS NULL AND user_id IS NULL AND brick_id IS NOT NULL)
);

-- Note: If the table uses division_id instead of group_id, you may need to handle that separately
-- This migration assumes group_id exists. If division_id exists instead, update accordingly.

-- =====================================================
-- 3. CREATE INDEXES FOR BRICK TARGETS
-- =====================================================
CREATE INDEX IF NOT EXISTS idx_monthly_targets_brick
    ON monthly_targets(brick_id, year, month) WHERE deleted_at IS NULL;

-- Unique constraint: One target per brick per month/year
CREATE UNIQUE INDEX IF NOT EXISTS idx_monthly_targets_brick_unique
    ON monthly_targets(brick_id, year, month)
    WHERE deleted_at IS NULL AND brick_id IS NOT NULL;

-- =====================================================
-- 4. COMMENTS FOR DOCUMENTATION
-- =====================================================
COMMENT ON COLUMN monthly_targets.brick_id IS 'Foreign key to bricks table - for brick-level monthly targets. Only one of group_id, user_id, or brick_id can be set.';

