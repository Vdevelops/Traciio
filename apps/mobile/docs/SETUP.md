# Setup Guide - Mobile App

## Prerequisites

- Flutter SDK (channel **stable**)
- Android Studio / VS Code dengan Flutter plugin
- Device/emulator Android atau iOS
- Backend API (Go + Gin) dari repo ini

## Installation

### 1. Install Dependencies

```bash
cd apps/mobile
flutter pub get
```

### 2. Environment Configuration

Mobile app menggunakan environment variables yang di-set melalui `--dart-define` saat build atau run.

#### Available Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `API_BASE_URL` | Base URL untuk backend API | Platform-specific | Yes |

#### Platform-Specific Defaults

Jika `API_BASE_URL` tidak di-set, aplikasi akan menggunakan default berdasarkan platform:

- **Android Emulator**: `http://10.0.2.2:8080` (alias untuk host machine)
- **iOS Simulator**: `http://localhost:8080`
- **Physical Device**: Harus set manual dengan IP address PC

## Running the App

### Development (Android Emulator)

```bash
flutter run --dart-define=API_BASE_URL=http://10.0.2.2:8080
```

**Catatan:**
- `10.0.2.2` adalah alias khusus Android emulator untuk mengakses host machine
- Jangan gunakan `localhost` atau `127.0.0.1` di Android emulator

### Development (iOS Simulator)

```bash
flutter run --dart-define=API_BASE_URL=http://localhost:8080
```

### Development (Physical Device)

```bash
# Ganti 192.168.1.100 dengan IP address PC/server Anda
flutter run --dart-define=API_BASE_URL=http://192.168.1.100:8080
```

**Catatan:**
- Pastikan PC dan device dalam network WiFi yang sama
- Pastikan firewall tidak memblokir port 8080
- Test koneksi dengan browser di device: `http://<PC_IP>:8080/health`

### Using Scripts (Recommended)

Untuk memudahkan, gunakan script di `package.json`:

```bash
# Android Emulator
npm run dev:android

# iOS Simulator
npm run dev:ios

# Physical Device
npm run dev:device
```

## Building for Production

**PENTING**: Saat build untuk production, WAJIB menggunakan URL API production agar aplikasi yang di-install menggunakan API production.

### Android Production Build

```bash
flutter build apk --release --dart-define=API_BASE_URL=https://api.gilabs.id
```

### iOS Production Build

```bash
flutter build ios --release --dart-define=API_BASE_URL=https://api.gilabs.id
```

**Catatan Penting:**
- URL `https://api.gilabs.id` adalah URL API production yang sudah dikonfigurasi
- Pastikan menggunakan `--release` flag untuk build production
- URL ini akan ter-embed di aplikasi saat build, jadi aplikasi yang di-install akan otomatis menggunakan API production
- Jangan lupa set `API_BASE_URL` saat build, jika tidak aplikasi akan menggunakan default (localhost) yang tidak akan bekerja di production

### Verifying Build

Setelah build, pastikan aplikasi menggunakan API production dengan:

1. Install aplikasi di device
2. Buka aplikasi dan coba login
3. Pastikan aplikasi bisa connect ke API production (`https://api.gilabs.id`)

## WebSocket Configuration

WebSocket URL otomatis di-generate dari `API_BASE_URL`:

- `http://` → `ws://`
- `https://` → `wss://`

Contoh:
- `API_BASE_URL=http://10.0.2.2:8080`
  - WebSocket: `ws://10.0.2.2:8080/api/v1/ws/notifications?token=...`

## Backend API Setup

**Pastikan backend API sudah running:**

- Jalankan API dari `apps/api` sesuai `apps/api/SETUP.md`
- Untuk development, bisa pakai: `pnpm run dev:web-api-docker` (dari root repo)
- Pastikan endpoint dan skema respons sesuai standar di `docs/api-standart/`

## Troubleshooting

### Android Emulator Connection Issues

Jika Android emulator tidak bisa connect ke backend:

1. Pastikan backend API running di host machine
2. Gunakan `10.0.2.2` bukan `localhost` atau `127.0.0.1`
3. Test dengan browser di emulator: `http://10.0.2.2:8080/health`

### Physical Device Connection Issues

Jika physical device tidak bisa connect:

1. Pastikan PC dan device dalam WiFi network yang sama
2. Pastikan firewall tidak memblokir port 8080
3. Gunakan IP address PC (bukan `localhost`)
4. Test dengan browser di device: `http://<PC_IP>:8080/health`

### Build Issues

Jika build gagal:

1. Pastikan Flutter SDK up-to-date: `flutter upgrade`
2. Clean build: `flutter clean && flutter pub get`
3. Check dependencies: `flutter pub outdated`
4. Check for errors: `flutter analyze`

---

Untuk informasi lebih detail, lihat [README.md](../README.md)

