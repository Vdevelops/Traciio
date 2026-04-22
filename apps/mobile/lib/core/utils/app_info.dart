import 'package:package_info_plus/package_info_plus.dart';

/// Utility class to get app information (version, build number, etc.)
class AppInfo {
  static PackageInfo? _packageInfo;
  static bool _initialized = false;

  /// Initialize app info (should be called once at app startup)
  static Future<void> initialize() async {
    if (!_initialized) {
      _packageInfo = await PackageInfo.fromPlatform();
      _initialized = true;
    }
  }

  /// Get app version (e.g., "1.0.0")
  static String get version => _packageInfo?.version ?? '1.0.0';

  /// Get build number (e.g., "1")
  static String get buildNumber => _packageInfo?.buildNumber ?? '1';

  /// Get full version string (e.g., "1.0.0+1")
  static String get versionString => '$version+$buildNumber';

  /// Get app name
  static String get appName => _packageInfo?.appName ?? 'CRM Healthcare';

  /// Get package name
  static String get packageName => _packageInfo?.packageName ?? 'com.example.mobile';

  /// Check if app info is initialized
  static bool get isInitialized => _initialized;
}

