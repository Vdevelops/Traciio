-- Migration: Add brick_id to deals table
-- Description: Links deals (pipeline) to bricks for territorial tracking and analytics

ALTER TABLE deals
ADD COLUMN IF NOT EXISTS brick_id UUID,
ADD CONSTRAINT fk_deals_brick FOREIGN KEY (brick_id)
    REFERENCES bricks(id) ON DELETE SET NULL;

-- Create index for performance
CREATE INDEX IF NOT EXISTS idx_deals_brick_id ON deals(brick_id) WHERE deleted_at IS NULL;

-- Add composite index for common queries (brick_id + status)
CREATE INDEX IF NOT EXISTS idx_deals_brick_status ON deals(brick_id, status) WHERE deleted_at IS NULL;

-- Add comment
COMMENT ON COLUMN deals.brick_id IS 'Brick/Area assignment for territorial tracking and analytics';

