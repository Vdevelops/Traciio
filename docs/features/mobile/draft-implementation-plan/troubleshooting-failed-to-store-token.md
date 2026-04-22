# Troubleshooting Guide: Google Calendar OAuth "failed_to_store_token" Error

## Overview

Dokumen ini menjelaskan penyebab dan solusi untuk error "failed_to_store_token" yang terjadi saat implementasi Google Calendar OAuth di mobile app.

**Status:** Error Analysis  
**Date:** 2026-03-07  
**Related Feature:** Google Calendar OAuth Integration  
**Severity:** HIGH - Blocking OAuth flow

---

## Error Description

### Symptoms

- User berhasil login ke Google (OAuth authorization berhasil)
- User di-redirect kembali ke mobile app
- Mobile app menerima error: `crmhealth://google-calendar/callback?error=failed_to_store_token`
- Status Google Calendar tetap "Not Connected"
- Toast menampilkan: "Failed to store token"

### Error Flow

```
1. Mobile → GET /auth-url (success) ✓
2. Browser → User authorize di Google ✓
3. Google → Redirect ke backend callback ✓
4. Backend → Exchange code → token ✓
5. Backend → StoreToken(userID, token) ✗ FAILED
6. Backend → Redirect ke mobile dengan error
```

---

## Root Cause Analysis

Error terjadi di **step 5** - Backend gagal menyimpan token ke database setelah berhasil exchange dari Google OAuth.

### Kemungkinan Penyebab

#### 1. Database Table Tidak Ada (Most Likely)

**Deskripsi:** Tabel `google_calendar_tokens` belum dibuat di database.

**Evidence:**

```
Error: relation "google_calendar_tokens" does not exist
```

**Solution:**

```sql
-- Cek apakah tabel ada
\d google_calendar_tokens;

-- Jika tidak ada, jalankan migration
-- Atau buat tabel manual:
CREATE TABLE google_calendar_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    access_token TEXT NOT NULL,
    refresh_token TEXT NOT NULL,
    token_type VARCHAR(50) DEFAULT 'Bearer',
    expires_at TIMESTAMP,
    scope TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id)
);
```

#### 2. Encryption Key Invalid

**Deskripsi:** Environment variable `ENCRYPTION_KEY` tidak di-set atau tidak valid.

**Evidence:**

```
[Google Calendar Service] StoreToken - ERROR: Failed to encrypt access token: crypto/aes: invalid key size
```

**Solution:**

```bash
# Check environment variable
env | grep ENCRYPTION_KEY

# Set key yang valid (32 bytes untuk AES-256)
export ENCRYPTION_KEY="your-32-byte-encryption-key-here"

# Atau di .env file:
ENCRYPTION_KEY=0123456789abcdef0123456789abcdef
```

#### 3. Database Connection Error

**Deskripsi:** Backend tidak bisa connect ke database.

**Evidence:**

```
[Google Calendar Service] StoreToken - ERROR: Failed to find existing token: connection refused
```

**Solution:**

```bash
# Check database connection
psql $DATABASE_URL -c "SELECT 1;"

# Check database status
sudo systemctl status postgresql
# atau
docker ps | grep postgres
```

#### 4. Foreign Key Constraint Violation

**Deskripsi:** User ID tidak ditemukan di tabel users.

**Evidence:**

```
ERROR: insert or update on table "google_calendar_tokens" violates foreign key constraint "fk_user_id"
```

**Solution:**

```sql
-- Verifikasi user ada di database
SELECT id, email FROM users WHERE id = 'fd247498-8a67-41c0-acd9-2df6cf38dccd';

-- Jika tidak ada, ada masalah dengan auth flow
```

#### 5. Database Permission Error

**Deskripsi:** Database user tidak punya permission untuk write.

**Evidence:**

```
ERROR: permission denied for table google_calendar_tokens
```

**Solution:**

```sql
-- Grant permissions
GRANT ALL PRIVILEGES ON TABLE google_calendar_tokens TO your_db_user;
GRANT USAGE, SELECT ON SEQUENCE google_calendar_tokens_id_seq TO your_db_user;
```

---

## Diagnostic Steps

### Step 1: Check Backend Logs

```bash
# SSH ke server
ssh user@staging-server

# Check logs real-time
docker logs <api-container> -f | grep -E "(StoreToken|Google Calendar)"

# Atau check log file
tail -f /var/log/api.log | grep "failed_to_store_token"
```

### Step 2: Verify Environment Variables

```bash
# Check all required env vars
echo "ENCRYPTION_KEY: $ENCRYPTION_KEY"
echo "DATABASE_URL: $DATABASE_URL"
echo "GOOGLE_CALENDAR_REDIRECT_URL: $GOOGLE_CALENDAR_REDIRECT_URL"

# Verify encryption key length
echo -n "$ENCRYPTION_KEY" | wc -c
# Should output: 32 (for AES-256)
```

### Step 3: Verify Database Schema

```sql
-- Connect to database
psql $DATABASE_URL

-- List all tables
\dt

-- Check table structure
\d google_calendar_tokens

-- Check if table has data
SELECT COUNT(*) FROM google_calendar_tokens;

-- Check for specific user
SELECT * FROM google_calendar_tokens WHERE user_id = 'fd247498-8a67-41c0-acd9-2df6cf38dccd';
```

### Step 4: Run Database Migration

```bash
# If using GORM auto-migration, restart backend
# If using manual migration:
psql $DATABASE_URL < migrations/001_create_google_calendar_tokens.sql

# Or use migration tool
make migrate-up
# atau
pnpm migrate:up
```

---

## Log Analysis Examples

### Good Log (Success Flow)

```
[Google Calendar Service] StoreToken - userID: fd247498-8a67-41c0-acd9-2df6cf38dccd, AccessToken length: 2048, RefreshToken length: 1024
[Google Calendar Service] StoreToken - encrypting tokens...
[Google Calendar Service] StoreToken - tokens encrypted successfully
[Google Calendar Service] StoreToken - checking for existing token...
[Google Calendar Service] StoreToken - creating new token...
[Google Calendar Service] StoreToken - token created successfully
[Google Calendar Handler] StoreToken SUCCESS for user fd247498-8a67-41c0-acd9-2df6cf38dccd
```

### Bad Log (Error Flow)

```
[Google Calendar Service] StoreToken - userID: fd247498-8a67-41c0-acd9-2df6cf38dccd, AccessToken length: 2048, RefreshToken length: 1024
[Google Calendar Service] StoreToken - encrypting tokens...
[Google Calendar Service] StoreToken - ERROR: Failed to encrypt access token: crypto/aes: invalid key size
[Google Calendar Handler] StoreToken ERROR for user fd247498-8a67-41c0-acd9-2df6cf38dccd: encryption failed
```

---

## Solutions by Environment

### Local Development

1. **Check if database is running:**

```bash
docker ps | grep postgres
```

2. **Run migrations:**

```bash
cd apps/api
make migrate-up
```

3. **Check .env file:**

```bash
cat .env | grep -E "(ENCRYPTION_KEY|DATABASE_URL)"
```

4. **Restart backend:**

```bash
pnpm dev
```

### Staging Environment

1. **Check deployment status:**

```bash
kubectl get pods -n staging
# atau
docker ps -a | grep staging
```

2. **Check environment variables in container:**

```bash
docker exec <staging-container> env | grep ENCRYPTION_KEY
```

3. **Run migration in staging:**

```bash
kubectl exec -it <staging-pod> -- ./migrate up
# atau
docker exec <staging-container> ./migrate up
```

4. **Check logs:**

```bash
kubectl logs -f <staging-pod> | grep "Google Calendar"
# atau
docker logs -f <staging-container> | grep "failed_to_store_token"
```

### Production Environment

⚠️ **WARNING: Be careful with production!**

1. **Backup database first:**

```bash
pg_dump $DATABASE_URL > backup_$(date +%Y%m%d_%H%M%S).sql
```

2. **Check current status:**

```bash
# Read logs only, no changes
kubectl logs <prod-pod> --tail=100 | grep "Google Calendar"
```

3. **Schedule maintenance window for migration**

---

## Prevention

### 1. Pre-deployment Checklist

- [ ] Database migration dijalankan
- [ ] Environment variables ter-set dengan benar
- [ ] Unit test untuk StoreToken passing
- [ ] Integration test OAuth flow passing

### 2. Health Check Endpoint

Tambahkan health check untuk Google Calendar service:

```go
// Health check endpoint
func (h *GoogleCalendarAuthHandler) HealthCheck(c *gin.Context) {
    // Check database connection
    _, err := h.tokenService.GetToken("health-check-user")
    if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        c.JSON(503, gin.H{"status": "unhealthy", "error": "database connection failed"})
        return
    }

    // Check encryption
    testToken := &oauth2.Token{AccessToken: "test", RefreshToken: "test"}
    err = h.tokenService.StoreToken("health-check-user", testToken)
    if err != nil {
        c.JSON(503, gin.H{"status": "unhealthy", "error": "encryption failed"})
        return
    }

    c.JSON(200, gin.H{"status": "healthy"})
}
```

### 3. Monitoring

Tambahkan alert untuk error rate:

```yaml
# Prometheus alert
- alert: GoogleCalendarStoreTokenError
  expr: rate(google_calendar_store_token_errors_total[5m]) > 0
  for: 1m
  labels:
    severity: critical
  annotations:
    summary: "Google Calendar token storage failing"
    description: "failed_to_store_token error detected"
```

---

## Related Files

- `apps/api/internal/api/handlers/google_calendar_auth_handler.go`
- `apps/api/internal/service/google_calendar_token/service.go`
- `apps/api/internal/repository/postgres/google_calendar_token/repository.go`
- `apps/mobile/lib/main.dart`
- `apps/mobile/lib/features/google_calendar/application/google_calendar_provider.dart`

---

## References

- [Google Calendar OAuth Documentation](../google-calendar/README.md)
- [Setup Guide for Dev1](../guides/google-calendar-oauth-setup.md)
- [Architecture Options](google-calendar-oauth-options.md)

---

## Changelog

- **2026-03-07** - Initial analysis and documentation
- **2026-03-07** - Added diagnostic steps and solutions
- **2026-03-07** - Added prevention measures

---

## Author

- **Analysis:** AI Assistant
- **Date:** 2026-03-07
- **Status:** Active Issue - Pending Resolution

---

## Next Steps

1. ⏳ Check backend logs for specific error message
2. ⏳ Verify database table exists
3. ⏳ Verify encryption key is set correctly
4. ⏳ Run database migration if needed
5. ⏳ Re-test OAuth flow

**Assignee:** Backend Team / DevOps  
**Priority:** HIGH  
**Estimated Fix Time:** 30 minutes - 2 hours (depending on root cause)
