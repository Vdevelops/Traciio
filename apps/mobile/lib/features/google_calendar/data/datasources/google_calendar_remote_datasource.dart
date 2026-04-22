import 'package:dio/dio.dart';

import '../../domain/models/google_calendar_event.dart';
import '../../domain/models/google_calendar_status.dart';

/// Remote datasource for Google Calendar API
class GoogleCalendarRemoteDataSource {
  final Dio _dio;

  GoogleCalendarRemoteDataSource(this._dio);

  /// Get connection status
  Future<GoogleCalendarStatus> getConnectionStatus() async {
    try {
      final response = await _dio.get('/api/v1/google-calendar/status');

      if (response.data['success'] == true) {
        return GoogleCalendarStatus.fromJson(response.data['data']);
      } else {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to get status',
        );
      }
    } on DioException catch (e) {
      throw _handleDioError(e);
    }
  }

  /// Get OAuth authorization URL
  Future<String> getAuthUrl() async {
    try {
      final response = await _dio.get(
        '/api/v1/google-calendar/auth-url',
        queryParameters: {'platform': 'mobile'},
      );

      if (response.data['success'] == true) {
        return response.data['data']['auth_url'];
      } else {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to get auth URL',
        );
      }
    } on DioException catch (e) {
      throw _handleDioError(e);
    }
  }

  /// Disconnect Google Calendar
  Future<void> disconnect() async {
    try {
      final response = await _dio.delete('/api/v1/google-calendar/disconnect');

      if (response.data['success'] != true) {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to disconnect',
        );
      }
    } on DioException catch (e) {
      throw _handleDioError(e);
    }
  }

  /// Sync schedule to Google Calendar
  Future<GoogleCalendarEvent> syncSchedule(String scheduleId) async {
    try {
      final response = await _dio.post(
        '/api/v1/schedules/$scheduleId/sync-google-calendar',
      );

      if (response.data['success'] == true) {
        return GoogleCalendarEvent.fromJson(response.data['data']);
      } else {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to sync schedule',
        );
      }
    } on DioException catch (e) {
      throw _handleDioError(e);
    }
  }

  /// Unsync schedule from Google Calendar
  Future<void> unsyncSchedule(String scheduleId) async {
    try {
      final response = await _dio.post(
        '/api/v1/schedules/$scheduleId/unsync-google-calendar',
      );

      if (response.data['success'] != true) {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to unsync schedule',
        );
      }
    } on DioException catch (e) {
      throw _handleDioError(e);
    }
  }

  /// Exchange authorization code for token (Option 2 - Direct Deep Link)
  Future<void> exchangeCode(String code, String state) async {
    try {
      final response = await _dio.post(
        '/api/v1/google-calendar/exchange-code',
        data: {'code': code, 'state': state},
      );

      if (response.data['success'] != true) {
        throw Exception(
          response.data['error']?['message'] ?? 'Failed to exchange code',
        );
      }
    } on DioException catch (e) {
      throw _handleDioError(e);
    }
  }

  Exception _handleDioError(DioException e) {
    if (e.response != null) {
      final data = e.response!.data;
      if (data is Map && data['error'] != null) {
        return Exception(data['error']['message'] ?? 'Unknown error');
      }
      return Exception('Server error: ${e.response?.statusCode}');
    }
    return Exception('Network error: ${e.message}');
  }
}
