# Offline Support Implementation - Stale-While-Revalidate Pattern

## CRM Healthcare Mobile App - Flutter

**Module**: Infrastructure  
**Version**: 1.1  
**Status**: ✅ **Completed**  
**Last Updated**: March 2026

---

## Table of Contents

1. [Overview](#overview)
2. [Implementation Status](#implementation-status)
3. [Architecture Pattern](#architecture-pattern)
4. [How It Works](#how-it-works)
5. [Implementation Details](#implementation-details)
6. [Features Implemented](#features-implemented)
7. [Code Structure](#code-structure)
8. [Best Practices](#best-practices)
9. [Testing](#testing)

---

## Overview

CRM Healthcare Mobile App menggunakan **Stale-While-Revalidate (SWR)** pattern untuk offline support. Pattern ini memastikan:

- ✅ **Instant Loading**: Data dari cache ditampilkan segera (tidak ada loading screen)
- ✅ **Fresh Data**: Data terbaru di-fetch di background dan UI di-update otomatis
- ✅ **Offline Support**: Data cache tetap tersedia saat offline
- ✅ **Auto-Sync**: Sinkronisasi otomatis saat koneksi kembali atau app resume
- ✅ **No Manual Refresh**: Semua proses sinkronisasi berjalan otomatis

---

## Implementation Status

### ✅ Completed Features

| Feature                | Repository | Provider | Screen | Status    |
| ---------------------- | ---------- | -------- | ------ | --------- |
| **Schedules**          | ✅         | ✅       | ✅     | Completed |
| **Leads**              | ✅         | ✅       | ✅     | Completed |
| **Accounts**           | ✅         | ✅       | ✅     | Completed |
| **Tasks**              | ✅         | ✅       | ✅     | Completed |
| **Pipeline (Deals)**   | ✅         | ✅       | ✅     | Completed |
| **Contacts**           | ✅         | ✅       | ✅     | Completed |
| **Visit Reports**      | ✅         | ✅       | ✅     | Completed |
| **Route Optimization** | ✅         | ✅       | ✅     | Completed |
| **Dashboard**          | ✅         | ✅       | ✅     | Completed |
| **Profile**            | ✅         | ✅       | N/A    | Completed |

### Core Components

| Component               | File                                      | Status |
| ----------------------- | ----------------------------------------- | ------ |
| **AutoSyncManager**     | `core/sync/auto_sync_manager.dart`        | ✅     |
| **ConnectivityService** | `core/network/connectivity_service.dart`  | ✅     |
| **OfflineStorage**      | `core/storage/offline_storage.dart`       | ✅     |
| **SyncStatusIndicator** | `core/widgets/sync_status_indicator.dart` | ✅     |

---

## Architecture Pattern

### Stale-While-Revalidate Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                        User Opens Screen                        │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│ Step 1: LOAD FROM CACHE                                        │
│ - Ambil data dari Hive storage                                  │
│ - Tampilkan SEKARANG (instant UI)                              │
│ - Tidak ada loading spinner                                     │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│ Step 2: BACKGROUND FETCH (Async)                               │
│ - Fetch data terbaru dari API                                   │
│ - Update cache dengan data baru                                 │
│ - Callback ke UI untuk update state                            │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│ Step 3: UI AUTO-UPDATE                                         │
│ - State di-update dengan data terbaru                          │
│ - UI merefresh otomatis                                        │
│ - User melihat data terbaru tanpa interaksi manual             │
└─────────────────────────────────────────────────────────────────┘
```

### Auto-Sync Triggers

```
┌─────────────────────────────────────────────────────────────────┐
│                    Auto-Sync Triggers                          │
├─────────────────────────────────────────────────────────────────┤
│ 1. 🔄 App Resume                                               │
│    - User kembali ke app dari background                       │
│    - Trigger: didChangeAppLifecycleState                       │
│                                                                 │
│ 2. 📶 Connection Restored                                      │
│    - Device kembali online                                     │
│    - Trigger: ConnectivityService listener                     │
│                                                                 │
│ 3. 🚀 Screen Open                                              │
│    - User buka screen dengan internet ON                       │
│    - Trigger: Initial load + background fetch                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## How It Works

### 1. Repository Layer

Setiap repository memiliki method `get{Feature}s()` dengan parameter `onBackgroundUpdate`:

```dart
Future<LeadListResponse> getLeads({
  int page = 1,
  int perPage = 10,
  String? search,
  bool forceRefresh = false,
  Function(LeadListResponse)? onBackgroundUpdate,
}) async {
  // 1. Return cached data immediately
  if (!forceRefresh && hasCache) {
    final cachedData = await OfflineStorage.getLeads();

    // 2. Start background fetch (don't await!)
    if (_connectivity.isOnline) {
      _fetchAndUpdateInBackground(
        onUpdate: onBackgroundUpdate,
      );
    }

    return cachedData; // Return immediately
  }

  // 3. No cache or force refresh: fetch from API
  return await _fetchFromApi(...);
}
```

### 2. Background Fetch

```dart
Future<void> _fetchAndUpdateInBackground({
  Function(LeadListResponse)? onUpdate,
}) async {
  try {
    // Fetch fresh data from API
    final freshData = await _fetchFromApi(page: 1, perPage: 10);

    // Save to cache
    await OfflineStorage.saveLeads(
      freshData.items.map((lead) => lead.toJson()).toList()
    );

    // Notify UI to update
    if (onUpdate != null) {
      onUpdate(freshData);
    }
  } catch (e) {
    // Silent fail - don't show error for background fetch
  }
}
```

### 3. Provider Layer

Provider menggunakan `WidgetsBindingObserver` untuk lifecycle monitoring:

```dart
class LeadListNotifier extends Notifier<LeadListState>
    with WidgetsBindingObserver {

  @override
  LeadListState build() {
    // Register with AutoSyncManager
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref
          .read(autoSyncManagerProvider.notifier)
          .registerFeature('leads', () => silentRefresh());
    });

    return const LeadListState();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    if (state == AppLifecycleState.resumed) {
      // Auto sync when app comes back to foreground
      ref.read(autoSyncManagerProvider.notifier).syncFeature('leads');
    }
  }

  void _updateStateWithFreshData(LeadListResponse freshData) {
    // Create new list instance (important for Riverpod)
    final newLeads = List<Lead>.from(freshData.items);

    state = LeadListState(
      leads: newLeads,
      pagination: freshData.pagination,
      // ... preserve other state
    );
  }
}
```

### 4. Screen Layer

Setiap screen menampilkan `SyncStatusIndicator` di AppBar:

```dart
AppBar(
  title: Text('Leads'),
  actions: [
    // Shows loading indicator during sync
    SyncStatusIndicator(featureKey: 'leads'),
  ],
)
```

---

## Implementation Details

### Critical Implementation Notes

#### 1. List Instance Recreation

**IMPORTANT**: Riverpod menggunakan identity comparison untuk collections. Saat update state dengan data baru, selalu buat instance List baru:

```dart
// ❌ WRONG - Won't trigger rebuild
state = state.copyWith(leads: freshData.items);

// ✅ CORRECT - Creates new instance
final newLeads = List<Lead>.from(freshData.items);
state = LeadListState(leads: newLeads, ...);
```

#### 2. Model Serialization

Pastikan semua field object diserialisasi di method `toJson()`:

```dart
// ❌ WRONG - Missing nested object
Map<String, dynamic> toJson() {
  return {
    'id': id,
    'name': name,
    // leadStatus is missing!
  };
}

// ✅ CORRECT - All fields included
Map<String, dynamic> toJson() {
  return {
    'id': id,
    'name': name,
    'lead_status_ref': leadStatus?.toJson(), // Include nested object
  };
}
```

#### 3. Cache Conditions

Cache hanya disimpan untuk:

- First page only (`page == 1`)
- No active filters (search, status, etc.)
- No pagination beyond first page

```dart
if (page == 1 &&
    (search == null || search.isEmpty) &&
    status == null &&
    ...) {
  await OfflineStorage.saveLeads(leadsJson);
}
```

---

## Features Implemented

### 1. Schedules

- **Repository**: `schedules/data/schedule_repository.dart`
- **Provider**: `schedules/application/schedule_provider.dart`
- **Screen**: `schedules/presentation/screens/schedule_list_screen.dart`
- **Storage**: `OfflineStorage.getSchedules()`, `OfflineStorage.saveSchedules()`

### 2. Leads

- **Repository**: `leads/data/lead_repository.dart`
- **Provider**: `leads/application/lead_provider.dart`
- **Screen**: `leads/presentation/lead_list_screen.dart`
- **Storage**: `OfflineStorage.getLeads()`, `OfflineStorage.saveLeads()`
- **Note**: Fixed leadStatus serialization in `lead_model.dart`

### 3. Accounts

- **Repository**: `accounts/data/account_repository.dart`
- **Provider**: `accounts/application/account_provider.dart`
- **Screen**: `accounts/presentation/account_list_screen.dart`
- **Storage**: `OfflineStorage.getAccounts()`, `OfflineStorage.saveAccounts()`

### 4. Tasks

- **Repository**: `tasks/data/task_repository.dart`
- **Provider**: `tasks/application/task_provider.dart`
- **Screen**: `tasks/presentation/task_list_screen.dart`
- **Storage**: `OfflineStorage.getTasks()`, `OfflineStorage.saveTasks()`

### 5. Pipeline (Deals)

- **Repository**: `pipeline/data/pipeline_repository.dart`
- **Provider**: `pipeline/application/pipeline_provider.dart`
- **Screen**: `pipeline/presentation/screens/pipeline_screen.dart`
- **Storage**: `OfflineStorage.getDeals()`, `OfflineStorage.saveDeals()`

### 6. Contacts

- **Repository**: `contacts/data/contact_repository.dart`
- **Provider**: `contacts/application/contact_provider.dart`
- **Screen**: `contacts/presentation/contact_list_screen.dart`
- **Storage**: `OfflineStorage.getContacts()`, `OfflineStorage.saveContacts()`

### 7. Visit Reports

- **Repository**: `visit_reports/data/visit_report_repository.dart`
- **Provider**: `visit_reports/application/visit_report_provider.dart`
- **Screen**: `visit_reports/presentation/visit_report_list_screen.dart`
- **Storage**: `OfflineStorage.getVisitReports()`, `OfflineStorage.saveVisitReports()`

### 8. Route Optimization

- **Repository**: `route_optimization/data/route_optimization_repository.dart`
- **Provider**: `route_optimization/application/route_optimization_provider.dart`
- **Screen**: `route_optimization/presentation/route_list_screen.dart`
- **Storage**: `RouteOptimizationCache.getCachedRouteList()`, `RouteOptimizationCache.cacheRouteList()`
- **Note**: Fixed `account` field serialization in `waypoint.dart`

### 9. Dashboard

- **Repository**: `dashboard/data/dashboard_repository.dart`
- **Provider**: `dashboard/application/dashboard_provider.dart`
- **Screen**: `dashboard/presentation/dashboard_screen.dart`
- **Storage**: `OfflineStorage.getMobileDashboardOverview()`, `OfflineStorage.saveMobileDashboardOverview()`
- **Storage**: `OfflineStorage.getMobileVisits()`, `OfflineStorage.saveMobileVisits()`
- **Storage**: `OfflineStorage.getMobileTasks()`

### 10. Profile

- **Repository**: `profile/data/profile_repository.dart`
- **Provider**: Uses `FutureProvider` with offline support via `profileRepository.getMyProfile()`
- **Screen**: Integrated in Dashboard header
- **Storage**: `OfflineStorage.getUserProfile()`, `OfflineStorage.saveUserProfile()`

---

## Code Structure

### Pattern Implementation Checklist

Untuk setiap feature, pastikan implementasi berikut sudah lengkap:

#### Repository (`{feature}_repository.dart`)

```dart
Future<FeatureListResponse> getFeatures({
  bool forceRefresh = false,
  Function(FeatureListResponse)? onBackgroundUpdate,
}) async {
  // 1. Return cached data
  if (!forceRefresh && hasCache) {
    final cached = await OfflineStorage.getFeatures();

    // 2. Background fetch
    if (_connectivity.isOnline) {
      _fetchAndUpdateInBackground(onUpdate: onBackgroundUpdate);
    }

    return cached;
  }

  // 3. Fetch from API
  return await _fetchFromApi();
}

Future<void> _fetchAndUpdateInBackground({
  Function(FeatureListResponse)? onUpdate,
}) async {
  final fresh = await _fetchFromApi();
  await OfflineStorage.saveFeatures(fresh.items.map(...));
  onUpdate?.call(fresh);
}
```

#### Provider (`{feature}_provider.dart`)

```dart
class FeatureListNotifier extends Notifier<FeatureListState>
    with WidgetsBindingObserver {

  @override
  FeatureListState build() {
    // Register with AutoSyncManager
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(autoSyncManagerProvider.notifier)
          .registerFeature('feature', () => silentRefresh());
    });
    return const FeatureListState();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    if (state == AppLifecycleState.resumed) {
      ref.read(autoSyncManagerProvider.notifier).syncFeature('feature');
    }
  }

  Future<void> loadFeatures() async {
    final response = await _repository.getFeatures(
      onBackgroundUpdate: _updateStateWithFreshData,
    );
    // ... update state
  }

  void _updateStateWithFreshData(FeatureListResponse fresh) {
    state = FeatureListState(
      items: List<Feature>.from(fresh.items), // New instance!
      // ...
    );
  }
}
```

#### Screen (`{feature}_list_screen.dart`)

```dart
AppBar(
  actions: [
    SyncStatusIndicator(featureKey: 'feature'),
  ],
)
```

---

## Best Practices

### 1. Cache Invalidation

Cache di-clear hanya saat:

- User melakukan pull-to-refresh
- Data berhasil di-create/update/delete
- User mengganti filter

```dart
Future<void> refresh() async {
  _cache.clearPrefix('list:features');
  await loadFeatures(page: 1, refresh: true, forceRefresh: true);
}
```

### 2. Error Handling

Background fetch **tidak boleh** menampilkan error ke user:

```dart
try {
  final freshData = await _fetchFromApi();
  onUpdate?.call(freshData);
} catch (e) {
  // Silent fail - user masih melihat cached data
}
```

### 3. Sync Status Indicator

Tampilkan indikator loading saat background sync berlangsung:

```dart
// SyncStatusIndicator shows spinner when sync is active
// No refresh button - fully automated
```

### 4. Model Updates

Saat menambahkan field baru ke model, pastikan:

1. Field diserialisasi di `toJson()`
2. Field di-deserialisasi di `fromJson()`
3. Hive adapter di-update (jika menggunakan Hive)

---

## Testing

### Manual Testing Checklist

1. **First Open (No Cache)**
   - [ ] Buka app pertama kali (cache kosong)
   - [ ] Pastikan data load dari API
   - [ ] Cache tersimpan di storage

2. **Second Open (With Cache)**
   - [ ] Tutup app (force stop)
   - [ ] Buka lagi
   - [ ] Data dari cache tampil langsung (instant)
   - [ ] Background fetch berjalan
   - [ ] UI update dengan data terbaru

3. **Data Changes**
   - [ ] Ubah data di backend/admin panel
   - [ ] Buka screen di mobile
   - [ ] Data cache tampil dulu
   - [ ] Setelah background fetch: UI update dengan perubahan

4. **Offline Mode**
   - [ ] Matikan internet
   - [ ] Buka screen
   - [ ] Data cache tampil
   - [ ] Sync indicator show "offline"

5. **Auto-Sync**
   - [ ] Minimize app
   - [ ] Ubah data di backend
   - [ ] Resume app
   - [ ] Data auto-sync dan UI update

6. **Connection Restore**
   - [ ] Buka app offline
   - [ ] Nyalakan internet
   - [ ] Auto-sync trigger
   - [ ] Data update otomatis

### Debug Logs (Development)

Untuk debugging, tambahkan log berikut:

```dart
// Repository
debugPrint('Returning cached data: ${cachedData.items.length} items');
debugPrint('Background fetch completed: ${freshData.items.length} items');

// Provider
debugPrint('State updated with fresh data: ${freshData.items.length} items');
```

**Note**: Hapus debug logs sebelum production release.

---

## Summary

CRM Healthcare Mobile App sekarang memiliki **robust offline support** dengan Stale-While-Revalidate pattern yang diimplementasikan di semua fitur utama:

✅ **Schedules** - Jadwal kunjungan tersedia offline  
✅ **Leads** - Data prospek tersedia offline  
✅ **Accounts** - Data akun tersedia offline  
✅ **Tasks** - Task list tersedia offline  
✅ **Pipeline** - Deals dan stages tersedia offline  
✅ **Contacts** - Kontak tersedia offline  
✅ **Visit Reports** - Laporan kunjungan tersedia offline  
✅ **Route Optimization** - Route dan waypoints tersedia offline  
✅ **Dashboard** - Overview, visits, tasks tersedia offline  
✅ **Profile** - Data user tersedia offline

### Keuntungan Pattern Ini:

1. **UX Terbaik**: User tidak pernah melihat loading screen
2. **Data Selalu Fresh**: Background fetch menjaga data up-to-date
3. **Offline-First**: Bekerja tanpa internet
4. **Auto-Sync**: Tidak perlu manual refresh
5. **Battery Efficient**: Tidak fetch data kalau tidak perlu

---

**Next Steps**:

- Monitor performance di production
- Consider implementing request queue untuk write operations (POST/PUT/DELETE)
- Add cache TTL untuk data yang jarang berubah

**Maintainer**: Development Team  
**Last Review**: March 2026
