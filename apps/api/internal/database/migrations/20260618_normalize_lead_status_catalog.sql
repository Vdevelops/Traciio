-- Normalize lead statuses to the canonical CRM flow:
-- new -> contacted -> interested -> qualified -> proposal_sent -> converted/lost

WITH desired_statuses AS (
    SELECT *
    FROM (
        VALUES
            ('New', 'new', 'Brand new lead, not yet contacted', 5, '#94A3B8', 1, TRUE, FALSE),
            ('Contacted', 'contacted', 'Initial contact has been made', 20, '#60A5FA', 2, FALSE, FALSE),
            ('Interested', 'interested', 'Lead has shown strong interest in our products or services', 45, '#14B8A6', 3, FALSE, FALSE),
            ('Qualified', 'qualified', 'Lead meets qualification criteria and is ready for proposal', 70, '#22C55E', 4, FALSE, FALSE),
            ('Proposal Sent', 'proposal_sent', 'Proposal has been sent and is awaiting a decision', 85, '#0EA5E9', 5, FALSE, FALSE),
            ('Converted', 'converted', 'Lead has been converted to opportunity', 100, '#2563EB', 6, FALSE, TRUE),
            ('Lost', 'lost', 'Lead was lost to competitor or no longer interested', 0, '#DC2626', 7, FALSE, FALSE)
    ) AS v(name, code, description, score, color, display_order, is_default, is_converted)
)
INSERT INTO lead_statuses (name, code, description, score, color, "order", is_active, is_default, is_converted, created_at, updated_at)
SELECT name, code, description, score, color, display_order, TRUE, is_default, is_converted, NOW(), NOW()
FROM desired_statuses ds
WHERE NOT EXISTS (
    SELECT 1
    FROM lead_statuses ls
    WHERE ls.deleted_at IS NULL
      AND LOWER(ls.code) = ds.code
);

WITH desired_statuses AS (
    SELECT *
    FROM (
        VALUES
            ('New', 'new', 'Brand new lead, not yet contacted', 5, '#94A3B8', 1, TRUE, FALSE),
            ('Contacted', 'contacted', 'Initial contact has been made', 20, '#60A5FA', 2, FALSE, FALSE),
            ('Interested', 'interested', 'Lead has shown strong interest in our products or services', 45, '#14B8A6', 3, FALSE, FALSE),
            ('Qualified', 'qualified', 'Lead meets qualification criteria and is ready for proposal', 70, '#22C55E', 4, FALSE, FALSE),
            ('Proposal Sent', 'proposal_sent', 'Proposal has been sent and is awaiting a decision', 85, '#0EA5E9', 5, FALSE, FALSE),
            ('Converted', 'converted', 'Lead has been converted to opportunity', 100, '#2563EB', 6, FALSE, TRUE),
            ('Lost', 'lost', 'Lead was lost to competitor or no longer interested', 0, '#DC2626', 7, FALSE, FALSE)
    ) AS v(name, code, description, score, color, display_order, is_default, is_converted)
)
UPDATE lead_statuses ls
SET name = ds.name,
    code = ds.code,
    description = ds.description,
    score = ds.score,
    color = ds.color,
    "order" = ds.display_order,
    is_active = TRUE,
    is_default = ds.is_default,
    is_converted = ds.is_converted,
    updated_at = NOW()
FROM desired_statuses ds
WHERE ls.deleted_at IS NULL
  AND LOWER(ls.code) = ds.code;

UPDATE lead_statuses
SET is_active = FALSE,
    is_default = FALSE,
    is_converted = FALSE,
    updated_at = NOW()
WHERE deleted_at IS NULL
  AND LOWER(code) NOT IN ('new', 'contacted', 'interested', 'qualified', 'proposal_sent', 'converted', 'lost');

UPDATE leads
SET lead_status = CASE LOWER(COALESCE(lead_status, ''))
    WHEN 'engaged' THEN 'interested'
    WHEN 'nurturing' THEN 'contacted'
    WHEN 'unqualified' THEN 'lost'
    WHEN 'disqualified' THEN 'lost'
    ELSE LOWER(lead_status)
END,
updated_at = NOW()
WHERE deleted_at IS NULL
  AND lead_status IS NOT NULL;

UPDATE leads l
SET lead_status_id = ls.id,
    lead_score = COALESCE(ls.score, l.lead_score),
    updated_at = NOW()
FROM lead_statuses ls
WHERE l.deleted_at IS NULL
  AND ls.deleted_at IS NULL
  AND ls.is_active = TRUE
  AND LOWER(COALESCE(l.lead_status, '')) = LOWER(ls.code)
  AND (
      l.lead_status_id IS DISTINCT FROM ls.id
      OR l.lead_score IS DISTINCT FROM ls.score
  );
