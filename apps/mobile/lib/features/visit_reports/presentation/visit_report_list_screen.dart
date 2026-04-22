import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'dart:async';

import '../application/visit_report_provider.dart';
import '../application/visit_report_state.dart';
import '../../../core/routing/app_router.dart';
import '../../../core/l10n/app_localizations.dart';
import '../../../core/widgets/error_widget.dart';
import '../../../core/widgets/loading_widget.dart';
import '../../../core/widgets/skeleton_widget.dart';
import '../../../core/widgets/sync_status_indicator.dart';
import 'widgets/visit_report_card.dart';

class VisitReportListScreen extends ConsumerStatefulWidget {
  const VisitReportListScreen({
    super.key,
    this.hideAppBar = false,
    this.searchController,
  });

  final bool hideAppBar;
  final TextEditingController? searchController;

  @override
  ConsumerState<VisitReportListScreen> createState() =>
      _VisitReportListScreenState();
}

class _VisitReportListScreenState extends ConsumerState<VisitReportListScreen> {
  // Removed ScrollController for NestedScrollView support

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(visitReportListProvider.notifier).loadVisitReports();
    });
  }

  @override
  void dispose() {
    super.dispose();
  }

  bool _onScrollNotification(ScrollNotification notification) {
    if (notification is ScrollUpdateNotification) {
      if (notification.metrics.pixels >=
          notification.metrics.maxScrollExtent * 0.8) {
        ref.read(visitReportListProvider.notifier).loadMore();
      }
    }
    return false;
  }

  Future<void> _onRefresh() async {
    await ref.read(visitReportListProvider.notifier).refresh();
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(visitReportListProvider);
    final theme = Theme.of(context);
    final l10n = AppLocalizations.of(context)!;

    final body = NotificationListener<ScrollNotification>(
      onNotification: _onScrollNotification,
      child: RefreshIndicator(
        onRefresh: _onRefresh,
        child: _buildContent(context, state, theme),
      ),
    );

    if (widget.hideAppBar) {
      return body;
    }

    return Scaffold(
      appBar: AppBar(
        title: Text(l10n.visitReports),
        elevation: 0,
        actions: [
          // Sync Status Indicator
          const Padding(
            padding: EdgeInsets.only(right: 8),
            child: Center(
              child: SyncStatusIndicator(
                featureKey: 'visit_reports',
                iconSize: 20,
              ),
            ),
          ),
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () {
              Navigator.pushNamed(context, '${AppRoutes.visitReports}/create');
            },
            tooltip: l10n.createVisitReport,
          ),
        ],
      ),
      body: body,
    );
  }

  Widget _buildContent(
    BuildContext context,
    VisitReportListState state,
    ThemeData theme,
  ) {
    final l10n = AppLocalizations.of(context)!;

    if (state.isLoading && state.visitReports.isEmpty) {
      return const LoadingWidget();
    }

    if (state.errorMessage != null && state.visitReports.isEmpty) {
      return ErrorStateWidget(
        message: state.errorMessage!,
        onRetry: () {
          ref.read(visitReportListProvider.notifier).refresh();
        },
      );
    }

    if (state.visitReports.isEmpty) {
      return EmptyStateWidget(
        message: l10n.noVisitReportsFound,
        subtitle: l10n.tapToCreateVisitReport,
        icon: Icons.assignment_outlined,
      );
    }

    // Show skeleton screens if loading first page
    if (state.isLoading && state.visitReports.isEmpty) {
      return ListView.builder(
        itemCount: 5, // Show 5 skeleton items
        itemBuilder: (context, index) {
          return const SkeletonListItem(height: 100);
        },
      );
    }

    return ListView.builder(
      // Controller removed for NestedScrollView support
      padding: const EdgeInsets.fromLTRB(20, 16, 20, 100),
      itemCount: state.visitReports.length + (state.isLoadingMore ? 1 : 0),
      // Optimize: Use addAutomaticKeepAlives and addRepaintBoundaries for better performance
      addAutomaticKeepAlives: false,
      addRepaintBoundaries: true,
      cacheExtent: 500, // Cache more items for smoother scrolling
      itemBuilder: (context, index) {
        if (index == state.visitReports.length) {
          return const Padding(
            padding: EdgeInsets.all(16),
            child: LoadingWidget(size: 24),
          );
        }

        final visitReport = state.visitReports[index];
        return VisitReportCard(
          key: ValueKey(
            visitReport.id,
          ), // Use stable key for better performance
          visitReport: visitReport,
          onTap: () {
            Navigator.pushNamed(
              context,
              '${AppRoutes.visitReports}/${visitReport.id}',
              arguments: visitReport.id,
            );
          },
        );
      },
    );
  }
}
