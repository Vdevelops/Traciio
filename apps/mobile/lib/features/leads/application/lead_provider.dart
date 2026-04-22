import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/network/api_client.dart';
import '../../../core/network/connectivity_service.dart';
import '../../../core/sync/auto_sync_manager.dart';
import '../data/models/lead_model.dart';
import '../data/lead_repository.dart';
import 'lead_state.dart';

final leadRepositoryProvider = Provider<LeadRepository>((ref) {
  return LeadRepository(ApiClient.dio, ref.watch(connectivityServiceProvider));
});

final leadListProvider = NotifierProvider<LeadListNotifier, LeadListState>(
  LeadListNotifier.new,
);

final leadDetailProvider = FutureProvider.family<Lead, String>((ref, id) async {
  final repository = ref.read(leadRepositoryProvider);
  return repository.getLeadById(id);
});

final leadFormDataProvider = FutureProvider<LeadFormData>((ref) async {
  final repository = ref.read(leadRepositoryProvider);
  return repository.getFormData();
});

final leadFormProvider = NotifierProvider<LeadFormNotifier, LeadFormState>(
  LeadFormNotifier.new,
);

class LeadListNotifier extends Notifier<LeadListState>
    with WidgetsBindingObserver {
  late final LeadRepository _repository;

  @override
  LeadListState build() {
    _repository = ref.read(leadRepositoryProvider);

    // Listen to app lifecycle changes for auto-sync
    WidgetsBinding.instance.addObserver(this);
    ref.onDispose(() {
      WidgetsBinding.instance.removeObserver(this);
    });

    // Register with AutoSyncManager for centralized sync
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref
          .read(autoSyncManagerProvider.notifier)
          .registerFeature('leads', () => silentRefresh());
    });

    return const LeadListState();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    super.didChangeAppLifecycleState(state);
    if (state == AppLifecycleState.resumed) {
      // Auto sync when app comes back to foreground
      ref.read(autoSyncManagerProvider.notifier).syncFeature('leads');
    }
  }

  Future<void> loadLeads({
    int page = 1,
    bool refresh = false,
    String? search,
    String? status,
    String? source,
    String? industry,
    String? province,
  }) async {
    final searchQuery = search ?? state.searchQuery;
    final statusFilter = status ?? state.selectedStatus;
    final sourceFilter = source ?? state.selectedSource;
    final industryFilter = industry ?? state.selectedIndustry;
    final provinceFilter = province ?? state.selectedProvince;

    final isSearchQueryChanged = searchQuery != state.searchQuery;

    if (refresh || page == 1 || isSearchQueryChanged) {
      state = state.copyWith(
        leads: const [],
        isLoading: true,
        isLoadingMore: false,
        errorMessage: null,
        clearLeads: true,
      );
    } else {
      state = state.copyWith(isLoadingMore: true);
    }

    try {
      final response = await _repository.getLeads(
        page: page,
        perPage: 10,
        search: searchQuery.isNotEmpty ? searchQuery : null,
        status: statusFilter.isNotEmpty ? statusFilter : null,
        source: sourceFilter.isNotEmpty ? sourceFilter : null,
        industry: industryFilter.isNotEmpty ? industryFilter : null,
        province: provinceFilter.isNotEmpty ? provinceFilter : null,
        forceRefresh: refresh,
        onBackgroundUpdate: (freshData) {
          // Update UI when fresh data arrives from background fetch
          _updateStateWithFreshData(freshData);
        },
      );

      if (refresh || page == 1) {
        state = state.copyWith(
          leads: response.items,
          pagination: response.pagination,
          searchQuery: searchQuery,
          selectedStatus: statusFilter,
          selectedSource: sourceFilter,
          selectedIndustry: industryFilter,
          selectedProvince: provinceFilter,
          isLoading: false,
          isLoadingMore: false,
          errorMessage: null,
        );
      } else {
        state = state.copyWith(
          leads: [...state.leads, ...response.items],
          pagination: response.pagination,
          isLoadingMore: false,
          errorMessage: null,
        );
      }
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        isLoadingMore: false,
        errorMessage: e.toString().replaceFirst('Exception: ', ''),
      );
    }
  }

  Future<void> refresh() async {
    await loadLeads(page: 1, refresh: true);
  }

  /// Silent refresh (dipanggil oleh AutoSyncManager)
  Future<void> silentRefresh() async {
    try {
      final response = await _repository.getLeads(
        page: 1,
        search: state.searchQuery.isNotEmpty ? state.searchQuery : null,
        status: state.selectedStatus.isNotEmpty ? state.selectedStatus : null,
        source: state.selectedSource.isNotEmpty ? state.selectedSource : null,
        industry: state.selectedIndustry.isNotEmpty
            ? state.selectedIndustry
            : null,
        province: state.selectedProvince.isNotEmpty
            ? state.selectedProvince
            : null,
        forceRefresh: true,
      );

      state = state.copyWith(
        leads: response.items,
        pagination: response.pagination,
      );
    } catch (e) {
      // Silent refresh failed, ignore
    }
  }

  /// Update state with fresh data from background fetch
  void _updateStateWithFreshData(LeadListResponse freshData) {
    // Create new list instance to ensure Riverpod detects the change
    // Riverpod uses identity comparison for collections
    final newLeads = List<Lead>.from(freshData.items);

    state = LeadListState(
      leads: newLeads,
      isLoading: state.isLoading,
      isLoadingMore: state.isLoadingMore,
      errorMessage: state.errorMessage,
      pagination: freshData.pagination,
      searchQuery: state.searchQuery,
      selectedStatus: state.selectedStatus,
      selectedSource: state.selectedSource,
      selectedIndustry: state.selectedIndustry,
      selectedProvince: state.selectedProvince,
    );
  }

  Future<void> loadMore() async {
    if (state.isLoading || state.isLoadingMore) return;
    final pagination = state.pagination;
    if (pagination == null || pagination.page >= pagination.totalPages) return;

    await loadLeads(page: pagination.page + 1);
  }

  void updateSearchQuery(String query) {
    state = state.copyWith(
      searchQuery: query,
      leads: query != state.searchQuery ? const [] : state.leads,
    );
  }

  void updateStatusFilter(String? status) {
    state = state.copyWith(selectedStatus: status ?? '');
    loadLeads(page: 1, refresh: true, status: status ?? '');
  }

  void updateSourceFilter(String? source) {
    state = state.copyWith(selectedSource: source ?? '');
    loadLeads(page: 1, refresh: true, source: source ?? '');
  }

  void updateIndustryFilter(String? industry) {
    state = state.copyWith(selectedIndustry: industry ?? '');
    loadLeads(page: 1, refresh: true, industry: industry ?? '');
  }

  void updateProvinceFilter(String? province) {
    state = state.copyWith(selectedProvince: province ?? '');
    loadLeads(page: 1, refresh: true, province: province ?? '');
  }

  void clearFilters() {
    state = state.copyWith(
      selectedStatus: '',
      selectedSource: '',
      selectedIndustry: '',
      selectedProvince: '',
      searchQuery: '',
    );
    loadLeads(page: 1, refresh: true);
  }

  Future<bool> convertLead(String id) async {
    try {
      await _repository.convertLead(id, {});
      // Refresh the list to remove the converted lead or update its status
      await refresh();
      return true;
    } catch (e) {
      state = state.copyWith(
        errorMessage: e.toString().replaceFirst('Exception: ', ''),
      );
      return false;
    }
  }
}

class LeadFormNotifier extends Notifier<LeadFormState> {
  late final LeadRepository _repository;

  @override
  LeadFormState build() {
    _repository = ref.read(leadRepositoryProvider);
    return const LeadFormState();
  }

  Future<void> loadFormData() async {
    state = state.copyWith(isLoading: true, errorMessage: null);
    try {
      final formData = await _repository.getFormData();
      state = state.copyWith(isLoading: false, formData: formData);
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        errorMessage: e.toString().replaceFirst('Exception: ', ''),
      );
    }
  }

  Future<Lead?> createLead(Map<String, dynamic> data) async {
    state = state.copyWith(isLoading: true, errorMessage: null);
    try {
      final lead = await _repository.createLead(data);
      await ref.read(leadListProvider.notifier).refresh();
      state = state.copyWith(isLoading: false);
      return lead;
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        errorMessage: e.toString().replaceFirst('Exception: ', ''),
      );
      return null;
    }
  }

  Future<Lead?> updateLead(String id, Map<String, dynamic> data) async {
    state = state.copyWith(isLoading: true, errorMessage: null);
    try {
      final lead = await _repository.updateLead(id, data);
      final _ = ref.refresh(leadDetailProvider(id));
      await ref.read(leadListProvider.notifier).refresh();
      state = state.copyWith(isLoading: false);
      return lead;
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        errorMessage: e.toString().replaceFirst('Exception: ', ''),
      );
      return null;
    }
  }

  Future<bool> convertLead(String id, Map<String, dynamic> data) async {
    state = state.copyWith(isLoading: true, errorMessage: null);
    try {
      await _repository.convertLead(id, data);
      final _ = ref.refresh(leadDetailProvider(id));
      await ref.read(leadListProvider.notifier).refresh();
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
