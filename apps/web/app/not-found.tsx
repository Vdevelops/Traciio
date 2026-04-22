"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Home, LogIn } from "lucide-react";
import { LoadingSpinner } from "@/components/ui/loading-spinner";
import { authService } from "@/features/auth/services/authService";
import { useAuthStore } from "@/features/auth/stores/useAuthStore";
import { setSecureCookie } from "@/lib/cookie";
import { fullAuthCleanup } from "@/features/auth/utils/clear-auth-cookies";

export default function NotFound() {
  const [authStatus, setAuthStatus] = useState<'checking' | 'authenticated' | 'unauthenticated'>('checking');
  const { setUser, setToken, setSessionVerified, logout } = useAuthStore();

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
      <main className="flex min-h-screen items-center justify-center bg-sidebar px-4">
        <div className="text-center space-y-4">
          <LoadingSpinner size="lg" />
          <p className="text-sm text-muted-foreground">Verifying session...</p>
        </div>
      </main>
    );
  }

  return (
    <main className="flex min-h-screen items-center justify-center bg-sidebar px-4">
      <div className="flex flex-col items-center text-center space-y-4">
        {/* 404 Number with accent styling */}
        <div className="relative">
          <span className="text-8xl sm:text-9xl font-bold text-primary/20 select-none">
            404
          </span>
        </div>

        <p className="text-xs font-medium uppercase tracking-[0.3em] text-muted-foreground">
          404
        </p>

        <h1 className="text-lg sm:text-xl font-semibold text-foreground">
          Page not found
        </h1>

        <p className="max-w-md text-sm text-muted-foreground leading-relaxed">
          The page you&apos;re looking for doesn&apos;t exist or may have been moved.
        </p>

        <div className="flex items-center gap-3 pt-2">
          {authStatus === 'authenticated' ? (
            <Button asChild size="default" className="gap-2">
              <Link href="/en/dashboard">
                <Home className="h-4 w-4" />
                Back to dashboard
              </Link>
            </Button>
          ) : (
            <Button asChild size="default" className="gap-2">
              <Link href="/en/login">
                <LogIn className="h-4 w-4" />
                Go to login
              </Link>
            </Button>
          )}
        </div>
      </div>
    </main>
  );
}


