-- Link lead-side timeline records to their converted deals so deal details can
-- display product interests captured before conversion.

WITH latest_deals AS (
    SELECT DISTINCT ON (lead_id)
        id,
        lead_id,
        account_id,
        created_at
    FROM deals
    WHERE deleted_at IS NULL
      AND lead_id IS NOT NULL
    ORDER BY lead_id, created_at DESC
)
UPDATE activities AS a
SET deal_id = d.id,
    account_id = COALESCE(a.account_id, d.account_id),
    updated_at = NOW()
FROM latest_deals AS d
WHERE a.lead_id = d.lead_id
  AND a.deleted_at IS NULL
  AND (a.deal_id IS NULL OR a.deal_id <> d.id);

WITH latest_deals AS (
    SELECT DISTINCT ON (lead_id)
        id,
        lead_id,
        account_id,
        created_at
    FROM deals
    WHERE deleted_at IS NULL
      AND lead_id IS NOT NULL
    ORDER BY lead_id, created_at DESC
)
UPDATE visit_reports AS vr
SET deal_id = d.id,
    account_id = COALESCE(vr.account_id, d.account_id),
    updated_at = NOW()
FROM latest_deals AS d
WHERE vr.lead_id = d.lead_id
  AND vr.deleted_at IS NULL
  AND (vr.deal_id IS NULL OR vr.deal_id <> d.id);

WITH latest_deals AS (
    SELECT DISTINCT ON (lead_id)
        id,
        lead_id,
        account_id,
        created_at
    FROM deals
    WHERE deleted_at IS NULL
      AND lead_id IS NOT NULL
    ORDER BY lead_id, created_at DESC
)
UPDATE tasks AS t
SET deal_id = d.id,
    account_id = COALESCE(t.account_id, d.account_id),
    updated_at = NOW()
FROM latest_deals AS d
WHERE t.lead_id = d.lead_id
  AND t.deleted_at IS NULL
  AND (t.deal_id IS NULL OR t.deal_id <> d.id);
