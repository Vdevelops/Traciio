-- Migration: Add brick_id to visit_reports table
-- Description: Links visit reports to bricks for territorial tracking and analytics

ALTER TABLE visit_reports
ADD COLUMN IF NOT EXISTS brick_id UUID,
ADD CONSTRAINT fk_visit_reports_brick FOREIGN KEY (brick_id)
    REFERENCES bricks(id) ON DELETE SET NULL;

-- Create index for performance
CREATE INDEX IF NOT EXISTS idx_visit_reports_brick_id ON visit_reports(brick_id) WHERE deleted_at IS NULL;

-- Add composite index for common queries (brick_id + visit_date)
CREATE INDEX IF NOT EXISTS idx_visit_reports_brick_date ON visit_reports(brick_id, visit_date) WHERE deleted_at IS NULL;

-- Add comment
COMMENT ON COLUMN visit_reports.brick_id IS 'Brick/Area assignment for territorial tracking and analytics';

