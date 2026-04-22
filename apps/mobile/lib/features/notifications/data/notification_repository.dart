import 'package:dio/dio.dart';

import 'models/notification.dart';

class NotificationRepository {
  NotificationRepository(this._dio);

  final Dio _dio;

  Future<NotificationListResponse> getNotifications({
    int page = 1,
    int perPage = 20,
    NotificationType? type,
    bool? isRead,
  }) async {
    try {
      final queryParams = <String, dynamic>{
        'page': page,
        'per_page': perPage,
      };

      if (type != null) {
        queryParams['type'] = type.value;
      }
      if (isRead != null) {
        queryParams['is_read'] = isRead;
      }

      final response = await _dio.get(
        '/api/v1/notifications',
        queryParameters: queryParams,
      );

      if (response.data is Map<String, dynamic>) {
        final responseData = response.data as Map<String, dynamic>;
        if (responseData['success'] == true) {
          return NotificationListResponse.fromJson(responseData);
        } else {
          throw Exception(
            responseData['error']?['message'] ?? 'Failed to fetch notifications',
          );
        }
      } else {
        throw Exception('Invalid response format');
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final statusCode = e.response!.statusCode;
        if (statusCode == 404) {
          throw Exception(
            'Notification endpoint not found. Please ensure the backend server is running and the notification routes are registered.',
          );
        }
        final errorData = e.response!.data;
        if (errorData is Map<String, dynamic>) {
          throw Exception(
            errorData['error']?['message'] ?? 'Failed to fetch notifications',
          );
        } else if (errorData is String) {
          throw Exception(errorData);
        } else {
          throw Exception('Failed to fetch notifications: ${e.response?.statusMessage ?? e.message}');
        }
      } else {
        throw Exception('Network error: ${e.message}');
      }
    } catch (e) {
      throw Exception('Failed to fetch notifications: $e');
    }
  }

  Future<int> getUnreadCount() async {
    try {
      final response = await _dio.get('/api/v1/notifications/unread-count');

      if (response.data is Map<String, dynamic>) {
        final responseData = response.data as Map<String, dynamic>;
        if (responseData['success'] == true) {
          return UnreadCountResponse.fromJson(responseData).unreadCount;
        } else {
          throw Exception(
            responseData['error']?['message'] ?? 'Failed to fetch unread count',
          );
        }
      } else {
        throw Exception('Invalid response format');
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final statusCode = e.response!.statusCode;
        if (statusCode == 404) {
          throw Exception(
            'Notification endpoint not found. Please ensure the backend server is running and the notification routes are registered.',
          );
        }
        final errorData = e.response!.data;
        if (errorData is Map<String, dynamic>) {
          throw Exception(
            errorData['error']?['message'] ?? 'Failed to fetch unread count',
          );
        } else if (errorData is String) {
          throw Exception(errorData);
        } else {
          throw Exception('Failed to fetch unread count: ${e.response?.statusMessage ?? e.message}');
        }
      } else {
        throw Exception('Network error: ${e.message}');
      }
    } catch (e) {
      throw Exception('Failed to fetch unread count: $e');
    }
  }

  Future<Notification> markAsRead(String id) async {
    try {
      final response = await _dio.put('/api/v1/notifications/$id/read');

      if (response.data is Map<String, dynamic>) {
        final responseData = response.data as Map<String, dynamic>;
        if (responseData['success'] == true) {
          return Notification.fromJson(responseData['data'] as Map<String, dynamic>);
        } else {
          throw Exception(
            responseData['error']?['message'] ?? 'Failed to mark notification as read',
          );
        }
      } else {
        throw Exception('Invalid response format');
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final errorData = e.response!.data;
        throw Exception(
          errorData['error']?['message'] ?? 'Failed to mark notification as read',
        );
      } else {
        throw Exception('Network error: ${e.message}');
      }
    } catch (e) {
      throw Exception('Failed to mark notification as read: $e');
    }
  }

  Future<void> markAllAsRead() async {
    try {
      final response = await _dio.put('/api/v1/notifications/read-all');

      if (response.data is Map<String, dynamic>) {
        final responseData = response.data as Map<String, dynamic>;
        if (responseData['success'] != true) {
          throw Exception(
            responseData['error']?['message'] ?? 'Failed to mark all notifications as read',
          );
        }
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final errorData = e.response!.data;
        throw Exception(
          errorData['error']?['message'] ?? 'Failed to mark all notifications as read',
        );
      } else {
        throw Exception('Network error: ${e.message}');
      }
    } catch (e) {
      throw Exception('Failed to mark all notifications as read: $e');
    }
  }

  Future<void> deleteNotification(String id) async {
    try {
      final response = await _dio.delete('/api/v1/notifications/$id');

      // Handle response - DELETE might return 204 No Content or success response
      if (response.statusCode == 204 || response.statusCode == 200) {
        // Success - no need to check response.data
        return;
      }

      // If response has data, check for success flag
      if (response.data != null) {
        final responseData = response.data;
        if (responseData is Map<String, dynamic>) {
          if (responseData['success'] == false) {
            throw Exception(
              responseData['error']?['message'] ?? 'Failed to delete notification',
            );
          }
        }
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final errorData = e.response!.data;
        if (errorData is Map<String, dynamic>) {
          throw Exception(
            errorData['error']?['message'] ?? 'Failed to delete notification',
          );
        } else {
          throw Exception('Failed to delete notification: ${e.response?.statusMessage ?? e.message}');
        }
      } else {
        throw Exception('Network error: ${e.message}');
      }
    } catch (e) {
      // Re-throw if it's already an Exception
      if (e is Exception) {
        rethrow;
      }
      throw Exception('Failed to delete notification: $e');
    }
  }
}

