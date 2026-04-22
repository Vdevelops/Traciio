# Core - State Management

## CRM Healthcare Mobile App - Flutter

**Module**: Core Infrastructure  
**Sprint**: Sprint 0  
**Version**: 1.0  
**Status**: ✅ **Completed**  
**Last Updated**: January 2025

---

## Table of Contents

1. [Ringkasan Fitur](#ringkasan-fitur)
2. [Fitur Utama](#fitur-utama)
3. [Business Rules](#business-rules)
4. [Keputusan Teknis & Trade-offs](#keputusan-teknis--trade-offs)
5. [Struktur Folder](#struktur-folder)
6. [State Patterns](#state-patterns)
7. [Configuration](#configuration)
8. [Usage Examples](#usage-examples)
9. [Cara Test Manual](#cara-test-manual)
10. [Dependencies](#dependencies)
11. [Notes & Improvements](#notes--improvements)

---

## Ringkasan Fitur

Sistem **State Management** mobile app CRM Healthcare menggunakan **Riverpod** dengan pattern **StateNotifier** untuk business logic. Sistem ini menyediakan type-safe, testable, dan scalable state management untuk seluruh aplikasi dengan proper separation of concerns antara UI dan business logic.

### Goals

- **Type Safety**: Compile-time type checking untuk state dan actions
- **Testability**: Business logic terpisah dari UI, mudah di-unit test
- **Scalability**: Pattern yang konsisten untuk semua features
- **Performance**: Efficient rebuilds dengan selective listeners
- **Debuggability**: Time-travel debugging dan state inspection

---

## Fitur Utama

### 1. StateNotifier Pattern

Setiap feature menggunakan StateNotifier untuk manage state:

```dart
class FeatureNotifier extends StateNotifier<FeatureState> {
  final Repository _repository;

  FeatureNotifier(this._repository) : super(FeatureState.initial());

  Future<void> loadData() async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final data = await _repository.fetchData();
      state = state.copyWith(data: data, isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }
}
```

### 2. Immutable State Classes

State classes menggunakan immutable pattern dengan copyWith:

```dart
@freezed
class FeatureState with _$FeatureState {
  const factory FeatureState({
    @Default([]) List<Item> items,
    @Default(false) bool isLoading,
    String? error,
    @Default(1) int currentPage,
  }) = _FeatureState;
}
```

### 3. Provider Types

**StateNotifierProvider**: Untuk business logic dan state management
**FutureProvider**: Untuk async operations (single value)
**StreamProvider**: Untuk real-time data streams
**Provider**: Untuk dependencies dan computed values

### 4. State Consumers

**ConsumerWidget**: Rebuild entire widget when state changes
**Consumer**: Rebuild only specific part of widget tree
**Selector**: Rebuild only when selected part of state changes
**ConsumerStatefulWidget**: Stateful widget dengan Riverpod integration

---

## Business Rules

### 1. State Lifecycle

```
Initial → Loading → Success/Error
  ↑                        │
  └──── Refresh/Reload ────┘
```

**States**:

- **Initial**: State awal, belum ada data
- **Loading**: Sedang fetch data dari API atau local
- **Success**: Data berhasil di-load, tersedia untuk UI
- **Error**: Terjadi error, error message tersedia

### 2. State Update Rules

1. **Immutable Updates**: Selalu buat state baru, jangan mutate state existing
2. **Single Source of Truth**: Setiap data hanya ada di satu provider
3. **Unidirectional Data Flow**: UI → Event → Notifier → State → UI
4. **Async Handling**: Handle loading dan error states untuk setiap async operation

### 3. State Composition

Complex states composed dari multiple providers:

```dart
// Global providers
final authProvider = StateNotifierProvider<AuthNotifier, AuthState>(...);
final connectivityProvider = StreamProvider<bool>(...);

// Feature providers dengan dependencies
final accountsProvider = StateNotifierProvider<AccountsNotifier, AccountsState>((ref) {
  final repository = ref.watch(accountRepositoryProvider);
  final authState = ref.watch(authProvider);
  return AccountsNotifier(repository, authState);
});
```

### 4. State Persistence Rules

1. **Transient State**: State yang tidak perlu persist (loading, error messages)
2. **Persistent State**: State yang perlu persist (auth tokens, user preferences)
3. **Cached State**: State yang bisa di-regenerate (API responses, cached dengan TTL)

---

## Keputusan Teknis & Trade-offs

### Mengapa Riverpod, bukan Provider atau Bloc?

**Keputusan**: Menggunakan Riverpod daripada Provider package atau Bloc pattern.

**Alasan**:

- **Type Safety**: Compile-time type checking untuk providers
- **Performance**: Better performance dengan selective rebuilds
- **Testing**: Mudah override providers di tests
- **Scoping**: AutoDispose providers untuk memory management
- **Code Generation**: Support untuk code generation (Riverpod Generator)

**Trade-off**: Learning curve dan migration dari Provider. **Mitigasi**: Pattern yang consistent dan dokumentasi yang baik.

### Mengapa StateNotifier, bukan ChangeNotifier?

**Keputusan**: Menggunakan StateNotifier daripada ChangeNotifier.

**Alasan**:

- **Immutability**: Enforce immutable state
- **Type Safety**: State type terdefine dengan jelas
- **Debugging**: Easier to track state changes
- **Consistency**: Single way to update state (via state property)

**Trade-off**: Slightly more boilerplate untuk state classes. **Mitigasi**: Menggunakan Freezed untuk generate copyWith dan other methods.

### Mengapa Freezed untuk State Classes?

**Keputusan**: Menggunakan Freezed package untuk state classes.

**Alasan**:

- **copyWith**: Auto-generated copyWith method
- **Equality**: Auto-generated == dan hashCode
- **toString**: Better debugging dengan toString
- **JSON Serialization**: Support untuk fromJson/toJson

**Trade-off**: Code generation step required. **Mitigasi**: Integrasi dengan build_runner di development workflow.

---

## Struktur Folder

```
apps/mobile/lib/
├── features/
│   └── [feature_name]/
│       ├── data/
│       │   ├── models/
│       │   │   └── [feature]_model.dart       # Data models
│       │   └── [feature]_repository.dart      # Repository implementation
│       ├── application/
│       │   ├── [feature]_provider.dart        # StateNotifierProvider
│       │   └── [feature]_state.dart           # State class (Freezed)
│       └── presentation/
│           ├── screens/
│           │   └── [feature]_screen.dart      # Screen widgets
│           └── widgets/
│               └── [feature]_widgets.dart     # Reusable widgets
├── core/
│   └── providers/
│       ├── auth_provider.dart                 # Global auth state
│       ├── connectivity_provider.dart         # Network connectivity
│       └── app_providers.dart                 # Provider definitions
```

---

## State Patterns

### 1. List State Pattern

**File**: `features/accounts/application/accounts_state.dart`

```dart
@freezed
class AccountsState with _$AccountsState {
  const factory AccountsState({
    @Default([]) List<Account> accounts,
    @Default(false) bool isLoading,
    @Default(false) bool isLoadingMore,
    String? error,
    @Default(1) int currentPage,
    @Default(false) bool hasMore,
    @Default('') String searchQuery,
  }) = _AccountsState;

  factory AccountsState.initial() => const AccountsState();
}
```

**File**: `features/accounts/application/accounts_provider.dart`

```dart
final accountsProvider = StateNotifierProvider<AccountsNotifier, AccountsState>((ref) {
  final repository = ref.watch(accountRepositoryProvider);
  return AccountsNotifier(repository);
});

class AccountsNotifier extends StateNotifier<AccountsState> {
  final AccountRepository _repository;

  AccountsNotifier(this._repository) : super(AccountsState.initial()) {
    loadAccounts();
  }

  Future<void> loadAccounts({bool refresh = false}) async {
    if (state.isLoading) return;

    state = state.copyWith(
      isLoading: true,
      error: null,
      currentPage: refresh ? 1 : state.currentPage,
    );

    try {
      final result = await _repository.getAccounts(
        page: state.currentPage,
        search: state.searchQuery,
      );

      final accounts = refresh
          ? result.items
          : [...state.accounts, ...result.items];

      state = state.copyWith(
        accounts: accounts,
        isLoading: false,
        hasMore: result.hasMore,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  Future<void> loadMore() async {
    if (!state.hasMore || state.isLoadingMore) return;

    state = state.copyWith(
      isLoadingMore: true,
      currentPage: state.currentPage + 1,
    );

    await loadAccounts();
    state = state.copyWith(isLoadingMore: false);
  }

  void search(String query) {
    state = state.copyWith(searchQuery: query);
    loadAccounts(refresh: true);
  }
}
```

### 2. Detail State Pattern

```dart
@freezed
class AccountDetailState with _$AccountDetailState {
  const factory AccountDetailState({
    Account? account,
    @Default(false) bool isLoading,
    String? error,
  }) = _AccountDetailState;
}

class AccountDetailNotifier extends StateNotifier<AccountDetailState> {
  final AccountRepository _repository;
  final String _accountId;

  AccountDetailNotifier(this._repository, this._accountId)
      : super(const AccountDetailState()) {
    loadAccount();
  }

  Future<void> loadAccount() async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      final account = await _repository.getAccountById(_accountId);
      state = state.copyWith(account: account, isLoading: false);
    } catch (e) {
      state = state.copyWith(error: e.toString(), isLoading: false);
    }
  }
}

// Provider dengan parameter
final accountDetailProvider = StateNotifierProvider.family
    <AccountDetailNotifier, AccountDetailState, String>((ref, accountId) {
  final repository = ref.watch(accountRepositoryProvider);
  return AccountDetailNotifier(repository, accountId);
});
```

### 3. Form State Pattern

```dart
@freezed
class VisitReportFormState with _$VisitReportFormState {
  const factory VisitReportFormState({
    @Default('') String accountId,
    @Default('') String notes,
    DateTime? visitDate,
    @Default([]) List<String> photoPaths,
    @Default(false) bool isSubmitting,
    String? error,
    @Default(false) bool isSuccess,
  }) = _VisitReportFormState;
}

class VisitReportFormNotifier extends StateNotifier<VisitReportFormState> {
  final VisitReportRepository _repository;

  VisitReportFormNotifier(this._repository)
      : super(const VisitReportFormState());

  void setAccount(String accountId) {
    state = state.copyWith(accountId: accountId);
  }

  void setNotes(String notes) {
    state = state.copyWith(notes: notes);
  }

  void addPhoto(String path) {
    state = state.copyWith(photoPaths: [...state.photoPaths, path]);
  }

  void removePhoto(String path) {
    state = state.copyWith(
      photoPaths: state.photoPaths.where((p) => p != path).toList(),
    );
  }

  Future<bool> submit() async {
    if (state.accountId.isEmpty) return false;

    state = state.copyWith(isSubmitting: true, error: null);

    try {
      await _repository.createVisitReport(
        accountId: state.accountId,
        notes: state.notes,
        visitDate: state.visitDate ?? DateTime.now(),
        photos: state.photoPaths,
      );

      state = state.copyWith(isSubmitting: false, isSuccess: true);
      return true;
    } catch (e) {
      state = state.copyWith(
        isSubmitting: false,
        error: e.toString(),
      );
      return false;
    }
  }

  void reset() {
    state = const VisitReportFormState();
  }
}
```

---

## Configuration

### Provider Initialization

**File**: `core/providers/app_providers.dart`

```dart
// Repository providers
final accountRepositoryProvider = Provider<AccountRepository>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  final localStorage = ref.watch(localStorageProvider);
  return AccountRepository(apiClient, localStorage);
});

// Global providers
final authProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  final repository = ref.watch(authRepositoryProvider);
  return AuthNotifier(repository);
});

final connectivityProvider = StreamProvider<bool>((ref) {
  final service = ref.watch(connectivityServiceProvider);
  return service.onConnectivityChanged;
});

// Feature providers
final accountsProvider = StateNotifierProvider
    <AccountsNotifier, AccountsState>((ref) {
  final repository = ref.watch(accountRepositoryProvider);
  final connectivity = ref.watch(connectivityProvider);
  return AccountsNotifier(repository, connectivity);
});

// Provider dengan autoDispose untuk memory management
final accountDetailProvider = StateNotifierProvider.family
    <AccountDetailNotifier, AccountDetailState, String>((ref, accountId) {
  final repository = ref.watch(accountRepositoryProvider);
  return AccountDetailNotifier(repository, accountId);
});
```

### Main.dart Setup

```dart
void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // Initialize Hive
  await HiveStorage.init();

  runApp(
    ProviderScope(
      child: const MyApp(),
    ),
  );
}

class MyApp extends ConsumerWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = ref.watch(routerProvider);

    return MaterialApp.router(
      routerConfig: router,
      title: 'CRM Healthcare',
      theme: AppTheme.light,
      darkTheme: AppTheme.dark,
    );
  }
}
```

---

## Usage Examples

### 1. Basic Consumer

```dart
class AccountsScreen extends ConsumerWidget {
  const AccountsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(accountsProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('Accounts')),
      body: state.when(
        loading: () => const LoadingIndicator(),
        error: (error) => ErrorWidget(error),
        data: (accounts) => AccountsList(accounts: accounts),
      ),
    );
  }
}
```

### 2. Selective Consumer

```dart
class AccountCount extends ConsumerWidget {
  const AccountCount({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // Hanya rebuild ketika count berubah, bukan seluruh state
    final count = ref.watch(
      accountsProvider.select((state) => state.accounts.length),
    );

    return Text('$count accounts');
  }
}
```

### 3. Action Handler

```dart
class RefreshButton extends ConsumerWidget {
  const RefreshButton({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return IconButton(
      icon: const Icon(Icons.refresh),
      onPressed: () {
        // Call method pada notifier
        ref.read(accountsProvider.notifier).loadAccounts(refresh: true);
      },
    );
  }
}
```

### 4. Multiple Providers

```dart
class DashboardScreen extends ConsumerWidget {
  const DashboardScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final dashboardAsync = ref.watch(dashboardProvider);
    final isOnline = ref.watch(connectivityProvider).value ?? true;

    return Scaffold(
      body: Column(
        children: [
          if (!isOnline) const OfflineIndicator(),
          dashboardAsync.when(
            loading: () => const LoadingIndicator(),
            error: (err, _) => ErrorWidget(err.toString()),
            data: (dashboard) => DashboardContent(data: dashboard),
          ),
        ],
      ),
    );
  }
}
```

---

## Cara Test Manual

### Test State Updates

1. **Loading State**:
   - Buka screen dengan data
   - Verifikasi: Loading indicator muncul
   - Verifikasi: Tidak ada error message

2. **Success State**:
   - Tunggu data di-load
   - Verifikasi: Data ditampilkan dengan benar
   - Verifikasi: Loading indicator hilang

3. **Error State**:
   - Matikan network
   - Coba load data
   - Verifikasi: Error message muncul
   - Verifikasi: Retry button tersedia

4. **State Persistence**:
   - Navigate ke detail screen
   - Switch ke tab lain
   - Kembali ke tab awal
   - Verifikasi: State masih terjaga

5. **Provider Disposal**:
   - Navigate ke detail screen (provider created)
   - Navigate back (provider disposed)
   - Navigate lagi ke detail
   - Verifikasi: Fresh data di-load (provider recreated)

---

## Dependencies

### Internal

- `core/providers/app_providers.dart` - Global provider definitions
- `core/network/api_client.dart` - API client untuk repository
- `core/storage/offline_storage.dart` - Local storage untuk caching

### External

- `flutter_riverpod: ^2.4.0` - State management library
- `riverpod_annotation: ^2.1.0` - Code generation annotations
- `freezed_annotation: ^2.4.0` - Freezed annotations
- `freezed: ^2.4.0` (dev) - Code generation untuk immutable classes
- `build_runner: ^2.4.0` (dev) - Build runner untuk code generation

---

## Notes & Improvements

### Known Limitations

1. **Code Generation**: Build runner memperlambat development cycle. Perlu run `build_runner` setiap kali ubah state classes.

2. **State Restoration**: Belum implement state restoration untuk process death.

3. **No Time-Travel Debugging**: Belum setup Riverpod DevTools untuk debugging.

### Future Improvements

1. **Riverpod Generator**: Migrate ke Riverpod Generator untuk less boilerplate

2. **State Restoration**: Implement state restoration untuk handle process death

3. **DevTools Integration**: Setup Riverpod DevTools untuk debugging

4. **State Persistence**: Auto-persist important states ke local storage

5. **Optimistic Updates**: Implement optimistic updates untuk better UX

6. **State Normalization**: Normalize state untuk complex relationships (accounts-contacts)

---

**Document Status**: Active  
**Last Updated**: January 2025  
**Maintained By**: Dev3 (Mobile Development Team)
