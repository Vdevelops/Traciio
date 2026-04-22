# Google Calendar OAuth - Mobile Implementation Options

## Overview

Dokumen ini menjelaskan dua pendekatan untuk mengimplementasikan Google Calendar OAuth di aplikasi mobile CRM Healthcare.

---

## Option 1: Backend Redirect to Mobile Deep Link (Current Implementation)

**Status:** ✅ Implemented

### Alur Kerja

```
1. Mobile App          → GET /api/v1/google-calendar/auth-url?platform=mobile
2. Backend             → Generate OAuth URL dengan state (encoded: userID + platform)
3. Browser (External)  → User authorize di Google
4. Google              → Redirect ke Backend callback (HTTPS)
5. Backend             → Process token, redirect ke mobile deep link
6. Mobile App          → Handle deep link crmhealth://google-calendar/callback
```

### Keuntungan

- ✅ Tidak perlu konfigurasi tambahan di Google Cloud Console
- ✅ Works dengan redirect URL HTTPS yang sudah terdaftar
- ✅ Backend tetap handle security (token exchange)
- ✅ Mobile app tidak perlu handle OAuth complexity

### Kekurangan

- User perlu kembali ke app setelah authorize (tapi otomatis via deep link)
- Extra redirect step

### Files yang Dimodifikasi (Option 1)

- `apps/api/internal/api/handlers/google_calendar_auth_handler.go`
- `apps/api/internal/service/google_calendar_token/service.go`
- `apps/mobile/lib/main.dart` (deep link handling)
- `apps/mobile/android/app/src/main/AndroidManifest.xml`
- `apps/mobile/ios/Runner/Info.plist`

---

## Option 2: Direct Deep Link Redirect (Recommended for Future)

**Status:** 📝 Draft - Documentation Only

### Alur Kerja

```
1. Mobile App          → GET /api/v1/google-calendar/auth-url?platform=mobile
2. Backend             → Generate OAuth URL dengan redirect_uri=crmhealth://google-calendar/callback
3. Browser (External)  → User authorize di Google
4. Google              → Direct redirect ke Mobile App (crmhealth://)
5. Mobile App          → Extract auth code dari URL
6. Mobile App          → POST code ke Backend untuk exchange token
```

### Keuntungan

- ✅ Lebih seamless, tidak ada intermediate redirect
- ✅ Lebih cepat (satu redirect less)
- ✅ Mobile app lebih control atas flow
- ✅ Better UX, user langsung kembali ke app

### Kekurangan

- ❌ **Perlu akses Google Cloud Console** untuk tambahkan custom scheme sebagai authorized redirect URI
- ❌ Perlu modify OAuth security model (mobile app handle auth code)
- ❌ Lebih kompleks di mobile side (perlu handle auth code exchange)

### Files yang Perlu Dimodifikasi (Option 2)

#### Backend Changes

- `apps/api/internal/api/handlers/google_calendar_auth_handler.go`
  - Modify `GetAuthURL` to return different redirect_uri for mobile
  - Add new endpoint to exchange code for token (mobile-specific)
- `apps/api/internal/service/google_calendar_token/service.go`
  - Add method to validate and exchange code from mobile
  - Potentially split token exchange logic

#### Google Cloud Console Configuration

- **Required:** Access to Google Cloud Console with project: `1051532630602-lcicb5e4bpldcbjoslmb98abj8o6g5gs`
- **Action:** Add authorized redirect URI:
  ```
  crmhealth://google-calendar/callback
  ```
- **Location:** APIs & Services > Credentials > OAuth 2.0 Client IDs > Web client

#### Mobile Changes

- `apps/mobile/lib/main.dart`
  - Handle deep link to extract auth code
  - Send code to backend for exchange
- `apps/mobile/lib/features/google_calendar/application/google_calendar_provider.dart`
  - Add method to exchange code for token

### Security Considerations for Option 2

⚠️ **Important:** When implementing Option 2, consider:

1. **PKCE (Proof Key for Code Exchange)**
   - Use PKCE extension to prevent authorization code interception
   - Backend must validate code_challenge

2. **State Parameter**
   - Tetap gunakan state parameter untuk CSRF protection
   - State harus di-validate di callback

3. **Code Expiration**
   - Auth code dari Google expired dalam waktu singkat (10 menit)
   - Mobile app harus segera kirim ke backend

### Implementation Plan for Option 2

#### Phase 1: Google Cloud Console Setup (Dev1 Required)

- [ ] Access Google Cloud Console
- [ ] Navigate to project credentials
- [ ] Add `crmhealth://google-calendar/callback` to authorized redirect URIs
- [ ] Save configuration
- [ ] Test OAuth flow di browser untuk verifikasi

#### Phase 2: Backend API Changes

- [ ] Modify `GetAuthURL` handler to check platform parameter
- [ ] For mobile requests, generate URL with custom scheme redirect
- [ ] Create new endpoint: `POST /api/v1/google-calendar/exchange-code`
- [ ] Implement code exchange logic dengan PKCE support
- [ ] Add security validation (state verification)

#### Phase 3: Mobile Implementation

- [ ] Update deep link handler untuk extract auth code
- [ ] Add API call to exchange code
- [ ] Handle error cases (invalid code, expired, dll)
- [ ] Update UI flow (show loading saat exchange)

#### Phase 4: Testing

- [ ] Test happy path: authorize → callback → exchange → success
- [ ] Test error cases: user deny, expired code, network error
- [ ] Test di both Android dan iOS
- [ ] Security testing (CSRF, code interception)

### Code Example for Option 2

#### Backend Handler Modification

```go
// GetAuthURL returns the OAuth2 authorization URL
func (h *GoogleCalendarAuthHandler) GetAuthURL(c *gin.Context) {
    userID := "" // ... get from context

    platform := c.Query("platform")

    // Build state with platform info
    state := buildState(userID, platform)

    // Get OAuth config dengan redirect yang sesuai
    oauthConfig := h.tokenService.GetOAuth2ConfigForPlatform(platform)
    authURL := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)

    response.SuccessResponse(c, map[string]interface{}{
        "auth_url": authURL,
        "state":    state,
    }, nil)
}

// GetOAuth2ConfigForPlatform returns config dengan redirect URI yang sesuai
func (s *Service) GetOAuth2ConfigForPlatform(platform string) *oauth2.Config {
    config := s.GetOAuth2Config()

    if platform == "mobile" {
        config.RedirectURL = "crmhealth://google-calendar/callback"
    }
    // else use default HTTPS redirect

    return config
}

// ExchangeCodeForToken - New endpoint for mobile
func (h *GoogleCalendarAuthHandler) ExchangeCodeForToken(c *gin.Context) {
    var req struct {
        Code  string `json:"code" binding:"required"`
        State string `json:"state" binding:"required"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        errors.BadRequestResponse(c, "INVALID_REQUEST", err.Error())
        return
    }

    // Validate state and extract userID
    userID, platform, err := validateAndExtractState(req.State)
    if err != nil {
        errors.UnauthorizedResponse(c, "INVALID_STATE", "Invalid state parameter")
        return
    }

    if platform != "mobile" {
        errors.BadRequestResponse(c, "INVALID_PLATFORM", "This endpoint is for mobile only")
        return
    }

    // Exchange code for token
    ctx := context.Background()
    token, err := h.tokenService.ExchangeCode(ctx, req.Code)
    if err != nil {
        errors.InternalServerErrorResponse(c, "Failed to exchange code")
        return
    }

    // Store token
    if err := h.tokenService.StoreToken(userID, token); err != nil {
        errors.InternalServerErrorResponse(c, "Failed to store token")
        return
    }

    response.SuccessResponse(c, map[string]interface{}{
        "connected": true,
    }, nil)
}
```

#### Mobile Deep Link Handler

```dart
void _handleDeepLink(Uri uri) {
  if (uri.scheme == 'crmhealth' &&
      uri.host == 'google-calendar' &&
      uri.path == '/callback') {

    final code = uri.queryParameters['code'];
    final state = uri.queryParameters['state'];
    final error = uri.queryParameters['error'];

    if (error != null) {
      // Handle error
      return;
    }

    if (code != null && state != null) {
      // Exchange code for token via backend
      ref.read(googleCalendarProvider.notifier).exchangeCode(code, state);
    }
  }
}
```

### Migration Path from Option 1 to Option 2

Jika ingin migrate dari Option 1 ke Option 2:

1. **Backward Compatibility:**
   - Keep existing endpoints working (don't break web flow)
   - Add new endpoints specifically for mobile direct flow

2. **Gradual Rollout:**
   - Deploy backend changes first (support both flows)
   - Update mobile app dengan feature flag
   - Test dengan subset of users
   - Full rollout setelah stabil

3. **Cleanup:**
   - Setelah semua user update mobile app, bisa remove old flow
   - Atau keep both untuk redundancy

---

## Decision Matrix

| Criteria                    | Option 1 | Option 2                           |
| --------------------------- | -------- | ---------------------------------- |
| **Implementation Speed**    | ⚡ Fast  | 🐢 Slow (need Google Cloud access) |
| **User Experience**         | 😊 Good  | 🤩 Excellent                       |
| **Security**                | 🔒 High  | 🔒 High (with PKCE)                |
| **Maintenance**             | 🟢 Easy  | 🟡 Medium                          |
| **Google Cloud Dependency** | ❌ No    | ✅ Yes (Dev1 needed)               |
| **Testing Complexity**      | 🟢 Low   | 🟡 Medium                          |

---

## Current Recommendation

**Untuk Demo Besok:** Gunakan Option 1 (sudah diimplementasikan)

**Untuk Production:** Pertimbangkan Option 2 untuk UX yang lebih baik, tapi:

- Koordinasi dengan Dev1 untuk Google Cloud Console access
- Allocate waktu untuk proper testing
- Consider security audit untuk mobile token exchange flow

---

## Notes

- **Created:** 2026-03-05
- **Last Updated:** 2026-03-05
- **Author:** AI Assistant
- **Status:** Draft - Awaiting decision on Option 2 implementation

---

## References

- [RFC 7636 - PKCE](https://tools.ietf.org/html/rfc7636)
- [Google OAuth 2.0 for Mobile Apps](https://developers.google.com/identity/protocols/oauth2/native-app)
- [Flutter Deep Linking Documentation](https://docs.flutter.dev/ui/navigation/deep-linking)
- [OAuth 2.0 Security Best Practices](https://tools.ietf.org/html/draft-ietf-oauth-security-topics)
