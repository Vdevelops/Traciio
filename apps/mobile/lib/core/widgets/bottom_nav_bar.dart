import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../l10n/app_localizations.dart';
import '../permissions/permission_checker.dart';
import '../permissions/permission_provider.dart';

enum NavItem { home, accounts, reports, routeOptimization, profile }

class BottomNavBar extends ConsumerWidget {
  const BottomNavBar({
    super.key,
    required this.currentIndex,
    required this.onTap,
  });

  final int currentIndex;
  final ValueChanged<int> onTap;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    final l10n = AppLocalizations.of(context);

    // Get permissions from auth state
    final permissions = ref.watch(userPermissionsProvider);

    // Define menu items with their resource names
    final menuItems = [
      _MenuItemData(
        index: 0,
        icon: Icons
            .dashboard_outlined, // Changed to dashboard/grid icon per image
        activeIcon: Icons.dashboard,
        label: l10n?.home ?? 'Home',
        route: '/dashboard',
        resource: 'dashboard',
      ),
      _MenuItemData(
        index: 1,
        icon: Icons.timeline, // Pipeline icon (analytics/graph)
        activeIcon: Icons.timeline_sharp,
        label: l10n?.pipeline ?? 'Pipeline',
        route: '/pipeline',
        resource: 'pipeline',
      ),
      _MenuItemData(
        index: 2,
        icon: Icons.bar_chart_outlined, // Leads icon (bar chart)
        activeIcon: Icons.bar_chart,
        label: l10n?.leads ?? 'Leads',
        route: '/leads',
        resource: 'leads',
      ),
      _MenuItemData(
        index: 3,
        icon: Icons.map_outlined, // Route icon (map)
        activeIcon: Icons.map,
        label: l10n?.routeMenu ?? 'Routes',
        route: '/route-optimization',
        resource: 'route-optimization',
      ),
      _MenuItemData(
        index: 4,
        icon: Icons.calendar_today_outlined, // Schedule icon (calendar)
        activeIcon: Icons.calendar_today,
        label: l10n?.schedules ?? 'Schedule',
        route: '/schedules',
        resource: 'schedules',
      ),
    ];

    // Filter menu items based on permissions
    final visibleItems = menuItems.where((item) {
      // Check if user has view permission for this resource
      return PermissionChecker.canView(permissions, item.resource);
    }).toList();

    // Build navbar with filtered items (no loading state needed - permissions available instantly)
    return _buildNavBar(context, theme, visibleItems, currentIndex, onTap);
  }

  Widget _buildNavBar(
    BuildContext context,
    ThemeData theme,
    List<_MenuItemData> items,
    int currentIndex,
    ValueChanged<int> onTap,
  ) {
    return Container(
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        border: Border(
          top: BorderSide(
            color: theme.colorScheme.onSurface.withValues(alpha: 0.1),
            width: 1,
          ),
        ),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.05),
            blurRadius: 3,
            offset: const Offset(0, -1),
          ),
        ],
      ),
      child: SafeArea(
        top: false,
        bottom: true,
        child: Container(
          height: 75,
          padding: const EdgeInsets.symmetric(vertical: 10, horizontal: 8),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceAround,
            children: items.map((item) {
              final visibleIndex = items.indexOf(item);
              final isActive = items[visibleIndex].index == currentIndex;

              return Expanded(
                child: _NavItem(
                  icon: item.icon,
                  activeIcon: item.activeIcon,
                  label: item.label,
                  isActive: isActive,
                  onTap: () => onTap(item.index),
                ),
              );
            }).toList(),
          ),
        ),
      ),
    );
  }
}

class _MenuItemData {
  const _MenuItemData({
    required this.index,
    required this.icon,
    required this.activeIcon,
    required this.label,
    required this.route,
    required this.resource,
  });

  final int index;
  final IconData icon;
  final IconData activeIcon;
  final String label;
  final String route;
  final String
  resource; // Resource name for permission checking (e.g., 'dashboard', 'tasks')
}

class _NavItem extends StatefulWidget {
  const _NavItem({
    required this.icon,
    required this.activeIcon,
    required this.label,
    required this.isActive,
    required this.onTap,
  });

  final IconData icon;
  final IconData activeIcon;
  final String label;
  final bool isActive;
  final VoidCallback onTap;

  @override
  State<_NavItem> createState() => _NavItemState();
}

class _NavItemState extends State<_NavItem>
    with SingleTickerProviderStateMixin {
  late AnimationController _controller;
  late Animation<double> _scaleAnimation;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      duration: const Duration(milliseconds: 200),
      vsync: this,
    );
    _scaleAnimation = Tween<double>(
      begin: 1.0,
      end: 0.95,
    ).animate(CurvedAnimation(parent: _controller, curve: Curves.easeInOut));
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  void _handleTap() {
    _controller.forward().then((_) {
      _controller.reverse();
    });
    widget.onTap();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return GestureDetector(
      onTap: _handleTap,
      child: ScaleTransition(
        scale: _scaleAnimation,
        child: TweenAnimationBuilder<double>(
          duration: const Duration(milliseconds: 400),
          curve: Curves.elasticOut,
          tween: Tween<double>(
            begin: widget.isActive ? 1.0 : 0.95,
            end: widget.isActive ? 1.05 : 1.0,
          ),
          builder: (context, scale, child) {
            return Transform.scale(
              scale: scale,
              child: AnimatedContainer(
                duration: const Duration(milliseconds: 350),
                curve: Curves.easeOutCubic,
                margin: const EdgeInsets.symmetric(horizontal: 4),
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 6),
                decoration: BoxDecoration(
                  color: widget.isActive
                      ? theme.colorScheme.primary.withValues(alpha: 0.15)
                      : Colors.transparent,
                  borderRadius: BorderRadius.circular(20),
                ),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    AnimatedSwitcher(
                      duration: const Duration(milliseconds: 300),
                      transitionBuilder: (child, animation) {
                        return ScaleTransition(
                          scale: animation,
                          child: RotationTransition(
                            turns: Tween<double>(
                              begin: 0.8,
                              end: 1.0,
                            ).animate(animation),
                            child: FadeTransition(
                              opacity: animation,
                              child: child,
                            ),
                          ),
                        );
                      },
                      child: Icon(
                        widget.isActive ? widget.activeIcon : widget.icon,
                        key: ValueKey(widget.isActive),
                        color: widget.isActive
                            ? theme.colorScheme.primary
                            : theme.colorScheme.onSurface.withValues(
                                alpha: 0.6,
                              ),
                        size: widget.isActive ? 24 : 22,
                      ),
                    ),
                    const SizedBox(height: 3),
                    AnimatedDefaultTextStyle(
                      duration: const Duration(milliseconds: 300),
                      curve: Curves.easeOutCubic,
                      style: TextStyle(
                        fontSize: widget.isActive ? 10.5 : 9.5,
                        fontWeight: widget.isActive
                            ? FontWeight.w600
                            : FontWeight.normal,
                        color: widget.isActive
                            ? theme.colorScheme.primary
                            : theme.colorScheme.onSurface.withValues(
                                alpha: 0.6,
                              ),
                        height: 1.2,
                      ),
                      child: Text(
                        widget.label,
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                        textAlign: TextAlign.center,
                      ),
                    ),
                  ],
                ),
              ),
            );
          },
        ),
      ),
    );
  }
}
