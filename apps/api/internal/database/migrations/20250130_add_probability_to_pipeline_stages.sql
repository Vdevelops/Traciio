-- Migration: Add Probability Column to Pipeline Stages
-- Date: 2025-01-30
-- Description: Adds probability column to pipeline_stages table to store probability percentage (0-100) for each stage

-- Add probability column if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_schema = CURRENT_SCHEMA() 
        AND table_name = 'pipeline_stages' 
        AND column_name = 'probability'
    ) THEN
        ALTER TABLE pipeline_stages 
        ADD COLUMN probability INTEGER NOT NULL DEFAULT 0;
        
        -- Add check constraint to ensure probability is between 0 and 100
        ALTER TABLE pipeline_stages 
        ADD CONSTRAINT chk_pipeline_stages_probability_range 
        CHECK (probability >= 0 AND probability <= 100);
        
        -- Update existing stages with default probability based on order
        -- Qualification (order 1) = 25%
        UPDATE pipeline_stages SET probability = 25 WHERE code = 'qualification' AND probability = 0;
        
        -- Proposal (order 2) = 50%
        UPDATE pipeline_stages SET probability = 50 WHERE code = 'proposal' AND probability = 0;
        
        -- Negotiation (order 3) = 75%
        UPDATE pipeline_stages SET probability = 75 WHERE code = 'negotiation' AND probability = 0;
        
        -- Closed Won (order 4) = 100%
        UPDATE pipeline_stages SET probability = 100 WHERE code = 'closed_won' AND probability = 0;
        
        -- Closed Lost (order 5) = 0%
        UPDATE pipeline_stages SET probability = 0 WHERE code = 'closed_lost' AND probability = 0;
        
        -- For any other stages, set probability based on order * 20 (fallback formula)
        UPDATE pipeline_stages 
        SET probability = LEAST(order * 20, 100) 
        WHERE probability = 0 AND code NOT IN ('qualification', 'proposal', 'negotiation', 'closed_won', 'closed_lost');
        
        COMMENT ON COLUMN pipeline_stages.probability IS 'Probability percentage (0-100) for deals in this stage';
    END IF;
END $$;




