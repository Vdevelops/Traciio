import 'package:dio/dio.dart';

import '../../../core/network/connectivity_service.dart';
import '../../../core/storage/offline_storage.dart';
import 'models/lead_model.dart';

export 'models/lead_model.dart' show LeadListResponse, Pagination;

class LeadRepository {
  final Dio _dio;
  final ConnectivityService _connectivity;

  LeadRepository(this._dio, this._connectivity);

  /// Get leads dengan Stale-While-Revalidate pattern
  ///
  /// Flow:
  /// 1. Return cached data immediately (if available)
  /// 2. Fetch fresh data in background (if online)
  /// 3. Update UI via callback when fresh data arrives
  Future<LeadListResponse> getLeads({
    int page = 1,
    int perPage = 10,
    String? search,
    String? status,
    String? source,
    String? industry,
    String? province,
    bool forceRefresh = false,
    Function(LeadListResponse)? onBackgroundUpdate,
  }) async {
    // For pagination > 1, always fetch from API (no cache for pages beyond 1)
    if (page > 1) {
      return await _fetchFromApi(
        page: page,
        perPage: perPage,
        search: search,
        status: status,
        source: source,
        industry: industry,
        province: province,
      );
    }

    // Get cached data first
    LeadListResponse? cachedData;
    if (!forceRefresh &&
        (search == null || search.isEmpty) &&
        status == null &&
        source == null &&
        industry == null &&
        province == null) {
      final cachedLeads = await OfflineStorage.getLeads();
      if (cachedLeads != null && cachedLeads.isNotEmpty) {
        try {
          final leads = cachedLeads.map((json) => Lead.fromJson(json)).toList();
          cachedData = LeadListResponse(
            items: leads,
            pagination: Pagination(
              page: 1,
              perPage: leads.length,
              total: leads.length,
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
          source: source,
          industry: industry,
          province: province,
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
        source: source,
        industry: industry,
        province: province,
      );
    }

    // Offline mode: return cached or throw error
    if (cachedData != null) {
      return cachedData;
    }

    throw Exception('No internet connection and no cached data available');
  }

  /// Fetch from API and update cache
  Future<LeadListResponse> _fetchFromApi({
    required int page,
    required int perPage,
    String? search,
    String? status,
    String? source,
    String? industry,
    String? province,
  }) async {
    final queryParams = <String, dynamic>{'page': page, 'per_page': perPage};

    if (search != null && search.isNotEmpty) queryParams['search'] = search;
    if (status != null && status.isNotEmpty) queryParams['status'] = status;
    if (source != null && source.isNotEmpty) {
      queryParams['lead_source'] = source;
    }
    if (industry != null && industry.isNotEmpty) {
      queryParams['industry'] = industry;
    }
    if (province != null && province.isNotEmpty) {
      queryParams['province'] = province;
    }

    final response = await _dio.get(
      '/api/v1/mobile/leads',
      queryParameters: queryParams,
    );

    if (response.data['success'] == true) {
      final leadListResponse = LeadListResponse.fromJson(response.data);

      // Save to cache (only for first page and no filters)
      if (page == 1 &&
          (search == null || search.isEmpty) &&
          status == null &&
          source == null &&
          industry == null &&
          province == null) {
        final leadsJson = leadListResponse.items
            .map((lead) => lead.toJson())
            .toList();
        await OfflineStorage.saveLeads(leadsJson);
      }

      return leadListResponse;
    } else {
      throw Exception(
        response.data['error']?['message'] ?? 'Failed to fetch leads',
      );
    }
  }

  /// Background fetch untuk Stale-While-Revalidate pattern
  Future<void> _fetchAndUpdateInBackground({
    required int perPage,
    Function(LeadListResponse)? onUpdate,
  }) async {
    try {
      final freshData = await _fetchFromApi(page: 1, perPage: perPage);

      // Notify UI untuk update
      if (onUpdate != null) {
        onUpdate(freshData);
      }
    } catch (e) {
      // Silent fail - don't show error for background fetch
    }
  }

  Future<LeadFormData> getFormData() async {
    try {
      final response = await _dio.get('/api/v1/mobile/leads/form-data');

      if (response.data['success'] == true) {
        return LeadFormData.fromJson(response.data['data']);
      } else {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to fetch form data',
        );
      }
    } on DioException catch (e) {
      throw _handleDioError(e);
    } catch (e) {
      throw Exception('Failed to fetch form data: $e');
    }
  }

  Future<Lead> getLeadById(String id) async {
    final cached = await OfflineStorage.getLeadDetail(id);

    if (_connectivity.isOnline) {
      try {
        final response = await _dio.get('/api/v1/mobile/leads/$id');

        if (response.data['success'] == true) {
          final lead = Lead.fromJson(response.data['data']);
          await OfflineStorage.saveLeadDetail(id, lead.toJson());
          return lead;
        } else {
          throw Exception(
            response.data['error']?['message'] ?? 'Failed to fetch lead',
          );
        }
      } on DioException catch (e) {
        if (cached != null) return Lead.fromJson(cached);
        throw _handleDioError(e);
      } catch (e) {
        if (cached != null) return Lead.fromJson(cached);
        throw Exception('Failed to fetch lead: $e');
      }
    }

    if (cached != null) return Lead.fromJson(cached);
    throw Exception('No internet connection and no cached data available');
  }

  Future<Lead> createLead(Map<String, dynamic> data) async {
    try {
      final response = await _dio.post('/api/v1/mobile/leads', data: data);

      if (response.data['success'] == true) {
        return Lead.fromJson(response.data['data']);
      } else {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to create lead',
        );
      }
    } on DioException catch (e) {
      throw _handleDioError(e);
    } catch (e) {
      throw Exception('Failed to create lead: $e');
    }
  }

  Future<Lead> updateLead(String id, Map<String, dynamic> data) async {
    try {
      final response = await _dio.put('/api/v1/mobile/leads/$id', data: data);

      if (response.data['success'] == true) {
        return Lead.fromJson(response.data['data']);
      } else {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to update lead',
        );
      }
    } on DioException catch (e) {
      throw _handleDioError(e);
    } catch (e) {
      throw Exception('Failed to update lead: $e');
    }
  }

  Future<void> deleteLead(String id) async {
    try {
      final response = await _dio.delete('/api/v1/mobile/leads/$id');

      if (response.data['success'] != true) {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to delete lead',
        );
      }
    } on DioException catch (e) {
      throw _handleDioError(e);
    } catch (e) {
      throw Exception('Failed to delete lead: $e');
    }
  }

  Future<void> convertLead(String id, Map<String, dynamic> data) async {
    try {
      final response = await _dio.post(
        '/api/v1/mobile/leads/$id/convert',
        data: data,
      );

      if (response.data['success'] != true) {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to convert lead',
        );
      }
    } on DioException catch (e) {
      throw _handleDioError(e);
    } catch (e) {
      throw Exception('Failed to convert lead: $e');
    }
  }

  Exception _handleDioError(DioException e) {
    if (e.response != null) {
      final errorData = e.response!.data;
      return Exception(errorData['error']?['message'] ?? 'Network error');
    }
    return Exception('Network error: ${e.message}');
  }
}
