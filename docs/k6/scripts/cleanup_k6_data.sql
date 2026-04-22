-- ============================================================
-- cleanup_k6_data.sql
-- Clean up test data created by K6 load tests
-- Run this AFTER load testing to remove synthetic data
-- ============================================================

-- =============================================
-- Step 1: Preview — see how much test data exists
-- =============================================
SELECT '=== BEFORE CLEANUP ===' AS phase;

-- Count K6-generated leads by notes pattern
SELECT notes, COUNT(*) AS total
FROM leads
WHERE first_name = 'K6'
GROUP BY notes;

-- Count K6-generated leads by lead_source/first_name
SELECT lead_source, first_name, COUNT(*) AS count
FROM leads
WHERE first_name = 'K6'
GROUP BY lead_source, first_name;

-- =============================================
-- Step 2: Delete test data (within transaction)
-- =============================================

BEGIN;

-- Delete leads created by K6 smoke test
DELETE FROM leads WHERE first_name = 'K6' AND notes LIKE 'K6 smoke test%';

-- Delete leads created by K6 load test
DELETE FROM leads WHERE first_name = 'K6' AND notes LIKE 'K6 load test%';

-- Delete leads created by K6 stress test
DELETE FROM leads WHERE first_name = 'K6' AND notes LIKE 'K6 stress test%';

-- Delete leads created by K6 soak test
DELETE FROM leads WHERE first_name = 'K6' AND notes LIKE 'K6 soak test%';

-- Catch-all: Delete any remaining K6-generated leads
DELETE FROM leads WHERE first_name = 'K6';

-- =============================================
-- Step 3: Verify cleanup
-- =============================================
SELECT '=== AFTER CLEANUP ===' AS phase;

SELECT COUNT(*) AS remaining_k6_leads
FROM leads
WHERE first_name = 'K6';

COMMIT;

-- =============================================
-- Optional: Reclaim space after large deletes
-- =============================================
-- VACUUM ANALYZE leads;
