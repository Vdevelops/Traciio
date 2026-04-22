-- Create schedules table
-- This table stores user-specific schedules connected to tasks
-- Each schedule is linked to a task and belongs to the user assigned to that task

CREATE TABLE IF NOT EXISTS schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    scheduled_at TIMESTAMP NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'completed', 'cancelled')),
    google_calendar_event_id VARCHAR(255),
    google_calendar_sync_status VARCHAR(20) NOT NULL DEFAULT 'not_synced' CHECK (google_calendar_sync_status IN ('not_synced', 'synced', 'sync_failed')),
    google_calendar_synced_at TIMESTAMP,
    reminder_minutes_before INTEGER CHECK (reminder_minutes_before >= 0 AND reminder_minutes_before <= 10080),
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_schedules_task_id ON schedules(task_id);
CREATE INDEX IF NOT EXISTS idx_schedules_user_id ON schedules(user_id);
CREATE INDEX IF NOT EXISTS idx_schedules_scheduled_at ON schedules(scheduled_at);
CREATE INDEX IF NOT EXISTS idx_schedules_status ON schedules(status);
CREATE INDEX IF NOT EXISTS idx_schedules_google_calendar_event_id ON schedules(google_calendar_event_id);
CREATE INDEX IF NOT EXISTS idx_schedules_google_calendar_sync_status ON schedules(google_calendar_sync_status);
CREATE INDEX IF NOT EXISTS idx_schedules_created_by ON schedules(created_by);
CREATE INDEX IF NOT EXISTS idx_schedules_deleted_at ON schedules(deleted_at);

-- Add comment
COMMENT ON TABLE schedules IS 'User-specific schedules connected to tasks for Google Calendar integration';
COMMENT ON COLUMN schedules.task_id IS 'Reference to the task this schedule is connected to';
COMMENT ON COLUMN schedules.user_id IS 'User who owns this schedule (from task.assigned_to)';
COMMENT ON COLUMN schedules.scheduled_at IS 'When the schedule/reminder should occur';
COMMENT ON COLUMN schedules.google_calendar_event_id IS 'Google Calendar event ID if synced';
COMMENT ON COLUMN schedules.google_calendar_sync_status IS 'Status of Google Calendar sync';
COMMENT ON COLUMN schedules.reminder_minutes_before IS 'Minutes before task due_date to remind (0-10080, max 7 days)';

