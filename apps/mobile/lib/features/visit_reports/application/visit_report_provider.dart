import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/cache/list_cache.dart';
import '../../../core/network/api_client.dart';
import '../../../core/network/connectivity_service.dart';
import '../../../core/sync/auto_sync_manager.dart';
import '../data/visit_report_repository.dart';
import '../data/models/visit_report.dart';
import 'visit_report_state.dart';

final visitReportRepositoryProvider = Provider<VisitReportRepository>((ref) {
  final connectivity = ref.watch(connectivityServiceProvider);
  return VisitReportRepository(ApiClient.dio, connectivity);
});

final visitReportListProvider =
    NotifierProvider<VisitReportListNotifier, VisitReportListState>(
      VisitReportListNotifier.new,
    );

final visitReportDetailProvider = FutureProvider.family<VisitReport, String>((
  ref,
  id,
) async {
  final repository = ref.read(visitReportRepositoryProvider);
  return repository.getVisitReportById(id);
});

final visitReportFormProvider =
    NotifierProvider<VisitReportFormNotifier, VisitReportFormState>(
      VisitReportFormNotifier.new,
    );

/// Provider for form data (accounts, contacts, deals, leads)
final visitReportFormDataProvider = FutureProvider<Map<String, dynamic>>((
  ref,
) async {
  final repository = ref.read(visitReportRepositoryProvider);
  return repository.getFormData();
});

class VisitReportListNotifier extends Notifier<VisitReportListState>
    with WidgetsBindingObserver {
  late final VisitReportRepository _repository;
  late final ConnectivityService _connectivity;
  final ListCache _cache = ListCache();

  @override
  VisitReportListState build() {
    _repository = ref.read(visitReportRepositoryProvider);
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
          .registerFeature('visit_reports', () => silentRefresh());
    });

    return const VisitReportListState();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    super.didChangeAppLifecycleState(state);
    if (state == AppLifecycleState.resumed) {
      // Auto sync when app comes back to foreground
      ref.read(autoSyncManagerProvider.notifier).syncFeature('visit_reports');
    }
  }

  Future<void> loadVisitReports({
    int page = 1,
    bool refresh = false,
    String? search,
    bool forceRefresh = false,
    bool forRouteOptimization = false,
    String? status,
  }) async {
    final searchQuery = search ?? state.searchQuery;
    final cacheKey = ListCache.cacheKey(
      'visit-reports',
      page: page,
      search: searchQuery.isNotEmpty ? searchQuery : null,
    );

    // Try to load from cache first (optimistic UI) - only for first page
    if (!forceRefresh && !refresh && page == 1) {
      final cachedVisitReports = _cache.get<VisitReport>(
        cacheKey,
        ttl: const Duration(seconds: 60),
        expectedMetadata: searchQuery.isNotEmpty
            ? {'search': searchQuery}
            : null,
      );

      if (cachedVisitReports != null && cachedVisitReports.isNotEmpty) {
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
          visitReports: cachedVisitReports,
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
      final response = await _repository.getVisitReports(
        page: page,
        perPage: 20,
        search: searchQuery.isNotEmpty ? searchQuery : null,
        forceRefresh: forceRefresh,
        forRouteOptimization: forRouteOptimization,
        status: status,
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
          visitReports: response.items,
          pagination: response.pagination,
          searchQuery: searchQuery,
          isLoading: false,
          isLoadingMore: false,
          errorMessage: null,
          isOffline: !_connectivity.isOnline,
        );
      } else {
        state = state.copyWith(
          visitReports: [...state.visitReports, ...response.items],
          pagination: response.pagination,
          isLoadingMore: false,
          errorMessage: null,
          isOffline: !_connectivity.isOnline,
        );
      }
    } catch (e) {
      // On error, try to use cached data as fallback
      if (page == 1) {
        final cachedVisitReports = _cache.get<VisitReport>(cacheKey);
        if (cachedVisitReports != null && cachedVisitReports.isNotEmpty) {
          state = state.copyWith(
            visitReports: cachedVisitReports,
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
      final response = await _repository.getVisitReports(
        page: 1,
        perPage: 20,
        search: state.searchQuery.isNotEmpty ? state.searchQuery : null,
        forceRefresh: true,
      );

      // Update state with fresh data
      state = state.copyWith(
        visitReports: response.items,
        pagination: response.pagination,
        errorMessage: null,
      );

      // Cache the response
      final cacheKey = ListCache.cacheKey(
        'visit-reports',
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
  void _updateStateWithFreshData(VisitReportListResponse freshData) {
    state = state.copyWith(
      visitReports: freshData.items,
      pagination: freshData.pagination,
      errorMessage: null,
    );

    // Cache the response
    final cacheKey = ListCache.cacheKey(
      'visit-reports',
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
    // Clear cache for visit reports
    _cache.clearPrefix('list:visit-reports');
    await loadVisitReports(page: 1, refresh: true, forceRefresh: true);
  }

  Future<void> loadMore() async {
    if (state.isLoading || state.isLoadingMore) return;
    final pagination = state.pagination;
    if (pagination == null || !pagination.hasNextPage) return;

    await loadVisitReports(page: pagination.page + 1);
  }

  void updateSearchQuery(String query) {
    state = state.copyWith(searchQuery: query);
  }

  /// Clear cache - exposed for VisitReportFormNotifier
  void clearCache() {
    _cache.clearPrefix('list:visit-reports');
  }
}

class VisitReportFormNotifier extends Notifier<VisitReportFormState> {
  late final VisitReportRepository _repository;

  @override
  VisitReportFormState build() {
    _repository = ref.read(visitReportRepositoryProvider);
    return const VisitReportFormState();
  }

  Future<VisitReport?> createVisitReport({
    String? accountId,
    String? contactId,
    String? dealId,
    String? leadId,
    required String visitDate,
    required String purpose,
    String? notes,
  }) async {
    state = state.copyWith(isSubmitting: true, errorMessage: null);

    try {
      final visitReport = await _repository.createVisitReport(
        accountId: accountId,
        contactId: contactId,
        dealId: dealId,
        leadId: leadId,
        visitDate: visitDate,
        purpose: purpose,
        notes: notes,
      );

      state = state.copyWith(isSubmitting: false);
      // Clear cache and refresh list
      ref.read(visitReportListProvider.notifier).clearCache();
      ref.read(visitReportListProvider.notifier).refresh();
      return visitReport;
    } catch (e) {
      state = state.copyWith(
        isSubmitting: false,
        errorMessage: e.toString().replaceFirst('Exception: ', ''),
      );
      return null;
    }
  }

  Future<VisitReport?> updateVisitReport({
    required String id,
    String? accountId,
    String? contactId,
    String? dealId,
    String? leadId,
    String? visitDate,
    String? purpose,
    String? notes,
  }) async {
    state = state.copyWith(isLoading: true, errorMessage: null);

    try {
      final visitReport = await _repository.updateVisitReport(
        id: id,
        accountId: accountId,
        contactId: contactId,
        dealId: dealId,
        leadId: leadId,
        visitDate: visitDate,
        purpose: purpose,
        notes: notes,
      );

      // Clear cache and update with new data
      await _repository.clearVisitReportDetailCache(id);
      await _repository.saveVisitReportDetailToCache(id, visitReport);

      // Invalidate detail provider to refresh
      ref.invalidate(visitReportDetailProvider(id));
      // Clear cache and refresh list
      ref.read(visitReportListProvider.notifier).clearCache();
      ref.read(visitReportListProvider.notifier).refresh();

      state = state.copyWith(isLoading: false);
      return visitReport;
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        errorMessage: e.toString().replaceFirst('Exception: ', ''),
      );
      return null;
    }
  }

  Future<bool> deleteVisitReport(String id) async {
    state = state.copyWith(isLoading: true, errorMessage: null);

    try {
      await _repository.deleteVisitReport(id);

      // Clear cache
      await _repository.clearVisitReportDetailCache(id);

      // Invalidate detail provider
      ref.invalidate(visitReportDetailProvider(id));
      // Clear cache and refresh list
      ref.read(visitReportListProvider.notifier).clearCache();
      ref.read(visitReportListProvider.notifier).refresh();

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

  Future<VisitReport?> submitVisitReport({
    required String id,
    String? outcome,
    String? nextSteps,
  }) async {
    state = state.copyWith(isLoading: true, errorMessage: null);

    try {
      final visitReport = await _repository.submitVisitReport(
        id: id,
        outcome: outcome,
        nextSteps: nextSteps,
      );

      // Clear cache and update with new data
      await _repository.clearVisitReportDetailCache(id);
      await _repository.saveVisitReportDetailToCache(id, visitReport);

      // Invalidate detail provider to refresh
      ref.invalidate(visitReportDetailProvider(id));
      // Clear cache and refresh list
      ref.read(visitReportListProvider.notifier).clearCache();
      ref.read(visitReportListProvider.notifier).refresh();

      state = state.copyWith(isLoading: false);
      return visitReport;
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        errorMessage: e.toString().replaceFirst('Exception: ', ''),
      );
      return null;
    }
  }

  Future<VisitReport?> checkIn({
    required String visitReportId,
    required double latitude,
    required double longitude,
    required File photoFile, // Selfie picture is required for mobile check-in
    double? accuracy,
    double? photoLatitude,
    double? photoLongitude,
    int? photoTimestamp,
  }) async {
    state = state.copyWith(isLoading: true, errorMessage: null);

    try {
      final visitReport = await _repository.checkIn(
        visitReportId: visitReportId,
        latitude: latitude,
        longitude: longitude,
        photoFile: photoFile,
        accuracy: accuracy,
        photoLatitude: photoLatitude,
        photoLongitude: photoLongitude,
        photoTimestamp: photoTimestamp,
      );

      // Clear cache and update with new data
      await _repository.clearVisitReportDetailCache(visitReportId);
      await _repository.saveVisitReportDetailToCache(
        visitReportId,
        visitReport,
      );

      // Invalidate detail provider to refresh
      ref.invalidate(visitReportDetailProvider(visitReportId));

      state = state.copyWith(isLoading: false);
      return visitReport;
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        errorMessage: e.toString().replaceFirst('Exception: ', ''),
      );
      return null;
    }
  }

  Future<VisitReport?> checkOut({
    required String visitReportId,
    required double latitude,
    required double longitude,
    File? photoFile, // Optional photo for check-out (like web version)
    double? accuracy,
    double? photoLatitude,
    double? photoLongitude,
    int? photoTimestamp,
  }) async {
    state = state.copyWith(isLoading: true, errorMessage: null);

    try {
      final visitReport = await _repository.checkOut(
        visitReportId: visitReportId,
        latitude: latitude,
        longitude: longitude,
        photoFile: photoFile,
        accuracy: accuracy,
        photoLatitude: photoLatitude,
        photoLongitude: photoLongitude,
        photoTimestamp: photoTimestamp,
      );

      // Clear cache and update with new data
      await _repository.clearVisitReportDetailCache(visitReportId);
      await _repository.saveVisitReportDetailToCache(
        visitReportId,
        visitReport,
      );

      // Invalidate detail provider to refresh
      ref.invalidate(visitReportDetailProvider(visitReportId));

      state = state.copyWith(isLoading: false);
      return visitReport;
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        errorMessage: e.toString().replaceFirst('Exception: ', ''),
      );
      return null;
    }
  }

  Future<bool> uploadPhoto({
    required String visitReportId,
    required File photoFile,
  }) async {
    state = state.copyWith(isLoading: true, errorMessage: null);

    try {
      await _repository.uploadPhoto(
        visitReportId: visitReportId,
        photoFile: photoFile,
      );

      // Clear cache to force refresh from API
      await _repository.clearVisitReportDetailCache(visitReportId);

      // Invalidate detail provider to refresh
      ref.invalidate(visitReportDetailProvider(visitReportId));

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
