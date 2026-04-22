-- Migration: Add brick_id to accounts table
-- Description: Links accounts (apotek/rumah sakit) to bricks for territorial tracking

ALTER TABLE accounts
ADD COLUMN IF NOT EXISTS brick_id UUID,
ADD CONSTRAINT fk_accounts_brick FOREIGN KEY (brick_id)
    REFERENCES bricks(id) ON DELETE SET NULL;

-- Create index for performance
CREATE INDEX IF NOT EXISTS idx_accounts_brick_id ON accounts(brick_id) WHERE deleted_at IS NULL;

-- Add comment
COMMENT ON COLUMN accounts.brick_id IS 'Brick/Area assignment for territorial management';

