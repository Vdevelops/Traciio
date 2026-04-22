---
trigger: always_on
---

# Flutter Riverpod Expert — 2025 Best Practices

Expert knowledge in Flutter Riverpod state management following 2025 best practices, adapted for CRM Healthcare project conventions.

## Core Principles

1. **StateNotifierProvider** is the standard in this project (not code generation)
2. **One provider per feature** — `features/<feature>/application/<feature>_provider.dart`
3. **Immutable state classes** — `features/<feature>/application/<feature>_state.dart`
4. **Repository pattern** — Separate data layer from state management
5. **Performance first** — Use `select()` to optimize rebuilds

## Provider Pattern (Project Standard)

This project uses **manual StateNotifier** (not `@riverpod` code generation):

```dart
// features/<feature>/application/<feature>_provider.dart
import 'package:flutter_riverpod/flutter_riverpod.dart';

final featureProvider = StateNotifierProvider<FeatureNotifier, FeatureState>((ref) {
  final repository = ref.read(featureRepositoryProvider);
  return FeatureNotifier(repository);
});

class FeatureNotifier extends StateNotifier<FeatureState> {
  final FeatureRepository _repository;

  FeatureNotifier(this._repository) : super(const FeatureState());

  Future<void> loadData() async {
    state = state.copyWith(isLoading: true, errorMessage: null);
    try {
      final data = await _repository.getData();
      state = state.copyWith(isLoading: false, data: data);
    } catch (e) {
      state = state.copyWith(isLoading: false, errorMessage: e.toString());
    }
  }
}
```

## State Class Pattern (Project Standard)

```dart
// features/<feature>/application/<feature>_state.dart
class FeatureState {
  final bool isLoading;
  final String? errorMessage;
  final List<Item> items;

  const FeatureState({
    this.isLoading = false,
    this.errorMessage,
    this.items = const [],
  });

  FeatureState copyWith({
    bool? isLoading,
    String? errorMessage,
    List<Item>? items,
  }) {
    return FeatureState(
      isLoading: isLoading ?? this.isLoading,
      errorMessage: errorMessage,
      items: items ?? this.items,
    );
  }
}
```

## Performance Optimization Patterns

### Use ref.select() for Specific Fields
```dart
// AVOID: Rebuilds on ANY state change
final state = ref.watch(dashboardProvider);

// PREFER: Only rebuilds when specific field changes
final isLoading = ref.watch(dashboardProvider.select((s) => s.isLoading));
final visits = ref.watch(dashboardProvider.select((s) => s.visits));
```

### ref.watch() vs ref.read() vs ref.listen()

- **ref.watch()** — Use in `build()` method for reactive UI updates
- **ref.read()** — Use in event handlers (onPressed, onTap) for one-time reads
- **ref.listen()** — Use for side effects (navigation, snackbars, logging)

```dart
// In build():
final tasks = ref.watch(taskProvider.select((s) => s.tasks));

// In event handler:
onPressed: () => ref.read(taskProvider.notifier).deleteTask(id);

// For side effects:
ref.listen<TaskState>(taskProvider, (prev, next) {
  if (next.errorMessage != null) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(next.errorMessage!)),
    );
  }
});
```

### Avoid Watching in Loops
```dart
// BAD: Causes performance issues
ListView.builder(
  itemBuilder: (context, index) {
    final item = ref.watch(itemProvider(ids[index])); // DON'T!
    return ListTile(...);
  },
);

// GOOD: Separate widget for each item
class ItemTile extends ConsumerWidget {
  final String itemId;
  const ItemTile({required this.itemId});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final item = ref.watch(itemProvider(itemId));
    return ListTile(title: Text(item.name));
  }
}
```

## Repository Pattern (Project Standard)

### Offline-First Repository
```dart
class FeatureRepository {
  final ApiClient _apiClient;
  final OfflineStorage _offlineStorage;
  final ConnectivityService _connectivity;

  Future<List<Item>> getItems() async {
    // 1. Try cache first
    final cached = await _offlineStorage.getData<List<Item>>('feature_items');
    if (cached != null) return cached;

    // 2. Try network
    if (await _connectivity.hasConnection) {
      try {
        final response = await _apiClient.get('/api/v1/mobile/feature/items');
        final items = _parseResponse(response);
        await _offlineStorage.saveData('feature_items', items);
        return items;
      } catch (e) {
        // 3. Fallback to expired cache
        return cached ?? [];
      }
    }

    return cached ?? [];
  }
}
```

## Error Handling

### AsyncValue-style with State
```dart
// In provider:
Future<void> loadData() async {
  state = state.copyWith(isLoading: true, errorMessage: null);
  try {
    final data = await _repository.getData();
    state = state.copyWith(isLoading: false, data: data);
  } on DioException catch (e) {
    final message = _extractErrorMessage(e);
    state = state.copyWith(isLoading: false, errorMessage: message);
  } catch (e) {
    state = state.copyWith(isLoading: false, errorMessage: 'Unexpected error');
  }
}

// In UI:
final state = ref.watch(featureProvider);
if (state.isLoading) return const CircularProgressIndicator();
if (state.errorMessage != null) return ErrorWidget(message: state.errorMessage!);
return DataWidget(data: state.data);
```

## Testing

### Provider Testing
```dart
test('loads data successfully', () async {
  final container = ProviderContainer(
    overrides: [
      featureRepositoryProvider.overrideWithValue(MockRepository()),
    ],
  );

  await container.read(featureProvider.notifier).loadData();
  final state = container.read(featureProvider);

  expect(state.isLoading, false);
  expect(state.data.length, 2);
});
```

### Widget Testing
```dart
testWidgets('displays items', (tester) async {
  await tester.pumpWidget(
    ProviderScope(
      overrides: [
        featureRepositoryProvider.overrideWithValue(MockRepository()),
      ],
      child: const MaterialApp(home: FeatureScreen()),
    ),
  );
  await tester.pumpAndSettle();
  expect(find.text('Item 1'), findsOneWidget);
});
```

## Common Anti-Patterns to AVOID

1. **Using `ref.read()` in `build()`** — Won't rebuild when state changes
2. **Mutating state directly** — Always use `state = state.copyWith(...)`
3. **Multiple sources of truth** — Don't mix local `setState` with Riverpod
4. **Not disposing resources** — Use `ref.onDispose()` for cleanup
5. **Watching entire state** — Use `select()` for specific fields
6. **Creating providers inside widgets** — Define providers at top level
