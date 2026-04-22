-- Create google_calendar_tokens table
-- This table stores encrypted Google Calendar OAuth2 tokens for users

CREATE TABLE IF NOT EXISTS google_calendar_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    access_token TEXT NOT NULL,  -- Encrypted access token
    refresh_token TEXT NOT NULL, -- Encrypted refresh token
    token_type VARCHAR(50) NOT NULL DEFAULT 'Bearer',
    expires_at TIMESTAMP NOT NULL,
    scope TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_google_calendar_tokens_user_id ON google_calendar_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_google_calendar_tokens_expires_at ON google_calendar_tokens(expires_at);
CREATE INDEX IF NOT EXISTS idx_google_calendar_tokens_deleted_at ON google_calendar_tokens(deleted_at);

-- Add comment
COMMENT ON TABLE google_calendar_tokens IS 'Encrypted Google Calendar OAuth2 tokens for users';
COMMENT ON COLUMN google_calendar_tokens.access_token IS 'Encrypted OAuth2 access token';
COMMENT ON COLUMN google_calendar_tokens.refresh_token IS 'Encrypted OAuth2 refresh token';
COMMENT ON COLUMN google_calendar_tokens.expires_at IS 'Token expiration time';

