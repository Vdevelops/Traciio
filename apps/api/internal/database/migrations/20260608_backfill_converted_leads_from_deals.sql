-- Backfill leads that were converted through pipeline deal creation before
-- lead_status_id synchronization existed in that flow.

WITH converted_status AS (
    SELECT id, score
    FROM lead_statuses
    WHERE deleted_at IS NULL
      AND (is_converted = TRUE OR code IN ('CONVERTED', 'converted'))
    ORDER BY is_converted DESC, "order" ASC
    LIMIT 1
),
latest_deals AS (
    SELECT DISTINCT ON (lead_id)
        id,
        lead_id,
        account_id,
        contact_id,
        created_by,
        created_at
    FROM deals
    WHERE deleted_at IS NULL
      AND lead_id IS NOT NULL
    ORDER BY lead_id, created_at DESC
)
UPDATE leads AS l
SET lead_status = 'converted',
    lead_status_id = (SELECT id FROM converted_status),
    lead_score = COALESCE((SELECT score FROM converted_status), l.lead_score),
    opportunity_id = COALESCE(l.opportunity_id, d.id),
    account_id = COALESCE(l.account_id, d.account_id),
    contact_id = COALESCE(l.contact_id, d.contact_id),
    converted_at = COALESCE(l.converted_at, d.created_at),
    converted_by = COALESCE(l.converted_by, d.created_by),
    updated_at = NOW()
FROM latest_deals AS d
WHERE l.id = d.lead_id
  AND l.deleted_at IS NULL
  AND (
      LOWER(l.lead_status) <> 'converted'
      OR l.lead_status_id IS DISTINCT FROM (SELECT id FROM converted_status)
      OR l.opportunity_id IS NULL
      OR l.converted_at IS NULL
  );
