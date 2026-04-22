# Postman Setup Guide

## Import Collection ke Postman

### Method 1: Import File
1. Buka Postman
2. Klik **Import** di sidebar kiri
3. Pilih file `CRM-Healthcare-API.postman_collection.json`
4. Klik **Import**

### Method 2: Import via URL (jika collection di-host)
1. Buka Postman
2. Klik **Import**
3. Pilih tab **Link**
4. Masukkan URL collection
5. Klik **Continue** → **Import**

## Setup Environment

1. Klik **Environments** di sidebar kiri
2. Klik **+** untuk membuat environment baru
3. Beri nama: `CRM Healthcare - Local`
4. Tambahkan variables:

| Variable | Initial Value | Current Value |
|----------|---------------|---------------|
| `base_url` | `http://localhost:8080` | `http://localhost:8080` |
| `token` | (kosong) | (akan di-set otomatis) |
| `refresh_token` | (kosong) | (akan di-set otomatis) |
| `user_id` | (kosong) | (akan di-set otomatis) |

5. Klik **Save**
6. Pilih environment tersebut di dropdown environment (kanan atas)

## Testing Flow

### 1. Health Check
- Jalankan request **Health Check → Health**
- Pastikan response: `{"status": "ok", "message": "API is running"}`

### 2. Login
- Jalankan request **Authentication → Login**
- Body request:
  ```json
  {
    "email": "admin@example.com",
    "password": "admin123"
  }
  ```
- Setelah berhasil, token akan otomatis tersimpan di environment variable
- Response akan berisi:
  - `token`: JWT access token
  - `refresh_token`: Refresh token untuk renew access token
  - `user`: User information

### 3. Authenticated Requests
- Setelah login, semua request yang memerlukan authentication akan otomatis menggunakan token
- Token di-set di header: `Authorization: Bearer {{token}}`

### 4. Refresh Token
- Jika token expired, gunakan request **Authentication → Refresh Token**
- Refresh token akan otomatis di-update di environment

### 5. Logout
- Gunakan request **Authentication → Logout** untuk logout
- Token akan di-invalidate

## Tips

1. **Auto-save Token**: Collection sudah dikonfigurasi untuk auto-save token setelah login berhasil
2. **Environment Variables**: Gunakan `{{variable_name}}` untuk menggunakan environment variables
3. **Request ID**: Setiap response memiliki `request_id` untuk tracking di logs

## Troubleshooting

### Token tidak tersimpan
- Pastikan environment sudah dipilih
- Cek Test script di request Login (harus ada script untuk set token)

### 401 Unauthorized
- Pastikan token masih valid
- Coba refresh token atau login ulang

### Connection Error
- Pastikan API server berjalan di `http://localhost:8080`
- Cek firewall/network settings

