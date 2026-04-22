"use client";

import { useEffect, useRef } from "react";
import { useRouter } from "@/i18n/routing";
import { useAuthStore } from "../stores/useAuthStore";
import { usePermissions } from "@/features/auth/providers/permissions-provider";
import { useLogout } from "./useLogout";
import { toast } from "sonner";
import { useTranslations } from "next-intl";

/**
 * Hook to validate user role and auto logout if role is invalid or missing
 * This hook should be used in layouts/components that require authenticated users
 */
export function useRoleValidation() {
  const router = useRouter();
  const { user } = useAuthStore();
  const logout = useLogout();
  const { permissions, error, isLoading } = usePermissions();
  const hasLoggedOutRef = useRef(false);
  const t = useTranslations("auth");

  useEffect(() => {
    // Don't check if already logged out or no user
    if (hasLoggedOutRef.current || !user?.id || isLoading) {
      return;
    }

    // Check for errors that indicate role is missing/invalid
    if (error) {
      const axiosError = error as {
        response?: {
          status?: number;
          data?: {
            error?: {
              code?: string;
              message?: string;
            };
          };
        };
      };
      
      // Check if error indicates role/user is missing
      if (
        axiosError.response?.data?.error?.code === "USER_NOT_FOUND" ||
        axiosError.response?.status === 404 ||
        (axiosError.response?.status === 401 && 
         axiosError.response?.data?.error?.message?.toLowerCase().includes("role"))
      ) {
        hasLoggedOutRef.current = true;
        
        // Show toast notification
        toast.error(t("roleMissing.title") || "Role Missing", {
          description: t("roleMissing.description") || "Your role has been removed. Please contact administrator.",
          duration: 5000,
        });

        // Auto logout after a short delay
        setTimeout(() => {
          logout();
        }, 2000);
        return;
      }
    }

    // If permissions loaded successfully, check if user has valid role
    if (permissions && !error) {
      // User has valid role and permissions
      hasLoggedOutRef.current = false;
      return;
    }
  }, [user, error, permissions, isLoading, logout, router, t]);

  return {
    isValid: !error && !!permissions,
    isLoading,
  };
}

