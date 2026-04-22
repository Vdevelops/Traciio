# Core - API Integration

## CRM Healthcare Mobile App - Flutter

**Module**: Core Infrastructure  
**Sprint**: Sprint 0  
**Version**: 1.0  
**Status**: ✅ **Completed**  
**Last Updated**: January 2025

---

## Table of Contents

1. [Ringkasan Fitur](#ringkasan-fitur)
2. [Fitur Utama](#fitur-utama)
3. [Business Rules](#business-rules)
4. [Keputusan Teknis & Trade-offs](#keputusan-teknis--trade-offs)
5. [Struktur Folder](#struktur-folder)
6. [API Endpoints](#api-endpoints)
7. [Configuration](#configuration)
8. [Response Parsing](#response-parsing)
9. [Cara Test Manual](#cara-test-manual)
10. [Dependencies](#dependencies)
11. [Notes & Improvements](#notes--improvements)

---

## Ringkasan Fitur

Sistem **API Integration** mobile app CRM Healthcare menggunakan **Dio** sebagai HTTP client dengan custom interceptors untuk authentication, logging, error handling, dan response parsing. Sistem ini mensupport REST API endpoints dari backend Go/Gin dengan flexible response parsing untuk handle berbagai format response.

### Goals

- **HTTP Client**: Robust HTTP client dengan timeout dan retry logic
- **Authentication**: Automatic token injection ke setiap request
- **Error Handling**: Centralized error handling via interceptors
- **Response Parsing**: Flexible parsing untuk web dan mobile endpoints
- **Offline Support**: Cache responses untuk offline access

---

## Fitur Utama

### 1. Dio HTTP Client

**Configuration**:

```dart
final dio = Dio(BaseOptions(
  baseUrl: Env.apiBaseUrl,
  connectTimeout: Duration(seconds: 30),
  receiveTimeout: Duration(seconds: 30),
  contentType: 'application/json',
));
```

### 2. Interceptors

**Auth Interceptor**: Inject Bearer token ke setiap request
**Logging Interceptor**: Log requests dan responses (development only)
**Error Interceptor**: Handle errors dan transform ke AppError
**Retry Interceptor**: Auto-retry failed requests

### 3. Response Parsing

**Standard Response Format**:

```json
{
  "success": true,
  "data": {},
  "meta": {
    "total": 100,
    "page": 1,
    "per_page": 20
  },
  "timestamp": "2025-01-15T10:30:45+07:00",
  "request_id": "req_abc123xyz"
}
```

**Flexible Parsing**: Support untuk direct array response `[...]` atau object response `{ data: [...] }`

### 4. Offline Caching

**Cache Strategy**: Cache-first dengan background refresh
**Cache TTL**: 5-15 minutes untuk data yang berubah
**Cache Storage**: Hive untuk structured data

---

## Business Rules

### 1. Request Timeout Rules

| Operation Type | Connect Timeout | Receive Timeout | Retry |
| -------------- | --------------- | --------------- | ----- |
| GET (List)     | 30s             | 30s             | 1x    |
| GET (Detail)   | 15s             | 15s             | 2x    |
| POST           | 30s             | 60s             | 0x    |
| PUT/PATCH      | 30s             | 60s             | 0x    |
| DELETE         | 15s             | 15s             | 0x    |
| File Upload    | 60s             | 120s            | 1x    |

### 2. Authentication Rules

- **Token Injection**: Automatic inject Bearer token ke header `Authorization`
- **Token Refresh**: Auto-refresh token jika mendapat 401 response
- **Retry Queue**: Queue concurrent requests saat refreshing token

### 3. Error Response Rules

**Standard Error Format**:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error message (bilingual)",
    "details": {},
    "field_errors": []
  },
  "timestamp": "2025-01-15T10:30:45+07:00",
  "request_id": "req_abc123xyz"
}
```

**Error Handling Priority**:

1. Check HTTP status code
2. Parse error message dari response
3. Transform ke AppError
4. Show error ke user via ErrorHandler

### 4. Pagination Rules

**Default Pagination**:

- Page size: 20 items
- Max page size: 100 items
- Cursor-based untuk large datasets (jika available)

**Query Parameters**:

```
?page=1&limit=20&search=keyword&sort=created_at&order=desc
```

---

## Keputusan Teknis & Trade-offs

### Mengapa Dio, bukan http package?

**Keputusan**: Menggunakan Dio daripada http package.

**Alasan**:

- **Interceptors**: Native support untuk request/response interceptors
- **Error Handling**: Better error handling dengan DioException
- **Configuration**: Global configuration options (timeouts, headers)
- **Advanced Features**: Form data, file upload, request cancellation

**Trade-off**: Larger package size. **Mitigasi**: Acceptable untuk feature-rich app.

### Mengapa Interceptors?

**Keputusan**: Menggunakan interceptors untuk cross-cutting concerns.

**Alasan**:

- **Separation of Concerns**: Auth, logging, error handling terpisah dari business logic
- **Reusability**: Same interceptors untuk semua API calls
- **Maintainability**: Single place untuk update logic

### Flexible Response Parsing

**Keputusan**: Support multiple response formats (array vs object dengan data field).

**Alasan**:

- **Compatibility**: Backend mungkin return different formats
- **Flexibility**: Handle legacy endpoints dengan berbagai format
- **Safety**: Fallback parsing untuk robust error handling

---

## Struktur Folder

```
apps/mobile/lib/
├── core/
│   ├── network/
│   │   ├── api_client.dart              # Dio client configuration
│   │   ├── interceptors/
│   │   │   ├── auth_interceptor.dart    # Token injection
│   │   │   ├── logging_interceptor.dart # Request/response logging
│   │   │   ├── error_interceptor.dart   # Error handling
│   │   │   └── retry_interceptor.dart   # Auto-retry logic
│   │   ├── api_response.dart            # Response wrapper classes
│   │   └── api_exception.dart           # Custom exceptions
│   └── config/
│       └── env.dart                     # Environment configuration
├── features/
│   └── [feature]/
│       └── data/
│           └── [feature]_repository.dart # API calls per feature
```

---

## API Endpoints

### Authentication Endpoints

| Method | Endpoint                          | Description          | Auth |
| ------ | --------------------------------- | -------------------- | ---- |
| POST   | `/api/v1/auth/login`              | User login           | No   |
| POST   | `/api/v1/auth/refresh`            | Refresh access token | No   |
| POST   | `/api/v1/auth/logout`             | User logout          | Yes  |
| GET    | `/api/v1/auth/mobile/permissions` | Get user permissions | Yes  |

### Account Endpoints

| Method | Endpoint                        | Description          | Auth |
| ------ | ------------------------------- | -------------------- | ---- |
| GET    | `/api/v1/accounts`              | List accounts        | Yes  |
| GET    | `/api/v1/accounts/:id`          | Get account detail   | Yes  |
| GET    | `/api/v1/accounts/:id/contacts` | Get account contacts | Yes  |

### Contact Endpoints

| Method | Endpoint               | Description        | Auth |
| ------ | ---------------------- | ------------------ | ---- |
| GET    | `/api/v1/contacts`     | List contacts      | Yes  |
| GET    | `/api/v1/contacts/:id` | Get contact detail | Yes  |

### Visit Report Endpoints

| Method | Endpoint                              | Description             | Auth |
| ------ | ------------------------------------- | ----------------------- | ---- |
| GET    | `/api/v1/visit-reports`               | List visit reports      | Yes  |
| GET    | `/api/v1/visit-reports/:id`           | Get visit report detail | Yes  |
| POST   | `/api/v1/visit-reports`               | Create visit report     | Yes  |
| PUT    | `/api/v1/visit-reports/:id`           | Update visit report     | Yes  |
| POST   | `/api/v1/visit-reports/:id/check-in`  | Check-in                | Yes  |
| POST   | `/api/v1/visit-reports/:id/check-out` | Check-out               | Yes  |
| POST   | `/api/v1/visit-reports/:id/photos`    | Upload photo            | Yes  |

### Task Endpoints

| Method | Endpoint                     | Description     | Auth |
| ------ | ---------------------------- | --------------- | ---- |
| GET    | `/api/v1/tasks`              | List tasks      | Yes  |
| GET    | `/api/v1/tasks/:id`          | Get task detail | Yes  |
| POST   | `/api/v1/tasks`              | Create task     | Yes  |
| PUT    | `/api/v1/tasks/:id`          | Update task     | Yes  |
| POST   | `/api/v1/tasks/:id/complete` | Complete task   | Yes  |

### Dashboard Endpoints (Mobile-Specific)

| Method | Endpoint                            | Description        | Auth |
| ------ | ----------------------------------- | ------------------ | ---- |
| GET    | `/api/v1/dashboard/mobile/overview` | Dashboard overview | Yes  |
| GET    | `/api/v1/dashboard/mobile/visits`   | Recent visits      | Yes  |
| GET    | `/api/v1/dashboard/mobile/tasks`    | Upcoming tasks     | Yes  |

---

## Configuration

### API Client Setup

**File**: `core/network/api_client.dart`

```dart
class ApiClient {
  late final Dio _dio;

  ApiClient() {
    _dio = Dio(BaseOptions(
      baseUrl: Env.apiBaseUrl,
      connectTimeout: const Duration(seconds: 30),
      receiveTimeout: const Duration(seconds: 30),
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
    ));

    // Add interceptors
    _dio.interceptors.addAll([
      AuthInterceptor(),
      LoggingInterceptor(),
      ErrorInterceptor(),
      RetryInterceptor(dio: _dio),
    ]);
  }

  Future<Response> get(String path, {
    Map<String, dynamic>? queryParameters,
    Options? options,
  }) async {
    return _dio.get(path, queryParameters: queryParameters, options: options);
  }

  Future<Response> post(String path, {
    dynamic data,
    Map<String, dynamic>? queryParameters,
  }) async {
    return _dio.post(path, data: data, queryParameters: queryParameters);
  }

  Future<Response> put(String path, {
    dynamic data,
  }) async {
    return _dio.put(path, data: data);
  }

  Future<Response> delete(String path) async {
    return _dio.delete(path);
  }

  Future<Response> uploadFile(String path, File file, {
    String fieldName = 'file',
    Map<String, dynamic>? extraData,
    ProgressCallback? onSendProgress,
  }) async {
    final formData = FormData.fromMap({
      fieldName: await MultipartFile.fromFile(file.path),
      if (extraData != null) ...extraData,
    });

    return _dio.post(
      path,
      data: formData,
      onSendProgress: onSendProgress,
    );
  }
}
```

### Auth Interceptor

**File**: `core/network/interceptors/auth_interceptor.dart`

```dart
class AuthInterceptor extends Interceptor {
  final TokenStorage _tokenStorage;
  bool _isRefreshing = false;
  final _pendingRequests = <Completer<void>>[];

  AuthInterceptor(this._tokenStorage);

  @override
  Future<void> onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    // Skip auth untuk public endpoints
    if (_isPublicEndpoint(options.path)) {
      handler.next(options);
      return;
    }

    // Get access token
    final token = await _tokenStorage.getAccessToken();
    if (token != null) {
      options.headers['Authorization'] = 'Bearer $token';
    }

    handler.next(options);
  }

  @override
  Future<void> onError(
    DioException err,
    ErrorInterceptorHandler handler,
  ) async {
    if (err.response?.statusCode == 401) {
      final options = err.requestOptions;

      // Wait jika sedang refresh token
      if (_isRefreshing) {
        final completer = Completer<void>();
        _pendingRequests.add(completer);
        await completer.future;

        // Retry dengan token baru
        final newToken = await _tokenStorage.getAccessToken();
        options.headers['Authorization'] = 'Bearer $newToken';

        try {
          final response = await Dio().fetch(options);
          handler.resolve(response);
          return;
        } catch (e) {
          handler.next(err);
          return;
        }
      }

      // Refresh token
      _isRefreshing = true;

      try {
        final refreshed = await _refreshToken();
        if (refreshed) {
          // Retry original request
          final newToken = await _tokenStorage.getAccessToken();
          options.headers['Authorization'] = 'Bearer $newToken';

          final response = await Dio().fetch(options);
          handler.resolve(response);

          // Complete pending requests
          for (final completer in _pendingRequests) {
            completer.complete();
          }
          _pendingRequests.clear();
        } else {
          // Refresh failed, logout user
          _handleTokenRefreshFailure();
          handler.next(err);
        }
      } finally {
        _isRefreshing = false;
      }
    } else {
      handler.next(err);
    }
  }

  bool _isPublicEndpoint(String path) {
    final publicPaths = ['/auth/login', '/auth/refresh'];
    return publicPaths.any((p) => path.contains(p));
  }

  Future<bool> _refreshToken() async {
    try {
      final refreshToken = await _tokenStorage.getRefreshToken();
      if (refreshToken == null) return false;

      final response = await Dio().post(
        '${Env.apiBaseUrl}/auth/refresh',
        data: {'refresh_token': refreshToken},
      );

      if (response.statusCode == 200) {
        final data = response.data['data'];
        await _tokenStorage.saveTokens(
          data['access_token'],
          data['refresh_token'],
        );
        return true;
      }
      return false;
    } catch (e) {
      return false;
    }
  }

  void _handleTokenRefreshFailure() {
    // Clear tokens dan trigger logout
    _tokenStorage.clearTokens();
    // Emit event untuk logout
  }
}
```

### Logging Interceptor

**File**: `core/network/interceptors/logging_interceptor.dart`

```dart
class LoggingInterceptor extends Interceptor {
  final bool enabled;

  LoggingInterceptor({this.enabled = true});

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) {
    if (!enabled) {
      handler.next(options);
      return;
    }

    print('┌─────────────────────────────────────────────────────────────');
    print('│ REQUEST: ${options.method} ${options.path}');
    print('│ Headers: ${options.headers}');
    print('│ Query: ${options.queryParameters}');
    if (options.data != null) {
      print('│ Body: ${options.data}');
    }
    print('└─────────────────────────────────────────────────────────────');

    handler.next(options);
  }

  @override
  void onResponse(Response response, ResponseInterceptorHandler handler) {
    if (!enabled) {
      handler.next(response);
      return;
    }

    print('┌─────────────────────────────────────────────────────────────');
    print('│ RESPONSE: ${response.statusCode} ${response.requestOptions.path}');
    print('│ Data: ${response.data}');
    print('└─────────────────────────────────────────────────────────────');

    handler.next(response);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    if (!enabled) {
      handler.next(err);
      return;
    }

    print('┌─────────────────────────────────────────────────────────────');
    print('│ ERROR: ${err.requestOptions.path}');
    print('│ Status: ${err.response?.statusCode}');
    print('│ Message: ${err.message}');
    print('└─────────────────────────────────────────────────────────────');

    handler.next(err);
  }
}
```

---

## Response Parsing

### Flexible Response Parser

**File**: `core/network/api_response.dart`

```dart
class ApiResponse<T> {
  final bool success;
  final T? data;
  final ApiError? error;
  final Map<String, dynamic>? meta;
  final String? timestamp;
  final String? requestId;

  ApiResponse({
    required this.success,
    this.data,
    this.error,
    this.meta,
    this.timestamp,
    this.requestId,
  });

  factory ApiResponse.fromJson(
    Map<String, dynamic> json,
    T Function(dynamic) fromJsonT,
  ) {
    return ApiResponse(
      success: json['success'] ?? false,
      data: json['data'] != null ? fromJsonT(json['data']) : null,
      error: json['error'] != null ? ApiError.fromJson(json['error']) : null,
      meta: json['meta'],
      timestamp: json['timestamp'],
      requestId: json['request_id'],
    );
  }

  // Parse list dengan flexible format
  static List<T> parseList<T>(
    dynamic json,
    T Function(dynamic) fromJsonT,
  ) {
    // Format 1: Direct array
    if (json is List) {
      return json.map(fromJsonT).toList();
    }

    // Format 2: Object dengan data.items
    if (json is Map<String, dynamic>) {
      if (json['data'] is List) {
        return (json['data'] as List).map(fromJsonT).toList();
      }
      if (json['data'] is Map && json['data']['items'] is List) {
        return (json['data']['items'] as List).map(fromJsonT).toList();
      }
    }

    return [];
  }
}

class ApiError {
  final String code;
  final String message;
  final Map<String, dynamic>? details;
  final List<FieldError>? fieldErrors;

  ApiError({
    required this.code,
    required this.message,
    this.details,
    this.fieldErrors,
  });

  factory ApiError.fromJson(Map<String, dynamic> json) {
    return ApiError(
      code: json['code'] ?? 'UNKNOWN_ERROR',
      message: json['message'] ?? 'Unknown error occurred',
      details: json['details'],
      fieldErrors: json['field_errors'] != null
          ? (json['field_errors'] as List)
              .map((e) => FieldError.fromJson(e))
              .toList()
          : null,
    );
  }
}

class FieldError {
  final String field;
  final String message;

  FieldError({required this.field, required this.message});

  factory FieldError.fromJson(Map<String, dynamic> json) {
    return FieldError(
      field: json['field'] ?? '',
      message: json['message'] ?? '',
    );
  }
}
```

---

## Cara Test Manual

### Test API Calls

1. **Successful GET**:
   - Buka accounts list
   - Verifikasi: Data loaded successfully
   - Check logs: Request dan response logged

2. **POST dengan Body**:
   - Create visit report
   - Verifikasi: Request body terkirim dengan benar
   - Verifikasi: Response parsed correctly

3. **Token Injection**:
   - Login dan periksa header di logs
   - Verifikasi: Authorization header ada dengan Bearer token

4. **Token Refresh**:
   - Wait sampai token expired
   - Trigger API call
   - Verifikasi: Token refresh request terkirim
   - Verifikasi: Original request retry dengan token baru

5. **Error Handling**:
   - Trigger 404 error
   - Verifikasi: Error parsed correctly
   - Verifikasi: Error message displayed ke user

6. **Timeout**:
   - Simulate slow network
   - Verifikasi: Timeout error muncul setelah 30s

### Test File Upload

1. **Photo Upload**:
   - Upload photo di visit report
   - Verifikasi: File terkirim dengan multipart/form-data
   - Verifikasi: Progress callback dipanggil

---

## Dependencies

### Internal

- `core/storage/token_storage.dart` - Token management
- `core/errors/error_handler.dart` - Error handling
- `core/config/env.dart` - Environment configuration

### External

- `dio: ^5.0.0` - HTTP client
- `path_provider: ^2.1.0` - File paths untuk uploads
- `mime: ^1.0.0` - MIME type detection untuk uploads

---

## Notes & Improvements

### Known Limitations

1. **No Request Cancellation**: Belum implement request cancellation untuk long-running requests.

2. **No Request Queue**: Belum implement request queue untuk handle offline scenarios.

3. **Limited Caching**: Cache strategy masih basic, belum sophisticated caching.

### Future Improvements

1. **Request Cancellation**: Support untuk cancel pending requests

2. **Request Queue**: Queue requests saat offline dan sync saat online

3. **Smart Caching**: Implement sophisticated cache dengan TTL dan cache invalidation

4. **GraphQL Support**: Siapkan struktur untuk future GraphQL migration

5. **API Versioning**: Handle API versioning di client side

6. **Request Deduplication**: Deduplicate identical concurrent requests

---

**Document Status**: Active  
**Last Updated**: January 2025  
**Maintained By**: Dev3 (Mobile Development Team)
