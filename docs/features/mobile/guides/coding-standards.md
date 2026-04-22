# Guide - Coding Standards

## CRM Healthcare Mobile App - Flutter

**Module**: Development Guide  
**Sprint**: All Sprints  
**Version**: 1.1  
**Status**: ✅ **Completed**  
**Last Updated**: March 2026

---

## Table of Contents

1. [Code Style](#code-style)
2. [Naming Conventions](#naming-conventions)
3. [File Organization](#file-organization)
4. [State Management](#state-management)
5. [Error Handling](#error-handling)
6. [Testing Standards](#testing-standards)
7. [Documentation](#documentation)

---

## Code Style

### Dart Formatting

- Use `dart format` (line length: 80)
- Trailing commas untuk multi-line
- Consistent indentation (2 spaces)

### Lint Rules

```yaml
# analysis_options.yaml
include: package:flutter_lints/flutter.yaml

linter:
  rules:
    prefer_single_quotes: true
    avoid_print: true
    prefer_const_constructors: true
    always_specify_types: false
```

### Import Order

```dart
// Dart imports
import 'dart:async';
import 'dart:convert';

// Package imports
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

// Project imports
import 'package:crm_healthcare/core/config/env.dart';
import 'package:crm_healthcare/features/accounts/data/account_repository.dart';

// Relative imports
import '../widgets/account_card.dart';
```

---

## Naming Conventions

### Files

| Type         | Convention                   | Example                    |
| ------------ | ---------------------------- | -------------------------- |
| Screens      | `snake_case_screen.dart`     | `account_list_screen.dart` |
| Widgets      | `snake_case_widget.dart`     | `account_card.dart`        |
| Models       | `snake_case_model.dart`      | `account_model.dart`       |
| Repositories | `snake_case_repository.dart` | `account_repository.dart`  |
| Providers    | `snake_case_provider.dart`   | `account_provider.dart`    |

### Classes

| Type         | Convention             | Example             |
| ------------ | ---------------------- | ------------------- |
| Screens      | `PascalCaseScreen`     | `AccountListScreen` |
| Widgets      | `PascalCase`           | `AccountCard`       |
| Models       | `PascalCase`           | `Account`           |
| Repositories | `PascalCaseRepository` | `AccountRepository` |
| Providers    | `PascalCaseProvider`   | `AccountProvider`   |
| Notifiers    | `PascalCaseNotifier`   | `AccountNotifier`   |

### Variables & Functions

```dart
// Variables: camelCase
final accountName = 'RS Medika';
final isLoading = false;

// Constants: SCREAMING_SNAKE_CASE
const API_BASE_URL = 'https://api.example.com';

// Functions: camelCase verb prefix
Future<void> loadAccounts() async {}
bool isValidEmail(String email) {}

// Private members: _leadingUnderscore
String _internalHelper() {}
```

---

## File Organization

### Feature Structure

```
features/
  accounts/
    data/
      models/
        account_model.dart
        account_filter.dart
      account_repository.dart
    application/
      account_list_provider.dart
      account_detail_provider.dart
    presentation/
      screens/
        account_list_screen.dart
        account_detail_screen.dart
      widgets/
        account_card.dart
        account_list_item.dart
```

### Max File Length

- Target: < 300 lines
- Maximum: 500 lines
- Split jika terlalu panjang

---

## State Management

### Provider Patterns

```dart
// Good: AutoDispose untuk memory management
final accountsProvider = StateNotifierProvider.autoDispose
    <AccountsNotifier, AccountsState>((ref) {
  return AccountsNotifier(ref.watch(accountRepositoryProvider));
});

// Good: Family untuk parameterized providers
final accountDetailProvider = StateNotifierProvider.family
    <AccountDetailNotifier, AccountDetailState, String>((ref, id) {
  return AccountDetailNotifier(id);
});
```

### State Classes

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

---

## Error Handling

### Always Handle Errors

```dart
// Good
try {
  final accounts = await repository.getAccounts();
  state = state.copyWith(accounts: accounts);
} catch (e) {
  state = state.copyWith(error: e.toString());
}

// Bad
final accounts = await repository.getAccounts(); // No error handling
```

### Use Result Types

```dart
// Good: Using Either type
Future<Either<Failure, List<Account>>> getAccounts() async {
  try {
    final accounts = await _apiClient.getAccounts();
    return Right(accounts);
  } on DioException catch (e) {
    return Left(ServerFailure(e.message));
  }
}
```

---

## Testing Standards

### Test Structure

```dart
group('AccountRepository', () {
  // Setup
  setUp(() {});

  // Tests
  test('should return accounts on success', () async {});
  test('should throw error on failure', () async {});

  // Teardown
  tearDown(() {});
});
```

### Test Naming

```dart
// Pattern: should [expected behavior] when [condition]
test('should return empty list when no accounts exist', () {});
test('should throw ServerException when API returns 500', () {});
```

---

## UI Standards

### Shadow & Elevation

Semua shadow di mobile app harus konsisten dengan web app (Tailwind `shadow-sm`).

**Standard BoxShadow** (equivalent to Tailwind `shadow-sm`):

```dart
// Light mode shadow
BoxShadow(
  color: Colors.black.withValues(alpha: 0.05),
  blurRadius: 3,
  offset: const Offset(0, 1),
)

// Dark mode - no shadow
boxShadow: isDarkMode ? [] : [/* light mode shadow */]
```

**Card Elevation**:

```dart
// Material Card widgets
Card(elevation: 0.5, ...)  // Bukan elevation: 2

// FAB (Floating Action Button)
FloatingActionButton(elevation: 0.5, ...)

// CardTheme di app_theme.dart
cardTheme: CardThemeData(elevation: 0, ...)
```

**Colored/Accent Shadows** (untuk primary cards, markers, avatars):

```dart
// Colored shadow (slightly stronger for visibility)
BoxShadow(
  color: theme.colorScheme.primary.withValues(alpha: 0.15),
  blurRadius: 3,
  offset: const Offset(0, 1),
)

// Map markers
BoxShadow(
  color: markerColor.withValues(alpha: 0.2),
  blurRadius: 3,
  offset: const Offset(0, 1),
)
```

**Rules**:

- ❌ Jangan gunakan `blurRadius > 4` untuk card shadows
- ❌ Jangan gunakan `alpha > 0.1` untuk standard shadows
- ❌ Jangan gunakan `elevation > 1` untuk Card widgets
- ✅ Gunakan conditional shadow untuk dark mode (empty list)
- ✅ Colored shadows boleh sedikit lebih kuat (alpha 0.15-0.2)

---

## Documentation

### Code Comments

```dart
/// Retrieves all accounts for the current user.
///
/// Returns a list of [Account] objects.
/// Throws [ServerException] if the API call fails.
Future<List<Account>> getAccounts() async {
  // Implementation
}
```

### File Headers

```dart
// features/accounts/data/account_repository.dart
// Repository for managing account data operations.
//
// This repository handles both remote API calls and local caching.
```

---

**Document Status**: Active  
**Last Updated**: January 2025
