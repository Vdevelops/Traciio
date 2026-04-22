# Data Layer Guide

Complete guide for implementing the Data Layer in Clean Architecture.

## Overview

The Data Layer is responsible for:
- **Models**: Data transfer objects with JSON serialization
- **DataSources**: API services, database access, local storage
- **Repositories**: Implementation of domain repository interfaces

## Data Models

### Basic Model with Freezed

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
    @JsonKey(name: 'avatar_url') String? avatarUrl,
  }) = _UserModel;

  factory UserModel.fromJson(Map<String, dynamic> json) =>
      _$UserModelFromJson(json);

  // Convert to domain entity
  User toEntity() => User(
    id: id,
    name: name,
    email: email,
    createdAt: createdAt,
    avatarUrl: avatarUrl,
  );
}
```

### Model with Nested Objects

```dart
@freezed
sealed class OrderModel with _$OrderModel {
  const OrderModel._();

  const factory OrderModel({
    required String id,
    required UserModel user,
    required List<OrderItemModel> items,
    required OrderStatusModel status,
    @JsonKey(name: 'total_amount') required double totalAmount,
  }) = _OrderModel;

  factory OrderModel.fromJson(Map<String, dynamic> json) =>
      _$OrderModelFromJson(json);

  Order toEntity() => Order(
    id: id,
    user: user.toEntity(),
    items: items.map((i) => i.toEntity()).toList(),
    status: status.toEntity(),
    totalAmount: totalAmount,
  );
}

@freezed
sealed class OrderItemModel with _$OrderItemModel {
  const factory OrderItemModel({
    required String productId,
    required int quantity,
    required double price,
  }) = _OrderItemModel;

  factory OrderItemModel.fromJson(Map<String, dynamic> json) =>
      _$OrderItemModelFromJson(json);
}
```

### Model with Enums

```dart
// Enum definition
enum OrderStatusModel {
  @JsonValue('pending')
  pending,
  @JsonValue('processing')
  processing,
  @JsonValue('shipped')
  shipped,
  @JsonValue('delivered')
  delivered,
  @JsonValue('cancelled')
  cancelled,
}

// Model using enum
@freezed
sealed class OrderModel with _$OrderModel {
  const factory OrderModel({
    required String id,
    required OrderStatusModel status,
  }) = _OrderModel;

  factory OrderModel.fromJson(Map<String, dynamic> json) =>
      _$OrderModelFromJson(json);
}
```

## Retrofit API Services

### Basic CRUD Service

```dart
import 'package:retrofit/retrofit.dart';
import 'package:dio/dio.dart';
import '../models/user_model.dart';

part 'user_api_service.g.dart';

@RestApi()
abstract class UserApiService {
  factory UserApiService(Dio dio) = _UserApiService;

  // GET single item
  @GET('/users/{id}')
  Future<UserModel> getUser(@Path('id') String id);

  // GET list
  @GET('/users')
  Future<List<UserModel>> getUsers();

  // POST create
  @POST('/users')
  Future<UserModel> createUser(@Body() Map<String, dynamic> body);

  // PUT update
  @PUT('/users/{id}')
  Future<UserModel> updateUser(
    @Path('id') String id,
    @Body() Map<String, dynamic> body,
  );

  // PATCH partial update
  @PATCH('/users/{id}')
  Future<UserModel> patchUser(
    @Path('id') String id,
    @Body() Map<String, dynamic> body,
  );

  // DELETE
  @DELETE('/users/{id}')
  Future<void> deleteUser(@Path('id') String id);
}
```

### Advanced Query Parameters

```dart
@RestApi()
abstract class ProductApiService {
  factory ProductApiService(Dio dio) = _ProductApiService;

  // Single query params
  @GET('/products')
  Future<List<ProductModel>> getProducts(
    @Query('category') String? category,
    @Query('page') int page,
    @Query('limit') int limit,
    @Query('sort') String sort,
  );

  // Query params as object
  @GET('/products')
  Future<List<ProductModel>> searchProducts(@Queries() Map<String, dynamic> filters);

  // Multiple values for same key
  @GET('/products')
  Future<List<ProductModel>> getProductsByTags(
    @Query('tags') List<String> tags,
  );

  // Path with query
  @GET('/users/{userId}/posts')
  Future<List<PostModel>> getUserPosts(
    @Path('userId') String userId,
    @Query('page') int page,
  );
}
```

### Request/Response Headers

```dart
@RestApi()
abstract class SecureApiService {
  factory SecureApiService(Dio dio) = _SecureApiService;

  // Custom header
  @GET('/protected')
  Future<Data> getProtected(
    @Header('Authorization') String token,
  );

  // Multiple headers
  @GET('/data')
  Future<Data> getData(
    @Header('Authorization') String token,
    @Header('Accept-Language') String language,
  );
}
```

### Form Data & Multipart

```dart
@RestApi()
abstract class UploadApiService {
  factory UploadApiService(Dio dio) = _UploadApiService;

  // Form URL encoded
  @FormUrlEncoded()
  @POST('/auth/login')
  Future<TokenModel> login(
    @Field('email') String email,
    @Field('password') String password,
  );

  // Multipart file upload
  @MULTIPART
  @POST('/upload')
  Future<UploadResponse> uploadFile(
    @Part(name: 'file') File file,
    @Part(name: 'description') String description,
  );

  // Multiple files
  @MULTIPART
  @POST('/upload/multiple')
  Future<UploadResponse> uploadFiles(
    @Part(name: 'files') List<File> files,
  );
}
```

### Response Types

```dart
@RestApi()
abstract class ApiResponseService {
  factory ApiResponseService(Dio dio) = _ApiResponseService;

  // Returns the full response with status code
  @GET('/data')
  Future<Response<DataModel>> getDataWithResponse();

  // Returns just the body
  @GET('/data')
  Future<DataModel> getData();

  // Returns raw string
  @GET('/text')
  Future<String> getText();

  // Returns raw bytes
  @GET('/download')
  Future<List<int>> downloadFile();
}
```

## Repository Implementation

### Basic Repository Pattern

```dart
import 'package:fpdart/fpdart.dart';
import 'package:dio/dio.dart';
import '../../../../core/errors/failures.dart';
import '../../../../core/errors/network_exceptions.dart';
import '../../domain/entities/user.dart';
import '../../domain/repositories/user_repository.dart';
import '../datasources/user_api_service.dart';

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
}
```

### Repository with Caching

```dart
class UserRepositoryImpl implements UserRepository {
  final UserApiService apiService;
  final UserLocalDataSource localDataSource;

  UserRepositoryImpl(this.apiService, this.localDataSource);

  @override
  Future<Either<Failure, List<User>>> getUsers({bool forceRefresh = false}) async {
    try {
      // Try cache first
      if (!forceRefresh) {
        final cached = await localDataSource.getCachedUsers();
        if (cached != null) {
          return Right(cached);
        }
      }

      // Fetch from API
      final userModels = await apiService.getUsers();
      final users = userModels.map((m) => m.toEntity()).toList();

      // Cache the result
      await localDataSource.cacheUsers(users);

      return Right(users);
    } on DioException catch (e) {
      // If network fails, try to return cached data
      final cached = await localDataSource.getCachedUsers();
      if (cached != null) {
        return Right(cached);
      }

      final exception = NetworkExceptions.fromDioError(e);
      return Left(Failure.network(exception.message));
    } catch (e) {
      return Left(Failure.unexpected(e.toString()));
    }
  }
}
```

### Repository with Pagination

```dart
class UserRepositoryImpl implements UserRepository {
  final UserApiService apiService;

  @override
  Future<Either<Failure, PaginatedResult<User>>> getUsers({
    required int page,
    required int limit,
  }) async {
    try {
      final response = await apiService.getUsers(page: page, limit: limit);

      return Right(PaginatedResult(
        data: response.data.map((m) => m.toEntity()).toList(),
        currentPage: response.currentPage,
        totalPages: response.totalPages,
        totalItems: response.totalItems,
      ));
    } on DioException catch (e) {
      final exception = NetworkExceptions.fromDioError(e);
      return Left(Failure.network(exception.message));
    } catch (e) {
      return Left(Failure.unexpected(e.toString()));
    }
  }
}
```

## Model Conversions

### Entity to Model

```dart
// For sending data to API
extension UserEntityExtension on User {
  UserModel toModel() => UserModel(
    id: id,
    name: name,
    email: email,
    createdAt: createdAt,
  );

  Map<String, dynamic> toJson() => {
    'id': id,
    'name': name,
    'email': email,
    'created_at': createdAt?.toIso8601String(),
  };
}
```

### Model to Entity

```dart
// Already included in freezed class
@freezed
sealed class UserModel with _$UserModel {
  const UserModel._();

  // ... model definition ...

  User toEntity() => User(
    id: id,
    name: name,
    email: email,
    createdAt: createdAt,
  );
}
```

## Best Practices

1. **Always use freezed** for models to ensure immutability
2. **Keep models simple** - just JSON serialization and toEntity()
3. **Handle all DioExceptions** in repositories with NetworkExceptions
4. **Convert to entities** at repository boundary
5. **Use Either<Failure, T>** for all repository methods
6. **Name API services** with `ApiService` suffix
7. **One API service per feature/domain**