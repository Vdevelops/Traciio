-- Migration: Create brick_target_distributions table
-- Date: 2025-02-01
-- Description: Adds target distribution system for distributing brick targets to sales within the same brick

-- =====================================================
-- 1. CREATE BRICK_TARGET_DISTRIBUTIONS TABLE
-- =====================================================
CREATE TABLE IF NOT EXISTS brick_target_distributions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Reference
    brick_id UUID NOT NULL,                        -- Brick yang memiliki target
    brick_target_id UUID NOT NULL,                 -- Reference ke monthly_targets (brick_id)
    sales_user_id UUID NOT NULL,                   -- Sales yang menerima distribusi target

    -- Distribusi Target
    distributed_amount BIGINT NOT NULL DEFAULT 0,  -- Jumlah target yang didistribusikan (dalam sen)

    -- Metadata
    distributed_by UUID NOT NULL,                  -- Manager yang melakukan distribusi
    distributed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Foreign Keys
    CONSTRAINT fk_brick_target_dist_brick FOREIGN KEY (brick_id)
        REFERENCES bricks(id) ON DELETE CASCADE,
    CONSTRAINT fk_brick_target_dist_target FOREIGN KEY (brick_target_id)
        REFERENCES monthly_targets(id) ON DELETE CASCADE,
    CONSTRAINT fk_brick_target_dist_sales FOREIGN KEY (sales_user_id)
        REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_brick_target_dist_distributor FOREIGN KEY (distributed_by)
        REFERENCES users(id) ON DELETE RESTRICT,

    -- Constraints
    CONSTRAINT chk_distributed_amount CHECK (distributed_amount >= 0)
);

-- =====================================================
-- 2. CREATE INDEXES FOR BRICK_TARGET_DISTRIBUTIONS
-- =====================================================
CREATE INDEX IF NOT EXISTS idx_brick_target_dist_brick_id
    ON brick_target_distributions(brick_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_brick_target_dist_target_id
    ON brick_target_distributions(brick_target_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_brick_target_dist_sales_id
    ON brick_target_distributions(sales_user_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_brick_target_dist_period
    ON brick_target_distributions(brick_id, brick_target_id) WHERE deleted_at IS NULL;

-- Unique constraint: One distribution per sales per brick target
CREATE UNIQUE INDEX IF NOT EXISTS idx_brick_target_dist_unique
    ON brick_target_distributions(sales_user_id, brick_target_id)
    WHERE deleted_at IS NULL;

-- =====================================================
-- 3. UPDATE TRIGGER FOR UPDATED_AT
-- =====================================================
-- Create or replace trigger function
CREATE OR REPLACE FUNCTION update_brick_target_distributions_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger
DROP TRIGGER IF EXISTS trigger_brick_target_distributions_updated_at ON brick_target_distributions;
CREATE TRIGGER trigger_brick_target_distributions_updated_at
    BEFORE UPDATE ON brick_target_distributions
    FOR EACH ROW
    EXECUTE FUNCTION update_brick_target_distributions_updated_at();

-- =====================================================
-- 4. COMMENTS FOR DOCUMENTATION
-- =====================================================
COMMENT ON TABLE brick_target_distributions IS 'Target distribution from brick to sales. Manager can distribute brick monthly target to multiple sales within the same brick.';
COMMENT ON COLUMN brick_target_distributions.brick_target_id IS 'Reference to monthly_targets table where brick_id is set';
COMMENT ON COLUMN brick_target_distributions.distributed_amount IS 'Amount distributed to sales in smallest currency unit (sen). Total distributed_amount does not need to equal brick target_amount (partial distribution allowed).';
COMMENT ON COLUMN brick_target_distributions.distributed_by IS 'Manager who performed the distribution (must be manager of the brick)';

