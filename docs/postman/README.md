# Postman Collection - CRM Healthcare API

Dokumentasi Postman collection untuk CRM Healthcare API.

## Setup

1. **Import Collection**
   - Buka Postman
   - Klik Import
   - Pilih file `CRM-Healthcare-API.postman_collection.json`

2. **Setup Environment Variables**
   - Buat environment baru dengan nama "CRM Healthcare - Local"
   - Set variable:
     - `base_url`: `http://localhost:8080`
   - Set variable:
     - `token`: (akan di-set otomatis setelah login)
     - `refresh_token`: (akan di-set otomatis setelah login)
     - `user_id`: (akan di-set otomatis setelah login)

## Usage

### 1. Health Check
- Jalankan request "Health" atau "Ping" untuk memastikan API berjalan

### 2. Login
- Jalankan request "Login" dengan credentials:
  ```json
  {
    "email": "admin@example.com",
    "password": "admin123"
  }
  ```
- Token akan otomatis disimpan di environment variable setelah login berhasil

### 3. Authenticated Requests
- Semua request yang memerlukan authentication akan otomatis menggunakan token dari environment variable
- Token akan di-set di header `Authorization: Bearer {{token}}`

## Collection Structure

### Authentication
- **Login**: POST `/api/v1/auth/login`
- **Refresh Token**: POST `/api/v1/auth/refresh`
- **Logout**: POST `/api/v1/auth/logout`

### Health Check
- **Health**: GET `/health`
- **Ping**: GET `/ping`

## Response Format

Semua response mengikuti format standar:

### Success Response
```json
{
  "success": true,
  "data": { ... },
  "timestamp": "2024-01-15T10:30:45+07:00",
  "request_id": "req_abc123xyz"
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error message"
  },
  "timestamp": "2024-01-15T10:30:45+07:00",
  "request_id": "req_abc123xyz"
}
```

## Notes

- Collection ini akan terus diupdate sesuai dengan perkembangan API
- Untuk detail lengkap, lihat dokumentasi di `/docs/api-standart/`

