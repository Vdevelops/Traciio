# Role-Based Access Control (RBAC) Implementation Guide
## CRM Healthcare Mobile App - Flutter

**Version**: 1.0  
**Last Updated**: 2025-01-15  
**Status**: ✅ **Completed**

---

## 📋 Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Permission Types](#permission-types)
4. [API Integration](#api-integration)
5. [Implementation Details](#implementation-details)
6. [Usage Examples](#usage-examples)
7. [File Structure](#file-structure)
8. [Testing](#testing)
9. [Troubleshooting](#troubleshooting)

---

## Overview

### Goals
- **Menu Permission**: Control which menu items appear in the bottom navigation bar based on user permissions
- **Action Permission**: Control which CRUD operations (Create, Edit, Delete) are available to users
- **Route Protection**: Prevent unauthorized access to pages
- **Offline Support**: Cache permissions for offline access

### Key Features
- ✅ Dynamic bottom navigation bar based on user permissions
- ✅ Action buttons (Create, Edit, Delete) visibility based on permissions
- ✅ Route protection with automatic redirect to dashboard
- ✅ Permission caching in Hive for offline support
- ✅ Automatic permission refresh in background

---

## Architecture

### Permission Flow

```
┌─────────────────────────────────────────────────────────┐
│                    Mobile App (Flutter)                  │
├─────────────────────────────────────────────────────────┤
│                                                           │
│  ┌──────────────┐      ┌──────────────┐                 │
│  │   UI Layer   │◄─────►│  Permission  │                 │
│  │  (Widgets)   │      │   Provider   │                 │
│  └──────────────┘      └──────────────┘                 │
│         │                      │                          │
│         │                      ▼                          │
│         │            ┌──────────────────┐                 │
│         │            │ Permission       │                 │
│         │            │ Helper           │                 │
│         │            │ (Utilities)      │                 │
│         │            └──────────────────┘                 │
│         │                      │                          │
│         │         ┌────────────┴────────────┐            │
│         │         │                         │            │
│         ▼         ▼                         ▼            │
│  ┌──────────┐ ┌──────────┐        ┌──────────────┐    │
│  │   Hive   │ │  Cache    │        │  API Client  │    │
│  │  (Local  │ │  (Memory)  │        │  (Dio)       │    │
│  │   DB)   │ │            │        │              │    │
│  └──────────┘ └──────────┘        └──────────────┘    │
│         │                                    │            │
│         │                                    ▼            │
│         │                            ┌──────────────┐    │
│         │                            │   Backend    │    │
│         │                            │    API       │    │
│         │                            │ /users/:id/  │    │
│         │                            │ permissions  │    │
│         │                            └──────────────┘    │
│                                                           │
└─────────────────────────────────────────────────────────┘
```

### Components

1. **Permission Provider** (`permission_provider.dart`)
   - Fetches permissions from API
   - Caches permissions in Hive
   - Provides permissions to UI via Riverpod

2. **Permission Helper** (`permission_helper.dart`)
   - Utility methods for permission checks
   - Menu permission validation
   - Action permission validation

3. **Permission Hooks** (`use_has_permission.dart`)
   - React-style hooks for easy permission checks in widgets
   - `useHasCreatePermission()`, `useHasEditPermission()`, `useHasDeletePermission()`

4. **Auth Gate** (`auth_gate.dart`)
   - Route protection component
   - Redirects unauthorized users to dashboard

---

## Permission Types

### 1. Menu Permission (VIEW)

Menu permissions control which menu items appear in the bottom navigation bar.

#### How It Works
- Checks for `VIEW_*` action codes in menu actions
- If user has `VIEW_DASHBOARD`, dashboard menu appears
- If user has `VIEW_ACCOUNTS`, accounts menu appears
- Profile menu is always accessible (no permission check)

#### Example API Response

```json
{
  "success": true,
  "data": {
    "menus": [
      {
        "id": "dashboard",
        "name": "Dashboard",
        "icon": "home",
        "url": "/dashboard",
        "actions": [
          {
            "id": "view_dashboard",
            "code": "VIEW_DASHBOARD",
            "name": "View Dashboard",
            "access": true
          }
        ]
      },
      {
        "id": "accounts",
        "name": "Accounts",
        "icon": "business",
        "url": "/accounts",
        "actions": [
          {
            "id": "view_accounts",
            "code": "VIEW_ACCOUNTS",
            "name": "View Accounts",
            "access": true
          }
        ]
      }
    ]
  }
}
```

### 2. Action Permission (CREATE, EDIT, DELETE)

Action permissions control which CRUD operations are available to users.

#### Action Codes

| Action | Code Pattern | Example |
|--------|-------------|---------|
| **CREATE** | `CREATE_<ENTITY>` or `CREATE_<ENTITY>S` | `CREATE_TASKS`, `CREATE_ACCOUNTS` |
| **EDIT** | `EDIT_<ENTITY>` or `UPDATE_<ENTITY>` | `EDIT_TASKS`, `UPDATE_TASKS` |
| **DELETE** | `DELETE_<ENTITY>` or `DELETE_<ENTITY>S` | `DELETE_TASKS`, `DELETE_ACCOUNTS` |

#### Entity Name Mapping

| Mobile Route | Entity Name | Action Codes |
|--------------|-------------|--------------|
| `/tasks` | `TASKS` | `CREATE_TASKS`, `EDIT_TASKS`, `DELETE_TASKS` |
| `/visit-reports` | `VISIT_REPORTS` | `CREATE_VISIT_REPORTS`, `EDIT_VISIT_REPORTS`, `DELETE_VISIT_REPORTS` |
| `/accounts` | `ACCOUNTS` | `CREATE_ACCOUNTS`, `EDIT_ACCOUNTS`, `DELETE_ACCOUNTS` |
| `/contacts` | `CONTACTS` | `CREATE_CONTACTS`, `EDIT_CONTACTS`, `DELETE_CONTACTS` |

#### Example API Response

```json
{
  "success": true,
  "data": {
    "menus": [
      {
        "id": "tasks",
        "name": "Tasks",
        "icon": "task",
        "url": "/tasks",
        "actions": [
          {
            "id": "view_tasks",
            "code": "VIEW_TASKS",
            "name": "View Tasks",
            "access": true
          },
          {
            "id": "create_tasks",
            "code": "CREATE_TASKS",
            "name": "Create Tasks",
            "access": true
          },
          {
            "id": "edit_tasks",
            "code": "EDIT_TASKS",
            "name": "Edit Tasks",
            "access": false
          },
          {
            "id": "delete_tasks",
            "code": "DELETE_TASKS",
            "name": "Delete Tasks",
            "access": false
          }
        ]
      }
    ]
  }
}
```

---

## API Integration

### Endpoint

**Mobile Endpoint** (Current):
```
GET /api/v1/auth/mobile/permissions
```

**Legacy Endpoint** (Deprecated):
```
GET /api/v1/users/:id/permissions
```

### Request Headers

```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Note**: Mobile endpoint does not require `userId` in the path. User is identified from the authentication token.

### Response Format (Mobile API)

```json
{
  "success": true,
  "data": {
    "menus": [
      {
        "menu": "dashboard",
        "actions": ["VIEW"]
      },
      {
        "menu": "task",
        "actions": ["VIEW", "CREATE", "EDIT", "DELETE"]
      },
      {
        "menu": "accounts",
        "actions": ["VIEW", "CREATE", "EDIT", "DELETE"]
      },
      {
        "menu": "contacts",
        "actions": ["VIEW", "CREATE", "EDIT", "DELETE"]
      },
      {
        "menu": "visit_reports",
        "actions": ["VIEW", "CREATE", "EDIT", "DELETE"]
      }
    ]
  },
  "timestamp": "2024-01-15T10:30:45+07:00",
  "request_id": "req_abc123xyz"
}
```

### Menu Name Mapping

The mobile API uses simplified menu names that are mapped to internal URLs:

| API Menu Name | Internal URL | Display Name |
|---------------|--------------|--------------|
| `dashboard` | `/dashboard` | Dashboard |
| `task` | `/tasks` | Tasks |
| `accounts` | `/accounts` | Accounts |
| `contacts` | `/contacts` | Contacts |
| `visit_reports` | `/visit-reports` | Visit Reports |

### Action Code Normalization

The mobile API returns simple action codes (e.g., `"VIEW"`, `"CREATE"`) that are normalized to full action codes:

| API Action | Normalized Code | Example |
|------------|----------------|---------|
| `VIEW` | `VIEW_<MENU>` | `VIEW_DASHBOARD`, `VIEW_TASK` |
| `CREATE` | `CREATE_<MENU>` | `CREATE_TASK`, `CREATE_ACCOUNTS` |
| `EDIT` | `EDIT_<MENU>` | `EDIT_TASK`, `EDIT_ACCOUNTS` |
| `DELETE` | `DELETE_<MENU>` | `DELETE_TASK`, `DELETE_ACCOUNTS` |

### Error Response

```json
{
  "success": false,
  "error": {
    "message": "Error message",
    "code": "ERROR_CODE"
  }
}
```

---

## Implementation Details

### 1. Permission Models

**File**: `apps/mobile/lib/features/permissions/data/models/permission.dart`

```dart
class Action {
  final String id;
  final String code;      // e.g., "CREATE_TASKS", "VIEW_ACCOUNTS"
  final String name;
  final bool access;      // true if user has permission
}

class MenuWithActions {
  final String id;
  final String name;
  final String icon;
  final String url;       // e.g., "/dashboard", "/accounts"
  final List<MenuWithActions>? children;
  final List<Action>? actions;
}

class UserPermissionsResponse {
  final List<MenuWithActions> menus;
}
```

### 2. Permission Repository

**File**: `apps/mobile/lib/features/permissions/data/permission_repository.dart`

```dart
class PermissionRepository {
  /// Get mobile permissions for authenticated user
  /// Uses mobile-specific endpoint: /api/v1/auth/mobile/permissions
  /// No userId required - uses token to identify user
  Future<UserPermissionsResponse> getMobilePermissions() async {
    final response = await _dio.get<Map<String, dynamic>>(
      '/api/v1/auth/mobile/permissions',
    );
    
    // Parse mobile API response format
    return UserPermissionsResponse.fromMobileJson(response.data!['data']);
  }
  
  /// Legacy method - kept for backward compatibility
  @Deprecated('Use getMobilePermissions() instead')
  Future<UserPermissionsResponse> getUserPermissions(String userId) async {
    return getMobilePermissions();
  }
}
```

### 3. Permission Provider

**File**: `apps/mobile/lib/features/permissions/application/permission_provider.dart`

```dart
final userPermissionsProvider = FutureProvider.autoDispose<UserPermissionsResponse>((ref) async {
  final repository = ref.read(permissionRepositoryProvider);
  final authState = ref.watch(authProvider);
  
  // Check authentication
  if (authState.status != AuthStatus.authenticated) {
    throw Exception('User not authenticated');
  }
  
  // Get user ID for cache key (mobile endpoint doesn't require userId)
  final userId = authState.user?.id;
  final cacheKey = userId != null ? 'user_permissions_$userId' : 'user_permissions_mobile';

  // Try cache first
  final cached = await OfflineStorage.get<UserPermissionsResponse>(
    cacheKey,
    (json) => UserPermissionsResponse.fromJson(json),
  );
  
  if (cached != null) {
    // Return cached, refresh in background
    Future.microtask(() async {
      try {
        final fresh = await repository.getMobilePermissions();
        await OfflineStorage.set(cacheKey, fresh.toJson());
      } catch (e) {
        // Ignore refresh errors
      }
    });
    return cached;
  }

  // Fetch from mobile API endpoint (no userId required)
  final permissions = await repository.getMobilePermissions();
  await OfflineStorage.set(cacheKey, permissions.toJson());
  return permissions;
});
```

### 4. Permission Helper

**File**: `apps/mobile/lib/features/permissions/utils/permission_helper.dart`

#### Menu Permission Check

```dart
static bool hasViewPermission(MenuWithActions menu) {
  // Check if menu has VIEW_* action with access = true
  if (menu.actions != null && menu.actions!.isNotEmpty) {
    final viewAction = menu.actions!.firstWhere(
      (action) => action.code.startsWith('VIEW_') && action.access,
      orElse: () => const Action(id: '', code: '', name: '', access: false),
    );
    if (viewAction.access) {
      return true;
    }
  }
  
  // Check children recursively
  if (menu.children != null && menu.children!.isNotEmpty) {
    return menu.children!.any((child) => hasViewPermission(child));
  }
  
  return false;
}
```

#### Action Permission Check

```dart
static bool hasActionPermission(
  List<MenuWithActions> menus,
  String url,
  String actionCode,
) {
  final menu = getMenuByUrl(menus, url);
  if (menu == null) return false;

  // Check if menu has the action with access = true
  if (menu.actions != null && menu.actions!.isNotEmpty) {
    final action = menu.actions!.firstWhere(
      (a) => a.code == actionCode && a.access,
      orElse: () => const Action(id: '', code: '', name: '', access: false),
    );
    if (action.access) return true;
  }

  // Check children recursively
  if (menu.children != null && menu.children!.isNotEmpty) {
    for (final child in menu.children!) {
      if (hasActionPermission([child], child.url, actionCode)) {
        return true;
      }
    }
  }

  return false;
}
```

### 5. Permission Hooks

**File**: `apps/mobile/lib/features/permissions/hooks/use_has_permission.dart`

```dart
/// Check CREATE permission
bool useHasCreatePermission(WidgetRef ref, String url) {
  final permissionsAsync = ref.watch(userPermissionsProvider);
  final menus = permissionsAsync.value?.menus ?? [];
  
  if (permissionsAsync.isLoading || permissionsAsync.hasError) {
    return false;
  }

  return PermissionHelper.hasCreatePermission(menus, url);
}

/// Check EDIT permission
bool useHasEditPermission(WidgetRef ref, String url) {
  // Similar implementation
}

/// Check DELETE permission
bool useHasDeletePermission(WidgetRef ref, String url) {
  // Similar implementation
}
```

### 6. Bottom Navigation Bar

**File**: `apps/mobile/lib/core/widgets/bottom_nav_bar.dart`

```dart
class BottomNavBar extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final permissionsAsync = ref.watch(userPermissionsProvider);
    
    List<BottomNavigationBarItem> navItems = [];
    
    permissionsAsync.when(
      data: (permissionsData) {
        final menus = permissionsData.menus;
        
        // Filter navigation items based on permissions
        final allNavItems = [
          _NavConfig(route: AppRoutes.dashboard, webUrl: '/dashboard', ...),
          _NavConfig(route: AppRoutes.accounts, webUrl: '/accounts', ...),
          // ...
        ];
        
        for (var config in allNavItems) {
          if (config.alwaysAccessible ||
              PermissionHelper.canAccessUrl(menus, config.route)) {
            navItems.add(/* Add to nav items */);
          }
        }
      },
      loading: () => /* Show default items */,
      error: (err, stack) => /* Show default items */,
    );
    
    return BottomNavigationBar(items: navItems);
  }
}
```

### 7. Route Protection

**File**: `apps/mobile/lib/core/widgets/auth_gate.dart`

```dart
class AuthGate extends ConsumerWidget {
  final Widget child;
  final String? requiredRoute;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final authState = ref.watch(authProvider);
    final permissionsAsync = ref.watch(userPermissionsProvider);

    if (authState.status == AuthStatus.authenticated) {
      if (requiredRoute != null) {
        // Profile is always accessible
        if (requiredRoute == AppRoutes.profile) {
          return child;
        }

        return permissionsAsync.when(
          data: (permissionsData) {
            final menus = permissionsData.menus;
            if (PermissionHelper.canAccessUrl(menus, requiredRoute!)) {
              return child;
            } else {
              // Redirect to dashboard
              WidgetsBinding.instance.addPostFrameCallback((_) {
                Navigator.of(context).pushReplacementNamed(AppRoutes.dashboard);
              });
              return const SizedBox.shrink();
            }
          },
          loading: () => const CircularProgressIndicator(),
          error: (err, stack) => /* Handle error */,
        );
      }
      return child;
    }

    return const LoginScreen();
  }
}
```

---

## Usage Examples

### 1. Check CREATE Permission in Widget

```dart
class TaskListScreen extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final hasCreatePermission = useHasCreatePermission(ref, '/tasks');
    
    return Scaffold(
      floatingActionButton: hasCreatePermission
          ? FloatingActionButton(
              onPressed: () {
                // Navigate to create task
              },
              child: const Icon(Icons.add),
            )
          : null,
      // ...
    );
  }
}
```

### 2. Check EDIT and DELETE Permissions

```dart
class TaskDetailScreen extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Builder(
      builder: (context) {
        final hasEditPermission = useHasEditPermission(ref, '/tasks');
        final hasDeletePermission = useHasDeletePermission(ref, '/tasks');
        
        return Column(
          children: [
            if (hasEditPermission) ...[
              FilledButton.icon(
                onPressed: () => _handleEdit(task),
                icon: const Icon(Icons.edit_outlined),
                label: Text('Edit'),
              ),
            ],
            if (hasDeletePermission) ...[
              OutlinedButton.icon(
                onPressed: () => _handleDelete(task),
                icon: const Icon(Icons.delete_outline),
                label: Text('Delete'),
              ),
            ],
          ],
        );
      },
    );
  }
}
```

### 3. Dynamic Navigation Bar

```dart
class MainScaffold extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final permissionsAsync = ref.watch(userPermissionsProvider);
    
    return Scaffold(
      bottomNavigationBar: BottomNavBar(
        currentIndex: currentIndex,
        onTap: (index) {
          // Navigation logic with permission checks
        },
      ),
      // ...
    );
  }
}
```

### 4. Route Protection

```dart
// In app_router.dart
AppRoutes.accounts: (_) =>
    const AuthGate(
      child: AccountsScreen(),
      requiredRoute: AppRoutes.accounts,
    ),
```

---

## File Structure

```
apps/mobile/lib/
├── features/
│   └── permissions/
│       ├── data/
│       │   ├── models/
│       │   │   └── permission.dart          # Action, MenuWithActions, UserPermissionsResponse
│       │   └── permission_repository.dart   # API calls
│       ├── application/
│       │   └── permission_provider.dart     # Riverpod provider with caching
│       ├── hooks/
│       │   └── use_has_permission.dart       # Permission hooks
│       └── utils/
│           └── permission_helper.dart       # Utility methods
├── core/
│   ├── widgets/
│   │   ├── auth_gate.dart                   # Route protection
│   │   └── bottom_nav_bar.dart               # Dynamic navigation bar
│   └── routing/
│       └── app_router.dart                   # Routes with AuthGate
└── core/
    └── storage/
        └── offline_storage.dart              # Permission caching
```

---

## Testing

### Manual Testing Checklist

#### Menu Permission
- [ ] Login as user with limited permissions
- [ ] Verify only accessible menus appear in bottom navigation
- [ ] Verify profile menu always appears
- [ ] Try navigating to restricted page directly (should redirect to dashboard)

#### Action Permission
- [ ] Login as user without CREATE permission
- [ ] Verify FloatingActionButton does not appear
- [ ] Login as user with CREATE permission
- [ ] Verify FloatingActionButton appears
- [ ] Test EDIT permission (button visibility)
- [ ] Test DELETE permission (button visibility)

#### Offline Support
- [ ] Login and verify permissions are cached
- [ ] Go offline
- [ ] Verify permissions still work (from cache)
- [ ] Go online
- [ ] Verify permissions refresh in background

### Unit Testing

```dart
// Example test for PermissionHelper
void main() {
  group('PermissionHelper', () {
    test('hasViewPermission returns true when VIEW action exists', () {
      final menu = MenuWithActions(
        id: 'dashboard',
        name: 'Dashboard',
        icon: 'home',
        url: '/dashboard',
        actions: [
          Action(
            id: 'view_dashboard',
            code: 'VIEW_DASHBOARD',
            name: 'View Dashboard',
            access: true,
          ),
        ],
      );
      
      expect(PermissionHelper.hasViewPermission(menu), isTrue);
    });
    
    test('hasCreatePermission returns true when CREATE action exists', () {
      final menus = [
        MenuWithActions(
          id: 'tasks',
          name: 'Tasks',
          icon: 'task',
          url: '/tasks',
          actions: [
            Action(
              id: 'create_tasks',
              code: 'CREATE_TASKS',
              name: 'Create Tasks',
              access: true,
            ),
          ],
        ),
      ];
      
      expect(
        PermissionHelper.hasCreatePermission(menus, '/tasks'),
        isTrue,
      );
    });
  });
}
```

---

## Troubleshooting

### Issue: Permissions not loading

**Symptoms**: Bottom navigation shows all items or no items

**Solutions**:
1. Check API endpoint is accessible: `GET /api/v1/users/:id/permissions`
2. Verify user ID is correct in `authProvider`
3. Check network connectivity
4. Verify token is valid and not expired
5. Check console for error messages

### Issue: Action buttons not showing/hiding correctly

**Symptoms**: Buttons appear when they shouldn't or don't appear when they should

**Solutions**:
1. Verify action codes match API response (e.g., `CREATE_TASKS` vs `CREATE_TASK`)
2. Check URL mapping in `PermissionHelper._getEntityNameFromUrl()`
3. Verify permission is cached correctly in Hive
4. Check if permission refresh is working

### Issue: Route protection not working

**Symptoms**: User can access restricted pages

**Solutions**:
1. Verify `AuthGate` wraps all protected routes
2. Check `requiredRoute` parameter is set correctly
3. Verify `PermissionHelper.canAccessUrl()` is working
4. Check if permissions are loaded before route check

### Issue: Permissions not cached offline

**Symptoms**: Permissions don't work when offline

**Solutions**:
1. Verify `OfflineStorage.set()` is called after fetching permissions
2. Check Hive box is initialized: `OfflineStorage.init()`
3. Verify cache key format: `user_permissions_$userId`
4. Check if permissions are fetched at least once when online

---

## Best Practices

### 1. Always Check Permissions Before Showing UI

```dart
// ✅ GOOD
final hasCreate = useHasCreatePermission(ref, '/tasks');
if (hasCreate) {
  return FloatingActionButton(...);
}

// ❌ BAD
return FloatingActionButton(...); // Always shows
```

### 2. Use Hooks for Permission Checks

```dart
// ✅ GOOD
final hasEdit = useHasEditPermission(ref, '/tasks');

// ❌ BAD
final permissions = ref.watch(userPermissionsProvider);
final hasEdit = permissions.when(
  data: (data) => PermissionHelper.hasEditPermission(data.menus, '/tasks'),
  // ... handle loading/error
);
```

### 3. Cache Permissions for Offline Support

```dart
// ✅ GOOD - Already implemented in permission_provider.dart
final cached = await OfflineStorage.get(...);
if (cached != null) {
  return cached; // Return immediately
}
// Fetch from API and cache
```

### 4. Handle Loading and Error States

```dart
// ✅ GOOD
permissionsAsync.when(
  data: (permissions) => /* Show UI */,
  loading: () => CircularProgressIndicator(),
  error: (err, stack) => ErrorWidget(err),
);
```

### 5. Profile is Always Accessible

```dart
// ✅ GOOD
if (requiredRoute == AppRoutes.profile) {
  return child; // Profile always accessible
}
```

---

## Future Enhancements

1. **Permission Groups**: Support for permission groups/roles
2. **Permission Inheritance**: Child menus inherit parent permissions
3. **Permission Audit Log**: Track permission changes
4. **Real-time Permission Updates**: WebSocket for permission changes
5. **Permission Testing Tools**: UI for testing permissions

---

## References

- [Riverpod Documentation](https://riverpod.dev/)
- [Hive Documentation](https://docs.hivedb.dev/)
- [Flutter State Management](https://flutter.dev/docs/development/data-and-backend/state-mgmt)

---

**Document Status**: Active  
**Last Updated**: 2025-01-15  
**Maintained By**: Development Team

