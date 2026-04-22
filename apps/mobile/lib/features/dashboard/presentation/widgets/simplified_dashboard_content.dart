import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../../core/l10n/app_localizations.dart';
import '../../../../core/widgets/loading_widget.dart';
import '../../application/dashboard_provider.dart';
import 'upcoming_tasks_section.dart';
import 'target_progress_card.dart';
import 'quick_navigation_widget.dart';
import 'dashboard_greeting_header.dart';

/// Simplified dashboard content untuk sales rep
/// Focus pada quick access ke fitur utama (Visit Reports, Tasks)
class SimplifiedDashboardContent extends ConsumerWidget {
  const SimplifiedDashboardContent({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final l10n = AppLocalizations.of(context)!;
    final dashboardState = ref.watch(dashboardProvider);

    // Load dashboard data if not already loading
    if (!dashboardState.isLoading &&
        dashboardState.overview == null &&
        dashboardState.errorMessage == null) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        if (context.mounted) {
          ref.read(dashboardProvider.notifier).loadDashboard();
        }
      });
    }

    // Note: This widget is now used inside SingleChildScrollView, so we use Column instead of ListView
    // or ListView with shrinkWrap to avoid nested scroll conflicts
    return Padding(
      padding: const EdgeInsets.fromLTRB(
        20,
        0,
        20,
        80,
      ), // Bottom padding untuk floating navbar (height 70 + minimal spacing)
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          // Greeting Header dengan nama user dan subtitle harian
          const DashboardGreetingHeader(),

          // Quick Stats Section
          if (dashboardState.isLoadingOverview && dashboardState.overview == null) ...[
            const SizedBox(height: 200, child: LoadingWidget()),
          ] else if (dashboardState.errorMessage != null) ...[
            // Error State
            Container(
              padding: const EdgeInsets.all(24),
              child: Column(
                children: [
                  Icon(
                    Icons.error_outline,
                    size: 48,
                    color: theme.colorScheme.error,
                  ),
                  const SizedBox(height: 16),
                  Text(
                    l10n.errorLoadingDashboard,
                    style: theme.textTheme.titleMedium?.copyWith(
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    dashboardState.errorMessage ?? l10n.unknownError,
                    style: theme.textTheme.bodyMedium,
                    textAlign: TextAlign.center,
                  ),
                  const SizedBox(height: 16),
                  FilledButton(
                    onPressed: () {
                      ref.read(dashboardProvider.notifier).refresh();
                    },
                    child: Text(l10n.retry),
                  ),
                ],
              ),
            ),
          ] else if (dashboardState.overview != null) ...[
            // Performance Goal Card (Orange theme)
            if (dashboardState.overview!.target.targetAmount > 0)
              TargetProgressCard(data: dashboardState.overview!.target),

            const SizedBox(height: 20),

            // Quick Navigation Widget (Small buttons untuk accounts, contacts, visits)
            const QuickNavigationWidget(),

            const SizedBox(height: 24),

            // Upcoming Tasks Section
            const UpcomingTasksSection(),

            const SizedBox(height: 24),
          ] else ...[
            // Empty State - No data available
            Container(
              padding: const EdgeInsets.all(24),
              child: Column(
                children: [
                  Icon(
                    Icons.dashboard_outlined,
                    size: 64,
                    color: theme.colorScheme.onSurface.withValues(alpha: 0.3),
                  ),
                  const SizedBox(height: 16),
                  Text(
                    l10n.noDataAvailable,
                    style: theme.textTheme.titleMedium?.copyWith(
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    l10n.pullDownToRefresh,
                    style: theme.textTheme.bodyMedium?.copyWith(
                      color: theme.colorScheme.onSurface.withValues(alpha: 0.6),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ],
      ),
    );
  }
}
