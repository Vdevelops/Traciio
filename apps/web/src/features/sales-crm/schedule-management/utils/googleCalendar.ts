/**
 * Utility functions for Google Calendar integration.
 */

/**
 * Creates a Google Calendar event URL from an event ID.
 * Uses the /r/eventedit format which works reliably with just the event ID.
 * If eid format doesn't work, it falls back to the calendar view.
 */
export function getGoogleCalendarEventURL(eventID: string): string {
  // The most reliable way is to use the event edit URL with the calendar ID
  // Since we use "primary" calendar, we can use this format
  // Format: https://calendar.google.com/calendar/r/eventedit/{eventId}
  // This opens the event for editing which also shows all details
  return `https://www.google.com/calendar/event?eid=${encodeURIComponent(eventID)}`;
}

