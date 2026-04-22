import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/cache/list_cache.dart';
import '../../../core/network/api_client.dart';
import '../../../core/network/connectivity_service.dart';
import '../../../core/sync/auto_sync_manager.dart';
import '../data/account_repository.dart';
import '../data/category_repository.dart';
import '../data/models/account.dart';
import 'account_state.dart';

final accountRepositoryProvider = Provider<AccountRepository>((ref) {
  final connectivity = ref.watch(connectivityServiceProvider);
  return AccountRepository(ApiClient.dio, connectivity);
});

final categoryRepositoryProvider = Provider<CategoryRepository>((ref) {
  return CategoryRepository(ApiClient.dio);
});

final categoriesProvider = FutureProvider<List<Category>>((ref) async {
  final repository = ref.read(categoryRepositoryProvider);
  return repository.getCategories();
});

final accountListProvider =
    NotifierProvider<AccountListNotifier, AccountListState>(
      AccountListNotifier.new,
    );

final accountDetailProvider = FutureProvider.family<Account, String>((
  ref,
  id,
) async {
  final repository = ref.read(accountRepositoryProvider);
  return repository.getAccountById(id);
});

class AccountListNotifier extends Notifier<AccountListState>
    with WidgetsBindingObserver {
  late final AccountRepository _repository;
  late final ConnectivityService _connectivity;
  final ListCache _cache = ListCache();

  @override
  AccountListState build() {
    _repository = ref.read(accountRepositoryProvider);
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
          .registerFeature('accounts', () => silentRefresh());
    });

    return const AccountListState();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    super.didChangeAppLifecycleState(state);
    if (state == AppLifecycleState.resumed) {
      // Auto sync when app comes back to foreground
      ref.read(autoSyncManagerProvider.notifier).syncFeature('accounts');
    }
  }

  Future<void> loadAccounts({
    int page = 1,
    bool refresh = false,
    String? search,
    bool forceRefresh = false,
  }) async {
    final searchQuery = search ?? state.searchQuery;
    final cacheKey = ListCache.cacheKey(
      'accounts',
      page: page,
      search: searchQuery.isNotEmpty ? searchQuery : null,
    );

    // Try to load from cache first (optimistic UI) - only for first page
    if (!forceRefresh && !refresh && page == 1) {
      final cachedAccounts = _cache.get<Account>(
        cacheKey,
        ttl: const Duration(seconds: 60),
        expectedMetadata: searchQuery.isNotEmpty
            ? {'search': searchQuery}
            : null,
      );

      if (cachedAccounts != null && cachedAccounts.isNotEmpty) {
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
          accounts: cachedAccounts,
          searchQuery: searchQuery,
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
      final response = await _repository.getAccounts(
        page: page,
        perPage: 20,
        search: searchQuery.isNotEmpty ? searchQuery : null,
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
        },
      );

      if (refresh || page == 1) {
        state = state.copyWith(
          accounts: response.items,
          pagination: response.pagination,
          searchQuery: searchQuery,
          isLoading: false,
          isLoadingMore: false,
          errorMessage: null,
          isOffline: !_connectivity.isOnline,
        );
      } else {
        state = state.copyWith(
          accounts: [...state.accounts, ...response.items],
          pagination: response.pagination,
          isLoadingMore: false,
          errorMessage: null,
          isOffline: !_connectivity.isOnline,
        );
      }
    } catch (e) {
      // On error, try to use cached data as fallback
      if (page == 1) {
        final cachedAccounts = _cache.get<Account>(cacheKey);
        if (cachedAccounts != null && cachedAccounts.isNotEmpty) {
          state = state.copyWith(
            accounts: cachedAccounts,
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
      final response = await _repository.getAccounts(
        page: 1,
        perPage: 20,
        search: state.searchQuery.isNotEmpty ? state.searchQuery : null,
        forceRefresh: true,
      );

      // Update state with fresh data
      state = state.copyWith(
        accounts: response.items,
        pagination: response.pagination,
        errorMessage: null,
      );

      // Cache the response
      final cacheKey = ListCache.cacheKey(
        'accounts',
        page: 1,
        search: state.searchQuery.isNotEmpty ? state.searchQuery : null,
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
        },
      );
    } catch (e) {
      // Silent fail
    }
  }

  /// Update state with fresh data from background fetch
  void _updateStateWithFreshData(AccountListResponse freshData) {
    state = state.copyWith(
      accounts: freshData.items,
      pagination: freshData.pagination,
      errorMessage: null,
    );

    // Cache the response
    final cacheKey = ListCache.cacheKey(
      'accounts',
      page: 1,
      search: state.searchQuery.isNotEmpty ? state.searchQuery : null,
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
      },
    );
  }

  Future<void> refresh() async {
    // Clear cache for accounts
    _cache.clearPrefix('list:accounts');
    await loadAccounts(page: 1, refresh: true, forceRefresh: true);
  }

  Future<void> loadMore() async {
    if (state.isLoading || state.isLoadingMore) return;
    final pagination = state.pagination;
    if (pagination == null || !pagination.hasNextPage) return;

    await loadAccounts(page: pagination.page + 1);
  }

  void updateSearchQuery(String query) {
    state = state.copyWith(searchQuery: query);
  }

  Future<Account?> createAccount({
    required String name,
    required String categoryId,
    String? address,
    String? city,
    String? province,
    String? phone,
    String? email,
    String? status,
    String? assignedTo,
  }) async {
    try {
      state = state.copyWith(isLoading: true, errorMessage: null);
      final account = await _repository.createAccount(
        name: name,
        categoryId: categoryId,
        address: address,
        city: city,
        province: province,
        phone: phone,
        email: email,
        status: status ?? 'active',
        assignedTo: assignedTo,
      );
      state = state.copyWith(isLoading: false);
      // Clear cache and refresh list
      _cache.clearPrefix('list:accounts');
      await loadAccounts(page: 1, refresh: true, forceRefresh: true);
      return account;
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        errorMessage: e.toString().replaceFirst('Exception: ', ''),
      );
      return null;
    }
  }

  Future<Account?> updateAccount({
    required String id,
    String? name,
    String? categoryId,
    String? address,
    String? city,
    String? province,
    String? phone,
    String? email,
    String? status,
    String? assignedTo,
  }) async {
    try {
      state = state.copyWith(isLoading: true, errorMessage: null);
      final account = await _repository.updateAccount(
        id: id,
        name: name,
        categoryId: categoryId,
        address: address,
        city: city,
        province: province,
        phone: phone,
        email: email,
        status: status,
        assignedTo: assignedTo,
      );
      state = state.copyWith(isLoading: false);
      // Clear cache and refresh list
      _cache.clearPrefix('list:accounts');
      await loadAccounts(page: 1, refresh: true, forceRefresh: true);
      return account;
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        errorMessage: e.toString().replaceFirst('Exception: ', ''),
      );
      return null;
    }
  }

  Future<bool> deleteAccount(String id) async {
    try {
      state = state.copyWith(isLoading: true, errorMessage: null);
      await _repository.deleteAccount(id);
      state = state.copyWith(isLoading: false);
      // Clear cache and refresh list
      _cache.clearPrefix('list:accounts');
      await loadAccounts(page: 1, refresh: true, forceRefresh: true);
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
