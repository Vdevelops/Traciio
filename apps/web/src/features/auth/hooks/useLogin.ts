import { useRouter } from "@/i18n/routing";
import { useAuthStore } from "../stores/useAuthStore";
import { authService } from "../services/authService";
import type { LoginFormData } from "../schemas/login.schema";
import type { AuthError } from "../types/errors";
import { useState } from "react";
import { setSecureCookie } from "@/lib/cookie";
import { setCSRFToken } from "@/lib/csrf";

export function useLogin() {
  const router = useRouter();
  const { setUser, setToken, setSessionVerified, isLoading: storeIsLoading, error: storeError, clearError } = useAuthStore();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleLogin = async (data: LoginFormData) => {
    setIsLoading(true);
    setError(null);
    try {
      const response = await authService.login({ email: data.email, password: data.password });
      if (response.success && response.data) {
        const { user, token, refresh_token, csrf_token } = response.data;
        if (typeof window !== "undefined") {
          localStorage.setItem("token", token);
          localStorage.setItem("refreshToken", refresh_token);
          // Set secure cookie for WebSocket and middleware
          setSecureCookie("token", token);
          // Store CSRF token if provided
          if (csrf_token) {
            setCSRFToken(csrf_token);
          }
        }
        setUser(user);
        setToken(token);
        // Set isSessionVerified = true since backend just confirmed credentials
        setSessionVerified(true);
        useAuthStore.setState({ refreshToken: refresh_token, isAuthenticated: true, error: null });
        // Redirect ke dashboard dengan locale aktif
        router.push("/dashboard");
      }
    } catch (err) {
      const authError = err as AuthError;
      const errorMessage =
        authError.response?.data?.error?.message ||
        authError.message ||
        "Login failed";
      setError(errorMessage);
      useAuthStore.setState({ isAuthenticated: false, isSessionVerified: false, error: errorMessage });
      throw err;
    } finally {
      setIsLoading(false);
    }
  };

  return {
    handleLogin,
    isLoading: isLoading || storeIsLoading,
    error: error || storeError,
    clearError: () => {
      setError(null);
      clearError();
    },
  };
}

