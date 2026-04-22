import 'package:dio/dio.dart';

import '../../../core/network/connectivity_service.dart';
import '../../../core/storage/offline_storage.dart';
import 'models/deal_model.dart';

class PipelineRepository {
  PipelineRepository(this._dio, this._connectivity);

  final Dio _dio;
  final ConnectivityService _connectivity;

  static const String _dealsPath = '/api/v1/mobile/deals';
  static const String _stagesPath = '/api/v1/mobile/pipelines/stages';
  static const String _formDataPath = '/api/v1/mobile/pipelines/form-data';

  Future<DealListResponse> getDeals({
    int page = 1,
    int perPage = 20,
    String? search,
    String? stageId,
    String? accountId,
    String? assignedTo,
    String? status,
    bool forceRefresh = false,
    Function(DealListResponse)? onBackgroundUpdate,
  }) async {
    // 1. Try to load from cache first (offline-first) - only for first page and no filters
    if (!forceRefresh &&
        page == 1 &&
        (search == null || search.isEmpty) &&
        (stageId == null || stageId.isEmpty) &&
        (accountId == null || accountId.isEmpty) &&
        (assignedTo == null || assignedTo.isEmpty) &&
        (status == null || status.isEmpty)) {
      final cachedDeals = await OfflineStorage.getDeals();
      if (cachedDeals != null && cachedDeals.isNotEmpty) {
        try {
          final deals = cachedDeals.map((json) => Deal.fromJson(json)).toList();
          final cachedResponse = DealListResponse(
            items: deals,
            pagination: Pagination(
              page: 1,
              perPage: deals.length,
              total: deals.length,
              totalPages: 1,
            ),
          );

          // Trigger background refresh if online
          if (_connectivity.isOnline && !forceRefresh) {
            _fetchAndUpdateInBackground(
              page: page,
              perPage: perPage,
              search: search,
              stageId: stageId,
              accountId: accountId,
              assignedTo: assignedTo,
              status: status,
              onBackgroundUpdate: onBackgroundUpdate,
            );
          }

          return cachedResponse;
        } catch (e) {
          // If parsing fails, continue to API call
        }
      }
    }

    // 2. If online, fetch from API
    if (_connectivity.isOnline) {
      try {
        final queryParams = <String, dynamic>{
          'page': page,
          'per_page': perPage,
        };

        if (search != null && search.isNotEmpty) {
          queryParams['search'] = search;
        }
        if (stageId != null && stageId.isNotEmpty) {
          queryParams['stage_id'] = stageId;
        }
        if (accountId != null && accountId.isNotEmpty) {
          queryParams['account_id'] = accountId;
        }
        if (assignedTo != null && assignedTo.isNotEmpty) {
          queryParams['assigned_to'] = assignedTo;
        }
        if (status != null && status.isNotEmpty) {
          queryParams['status'] = status;
        }

        final response = await _dio.get(
          _dealsPath,
          queryParameters: queryParams,
        );

        if (response.data['success'] == true) {
          final dealListResponse = DealListResponse.fromJson(response.data);

          // 3. Save to cache (only for first page and no filters)
          if (page == 1 &&
              (search == null || search.isEmpty) &&
              (stageId == null || stageId.isEmpty) &&
              (accountId == null || accountId.isEmpty) &&
              (assignedTo == null || assignedTo.isEmpty) &&
              (status == null || status.isEmpty)) {
            final dealsJson = dealListResponse.items
                .map((deal) => deal.toJson())
                .toList();
            await OfflineStorage.saveDeals(dealsJson);
          }

          return dealListResponse;
        } else {
          throw Exception(
            response.data['error']?['message'] ?? 'Failed to fetch deals',
          );
        }
      } on DioException catch (e) {
        if (e.response != null) {
          final error = e.response!.data;
          if (error is Map<String, dynamic> && error['error'] != null) {
            throw Exception(
              error['error']['message'] ?? 'Failed to fetch deals',
            );
          }
        }
        throw Exception('Failed to fetch deals: ${e.message}');
      } catch (e) {
        throw Exception('Failed to fetch deals: $e');
      }
    }

    throw Exception('No internet connection and no cached data available');
  }

  /// Fetch deals in background and update cache + UI
  Future<void> _fetchAndUpdateInBackground({
    required int page,
    required int perPage,
    String? search,
    String? stageId,
    String? accountId,
    String? assignedTo,
    String? status,
    Function(DealListResponse)? onBackgroundUpdate,
  }) async {
    try {
      final queryParams = <String, dynamic>{'page': page, 'per_page': perPage};

      if (search != null && search.isNotEmpty) {
        queryParams['search'] = search;
      }
      if (stageId != null && stageId.isNotEmpty) {
        queryParams['stage_id'] = stageId;
      }
      if (accountId != null && accountId.isNotEmpty) {
        queryParams['account_id'] = accountId;
      }
      if (assignedTo != null && assignedTo.isNotEmpty) {
        queryParams['assigned_to'] = assignedTo;
      }
      if (status != null && status.isNotEmpty) {
        queryParams['status'] = status;
      }

      final response = await _dio.get(_dealsPath, queryParameters: queryParams);

      if (response.data['success'] == true) {
        final dealListResponse = DealListResponse.fromJson(response.data);

        // Save to cache
        if (page == 1 &&
            (search == null || search.isEmpty) &&
            (stageId == null || stageId.isEmpty) &&
            (accountId == null || accountId.isEmpty) &&
            (assignedTo == null || assignedTo.isEmpty) &&
            (status == null || status.isEmpty)) {
          final dealsJson = dealListResponse.items
              .map((deal) => deal.toJson())
              .toList();
          await OfflineStorage.saveDeals(dealsJson);
        }

        // Notify UI to update with fresh data
        onBackgroundUpdate?.call(dealListResponse);
      }
    } catch (e) {
      // Silently fail in background
    }
  }

  Future<Deal> getDealById(String id) async {
    final cachedDeal = await OfflineStorage.getDealDetail(id);
    if (cachedDeal != null) {
      try {
        final deal = Deal.fromJson(cachedDeal);
        if (_connectivity.isOnline) {
          _fetchAndUpdateDealDetail(id).catchError((_) {});
        }
        return deal;
      } catch (_) {}
    }

    if (_connectivity.isOnline) {
      try {
        final response = await _dio.get('$_dealsPath/$id');

        if (response.data['success'] == true) {
          final deal = Deal.fromJson(
            response.data['data'] as Map<String, dynamic>,
          );
          await OfflineStorage.saveDealDetail(id, deal.toJson());
          return deal;
        } else {
          throw Exception(
            response.data['error']?['message'] ?? 'Failed to fetch deal detail',
          );
        }
      } on DioException catch (e) {
        if (e.response != null) {
          final error = e.response!.data;
          if (error is Map<String, dynamic> && error['error'] != null) {
            throw Exception(
              error['error']['message'] ?? 'Failed to fetch deal detail',
            );
          }
        }
        throw Exception('Failed to fetch deal detail: ${e.message}');
      } catch (e) {
        throw Exception('Failed to fetch deal detail: $e');
      }
    }

    throw Exception('No internet connection and no cached data available');
  }

  Future<void> _fetchAndUpdateDealDetail(String id) async {
    try {
      final response = await _dio.get('$_dealsPath/$id');
      if (response.data['success'] == true) {
        final deal = Deal.fromJson(
          response.data['data'] as Map<String, dynamic>,
        );
        await OfflineStorage.saveDealDetail(id, deal.toJson());
      }
    } catch (_) {}
  }

  Future<List<PipelineStage>> getStages({bool forceRefresh = false}) async {
    // 1. If online, try to fetch from API first (Network First)
    // This ensures we always get the latest Stage IDs/Order from backend
    if (_connectivity.isOnline) {
      try {
        final response = await _dio.get(_stagesPath);

        if (response.data['success'] == true) {
          final List<dynamic> stagesData = response.data['data'];
          final stages = stagesData
              .map((json) => PipelineStage.fromJson(json))
              .toList();

          // Save to cache for offline usage
          await OfflineStorage.saveStages(
            stages.map((s) => s.toJson()).toList(),
          );
          return stages;
        } else {
          // If server returns explicit error, throw it (don't fallback to cache silently maybe?)
          // Or we can fallback. Let's fallback only on connection/server errors, not logic errors.
          throw Exception(
            response.data['error']?['message'] ??
                'Failed to fetch pipeline stages',
          );
        }
      } catch (e) {
        // If API fails, logging it might be good.
        // Fallthrough to usage of cache
      }
    }

    // 2. Fallback to cache (Offline or API failed)
    final cachedStages = await OfflineStorage.getStages();
    if (cachedStages != null && cachedStages.isNotEmpty) {
      try {
        return cachedStages
            .map((json) => PipelineStage.fromJson(json))
            .toList();
      } catch (_) {}
    }

    // 3. If no cache and offline/failed
    if (!_connectivity.isOnline) {
      throw Exception('No internet connection and no cached data available');
    }

    // If we are here, it means we are online, API failed (caught), and no cache.
    // We should rethrow the API error if possible, but we swallowed it.
    // Let's return empty list or throw generic error.
    throw Exception('Failed to load pipeline stages');
  }

  Future<Deal> createDeal(Map<String, dynamic> data) async {
    try {
      final response = await _dio.post(_dealsPath, data: data);

      if (response.data['success'] == true) {
        return Deal.fromJson(response.data['data'] as Map<String, dynamic>);
      } else {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to create deal',
        );
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final error = e.response!.data;
        if (error is Map<String, dynamic> && error['error'] != null) {
          throw Exception(error['error']['message'] ?? 'Failed to create deal');
        }
      }
      throw Exception('Failed to create deal: ${e.message}');
    } catch (e) {
      throw Exception('Failed to create deal: $e');
    }
  }

  Future<Deal> updateDeal(String id, Map<String, dynamic> data) async {
    try {
      final response = await _dio.put('$_dealsPath/$id', data: data);

      if (response.data['success'] == true) {
        return Deal.fromJson(response.data['data'] as Map<String, dynamic>);
      } else {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to update deal',
        );
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final error = e.response!.data;
        if (error is Map<String, dynamic> && error['error'] != null) {
          throw Exception(error['error']['message'] ?? 'Failed to update deal');
        }
      }
      throw Exception('Failed to update deal: ${e.message}');
    } catch (e) {
      throw Exception('Failed to update deal: $e');
    }
  }

  Future<void> deleteDeal(String id) async {
    try {
      final response = await _dio.delete('$_dealsPath/$id');

      if (response.data['success'] != true) {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to delete deal',
        );
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final error = e.response!.data;
        if (error is Map<String, dynamic> && error['error'] != null) {
          throw Exception(error['error']['message'] ?? 'Failed to delete deal');
        }
      }
      throw Exception('Failed to delete deal: ${e.message}');
    } catch (e) {
      throw Exception('Failed to delete deal: $e');
    }
  }

  /// Get form data for deal creation
  Future<Map<String, dynamic>> getFormData({bool forceRefresh = false}) async {
    try {
      if (!forceRefresh) {
        final accounts = await OfflineStorage.getAccounts();
        final contacts = await OfflineStorage.getContacts();
        final stages = await OfflineStorage.getStages();
        final products = await OfflineStorage.getProducts();

        if (accounts != null &&
            contacts != null &&
            stages != null &&
            products != null) {
          return {
            'accounts': accounts,
            'contacts': contacts,
            'pipeline_stages': stages,
            'products': products,
          };
        }
      }

      final response = await _dio.get(_formDataPath);
      final data = response.data['data'] as Map<String, dynamic>;

      // Cache the individual parts if possible
      if (data['accounts'] != null) {
        await OfflineStorage.saveAccounts(
          List<Map<String, dynamic>>.from(data['accounts']),
        );
      }
      if (data['contacts'] != null) {
        await OfflineStorage.saveContacts(
          List<Map<String, dynamic>>.from(data['contacts']),
        );
      }
      if (data['pipeline_stages'] != null) {
        await OfflineStorage.saveStages(
          List<Map<String, dynamic>>.from(data['pipeline_stages']),
        );
      }
      if (data['products'] != null) {
        await OfflineStorage.saveProducts(
          List<Map<String, dynamic>>.from(data['products']),
        );
      }

      return data;
    } catch (e) {
      // Return cached data if available
      final accounts = await OfflineStorage.getAccounts();
      final contacts = await OfflineStorage.getContacts();
      final stages = await OfflineStorage.getStages();
      final products = await OfflineStorage.getProducts();

      return {
        'accounts': accounts ?? [],
        'contacts': contacts ?? [],
        'pipeline_stages': stages ?? [],
        'products': products ?? [],
      };
    }
  }
}
