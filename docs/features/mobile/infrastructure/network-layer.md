# Infrastructure - Network Layer

## CRM Healthcare Mobile App - Flutter

**Module**: Infrastructure  
**Sprint**: Sprint 0  
**Version**: 1.1  
**Status**: ✅ **Completed**  
**Last Updated**: March 2026

---

## Table of Contents

1. [Ringkasan Fitur](#ringkasan-fitur)
2. [Network Components](#network-components)
3. [Business Rules](#business-rules)
4. [Keputusan Teknis & Trade-offs](#keputusan-teknis--trade-offs)
5. [Struktur Folder](#struktur-folder)
6. [Implementation](#implementation)
7. [Configuration](#configuration)
8. [Usage Examples](#usage-examples)
9. [Cara Test Manual](#cara-test-manual)
10. [Dependencies](#dependencies)
11. [Notes & Improvements](#notes--improvements)

---

## Ringkasan Fitur

Network Layer menyediakan komunikasi antara mobile app dengan backend API menggunakan **Dio** HTTP client. Layer ini mencakup connectivity monitoring, retry logic, timeout handling, dan request queueing untuk offline support.

### Goals

- **HTTP Client**: Robust client dengan interceptors
- **Connectivity**: Monitor network status
- **Retry Logic**: Auto-retry failed requests
- **Offline Queue**: Queue requests saat offline
- **Performance**: Optimized timeouts dan caching

---

## Network Components

### 1. Connectivity Service

Monitor network status dan trigger actions saat online/offline.

### 2. Retry Logic

Exponential backoff untuk failed requests:

- Retry 1: 1 detik
- Retry 2: 2 detik
- Retry 3: 4 detik
- Max retries: 3

### 3. Request Queue

Queue untuk POST/PUT/DELETE saat offline:

- Auto-sync saat online
- Persist queue ke storage
- Error handling untuk failed sync

### 4. Timeout Configuration

**Connection Timeout**: 30 detik  
**Receive Timeout**: 30 detik  
**Send Timeout**: 60 detik (untuk uploads)

---

## Business Rules

### 1. Request Types

**Auto-retry**:

- GET requests
- Idempotent operations

**No Retry**:

- POST dengan side effects
- Payment transactions

**Queue saat Offline**:

- POST/PUT/DELETE
- User-initiated actions

### 2. Connectivity States

**Online**: All operations allowed  
**Offline**: Queue non-GET requests  
**Poor Connection**: Retry dengan backoff

### 3. Error Handling

**Network Error**: Retry dengan backoff  
**Timeout Error**: Retry once immediately  
**Server Error**: No retry, report error  
**Auth Error**: Refresh token, then retry. Jika refresh gagal, `onLogout` callback dipanggil yang:
  1. Clear auth tokens via `AuthNotifier.logout()`
  2. Navigate ke login screen via global `MyApp.navigatorKey` (agar bisa navigate dari non-widget context)

---

## Keputusan Teknis & Trade-offs

### Exponential vs Linear Backoff

**Keputusan**: Exponential backoff.

**Alasan**:

- Reduce server load saat recovery
- Network issues usually temporary
- Industry standard pattern

---

## Struktur Folder

```
apps/mobile/lib/
├── core/
│   ├── network/
│   │   ├── api_client.dart
│   │   ├── connectivity_service.dart
│   │   ├── retry_interceptor.dart
│   │   └── request_queue.dart
│   └── services/
│       └── sync_service.dart
```

---

## Implementation

### Connectivity Service

```dart
class ConnectivityService {
  final Connectivity _connectivity = Connectivity();
  StreamSubscription? _subscription;

  final _controller = StreamController<bool>.broadcast();
  Stream<bool> get onConnectivityChanged => _controller.stream;

  bool _isOnline = true;
  bool get isOnline => _isOnline;

  Future<void> init() async {
    final result = await _connectivity.checkConnectivity();
    _updateStatus(result);

    _subscription = _connectivity.onConnectivityChanged.listen(_updateStatus);
  }

  void _updateStatus(ConnectivityResult result) {
    _isOnline = result != ConnectivityResult.none;
    _controller.add(_isOnline);
  }

  void dispose() {
    _subscription?.cancel();
    _controller.close();
  }
}
```

### Retry Interceptor

```dart
class RetryInterceptor extends Interceptor {
  final Dio _dio;
  final int maxRetries;

  RetryInterceptor(this._dio, {this.maxRetries = 3});

  @override
  Future<void> onError(
    DioException err,
    ErrorInterceptorHandler handler,
  ) async {
    if (_shouldRetry(err) && err.requestOptions.extra['retries'] < maxRetries) {
      final retries = (err.requestOptions.extra['retries'] ?? 0) + 1;
      err.requestOptions.extra['retries'] = retries;

      // Exponential backoff
      final delay = Duration(seconds: pow(2, retries - 1).toInt());
      await Future.delayed(delay);

      try {
        final response = await _dio.fetch(err.requestOptions);
        handler.resolve(response);
        return;
      } catch (e) {
        // Continue to error
      }
    }

    handler.next(err);
  }

  bool _shouldRetry(DioException error) {
    return error.type == DioExceptionType.connectionError ||
           error.type == DioExceptionType.connectionTimeout ||
           error.type == DioExceptionType.receiveTimeout;
  }
}
```

---

## Cara Test Manual

1. **Connectivity Detection**:
   - Matikan WiFi
   - Verifikasi: Offline status detected
   - Nyalakan WiFi
   - Verifikasi: Online status detected

2. **Retry Logic**:
   - Simulate network error
   - Verifikasi: Retry attempts dengan delay

3. **Offline Queue**:
   - Create task saat offline
   - Verifikasi: Queued locally
   - Go online
   - Verifikasi: Auto-sync

---

## Dependencies

```yaml
dependencies:
  dio: ^5.4.0
  connectivity_plus: ^5.0.0
```

---

**Document Status**: Active  
**Last Updated**: March 2026
