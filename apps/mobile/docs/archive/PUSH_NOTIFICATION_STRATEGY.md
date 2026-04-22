# Push Notification Strategy - Mobile App

## Overview

Dokumen ini menjelaskan strategi implementasi push notification untuk aplikasi mobile CRM Healthcare menggunakan Firebase Cloud Messaging (FCM).

## Teknologi Stack

### 1. Firebase Cloud Messaging (FCM)
- **Platform**: Firebase Cloud Messaging
- **Alasan**: 
  - Free dan reliable
  - Support Android & iOS
  - Easy integration dengan Flutter
  - Scalable untuk production
  - Built-in analytics

### 2. Flutter Packages
- `firebase_messaging`: Official FCM plugin untuk Flutter
- `firebase_core`: Firebase core functionality
- `flutter_local_notifications`: Untuk menampilkan local notifications
- `permission_handler`: Untuk request notification permissions

## Arsitektur & Flow

### 1. Setup Flow

```
1. Setup Firebase Project
   ├── Create Firebase project
   ├── Add Android app (package name)
   ├── Add iOS app (bundle ID)
   ├── Download google-services.json (Android)
   └── Download GoogleService-Info.plist (iOS)

2. Configure Flutter App
   ├── Add firebase packages
   ├── Initialize Firebase
   ├── Request notification permissions
   └── Get FCM token

3. Register Token ke Backend
   ├── Send FCM token ke API
   ├── Backend store token per user
   └── Associate token dengan user ID
```

### 2. Notification Flow

```
Backend → FCM → Device
   │
   ├── Foreground: App menerima via onMessage
   ├── Background: System tray notification
   └── Terminated: System tray notification
```

### 3. Data Flow

```
User Action → Backend API → FCM Service → Device
   │
   └── Notification received
       ├── Update local state
       ├── Refresh notification list
       └── Show in-app notification
```

## Strategi Implementasi

### Phase 1: Setup & Infrastructure

#### 1.1 Firebase Setup
- [ ] Create Firebase project di Firebase Console
- [ ] Add Android app dengan package name
- [ ] Add iOS app dengan bundle ID
- [ ] Download konfigurasi files
- [ ] Setup Firebase Cloud Messaging API

#### 1.2 Flutter Dependencies
```yaml
dependencies:
  firebase_core: ^3.0.0
  firebase_messaging: ^15.0.0
  flutter_local_notifications: ^17.0.0
  permission_handler: ^11.0.0
```

#### 1.3 Platform Configuration
- **Android**: 
  - Add `google-services.json` ke `android/app/`
  - Update `build.gradle` files
  - Configure notification channels
  - Setup notification icons

- **iOS**:
  - Add `GoogleService-Info.plist` ke `ios/Runner/`
  - Enable Push Notifications capability
  - Configure APNs (Apple Push Notification service)
  - Request permissions di Info.plist

### Phase 2: Core Implementation

#### 2.1 FCM Service Layer

**File**: `lib/core/services/push_notification_service.dart`

```dart
class PushNotificationService {
  // Initialize FCM
  Future<void> initialize() async {
    // Request permissions
    // Get FCM token
    // Setup message handlers
    // Register token ke backend
  }
  
  // Get FCM token
  Future<String?> getToken() async { }
  
  // Register token ke backend
  Future<void> registerToken(String token) async { }
  
  // Handle foreground messages
  void onMessage(RemoteMessage message) { }
  
  // Handle background messages
  static Future<void> onBackgroundMessage(RemoteMessage message) async { }
  
  // Handle notification tap
  void onNotificationTap(String? data) { }
}
```

#### 2.2 Notification Handler

**File**: `lib/core/services/notification_handler.dart`

```dart
class NotificationHandler {
  // Show local notification
  Future<void> showNotification(RemoteMessage message) async { }
  
  // Navigate based on notification data
  void handleNotificationNavigation(Map<String, dynamic> data) { }
  
  // Update notification badge count
  void updateBadgeCount(int count) { }
}
```

#### 2.3 Token Management

**File**: `lib/core/services/fcm_token_manager.dart`

```dart
class FCMTokenManager {
  // Store token locally
  Future<void> saveToken(String token) async { }
  
  // Get stored token
  Future<String?> getStoredToken() async { }
  
  // Check if token needs refresh
  bool needsRefresh(String? currentToken) { }
  
  // Refresh token
  Future<String?> refreshToken() async { }
}
```

### Phase 3: Integration dengan Existing Features

#### 3.1 Integration dengan Notification Feature

**File**: `lib/features/notifications/application/notification_provider.dart`

```dart
// Update notification count saat push notification diterima
void onPushNotificationReceived(RemoteMessage message) {
  // Increment unread count
  // Refresh notification list
  // Show in-app notification
}
```

#### 3.2 Integration dengan Auth

**File**: `lib/features/auth/application/auth_provider.dart`

```dart
// Register FCM token setelah login
Future<void> login(...) async {
  // ... existing login logic
  // Register FCM token
  await pushNotificationService.registerToken();
}

// Unregister FCM token saat logout
Future<void> logout() async {
  // Unregister token dari backend
  await pushNotificationService.unregisterToken();
  // ... existing logout logic
}
```

### Phase 4: Backend Integration

#### 4.1 API Endpoints

**Required Backend Endpoints**:

1. **Register FCM Token**
   ```
   POST /api/v1/notifications/register-token
   Body: { "fcm_token": "string", "device_type": "android|ios" }
   ```

2. **Unregister FCM Token**
   ```
   DELETE /api/v1/notifications/unregister-token
   Body: { "fcm_token": "string" }
   ```

3. **Update FCM Token** (untuk token refresh)
   ```
   PUT /api/v1/notifications/update-token
   Body: { "old_token": "string", "new_token": "string" }
   ```

#### 4.2 Backend FCM Service

Backend perlu:
- Store FCM tokens per user
- Send notifications via FCM Admin SDK
- Handle token refresh
- Clean up invalid tokens

### Phase 5: Notification Types & Handling

#### 5.1 Notification Types

Berdasarkan existing notification types:
- `task` - Task assignment, completion, reminder
- `visit` - Visit report created, approved, rejected
- `deal` - Deal status changed, won, lost
- `account` - Account updated, new contact added
- `general` - General notifications

#### 5.2 Notification Data Structure

```json
{
  "notification": {
    "title": "New Task Assigned",
    "body": "You have been assigned a new task"
  },
  "data": {
    "type": "task",
    "id": "task-uuid",
    "action": "assigned",
    "screen": "tasks",
    "screen_id": "task-uuid"
  }
}
```

#### 5.3 Navigation Mapping

```dart
Map<String, String> notificationRoutes = {
  'task': '/tasks/{id}',
  'visit': '/visit-reports/{id}',
  'deal': '/deals/{id}',
  'account': '/accounts/{id}',
  'contact': '/contacts/{id}',
};
```

### Phase 6: UI/UX Implementation

#### 6.1 In-App Notification Banner

**File**: `lib/core/widgets/notification_banner.dart`

- Show banner saat app di foreground
- Auto-dismiss setelah 5 detik
- Tap untuk navigate ke detail
- Swipe to dismiss

#### 6.2 Notification Badge Update

- Update badge count di AppBar
- Real-time update saat notification diterima
- Clear badge saat semua notifications dibaca

#### 6.3 Notification Settings

**File**: `lib/features/profile/presentation/notification_settings_screen.dart`

- Enable/disable push notifications
- Notification types preferences
- Quiet hours settings
- Sound & vibration settings

### Phase 7: Background & Foreground Handling

#### 7.1 Foreground Handling

```dart
FirebaseMessaging.onMessage.listen((RemoteMessage message) {
  // Show in-app notification banner
  // Update notification count
  // Play sound/vibration
  // Refresh notification list
});
```

#### 7.2 Background Handling

```dart
@pragma('vm:entry-point')
Future<void> firebaseMessagingBackgroundHandler(RemoteMessage message) async {
  // Show system notification
  // Update local database
  // Schedule local notification
}
```

#### 7.3 Terminated App Handling

- System automatically shows notification
- On tap: Open app and navigate to detail
- Handle deep linking

### Phase 8: Testing Strategy

#### 8.1 Unit Testing
- Test FCM token registration
- Test notification parsing
- Test navigation logic
- Test cache updates

#### 8.2 Integration Testing
- Test end-to-end notification flow
- Test token refresh
- Test notification tap navigation
- Test background/foreground scenarios

#### 8.3 Manual Testing Checklist
- [ ] Receive notification saat app di foreground
- [ ] Receive notification saat app di background
- [ ] Receive notification saat app terminated
- [ ] Tap notification navigates correctly
- [ ] Notification badge updates
- [ ] Token refresh works
- [ ] Multiple devices support
- [ ] Notification permissions handling

## Best Practices

### 1. Token Management
- **Store token locally** untuk comparison
- **Refresh token** saat app start jika berbeda
- **Unregister token** saat logout
- **Handle token refresh** dari FCM

### 2. Notification Handling
- **Show in-app notification** saat foreground
- **System notification** saat background/terminated
- **Update badge count** real-time
- **Navigate to detail** saat notification tapped

### 3. Performance
- **Lazy load** notification data
- **Cache** notification list
- **Batch update** badge count
- **Debounce** notification refresh

### 4. Security
- **Validate** notification data
- **Sanitize** notification content
- **Verify** notification source
- **Encrypt** sensitive data

### 5. Error Handling
- **Handle** token registration failures
- **Retry** failed token registrations
- **Fallback** untuk notification display
- **Log** errors untuk debugging

## Implementation Checklist

### Setup Phase
- [ ] Create Firebase project
- [ ] Configure Android app
- [ ] Configure iOS app
- [ ] Add Flutter dependencies
- [ ] Initialize Firebase in app

### Core Implementation
- [ ] Create PushNotificationService
- [ ] Create NotificationHandler
- [ ] Create FCMTokenManager
- [ ] Implement token registration
- [ ] Implement foreground handler
- [ ] Implement background handler
- [ ] Implement notification tap handler

### Integration
- [ ] Integrate dengan auth (register/unregister)
- [ ] Integrate dengan notification feature
- [ ] Update notification badge
- [ ] Implement navigation logic

### Backend Integration
- [ ] Create register token endpoint
- [ ] Create unregister token endpoint
- [ ] Create update token endpoint
- [ ] Implement FCM sending logic
- [ ] Store tokens per user

### UI/UX
- [ ] Create notification banner widget
- [ ] Update notification settings screen
- [ ] Implement notification preferences
- [ ] Add notification icons

### Testing
- [ ] Unit tests
- [ ] Integration tests
- [ ] Manual testing
- [ ] Performance testing

## Security Considerations

### 1. Token Security
- Store tokens securely (encrypted storage)
- Validate token format
- Handle token expiration
- Prevent token leakage

### 2. Notification Content
- Sanitize notification content
- Validate notification data
- Prevent XSS attacks
- Encrypt sensitive data

### 3. Backend Security
- Authenticate FCM requests
- Validate user permissions
- Rate limit notification sending
- Log all notification activities

## Monitoring & Analytics

### 1. Metrics to Track
- Token registration success rate
- Notification delivery rate
- Notification open rate
- Token refresh frequency
- Error rates

### 2. Logging
- Log token registrations
- Log notification receipts
- Log navigation actions
- Log errors

### 3. Firebase Analytics
- Track notification events
- Track user engagement
- Track notification preferences
- Track error occurrences

## Troubleshooting

### Common Issues

1. **Token not received**
   - Check Firebase configuration
   - Verify permissions granted
   - Check network connectivity
   - Review Firebase console logs

2. **Notifications not showing**
   - Check notification permissions
   - Verify notification channel setup (Android)
   - Check APNs configuration (iOS)
   - Review notification payload

3. **Navigation not working**
   - Verify notification data structure
   - Check route definitions
   - Review navigation logic
   - Test deep linking

4. **Token refresh issues**
   - Check token storage
   - Verify backend token update endpoint
   - Review token comparison logic
   - Check FCM token refresh callbacks

## Future Enhancements

### Phase 2 Features
- [ ] Rich notifications (images, actions)
- [ ] Notification grouping
- [ ] Notification scheduling
- [ ] Notification templates
- [ ] Notification analytics dashboard
- [ ] A/B testing untuk notifications
- [ ] Notification preferences per type
- [ ] Quiet hours implementation
- [ ] Notification sound customization

## References

- [Firebase Cloud Messaging Documentation](https://firebase.google.com/docs/cloud-messaging)
- [Flutter Firebase Messaging Plugin](https://pub.dev/packages/firebase_messaging)
- [Flutter Local Notifications](https://pub.dev/packages/flutter_local_notifications)
- [FCM Best Practices](https://firebase.google.com/docs/cloud-messaging/best-practices)

