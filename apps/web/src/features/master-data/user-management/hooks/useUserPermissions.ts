 "use client";

import { useQuery } from "@tanstack/react-query";
import { useRouter } from "@/i18n/routing";
import { useEffect, useRef } from "react";
import { userService } from "../services/userService";
import { useAuthStore } from "@/features/auth/stores/useAuthStore";
import { useLogout } from "@/features/auth/hooks/useLogout";
import { AxiosError } from "axios";
import { toast } from "sonner";
import { useTranslations } from "next-intl";

export function useUserPermissions() {
  // Ensure hooks are called in consistent order - router first, then store
  const router = useRouter();
  const { user } = useAuthStore();
  const handleLogout = useLogout();
  const hasShownToastRef = useRef(false);
  const t = useTranslations("auth");
  
  const query = useQuery({
    queryKey: ["user-permissions", user?.id],
    queryFn: async () => {
      if (!user?.id) throw new Error("User not authenticated");
      const response = await userService.getPermissions(user.id);
      return response;
    },
    enabled: !!user?.id,
    retry: false, // Don't retry on 404 - user doesn't exist
    staleTime: 5 * 60 * 1000, // 5 minutes - WS "permissions.updated" handles real-time invalidation
    refetchInterval: false, // Disable polling, rely on WS (permissions.updated) event
    refetchOnWindowFocus: false, // Rely on global default; WS handles real-time updates
    refetchOnReconnect: true, // Refetch when network reconnects
  });

  // Handle 404 error - user not found, clear auth and redirect
  useEffect(() => {
    if (query.error) {
      const error = query.error as AxiosError<{
        success: false;
        error: {
          code: string;
          message: string;
        };
      }>;
      
      // If user not found (404) or USER_NOT_FOUND error, clear auth state
      if (
        error.response?.status === 404 ||
        error.response?.data?.error?.code === "USER_NOT_FOUND"
      ) {
        // Show toast only once
        if (!hasShownToastRef.current) {
          hasShownToastRef.current = true;
          toast.error(t("roleMissing.title") || "Role Missing", {
            description: t("roleMissing.description") || "Your role has been removed. Please contact administrator.",
            duration: 5000,
          });
        }
        
        // Clear auth state and redirect to locale-scoped login
        setTimeout(() => {
          handleLogout();
        }, 2000);
      }
    } else {
      // Reset toast flag on success
      hasShownToastRef.current = false;
    }
  }, [query.error, handleLogout, t]);

  return query;
}

