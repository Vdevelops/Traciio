import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/cache/list_cache.dart';
import '../../../core/network/api_client.dart';
import '../../../core/network/connectivity_service.dart';
import '../../../core/sync/auto_sync_manager.dart';
import '../data/models/deal_model.dart';
import '../data/pipeline_repository.dart';
import 'pipeline_state.dart';

final pipelineRepositoryProvider = Provider<PipelineRepository>((ref) {
  final connectivity = ref.watch(connectivityServiceProvider);
  return PipelineRepository(ApiClient.dio, connectivity);
});

final pipelineStagesProvider = FutureProvider<List<PipelineStage>>((ref) async {
  final repository = ref.read(pipelineRepositoryProvider);
  return repository.getStages();
});

final pipelineProvider = NotifierProvider<PipelineNotifier, PipelineState>(
  PipelineNotifier.new,
);

class PipelineNotifier extends Notifier<PipelineState>
    with WidgetsBindingObserver {
  late final PipelineRepository _repository;
  late final ConnectivityService _connectivity;
  final ListCache _cache = ListCache();

  @override
  PipelineState build() {
    _repository = ref.read(pipelineRepositoryProvider);
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
          .registerFeature('pipeline', () => silentRefresh());
    });

    return const PipelineState();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    super.didChangeAppLifecycleState(state);
    if (state == AppLifecycleState.resumed) {
      // Auto sync when app comes back to foreground
      ref.read(autoSyncManagerProvider.notifier).syncFeature('pipeline');
    }
  }

  Future<void> loadStages() async {
    try {
      final stages = await _repository.getStages();

      // Validate if current selectedStageId still exists in the new list
      // This handles cases where backend is reset and UUIDs change
      final isValidStage =
          state.selectedStageId != null &&
          stages.any((s) => s.id == state.selectedStageId);

      if (!isValidStage && stages.isNotEmpty) {
        state = state.copyWith(
          stages: stages,
          selectedStageId: stages.first.id,
        );
        // Also reload deals for the new valid stage
        loadDeals(page: 1, refresh: true, stageId: stages.first.id);
      } else {
        state = state.copyWith(stages: stages);
      }
    } catch (e) {
      state = state.copyWith(
        errorMessage: e.toString().replaceFirst('Exception: ', ''),
      );
    }
  }

  Future<void> loadFormData({bool forceRefresh = false}) async {
    try {
      final formData = await _repository.getFormData(
        forceRefresh: forceRefresh,
      );
      final productsJson = formData['products'] as List?;
      final products =
          productsJson?.map((e) => Product.fromJson(e)).toList() ?? [];

      state = state.copyWith(formData: formData, products: products);
    } catch (e) {
      // Error loading form data, fallback to cache is handled in repository
    }
  }

  Future<void> loadDeals({
    int page = 1,
    bool refresh = false,
    String? search,
    String? stageId,
    bool forceRefresh = false,
  }) async {
    final searchQuery = search ?? state.searchQuery;
    final currentStageId = stageId ?? state.selectedStageId;

    if (currentStageId == null) return;

    final cacheKey = ListCache.cacheKey(
      'deals',
      page: page,
      search: searchQuery.isNotEmpty ? searchQuery : null,
      filters: {'stage_id': currentStageId},
    );

    // Try to load from cache first (optimistic UI) - only for first page
    if (!forceRefresh && !refresh && page == 1) {
      final cachedDeals = _cache.get<Deal>(
        cacheKey,
        ttl: const Duration(seconds: 60),
        expectedMetadata: searchQuery.isNotEmpty
            ? {'search': searchQuery}
            : null,
      );

      if (cachedDeals != null && cachedDeals.isNotEmpty) {
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
          deals: cachedDeals,
          searchQuery: searchQuery,
          selectedStageId: currentStageId,
          isLoading: false,
          isLoadingMore: false,
          errorMessage: null,
          pagination: cachedPagination,
        );
      }
    }

    // Set loading state
    if (refresh || page == 1) {
      state = state.copyWith(
        isLoading: true,
        isLoadingMore: false,
        errorMessage: null,
      );
    } else {
      state = state.copyWith(isLoadingMore: true);
    }

    try {
      final response = await _repository.getDeals(
        page: page,
        perPage: 20,
        search: searchQuery.isNotEmpty ? searchQuery : null,
        stageId: currentStageId,
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
          'stage_id': currentStageId,
        },
      );

      if (refresh || page == 1) {
        state = state.copyWith(
          deals: response.items,
          pagination: response.pagination,
          searchQuery: searchQuery,
          selectedStageId: currentStageId,
          isLoading: false,
          isLoadingMore: false,
          errorMessage: null,
          isOffline: !_connectivity.isOnline,
        );
      } else {
        state = state.copyWith(
          deals: [...state.deals, ...response.items],
          pagination: response.pagination,
          isLoadingMore: false,
          errorMessage: null,
          isOffline: !_connectivity.isOnline,
        );
      }
    } catch (e) {
      // On error, try to use cached data as fallback
      if (page == 1) {
        final cachedDeals = _cache.get<Deal>(cacheKey);
        if (cachedDeals != null && cachedDeals.isNotEmpty) {
          state = state.copyWith(
            deals: cachedDeals,
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
    final currentStageId = state.selectedStageId;
    if (currentStageId == null) return;

    try {
      final response = await _repository.getDeals(
        page: 1,
        perPage: 20,
        search: state.searchQuery.isNotEmpty ? state.searchQuery : null,
        stageId: currentStageId,
        forceRefresh: true,
      );

      // Update state with fresh data
      state = state.copyWith(
        deals: response.items,
        pagination: response.pagination,
        errorMessage: null,
      );

      // Cache the response
      final cacheKey = ListCache.cacheKey(
        'deals',
        page: 1,
        search: state.searchQuery.isNotEmpty ? state.searchQuery : null,
        filters: {'stage_id': currentStageId},
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
          'stage_id': currentStageId,
        },
      );
    } catch (e) {
      // Silent fail
    }
  }

  /// Update state with fresh data from background fetch
  void _updateStateWithFreshData(DealListResponse freshData) {
    final currentStageId = state.selectedStageId;
    if (currentStageId == null) return;

    state = state.copyWith(
      deals: freshData.items,
      pagination: freshData.pagination,
      errorMessage: null,
    );

    // Cache the response
    final cacheKey = ListCache.cacheKey(
      'deals',
      page: 1,
      search: state.searchQuery.isNotEmpty ? state.searchQuery : null,
      filters: {'stage_id': currentStageId},
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
        'stage_id': currentStageId,
      },
    );
  }

  Future<void> refresh() async {
    _cache.clearPrefix('list:deals');
    await loadStages();
    await loadDeals(page: 1, refresh: true, forceRefresh: true);
  }

  Future<void> loadMore() async {
    if (state.isLoading || state.isLoadingMore) return;
    final pagination = state.pagination;
    if (pagination == null || !pagination.hasNextPage) return;

    await loadDeals(page: pagination.page + 1);
  }

  Future<void> saveDeal({
    required String title,
    required int value,
    String? accountId,
    String? contactId,
    required String stageId,
    String? notes,
    List<Map<String, dynamic>>? productItems,
    String? dealId,
  }) async {
    final data = {
      'title': title,
      'value': value,
      'account_id': accountId,
      'contact_id': contactId,
      'stage_id': stageId,
      'notes': notes,
      ...?(productItems != null ? {'product_items': productItems} : null),
    };

    try {
      if (dealId == null) {
        await _repository.createDeal(data);
      } else {
        await _repository.updateDeal(dealId, data);
      }
      // Refresh the list after saving
      await refresh();
    } catch (e) {
      throw Exception('Failed to save deal: $e');
    }
  }

  void updateSearchQuery(String query) {
    state = state.copyWith(searchQuery: query);
  }

  void selectStage(String stageId) {
    if (state.selectedStageId == stageId) return;
    state = state.copyWith(selectedStageId: stageId, deals: []);
    loadDeals(page: 1, refresh: true);
  }

  void clearCache() {
    _cache.clearPrefix('list:deals');
  }
}
