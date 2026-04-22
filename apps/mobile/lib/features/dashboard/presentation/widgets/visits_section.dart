import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../application/dashboard_provider.dart';
import '../../../../core/routing/app_router.dart';
import '../../../../core/widgets/auth_gate.dart';
import '../../../../core/permissions/permission_checker.dart';
import '../../../../core/permissions/permission_provider.dart';
import '../../../../core/l10n/app_localizations.dart';
import '../../../visit_reports/presentation/reports_screen.dart';
// import '../../data/models/dashboard.dart';
import 'visit_card.dart';

class VisitsSection extends ConsumerWidget {
  const VisitsSection({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final l10n = AppLocalizations.of(context)!;
    final state = ref.watch(dashboardProvider);
    final notifier = ref.read(dashboardProvider.notifier);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.symmetric(horizontal: 4),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Row(
                children: [
                  Container(
                    padding: const EdgeInsets.all(8),
                    decoration: BoxDecoration(
                      color: colorScheme.primary.withValues(alpha: 0.15),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Icon(
                      Icons.location_on,
                      size: 20,
                      color: colorScheme.primary,
                    ),
                  ),
                  const SizedBox(width: 12),
                  Text(
                    'Visits', // Tetap "Visits" di semua bahasa
                    style: theme.textTheme.titleLarge?.copyWith(
                      fontWeight: FontWeight.bold,
                      color: colorScheme.onSurface,
                    ),
                  ),
                ],
              ),
              TextButton(
                onPressed: () {
                  _navigateToVisitsMenu(context, ref);
                },
                style: TextButton.styleFrom(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 12,
                    vertical: 8,
                  ),
                ),
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Text(
                      l10n.seeAll,
                      style: theme.textTheme.labelLarge?.copyWith(
                        color: colorScheme.primary,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(width: 4),
                    Icon(
                      Icons.arrow_forward_ios,
                      size: 14,
                      color: colorScheme.primary,
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
        const SizedBox(height: 16),
        // Tabs
        SingleChildScrollView(
          scrollDirection: Axis.horizontal,
          padding: const EdgeInsets.symmetric(horizontal: 4),
          child: Row(
            children: [
              _buildTab(
                context,
                l10n.draft,
                state.visitStatusFilter == 'draft',
                () => notifier.setVisitStatusFilter('draft'),
              ),
              const SizedBox(width: 8),
              _buildTab(
                context,
                l10n.submitted,
                state.visitStatusFilter == 'submitted',
                () => notifier.setVisitStatusFilter('submitted'),
              ),
              const SizedBox(width: 8),
              _buildTab(
                context,
                l10n.approved,
                state.visitStatusFilter == 'approved',
                () => notifier.setVisitStatusFilter('approved'),
              ),
              const SizedBox(width: 8),
              _buildTab(
                context,
                l10n.rejected,
                state.visitStatusFilter == 'rejected',
                () => notifier.setVisitStatusFilter('rejected'),
              ),
            ],
          ),
        ),
        const SizedBox(height: 16),
        // List
        SizedBox(
          height: 180,
          child: state.isLoadingVisits
              ? const Center(child: CircularProgressIndicator())
              : (state.visits == null || state.visits!.isEmpty)
              ? Container(
                  margin: const EdgeInsets.symmetric(horizontal: 4),
                  padding: const EdgeInsets.all(24),
                  decoration: BoxDecoration(
                    color: theme.colorScheme.surfaceContainerHighest
                        .withValues(alpha: 0.3),
                    borderRadius: BorderRadius.circular(16),
                  ),
                  child: Center(
                    child: Column(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Icon(
                          Icons.location_off,
                          size: 48,
                          color: colorScheme.onSurface.withValues(alpha: 0.3),
                        ),
                        const SizedBox(height: 12),
                        Text(
                          l10n.noVisitsFound,
                          style: theme.textTheme.bodyMedium?.copyWith(
                            color: colorScheme.onSurface.withValues(alpha: 0.6),
                          ),
                        ),
                      ],
                    ),
                  ),
                )
              : ListView.separated(
                  padding: const EdgeInsets.symmetric(horizontal: 4),
                  scrollDirection: Axis.horizontal,
                  itemCount: state.visits!.length.clamp(0, 10), // Limit to 10 items
                  separatorBuilder: (_, _) => const SizedBox(width: 12),
                  // Optimize: Use addRepaintBoundaries and cacheExtent
                  addRepaintBoundaries: true,
                  cacheExtent: 500,
                  itemBuilder: (context, index) {
                    return VisitCard(
                      key: ValueKey(state.visits![index].id),
                      visit: state.visits![index],
                    );
                  },
                ),
        ),
      ],
    );
  }

  void _navigateToVisitsMenu(BuildContext context, WidgetRef ref) {
    // Helper function to create a route without animation (same as MainScaffold._handleNavTap)
    Route<T> createNoAnimationRoute<T>(Widget page, String routeName) {
      return PageRouteBuilder<T>(
        settings: RouteSettings(name: routeName),
        pageBuilder: (context, animation, secondaryAnimation) => page,
        transitionDuration: Duration.zero,
        reverseTransitionDuration: Duration.zero,
      );
    }

    // Get permissions to check access
    final permissions = ref.read(userPermissionsProvider);
    final resource = 'visit-reports';

    // Check permission
    final hasPermission = PermissionChecker.canView(permissions, resource);
    if (!hasPermission) {
      // Show error message if no permission
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('You do not have permission to access this page'),
          duration: Duration(seconds: 3),
        ),
      );
      return;
    }

    // Navigate to ReportsScreen (visits menu) with MainScaffold currentIndex: 2
    // ReportsScreen already wraps itself with MainScaffold with currentIndex: 2
    Navigator.pushReplacement(
      context,
      createNoAnimationRoute(
        const AuthGate(
          requiredRoute: AppRoutes.visitReports,
          child: ReportsScreen(),
        ),
        AppRoutes.visitReports,
      ),
    );
  }

  Widget _buildTab(
    BuildContext context,
    String label,
    bool isSelected,
    VoidCallback onTap,
  ) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    
    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(16),
        child: Container(
          padding: const EdgeInsets.symmetric(horizontal: 18, vertical: 10),
          decoration: BoxDecoration(
            color: isSelected
                ? colorScheme.primary
                : colorScheme.surfaceContainerHighest.withValues(alpha: 0.5),
            borderRadius: BorderRadius.circular(16),
            border: isSelected
                ? null
                : Border.all(
                    color: colorScheme.outline.withValues(alpha: 0.2),
                    width: 1,
                  ),
          ),
          child: Text(
            label,
            style: theme.textTheme.labelLarge?.copyWith(
              color: isSelected
                  ? colorScheme.onPrimary
                  : colorScheme.onSurface.withValues(alpha: 0.7),
              fontWeight: isSelected ? FontWeight.bold : FontWeight.w500,
            ),
          ),
        ),
      ),
    );
  }
}
