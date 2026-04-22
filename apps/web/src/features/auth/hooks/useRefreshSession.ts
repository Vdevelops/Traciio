import { useAuthStore } from "../stores/useAuthStore";
import { authService } from "../services/authService";
import { useLogout } from "./useLogout";
import { setSecureCookie } from "@/lib/cookie";

export function useRefreshSession() {
  const { refreshToken, setUser, setToken, setSessionVerified } = useAuthStore();
  const handleLogout = useLogout();

  const refreshSession = async () => {
    const currentRefreshToken = refreshToken;
    if (!currentRefreshToken) {
      throw new Error("No refresh token available");
    }
    try {
      const response = await authService.refreshToken(currentRefreshToken);
      if (response.success && response.data) {
        const { user, token, refresh_token } = response.data;
        if (typeof window !== "undefined") {
          localStorage.setItem("token", token);
          localStorage.setItem("refreshToken", refresh_token);
          // Update secure cookie
          setSecureCookie("token", token);
        }
        setUser(user);
        setToken(token);
        // Update refresh token in store and mark session as verified
        setSessionVerified(true);
        useAuthStore.setState({ refreshToken: refresh_token, isAuthenticated: true });
      }
    } catch (error) {
      // Refresh failed, logout user
      await handleLogout();
      throw error;
    }
  };

  return { refreshSession };
}

