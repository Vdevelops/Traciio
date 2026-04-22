-- Migration: Add Google Calendar fields to Tasks table
-- Date: 2025-02-02
-- Description: Adds Google Calendar integration fields to tasks for automatic calendar sync

-- =====================================================
-- 1. ADD GOOGLE CALENDAR FIELDS TO TASKS TABLE
-- =====================================================
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS google_calendar_event_id VARCHAR(255);
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS google_calendar_sync_status VARCHAR(20) NOT NULL DEFAULT 'not_synced' CHECK (google_calendar_sync_status IN ('not_synced', 'synced', 'sync_failed'));
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS google_calendar_synced_at TIMESTAMP;

-- =====================================================
-- 2. CREATE INDEXES FOR PERFORMANCE
-- =====================================================
CREATE INDEX IF NOT EXISTS idx_tasks_google_calendar_event_id 
    ON tasks(google_calendar_event_id) WHERE deleted_at IS NULL AND google_calendar_event_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_tasks_google_calendar_sync_status 
    ON tasks(google_calendar_sync_status) WHERE deleted_at IS NULL;

-- =====================================================
-- 3. COMMENTS FOR DOCUMENTATION
-- =====================================================
COMMENT ON COLUMN tasks.google_calendar_event_id IS 'Google Calendar event ID if task is synced to Google Calendar';
COMMENT ON COLUMN tasks.google_calendar_sync_status IS 'Status of Google Calendar sync: not_synced, synced, sync_failed';
COMMENT ON COLUMN tasks.google_calendar_synced_at IS 'Timestamp when task was last synced to Google Calendar';

