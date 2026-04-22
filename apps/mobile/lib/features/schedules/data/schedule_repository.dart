import 'package:dio/dio.dart';

import '../../../core/network/connectivity_service.dart';
import '../../../core/storage/offline_storage.dart';
import '../../../core/network/pagination.dart';
import 'models/schedule_model.dart';

class ScheduleRepository {
  final Dio _dio;
  final ConnectivityService _connectivity;

  ScheduleRepository(this._dio, this._connectivity);

  /// Get schedules dengan Stale-While-Revalidate pattern
  ///
  /// Flow:
  /// 1. Return cached data immediately (if available)
  /// 2. Fetch fresh data in background (if online)
  /// 3. Update UI via callback when fresh data arrives
  Future<ScheduleListResponse> getSchedules({
    int page = 1,
    int perPage = 20,
    String? search,
    String? status,
    DateTime? scheduledFrom,
    DateTime? scheduledTo,
    bool forceRefresh = false,
    Function(ScheduleListResponse)? onBackgroundUpdate,
  }) async {
    // For pagination > 1, always fetch from API (no cache for pages beyond 1)
    if (page > 1) {
      return await _fetchFromApi(
        page: page,
        perPage: perPage,
        search: search,
        status: status,
        scheduledFrom: scheduledFrom,
        scheduledTo: scheduledTo,
      );
    }

    // Get cached data first
    ScheduleListResponse? cachedData;
    if (!forceRefresh &&
        (search == null || search.isEmpty) &&
        status == null &&
        scheduledFrom == null &&
        scheduledTo == null) {
      final cachedSchedules = await OfflineStorage.getSchedules();
      if (cachedSchedules != null && cachedSchedules.isNotEmpty) {
        try {
          final items = cachedSchedules
              .map((json) => ScheduleModel.fromJson(json))
              .toList();
          cachedData = ScheduleListResponse(
            items: items,
            pagination: Pagination(
              page: 1,
              perPage: items.length,
              total: items.length,
              totalPages: 1,
            ),
          );
        } catch (e) {
          // If parsing fails, continue to API call
        }
      }
    }

    // If online, fetch fresh data
    if (_connectivity.isOnline) {
      // If force refresh, always fetch from API
      if (forceRefresh) {
        return await _fetchFromApi(
          page: page,
          perPage: perPage,
          search: search,
          status: status,
          scheduledFrom: scheduledFrom,
          scheduledTo: scheduledTo,
        );
      }

      // Return cached immediately if available
      if (cachedData != null) {
        // Start background fetch (don't await!)
        _fetchAndUpdateInBackground(
          perPage: perPage,
          onUpdate: onBackgroundUpdate,
        );

        return cachedData;
      }

      // No cache, fetch from API
      return await _fetchFromApi(
        page: page,
        perPage: perPage,
        search: search,
        status: status,
        scheduledFrom: scheduledFrom,
        scheduledTo: scheduledTo,
      );
    }

    // Offline mode: return cached or throw error
    if (cachedData != null) {
      return cachedData;
    }

    throw Exception('No internet connection and no cached data available');
  }

  /// Fetch from API and update cache
  Future<ScheduleListResponse> _fetchFromApi({
    required int page,
    required int perPage,
    String? search,
    String? status,
    DateTime? scheduledFrom,
    DateTime? scheduledTo,
  }) async {
    final queryParams = <String, dynamic>{'page': page, 'per_page': perPage};

    if (search != null && search.isNotEmpty) {
      queryParams['search'] = search;
    }
    if (status != null && status.isNotEmpty) {
      queryParams['status'] = status;
    }
    if (scheduledFrom != null) {
      queryParams['scheduled_at_from'] =
          '${scheduledFrom.toUtc().toIso8601String().split('.').first}Z';
    }
    if (scheduledTo != null) {
      queryParams['scheduled_at_to'] =
          '${scheduledTo.toUtc().toIso8601String().split('.').first}Z';
    }

    final response = await _dio.get(
      '/api/v1/schedules',
      queryParameters: queryParams,
    );

    if (response.data is Map<String, dynamic>) {
      final responseData = response.data as Map<String, dynamic>;
      if (responseData['success'] == true) {
        final scheduleListResponse = ScheduleListResponse.fromJson(
          responseData,
        );

        // Save to cache (only for first page and no filters)
        if (page == 1 &&
            (search == null || search.isEmpty) &&
            status == null &&
            scheduledFrom == null &&
            scheduledTo == null) {
          final itemsJson = scheduleListResponse.items
              .map((item) => item.toJson())
              .toList();
          await OfflineStorage.saveSchedules(itemsJson);
        }

        return scheduleListResponse;
      } else {
        throw Exception(
          responseData['error']?['message'] ?? 'Failed to fetch schedules',
        );
      }
    } else {
      throw Exception('Invalid response format');
    }
  }

  /// Background fetch untuk Stale-While-Revalidate pattern
  ///
  /// Method ini dipanggil secara async (tanpa await) sehingga
  /// tidak blocking UI thread. Setelah data di-fetch dan disimpan,
  /// callback onUpdate dipanggil untuk update UI.
  Future<void> _fetchAndUpdateInBackground({
    required int perPage,
    Function(ScheduleListResponse)? onUpdate,
  }) async {
    try {
      final freshData = await _fetchFromApi(page: 1, perPage: perPage);

      // Notify UI untuk update
      if (onUpdate != null) {
        onUpdate(freshData);
      }
    } catch (e) {
      // Silent fail - don't show error for background fetch
      // Cache tetap digunakan, user akan mendapat data fresh saat next visit
    }
  }

  Future<ScheduleModel> getScheduleById(String id) async {
    final cached = await OfflineStorage.getScheduleDetail(id);

    if (_connectivity.isOnline) {
      try {
        final response = await _dio.get('/api/v1/schedules/$id');

        if (response.data is Map<String, dynamic>) {
          final responseData = response.data as Map<String, dynamic>;
          if (responseData['success'] == true) {
            final schedule = ScheduleModel.fromJson(
              responseData['data'] as Map<String, dynamic>,
            );
            await OfflineStorage.saveScheduleDetail(id, schedule.toJson());
            return schedule;
          } else {
            throw Exception(
              responseData['error']?['message'] ?? 'Failed to fetch schedule',
            );
          }
        } else {
          throw Exception('Invalid response format');
        }
      } on DioException catch (e) {
        if (cached != null) return ScheduleModel.fromJson(cached);
        if (e.response != null) {
          final errorData = e.response!.data;
          throw Exception(
            errorData['error']?['message'] ?? 'Failed to fetch schedule',
          );
        } else {
          throw Exception('Network error: ${e.message}');
        }
      } catch (e) {
        if (cached != null) return ScheduleModel.fromJson(cached);
        throw Exception('Failed to fetch schedule: $e');
      }
    }

    if (cached != null) return ScheduleModel.fromJson(cached);
    throw Exception('No internet connection and no cached data available');
  }

  Future<ScheduleModel> createSchedule(ScheduleRequest request) async {
    try {
      final response = await _dio.post(
        '/api/v1/schedules',
        data: request.toJson(),
      );

      if (response.data is Map<String, dynamic>) {
        final responseData = response.data as Map<String, dynamic>;
        if (responseData['success'] == true) {
          return ScheduleModel.fromJson(
            responseData['data'] as Map<String, dynamic>,
          );
        } else {
          throw Exception(
            responseData['error']?['message'] ?? 'Failed to create schedule',
          );
        }
      } else {
        throw Exception('Invalid response format');
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final errorData = e.response!.data;
        throw Exception(
          errorData['error']?['message'] ?? 'Failed to create schedule',
        );
      } else {
        throw Exception('Network error: ${e.message}');
      }
    } catch (e) {
      throw Exception('Failed to create schedule: $e');
    }
  }

  Future<ScheduleModel> updateSchedule(
    String id,
    ScheduleRequest request,
  ) async {
    try {
      final response = await _dio.put(
        '/api/v1/schedules/$id',
        data: request.toJson(),
      );

      if (response.data is Map<String, dynamic>) {
        final responseData = response.data as Map<String, dynamic>;
        if (responseData['success'] == true) {
          return ScheduleModel.fromJson(
            responseData['data'] as Map<String, dynamic>,
          );
        } else {
          throw Exception(
            responseData['error']?['message'] ?? 'Failed to update schedule',
          );
        }
      } else {
        throw Exception('Invalid response format');
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final errorData = e.response!.data;
        throw Exception(
          errorData['error']?['message'] ?? 'Failed to update schedule',
        );
      } else {
        throw Exception('Network error: ${e.message}');
      }
    } catch (e) {
      throw Exception('Failed to update schedule: $e');
    }
  }

  Future<void> deleteSchedule(String id) async {
    try {
      final response = await _dio.delete('/api/v1/schedules/$id');

      if (response.statusCode == 200 || response.statusCode == 204) {
        return;
      }

      if (response.data != null && response.data is Map<String, dynamic>) {
        final responseData = response.data as Map<String, dynamic>;
        if (responseData['success'] == false) {
          throw Exception(
            responseData['error']?['message'] ?? 'Failed to delete schedule',
          );
        }
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final errorData = e.response!.data;
        throw Exception(
          errorData['error']?['message'] ?? 'Failed to delete schedule',
        );
      } else {
        throw Exception('Network error: ${e.message}');
      }
    } catch (e) {
      throw Exception('Failed to delete schedule: $e');
    }
  }

  Future<ScheduleModel> updateStatus(String id, String status) async {
    try {
      final response = await _dio.patch(
        '/api/v1/schedules/$id/status',
        data: {'status': status},
      );

      if (response.data is Map<String, dynamic>) {
        final responseData = response.data as Map<String, dynamic>;
        if (responseData['success'] == true) {
          return ScheduleModel.fromJson(
            responseData['data'] as Map<String, dynamic>,
          );
        } else {
          throw Exception(
            responseData['error']?['message'] ?? 'Failed to update status',
          );
        }
      } else {
        throw Exception('Invalid response format');
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final errorData = e.response!.data;
        throw Exception(
          errorData['error']?['message'] ?? 'Failed to update status',
        );
      } else {
        throw Exception('Network error: ${e.message}');
      }
    } catch (e) {
      throw Exception('Failed to update status: $e');
    }
  }

  Future<ScheduleModel> syncToGoogleCalendar(String id) async {
    try {
      final response = await _dio.post(
        '/api/v1/schedules/$id/sync-google-calendar',
      );

      if (response.data is Map<String, dynamic>) {
        final responseData = response.data as Map<String, dynamic>;
        if (responseData['success'] == true) {
          return ScheduleModel.fromJson(
            responseData['data'] as Map<String, dynamic>,
          );
        } else {
          throw Exception(
            responseData['error']?['message'] ??
                'Failed to sync to Google Calendar',
          );
        }
      } else {
        throw Exception('Invalid response format');
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final errorData = e.response!.data;
        throw Exception(
          errorData['error']?['message'] ?? 'Failed to sync to Google Calendar',
        );
      } else {
        throw Exception('Network error: ${e.message}');
      }
    } catch (e) {
      throw Exception('Failed to sync to Google Calendar: $e');
    }
  }

  Future<ScheduleModel> unsyncFromGoogleCalendar(String id) async {
    try {
      final response = await _dio.post(
        '/api/v1/schedules/$id/unsync-google-calendar',
      );

      if (response.data is Map<String, dynamic>) {
        final responseData = response.data as Map<String, dynamic>;
        if (responseData['success'] == true) {
          return ScheduleModel.fromJson(
            responseData['data'] as Map<String, dynamic>,
          );
        } else {
          throw Exception(
            responseData['error']?['message'] ??
                'Failed to unsync from Google Calendar',
          );
        }
      } else {
        throw Exception('Invalid response format');
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final errorData = e.response!.data;
        throw Exception(
          errorData['error']?['message'] ??
              'Failed to unsync from Google Calendar',
        );
      } else {
        throw Exception('Network error: ${e.message}');
      }
    } catch (e) {
      throw Exception('Failed to unsync from Google Calendar: $e');
    }
  }
}
