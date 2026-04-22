-- Migration: Create Bricks table
-- Date: 2025-02-01
-- Description: Adds brick/area management for territory-based sales target management

-- =====================================================
-- 1. CREATE BRICKS TABLE
-- =====================================================
CREATE TABLE IF NOT EXISTS bricks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Identitas Brick
    name VARCHAR(255) NOT NULL,                    -- Nama brick (e.g., "Jakarta Pusat - Menteng")
    code VARCHAR(50) UNIQUE NOT NULL,              -- Kode brick (e.g., "JKT-PST-MTG")
    description TEXT,                              -- Deskripsi brick

    -- Lokasi & Wilayah
    province VARCHAR(100) NOT NULL,                -- Provinsi (e.g., "DKI Jakarta")
    regency VARCHAR(100) NOT NULL,                 -- Kabupaten/Kota (e.g., "Jakarta Pusat")
    district VARCHAR(100),                         -- Kecamatan (opsional)

    -- Manager Assignment
    manager_id UUID,                               -- Manager (Sales Manager) yang mengelola brick ini
    CONSTRAINT fk_bricks_manager FOREIGN KEY (manager_id)
        REFERENCES users(id) ON DELETE SET NULL,

    -- Status & Metadata
    status VARCHAR(20) NOT NULL DEFAULT 'active',  -- active, inactive
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

    -- Constraints
    CONSTRAINT chk_bricks_status CHECK (status IN ('active', 'inactive'))
);

-- =====================================================
-- 2. CREATE INDEXES FOR BRICKS
-- =====================================================
CREATE INDEX IF NOT EXISTS idx_bricks_code ON bricks(code) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_bricks_manager_id ON bricks(manager_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_bricks_regency ON bricks(regency) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_bricks_province ON bricks(province) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_bricks_status ON bricks(status) WHERE deleted_at IS NULL;

-- =====================================================
-- 3. UPDATE TRIGGER FOR UPDATED_AT
-- =====================================================
-- Create or replace trigger function
CREATE OR REPLACE FUNCTION update_bricks_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger
DROP TRIGGER IF EXISTS trigger_bricks_updated_at ON bricks;
CREATE TRIGGER trigger_bricks_updated_at
    BEFORE UPDATE ON bricks
    FOR EACH ROW
    EXECUTE FUNCTION update_bricks_updated_at();

-- =====================================================
-- 4. COMMENTS FOR DOCUMENTATION
-- =====================================================
COMMENT ON TABLE bricks IS 'Brick/Area management - fixed territory units for sales management. Historical data remains attached to brick even if sales staff changes.';
COMMENT ON COLUMN bricks.code IS 'Unique code identifier for the brick (case-insensitive unique)';
COMMENT ON COLUMN bricks.manager_id IS 'Sales Manager who manages this brick (optional, can be assigned later)';
COMMENT ON COLUMN bricks.regency IS 'Kabupaten/Kota name for geographic identification';
COMMENT ON COLUMN bricks.province IS 'Province name for geographic identification';

