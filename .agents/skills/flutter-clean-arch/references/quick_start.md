# Quick Start Guide

Step-by-step workflow for creating a new feature in Flutter Clean Architecture.

## Overview

Follow this sequence when creating a new feature:
1. Domain Layer → 2. Data Layer → 3. Presentation Layer

## Step 1: Domain Layer (Pure Business Logic)

### 1.1 Create Entity

`lib/features/user/domain/entities/user.dart`:

```dart
import 'package:freezed_annotation/freezed_annotation.dart';

part 'user.freezed.dart';

@freezed
sealed class User with _$User {
  const factory User({
    required String id,
    required String name,
    required String email,
    @JsonKey(name: 'created_at') DateTime? createdAt,
  }) = _User;
}
```

**Generate**: `dart run build_runner build --delete-conflicting-outputs`

### 1.2 Create Repository Interface

`lib/features/user/domain/repositories/user_repository.dart`:

```dart
import 'package:fpdart/fpdart.dart';
import '../../../../core/errors/failures.dart';
import '../entities/user.dart';

abstract class UserRepository {
  Future<Either<Failure, User>> getUser(String id);
  Future<Either<Failure, List<User>>> getUsers();
  Future<Either<Failure, User>> createUser(CreateUserParams params);
  Future<Either<Failure, User>> updateUser(String id, UpdateUserParams params);
  Future<Either<Failure, void>> deleteUser(String id);
}
```

### 1.3 Create Use Cases

`lib/features/user/domain/usecases/get_user.dart`:

```dart
import 'package:fpdart/fpdart.dart';
import '../../../../core/errors/failures.dart';
import '../entities/user.dart';
import '../repositories/user_repository.dart';

class GetUser {
  final UserRepository repository;

  GetUser(this.repository);

  Future<Either<Failure, User>> call(String id) {
    return repository.getUser(id);
  }
}
```

`lib/features/user/domain/usecases/get_users.dart`:

```dart
import 'package:fpdart/fpdart.dart';
import '../../../../core/errors/failures.dart';
import '../entities/user.dart';
import '../repositories/user_repository.dart';

class GetUsers {
  final UserRepository repository;

  GetUsers(this.repository);

  Future<Either<Failure, List<User>>> call() {
    return repository.getUsers();
  }
}
```

## Step 2: Data Layer

### 2.1 Create Data Model

`lib/features/user/data/models/user_model.dart`:

```dart
import 'package:freezed_annotation/freezed_annotation.dart';
import '../../domain/entities/user.dart';

part 'user_model.freezed.dart';
part 'user_model.g.dart';

@freezed
sealed class UserModel with _$UserModel {
  const UserModel._();

  const factory UserModel({
    required String id,
    required String name,
    required String email,
    @JsonKey(name: 'created_at') DateTime? createdAt,
  }) = _UserModel;

  factory UserModel.fromJson(Map<String, dynamic> json) =>
      _$UserModelFromJson(json);

  User toEntity() => User(
    id: id,
    name: name,
    email: email,
    createdAt: createdAt,
  );
}
```

### 2.2 Create Retrofit API Service

`lib/features/user/data/datasources/user_api_service.dart`:

```dart
import 'package:retrofit/retrofit.dart';
import 'package:dio/dio.dart';
import '../../../../core/network/dio_provider.dart';
import '../models/user_model.dart';

part 'user_api_service.g.dart';

@RestApi()
abstract class UserApiService {
  factory UserApiService(Dio dio) = _UserApiService;

  @GET('/users/{id}')
  Future<UserModel> getUser(@Path('id') String id);

  @GET('/users')
  Future<List<UserModel>> getUsers();

  @POST('/users')
  Future<UserModel> createUser(@Body() Map<String, dynamic> body);

  @PUT('/users/{id}')
  Future<UserModel> updateUser(
    @Path('id') String id,
    @Body() Map<String, dynamic> body,
  );

  @DELETE('/users/{id}')
  Future<void> deleteUser(@Path('id') String id);
}
```

### 2.3 Create Repository Implementation

`lib/features/user/data/repositories/user_repository_impl.dart`:

```dart
import 'package:fpdart/fpdart.dart';
import 'package:dio/dio.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/errors/network_exceptions.dart';
import '../../domain/entities/user.dart';
import '../../domain/repositories/user_repository.dart';
import '../datasources/user_api_service.dart';
import '../models/user_model.dart';

class UserRepositoryImpl implements UserRepository {
  final UserApiService apiService;

  UserRepositoryImpl(this.apiService);

  @override
  Future<Either<Failure, User>> getUser(String id) async {
    try {
      final userModel = await apiService.getUser(id);
      return Right(userModel.toEntity());
    } on DioException catch (e) {
      final exception = NetworkExceptions.fromDioError(e);
      return Left(Failure.network(exception.message));
    } catch (e) {
      return Left(Failure.unexpected(e.toString()));
    }
  }

  @override
  Future<Either<Failure, List<User>>> getUsers() async {
    try {
      final userModels = await apiService.getUsers();
      return Right(userModels.map((m) => m.toEntity()).toList());
    } on DioException catch (e) {
      final exception = NetworkExceptions.fromDioError(e);
      return Left(Failure.network(exception.message));
    } catch (e) {
      return Left(Failure.unexpected(e.toString()));
    }
  }

  @override
  Future<Either<Failure, User>> CreateUser(CreateUserParams params) async {
    try {
      final userModel = await apiService.createUser(params.toJson());
      return Right(userModel.toEntity());
    } on DioException catch (e) {
      final exception = NetworkExceptions.fromDioError(e);
      return Left(Failure.network(exception.message));
    } catch (e) {
      return Left(Failure.unexpected(e.toString()));
    }
  }

  @override
  Future<Either<Failure, User>> updateUser(String id, UpdateUserParams params) async {
    try {
      final userModel = await apiService.updateUser(id, params.toJson());
      return Right(userModel.toEntity());
    } on DioException catch (e) {
      final exception = NetworkExceptions.fromDioError(e);
      return Left(Failure.network(exception.message));
    } catch (e) {
      return Left(Failure.unexpected(e.toString()));
    }
  }

  @override
  Future<Either<Failure, void>> deleteUser(String id) async {
    try {
      await apiService.deleteUser(id);
      return const Right(null);
    } on DioException catch (e) {
      final exception = NetworkExceptions.fromDioError(e);
      return Left(Failure.network(exception.message));
    } catch (e) {
      return Left(Failure.unexpected(e.toString()));
    }
  }
}
```

## Step 3: Presentation Layer

### 3.1 Create Providers

`lib/features/user/presentation/providers/user_provider.dart`:

```dart
import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../../domain/usecases/get_user.dart';
import '../../domain/usecases/get_users.dart';
import '../../data/datasources/user_api_service.dart';
import '../../data/repositories/user_repository_impl.dart';
import '../../../../core/network/dio_provider.dart';

part 'user_provider.g.dart';

@riverpod
UserApiService userApiService(Ref ref) {
  return UserApiService(ref.watch(dioProvider));
}

@riverpod
UserRepositoryImpl userRepository(Ref ref) {
  return UserRepositoryImpl(ref.watch(userApiServiceProvider));
}

@riverpod
GetUser getUser(Ref ref) {
  return GetUser(ref.watch(userRepositoryProvider));
}

@riverpod
GetUsers getUsers(Ref ref) {
  return GetUsers(ref.watch(userRepositoryProvider));
}

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

@riverpod
class UsersNotifier extends _$UsersNotifier {
  @override
  FutureOr<List<User>> build() => [];

  Future<void> fetchUsers() async {
    state = const AsyncLoading();
    final result = await ref.read(getUsersProvider)();

    state = result.fold(
      (failure) => AsyncError(failure, StackTrace.current),
      (users) => AsyncData(users),
    );
  }
}
```

### 3.2 Create Screen

`lib/features/user/presentation/screens/user_list_screen.dart`:

```dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../providers/user_provider.dart';

class UserListScreen extends ConsumerWidget {
  const UserListScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final usersState = ref.watch(usersNotifierProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Users'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () => ref.read(usersNotifierProvider.notifier).fetchUsers(),
          ),
        ],
      ),
      body: usersState.when(
        data: (users) {
          if (users.isEmpty) {
            return const Center(child: Text('No users found'));
          }
          return ListView.builder(
            itemCount: users.length,
            itemBuilder: (context, index) {
              final user = users[index];
              return ListTile(
                title: Text(user.name),
                subtitle: Text(user.email),
                onTap: () {
                  // Navigate to user detail
                },
              );
            },
          );
        },
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, stack) => Center(
          child: Text('Error: $error'),
        ),
      ),
    );
  }
}
```

### 3.3 Create Detail Screen

`lib/features/user/presentation/screens/user_detail_screen.dart`:

```dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../providers/user_provider.dart';

class UserDetailScreen extends ConsumerWidget {
  final String userId;

  const UserDetailScreen({required this.userId, super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final userState = ref.watch(userNotifierProvider);

    useEffect(() {
      ref.read(userNotifierProvider.notifier).fetchUser(userId);
      return null;
    }, []);

    return Scaffold(
      appBar: AppBar(title: const Text('User Details')),
      body: userState.when(
        data: (user) {
          if (user == null) {
            return const Center(child: Text('No user loaded'));
          }
          return Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Name: ${user.name}', style: Theme.of(context).textTheme.titleLarge),
                const SizedBox(height: 8),
                Text('Email: ${user.email}'),
                if (user.createdAt != null) ...[
                  const SizedBox(height: 8),
                  Text('Created: ${user.createdAt}'),
                ],
              ],
            ),
          );
        },
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, stack) => Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text('Error: $error'),
              const SizedBox(height: 16),
              ElevatedButton(
                onPressed: () => ref.read(userNotifierProvider.notifier).fetchUser(userId),
                child: const Text('Retry'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
```

## Code Generation Commands

```bash
# After creating/modifying any annotated files
dart run build_runner build --delete-conflicting-outputs

# Or use watch mode during development
dart run build_runner watch --delete-conflicting-outputs
```

## Checklist

- [ ] Domain: Entity created with freezed
- [ ] Domain: Repository interface defined
- [ ] Domain: Use cases created
- [ ] Data: Model created with freezed + json_serializable
- [ ] Data: Retrofit API service created
- [ ] Data: Repository implementation with error handling
- [ ] Presentation: Providers created with Riverpod
- [ ] Presentation: Screen widgets created
- [ ] All code generated with build_runner
- [ ] Feature tested in the app