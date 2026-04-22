-- Migration: Add Performance Indexes for Foreign Keys and Search Columns
-- Date: 2025-01-31
-- Description: Adds indexes for foreign keys and frequently searched columns to improve query performance
-- This migration addresses enterprise-scale performance requirements

-- =====================================================
-- 1. VISIT REPORTS TABLE INDEXES
-- =====================================================
-- Foreign key indexes
CREATE INDEX IF NOT EXISTS idx_visit_reports_account_id 
    ON visit_reports(account_id) WHERE deleted_at IS NULL AND account_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_visit_reports_sales_rep_id 
    ON visit_reports(sales_rep_id) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_visit_reports_deal_id 
    ON visit_reports(deal_id) WHERE deleted_at IS NULL AND deal_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_visit_reports_lead_id 
    ON visit_reports(lead_id) WHERE deleted_at IS NULL AND lead_id IS NOT NULL;

-- Date range queries (frequently used in dashboard/reports)
CREATE INDEX IF NOT EXISTS idx_visit_reports_visit_date 
    ON visit_reports(visit_date) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_visit_reports_created_at 
    ON visit_reports(created_at) WHERE deleted_at IS NULL;

-- Status filtering (frequently used)
CREATE INDEX IF NOT EXISTS idx_visit_reports_status 
    ON visit_reports(status) WHERE deleted_at IS NULL;

-- Composite index for common query patterns
CREATE INDEX IF NOT EXISTS idx_visit_reports_sales_rep_date 
    ON visit_reports(sales_rep_id, visit_date) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_visit_reports_status_date 
    ON visit_reports(status, visit_date) WHERE deleted_at IS NULL;

-- =====================================================
-- 2. DEALS TABLE INDEXES
-- =====================================================
-- Foreign key indexes
CREATE INDEX IF NOT EXISTS idx_deals_account_id 
    ON deals(account_id) WHERE deleted_at IS NULL AND account_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_deals_contact_id 
    ON deals(contact_id) WHERE deleted_at IS NULL AND contact_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_deals_stage_id 
    ON deals(stage_id) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_deals_assigned_to 
    ON deals(assigned_to) WHERE deleted_at IS NULL;

-- Status filtering (frequently used in aggregations)
CREATE INDEX IF NOT EXISTS idx_deals_status 
    ON deals(status) WHERE deleted_at IS NULL;

-- Date range queries
CREATE INDEX IF NOT EXISTS idx_deals_created_at 
    ON deals(created_at) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_deals_actual_close_date 
    ON deals(actual_close_date) WHERE deleted_at IS NULL AND actual_close_date IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_deals_expected_close_date 
    ON deals(expected_close_date) WHERE deleted_at IS NULL AND expected_close_date IS NOT NULL;

-- Composite indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_deals_assigned_status 
    ON deals(assigned_to, status) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_deals_status_close_date 
    ON deals(status, actual_close_date) WHERE deleted_at IS NULL AND actual_close_date IS NOT NULL;

-- Value range queries (for filtering by deal value)
CREATE INDEX IF NOT EXISTS idx_deals_value 
    ON deals(value) WHERE deleted_at IS NULL;

-- =====================================================
-- 3. LEADS TABLE INDEXES
-- =====================================================
-- Foreign key indexes
CREATE INDEX IF NOT EXISTS idx_leads_account_id 
    ON leads(account_id) WHERE deleted_at IS NULL AND account_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_leads_contact_id 
    ON leads(contact_id) WHERE deleted_at IS NULL AND contact_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_leads_assigned_to 
    ON leads(assigned_to) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_leads_opportunity_id 
    ON leads(opportunity_id) WHERE deleted_at IS NULL AND opportunity_id IS NOT NULL;

-- Status filtering (frequently used)
CREATE INDEX IF NOT EXISTS idx_leads_lead_status 
    ON leads(lead_status) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_leads_lead_source 
    ON leads(lead_source) WHERE deleted_at IS NULL;

-- Date range queries
CREATE INDEX IF NOT EXISTS idx_leads_created_at 
    ON leads(created_at) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_leads_converted_at 
    ON leads(converted_at) WHERE deleted_at IS NULL AND converted_at IS NOT NULL;

-- Composite indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_leads_assigned_status 
    ON leads(assigned_to, lead_status) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_leads_source_status 
    ON leads(lead_source, lead_status) WHERE deleted_at IS NULL;

-- =====================================================
-- 4. ACTIVITIES TABLE INDEXES
-- =====================================================
-- Foreign key indexes
CREATE INDEX IF NOT EXISTS idx_activities_account_id 
    ON activities(account_id) WHERE deleted_at IS NULL AND account_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_activities_contact_id 
    ON activities(contact_id) WHERE deleted_at IS NULL AND contact_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_activities_deal_id 
    ON activities(deal_id) WHERE deleted_at IS NULL AND deal_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_activities_lead_id 
    ON activities(lead_id) WHERE deleted_at IS NULL AND lead_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_activities_user_id 
    ON activities(user_id) WHERE deleted_at IS NULL;

-- Type filtering (frequently used)
CREATE INDEX IF NOT EXISTS idx_activities_type 
    ON activities(type) WHERE deleted_at IS NULL;

-- Timestamp queries (for timeline)
CREATE INDEX IF NOT EXISTS idx_activities_timestamp 
    ON activities(timestamp) WHERE deleted_at IS NULL;

-- Composite index for common query patterns
CREATE INDEX IF NOT EXISTS idx_activities_account_timestamp 
    ON activities(account_id, timestamp) WHERE deleted_at IS NULL AND account_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_activities_type_timestamp 
    ON activities(type, timestamp) WHERE deleted_at IS NULL;

-- =====================================================
-- 5. TASKS TABLE INDEXES
-- =====================================================
-- Foreign key indexes
CREATE INDEX IF NOT EXISTS idx_tasks_assigned_to 
    ON tasks(assigned_to) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_tasks_assigned_from 
    ON tasks(assigned_from) WHERE deleted_at IS NULL AND assigned_from IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_tasks_account_id 
    ON tasks(account_id) WHERE deleted_at IS NULL AND account_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_tasks_contact_id 
    ON tasks(contact_id) WHERE deleted_at IS NULL AND contact_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_tasks_deal_id 
    ON tasks(deal_id) WHERE deleted_at IS NULL AND deal_id IS NOT NULL;

-- Status filtering (frequently used)
CREATE INDEX IF NOT EXISTS idx_tasks_status 
    ON tasks(status) WHERE deleted_at IS NULL;

-- Priority filtering
CREATE INDEX IF NOT EXISTS idx_tasks_priority 
    ON tasks(priority) WHERE deleted_at IS NULL;

-- Due date queries (for upcoming tasks)
CREATE INDEX IF NOT EXISTS idx_tasks_due_date 
    ON tasks(due_date) WHERE deleted_at IS NULL AND due_date IS NOT NULL;

-- Composite indexes for common query patterns
CREATE INDEX IF NOT EXISTS idx_tasks_assigned_status 
    ON tasks(assigned_to, status) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_tasks_status_due_date 
    ON tasks(status, due_date) WHERE deleted_at IS NULL AND due_date IS NOT NULL;

-- =====================================================
-- 6. ACCOUNTS TABLE INDEXES
-- =====================================================
-- Foreign key indexes
CREATE INDEX IF NOT EXISTS idx_accounts_category_id 
    ON accounts(category_id) WHERE deleted_at IS NULL AND category_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_accounts_assigned_to 
    ON accounts(assigned_to) WHERE deleted_at IS NULL AND assigned_to IS NOT NULL;

-- Status filtering (frequently used)
CREATE INDEX IF NOT EXISTS idx_accounts_status 
    ON accounts(status) WHERE deleted_at IS NULL;

-- Search columns (for LIKE queries - consider full-text search for better performance)
CREATE INDEX IF NOT EXISTS idx_accounts_name 
    ON accounts(name) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_accounts_city 
    ON accounts(city) WHERE deleted_at IS NULL AND city IS NOT NULL;

-- =====================================================
-- 7. CONTACTS TABLE INDEXES
-- =====================================================
-- Foreign key indexes
CREATE INDEX IF NOT EXISTS idx_contacts_account_id 
    ON contacts(account_id) WHERE deleted_at IS NULL AND account_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_contacts_assigned_to 
    ON contacts(assigned_to) WHERE deleted_at IS NULL AND assigned_to IS NOT NULL;

-- Search columns
CREATE INDEX IF NOT EXISTS idx_contacts_name 
    ON contacts(name) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_contacts_email 
    ON contacts(email) WHERE deleted_at IS NULL AND email IS NOT NULL;

-- =====================================================
-- 8. USERS TABLE INDEXES
-- =====================================================
-- Foreign key indexes
CREATE INDEX IF NOT EXISTS idx_users_role_id 
    ON users(role_id) WHERE deleted_at IS NULL AND role_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_users_group_id 
    ON users(group_id) WHERE deleted_at IS NULL AND group_id IS NOT NULL;

-- Status filtering (frequently used)
CREATE INDEX IF NOT EXISTS idx_users_status 
    ON users(status) WHERE deleted_at IS NULL;

-- Email search (for login)
CREATE INDEX IF NOT EXISTS idx_users_email 
    ON users(email) WHERE deleted_at IS NULL;

-- Composite index for common query patterns
CREATE INDEX IF NOT EXISTS idx_users_role_status 
    ON users(role_id, status) WHERE deleted_at IS NULL AND role_id IS NOT NULL;

-- =====================================================
-- 9. MONTHLY TARGETS TABLE INDEXES
-- =====================================================
-- Foreign key indexes
CREATE INDEX IF NOT EXISTS idx_monthly_targets_user_id 
    ON monthly_targets(user_id) WHERE deleted_at IS NULL AND user_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_monthly_targets_group_id 
    ON monthly_targets(group_id) WHERE deleted_at IS NULL AND group_id IS NOT NULL;

-- Date range queries (for effective target lookup)
CREATE INDEX IF NOT EXISTS idx_monthly_targets_year_month 
    ON monthly_targets(year, month) WHERE deleted_at IS NULL;

-- Composite index for effective target queries
CREATE INDEX IF NOT EXISTS idx_monthly_targets_user_year_month 
    ON monthly_targets(user_id, year, month) WHERE deleted_at IS NULL AND user_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_monthly_targets_group_year_month 
    ON monthly_targets(group_id, year, month) WHERE deleted_at IS NULL AND group_id IS NOT NULL;

-- =====================================================
-- 10. PRODUCT ITEMS (DEAL PRODUCT ITEMS) INDEXES
-- =====================================================
-- Foreign key indexes
CREATE INDEX IF NOT EXISTS idx_deal_product_items_deal_id 
    ON deal_product_items(deal_id) WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_deal_product_items_product_id 
    ON deal_product_items(product_id) WHERE deleted_at IS NULL;

-- =====================================================
-- 11. REMINDERS TABLE INDEXES
-- =====================================================
-- Foreign key indexes
CREATE INDEX IF NOT EXISTS idx_reminders_task_id 
    ON reminders(task_id) WHERE deleted_at IS NULL;

-- Date queries (for pending reminders)
CREATE INDEX IF NOT EXISTS idx_reminders_remind_at 
    ON reminders(remind_at) WHERE deleted_at IS NULL;

-- Composite index for pending reminders query
CREATE INDEX IF NOT EXISTS idx_reminders_pending 
    ON reminders(remind_at, is_sent) WHERE deleted_at IS NULL AND is_sent = false;

-- =====================================================
-- 12. NOTIFICATIONS TABLE INDEXES
-- =====================================================
-- Foreign key indexes
CREATE INDEX IF NOT EXISTS idx_notifications_user_id 
    ON notifications(user_id) WHERE deleted_at IS NULL;

-- Read status filtering
CREATE INDEX IF NOT EXISTS idx_notifications_is_read 
    ON notifications(is_read) WHERE deleted_at IS NULL;

-- Date queries (for recent notifications)
CREATE INDEX IF NOT EXISTS idx_notifications_created_at 
    ON notifications(created_at) WHERE deleted_at IS NULL;

-- Composite index for unread notifications
CREATE INDEX IF NOT EXISTS idx_notifications_user_unread 
    ON notifications(user_id, is_read, created_at) WHERE deleted_at IS NULL AND is_read = false;

-- =====================================================
-- NOTES
-- =====================================================
-- These indexes are optimized for:
-- 1. Foreign key lookups (JOIN operations)
-- 2. Status filtering (WHERE status = ...)
-- 3. Date range queries (WHERE date >= ... AND date <= ...)
-- 4. Composite queries (WHERE user_id = ... AND status = ...)
-- 5. Soft-delete queries (WHERE deleted_at IS NULL)

-- Partial indexes (WHERE deleted_at IS NULL) are used to:
-- - Reduce index size
-- - Improve query performance for active records
-- - Support soft-delete pattern

-- For very large tables, consider:
-- - Partitioning by date ranges
-- - Materialized views for complex aggregations
-- - Full-text search indexes for text search columns

