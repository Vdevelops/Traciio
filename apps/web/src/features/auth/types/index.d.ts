export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  success: boolean;
  data: {
    user: User;
    token: string;
    refresh_token: string;
    expires_in: number;
    csrf_token?: string; // CSRF token for protection
  };
  timestamp: string;
  request_id: string;
}

export interface User {
  id: string;
  email: string;
  name: string;
  avatar_url?: string;
  role: string;
  permissions: string[];
  created_at: string;
  updated_at: string;
}

export interface AuthState {
  user: User | null;
  token: string | null;
  refreshToken: string | null;
  isAuthenticated: boolean;
  isSessionVerified: boolean; // Backend verification flag - NOT persisted
  isLoading: boolean;
  error: string | null;
}

