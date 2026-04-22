# Guide - Testing

## CRM Healthcare Mobile App - Flutter

**Module**: Development Guide  
**Sprint**: All Sprints  
**Version**: 1.0  
**Status**: ✅ **Completed**  
**Last Updated**: January 2025

---

## Table of Contents

1. [Testing Strategy](#testing-strategy)
2. [Unit Testing](#unit-testing)
3. [Widget Testing](#widget-testing)
4. [Integration Testing](#integration-testing)
5. [Running Tests](#running-tests)
6. [Best Practices](#best-practices)

---

## Testing Strategy

### Testing Pyramid

```
    /\
   /  \  Integration Tests (Few)
  /----\
 /      \ Widget Tests (Some)
/________\ Unit Tests (Many)
```

### Coverage Targets

- **Unit Tests**: 70%+ coverage
- **Widget Tests**: Critical screens only
- **Integration Tests**: Main user flows

---

## Unit Testing

### Setup

```yaml
dev_dependencies:
  flutter_test:
    sdk: flutter
  mocktail: ^1.0.0
  build_runner: ^2.4.0
```

### Testing Repository

```dart
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';

class MockApiClient extends Mock implements ApiClient {}

void main() {
  group('AccountRepository', () {
    late AccountRepository repository;
    late MockApiClient mockApi;

    setUp(() {
      mockApi = MockApiClient();
      repository = AccountRepository(mockApi);
    });

    test('getAccounts returns list of accounts', () async {
      // Arrange
      when(() => mockApi.get('/api/v1/accounts'))
          .thenAnswer((_) async => Response(
                data: {'data': {'items': []}},
                statusCode: 200,
                requestOptions: RequestOptions(),
              ));

      // Act
      final accounts = await repository.getAccounts();

      // Assert
      expect(accounts, isEmpty);
      verify(() => mockApi.get('/api/v1/accounts')).called(1);
    });

    test('getAccounts throws exception on error', () async {
      // Arrange
      when(() => mockApi.get('/api/v1/accounts'))
          .thenThrow(DioException(
            requestOptions: RequestOptions(),
          ));

      // Act & Assert
      expect(() => repository.getAccounts(), throwsException);
    });
  });
}
```

### Testing State Notifier

```dart
void main() {
  group('AccountsNotifier', () {
    late AccountsNotifier notifier;
    late MockAccountRepository mockRepo;

    setUp(() {
      mockRepo = MockAccountRepository();
      notifier = AccountsNotifier(mockRepo);
    });

    test('initial state is correct', () {
      expect(notifier.state.accounts, isEmpty);
      expect(notifier.state.isLoading, isFalse);
      expect(notifier.state.error, isNull);
    });

    test('loadAccounts updates state', () async {
      // Arrange
      final accounts = [Account(id: '1', name: 'Test')];
      when(() => mockRepo.getAccounts()).thenAnswer((_) async => accounts);

      // Act
      await notifier.loadAccounts();

      // Assert
      expect(notifier.state.accounts, equals(accounts));
      expect(notifier.state.isLoading, isFalse);
    });
  });
}
```

---

## Widget Testing

### Testing Screen

```dart
void main() {
  testWidgets('AccountListScreen displays accounts', (tester) async {
    // Arrange
    await tester.pumpWidget(
      ProviderScope(
        overrides: [
          accountsProvider.overrideWith((ref) => mockAccountsNotifier),
        ],
        child: const MaterialApp(
          home: AccountsScreen(),
        ),
      ),
    );

    // Act
    await tester.pump();

    // Assert
    expect(find.text('Accounts'), findsOneWidget);
    expect(find.byType(AccountCard), findsWidgets);
  });

  testWidgets('AccountListScreen shows loading', (tester) async {
    // Arrange
    await tester.pumpWidget(
      ProviderScope(
        child: const MaterialApp(
          home: AccountsScreen(),
        ),
      ),
    );

    // Assert
    expect(find.byType(CircularProgressIndicator), findsOneWidget);
  });
}
```

---

## Integration Testing

### Setup

```yaml
dev_dependencies:
  integration_test:
    sdk: flutter
```

### Main Flow Test

```dart
// integration_test/app_test.dart
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:crm_healthcare/main.dart' as app;

void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  group('App Test', () {
    testWidgets('Login and view accounts', (tester) async {
      app.main();
      await tester.pumpAndSettle();

      // Login
      await tester.enterText(find.byKey(const Key('email')), 'test@example.com');
      await tester.enterText(find.byKey(const Key('password')), 'password123');
      await tester.tap(find.text('Login'));
      await tester.pumpAndSettle();

      // Navigate ke accounts
      await tester.tap(find.text('Accounts'));
      await tester.pumpAndSettle();

      // Verify
      expect(find.text('Accounts'), findsOneWidget);
    });
  });
}
```

---

## Running Tests

### Unit & Widget Tests

```bash
# Run all tests
flutter test

# Run specific file
flutter test test/features/accounts/data/account_repository_test.dart

# Run dengan coverage
flutter test --coverage

# Generate coverage report
genhtml coverage/lcov.info -o coverage/html
```

### Integration Tests

```bash
# Android
flutter test integration_test/app_test.dart

# iOS
flutter test integration_test/app_test.dart
```

---

## Best Practices

### 1. Test Naming

```dart
// Good
test('getAccounts returns empty list when no accounts', () {});

// Bad
test('test accounts', () {});
```

### 2. Arrange-Act-Assert Pattern

```dart
test('description', () {
  // Arrange
  final input = ...;

  // Act
  final result = ...;

  // Assert
  expect(result, expected);
});
```

### 3. Mock External Dependencies

- API calls
- Storage
- Location services

### 4. Test Edge Cases

- Empty states
- Error states
- Loading states
- Boundary conditions

---

**Document Status**: Active  
**Last Updated**: January 2025
