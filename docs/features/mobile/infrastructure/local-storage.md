# Infrastructure - Local Storage

## CRM Healthcare Mobile App - Flutter

**Module**: Infrastructure  
**Sprint**: Sprint 0  
**Version**: 1.0  
**Status**: ✅ **Completed**  
**Last Updated**: January 2025

---

## Table of Contents

1. [Ringkasan Fitur](#ringkasan-fitur)
2. [Storage Options](#storage-options)
3. [Business Rules](#business-rules)
4. [Keputusan Teknis & Trade-offs](#keputusan-teknis--trade-offs)
5. [Struktur Folder](#struktur-folder)
6. [Storage Implementation](#storage-implementation)
7. [Configuration](#configuration)
8. [Usage Examples](#usage-examples)
9. [Best Practices](#best-practices)
10. [Cara Test Manual](#cara-test-manual)
11. [Dependencies](#dependencies)
12. [Notes & Improvements](#notes--improvements)

---

## Ringkasan Fitur

Sistem **Local Storage** mobile app CRM Healthcare menggunakan kombinasi **Hive** untuk structured data dan **SharedPreferences** untuk simple key-value storage. Sistem ini menyediakan type-safe, performant, dan reliable local persistence untuk offline support dan data caching.

### Goals

- **Structured Data**: Store complex objects dengan type safety
- **Simple Storage**: Key-value pairs untuk preferences
- **Offline Support**: Cache data untuk offline access
- **Performance**: Fast read/write operations
- **Security**: Optional encryption untuk sensitive data

---

## Storage Options

### 1. Hive - Primary Storage

**Use Cases**:

- Accounts, Contacts, Tasks
- Visit Reports
- Dashboard data
- Sync queue
- User permissions

**Advantages**:

- Type-safe dengan code generation
- Binary format (fast I/O)
- Cross-platform
- Lazy loading support

### 2. SharedPreferences - Simple Storage

**Use Cases**:

- User preferences
- Theme settings
- Language selection
- Onboarding flags
- Simple booleans/strings

**Advantages**:

- Simple API
- Built-in Flutter support
- No setup required

### 3. Flutter Secure Storage - Sensitive Data

**Use Cases**:

- Auth tokens (production)
- Encryption keys
- Sensitive user data

**Advantages**:

- Platform keystore/keychain
- Encrypted storage
- Secure by default

---

## Business Rules

### 1. Storage Type Selection

**Use Hive when**:

- Complex objects dengan multiple fields
- Need type safety
- Need query/filter capabilities
- Offline-first features

**Use SharedPreferences when**:

- Simple primitives (bool, int, string)
- User preferences
- Feature flags
- Simple caching

**Use Secure Storage when**:

- Auth tokens
- API keys
- Sensitive user data
- Production apps

### 2. Data Retention

**Cache TTL**:

- Dashboard data: 5 minutes
- Account list: 15 minutes
- Account detail: 30 minutes
- User permissions: 60 minutes

**Permanent Storage**:

- User preferences
- Settings
- Onboarding status
- Auth tokens (until logout)

### 3. Storage Keys

**Naming Convention**:

- Prefix dengan feature: `account_`, `task_`
- Use snake_case: `user_preferences`
- Include version jika ada migration: `cache_v2`

### 4. Cleanup Rules

**Auto-cleanup**:

- Expired cache entries
- Old sync queue items (> 7 days)
- Orphaned data

**Manual Cleanup**:

- Logout: Clear all user-specific data
- Clear cache: User-initiated cache clear

---

## Keputusan Teknis & Trade-offs

### Mengapa Hive, bukan SQLite?

**Keputusan**: Menggunakan Hive daripada SQLite.

**Alasan**:

- **Simplicity**: No SQL queries needed
- **Performance**: Binary format lebih cepat
- **Type Safety**: Code generation untuk models
- **Flutter Native**: Better Flutter integration

**Trade-off**: Limited query capabilities. **Mitigasi**: Use in-memory filtering untuk complex queries.

### Mengapa Tidak Semua Data di Secure Storage?

**Keputusan**: Use secure storage hanya untuk sensitive data.

**Alasan**:

- **Performance**: Secure storage lebih lambat
- **Simplicity**: Hive lebih mudah untuk development
- **Use Case**: Not all data needs encryption

**Trade-off**: Data di Hive tidak encrypted. **Mitigasi**: Enable Hive encryption untuk production.

---

## Struktur Folder

```
apps/mobile/lib/
├── core/
│   └── storage/
│       ├── hive_storage.dart              # Hive initialization
│       ├── hive_adapters.dart             # Type adapters
│       ├── shared_prefs_storage.dart      # SharedPreferences wrapper
│       ├── secure_storage.dart            # Secure storage wrapper
│       ├── cache_manager.dart             # Cache TTL management
│       └── offline_storage.dart           # Offline data helper
│
├── features/
│   └── [feature]/
│       └── data/
│           └── models/
│               └── [model].g.dart         # Generated Hive adapters
```

---

## Storage Implementation

### Hive Setup

**File**: `core/storage/hive_storage.dart`

```dart
class HiveStorage {
  static bool _initialized = false;

  static Future<void> init() async {
    if (_initialized) return;

    // Initialize Hive
    await Hive.initFlutter();

    // Register adapters
    _registerAdapters();

    // Open boxes
    await _openBoxes();

    _initialized = true;
  }

  static void _registerAdapters() {
    Hive.registerAdapter(AccountAdapter());
    Hive.registerAdapter(ContactAdapter());
    Hive.registerAdapter(TaskAdapter());
    Hive.registerAdapter(VisitReportAdapter());
    Hive.registerAdapter(UserPermissionsAdapter());
  }

  static Future<void> _openBoxes() async {
    await Hive.openBox<Account>('accounts');
    await Hive.openBox<Contact>('contacts');
    await Hive.openBox<Task>('tasks');
    await Hive.openBox<VisitReport>('visit_reports');
    await Hive.openBox<UserPermissions>('permissions');
    await Hive.openBox<Map>('cache_meta'); // For TTL tracking
  }

  static Box<T> getBox<T>(String name) {
    return Hive.box<T>(name);
  }

  static Future<void> clearAll() async {
    await Hive.deleteFromDisk();
    _initialized = false;
  }
}
```

### Cache Manager dengan TTL

**File**: `core/storage/cache_manager.dart`

```dart
class CacheManager {
  final Box<Map> _metaBox;

  CacheManager(this._metaBox);

  Future<void> setWithTTL<T>(
    Box<T> box,
    String key,
    T value, {
    required Duration ttl,
  }) async {
    await box.put(key, value);

    // Store metadata
    await _metaBox.put('${key}_meta', {
      'cached_at': DateTime.now().toIso8601String(),
      'ttl_minutes': ttl.inMinutes,
    });
  }

  T? getWithTTL<T>(Box<T> box, String key) {
    final value = box.get(key);
    if (value == null) return null;

    // Check TTL
    final meta = _metaBox.get('${key}_meta');
    if (meta != null) {
      final cachedAt = DateTime.parse(meta['cached_at']);
      final ttl = Duration(minutes: meta['ttl_minutes']);

      if (DateTime.now().difference(cachedAt) > ttl) {
        // Expired, clear and return null
        box.delete(key);
        _metaBox.delete('${key}_meta');
        return null;
      }
    }

    return value;
  }

  Future<void> clearExpired() async {
    final now = DateTime.now();

    for (final key in _metaBox.keys) {
      if (key.toString().endsWith('_meta')) {
        final meta = _metaBox.get(key);
        if (meta != null) {
          final cachedAt = DateTime.parse(meta['cached_at']);
          final ttl = Duration(minutes: meta['ttl_minutes']);

          if (now.difference(cachedAt) > ttl) {
            final dataKey = key.toString().replaceAll('_meta', '');
            _metaBox.delete(key);
            // Note: Actual data cleanup done by feature repositories
          }
        }
      }
    }
  }
}
```

### SharedPreferences Wrapper

**File**: `core/storage/shared_prefs_storage.dart`

```dart
class SharedPrefsStorage {
  static SharedPreferences? _prefs;

  static Future<void> init() async {
    _prefs = await SharedPreferences.getInstance();
  }

  // String
  static Future<bool> setString(String key, String value) async {
    return await _prefs?.setString(key, value) ?? false;
  }

  static String? getString(String key) {
    return _prefs?.getString(key);
  }

  // Bool
  static Future<bool> setBool(String key, bool value) async {
    return await _prefs?.setBool(key, value) ?? false;
  }

  static bool? getBool(String key) {
    return _prefs?.getBool(key);
  }

  // Int
  static Future<bool> setInt(String key, int value) async {
    return await _prefs?.setInt(key, value) ?? false;
  }

  static int? getInt(String key) {
    return _prefs?.getInt(key);
  }

  // Remove
  static Future<bool> remove(String key) async {
    return await _prefs?.remove(key) ?? false;
  }

  // Clear all
  static Future<bool> clear() async {
    return await _prefs?.clear() ?? false;
  }
}
```

### Secure Storage Wrapper

**File**: `core/storage/secure_storage.dart`

```dart
class SecureStorage {
  static const _storage = FlutterSecureStorage(
    aOptions: AndroidOptions(
      encryptedSharedPreferences: true,
    ),
    iOptions: IOSOptions(
      accountName: 'flutter_secure_storage',
    ),
  );

  static Future<void> write(String key, String value) async {
    await _storage.write(key: key, value: value);
  }

  static Future<String?> read(String key) async {
    return await _storage.read(key: key);
  }

  static Future<void> delete(String key) async {
    await _storage.delete(key: key);
  }

  static Future<void> deleteAll() async {
    await _storage.deleteAll();
  }
}
```

---

## Usage Examples

### Hive Usage in Repository

```dart
class AccountRepository {
  final Box<Account> _accountBox;
  final CacheManager _cacheManager;

  AccountRepository(this._accountBox, this._cacheManager);

  Future<void> cacheAccounts(List<Account> accounts) async {
    for (final account in accounts) {
      await _cacheManager.setWithTTL(
        _accountBox,
        account.id,
        account,
        ttl: const Duration(minutes: 15),
      );
    }
  }

  List<Account> getCachedAccounts() {
    final accounts = <Account>[];

    for (final key in _accountBox.keys) {
      final account = _cacheManager.getWithTTL(_accountBox, key.toString());
      if (account != null) {
        accounts.add(account);
      }
    }

    return accounts;
  }
}
```

### SharedPreferences Usage

```dart
// Save user preference
await SharedPrefsStorage.setString('language', 'id');
await SharedPrefsStorage.setBool('dark_mode', true);

// Read user preference
final language = SharedPrefsStorage.getString('language') ?? 'id';
final isDarkMode = SharedPrefsStorage.getBool('dark_mode') ?? false;
```

### Secure Storage Usage

```dart
// Production: Store tokens securely
await SecureStorage.write('access_token', token);

// Read token
final token = await SecureStorage.read('access_token');
```

---

## Best Practices

### 1. Model Design

```dart
@HiveType(typeId: 1)
class Account extends HiveObject {
  @HiveField(0)
  late String id;

  @HiveField(1)
  late String name;

  @HiveField(2)
  String? email;

  @HiveField(3)
  DateTime? cachedAt;
}
```

**Rules**:

- Always extend HiveObject
- Use @HiveType dengan unique typeId
- Use @HiveField dengan sequential indices
- Include cache timestamp

### 2. Type ID Management

```dart
// constants.dart
class HiveTypeIds {
  static const int account = 1;
  static const int contact = 2;
  static const int task = 3;
  static const int visitReport = 4;
  static const int userPermissions = 5;
  // ... add more
}
```

### 3. Error Handling

```dart
try {
  final account = await accountBox.get(id);
} catch (e) {
  // Handle Hive errors
  print('Hive error: $e');
  // Fallback to API
}
```

---

## Cara Test Manual

1. **Data Persistence**:
   - Save data ke Hive
   - Kill app
   - Reopen app
   - Verifikasi: Data masih tersedia

2. **Cache Expiration**:
   - Save data dengan TTL 1 menit
   - Wait 2 menit
   - Try retrieve data
   - Verifikasi: Data expired dan null

3. **Clear Storage**:
   - Logout
   - Verifikasi: All user data cleared

---

## Dependencies

```yaml
dependencies:
  hive: ^2.2.3
  hive_flutter: ^1.1.0
  shared_preferences: ^2.2.0
  flutter_secure_storage: ^9.0.0

dev_dependencies:
  hive_generator: ^2.0.1
  build_runner: ^2.4.13
```

---

## Notes & Improvements

### Future Improvements

1. **Hive Encryption**: Enable encryption untuk production
2. **Storage Migration**: Versioning untuk schema changes
3. **Backup**: Auto-backup ke cloud storage
4. **Compression**: Compress large objects sebelum storage

---

**Document Status**: Active  
**Last Updated**: January 2025
