-- Seed default lead statuses with scores
-- This should be run after creating the lead_statuses table

-- Check if statuses already exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM lead_statuses WHERE code = 'NEW') THEN
        -- Insert default statuses
        INSERT INTO lead_statuses (name, code, description, score, color, "order", is_active, is_default, is_converted, created_at, updated_at)
        VALUES
            ('New', 'NEW', 'New lead, first contact not yet made', 5, '#94A3B8', 1, true, true, false, NOW(), NOW()),
            ('Contacted', 'CONTACTED', 'Initial contact has been made', 15, '#3B82F6', 2, true, false, false, NOW(), NOW()),
            ('Engaged', 'ENGAGED', 'Lead is showing interest', 30, '#8B5CF6', 3, true, false, false, NOW(), NOW()),
            ('Interested', 'INTERESTED', 'Lead has expressed clear interest', 50, '#06B6D4', 4, true, false, false, NOW(), NOW()),
            ('Qualified', 'QUALIFIED', 'Lead meets qualification criteria', 75, '#10B981', 5, true, false, false, NOW(), NOW()),
            ('Converted', 'CONVERTED', 'Lead has been converted to opportunity', 100, '#22C55E', 6, true, false, true, NOW(), NOW()),
            ('Nurturing', 'NURTURING', 'Lead needs more nurturing before conversion', 25, '#F59E0B', 7, true, false, false, NOW(), NOW()),
            ('Unqualified', 'UNQUALIFIED', 'Lead does not meet qualification criteria', 0, '#6B7280', 8, true, false, false, NOW(), NOW()),
            ('Disqualified', 'DISQUALIFIED', 'Lead has been disqualified', 0, '#EF4444', 9, true, false, false, NOW(), NOW()),
            ('Lost', 'LOST', 'Lead was lost to competitor or no longer interested', 0, '#DC2626', 10, true, false, false, NOW(), NOW());
        
        RAISE NOTICE 'Successfully seeded % lead statuses', (SELECT COUNT(*) FROM lead_statuses);
    ELSE
        RAISE NOTICE 'Lead statuses already exist, skipping seed';
    END IF;
END $$;
