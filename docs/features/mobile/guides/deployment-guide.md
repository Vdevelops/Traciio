# Guide - Deployment

## CRM Healthcare Mobile App - Flutter

**Module**: Development Guide  
**Sprint**: Sprint 6  
**Version**: 1.0  
**Status**: ✅ **Completed**  
**Last Updated**: January 2025

---

## Table of Contents

1. [Pre-Deployment Checklist](#pre-deployment-checklist)
2. [Android Deployment](#android-deployment)
3. [iOS Deployment](#ios-deployment)
4. [Environment Setup](#environment-setup)
5. [Build Process](#build-process)
6. [Distribution](#distribution)

---

## Pre-Deployment Checklist

### Code Quality

- [ ] All tests passing
- [ ] No compiler warnings
- [ ] Static analysis clean
- [ ] Code review completed
- [ ] Documentation updated

### Configuration

- [ ] Environment variables configured
- [ ] API base URL pointing to production
- [ ] Debug logging disabled
- [ ] Analytics enabled

### Assets

- [ ] App icon (all sizes)
- [ ] Splash screen
- [ ] Feature graphics
- [ ] Screenshots prepared

---

## Android Deployment

### 1. Setup Signing

**Create keystore**:

```bash
keytool -genkey -v -keystore crm-healthcare.keystore \
  -alias crm_healthcare \
  -keyalg RSA -keysize 2048 -validity 10000
```

**Configure signing** (`android/key.properties`):

```
storePassword=your_password
keyPassword=your_password
keyAlias=crm_healthcare
storeFile=crm-healthcare.keystore
```

### 2. Build APK

```bash
# Release APK
flutter build apk --release

# App Bundle (recommended for Play Store)
flutter build appbundle --release
```

Output:

- `build/app/outputs/flutter-apk/app-release.apk`
- `build/app/outputs/bundle/release/app-release.aab`

### 3. Play Store Upload

1. Login ke Google Play Console
2. Create new release
3. Upload app bundle
4. Fill store listing
5. Submit untuk review

---

## iOS Deployment

### 1. Setup Certificates

- Apple Developer Account
- Distribution certificate
- App ID provisioning profile

### 2. Configure Xcode

```bash
cd ios
open Runner.xcworkspace
```

Configure di Xcode:

- Bundle identifier
- Signing & Capabilities
- App Icons
- Launch Screen

### 3. Build IPA

```bash
# Archive
flutter build ipa --release

# Atau via Xcode
# Product -> Archive
```

### 4. App Store Connect

1. Upload via Xcode Organizer
2. Atau menggunakan Transporter app
3. Fill App Store information
4. Submit untuk review

---

## Environment Setup

### Production Environment

```dart
// lib/core/config/env.dart
class Env {
  static const String apiBaseUrl = String.fromEnvironment(
    'API_BASE_URL',
    defaultValue: 'https://api.crm-healthcare.com',
  );

  static const String environment = String.fromEnvironment(
    'ENVIRONMENT',
    defaultValue: 'production',
  );
}
```

### Build dengan Environment

```bash
# Android
flutter build apk --release \
  --dart-define=API_BASE_URL=https://api.crm-healthcare.com \
  --dart-define=ENVIRONMENT=production

# iOS
flutter build ipa --release \
  --dart-define=API_BASE_URL=https://api.crm-healthcare.com \
  --dart-define=ENVIRONMENT=production
```

---

## Build Process

### Automated Build (CI/CD)

```yaml
# .github/workflows/build.yml
name: Build Release

on:
  push:
    tags:
      - "v*"

jobs:
  build-android:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: subosito/flutter-action@v2
        with:
          flutter-version: "3.16.0"
      - run: flutter pub get
      - run: flutter build apk --release
      - name: Upload APK
        uses: actions/upload-artifact@v3
        with:
          name: release-apk
          path: build/app/outputs/flutter-apk/app-release.apk
```

---

## Distribution

### Internal Testing

**Android**:

- Firebase App Distribution
- Internal Testing Track (Play Store)

**iOS**:

- TestFlight
- Ad Hoc distribution

### Production Release

**Staged Rollout**:

1. Release to 10% of users
2. Monitor for issues
3. Increase to 50%
4. Full rollout

**Monitoring**:

- Crashlytics
- Analytics
- Performance metrics
- User feedback

---

## Version Management

### Semantic Versioning

Format: `MAJOR.MINOR.PATCH`

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes

### Version Update

```yaml
# pubspec.yaml
version: 1.2.3+4 # version_name+version_code
```

```bash
# Android: versionCode di android/app/build.gradle
# iOS: CFBundleVersion di ios/Runner/Info.plist
```

---

**Document Status**: Active  
**Last Updated**: January 2025
