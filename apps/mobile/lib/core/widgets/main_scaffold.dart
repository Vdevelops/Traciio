import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../features/dashboard/presentation/dashboard_screen.dart';
import '../../features/leads/presentation/lead_list_screen.dart';
import '../../features/pipeline/presentation/screens/pipeline_screen.dart';
import '../../features/route_optimization/presentation/route_list_screen.dart';
import '../../features/schedules/presentation/screens/schedule_list_screen.dart';
import '../permissions/permission_checker.dart';
import '../permissions/permission_provider.dart';
import '../routing/app_router.dart';
import 'auth_gate.dart';
import 'bottom_nav_bar.dart';
import 'offline_indicator.dart';

/// Custom FloatingActionButtonLocation untuk menghindari overlap dengan bottom navbar
class _CustomFloatingActionButtonLocation extends FloatingActionButtonLocation {
  const _CustomFloatingActionButtonLocation();

  @override
  Offset getOffset(ScaffoldPrelayoutGeometry scaffoldGeometry) {
    // Tinggi bottom navbar + padding untuk FAB
    const bottomNavBarHeight = 75.0;
    const fabPadding = 20.0; // Increased padding untuk lebih aman
    final fabSize = scaffoldGeometry.floatingActionButtonSize;
    
    // Posisi FAB di kanan bawah, dengan offset untuk menghindari navbar
    // Tambahkan extra spacing untuk memastikan tidak menyentuh navbar
    final x = scaffoldGeometry.scaffoldSize.width - fabSize.width - fabPadding;
    final y = scaffoldGeometry.scaffoldSize.height - 
              fabSize.height - 
              bottomNavBarHeight - 
              fabPadding - 
              16.0; // Extra spacing untuk memastikan tidak menyentuh
    
    return Offset(x, y);
  }

  @override
  String toString() => 'FloatingActionButtonLocation.customEndFloat';
}

class MainScaffold extends ConsumerWidget {
  const MainScaffold({
    super.key,
    required this.body,
    required this.currentIndex,
    this.onNavTap,
    this.title,
    this.actions,
    this.floatingActionButton,
  });

  final Widget body;
  final int currentIndex;
  final ValueChanged<int>? onNavTap;
  final String? title;
  final List<Widget>? actions;
  final Widget? floatingActionButton;

  void _handleNavTap(BuildContext context, WidgetRef ref, int index) {
    if (index == currentIndex) return;

    // Helper function to create a route without animation
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

    // Define routes - harus sesuai dengan menu items di BottomNavBar
    // Index 0: Dashboard, Index 1: Pipeline, Index 2: Leads, Index 3: Route, Index 4: Schedules
    final routes = [
      {
        'route': AppRoutes.dashboard,
        'resource': 'dashboard',
        'screen': const AuthGate(
          requiredRoute: AppRoutes.dashboard,
          child: DashboardScreen(),
        ),
      },
      {
        'route': AppRoutes.pipeline,
        'resource': 'pipeline',
        'screen': const AuthGate(
          requiredRoute: AppRoutes.pipeline,
          child: PipelineScreen(),
        ),
      },
      {
        'route': AppRoutes.leads,
        'resource': 'leads',
        'screen': const AuthGate(
          requiredRoute: AppRoutes.leads,
          child: LeadListScreen(),
        ),
      },
      {
        'route': AppRoutes.routeOptimization,
        'resource': 'route-optimization',
        'screen': const AuthGate(
          requiredRoute: AppRoutes.routeOptimization,
          child: RouteListScreen(),
        ),
      },
      {
        'route': AppRoutes.schedules,
        'resource': 'schedules',
        'screen': const AuthGate(
          requiredRoute: AppRoutes.schedules,
          child: ScheduleListScreen(),
        ),
      },
    ];

    if (index < 0 || index >= routes.length) {
      return;
    }

    final routeData = routes[index];
    final resource = routeData['resource'] as String;
    final alwaysAccessible = routeData['alwaysAccessible'] as bool? ?? false;

    // Check permission (profile is always accessible)
    if (!alwaysAccessible) {
      final hasPermission = PermissionChecker.canView(permissions, resource);
      if (!hasPermission) {
        // Redirect to dashboard if no permission
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('You do not have permission to access this page'),
            duration: Duration(seconds: 3),
          ),
        );
        Navigator.pushReplacement(
          context,
          createNoAnimationRoute(
            const AuthGate(
              requiredRoute: AppRoutes.dashboard,
              child: DashboardScreen(),
            ),
            AppRoutes.dashboard,
          ),
        );
        return;
      }
    }

    // Navigate to the route
    Navigator.pushReplacement(
      context,
      createNoAnimationRoute(
        routeData['screen'] as Widget,
        routeData['route'] as String,
      ),
    );
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      appBar: title != null
          ? AppBar(
              title: Text(title!),
              elevation: 0,
              actions: actions,
              automaticallyImplyLeading: false, // Remove hamburger menu
            )
          : null,
      drawer: null, // Remove drawer (hamburger menu)
      body: Stack(
        children: [
          Column(
            children: [
              const OfflineIndicator(),
              Expanded(child: body),
            ],
          ),
          Positioned(
            left: 0,
            right: 0,
            bottom: 0,
            child: BottomNavBar(
              currentIndex: currentIndex,
              onTap:
                  onNavTap ?? ((index) => _handleNavTap(context, ref, index)),
            ),
          ),
        ],
      ),
      floatingActionButton: floatingActionButton,
      floatingActionButtonLocation: floatingActionButton != null
          ? const _CustomFloatingActionButtonLocation()
          : null,
    );
  }
}
