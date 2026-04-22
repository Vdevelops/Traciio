# Core - Push Notifications

## CRM Healthcare Mobile App - Flutter

**Module**: Core Infrastructure  
**Sprint**: Sprint 3  
**Version**: 1.0  
**Status**: ⏳ **In Progress** (Menunggu Backend Push Notification Service)  
**Last Updated**: January 2025

---

## Table of Contents

1. [Ringkasan Fitur](#ringkasan-fitur)
2. [Fitur Utama](#fitur-utama)
3. [Business Rules](#business-rules)
4. [Keputusan Teknis & Trade-offs](#keputusan-teknis--trade-offs)
5. [Struktur Folder](#struktur-folder)
6. [API / Package Reference](#api--package-reference)
7. [Configuration](#configuration)
8. [Notification Flow](#notification-flow)
9. [Cara Test Manual](#cara-test-manual)
10. [Dependencies](#dependencies)
11. [Notes & Improvements](#notes--improvements)

---

## Ringkasan Fitur

Sistem **Push Notifications** mobile app CRM Healthcare menggunakan **Firebase Cloud Messaging (FCM)** untuk mengirim notifikasi real-time ke user devices. Notifikasi digunakan untuk task reminders, visit report updates, dan system notifications. Sistem ini juga support local notifications untuk reminders saat FCM tidak available.

### Goals

- **Real-time Updates**: Push notifications untuk task reminders dan updates
- **Task Reminders**: Notifikasi untuk task yang mendekati due date
- **Visit Report Updates**: Notifikasi untuk visit report approval/rejection
- **System Notifications**: General system announcements
- **Offline Support**: Local notifications sebagai fallback

---

## Fitur Utama

### 1. Firebase Cloud Messaging (FCM)

**Token Management**:

- Registrasi FCM token saat login
- Token refresh handling
- Unsubscribe saat logout

**Message Types**:

- **Notification Messages**: Display notifications automatically
- **Data Messages**: Handle in background, custom display logic

### 2. Local Notifications

**flutter_local_notifications**:

- Scheduled reminders
- Immediate local notifications
- Rich notifications dengan actions

### 3. Notification Types

**Task Reminders**:

```json
{
  "type": "task_reminder",
  "title": "Task Reminder",
  "body": "Follow up with Dr. Smith",
  "data": {
    "task_id": "uuid",
    "due_date": "2025-01-20T10:00:00Z"
  }
}
```

**Visit Report Status**:

```json
{
  "type": "visit_report_status",
  "title": "Visit Report Approved",
  "body": "Your visit report to RS Medika has been approved",
  "data": {
    "visit_report_id": "uuid",
    "status": "approved"
  }
}
```

**System Notifications**:

```json
{
  "type": "system",
  "title": "System Maintenance",
  "body": "System will be down for maintenance tonight",
  "data": {}
}
```

### 4. Notification Actions

**Task Actions**:

- ✅ Mark as Complete
- ⏰ Snooze (15 min, 1 hour, 1 day)
- 📱 View Task

**Visit Report Actions**:

- 👁️ View Report
- ✏️ Edit Report

---

## Business Rules

### 1. Notification Permissions

**First Launch**:

- Request notification permission saat pertama kali app dibuka
- Explain kenapa notifications diperlukan
- Allow user to skip (can enable later in settings)

**Settings**:

- User dapat enable/disable notifications per type
- Global notification toggle
- Sound/vibration preferences

### 2. Notification Delivery Rules

**FCM Priority**:

- High priority untuk task reminders (vibration + sound)
- Normal priority untuk system notifications

**Local Notification Fallback**:

- Schedule local notification jika FCM failed
- Fallback untuk offline scenarios

### 3. Notification Timing

**Task Reminders**:

- 1 hour before due date
- At due date
- 1 day before untuk high priority tasks

**Visit Report Updates**:

- Immediate setelah status change

### 4. Notification Handling

**Foreground**: Show in-app banner atau snackbar
**Background**: Show system notification
**Terminated**: Show system notification + data sync saat app dibuka

---

## Keputusan Teknis & Trade-offs

### Mengapa Firebase Cloud Messaging?

**Keputusan**: Menggunakan FCM daripada custom push notification service.

**Alasan**:

- **Reliability**: FCM reliable delivery
- **Free**: No cost untuk notification delivery
- **Cross-platform**: Support Android dan iOS dengan satu implementation
- **Rich Features**: Topic subscriptions, device groups, dll.

**Trade-off**: Dependency pada Google services. **Mitigasi**: Implement local notifications sebagai fallback.

### Mengapa flutter_local_notifications?

**Keputusan**: Menggunakan flutter_local_notifications untuk local notifications.

**Alasan**:

- **Feature Rich**: Support untuk scheduled notifications, actions, attachments
- **Cross-platform**: Consistent API untuk Android dan iOS
- **Active Maintenance**: Package actively maintained

### Foreground vs Background Handling

**Keputusan**: Different handling untuk foreground dan background.

**Alasan**:

- **Foreground**: User sedang menggunakan app, jadi show in-app notification
- **Background**: User tidak menggunakan app, jadi show system notification
- **User Experience**: Tidak mengganggu user experience saat actively using app

---

## Struktur Folder

```
apps/mobile/lib/
├── core/
│   └── notifications/
│       ├── fcm_service.dart              # Firebase Cloud Messaging setup
│       ├── local_notification_service.dart # Local notifications
│       ├── notification_handler.dart     # Handle incoming notifications
│       ├── notification_router.dart      # Navigate dari notifications
│       └── notification_preferences.dart # User preferences
├── features/
│   └── notifications/
│       ├── data/
│       │   └── notification_repository.dart # API calls untuk FCM token
│       ├── application/
│       │   └── notification_provider.dart   # State management
│       └── presentation/
│           └── screens/
│               └── notification_settings.dart # Settings screen
└── android/app/src/main/
    └── AndroidManifest.xml               # FCM configuration
```

---

## API / Package Reference

### Backend Endpoints

#### POST /api/v1/users/:id/fcm-token

Register FCM token untuk user.

**Request**:

```json
{
  "fcm_token": "c9d8x7v6b5n4m3k2j1h0g9f8e7d6c5b4a3",
  "device_type": "android",
  "device_name": "Samsung Galaxy S21"
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "message": "FCM token registered successfully"
  }
}
```

#### DELETE /api/v1/users/:id/fcm-token

Unregister FCM token (logout).

**Response**:

```json
{
  "success": true,
  "data": {
    "message": "FCM token unregistered successfully"
  }
}
```

### Firebase Cloud Messaging

**Setup**:

1. **Firebase Console**: Create project dan download `google-services.json` (Android) dan `GoogleService-Info.plist` (iOS)

2. **pubspec.yaml**:

```yaml
dependencies:
  firebase_core: ^2.24.0
  firebase_messaging: ^14.7.0
  flutter_local_notifications: ^16.0.0
```

3. **AndroidManifest.xml**:

```xml
<application>
  <activity>
    <intent-filter>
      <action android:name="FLUTTER_NOTIFICATION_CLICK" />
      <category android:name="android.intent.category.DEFAULT" />
    </intent-filter>
  </activity>
</application>
```

### Notification Payload Format

**FCM Payload**:

```json
{
  "notification": {
    "title": "Task Reminder",
    "body": "Follow up with Dr. Smith",
    "sound": "default",
    "badge": "1"
  },
  "data": {
    "type": "task_reminder",
    "task_id": "uuid",
    "click_action": "FLUTTER_NOTIFICATION_CLICK"
  },
  "priority": "high",
  "to": "device_fcm_token"
}
```

---

## Configuration

### FCM Service

**File**: `core/notifications/fcm_service.dart`

```dart
class FCMService {
  final FirebaseMessaging _fcm = FirebaseMessaging.instance;

  Future<void> initialize() async {
    // Request permission
    await _requestPermission();

    // Get FCM token
    final token = await _fcm.getToken();
    await _saveTokenToServer(token);

    // Listen untuk token refresh
    _fcm.onTokenRefresh.listen(_saveTokenToServer);

    // Setup message handlers
    FirebaseMessaging.onMessage.listen(_handleForegroundMessage);
    FirebaseMessaging.onMessageOpenedApp.listen(_handleNotificationOpen);
    FirebaseMessaging.onBackgroundMessage(_handleBackgroundMessage);
  }

  Future<void> _requestPermission() async {
    await _fcm.requestPermission(
      alert: true,
      badge: true,
      sound: true,
    );
  }

  Future<void> _saveTokenToServer(String? token) async {
    if (token == null) return;

    final repository = ref.read(notificationRepositoryProvider);
    await repository.registerFCMToken(token);
  }

  void _handleForegroundMessage(RemoteMessage message) {
    // Show local notification atau in-app banner
    LocalNotificationService.showNotification(
      title: message.notification?.title ?? '',
      body: message.notification?.body ?? '',
      payload: jsonEncode(message.data),
    );
  }

  void _handleNotificationOpen(RemoteMessage message) {
    // Navigate ke screen yang sesuai
    NotificationRouter.navigate(message.data);
  }
}

// Background message handler (top-level function)
Future<void> _handleBackgroundMessage(RemoteMessage message) async {
  await Firebase.initializeApp();
  // Handle background message
  print('Background message: ${message.messageId}');
}
```

### Local Notification Service

**File**: `core/notifications/local_notification_service.dart`

```dart
class LocalNotificationService {
  static final FlutterLocalNotificationsPlugin _notifications =
      FlutterLocalNotificationsPlugin();

  static Future<void> initialize() async {
    const androidSettings = AndroidInitializationSettings('@mipmap/ic_launcher');
    const iosSettings = DarwinInitializationSettings();

    const initSettings = InitializationSettings(
      android: androidSettings,
      iOS: iosSettings,
    );

    await _notifications.initialize(
      initSettings,
      onDidReceiveNotificationResponse: _onNotificationTap,
    );
  }

  static Future<void> showNotification({
    required String title,
    required String body,
    String? payload,
  }) async {
    const androidDetails = AndroidNotificationDetails(
      'crm_channel',
      'CRM Notifications',
      importance: Importance.high,
      priority: Priority.high,
      showWhen: true,
      enableVibration: true,
    );

    const iosDetails = DarwinNotificationDetails(
      presentAlert: true,
      presentBadge: true,
      presentSound: true,
    );

    const details = NotificationDetails(
      android: androidDetails,
      iOS: iosDetails,
    );

    await _notifications.show(
      DateTime.now().millisecond,
      title,
      body,
      details,
      payload: payload,
    );
  }

  static Future<void> scheduleNotification({
    required String title,
    required String body,
    required DateTime scheduledDate,
    String? payload,
  }) async {
    await _notifications.zonedSchedule(
      DateTime.now().millisecond,
      title,
      body,
      tz.TZDateTime.from(scheduledDate, tz.local),
      const NotificationDetails(
        android: AndroidNotificationDetails(
          'reminder_channel',
          'Reminders',
          importance: Importance.high,
        ),
        iOS: DarwinNotificationDetails(),
      ),
      androidAllowWhileIdle: true,
      uiLocalNotificationDateInterpretation:
          UILocalNotificationDateInterpretation.absoluteTime,
    );
  }

  static void _onNotificationTap(NotificationResponse response) {
    if (response.payload != null) {
      final data = jsonDecode(response.payload!);
      NotificationRouter.navigate(data);
    }
  }
}
```

### Notification Router

**File**: `core/notifications/notification_router.dart`

```dart
class NotificationRouter {
  static void navigate(Map<String, dynamic> data) {
    final type = data['type'];

    switch (type) {
      case 'task_reminder':
      case 'task_assigned':
        final taskId = data['task_id'];
        navigatorKey.currentState?.pushNamed(
          AppRoutes.taskDetail,
          arguments: {'taskId': taskId},
        );
        break;

      case 'visit_report_status':
        final visitReportId = data['visit_report_id'];
        navigatorKey.currentState?.pushNamed(
          AppRoutes.visitReportDetail,
          arguments: {'visitReportId': visitReportId},
        );
        break;

      default:
        // Navigate ke dashboard
        navigatorKey.currentState?.pushNamed(AppRoutes.dashboard);
    }
  }
}
```

---

## Notification Flow

### Task Reminder Flow

```
Backend Scheduler
      │
      ▼
Check Due Tasks (every hour)
      │
      ▼
Send FCM Message
      │
      ├──────────────┐
      │              │
      ▼              ▼
App Foreground   App Background
      │              │
      ▼              ▼
Show In-App    Show System
Banner         Notification
      │              │
      └──────────────┘
             │
             ▼
      User Taps
             │
             ▼
      Navigate to
      Task Detail
```

### Token Management Flow

```
User Login
     │
     ▼
Get FCM Token
     │
     ▼
Register ke Backend
     │
     ▼
Store Token Locally
     │
     ▼
Listen Token Refresh
     │
     ▼
Update Backend Token
     │
     ▼
User Logout
     │
     ▼
Unregister Token
```

---

## Cara Test Manual

### Test FCM Setup

1. **Token Registration**:
   - Login sebagai user
   - Check console: FCM token harus logged
   - Verifikasi: Token terkirim ke backend

2. **Receive Notification (Foreground)**:
   - Keep app di foreground
   - Kirim test notification dari Firebase Console
   - Verifikasi: In-app banner muncul

3. **Receive Notification (Background)**:
   - Minimize app
   - Kirim test notification
   - Verifikasi: System notification muncul
   - Tap notification
   - Verifikasi: Navigate ke screen yang benar

4. **Receive Notification (Terminated)**:
   - Kill app
   - Kirim test notification
   - Verifikasi: System notification muncul
   - Tap notification
   - Verifikasi: App opens dan navigate ke screen yang benar

### Test Local Notifications

1. **Schedule Reminder**:
   - Create task dengan due date 1 menit dari sekarang
   - Verifikasi: Reminder scheduled
   - Wait 1 menit
   - Verifikasi: Local notification muncul

2. **Notification Actions**:
   - Show notification dengan actions
   - Tap action button
   - Verifikasi: Action executed correctly

### Test Notification Settings

1. **Disable Notifications**:
   - Go to notification settings
   - Disable notifications
   - Kirim test notification
   - Verifikasi: Notification tidak muncul

2. **Enable Notifications**:
   - Enable notifications
   - Kirim test notification
   - Verifikasi: Notification muncul

---

## Dependencies

### Internal

- `core/routing/app_router.dart` - Navigation dari notifications
- `features/auth/application/auth_provider.dart` - User authentication state

### External

- `firebase_core: ^2.24.0` - Firebase core
- `firebase_messaging: ^14.7.0` - FCM
- `flutter_local_notifications: ^16.0.0` - Local notifications
- `timezone: ^0.9.0` - Timezone support untuk scheduled notifications

---

## Notes & Improvements

### Known Limitations

1. **Backend Dependency**: Menunggu backend push notification service untuk task reminders.

2. **No Rich Notifications**: Belum implement rich notifications dengan images.

3. **Limited Actions**: Notification actions masih basic.

4. **No Notification History**: Belum implement notification history/inbox.

### Future Improvements

1. **Rich Notifications**: Support untuk notification dengan images dan attachments

2. **Notification History**: Implement notification inbox/history

3. **Smart Reminders**: AI-powered reminder timing based on user behavior

4. **Notification Analytics**: Track notification open rates dan engagement

5. **Topic Subscriptions**: Subscribe ke topics untuk group notifications

6. **Silent Notifications**: Background data sync tanpa user notification

### Current Status

⏳ **Waiting for Backend**:

- Task reminder scheduler
- FCM token management API
- Notification payload templates

✅ **Completed**:

- FCM setup dan configuration
- Local notification service
- Notification routing
- Basic notification display

---

**Document Status**: In Progress  
**Last Updated**: January 2025  
**Maintained By**: Dev3 (Mobile Development Team)  
**Blocked By**: Backend Push Notification Service
