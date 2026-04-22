import 'package:flutter_dotenv/flutter_dotenv.dart';

/// Environment configuration class
///
/// Configuration values are loaded from .env file using flutter_dotenv.
///
/// Setup:
/// 1. Copy .env.example to .env
/// 2. Fill in your API_BASE_URL in .env file
/// 3. NEVER commit .env file to repository
///
/// Example .env content:
/// ```
/// API_BASE_URL=http://192.168.1.100:8080
/// ```
///
/// For production:
/// ```
/// API_BASE_URL=https://api.yourdomain.com
/// ```
class Env {
  const Env._();

  /// Base URL untuk backend API
  ///
  /// WAJIB di-set di file .env
  /// Jika tidak di-set, akan throw exception saat aplikasi berjalan
  static String get apiBaseUrl {
    final url = dotenv.env['API_BASE_URL'];

    if (url == null || url.isEmpty) {
      throw Exception(
        'API_BASE_URL tidak di-set di file .env.\n'
        'Mohon copy .env.example ke .env dan isi API_BASE_URL dengan endpoint API Anda.',
      );
    }

    return url;
  }

  /// Environment identifier (optional)
  static String get environment {
    return dotenv.env['ENVIRONMENT'] ?? 'development';
  }

  /// Check if running in production
  static bool get isProduction => environment == 'production';

  /// Check if running in development
  static bool get isDevelopment => environment == 'development';
}
