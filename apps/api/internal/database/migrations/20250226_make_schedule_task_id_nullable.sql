-- Make schedule task_id nullable to allow standalone schedules
ALTER TABLE schedules ALTER COLUMN task_id DROP NOT NULL;
COMMENT ON COLUMN schedules.task_id IS 'Reference to the task this schedule is connected to (nullable for standalone schedules)';

