import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/network/api_client.dart';
import '../data/models/notification.dart';
import '../data/notification_repository.dart';
import 'notification_state.dart';

final notificationRepositoryProvider = Provider<NotificationRepository>((ref) {
  return NotificationRepository(ApiClient.dio);
});

final notificationListProvider =
    NotifierProvider<NotificationListNotifier, NotificationListState>(NotificationListNotifier.new);

final notificationCountProvider =
    NotifierProvider<NotificationCountNotifier, NotificationCountState>(NotificationCountNotifier.new);

class NotificationListNotifier extends Notifier<NotificationListState> {
  late final NotificationRepository _repository;

  @override
  NotificationListState build() {
    _repository = ref.read(notificationRepositoryProvider);
    return const NotificationListState();
  }

  Future<void> loadNotifications({
    int page = 1,
    bool refresh = false,
    String? filter,
    NotificationType? type,
  }) async {
    if (refresh || page == 1) {
      state = state.copyWith(isLoading: true, errorMessage: null);
    } else {
      state = state.copyWith(isLoading: true);
    }

    try {
      final filterValue = filter ?? state.selectedFilter ?? 'all';
      final isReadFilter = filterValue == 'read'
          ? true
          : filterValue == 'unread'
              ? false
              : null;

      final response = await _repository.getNotifications(
        page: page,
        perPage: 20,
        type: type ?? state.selectedType,
        isRead: isReadFilter,
      );

      if (refresh || page == 1) {
        state = state.copyWith(
          notifications: response.items,
          pagination: response.pagination,
          selectedFilter: filterValue,
          selectedType: type ?? state.selectedType,
          isLoading: false,
          errorMessage: null,
        );
      } else {
        state = state.copyWith(
          notifications: [...state.notifications, ...response.items],
          pagination: response.pagination,
          isLoading: false,
          errorMessage: null,
        );
      }
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        errorMessage: e.toString().replaceFirst('Exception: ', ''),
      );
    }
  }

  Future<void> refresh() async {
    await loadNotifications(page: 1, refresh: true);
  }

  Future<void> loadMore() async {
    if (state.isLoading) return;
    final pagination = state.pagination;
    if (pagination == null || !pagination.hasNextPage) return;

    await loadNotifications(page: pagination.page + 1);
  }

  void updateFilter(String? filter) {
    state = state.copyWith(selectedFilter: filter);
    loadNotifications(page: 1, refresh: true, filter: filter);
  }

  void updateTypeFilter(NotificationType? type) {
    state = state.copyWith(selectedType: type);
    loadNotifications(page: 1, refresh: true, type: type);
  }

  Future<bool> markAsRead(String id) async {
    try {
      await _repository.markAsRead(id);
      // Update local state
      state = state.copyWith(
        notifications: state.notifications.map((n) {
          if (n.id == id) {
            return Notification(
              id: n.id,
              userId: n.userId,
              title: n.title,
              message: n.message,
              type: n.type,
              isRead: true,
              readAt: DateTime.now(),
              data: n.data,
              createdAt: n.createdAt,
              updatedAt: n.updatedAt,
            );
          }
          return n;
        }).toList(),
      );
      return true;
    } catch (e) {
      return false;
    }
  }

  Future<bool> markAllAsRead() async {
    try {
      await _repository.markAllAsRead();
      // Update local state
      state = state.copyWith(
        notifications: state.notifications.map((n) {
          return Notification(
            id: n.id,
            userId: n.userId,
            title: n.title,
            message: n.message,
            type: n.type,
            isRead: true,
            readAt: DateTime.now(),
            data: n.data,
            createdAt: n.createdAt,
            updatedAt: n.updatedAt,
          );
        }).toList(),
      );
      return true;
    } catch (e) {
      return false;
    }
  }

  Future<bool> deleteNotification(String id) async {
    try {
      await _repository.deleteNotification(id);
      // Remove from local state
      state = state.copyWith(
        notifications: state.notifications.where((n) => n.id != id).toList(),
      );
      return true;
    } catch (e) {
      return false;
    }
  }
}

class NotificationCountNotifier extends Notifier<NotificationCountState> {
  late final NotificationRepository _repository;

  @override
  NotificationCountState build() {
    _repository = ref.read(notificationRepositoryProvider);
    return const NotificationCountState();
  }

  Future<void> loadUnreadCount() async {
    state = state.copyWith(isLoading: true, errorMessage: null);

    try {
      final count = await _repository.getUnreadCount();
      state = state.copyWith(
        unreadCount: count,
        isLoading: false,
        errorMessage: null,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        errorMessage: e.toString().replaceFirst('Exception: ', ''),
      );
    }
  }

  Future<void> refresh() async {
    await loadUnreadCount();
  }
}

