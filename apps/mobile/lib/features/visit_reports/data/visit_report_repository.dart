import 'dart:io';

import 'package:dio/dio.dart';

import '../../../core/network/connectivity_service.dart';
import '../../../core/storage/offline_storage.dart';
import 'models/visit_report.dart';

class VisitReportRepository {
  VisitReportRepository(this._dio, this._connectivity);

  final Dio _dio;
  final ConnectivityService _connectivity;

  Future<VisitReportListResponse> getVisitReports({
    int page = 1,
    int perPage = 20,
    String? search,
    String? status,
    String? accountId,
    String? salesRepId,
    String? startDate,
    String? endDate,
    bool forceRefresh = false,
    bool forRouteOptimization = false,
    Function(VisitReportListResponse)? onBackgroundUpdate,
  }) async {
    // 1. Try to load from cache first (offline-first) - only for first page and no filters
    if (!forceRefresh &&
        page == 1 &&
        (search == null || search.isEmpty) &&
        status == null &&
        accountId == null &&
        startDate == null &&
        endDate == null) {
      final cachedVisitReports = await OfflineStorage.getVisitReports();
      if (cachedVisitReports != null && cachedVisitReports.isNotEmpty) {
        try {
          final visitReports = cachedVisitReports
              .map((json) => VisitReport.fromJson(json))
              .toList();
          final cachedResponse = VisitReportListResponse(
            items: visitReports,
            pagination: Pagination(
              page: 1,
              perPage: visitReports.length,
              total: visitReports.length,
              totalPages: 1,
            ),
          );

          // Trigger background refresh if online
          if (_connectivity.isOnline && !forceRefresh) {
            _fetchAndUpdateInBackground(
              page: page,
              perPage: perPage,
              search: search,
              status: status,
              accountId: accountId,
              salesRepId: salesRepId,
              startDate: startDate,
              endDate: endDate,
              forRouteOptimization: forRouteOptimization,
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
        if (status != null && status.isNotEmpty) {
          queryParams['status'] = status;
        }
        if (accountId != null && accountId.isNotEmpty) {
          queryParams['account_id'] = accountId;
        }
        if (salesRepId != null && salesRepId.isNotEmpty) {
          queryParams['sales_rep_id'] = salesRepId;
        }
        if (startDate != null && startDate.isNotEmpty) {
          queryParams['start_date'] = startDate;
        }
        if (endDate != null && endDate.isNotEmpty) {
          queryParams['end_date'] = endDate;
        }

        // Always use mobile endpoint
        // For route optimization, add query parameter to get all visit reports
        if (forRouteOptimization) {
          queryParams['for_route_optimization'] = 'true';
        }

        final response = await _dio.get(
          '/api/v1/mobile/visit-reports/my-visit-reports',
          queryParameters: queryParams,
        );

        if (response.data['success'] == true) {
          try {
            final visitReportListResponse = VisitReportListResponse.fromJson(
              response.data,
            );

            // 3. Save to cache (only for first page and no filters)
            if (page == 1 &&
                (search == null || search.isEmpty) &&
                status == null &&
                accountId == null &&
                startDate == null &&
                endDate == null) {
              final visitReportsJson = visitReportListResponse.items
                  .map((report) => report.toJson())
                  .toList();
              await OfflineStorage.saveVisitReports(visitReportsJson);
            }

            return visitReportListResponse;
          } catch (e) {
            throw Exception(
              'Failed to parse visit reports response: $e. Response: ${response.data}',
            );
          }
        } else {
          throw Exception(
            response.data['error']?['message'] ??
                'Failed to fetch visit reports',
          );
        }
      } on DioException catch (e) {
        // If API fails, try to return cached data if available
        if (page == 1 &&
            (search == null || search.isEmpty) &&
            status == null &&
            accountId == null &&
            startDate == null &&
            endDate == null) {
          final cachedVisitReports = await OfflineStorage.getVisitReports();
          if (cachedVisitReports != null && cachedVisitReports.isNotEmpty) {
            try {
              final visitReports = cachedVisitReports
                  .map((json) => VisitReport.fromJson(json))
                  .toList();
              return VisitReportListResponse(
                items: visitReports,
                pagination: Pagination(
                  page: 1,
                  perPage: visitReports.length,
                  total: visitReports.length,
                  totalPages: 1,
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
              error['error']['message'] ?? 'Failed to fetch visit reports',
            );
          }
        }
        throw Exception('Failed to fetch visit reports: ${e.message}');
      } catch (e) {
        // If other error, try cached data
        if (page == 1 &&
            (search == null || search.isEmpty) &&
            status == null &&
            accountId == null &&
            startDate == null &&
            endDate == null) {
          final cachedVisitReports = await OfflineStorage.getVisitReports();
          if (cachedVisitReports != null && cachedVisitReports.isNotEmpty) {
            try {
              final visitReports = cachedVisitReports
                  .map((json) => VisitReport.fromJson(json))
                  .toList();
              return VisitReportListResponse(
                items: visitReports,
                pagination: Pagination(
                  page: 1,
                  perPage: visitReports.length,
                  total: visitReports.length,
                  totalPages: 1,
                ),
              );
            } catch (_) {
              // Ignore parsing errors
            }
          }
        }
        throw Exception('Failed to fetch visit reports: $e');
      }
    }

    // 4. Offline: return cached data or throw error
    if (page == 1 &&
        (search == null || search.isEmpty) &&
        status == null &&
        accountId == null &&
        startDate == null &&
        endDate == null) {
      final cachedVisitReports = await OfflineStorage.getVisitReports();
      if (cachedVisitReports != null && cachedVisitReports.isNotEmpty) {
        try {
          final visitReports = cachedVisitReports
              .map((json) => VisitReport.fromJson(json))
              .toList();
          return VisitReportListResponse(
            items: visitReports,
            pagination: Pagination(
              page: 1,
              perPage: visitReports.length,
              total: visitReports.length,
              totalPages: 1,
            ),
          );
        } catch (e) {
          throw Exception('Failed to load cached visit reports: $e');
        }
      }
    }

    throw Exception('No internet connection and no cached data available');
  }

  /// Fetch visit reports in background and update cache + UI
  Future<void> _fetchAndUpdateInBackground({
    required int page,
    required int perPage,
    String? search,
    String? status,
    String? accountId,
    String? salesRepId,
    String? startDate,
    String? endDate,
    bool forRouteOptimization = false,
    Function(VisitReportListResponse)? onBackgroundUpdate,
  }) async {
    try {
      final queryParams = <String, dynamic>{'page': page, 'per_page': perPage};

      if (search != null && search.isNotEmpty) {
        queryParams['search'] = search;
      }
      if (status != null && status.isNotEmpty) {
        queryParams['status'] = status;
      }
      if (accountId != null && accountId.isNotEmpty) {
        queryParams['account_id'] = accountId;
      }
      if (salesRepId != null && salesRepId.isNotEmpty) {
        queryParams['sales_rep_id'] = salesRepId;
      }
      if (startDate != null && startDate.isNotEmpty) {
        queryParams['start_date'] = startDate;
      }
      if (endDate != null && endDate.isNotEmpty) {
        queryParams['end_date'] = endDate;
      }

      // Always use mobile endpoint
      if (forRouteOptimization) {
        queryParams['for_route_optimization'] = 'true';
      }

      final response = await _dio.get(
        '/api/v1/mobile/visit-reports/my-visit-reports',
        queryParameters: queryParams,
      );

      if (response.data['success'] == true) {
        final visitReportListResponse = VisitReportListResponse.fromJson(
          response.data,
        );

        // Save to cache
        if (page == 1 &&
            (search == null || search.isEmpty) &&
            status == null &&
            accountId == null &&
            startDate == null &&
            endDate == null) {
          final visitReportsJson = visitReportListResponse.items
              .map((report) => report.toJson())
              .toList();
          await OfflineStorage.saveVisitReports(visitReportsJson);
        }

        // Notify UI to update with fresh data
        onBackgroundUpdate?.call(visitReportListResponse);
      }
    } catch (e) {
      // Silently fail in background
    }
  }

  Future<VisitReport> getVisitReportById(String id) async {
    // 1. Try to load from cache first (offline-first)
    final cachedVisitReport = await OfflineStorage.getVisitReportDetail(id);
    if (cachedVisitReport != null) {
      try {
        final visitReport = VisitReport.fromJson(cachedVisitReport);
        // If online, fetch from API in background to update cache
        if (_connectivity.isOnline) {
          _fetchAndUpdateVisitReportDetail(id).catchError((_) {
            // Ignore errors, use cached data
          });
        }
        return visitReport;
      } catch (e) {
        // If parsing fails, continue to API call
      }
    }

    // 2. If online, fetch from API (use mobile endpoint for ownership validation)
    if (_connectivity.isOnline) {
      try {
        final response = await _dio.get('/api/v1/mobile/visit-reports/$id');

        if (response.data['success'] == true) {
          final visitReport = VisitReport.fromJson(
            response.data['data'] as Map<String, dynamic>,
          );

          // 3. Save to cache
          await OfflineStorage.saveVisitReportDetail(id, visitReport.toJson());

          return visitReport;
        } else {
          throw Exception(
            response.data['error']?['message'] ??
                'Failed to fetch visit report',
          );
        }
      } on DioException catch (e) {
        // If API fails, try to return cached data if available
        if (cachedVisitReport != null) {
          try {
            return VisitReport.fromJson(cachedVisitReport);
          } catch (_) {
            // Ignore parsing errors
          }
        }

        if (e.response != null) {
          final error = e.response!.data;
          if (error is Map<String, dynamic> && error['error'] != null) {
            throw Exception(
              error['error']['message'] ?? 'Failed to fetch visit report',
            );
          }
        }
        throw Exception('Failed to fetch visit report: ${e.message}');
      } catch (e) {
        // If other error, try cached data
        if (cachedVisitReport != null) {
          try {
            return VisitReport.fromJson(cachedVisitReport);
          } catch (_) {
            // Ignore parsing errors
          }
        }
        throw Exception('Failed to fetch visit report: $e');
      }
    }

    // 4. Offline: return cached data or throw error
    if (cachedVisitReport != null) {
      try {
        return VisitReport.fromJson(cachedVisitReport);
      } catch (e) {
        throw Exception('Failed to load cached visit report: $e');
      }
    }

    throw Exception('No internet connection and no cached data available');
  }

  /// Clear visit report detail cache
  Future<void> clearVisitReportDetailCache(String id) async {
    await OfflineStorage.clearVisitReportDetail(id);
  }

  /// Save visit report detail to cache
  Future<void> saveVisitReportDetailToCache(
    String id,
    VisitReport visitReport,
  ) async {
    await OfflineStorage.saveVisitReportDetail(id, visitReport.toJson());
  }

  /// Fetch visit report detail from API and update cache (background operation)
  Future<void> _fetchAndUpdateVisitReportDetail(String id) async {
    try {
      final response = await _dio.get('/api/v1/mobile/visit-reports/$id');
      if (response.data['success'] == true) {
        final visitReport = VisitReport.fromJson(
          response.data['data'] as Map<String, dynamic>,
        );
        await OfflineStorage.saveVisitReportDetail(id, visitReport.toJson());
      }
    } catch (_) {
      // Ignore errors in background operation
    }
  }

  /// Get form data for visit report creation (accounts, contacts, deals, leads)
  Future<Map<String, dynamic>> getFormData() async {
    try {
      final response = await _dio.get('/api/v1/mobile/visit-reports/form-data');

      if (response.data['success'] == true) {
        return response.data['data'] as Map<String, dynamic>;
      } else {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to fetch form data',
        );
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final error = e.response!.data;
        if (error is Map<String, dynamic> && error['error'] != null) {
          throw Exception(
            error['error']['message'] ?? 'Failed to fetch form data',
          );
        }
      }
      throw Exception('Failed to fetch form data: ${e.message}');
    } catch (e) {
      throw Exception('Failed to fetch form data: $e');
    }
  }

  Future<VisitReport> createVisitReport({
    String? accountId,
    String? contactId,
    String? dealId,
    String? leadId,
    required String visitDate,
    required String purpose,
    String? notes,
  }) async {
    try {
      // Use mobile endpoint - auto-sets sales_rep_id from authenticated user
      final response = await _dio.post(
        '/api/v1/mobile/visit-reports',
        data: {
          if (accountId != null && accountId.isNotEmpty)
            'account_id': accountId,
          if (contactId != null && contactId.isNotEmpty)
            'contact_id': contactId,
          if (dealId != null && dealId.isNotEmpty) 'deal_id': dealId,
          if (leadId != null && leadId.isNotEmpty) 'lead_id': leadId,
          'visit_date': visitDate,
          'purpose': purpose,
          if (notes != null && notes.isNotEmpty) 'notes': notes,
        },
      );

      if (response.data['success'] == true) {
        return VisitReport.fromJson(
          response.data['data'] as Map<String, dynamic>,
        );
      } else {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to create visit report',
        );
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final error = e.response!.data;
        if (error is Map<String, dynamic> && error['error'] != null) {
          throw Exception(
            error['error']['message'] ?? 'Failed to create visit report',
          );
        }
      }
      throw Exception('Failed to create visit report: ${e.message}');
    } catch (e) {
      throw Exception('Failed to create visit report: $e');
    }
  }

  Future<VisitReport> updateVisitReport({
    required String id,
    String? accountId,
    String? contactId,
    String? dealId,
    String? leadId,
    String? visitDate,
    String? purpose,
    String? notes,
  }) async {
    try {
      // Use mobile endpoint - validates ownership and status (only draft)
      final response = await _dio.put(
        '/api/v1/mobile/visit-reports/$id',
        data: {
          if (accountId != null && accountId.isNotEmpty)
            'account_id': accountId,
          if (contactId != null && contactId.isNotEmpty)
            'contact_id': contactId,
          if (dealId != null && dealId.isNotEmpty) 'deal_id': dealId,
          if (leadId != null && leadId.isNotEmpty) 'lead_id': leadId,
          if (visitDate != null && visitDate.isNotEmpty)
            'visit_date': visitDate,
          if (purpose != null && purpose.isNotEmpty) 'purpose': purpose,
          if (notes != null && notes.isNotEmpty) 'notes': notes,
        },
      );

      if (response.data['success'] == true) {
        return VisitReport.fromJson(
          response.data['data'] as Map<String, dynamic>,
        );
      } else {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to update visit report',
        );
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final error = e.response!.data;
        if (error is Map<String, dynamic> && error['error'] != null) {
          throw Exception(
            error['error']['message'] ?? 'Failed to update visit report',
          );
        }
      }
      throw Exception('Failed to update visit report: ${e.message}');
    } catch (e) {
      throw Exception('Failed to update visit report: $e');
    }
  }

  Future<void> deleteVisitReport(String id) async {
    try {
      // Use mobile endpoint - validates ownership
      await _dio.delete('/api/v1/mobile/visit-reports/$id');
    } on DioException catch (e) {
      if (e.response != null) {
        final error = e.response!.data;
        if (error is Map<String, dynamic> && error['error'] != null) {
          throw Exception(
            error['error']['message'] ?? 'Failed to delete visit report',
          );
        }
      }
      throw Exception('Failed to delete visit report: ${e.message}');
    } catch (e) {
      throw Exception('Failed to delete visit report: $e');
    }
  }

  Future<VisitReport> submitVisitReport({
    required String id,
    String? outcome,
    String? nextSteps,
  }) async {
    try {
      // Use mobile endpoint - validates ownership and status (only draft, must be checked in and out)
      final response = await _dio.patch(
        '/api/v1/mobile/visit-reports/$id/submit',
        data: {
          if (outcome != null && outcome.isNotEmpty) 'outcome': outcome,
          if (nextSteps != null && nextSteps.isNotEmpty)
            'next_steps': nextSteps,
        },
      );

      if (response.data['success'] == true) {
        return VisitReport.fromJson(
          response.data['data'] as Map<String, dynamic>,
        );
      } else {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to submit visit report',
        );
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final error = e.response!.data;
        if (error is Map<String, dynamic> && error['error'] != null) {
          throw Exception(
            error['error']['message'] ?? 'Failed to submit visit report',
          );
        }
      }
      throw Exception('Failed to submit visit report: ${e.message}');
    } catch (e) {
      throw Exception('Failed to submit visit report: $e');
    }
  }

  Future<VisitReport> checkIn({
    required String visitReportId,
    required double latitude,
    required double longitude,
    required File photoFile, // Selfie picture is required for mobile check-in
    double? accuracy,
    double? photoLatitude,
    double? photoLongitude,
    int? photoTimestamp,
  }) async {
    try {
      // Mobile check-in always requires selfie picture - use multipart/form-data
      final formData = FormData.fromMap({
        'location[latitude]': latitude.toString(),
        'location[longitude]': longitude.toString(),
        'photo': await MultipartFile.fromFile(
          photoFile.path,
          filename:
              'selfie-checkin-${DateTime.now().millisecondsSinceEpoch}.jpg',
        ),
      });

      // Add device GPS metadata if accuracy is provided
      if (accuracy != null) {
        formData.fields.addAll([
          MapEntry('device_gps[latitude]', latitude.toString()),
          MapEntry('device_gps[longitude]', longitude.toString()),
          MapEntry('device_gps[accuracy]', accuracy.toString()),
          MapEntry(
            'device_gps[timestamp]',
            (DateTime.now().millisecondsSinceEpoch ~/ 1000).toString(),
          ),
        ]);
      }

      // Add photo GPS metadata if provided (from EXIF)
      if (photoLatitude != null && photoLongitude != null) {
        formData.fields.addAll([
          MapEntry('photo_gps[latitude]', photoLatitude.toString()),
          MapEntry('photo_gps[longitude]', photoLongitude.toString()),
          if (photoTimestamp != null)
            MapEntry('photo_gps[timestamp]', photoTimestamp.toString()),
        ]);
      }

      final response = await _dio.post(
        '/api/v1/mobile/visit-reports/$visitReportId/check-in',
        data: formData,
      );

      if (response.data['success'] == true) {
        return VisitReport.fromJson(
          response.data['data'] as Map<String, dynamic>,
        );
      } else {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to check in',
        );
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final error = e.response!.data;
        if (error is Map<String, dynamic>) {
          // Check for error object
          if (error['error'] != null) {
            final errorObj = error['error'];
            if (errorObj is Map<String, dynamic>) {
              throw Exception(errorObj['message'] ?? 'Failed to check in');
            }
            throw Exception(errorObj.toString());
          }
          // Check for direct message
          if (error['message'] != null) {
            throw Exception(error['message']);
          }
        }
        // Return status code based error
        final statusCode = e.response!.statusCode;
        if (statusCode == 500) {
          throw Exception('Server error occurred. Please try again later.');
        }
        throw Exception(
          'Failed to check in: ${e.response!.statusMessage ?? e.message}',
        );
      }
      throw Exception('Failed to check in: ${e.message ?? "Network error"}');
    } catch (e) {
      throw Exception('Failed to check in: $e');
    }
  }

  Future<VisitReport> checkOut({
    required String visitReportId,
    required double latitude,
    required double longitude,
    File? photoFile, // Optional photo for check-out (like web version)
    double? accuracy,
    double? photoLatitude,
    double? photoLongitude,
    int? photoTimestamp,
  }) async {
    try {
      // If photo is provided, use multipart/form-data (like web version)
      if (photoFile != null) {
        final formData = FormData.fromMap({
          'location[latitude]': latitude.toString(),
          'location[longitude]': longitude.toString(),
          'photo': await MultipartFile.fromFile(
            photoFile.path,
            filename:
                'selfie-checkout-${DateTime.now().millisecondsSinceEpoch}.jpg',
          ),
        });

        // Add device GPS metadata if accuracy is provided
        if (accuracy != null) {
          formData.fields.addAll([
            MapEntry('device_gps[latitude]', latitude.toString()),
            MapEntry('device_gps[longitude]', longitude.toString()),
            MapEntry('device_gps[accuracy]', accuracy.toString()),
            MapEntry(
              'device_gps[timestamp]',
              (DateTime.now().millisecondsSinceEpoch ~/ 1000).toString(),
            ),
          ]);
        }

        // Add photo GPS metadata if provided (from EXIF)
        if (photoLatitude != null && photoLongitude != null) {
          formData.fields.addAll([
            MapEntry('photo_gps[latitude]', photoLatitude.toString()),
            MapEntry('photo_gps[longitude]', photoLongitude.toString()),
            if (photoTimestamp != null)
              MapEntry('photo_gps[timestamp]', photoTimestamp.toString()),
          ]);
        }

        final response = await _dio.post(
          '/api/v1/mobile/visit-reports/$visitReportId/check-out',
          data: formData,
        );

        if (response.data['success'] == true) {
          return VisitReport.fromJson(
            response.data['data'] as Map<String, dynamic>,
          );
        } else {
          throw Exception(
            response.data['error']?['message'] ?? 'Failed to check out',
          );
        }
      } else {
        // No photo - use JSON (like web version)
        final response = await _dio.post(
          '/api/v1/mobile/visit-reports/$visitReportId/check-out',
          data: {
            'location': {'latitude': latitude, 'longitude': longitude},
          },
        );

        if (response.data['success'] == true) {
          return VisitReport.fromJson(
            response.data['data'] as Map<String, dynamic>,
          );
        } else {
          throw Exception(
            response.data['error']?['message'] ?? 'Failed to check out',
          );
        }
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final error = e.response!.data;
        if (error is Map<String, dynamic> && error['error'] != null) {
          throw Exception(error['error']['message'] ?? 'Failed to check out');
        }
      }
      throw Exception('Failed to check out: ${e.message}');
    } catch (e) {
      throw Exception('Failed to check out: $e');
    }
  }

  Future<void> uploadPhoto({
    required String visitReportId,
    required File photoFile,
  }) async {
    try {
      final formData = FormData.fromMap({
        'photo': await MultipartFile.fromFile(
          photoFile.path,
          filename: photoFile.path.split('/').last,
        ),
      });

      // Use mobile endpoint - validates ownership
      final response = await _dio.post(
        '/api/v1/mobile/visit-reports/$visitReportId/photos',
        data: formData,
        // Dio will automatically set Content-Type to multipart/form-data
      );

      if (response.data['success'] != true) {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to upload photo',
        );
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final error = e.response!.data;
        if (error is Map<String, dynamic> && error['error'] != null) {
          throw Exception(
            error['error']['message'] ?? 'Failed to upload photo',
          );
        }
      }
      throw Exception('Failed to upload photo: ${e.message}');
    } catch (e) {
      throw Exception('Failed to upload photo: $e');
    }
  }
}
