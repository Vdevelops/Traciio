import 'package:dio/dio.dart';

import '../../../core/network/connectivity_service.dart';
import 'models/optimized_route.dart';
import 'models/optimize_route_request.dart';
import 'models/calculate_distance_request.dart';
import 'route_optimization_cache.dart';

class RouteOptimizationRepository {
  RouteOptimizationRepository(this._dio, this._connectivity);

  final Dio _dio;
  final ConnectivityService _connectivity;

  /// Optimize route using OSRM routing
  Future<OptimizedRoute> optimizeRoute(OptimizeRouteRequest request) async {
    if (!_connectivity.isOnline) {
      throw Exception('No internet connection. Please check your connection.');
    }

    try {
      final response = await _dio.post(
        '/api/v1/mobile/route-optimization/optimize',
        data: request.toJson(),
      );

      if (response.data['success'] == true) {
        try {
          return OptimizedRoute.fromJson(
            response.data['data'] as Map<String, dynamic>,
          );
        } catch (e) {
          throw Exception(
            'Failed to parse route optimization response: $e. Response: ${response.data}',
          );
        }
      } else {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to optimize route',
        );
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final error = e.response!.data;
        if (error is Map<String, dynamic> && error['error'] != null) {
          throw Exception(
            error['error']['message'] ?? 'Failed to optimize route',
          );
        }
      }
      throw Exception('Failed to optimize route: ${e.message}');
    }
  }

  /// Get my routes (only routes created by logged-in user)
  Future<RouteListResponse> getMyRoutes({
    int page = 1,
    int perPage = 20,
    bool forceRefresh = false,
    Function(RouteListResponse)? onBackgroundUpdate,
  }) async {
    // 1. Try cache first if not forcing refresh and page 1
    if (!forceRefresh && page == 1) {
      final cachedRoutes = await RouteOptimizationCache.getCachedRouteList(
        page,
      );
      if (cachedRoutes != null && cachedRoutes.isNotEmpty) {
        try {
          final routes = cachedRoutes
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

    // 2. Fetch from network
    if (_connectivity.isOnline) {
      try {
        final queryParams = <String, dynamic>{
          'page': page,
          'per_page': perPage,
        };

        final response = await _dio.get(
          '/api/v1/mobile/route-optimization/my-routes',
          queryParameters: queryParams,
        );

        if (response.data['success'] == true) {
          try {
            final routeListResponse = RouteListResponse.fromJson(response.data);

            // 3. Save to cache (only for first page)
            if (page == 1) {
              await RouteOptimizationCache.cacheRouteList(
                routeListResponse.items,
                page,
              );
            }

            return routeListResponse;
          } catch (e) {
            throw Exception(
              'Failed to parse routes response: $e. Response: ${response.data}',
            );
          }
        } else {
          throw Exception(
            response.data['error']?['message'] ?? 'Failed to fetch routes',
          );
        }
      } on DioException catch (e) {
        // If API fails, try cached data if available
        if (page == 1) {
          final cachedRoutes = await RouteOptimizationCache.getCachedRouteList(
            page,
          );
          if (cachedRoutes != null && cachedRoutes.isNotEmpty) {
            try {
              final routes = cachedRoutes
                  .map(
                    (json) =>
                        OptimizedRoute.fromJson(json as Map<String, dynamic>),
                  )
                  .toList();
              return RouteListResponse(
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
            } catch (_) {
              // Ignore parsing errors
            }
          }
        }

        if (e.response != null) {
          final error = e.response!.data;
          if (error is Map<String, dynamic> && error['error'] != null) {
            throw Exception(
              error['error']['message'] ?? 'Failed to fetch routes',
            );
          }
        }
        throw Exception('Failed to fetch routes: ${e.message}');
      }
    }

    // 3. Offline fallback
    if (page == 1) {
      final cachedRoutes = await RouteOptimizationCache.getCachedRouteList(
        page,
      );
      if (cachedRoutes != null && cachedRoutes.isNotEmpty) {
        try {
          final routes = cachedRoutes
              .map(
                (json) => OptimizedRoute.fromJson(json as Map<String, dynamic>),
              )
              .toList();
          return RouteListResponse(
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
        } catch (e) {
          throw Exception('Failed to load cached routes: $e');
        }
      }
    }

    throw Exception('No internet connection and no cached data available');
  }

  /// Fetch routes in background and update cache + UI
  Future<void> _fetchAndUpdateRoutesInBackground({
    required int page,
    required int perPage,
    Function(RouteListResponse)? onUpdate,
  }) async {
    try {
      final queryParams = <String, dynamic>{'page': page, 'per_page': perPage};

      final response = await _dio.get(
        '/api/v1/mobile/route-optimization/my-routes',
        queryParameters: queryParams,
      );

      if (response.data['success'] == true) {
        final routeListResponse = RouteListResponse.fromJson(response.data);

        // Save to cache
        if (page == 1) {
          await RouteOptimizationCache.cacheRouteList(
            routeListResponse.items,
            page,
          );
        }

        // Notify UI to update with fresh data
        onUpdate?.call(routeListResponse);
      }
    } catch (e) {
      // Silently fail in background
    }
  }

  /// Get route by ID (only if belongs to logged-in user)
  Future<OptimizedRoute> getRouteById(String routeId) async {
    if (!_connectivity.isOnline) {
      throw Exception('No internet connection. Please check your connection.');
    }

    try {
      final response = await _dio.get(
        '/api/v1/mobile/route-optimization/route/$routeId',
      );

      if (response.data['success'] == true) {
        try {
          return OptimizedRoute.fromJson(
            response.data['data'] as Map<String, dynamic>,
          );
        } catch (e) {
          throw Exception(
            'Failed to parse route response: $e. Response: ${response.data}',
          );
        }
      } else {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to fetch route',
        );
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final error = e.response!.data;
        if (error is Map<String, dynamic> && error['error'] != null) {
          final errorCode = error['error']['code'] as String?;
          if (errorCode == 'ROUTE_NOT_FOUND') {
            throw Exception('Route not found');
          }
          throw Exception(error['error']['message'] ?? 'Failed to fetch route');
        }
      }
      throw Exception('Failed to fetch route: ${e.message}');
    }
  }

  /// Calculate distance between two points using OSRM routing
  Future<CalculateDistanceResponse> calculateDistance(
    CalculateDistanceRequest request,
  ) async {
    if (!_connectivity.isOnline) {
      throw Exception('No internet connection. Please check your connection.');
    }

    try {
      final response = await _dio.post(
        '/api/v1/mobile/route-optimization/calculate-distance',
        data: request.toJson(),
      );

      if (response.data['success'] == true) {
        try {
          return CalculateDistanceResponse.fromJson(
            response.data['data'] as Map<String, dynamic>,
          );
        } catch (e) {
          throw Exception(
            'Failed to parse distance calculation response: $e. Response: ${response.data}',
          );
        }
      } else {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to calculate distance',
        );
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final error = e.response!.data;
        if (error is Map<String, dynamic> && error['error'] != null) {
          throw Exception(
            error['error']['message'] ?? 'Failed to calculate distance',
          );
        }
      }
      throw Exception('Failed to calculate distance: ${e.message}');
    }
  }

  /// Delete route by ID (only if belongs to logged-in user)
  Future<void> deleteRoute(String routeId) async {
    if (!_connectivity.isOnline) {
      throw Exception('No internet connection. Please check your connection.');
    }

    try {
      final response = await _dio.delete(
        '/api/v1/mobile/route-optimization/route/$routeId',
      );

      if (response.data['success'] != true) {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to delete route',
        );
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final error = e.response!.data;
        if (error is Map<String, dynamic> && error['error'] != null) {
          final errorCode = error['error']['code'] as String?;
          if (errorCode == 'ROUTE_NOT_FOUND') {
            throw Exception('Route not found');
          }
          throw Exception(
            error['error']['message'] ?? 'Failed to delete route',
          );
        }
      }
      throw Exception('Failed to delete route: ${e.message}');
    }
  }
}

class RouteListResponse {
  final List<OptimizedRoute> items;
  final Pagination pagination;

  RouteListResponse({required this.items, required this.pagination});

  factory RouteListResponse.fromJson(Map<String, dynamic> json) {
    return RouteListResponse(
      items: (json['data'] as List<dynamic>)
          .map((e) => OptimizedRoute.fromJson(e as Map<String, dynamic>))
          .toList(),
      pagination: Pagination.fromJson(
        json['meta']?['pagination'] as Map<String, dynamic>,
      ),
    );
  }
}

class Pagination {
  final int page;
  final int perPage;
  final int total;
  final int totalPages;
  final bool hasNext;
  final bool hasPrev;
  final int? nextPage;
  final int? prevPage;

  Pagination({
    required this.page,
    required this.perPage,
    required this.total,
    required this.totalPages,
    required this.hasNext,
    required this.hasPrev,
    this.nextPage,
    this.prevPage,
  });

  factory Pagination.fromJson(Map<String, dynamic> json) {
    return Pagination(
      page: json['page'] as int,
      perPage: json['per_page'] as int,
      total: json['total'] as int,
      totalPages: json['total_pages'] as int,
      hasNext: json['has_next'] as bool? ?? false,
      hasPrev: json['has_prev'] as bool? ?? false,
      nextPage: json['next_page'] as int?,
      prevPage: json['prev_page'] as int?,
    );
  }
}
