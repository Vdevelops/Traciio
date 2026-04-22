import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../features/auth/application/auth_provider.dart';
import '../../features/auth/application/auth_state.dart';
import '../permissions/permission_checker.dart';
import '../routing/app_router.dart';

class AuthGate extends ConsumerWidget {
  const AuthGate({
    super.key,
    required this.child,
    this.requiredRoute,
    this.requiredPermission,
  });

  final Widget child;
  final String? requiredRoute;
  /// Custom permission to check (overrides route-based permission check)
  /// Example: Use 'accounts.view' for contacts route
  final String? requiredPermission;

  /// Prevent multiple simultaneous navigations to login during logout
  static bool _isNavigatingToLogin = false;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final authState = ref.watch(authProvider);

    // Show loading saat check auth status (unknown)
    if (authState.status == AuthStatus.unknown) {
      return const Scaffold(
        backgroundColor: Colors.white,
        body: Center(
          child: CircularProgressIndicator(),
        ),
      );
    }

    // If not authenticated, navigate to login and clear stack
    if (authState.status != AuthStatus.authenticated) {
      if (!_isNavigatingToLogin) {
        _isNavigatingToLogin = true;
        WidgetsBinding.instance.addPostFrameCallback((_) {
          if (context.mounted) {
            Navigator.of(context).pushNamedAndRemoveUntil(
              AppRoutes.login,
              (route) => false,
            );
          }
          // Reset flag after frame completes
          Future.microtask(() => _isNavigatingToLogin = false);
        });
      }
      return const Scaffold(
        backgroundColor: Colors.white,
        body: Center(
          child: CircularProgressIndicator(),
        ),
      );
    }

    // If no route specified, just return child (for backward compatibility)
    if (requiredRoute == null) {
      return child;
    }

    // Get permissions from user
    final permissions = authState.user?.permissions ?? [];

    // Check if user can access this route
    // If custom permission is provided, use it; otherwise use route-based check
    final hasPermission = requiredPermission != null
        ? PermissionChecker.hasPermission(permissions, requiredPermission!)
        : PermissionChecker.canAccessRoute(permissions, requiredRoute!);

    if (!hasPermission) {
      // Redirect to dashboard if no permission
      WidgetsBinding.instance.addPostFrameCallback((_) {
        Navigator.of(context).pushReplacementNamed(AppRoutes.dashboard);
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('You do not have permission to access this page'),
            duration: Duration(seconds: 3),
          ),
        );
      });
      return const Scaffold(
        backgroundColor: Colors.white,
        body: Center(
          child: CircularProgressIndicator(),
        ),
      );
    }

    return child;
  }
}


