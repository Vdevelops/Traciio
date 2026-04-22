# Network Setup Guide

Complete guide for setting up Dio and Retrofit networking in Clean Architecture.

## Overview

Network setup includes:
- **Dio Provider**: Base HTTP client configuration
- **Interceptors**: Request/response processing
- **Network Exceptions**: Centralized error handling

## Dio Provider

### Basic Dio Provider

`lib/core/network/dio_provider.dart`:

```dart
import 'package:dio/dio.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../interceptors/logging_interceptor.dart';
import '../interceptors/auth_interceptor.dart';
import '../interceptors/error_interceptor.dart';

part 'dio_provider.g.dart';

@riverpod
Dio dio(Ref ref) {
  final options = BaseOptions(
    baseUrl: String.fromEnvironment(
      'API_BASE_URL',
      defaultValue: 'https://api.example.com',
    ),
    connectTimeout: const Duration(seconds: 30),
    receiveTimeout: const Duration(seconds: 30),
    sendTimeout: const Duration(seconds: 30),
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    },
  );

  final dio = Dio(options);

  // Add interceptors in order
  dio.interceptors.addAll([
    ref.watch(loggingInterceptorProvider),
    ref.watch(authInterceptorProvider),
    ref.watch(errorInterceptorProvider),
  ]);

  return dio;
}
```

### Environment-specific Configuration

```dart
@riverpod
Dio dio(Ref ref) {
  final env = String.fromEnvironment('ENV', defaultValue: 'dev');

  final options = BaseOptions(
    baseUrl: _getBaseUrl(env),
    connectTimeout: _getTimeout(env),
    receiveTimeout: _getTimeout(env),
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
      if (env == 'dev') 'X-Debug-Mode': 'true',
    },
  );

  final dio = Dio(options);

  // Add different interceptors based on environment
  if (env == 'dev') {
    dio.interceptors.add(ref.watch(loggingInterceptorProvider));
  }

  dio.interceptors.addAll([
    ref.watch(authInterceptorProvider),
    ref.watch(errorInterceptorProvider),
  ]);

  return dio;
}

String _getBaseUrl(String env) {
  switch (env) {
    case 'prod':
      return 'https://api.production.com';
    case 'staging':
      return 'https://api.staging.com';
    default:
      return 'https://api.dev.com';
  }
}

Duration _getTimeout(String env) {
  switch (env) {
    case 'prod':
      return const Duration(seconds: 15);
    default:
      return const Duration(seconds: 30);
  }
}
```

## Interceptors

### Logging Interceptor

`lib/core/network/interceptors/logging_interceptor.dart`:

```dart
import 'package:dio/dio.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:flutter/foundation.dart';

part 'logging_interceptor.g.dart';

@riverpod
LoggingInterceptor loggingInterceptor(Ref ref) {
  return LoggingInterceptor();
}

class LoggingInterceptor extends Interceptor {
  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) {
    if (kDebugMode) {
      print('==> REQUEST');
      print('Method: ${options.method}');
      print('URL: ${options.uri}');
      print('Headers: ${options.headers}');
      if (options.data != null) {
        print('Body: ${options.data}');
      }
      print('Query Parameters: ${options.queryParameters}');
    }
    handler.next(options);
  }

  @override
  void onResponse(Response response, ResponseInterceptorHandler handler) {
    if (kDebugMode) {
      print('<== RESPONSE');
      print('Status Code: ${response.statusCode}');
      print('URL: ${response.requestOptions.uri}');
      print('Data: ${response.data}');
      print('Headers: ${response.headers}');
    }
    handler.next(response);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    if (kDebugMode) {
      print('==> ERROR');
      print('Type: ${err.type}');
      print('Message: ${err.message}');
      print('Response: ${err.response}');
    }
    handler.next(err);
  }
}
```

### Auth Interceptor

`lib/core/network/interceptors/auth_interceptor.dart`:

```dart
import 'package:dio/dio.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../../../features/auth/domain/providers/auth_provider.dart';

part 'auth_interceptor.g.dart';

@riverpod
AuthInterceptor authInterceptor(Ref ref) {
  return AuthInterceptor(ref);
}

class AuthInterceptor extends Interceptor {
  final Ref ref;

  AuthInterceptor(this.ref);

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) async {
    // Add auth token if available
    final authToken = ref.read(authTokenProvider);

    if (authToken != null && authToken.accessToken.isNotEmpty) {
      options.headers['Authorization'] = 'Bearer ${authToken.accessToken}';
    }

    // Add refresh token header if needed
    if (authToken?.refreshToken != null) {
      options.headers['X-Refresh-Token'] = authToken?.refreshToken;
    }

    handler.next(options);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    // Handle 401 - token expired
    if (err.response?.statusCode == 401) {
      // Try to refresh token
      final refreshToken = ref.read(refreshTokenProvider.notifier);

      final refreshed = await refreshToken();

      if (refreshed) {
        // Retry the original request with new token
        final options = err.requestOptions;
        final newToken = ref.read(authTokenProvider)?.accessToken;

        if (newToken != null) {
          options.headers['Authorization'] = 'Bearer $newToken';

          try {
            final response = await Dio().fetch(options);
            handler.resolve(response);
            return;
          } catch (e) {
            // If retry fails, continue with error
          }
        }
      }

      // If refresh failed, clear auth and redirect to login
      ref.read(authTokenProvider.notifier).clear();
      // Navigate to login screen
    }

    handler.next(err);
  }
}
```

### Error Interceptor

`lib/core/network/interceptors/error_interceptor.dart`:

```dart
import 'package:dio/dio.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../../errors/network_exceptions.dart';

part 'error_interceptor.g.dart';

@riverpod
ErrorInterceptor errorInterceptor(Ref ref) {
  return ErrorInterceptor();
}

class ErrorInterceptor extends Interceptor {
  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    final exception = NetworkExceptions.fromDioError(err);

    // Log the error
    print('Network Error: ${exception.message}');

    // Show toast notification for certain errors
    if (err.type == DioExceptionType.connectionError ||
        err.type == DioExceptionType.connectionTimeout) {
      // Show connectivity error to user
      _showConnectivityError();
    }

    handler.next(err);
  }

  void _showConnectivityError() {
    // Implement toast/snackbar notification
  }
}
```

### Retry Interceptor

```dart
class RetryInterceptor extends Interceptor {
  final Dio dio;
  final int maxRetries;
  final RetryDecider retryDecider;

  RetryInterceptor({
    required this.dio,
    this.maxRetries = 3,
    RetryDecider? retryDecider,
  }) : retryDecider = retryDecider ?? _defaultRetryDecider;

  static bool _defaultRetryDecider(DioException error) {
    switch (error.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.sendTimeout:
      case DioExceptionType.receiveTimeout:
      case DioExceptionType.connectionError:
        return true;
      default:
        return false;
    }
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    if (_shouldRetry(err)) {
      final retries = err.requestOptions.extra['retry_count'] ?? 0;

      if (retries < maxRetries) {
        err.requestOptions.extra['retry_count'] = retries + 1;

        // Exponential backoff
        final delay = Duration(milliseconds: 1000 * (retries + 1));
        await Future.delayed(delay);

        try {
          final response = await dio.fetch(err.requestOptions);
          handler.resolve(response);
          return;
        } catch (e) {
          // If retry fails, continue with error
        }
      }
    }

    handler.next(err);
  }

  bool _shouldRetry(DioException error) {
    return retryDecider(error);
  }
}

typedef RetryDecider = bool Function(DioException error);
```

## Network Exceptions

### Comprehensive Network Exceptions

See `error_handling.md` for complete NetworkExceptions implementation.

Basic usage:

```dart
try {
  final response = await apiService.getData();
  return Right(response);
} on DioException catch (e) {
  final exception = NetworkExceptions.fromDioError(e);
  return Left(Failure.network(exception.message));
}
```

## API Endpoints Configuration

### Centralized Endpoint Constants

`lib/core/network/api_endpoints.dart`:

```dart
class ApiEndpoints {
  static const String v1 = '/v1';

  // Auth endpoints
  static const String login = '$v1/auth/login';
  static const String register = '$v1/auth/register';
  static const String logout = '$v1/auth/logout';
  static const String refreshToken = '$v1/auth/refresh';
  static const String forgotPassword = '$v1/auth/forgot-password';

  // User endpoints
  static const String users = '$v1/users';
  static String userById(String id) => '$v1/users/$id';
  static String userAvatar(String id) => '$v1/users/$id/avatar';

  // Product endpoints
  static const String products = '$v1/products';
  static String productById(String id) => '$v1/products/$id';
  static String productReviews(String id) => '$v1/products/$id/reviews';

  // Order endpoints
  static const String orders = '$v1/orders';
  static String orderById(String id) => '$v1/orders/$id';
}
```

### Dynamic Endpoints

```dart
class ApiEndpoints {
  static String search(String query, {int page = 1}) {
    return '/v1/search?q=$query&page=$page';
  }

  static String userPosts(String userId, {String sort = 'latest'}) {
    return '/v1/users/$userId/posts?sort=$sort';
  }

  static String categoryProducts(String categoryId, Map<String, dynamic> filters) {
    final queryString = filters.entries
        .map((e) => '${e.key}=${e.value}')
        .join('&');
    return '/v1/categories/$categoryId/products?$queryString';
  }
}
```

## Environment Configuration

### Running with Environment Variables

```bash
# Development
flutter run

# Production
flutter run --dart-define=ENV=prod --dart-define=API_BASE_URL=https://api.production.com

# Staging
flutter run --dart-define=ENV=staging --dart-define=API_BASE_URL=https://api.staging.com

# Custom API URL
flutter run --dart-define=API_BASE_URL=https://custom.api.com
```

### Environment Configuration File

`lib/core/config/app_config.dart`:

```dart
class AppConfig {
  static const String environment = String.fromEnvironment('ENV', defaultValue: 'dev');
  static const String apiBaseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'https://api.dev.com',
  );

  static bool get isDebug => environment == 'dev';
  static bool get isProduction => environment == 'prod';

  static Duration get connectTimeout {
    switch (environment) {
      case 'prod':
        return const Duration(seconds: 15);
      default:
        return const Duration(seconds: 30);
    }
  }
}
```

## Best Practices

1. **Use environment variables** for different environments
2. **Log requests in debug mode only**
3. **Handle token refresh in interceptor**
4. **Centralize endpoint constants**
5. **Use appropriate timeout values** based on environment
6. **Retry failed requests** for transient errors
7. **Don't log sensitive data** (tokens, passwords)
8. **Use interceptors** for cross-cutting concerns