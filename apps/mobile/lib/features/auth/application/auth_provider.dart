import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/network/api_client.dart';
import '../../../core/storage/local_storage.dart';
import '../data/auth_repository.dart';
import '../data/models/login_response.dart';
import 'auth_state.dart';

final localStorageProvider = FutureProvider<LocalStorage>((ref) async {
  return LocalStorage.create();
});

final authRepositoryProvider = Provider<AuthRepository>((ref) {
  return AuthRepository(ApiClient.dio);
});

final authProvider = NotifierProvider<AuthNotifier, AuthState>(AuthNotifier.new);

class AuthNotifier extends Notifier<AuthState> {
  late final AuthRepository _repository;

  @override
  AuthState build() {
    _repository = ref.read(authRepositoryProvider);
    // Check existing token saat app start (auto-login)
    Future.microtask(() => _checkAuthStatus());
    return AuthState.unknown();
  }

  /// Check authentication status saat app start
  /// Jika ada token dan rememberMe = true, set state ke authenticated
  Future<void> _checkAuthStatus() async {
    try {
      final storage = await ref.read(localStorageProvider.future);
      final token = storage.getAuthToken();
      final rememberMe = storage.getRememberMe();

      if (token != null && rememberMe) {
        // Token ada dan rememberMe = true → auto-login
        // Try to get user data from storage first
        final savedUser = storage.getUser();
        if (savedUser != null) {
          // User data exists in storage, use it
          state = AuthState.authenticated(user: savedUser);
        } else {
          // No user data in storage, try to refresh token to get user data
          final refreshed = await refreshToken();
          if (!refreshed) {
            // Refresh failed, logout
            state = AuthState.unauthenticated();
          }
          // If refresh succeeded, state is already set in refreshToken()
        }
      } else if (token != null && !rememberMe) {
        // Token ada tapi rememberMe = false → hapus token (session only)
        await storage.clearAuthToken();
        await storage.clearUser();
        state = AuthState.unauthenticated();
      } else {
        // Tidak ada token → unauthenticated
        state = AuthState.unauthenticated();
      }
    } catch (e) {
      // Jika error, anggap unauthenticated
      state = AuthState.unauthenticated();
    }
  }

  Future<void> login(
    String email,
    String password, {
    bool rememberMe = false,
  }) async {
    if (email.isEmpty || password.isEmpty) {
      state = state.copyWith(
        errorMessage: 'Email dan password wajib diisi',
        isLoading: false,
        status: AuthStatus.unauthenticated,
      );
      return;
    }

    state = state.copyWith(
      isLoading: true,
      errorMessage: null,
    );

    try {
      final loginResponse = await _repository.login(
        email: email,
        password: password,
      );

      // Save tokens to local storage
      final storage = await ref.read(localStorageProvider.future);
      await storage.saveAuthToken(loginResponse.token);
      await storage.saveRefreshToken(loginResponse.refreshToken);
      await storage.setRememberMe(rememberMe);
      // Save user data to storage for auto-login
      await storage.saveUser(loginResponse.user);

      // Set authenticated state with user data
      state = AuthState(
        status: AuthStatus.authenticated,
        isLoading: false,
        errorMessage: null,
        user: loginResponse.user,
      );
    } catch (error) {
      final errorMessage = _extractErrorMessage(error);
      state = state.copyWith(
        isLoading: false,
        status: AuthStatus.unauthenticated,
        errorMessage: errorMessage,
      );
    }
  }

  Future<void> logout() async {
      final storage = await ref.read(localStorageProvider.future);
    await storage.clearAuthToken();
    await storage.clearUser();
    state = AuthState.unauthenticated();
  }

  Future<bool> refreshToken() async {
    try {
      final storage = await ref.read(localStorageProvider.future);
      final refreshToken = storage.getRefreshToken();

      if (refreshToken == null) {
        return false;
      }

      final loginResponse = await _repository.refreshToken(refreshToken);

      // Save new tokens to local storage
      await storage.saveAuthToken(loginResponse.token);
      await storage.saveRefreshToken(loginResponse.refreshToken);
      // Save user data to storage
      await storage.saveUser(loginResponse.user);

      // Update state with new user data
      state = AuthState(
        status: AuthStatus.authenticated,
        isLoading: false,
        errorMessage: null,
        user: loginResponse.user,
      );

      return true;
    } catch (e) {
      // Refresh token failed, logout user
      await logout();
      return false;
    }
  }

  /// Update user data in auth state (e.g., after profile update)
  Future<void> updateUser(UserResponse user) async {
      final storage = await ref.read(localStorageProvider.future);
    await storage.saveUser(user);
    state = AuthState(
      status: AuthStatus.authenticated,
      isLoading: false,
      errorMessage: null,
      user: user,
    );
  }

  String _extractErrorMessage(dynamic error) {
    // Check if it's AuthException (from repository)
    // AuthException.toString() returns just the message, so we can use that
    if (error is Exception) {
      final errorString = error.toString();
      
      // AuthException.toString() returns just the message (no prefix)
      // Regular Exception.toString() returns "Exception: message"
      if (errorString.startsWith('Exception: ')) {
        return errorString.substring(11);
      }
      
      // If no prefix, it's likely AuthException or already just the message
      return errorString;
    }
    
    // If error is already a string message, return it directly
    if (error is String) {
      return error;
    }
    
    // Handle DioException (should be rare now since we handle it in repository)
    if (error is DioException) {
      if (error.response != null) {
        final responseData = error.response!.data;
        if (responseData is Map<String, dynamic>) {
          // Try to extract error message from response
          if (responseData['error'] != null) {
            final errorObj = responseData['error'];
            if (errorObj is Map<String, dynamic>) {
              final message = errorObj['message'] as String?;
              if (message != null && message.isNotEmpty) {
                return message;
              }
            } else if (errorObj is String) {
              return errorObj;
            }
          }
          
          // Check for message directly
          final message = responseData['message'] as String?;
          if (message != null && message.isNotEmpty) {
            return message;
          }
        }
      }
      
      // Fallback to DioException message
      return error.message ?? 'Network error occurred';
    }
    
    // Generic fallback - extract message from exception
    final errorString = error.toString();
    if (errorString.startsWith('Exception: ')) {
      return errorString.substring(11);
    }
    
    return errorString.isNotEmpty 
        ? errorString 
        : 'Login failed. Please try again.';
  }
}


