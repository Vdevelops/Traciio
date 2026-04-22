/// Generic cache service for list data with pagination support
class CachedListData<T> {
  final List<T> items;
  final DateTime timestamp;
  final Map<String, dynamic>? metadata; // For pagination, filters, etc.

  CachedListData({
    required this.items,
    required this.timestamp,
    this.metadata,
  });

  bool isExpired(Duration ttl) {
    return DateTime.now().difference(timestamp) > ttl;
  }
}

class ListCache {
  static final ListCache _instance = ListCache._internal();
  factory ListCache() => _instance;
  ListCache._internal();

  final Map<String, CachedListData> _cache = {};

  /// Get cached list data if not expired
  List<T>? get<T>(String key, {Duration? ttl, Map<String, dynamic>? expectedMetadata}) {
    final cached = _cache[key];
    if (cached != null) {
      final cacheTtl = ttl ?? const Duration(seconds: 60);
      
      // Check if expired
      if (cached.isExpired(cacheTtl)) {
        _cache.remove(key);
        return null;
      }

      // Check if metadata matches (for filters, search, etc.)
      if (expectedMetadata != null && cached.metadata != null) {
        final matches = expectedMetadata.entries.every((entry) {
          return cached.metadata?[entry.key] == entry.value;
        });
        if (!matches) {
          return null; // Metadata doesn't match, cache invalid
        }
      }

      return cached.items as List<T>;
    }
    return null;
  }

  /// Get cached metadata
  Map<String, dynamic>? getMetadata(String key) {
    return _cache[key]?.metadata;
  }

  /// Set cached list data
  void set<T>(
    String key,
    List<T> items, {
    Map<String, dynamic>? metadata,
  }) {
    _cache[key] = CachedListData<T>(
      items: items,
      timestamp: DateTime.now(),
      metadata: metadata,
    );
  }

  /// Append items to cached list (for pagination)
  void append<T>(String key, List<T> newItems) {
    final cached = _cache[key];
    if (cached != null) {
      final existingItems = cached.items as List<T>;
      _cache[key] = CachedListData<T>(
        items: [...existingItems, ...newItems],
        timestamp: cached.timestamp, // Keep original timestamp
        metadata: cached.metadata,
      );
    }
  }

  /// Clear all cache
  void clear() {
    _cache.clear();
  }

  /// Clear cache for specific key
  void clearKey(String key) {
    _cache.remove(key);
  }

  /// Clear cache for specific prefix (e.g., all accounts cache)
  void clearPrefix(String prefix) {
    final keysToRemove = _cache.keys
        .where((key) => key.startsWith(prefix))
        .toList();
    for (final key in keysToRemove) {
      _cache.remove(key);
    }
  }

  /// Generate cache key
  static String cacheKey(String type, {
    int? page,
    String? search,
    String? filter,
    Map<String, String>? filters,
  }) {
    final parts = ['list', type];
    
    if (page != null) {
      parts.add('page:$page');
    }
    
    if (search != null && search.isNotEmpty) {
      parts.add('search:${search.toLowerCase()}');
    }
    
    if (filter != null) {
      parts.add('filter:$filter');
    }
    
    if (filters != null && filters.isNotEmpty) {
      final filterParts = filters.entries
          .map((e) => '${e.key}:${e.value}')
          .toList()
        ..sort();
      parts.addAll(filterParts);
    }
    
    return parts.join(':');
  }
}

