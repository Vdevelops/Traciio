# Core - Navigation & Routing

## CRM Healthcare Mobile App - Flutter

**Module**: Core Infrastructure  
**Sprint**: Sprint 0  
**Version**: 1.1  
**Status**: ✅ **Completed**  
**Last Updated**: March 2026

---

## Table of Contents

1. [Ringkasan Fitur](#ringkasan-fitur)
2. [Fitur Utama](#fitur-utama)
3. [Business Rules](#business-rules)
4. [Keputusan Teknis & Trade-offs](#keputusan-teknis--trade-offs)
5. [Struktur Folder](#struktur-folder)
6. [API / Package Reference](#api--package-reference)
7. [Configuration](#configuration)
8. [Route Definitions](#route-definitions)
9. [Navigation Patterns](#navigation-patterns)
10. [Cara Test Manual](#cara-test-manual)
11. [Dependencies](#dependencies)
12. [Notes & Improvements](#notes--improvements)

---

## Ringkasan Fitur

Sistem **Navigation & Routing** mobile app CRM Healthcare menggunakan **Go Router** sebagai routing solution dengan named routes, deep linking support, dan route protection melalui AuthGate. Sistem ini menyediakan seamless navigation experience dengan bottom navigation bar, nested navigation, dan automatic handling untuk authentication state.

### Goals

- **Named Routes**: Type-safe route definitions dengan constants
- **Route Protection**: Auth guards untuk melindungi protected routes
- **Deep Linking**: Support untuk external links ke dalam app
- **Bottom Navigation**: Persistent bottom bar dengan state preservation
- **Navigation State**: Proper state management untuk back button handling

---

## Fitur Utama

### 1. Named Routes

Route definitions menggunakan class constants untuk type safety:

```dart
class AppRoutes {
  static const String login = '/login';
  static const String dashboard = '/dashboard';
  static const String accounts = '/accounts';
  static const String accountDetail = '/accounts/:id';
  static const String contacts = '/contacts';
  static const String tasks = '/tasks';
  static const String visitReports = '/visit-reports';
  static const String profile = '/profile';
}
```

### 2. Route Protection

**AuthGate**: Widget wrapper yang memeriksa authentication status sebelum mengakses route:

```
User Request → AuthGate → Check Auth → [Valid] → Show Screen
                                     ↓
                                [Invalid] → Redirect to Login
```

### 3. Bottom Navigation Bar

Persistent bottom navigation dengan 3-5 menu items berdasarkan user permissions:

- Dashboard (always visible)
- Accounts
- Tasks
- Visit Reports
- Profile (always visible)

### 4. Deep Linking

Support untuk external URLs:

- `crmhealth://dashboard` → Open dashboard
- `crmhealth://accounts/:id` → Open account detail
- `https://crm.example.com/accounts/:id` → Universal links

### 5. Nested Navigation

Navigation di dalam bottom navigation tabs dengan proper state preservation:

- Each tab memiliki navigation stack sendiri
- Switching tabs preserves navigation state
- Back button handling per tab

---

## Business Rules

### 1. Route Accessibility

| Route            | Requires Auth | Permission Check   | Notes                   |
| ---------------- | ------------- | ------------------ | ----------------------- |
| `/login`         | No            | No                 | Always accessible       |
| `/dashboard`     | Yes           | VIEW_DASHBOARD     | Default after login     |
| `/accounts`      | Yes           | VIEW_ACCOUNTS      | Main navigation item    |
| `/accounts/:id`  | Yes           | VIEW_ACCOUNTS      | Detail screen           |
| `/contacts`      | Yes           | VIEW_CONTACTS      | Sub-menu under accounts |
| `/tasks`         | Yes           | VIEW_TASKS         | Main navigation item    |
| `/visit-reports` | Yes           | VIEW_VISIT_REPORTS | Main navigation item    |
| `/profile`       | Yes           | No                 | Always accessible       |

### 2. Navigation Hierarchy

```
App
├── Login (initial jika tidak terautentikasi)
└── Main Scaffold (jika terautentikasi)
    ├── Dashboard Tab
    │   └── Dashboard Screen
    ├── Accounts Tab
    │   ├── Accounts List Screen
    │   ├── Account Detail Screen
    │   └── Contacts Screen
    ├── Tasks Tab
    │   ├── Tasks List Screen
    │   └── Task Detail Screen
    ├── Visit Reports Tab
    │   ├── Visit Reports List Screen
    │   ├── Visit Report Detail Screen
    │   └── Visit Report Form Screen
    └── Profile Tab
        └── Profile Screen
```

### 3. Navigation State Management

- **CurrentRoute**: Track current route untuk bottom nav highlighting
- **NavigationStack**: Stack untuk back button handling
- **TabHistory**: History per tab untuk proper back behavior
- **RouteArguments**: Passing data antar screens

### 4. Redirect Rules

1. **Unauthenticated Access**: Redirect ke `/login`
2. **Authenticated Access to Login**: Redirect ke `/dashboard`
3. **No Permission**: Redirect ke `/dashboard` dengan error message
4. **Unknown Route**: Redirect ke `/dashboard`

---

## Keputusan Teknis & Trade-offs

### Mengapa Go Router, bukan Navigator 1.0?

**Keputusan**: Menggunakan Go Router (Navigator 2.0) daripada traditional Navigator.

**Alasan**:

- **Deep Linking**: Native support untuk deep links dan URL-based navigation
- **Type Safety**: Named routes dengan compile-time safety
- **Web Support**: Jika ada rencana web version, Go Router sudah ready
- **State Restoration**: Better support untuk state restoration
- **Declarative**: Declarative routing configuration

**Trade-off**: Learning curve lebih tinggi dibanding Navigator 1.0. **Mitigasi**: Pattern yang consistent di seluruh app.

### Mengapa Nested Navigation per Tab?

**Keputusan**: Setiap tab memiliki navigation stack sendiri (CupertinoTabBar pattern).

**Alasan**:

- **User Experience**: User expects setiap tab memiliki history sendiri
- **State Preservation**: Navigating ke tab lain tidak menghilangkan state
- **Back Button**: Back button hanya affects current tab
- **Industry Standard**: Pattern yang umum di iOS dan modern Android apps

**Trade-off**: Kompleksitas lebih tinggi untuk manage multiple navigators. **Mitigasi**: Go Router handles ini dengan baik melalui ShellRoute.

### Static Routes vs Dynamic Routes

**Keputusan**: Static route constants (AppRoutes) dengan dynamic route builder.

**Alasan**:

- **Type Safety**: Compile-time checking untuk route names
- **Refactoring**: Mudah rename routes tanpa fear of breaking
- **Documentation**: Routes terdocumentasi di satu tempat
- **Testing**: Mudah mock routes di tests

---

## Struktur Folder

```
apps/mobile/lib/
├── core/
│   ├── routing/
│   │   ├── app_router.dart              # Main GoRouter configuration
│   │   ├── app_routes.dart              # Route constants (AppRoutes class)
│   │   ├── route_utils.dart             # Navigation helper methods
│   │   └── router_refresh_listener.dart # Auth state listener untuk redirect
│   └── widgets/
│       ├── auth_gate.dart               # Route protection wrapper
│       ├── bottom_nav_bar.dart          # Bottom navigation bar
│       └── main_scaffold.dart           # Main scaffold dengan bottom nav
└── features/
    └── [feature]/
        └── presentation/
            └── screens/
                └── [screen].dart        # Individual screens
```

---

## API / Package Reference

### Go Router Package

**Package**: `go_router: ^13.0.0`

**Main Components**:

#### GoRouter Configuration

```dart
final router = GoRouter(
  initialLocation: AppRoutes.login,
  refreshListenable: authNotifier, // Listen untuk auth state changes
  redirect: (context, state) {
    final isAuthenticated = authNotifier.isAuthenticated;
    final isLoggingIn = state.matchedLocation == AppRoutes.login;

    if (!isAuthenticated && !isLoggingIn) {
      return AppRoutes.login;
    }

    if (isAuthenticated && isLoggingIn) {
      return AppRoutes.dashboard;
    }

    return null; // No redirect
  },
  routes: [
    // Route definitions...
  ],
);
```

#### Route Types

**1. Simple Route**:

```dart
GoRoute(
  path: AppRoutes.dashboard,
  builder: (context, state) => const DashboardScreen(),
)
```

**2. Route with Parameters**:

```dart
GoRoute(
  path: AppRoutes.accountDetail,
  builder: (context, state) {
    final accountId = state.pathParameters['id']!;
    return AccountDetailScreen(accountId: accountId);
  },
)
```

**3. Route with Query Parameters**:

```dart
GoRoute(
  path: AppRoutes.accounts,
  builder: (context, state) {
    final searchQuery = state.uri.queryParameters['search'];
    return AccountsScreen(initialSearch: searchQuery);
  },
)
```

**4. ShellRoute (Nested Navigation)**:

```dart
ShellRoute(
  builder: (context, state, child) => MainScaffold(child: child),
  routes: [
    GoRoute(
      path: AppRoutes.dashboard,
      builder: (context, state) => const DashboardScreen(),
    ),
    GoRoute(
      path: AppRoutes.accounts,
      builder: (context, state) => const AccountsScreen(),
    ),
    // ... other routes
  ],
)
```

### Navigation Methods

#### Programmatic Navigation

```dart
// Navigate to route
context.go(AppRoutes.dashboard);

// Navigate with parameters
context.go('/accounts/123');

// Navigate dengan push (adds to stack)
context.push(AppRoutes.accountDetail.replaceAll(':id', '123'));

// Navigate dengan query parameters
context.go('${AppRoutes.accounts}?search=john');

// Pop current route
context.pop();

// Pop dengan result
context.pop(someData);
```

#### Deep Link Handling

```dart
// Android: AndroidManifest.xml
<intent-filter>
  <action android:name="android.intent.action.VIEW" />
  <category android:name="android.intent.category.DEFAULT" />
  <category android:name="android.intent.category.BROWSABLE" />
  <data android:scheme="crmhealth" />
</intent-filter>

// iOS: Info.plist
<key>CFBundleURLTypes</key>
<array>
  <dict>
    <key>CFBundleURLName</key>
    <string>com.example.crmhealth</string>
    <key>CFBundleURLSchemes</key>
    <array>
      <string>crmhealth</string>
    </array>
  </dict>
</array>
```

---

## Configuration

### Route Constants

**File**: `core/routing/app_routes.dart`

```dart
class AppRoutes {
  // Auth
  static const String login = '/login';
  static const String register = '/register';
  static const String forgotPassword = '/forgot-password';

  // Main
  static const String dashboard = '/dashboard';

  // Accounts
  static const String accounts = '/accounts';
  static const String accountDetail = '/accounts/:id';
  static const String accountCreate = '/accounts/create';
  static const String accountEdit = '/accounts/:id/edit';

  // Contacts
  static const String contacts = '/contacts';
  static const String contactDetail = '/contacts/:id';

  // Tasks
  static const String tasks = '/tasks';
  static const String taskDetail = '/tasks/:id';
  static const String taskCreate = '/tasks/create';
  static const String taskEdit = '/tasks/:id/edit';

  // Visit Reports
  static const String visitReports = '/visit-reports';
  static const String visitReportDetail = '/visit-reports/:id';
  static const String visitReportCreate = '/visit-reports/create';
  static const String visitReportEdit = '/visit-reports/:id/edit';

  // Profile
  static const String profile = '/profile';
  static const String settings = '/settings';

  // Helper methods
  static String accountDetailPath(String id) => '/accounts/$id';
  static String accountEditPath(String id) => '/accounts/$id/edit';
  static String taskDetailPath(String id) => '/tasks/$id';
  static String visitReportDetailPath(String id) => '/visit-reports/$id';
}
```

### Router Configuration

**File**: `core/routing/app_router.dart`

```dart
final routerProvider = Provider<GoRouter>((ref) {
  final authState = ref.watch(authProvider);

  return GoRouter(
    initialLocation: AppRoutes.login,
    refreshListenable: RouterRefreshListener(authState),
    redirect: (context, state) {
      final isAuthenticated = authState.status == AuthStatus.authenticated;
      final currentPath = state.matchedLocation;

      // Public routes yang tidak memerlukan auth
      final publicRoutes = [
        AppRoutes.login,
        AppRoutes.register,
        AppRoutes.forgotPassword,
      ];

      final isPublicRoute = publicRoutes.any((route) =>
        currentPath.startsWith(route));

      // Jika tidak terautentikasi dan mencoba akses protected route
      if (!isAuthenticated && !isPublicRoute) {
        return AppRoutes.login;
      }

      // Jika terautentikasi dan mencoba akses login
      if (isAuthenticated && currentPath == AppRoutes.login) {
        return AppRoutes.dashboard;
      }

      return null;
    },
    routes: [
      // Login route
      GoRoute(
        path: AppRoutes.login,
        builder: (context, state) => const LoginScreen(),
      ),

      // Main scaffold dengan bottom navigation
      ShellRoute(
        builder: (context, state, child) => MainScaffold(child: child),
        routes: [
          GoRoute(
            path: AppRoutes.dashboard,
            builder: (context, state) => const DashboardScreen(),
          ),
          GoRoute(
            path: AppRoutes.accounts,
            builder: (context, state) => const AccountsScreen(),
            routes: [
              GoRoute(
                path: ':id',
                builder: (context, state) {
                  final id = state.pathParameters['id']!;
                  return AccountDetailScreen(accountId: id);
                },
              ),
            ],
          ),
          // ... other routes
        ],
      ),
    ],
    errorBuilder: (context, state) => const ErrorScreen(),
  );
});
```

---

## Route Definitions

### Main Routes

| Route               | Path                 | Screen                  | Parameters           |
| ------------------- | -------------------- | ----------------------- | -------------------- |
| Login               | `/login`             | LoginScreen             | -                    |
| Dashboard           | `/dashboard`         | DashboardScreen         | -                    |
| Accounts            | `/accounts`          | AccountsScreen          | `?search=&page=`     |
| Account Detail      | `/accounts/:id`      | AccountDetailScreen     | `:id`                |
| Contacts            | `/contacts`          | ContactsScreen          | `?accountId=`        |
| Tasks               | `/tasks`             | TasksScreen             | `?status=&priority=` |
| Task Detail         | `/tasks/:id`         | TaskDetailScreen        | `:id`                |
| Visit Reports       | `/visit-reports`     | VisitReportsScreen      | `?status=`           |
| Visit Report Detail | `/visit-reports/:id` | VisitReportDetailScreen | `:id`                |
| Profile             | `/profile`           | ProfileScreen           | -                    |

### Route Guards

**AuthGate Implementation**:

```dart
class AuthGate extends ConsumerWidget {
  final Widget child;
  final String? requiredPermission;

  // Static flag untuk mencegah multiple navigasi saat logout
  static bool _isNavigatingToLogin = false;

  const AuthGate({
    super.key,
    required this.child,
    this.requiredPermission,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final authState = ref.watch(authProvider);
    final permissionsAsync = ref.watch(userPermissionsProvider);

    // Check authentication - navigate ke login dengan clear stack
    if (authState.status != AuthStatus.authenticated) {
      // Hindari multiple navigasi
      if (!_isNavigatingToLogin) {
        _isNavigatingToLogin = true;

        WidgetsBinding.instance.addPostFrameCallback((_) {
          navigatorKey.currentState?.pushNamedAndRemoveUntil(
            AppRoutes.login,
            (route) => false,
          );
          // Reset flag setelah navigasi
          Future.delayed(const Duration(milliseconds: 500), () {
            _isNavigatingToLogin = false;
          });
        });
      }
      // Tampilkan loading saat menunggu navigasi
      return const Scaffold(
        body: Center(
          child: CircularProgressIndicator(),
        ),
      );
    }

    // Check permission jika diperlukan
    if (requiredPermission != null) {
      return permissionsAsync.when(
        data: (permissions) {
          final hasPermission = permissions.hasPermission(requiredPermission!);
          if (!hasPermission) {
            WidgetsBinding.instance.addPostFrameCallback((_) {
              context.go(AppRoutes.dashboard);
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(
                  content: Text('Anda tidak memiliki akses ke halaman ini'),
                ),
              );
            });
            return const SizedBox.shrink();
          }
          return child;
        },
        loading: () => const CircularProgressIndicator(),
        error: (_, __) => const ErrorWidget(),
      );
    }

    return child;
  }
}
```

**Important Changes**:

1. **Static Flag `_isNavigatingToLogin`**: Mencegah multiple navigasi yang bersamaan saat logout
2. **Loading Indicator**: Menampilkan loading screen bukan `LoginScreen` inline
3. **`pushNamedAndRemoveUntil`**: Membersihkan navigation stack sepenuhnya saat logout
4. **Reset Flag**: Flag di-reset setelah 500ms untuk memungkinkan navigasi berikutnya

**Why These Changes**:

Sebelumnya, saat logout terjadi:

- Multiple `AuthGate` di stack bereaksi bersamaan
- Render `LoginScreen` inline + loading spinner + navigasi = glitch
- Navigation conflict menyebabkan UI tidak responsif

Setelah perubahan:

- Single clean navigation path
- Tidak ada render konflik
- Loading spinner yang jelas saat transisi
- Navigation stack sepenuhnya di-clear

---

## Navigation Patterns

### 1. Navigate to Detail Screen

```dart
// From list screen
void _onAccountTap(Account account) {
  context.push(AppRoutes.accountDetailPath(account.id));
}
```

### 2. Navigate with Result

```dart
// Navigate dan tunggu result
final result = await context.push(AppRoutes.taskCreate);

if (result == true) {
  // Task created successfully, refresh list
  ref.read(tasksProvider.notifier).refresh();
}

// Pop dengan result
context.pop(true);
```

### 3. Replace Current Route

```dart
// Replace login dengan dashboard setelah login success
context.go(AppRoutes.dashboard);

// Atau menggunakan pushReplacement
context.pushReplacement(AppRoutes.dashboard);
```

### 4. Clear Stack and Navigate

```dart
// Logout: Clear stack dan navigate ke login
context.go(AppRoutes.login);
```

### 5. Deep Link Navigation

```dart
// Handle deep link dari push notification
void handleDeepLink(String path) {
  if (path.startsWith('/tasks/')) {
    final taskId = path.split('/').last;
    context.push(AppRoutes.taskDetailPath(taskId));
  } else if (path.startsWith('/accounts/')) {
    final accountId = path.split('/').last;
    context.push(AppRoutes.accountDetailPath(accountId));
  }
}
```

---

## Cara Test Manual

### Test Route Navigation

1. **Navigation to Detail**:
   - Buka Accounts List
   - Tap salah satu account
   - Verifikasi: Navigate ke Account Detail dengan ID yang benar
   - Verifikasi: Back button kembali ke Accounts List

2. **Deep Linking**:
   - Buka browser dan navigate ke `crmhealth://accounts/123`
   - Verifikasi: App terbuka dan navigate ke Account Detail
   - Verifikasi: Data account ID benar

3. **Route Protection**:
   - Logout dari app
   - Coba navigate ke `/dashboard` secara manual (melalui adb atau URL)
   - Verifikasi: Redirect ke login screen

4. **Auth Redirect**:
   - Login sebagai user
   - Coba navigate ke `/login`
   - Verifikasi: Redirect ke dashboard

5. **Bottom Navigation**:
   - Navigate ke Accounts tab
   - Buka Account Detail
   - Switch ke Tasks tab
   - Switch kembali ke Accounts tab
   - Verifikasi: Masih di Account Detail (state preserved)

### Test Query Parameters

1. **Search Parameter**:
   - Navigate ke `/accounts?search=john`
   - Verifikasi: Accounts list terfilter dengan search query "john"

2. **Filter Parameters**:
   - Navigate ke `/tasks?status=pending&priority=high`
   - Verifikasi: Task list terfilter berdasarkan status dan priority

---

## Dependencies

### Internal

- `features/auth/application/auth_provider.dart` - Auth state untuk redirect logic
- `features/permissions/application/permission_provider.dart` - Permission checks
- `core/widgets/auth_gate.dart` - Route protection widget

### External

- `go_router: ^13.0.0` - Declarative routing
- `flutter_riverpod: ^2.4.0` - State management untuk auth listener

---

## Notes & Improvements

### Known Limitations

1. **No Route Transition Animations**: Belum customize page transitions. Default material transitions digunakan.

2. **Limited Deep Link Testing**: Deep link testing hanya tested di development environment.

3. **No Route Analytics**: Belum implement route tracking untuk analytics.

### Future Improvements

1. **Custom Transitions**: Add custom page transitions (slide, fade, shared element)

2. **Route Analytics**: Track page views untuk analytics

3. **Better Error Handling**: Custom error screens untuk 404 dan navigation errors

4. **Route Preloading**: Preload data untuk routes yang sering diakses

5. **Navigation Drawer**: Add drawer navigation untuk tablet/desktop layouts

6. **Hero Animations**: Implement hero animations untuk smooth transitions

---

**Document Status**: Active  
**Last Updated**: January 2025  
**Maintained By**: Dev3 (Mobile Development Team)
