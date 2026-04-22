-- Add latitude and longitude columns to accounts table for geocoding
ALTER TABLE accounts
ADD COLUMN IF NOT EXISTS latitude DECIMAL(10, 8),
ADD COLUMN IF NOT EXISTS longitude DECIMAL(11, 8);

-- Add index for spatial queries (optional but recommended for map queries)
CREATE INDEX IF NOT EXISTS idx_accounts_location ON accounts(latitude, longitude) WHERE latitude IS NOT NULL AND longitude IS NOT NULL;

