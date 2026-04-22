# Business - Notifications Center

## CRM Healthcare Mobile App - Flutter

**Module**: Business Domain  
**Sprint**: Sprint 3  
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
7. [Data Models](#data-models)
8. [Configuration](#configuration)
9. [Usage Examples](#usage-examples)
10. [Cara Test Manual](#cara-test-manual)
11. [Dependencies](#dependencies)
12. [Notes & Improvements](#notes--improvements)

---

## Ringkasan Fitur

Fitur **Notifications Center** menyediakan centralized inbox untuk semua notifikasi aplikasi. User dapat melihat history notifikasi, mark as read, dan navigate ke relevant screens dari notifikasi.

### Goals

- **Notification History**: Lihat semua notifikasi yang diterima
- **Read Management**: Mark as read/unread
- **Quick Actions**: Navigate ke relevant content
- **Filtering**: Filter by type dan status
- **Clear All**: Bulk operations

---

## Fitur Utama

### 1. Notification List

**Display**:

- List notifikasi dengan timestamp
- Unread indicator (badge)
- Icon berdasarkan type
- Swipe actions

### 2. Notification Types

**Types**:

- 📋 Task reminders
- ✅ Task assignments
- 📊 Visit report approvals
- 🔔 System announcements
- 📅 Schedule changes

### 3. Notification Actions

**Swipe Actions**:

- Mark as read/unread
- Delete
- Archive

**Tap Actions**:

- Navigate ke related screen
- Mark as read automatically

### 4. Badge Counter

**Features**:

- Show unread count di bottom nav
- Update real-time
- Clear when semua read

---

## Business Rules

### 1. Notification Retention

- Store 30 days history
- Auto-delete after 30 days
- Max 100 notifications per user

### 2. Read Status

- Tap notification = mark as read
- Can manually mark as unread
- Bulk mark as read available

### 3. Priority Levels

**High**: Task due, visit report approval  
**Normal**: Task assignment, general updates  
**Low**: System announcements

---

## Keputusan Teknis & Trade-offs

### Local vs Server Storage

**Keputusan**: Hybrid approach - cache locally, sync dengan server.

**Alasan**:

- Fast access untuk recent notifications
- Offline capability
- Sync across devices

---

## Struktur Folder

```
apps/mobile/lib/
├── features/
│   └── notifications/
│       ├── data/
│       │   ├── models/
│       │   │   └── notification_model.dart
│       │   └── notification_repository.dart
│       ├── application/
│       │   ├── notification_list_provider.dart
│       │   └── notification_badge_provider.dart
│       └── presentation/
│           ├── screens/
│           │   └── notification_center_screen.dart
│           └── widgets/
│               ├── notification_item.dart
│               └── notification_filter_sheet.dart
```

---

## API Endpoints

#### GET /api/v1/notifications

Get user notifications.

**Query Parameters**:

```
?page=1&limit=20&status=unread&type=task
```

**Response**:

```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "notif-uuid",
        "type": "task_reminder",
        "title": "Task Due Soon",
        "message": "Follow up with RS Medika is due in 1 hour",
        "data": {
          "task_id": "task-uuid",
          "screen": "task_detail"
        },
        "is_read": false,
        "created_at": "2025-01-20T13:00:00Z"
      }
    ],
    "unread_count": 5
  }
}
```

#### POST /api/v1/notifications/:id/read

Mark notification as read.

#### POST /api/v1/notifications/read-all

Mark all as read.

#### DELETE /api/v1/notifications/:id

Delete notification.

---

## Cara Test Manual

1. **Receive Notification**: Trigger notification dan verifikasi muncul di list
2. **Mark as Read**: Tap notification, verifikasi badge count berkurang
3. **Navigate**: Tap notification dengan action, verifikasi navigate ke screen yang benar
4. **Clear All**: Tap "Mark all as read", verifikasi semua marked
5. **Delete**: Swipe delete, verifikasi notification dihapus

---

**Document Status**: Active  
**Last Updated**: January 2025
