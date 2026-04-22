# Complete Feature Examples

This document provides complete implementations of common Flutter features using Clean Architecture.

**Note**: Data source implementations (HTTP clients, databases, etc.) are not specified. Implement according to your needs using Dio, http, GraphQL, Firebase, SQLite, etc.

## Example 1: Authentication Feature

### Domain Layer

**Entity** (domain/entities/auth_user.dart):
```dart
import 'package:freezed_annotation/freezed_annotation.dart';

part 'auth_user.freezed.dart';

@freezed
sealed class AuthUser with _$AuthUser {
  const factory AuthUser({
    required String id,
    required String email,
    required String name,
    String? avatarUrl,
    required DateTime createdAt,
  }) = _AuthUser;
}
```

**Repository** (domain/repositories/auth_repository.dart):
```dart
import 'package:fpdart/fpdart.dart';
import '../../core/errors/failures.dart';
import '../entities/auth_user.dart';

abstract class AuthRepository {
  Future<Either<Failure, AuthUser>> login(String email, String password);
  Future<Either<Failure, AuthUser>> register(String email, String password, String name);
  Future<Either<Failure, Unit>> logout();
  Future<Either<Failure, AuthUser>> getCurrentUser();
}
```

**UseCases** (domain/usecases/):

login_user.dart:
```dart
import 'package:fpdart/fpdart.dart';
import '../../core/errors/failures.dart';
import '../entities/auth_user.dart';
import '../repositories/auth_repository.dart';

class LoginUser {
  final AuthRepository repository;

  LoginUser(this.repository);

  Future<Either<Failure, AuthUser>> call(String email, String password) {
    return repository.login(email, password);
  }
}
```

### Data Layer

**Model** (data/models/auth_user_model.dart):
```dart
import 'package:freezed_annotation/freezed_annotation.dart';
import '../../domain/entities/auth_user.dart';

part 'auth_user_model.freezed.dart';
part 'auth_user_model.g.dart';

@freezed
sealed class AuthUserModel with _$AuthUserModel {
  const AuthUserModel._();

  const factory AuthUserModel({
    required String id,
    required String email,
    required String name,
    String? avatarUrl,
    required String createdAt,
  }) = _AuthUserModel;

  factory AuthUserModel.fromJson(Map<String, dynamic> json) =>
      _$AuthUserModelFromJson(json);

  AuthUser toEntity() => AuthUser(
    id: id,
    email: email,
    name: name,
    avatarUrl: avatarUrl,
    createdAt: DateTime.parse(createdAt),
  );
}
```

**DataSource Interface** (data/datasources/auth_datasource.dart):
```dart
import '../models/auth_user_model.dart';

/// Abstract data source - implement with your choice of HTTP client, database, etc.
abstract class AuthDataSource {
  Future<AuthUserModel> login(String email, String password);
  Future<AuthUserModel> register(String email, String password, String name);
  Future<void> logout();
  Future<String?> getToken();
}

/// Example implementation (replace with your actual data source)
class AuthRemoteDataSource implements AuthDataSource {
  // TODO: Implement with Dio, http, GraphQL, Firebase, etc.
  // This is a placeholder - you must provide the actual implementation
  @override
  Future<AuthUserModel> login(String email, String password) {
    throw UnimplementedError('Implement login with your data source');
  }

  @override
  Future<AuthUserModel> register(String email, String password, String name) {
    throw UnimplementedError('Implement register with your data source');
  }

  @override
  Future<void> logout() {
    throw UnimplementedError('Implement logout with your data source');
  }

  @override
  Future<String?> getToken() {
    throw UnimplementedError('Implement getToken with your data source');
  }
}
```

**Repository Implementation** (data/repositories/auth_repository_impl.dart):
```dart
import 'package:fpdart/fpdart.dart';
import '../../core/errors/failures.dart';
import '../../domain/entities/auth_user.dart';
import '../../domain/repositories/auth_repository.dart';
import '../datasources/auth_datasource.dart';

class AuthRepositoryImpl implements AuthRepository {
  final AuthDataSource dataSource;

  AuthRepositoryImpl(this.dataSource);

  @override
  Future<Either<Failure, AuthUser>> login(String email, String password) async {
    try {
      final userModel = await dataSource.login(email, password);
      return Right(userModel.toEntity());
    } on Exception catch (e) {
      return Left(ServerFailure(e.toString()));
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, AuthUser>> register(String email, String password, String name) async {
    try {
      final userModel = await dataSource.register(email, password, name);
      return Right(userModel.toEntity());
    } on Exception catch (e) {
      return Left(ServerFailure(e.toString()));
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, Unit>> logout() async {
    try {
      await dataSource.logout();
      return const Right(unit);
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, AuthUser>> getCurrentUser() async {
    try {
      final token = await dataSource.getToken();
      if (token == null) {
        return const Left(CacheFailure('No token found'));
      }
      // Fetch current user from API - implement as needed
      return const Left(ServerFailure('Not implemented'));
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }
}
```

### Presentation Layer

**Provider** (presentation/providers/auth_provider.dart):
```dart
import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../../domain/entities/auth_user.dart';
import '../../domain/usecases/login_user.dart';
import '../../data/datasources/auth_datasource.dart';
import '../../data/repositories/auth_repository_impl.dart';

part 'auth_provider.g.dart';

// Define your data source provider - implement with your choice
@riverpod
AuthDataSource authDataSource(Ref ref) {
  // Return your implementation: AuthRemoteDataSource(), AuthLocalDataSource(), etc.
  throw UnimplementedError('Provide your AuthDataSource implementation');
}

@riverpod
AuthRepositoryImpl authRepository(Ref ref) {
  return AuthRepositoryImpl(ref.watch(authDataSourceProvider));
}

@riverpod
LoginUser loginUser(Ref ref) {
  return LoginUser(ref.watch(authRepositoryProvider));
}

@riverpod
class AuthNotifier extends _$AuthNotifier {
  @override
  FutureOr<AuthUser?> build() => null;

  Future<void> login(String email, String password) async {
    state = const AsyncLoading();
    final result = await ref.read(loginUserProvider)(email, password);

    state = result.fold(
      (failure) => AsyncError(failure, StackTrace.current),
      (user) => AsyncData(user),
    );
  }

  Future<void> logout() async {
    final result = await ref.read(authRepositoryProvider).logout();
    result.fold(
      (failure) => state = AsyncError(failure, StackTrace.current),
      (_) => state = const AsyncData(null),
    );
  }
}
```

**Login Screen** (presentation/screens/login_screen.dart):
```dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../providers/auth_provider.dart';
import '../../core/errors/failures.dart';

class LoginScreen extends ConsumerStatefulWidget {
  const LoginScreen({super.key});

  @override
  ConsumerState<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends ConsumerState<LoginScreen> {
  final _formKey = GlobalKey<FormState>();
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();

  @override
  void dispose() {
    _emailController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  void _handleLogin() async {
    if (_formKey.currentState!.validate()) {
      await ref.read(authNotifierProvider.notifier).login(
        _emailController.text,
        _passwordController.text,
      );

      final authState = ref.read(authNotifierProvider);

      if (authState.hasValue && authState.value != null) {
        if (mounted) {
          context.go('/home');
        }
      } else if (authState.hasError) {
        if (mounted) {
          final failure = authState.error as Failure;
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text(failure.toString())),
          );
        }
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final authState = ref.watch(authNotifierProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('Login')),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Form(
          key: _formKey,
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              TextFormField(
                controller: _emailController,
                decoration: const InputDecoration(
                  labelText: 'Email',
                  border: OutlineInputBorder(),
                ),
                keyboardType: TextInputType.emailAddress,
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Please enter email';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _passwordController,
                decoration: const InputDecoration(
                  labelText: 'Password',
                  border: OutlineInputBorder(),
                ),
                obscureText: true,
                validator: (value) {
                  if (value == null || value.isEmpty) {
                    return 'Please enter password';
                  }
                  return null;
                },
              ),
              const SizedBox(height: 24),
              SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: authState.isLoading ? null : _handleLogin,
                  child: authState.isLoading
                      ? const SizedBox(
                          height: 20,
                          width: 20,
                          child: CircularProgressIndicator(strokeWidth: 2),
                        )
                      : const Text('Login'),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
```

## Example 2: Product List with Pagination

### Domain Layer

**Entity** (domain/entities/product.dart):
```dart
import 'package:freezed_annotation/freezed_annotation.dart';

part 'product.freezed.dart';

@freezed
sealed class Product with _$Product {
  const factory Product({
    required String id,
    required String name,
    required String description,
    required double price,
    String? imageUrl,
    required int stock,
  }) = _Product;
}

@freezed
sealed class ProductList with _$ProductList {
  const factory ProductList({
    required List<Product> items,
    required int totalCount,
    required int page,
    required bool hasMore,
  }) = _ProductList;
}
```

**Repository** (domain/repositories/product_repository.dart):
```dart
import 'package:fpdart/fpdart.dart';
import '../../core/errors/failures.dart';
import '../entities/product.dart';

abstract class ProductRepository {
  Future<Either<Failure, ProductList>> getProducts({
    required int page,
    required int limit,
  });
  Future<Either<Failure, Product>> getProduct(String id);
}
```

### Data Layer

**DataSource Interface** (data/datasources/product_datasource.dart):
```dart
import '../models/product_model.dart';

/// Abstract data source - implement with your choice of HTTP client, database, etc.
abstract class ProductDataSource {
  Future<ProductListModel> getProducts({required int page, required int limit});
  Future<ProductModel> getProduct(String id);
}

/// Example implementation (replace with your actual data source)
class ProductRemoteDataSource implements ProductDataSource {
  // TODO: Implement with Dio, http, GraphQL, etc.
  @override
  Future<ProductListModel> getProducts({required int page, required int limit}) {
    throw UnimplementedError('Implement getProducts with your data source');
  }

  @override
  Future<ProductModel> getProduct(String id) {
    throw UnimplementedError('Implement getProduct with your data source');
  }
}
```

**Model** (data/models/product_model.dart):
```dart
import 'package:freezed_annotation/freezed_annotation.dart';
import '../../domain/entities/product.dart';

part 'product_model.freezed.dart';
part 'product_model.g.dart';

@freezed
sealed class ProductModel with _$ProductModel {
  const factory ProductModel({
    required String id,
    required String name,
    required String description,
    required double price,
    String? imageUrl,
    required int stock,
  }) = _ProductModel;

  factory ProductModel.fromJson(Map<String, dynamic> json) =>
      _$ProductModelFromJson(json);

  Product toEntity() => Product(
    id: id,
    name: name,
    description: description,
    price: price,
    imageUrl: imageUrl,
    stock: stock,
  );
}

@freezed
sealed class ProductListModel with _$ProductListModel {
  const factory ProductListModel({
    required List<ProductModel> items,
    required int totalCount,
    required int page,
    required bool hasMore,
  }) = _ProductListModel;
}
```

**Repository Implementation** (data/repositories/product_repository_impl.dart):
```dart
import 'package:fpdart/fpdart.dart';
import '../../core/errors/failures.dart';
import '../../domain/entities/product.dart';
import '../../domain/repositories/product_repository.dart';
import '../datasources/product_datasource.dart';

class ProductRepositoryImpl implements ProductRepository {
  final ProductDataSource dataSource;

  ProductRepositoryImpl(this.dataSource);

  @override
  Future<Either<Failure, ProductList>> getProducts({
    required int page,
    required int limit,
  }) async {
    try {
      final productListModel = await dataSource.getProducts(page: page, limit: limit);
      return Right(ProductList(
        items: productListModel.items.map((m) => m.toEntity()).toList(),
        totalCount: productListModel.totalCount,
        page: productListModel.page,
        hasMore: productListModel.hasMore,
      ));
    } on Exception catch (e) {
      return Left(ServerFailure(e.toString()));
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }

  @override
  Future<Either<Failure, Product>> getProduct(String id) async {
    try {
      final productModel = await dataSource.getProduct(id);
      return Right(productModel.toEntity());
    } on Exception catch (e) {
      return Left(ServerFailure(e.toString()));
    } catch (e) {
      return Left(UnexpectedFailure(e.toString()));
    }
  }
}
```

### Presentation Layer with Infinite Scroll

**Provider** (presentation/providers/product_provider.dart):
```dart
import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../../domain/entities/product.dart';
import '../../data/datasources/product_datasource.dart';
import '../../data/repositories/product_repository_impl.dart';

part 'product_provider.g.dart';

@riverpod
ProductDataSource productDataSource(Ref ref) {
  // Return your implementation: ProductRemoteDataSource(), ProductLocalDataSource(), etc.
  throw UnimplementedError('Provide your ProductDataSource implementation');
}

@riverpod
ProductRepositoryImpl productRepository(Ref ref) {
  return ProductRepositoryImpl(ref.watch(productDataSourceProvider));
}

@riverpod
class ProductListNotifier extends _$ProductListNotifier {
  int _currentPage = 1;
  final int _limit = 20;

  @override
  FutureOr<List<Product>> build() async {
    return _fetchProducts();
  }

  Future<List<Product>> _fetchProducts() async {
    final result = await ref.read(productRepositoryProvider).getProducts(
      page: _currentPage,
      limit: _limit,
    );

    return result.fold(
      (failure) => throw failure,
      (productList) => productList.items,
    );
  }

  Future<void> loadMore() async {
    if (state.isLoading) return;

    _currentPage++;

    final result = await ref.read(productRepositoryProvider).getProducts(
      page: _currentPage,
      limit: _limit,
    );

    result.fold(
      (failure) {
        _currentPage--; // Revert on error
        state = AsyncError(failure, StackTrace.current);
      },
      (productList) {
        final currentProducts = state.value ?? [];
        state = AsyncData([...currentProducts, ...productList.items]);
      },
    );
  }

  Future<void> refresh() async {
    _currentPage = 1;
    state = const AsyncLoading();
    state = await AsyncValue.guard(() => _fetchProducts());
  }
}
```

**Screen with Infinite Scroll** (presentation/screens/product_list_screen.dart):
```dart
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../providers/product_provider.dart';

class ProductListScreen extends ConsumerStatefulWidget {
  const ProductListScreen({super.key});

  @override
  ConsumerState<ProductListScreen> createState() => _ProductListScreenState();
}

class _ProductListScreenState extends ConsumerState<ProductListScreen> {
  final _scrollController = ScrollController();

  @override
  void initState() {
    super.initState();
    _scrollController.addListener(_onScroll);
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  void _onScroll() {
    if (_scrollController.position.pixels >=
        _scrollController.position.maxScrollExtent * 0.9) {
      ref.read(productListNotifierProvider.notifier).loadMore();
    }
  }

  @override
  Widget build(BuildContext context) {
    final productsState = ref.watch(productListNotifierProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('Products')),
      body: productsState.when(
        data: (products) {
          if (products.isEmpty) {
            return const Center(child: Text('No products found'));
          }

          return RefreshIndicator(
            onRefresh: () => ref.read(productListNotifierProvider.notifier).refresh(),
            child: ListView.builder(
              controller: _scrollController,
              itemCount: products.length + 1,
              itemBuilder: (context, index) {
                if (index == products.length) {
                  return const Padding(
                    padding: EdgeInsets.all(16.0),
                    child: Center(child: CircularProgressIndicator()),
                  );
                }

                final product = products[index];
                return ListTile(
                  leading: product.imageUrl != null
                      ? Image.network(product.imageUrl!, width: 50, height: 50)
                      : const Icon(Icons.image),
                  title: Text(product.name),
                  subtitle: Text('\$${product.price.toStringAsFixed(2)}'),
                  trailing: Text('Stock: ${product.stock}'),
                );
              },
            ),
          );
        },
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, stack) => Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Text('Error: $error'),
              ElevatedButton(
                onPressed: () => ref.invalidate(productListNotifierProvider),
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