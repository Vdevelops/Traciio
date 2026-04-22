-- Migration: Optimize Task and Schedule Indexes for High-Volume Data
-- Date: 2025-02-01
-- Description: Adds composite indexes and trigram indexes for optimal performance with hundreds of thousands to millions of records

-- =====================================================
-- 1. TASKS TABLE - ADDITIONAL COMPOSITE INDEXES
-- =====================================================

-- Composite index for common query: assigned_to + status + due_date (for user task lists)
CREATE INDEX IF NOT EXISTS idx_tasks_assigned_status_due_date 
    ON tasks(assigned_to, status, due_date) 
    WHERE deleted_at IS NULL AND assigned_to IS NOT NULL AND due_date IS NOT NULL;

-- Composite index for filtering by status and priority
CREATE INDEX IF NOT EXISTS idx_tasks_status_priority 
    ON tasks(status, priority) 
    WHERE deleted_at IS NULL;

-- Composite index for date range queries with status
CREATE INDEX IF NOT EXISTS idx_tasks_due_date_status 
    ON tasks(due_date, status) 
    WHERE deleted_at IS NULL AND due_date IS NOT NULL;

-- Composite index for account-related tasks
CREATE INDEX IF NOT EXISTS idx_tasks_account_status 
    ON tasks(account_id, status) 
    WHERE deleted_at IS NULL AND account_id IS NOT NULL;

-- Composite index for deal-related tasks
CREATE INDEX IF NOT EXISTS idx_tasks_deal_status 
    ON tasks(deal_id, status) 
    WHERE deleted_at IS NULL AND deal_id IS NOT NULL;

-- Composite index for created_at ordering (for pagination)
CREATE INDEX IF NOT EXISTS idx_tasks_created_at_desc 
    ON tasks(created_at DESC) 
    WHERE deleted_at IS NULL;

-- Composite index for type and status filtering
CREATE INDEX IF NOT EXISTS idx_tasks_type_status 
    ON tasks(type, status) 
    WHERE deleted_at IS NULL;

-- =====================================================
-- 2. TASKS TABLE - TRIGRAM INDEXES FOR TEXT SEARCH
-- =====================================================

-- Enable pg_trgm extension if not already enabled
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Trigram index for title search (faster than LIKE for large datasets)
CREATE INDEX IF NOT EXISTS idx_tasks_title_trgm 
    ON tasks USING gin(title gin_trgm_ops) 
    WHERE deleted_at IS NULL;

-- Trigram index for description search
CREATE INDEX IF NOT EXISTS idx_tasks_description_trgm 
    ON tasks USING gin(description gin_trgm_ops) 
    WHERE deleted_at IS NULL AND description IS NOT NULL;

-- =====================================================
-- 3. SCHEDULES TABLE - COMPOSITE INDEXES
-- =====================================================

-- Composite index for user schedules with date range (most common query)
CREATE INDEX IF NOT EXISTS idx_schedules_user_scheduled_at 
    ON schedules(user_id, scheduled_at) 
    WHERE deleted_at IS NULL;

-- Composite index for user schedules with status
CREATE INDEX IF NOT EXISTS idx_schedules_user_status 
    ON schedules(user_id, status) 
    WHERE deleted_at IS NULL;

-- Composite index for scheduled_at with status (for calendar views)
CREATE INDEX IF NOT EXISTS idx_schedules_scheduled_at_status 
    ON schedules(scheduled_at, status) 
    WHERE deleted_at IS NULL;

-- Composite index for task_id queries (for task-schedule relationships)
CREATE INDEX IF NOT EXISTS idx_schedules_task_id_status 
    ON schedules(task_id, status) 
    WHERE deleted_at IS NULL AND task_id IS NOT NULL;

-- Composite index for date range queries with user
CREATE INDEX IF NOT EXISTS idx_schedules_user_date_range 
    ON schedules(user_id, scheduled_at) 
    WHERE deleted_at IS NULL;

-- Composite index for Google Calendar sync status queries
CREATE INDEX IF NOT EXISTS idx_schedules_sync_status_user 
    ON schedules(google_calendar_sync_status, user_id) 
    WHERE deleted_at IS NULL;

-- =====================================================
-- 4. SCHEDULES TABLE - TRIGRAM INDEXES FOR TEXT SEARCH
-- =====================================================

-- Trigram index for title search
CREATE INDEX IF NOT EXISTS idx_schedules_title_trgm 
    ON schedules USING gin(title gin_trgm_ops) 
    WHERE deleted_at IS NULL;

-- Trigram index for description search
CREATE INDEX IF NOT EXISTS idx_schedules_description_trgm 
    ON schedules USING gin(description gin_trgm_ops) 
    WHERE deleted_at IS NULL AND description IS NOT NULL;

-- =====================================================
-- 5. REMINDERS TABLE - ADDITIONAL INDEXES
-- =====================================================

-- Composite index for task reminders
CREATE INDEX IF NOT EXISTS idx_reminders_task_remind_at 
    ON reminders(task_id, remind_at) 
    WHERE deleted_at IS NULL;

-- Composite index for pending reminders by date
CREATE INDEX IF NOT EXISTS idx_reminders_pending_date 
    ON reminders(remind_at, is_sent) 
    WHERE deleted_at IS NULL AND is_sent = false;

-- =====================================================
-- NOTES
-- =====================================================
-- These indexes are optimized for:
-- 1. High-volume queries (hundreds of thousands to millions of records)
-- 2. Common query patterns in task and schedule services
-- 3. Text search performance (trigram indexes for LIKE queries)
-- 4. Composite queries (multiple WHERE conditions)
-- 5. Date range queries (for calendar views and filtering)
--
-- Performance improvements expected:
-- - Text search: 10-100x faster with trigram indexes
-- - Composite queries: 5-20x faster with composite indexes
-- - Date range queries: 3-10x faster with optimized indexes
--
-- Index maintenance:
-- - Trigram indexes use GIN (Generalized Inverted Index) which is larger but faster for search
-- - Composite indexes are ordered by query frequency (most selective first)
-- - Partial indexes (WHERE deleted_at IS NULL) reduce index size and improve performance


