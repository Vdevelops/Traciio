-- Migration: Normalize visit and task statuses to pending/completed
-- Date: 2026-06-02

-- Normalize historical visit statuses
UPDATE visit_reports
SET status = CASE
    WHEN status IN ('approved', 'rejected', 'cancelled', 'completed') THEN 'completed'
    ELSE 'pending'
END
WHERE status IS NOT NULL
  AND status NOT IN ('pending', 'completed');

ALTER TABLE visit_reports
ALTER COLUMN status SET DEFAULT 'pending';

ALTER TABLE visit_reports DROP CONSTRAINT IF EXISTS check_visit_reports_status;
ALTER TABLE visit_reports DROP CONSTRAINT IF EXISTS visit_reports_status_check;
ALTER TABLE visit_reports
ADD CONSTRAINT check_visit_reports_status
CHECK (status IN ('pending', 'completed'));

-- Normalize historical task statuses
UPDATE tasks
SET status = CASE
    WHEN status IN ('completed', 'approved', 'cancelled', 'rejected') THEN 'completed'
    ELSE 'pending'
END
WHERE status IS NOT NULL
  AND status NOT IN ('pending', 'completed');

ALTER TABLE tasks
ALTER COLUMN status SET DEFAULT 'pending';

ALTER TABLE tasks DROP CONSTRAINT IF EXISTS check_tasks_status;
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS tasks_status_check;
ALTER TABLE tasks
ADD CONSTRAINT check_tasks_status
CHECK (status IN ('pending', 'completed'));
