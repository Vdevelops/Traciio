# Presentation Layer Guide

Complete guide for implementing the Presentation Layer in Clean Architecture with Riverpod.

## Overview

The Presentation Layer consists of:
- **Providers**: Riverpod providers for dependency injection and state management
- **Screens**: Full-page widgets
- **Widgets**: Reusable UI components

## Provider Patterns

### Basic Provider Setup

```dart
import 'package:riverpod_annotation/riverpod_annotation.dart';
import '../../domain/usecases/get_user.dart';
import '../../data/datasources/user_api_service.dart';
import '../../data/repositories/user_repository_impl.dart';
import '../../../../core/network/dio_provider.dart';

part 'user_provider.g.dart';

// Dio provider (from core)
@riverpod
Dio dio(Ref ref) {
  // Defined in core/network/dio_provider.dart
  throw UnimplementedError('Use dioProvider from core');
}

// API Service provider
@riverpod
UserApiService userApiService(Ref ref) {
  return UserApiService(ref.watch(dioProvider));
}

// Repository provider
@riverpod
UserRepositoryImpl userRepository(Ref ref) {
  return UserRepositoryImpl(ref.watch(userApiServiceProvider));
}

// UseCase provider
@riverpod
GetUser getUser(Ref ref) {
  return GetUser(ref.watch(userRepositoryProvider));
}
```

### Notifier Pattern (State Management)

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

  Future<void> refreshUser() async {
    final currentUser = state.value;
    if (currentUser == null) return;

    final result = await ref.read(getUserProvider)(currentUser.id);

    state = result.fold(
      (failure) => AsyncError(failure, StackTrace.current),
      (user) => AsyncData(user),
    );
  }
}
```

### List Notifier Pattern

```dart
@riverpod
class UsersNotifier extends _$UsersNotifier {
  @override
  FutureOr<List<User>> build() async {
    // Load initial data
    final result = await ref.read(getUsersProvider)());
    return result.fold(
      (failure) => throw Exception(failure),
      (users) => users,
    );
  }

  Future<void> refresh() async {
    final result = await ref.read(getUsersProvider)());
    state = result.fold(
      (failure) => AsyncError(failure, StackTrace.current),
      (users) => AsyncData(users),
    );
  }

  Future<void> addUser(CreateUserParams params) async {
    final result = await ref.read(createUserProvider)(params);

    result.fold(
      (failure) => null, // Handle error
      (user) {
        // Add to current list
        state.whenData((users) => [...users, user]);
      },
    );
  }

  Future<void> removeUser(String id) async {
    final result = await ref.read(deleteUserProvider)(id);

    result.fold(
      (failure) => null, // Handle error
      (_) {
        // Remove from current list
        state.whenData((users) => users.where((u) => u.id != id).toList());
      },
    );
  }
}
```

### Async Notifier with Loading States

```dart
@riverpod
class UserDetailNotifier extends _$UserDetailNotifier {
  @override
  FutureOr<User?> build() => null;

  Future<void> loadUser(String id) async {
    state = const AsyncLoading();

    // Show loading indicator
    ref.read(isLoadingProvider.notifier).state = true;

    final result = await ref.read(getUserProvider)(id);

    ref.read(isLoadingProvider.notifier).state = false;

    state = result.fold(
      (failure) {
        // Show error message
        ref.read(errorMessageProvider.notifier).state = failure.message;
        return AsyncError(failure, StackTrace.current);
      },
      (user) => AsyncData(user),
    );
  }
}

// Simple state providers
@riverpod
class IsLoading extends _$IsLoading {
  @override
  bool build() => false;
}

@riverpod
class ErrorMessage extends _$ErrorMessage {
  @override
  String? build() => null;
}
```

### Filtered/Computed State

```dart
@riverpod
List<User> activeUsers(Ref ref) {
  final usersAsync = ref.watch(usersNotifierProvider);

  return usersAsync.when(
    data: (users) => users.where((u) => u.isActive).toList(),
    loading: () => [],
    error: (_, __) => [],
  );
}

@riverpod
User? findUserById(Ref ref, String userId) {
  final usersAsync = ref.watch(usersNotifierProvider);

  return usersAsync.when(
    data: (users) => users.tryFind((u) => u.id == userId),
    loading: () => null,
    error: (_, __) => null,
  );
}
```

## Screen Patterns

### Basic List Screen

```dart
class UserListScreen extends ConsumerWidget {
  const UserListScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final usersAsync = ref.watch(usersNotifierProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Users'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () => ref.read(usersNotifierProvider.notifier).refresh(),
          ),
        ],
      ),
      body: usersAsync.when(
        data: (users) {
          if (users.isEmpty) {
            return const EmptyStateWidget(
              message: 'No users found',
              icon: Icons.people_outline,
            );
          }
          return ListView.builder(
            itemCount: users.length,
            itemBuilder: (context, index) {
              final user = users[index];
              return UserListTile(user: user);
            },
          );
        },
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, stack) => ErrorStateWidget(
          message: error.toString(),
          onRetry: () => ref.read(usersNotifierProvider.notifier).refresh(),
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _navigateToAddUser(context),
        child: const Icon(Icons.add),
      ),
    );
  }
}
```

### Detail Screen with Auto-refresh

```dart
class UserDetailScreen extends ConsumerStatefulWidget {
  final String userId;

  const UserDetailScreen({required this.userId, super.key});

  @override
  ConsumerState<UserDetailScreen> createState() => _UserDetailScreenState();
}

class _UserDetailScreenState extends ConsumerState<UserDetailScreen> {
  @override
  void initState() {
    super.initState();
    // Load user data on init
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(userNotifierProvider.notifier).fetchUser(widget.userId);
    });
  }

  @override
  Widget build(BuildContext context) {
    final userAsync = ref.watch(userNotifierProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('User Details'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () => ref.read(userNotifierProvider.notifier).refreshUser(),
          ),
          IconButton(
            icon: const Icon(Icons.edit),
            onPressed: () => _navigateToEdit(context, widget.userId),
          ),
        ],
      ),
      body: userAsync.when(
        data: (user) {
          if (user == null) {
            return const Center(child: Text('No user loaded'));
          }
          return UserDetailView(user: user);
        },
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, stack) => ErrorStateWidget(
          message: error.toString(),
          onRetry: () => ref.read(userNotifierProvider.notifier).fetchUser(widget.userId),
        ),
      ),
    );
  }
}
```

### Form Screen with Validation

```dart
class CreateUserScreen extends ConsumerStatefulWidget {
  const CreateUserScreen({super.key});

  @override
  ConsumerState<CreateUserScreen> createState() => _CreateUserScreenState();
}

class _CreateUserScreenState extends ConsumerState<CreateUserScreen> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _emailController = TextEditingController();

  @override
  void dispose() {
    _nameController.dispose();
    _emailController.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;

    final params = CreateUserParams(
      name: _nameController.text,
      email: _emailController.text,
    );

    final result = await ref.read(createUserProvider)(params);

    if (!mounted) return;

    result.fold(
      (failure) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error: ${failure.message}')),
        );
      },
      (user) {
        // Navigate back or show success
        Navigator.of(context).pop(user);
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Create User')),
      body: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            TextFormField(
              controller: _nameController,
              decoration: const InputDecoration(labelText: 'Name'),
              validator: (value) => value?.isEmpty ?? true ? 'Name is required' : null,
            ),
            TextFormField(
              controller: _emailController,
              decoration: const InputDecoration(labelText: 'Email'),
              keyboardType: TextInputType.emailAddress,
              validator: (value) {
                if (value?.isEmpty ?? true) return 'Email is required';
                if (!value!.contains('@')) return 'Invalid email';
                return null;
              },
            ),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: _submit,
              child: const Text('Create'),
            ),
          ],
        ),
      ),
    );
  }
}
```

## Reusable Widgets

### List Tile Widget

```dart
class UserListTile extends StatelessWidget {
  final User user;
  final VoidCallback? onTap;

  const UserListTile({
    required this.user,
    this.onTap,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return ListTile(
      leading: CircleAvatar(
        child: Text(user.name[0].toUpperCase()),
      ),
      title: Text(user.name),
      subtitle: Text(user.email),
      trailing: Icon(Icons.chevron_right),
      onTap: onTap,
    );
  }
}
```

### Empty State Widget

```dart
class EmptyStateWidget extends StatelessWidget {
  final String message;
  final IconData icon;

  const EmptyStateWidget({
    required this.message,
    this.icon = Icons.inbox,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(icon, size: 64, color: Colors.grey),
          const SizedBox(height: 16),
          Text(
            message,
            style: Theme.of(context).textTheme.titleMedium?.copyWith(
              color: Colors.grey,
            ),
          ),
        ],
      ),
    );
  }
}
```

### Error State Widget

```dart
class ErrorStateWidget extends StatelessWidget {
  final String message;
  final VoidCallback onRetry;

  const ErrorStateWidget({
    required this.message,
    required this.onRetry,
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Icon(Icons.error_outline, size: 64, color: Colors.red),
          const SizedBox(height: 16),
          Text(message),
          const SizedBox(height: 16),
          ElevatedButton(
            onPressed: onRetry,
            child: const Text('Retry'),
          ),
        ],
      ),
    );
  }
}
```

## Best Practices

1. **Use @riverpod annotations** for all providers
2. **Keep providers in presentation/providers/** directory
3. **One provider file per feature**
4. **Use notifier pattern** for mutable state
5. **Separate UI widgets** into separate files
6. **Handle loading, error, and data states** properly
7. **Dispose controllers** in StatefulWidget
8. **Use ref.watch** for reading, **ref.read** for calling methods
9. **Avoid business logic** in widgets - use use cases
10. **Use const constructors** for widgets when possible