import 'dart:convert';

import '../../../core/storage/offline_storage.dart';
import 'models/optimized_route.dart';
import 'models/optimize_route_request.dart';

/// Cache for route optimization results
/// Provides in-memory and persistent caching for better performance
class RouteOptimizationCache {
  static const String _cacheBox = 'route_optimization_cache';
  static const Duration _defaultTtl = Duration(hours: 1);
  static const Duration _routeListTtl = Duration(minutes: 5);

  // In-memory cache for faster access
  static final Map<String, _CachedRoute> _memoryCache = {};

  /// Initialize cache box
  static Future<void> init() async {
    await OfflineStorage.set(_cacheBox, {});
  }

  /// Generate cache key from request
  static String _generateCacheKey(OptimizeRouteRequest request) {
    // Create hash from waypoints (sorted for consistency)
    final waypoints = List<Map<String, dynamic>>.from(
      request.waypoints.map((w) => w.toJson()),
    );
    
    // Sort waypoints by lat/lng for consistent hash
    waypoints.sort((a, b) {
      if (a['lat'] != b['lat']) {
        return (a['lat'] as num).compareTo(b['lat'] as num);
      }
      return (a['lng'] as num).compareTo(b['lng'] as num);
    });

    final keyData = {
      'start_lat': request.startLocation.lat,
      'start_lng': request.startLocation.lng,
      'waypoints': waypoints,
      if (request.startTime != null) 'start_time': request.startTime!.toIso8601String(),
      if (request.optimizationType != null) 'optimization_type': request.optimizationType,
    };

    // Simple hash function
    final keyString = jsonEncode(keyData);
    return 'route_${keyString.hashCode}';
  }

  /// Get cached route result
  static Future<OptimizedRoute?> getCachedRoute(
    OptimizeRouteRequest request, {
    Duration? ttl,
  }) async {
    final cacheKey = _generateCacheKey(request);
    final cacheTtl = ttl ?? _defaultTtl;

    // Check memory cache first
    final memoryCached = _memoryCache[cacheKey];
    if (memoryCached != null && !memoryCached.isExpired(cacheTtl)) {
      return memoryCached.route;
    }

    // Check persistent cache
    try {
      final cached = await OfflineStorage.get<Map<String, dynamic>>(
        'route_cache_$cacheKey',
        (json) => json,
      );

      if (cached != null) {
        final timestamp = DateTime.parse(cached['cached_at'] as String);
        if (DateTime.now().difference(timestamp) < cacheTtl) {
          final route = OptimizedRoute.fromJson(
            cached['route'] as Map<String, dynamic>,
          );
          
          // Update memory cache
          _memoryCache[cacheKey] = _CachedRoute(route, timestamp);
          
          return route;
        }
      }
    } catch (e) {
      // Cache error, continue
    }

    return null;
  }

  /// Cache route result
  static Future<void> cacheRoute(
    OptimizeRouteRequest request,
    OptimizedRoute route,
  ) async {
    final cacheKey = _generateCacheKey(request);
    final timestamp = DateTime.now();

    // Update memory cache
    _memoryCache[cacheKey] = _CachedRoute(route, timestamp);

    // Update persistent cache
    try {
      await OfflineStorage.set(
        'route_cache_$cacheKey',
        {
          'route': route.toJson(),
          'cached_at': timestamp.toIso8601String(),
        },
      );
    } catch (e) {
      // Cache error, continue (memory cache still works)
    }
  }

  /// Clear all route cache
  static Future<void> clearAll() async {
    _memoryCache.clear();
    // Note: Persistent cache will be cleared on next app restart
    // or can be cleared manually if needed
  }

  /// Clear expired cache entries
  static Future<void> clearExpired({Duration? ttl}) async {
    final cacheTtl = ttl ?? _defaultTtl;
    
    // Clear expired memory cache
    _memoryCache.removeWhere((key, cached) => cached.isExpired(cacheTtl));
  }

  /// Cache route list
  static Future<void> cacheRouteList(
    List<OptimizedRoute> routes,
    int page,
  ) async {
    final cacheKey = 'route_list_page_$page';
    final timestamp = DateTime.now();

    try {
      await OfflineStorage.set(
        cacheKey,
        {
          'routes': routes.map((r) => r.toJson()).toList(),
          'page': page,
          'cached_at': timestamp.toIso8601String(),
        },
      );
    } catch (e) {
      // Cache error, continue
    }
  }

  /// Get cached route list
  static Future<List<OptimizedRoute>?> getCachedRouteList(
    int page, {
    Duration? ttl,
  }) async {
    final cacheTtl = ttl ?? _routeListTtl;
    final cacheKey = 'route_list_page_$page';

    try {
      final cached = await OfflineStorage.get<Map<String, dynamic>>(
        cacheKey,
        (json) => json,
      );

      if (cached != null) {
        final timestamp = DateTime.parse(cached['cached_at'] as String);
        if (DateTime.now().difference(timestamp) < cacheTtl) {
          final routes = (cached['routes'] as List<dynamic>)
              .map((r) => OptimizedRoute.fromJson(r as Map<String, dynamic>))
              .toList();
          return routes;
        }
      }
    } catch (e) {
      // Cache error, continue
    }

    return null;
  }
}

/// Internal cached route data
class _CachedRoute {
  final OptimizedRoute route;
  final DateTime timestamp;

  _CachedRoute(this.route, this.timestamp);

  bool isExpired(Duration ttl) {
    return DateTime.now().difference(timestamp) > ttl;
  }
}

