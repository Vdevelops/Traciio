/// Permission checker untuk RBAC dengan dot notation format
/// 
/// Format permission: "resource.action"
/// Contoh: "dashboard.view", "tasks.create", "visit-reports.edit"
/// 
/// Supports:
/// - Exact match: "tasks.view"
/// - Resource wildcard: "tasks.*" (semua actions untuk tasks)
/// - Global wildcard: "*" (full access - superadmin)
class PermissionChecker {
  const PermissionChecker._();

  /// Check if user has specific permission
  /// 
  /// Examples:
  /// - hasPermission(["tasks.view"], "tasks.view") → true
  /// - hasPermission(["tasks.*"], "tasks.create") → true
  /// - hasPermission(["*"], "anything.here") → true
  static bool hasPermission(List<String> userPermissions, String requiredPermission) {
    if (userPermissions.isEmpty) {
      return false;
    }

    // Check for global wildcard (superadmin)
    if (userPermissions.contains('*')) {
      return true;
    }

    // Check for exact match
    if (userPermissions.contains(requiredPermission)) {
      return true;
    }

    // Check for resource wildcard (e.g., "tasks.*" allows "tasks.view", "tasks.create")
    // Also handle multi-part actions like "pipeline.opportunity-create"
    final parts = requiredPermission.split('.');
    if (parts.length >= 2) {
      final resourceWildcard = '${parts[0]}.*';
      if (userPermissions.contains(resourceWildcard)) {
        return true;
      }
    }

    return false;
  }

  /// Check if user has ANY of the specified permissions
  /// 
  /// Example:
  /// - hasAnyPermission(["tasks.view"], ["tasks.view", "tasks.edit"]) → true
  static bool hasAnyPermission(List<String> userPermissions, List<String> requiredPermissions) {
    return requiredPermissions.any((perm) => hasPermission(userPermissions, perm));
  }

  /// Check if user has ALL of the specified permissions
  /// 
  /// Example:
  /// - hasAllPermissions(["tasks.view", "tasks.edit"], ["tasks.view", "tasks.edit"]) → true
  static bool hasAllPermissions(List<String> userPermissions, List<String> requiredPermissions) {
    return requiredPermissions.every((perm) => hasPermission(userPermissions, perm));
  }

  /// Check VIEW permission for a resource
  static bool canView(List<String> userPermissions, String resource) {
    return hasPermission(userPermissions, '$resource.view');
  }

  /// Check CREATE permission for a resource
  static bool canCreate(List<String> userPermissions, String resource) {
    return hasPermission(userPermissions, '$resource.create');
  }

  /// Check EDIT permission for a resource
  static bool canEdit(List<String> userPermissions, String resource) {
    return hasPermission(userPermissions, '$resource.edit');
  }

  /// Check DELETE permission for a resource
  static bool canDelete(List<String> userPermissions, String resource) {
    return hasPermission(userPermissions, '$resource.delete');
  }

  /// Check COMPLETE permission for tasks
  static bool canCompleteTask(List<String> userPermissions) {
    return hasPermission(userPermissions, 'tasks.complete');
  }

  /// Check START permission for tasks
  static bool canStartTask(List<String> userPermissions) {
    return hasPermission(userPermissions, 'tasks.start');
  }

  /// Check CANCEL permission for tasks
  static bool canCancelTask(List<String> userPermissions) {
    return hasPermission(userPermissions, 'tasks.cancel');
  }

  /// Check MARK-READ permission for notifications
  static bool canMarkNotificationRead(List<String> userPermissions) {
    return hasPermission(userPermissions, 'notifications.mark-read');
  }

  /// Check DELETE permission for notifications
  static bool canDeleteNotification(List<String> userPermissions) {
    return hasPermission(userPermissions, 'notifications.delete');
  }

  /// Check CHANGE-PASSWORD permission for profile
  static bool canChangePassword(List<String> userPermissions) {
    return hasPermission(userPermissions, 'profile.change-password');
  }

  /// Check OPPORTUNITY-CREATE permission for pipeline
  static bool canCreateOpportunity(List<String> userPermissions) {
    return hasPermission(userPermissions, 'pipeline.opportunity-create');
  }

  /// Check OPPORTUNITY-EDIT permission for pipeline
  static bool canEditOpportunity(List<String> userPermissions) {
    return hasPermission(userPermissions, 'pipeline.opportunity-edit');
  }

  /// Get all permissions for a specific resource
  /// 
  /// Example: getResourcePermissions(["tasks.view", "tasks.create", "users.view"], "tasks")
  /// Returns: ["tasks.view", "tasks.create"]
  static List<String> getResourcePermissions(List<String> userPermissions, String resource) {
    return userPermissions
        .where((perm) => perm.startsWith('$resource.'))
        .toList();
  }

  /// Get all available actions for a resource
  /// 
  /// Example: getResourceActions(["tasks.view", "tasks.create"], "tasks")
  /// Returns: ["view", "create"]
  static List<String> getResourceActions(List<String> userPermissions, String resource) {
    return getResourcePermissions(userPermissions, resource)
        .map((perm) => perm.split('.').last)
        .toList();
  }

  /// Check if user can access a route based on resource
  /// 
  /// Route mapping:
  /// - /dashboard → dashboard.view
  /// - /tasks → tasks.view
  /// - /tasks/create → tasks.create
  /// - /tasks/:id/edit → tasks.edit
  /// - /visit-reports → visit-reports.view
  static bool canAccessRoute(List<String> userPermissions, String route) {
    // Remove leading slash
    final cleanRoute = route.startsWith('/') ? route.substring(1) : route;
    
    if (cleanRoute.isEmpty || cleanRoute == '/') {
      return true; // Root always accessible
    }

    // Parse route to determine resource and action
    final parts = cleanRoute.split('/');
    final resource = parts[0];

    // Determine action from route
    String action = 'view'; // Default action
    
    if (parts.length > 1) {
      final lastPart = parts.last;
      if (lastPart == 'create' || lastPart == 'new') {
        action = 'create';
      } else if (lastPart == 'edit') {
        action = 'edit';
      } else if (parts.contains('edit')) {
        action = 'edit';
      }
    }

    return hasPermission(userPermissions, '$resource.$action');
  }

  /// Get list of accessible resources (for navigation menu)
  /// 
  /// Returns unique list of resources user can view
  /// Example: ["dashboard", "tasks", "visit-reports"]
  static List<String> getAccessibleResources(List<String> userPermissions) {
    return userPermissions
        .where((perm) => perm.contains('.view') || perm.contains('.*') || perm == '*')
        .map((perm) {
          if (perm == '*') {
            return 'all';
          }
          return perm.split('.').first;
        })
        .toSet()
        .toList();
  }

  /// Check if user is superadmin (has * permission)
  static bool isSuperAdmin(List<String> userPermissions) {
    return userPermissions.contains('*');
  }
}
