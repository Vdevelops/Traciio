-- HIGH PRIORITY INDEXES untuk Performance Optimization (1000-5000 users)
-- Run these in order using: psql -d crm_healthcare -f database_indexes.sql

-- =============================================================================
-- VISIT REPORTS - Full Text Search Optimization
-- =============================================================================
-- Issue: LIKE queries cause full table scans
-- Fix: GIN index for full-text search

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_visit_reports_fts 
ON visit_reports USING GIN (to_tsvector('indonesian', purpose || ' ' || COALESCE(notes, '')));

-- Index untuk date range queries (dashboard aggregation)
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_visit_reports_date_user_status 
ON visit_reports(visit_date, sales_rep_id, status) 
WHERE visit_date IS NOT NULL;

-- Index untuk lead-based queries
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_visit_reports_lead_id 
ON visit_reports(lead_id) 
WHERE lead_id IS NOT NULL;

-- =============================================================================
-- DEALS - Aggregation Queries Optimization
-- =============================================================================
-- Issue: Dashboard revenue queries are slow
-- Fix: Composite indexes untuk scoped user queries

-- Untuk GetWonDealsValueInPeriodByUser dan batch version
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_deals_user_status_date 
ON deals(assigned_to, status, actual_close_date, created_at) 
WHERE status = 'won';

-- Untuk deal listing dengan filters
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_deals_status_pipeline 
ON deals(status, pipeline_stage_id, created_at DESC);

-- Untuk account-based deal queries
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_deals_account_status 
ON deals(account_id, status, created_at DESC);

-- =============================================================================
-- ACTIVITIES - Timeline dan Stats Queries
-- =============================================================================
-- Issue: Activity stats aggregation slow
-- Fix: Composite indexes untuk timestamp-based queries

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_activities_user_timestamp 
ON activities(user_id, timestamp DESC, type);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_activities_account_timestamp 
ON activities(account_id, timestamp DESC);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_activities_lead_timestamp 
ON activities(lead_id, timestamp DESC) 
WHERE lead_id IS NOT NULL;

-- =============================================================================
-- LEADS - List dan Conversion Queries
-- =============================================================================
-- Issue: Lead listing dengan status filters
-- Fix: Partial indexes untuk common status values

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_leads_status_created 
ON leads(status, created_at DESC);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_leads_owner_status 
ON leads(assigned_to, status, created_at DESC);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_leads_qualified 
ON leads(status, created_at DESC) 
WHERE status IN ('new', 'contacted', 'qualified', 'converted');

-- =============================================================================
-- ACCOUNTS - Search dan List Queries
-- =============================================================================
-- Issue: Account search slow
-- Fix: GIN index untuk name search

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_accounts_name_gin 
ON accounts USING GIN (to_tsvector('indonesian', name || ' ' || COALESCE(address, '') || ' ' || COALESCE(city, '')));

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_accounts_status_category 
ON accounts(status, category_id, created_at DESC);

-- =============================================================================
-- CONTACTS - Relationship Queries
-- =============================================================================

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_contacts_account_id 
ON contacts(account_id, created_at DESC);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_contacts_email 
ON contacts(email) 
WHERE email IS NOT NULL;

-- =============================================================================
-- ROUTE OPTIMIZATIONS - User-based Queries
-- =============================================================================

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_route_optimizations_user_date 
ON route_optimizations(user_id, created_at DESC);

-- =============================================================================
-- TASKS - User Assignment Queries
-- =============================================================================

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_tasks_user_status 
ON tasks(assigned_to, status, due_date);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_tasks_account_status 
ON tasks(account_id, status, due_date);

-- =============================================================================
-- NOTIFICATIONS - Unread Count Queries
-- =============================================================================

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_notifications_user_read 
ON notifications(user_id, is_read, created_at DESC);

-- =============================================================================
-- USERS - Role-based Queries
-- =============================================================================

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_role_status 
ON users(role_id, status, created_at DESC);

-- =============================================================================
-- BRICKS - Manager-based Queries (Dashboard)
-- =============================================================================

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_bricks_manager 
ON bricks(manager_id, created_at DESC);

-- =============================================================================
-- PIPELINE STAGES - Ordering
-- =============================================================================

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_pipeline_stages_pipeline_order 
ON pipeline_stages(pipeline_id, "order");

-- =============================================================================
-- FOREIGN KEY INDEXES (untuk JOIN performance)
-- =============================================================================

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_deal_product_items_deal 
ON deal_product_items(deal_id);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_monthly_targets_user 
ON monthly_targets(user_id, year, month);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_brick_target_distributions_brick 
ON brick_target_distributions(brick_id);

-- =============================================================================
-- VERIFICATION
-- =============================================================================

-- Check index sizes and usage after 24 hours:
-- SELECT schemaname, tablename, indexname, pg_size_pretty(pg_relation_size(indexname::regclass)) as size
-- FROM pg_indexes
-- WHERE schemaname = 'public'
-- ORDER BY pg_relation_size(indexname::regclass) DESC;

-- Check index usage:
-- SELECT relname, indexrelname, idx_scan, idx_tup_read, idx_tup_fetch
-- FROM pg_stat_user_indexes
-- WHERE schemaname = 'public'
-- ORDER BY idx_scan DESC;

-- Note: Run these during low-traffic period as CONCURRENTLY still requires brief locks
-- Expected improvement: 50-80% faster dashboard queries, 30-50% faster list queries
