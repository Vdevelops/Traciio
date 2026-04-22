import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/network/api_client.dart';
import '../../../core/network/connectivity_service.dart';
import '../../../core/sync/auto_sync_manager.dart';
import '../data/route_optimization_repository.dart';
import '../data/route_optimization_repository_optimized.dart';
import '../data/models/optimized_route.dart';
import '../data/models/optimize_route_request.dart';
import '../data/models/calculate_distance_request.dart';
import 'route_optimization_state.dart';

final routeOptimizationRepositoryProvider =
    Provider<RouteOptimizationRepository>((ref) {
      final connectivity = ref.watch(connectivityServiceProvider);
      // Use optimized repository with caching
      return OptimizedRouteOptimizationRepository(ApiClient.dio, connectivity);
    });

final routeListProvider = NotifierProvider<RouteListNotifier, RouteListState>(
  RouteListNotifier.new,
);

final routeDetailProvider = FutureProvider.family<OptimizedRoute, String>((
  ref,
  id,
) async {
  final repository = ref.read(routeOptimizationRepositoryProvider);
  return repository.getRouteById(id);
});

final routeOptimizationProvider =
    NotifierProvider<RouteOptimizationNotifier, RouteOptimizationState>(
      RouteOptimizationNotifier.new,
    );

final distanceCalculationProvider =
    NotifierProvider<DistanceCalculationNotifier, DistanceCalculationState>(
      DistanceCalculationNotifier.new,
    );

class RouteListNotifier extends Notifier<RouteListState>
    with WidgetsBindingObserver {
  late final RouteOptimizationRepository _repository;
  late final ConnectivityService _connectivity;

  @override
  RouteListState build() {
    _repository = ref.read(routeOptimizationRepositoryProvider);
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
          .registerFeature('route_optimization', () => silentRefresh());
    });

    return const RouteListState();
  }

  @override
  void didChangeAppLifecycleState(AppLifecycleState state) {
    super.didChangeAppLifecycleState(state);
    if (state == AppLifecycleState.resumed) {
      // Auto sync when app comes back to foreground
      ref
          .read(autoSyncManagerProvider.notifier)
          .syncFeature('route_optimization');
    }
  }

  Future<void> loadRoutes({
    int page = 1,
    bool refresh = false,
    bool forceRefresh = false,
  }) async {
    // Set loading state
    if (refresh || page == 1) {
      state = state.copyWith(
        isLoading: true,
        isLoadingMore: false,
        errorMessage: null,
        isOffline: !_connectivity.isOnline,
      );
    } else {
      state = state.copyWith(isLoadingMore: true, errorMessage: null);
    }

    try {
      final response = await _repository.getMyRoutes(
        page: page,
        perPage: 20,
        forceRefresh: forceRefresh,
        onBackgroundUpdate: page == 1 ? _updateStateWithFreshData : null,
      );

      if (page == 1) {
        state = state.copyWith(
          routes: response.items,
          pagination: response.pagination,
          isLoading: false,
          isLoadingMore: false,
          errorMessage: null,
          isOffline: false,
        );
      } else {
        state = state.copyWith(
          routes: [...state.routes, ...response.items],
          pagination: response.pagination,
          isLoading: false,
          isLoadingMore: false,
          errorMessage: null,
          isOffline: false,
        );
      }
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        isLoadingMore: false,
        errorMessage: e.toString().replaceAll('Exception: ', ''),
        isOffline: !_connectivity.isOnline,
      );
    }
  }

  Future<void> refresh() async {
    await loadRoutes(page: 1, refresh: true, forceRefresh: true);
  }

  /// Silent refresh (dipanggil oleh AutoSyncManager)
  Future<void> silentRefresh() async {
    try {
      final response = await _repository.getMyRoutes(
        page: 1,
        perPage: 20,
        forceRefresh: true,
      );

      state = state.copyWith(
        routes: response.items,
        pagination: response.pagination,
        errorMessage: null,
      );
    } catch (e) {
      // Silent fail
    }
  }

  /// Update state with fresh data from background fetch
  void _updateStateWithFreshData(RouteListResponse freshData) {
    // Create new list instance to ensure Riverpod detects the change
    final newRoutes = List<OptimizedRoute>.from(freshData.items);

    state = state.copyWith(
      routes: newRoutes,
      pagination: freshData.pagination,
      errorMessage: null,
    );
  }

  Future<void> loadMore() async {
    final pagination = state.pagination;
    if (pagination == null || !pagination.hasNext || state.isLoadingMore) {
      return;
    }

    await loadRoutes(page: pagination.page + 1);
  }

  Future<void> deleteRoute(String routeId) async {
    try {
      await _repository.deleteRoute(routeId);
      // Remove from list
      state = state.copyWith(
        routes: state.routes.where((r) => r.id != routeId).toList(),
      );
    } catch (e) {
      state = state.copyWith(
        errorMessage: e.toString().replaceAll('Exception: ', ''),
      );
      rethrow;
    }
  }
}

class RouteOptimizationNotifier extends Notifier<RouteOptimizationState> {
  late final RouteOptimizationRepository _repository;
  // ignore: unused_field
  late final ConnectivityService _connectivity;

  @override
  RouteOptimizationState build() {
    _repository = ref.read(routeOptimizationRepositoryProvider);
    _connectivity = ref.read(connectivityServiceProvider);
    return const RouteOptimizationState();
  }

  Future<void> optimizeRoute(
    OptimizeRouteRequest request, {
    bool useCache = true,
    bool forceRefresh = false,
  }) async {
    state = state.copyWith(
      isOptimizing: true,
      errorMessage: null,
      clearOptimizedRoute: true,
    );

    try {
      // Use optimized repository if available
      final optimizedRoute = switch (_repository) {
        OptimizedRouteOptimizationRepository() =>
          await _repository.optimizeRoute(
            request,
            useCache: useCache,
            forceRefresh: forceRefresh,
          ),
        _ => await _repository.optimizeRoute(request),
      };

      state = state.copyWith(
        isOptimizing: false,
        optimizedRoute: optimizedRoute,
        errorMessage: null,
      );
    } catch (e) {
      state = state.copyWith(
        isOptimizing: false,
        errorMessage: e.toString().replaceAll('Exception: ', ''),
      );
    }
  }

  void clearOptimizedRoute() {
    state = state.copyWith(clearOptimizedRoute: true);
  }
}

class DistanceCalculationNotifier extends Notifier<DistanceCalculationState> {
  late final RouteOptimizationRepository _repository;
  // ignore: unused_field
  late final ConnectivityService _connectivity;

  @override
  DistanceCalculationState build() {
    _repository = ref.read(routeOptimizationRepositoryProvider);
    _connectivity = ref.read(connectivityServiceProvider);
    return const DistanceCalculationState();
  }

  Future<void> calculateDistance(CalculateDistanceRequest request) async {
    state = state.copyWith(
      isCalculating: true,
      errorMessage: null,
      clearResult: true,
    );

    try {
      final result = await _repository.calculateDistance(request);
      state = state.copyWith(
        isCalculating: false,
        distance: result.distance,
        duration: result.duration,
        distanceFormatted: result.distanceFormatted,
        durationFormatted: result.durationFormatted,
        errorMessage: null,
      );
    } catch (e) {
      state = state.copyWith(
        isCalculating: false,
        errorMessage: e.toString().replaceAll('Exception: ', ''),
      );
    }
  }

  void clearResult() {
    state = state.copyWith(clearResult: true);
  }
}
