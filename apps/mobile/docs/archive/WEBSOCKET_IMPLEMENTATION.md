# WebSocket Implementation untuk Real-time Notifications

## Overview

Dokumen ini menjelaskan cara mengimplementasikan WebSocket untuk real-time notifications pada mobile app, mengikuti pola yang sama dengan versi web.

## Prerequisites

1. **Package yang diperlukan:**
   ```yaml
   dependencies:
     web_socket_channel: ^2.4.0
   ```

2. **Backend WebSocket Endpoint:**
   - URL: `ws://<base_url>/api/v1/ws/notifications?token=<access_token>`
   - Protocol: WebSocket
   - Authentication: Token via query parameter

## Struktur Implementasi

### 1. WebSocket Service

Buat file `apps/mobile/lib/features/notifications/data/notification_websocket_service.dart`:

```dart
import 'dart:async';
import 'dart:convert';
import 'package:web_socket_channel/web_socket_channel.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../application/notification_provider.dart';
import '../data/models/notification.dart' as models;

class NotificationWebSocketService {
  NotificationWebSocketService(this._ref);
  
  final Ref _ref;
  WebSocketChannel? _channel;
  Timer? _reconnectTimer;
  int _reconnectAttempts = 0;
  static const int _maxReconnectAttempts = 5;
  static const Duration _reconnectDelay = Duration(seconds: 5);

  /// Connect to WebSocket server
  Future<void> connect(String token, String baseUrl) async {
    try {
      // Close existing connection if any
      await disconnect();

      // Build WebSocket URL
      final wsUrl = baseUrl
          .replaceFirst('http://', 'ws://')
          .replaceFirst('https://', 'wss://');
      final url = '$wsUrl/api/v1/ws/notifications?token=${Uri.encodeComponent(token)}';

      // Create WebSocket connection
      _channel = WebSocketChannel.connect(Uri.parse(url));

      // Listen to messages
      _channel!.stream.listen(
        _handleMessage,
        onError: _handleError,
        onDone: _handleDisconnect,
        cancelOnError: false,
      );

      _reconnectAttempts = 0;
    } catch (e) {
      print('WebSocket connection error: $e');
      _scheduleReconnect(token, baseUrl);
    }
  }

  /// Handle incoming WebSocket messages
  void _handleMessage(dynamic message) {
    try {
      final data = jsonDecode(message as String) as Map<String, dynamic>;
      final type = data['type'] as String?;

      switch (type) {
        case 'notification.created':
          _handleNotificationCreated(data['data'] as Map<String, dynamic>);
          break;
        case 'notification.updated':
          _handleNotificationUpdated(data['data'] as Map<String, dynamic>);
          break;
        case 'notification.deleted':
          _handleNotificationDeleted(data['data'] as Map<String, dynamic>);
          break;
        default:
          print('Unknown WebSocket message type: $type');
      }
    } catch (e) {
      print('Error parsing WebSocket message: $e');
    }
  }

  /// Handle notification.created event
  void _handleNotificationCreated(Map<String, dynamic> data) {
    try {
      final notification = models.Notification.fromJson(data);
      
      // Invalidate queries to refresh data
      _ref.invalidate(notificationListProvider);
      _ref.invalidate(notificationCountProvider);
      
      // Optionally show a toast/notification
      // You can use flutter_local_notifications or similar package
    } catch (e) {
      print('Error handling notification.created: $e');
    }
  }

  /// Handle notification.updated event
  void _handleNotificationUpdated(Map<String, dynamic> data) {
    try {
      // Invalidate queries to refresh data
      _ref.invalidate(notificationListProvider);
      _ref.invalidate(notificationCountProvider);
    } catch (e) {
      print('Error handling notification.updated: $e');
    }
  }

  /// Handle notification.deleted event
  void _handleNotificationDeleted(Map<String, dynamic> data) {
    try {
      // Invalidate queries to refresh data
      _ref.invalidate(notificationListProvider);
      _ref.invalidate(notificationCountProvider);
    } catch (e) {
      print('Error handling notification.deleted: $e');
    }
  }

  /// Handle WebSocket errors
  void _handleError(dynamic error) {
    print('WebSocket error: $error');
    // Don't reconnect on error, let onDone handle it
  }

  /// Handle WebSocket disconnection
  void _handleDisconnect() {
    print('WebSocket disconnected');
    _channel = null;
    // Schedule reconnection if not manually disconnected
    if (_reconnectAttempts < _maxReconnectAttempts) {
      // Get token and baseUrl from auth/storage
      // _scheduleReconnect(token, baseUrl);
    }
  }

  /// Schedule reconnection attempt
  void _scheduleReconnect(String token, String baseUrl) {
    if (_reconnectAttempts >= _maxReconnectAttempts) {
      print('Max reconnection attempts reached');
      return;
    }

    _reconnectAttempts++;
    _reconnectTimer?.cancel();
    _reconnectTimer = Timer(_reconnectDelay, () {
      connect(token, baseUrl);
    });
  }

  /// Disconnect from WebSocket server
  Future<void> disconnect() async {
    _reconnectTimer?.cancel();
    _reconnectTimer = null;
    await _channel?.sink.close();
    _channel = null;
    _reconnectAttempts = 0;
  }

  /// Check if WebSocket is connected
  bool get isConnected => _channel != null;
}
```

### 2. WebSocket Provider

Buat file `apps/mobile/lib/features/notifications/application/notification_websocket_provider.dart`:

```dart
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../auth/application/auth_provider.dart';
import '../../../core/network/api_client.dart';
import '../data/notification_websocket_service.dart';

final notificationWebSocketServiceProvider = Provider<NotificationWebSocketService?>((ref) {
  final authState = ref.watch(authProvider);
  
  if (!authState.isAuthenticated) {
    return null;
  }

  return NotificationWebSocketService(ref);
});

final notificationWebSocketProvider = StateNotifierProvider<NotificationWebSocketNotifier, bool>((ref) {
  final service = ref.watch(notificationWebSocketServiceProvider);
  return NotificationWebSocketNotifier(service);
});

class NotificationWebSocketNotifier extends StateNotifier<bool> {
  NotificationWebSocketNotifier(this._service) : super(false) {
    _init();
  }

  final NotificationWebSocketService? _service;

  Future<void> _init() async {
    if (_service == null) return;

    // Get token from auth state
    final authState = ref.read(authProvider);
    if (!authState.isAuthenticated || authState.token == null) {
      return;
    }

    // Get base URL from ApiClient
    final baseUrl = ApiClient.baseUrl; // You may need to expose this

    try {
      await _service!.connect(authState.token!, baseUrl);
      state = true;
    } catch (e) {
      print('Failed to connect WebSocket: $e');
      state = false;
    }
  }

  Future<void> reconnect() async {
    await _init();
  }

  Future<void> disconnect() async {
    await _service?.disconnect();
    state = false;
  }
}
```

### 3. Integrasi dengan Main App

Update `apps/mobile/lib/main.dart`:

```dart
// ... existing imports ...
import 'features/notifications/application/notification_websocket_provider.dart';

class MyApp extends ConsumerStatefulWidget {
  // ... existing code ...
}

class _MyAppState extends ConsumerState<MyApp> with WidgetsBindingObserver {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addObserver(this);
    
    // Initialize WebSocket when app starts
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _initWebSocket();
    });
  }

  @override
  void dispose() {
    WidgetsBinding.instance.removeObserver(this);
    // Disconnect WebSocket when app closes
    ref.read(notificationWebSocketProvider.notifier).disconnect();
    super.dispose();
  }

  void _initWebSocket() {
    final authState = ref.read(authProvider);
    if (authState.isAuthenticated) {
      ref.read(notificationWebSocketProvider.notifier).reconnect();
    }
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    super.didChangeAppLifecycleState(state);
    
    // Reconnect when app comes to foreground
    if (state == AppLifecycleState.resumed) {
      _initWebSocket();
    } else if (state == AppLifecycleState.paused) {
      // Optionally disconnect when app goes to background
      // ref.read(notificationWebSocketProvider.notifier).disconnect();
    }
  }

  // ... rest of the code ...
}
```

### 4. Integrasi dengan Auth Provider

Update `apps/mobile/lib/features/auth/application/auth_provider.dart`:

```dart
// After successful login:
Future<void> login(String email, String password) async {
  // ... existing login code ...
  
  // Connect WebSocket after login
  WidgetsBinding.instance.addPostFrameCallback((_) {
    ref.read(notificationWebSocketProvider.notifier).reconnect();
  });
}

// After logout:
Future<void> logout() async {
  // ... existing logout code ...
  
  // Disconnect WebSocket on logout
  ref.read(notificationWebSocketProvider.notifier).disconnect();
}
```

## Testing

1. **Test Connection:**
   - Login ke aplikasi
   - Check console logs untuk konfirmasi WebSocket connection
   - Verify bahwa notification list auto-refresh saat ada notification baru

2. **Test Reconnection:**
   - Putuskan koneksi internet
   - Hubungkan kembali
   - Verify bahwa WebSocket otomatis reconnect

3. **Test Events:**
   - Buat notification baru dari backend/Postman
   - Verify bahwa notification muncul di list tanpa perlu refresh manual
   - Verify bahwa unread count update otomatis

## Troubleshooting

1. **WebSocket tidak connect:**
   - Check token validity
   - Check base URL format (ws:// untuk HTTP, wss:// untuk HTTPS)
   - Check network connectivity
   - Check backend WebSocket endpoint availability

2. **Messages tidak diterima:**
   - Check message format dari backend
   - Check event type handling
   - Check console logs untuk error messages

3. **Reconnection tidak bekerja:**
   - Check max reconnection attempts
   - Check reconnection delay
   - Verify token masih valid saat reconnect

## Catatan Penting

1. **Token Expiration:**
   - WebSocket connection akan terputus saat token expired
   - Implement token refresh mechanism sebelum reconnect
   - Atau handle 401 error dan redirect ke login

2. **Battery Optimization:**
   - WebSocket connection dapat menguras baterai
   - Pertimbangkan untuk disconnect saat app di background
   - Atau gunakan background service untuk maintain connection

3. **Network Changes:**
   - Handle network state changes (WiFi to mobile data, etc.)
   - Reconnect saat network kembali tersedia

4. **Error Handling:**
   - Implement proper error handling untuk semua WebSocket operations
   - Show user-friendly error messages jika diperlukan
   - Log errors untuk debugging

## Referensi

- Web version implementation: `apps/web/src/features/notifications/hooks/useWebSocket.ts`
- WebSocket package: https://pub.dev/packages/web_socket_channel
- Flutter WebSocket guide: https://docs.flutter.dev/cookbook/networking/web-sockets

