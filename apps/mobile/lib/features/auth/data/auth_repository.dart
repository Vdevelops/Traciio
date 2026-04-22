import 'package:dio/dio.dart';

import 'models/login_response.dart';

/// Custom exception class for auth errors with message from API
class AuthException implements Exception {
  AuthException(this.message);
  
  final String message;
  
  @override
  String toString() => message;
}

class AuthRepository {
  const AuthRepository(this._dio);

  final Dio _dio;

  Future<LoginResponse> login({
    required String email,
    required String password,
  }) async {
    try {
      final response = await _dio.post<Map<String, dynamic>>(
        '/api/v1/auth/mobile/login',
        data: {
          'email': email,
          'password': password,
        },
      );

      // Parse API response format: { success: true, data: { ... } }
      if (response.data == null) {
        throw AuthException('Invalid response from server');
      }

      final responseData = response.data!;
      if (responseData['success'] != true) {
        // Extract error details from API response
        // API format: { success: false, error: { message: "...", code: "...", details: "..." } }
        final error = responseData['error'];
        String? details;
        
        if (error is Map<String, dynamic>) {
          // Priority: details > message
          final detailsValue = error['details'];
          if (detailsValue != null) {
            details = _extractStringFromValue(detailsValue);
          }
          if (details == null || details.isEmpty) {
            final messageValue = error['message'];
            if (messageValue != null) {
              details = _extractStringFromValue(messageValue);
            }
          }
        } else if (error is String) {
          details = error;
        }
        
        // Use details from API, or fallback
        throw AuthException(details ?? 'Login failed');
      }

      final data = responseData['data'] as Map<String, dynamic>;
      return LoginResponse.fromJson(data);
    } on DioException catch (e) {
      // Handle DioException - extract error message from response
      if (e.response != null) {
        final responseData = e.response!.data;
        
        // Try to extract error message from response
        String? errorMessage;
        
        if (responseData is Map<String, dynamic>) {
          // Check for error object with details
          if (responseData['error'] != null) {
            final error = responseData['error'];
            if (error is Map<String, dynamic>) {
              // Priority: details > message
              final detailsValue = error['details'];
              if (detailsValue != null) {
                errorMessage = _extractStringFromValue(detailsValue);
              }
              if (errorMessage == null || errorMessage.isEmpty) {
                final messageValue = error['message'];
                if (messageValue != null) {
                  errorMessage = _extractStringFromValue(messageValue);
                }
              }
            } else if (error is String) {
              errorMessage = error;
            }
          }
          
          // If no error object, check for details or message directly
          if (errorMessage == null) {
            final detailsValue = responseData['details'];
            if (detailsValue != null) {
              errorMessage = _extractStringFromValue(detailsValue);
            }
            if (errorMessage == null || errorMessage.isEmpty) {
              final messageValue = responseData['message'];
              if (messageValue != null) {
                errorMessage = _extractStringFromValue(messageValue);
              }
            }
          }
        }
        
        // If we found an error message, throw it
        if (errorMessage != null && errorMessage.isNotEmpty) {
          throw AuthException(errorMessage);
        }
        
        // Fallback to status code based messages
        final statusCode = e.response!.statusCode;
        if (statusCode == 401) {
          throw AuthException('Invalid email or password');
        } else if (statusCode == 403) {
          throw AuthException('Access forbidden. Your role may not be allowed to login via mobile.');
        } else if (statusCode == 404) {
          throw AuthException('Login endpoint not found');
        } else if (statusCode == 422) {
          throw AuthException('Validation error. Please check your input');
        } else if (statusCode == 500) {
          throw AuthException('Server error. Please try again later');
        }
      }
      
      // Handle network errors
      if (e.type == DioExceptionType.connectionTimeout ||
          e.type == DioExceptionType.sendTimeout ||
          e.type == DioExceptionType.receiveTimeout) {
        throw AuthException('Connection timeout. Please check your internet connection');
      }
      
      if (e.type == DioExceptionType.connectionError) {
        throw AuthException('Unable to connect to server. Please check your internet connection');
      }
      
      // Generic network error
      throw AuthException('Network error: ${e.message ?? "Unknown error"}');
    } catch (e) {
      // Re-throw AuthException as is
      if (e is AuthException) {
        rethrow;
      }
      // Wrap other exceptions
      throw AuthException(e.toString());
    }
  }

  Future<LoginResponse> refreshToken(String refreshToken) async {
    try {
      final response = await _dio.post<Map<String, dynamic>>(
        '/api/v1/auth/refresh',
        data: {
          'refresh_token': refreshToken,
        },
      );

      // Parse API response format: { success: true, data: { ... } }
      if (response.data == null) {
        throw Exception('Invalid response from server');
      }

      final responseData = response.data!;
      if (responseData['success'] != true) {
        final error = responseData['error'];
        String? details;
        
        if (error is Map<String, dynamic>) {
          // Priority: details > message
          final detailsValue = error['details'];
          if (detailsValue != null) {
            details = _extractStringFromValue(detailsValue);
          }
          if (details == null || details.isEmpty) {
            final messageValue = error['message'];
            if (messageValue != null) {
              details = _extractStringFromValue(messageValue);
            }
          }
        } else if (error is String) {
          details = error;
        }
        
        throw AuthException(details ?? 'Failed to refresh token');
      }

      final data = responseData['data'] as Map<String, dynamic>;
      return LoginResponse.fromJson(data);
    } on DioException catch (e) {
      if (e.response != null) {
        final responseData = e.response!.data;
        String? details;
        
        if (responseData is Map<String, dynamic>) {
          if (responseData['error'] != null) {
            final error = responseData['error'];
            if (error is Map<String, dynamic>) {
              // Priority: details > message
              final detailsValue = error['details'];
              if (detailsValue != null) {
                details = _extractStringFromValue(detailsValue);
              }
              if (details == null || details.isEmpty) {
                final messageValue = error['message'];
                if (messageValue != null) {
                  details = _extractStringFromValue(messageValue);
                }
              }
            } else if (error is String) {
              details = error;
            }
          }
          
          if (details == null) {
            final detailsValue = responseData['details'];
            if (detailsValue != null) {
              details = _extractStringFromValue(detailsValue);
            }
            if (details == null || details.isEmpty) {
              final messageValue = responseData['message'];
              if (messageValue != null) {
                details = _extractStringFromValue(messageValue);
              }
            }
          }
        }
        
        if (details != null && details.isNotEmpty) {
          throw AuthException(details);
        }
      }
      throw AuthException('Network error: ${e.message ?? "Unknown error"}');
    } catch (e) {
      if (e is AuthException) {
        rethrow;
      }
      throw AuthException(e.toString());
    }
  }
  
  /// Extract string value from dynamic value
  /// Handles cases where value might be String, Map, List, or other types
  /// Returns only the text content, not JSON format
  String? _extractStringFromValue(dynamic value) {
    if (value == null) {
      return null;
    }
    
    if (value is String) {
      return value.isEmpty ? null : value;
    }
    
    if (value is Map) {
      // If it's a Map, try to find common error message fields
      // Priority: reason > message > details > error
      final fields = ['reason', 'message', 'details', 'error', 'msg'];
      for (final field in fields) {
        if (value.containsKey(field)) {
          final fieldValue = value[field];
          if (fieldValue is String && fieldValue.isNotEmpty) {
            return fieldValue;
          }
          // If field value is also a Map, recursively extract
          if (fieldValue is Map) {
            final extracted = _extractStringFromValue(fieldValue);
            if (extracted != null && extracted.isNotEmpty) {
              return extracted;
            }
          }
        }
      }
      
      // If Map has only one entry, use its value
      if (value.length == 1) {
        final firstValue = value.values.first;
        if (firstValue is String && firstValue.isNotEmpty) {
          return firstValue;
        }
        // Recursively extract if value is nested
        if (firstValue is Map || firstValue is List) {
          return _extractStringFromValue(firstValue);
        }
      }
      
      // If multiple entries, try to find the most relevant text field
      // or join all string values
      final stringValues = <String>[];
      for (final entry in value.entries) {
        if (entry.value is String && (entry.value as String).isNotEmpty) {
          stringValues.add(entry.value as String);
        }
      }
      if (stringValues.isNotEmpty) {
        return stringValues.join('. ');
      }
      
      // Last resort: return null to use fallback message
      return null;
    }
    
    if (value is List) {
      // If it's a List, extract strings from elements
      final stringValues = <String>[];
      for (final element in value) {
        if (element is String && element.isNotEmpty) {
          stringValues.add(element);
        } else if (element is Map || element is List) {
          final extracted = _extractStringFromValue(element);
          if (extracted != null && extracted.isNotEmpty) {
            stringValues.add(extracted);
          }
        }
      }
      if (stringValues.isNotEmpty) {
        return stringValues.join('. ');
      }
      return null;
    }
    
    // For other types, convert to string only if it's meaningful
    final stringValue = value.toString();
    // Avoid showing object representations like "Instance of..."
    if (stringValue.startsWith('Instance of') || 
        stringValue.startsWith('_')) {
      return null;
    }
    return stringValue;
  }
}


