/// Google Calendar event sync data model
class GoogleCalendarEvent {
  final String? eventId;
  final String? eventLink;
  final String syncStatus; // not_synced | synced | sync_failed
  final DateTime? syncedAt;

  const GoogleCalendarEvent({
    this.eventId,
    this.eventLink,
    this.syncStatus = 'not_synced',
    this.syncedAt,
  });

  factory GoogleCalendarEvent.fromJson(Map<String, dynamic> json) {
    return GoogleCalendarEvent(
      eventId: json['google_calendar_event_id'],
      eventLink: json['google_calendar_event_link'],
      syncStatus: json['google_calendar_sync_status'] ?? 'not_synced',
      syncedAt: json['google_calendar_synced_at'] != null
          ? DateTime.parse(json['google_calendar_synced_at'])
          : null,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'google_calendar_event_id': eventId,
      'google_calendar_event_link': eventLink,
      'google_calendar_sync_status': syncStatus,
      'google_calendar_synced_at': syncedAt?.toIso8601String(),
    };
  }

  bool get isSynced => syncStatus == 'synced';
  bool get isSyncFailed => syncStatus == 'sync_failed';
  bool get isNotSynced => syncStatus == 'not_synced';
}
