# Business - Profile Management

## CRM Healthcare Mobile App - Flutter

**Module**: Business Domain  
**Sprint**: Sprint 0  
**Version**: 1.1  
**Status**: ✅ **Completed**  
**Last Updated**: March 2026

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

Fitur **Profile Management** memungkinkan user untuk melihat dan mengelola informasi profil mereka, mengubah pengaturan aplikasi, serta mengakses fitur logout. Profile screen juga menampilkan statistik performa sales rep dan quick access ke settings.

### Goals

- **Profile View**: Lihat informasi user dan performa
- **Settings**: Kelola pengaturan aplikasi
- **Preferences**: Language, notifications, theme
- **Security**: Change password, biometric auth
- **Logout**: Secure logout dengan token cleanup

---

## Fitur Utama

### 1. Profile Header

**Display Information**:

- Avatar/profile photo
- Full name
- Email
- Role (Sales Rep, Supervisor, Admin)
- Organization/Team info
- Member since date

### 2. Performance Summary

**Quick Stats**:

- Total visits this month
- Tasks completed
- Performance rating
- Achievement badges

### 3. Settings Menu

**Menu Items**:

- 🔔 Notification Settings
- 🌐 Language (ID/EN)
- 🎨 Theme (Light/Dark/System)
- 🔒 Security (Password, Biometric)
- 📅 Google Calendar Integration
- 📱 About App
- ❓ Help & Support

#### Google Calendar Integration

**Description**: Enable two-way sync dengan Google Calendar untuk schedules dan appointments.

**Features**:

- **Connect**: OAuth2 authorization dengan Google
- **Sync Options**: Auto-sync atau manual sync
- **Sync Direction**: Bidirectional (CRM ↔ Google Calendar)
- **Privacy**: Token disimpan encrypted di backend

**Flow**:

1. User klik "Connect Google Calendar" di Profile Settings
2. External browser terbuka untuk OAuth authorization
3. Google redirect ke mobile app via deep link (`crmhealth://`)
4. Mobile app exchange authorization code dengan backend
5. Connection established, schedules dapat di-sync

**Implementation**: See [Google Calendar Documentation](../google-calendar/README.md)

### 4. Account Actions

**Actions**:

- Edit Profile
- Change Password
- Logout

---

## Business Rules

### 1. Profile Data Rules

**Read-only Fields**:

- Email (admin only can change)
- Role
- Organization
- Employee ID

**Editable Fields**:

- Full name
- Phone number
- Profile photo
- Address

### 2. Profile Photo Rules

**Upload Requirements**:

- Max size: 2MB
- Format: JPG, PNG
- Min resolution: 200x200px
- Square aspect ratio recommended

### 3. Change Password Rules

**Requirements**:

- Current password required
- New password min 8 characters
- Must contain uppercase, lowercase, number
- Cannot reuse last 3 passwords

### 4. Logout Rules

**Process**:

1. Tampilkan confirmation dialog ("Apakah Anda yakin ingin logout?")
2. Jika confirmed:
   - Clear auth tokens dan user data via `AuthNotifier.logout()`
   - Navigate ke login screen via `pushNamedAndRemoveUntil` (clear seluruh navigation stack)

**Logout Implementation**:

```dart
// Profile Screen - Logout Button Handler
onLogout: () async {
  final confirmed = await showLogoutConfirmation();
  if (confirmed) {
    // 1. Clear auth state
    await ref.read(authProvider.notifier).logout();

    // 2. Navigate ke login dengan clear stack
    // Navigation ini dilakukan DI SINI, bukan di AuthGate
    navigatorKey.currentState?.pushNamedAndRemoveUntil(
      AppRoutes.login,
      (route) => false,
    );
  }
}
```

**Important Changes** (March 2026):

✅ **Dihapus**: Redundant auth state guard di `build()` method profile screen  
✅ **Sekarang**: Logout navigation ditangani oleh logout button handler saja  
✅ **Alasan**: Mencegah multiple navigation conflict dengan `AuthGate`

**Sebelumnya** (Problem):

- Logout button handler → navigate ke login
- Auth state guard di `build()` → juga navigate ke login
- Dashboard's `AuthGate` (di bawah di stack) → render `LoginScreen` inline
- Result: **Navigation glitch**, loading spinner + inline login screen muncul bersamaan

**Sekarang** (Solution):

- Logout button handler → single clean navigation via `pushNamedAndRemoveUntil`
- Profile screen → tidak ada auth guard di `build()`, langsung pop
- `AuthGate` di routes lain → menampilkan loading indicator (bukan login screen)
- Result: **Smooth logout**, tidak ada render conflict

---

## Keputusan Teknis & Trade-offs

### Profile Photo Storage

**Keputusan**: Upload to backend storage, bukan local.

**Alasan**:

- Consistency across devices
- Backup capability
- Shareable profile photos

---

## Struktur Folder

```
apps/mobile/lib/
├── features/
│   └── profile/
│       ├── data/
│       │   ├── models/
│       │   │   └── user_profile_model.dart
│       │   └── profile_repository.dart
│       ├── application/
│       │   ├── profile_provider.dart
│       │   └── settings_provider.dart
│       └── presentation/
│           ├── screens/
│           │   └── profile_screen.dart
│           └── widgets/
│               ├── profile_header.dart
│               ├── performance_card.dart
│               └── settings_menu.dart
```

---

## API Endpoints

#### GET /api/v1/users/me

Get current user profile.

**Response**:

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "+6281234567890",
    "avatar_url": "https://...",
    "role": "sales_rep",
    "organization": {
      "id": "org-uuid",
      "name": "PT Healthcare"
    },
    "team": {
      "id": "team-uuid",
      "name": "Jakarta Team"
    },
    "created_at": "2024-01-15T10:00:00Z"
  }
}
```

#### PUT /api/v1/users/me

Update profile.

**Request**:

```json
{
  "name": "John Doe Updated",
  "phone": "+6281234567890"
}
```

#### POST /api/v1/users/me/avatar

Upload profile photo.

#### POST /api/v1/auth/change-password

Change password.

**Request**:

```json
{
  "current_password": "oldpass",
  "new_password": "newpass123"
}
```

#### GET /api/v1/google-calendar/status

Check Google Calendar connection status.

**Response**:

```json
{
  "success": true,
  "data": {
    "connected": true,
    "email": "user@gmail.com"
  }
}
```

#### GET /api/v1/google-calendar/auth-url

Get OAuth authorization URL.

**Query Parameters**:

```
?platform=mobile
```

**Response**:

```json
{
  "success": true,
  "data": {
    "auth_url": "https://accounts.google.com/o/oauth2/auth?...",
    "state": "eyJ1c2VyX2lkIjoi..."
  }
}
```

#### POST /api/v1/google-calendar/exchange-code

Exchange authorization code for token (mobile only).

**Request**:

```json
{
  "code": "4/0AX4XfW...",
  "state": "eyJ1c2VyX2lkIjoi..."
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "connected": true,
    "message": "Google Calendar connected successfully"
  }
}
```

#### DELETE /api/v1/google-calendar/disconnect

Disconnect Google Calendar.

**Response**:

```json
{
  "success": true,
  "data": {
    "message": "Google Calendar disconnected successfully"
  }
}
```

---

## Cara Test Manual

1. **View Profile**: Verifikasi semua data ditampilkan dengan benar
2. **Edit Profile**: Ubah nama dan phone, verifikasi changes tersimpan
3. **Upload Photo**: Upload new photo, verifikasi update di backend
4. **Change Password**: Test password change flow
5. **Logout**: Test logout dengan confirmation dialog

---

## Dependencies

### Internal Dependencies

- `auth` - Authentication state management
- `google_calendar` - Google Calendar integration (see [Google Calendar Docs](../google-calendar/README.md))

### External Dependencies

- `url_launcher` - Open external browser untuk OAuth
- `app_links` - Handle deep link callbacks

### Related Documentation

- [Google Calendar Integration](../google-calendar/README.md) - Detailed implementation docs
- [Google Calendar Setup Guide](../guides/google-calendar-oauth-setup.md) - Dev1 setup instructions

---

**Document Status**: Active  
**Last Updated**: March 2026
