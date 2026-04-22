# Guide - Project Setup

## CRM Healthcare Mobile App - Flutter

**Module**: Development Guide  
**Sprint**: Sprint 0  
**Version**: 1.0  
**Status**: ✅ **Completed**  
**Last Updated**: January 2025

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Flutter Installation](#flutter-installation)
3. [Project Setup](#project-setup)
4. [Environment Configuration](#environment-configuration)
5. [Running the App](#running-the-app)
6. [Troubleshooting](#troubleshooting)

---

## Prerequisites

### System Requirements

- **OS**: Windows 10/11, macOS 10.14+, or Linux
- **RAM**: 8GB minimum, 16GB recommended
- **Storage**: 10GB free space
- **IDE**: Android Studio, VS Code, atau IntelliJ IDEA

### Required Tools

1. **Flutter SDK** (3.16.0 atau lebih baru)
2. **Dart SDK** (bundled dengan Flutter)
3. **Android Studio** (untuk Android development)
4. **Xcode** (untuk iOS development, macOS only)
5. **Git**

---

## Flutter Installation

### 1. Download Flutter

```bash
# Windows (gunakan PowerShell)
flutter --version  # Check if already installed

# Jika belum, download dari:
# https://docs.flutter.dev/get-started/install
```

### 2. Update PATH

**Windows**:

```bash
# Tambahkan ke System Environment Variables
C:\flutter\bin
```

**macOS/Linux**:

```bash
# Tambahkan ke ~/.bashrc atau ~/.zshrc
export PATH="$PATH:`pwd`/flutter/bin"
```

### 3. Verify Installation

```bash
flutter doctor
```

Expected output:

```
[✓] Flutter (Channel stable, 3.16.0, ...)
[✓] Android toolchain
[✓] Xcode (untuk macOS)
[✓] Chrome
[✓] Android Studio
```

---

## Project Setup

### 1. Clone Repository

```bash
# Clone monorepo
git clone https://github.com/your-org/crm-healthcare.git
cd crm-healthcare/apps/mobile
```

### 2. Install Dependencies

```bash
# Install Flutter dependencies
flutter pub get

# Generate code (Hive adapters, Freezed, dll)
flutter pub run build_runner build --delete-conflicting-outputs
```

### 3. Setup Environment

```bash
# Copy environment template
copy .env.example .env  # Windows
cp .env.example .env    # macOS/Linux

# Edit .env dengan values yang sesuai
API_BASE_URL=http://localhost:8080
ENVIRONMENT=development
```

---

## Environment Configuration

### Android Setup

1. **Install Android SDK** via Android Studio
2. **Create Emulator** atau connect physical device
3. **Enable USB Debugging** di device (untuk physical device)

### iOS Setup (macOS only)

1. **Install Xcode** dari App Store
2. **Install CocoaPods**:
   ```bash
   sudo gem install cocoapods
   ```
3. **Setup Simulator** atau connect iPhone

---

## Running the App

### Development Mode

```bash
# Android
flutter run

# iOS
flutter run

# Dengan specific device
flutter devices  # List devices
flutter run -d emulator-5554
```

### Hot Reload

Tekan `r` di terminal untuk hot reload  
Tekan `R` untuk hot restart  
Tekan `q` untuk quit

### Build Modes

```bash
# Debug (development)
flutter run --debug

# Profile (performance testing)
flutter run --profile

# Release (production)
flutter run --release
```

---

## Troubleshooting

### Common Issues

**1. Gradle Sync Failed**:

```bash
# Clean dan rebuild
cd android
./gradlew clean
flutter clean
flutter pub get
```

**2. iOS Build Error**:

```bash
cd ios
pod install --repo-update
```

**3. Dependency Conflicts**:

```bash
flutter pub upgrade
flutter pub run build_runner build --delete-conflicting-outputs
```

**4. Emulator Not Found**:

- Pastikan emulator sudah dibuat di Android Studio
- Atau enable USB debugging untuk physical device

---

**Document Status**: Active  
**Last Updated**: January 2025
