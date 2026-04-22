-- Migration: Create Divisions table and add division_id to users table
-- Date: 2025-01-29
-- Description: Adds division management with one-to-many relationship to users

-- =====================================================
-- 1. CREATE DIVISIONS TABLE
-- =====================================================
CREATE TABLE IF NOT EXISTS divisions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Division Information
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) NOT NULL,
    description TEXT,
    
    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    
    -- Timestamps
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- =====================================================
-- 2. CREATE INDEXES FOR DIVISIONS
-- =====================================================
CREATE UNIQUE INDEX IF NOT EXISTS idx_divisions_code 
    ON divisions(code) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_divisions_status 
    ON divisions(status) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_divisions_name 
    ON divisions(name) WHERE deleted_at IS NULL;

-- =====================================================
-- 3. ADD DIVISION_ID TO USERS TABLE
-- =====================================================
ALTER TABLE users ADD COLUMN IF NOT EXISTS division_id UUID;

-- Create foreign key constraint
ALTER TABLE users 
    ADD CONSTRAINT fk_users_division 
    FOREIGN KEY (division_id) 
    REFERENCES divisions(id) 
    ON DELETE SET NULL;

-- Create index for division_id
CREATE INDEX IF NOT EXISTS idx_users_division_id 
    ON users(division_id) WHERE deleted_at IS NULL;

-- =====================================================
-- 4. UPDATE TRIGGER FOR UPDATED_AT
-- =====================================================
-- Create or replace trigger function
CREATE OR REPLACE FUNCTION update_divisions_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger
DROP TRIGGER IF EXISTS trigger_divisions_updated_at ON divisions;
CREATE TRIGGER trigger_divisions_updated_at
    BEFORE UPDATE ON divisions
    FOR EACH ROW
    EXECUTE FUNCTION update_divisions_updated_at();

-- =====================================================
-- 5. COMMENTS FOR DOCUMENTATION
-- =====================================================
COMMENT ON TABLE divisions IS 'Division management - one division can have many users';
COMMENT ON COLUMN divisions.code IS 'Unique code identifier for the division';
COMMENT ON COLUMN users.division_id IS 'Foreign key to divisions table - one user belongs to one division';





