# Google Calendar OAuth Setup Guide for Dev1

## Overview

Guide ini menjelaskan langkah-langkah yang harus dilakukan oleh **Dev1** untuk mengkonfigurasi Google Calendar OAuth di Google Cloud Console.

**Priority:** HIGH  
**Blocking:** Mobile Google Calendar OAuth tidak akan berfungsi tanpa konfigurasi ini  
**Estimated Time:** 5-10 menit

---

## Prerequisites

- [ ] Akses ke Google Cloud Console dengan permission Owner/Editor
- [ ] Project ID: `1051532630602-lcicb5e4bpldcbjoslmb98abj8o6g5gs`
- [ ] OAuth 2.0 Client ID untuk Web application sudah dibuat

---

## Setup Steps

### Step 1: Login ke Google Cloud Console

1. Buka [Google Cloud Console](https://console.cloud.google.com/)
2. Login dengan akun yang memiliki akses ke project
3. Pastikan project yang aktif adalah:
   ```
   1051532630602-lcicb5e4bpldcbjoslmb98abj8o6g5gs
   ```

### Step 2: Navigate ke Credentials

1. Klik menu **☰ (Hamburger menu)** di kiri atas
2. Pilih **APIs & Services** > **Credentials**
3. Atau langsung URL:
   ```
   https://console.cloud.google.com/apis/credentials
   ```

### Step 3: Edit OAuth 2.0 Client ID

1. Cari section **OAuth 2.0 Client IDs**
2. Cari client dengan nama **"Web client"** atau **"Web application"**
3. Klik **Edit** (icon pencil) pada client tersebut

**Expected Client ID:**

```
1051532630602-lcicb5e4bpldcbjoslmb98abj8o6g5gs.apps.googleusercontent.com
```

### Step 4: Add Authorized Redirect URI (HTTPS)

**Catatan Penting**: Menggunakan metode **HTTPS Web Redirect + Server Forward** (Recommended by Google)

1. Scroll ke section **Authorized redirect URIs**
2. Klik tombol **+ ADD URI**
3. Masukkan URL berikut:
   ```
   https://api.gilabs.id/api/v1/google-calendar/callback
   ```
4. Klik **SAVE**

**Catatan**: URL ini SAMA untuk web dan mobile. Backend akan mendeteksi platform dari state parameter dan melakukan forward ke mobile app jika diperlukan.

### Step 5: Verification

1. Pastikan URL sudah muncul di daftar **Authorized redirect URIs**
2. Contoh tampilan yang benar:
   ```
   Authorized redirect URIs
   ☑ https://api.gilabs.id/api/v1/google-calendar/callback    ← UNTUK WEB & MOBILE
   ```

**Tidak perlu menambahkan custom scheme** (`crmhealth://`) karena menggunakan metode server forward.

---

## Screenshots Reference

### Step 2: Credentials Page

```
┌─────────────────────────────────────────────────────────────┐
│ APIs & Services > Credentials                                │
├─────────────────────────────────────────────────────────────┤
│ OAuth 2.0 Client IDs                                        │
│ ┌───────────────────────────────────────────────────────────┐│
│ │ Web client 1                                              ││
│ │ Client ID: 1051532630602-...                              ││
│ │ [Edit] [Delete]                                           ││
│ └───────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────┘
```

### Step 4: Add Redirect URI

```
┌─────────────────────────────────────────────────────────────┐
│ Authorized redirect URIs                                    │
├─────────────────────────────────────────────────────────────┤
│ ☑ https://api.gilabs.id/api/v1/google-calendar/callback     │
│ ☑ crmhealth://google-calendar/callback    ← TAMBAHKAN INI  │
│ [+ ADD URI]                                                 │
└─────────────────────────────────────────────────────────────┘
```

---

## Testing Verification

Setelah konfigurasi selesai, tim developer akan melakukan testing. Dev1 tidak perlu melakukan testing teknis, tapi bisa verify dengan:

1. Check bahwa URI sudah tersimpan di console
2. Confirm dengan tim mobile/backend bahwa konfigurasi sudah OK

---

## Rollback (jika diperlukan)

Jika perlu menghapus konfigurasi:

1. Kembali ke **Credentials > OAuth 2.0 Client IDs**
2. Edit Web client
3. Cari `crmhealth://google-calendar/callback` di daftar
4. Klik **X** (delete) pada URI tersebut
5. Click **SAVE**

---

## FAQ

### Q: Menggunakan HTTPS redirect bukan custom scheme?

**A:** Ya, mengikuti rekomendasi Google untuk production:

- **Keamanan**: HTTPS lebih aman dan terpercaya
- **Validasi**: Google hanya menerima domain HTTPS yang terverifikasi
- **Server Forward**: Backend menerima callback, exchange code, lalu forward ke mobile app
- **No Custom Scheme Registration**: Tidak perlu register custom scheme di Google Cloud Console

### Q: Bagaimana alur kerjanya?

**A:**

1. Mobile request auth URL dengan `platform=mobile`
2. Backend generate OAuth URL dengan redirect ke `https://api.gilabs.id/api/v1/google-calendar/callback`
3. User authorize di Google
4. Google redirect ke HTTPS endpoint
5. Backend detect mobile dari state, exchange code → token → store
6. Backend redirect ke mobile app via deep link `crmhealth://google-calendar/callback?success=true`
7. Mobile app menerima notifikasi success

### Q: Apakah perlu publish app terlebih dahulu?

**A:** Tidak, development/testing bisa dilakukan tanpa publish.

### Q: Berapa lama perubahan berlaku?

**A:** Biasanya langsung berlaku (real-time), tapi kadang butuh beberapa menit untuk propagate.

### Q: Apakah perlu restart service?

**A:** Tidak, perubahan di Google Cloud Console independent dari backend service.

---

## Contact

Jika ada kendala:

- **Mobile Team:** [Contact Mobile Lead]
- **Backend Team:** [Contact Backend Lead]
- **Documentation:** `docs/features/mobile/google-calendar/README.md`

---

## Checklist

- [ ] Login ke Google Cloud Console
- [ ] Navigate ke APIs & Services > Credentials
- [ ] Edit OAuth 2.0 Client ID (Web client)
- [ ] Add redirect URI: `https://api.gilabs.id/api/v1/google-calendar/callback`
- [ ] Save configuration
- [ ] Verify URI muncul di daftar
- [ ] Notify team bahwa setup selesai

## Environment Variables (Backend)

Pastikan backend environment variable sudah di-set:

```env
GOOGLE_CALENDAR_REDIRECT_URL=https://api.gilabs.id/api/v1/google-calendar/callback
```

**Note**: Mobile menggunakan URL yang sama dengan web (HTTPS). Backend akan mendeteksi platform dari state parameter.

---

**Last Updated:** 2026-03-05  
**Document Version:** 1.0  
**Author:** AI Assistant
