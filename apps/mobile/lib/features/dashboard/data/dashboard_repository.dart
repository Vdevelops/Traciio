import 'package:dio/dio.dart';

import '../../../core/network/api_client.dart';
import '../../../core/network/connectivity_service.dart';
import '../../../core/storage/offline_storage.dart';
import 'models/dashboard.dart';

class DashboardRepository {
  final Dio _dio;
  final ConnectivityService _connectivity;

  DashboardRepository({Dio? dio, ConnectivityService? connectivity})
    : _dio = dio ?? ApiClient.dio,
      _connectivity = connectivity ?? ConnectivityService();

  /// Fetch mobile dashboard overview (Target Summary, Quick Stats, etc.)
  Future<MobileDashboardOverview> getMobileOverview({
    String? period, // e.g., 'today', 'week', 'month'
    bool forceRefresh = false,
  }) async {
    // 1. Try cache if not forcing refresh
    if (!forceRefresh) {
      final cached = await OfflineStorage.getMobileDashboardOverview(
        period: period,
      );
      if (cached != null) {
        try {
          final overview = MobileDashboardOverview.fromJson(cached);

          // If online, fetch background update
          if (_connectivity.isOnline) {
            _fetchMobileOverview(period: period)
                .then((data) {
                  OfflineStorage.saveMobileDashboardOverview(
                    data.toJson(),
                    period: period,
                  );
                })
                .catchError((_) {
                  // Ignore background fetch errors
                });
          }
          return overview;
        } catch (e) {
          // parsing error, ignore and try network
        }
      }
    }

    // 2. Fetch from network
    if (_connectivity.isOnline) {
      try {
        final overview = await _fetchMobileOverview(period: period);
        await OfflineStorage.saveMobileDashboardOverview(
          overview.toJson(),
          period: period,
        );
        return overview;
      } catch (e) {
        // If network fails, try cache as fallback
        final cached = await OfflineStorage.getMobileDashboardOverview(
          period: period,
        );
        if (cached != null) {
          return MobileDashboardOverview.fromJson(cached);
        }
        rethrow;
      }
    }

    // 3. Fallback to cache if offline
    final cached = await OfflineStorage.getMobileDashboardOverview(
      period: period,
    );
    if (cached != null) {
      return MobileDashboardOverview.fromJson(cached);
    }

    throw Exception('No internet connection and no cached data available');
  }

  Future<MobileDashboardOverview> _fetchMobileOverview({String? period}) async {
    final queryParams = <String, dynamic>{};
    if (period != null) queryParams['period'] = period;

    final response = await _dio.get(
      '/api/v1/dashboard/mobile/overview',
      queryParameters: queryParams.isNotEmpty ? queryParams : null,
    );

    final data = _extractData(response);

    if (data is! Map<String, dynamic>) {
      throw Exception(
        'Invalid overview data format: expected Map, got ${data.runtimeType}',
      );
    }

    return MobileDashboardOverview.fromJson(data);
  }

  /// Fetch mobile dashboard visits
  Future<List<MobileVisit>> getMobileVisits({
    String? period,
    String? status, // 'all', 'planned', 'completed', etc.
    bool forceRefresh = false,
  }) async {
    // 1. Try cache if not forcing refresh
    if (!forceRefresh) {
      final cached = await OfflineStorage.getMobileVisits(period: period);
      if (cached != null) {
        try {
          // cached is List<dynamic>
          final visits = cached
              .map((e) => MobileVisit.fromJson(e as Map<String, dynamic>))
              .toList();

          // Helper to check if we should show cached data based on status filter
          // Note: Simple filtering on cached data if possible, or just return all if 'all'
          // For now, returning cached list. If status is specific, we might need to filter client-side
          // or rely on network for filtered data.
          // Assuming cache stores 'all' visits for the period.

          // Filter cached visits by status
          // Status filter now uses API status directly: draft, submitted, approved, rejected
          List<MobileVisit> filteredVisits = visits;
          if (status != null) {
            filteredVisits = visits
                .where((v) => v.status.toLowerCase() == status.toLowerCase())
                .toList();
          }

          if (_connectivity.isOnline) {
            _fetchMobileVisits(period: period, status: status)
                .then((data) {
                  // Only cache if fetching 'all', or implement smarter caching
                  // For simplicity: saving whatever we fetch might overwrite 'all' with partial data if we are not careful.
                  // Strategy: Always fetch 'all' in background to update cache?
                  // Or just cache what we get and handle it.
                  // Let's assume for now we cache whatever list we get for the period.
                  OfflineStorage.saveMobileVisits(
                    data.map((e) => e.toJson()).toList(),
                    period: period,
                  );
                })
                .catchError((_) {});
          }
          return filteredVisits;
        } catch (e) {
          // Ignore cache errors, fallback to network
        }
      }
    }

    // 2. Fetch from network
    if (_connectivity.isOnline) {
      try {
        final visits = await _fetchMobileVisits(period: period, status: status);
        // Only update cache if we fetched all visits (no status filter)
        // If we fetch filtered status, we shouldn't overwrite the cache with subset.
        // For safety, let's only save to cache if status is null.
        if (status == null) {
          await OfflineStorage.saveMobileVisits(
            visits.map((e) => e.toJson()).toList(),
            period: period,
          );
        }
        return visits;
      } catch (e) {
        final cached = await OfflineStorage.getMobileVisits(period: period);
        if (cached != null) {
          final visits = cached
              .map((e) => MobileVisit.fromJson(e as Map<String, dynamic>))
              .toList();
          if (status != null) {
            // Status filter now uses API status directly: draft, submitted, approved, rejected
            return visits
                .where((v) => v.status.toLowerCase() == status.toLowerCase())
                .toList();
          }
          return visits;
        }
        rethrow;
      }
    }

    // 3. Fallback
    final cached = await OfflineStorage.getMobileVisits(period: period);
    if (cached != null) {
      final visits = cached
          .map((e) => MobileVisit.fromJson(e as Map<String, dynamic>))
          .toList();
      if (status != null) {
        // Status filter now uses API status directly: draft, submitted, approved, rejected
        return visits
                .where((v) => v.status.toLowerCase() == status.toLowerCase())
                .toList();
      }
      return visits;
    }

    throw Exception('No internet connection and no cached data available');
  }

  Future<List<MobileVisit>> _fetchMobileVisits({
    String? period,
    String? status,
  }) async {
    final queryParams = <String, dynamic>{};
    if (period != null) queryParams['period'] = period;
    
    // Status filter now uses API status directly: draft, submitted, approved, rejected
    if (status != null) {
      queryParams['status'] = status;
    }

    // Add limit parameter (max 10 items for dashboard)
    queryParams['limit'] = 10;
    
    final response = await _dio.get(
      '/api/v1/dashboard/mobile/visits',
      queryParameters: queryParams.isNotEmpty ? queryParams : null,
    );

    final data = _extractData(response);

    // API returns {visits: [], total: int, has_more: bool}
    if (data is Map<String, dynamic>) {
      if (data['visits'] is List) {
        return (data['visits'] as List)
            .map((e) => MobileVisit.fromJson(e as Map<String, dynamic>))
            .toList();
      }
      throw Exception(
        'Invalid visits data format: expected visits array in Map, got ${data.runtimeType}',
      );
    }

    // Fallback: if data is directly a List (for backward compatibility)
    if (data is List) {
      return data
          .map((e) => MobileVisit.fromJson(e as Map<String, dynamic>))
          .toList();
    }

    throw Exception(
      'Invalid visits data format: expected Map or List, got ${data.runtimeType}',
    );
  }

  /// Fetch mobile dashboard tasks
  Future<List<MobileTask>> getMobileTasks({
    String? period,
    bool forceRefresh = false,
  }) async {
    if (!forceRefresh) {
      final cached = await OfflineStorage.getMobileTasks(period: period);
      if (cached != null) {
        try {
          final tasks = cached
              .map((e) => MobileTask.fromJson(e as Map<String, dynamic>))
              .toList();
          if (_connectivity.isOnline) {
            _fetchMobileTasks(period: period)
                .then((data) {
                  OfflineStorage.saveMobileTasks(
                    data.map((e) => e.toJson()).toList(),
                    period: period,
                  );
                })
                .catchError((_) {});
          }
          return tasks;
        } catch (e) {
          // Ignore cache errors, fallback to network
        }
      }
    }

    if (_connectivity.isOnline) {
      try {
        final tasks = await _fetchMobileTasks(period: period);
        await OfflineStorage.saveMobileTasks(
          tasks.map((e) => e.toJson()).toList(),
          period: period,
        );
        return tasks;
      } catch (e) {
        final cached = await OfflineStorage.getMobileTasks(period: period);
        if (cached != null) {
          return cached
              .map((e) => MobileTask.fromJson(e as Map<String, dynamic>))
              .toList();
        }
        rethrow;
      }
    }

    final cached = await OfflineStorage.getMobileTasks(period: period);
    if (cached != null) {
      return cached
          .map((e) => MobileTask.fromJson(e as Map<String, dynamic>))
          .toList();
    }

    throw Exception('No internet connection and no cached data available');
  }

  Future<List<MobileTask>> _fetchMobileTasks({String? period}) async {
    final queryParams = <String, dynamic>{
      'limit': 5, // Default limit 5 untuk upcoming tasks
    };
    // Note: Tidak perlu period atau filter karena API hanya mengembalikan upcoming tasks

    final response = await _dio.get(
      '/api/v1/dashboard/mobile/tasks',
      queryParameters: queryParams,
    );

    final data = _extractData(response);

    // API returns {tasks: [], total: int, has_more: bool}
    if (data is Map<String, dynamic>) {
      if (data['tasks'] is List) {
        return (data['tasks'] as List)
            .map((e) => MobileTask.fromJson(e as Map<String, dynamic>))
            .toList();
      }
      throw Exception(
        'Invalid tasks data format: expected tasks array in Map, got ${data.runtimeType}',
      );
    }

    // Fallback: if data is directly a List (for backward compatibility)
    if (data is List) {
      return data
          .map((e) => MobileTask.fromJson(e as Map<String, dynamic>))
          .toList();
    }

    throw Exception(
      'Invalid tasks data format: expected Map or List, got ${data.runtimeType}',
    );
  }

  dynamic _extractData(Response response) {
    if (response.data is Map<String, dynamic>) {
      final responseData = response.data as Map<String, dynamic>;
      if (responseData['success'] == true) {
        return responseData['data'];
      }
      throw Exception(responseData['error']?['message'] ?? 'API Error');
    }
    if (response.data is List) {
      return response.data;
    }
    throw Exception('Invalid response format');
  }
}
