import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/cache/list_cache.dart';
import '../../../core/network/api_client.dart';
import '../../../core/network/connectivity_service.dart';
import '../../../core/sync/auto_sync_manager.dart';
import '../data/contact_repository.dart';
import '../data/role_repository.dart';
import '../data/models/contact.dart';
import 'contact_state.dart';

final contactRepositoryProvider = Provider<ContactRepository>((ref) {
  final connectivity = ref.watch(connectivityServiceProvider);
  return ContactRepository(ApiClient.dio, connectivity);
});

final roleRepositoryProvider = Provider<RoleRepository>((ref) {
  return RoleRepository(ApiClient.dio);
});

final rolesProvider = FutureProvider<List<Role>>((ref) async {
  final repository = ref.read(roleRepositoryProvider);
  return repository.getRoles();
});

final contactListProvider =
    NotifierProvider<ContactListNotifier, ContactListState>(
      ContactListNotifier.new,
    );

final contactDetailProvider = FutureProvider.family<Contact, String>((
  ref,
  id,
) async {
  final repository = ref.read(contactRepositoryProvider);
  return repository.getContactById(id);
});

class ContactListNotifier extends Notifier<ContactListState>
    with WidgetsBindingObserver {
  late final ContactRepository _repository;
  late final ConnectivityService _connectivity;
  final ListCache _cache = ListCache();

  @override
  ContactListState build() {
    _repository = ref.read(contactRepositoryProvider);
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
          .registerFeature('contacts', () => silentRefresh());
    });

    return const ContactListState();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    super.didChangeAppLifecycleState(state);
    if (state == AppLifecycleState.resumed) {
      // Auto sync when app comes back to foreground
      ref.read(autoSyncManagerProvider.notifier).syncFeature('contacts');
    }
  }

  Future<void> loadContacts({
    int page = 1,
    bool refresh = false,
    String? search,
    String? accountId,
    bool forceRefresh = false,
  }) async {
    final searchQuery = search ?? state.searchQuery;
    final filterAccountId = accountId ?? state.accountId;
    final cacheKey = ListCache.cacheKey(
      'contacts',
      page: page,
      search: searchQuery.isNotEmpty ? searchQuery : null,
      filters: filterAccountId != null ? {'account_id': filterAccountId} : null,
    );

    // Try to load from cache first (optimistic UI) - only for first page
    if (!forceRefresh && !refresh && page == 1) {
      final cachedContacts = _cache.get<Contact>(
        cacheKey,
        ttl: const Duration(seconds: 60),
        expectedMetadata: {
          if (searchQuery.isNotEmpty) 'search': searchQuery,
          ...?(filterAccountId != null
              ? {'account_id': filterAccountId}
              : null),
        },
      );

      if (cachedContacts != null && cachedContacts.isNotEmpty) {
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
          contacts: cachedContacts,
          searchQuery: searchQuery,
          accountId: filterAccountId,
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
      final response = await _repository.getContacts(
        page: page,
        perPage: 20,
        search: searchQuery.isNotEmpty ? searchQuery : null,
        accountId: filterAccountId,
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
          ...?(filterAccountId != null
              ? {'account_id': filterAccountId}
              : null),
        },
      );

      if (refresh || page == 1) {
        state = state.copyWith(
          contacts: response.items,
          pagination: response.pagination,
          searchQuery: searchQuery,
          accountId: filterAccountId,
          isLoading: false,
          isLoadingMore: false,
          errorMessage: null,
          isOffline: !_connectivity.isOnline,
        );
      } else {
        state = state.copyWith(
          contacts: [...state.contacts, ...response.items],
          pagination: response.pagination,
          isLoadingMore: false,
          errorMessage: null,
          isOffline: !_connectivity.isOnline,
        );
      }
    } catch (e) {
      // On error, try to use cached data as fallback
      if (page == 1) {
        final cachedContacts = _cache.get<Contact>(cacheKey);
        if (cachedContacts != null && cachedContacts.isNotEmpty) {
          state = state.copyWith(
            contacts: cachedContacts,
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
      final response = await _repository.getContacts(
        page: 1,
        perPage: 20,
        search: state.searchQuery.isNotEmpty ? state.searchQuery : null,
        accountId: state.accountId,
        forceRefresh: true,
      );

      // Update state with fresh data
      state = state.copyWith(
        contacts: response.items,
        pagination: response.pagination,
        errorMessage: null,
      );

      // Cache the response
      final cacheKey = ListCache.cacheKey(
        'contacts',
        page: 1,
        search: state.searchQuery.isNotEmpty ? state.searchQuery : null,
        filters: state.accountId != null
            ? {'account_id': state.accountId!}
            : null,
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
          ...?(state.accountId != null
              ? {'account_id': state.accountId!}
              : null),
        },
      );
    } catch (e) {
      // Silent fail
    }
  }

  /// Update state with fresh data from background fetch
  void _updateStateWithFreshData(ContactListResponse freshData) {
    state = state.copyWith(
      contacts: freshData.items,
      pagination: freshData.pagination,
      errorMessage: null,
    );

    // Cache the response
    final cacheKey = ListCache.cacheKey(
      'contacts',
      page: 1,
      search: state.searchQuery.isNotEmpty ? state.searchQuery : null,
      filters: state.accountId != null
          ? {'account_id': state.accountId!}
          : null,
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
        ...?(state.accountId != null ? {'account_id': state.accountId!} : null),
      },
    );
  }

  Future<void> refresh() async {
    // Clear cache for contacts
    _cache.clearPrefix('list:contacts');
    await loadContacts(page: 1, refresh: true, forceRefresh: true);
  }

  Future<void> loadMore() async {
    if (state.isLoading || state.isLoadingMore) return;
    final pagination = state.pagination;
    if (pagination == null || !pagination.hasNextPage) return;

    await loadContacts(page: pagination.page + 1);
  }

  void updateSearchQuery(String query) {
    state = state.copyWith(searchQuery: query);
  }

  void setAccountFilter(String? accountId) {
    state = state.copyWith(accountId: accountId);
  }

  Future<Contact?> createContact({
    required String accountId,
    required String name,
    required String roleId,
    String? phone,
    String? email,
    String? position,
    String? notes,
  }) async {
    try {
      state = state.copyWith(isLoading: true, errorMessage: null);
      final contact = await _repository.createContact(
        accountId: accountId,
        name: name,
        roleId: roleId,
        phone: phone,
        email: email,
        position: position,
        notes: notes,
      );
      state = state.copyWith(isLoading: false);
      // Refresh list
      await loadContacts(page: 1, refresh: true);
      return contact;
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        errorMessage: e.toString().replaceFirst('Exception: ', ''),
      );
      return null;
    }
  }

  Future<Contact?> updateContact({
    required String id,
    String? accountId,
    String? name,
    String? roleId,
    String? phone,
    String? email,
    String? position,
    String? notes,
  }) async {
    try {
      state = state.copyWith(isLoading: true, errorMessage: null);
      final contact = await _repository.updateContact(
        id: id,
        accountId: accountId,
        name: name,
        roleId: roleId,
        phone: phone,
        email: email,
        position: position,
        notes: notes,
      );
      state = state.copyWith(isLoading: false);
      // Refresh list
      await loadContacts(page: 1, refresh: true);
      return contact;
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        errorMessage: e.toString().replaceFirst('Exception: ', ''),
      );
      return null;
    }
  }

  Future<bool> deleteContact(String id) async {
    try {
      state = state.copyWith(isLoading: true, errorMessage: null);
      await _repository.deleteContact(id);
      state = state.copyWith(isLoading: false);
      // Clear cache and refresh list
      _cache.clearPrefix('list:contacts');
      await loadContacts(page: 1, refresh: true, forceRefresh: true);
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
