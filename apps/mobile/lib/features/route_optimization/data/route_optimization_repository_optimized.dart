import 'dart:async';

import 'package:dio/dio.dart';

import '../../../core/network/connectivity_service.dart';
import 'models/optimized_route.dart';
import 'models/optimize_route_request.dart';
import 'route_optimization_cache.dart';
import 'route_optimization_repository.dart';

/// Optimized route optimization repository with caching and debouncing
class OptimizedRouteOptimizationRepository extends RouteOptimizationRepository {
  OptimizedRouteOptimizationRepository(super.dio, super.connectivity)
    : _connectivity = connectivity;

  final ConnectivityService _connectivity;

  // Request debouncing
  Timer? _optimizeDebounceTimer;
  Completer<OptimizedRoute>? _pendingOptimizeRequest;

  /// Optimize route with caching and debouncing
  @override
  Future<OptimizedRoute> optimizeRoute(
    OptimizeRouteRequest request, {
    bool useCache = true,
    bool forceRefresh = false,
  }) async {
    // Check cache first (if enabled and not forcing refresh)
    if (useCache && !forceRefresh) {
      final cached = await RouteOptimizationCache.getCachedRoute(request);
      if (cached != null) {
        // Return cached result immediately
        // Refresh in background if online
        if (_connectivity.isOnline) {
          _refreshRouteInBackground(request);
        }
        return cached;
      }
    }

    // Debounce rapid requests (wait 500ms for additional requests)
    if (_pendingOptimizeRequest != null &&
        !_pendingOptimizeRequest!.isCompleted) {
      // Cancel previous debounce timer
      _optimizeDebounceTimer?.cancel();

      // Wait for pending request or timeout
      try {
        return await _pendingOptimizeRequest!.future.timeout(
          const Duration(seconds: 30),
        );
      } catch (e) {
        // Timeout or error, continue with new request
      }
    }

    // Create new request
    final completer = Completer<OptimizedRoute>();
    _pendingOptimizeRequest = completer;

    // Debounce: wait 500ms before making request
    _optimizeDebounceTimer = Timer(const Duration(milliseconds: 500), () async {
      try {
        // Check cache again after debounce (might have been cached by another request)
        if (useCache && !forceRefresh) {
          final cached = await RouteOptimizationCache.getCachedRoute(request);
          if (cached != null) {
            completer.complete(cached);
            _pendingOptimizeRequest = null;
            return;
          }
        }

        // Make API request
        final route = await super.optimizeRoute(request);

        // Cache the result
        if (useCache) {
          await RouteOptimizationCache.cacheRoute(request, route);
        }

        completer.complete(route);
      } catch (e) {
        completer.completeError(e);
      } finally {
        _pendingOptimizeRequest = null;
      }
    });

    return completer.future;
  }

  /// Refresh route in background (for cache updates)
  void _refreshRouteInBackground(OptimizeRouteRequest request) {
    // Don't await, run in background
    Future.microtask(() async {
      try {
        final route = await super.optimizeRoute(request);
        await RouteOptimizationCache.cacheRoute(request, route);
      } catch (e) {
        // Ignore background refresh errors
      }
    });
  }

  /// Get my routes with caching
  @override
  Future<RouteListResponse> getMyRoutes({
    int page = 1,
    int perPage = 20,
    bool forceRefresh = false,
    Function(RouteListResponse)? onBackgroundUpdate,
  }) async {
    // Check cache first (if not forcing refresh and page 1)
    if (!forceRefresh && page == 1) {
      final cached = await RouteOptimizationCache.getCachedRouteList(page);
      if (cached != null && cached.isNotEmpty) {
        try {
          final routes = cached
              .map(
                (json) => OptimizedRoute.fromJson(json as Map<String, dynamic>),
              )
              .toList();
          final cachedResponse = RouteListResponse(
            items: routes,
            pagination: Pagination(
              page: 1,
              perPage: routes.length,
              total: routes.length,
              totalPages: 1,
              hasNext: false,
              hasPrev: false,
            ),
          );

          // Trigger background refresh if online
          if (_connectivity.isOnline && !forceRefresh) {
            _fetchAndUpdateRoutesInBackground(
              page: page,
              perPage: perPage,
              onUpdate: onBackgroundUpdate,
            );
          }

          return cachedResponse;
        } catch (e) {
          // If parsing fails, continue to API call
        }
      }
    }

    // Make API request
    final response = await super.getMyRoutes(
      page: page,
      perPage: perPage,
      forceRefresh: forceRefresh,
    );

    // Cache the result (only for first page)
    if (page == 1) {
      await RouteOptimizationCache.cacheRouteList(response.items, page);
    }

    return response;
  }

  /// Refresh route list in background
  void _fetchAndUpdateRoutesInBackground({
    required int page,
    required int perPage,
    Function(RouteListResponse)? onUpdate,
  }) {
    Future.microtask(() async {
      try {
        final response = await super.getMyRoutes(
          page: page,
          perPage: perPage,
          forceRefresh: true,
        );
        await RouteOptimizationCache.cacheRouteList(response.items, page);
        // Notify UI to update with fresh data
        onUpdate?.call(response);
      } catch (e) {
        // Ignore background refresh errors
      }
    });
  }

  /// Get route by ID with caching
  @override
  Future<OptimizedRoute> getRouteById(String routeId) async {
    // Try cache first (route detail cache can be implemented if needed)
    // For now, just use API with retry logic
    return _getRouteByIdWithRetry(routeId);
  }

  /// Get route by ID with retry logic
  Future<OptimizedRoute> _getRouteByIdWithRetry(
    String routeId, {
    int maxRetries = 2,
  }) async {
    int attempts = 0;
    while (attempts < maxRetries) {
      try {
        return await super.getRouteById(routeId);
      } on DioException catch (e) {
        attempts++;

        // Don't retry on client errors (4xx)
        if (e.response?.statusCode != null &&
            e.response!.statusCode! >= 400 &&
            e.response!.statusCode! < 500) {
          rethrow;
        }

        // Retry on server errors or network errors
        if (attempts >= maxRetries) {
          rethrow;
        }

        // Wait before retry (exponential backoff)
        await Future.delayed(Duration(seconds: attempts));
      } catch (e) {
        if (attempts >= maxRetries) {
          rethrow;
        }
        attempts++;
        await Future.delayed(Duration(seconds: attempts));
      }
    }

    throw Exception('Failed to get route after $maxRetries attempts');
  }

  /// Cleanup resources
  void dispose() {
    _optimizeDebounceTimer?.cancel();
    _pendingOptimizeRequest = null;
  }
}
