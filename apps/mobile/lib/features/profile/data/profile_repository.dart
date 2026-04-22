import 'package:dio/dio.dart';

import '../../../core/network/connectivity_service.dart';
import '../../../core/storage/offline_storage.dart';
import 'models/profile.dart';

/// Custom exception for profile-related errors
class ProfileException implements Exception {
  const ProfileException(this.message, {this.statusCode, this.details});

  final String message;
  final int? statusCode;
  final Map<String, dynamic>? details;

  @override
  String toString() => message;
}

class ProfileRepository {
  const ProfileRepository(this._dio, this._connectivity);

  final Dio _dio;
  final ConnectivityService _connectivity;

  /// Get complete profile for authenticated user with offline support
  /// Uses mobile endpoint: /api/v1/auth/mobile/profile
  Future<ProfileResponse> getMyProfile({bool forceRefresh = false}) async {
    // 1. Try cache if not forcing refresh
    if (!forceRefresh) {
      final cached = await OfflineStorage.getUserProfile();
      if (cached != null) {
        try {
          final profile = ProfileResponse.fromJson(cached);

          // If online, fetch background update
          if (_connectivity.isOnline) {
            _fetchProfileInBackground()
                .then((data) {
                  OfflineStorage.saveUserProfile(data.toJson());
                })
                .catchError((_) {
                  // Ignore background fetch errors
                });
          }
          return profile;
        } catch (e) {
          // parsing error, ignore and try network
        }
      }
    }

    // 2. Fetch from network
    if (_connectivity.isOnline) {
      try {
        final profile = await _fetchProfileFromApi();
        await OfflineStorage.saveUserProfile(profile.toJson());
        return profile;
      } catch (e) {
        // If network fails, try cache as fallback
        final cached = await OfflineStorage.getUserProfile();
        if (cached != null) {
          return ProfileResponse.fromJson(cached);
        }
        rethrow;
      }
    }

    // 3. Fallback to cache if offline
    final cached = await OfflineStorage.getUserProfile();
    if (cached != null) {
      return ProfileResponse.fromJson(cached);
    }

    throw ProfileException(
      'No internet connection and no cached data available',
    );
  }

  Future<ProfileResponse> _fetchProfileFromApi() async {
    try {
      final response = await _dio.get<Map<String, dynamic>>(
        '/api/v1/auth/mobile/profile',
      );

      // Parse API response format: { success: true, data: { user: {...}, stats: {...}, ... } }
      if (response.data == null) {
        throw Exception('Invalid response from server');
      }

      final responseData = response.data!;
      if (responseData['success'] != true) {
        final error = responseData['error'];
        String? message;
        if (error is Map<String, dynamic>) {
          message = error['message'] as String?;
        } else if (error is String) {
          message = error;
        }
        throw Exception(message ?? 'Failed to fetch profile');
      }

      // Safely extract data field
      final data = responseData['data'];
      if (data == null) {
        throw Exception('Invalid response: data field is null');
      }

      // Handle case where data might be a Map or already the profile object
      Map<String, dynamic> profileData;
      if (data is Map<String, dynamic>) {
        profileData = data;
      } else if (data is String) {
        // If data is a string, try to parse it (shouldn't happen, but defensive)
        throw Exception('Invalid response format: data is a string');
      } else {
        throw Exception('Invalid response format: data is not a map');
      }

      return ProfileResponse.fromJson(profileData);
    } on DioException catch (e) {
      // Handle 404 and other HTTP errors
      if (e.response != null) {
        final statusCode = e.response!.statusCode;
        final responseData = e.response!.data;

        // Handle 404 - endpoint not found
        if (statusCode == 404) {
          throw ProfileException(
            'Profile endpoint not found. Please check if the API is deployed correctly.',
            statusCode: statusCode,
          );
        }

        // Handle response data - could be Map or String
        if (responseData is Map<String, dynamic>) {
          if (responseData['error'] != null) {
            final error = responseData['error'];
            String? message;

            if (error is Map<String, dynamic>) {
              message = error['message'] as String?;
              // Try to extract from details if message is not available
              if (message == null || message.isEmpty) {
                final details = error['details'];
                if (details != null) {
                  message = _extractStringFromValue(details);
                }
              }
            } else if (error is String) {
              message = error;
            }

            throw ProfileException(
              message ?? 'Failed to fetch profile',
              statusCode: statusCode,
            );
          }
        } else if (responseData is String) {
          // Response is a plain string (e.g., "404 page not found")
          throw ProfileException(
            'Server error: $responseData',
            statusCode: statusCode,
          );
        }
      }

      // Network or other errors
      throw ProfileException('Network error: ${e.message ?? 'Unknown error'}');
    } catch (e) {
      rethrow;
    }
  }

  /// Fetch profile in background for cache update
  Future<ProfileResponse> _fetchProfileInBackground() async {
    try {
      final profile = await _fetchProfileFromApi();
      await OfflineStorage.saveUserProfile(profile.toJson());
      return profile;
    } catch (e) {
      // Silent fail for background fetch
      rethrow;
    }
  }

  /// Update profile information for authenticated user
  /// Uses mobile endpoint: /api/v1/auth/mobile/profile
  Future<ProfileUser> updateMyProfile(UpdateProfileRequest request) async {
    try {
      final response = await _dio.put<Map<String, dynamic>>(
        '/api/v1/auth/mobile/profile',
        data: request.toJson(),
      );

      // Parse API response format: { success: true, data: { id: "...", name: "...", ... } }
      if (response.data == null) {
        throw Exception('Invalid response from server');
      }

      final responseData = response.data!;
      if (responseData['success'] != true) {
        final error = responseData['error'] as Map<String, dynamic>?;
        final message =
            error?['message'] as String? ?? 'Failed to update profile';
        throw Exception(message);
      }

      final data = responseData['data'] as Map<String, dynamic>;
      return ProfileUser.fromJson(data);
    } on DioException catch (e) {
      if (e.response != null) {
        final statusCode = e.response!.statusCode;
        final responseData = e.response!.data;

        // Handle 404
        if (statusCode == 404) {
          throw ProfileException(
            'Profile endpoint not found. Please check if the API is deployed correctly.',
            statusCode: statusCode,
          );
        }

        // Handle response data - could be Map or String
        if (responseData is Map<String, dynamic>) {
          if (responseData['error'] != null) {
            final error = responseData['error'];
            String? message;

            if (error is Map<String, dynamic>) {
              message = error['message'] as String?;

              // Check for validation errors
              if (error['field_errors'] != null) {
                final fieldErrors = error['field_errors'];
                if (fieldErrors is List && fieldErrors.isNotEmpty) {
                  final firstError = fieldErrors.first;
                  if (firstError is Map<String, dynamic>) {
                    message = firstError['message'] as String? ?? message;
                  }
                }
              }

              // Try to extract from details if message is not available
              if (message == null || message.isEmpty) {
                final details = error['details'];
                if (details != null) {
                  message = _extractStringFromValue(details);
                }
              }
            } else if (error is String) {
              message = error;
            }

            throw ProfileException(
              message ?? 'Failed to update profile',
              statusCode: statusCode,
            );
          }
        } else if (responseData is String) {
          throw ProfileException(
            'Server error: $responseData',
            statusCode: statusCode,
          );
        }
      }
      throw ProfileException('Network error: ${e.message ?? 'Unknown error'}');
    } catch (e) {
      rethrow;
    }
  }

  /// Change password for authenticated user
  /// Uses mobile endpoint: /api/v1/auth/mobile/password
  Future<void> changeMyPassword(ChangePasswordRequest request) async {
    try {
      final response = await _dio.put<Map<String, dynamic>>(
        '/api/v1/auth/mobile/password',
        data: request.toJson(),
      );

      // Success response is 204 No Content, but Dio might return null
      // Check for error response format: { success: false, error: {...} }
      if (response.data != null) {
        final responseData = response.data!;
        if (responseData['success'] != true) {
          final error = responseData['error'] as Map<String, dynamic>?;
          final message =
              error?['message'] as String? ?? 'Failed to change password';

          // Check for validation errors
          if (error?['field_errors'] != null) {
            final fieldErrors = error!['field_errors'] as List<dynamic>?;
            if (fieldErrors != null && fieldErrors.isNotEmpty) {
              final firstError = fieldErrors.first as Map<String, dynamic>;
              final fieldMessage = firstError['message'] as String? ?? message;
              throw Exception(fieldMessage);
            }
          }

          throw Exception(message);
        }
      }
    } on DioException catch (e) {
      if (e.response != null) {
        final statusCode = e.response!.statusCode;
        final responseData = e.response!.data;

        // Handle 404
        if (statusCode == 404) {
          throw ProfileException(
            'Password endpoint not found. Please check if the API is deployed correctly.',
            statusCode: statusCode,
          );
        }

        // Handle response data - could be Map or String
        if (responseData is Map<String, dynamic>) {
          if (responseData['error'] != null) {
            final error = responseData['error'];
            String? message;

            if (error is Map<String, dynamic>) {
              message = error['message'] as String?;

              // Check for validation errors
              if (error['field_errors'] != null) {
                final fieldErrors = error['field_errors'];
                if (fieldErrors is List && fieldErrors.isNotEmpty) {
                  final firstError = fieldErrors.first;
                  if (firstError is Map<String, dynamic>) {
                    message = firstError['message'] as String? ?? message;
                  }
                }
              }

              // Try to extract from details if message is not available
              if (message == null || message.isEmpty) {
                final details = error['details'];
                if (details != null) {
                  message = _extractStringFromValue(details);
                }
              }
            } else if (error is String) {
              message = error;
            }

            throw ProfileException(
              message ?? 'Failed to change password',
              statusCode: statusCode,
            );
          }
        } else if (responseData is String) {
          throw ProfileException(
            'Server error: $responseData',
            statusCode: statusCode,
          );
        }
      }
      throw ProfileException('Network error: ${e.message ?? 'Unknown error'}');
    } catch (e) {
      rethrow;
    }
  }

  /// Helper method to safely extract string from dynamic value
  String _extractStringFromValue(dynamic value) {
    if (value is String) {
      return value;
    } else if (value is Map) {
      // If it's a map, try to extract a common error field
      return value['message']?.toString() ??
          value['details']?.toString() ??
          value.toString();
    } else {
      return value.toString();
    }
  }
}
