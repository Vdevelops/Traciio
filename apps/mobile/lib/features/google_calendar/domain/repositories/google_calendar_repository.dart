import '../models/google_calendar_event.dart';
import '../models/google_calendar_status.dart';

/// Abstract repository for Google Calendar operations
abstract class GoogleCalendarRepository {
  /// Get Google Calendar connection status
  Future<GoogleCalendarStatus> getConnectionStatus();

  /// Get OAuth authorization URL
  Future<String> getAuthUrl();

  /// Disconnect Google Calendar
  Future<void> disconnect();

  /// Sync schedule to Google Calendar
  Future<GoogleCalendarEvent> syncSchedule(String scheduleId);

  /// Unsync schedule from Google Calendar
  Future<void> unsyncSchedule(String scheduleId);

  /// Exchange authorization code for token (Option 2 - Direct Deep Link)
  Future<void> exchangeCode(String code, String state);
}
