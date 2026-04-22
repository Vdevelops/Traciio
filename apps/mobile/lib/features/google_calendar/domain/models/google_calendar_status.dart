/// Google Calendar connection status model
class GoogleCalendarStatus {
  final bool isConnected;
  final String? email;
  final DateTime? connectedAt;
  final String? errorMessage;

  const GoogleCalendarStatus({
    this.isConnected = false,
    this.email,
    this.connectedAt,
    this.errorMessage,
  });

  factory GoogleCalendarStatus.fromJson(Map<String, dynamic> json) {
    return GoogleCalendarStatus(
      isConnected: json['is_connected'] ?? false,
      email: json['email'],
      connectedAt: json['connected_at'] != null
          ? DateTime.parse(json['connected_at'])
          : null,
      errorMessage: json['error_message'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'is_connected': isConnected,
      'email': email,
      'connected_at': connectedAt?.toIso8601String(),
      'error_message': errorMessage,
    };
  }

  GoogleCalendarStatus copyWith({
    bool? isConnected,
    String? email,
    DateTime? connectedAt,
    String? errorMessage,
  }) {
    return GoogleCalendarStatus(
      isConnected: isConnected ?? this.isConnected,
      email: email ?? this.email,
      connectedAt: connectedAt ?? this.connectedAt,
      errorMessage: errorMessage ?? this.errorMessage,
    );
  }
}
