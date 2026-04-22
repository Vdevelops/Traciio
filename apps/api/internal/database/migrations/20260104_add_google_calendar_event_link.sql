-- Migration: Add google_calendar_event_link column to track direct Google Calendar event URL
-- This column stores the HtmlLink returned by Google Calendar API for reliable event viewing

-- Add google_calendar_event_link column to schedules table
ALTER TABLE schedules ADD COLUMN IF NOT EXISTS google_calendar_event_link VARCHAR(500);

-- Add google_calendar_event_link column to tasks table
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS google_calendar_event_link VARCHAR(500);

-- Add comments
COMMENT ON COLUMN schedules.google_calendar_event_link IS 'Direct URL to view the event in Google Calendar (from HtmlLink API response)';
COMMENT ON COLUMN tasks.google_calendar_event_link IS 'Direct URL to view the event in Google Calendar (from HtmlLink API response)';
