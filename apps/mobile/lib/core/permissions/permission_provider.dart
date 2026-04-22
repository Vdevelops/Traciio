import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../features/auth/application/auth_provider.dart';

/// Provider untuk mendapatkan user permissions dari AuthState
/// 
/// Permissions sudah tersedia langsung setelah login, tidak perlu fetch terpisah
final userPermissionsProvider = Provider<List<String>>((ref) {
  final authState = ref.watch(authProvider);
  return authState.user?.permissions ?? [];
});

/// Provider untuk check apakah user punya permission tertentu
/// 
/// Usage:
/// ```dart
/// final canCreateTask = ref.watch(hasPermissionProvider('tasks.create'));
/// ```
final hasPermissionProvider = Provider.family<bool, String>((ref, permission) {
  final permissions = ref.watch(userPermissionsProvider);
  
  if (permissions.isEmpty) {
    return false;
  }

  // Check for global wildcard (superadmin)
  if (permissions.contains('*')) {
    return true;
  }

  // Check for exact match
  if (permissions.contains(permission)) {
    return true;
  }

  // Check for resource wildcard (e.g., "tasks.*" allows "tasks.view", "tasks.create")
  final parts = permission.split('.');
  if (parts.length == 2) {
    final resourceWildcard = '${parts[0]}.*';
    if (permissions.contains(resourceWildcard)) {
      return true;
    }
  }

  return false;
});

/// Provider untuk check apakah user bisa view resource tertentu
/// 
/// Usage:
/// ```dart
/// final canViewTasks = ref.watch(canViewProvider('tasks'));
/// ```
final canViewProvider = Provider.family<bool, String>((ref, resource) {
  return ref.watch(hasPermissionProvider('$resource.view'));
});

/// Provider untuk check apakah user bisa create resource tertentu
final canCreateProvider = Provider.family<bool, String>((ref, resource) {
  return ref.watch(hasPermissionProvider('$resource.create'));
});

/// Provider untuk check apakah user bisa edit resource tertentu
final canEditProvider = Provider.family<bool, String>((ref, resource) {
  return ref.watch(hasPermissionProvider('$resource.edit'));
});

/// Provider untuk check apakah user bisa delete resource tertentu
final canDeleteProvider = Provider.family<bool, String>((ref, resource) {
  return ref.watch(hasPermissionProvider('$resource.delete'));
});

/// Provider untuk check apakah user bisa complete task
final canCompleteTaskProvider = Provider<bool>((ref) {
  return ref.watch(hasPermissionProvider('tasks.complete'));
});

/// Provider untuk check apakah user bisa start task
final canStartTaskProvider = Provider<bool>((ref) {
  return ref.watch(hasPermissionProvider('tasks.start'));
});

/// Provider untuk check apakah user bisa cancel task
final canCancelTaskProvider = Provider<bool>((ref) {
  return ref.watch(hasPermissionProvider('tasks.cancel'));
});

/// Provider untuk check apakah user bisa mark notification as read
final canMarkNotificationReadProvider = Provider<bool>((ref) {
  return ref.watch(hasPermissionProvider('notifications.mark-read'));
});

/// Provider untuk check apakah user bisa delete notification
final canDeleteNotificationProvider = Provider<bool>((ref) {
  return ref.watch(hasPermissionProvider('notifications.delete'));
});

/// Provider untuk check apakah user bisa change password
final canChangePasswordProvider = Provider<bool>((ref) {
  return ref.watch(hasPermissionProvider('profile.change-password'));
});

/// Provider untuk check apakah user bisa create opportunity
final canCreateOpportunityProvider = Provider<bool>((ref) {
  return ref.watch(hasPermissionProvider('pipeline.opportunity-create'));
});

/// Provider untuk check apakah user bisa edit opportunity
final canEditOpportunityProvider = Provider<bool>((ref) {
  return ref.watch(hasPermissionProvider('pipeline.opportunity-edit'));
});

/// Provider untuk mendapatkan list resource yang bisa di-view oleh user
/// 
/// Useful untuk dynamic navigation menu
final accessibleResourcesProvider = Provider<List<String>>((ref) {
  final permissions = ref.watch(userPermissionsProvider);
  
  return permissions
      .where((perm) => perm.contains('.view') || perm.contains('.*') || perm == '*')
      .map((perm) {
        if (perm == '*') {
          return 'all';
        }
        return perm.split('.').first;
      })
      .toSet()
      .toList();
});

/// Provider untuk check apakah user adalah superadmin
final isSuperAdminProvider = Provider<bool>((ref) {
  final permissions = ref.watch(userPermissionsProvider);
  return permissions.contains('*');
});
