import 'dart:async';
import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../application/notification_provider.dart';
import '../application/notification_state.dart';
import '../data/models/notification.dart' as models;
import '../../../core/l10n/app_localizations.dart';
import '../../../core/widgets/error_widget.dart';
import '../../../core/widgets/loading_widget.dart';
import '../../../core/permissions/permission_provider.dart';
import 'widgets/notification_card.dart';

class NotificationListScreen extends ConsumerStatefulWidget {
  const NotificationListScreen({super.key});

  @override
  ConsumerState<NotificationListScreen> createState() =>
      _NotificationListScreenState();
}

class _NotificationListScreenState
    extends ConsumerState<NotificationListScreen> {
  final ScrollController _scrollController = ScrollController();

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(notificationListProvider.notifier).loadNotifications();
    });

    _scrollController.addListener(_onScroll);
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  void _onScroll() {
    if (_scrollController.position.pixels >=
        _scrollController.position.maxScrollExtent * 0.8) {
      ref.read(notificationListProvider.notifier).loadMore();
    }
  }

  Future<void> _onRefresh() async {
    await ref.read(notificationListProvider.notifier).refresh();
    await ref.read(notificationCountProvider.notifier).refresh();
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(notificationListProvider);
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.notifications),
        elevation: 0,
        actions: [
          Builder(
            builder: (context) {
              final canMarkRead = ref.watch(canMarkNotificationReadProvider);
              if (canMarkRead && state.notifications.any((n) => !n.isRead)) {
                return IconButton(
                  icon: const Icon(Icons.done_all),
                  tooltip: l10n.markAllAsRead,
                  onPressed: () =>
                      _handleMarkAllAsRead(context, l10n, colorScheme),
                );
              }
              return const SizedBox.shrink();
            },
          ),
        ],
      ),
      body: Column(
        children: [
          // Filter Chips
          _buildFilters(context, state, theme, colorScheme, l10n),
          // Content
          Expanded(
            child: RefreshIndicator(
              onRefresh: _onRefresh,
              child: _buildContent(context, state, theme, colorScheme, l10n),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildFilters(
    BuildContext context,
    NotificationListState state,
    ThemeData theme,
    ColorScheme colorScheme,
    AppLocalizations l10n,
  ) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: Row(
        children: [
          Expanded(
            child: _FilterChip(
              label: l10n.filter,
              value: state.selectedFilter == null
                  ? l10n.all
                  : state.selectedFilter == 'unread'
                  ? l10n.unread
                  : state.selectedFilter == 'read'
                  ? l10n.read
                  : l10n.all,
              onTap: () => _showFilterSheet(context, state, l10n),
            ),
          ),
        ],
      ),
    );
  }

  void _showFilterSheet(
    BuildContext context,
    NotificationListState state,
    AppLocalizations l10n,
  ) {
    showModalBottomSheet(
      context: context,
      builder: (_) => _FilterSheet(
        selectedFilter: state.selectedFilter,
        onSelect: (filter) {
          ref.read(notificationListProvider.notifier).updateFilter(filter);
          Navigator.pop(context);
        },
        l10n: l10n,
      ),
    );
  }

  Widget _buildContent(
    BuildContext context,
    NotificationListState state,
    ThemeData theme,
    ColorScheme colorScheme,
    AppLocalizations l10n,
  ) {
    if (state.isLoading && state.notifications.isEmpty) {
      return const LoadingWidget();
    }

    if (state.errorMessage != null) {
      final is404Error =
          state.errorMessage!.contains('404') ||
          state.errorMessage!.contains('not found') ||
          state.errorMessage!.contains('endpoint not found');

      return ErrorStateWidget(
        message: is404Error
            ? 'Notification endpoint not available.\nPlease ensure the backend server is running\nand notification routes are registered.'
            : state.errorMessage!,
        icon: is404Error ? Icons.wifi_off : null,
        onRetry: () => ref.read(notificationListProvider.notifier).refresh(),
      );
    }

    if (state.notifications.isEmpty) {
      return EmptyStateWidget(
        message: l10n.noNotificationsFound,
        icon: Icons.notifications_none,
      );
    }

    return ListView.builder(
      controller: _scrollController,
      itemCount:
          state.notifications.length +
          (state.pagination?.hasNextPage == true ? 1 : 0),
      itemBuilder: (context, index) {
        if (index == state.notifications.length) {
          return const Padding(
            padding: EdgeInsets.all(16.0),
            child: LoadingWidget(size: 24),
          );
        }
        final notification = state.notifications[index];
        return NotificationCard(
          notification: notification,
          onTap: () => _handleNotificationTap(context, notification),
          onMarkAsRead: () =>
              _handleMarkAsRead(context, notification, l10n, colorScheme),
          onDelete: () =>
              _handleDelete(context, notification, l10n, colorScheme),
        );
      },
    );
  }

  Future<void> _handleNotificationTap(
    BuildContext context,
    models.Notification notification,
  ) async {
    // Store Navigator before async operations
    final navigator = Navigator.of(context);

    // Mark as read if unread
    if (!notification.isRead) {
      await ref
          .read(notificationListProvider.notifier)
          .markAsRead(notification.id);
      await ref.read(notificationCountProvider.notifier).refresh();
    }

    // Navigate based on notification type and data
    if (!mounted) return;
    await _navigateToDetail(navigator, notification);
  }

  Future<void> _navigateToDetail(
    NavigatorState navigator,
    models.Notification notification,
  ) async {
    // Parse data if available
    String? taskId;

    if (notification.data != null && notification.data!.isNotEmpty) {
      try {
        // Parse JSON data to extract IDs
        final dataMap = jsonDecode(notification.data!) as Map<String, dynamic>;
        taskId = dataMap['task_id'] as String?;
        // reminderId can be extracted if needed in the future
        // final reminderId = dataMap['reminder_id'] as String?;
      } catch (e) {
        // Ignore parsing errors
      }
    }

    // Navigate based on notification type
    switch (notification.type) {
      case models.NotificationType.task:
        if (taskId != null) {
          await navigator.pushNamed('/tasks/$taskId');
        }
        break;
      case models.NotificationType.reminder:
        // Navigate to task if reminder is for a task
        if (taskId != null) {
          await navigator.pushNamed('/tasks/$taskId');
        }
        break;
      case models.NotificationType.deal:
        // Navigate to deal detail if available
        // Note: Deal navigation will be implemented when deal feature is available
        break;
      case models.NotificationType.activity:
        // Navigate to activity detail if available
        // Note: Activity navigation will be implemented when activity feature is available
        break;
    }
  }

  Future<void> _handleMarkAsRead(
    BuildContext context,
    models.Notification notification,
    AppLocalizations l10n,
    ColorScheme colorScheme,
  ) async {
    final messenger = ScaffoldMessenger.of(context);
    final success = await ref
        .read(notificationListProvider.notifier)
        .markAsRead(notification.id);

    if (!mounted) return;
    await ref.read(notificationCountProvider.notifier).refresh();
    if (!mounted) return;
    // messenger is safe to use here as it was captured before async gap

    if (success) {
      messenger.showSnackBar(
        SnackBar(
          content: Text(l10n.notificationMarkedAsRead),
          backgroundColor: Colors.green,
        ),
      );
    } else {
      messenger.showSnackBar(
        SnackBar(
          content: Text(l10n.failedToMarkAsRead),
          backgroundColor: Colors.red,
        ),
      );
    }
  }

  Future<void> _handleMarkAllAsRead(
    BuildContext context,
    AppLocalizations l10n,
    ColorScheme colorScheme,
  ) async {
    final messenger = ScaffoldMessenger.of(context);
    final success = await ref
        .read(notificationListProvider.notifier)
        .markAllAsRead();

    if (!mounted) return;
    await ref.read(notificationCountProvider.notifier).refresh();
    if (!mounted) return;
    // messenger is safe to use here as it was captured before async gap

    if (success) {
      messenger.showSnackBar(
        SnackBar(
          content: Text(l10n.allNotificationsMarkedAsRead),
          backgroundColor: Colors.green,
        ),
      );
    } else {
      messenger.showSnackBar(
        SnackBar(
          content: Text(l10n.failedToMarkAllAsRead),
          backgroundColor: Colors.red,
        ),
      );
    }
  }

  Future<void> _handleDelete(
    BuildContext context,
    models.Notification notification,
    AppLocalizations l10n,
    ColorScheme colorScheme,
  ) async {
    final messenger = ScaffoldMessenger.of(context);
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(
          l10n.deleteNotification,
        ), // Note: using l10n from parameter is fine, but we'll use local one forSnackBar if needed
        content: Text(l10n.deleteNotificationConfirmation),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: Text(l10n.cancel),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(context, true),
            style: FilledButton.styleFrom(backgroundColor: colorScheme.error),
            child: Text(l10n.delete),
          ),
        ],
      ),
    );

    if (!mounted) return;
    if (confirmed == true) {
      final success = await ref
          .read(notificationListProvider.notifier)
          .deleteNotification(notification.id);

      if (!mounted) return;
      await ref.read(notificationCountProvider.notifier).refresh();
      if (!mounted) return;

      if (success) {
        messenger.showSnackBar(
          SnackBar(
            content: Text(l10n.notificationDeleted),
            backgroundColor: Colors.green,
          ),
        );
      } else {
        messenger.showSnackBar(
          SnackBar(
            content: Text(l10n.failedToDeleteNotification),
            backgroundColor: Colors.red,
          ),
        );
      }
    }
  }
}

class _FilterChip extends StatelessWidget {
  const _FilterChip({
    required this.label,
    required this.value,
    required this.onTap,
  });

  final String label;
  final String value;
  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    return InkWell(
      onTap: onTap,
      borderRadius: BorderRadius.circular(8),
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        decoration: BoxDecoration(
          color: colorScheme.surface,
          borderRadius: BorderRadius.circular(8),
          border: Border.all(color: colorScheme.outline.withValues(alpha: 0.2)),
        ),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text(
              '$label: ',
              style: theme.textTheme.bodySmall?.copyWith(
                color: colorScheme.onSurface.withValues(alpha: 0.7),
              ),
            ),
            Text(
              value,
              style: theme.textTheme.bodySmall?.copyWith(
                fontWeight: FontWeight.w600,
                color: colorScheme.onSurface,
              ),
            ),
            const SizedBox(width: 4),
            Icon(
              Icons.keyboard_arrow_down,
              size: 16,
              color: colorScheme.onSurface.withValues(alpha: 0.7),
            ),
          ],
        ),
      ),
    );
  }
}

class _FilterSheet extends StatelessWidget {
  const _FilterSheet({
    this.selectedFilter,
    required this.onSelect,
    required this.l10n,
  });

  final String? selectedFilter;
  final ValueChanged<String?> onSelect;
  final AppLocalizations l10n;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final filters = <String?>[null, 'unread', 'read'];

    return Container(
      padding: const EdgeInsets.all(16),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            l10n.filterNotifications,
            style: theme.textTheme.titleLarge?.copyWith(
              fontWeight: FontWeight.bold,
              color: colorScheme.onSurface,
            ),
          ),
          const SizedBox(height: 16),
          ...filters.map((filter) {
            final label = filter == null
                ? l10n.all
                : filter == 'unread'
                ? l10n.unread
                : l10n.read;
            final isSelected = selectedFilter == filter;
            return ListTile(
              title: Text(label),
              trailing: isSelected
                  ? Icon(Icons.check, color: colorScheme.primary)
                  : null,
              onTap: () => onSelect(filter),
            );
          }),
        ],
      ),
    );
  }
}
