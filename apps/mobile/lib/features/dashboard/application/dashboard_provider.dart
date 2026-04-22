import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/network/connectivity_service.dart';
import '../../../core/sync/auto_sync_manager.dart';
import '../data/dashboard_repository.dart';
import 'dashboard_state.dart';

final dashboardRepositoryProvider = Provider<DashboardRepository>((ref) {
  return DashboardRepository();
});

final dashboardProvider = NotifierProvider<DashboardNotifier, DashboardState>(
  DashboardNotifier.new,
);

class DashboardNotifier extends Notifier<DashboardState>
    with WidgetsBindingObserver {
  late final DashboardRepository _repository;
  late final ConnectivityService _connectivity;

  @override
  DashboardState build() {
    _repository = ref.read(dashboardRepositoryProvider);
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
          .registerFeature('dashboard', () => silentRefresh());
    });

    return DashboardState();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    super.didChangeAppLifecycleState(state);
    if (state == AppLifecycleState.resumed) {
      // Auto sync when app comes back to foreground
      ref.read(autoSyncManagerProvider.notifier).syncFeature('dashboard');
    }
  }

  Future<void> loadDashboard() async {
    state = state.copyWith(isLoading: true, errorMessage: null);
    try {
      await Future.wait([
        _loadOverview(state.selectedPeriod, false),
        _loadVisits(state.selectedPeriod, state.visitStatusFilter, false),
        _loadTasks(state.selectedPeriod, false),
      ]);
      state = state.copyWith(
        isLoading: false,
        isOffline: !_connectivity.isOnline,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        errorMessage: _extractErrorMessage(e),
        isOffline: !_connectivity.isOnline,
      );
    }
  }

  Future<void> _loadOverview(String period, bool forceRefresh) async {
    try {
      final overview = await _repository.getMobileOverview(
        period: period,
        forceRefresh: forceRefresh,
      );
      state = state.copyWith(overview: overview, isLoadingOverview: false);
    } catch (e) {
      state = state.copyWith(
        isLoadingOverview: false,
        errorMessage: _extractErrorMessage(e),
      );
    }
  }

  Future<void> _loadVisits(
    String period,
    String status,
    bool forceRefresh,
  ) async {
    try {
      final visits = await _repository.getMobileVisits(
        period: period,
        status: status,
        forceRefresh: forceRefresh,
      );
      state = state.copyWith(visits: visits, isLoadingVisits: false);
    } catch (e) {
      state = state.copyWith(isLoadingVisits: false);
    }
  }

  Future<void> _loadTasks(String period, bool forceRefresh) async {
    try {
      final tasks = await _repository.getMobileTasks(
        period: period,
        forceRefresh: forceRefresh,
      );
      state = state.copyWith(tasks: tasks, isLoadingTasks: false);
    } catch (e) {
      state = state.copyWith(isLoadingTasks: false);
    }
  }

  /// Silent refresh (dipanggil oleh AutoSyncManager)
  Future<void> silentRefresh() async {
    try {
      await Future.wait([
        _loadOverview(state.selectedPeriod, true),
        _loadVisits(state.selectedPeriod, state.visitStatusFilter, true),
        _loadTasks(state.selectedPeriod, true),
      ]);
    } catch (e) {
      // Silent fail
    }
  }

  Future<void> refresh() async {
    await loadDashboard();
  }

  void setPeriod(String period) {
    if (state.selectedPeriod != period) {
      state = state.copyWith(selectedPeriod: period);
      loadDashboard();
    }
  }

  void setVisitStatusFilter(String status) {
    if (state.visitStatusFilter != status) {
      state = state.copyWith(visitStatusFilter: status);
      _loadVisits(state.selectedPeriod, status, true);
    }
  }

  String _extractErrorMessage(dynamic error) {
    if (error is DioException) {
      if (error.response != null) {
        final data = error.response!.data;
        if (data is Map<String, dynamic> && data.containsKey('error')) {
          return data['error'].toString();
        }
        return 'Server error: ${error.response!.statusCode}';
      }
      if (error.type == DioExceptionType.connectionTimeout ||
          error.type == DioExceptionType.receiveTimeout) {
        return 'Connection timeout. Please check your internet connection.';
      }
      if (error.type == DioExceptionType.connectionError) {
        return 'No internet connection. Please check your network.';
      }
      return 'Network error: ${error.message}';
    }
    return error.toString();
  }
}
