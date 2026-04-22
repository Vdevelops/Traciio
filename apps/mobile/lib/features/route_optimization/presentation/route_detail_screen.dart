import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/widgets/loading_widget.dart';
import '../../../core/widgets/error_widget.dart' as error_widget;
import '../../../core/l10n/app_localizations.dart';
import '../application/route_optimization_provider.dart';
import 'widgets/route_map_widget.dart';
import 'widgets/route_header.dart';

class RouteDetailScreen extends ConsumerWidget {
  const RouteDetailScreen({super.key, required this.routeId});

  final String routeId;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final routeAsync = ref.watch(routeDetailProvider(routeId));
    final theme = Theme.of(context);
    final l10n = AppLocalizations.of(context)!;

    return Scaffold(
      body: SafeArea(
        child: Column(
          children: [
            // Header
            Padding(
              padding: const EdgeInsets.fromLTRB(16, 16, 16, 0),
              child: RouteHeader(title: l10n.routeDetails),
            ),
            // Content
            Expanded(
              child: routeAsync.when(
                data: (route) => Column(
                  children: [
                    // Map Section
                    Expanded(
                      flex: 2,
                      child: RouteMapWidget(
                        route: route,
                        height: MediaQuery.of(context).size.height * 0.4,
                      ),
                    ),
                    // Content Section
                    Expanded(
                      flex: 3,
                      child: SingleChildScrollView(
                        padding: const EdgeInsets.all(16),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            // Route Name Card
                            Card(
                              elevation: 0.5,
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(12),
                              ),
                              child: Padding(
                                padding: const EdgeInsets.all(16),
                                child: Row(
                                  children: [
                                    Container(
                                      padding: const EdgeInsets.all(12),
                                      decoration: BoxDecoration(
                                        color:
                                            theme.colorScheme.primaryContainer,
                                        borderRadius: BorderRadius.circular(12),
                                      ),
                                      child: Icon(
                                        Icons.route,
                                        color: theme
                                            .colorScheme
                                            .onPrimaryContainer,
                                      ),
                                    ),
                                    const SizedBox(width: 16),
                                    Expanded(
                                      child: Column(
                                        crossAxisAlignment:
                                            CrossAxisAlignment.start,
                                        children: [
                                          Text(
                                            l10n.routeName,
                                            style: theme.textTheme.bodySmall
                                                ?.copyWith(
                                                  color: theme
                                                      .colorScheme
                                                      .onSurface
                                                      .withValues(alpha: 0.6),
                                                ),
                                          ),
                                          const SizedBox(height: 4),
                                          Text(
                                            route.routeName ??
                                                l10n.unnamedRoute,
                                            style: theme.textTheme.titleLarge
                                                ?.copyWith(
                                                  fontWeight: FontWeight.bold,
                                                ),
                                          ),
                                        ],
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                            ),
                            const SizedBox(height: 16),
                            // Route Stats
                            Row(
                              children: [
                                Expanded(
                                  child: _buildInfoCard(
                                    context,
                                    Icons.straighten,
                                    l10n.distance,
                                    route.totalDistanceFormatted ?? 'N/A',
                                    theme,
                                    Colors.blue,
                                  ),
                                ),
                                const SizedBox(width: 12),
                                Expanded(
                                  child: _buildInfoCard(
                                    context,
                                    Icons.access_time,
                                    l10n.duration,
                                    route.totalDurationFormatted ?? 'N/A',
                                    theme,
                                    Colors.green,
                                  ),
                                ),
                              ],
                            ),
                            const SizedBox(height: 12),
                            _buildInfoCard(
                              context,
                              Icons.location_on,
                              l10n.stops,
                              '${route.waypoints.length - 1}',
                              theme,
                              Colors.purple,
                            ),
                            const SizedBox(height: 24),
                            // Waypoints List
                            Row(
                              children: [
                                Icon(
                                  Icons.place,
                                  color: theme.colorScheme.primary,
                                  size: 20,
                                ),
                                const SizedBox(width: 8),
                                Text(
                                  l10n.waypoints,
                                  style: theme.textTheme.titleMedium?.copyWith(
                                    fontWeight: FontWeight.bold,
                                  ),
                                ),
                              ],
                            ),
                            const SizedBox(height: 12),
                            ...route.waypoints.asMap().entries.map((entry) {
                              final index = entry.key;
                              final waypoint = entry.value;
                              final isStart = waypoint.order == 0 || index == 0;
                              return Padding(
                                padding: const EdgeInsets.only(bottom: 12),
                                child: Card(
                                  elevation: 0.5,
                                  shape: RoundedRectangleBorder(
                                    borderRadius: BorderRadius.circular(12),
                                  ),
                                  child: ListTile(
                                    contentPadding: const EdgeInsets.symmetric(
                                      horizontal: 16,
                                      vertical: 8,
                                    ),
                                    leading: Container(
                                      width: 48,
                                      height: 48,
                                      decoration: BoxDecoration(
                                        gradient: LinearGradient(
                                          begin: Alignment.topLeft,
                                          end: Alignment.bottomRight,
                                          colors: isStart
                                              ? [
                                                  Colors.green.shade400,
                                                  Colors.green.shade600,
                                                ]
                                              : [
                                                  theme
                                                      .colorScheme
                                                      .primaryContainer,
                                                  theme.colorScheme.primary,
                                                ],
                                        ),
                                        shape: BoxShape.circle,
                                        boxShadow: [
                                          BoxShadow(
                                            color:
                                                (isStart
                                                        ? Colors.green
                                                        : theme
                                                              .colorScheme
                                                              .primary)
                                                    .withValues(alpha: 0.15),
                                            blurRadius: 3,
                                            offset: const Offset(0, 1),
                                          ),
                                        ],
                                      ),
                                      child: Center(
                                        child: Text(
                                          isStart ? 'S' : '$index',
                                          style: const TextStyle(
                                            color: Colors.white,
                                            fontWeight: FontWeight.bold,
                                            fontSize: 16,
                                          ),
                                        ),
                                      ),
                                    ),
                                    title: Text(
                                      isStart
                                          ? l10n.startLocation
                                          : waypoint.accountName ??
                                                '${l10n.destination} $index',
                                      style: theme.textTheme.titleSmall
                                          ?.copyWith(
                                            fontWeight: FontWeight.bold,
                                          ),
                                    ),
                                    subtitle: Column(
                                      crossAxisAlignment:
                                          CrossAxisAlignment.start,
                                      children: [
                                        const SizedBox(height: 4),
                                        if (waypoint.address != null)
                                          Text(
                                            waypoint.address!,
                                            style: theme.textTheme.bodySmall,
                                          )
                                        else
                                          Text(
                                            '${waypoint.lat.toStringAsFixed(6)}, ${waypoint.lng.toStringAsFixed(6)}',
                                            style: theme.textTheme.bodySmall
                                                ?.copyWith(
                                                  color: theme
                                                      .colorScheme
                                                      .onSurface
                                                      .withValues(alpha: 0.6),
                                                ),
                                          ),
                                        if (waypoint.accountName != null &&
                                            !isStart) ...[
                                          const SizedBox(height: 4),
                                          Container(
                                            padding: const EdgeInsets.symmetric(
                                              horizontal: 8,
                                              vertical: 4,
                                            ),
                                            decoration: BoxDecoration(
                                              color: theme
                                                  .colorScheme
                                                  .secondaryContainer,
                                              borderRadius:
                                                  BorderRadius.circular(8),
                                            ),
                                            child: Text(
                                              waypoint.accountName!,
                                              style: theme.textTheme.bodySmall
                                                  ?.copyWith(
                                                    color: theme
                                                        .colorScheme
                                                        .onSecondaryContainer,
                                                    fontWeight: FontWeight.w500,
                                                  ),
                                            ),
                                          ),
                                        ],
                                      ],
                                    ),
                                  ),
                                ),
                              );
                            }),
                            // Route Steps (if available)
                            if (route.routeSteps != null &&
                                route.routeSteps!.isNotEmpty) ...[
                              const SizedBox(height: 16),
                              Text(
                                l10n.routeSteps,
                                style: theme.textTheme.titleMedium?.copyWith(
                                  fontWeight: FontWeight.bold,
                                ),
                              ),
                              const SizedBox(height: 12),
                              ...route.routeSteps!.map((step) {
                                return Padding(
                                  padding: const EdgeInsets.only(bottom: 8),
                                  child: Card(
                                    child: ListTile(
                                      leading: CircleAvatar(
                                        backgroundColor:
                                            theme.colorScheme.primary,
                                        child: Text(
                                          '${step.step}',
                                          style: const TextStyle(
                                            color: Colors.white,
                                            fontWeight: FontWeight.bold,
                                          ),
                                        ),
                                      ),
                                      title: Text(
                                        step.instruction ?? l10n.continueRoute,
                                        style: const TextStyle(
                                          fontWeight: FontWeight.w500,
                                        ),
                                      ),
                                      subtitle: Text(
                                        '${step.distanceFormatted ?? step.distance.toStringAsFixed(2)} km • ${step.durationFormatted ?? "${step.duration}s"}',
                                      ),
                                    ),
                                  ),
                                );
                              }),
                            ],
                          ],
                        ),
                      ),
                    ),
                  ],
                ),
                loading: () => const LoadingWidget(),
                error: (error, stack) => error_widget.ErrorStateWidget(
                  message: error.toString(),
                  onRetry: () {
                    ref.invalidate(routeDetailProvider(routeId));
                  },
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildInfoCard(
    BuildContext context,
    IconData icon,
    String label,
    String value,
    ThemeData theme,
    Color iconColor,
  ) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: theme.colorScheme.surfaceContainerHighest,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(
          color: theme.colorScheme.outline.withValues(alpha: 0.1),
        ),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Container(
                padding: const EdgeInsets.all(8),
                decoration: BoxDecoration(
                  color: iconColor.withValues(alpha: 0.1),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Icon(icon, size: 20, color: iconColor),
              ),
              const SizedBox(width: 8),
              Expanded(
                child: Text(
                  label,
                  style: theme.textTheme.bodySmall?.copyWith(
                    color: theme.colorScheme.onSurface.withValues(alpha: 0.6),
                    fontWeight: FontWeight.w500,
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(height: 8),
          Text(
            value,
            style: theme.textTheme.titleMedium?.copyWith(
              fontWeight: FontWeight.bold,
            ),
          ),
        ],
      ),
    );
  }
}
