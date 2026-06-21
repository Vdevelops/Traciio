-- Seed default lead statuses with scores
-- This should be run after creating the lead_statuses table

-- Check if statuses already exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM lead_statuses WHERE LOWER(code) = 'new') THEN
        -- Insert default statuses
        INSERT INTO lead_statuses (name, code, description, score, color, "order", is_active, is_default, is_converted, created_at, updated_at)
        VALUES
            ('New', 'new', 'Brand new lead, not yet contacted', 5, '#94A3B8', 1, true, true, false, NOW(), NOW()),
            ('Contacted', 'contacted', 'Initial contact has been made', 20, '#60A5FA', 2, true, false, false, NOW(), NOW()),
            ('Interested', 'interested', 'Lead has shown strong interest in our products or services', 45, '#14B8A6', 3, true, false, false, NOW(), NOW()),
            ('Qualified', 'qualified', 'Lead meets qualification criteria and is ready for proposal', 70, '#22C55E', 4, true, false, false, NOW(), NOW()),
            ('Proposal Sent', 'proposal_sent', 'Proposal has been sent and is awaiting a decision', 85, '#0EA5E9', 5, true, false, false, NOW(), NOW()),
            ('Converted', 'converted', 'Lead has been converted to opportunity', 100, '#2563EB', 6, true, false, true, NOW(), NOW()),
            ('Lost', 'lost', 'Lead was lost to competitor or no longer interested', 0, '#DC2626', 7, true, false, false, NOW(), NOW());
        
        RAISE NOTICE 'Successfully seeded % lead statuses', (SELECT COUNT(*) FROM lead_statuses);
    ELSE
        RAISE NOTICE 'Lead statuses already exist, skipping seed';
    END IF;
END $$;
