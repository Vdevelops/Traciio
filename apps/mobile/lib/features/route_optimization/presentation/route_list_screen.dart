import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/routing/app_router.dart';
import '../../../core/l10n/app_localizations.dart';
import '../../../core/widgets/skeleton_widget.dart';
import '../../../core/widgets/error_widget.dart' as error_widget;
import '../../../core/widgets/main_scaffold.dart';
import '../../../core/widgets/sync_status_indicator.dart';
import '../../../core/permissions/permission_provider.dart';
import '../application/route_optimization_provider.dart';
import '../application/route_optimization_state.dart';
import 'route_form_screen.dart';
import 'widgets/route_card.dart';

class RouteListScreen extends ConsumerStatefulWidget {
  const RouteListScreen({super.key});

  @override
  ConsumerState<RouteListScreen> createState() => _RouteListScreenState();
}

class _RouteListScreenState extends ConsumerState<RouteListScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(routeListProvider.notifier).loadRoutes();
    });
  }

  bool _onScrollNotification(ScrollNotification notification) {
    if (notification is ScrollUpdateNotification) {
      if (notification.metrics.pixels >=
          notification.metrics.maxScrollExtent * 0.8) {
        ref.read(routeListProvider.notifier).loadMore();
      }
    }
    return false;
  }

  Future<void> _onRefresh() async {
    await ref.read(routeListProvider.notifier).refresh();
  }

  void _navigateToCreateRoute() {
    Navigator.push(
      context,
      MaterialPageRoute(builder: (context) => const RouteFormScreen()),
    ).then((_) {
      // Refresh list after creating
      ref.read(routeListProvider.notifier).refresh();
    });
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(routeListProvider);
    final theme = Theme.of(context);
    final l10n = AppLocalizations.of(context)!;

    // Check if user has permission to create route
    final hasCreatePermission = ref.watch(
      canCreateProvider('route-optimization'),
    );

    return MainScaffold(
      currentIndex: 3,
      title: l10n.routeOptimization,
      actions: const [
        // Sync Status Indicator
        Padding(
          padding: EdgeInsets.only(right: 16),
          child: Center(
            child: SyncStatusIndicator(
              featureKey: 'route_optimization',
              iconSize: 20,
            ),
          ),
        ),
      ],
      floatingActionButton: hasCreatePermission
          ? FloatingActionButton(
              elevation: 0.5,
              onPressed: _navigateToCreateRoute,
              child: const Icon(Icons.add),
            )
          : null,
      body: NotificationListener<ScrollNotification>(
        onNotification: _onScrollNotification,
        child: RefreshIndicator(
          onRefresh: _onRefresh,
          child: _buildContent(context, state, theme, l10n),
        ),
      ),
    );
  }

  Widget _buildContent(
    BuildContext context,
    RouteListState state,
    ThemeData theme,
    AppLocalizations l10n,
  ) {
    if (state.isLoading && state.routes.isEmpty) {
      return ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: 5,
        itemBuilder: (context, index) => const Padding(
          padding: EdgeInsets.only(bottom: 12),
          child: SkeletonWidget(height: 120),
        ),
      );
    }

    if (state.errorMessage != null && state.routes.isEmpty) {
      return error_widget.ErrorStateWidget(
        message: state.errorMessage!,
        onRetry: () {
          ref.read(routeListProvider.notifier).refresh();
        },
      );
    }

    if (state.routes.isEmpty) {
      final hasCreatePermission = ref.watch(
        canCreateProvider('route-optimization'),
      );

      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              Icons.route,
              size: 64,
              color: theme.colorScheme.onSurface.withValues(alpha: 0.3),
            ),
            const SizedBox(height: 16),
            Text(
              l10n.noRoutesFound,
              style: theme.textTheme.titleMedium?.copyWith(
                color: theme.colorScheme.onSurface.withValues(alpha: 0.6),
              ),
            ),
            const SizedBox(height: 8),
            Text(
              l10n.createFirstOptimizedRoute,
              style: theme.textTheme.bodyMedium?.copyWith(
                color: theme.colorScheme.onSurface.withValues(alpha: 0.5),
              ),
            ),
            if (hasCreatePermission) ...[
              const SizedBox(height: 24),
              ElevatedButton.icon(
                onPressed: _navigateToCreateRoute,
                icon: const Icon(Icons.add),
                label: Text(l10n.createRoute),
                style: ElevatedButton.styleFrom(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 24,
                    vertical: 12,
                  ),
                ),
              ),
            ],
          ],
        ),
      );
    }

    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: state.routes.length + (state.isLoadingMore ? 1 : 0),
      itemBuilder: (context, index) {
        if (index == state.routes.length) {
          return const Padding(
            padding: EdgeInsets.all(16),
            child: Center(child: CircularProgressIndicator()),
          );
        }

        final route = state.routes[index];
        return Padding(
          padding: const EdgeInsets.only(bottom: 12),
          child: RouteCard(
            route: route,
            onTap: () {
              Navigator.pushNamed(
                context,
                '${AppRoutes.routeOptimization}/${route.id}',
              );
            },
          ),
        );
      },
    );
  }
}
