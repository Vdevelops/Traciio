# CRM Healthcare Mobile App (Flutter)

Aplikasi mobile untuk **CRM Healthcare/Pharmaceutical Platform** yang digunakan oleh **Sales Representative**.  
Project ini adalah bagian dari monorepo `crm-healthcare` dan mengikuti perencanaan sprint di `docs/SPRINT_PLANNING_DEV3.md`.

## 📱 Tujuan & Scope

- Menyediakan aplikasi mobile untuk sales rep dengan fitur utama:
  - Autentikasi (login/logout)
  - Account & Contact (list & detail)
  - Visit Report (dengan GPS & foto) — sprint berikutnya
  - Task & Reminder — sprint berikutnya
  - Dashboard ringkas — sprint berikutnya
- Aplikasi ini menjadi **client** dari backend API (`apps/api`) dan menjaga konsistensi UI/UX dengan web app (`apps/web`).

## 🧱 Teknologi

- **Flutter** (stable) + **Dart ≥ 3**
- **State Management**: `flutter_riverpod`
- **HTTP Client**: `dio`
- **Local Storage**: `shared_preferences` (untuk token auth dan konfigurasi ringan)
- **Arsitektur**: feature-based dengan pemisahan `core/` dan `features/` sesuai `mobile-dev3` rules

## 📋 System Requirements

**Quick Reference:**
- **Minimum**: Android 7.0+ (API 24) / iOS 12.0+, 2 GB RAM, 500 MB storage
- **Recommended**: Android 10.0+ / iOS 15.0+, 4 GB RAM, 1 GB storage
- **Rekomendasi HP**: Rp 3-6 juta (mid-range) untuk penggunaan optimal jangka panjang

**📖 Untuk informasi lengkap, lihat [docs/SYSTEM_REQUIREMENTS.md](docs/SYSTEM_REQUIREMENTS.md)**

## 🗂 Struktur Project (Mobile)

Lokasi project mobile di dalam monorepo:

```text
apps/
  mobile/
    lib/
      main.dart
      core/
        config/        # Env & konfigurasi global
        routing/       # AppRouter & AppRoutes
        theme/         # AppTheme (light/dark)
        network/       # ApiClient (Dio + interceptors)
        storage/       # LocalStorage (shared_preferences)
        widgets/       # Widget shared (mis. AuthGate)
      features/
        auth/
          data/        # AuthRepository, models (nanti)
          application/ # AuthState, AuthNotifier, providers
          presentation/# LoginScreen & UI terkait auth
        accounts/
          data/
          application/
          presentation/
        contacts/
        visit_reports/
        tasks/
        dashboard/
```

Struktur ini mengikuti **Sprint 0** dan akan diisi bertahap sesuai sprint berikutnya (Account & Contact, Visit Report, Task & Reminder, Dashboard).

## 🚀 Quick Start

### Installation

```bash
cd apps/mobile
flutter pub get
```

### Running the App

```bash
# Development (Android Emulator)
flutter run --dart-define=API_BASE_URL=http://10.0.2.2:8080

# Development (iOS Simulator)
flutter run --dart-define=API_BASE_URL=http://localhost:8080

# Development (Physical Device)
flutter run --dart-define=API_BASE_URL=http://<YOUR_PC_IP>:8080
```

### Building for Production

```bash
# Android
flutter build apk --release --dart-define=API_BASE_URL=https://api.gilabs.id

# iOS
flutter build ios --release --dart-define=API_BASE_URL=https://api.gilabs.id
```

**📖 Untuk panduan setup lengkap, lihat [docs/SETUP.md](docs/SETUP.md)**

## 📚 Documentation

- **[SETUP.md](docs/SETUP.md)** - Setup guide dan environment configuration
- **[SYSTEM_REQUIREMENTS.md](docs/SYSTEM_REQUIREMENTS.md)** - System requirements dan rekomendasi spesifikasi HP
- **[GPS_SETUP.md](docs/GPS_SETUP.md)** - GPS setup dan testing guide
- **[OPTIMIZATION.md](docs/OPTIMIZATION.md)** - Optimization summary dan results
- **[UI_UX.md](docs/UI_UX.md)** - UI/UX simplification documentation

## 🧪 Lint & Testing

Jalankan sebelum push/perubahan besar:

```bash
cd apps/mobile
flutter analyze
flutter test
```

Pastikan tidak ada error lint besar, terutama saat menyelesaikan acceptance criteria sprint.

## 🤝 Koordinasi dengan Tim Lain

Sesuai `docs/SPRINT_PLANNING.md`:

- **Developer 1 (Web)**: koordinasi untuk konsistensi UI/UX (warna, layout, copy teks).
- **Developer 2 (Backend)**: koordinasi untuk desain API (endpoint, payload, error codes).
- Gunakan Postman collection di `docs/postman/CRM-Healthcare-API.postman_collection.json` sebagai referensi kontrak API saat mulai integrasi.

---

Untuk detail lengkap mengenai scope dan timeline Dev 3, lihat:

- `docs/SPRINT_PLANNING_DEV3.md`
- `docs/SPRINT_PLANNING.md`
