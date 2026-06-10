-- Add map coordinates for lead location capture.

ALTER TABLE leads ADD COLUMN IF NOT EXISTS latitude DECIMAL(10,8);
ALTER TABLE leads ADD COLUMN IF NOT EXISTS longitude DECIMAL(11,8);

CREATE INDEX IF NOT EXISTS idx_leads_location
    ON leads(latitude, longitude)
    WHERE deleted_at IS NULL
      AND latitude IS NOT NULL
      AND longitude IS NOT NULL;

COMMENT ON COLUMN leads.latitude IS 'Latitude coordinate captured for lead location.';
COMMENT ON COLUMN leads.longitude IS 'Longitude coordinate captured for lead location.';
