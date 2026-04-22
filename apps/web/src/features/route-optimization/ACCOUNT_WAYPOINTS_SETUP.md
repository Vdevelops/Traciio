# Account Waypoints Setup

## Overview

Fitur route optimization sekarang mendukung penggunaan account sebagai waypoints dengan geocoding otomatis.

## Setup

### Backend
Route geocoding sudah terdaftar di:
- Path: `/api/v1/geocoding/geocode`
- Method: POST
- Auth: Required (JWT token)

### Frontend
Service geocoding tersedia di:
- File: `apps/web/src/features/route-optimization/services/geocodingService.ts`
- Usage: `geocodingService.geocode(address)`

## Troubleshooting

### Error 404: Request failed with status code 404

**Kemungkinan penyebab:**
1. **Server belum di-restart** setelah route ditambahkan
   - **Solusi**: Restart backend server
   ```bash
   cd apps/api
   go run cmd/server/main.go
   ```

2. **Route belum terdaftar**
   - **Cek**: Pastikan `routes.SetupGeocodingRoutes(v1, geocodingHandler, jwtManager)` ada di `main.go`
   - **Lokasi**: `apps/api/cmd/server/main.go` line ~541

3. **Path endpoint salah**
   - **Expected**: `/api/v1/geocoding/geocode`
   - **Cek**: Pastikan baseURL di `api-client.ts` sudah benar

### Error: Failed to geocode address

**Kemungkinan penyebab:**
1. **Address tidak valid atau tidak ditemukan**
   - **Solusi**: Pastikan account memiliki address yang valid
   - System akan menampilkan error toast untuk account yang gagal

2. **Nominatim API rate limit**
   - **Solusi**: Tunggu beberapa saat dan coba lagi
   - Backend sudah memiliki retry logic

3. **Network error**
   - **Solusi**: Cek koneksi internet
   - System akan menampilkan error message yang jelas

## Testing

1. **Test endpoint langsung:**
   ```bash
   curl -X POST http://localhost:8080/api/v1/geocoding/geocode \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -d '{"address": "Jakarta, Indonesia"}'
   ```

2. **Test dari UI:**
   - Buka route optimization form
   - Klik "Select Waypoints"
   - Pilih tab "Accounts"
   - Pilih account dengan address
   - Klik "Add Waypoints"
   - System akan geocode addresses secara otomatis

## Features

✅ **Account Selection**: User bisa memilih accounts sebagai waypoints
✅ **Automatic Geocoding**: Address di-geocode otomatis ke coordinates
✅ **Error Handling**: Error ditampilkan dengan toast notifications
✅ **Loading States**: Loading indicator saat geocoding
✅ **Validation**: Account tanpa address tidak bisa dipilih

## API Response Format

**Request:**
```json
{
  "address": "Jl. Kebayoran Baru No. 78, Jakarta"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "latitude": -6.2088,
    "longitude": 106.8456,
    "address": "Jl. Kebayoran Baru No. 78, Jakarta"
  },
  "timestamp": "2024-01-01T00:00:00+07:00",
  "request_id": "req_..."
}
```

