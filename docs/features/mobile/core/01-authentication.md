# Core - Authentication & Authorization

## CRM Healthcare Mobile App - Flutter

**Module**: Core Infrastructure  
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
6. [API / Package Reference](#api--package-reference)
7. [Configuration](#configuration)
8. [Authentication Flow](#authentication-flow)
9. [Token Management](#token-management)
10. [Cara Test Manual](#cara-test-manual)
11. [Dependencies](#dependencies)
12. [Notes & Improvements](#notes--improvements)

---

## Ringkasan Fitur

Sistem **Authentication & Authorization** mobile app CRM Healthcare menyediakan mekanisme login yang aman dengan JWT (JSON Web Token), automatic token refresh, dan integrasi dengan Role-Based Access Control (RBAC) system. Sistem ini memastikan hanya user yang terautentikasi yang dapat mengakses aplikasi dengan permission yang sesuai.

### Goals

- **Secure Authentication**: Login dengan email/username dan password menggunakan JWT
- **Token Management**: Automatic access token refresh menggunakan refresh token
- **Route Protection**: Auth guards untuk melindungi route yang memerlukan authentication
- **Session Management**: Handle expired tokens dan auto-logout
- **Offline Support**: Cache credentials untuk remember me functionality

---

## Fitur Utama

### 1. Login System

**Endpoint**: `POST /api/v1/auth/login`

- Login dengan email/username dan password
- Response berisi access token (JWT) dan refresh token
- Support "Remember Me" untuk persisten login
- Bilingual error messages (ID/EN)

### 2. JWT Token Management

**Access Token**:

- Short-lived (15-30 menit expiration)
- Bearer token untuk setiap API request
- Disimpan di memory (SharedPreferences/Hive) - NOT secure storage untuk development
- Auto-refresh sebelum expired

**Refresh Token**:

- Long-lived (7-30 hari expiration)
- Rotation mechanism (setiap refresh menghasilkan token baru)
- Stored di database backend untuk revocation support
- Used untuk mendapatkan access token baru

### 3. Automatic Token Refresh

```
┌──────────────┐      ┌──────────────────┐      ┌──────────────┐
│   API Call   │─────►│ Token Expired?   │      │   Backend    │
│   (401)      │      │ (401 Response)   │      │   /refresh   │
└──────────────┘      └──────────────────┘      └──────────────┘
                             │                           │
                             ▼                           ▼
                    ┌──────────────────┐      ┌──────────────┐
                    │ Call Refresh     │─────►│ Generate New │
                    │ Token Endpoint   │      │ Access Token │
                    └──────────────────┘      └──────────────┘
                             │                           │
                             ▼                           ▼
                    ┌──────────────────┐      ┌──────────────┐
                    │ Retry Original   │◄─────│ Return New   │
                    │ Request with New │      │ Token        │
                    │ Access Token     │      │              │
                    └──────────────────┘      └──────────────┘
```

### 4. Auth Guards & Route Protection

- **AuthGate Widget**: Wrapper untuk protected screens
- **Redirect Logic**: Redirect ke login screen jika tidak terautentikasi
- **Token Validation**: Validasi token sebelum mengakses route
- **Permission Checks**: Integrasi dengan RBAC untuk permission validation

### 5. Logout System

**Trigger Points**:

- **Manual Logout**: Tombol logout di profile screen
- **Auto Logout**: Token refresh gagal (401 dari API interceptor)

**Implementation**:

- Clear auth tokens dan user data dari local storage
- Set `AuthState` ke `unauthenticated`
- Navigate ke login screen via `pushNamedAndRemoveUntil` (clear seluruh nav stack)
- Untuk interceptor-triggered logout, menggunakan global `navigatorKey` dari `MyApp` agar bisa navigate dari non-widget context

```dart
// Profile screen logout
await ref.read(authProvider.notifier).logout();
if (context.mounted) {
  Navigator.of(context).pushNamedAndRemoveUntil(
    AppRoutes.login,
    (route) => false,
  );
}

// API interceptor logout (via global navigatorKey)
await authNotifier.logout();
MyApp.navigatorKey.currentState?.pushNamedAndRemoveUntil(
  AppRoutes.login,
  (route) => false,
);
```

**Important**: Profile screen juga memiliki guard di `build()` — jika `authState.status == AuthStatus.unauthenticated`, otomatis redirect ke login untuk menghindari error dari providers yang depend on auth state.

---

## Business Rules

### 1. Login Requirements

- **Email/Username**: Required, case-insensitive
- **Password**: Required, minimum 6 characters (sesuai backend validation)
- **Organization**: Optional (untuk multi-tenant support)

### 2. Token Lifecycle

| Token Type    | Expiration    | Storage                           | Refresh                     |
| ------------- | ------------- | --------------------------------- | --------------------------- |
| Access Token  | 15-30 minutes | SharedPreferences/Hive            | Manual via refresh endpoint |
| Refresh Token | 7-30 days     | SharedPreferences/Hive + Database | Automatic rotation          |

### 3. Authentication States

```dart
enum AuthStatus {
  unknown,         // App startup, checking stored tokens
  authenticated,   // User logged in, valid tokens
  unauthenticated, // User logged out / tokens expired
}
```

**Note**: Loading dan error states di-handle via `isLoading` dan `errorMessage` fields di `AuthState`, bukan via `AuthStatus` enum. Ini memudahkan kombinasi state (misal: authenticated tapi sedang loading profile update).

### 4. Error Handling

**Bilingual Error Messages**:

- Indonesia: "Email atau password salah"
- English: "Invalid email or password"

**Error Types**:

- Invalid credentials (401)
- Account suspended/disabled (403)
- Token expired (401 dengan specific error code)
- Network error (connection timeout)
- Server error (500)

### 5. Security Rules

- **Token Storage**: Jangan hardcode tokens di code
- **HTTPS Only**: All API calls menggunakan HTTPS
- **Auto-logout**: Logout otomatis jika refresh token invalid
- **Concurrent Sessions**: Backend mendukung multiple sessions per user

---

## Keputusan Teknis & Trade-offs

### Mengapa SharedPreferences/Hive, bukan Flutter Secure Storage?

**Keputusan**: Menggunakan SharedPreferences/Hive untuk token storage, bukan flutter_secure_storage.

**Alasan**:

- **Simplicity**: Lebih mudah untuk development dan testing
- **Performance**: Faster read/write operations
- **Cross-platform**: Consistent behavior di Android dan iOS
- **Development Speed**: Tidak memerlukan platform-specific setup

**Trade-off**: Kurang secure dibanding flutter_secure_storage. **Mitigasi**: Untuk production, migrasi ke flutter_secure_storage atau encrypted Hive.

### Mengapa Riverpod StateNotifier untuk Auth State?

**Keputusan**: Menggunakan StateNotifierProvider daripada ChangeNotifier atau Bloc.

**Alasan**:

- **Type Safety**: Compile-time safety untuk state transitions
- **Testability**: Mudah di-unit test
- **Performance**: Efficient rebuild dengan selective listeners
- **Consistency**: Pattern yang sama digunakan di seluruh app

### Token Refresh Strategy

**Keputusan**: Reactive token refresh (refresh on 401) daripada proactive (refresh before expiry).

**Alasan**:

- **Simplicity**: Tidak perlu timer/scheduler
- **Efficiency**: Hanya refresh saat diperlukan
- **Reliability**: Handle edge cases lebih baik (misal: token invalidated di backend)

**Trade-off**: Sedikit latency saat token expired (satu request gagal, kemudian refresh dan retry). **Mitigasi**: Implement queue untuk concurrent requests saat refreshing.

---

## Struktur Folder

```
apps/mobile/lib/
├── features/
│   └── auth/
│       ├── data/
│       │   ├── models/
│       │   │   ├── auth_token.dart           # Token model (access_token, refresh_token, expires_at)
│       │   │   ├── login_request.dart        # Login request DTO
│       │   │   └── login_response.dart       # Login response DTO
│       │   └── auth_repository.dart          # API calls (login, logout, refresh)
│       ├── application/
│       │   ├── auth_provider.dart            # StateNotifier untuk auth state
│       │   └── auth_state.dart               # AuthState class dengan status
│       └── presentation/
│           ├── screens/
│           │   └── login_screen.dart         # Login UI
│           └── widgets/
│               └── auth_guard.dart           # Route protection widget
├── core/
│   ├── network/
│   │   ├── api_client.dart                   # Dio client dengan interceptors
│   │   └── auth_interceptor.dart             # Automatic token injection
│   └── storage/
│       ├── token_storage.dart                # Token persistence
│       └── secure_storage.dart               # Encrypted storage (optional)
```

---

## API / Package Reference

### Authentication Endpoints

#### POST /api/v1/auth/login

Login dengan credentials.

**Request Body**:

```json
{
  "email": "sales@example.com",
  "password": "password123",
  "organization_code": "ORG001"
}
```

**Response (Success)**:

```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "email": "sales@example.com",
      "name": "John Doe",
      "role": "sales_rep"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 1800
  },
  "timestamp": "2025-01-15T10:30:45+07:00",
  "request_id": "req_abc123xyz"
}
```

**Response (Error)**:

```json
{
  "success": false,
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Email atau password salah",
    "details": {}
  },
  "timestamp": "2025-01-15T10:30:45+07:00"
}
```

#### POST /api/v1/auth/refresh

Refresh access token menggunakan refresh token.

**Request Body**:

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response**:

```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 1800
  }
}
```

#### POST /api/v1/auth/logout

Logout dan revoke refresh token.

**Headers**:

```
Authorization: Bearer <access_token>
```

**Response**:

```json
{
  "success": true,
  "data": {
    "message": "Logout berhasil"
  }
}
```

### Mobile-Specific Endpoints

#### GET /api/v1/auth/mobile/permissions

Get user permissions untuk mobile app.

**Headers**:

```
Authorization: Bearer <access_token>
```

**Response**:

```json
{
  "success": true,
  "data": {
    "menus": [
      {
        "menu": "dashboard",
        "actions": ["VIEW"]
      },
      {
        "menu": "task",
        "actions": ["VIEW", "CREATE", "EDIT", "DELETE"]
      }
    ]
  }
}
```

---

## Configuration

### Environment Variables

| Variable                  | Default                 | Description                            |
| ------------------------- | ----------------------- | -------------------------------------- |
| `API_BASE_URL`            | `http://localhost:8080` | Base URL untuk backend API             |
| `TOKEN_STORAGE_KEY`       | `auth_tokens`           | Key untuk token storage                |
| `REFRESH_TOKEN_THRESHOLD` | `300`                   | Refresh token jika kurang dari X detik |

### pubspec.yaml Dependencies

```yaml
dependencies:
  dio: ^5.0.0
  flutter_riverpod: ^2.4.0
  shared_preferences: ^2.2.0
  hive: ^2.2.3
  hive_flutter: ^1.1.0
  jwt_decoder: ^2.0.1
```

---

## Authentication Flow

### Login Flow

```
┌──────────────┐
│  Login Screen │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Validate    │
│  Input       │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  POST /login │
└──────┬───────┘
       │
       ├──────────────────┐
       │                  │
       ▼                  ▼
┌──────────────┐  ┌──────────────┐
│   Success    │  │    Error     │
└──────┬───────┘  └──────┬───────┘
       │                  │
       ▼                  ▼
┌──────────────┐  ┌──────────────┐
│ Save Tokens  │  │ Show Error   │
│ to Storage   │  │ Message      │
└──────┬───────┘  └──────────────┘
       │
       ▼
┌──────────────┐
│ Navigate to  │
│ Dashboard    │
└──────────────┘
```

### Token Refresh Flow

**Token Refresh dengan AuthGate Navigation**:

```
┌──────────────┐
│  API Call    │
│  Returns 401 │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ Check if     │
│ Refresh Token│
│ Exists       │
└──────┬───────┘
       │
       ├──────────┐
       │          │
       ▼          ▼
┌──────────┐ ┌──────────┐
│   Yes    │ │    No    │
└────┬─────┘ └────┬─────┘
     │            │
     ▼            │
┌──────────┐      │
│ POST     │      │
│ /refresh │      │
└────┬─────┘      │
     │            │
     ├──────────┐ │
     │          │ │
     ▼          ▼ │
┌──────────┐ ┌──────────┐
│ Success  │ │  Error   │
└────┬─────┘ └────┬─────┘
     │            │
     ▼            │
┌──────────┐      │
│ Update   │      │
│ Tokens   │      │
└────┬─────┘      │
     │            │
     ▼            │
┌──────────┐      │
│ Retry    │      │
│ Original │      │
│ Request  │      │
└──────────┘      │
                  │
                  ▼
         ┌──────────────────┐
         │ AuthGate detects │
         │ unauthenticated  │
         │ (saat refresh    │
         │ token invalid)   │
         └────────┬─────────┘
                  │
                  ▼
         ┌──────────────────┐
         │ Set flag:        │
         │ _isNavigating    │
         │ ToLogin = true   │
         └────────┬─────────┘
                  │
                  ▼
         ┌──────────────────┐
         │ Show loading     │
         │ indicator        │
         │ (bukan login     │
         │ screen inline!)  │
         └────────┬─────────┘
                  │
                  ▼
         ┌──────────────────┐
         │ Schedule         │
         │ pushNamedAnd     │
         │ RemoveUntil      │
         └────────┬─────────┘
                  │
                  ▼
         ┌──────────────────┐
         │ Login Screen     │
         └──────────────────┘
```

**Key Improvements**:

1. **No Inline LoginScreen**: `AuthGate` tidak lagi render `LoginScreen()` inline
2. **Loading Indicator**: Menampilkan loading spinner saat menunggu navigasi
3. **Flag Protection**: `_isNavigatingToLogin` mencegah multiple navigasi
4. **Clean Navigation**: `pushNamedAndRemoveUntil` membersihkan stack sepenuhnya

---

## Token Management

### Token Storage Strategy

**Secure Storage (Recommended untuk Production)**:

```dart
class TokenStorage {
  static const String _accessTokenKey = 'access_token';
  static const String _refreshTokenKey = 'refresh_token';

  Future<void> saveTokens(String accessToken, String refreshToken) async {
    await _storage.write(key: _accessTokenKey, value: accessToken);
    await _storage.write(key: _refreshTokenKey, value: refreshToken);
  }

  Future<String?> getAccessToken() async {
    return await _storage.read(key: _accessTokenKey);
  }

  Future<void> clearTokens() async {
    await _storage.deleteAll();
  }
}
```

**Hive Storage (Current Implementation)**:

```dart
class TokenStorage {
  static const String _boxName = 'auth_tokens';

  Future<void> saveTokens(String accessToken, String refreshToken) async {
    final box = await Hive.openBox(_boxName);
    await box.put('access_token', accessToken);
    await box.put('refresh_token', refreshToken);
  }
}
```

### Token Validation

```dart
bool isTokenValid(String token) {
  try {
    final decodedToken = JwtDecoder.decode(token);
    final expirationDate = JwtDecoder.getExpirationDate(token);
    return expirationDate.isAfter(DateTime.now());
  } catch (e) {
    return false;
  }
}

bool shouldRefreshToken(String token) {
  try {
    final expirationDate = JwtDecoder.getExpirationDate(token);
    final timeUntilExpiry = expirationDate.difference(DateTime.now());
    return timeUntilExpiry.inSeconds < 300; // Refresh jika < 5 menit
  } catch (e) {
    return true;
  }
}
```

---

## Cara Test Manual

### Test Login Flow

1. **Successful Login**:
   - Masukkan email dan password yang valid
   - Tap "Login"
   - Verifikasi: Navigate ke dashboard
   - Verifikasi: Tokens tersimpan di storage

2. **Failed Login - Invalid Credentials**:
   - Masukkan email dan password yang salah
   - Tap "Login"
   - Verifikasi: Error message muncul (ID/EN)
   - Verifikasi: Tetap di login screen

3. **Failed Login - Network Error**:
   - Matikan network/WiFi
   - Tap "Login"
   - Verifikasi: Network error message muncul

### Test Token Refresh

1. **Automatic Refresh**:
   - Login dan tunggu sampai access token hampir expired
   - Lakukan API call
   - Verifikasi: Request berhasil tanpa manual login
   - Verifikasi: New tokens tersimpan

2. **Refresh Token Expired**:
   - Wait sampai refresh token expired (atau revoke dari backend)
   - Lakukan API call
   - Verifikasi: Navigate ke login screen
   - Verifikasi: Error message muncul

### Test Logout

1. **Logout Success** (Improved Flow):
   - Login sebagai user
   - Navigate ke profile/settings
   - Tap "Logout"
   - **Verifikasi Baru**:
     - Tidak ada glitch/render conflict
     - Profile screen langsung pop (tidak ada loading spinner lama)
     - Smooth transition ke login screen
     - Navigation stack sepenuhnya di-clear
   - Verifikasi: Tokens dihapus dari storage

2. **Logout with Network Error**:
   - Matikan network
   - Tap "Logout"
   - Verifikasi: Tetap logout (clear local tokens)
   - Verifikasi: Navigate ke login screen
   - **Verifikasi Baru**: Navigation tetap smooth meskipun offline

3. **Logout dari Protected Screen**:
   - Buka screen yang menggunakan `AuthGate` (misal: Dashboard, Accounts)
   - Tap logout dari profile (di-push di atas screen tersebut)
   - **Verifikasi**:
     - Tidak ada `LoginScreen` yang render inline di `AuthGate`
     - Hanya ada satu navigation action
     - Tidak ada multiple loading indicators

4. **Rapid Logout/Login**:
   - Login → Logout → Login lagi (cepat)
   - **Verifikasi**: Tidak ada race condition atau navigation error

---

## Dependencies

### Internal

- `features/permissions/` - RBAC integration
- `core/network/api_client.dart` - Dio client configuration
- `core/storage/offline_storage.dart` - Token persistence

### External

- `dio: ^5.0.0` - HTTP client dengan interceptors
- `flutter_riverpod: ^2.4.0` - State management
- `hive: ^2.2.3` - Local storage
- `jwt_decoder: ^2.0.1` - JWT token decoding

---

## Notes & Improvements

### Recent Fixes

**Logout Navigation Glitch (March 2026)** ✅ **FIXED**

**Problem**: Saat logout, multiple navigation sources bereaksi bersamaan:

- Profile screen logout handler → navigate ke login
- Profile screen auth guard di `build()` → juga navigate ke login
- Dashboard's `AuthGate` → render `LoginScreen` inline
- Result: UI glitch dengan loading spinner + inline login + navigation conflict

**Solution**:

- Dihapus: Redundant auth state guard di profile screen `build()`
- Diubah: `AuthGate` tidak lagi render `LoginScreen()` inline
- Ditambahkan: Static flag `_isNavigatingToLogin` untuk mencegah multiple navigasi
- Ditambahkan: Loading indicator di `AuthGate` saat menunggu navigasi
- Sekarang: Single clean navigation path dari logout button handler saja

**Files Changed**:

- `core/widgets/auth_gate.dart`
- `features/profile/presentation/profile_screen.dart`
- `core/routing/app_router.dart`

### Known Limitations

1. **Token Storage Security**: Menggunakan Hive/SharedPreferences (tidak encrypted). Untuk production, gunakan flutter_secure_storage atau encrypted Hive.

2. **Single Device Login**: Tidak ada mekanisme untuk invalidate sessions di device lain. User dapat login di multiple devices.

3. **No Biometric Auth**: Belum implementasi fingerprint/Face ID untuk quick login.

### Future Improvements

1. **Biometric Authentication**: Tambahkan fingerprint/Face ID untuk quick login
2. **PIN Login**: Alternative login method dengan 6-digit PIN
3. **Session Management UI**: Show active sessions dan allow remote logout
4. **Social Login**: Google/Apple sign-in integration
5. **Two-Factor Authentication**: SMS/Email OTP untuk enhanced security
6. **Password Reset**: Self-service password reset via email

### Security Recommendations untuk Production

1. Gunakan `flutter_secure_storage` untuk token storage
2. Implement certificate pinning untuk API calls
3. Encrypt sensitive data di local storage
4. Add root/jailbreak detection
5. Implement app attestation (DeviceCheck/App Attest)

---

**Document Status**: Active  
**Last Updated**: March 2026  
**Maintained By**: Dev3 (Mobile Development Team)
