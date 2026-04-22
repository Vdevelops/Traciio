import '../data/models/login_response.dart';

enum AuthStatus {
  unknown,
  authenticated,
  unauthenticated,
}

class AuthState {
  const AuthState({
    required this.status,
    this.isLoading = false,
    this.errorMessage,
    this.user,
  });

  final AuthStatus status;
  final bool isLoading;
  final String? errorMessage;
  final UserResponse? user;

  factory AuthState.unknown() {
    return const AuthState(status: AuthStatus.unknown);
  }

  factory AuthState.unauthenticated() {
    return const AuthState(status: AuthStatus.unauthenticated);
  }

  factory AuthState.authenticated({UserResponse? user}) {
    return AuthState(status: AuthStatus.authenticated, user: user);
  }

  AuthState copyWith({
    AuthStatus? status,
    bool? isLoading,
    String? errorMessage,
    UserResponse? user,
  }) {
    return AuthState(
      status: status ?? this.status,
      isLoading: isLoading ?? this.isLoading,
      errorMessage: errorMessage,
      user: user ?? this.user,
    );
  }
}


