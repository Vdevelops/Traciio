# Google Calendar OAuth - Mobile Implementation

## Overview

Dokumentasi implementasi Google Calendar OAuth untuk aplikasi mobile CRM Healthcare menggunakan pendekatan **HTTPS Web Redirect + Server Forward (Option 1 - Recommended)**.

**Status:** ✅ Implemented  
**Date:** 2026-03-05  
**Type:** Mobile OAuth Integration  
**Architecture:** HTTPS Redirect with Backend Token Exchange

---

## Architecture

### Flow Diagram

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│ Mobile App  │────▶│   Backend   │────▶│   Google    │────▶│   Backend   │────▶│ Mobile App  │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
      │                   │                   │                   │                   │
      │ GET /auth-url     │                   │                   │                   │
      │ ?platform=mobile  │                   │                   │                   │
      │──────────────────▶│                   │                   │                   │
      │                   │                   │                   │                   │
      │   {auth_url}      │                   │                   │                   │
      │◀──────────────────│                   │                   │                   │
      │                   │                   │                   │                   │
      │ Open Browser      │                   │                   │                   │
      │ (External)        │                   │                   │                   │
      │──────────────────────────────────────▶│                   │                   │
      │                   │                   │                   │                   │
      │                   │                   │ User Authorize    │                   │
      │                   │                   │                   │                   │
      │                   │                   │ Redirect to       │                   │
      │                   │                   │ HTTPS URL         │                   │
      │                   │◀──────────────────────────────────────│                   │
      │                   │                   │                   │                   │
      │                   │ Exchange Code     │                   │                   │
      │                   │ Store Token       │                   │                   │
      │                   │                   │                   │                   │
      │                   │ 302 Redirect      │                   │                   │
      │◀──────────────────────────────────────────────────────────────────────────────│
      │                   │ crmhealth://      │                   │                   │
      │                   │ ?success=true     │                   │                   │
```

### Keuntungan Pendekatan Ini

1. **Secure**: Menggunakan HTTPS redirect yang terverifikasi oleh Google
2. **Server Forward**: Backend menangani token exchange, tidak expose ke client
3. **Single Redirect URI**: Web dan mobile menggunakan URL yang sama
4. **Production Ready**: Sesuai rekomendasi Google untuk aplikasi production

---

## Implementation Details

### Backend Changes

#### 1. Service Layer (`apps/api/internal/service/google_calendar_token/service.go`)

**New Methods:**

- `GetOAuth2ConfigForPlatform(platform string)` - Returns OAuth config dengan HTTPS redirect URI
- `GetAuthURLForPlatform(state, platform string)` - Generate auth URL dengan HTTPS redirect

```go
func (s *Service) GetOAuth2ConfigForPlatform(platform string) *oauth2.Config {
    config := s.GetOAuth2Config()
    if platform == "mobile" {
        // Mobile menggunakan HTTPS redirect yang sama dengan web
        // Backend akan handle exchange dan forward ke mobile
        config.RedirectURL = "https://api.gilabs.id/api/v1/google-calendar/callback"
    }
    return config
}
```

#### 2. Handler Layer (`apps/api/internal/api/handlers/google_calendar_auth_handler.go`)

**Modified:**

- `GetAuthURL()` - Menggunakan `GetAuthURLForPlatform()` dengan HTTPS redirect
- `HandleCallback()` - Handle mobile flow dengan exchange code di backend

**Key Changes:**

- Mobile: Backend exchange code → store token → redirect ke mobile app
- Web: Redirect ke frontend dengan code (frontend handle exchange)

```go
func (h *GoogleCalendarAuthHandler) HandleCallback(c *gin.Context) {
    code := c.Query("code")
    stateStr := c.Query("state")

    oauthState := decodeState(stateStr)

    if oauthState.Platform == "mobile" {
        // Mobile: Exchange code immediately and redirect to app
        ctx := context.Background()
        token, err := h.tokenService.HandleOAuth2Callback(ctx, code)

        if err := h.tokenService.StoreToken(oauthState.UserID, token); err != nil {
            // Redirect to mobile with error
            redirectURL := buildRedirectURL(oauthState, "failed_to_store_token", "")
            c.Redirect(http.StatusFound, redirectURL)
            return
        }

        // Redirect to mobile app with success
        redirectURL := buildRedirectURL(oauthState, "", "success")
        c.Redirect(http.StatusFound, redirectURL)
        return
    }

    // Web: Redirect to frontend (frontend will handle exchange)
    frontendURL := getFrontendURL(c)
    redirectURL := frontendURL + "/google-calendar/callback?code=" + code + "&state=" + stateStr
    c.Redirect(http.StatusFound, redirectURL)
}
```

### Mobile Changes

#### 1. Deep Link Handler (`apps/mobile/lib/main.dart`)

**Simplified Flow**: Backend sudah handle token exchange, mobile hanya perlu refresh status.

```dart
void _handleGoogleCalendarCallback(Uri uri) async {
  final success = uri.queryParameters['success'] == 'true';
  final error = uri.queryParameters['error'];

  if (error != null) {
    // Show error message
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('Error: $error'), backgroundColor: Colors.red),
    );
    return;
  }

  if (success) {
    // Backend already handled token exchange, just refresh status
    await ref.read(googleCalendarNotifierProvider.notifier).refreshStatus();

    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(
        content: Text('Google Calendar connected successfully!'),
        backgroundColor: Colors.green,
      ),
    );
  }
}
```

**Key Point**: Mobile tidak perlu lagi call `exchangeCode()` karena backend sudah handle.

```dart
Future<void> exchangeCode(String code, String state) async {
  final response = await _dio.post(
    '/api/v1/google-calendar/exchange-code',
    data: {
      'code': code,
      'state': state,
    },
  );
  // Handle response
}
```

#### 4. Provider Layer

```dart
Future<bool> exchangeCode(String code, String oauthState) async {
  this.state = const AsyncValue.loading();
  try {
    await _repository.exchangeCode(code, oauthState);
    await refreshStatus();
    return true;
  } catch (e, stackTrace) {
    this.state = AsyncValue.error(e, stackTrace);
    return false;
  }
}
```

---

## API Endpoints

### 1. Get Auth URL

**Request:**

```http
GET /api/v1/google-calendar/auth-url?platform=mobile
Authorization: Bearer <token>
```

**Response:**

```json
{
  "success": true,
  "data": {
    "auth_url": "https://accounts.google.com/o/oauth2/auth?client_id=...&redirect_uri=https://api.gilabs.id/api/v1/google-calendar/callback&state=...",
    "state": "eyJ1c2VyX2lkIjoi..."
  }
}
```

**Note**: `redirect_uri` adalah HTTPS URL (sama untuk web dan mobile). Backend akan mendeteksi platform dari state parameter.

---

## Configuration

### Required Google Cloud Console Setup

⚠️ **Must be done by Dev1 before testing:**

**Method**: HTTPS Web Redirect + Server Forward (Recommended by Google)

1. Login ke [Google Cloud Console](https://console.cloud.google.com/)
2. Navigate ke project: `1051532630602-lcicb5e4bpldcbjoslmb98abj8o6g5gs`
3. Menu: **APIs & Services > Credentials**
4. Edit **OAuth 2.0 Client ID** (Web application)
5. Add **Authorized Redirect URI**:
   ```
   https://api.gilabs.id/api/v1/google-calendar/callback
   ```
6. Click **SAVE**

**Note**:

- Web dan mobile menggunakan URL yang **sama** (HTTPS)
- Tidak perlu register custom scheme (`crmhealth://`) di Google Cloud Console
- Backend akan handle forwarding ke mobile app

### Mobile Deep Link Configuration

#### Android (`android/app/src/main/AndroidManifest.xml`)

```xml
<intent-filter android:autoVerify="true">
    <action android:name="android.intent.action.VIEW" />
    <category android:name="android.intent.category.DEFAULT" />
    <category android:name="android.intent.category.BROWSABLE" />
    <data android:scheme="crmhealth" android:path="/google-calendar/callback" />
</intent-filter>
```

#### iOS (`ios/Runner/Info.plist`)

```xml
<key>CFBundleURLTypes</key>
<array>
    <dict>
        <key>CFBundleURLName</key>
        <string>id.gilabs.crmhealthcare</string>
        <key>CFBundleURLSchemes</key>
        <array>
            <string>crmhealth</string>
        </array>
    </dict>
</array>
```

---

## Files Modified

### Backend

- `apps/api/internal/service/google_calendar_token/service.go`
- `apps/api/internal/api/handlers/google_calendar_auth_handler.go`
- `apps/api/internal/api/routes/schedule_routes.go`

### Mobile

- `apps/mobile/lib/main.dart`
- `apps/mobile/lib/features/google_calendar/application/google_calendar_provider.dart`
- `apps/mobile/lib/features/google_calendar/domain/repositories/google_calendar_repository.dart`
- `apps/mobile/lib/features/google_calendar/data/repositories/google_calendar_repository_impl.dart`
- `apps/mobile/lib/features/google_calendar/data/datasources/google_calendar_remote_datasource.dart`
- `apps/mobile/lib/features/profile/presentation/profile_screen.dart`

---

## Testing

### Prerequisites

- [ ] Google Cloud Console configured with redirect URI
- [ ] Backend deployed dengan latest code
- [ ] Mobile app rebuild dengan latest code

### Test Cases

#### Happy Path

1. Buka Profile screen
2. Klik "Connect Google Calendar"
3. Authorize di browser
4. Verify redirect ke mobile app
5. Verify success message
6. Verify status "Connected"

#### Error Cases

- User deny permission
- Invalid state parameter
- Expired authorization code
- Network error during exchange

---

## Security

### Implemented

- ✅ State parameter dengan base64 encoding
- ✅ User ID validation di state
- ✅ Platform detection
- ✅ Code expiration handled by Google

### Considerations

- Custom scheme bisa di-hijack (gunakan HTTPS App Links untuk production)
- State parameter wajib divalidasi
- Code harus segera ditukar (10 menit expiration)

---

## Troubleshooting

### Common Issues

**Issue:** Deep link tidak terbuka

- **Check:** AndroidManifest.xml / Info.plist configuration
- **Fix:** Pastikan intent-filter dan URL scheme benar

**Issue:** "Invalid redirect_uri" error

- **Check:** Google Cloud Console redirect URI list
- **Fix:** Tambahkan `crmhealth://google-calendar/callback`

**Issue:** Exchange code failed

- **Check:** Network connectivity
- **Check:** State parameter valid
- **Fix:** Verifikasi state encoding/decoding

---

## References

- [RFC 7636 - PKCE](https://tools.ietf.org/html/rfc7636)
- [Google OAuth 2.0 for Mobile Apps](https://developers.google.com/identity/protocols/oauth2/native-app)
- [Flutter Deep Linking](https://docs.flutter.dev/ui/navigation/deep-linking)

---

## Author

- **Implementation:** AI Assistant
- **Date:** 2026-03-05
- **Version:** 1.0.0
