import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/network/pagination.dart';

import '../../../core/cache/list_cache.dart';
import '../../../core/network/api_client.dart';
import '../../../core/network/connectivity_service.dart';
import '../../../core/sync/auto_sync_manager.dart';
import '../data/models/task.dart';
import '../data/task_repository.dart';
import 'task_state.dart';

final taskRepositoryProvider = Provider<TaskRepository>((ref) {
  final connectivity = ref.watch(connectivityServiceProvider);
  return TaskRepository(ApiClient.dio, connectivity);
});

final taskListProvider = NotifierProvider<TaskListNotifier, TaskListState>(
  TaskListNotifier.new,
);

final taskDetailProvider = FutureProvider.family<Task, String>((ref, id) async {
  final repository = ref.read(taskRepositoryProvider);
  // Always fetch fresh data from API when online (cache only used when offline)
  return repository.getTaskById(id);
});

final taskFormProvider = NotifierProvider<TaskFormNotifier, TaskFormState>(
  TaskFormNotifier.new,
);

class TaskListNotifier extends Notifier<TaskListState>
    with WidgetsBindingObserver {
  late final TaskRepository _repository;
  late final ConnectivityService _connectivity;
  final ListCache _cache = ListCache();

  @override
  TaskListState build() {
    _repository = ref.read(taskRepositoryProvider);
    _connectivity = ref.read(connectivityServiceProvider);

    // Listen to app lifecycle changes for auto-sync
    WidgetsBinding.instance.addObserver(this);
    ref.onDispose(() {
      WidgetsBinding.instance.removeObserver(this);
    });

    // Register with AutoSyncManager for centralized sync
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref
          .read(autoSyncManagerProvider.notifier)
          .registerFeature('tasks', () => silentRefresh());
    });

    return const TaskListState();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    super.didChangeAppLifecycleState(state);
    if (state == AppLifecycleState.resumed) {
      // Auto sync when app comes back to foreground
      ref.read(autoSyncManagerProvider.notifier).syncFeature('tasks');
    }
  }

  Future<void> loadTasks({
    int page = 1,
    bool refresh = false,
    String? search,
    String? status,
    String? priority,
    String? type,
    DateTime? dueDateFrom,
    DateTime? dueDateTo,
    bool forceRefresh = false,
  }) async {
    final searchQuery = search ?? state.searchQuery;
    final statusFilter = status ?? state.selectedStatus;
    final priorityFilter = priority ?? state.selectedPriority;
    final typeFilter = type ?? state.selectedType;
    final dueDateFromFilter = dueDateFrom ?? state.dueDateFrom;
    final dueDateToFilter = dueDateTo ?? state.dueDateTo;

    final cacheKey = ListCache.cacheKey(
      'tasks',
      page: page,
      search: searchQuery.isNotEmpty ? searchQuery : null,
      filters: {
        ...?(statusFilter != null ? {'status': statusFilter} : null),
        ...?(priorityFilter != null ? {'priority': priorityFilter} : null),
        ...?(typeFilter != null ? {'type': typeFilter} : null),
        ...?(dueDateFromFilter != null
            ? {'due_date_from': dueDateFromFilter.toIso8601String()}
            : null),
        ...?(dueDateToFilter != null
            ? {'due_date_to': dueDateToFilter.toIso8601String()}
            : null),
      },
    );

    // Try to load from cache first (optimistic UI) - only for first page
    if (!forceRefresh && !refresh && page == 1) {
      final cachedTasks = _cache.get<Task>(
        cacheKey,
        ttl: const Duration(seconds: 60),
        expectedMetadata: {
          if (searchQuery.isNotEmpty) 'search': searchQuery,
          ...?(statusFilter != null ? {'status': statusFilter} : null),
          ...?(priorityFilter != null ? {'priority': priorityFilter} : null),
          ...?(typeFilter != null ? {'type': typeFilter} : null),
          ...?(dueDateFromFilter != null
              ? {'due_date_from': dueDateFromFilter.toIso8601String()}
              : null),
          ...?(dueDateToFilter != null
              ? {'due_date_to': dueDateToFilter.toIso8601String()}
              : null),
        },
      );

      if (cachedTasks != null && cachedTasks.isNotEmpty) {
        // Show cached data immediately
        final cachedMetadata = _cache.getMetadata(cacheKey);
        Pagination? cachedPagination;
        if (cachedMetadata?['pagination'] != null) {
          try {
            cachedPagination = Pagination.fromJson(
              cachedMetadata!['pagination'] as Map<String, dynamic>,
            );
          } catch (e) {
            // Ignore pagination parsing error
          }
        }
        state = state.copyWith(
          tasks: cachedTasks,
          searchQuery: searchQuery,
          selectedStatus: statusFilter,
          selectedPriority: priorityFilter,
          selectedType: typeFilter,
          dueDateFrom: dueDateFromFilter,
          dueDateTo: dueDateToFilter,
          isLoading: false,
          isLoadingMore: false,
          errorMessage: null,
          pagination: cachedPagination,
        );
      }
    }

    // Set loading state
    // If search query changed or refresh, clear tasks first
    final isSearchQueryChanged = searchQuery != state.searchQuery;
    if (refresh || page == 1 || isSearchQueryChanged) {
      state = state.copyWith(
        tasks: const [], // Clear tasks when search query changes or refresh
        isLoading: true,
        isLoadingMore: false,
        errorMessage: null,
        clearTasks: true, // Use clearTasks flag to ensure tasks are cleared
      );
    } else {
      state = state.copyWith(isLoadingMore: true);
    }

    try {
      // Mobile API automatically filters by logged-in user
      final response = await _repository.getTasks(
        page: page,
        perPage: 20,
        search: searchQuery.isNotEmpty ? searchQuery : null,
        status: statusFilter,
        priority: priorityFilter,
        type: typeFilter,
        dueDateFrom: dueDateFromFilter,
        dueDateTo: dueDateToFilter,
        forceRefresh: forceRefresh,
        onBackgroundUpdate: (freshData) {
          // Update UI when fresh data arrives from background fetch
          _updateStateWithFreshData(freshData);
        },
      );

      // Cache the response
      _cache.set(
        cacheKey,
        response.items,
        metadata: {
          'pagination': {
            'page': response.pagination.page,
            'perPage': response.pagination.perPage,
            'total': response.pagination.total,
            'totalPages': response.pagination.totalPages,
          },
          'search': searchQuery,
          ...?(statusFilter != null ? {'status': statusFilter} : null),
          ...?(priorityFilter != null ? {'priority': priorityFilter} : null),
          ...?(typeFilter != null ? {'type': typeFilter} : null),
          ...?(dueDateFromFilter != null
              ? {'due_date_from': dueDateFromFilter.toIso8601String()}
              : null),
          ...?(dueDateToFilter != null
              ? {'due_date_to': dueDateToFilter.toIso8601String()}
              : null),
        },
      );

      if (refresh || page == 1) {
        state = state.copyWith(
          tasks: response.items,
          pagination: response.pagination,
          searchQuery: searchQuery,
          selectedStatus: statusFilter,
          selectedPriority: priorityFilter,
          selectedType: typeFilter,
          dueDateFrom: dueDateFromFilter,
          dueDateTo: dueDateToFilter,
          isLoading: false,
          isLoadingMore: false,
          errorMessage: null,
          isOffline: !_connectivity.isOnline,
        );
      } else {
        state = state.copyWith(
          tasks: [...state.tasks, ...response.items],
          pagination: response.pagination,
          isLoadingMore: false,
          errorMessage: null,
          isOffline: !_connectivity.isOnline,
        );
      }
    } catch (e) {
      // On error, try to use cached data as fallback
      if (page == 1) {
        final cachedTasks = _cache.get<Task>(cacheKey);
        if (cachedTasks != null && cachedTasks.isNotEmpty) {
          state = state.copyWith(
            tasks: cachedTasks,
            isLoading: false,
            isLoadingMore: false,
            errorMessage: null,
            isOffline: !_connectivity.isOnline,
          );
          return;
        }
      }

      state = state.copyWith(
        isLoading: false,
        isLoadingMore: false,
        errorMessage: e.toString().replaceFirst('Exception: ', ''),
        isOffline: !_connectivity.isOnline,
      );
    }
  }

  /// Silent refresh (dipanggil oleh AutoSyncManager)
  Future<void> silentRefresh() async {
    try {
      final response = await _repository.getTasks(
        page: 1,
        perPage: 20,
        search: state.searchQuery.isNotEmpty ? state.searchQuery : null,
        status: state.selectedStatus,
        priority: state.selectedPriority,
        type: state.selectedType,
        dueDateFrom: state.dueDateFrom,
        dueDateTo: state.dueDateTo,
        forceRefresh: true,
      );

      // Update state with fresh data
      state = state.copyWith(
        tasks: response.items,
        pagination: response.pagination,
        errorMessage: null,
      );

      // Cache the response
      final cacheKey = ListCache.cacheKey(
        'tasks',
        page: 1,
        search: state.searchQuery.isNotEmpty ? state.searchQuery : null,
        filters: {
          ...?(state.selectedStatus != null
              ? {'status': state.selectedStatus!}
              : null),
          ...?(state.selectedPriority != null
              ? {'priority': state.selectedPriority!}
              : null),
          ...?(state.selectedType != null
              ? {'type': state.selectedType!}
              : null),
          ...?(state.dueDateFrom != null
              ? {'due_date_from': state.dueDateFrom!.toIso8601String()}
              : null),
          ...?(state.dueDateTo != null
              ? {'due_date_to': state.dueDateTo!.toIso8601String()}
              : null),
        },
      );
      _cache.set(
        cacheKey,
        response.items,
        metadata: {
          'pagination': {
            'page': response.pagination.page,
            'perPage': response.pagination.perPage,
            'total': response.pagination.total,
            'totalPages': response.pagination.totalPages,
          },
          'search': state.searchQuery,
          ...?(state.selectedStatus != null
              ? {'status': state.selectedStatus!}
              : null),
          ...?(state.selectedPriority != null
              ? {'priority': state.selectedPriority!}
              : null),
          ...?(state.selectedType != null
              ? {'type': state.selectedType!}
              : null),
          ...?(state.dueDateFrom != null
              ? {'due_date_from': state.dueDateFrom!.toIso8601String()}
              : null),
          ...?(state.dueDateTo != null
              ? {'due_date_to': state.dueDateTo!.toIso8601String()}
              : null),
        },
      );
    } catch (e) {
      // Silent fail
    }
  }

  /// Update state with fresh data from background fetch
  void _updateStateWithFreshData(TaskListResponse freshData) {
    // Create new list instance to ensure Riverpod detects the change
    final newTasks = List<Task>.from(freshData.items);

    state = TaskListState(
      tasks: newTasks,
      pagination: freshData.pagination,
      searchQuery: state.searchQuery,
      selectedStatus: state.selectedStatus,
      selectedPriority: state.selectedPriority,
      selectedType: state.selectedType,
      dueDateFrom: state.dueDateFrom,
      dueDateTo: state.dueDateTo,
      isLoading: state.isLoading,
      isLoadingMore: state.isLoadingMore,
      errorMessage: null,
    );

    // Cache the response
    final cacheKey = ListCache.cacheKey(
      'tasks',
      page: 1,
      search: state.searchQuery.isNotEmpty ? state.searchQuery : null,
      filters: {
        ...?(state.selectedStatus != null
            ? {'status': state.selectedStatus!}
            : null),
        ...?(state.selectedPriority != null
            ? {'priority': state.selectedPriority!}
            : null),
        ...?(state.selectedType != null ? {'type': state.selectedType!} : null),
        ...?(state.dueDateFrom != null
            ? {'due_date_from': state.dueDateFrom!.toIso8601String()}
            : null),
        ...?(state.dueDateTo != null
            ? {'due_date_to': state.dueDateTo!.toIso8601String()}
            : null),
      },
    );
    _cache.set(
      cacheKey,
      freshData.items,
      metadata: {
        'pagination': {
          'page': freshData.pagination.page,
          'perPage': freshData.pagination.perPage,
          'total': freshData.pagination.total,
          'totalPages': freshData.pagination.totalPages,
        },
        'search': state.searchQuery,
        ...?(state.selectedStatus != null
            ? {'status': state.selectedStatus!}
            : null),
        ...?(state.selectedPriority != null
            ? {'priority': state.selectedPriority!}
            : null),
        ...?(state.selectedType != null ? {'type': state.selectedType!} : null),
        ...?(state.dueDateFrom != null
            ? {'due_date_from': state.dueDateFrom!.toIso8601String()}
            : null),
        ...?(state.dueDateTo != null
            ? {'due_date_to': state.dueDateTo!.toIso8601String()}
            : null),
      },
    );
  }

  Future<void> refresh() async {
    // Clear cache for tasks
    _cache.clearPrefix('list:tasks');
    await loadTasks(page: 1, refresh: true, forceRefresh: true);
  }

  Future<void> loadMore() async {
    if (state.isLoading || state.isLoadingMore) return;
    final pagination = state.pagination;
    if (pagination == null || !pagination.hasNextPage) return;

    await loadTasks(page: pagination.page + 1);
  }

  void updateSearchQuery(String query) {
    // Clear tasks when search query changes
    state = state.copyWith(
      searchQuery: query,
      tasks: query != state.searchQuery ? const [] : state.tasks,
    );
  }

  void updateStatusFilter(String? status) {
    state = state.copyWith(selectedStatus: status);
    _cache.clearPrefix('list:tasks');
    loadTasks(page: 1, refresh: true, status: status, forceRefresh: true);
  }

  void updatePriorityFilter(String? priority) {
    state = state.copyWith(selectedPriority: priority);
    _cache.clearPrefix('list:tasks');
    loadTasks(page: 1, refresh: true, priority: priority, forceRefresh: true);
  }

  void updateTypeFilter(String? type) {
    state = state.copyWith(selectedType: type);
    _cache.clearPrefix('list:tasks');
    loadTasks(page: 1, refresh: true, type: type, forceRefresh: true);
  }

  void updateDueDateFilter(DateTime? from, DateTime? to) {
    state = state.copyWith(dueDateFrom: from, dueDateTo: to);
    _cache.clearPrefix('list:tasks');
    loadTasks(
      page: 1,
      refresh: true,
      dueDateFrom: from,
      dueDateTo: to,
      forceRefresh: true,
    );
  }

  void clearFilters() {
    state = state.copyWith(
      selectedStatus: null,
      selectedPriority: null,
      selectedType: null,
      dueDateFrom: null,
      dueDateTo: null,
      searchQuery: '',
    );
    _cache.clearPrefix('list:tasks');
    loadTasks(page: 1, refresh: true, forceRefresh: true);
  }

  /// Clear cache - exposed for TaskFormNotifier
  void clearCache() {
    _cache.clearPrefix('list:tasks');
  }
}

class TaskFormNotifier extends Notifier<TaskFormState> {
  late final TaskRepository _repository;

  @override
  TaskFormState build() {
    _repository = ref.read(taskRepositoryProvider);
    return const TaskFormState();
  }

  Future<Task> markInProgress(String id) async {
    state = state.copyWith(isLoading: true, errorMessage: null);

    try {
      final task = await _repository.markInProgress(id);

      // Clear cache and invalidate providers to refresh
      ref.read(taskListProvider.notifier).clearCache();
      ref.invalidate(taskDetailProvider(id));
      ref.invalidate(taskListProvider);

      state = state.copyWith(isLoading: false);
      return task;
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        errorMessage: e.toString().replaceFirst('Exception: ', ''),
      );
      rethrow;
    }
  }

  Future<bool> completeTask(String id) async {
    state = state.copyWith(isLoading: true, errorMessage: null);

    try {
      await _repository.completeTask(id);

      // Clear cache and invalidate providers to refresh
      ref.read(taskListProvider.notifier).clearCache();
      ref.invalidate(taskDetailProvider(id));
      ref.invalidate(taskListProvider);

      state = state.copyWith(isLoading: false);
      return true;
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        errorMessage: e.toString().replaceFirst('Exception: ', ''),
      );
      return false;
    }
  }

  Future<Reminder?> createReminder({
    required String taskId,
    required DateTime remindAt,
    String reminderType = 'in_app',
    String? message,
  }) async {
    state = state.copyWith(isLoading: true, errorMessage: null);

    try {
      final reminder = await _repository.createReminder(
        taskId: taskId,
        remindAt: remindAt,
        reminderType: reminderType,
        message: message,
      );

      // Invalidate task detail to refresh
      ref.invalidate(taskDetailProvider(taskId));

      state = state.copyWith(isLoading: false);
      return reminder;
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        errorMessage: e.toString().replaceFirst('Exception: ', ''),
      );
      return null;
    }
  }

  Future<bool> deleteReminder(String id, String taskId) async {
    state = state.copyWith(isLoading: true, errorMessage: null);

    try {
      await _repository.deleteReminder(id);

      // Invalidate task detail to refresh
      ref.invalidate(taskDetailProvider(taskId));

      state = state.copyWith(isLoading: false);
      return true;
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        errorMessage: e.toString().replaceFirst('Exception: ', ''),
      );
      return false;
    }
  }
}
