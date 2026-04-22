"use client";

import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { LoadingSpinner } from "@/components/ui/loading-spinner";
import { authService } from "@/features/auth/services/authService";
import { useAuthStore } from "@/features/auth/stores/useAuthStore";
import { setSecureCookie } from "@/lib/cookie";
import { fullAuthCleanup } from "@/features/auth/utils/clear-auth-cookies";

export default function RouteError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  const [authStatus, setAuthStatus] = useState<'checking' | 'authenticated' | 'unauthenticated'>('checking');
  const { setUser, setToken, setSessionVerified, logout } = useAuthStore();

  useEffect(() => {
    // Log the error to an error reporting service
    console.error("Route error:", error);
  }, [error]);

  useEffect(() => {
    const verifySession = async () => {
      if (globalThis.window === undefined) {
        setAuthStatus('unauthenticated');
        return;
      }

      const token = localStorage.getItem('token');
      const refreshToken = localStorage.getItem('refreshToken');

      if (!token || !refreshToken) {
        setAuthStatus('unauthenticated');
        return;
      }

      try {
        const response = await authService.refreshToken(refreshToken);

        if (response.success && response.data) {
          const { user, token: newToken, refresh_token: newRefreshToken } = response.data;

          localStorage.setItem('token', newToken);
          localStorage.setItem('refreshToken', newRefreshToken);
          setSecureCookie('token', newToken);

          setUser(user);
          setToken(newToken);
          useAuthStore.setState({ 
            refreshToken: newRefreshToken, 
            isAuthenticated: true 
          });
          setSessionVerified(true);

          setAuthStatus('authenticated');
        } else {
          throw new Error('Session verification failed');
        }
      } catch {
        logout();
        await fullAuthCleanup();
        setAuthStatus('unauthenticated');
      }
    };

    verifySession();
  }, [setUser, setToken, setSessionVerified, logout]);

  // Show loading spinner while verifying session
  if (authStatus === 'checking') {
    return (
      <div className="flex min-h-screen items-center justify-center p-6 bg-background">
        <div className="text-center space-y-4">
          <LoadingSpinner size="lg" />
          <p className="text-sm text-muted-foreground">Verifying session...</p>
        </div>
      </div>
    );
  }

  const handleHomeNavigation = () => {
    if (authStatus === 'authenticated') {
      globalThis.location.href = '/dashboard';
    } else {
      globalThis.location.href = '/login';
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center p-6">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>Something went wrong</CardTitle>
          <CardDescription>
            An unexpected error occurred. Please try again.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {process.env.NODE_ENV === "development" && (
            <div className="rounded-md bg-destructive/10 p-4">
              <p className="text-sm font-mono text-destructive">
                {error.message}
              </p>
              {error.digest && (
                <p className="mt-2 text-xs text-muted-foreground">
                  Error ID: {error.digest}
                </p>
              )}
            </div>
          )}
          <div className="flex gap-2">
            <Button onClick={reset} className="flex-1">
              Try again
            </Button>
            <Button
              onClick={handleHomeNavigation}
              variant="outline"
              className="flex-1"
            >
              {authStatus === 'authenticated' ? 'Go to Dashboard' : 'Go to Login'}
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
