# Guide - Architecture

## CRM Healthcare Mobile App - Flutter

**Module**: Development Guide  
**Sprint**: Sprint 0  
**Version**: 1.0  
**Status**: ✅ **Completed**  
**Last Updated**: January 2025

---

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Layer Structure](#layer-structure)
3. [Data Layer](#data-layer)
4. [Domain Layer](#domain-layer)
5. [Presentation Layer](#presentation-layer)
6. [Dependency Flow](#dependency-flow)
7. [Best Practices](#best-practices)

---

## Architecture Overview

Project menggunakan **Clean Architecture** dengan **Feature-First** folder structure. Arsitektur ini memisahkan concerns dan memudahkan testing serta maintenance.

### Key Principles

1. **Separation of Concerns**: Clear boundaries antara layers
2. **Dependency Rule**: Dependencies mengarah ke dalam (core &lt;- features)
3. **Testability**: Each layer dapat di-test independently
4. **Scalability**: Easy to add new features

---

## Layer Structure

```
lib/
├── core/                      # Shared infrastructure
│   ├── config/               # Environment, constants
│   ├── network/              # API client, interceptors
│   ├── storage/              # Hive, SharedPrefs
│   ├── routing/              # Navigation
│   ├── errors/               # Error handling
│   └── widgets/              # Shared widgets
│
├── features/                  # Business features
│   ├── auth/
│   ├── accounts/
│   ├── tasks/
│   └── ...
│
└── main.dart
```

---

## Data Layer

### Responsibility

- API calls
- Local storage
- Data models
- Repository pattern

### Structure per Feature

```
features/[feature]/data/
├── models/
│   └── [model].dart         # Data classes
├── repositories/
│   └── [feature]_repository.dart
└── datasources/
    ├── [feature]_remote_datasource.dart
    └── [feature]_local_datasource.dart
```

### Example

```dart
class AccountRepository {
  final AccountRemoteDataSource _remote;
  final AccountLocalDataSource _local;

  Future<List<Account>> getAccounts() async {
    try {
      final accounts = await _remote.getAccounts();
      await _local.cacheAccounts(accounts);
      return accounts;
    } catch (e) {
      return await _local.getCachedAccounts();
    }
  }
}
```

---

## Domain Layer

### Responsibility

- Business logic
- Use cases
- Repository interfaces

### Structure

```
features/[feature]/domain/
├── entities/                 # Core business objects
├── repositories/            # Repository interfaces
└── usecases/                # Business operations
```

---

## Presentation Layer

### Responsibility

- UI components
- State management
- User interactions

### Structure

```
features/[feature]/presentation/
├── screens/                 # Full screens
├── widgets/                 # Reusable widgets
└── providers/               # State management
```

### State Management dengan Riverpod

```dart
// Provider
final accountsProvider = StateNotifierProvider<AccountsNotifier, AccountsState>(
  (ref) => AccountsNotifier(ref.watch(accountRepositoryProvider)),
);

// Notifier
class AccountsNotifier extends StateNotifier<AccountsState> {
  final AccountRepository _repository;

  AccountsNotifier(this._repository) : super(AccountsState.initial());

  Future<void> loadAccounts() async {
    state = state.copyWith(isLoading: true);
    final accounts = await _repository.getAccounts();
    state = state.copyWith(accounts: accounts, isLoading: false);
  }
}
```

---

## Dependency Flow

```
Presentation -> Domain -> Data
     ↑
   Core (infrastructure)
```

**Rules**:

- Presentation depends on Domain
- Domain depends on nothing (interfaces only)
- Data implements Domain interfaces
- All depend on Core

---

## Best Practices

### 1. Feature-First Structure

```
features/
  auth/
    data/
    domain/
    presentation/
  accounts/
    data/
    domain/
    presentation/
```

### 2. Repository Pattern

- Abstract data source
- Enable testing dengan mocks
- Support multiple data sources

### 3. Immutable State

```dart
@freezed
class AccountsState with _$AccountsState {
  const factory AccountsState({
    @Default([]) List<Account> accounts,
    @Default(false) bool isLoading,
    String? error,
  }) = _AccountsState;
}
```

### 4. Error Handling

- Centralized error handling
- User-friendly messages
- Proper error propagation

### 5. Testing Strategy

- Unit tests untuk domain
- Widget tests untuk presentation
- Integration tests untuk flows

---

**Document Status**: Active  
**Last Updated**: January 2025
