# Advanced Riverpod Patterns

This document covers advanced Riverpod patterns for complex state management scenarios.

## Pattern 1: Family Providers for Dynamic Parameters

Use family providers when you need providers with different parameters:

```dart
@riverpod
class ProductDetails extends _$ProductDetails {
  @override
  FutureOr<Product> build(String productId) async {
    final repository = ref.watch(productRepositoryProvider);
    final result = await repository.getProduct(productId);
    
    return result.fold(
      (failure) => throw failure,
      (product) => product,
    );
  }
  
  Future<void> refresh() async {
    state = const AsyncLoading();
    state = await AsyncValue.guard(() async {
      final repository = ref.watch(productRepositoryProvider);
      final result = await repository.getProduct(productId);
      return result.fold(
        (failure) => throw failure,
        (product) => product,
      );
    });
  }
}

// Usage in widget
class ProductDetailScreen extends ConsumerWidget {
  final String productId;
  
  const ProductDetailScreen({required this.productId, super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final productState = ref.watch(productDetailsProvider(productId));
    
    return productState.when(
      data: (product) => Text(product.name),
      loading: () => CircularProgressIndicator(),
      error: (err, stack) => Text('Error: $err'),
    );
  }
}
```

## Pattern 2: Combining Multiple Providers

Combine data from multiple sources:

```dart
@riverpod
Future<UserProfile> userProfile(Ref ref, String userId) async {
  // Wait for multiple async operations
  final user = await ref.watch(userProvider(userId).future);
  final posts = await ref.watch(userPostsProvider(userId).future);
  final followers = await ref.watch(userFollowersProvider(userId).future);

  return UserProfile(
    user: user,
    posts: posts,
    followerCount: followers.length,
  );
}

// Alternative: Using AsyncValue.guard with multiple providers
@riverpod
class CombinedData extends _$CombinedData {
  @override
  FutureOr<Dashboard> build() async {
    return AsyncValue.guard(() async {
      final stats = await ref.watch(statsProvider.future);
      final recentActivity = await ref.watch(recentActivityProvider.future);
      final notifications = await ref.watch(notificationsProvider.future);
      
      return Dashboard(
        stats: stats,
        recentActivity: recentActivity,
        notifications: notifications,
      );
    });
  }
}
```

## Pattern 3: Dependent Providers Chain

Create a chain of dependent providers:

```dart
// 1. Auth token provider
@riverpod
Future<String?> authToken(Ref ref) async {
  final storage = ref.watch(secureStorageProvider);
  return await storage.read(key: 'auth_token');
}

// 2. Authorized data source (depends on auth token)
@riverpod
UserDataSource userDataSource(Ref ref) {
  // Your data source implementation that uses auth token
  // Replace with your actual implementation
  throw UnimplementedError('Provide your UserDataSource implementation');
}

// 3. Repository (depends on data source)
@riverpod
UserRepository userRepository(Ref ref) {
  return UserRepositoryImpl(ref.watch(userDataSourceProvider));
}
```

## Pattern 4: Listening to Provider Changes

React to provider state changes:

```dart
class MyScreen extends ConsumerStatefulWidget {
  const MyScreen({super.key});

  @override
  ConsumerState<MyScreen> createState() => _MyScreenState();
}

class _MyScreenState extends ConsumerState<MyScreen> {
  @override
  void initState() {
    super.initState();
    
    // Listen to auth state changes
    ref.listenManual(authNotifierProvider, (previous, next) {
      next.whenData((user) {
        if (user == null) {
          // User logged out, navigate to login
          context.go('/login');
        }
      });
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(child: Text('Home Screen')),
    );
  }
}

// Alternative: Using ref.listen in build method
@override
Widget build(BuildContext context, WidgetRef ref) {
  ref.listen(cartProvider, (previous, next) {
    next.whenData((cart) {
      if (cart.items.isEmpty) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Cart is empty')),
        );
      }
    });
  });
  
  return Scaffold(...);
}
```

## Pattern 5: Optimistic Updates

Update UI immediately, then sync with server:

```dart
@riverpod
class TodoList extends _$TodoList {
  @override
  FutureOr<List<Todo>> build() async {
    final repository = ref.watch(todoRepositoryProvider);
    final result = await repository.getTodos();
    return result.fold(
      (failure) => throw failure,
      (todos) => todos,
    );
  }
  
  Future<void> addTodo(Todo todo) async {
    // Optimistic update
    final currentTodos = state.value ?? [];
    state = AsyncData([...currentTodos, todo]);
    
    // Sync with server
    final repository = ref.read(todoRepositoryProvider);
    final result = await repository.createTodo(todo);
    
    result.fold(
      (failure) {
        // Revert on failure
        state = AsyncData(currentTodos);
        // Show error
      },
      (createdTodo) {
        // Update with server response
        final updatedTodos = currentTodos.map((t) => 
          t.id == todo.id ? createdTodo : t
        ).toList();
        state = AsyncData(updatedTodos);
      },
    );
  }
  
  Future<void> toggleTodo(String id) async {
    final currentTodos = state.value ?? [];
    
    // Optimistic update
    final optimisticTodos = currentTodos.map((todo) {
      if (todo.id == id) {
        return todo.copyWith(completed: !todo.completed);
      }
      return todo;
    }).toList();
    state = AsyncData(optimisticTodos);
    
    // Sync with server
    final repository = ref.read(todoRepositoryProvider);
    final result = await repository.toggleTodo(id);
    
    result.fold(
      (failure) => state = AsyncData(currentTodos), // Revert
      (_) {}, // Keep optimistic update
    );
  }
}
```

## Pattern 6: Debouncing and Throttling

Debounce search queries:

```dart
import 'dart:async';

@riverpod
class SearchNotifier extends _$SearchNotifier {
  Timer? _debounceTimer;
  
  @override
  FutureOr<List<Product>> build() => [];
  
  void search(String query) {
    _debounceTimer?.cancel();
    
    if (query.isEmpty) {
      state = const AsyncData([]);
      return;
    }
    
    _debounceTimer = Timer(const Duration(milliseconds: 500), () async {
      state = const AsyncLoading();
      
      final repository = ref.read(productRepositoryProvider);
      final result = await repository.searchProducts(query);
      
      state = result.fold(
        (failure) => AsyncError(failure, StackTrace.current),
        (products) => AsyncData(products),
      );
    });
  }
  
  @override
  void dispose() {
    _debounceTimer?.cancel();
    super.dispose();
  }
}

// Usage in widget
class SearchScreen extends ConsumerWidget {
  const SearchScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final searchResults = ref.watch(searchNotifierProvider);
    
    return Column(
      children: [
        TextField(
          onChanged: (value) => 
            ref.read(searchNotifierProvider.notifier).search(value),
          decoration: const InputDecoration(hintText: 'Search...'),
        ),
        Expanded(
          child: searchResults.when(
            data: (products) => ListView.builder(
              itemCount: products.length,
              itemBuilder: (context, index) => ListTile(
                title: Text(products[index].name),
              ),
            ),
            loading: () => const Center(child: CircularProgressIndicator()),
            error: (err, stack) => Center(child: Text('Error: $err')),
          ),
        ),
      ],
    );
  }
}
```

## Pattern 7: Caching with Auto-Invalidation

Cache data with automatic refresh:

```dart
@riverpod
Future<WeatherData> weather(
  Ref ref,
  String city,
) async {
  // Auto-dispose after 5 minutes of inactivity
  final link = ref.keepAlive();
  final timer = Timer(const Duration(minutes: 5), () {
    link.close();
  });
  ref.onDispose(() => timer.cancel());

  // Fetch data
  final repository = ref.watch(weatherRepositoryProvider);
  final result = await repository.getWeather(city);

  return result.fold(
    (failure) => throw failure,
    (weather) => weather,
  );
}

// Force refresh every 10 minutes
@riverpod
class LiveWeather extends _$LiveWeather {
  Timer? _refreshTimer;
  
  @override
  FutureOr<WeatherData> build(String city) async {
    _startAutoRefresh();
    return _fetchWeather();
  }
  
  void _startAutoRefresh() {
    _refreshTimer?.cancel();
    _refreshTimer = Timer.periodic(
      const Duration(minutes: 10),
      (_) => refresh(),
    );
  }
  
  Future<WeatherData> _fetchWeather() async {
    final repository = ref.watch(weatherRepositoryProvider);
    final result = await repository.getWeather(city);
    return result.fold(
      (failure) => throw failure,
      (weather) => weather,
    );
  }
  
  Future<void> refresh() async {
    state = const AsyncLoading();
    state = await AsyncValue.guard(() => _fetchWeather());
  }
  
  @override
  void dispose() {
    _refreshTimer?.cancel();
    super.dispose();
  }
}
```

## Pattern 8: Global State Management

Manage app-wide state:

```dart
@riverpod
class AppSettings extends _$AppSettings {
  @override
  AppSettingsState build() {
    _loadSettings();
    return const AppSettingsState();
  }
  
  Future<void> _loadSettings() async {
    final prefs = await SharedPreferences.getInstance();
    state = AppSettingsState(
      isDarkMode: prefs.getBool('dark_mode') ?? false,
      language: prefs.getString('language') ?? 'en',
      notificationsEnabled: prefs.getBool('notifications') ?? true,
    );
  }
  
  Future<void> setDarkMode(bool value) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setBool('dark_mode', value);
    state = state.copyWith(isDarkMode: value);
  }
  
  Future<void> setLanguage(String language) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('language', language);
    state = state.copyWith(language: language);
  }
}

@freezed
class AppSettingsState with _$AppSettingsState {
  const factory AppSettingsState({
    @Default(false) bool isDarkMode,
    @Default('en') String language,
    @Default(true) bool notificationsEnabled,
  }) = _AppSettingsState;
}

// Usage in MaterialApp
class MyApp extends ConsumerWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final settings = ref.watch(appSettingsProvider);
    
    return MaterialApp(
      theme: settings.isDarkMode ? ThemeData.dark() : ThemeData.light(),
      locale: Locale(settings.language),
      home: const HomeScreen(),
    );
  }
}
```

## Pattern 9: Error Recovery

Handle and recover from errors gracefully:

```dart
@riverpod
class ResilientData extends _$ResilientData {
  int _retryCount = 0;
  static const maxRetries = 3;
  
  @override
  FutureOr<DataModel> build() async {
    return _fetchWithRetry();
  }
  
  Future<DataModel> _fetchWithRetry() async {
    try {
      final repository = ref.watch(dataRepositoryProvider);
      final result = await repository.getData();
      
      return result.fold(
        (failure) async {
          if (_retryCount < maxRetries) {
            _retryCount++;
            // Exponential backoff
            await Future.delayed(Duration(seconds: _retryCount * 2));
            return _fetchWithRetry();
          }
          throw failure;
        },
        (data) {
          _retryCount = 0; // Reset on success
          return data;
        },
      );
    } catch (e) {
      // Try to return cached data if available
      final cache = await _getCachedData();
      if (cache != null) return cache;
      rethrow;
    }
  }
  
  Future<DataModel?> _getCachedData() async {
    // Implement cache retrieval
    return null;
  }
  
  Future<void> refresh() async {
    _retryCount = 0;
    state = const AsyncLoading();
    state = await AsyncValue.guard(() => _fetchWithRetry());
  }
}
```

## Pattern 10: Pagination State Management

Advanced pagination with load more:

```dart
@freezed
class PaginatedState<T> with _$PaginatedState<T> {
  const factory PaginatedState({
    @Default([]) List<T> items,
    @Default(false) bool isLoadingMore,
    @Default(false) bool hasReachedEnd,
    @Default(1) int currentPage,
  }) = _PaginatedState;
}

@riverpod
class PaginatedProducts extends _$PaginatedProducts {
  static const _pageSize = 20;
  
  @override
  FutureOr<PaginatedState<Product>> build() async {
    final products = await _fetchPage(1);
    return PaginatedState(
      items: products,
      currentPage: 1,
      hasReachedEnd: products.length < _pageSize,
    );
  }
  
  Future<List<Product>> _fetchPage(int page) async {
    final repository = ref.watch(productRepositoryProvider);
    final result = await repository.getProducts(
      page: page,
      limit: _pageSize,
    );
    
    return result.fold(
      (failure) => throw failure,
      (productList) => productList.items,
    );
  }
  
  Future<void> loadMore() async {
    final currentState = state.value;
    if (currentState == null || 
        currentState.isLoadingMore || 
        currentState.hasReachedEnd) {
      return;
    }
    
    state = AsyncData(currentState.copyWith(isLoadingMore: true));
    
    try {
      final nextPage = currentState.currentPage + 1;
      final newProducts = await _fetchPage(nextPage);
      
      state = AsyncData(PaginatedState(
        items: [...currentState.items, ...newProducts],
        currentPage: nextPage,
        hasReachedEnd: newProducts.length < _pageSize,
        isLoadingMore: false,
      ));
    } catch (e) {
      state = AsyncData(currentState.copyWith(isLoadingMore: false));
      rethrow;
    }
  }
  
  Future<void> refresh() async {
    state = const AsyncLoading();
    final products = await _fetchPage(1);
    state = AsyncData(PaginatedState(
      items: products,
      currentPage: 1,
      hasReachedEnd: products.length < _pageSize,
    ));
  }
}
```
