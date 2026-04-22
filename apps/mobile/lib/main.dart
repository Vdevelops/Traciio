import 'dart:async';

import 'package:app_links/app_links.dart';
import 'package:flutter/material.dart';
import 'package:flutter_dotenv/flutter_dotenv.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'core/l10n/app_localizations.dart';
import 'core/l10n/locale_provider.dart';
import 'core/network/api_client.dart';
import 'core/routing/app_router.dart';
import 'core/storage/hive_storage.dart';
import 'core/storage/offline_storage.dart';
import 'core/theme/app_theme.dart';
import 'features/route_optimization/data/route_optimization_cache.dart';
import 'core/theme/theme_provider.dart';
import 'core/utils/app_info.dart';
import 'features/auth/application/auth_provider.dart';
import 'features/google_calendar/application/google_calendar_provider.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // Load environment variables from .env file
  // File .env harus ada di root project dan di .gitignore
  await dotenv.load(fileName: ".env");

  // Initialize app info
  await AppInfo.initialize();
  // Initialize Hive storage for offline support
  await HiveStorage.init();
  await OfflineStorage.init();
  // Initialize route optimization cache
  await RouteOptimizationCache.init();

  runApp(const ProviderScope(child: MyApp()));
}

class MyApp extends ConsumerStatefulWidget {
  const MyApp({super.key});

  @override
  ConsumerState<MyApp> createState() => _MyAppState();
}

class _MyAppState extends ConsumerState<MyApp> {
  /// Global navigator key for navigating from non-widget contexts (e.g., API interceptor logout)
  static final navigatorKey = GlobalKey<NavigatorState>();

  /// Global scaffold messenger key for showing snackbars from outside widget tree
  static final scaffoldMessengerKey = GlobalKey<ScaffoldMessengerState>();

  AppLinks? _appLinks;
  StreamSubscription<Uri?>? _linkSubscription;

  /// Store pending deep link to handle after first frame
  Uri? _pendingDeepLink;

  @override
  void initState() {
    super.initState();
    _initDeepLinks();
  }

  @override
  void dispose() {
    _linkSubscription?.cancel();
    super.dispose();
  }

  /// Initialize deep link handling
  Future<void> _initDeepLinks() async {
    try {
      _appLinks = AppLinks();

      // Handle app being opened via deep link (cold start)
      // Store it and process after first frame is built
      final uri = await _appLinks!.getInitialLink();
      if (uri != null) {
        _pendingDeepLink = uri;
        // Process after first frame when ScaffoldMessenger is available
        WidgetsBinding.instance.addPostFrameCallback((_) {
          if (_pendingDeepLink != null) {
            _handleDeepLink(_pendingDeepLink!);
            _pendingDeepLink = null;
          }
        });
      }

      // Listen for deep links while app is running (warm start)
      _linkSubscription = _appLinks!.uriLinkStream.listen(
        (uri) {
          _handleDeepLink(uri);
        },
        onError: (err) {
          debugPrint('Deep link error: $err');
        },
      );
    } catch (e) {
      debugPrint('AppLinks initialization error: $e');
      // Silently fail if app_links plugin isn't available
      // This can happen on some platforms or if plugin registration fails
    }
  }

  /// Handle incoming deep links
  void _handleDeepLink(Uri uri) {
    debugPrint('Handling deep link: $uri');

    // Handle Google Calendar OAuth callback
    if (uri.scheme == 'crmhealth' &&
        uri.host == 'google-calendar' &&
        uri.path == '/callback') {
      _handleGoogleCalendarCallback(uri);
    }
  }

  /// Handle Google Calendar OAuth callback
  /// HTTPS Web Redirect + Server Forward (Option 1)
  /// Backend already exchanged code and stored token, just need to refresh status
  void _handleGoogleCalendarCallback(Uri uri) async {
    final success = uri.queryParameters['success'] == 'true';
    final error = uri.queryParameters['error'];

    if (error != null) {
      // Show error message using scaffoldMessengerKey
      scaffoldMessengerKey.currentState?.showSnackBar(
        SnackBar(content: Text('Error: $error'), backgroundColor: Colors.red),
      );
      return;
    }

    if (success) {
      // Show success message first
      scaffoldMessengerKey.currentState?.showSnackBar(
        const SnackBar(
          content: Text('Google Calendar connected successfully!'),
          backgroundColor: Colors.green,
        ),
      );

      // Wait a moment for backend to finish storing token, then invalidate and refresh
      await Future.delayed(const Duration(seconds: 2));
      ref.invalidate(googleCalendarNotifierProvider);
      ref.invalidate(googleCalendarStatusProvider);
    } else {
      // Handle failure case
      scaffoldMessengerKey.currentState?.showSnackBar(
        const SnackBar(
          content: Text('Failed to connect Google Calendar'),
          backgroundColor: Colors.red,
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final themeMode = ref.watch(themeModeProvider);
    final locale = ref.watch(localeProvider);

    // Setup API client interceptors
    ApiClient.setupInterceptors(
      onRefreshToken: () async {
        final authNotifier = ref.read(authProvider.notifier);
        return await authNotifier.refreshToken();
      },
      onLogout: () async {
        final authNotifier = ref.read(authProvider.notifier);
        await authNotifier.logout();
        // Navigate to login and clear navigation stack
        navigatorKey.currentState?.pushNamedAndRemoveUntil(
          AppRoutes.login,
          (route) => false,
        );
      },
    );

    return MaterialApp(
      navigatorKey: navigatorKey,
      scaffoldMessengerKey: scaffoldMessengerKey,
      title: 'SalesView',
      theme: AppTheme.light,
      darkTheme: AppTheme.dark,
      themeMode: themeMode,
      locale: locale,
      supportedLocales: AppLocalizations.supportedLocales,
      localizationsDelegates: const [
        AppLocalizations.delegate,
        GlobalMaterialLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
      ],
      initialRoute: AppRouter.initialRoute,
      routes: AppRouter.routes,
      onGenerateRoute: AppRouter.onGenerateRoute,
    );
  }
}
