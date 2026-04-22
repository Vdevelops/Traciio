# Testing Guide for Clean Architecture Flutter Apps

This guide covers testing strategies for Flutter apps built with Clean Architecture, Riverpod, and fpdart.

## Testing Pyramid

```
       E2E Tests (Few)
      /              \
     /                \
    Integration Tests
   /                    \
  /                      \
 Unit Tests (Many)
```

- **Unit Tests**: 70% - Test domain logic, use cases, repositories
- **Widget Tests**: 20% - Test UI components and widgets
- **Integration Tests**: 10% - Test full feature flows

## Setup

Add testing dependencies to `pubspec.yaml`:

```yaml
dev_dependencies:
  flutter_test:
    sdk: flutter
  mockito: ^5.4.0
  build_runner: ^2.4.0
  mocktail: ^1.0.0
  flutter_riverpod: ^2.4.0
  riverpod_test: ^2.0.0
```

## Unit Testing Domain Layer

### Testing Entities

Entities are simple data classes with freezed, so testing focuses on equality and copying:

```dart
// test/domain/entities/user_test.dart
import 'package:flutter_test/flutter_test.dart';
import 'package:your_app/features/user/domain/entities/user.dart';

void main() {
  group('User Entity', () {
    test('should create user with required fields', () {
      const user = User(
        id: '1',
        name: 'John Doe',
        email: 'john@example.com',
      );

      expect(user.id, '1');
      expect(user.name, 'John Doe');
      expect(user.email, 'john@example.com');
    });

    test('should support equality comparison', () {
      const user1 = User(id: '1', name: 'John', email: 'john@test.com');
      const user2 = User(id: '1', name: 'John', email: 'john@test.com');

      expect(user1, equals(user2));
    });

    test('should support copyWith', () {
      const user = User(id: '1', name: 'John', email: 'john@test.com');
      final updated = user.copyWith(name: 'Jane');

      expect(updated.name, 'Jane');
      expect(updated.id, '1');
      expect(updated.email, 'john@test.com');
    });
  });
}
```

### Testing Use Cases

Use cases contain business logic and should be thoroughly tested:

```dart
// test/domain/usecases/get_user_test.dart
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';
import 'package:fpdart/fpdart.dart';
import 'package:your_app/core/errors/failures.dart';
import 'package:your_app/features/user/domain/entities/user.dart';
import 'package:your_app/features/user/domain/repositories/user_repository.dart';
import 'package:your_app/features/user/domain/usecases/get_user.dart';

@GenerateMocks([UserRepository])
import 'get_user_test.mocks.dart';

void main() {
  late GetUser useCase;
  late MockUserRepository mockRepository;

  setUp(() {
    mockRepository = MockUserRepository();
    useCase = GetUser(mockRepository);
  });

  const testUser = User(
    id: '1',
    name: 'Test User',
    email: 'test@example.com',
  );

  group('GetUser', () {
    test('should get user from repository when successful', () async {
      // Arrange
      when(mockRepository.getUser('1'))
          .thenAnswer((_) async => const Right(testUser));

      // Act
      final result = await useCase('1');

      // Assert
      expect(result, const Right(testUser));
      verify(mockRepository.getUser('1')).called(1);
      verifyNoMoreInteractions(mockRepository);
    });

    test('should return failure when repository fails', () async {
      // Arrange
      const failure = ServerFailure('Server error');
      when(mockRepository.getUser('1'))
          .thenAnswer((_) async => const Left(failure));

      // Act
      final result = await useCase('1');

      // Assert
      expect(result, const Left(failure));
      verify(mockRepository.getUser('1')).called(1);
    });
  });
}
```

## Unit Testing Data Layer

### Testing Models

Test JSON serialization and entity conversion:

```dart
// test/data/models/user_model_test.dart
import 'dart:convert';
import 'package:flutter_test/flutter_test.dart';
import 'package:your_app/features/user/data/models/user_model.dart';
import 'package:your_app/features/user/domain/entities/user.dart';

void main() {
  const testUserModel = UserModel(
    id: '1',
    name: 'Test User',
    email: 'test@example.com',
  );

  group('UserModel', () {
    test('should be a subclass of User entity', () {
      expect(testUserModel.toEntity(), isA<User>());
    });

    test('should convert from JSON correctly', () {
      final Map<String, dynamic> jsonMap = {
        'id': '1',
        'name': 'Test User',
        'email': 'test@example.com',
      };

      final result = UserModel.fromJson(jsonMap);

      expect(result, testUserModel);
    });

    test('should convert to JSON correctly', () {
      final result = testUserModel.toJson();

      final expectedMap = {
        'id': '1',
        'name': 'Test User',
        'email': 'test@example.com',
      };

      expect(result, expectedMap);
    });

    test('should convert to entity correctly', () {
      final entity = testUserModel.toEntity();

      expect(entity.id, testUserModel.id);
      expect(entity.name, testUserModel.name);
      expect(entity.email, testUserModel.email);
    });
  });
}
```

### Testing Data Sources

Mock data source interfaces for testing:

```dart
// test/data/datasources/user_datasource_test.dart
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';
import 'package:your_app/features/user/data/datasources/user_datasource.dart';
import 'package:your_app/features/user/data/models/user_model.dart';

@GenerateMocks([UserDataSource])
import 'user_datasource_test.mocks.dart';

void main() {
  late MockUserDataSource mockDataSource;

  setUp(() {
    mockDataSource = MockUserDataSource();
  });

  group('UserDataSource', () {
    const userId = '1';
    final testUserModel = UserModel(
      id: '1',
      name: 'Test User',
      email: 'test@example.com',
    );

    test('should return UserModel when getUser is successful', () async {
      // Arrange
      when(mockDataSource.getUser(userId))
          .thenAnswer((_) async => testUserModel);

      // Act
      final result = await mockDataSource.getUser(userId);

      // Assert
      expect(result, equals(testUserModel));
      verify(mockDataSource.getUser(userId)).called(1);
    });

    test('should throw exception when getUser fails', () async {
      // Arrange
      when(mockDataSource.getUser(userId))
          .thenThrow(Exception('Failed to fetch user'));

      // Act & Assert
      expect(
        () => mockDataSource.getUser(userId),
        throwsA(isA<Exception>()),
      );
    });
  });
}
```

**Note**: This approach tests the abstract interface. Your actual data source implementation
(Dio, http, GraphQL, database, etc.) should be tested separately based on its specific implementation.

### Testing Repository Implementation

Test error handling and data transformation:

```dart
// test/data/repositories/user_repository_impl_test.dart
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';
import 'package:fpdart/fpdart.dart';
import 'package:your_app/core/errors/failures.dart';
import 'package:your_app/features/user/data/datasources/user_datasource.dart';
import 'package:your_app/features/user/data/models/user_model.dart';
import 'package:your_app/features/user/data/repositories/user_repository_impl.dart';

@GenerateMocks([UserDataSource])
import 'user_repository_impl_test.mocks.dart';

void main() {
  late UserRepositoryImpl repository;
  late MockUserDataSource mockDataSource;

  setUp(() {
    mockDataSource = MockUserDataSource();
    repository = UserRepositoryImpl(mockDataSource);
  });

  const testUserModel = UserModel(
    id: '1',
    name: 'Test User',
    email: 'test@example.com',
  );

  group('getUser', () {
    test('should return User when datasource call is successful', () async {
      // Arrange
      when(mockDataSource.getUser('1'))
          .thenAnswer((_) async => testUserModel);

      // Act
      final result = await repository.getUser('1');

      // Assert
      expect(result.isRight(), true);
      result.fold(
        (failure) => fail('Should not return failure'),
        (user) {
          expect(user.id, '1');
          expect(user.name, 'Test User');
        },
      );
    });

    test('should return ServerFailure when exception occurs', () async {
      // Arrange
      when(mockDataSource.getUser('1')).thenThrow(
        Exception('Network error'),
      );

      // Act
      final result = await repository.getUser('1');

      // Assert
      expect(result.isLeft(), true);
      result.fold(
        (failure) => expect(failure, isA<ServerFailure>()),
        (user) => fail('Should not return user'),
      );
    });

    test('should return UnexpectedFailure on unknown exception', () async {
      // Arrange
      when(mockDataSource.getUser('1'))
          .thenThrow(Exception('Unknown error'));

      // Act
      final result = await repository.getUser('1');

      // Assert
      expect(result.isLeft(), true);
      result.fold(
        (failure) => expect(failure, isA<UnexpectedFailure>()),
        (user) => fail('Should not return user'),
      );
    });
  });
}
```

## Testing Presentation Layer

### Testing Riverpod Providers

Use `ProviderContainer` for testing providers:

```dart
// test/presentation/providers/user_provider_test.dart
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:mockito/annotations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fpdart/fpdart.dart';
import 'package:your_app/core/errors/failures.dart';
import 'package:your_app/features/user/domain/entities/user.dart';
import 'package:your_app/features/user/domain/usecases/get_user.dart';
import 'package:your_app/features/user/presentation/providers/user_provider.dart';

@GenerateMocks([GetUser])
import 'user_provider_test.mocks.dart';

void main() {
  late MockGetUser mockGetUser;
  late ProviderContainer container;

  setUp(() {
    mockGetUser = MockGetUser();
    container = ProviderContainer(
      overrides: [
        getUserProvider.overrideWithValue(mockGetUser),
      ],
    );
  });

  tearDown(() {
    container.dispose();
  });

  const testUser = User(
    id: '1',
    name: 'Test User',
    email: 'test@example.com',
  );

  group('UserNotifier', () {
    test('should start with null state', () {
      final state = container.read(userNotifierProvider);
      
      expect(state.value, null);
    });

    test('should fetch user successfully', () async {
      // Arrange
      when(mockGetUser('1')).thenAnswer((_) async => const Right(testUser));

      // Act
      await container.read(userNotifierProvider.notifier).fetchUser('1');

      // Assert
      final state = container.read(userNotifierProvider);
      expect(state.value, testUser);
      expect(state.hasError, false);
    });

    test('should handle failure when fetching user', () async {
      // Arrange
      const failure = ServerFailure('Server error');
      when(mockGetUser('1')).thenAnswer((_) async => const Left(failure));

      // Act
      await container.read(userNotifierProvider.notifier).fetchUser('1');

      // Assert
      final state = container.read(userNotifierProvider);
      expect(state.hasError, true);
      expect(state.error, isA<Failure>());
    });
  });
}
```

### Widget Testing

Test widgets with Riverpod:

```dart
// test/presentation/screens/user_screen_test.dart
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:mockito/mockito.dart';
import 'package:your_app/features/user/domain/entities/user.dart';
import 'package:your_app/features/user/presentation/providers/user_provider.dart';
import 'package:your_app/features/user/presentation/screens/user_screen.dart';

void main() {
  const testUser = User(
    id: '1',
    name: 'Test User',
    email: 'test@example.com',
  );

  Widget createWidgetUnderTest() {
    return ProviderScope(
      overrides: [
        userNotifierProvider.overrideWith((ref) {
          return UserNotifierMock();
        }),
      ],
      child: const MaterialApp(
        home: UserScreen(userId: '1'),
      ),
    );
  }

  group('UserScreen', () {
    testWidgets('should display loading indicator initially', (tester) async {
      await tester.pumpWidget(createWidgetUnderTest());
      
      expect(find.byType(CircularProgressIndicator), findsOneWidget);
    });

    testWidgets('should display user data when loaded', (tester) async {
      await tester.pumpWidget(
        ProviderScope(
          overrides: [
            userNotifierProvider.overrideWith((ref) {
              return UserNotifierMock()..state = AsyncData(testUser);
            }),
          ],
          child: const MaterialApp(
            home: UserScreen(userId: '1'),
          ),
        ),
      );
      await tester.pump();

      expect(find.text('Test User'), findsOneWidget);
      expect(find.text('test@example.com'), findsOneWidget);
    });

    testWidgets('should display error message on failure', (tester) async {
      const error = 'Failed to load user';
      
      await tester.pumpWidget(
        ProviderScope(
          overrides: [
            userNotifierProvider.overrideWith((ref) {
              return UserNotifierMock()
                ..state = AsyncError(error, StackTrace.current);
            }),
          ],
          child: const MaterialApp(
            home: UserScreen(userId: '1'),
          ),
        ),
      );
      await tester.pump();

      expect(find.textContaining('Error'), findsOneWidget);
    });

    testWidgets('should call fetchUser when refresh button is tapped', 
        (tester) async {
      final mockNotifier = UserNotifierMock();
      
      await tester.pumpWidget(
        ProviderScope(
          overrides: [
            userNotifierProvider.overrideWith((ref) => mockNotifier),
          ],
          child: const MaterialApp(
            home: UserScreen(userId: '1'),
          ),
        ),
      );

      await tester.tap(find.byType(FloatingActionButton));
      await tester.pump();

      verify(mockNotifier.fetchUser('1')).called(1);
    });
  });
}

class UserNotifierMock extends StateNotifier<AsyncValue<User?>> 
    with Mock 
    implements UserNotifier {
  UserNotifierMock() : super(const AsyncValue.loading());
  
  @override
  Future<void> fetchUser(String id) async {}
}
```

## Integration Testing

Test complete user flows:

```dart
// integration_test/app_test.dart
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:your_app/main.dart' as app;

void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  group('User Flow Integration Tests', () {
    testWidgets('complete login and view profile flow', (tester) async {
      app.main();
      await tester.pumpAndSettle();

      // Should start on login screen
      expect(find.text('Login'), findsOneWidget);

      // Enter credentials
      await tester.enterText(
        find.byType(TextField).first,
        'test@example.com',
      );
      await tester.enterText(
        find.byType(TextField).last,
        'password123',
      );

      // Tap login button
      await tester.tap(find.widgetWithText(ElevatedButton, 'Login'));
      await tester.pumpAndSettle();

      // Should navigate to home screen
      expect(find.text('Home'), findsOneWidget);

      // Navigate to profile
      await tester.tap(find.byIcon(Icons.person));
      await tester.pumpAndSettle();

      // Verify profile screen
      expect(find.text('Profile'), findsOneWidget);
    });
  });
}
```

## Test Coverage

Run tests with coverage:

```bash
# Run all tests
flutter test

# Run with coverage
flutter test --coverage

# Generate HTML coverage report
genhtml coverage/lcov.info -o coverage/html

# Open coverage report
open coverage/html/index.html
```

## Best Practices

1. **Test Naming**: Use descriptive names that explain what is being tested
2. **AAA Pattern**: Arrange, Act, Assert in every test
3. **Mock External Dependencies**: Mock APIs, databases, and third-party services
4. **Test Edge Cases**: Test error scenarios, empty states, and boundary conditions
5. **Isolation**: Each test should be independent and not rely on other tests
6. **Use const**: Use const constructors when possible for better performance
7. **Clean Up**: Always dispose resources in tearDown()
8. **Coverage Goals**: Aim for >80% code coverage
9. **Fast Tests**: Keep unit tests fast (<100ms each)
10. **Readable Assertions**: Use clear, specific assertions

## Common Testing Utilities

```dart
// test/helpers/test_helpers.dart

// Pump and settle with timeout
Future<void> pumpWithTimeout(WidgetTester tester, [Duration? duration]) {
  return tester.pumpAndSettle(duration ?? const Duration(seconds: 5));
}

// Find by key helper
Finder findByKey(String key) => find.byKey(Key(key));

// Create test user
User createTestUser({
  String id = '1',
  String name = 'Test User',
  String email = 'test@example.com',
}) {
  return User(id: id, name: name, email: email);
}

// Mock response helper
Response<T> createMockResponse<T>(
  T data, {
  int statusCode = 200,
  String path = '',
}) {
  return Response(
    data: data,
    statusCode: statusCode,
    requestOptions: RequestOptions(path: path),
  );
}
```
