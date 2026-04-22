import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/l10n/app_localizations.dart';
import '../../../../core/routing/app_router.dart';
import '../../../../core/widgets/error_widget.dart';
import '../../../../core/widgets/loading_widget.dart';
import '../../../../core/widgets/main_scaffold.dart';
import '../../../../core/widgets/search_modal.dart';
import '../../../../core/widgets/sync_status_indicator.dart';
import '../../application/schedule_provider.dart';
import '../../application/schedule_state.dart';
import '../widgets/schedule_card.dart';
import '../widgets/schedule_filter_sheet.dart';

class ScheduleListScreen extends ConsumerStatefulWidget {
  const ScheduleListScreen({super.key});

  @override
  ConsumerState<ScheduleListScreen> createState() => _ScheduleListScreenState();
}

class _ScheduleListScreenState extends ConsumerState<ScheduleListScreen> {
  // The _searchController is no longer needed for the inline TextField,
  // but if it's used elsewhere for state management related to search,
  // it might still be relevant. Given the new search modal approach,
  // it's likely not needed for UI control anymore.
  // However, the instruction only asks to remove the TextField, not the controller.
  // For now, I'll keep the controller as the instruction doesn't explicitly remove it,
  // but it will become unused.
  final TextEditingController _searchController = TextEditingController();

  @override
  void initState() {
    super.initState();
    Future.microtask(() {
      ref.read(scheduleListProvider.notifier).loadSchedules();
    });
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  void _showStatusFilter() {
    if (!mounted) return;
    final state = ref.read(scheduleListProvider);
    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) {
        return ScheduleStatusFilterSheet(
          selectedStatus: state.selectedStatus,
          onSelect: (status) {
            if (!mounted) return;
            ref.read(scheduleListProvider.notifier).setStatus(status);
            Navigator.pop(context);
          },
        );
      },
    );
  }

  void _showDateFilter() {
    if (!mounted) return;
    final state = ref.read(scheduleListProvider);
    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) {
        return ScheduleDateFilterSheet(
          scheduledFrom: state.scheduledFrom,
          scheduledTo: state.scheduledTo,
          onSelect: (dates) {
            if (!mounted) return;
            ref
                .read(scheduleListProvider.notifier)
                .setDateRange(dates['from'], dates['to']);
          },
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(scheduleListProvider);
    final l10n = AppLocalizations.of(context)!;
    final theme = Theme.of(context);

    return MainScaffold(
      currentIndex: 4, // Dashboard=0, Pipeline=1, Leads=2, Route=3, Schedules=4
      title: l10n.schedules,
      actions: [
        // Sync Status Indicator
        const Padding(
          padding: EdgeInsets.only(right: 8),
          child: Center(
            child: SyncStatusIndicator(featureKey: 'schedules', iconSize: 20),
          ),
        ),
        IconButton(
          icon: const Icon(Icons.search),
          onPressed: () {
            if (!mounted) return;
            showSearchModal(
              context,
              hintText: l10n.searchSchedules,
              initialQuery: ref.read(scheduleListProvider).search,
              onSearch: (query) {
                if (!mounted) return;
                ref.read(scheduleListProvider.notifier).setSearch(query);
              },
              onClear: () {
                if (!mounted) return;
                ref.read(scheduleListProvider.notifier).setSearch('');
              },
            );
          },
        ),
        IconButton(
          icon: Icon(
            Icons.filter_list,
            color: state.selectedStatus != null
                ? theme.colorScheme.primary
                : null,
          ),
          onPressed: _showStatusFilter,
        ),
        IconButton(
          icon: Icon(
            Icons.calendar_today,
            color: state.scheduledFrom != null || state.scheduledTo != null
                ? theme.colorScheme.primary
                : null,
          ),
          onPressed: _showDateFilter,
        ),
      ],
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          Navigator.pushNamed(context, AppRoutes.scheduleForm);
        },
        child: const Icon(Icons.add),
      ),
      body: Column(
        children: [
          if (state.hasActiveFilters)
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 16.0),
              child: Row(
                children: [
                  Text(
                    l10n.filter,
                    style: theme.textTheme.labelMedium?.copyWith(
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(width: 8),
                  Expanded(
                    child: SingleChildScrollView(
                      scrollDirection: Axis.horizontal,
                      child: Row(
                        children: [
                          if (state.selectedStatus != null)
                            Padding(
                              padding: const EdgeInsets.only(right: 8.0),
                              child: Chip(
                                label: Text(state.selectedStatus!),
                                onDeleted: () {
                                  if (!mounted) return;
                                  ref
                                      .read(scheduleListProvider.notifier)
                                      .setStatus(null);
                                },
                              ),
                            ),
                          if (state.scheduledFrom != null ||
                              state.scheduledTo != null)
                            Padding(
                              padding: const EdgeInsets.only(right: 8.0),
                              child: Chip(
                                label: Text(l10n.visitDate),
                                onDeleted: () {
                                  if (!mounted) return;
                                  ref
                                      .read(scheduleListProvider.notifier)
                                      .setDateRange(null, null);
                                },
                              ),
                            ),
                          TextButton(
                            onPressed: () {
                              if (!mounted) return;
                              _searchController.clear();
                              ref
                                  .read(scheduleListProvider.notifier)
                                  .clearFilters();
                            },
                            child: Text(l10n.clearFilters),
                          ),
                        ],
                      ),
                    ),
                  ),
                ],
              ),
            ),
          Expanded(
            child: RefreshIndicator(
              onRefresh: () async {
                if (!mounted) return;
                await ref
                    .read(scheduleListProvider.notifier)
                    .loadSchedules(forceRefresh: true);
              },
              child: _buildList(state, l10n),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildList(ScheduleListState state, AppLocalizations l10n) {
    if (state.isLoading && state.schedules.isEmpty) {
      return const LoadingWidget();
    }

    if (state.errorMessage != null && state.schedules.isEmpty) {
      return ErrorStateWidget(
        message: state.errorMessage!,
        onRetry: () {
          if (!mounted) return;
          ref
              .read(scheduleListProvider.notifier)
              .loadSchedules(forceRefresh: true);
        },
      );
    }

    if (state.schedules.isEmpty) {
      return EmptyStateWidget(
        message: l10n.noSchedulesFound,
        subtitle: l10n.tapToCreateSchedule,
        icon: Icons.event_note,
      );
    }

    return NotificationListener<ScrollNotification>(
      onNotification: (notification) {
        if (notification is ScrollEndNotification &&
            notification.metrics.extentAfter < 200) {
          if (mounted) {
            ref.read(scheduleListProvider.notifier).loadMore();
          }
        }
        return false;
      },
      child: ListView.builder(
        padding: const EdgeInsets.fromLTRB(16, 16, 16, 100),
        itemCount: state.schedules.length + (state.hasNextPage ? 1 : 0),
        itemBuilder: (context, index) {
          if (index == state.schedules.length) {
            return const Center(
              child: Padding(
                padding: EdgeInsets.all(8.0),
                child: CircularProgressIndicator(),
              ),
            );
          }

          final schedule = state.schedules[index];
          return ScheduleCard(
            schedule: schedule,
            onTap: () {
              Navigator.pushNamed(
                context,
                AppRoutes.scheduleDetail,
                arguments: schedule.id,
              );
            },
          );
        },
      ),
    );
  }
}
