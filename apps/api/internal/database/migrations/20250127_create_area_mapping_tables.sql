-- Enable PostGIS extension
CREATE EXTENSION IF NOT EXISTS postgis;

-- Area captures (GPS points from visits)
CREATE TABLE IF NOT EXISTS area_captures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    visit_report_id UUID NOT NULL,
    capture_type VARCHAR(20) NOT NULL CHECK (capture_type IN ('check_in', 'check_out', 'area')),
    location GEOGRAPHY(POINT, 4326) NOT NULL,
    address TEXT,
    accuracy DECIMAL(10,2),
    captured_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_area_captures_visit_report FOREIGN KEY (visit_report_id) 
        REFERENCES visit_reports(id) ON DELETE CASCADE
);

-- Create indexes for area_captures
CREATE INDEX IF NOT EXISTS idx_area_captures_location ON area_captures USING GIST(location);
CREATE INDEX IF NOT EXISTS idx_area_captures_visit_report_id ON area_captures(visit_report_id);
CREATE INDEX IF NOT EXISTS idx_area_captures_captured_at ON area_captures(captured_at);
CREATE INDEX IF NOT EXISTS idx_area_captures_capture_type ON area_captures(capture_type);

-- Territories (polygon areas)
CREATE TABLE IF NOT EXISTS territories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    polygon GEOGRAPHY(POLYGON, 4326) NOT NULL,
    assigned_to UUID,
    color VARCHAR(50) DEFAULT '#3B82F6',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_territories_assigned_to FOREIGN KEY (assigned_to) 
        REFERENCES users(id) ON DELETE SET NULL
);

-- Create indexes for territories
CREATE INDEX IF NOT EXISTS idx_territories_polygon ON territories USING GIST(polygon);
CREATE INDEX IF NOT EXISTS idx_territories_assigned_to ON territories(assigned_to);
CREATE INDEX IF NOT EXISTS idx_territories_name ON territories(name);

-- Coverage analysis
CREATE TABLE IF NOT EXISTS coverage_analysis (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    territory_id UUID,
    user_id UUID,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    visit_count INTEGER NOT NULL DEFAULT 0,
    coverage_percent DECIMAL(5,2),
    analyzed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_coverage_analysis_territory FOREIGN KEY (territory_id) 
        REFERENCES territories(id) ON DELETE SET NULL,
    CONSTRAINT fk_coverage_analysis_user FOREIGN KEY (user_id) 
        REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chk_coverage_percent CHECK (coverage_percent >= 0 AND coverage_percent <= 100)
);

-- Create indexes for coverage_analysis
CREATE INDEX IF NOT EXISTS idx_coverage_analysis_territory_id ON coverage_analysis(territory_id);
CREATE INDEX IF NOT EXISTS idx_coverage_analysis_user_id ON coverage_analysis(user_id);
CREATE INDEX IF NOT EXISTS idx_coverage_analysis_period ON coverage_analysis(period_start, period_end);
CREATE INDEX IF NOT EXISTS idx_coverage_analysis_analyzed_at ON coverage_analysis(analyzed_at);

-- Add comment for PostGIS extension
COMMENT ON EXTENSION postgis IS 'PostGIS geometry and geography spatial types and functions';
