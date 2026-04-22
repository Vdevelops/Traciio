class CachedData {
  final dynamic data;
  final DateTime timestamp;

  CachedData(this.data, this.timestamp);

  bool isExpired(Duration ttl) {
    return DateTime.now().difference(timestamp) > ttl;
  }
}

class DashboardCache {
  static final DashboardCache _instance = DashboardCache._internal();
  factory DashboardCache() => _instance;
  DashboardCache._internal();

  final Map<String, CachedData> _cache = {};

  /// Get cached data if not expired
  T? get<T>(String key, {Duration? ttl}) {
    final cached = _cache[key];
    if (cached != null) {
      final cacheTtl = ttl ?? const Duration(seconds: 30);
      if (!cached.isExpired(cacheTtl)) {
        return cached.data as T;
      }
      // Remove expired cache
      _cache.remove(key);
    }
    return null;
  }

  /// Set cached data
  void set(String key, dynamic data) {
    _cache[key] = CachedData(data, DateTime.now());
  }

  /// Clear all cache
  void clear() {
    _cache.clear();
  }

  /// Clear cache for specific period
  void clearPeriod(String period) {
    final keysToRemove = _cache.keys
        .where((key) => key.contains('period:$period'))
        .toList();
    for (final key in keysToRemove) {
      _cache.remove(key);
    }
  }

  /// Generate cache key
  static String cacheKey(String type, String period) {
    return 'dashboard:$type:period:$period';
  }
}

