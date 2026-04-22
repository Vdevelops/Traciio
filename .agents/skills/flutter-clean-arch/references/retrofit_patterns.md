# Retrofit API Patterns

Complete guide for REST API patterns using Retrofit.

## Overview

Retrofit provides type-safe REST API calls with annotations. This guide covers all common patterns.

## Basic Patterns

### GET Request

```dart
@RestApi()
abstract class UserApiService {
  factory UserApiService(Dio dio) = _UserApiService;

  // Simple GET
  @GET('/users/{id}')
  Future<UserModel> getUser(@Path('id') String id);

  // GET with query parameters
  @GET('/users')
  Future<List<UserModel>> getUsers({
    @Query('page') int? page,
    @Query('limit') int? limit,
  });

  // GET with required query params
  @GET('/search')
  Future<List<UserModel>> searchUsers(
    @Query('q') String query,
    @Query('page') int page,
  );
}
```

### POST Request

```dart
@POST('/users')
Future<UserModel> createUser(@Body() Map<String, dynamic> body);

// Or with typed request body
@POST('/users')
Future<UserModel> createUser(@Body() CreateUserRequest request);

// Request body class
class CreateUserRequest {
  final String name;
  final String email;

  CreateUserRequest({required this.name, required this.email});

  Map<String, dynamic> toJson() => {
    'name': name,
    'email': email,
  };
}
```

### PUT/PATCH Request

```dart
// Full update with PUT
@PUT('/users/{id}')
Future<UserModel> updateUser(
  @Path('id') String id,
  @Body() UpdateUserRequest request,
);

// Partial update with PATCH
@PATCH('/users/{id}')
Future<UserModel> patchUser(
  @Path('id') String id,
  @Body() Map<String, dynamic> fields,
);
```

### DELETE Request

```dart
@DELETE('/users/{id}')
Future<void> deleteUser(@Path('id') String id);

// Delete with response
@DELETE('/users/{id}')
Future<UserModel> deleteUserWithResponse(@Path('id') String id);
```

## Query Parameters

### Single Query Parameters

```dart
@GET('/products')
Future<List<ProductModel>> getProducts(
  @Query('category') String? category,
  @Query('sort') String sort,
  @Query('page') int page,
  @Query('limit') int limit,
);
```

### Multiple Query Parameters (Map)

```dart
@GET('/products')
Future<List<ProductModel>> searchProducts(
  @Queries() Map<String, dynamic> filters,
);

// Usage
final products = await apiService.searchProducts({
  'category': 'electronics',
  'price_min': '100',
  'price_max': '1000',
  'brand': 'Apple',
});
```

### Array Query Parameters

```dart
@GET('/products')
Future<List<ProductModel>> getProductsByTags(
  @Query('tags') List<String> tags,
);

// Generates: /products?tags=electronics&tags=phones&tags=apple
```

### Optional Query Parameters

```dart
@GET('/users')
Future<List<UserModel>> getUsers({
  @Query('page') int? page,
  @Query('limit') int? limit,
  @Query('search') String? search,
});
```

## Path Parameters

### Single Path Parameter

```dart
@GET('/users/{id}')
Future<UserModel> getUser(@Path('id') String id);

@GET('/users/{userId}/posts/{postId}')
Future<PostModel> getPost(
  @Path('userId') String userId,
  @Path('postId') String postId,
);
```

### Multiple Path Parameters

```dart
@GET('/organizations/{orgId}/teams/{teamId}/members/{memberId}')
Future<MemberModel> getMember(
  @Path('orgId') String orgId,
  @Path('teamId') String teamId,
  @Path('memberId') String memberId,
);
```

## Headers

### Static Headers (in Dio configuration)

```dart
@riverpod
Dio dio(Ref ref) {
  final options = BaseOptions(
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
      'X-API-Version': 'v1',
    },
  );
  return Dio(options);
}
```

### Dynamic Headers per Request

```dart
@GET('/protected')
Future<Data> getProtectedData(
  @Header('Authorization') String token,
);

@GET('/data')
Future<Data> getDataWithHeaders(
  @Header('Authorization') String token,
  @Header('Accept-Language') String language,
);
```

### Conditional Headers

```dart
@GET('/download')
Future<List<int>> downloadFile(
  @Header('Range') String? range,
);

// Usage
await apiService.downloadFile('bytes=0-1023');
await apiService.downloadFile(null); // No Range header
```

## Request Body

### JSON Body

```dart
@POST('/users')
Future<UserModel> createUser(@Body() CreateUserRequest request);

class CreateUserRequest {
  final String name;
  final String email;

  CreateUserRequest({required this.name, required this.email});

  Map<String, dynamic> toJson() => {
    'name': name,
    'email': email,
  };
}
```

### Map Body

```dart
@POST('/users')
Future<UserModel> createUser(@Body() Map<String, dynamic> body);

// Usage
await apiService.createUser({'name': 'John', 'email': 'john@example.com'});
```

### Form URL Encoded

```dart
@FormUrlEncoded()
@POST('/auth/login')
Future<TokenModel> login(
  @Field('email') String email,
  @Field('password') String password,
);

@FormUrlEncoded()
@POST('/auth/refresh')
Future<TokenModel> refreshToken(
  @Field('refresh_token') String refreshToken,
);
```

### Multipart (File Upload)

```dart
@MULTIPART
@POST('/upload')
Future<UploadResponse> uploadFile(
  @Part(name: 'file') File file,
);

@MULTIPART
@POST('/upload')
Future<UploadResponse> uploadFileWithMetadata(
  @Part(name: 'file') File file,
  @Part(name: 'description') String description,
  @Part(name: 'category') String category,
);

// Multiple files
@MULTIPART
@POST('/upload/multiple')
Future<UploadResponse> uploadMultipleFiles(
  @Part(name: 'files') List<File> files,
);

// With progress
@POST('/upload')
Future<void> uploadWithProgress(
  @Part() File file,
);

// Usage with progress callback
final formData = FormData.fromMap({
  'file': await MultipartFile.fromFile(
    file.path,
    filename: 'upload.jpg',
  ),
});

await dio.post('/upload', data: formData, onSendProgress: (sent, total) {
  print('Progress: ${(sent / total * 100).toStringAsFixed(0)}%');
});
```

## Response Types

### Direct Response Body

```dart
@GET('/users/{id}')
Future<UserModel> getUser(@Path('id') String id);
```

### Full Response with Metadata

```dart
@GET('/users/{id}')
Future<Response<UserModel>> getUserWithResponse(@Path('id') String id);

// Usage
final response = await apiService.getUserWithResponse('123');
print('Status: ${response.statusCode}');
print('Headers: ${response.headers}');
final user = response.data;
```

### String Response

```dart
@GET('/text')
Future<String> getText();
```

### Raw Bytes Response

```dart
@GET('/download')
Future<List<int>> downloadFile();
```

### Void Response (no content)

```dart
@DELETE('/users/{id}')
Future<void> deleteUser(@Path('id') String id);
```

## Advanced Patterns

### Pagination

```dart
@GET('/products')
Future<PaginatedResponse<ProductModel>> getProducts(
  @Query('page') int page,
  @Query('limit') int limit,
);

class PaginatedResponse<T> {
  final List<T> data;
  final int currentPage;
  final int totalPages;
  final int totalItems;

  PaginatedResponse({
    required this.data,
    required this.currentPage,
    required this.totalPages,
    required this.totalItems,
  });

  factory PaginatedResponse.fromJson(
    Map<String, dynamic> json,
    T Function(Object?) fromJsonT,
  ) {
    return PaginatedResponse(
      data: (json['data'] as List)
          .map((item) => fromJsonT(item))
          .toList(),
      currentPage: json['current_page'],
      totalPages: json['total_pages'],
      totalItems: json['total_items'],
    );
  }
}
```

### Filtering and Sorting

```dart
@GET('/products')
Future<List<ProductModel>> getProducts(
  @Query('category') String? category,
  @Query('min_price') double? minPrice,
  @Query('max_price') double? maxPrice,
  @Query('sort_by') String? sortBy,
  @Query('sort_order') String? sortOrder, // 'asc' or 'desc'
);

// Usage
await apiService.getProducts(
  category: 'electronics',
  minPrice: 100,
  maxPrice: 1000,
  sortBy: 'price',
  sortOrder: 'asc',
);
```

### Nested Resources

```dart
@GET('/users/{userId}/posts')
Future<List<PostModel>> getUserPosts(
  @Path('userId') String userId,
  @Query('page') int? page,
);

@POST('/users/{userId}/posts')
Future<PostModel> createUserPost(
  @Path('userId') String userId,
  @Body() CreatePostRequest request,
);

@DELETE('/users/{userId}/posts/{postId}')
Future<void> deleteUserPost(
  @Path('userId') String userId,
  @Path('postId') String postId,
);
```

### Batch Operations

```dart
@POST('/users/batch')
Future<List<UserModel>> createUsers(
  @Body() List<CreateUserRequest> requests,
);

@PUT('/users/batch')
Future<List<UserModel>> updateUsers(
  @Body() List<UpdateUserRequest> requests,
);

@DELETE('/users/batch')
Future<void> deleteUsers(
  @Body() List<String> userIds,
);
```

### Search and Autocomplete

```dart
@GET('/search')
Future<SearchResultModel> search(
  @Query('q') String query,
  @Query('type') SearchType type,
  @Query('page') int? page,
);

enum SearchType { all, users, products, posts }

@GET('/autocomplete')
Future<List<String>> autocomplete(
  @Query('q') String query,
  @Query('limit') int limit,
);
```

## Extra Annotations

### Timeout per Request

```dart
@GET('/slow-endpoint')
Future<Data> getSlowData();

// Not directly supported in Retrofit
// Configure timeout in Dio options instead
```

### Cache Control

```dart
@GET('/data')
@Headers({'Cache-Control': 'no-cache'})
Future<Data> getData();

@GET('/cached-data')
@Headers({'Cache-Control': 'max-age=3600'})
Future<Data> getCachedData();
```

## Best Practices

1. **One API service per feature/domain**
2. **Use typed request classes** for complex bodies
3. **Document expected responses** in comments
4. **Handle all status codes** in error interceptor
5. **Use appropriate HTTP methods** (GET for read, POST for create, etc.)
6. **Name methods clearly** based on their action
7. **Use @Path for IDs**, @Query for filters
8. **Return domain entities** from repositories, not models
9. **Keep services focused** on single responsibility
10. **Use futures** for all async operations