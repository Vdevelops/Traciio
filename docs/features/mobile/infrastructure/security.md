# Infrastructure - Security

## CRM Healthcare Mobile App - Flutter

**Module**: Infrastructure  
**Sprint**: Sprint 0  
**Version**: 1.0  
**Status**: ✅ **Completed**  
**Last Updated**: January 2025

---

## Table of Contents

1. [Ringkasan Fitur](#ringkasan-fitur)
2. [Security Areas](#security-areas)
3. [Business Rules](#business-rules)
4. [Keputusan Teknis & Trade-offs](#keputusan-teknis--trade-offs)
5. [Implementation](#implementation)
6. [Best Practices](#best-practices)
7. [Dependencies](#dependencies)
8. [Notes & Improvements](#notes--improvements)

---

## Ringkasan Fitur

Dokumen ini menjelaskan security best practices dan implementasi untuk CRM Healthcare mobile app, mencakup data protection, secure communication, dan authentication security.

### Goals

- **Data Protection**: Encrypt sensitive data
- **Secure Communication**: HTTPS dengan certificate pinning
- **Authentication**: Secure token storage dan handling
- **Input Validation**: Prevent injection attacks
- **Code Security**: Obfuscation dan anti-tampering

---

## Security Areas

### 1. Data at Rest

**Local Storage**:

- Hive dengan encryption
- Secure storage untuk tokens
- No hardcoded secrets

**Implementation**:

```dart
// Hive encryption
final key = await _generateKey();
final encryptedBox = await Hive.openBox(
  'secure_box',
  encryptionCipher: HiveAesCipher(key),
);
```

### 2. Data in Transit

**HTTPS Only**:

- All API calls via HTTPS
- Certificate pinning untuk production
- No HTTP fallback

### 3. Authentication Security

**Token Storage**:

- Flutter Secure Storage
- Automatic token refresh
- Secure logout (clear all)

**Biometric Auth** (optional):

- Fingerprint/Face ID
- Secure enrollment

### 4. Input Validation

**Client-side**:

- Form validation
- Sanitize user inputs
- Prevent XSS

**Server-side**:

- Never trust client input
- Always validate di backend

---

## Business Rules

### 1. Token Handling

- Store di secure storage (production)
- Refresh before expiry
- Clear on logout
- Auto-logout saat token invalid

### 2. Password Policy

- Min 8 characters
- Upper, lower, number
- No common passwords
- Change periodically

### 3. Session Management

- Timeout after 30 minutes inactivity
- Single session per device (optional)
- Force logout on password change

---

## Keputusan Teknis & Trade-offs

### Certificate Pinning

**Keputusan**: Implement untuk production.

**Alasan**:

- Prevent MITM attacks
- Ensure server authenticity

**Trade-off**: Maintenance overhead saat certificate renew. **Mitigasi**: Proper certificate management process.

---

## Implementation

### Secure Storage Helper

```dart
class TokenSecureStorage {
  static const _storage = FlutterSecureStorage();

  static Future<void> saveTokens(String access, String refresh) async {
    await _storage.write(key: 'access_token', value: access);
    await _storage.write(key: 'refresh_token', value: refresh);
  }

  static Future<Map<String, String?>> getTokens() async {
    return {
      'access': await _storage.read(key: 'access_token'),
      'refresh': await _storage.read(key: 'refresh_token'),
    };
  }

  static Future<void> clearTokens() async {
    await _storage.deleteAll();
  }
}
```

### Input Sanitization

```dart
class InputValidator {
  static String? sanitize(String? input) {
    if (input == null) return null;
    return input.trim().replaceAll(RegExp(r'[<>]'), '');
  }

  static bool isValidEmail(String email) {
    return RegExp(r'^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$').hasMatch(email);
  }
}
```

---

## Best Practices

1. **Never commit secrets**: Use environment variables
2. **Use latest dependencies**: Security updates
3. **Obfuscate code**: Enable Flutter obfuscation
4. **Root detection**: Optional, untuk production
5. **Secure logging**: No sensitive data di logs

---

## Dependencies

```yaml
dependencies:
  flutter_secure_storage: ^9.0.0
  crypto: ^3.0.3
```

---

**Document Status**: Active  
**Last Updated**: January 2025
