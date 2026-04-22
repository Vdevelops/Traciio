-- Migration: Add brick_id to users table
-- Date: 2025-02-01
-- Description: Adds brick assignment for sales users (sales can be assigned to a brick)

-- =====================================================
-- 1. ADD BRICK_ID TO USERS TABLE
-- =====================================================
ALTER TABLE users ADD COLUMN IF NOT EXISTS brick_id UUID;

-- Create foreign key constraint
ALTER TABLE users 
    ADD CONSTRAINT fk_users_brick 
    FOREIGN KEY (brick_id) 
    REFERENCES bricks(id) 
    ON DELETE SET NULL;

-- Create index for brick_id
CREATE INDEX IF NOT EXISTS idx_users_brick_id 
    ON users(brick_id) WHERE deleted_at IS NULL;

-- =====================================================
-- 2. COMMENTS FOR DOCUMENTATION
-- =====================================================
COMMENT ON COLUMN users.brick_id IS 'Foreign key to bricks table - sales user can be assigned to one brick. manager_id in bricks table is for brick manager, brick_id in users is for sales assignment.';

