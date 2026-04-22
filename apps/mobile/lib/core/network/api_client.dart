import 'package:dio/dio.dart';

import '../config/env.dart';
import '../storage/local_storage.dart';

class ApiClient {
  const ApiClient._();

  static final Dio dio = Dio(
    BaseOptions(
      baseUrl: Env.apiBaseUrl,
      connectTimeout: const Duration(seconds: 30),
      receiveTimeout: const Duration(seconds: 60),
      sendTimeout: const Duration(seconds: 30),
      contentType: 'application/json',
    ),
  );

  static void setupInterceptors({
    required Future<bool> Function() onRefreshToken,
    required Future<void> Function() onLogout,
  }) {
    dio.interceptors.clear();
    dio.interceptors.addAll([
      _AuthInterceptor(
        onRefreshToken: onRefreshToken,
        onLogout: onLogout,
      ),
      LogInterceptor(
        requestBody: true,
        responseBody: true,
      ),
    ]);
  }
}

class _AuthInterceptor extends Interceptor {
  _AuthInterceptor({
    required this.onRefreshToken,
    required this.onLogout,
  });

  final Future<bool> Function() onRefreshToken;
  final Future<void> Function() onLogout;
  bool _isRefreshing = false;
  final List<_PendingRequest> _pendingRequests = [];

  @override
  void onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    // Get token from local storage
    final storage = await LocalStorage.create();
    final token = storage.getAuthToken();

    if (token != null) {
      options.headers['Authorization'] = 'Bearer $token';
    }

    handler.next(options);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    // Check if error is 401 (Unauthorized) - handle all 401 errors as potential token issues
    if (err.response?.statusCode == 401) {
      final responseData = err.response!.data;
      String? errorCode;
      String? errorMessage;

      if (responseData is Map<String, dynamic> &&
          responseData['error'] != null) {
        final error = responseData['error'] as Map<String, dynamic>;
        errorCode = error['code'] as String?;
        errorMessage = error['message'] as String? ?? 'An error occurred';
      } else if (responseData is String) {
        errorMessage = responseData;
      }

      // Check if it's a token-related error (expired, invalid, unauthorized)
      final isTokenError = errorCode == 'TOKEN_EXPIRED' ||
          errorCode == 'UNAUTHORIZED' ||
          errorMessage?.toLowerCase().contains('expired') == true ||
          errorMessage?.toLowerCase().contains('invalid token') == true ||
          errorMessage?.toLowerCase().contains('unauthorized') == true;

      if (isTokenError) {
        // Add request to pending queue
        final pendingRequest = _PendingRequest(
          options: err.requestOptions,
          handler: handler,
        );
        _pendingRequests.add(pendingRequest);

        // Try to refresh token if not already refreshing
        if (!_isRefreshing) {
          _isRefreshing = true;
          final success = await _refreshToken();

          // Process all pending requests
          final pending = List<_PendingRequest>.from(_pendingRequests);
          _pendingRequests.clear();

          if (success) {
            // Retry all pending requests with new token
            for (final request in pending) {
              await _retryRequest(request);
            }
          } else {
            // Refresh failed, reject all pending requests
            for (final request in pending) {
              request.handler.next(
                DioException(
                  requestOptions: request.options,
                  error: 'Token refresh failed. Please login again.',
                ),
              );
            }
          }

          _isRefreshing = false;
          return;
        }
        // If already refreshing, wait for it to complete
        return;
      }
    }

    // Handle other API error response format
    if (err.response != null) {
      final responseData = err.response!.data;
      if (responseData is Map<String, dynamic> &&
          responseData['error'] != null) {
        final error = responseData['error'] as Map<String, dynamic>;
        final message = error['message'] as String? ?? 'An error occurred';
        err = err.copyWith(
          message: message,
        );
      }
    }
    handler.next(err);
  }

  Future<bool> _refreshToken() async {
    try {
      return await onRefreshToken();
    } catch (e) {
      await onLogout();
      return false;
    }
  }

  Future<void> _retryRequest(_PendingRequest pendingRequest) async {
    try {
      // Get new token
      final storage = await LocalStorage.create();
      final token = storage.getAuthToken();

      // Update request headers with new token
      pendingRequest.options.headers['Authorization'] = 'Bearer $token';

      // Create new request options
      final newOptions = Options(
        method: pendingRequest.options.method,
        headers: pendingRequest.options.headers,
      );

      // Retry the request
      final response = await ApiClient.dio.request(
        pendingRequest.options.path,
        data: pendingRequest.options.data,
        queryParameters: pendingRequest.options.queryParameters,
        options: newOptions,
        cancelToken: pendingRequest.options.cancelToken,
        onReceiveProgress: pendingRequest.options.onReceiveProgress,
        onSendProgress: pendingRequest.options.onSendProgress,
      );

      // Resolve with successful response
      pendingRequest.handler.resolve(response);
    } catch (e) {
      // Reject with error
      pendingRequest.handler.next(
        DioException(
          requestOptions: pendingRequest.options,
          error: e,
        ),
      );
    }
  }
}

class _PendingRequest {
  _PendingRequest({
    required this.options,
    required this.handler,
  });

  final RequestOptions options;
  final ErrorInterceptorHandler handler;
}


