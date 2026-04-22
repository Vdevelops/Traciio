# Error Handling Guide

Complete guide for error handling using fpdart Either and custom Failure types.

## Overview

Error handling in Clean Architecture:
- **Either<Failure, T>**: Functional error handling from fpdart
- **Failure**: Custom sealed class for error types
- **NetworkExceptions**: Dio-specific error handling

## Failure Types

### Base Failure Class

`lib/core/errors/failures.dart`:

```dart
import 'package:freezed_annotation/freezed_annotation.dart';

part 'failures.freezed.dart';

@freezed
sealed class Failure with _$Failure {
  const factory Failure.server(String message) = ServerFailure;
  const factory Failure.network(String message) = NetworkFailure;
  const factory Failure.cache(String message) = CacheFailure;
  const factory Failure.unauthorized(String message) = UnauthorizedFailure;
  const factory Failure.notFound(String message) = NotFoundFailure;
  const factory Failure.validation(String message) = ValidationFailure;
  const factory Failure.unexpected(String message) = UnexpectedFailure;

  const Failure._();

  String get message {
    return when(
      server: (msg) => msg,
      network: (msg) => msg,
      cache: (msg) => msg,
      unauthorized: (msg) => msg,
      notFound: (msg) => msg,
      validation: (msg) => msg,
      unexpected: (msg) => msg,
    );
  }
}
```

### Using Either in Repository Interfaces

```dart
import 'package:fpdart/fpdart.dart';
import '../../../../core/errors/failures.dart';
import '../entities/user.dart';

abstract class UserRepository {
  Future<Either<Failure, User>> getUser(String id);
  Future<Either<Failure, List<User>>> getUsers();
  Future<Either<Failure, User>> createUser(CreateUserParams params);
  Future<Either<Failure, void>> deleteUser(String id);
}
```

## Network Exceptions

### Network Exceptions Class

`lib/core/errors/network_exceptions.dart`:

```dart
import 'package:dio/dio.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

part 'network_exceptions.freezed.dart';

@freezed
sealed class NetworkExceptions with _$NetworkExceptions {
  const factory NetworkExceptions.requestCancelled() = RequestCancelled;
  const factory NetworkExceptions.unauthorizedRequest() = UnauthorizedRequest;
  const factory NetworkExceptions.badRequest() = BadRequest;
  const factory NetworkExceptions.notFound(String resource) = NotFound;
  const factory NetworkExceptions.requestTimeout() = RequestTimeout;
  const factory NetworkExceptions.sendTimeout() = SendTimeout;
  const factory NetworkExceptions.unprocessableEntity() = UnprocessableEntity;
  const factory NetworkExceptions.conflict() = Conflict;
  const factory NetworkExceptions.internalServerError() = InternalServerError;
  const factory NetworkExceptions.noInternetConnection() = NoInternetConnection;
  const factory NetworkExceptions.unexpectedError(String message) = UnexpectedError;

  const NetworkExceptions._();

  String get message {
    return when(
      requestCancelled: () => 'Request was cancelled',
      unauthorizedRequest: () => 'Unauthorized access',
      badRequest: () => 'Invalid request',
      notFound: (resource) => '$resource not found',
      requestTimeout: () => 'Request timeout',
      sendTimeout: () => 'Send timeout',
      unprocessableEntity: () => 'Unable to process request',
      conflict: () => 'Conflict with current state',
      internalServerError: () => 'Internal server error',
      noInternetConnection: () => 'No internet connection',
      unexpectedError: (msg) => msg,
    );
  }

  static NetworkExceptions fromDioError(DioException error) {
    switch (error.type) {
      case DioExceptionType.cancel:
        return const NetworkExceptions.requestCancelled();
      case DioExceptionType.connectionTimeout:
        return const NetworkExceptions.requestTimeout();
      case DioExceptionType.sendTimeout:
        return const NetworkExceptions.sendTimeout();
      case DioExceptionType.receiveTimeout:
        return const NetworkExceptions.requestTimeout();
      case DioExceptionType.badResponse:
        return _handleStatusCode(error.response?.statusCode);
      case DioExceptionType.connectionError:
        return const NetworkExceptions.noInternetConnection();
      default:
        return NetworkExceptions.unexpectedError(
          error.message ?? 'Unknown error occurred',
        );
    }
  }

  static NetworkExceptions _handleStatusCode(int? statusCode) {
    switch (statusCode) {
      case 400:
        return const NetworkExceptions.badRequest();
      case 401:
      case 403:
        return const NetworkExceptions.unauthorizedRequest();
      case 404:
        return const NetworkExceptions.notFound('Resource');
      case 409:
        return const NetworkExceptions.conflict();
      case 408:
        return const NetworkExceptions.requestTimeout();
      case 422:
        return const NetworkExceptions.unprocessableEntity();
      case 500:
      case 502:
      case 503:
        return const NetworkExceptions.internalServerError();
      default:
        return NetworkExceptions.unexpectedError(
          'Received error code: $statusCode',
        );
    }
  }

  Failure toFailure() {
    return when(
      requestCancelled: () => Failure.unexpected(message),
      unauthorizedRequest: () => Failure.unauthorized(message),
      badRequest: () => Failure.validation(message),
      notFound: (_) => Failure.notFound(message),
      requestTimeout: () => Failure.network(message),
      sendTimeout: () => Failure.network(message),
      unprocessableEntity: () => Failure.validation(message),
      conflict: () => Failure.unexpected(message),
      internalServerError: () => Failure.server(message),
      noInternetConnection: () => Failure.network(message),
      unexpectedError: (msg) => Failure.unexpected(msg),
    );
  }
}
```

## Either Patterns

### Pattern 1: Using fold()

Most common pattern for handling Either:

```dart
final result = await useCase(params);

result.fold(
  (failure) {
    // Handle failure
    showError(failure.message);
  },
  (data) {
    // Handle success
    showData(data);
  },
);
```

### Pattern 2: Using match()

Similar to fold but returns a value:

```dart
final result = await useCase(params);

final message = result.match(
  (failure) => 'Error: ${failure.message}',
  (data) => 'Success: $data',
);

print(message);
```

### Pattern 3: Using getOrElse()

Get value or provide default:

```dart
final result = await useCase(params);

final user = result.getOrElse(() => User.empty());
```

### Pattern 4: Using swap()

Flip Either (success becomes failure, failure becomes success):

```dart
final result = await useCase(params);

// If you want to treat success as failure for some reason
final swapped = result.swap();
```

### Pattern 5: Chaining with andThen()

Chain operations that return Either:

```dart
final result = await getUser(userId)
    .andThen((user) => validateUser(user))
    .andThen((user) => saveUser(user));

result.fold(
  (failure) => showError(failure.message),
  (user) => showSuccess(user),
);
```

### Pattern 6: AsyncValue in Riverpod

Convert Either to AsyncValue:

```dart
@riverpod
class UserNotifier extends _$UserNotifier {
  @override
  FutureOr<User?> build() => null;

  Future<void> fetchUser(String id) async {
    state = const AsyncLoading();

    final result = await ref.read(getUserProvider)(id);

    state = result.fold(
      (failure) => AsyncError(failure, StackTrace.current),
      (user) => AsyncData(user),
    );
  }
}
```

## Repository Error Handling

### Basic Error Handling

```dart
@override
Future<Either<Failure, User>> getUser(String id) async {
  try {
    final userModel = await apiService.getUser(id);
    return Right(userModel.toEntity());
  } on DioException catch (e) {
    final exception = NetworkExceptions.fromDioError(e);
    return Left(exception.toFailure());
  } catch (e) {
    return Left(Failure.unexpected(e.toString()));
  }
}
```

### With Specific Error Handling

```dart
@override
Future<Either<Failure, User>> createUser(CreateUserParams params) async {
  try {
    // Validate params first
    final validationError = _validateParams(params);
    if (validationError != null) {
      return Left(Failure.validation(validationError));
    }

    final userModel = await apiService.createUser(params.toJson());
    return Right(userModel.toEntity());
  } on DioException catch (e) {
    final exception = NetworkExceptions.fromDioError(e);

    // Handle specific cases
    if (exception == const NetworkExceptions.unauthorizedRequest()) {
      return Left(Failure.unauthorized('Session expired'));
    }

    if (exception == const NetworkExceptions.conflict()) {
      return Left(Failure.validation('User already exists'));
    }

    return Left(exception.toFailure());
  } catch (e) {
    return Left(Failure.unexpected(e.toString()));
  }
}

String? _validateParams(CreateUserParams params) {
  if (params.name.isEmpty) return 'Name is required';
  if (!params.email.contains('@')) return 'Invalid email';
  return null;
}
```

## UI Error Handling

### Displaying Errors in Widgets

```dart
class UserScreen extends ConsumerWidget {
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final userAsync = ref.watch(userNotifierProvider);

    return Scaffold(
      body: userAsync.when(
        data: (user) => UserView(user: user),
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, stack) {
          // Extract failure message
          final message = error is Failure ? error.message : 'An error occurred';

          return ErrorStateWidget(
            message: message,
            onRetry: () => ref.read(userNotifierProvider.notifier).refresh(),
          );
        },
      ),
    );
  }
}
```

### Error SnackBar

```dart
void _showError(BuildContext context, Failure failure) {
  final color = switch (failure) {
    ServerFailure() => Colors.red,
    NetworkFailure() => Colors.orange,
    ValidationFailure() => Colors.yellow,
    _ => Colors.grey,
  };

  ScaffoldMessenger.of(context).showSnackBar(
    SnackBar(
      content: Text(failure.message),
      backgroundColor: color,
      action: SnackBarAction(
        label: 'Dismiss',
        textColor: Colors.white,
        onPressed: () {},
      ),
    ),
  );
}
```

## Custom Error Extensions

### Error Code Mapping

```dart
extension FailureExtension on Failure {
  int get code {
    return when(
      server: (_) => 500,
      network: (_) => 0,
      cache: (_) => 1,
      unauthorized: (_) => 401,
      notFound: (_) => 404,
      validation: (_) => 400,
      unexpected: (_) => -1,
    );
  }

  String get localizedMessage {
    // Implement localization
    return message;
  }

  bool get isRetryable {
    return map(
      server: (_) => true,
      network: (_) => true,
      cache: (_) => true,
      unauthorized: (_) => false,
      notFound: (_) => false,
      validation: (_) => false,
      unexpected: (_) => false,
    );
  }
}
```

## Best Practices

1. **Always use Either<Failure, T>** for repository methods
2. **Handle all exceptions** in repositories
3. **Convert DioException** to NetworkExceptions then to Failure
4. **Use fold()** for handling Either results
5. **Provide meaningful error messages**
6. **Distinguish between retryable and non-retryable errors**
7. **Log errors** for debugging
8. **Show user-friendly messages** in UI
9. **Don't expose internal errors** to users
10. **Use typed failures** instead of generic exceptions