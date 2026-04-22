"use client";

import React, { useEffect, useState } from "react";
import { useTranslations } from "next-intl";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { motion } from "framer-motion";
import { Eye, EyeOff } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Checkbox } from "@/components/ui/checkbox";
import { AuthLayout } from "./auth-layout";
import { useAuthStore } from "../stores/useAuthStore";
import { loginSchema, type LoginFormData } from "../schemas/login.schema";
import { useLogin } from "../hooks/useLogin";
import type { AuthError } from "../types/errors";
import { useRateLimitCountdown } from "@/lib/hooks/useRateLimitCountdown";
import { useRateLimitStore } from "@/lib/stores/useRateLimitStore";
import { fullAuthCleanup } from "../utils/clear-auth-cookies";
import { useRouter } from "@/i18n/routing";
import { authService } from "../services/authService";
import { setSecureCookie } from "@/lib/cookie";
import { LoadingSpinner } from "@/components/ui/loading-spinner";

export function LoginForm() {
  const t = useTranslations("auth.login");
  const router = useRouter();
  const { setUser, setToken, setSessionVerified, logout } = useAuthStore();
  const { handleLogin, isLoading, error, clearError } = useLogin();
  const [showPassword, setShowPassword] = useState(false);
  
  // Session verification state
  // - 'checking': Initial state, verifying session with backend
  // - 'authenticated': Session valid, redirect to dashboard
  // - 'unauthenticated': No session or invalid, show login form
  const [authStatus, setAuthStatus] = useState<'checking' | 'authenticated' | 'unauthenticated'>('checking');
  
  // Rate limit countdown hook - shows toast notification with countdown
  useRateLimitCountdown();
  
  // Get countdown text for display in form - update every second
  const resetTime = useRateLimitStore((state) => state.resetTime);
  const getCountdownText = useRateLimitStore((state) => state.getCountdownText);
  
  // Use tick state to trigger re-render every second for countdown updates
  // This avoids calling Date.now() during render and avoids synchronous setState in effects
  // The tick value is not used, only setTick is used to trigger re-renders
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [tick, setTick] = useState(0);
  
  useEffect(() => {
    if (!resetTime) {
      return;
    }
    
    // Update tick every second to trigger re-render and recalculate countdown
    // setTick from useState is stable and doesn't need to be in dependencies
    const interval = setInterval(() => {
      setTick((prev) => prev + 1);
    }, 1000);
    
    return () => clearInterval(interval);
  }, [resetTime]);

  /**
   * Verify session with backend on mount
   * - If refresh token is valid (200 OK) -> redirect to dashboard
   * - If 401/403 or no token -> show login form
   * - While checking -> show loading spinner
   */
  useEffect(() => {
    const verifySession = async () => {
      if (globalThis.window === undefined) {
        setAuthStatus('unauthenticated');
        return;
      }

      const token = localStorage.getItem('token');
      const refreshToken = localStorage.getItem('refreshToken');

      // No tokens stored - show login form immediately
      if (!token || !refreshToken) {
        setAuthStatus('unauthenticated');
        return;
      }

      try {
        // Call refresh token endpoint to verify session
        const response = await authService.refreshToken(refreshToken);

        if (response.success && response.data) {
          const { user, token: newToken, refresh_token: newRefreshToken } = response.data;

          // Update localStorage with new tokens
          localStorage.setItem('token', newToken);
          localStorage.setItem('refreshToken', newRefreshToken);
          setSecureCookie('token', newToken);

          // Update store with verified user data
          setUser(user);
          setToken(newToken);
          useAuthStore.setState({ 
            refreshToken: newRefreshToken, 
            isAuthenticated: true 
          });
          setSessionVerified(true);

          // Session valid - redirect to dashboard
          setAuthStatus('authenticated');
          router.push('/dashboard');
        } else {
          throw new Error('Session verification failed');
        }
      } catch {
        // Session invalid (401/403) - clear auth and show login form
        logout();
        await fullAuthCleanup();
        setAuthStatus('unauthenticated');
      }
    };

    verifySession();
  }, [router, setUser, setToken, setSessionVerified, logout]);
  
  // Calculate countdown text and rate limited status
  // getCountdownText() is safe to call here because it's called during render
  // and the tick state ensures it updates every second
  const countdownText = resetTime ? (() => {
    const text = getCountdownText();
    if (text === "a moment") {
      return null;
    }
    return text;
  })() : null;
  
  const isRateLimited = countdownText !== null;

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    setError,
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: "",
      password: "",
      rememberMe: false,
    },
  });

  useEffect(() => {
    if (error) {
      setError("root", {
        message: error,
      });
      clearError();
    }
  }, [error, setError, clearError]);

  const onSubmit = async (data: LoginFormData) => {
    try {
      await handleLogin(data);
      // rememberMe value is available here if needed later
    } catch (err) {
      const errorValue = err as AuthError;
      const errorMessage =
        errorValue.response?.data?.error?.message ||
        errorValue.message ||
        "Login failed. Please try again.";
      setError("root", {
        message: errorMessage,
      });
    }
  };

  const isFormLoading = isLoading || isSubmitting;

  // Show loading spinner while verifying session with backend
  if (authStatus === 'checking') {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background">
        <div className="text-center space-y-4">
          <LoadingSpinner size="lg" />
          <p className="text-sm text-muted-foreground">Verifying session...</p>
        </div>
      </div>
    );
  }

  // User is authenticated - show nothing while redirecting to dashboard
  if (authStatus === 'authenticated') {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background">
        <div className="text-center space-y-4">
          <LoadingSpinner size="lg" />
          <p className="text-sm text-muted-foreground">Redirecting to dashboard...</p>
        </div>
      </div>
    );
  }

  // User is unauthenticated - show login form
  return (
    <AuthLayout>
      <motion.div
        initial={{ opacity: 0, y: 10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3 }}
        className="w-full"
      >
        <Card className="border border-border/60 bg-card/90 shadow-sm">
          <CardHeader className="space-y-2 px-6 pb-2 pt-6">
            <CardTitle className="text-2xl">{t("title")}</CardTitle>
            <CardDescription className="text-sm text-muted-foreground">
              {t("description")}
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-5 px-6 pb-6 pt-2">
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
              <FieldGroup className="space-y-4">
                <Field className="space-y-2">
                  <FieldLabel htmlFor="email">{t("emailLabel")}</FieldLabel>
                  <Input
                    id="email"
                    type="email"
                    placeholder={t("emailPlaceholder")}
                    {...register("email")}
                    disabled={isFormLoading}
                    aria-invalid={!!errors.email}
                    className="h-11"
                  />
                  {errors.email && (
                    <FieldError>{errors.email.message}</FieldError>
                  )}
                </Field>

                <Field className="space-y-2">
                  <div className="flex items-center justify-between">
                    <FieldLabel htmlFor="password">{t("passwordLabel")}</FieldLabel>
                    <button
                      type="button"
                      className="text-xs font-medium text-primary hover:underline"
                    >
                      {t("forgotPassword")}
                    </button>
                  </div>
                  <div className="relative">
                    <Input
                      id="password"
                      type={showPassword ? "text" : "password"}
                      placeholder={t("passwordPlaceholder")}
                      {...register("password")}
                      disabled={isFormLoading}
                      aria-invalid={!!errors.password}
                      className="h-11 pr-10"
                    />
                    <button
                      type="button"
                      onClick={() => setShowPassword(!showPassword)}
                      disabled={isFormLoading}
                      className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                      aria-label={showPassword ? "Hide password" : "Show password"}
                    >
                      {showPassword ? (
                        <EyeOff className="h-4 w-4" />
                      ) : (
                        <Eye className="h-4 w-4" />
                      )}
                    </button>
                  </div>
                  {errors.password && (
                    <FieldError>{errors.password.message}</FieldError>
                  )}
                </Field>

              <Field>
                <label className="flex cursor-pointer items-center gap-2">
                  <Checkbox
                    {...register("rememberMe")}
                    disabled={isFormLoading}
                  />
                  <span className="text-sm text-muted-foreground">
                    {t("rememberMe")}
                  </span>
                </label>
              </Field>

                {errors.root && (
                  <Field>
                    <FieldError>{errors.root.message}</FieldError>
                  </Field>
                )}

                {/* Rate limit countdown display */}
                {countdownText && resetTime && (
                  <Field>
                    <div className="rounded-md bg-muted/50 px-3 py-2 text-sm text-muted-foreground">
                      <p className="font-medium">
                        {t("rateLimitMessage", { countdown: countdownText }) || 
                         `Too many login attempts. Please try again in ${countdownText}.`}
                      </p>
                    </div>
                  </Field>
                )}

                <Field className="pt-1">
                  <Button
                    type="submit"
                    disabled={isFormLoading || isRateLimited}
                    className="h-11 w-full text-sm font-medium tracking-wide"
                  >
                    {isFormLoading ? t("submitting") : t("submit")}
                  </Button>
                </Field>
              </FieldGroup>
            </form>
          </CardContent>
        </Card>
      </motion.div>
    </AuthLayout>
  );
}
