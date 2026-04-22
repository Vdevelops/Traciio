import '../data/models/notification.dart';

class NotificationListState {
  final List<Notification> notifications;
  final Pagination? pagination;
  final bool isLoading;
  final String? errorMessage;
  final String? selectedFilter; // 'all', 'unread', 'read'
  final NotificationType? selectedType;

  const NotificationListState({
    this.notifications = const [],
    this.pagination,
    this.isLoading = false,
    this.errorMessage,
    this.selectedFilter,
    this.selectedType,
  });

  NotificationListState copyWith({
    List<Notification>? notifications,
    Pagination? pagination,
    bool? isLoading,
    String? errorMessage,
    String? selectedFilter,
    NotificationType? selectedType,
  }) {
    return NotificationListState(
      notifications: notifications ?? this.notifications,
      pagination: pagination ?? this.pagination,
      isLoading: isLoading ?? this.isLoading,
      errorMessage: errorMessage ?? this.errorMessage,
      selectedFilter: selectedFilter ?? this.selectedFilter,
      selectedType: selectedType ?? this.selectedType,
    );
  }
}

class NotificationCountState {
  final int unreadCount;
  final bool isLoading;
  final String? errorMessage;

  const NotificationCountState({
    this.unreadCount = 0,
    this.isLoading = false,
    this.errorMessage,
  });

  NotificationCountState copyWith({
    int? unreadCount,
    bool? isLoading,
    String? errorMessage,
  }) {
    return NotificationCountState(
      unreadCount: unreadCount ?? this.unreadCount,
      isLoading: isLoading ?? this.isLoading,
      errorMessage: errorMessage ?? this.errorMessage,
    );
  }
}

