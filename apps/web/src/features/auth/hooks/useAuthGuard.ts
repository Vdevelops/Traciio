import { useEffect, useState, useCallback, useRef } from "react";
import { usePathname } from "next/navigation";
import { useAuthStore } from "../stores/useAuthStore";
import { authService } from "../services/authService";
import { fullAuthCleanup } from "../utils/clear-auth-cookies";
import { setSecureCookie } from "@/lib/cookie";

/**
 * Auth Guard Hook with Session Verification
 * 
 * Prevents zombie auth state by verifying session with backend on app load.
 * Returns isSessionVerified which should be checked along with isAuthenticated
 * before allowing access to protected routes.
 */
export function useAuthGuard() {
  const pathname = usePathname();
  const { 
    isAuthenticated, 
    isSessionVerified, 
    setUser, 
    setToken, 
    setSessionVerified,
    logout 
  } = useAuthStore();

  const [isVerifying, setIsVerifying] = useState(false);
  const verificationAttempted = useRef(false);

  /**
   * Verify session with backend by calling refresh token endpoint
   * This confirms that the stored tokens are still valid
   */
  const verifySession = useCallback(async () => {
    // Skip if already verified or already verifying
    if (isSessionVerified || isVerifying) {
      return;
    }

    const refreshToken = localStorage.getItem("refreshToken");
    const token = localStorage.getItem("token");

    // No tokens - nothing to verify
    if (!refreshToken || !token) {
      setSessionVerified(false);
      return;
    }

    setIsVerifying(true);

    try {
      // Use refresh token to verify session is still valid
      const response = await authService.refreshToken(refreshToken);

      if (response.success && response.data) {
        const { user, token: newToken, refresh_token: newRefreshToken } = response.data;

        // Update localStorage with new tokens
        localStorage.setItem("token", newToken);
        localStorage.setItem("refreshToken", newRefreshToken);
        setSecureCookie("token", newToken);

        // Update store with verified user data
        setUser(user);
        setToken(newToken);
        useAuthStore.setState({ 
          refreshToken: newRefreshToken, 
          isAuthenticated: true 
        });

        // Mark session as verified - backend confirmed tokens are valid
        setSessionVerified(true);
      } else {
        throw new Error("Session verification failed");
      }
    } catch {
      // Session invalid - clear everything and redirect to login
      logout();
      await fullAuthCleanup();
      setSessionVerified(false);
    } finally {
      setIsVerifying(false);
    }
  }, [isSessionVerified, isVerifying, setUser, setToken, setSessionVerified, logout]);

  // Verify session on mount if user appears authenticated but not verified
  useEffect(() => {
    if (typeof window === "undefined") {
      return;
    }

    const token = localStorage.getItem("token");
    const refreshToken = localStorage.getItem("refreshToken");

    // Only verify if we have tokens but session not yet verified
    if (token && refreshToken && !isSessionVerified && !verificationAttempted.current) {
      verificationAttempted.current = true;
      verifySession();
    } else if (!token || !refreshToken) {
      // No tokens - mark as not verified
      setSessionVerified(false);
    }
  }, [pathname, isSessionVerified, verifySession, setSessionVerified]);

  const hasToken =
    typeof window !== "undefined" && !!localStorage.getItem("token");
  const isLoading = isVerifying || (!isSessionVerified && hasToken);

  return {
    // Both conditions required for secure auth check
    isAuthenticated: isAuthenticated && isSessionVerified,
    // Raw values for components that need granular control
    isAuthenticatedRaw: isAuthenticated || hasToken,
    isSessionVerified,
    isLoading,
    // Manual trigger for session verification
    verifySession,
  };
}
