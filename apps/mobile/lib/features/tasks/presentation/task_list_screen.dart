import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../application/task_provider.dart';
import '../application/task_state.dart';
import '../../../core/routing/app_router.dart';
import '../../../core/l10n/app_localizations.dart';
import '../../../core/widgets/error_widget.dart';
import '../../../core/widgets/loading_widget.dart';
import '../../../core/widgets/skeleton_widget.dart';
import '../../../core/widgets/sync_status_indicator.dart';
import 'widgets/task_card.dart';
import 'widgets/task_filters.dart';
import 'widgets/task_search_modal.dart';

class TaskListScreen extends ConsumerStatefulWidget {
  const TaskListScreen({super.key, this.hideAppBar = false});

  final bool hideAppBar;

  @override
  ConsumerState<TaskListScreen> createState() => _TaskListScreenState();
}

class _TaskListScreenState extends ConsumerState<TaskListScreen> {
  // Removed custom ScrollController to allow NestedScrollView to handle scrolling

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(taskListProvider.notifier).loadTasks();
    });
  }

  bool _onScrollNotification(ScrollNotification notification) {
    if (notification is ScrollUpdateNotification) {
      if (notification.metrics.pixels >=
          notification.metrics.maxScrollExtent * 0.8) {
        ref.read(taskListProvider.notifier).loadMore();
      }
    }
    return false;
  }

  Future<void> _onRefresh() async {
    await ref.read(taskListProvider.notifier).refresh();
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(taskListProvider);
    final theme = Theme.of(context);
    final l10n = AppLocalizations.of(context)!;

    final hasActiveSearch = state.searchQuery.isNotEmpty;

    Widget bodyContent = NotificationListener<ScrollNotification>(
      onNotification: _onScrollNotification,
      child: RefreshIndicator(
        onRefresh: _onRefresh,
        child: _buildContent(context, state, theme),
      ),
    );

    final body = Column(
      children: [
        if (hasActiveSearch)
          Padding(
            padding: const EdgeInsets.fromLTRB(16, 8, 16, 0),
            child: Row(
              children: [
                Expanded(
                  child: Text(
                    '${l10n.searchResultsFor} "${state.searchQuery}"',
                    style: theme.textTheme.bodyMedium?.copyWith(
                      color: theme.colorScheme.onSurface.withValues(alpha: 0.7),
                    ),
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
                const SizedBox(width: 8),
                TextButton.icon(
                  onPressed: () {
                    ref.read(taskListProvider.notifier).updateSearchQuery('');
                    ref
                        .read(taskListProvider.notifier)
                        .loadTasks(page: 1, refresh: true, search: '');
                  },
                  icon: Icon(
                    Icons.close,
                    size: 16,
                    color: theme.colorScheme.error,
                  ),
                  label: Text(
                    l10n.clearSearch,
                    style: theme.textTheme.labelLarge?.copyWith(
                      color: theme.colorScheme.error,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  style: TextButton.styleFrom(
                    visualDensity: VisualDensity.compact,
                    padding: const EdgeInsets.symmetric(
                      horizontal: 12,
                      vertical: 4,
                    ),
                    backgroundColor: theme.colorScheme.error.withValues(
                      alpha: 0.1,
                    ),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(20),
                    ),
                  ),
                ),
              ],
            ),
          ),
        Expanded(child: bodyContent),
      ],
    );

    if (widget.hideAppBar) {
      return body;
    }

    final hasActiveFilters =
        state.selectedStatus != null ||
        state.selectedPriority != null ||
        state.selectedType != null ||
        state.dueDateFrom != null ||
        state.dueDateTo != null;

    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.tasks),
        actions: [
          // Sync Status Indicator
          const Padding(
            padding: EdgeInsets.only(right: 8),
            child: Center(
              child: SyncStatusIndicator(featureKey: 'tasks', iconSize: 20),
            ),
          ),
          // Search button
          IconButton(
            icon: const Icon(Icons.search),
            onPressed: () => showTaskSearchModal(context),
            tooltip: l10n.searchTasks,
          ),
          // Filter button
          PopupMenuButton<String>(
            icon: Stack(
              children: [
                const Icon(Icons.filter_list),
                if (hasActiveFilters)
                  Positioned(
                    right: 0,
                    top: 0,
                    child: Container(
                      width: 8,
                      height: 8,
                      decoration: BoxDecoration(
                        color: theme.colorScheme.error,
                        shape: BoxShape.circle,
                      ),
                    ),
                  ),
              ],
            ),
            tooltip: l10n.filter,
            onSelected: (value) {
              final state = ref.read(taskListProvider);
              switch (value) {
                case 'status':
                  showModalBottomSheet(
                    context: context,
                    shape: const RoundedRectangleBorder(
                      borderRadius: BorderRadius.vertical(
                        top: Radius.circular(20),
                      ),
                    ),
                    builder: (context) => StatusFilterSheet(
                      selectedStatus: state.selectedStatus,
                      onSelect: (status) {
                        Navigator.pop(context);
                        ref
                            .read(taskListProvider.notifier)
                            .updateStatusFilter(status);
                      },
                    ),
                  );
                  break;
                case 'priority':
                  showModalBottomSheet(
                    context: context,
                    shape: const RoundedRectangleBorder(
                      borderRadius: BorderRadius.vertical(
                        top: Radius.circular(20),
                      ),
                    ),
                    builder: (context) => PriorityFilterSheet(
                      selectedPriority: state.selectedPriority,
                      onSelect: (priority) {
                        Navigator.pop(context);
                        ref
                            .read(taskListProvider.notifier)
                            .updatePriorityFilter(priority);
                      },
                    ),
                  );
                  break;
                case 'type':
                  showModalBottomSheet(
                    context: context,
                    shape: const RoundedRectangleBorder(
                      borderRadius: BorderRadius.vertical(
                        top: Radius.circular(20),
                      ),
                    ),
                    builder: (context) => TypeFilterSheet(
                      selectedType: state.selectedType,
                      onSelect: (type) {
                        Navigator.pop(context);
                        ref
                            .read(taskListProvider.notifier)
                            .updateTypeFilter(type);
                      },
                    ),
                  );
                  break;
                case 'due_date':
                  showModalBottomSheet(
                    context: context,
                    shape: const RoundedRectangleBorder(
                      borderRadius: BorderRadius.vertical(
                        top: Radius.circular(20),
                      ),
                    ),
                    builder: (context) => DueDateFilterSheet(
                      dueDateFrom: state.dueDateFrom,
                      dueDateTo: state.dueDateTo,
                      onSelect: (dates) {
                        Navigator.pop(context);
                        ref
                            .read(taskListProvider.notifier)
                            .updateDueDateFilter(dates['from'], dates['to']);
                      },
                    ),
                  );
                  break;
                case 'clear':
                  ref.read(taskListProvider.notifier).clearFilters();
                  break;
              }
            },
            itemBuilder: (context) => [
              PopupMenuItem(
                value: 'status',
                child: Row(
                  children: [
                    const Icon(Icons.info_outlined, size: 20),
                    const SizedBox(width: 12),
                    Text(l10n.filterByStatus),
                    if (state.selectedStatus != null) const Spacer(),
                    if (state.selectedStatus != null)
                      Icon(
                        Icons.check,
                        size: 16,
                        color: theme.colorScheme.primary,
                      ),
                  ],
                ),
              ),
              PopupMenuItem(
                value: 'priority',
                child: Row(
                  children: [
                    const Icon(Icons.flag_outlined, size: 20),
                    const SizedBox(width: 12),
                    Text(l10n.filterByPriority),
                    if (state.selectedPriority != null) const Spacer(),
                    if (state.selectedPriority != null)
                      Icon(
                        Icons.check,
                        size: 16,
                        color: theme.colorScheme.primary,
                      ),
                  ],
                ),
              ),
              PopupMenuItem(
                value: 'type',
                child: Row(
                  children: [
                    const Icon(Icons.label_outlined, size: 20),
                    const SizedBox(width: 12),
                    Text(l10n.filterByType),
                    if (state.selectedType != null) const Spacer(),
                    if (state.selectedType != null)
                      Icon(
                        Icons.check,
                        size: 16,
                        color: theme.colorScheme.primary,
                      ),
                  ],
                ),
              ),
              PopupMenuItem(
                value: 'due_date',
                child: Row(
                  children: [
                    const Icon(Icons.calendar_today_outlined, size: 20),
                    const SizedBox(width: 12),
                    Text(l10n.filterByDueDate),
                    if (state.dueDateFrom != null || state.dueDateTo != null)
                      const Spacer(),
                    if (state.dueDateFrom != null || state.dueDateTo != null)
                      Icon(
                        Icons.check,
                        size: 16,
                        color: theme.colorScheme.primary,
                      ),
                  ],
                ),
              ),
              if (hasActiveFilters) ...[
                const PopupMenuDivider(),
                PopupMenuItem(
                  value: 'clear',
                  child: Row(
                    children: [
                      Icon(
                        Icons.clear_all,
                        size: 20,
                        color: theme.colorScheme.error,
                      ),
                      const SizedBox(width: 12),
                      Text(
                        l10n.clearFilters,
                        style: TextStyle(color: theme.colorScheme.error),
                      ),
                    ],
                  ),
                ),
              ],
            ],
          ),
        ],
      ),
      body: body,
    );
  }

  Widget _buildContent(
    BuildContext context,
    TaskListState state,
    ThemeData theme,
  ) {
    final l10n = AppLocalizations.of(context)!;

    if (state.isLoading && state.tasks.isEmpty) {
      return ListView.builder(
        itemCount: 5,
        padding: const EdgeInsets.all(16),
        itemBuilder: (context, index) {
          return const SkeletonListItem(height: 100);
        },
      );
    }

    if (state.errorMessage != null) {
      return ErrorStateWidget(
        message: state.errorMessage!,
        onRetry: () => ref.read(taskListProvider.notifier).refresh(),
      );
    }

    if (state.tasks.isEmpty) {
      // Show different message for search vs no search
      if (state.searchQuery.isNotEmpty) {
        return EmptyStateWidget(
          message: '${l10n.noSearchResults} "${state.searchQuery}"',
          subtitle: l10n.tapToCreateTask,
          icon: Icons.search_off,
        );
      } else {
        return EmptyStateWidget(
          message: l10n.noTasksFound,
          subtitle: l10n.tapToCreateTask,
          icon: Icons.check_box,
        );
      }
    }

    return ListView.builder(
      // controller removed to allow NestedScrollView scroll
      padding: const EdgeInsets.fromLTRB(
        20,
        16,
        20,
        100,
      ), // Match Dashboard padding
      itemCount: state.tasks.length + (state.isLoadingMore ? 1 : 0),
      itemBuilder: (context, index) {
        if (index == state.tasks.length) {
          return const Padding(
            padding: EdgeInsets.all(16.0),
            child: LoadingWidget(size: 24),
          );
        }
        final task = state.tasks[index];
        return TaskCard(
          task: task,
          onTap: () async {
            await Navigator.of(
              context,
            ).pushNamed('${AppRoutes.tasks}/${task.id}');
            if (mounted) {
              ref.read(taskListProvider.notifier).refresh();
            }
          },
        );
      },
    );
  }
}
