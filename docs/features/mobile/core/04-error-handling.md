# Core - Error Handling

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
6. [Error Types & Handling](#error-types--handling)
7. [Configuration](#configuration)
8. [Usage Examples](#usage-examples)
9. [Cara Test Manual](#cara-test-manual)
10. [Dependencies](#dependencies)
11. [Notes & Improvements](#notes--improvements)

---

## Ringkasan Fitur

Sistem **Error Handling** mobile app CRM Healthcare menyediakan centralized error management dengan bilingual error messages (Bahasa Indonesia dan English), user-friendly error display, dan proper error logging untuk debugging. Sistem ini menangani berbagai jenis error: network errors, validation errors, server errors, dan authentication errors.

### Goals

- **Centralized Handling**: Single point untuk handle semua errors
- **Bilingual Messages**: Error messages dalam ID dan EN
- **User-Friendly**: Error messages yang mudah dipahami user
- **Debuggable**: Proper error logging untuk development
- **Graceful Degradation**: App tetap functional meskipun ada errors

---

## Fitur Utama

### 1. Error Types

**Network Errors**:

- Connection timeout
- No internet connection
- Server unreachable
- SSL/TLS errors

**Validation Errors**:

- Form validation errors
- Missing required fields
- Invalid data format

**Server Errors**:

- 500 Internal Server Error
- 502 Bad Gateway
- 503 Service Unavailable

**Authentication Errors**:

- Invalid credentials
- Token expired
- Unauthorized access
- Session expired

**Business Logic Errors**:

- Duplicate data
- Resource not found
- Constraint violations

### 2. Bilingual Error Messages

Semua error messages support Bahasa Indonesia dan English:

```dart
class ErrorMessages {
  static const Map<String, Map<String, String>> _messages = {
    'INVALID_CREDENTIALS': {
      'id': 'Email atau password salah',
      'en': 'Invalid email or password',
    },
    'NETWORK_ERROR': {
      'id': 'Tidak dapat terhubung ke server. Periksa koneksi internet Anda.',
      'en': 'Cannot connect to server. Please check your internet connection.',
    },
    'TIMEOUT_ERROR': {
      'id': 'Koneksi timeout. Silakan coba lagi.',
      'en': 'Connection timeout. Please try again.',
    },
    // ... more messages
  };
}
```

### 3. Error Display Patterns

**Snackbar**: Untuk transient errors (toast-like)
**Dialog**: Untuk errors yang memerlukan user action
**Inline**: Untuk form validation errors
**Full Screen**: Untuk critical errors (no data, no connection)

### 4. Error Interceptor

Dio interceptor untuk automatic error handling:

```dart
class ErrorInterceptor extends Interceptor {
  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    final error = _parseError(err);
    _logError(error);
    handler.next(err);
  }
}
```

---

## Business Rules

### 1. Error Priority

1. **Critical**: Authentication errors (force logout)
2. **High**: Network errors (retry suggestion)
3. **Medium**: Validation errors (inline display)
4. **Low**: Non-blocking errors (silent logging)

### 2. Error Display Rules

| Error Type     | Display Pattern  | Duration          | Action Required   |
| -------------- | ---------------- | ----------------- | ----------------- |
| Network        | Snackbar + Retry | 5 seconds         | Optional retry    |
| Validation     | Inline           | Persistent        | Fix input         |
| Authentication | Dialog           | Until dismissed   | Re-login          |
| Server         | Snackbar         | 5 seconds         | Contact support   |
| Business       | Snackbar/Dialog  | Context-dependent | Context-dependent |

### 3. Error Logging Rules

**Development Mode**:

- Log all errors ke console
- Include stack trace
- Show detailed error information

**Production Mode**:

- Log only critical errors
- Send ke error tracking service (Sentry/Firebase Crashlytics)
- Sanitize sensitive information

### 4. Retry Logic

**Automatic Retry**:

- Network errors: Retry 3x dengan exponential backoff
- Timeout errors: Retry immediately 1x
- Server errors: No automatic retry

**User-Initiated Retry**:

- Provide retry button untuk network errors
- Pull-to-refresh untuk list screens
- Tap to retry untuk failed operations

---

## Keputusan Teknis & Trade-offs

### Mengapa Centralized Error Handler?

**Keputusan**: Menggunakan centralized error handler daripada handle errors di setiap screen.

**Alasan**:

- **Consistency**: Consistent error handling di seluruh app
- **Maintainability**: Single place untuk update error handling logic
- **Reusability**: Error handling logic bisa reused
- **Testing**: Mudah test error scenarios

**Trade-off**: Less flexibility untuk screen-specific error handling. **Mitigasi**: Provide hooks untuk custom error handling per screen.

### Mengapa Bilingual Messages?

**Keputusan**: Support Bahasa Indonesia dan English.

**Alasan**:

- **User Base**: Target users di Indonesia
- **Accessibility**: Users yang lebih nyaman dengan Bahasa Indonesia
- **Professional**: Show respect untuk local language

**Trade-off**: Maintenance overhead untuk translate messages. **Mitigasi**: Centralized message repository, easy to update.

### Mengapa Dio Interceptor?

**Keputusan**: Handle errors via Dio interceptor.

**Alasan**:

- **Automatic**: Handle errors tanpa explicit try-catch di setiap call
- **Centralized**: Single place untuk handle API errors
- **Consistent**: Consistent error format

**Trade-off**: Less visibility untuk specific error handling. **Mitigasi**: Allow custom error handling via callbacks.

---

## Struktur Folder

```
apps/mobile/lib/
├── core/
│   ├── errors/
│   │   ├── app_error.dart              # Base error class
│   │   ├── error_handler.dart          # Centralized error handler
│   │   ├── error_messages.dart         # Bilingual error messages
│   │   └── error_logger.dart           # Error logging service
│   ├── network/
│   │   ├── api_client.dart             # Dio client setup
│   │   └── error_interceptor.dart      # Dio error interceptor
│   └── widgets/
│       ├── error_widget.dart           # Reusable error display
│       ├── error_snackbar.dart         # Error snackbar component
│       └── error_dialog.dart           # Error dialog component
└── features/
    └── [feature]/
        └── presentation/
            └── screens/
                └── [screen].dart       # Screen-specific error handling
```

---

## Error Types & Handling

### 1. Network Errors

**Error Class**:

```dart
class NetworkError extends AppError {
  final String? url;

  NetworkError({this.url, super.originalError})
      : super(
          code: 'NETWORK_ERROR',
          messageId: 'Tidak dapat terhubung ke server',
          messageEn: 'Cannot connect to server',
        );
}
```

**Handling**:

```dart
void handleNetworkError(NetworkError error) {
  showSnackbar(
    message: error.localizedMessage,
    action: SnackBarAction(
      label: 'Retry',
      onPressed: () => retryLastRequest(),
    ),
  );
}
```

### 2. Validation Errors

**Error Class**:

```dart
class ValidationError extends AppError {
  final Map<String, List<String>> fieldErrors;

  ValidationError({required this.fieldErrors})
      : super(
          code: 'VALIDATION_ERROR',
          messageId: 'Data tidak valid',
          messageEn: 'Invalid data',
        );
}
```

**Handling**:

```dart
void handleValidationError(ValidationError error) {
  // Display inline errors
  for (final entry in error.fieldErrors.entries) {
    final fieldName = entry.key;
    final errors = entry.value;
    formKey.currentState?.fields[fieldName]?.invalidate(errors.first);
  }
}
```

### 3. Authentication Errors

**Error Class**:

```dart
class AuthenticationError extends AppError {
  final bool shouldLogout;

  AuthenticationError({
    required super.code,
    required super.messageId,
    required super.messageEn,
    this.shouldLogout = false,
  });
}
```

**Handling**:

```dart
void handleAuthenticationError(AuthenticationError error) {
  if (error.shouldLogout) {
    showDialog(
      context: context,
      builder: (_) => AlertDialog(
        title: Text('Session Expired'),
        content: Text(error.localizedMessage),
        actions: [
          TextButton(
            onPressed: () {
              Navigator.pop(context);
              logout();
            },
            child: Text('Login Again'),
          ),
        ],
      ),
    );
  }
}
```

### 4. Server Errors

**Error Class**:

```dart
class ServerError extends AppError {
  final int statusCode;

  ServerError({
    required this.statusCode,
    super.originalError,
  }) : super(
    code: 'SERVER_ERROR_$statusCode',
    messageId: 'Terjadi kesalahan pada server',
    messageEn: 'Server error occurred',
  );
}
```

**Handling**:

```dart
void handleServerError(ServerError error) {
  showSnackbar(
    message: '${error.localizedMessage} (${error.statusCode})',
    duration: Duration(seconds: 5),
  );
}
```

---

## Configuration

### Error Handler Setup

**File**: `core/errors/error_handler.dart`

```dart
class ErrorHandler {
  static void initialize() {
    // Setup error zone untuk catch async errors
    FlutterError.onError = (FlutterErrorDetails details) {
      _handleFlutterError(details);
    };
  }

  static void handleError(dynamic error, StackTrace? stackTrace) {
    final appError = _convertToAppError(error);

    // Log error
    ErrorLogger.log(appError, stackTrace);

    // Show error to user
    _showError(appError);
  }

  static AppError _convertToAppError(dynamic error) {
    if (error is AppError) return error;

    if (error is DioException) {
      return _handleDioError(error);
    }

    if (error is SocketException) {
      return NetworkError(originalError: error);
    }

    if (error is TimeoutException) {
      return TimeoutError(originalError: error);
    }

    return UnknownError(originalError: error);
  }

  static AppError _handleDioError(DioException error) {
    switch (error.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.sendTimeout:
      case DioExceptionType.receiveTimeout:
        return TimeoutError(originalError: error);

      case DioExceptionType.connectionError:
      case DioExceptionType.unknown:
        return NetworkError(url: error.requestOptions.path, originalError: error);

      case DioExceptionType.badResponse:
        final statusCode = error.response?.statusCode;
        if (statusCode == 401) {
          return AuthenticationError(
            code: 'UNAUTHORIZED',
            messageId: 'Sesi telah berakhir',
            messageEn: 'Session expired',
            shouldLogout: true,
          );
        }
        if (statusCode == 403) {
          return AuthenticationError(
            code: 'FORBIDDEN',
            messageId: 'Akses ditolak',
            messageEn: 'Access denied',
          );
        }
        if (statusCode == 422) {
          return ValidationError(
            fieldErrors: error.response?.data?['errors'] ?? {},
          );
        }
        return ServerError(statusCode: statusCode ?? 500, originalError: error);

      default:
        return UnknownError(originalError: error);
    }
  }

  static void _showError(AppError error) {
    final context = navigatorKey.currentContext;
    if (context == null) return;

    switch (error.displayType) {
      case ErrorDisplayType.snackbar:
        showErrorSnackbar(context, error);
        break;
      case ErrorDisplayType.dialog:
        showErrorDialog(context, error);
        break;
      case ErrorDisplayType.none:
        // Silent error
        break;
    }
  }
}
```

### Error Messages

**File**: `core/errors/error_messages.dart`

```dart
class ErrorMessages {
  static const Map<String, Map<String, String>> _messages = {
    // Network Errors
    'NETWORK_ERROR': {
      'id': 'Tidak dapat terhubung ke server. Periksa koneksi internet Anda.',
      'en': 'Cannot connect to server. Please check your internet connection.',
    },
    'TIMEOUT_ERROR': {
      'id': 'Koneksi timeout. Silakan coba lagi.',
      'en': 'Connection timeout. Please try again.',
    },

    // Authentication Errors
    'INVALID_CREDENTIALS': {
      'id': 'Email atau password salah',
      'en': 'Invalid email or password',
    },
    'SESSION_EXPIRED': {
      'id': 'Sesi Anda telah berakhir. Silakan login kembali.',
      'en': 'Your session has expired. Please login again.',
    },
    'UNAUTHORIZED': {
      'id': 'Anda tidak memiliki akses untuk melakukan ini',
      'en': 'You do not have permission to do this',
    },

    // Validation Errors
    'VALIDATION_ERROR': {
      'id': 'Data yang dimasukkan tidak valid',
      'en': 'The entered data is invalid',
    },
    'REQUIRED_FIELD': {
      'id': 'Field ini wajib diisi',
      'en': 'This field is required',
    },

    // Server Errors
    'SERVER_ERROR': {
      'id': 'Terjadi kesalahan pada server. Silakan coba lagi nanti.',
      'en': 'Server error occurred. Please try again later.',
    },
    'NOT_FOUND': {
      'id': 'Data tidak ditemukan',
      'en': 'Data not found',
    },

    // Business Logic Errors
    'DUPLICATE_DATA': {
      'id': 'Data sudah ada dalam sistem',
      'en': 'Data already exists in the system',
    },

    // Unknown Errors
    'UNKNOWN_ERROR': {
      'id': 'Terjadi kesalahan. Silakan coba lagi.',
      'en': 'An error occurred. Please try again.',
    },
  };

  static String getMessage(String code, {String locale = 'id'}) {
    return _messages[code]?[locale] ?? _messages['UNKNOWN_ERROR']![locale]!;
  }
}
```

---

## Usage Examples

### 1. Handle Error di Screen

```dart
class AccountsScreen extends ConsumerWidget {
  const AccountsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(accountsProvider);

    return Scaffold(
      body: state.when(
        loading: () => const LoadingIndicator(),
        error: (error) => ErrorWidget(
          error: error,
          onRetry: () => ref.read(accountsProvider.notifier).loadAccounts(),
        ),
        data: (accounts) => AccountsList(accounts: accounts),
      ),
    );
  }
}
```

### 2. Handle Error dengan Snackbar

```dart
void showErrorSnackbar(BuildContext context, AppError error) {
  ScaffoldMessenger.of(context).showSnackBar(
    SnackBar(
      content: Text(error.localizedMessage),
      backgroundColor: Colors.red,
      behavior: SnackBarBehavior.floating,
      duration: const Duration(seconds: 5),
      action: error.isRetryable
          ? SnackBarAction(
              label: 'Retry',
              textColor: Colors.white,
              onPressed: () => retryOperation(),
            )
          : null,
    ),
  );
}
```

### 3. Handle Error dengan Dialog

```dart
void showErrorDialog(BuildContext context, AppError error) {
  showDialog(
    context: context,
    barrierDismissible: false,
    builder: (context) => AlertDialog(
      title: Row(
        children: [
          Icon(Icons.error_outline, color: Colors.red),
          SizedBox(width: 8),
          Text('Error'),
        ],
      ),
      content: Text(error.localizedMessage),
      actions: [
        TextButton(
          onPressed: () => Navigator.pop(context),
          child: Text('OK'),
        ),
        if (error.isRetryable)
          ElevatedButton(
            onPressed: () {
              Navigator.pop(context);
              retryOperation();
            },
            child: Text('Retry'),
          ),
      ],
    ),
  );
}
```

### 4. Try-Catch dengan Error Handler

```dart
Future<void> performOperation() async {
  try {
    await apiClient.post('/endpoint', data: data);
  } catch (e, stackTrace) {
    ErrorHandler.handleError(e, stackTrace);
  }
}
```

---

## Cara Test Manual

### Test Error Scenarios

1. **Network Error**:
   - Matikan WiFi/mobile data
   - Coba load data
   - Verifikasi: Snackbar error muncul dengan message ID
   - Verifikasi: Retry button tersedia

2. **Timeout Error**:
   - Simulate slow network
   - Coba API call
   - Verifikasi: Timeout error message muncul

3. **Validation Error**:
   - Submit form dengan data invalid
   - Verifikasi: Inline errors muncul
   - Verifikasi: Field yang error di-highlight

4. **Authentication Error**:
   - Wait sampai token expired
   - Coba API call
   - Verifikasi: Dialog session expired muncul
   - Verifikasi: Redirect ke login setelah dismiss

5. **Server Error**:
   - Trigger server error (via mocking atau downtime)
   - Verifikasi: Server error message muncul

6. **Bilingual Messages**:
   - Change app language ke English
   - Trigger error
   - Verifikasi: Error message dalam English

---

## Dependencies

### Internal

- `core/network/api_client.dart` - Dio client untuk error intercepting
- `core/storage/local_storage.dart` - Untuk cache error logs

### External

- `dio: ^5.0.0` - HTTP client dengan error handling
- `logger: ^2.0.0` - Logging framework
- `sentry_flutter: ^7.0.0` (optional) - Error tracking untuk production

---

## Notes & Improvements

### Known Limitations

1. **No Retry Queue**: Belum implement retry queue untuk failed operations.

2. **Limited Error Analytics**: Belum integrate dengan comprehensive error analytics.

3. **No User Feedback**: Belum mekanisme untuk user report errors.

### Future Improvements

1. **Retry Queue**: Implement queue untuk auto-retry failed operations

2. **Error Analytics**: Integrasi dengan Firebase Crashlytics atau Sentry

3. **User Feedback**: Add "Report Error" button untuk user feedback

4. **Smart Retry**: Implement smart retry dengan exponential backoff dan circuit breaker

5. **Offline Queue**: Queue operations saat offline dan auto-retry saat online

6. **Error Recovery**: Better error recovery strategies (contoh: load cached data saat network error)

---

**Document Status**: Active  
**Last Updated**: January 2025  
**Maintained By**: Dev3 (Mobile Development Team)
