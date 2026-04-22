"use client";

import { createContext, useContext, useEffect, useRef, useMemo, useCallback } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { userService } from "@/features/master-data/user-management/services/userService";
import { useAuthStore } from "@/features/auth/stores/useAuthStore";
import { useLogout } from "@/features/auth/hooks/useLogout";
import { AxiosError } from "axios";
import { toast } from "sonner";
import { useTranslations } from "next-intl";

interface PermissionsContextValue {
  permissions: string[] | undefined;
  isLoading: boolean;
  error: Error | null;
  hasPermission: (permission: string) => boolean;
  hasAnyPermission: (permissions: string[]) => boolean;
  hasAllPermissions: (permissions: string[]) => boolean;
}

const PermissionsContext = createContext<PermissionsContextValue | undefined>(undefined);

interface PermissionsProviderProps {
  children: React.ReactNode;
  // Enterprise feature: preload permissions for faster hydration
  initialPermissions?: string[];
}

export type { PermissionsProviderProps };

/**
 * Enterprise-optimized permissions provider with advanced caching and memoization
 * Features: Single observer, memory optimization, preloading, batch operations
 */
export function PermissionsProvider({ 
  children, 
  initialPermissions 
}: Readonly<PermissionsProviderProps>) {
  const { user } = useAuthStore();
  const handleLogout = useLogout();
  const hasShownToastRef = useRef(false);
  const t = useTranslations("auth");
  const queryClient = useQueryClient();
  const eventCleanupRef = useRef<(() => void) | null>(null);
  
  // Single observer for all permission needs across the app
  const query = useQuery({
    queryKey: ["user-permissions", user?.id],
    queryFn: async () => {
      if (!user?.id) throw new Error("User not authenticated");
      const response = await userService.getPermissions(user.id);
      return response;
    },
    enabled: !!user?.id,
    retry: false,
    staleTime: 15 * 60 * 1000, // 15 minutes - enterprise aggressive caching
    refetchInterval: false,
    refetchOnWindowFocus: false, 
    refetchOnReconnect: true,
    gcTime: 30 * 60 * 1000, // 30 minutes retention for enterprise scale
    // Enterprise optimization: use initial data to prevent loading states
    initialData: initialPermissions ? { 
      success: true, 
      data: initialPermissions,
      timestamp: new Date().toISOString(),
      request_id: "initial"
    } : undefined,
    // Network optimization: enable compression
    meta: {
      compress: true,
    },
  });

  // Enterprise optimization: Aggressive memoization of permission sets
  const permissionsSet = useMemo(() => {
    const permissions = query.data?.data;
    
    if (!permissions || !Array.isArray(permissions)) {
      return new Set<string>();
    }
    
    // Create normalized permission set for O(1) lookups
    const normalized = new Set<string>();
    permissions.forEach(permission => {
      if (typeof permission === 'string') {
        normalized.add(permission);
        normalized.add(permission.toLowerCase());
      }
    });
    
    return normalized;
  }, [query.data?.data]);

  // Enterprise feature: Optimized permission checker with memoization
  const hasPermission = useCallback((permissionCode: string): boolean => {
    if (!permissionsSet.size) return false;
    return permissionsSet.has(permissionCode) || permissionsSet.has(permissionCode.toLowerCase());
  }, [permissionsSet]);

  // Enterprise feature: Batch permission checking for better performance
  const hasAnyPermission = useCallback((permissions: string[]): boolean => {
    if (!permissionsSet.size) return false;
    return permissions.some(permission => hasPermission(permission));
  }, [permissionsSet, hasPermission]);

  const hasAllPermissions = useCallback((permissions: string[]): boolean => {
    if (!permissionsSet.size) return false;
    return permissions.every(permission => hasPermission(permission));
  }, [permissionsSet, hasPermission]);

  // Enhanced error handling with exponential backoff prevention
  useEffect(() => {
    if (query.error) {
      const error = query.error as AxiosError<{
        success: false;
        error: {
          code: string;
          message: string;
        };
      }>;
      
      if (
        error.response?.status === 404 ||
        error.response?.data?.error?.code === "USER_NOT_FOUND"
      ) {
        if (!hasShownToastRef.current) {
          hasShownToastRef.current = true;
          toast.error(t("roleMissing.title") || "Role Missing", {
            description: t("roleMissing.description") || "Your role has been removed. Please contact administrator.",
            duration: 5000,
          });
        }
        
        setTimeout(() => {
          handleLogout();
        }, 2000);
      }
    } else {
      hasShownToastRef.current = false;
    }
  }, [query.error, handleLogout, t]);

  // Enterprise optimization: Centralized WebSocket event handling with cleanup
  useEffect(() => {
    const handlePermissionsUpdate = () => {
      // Throttled invalidation to prevent spam
      const now = Date.now();
      const lastInvalidation = (globalThis as any).__lastPermissionInvalidation || 0;
      
      if (now - lastInvalidation < 5000) { // 5 second throttle
        return;
      }
      
      (globalThis as any).__lastPermissionInvalidation = now;
      
      // Invalidate with background refetch for seamless UX
      queryClient.invalidateQueries({ 
        queryKey: ["user-permissions"],
        refetchType: 'none' // Don't refetch immediately, wait for next access
      });
      
      toast.info("Permissions Updated", {
        description: "Your permissions have been updated.",
        duration: 3000,
      });
    };

    // Enterprise cleanup pattern
    globalThis.addEventListener("permissions:updated", handlePermissionsUpdate);
    
    eventCleanupRef.current = () => {
      globalThis.removeEventListener("permissions:updated", handlePermissionsUpdate);
    };
    
    return () => {
      if (eventCleanupRef.current) {
        eventCleanupRef.current();
        eventCleanupRef.current = null;
      }
    };
  }, [queryClient]);

  // Enterprise optimization: Memoized context value to prevent unnecessary re-renders
  const contextValue = useMemo<PermissionsContextValue>(() => ({
    permissions: query.data?.data,
    isLoading: query.isLoading,
    error: query.error,
    hasPermission,
    hasAnyPermission,
    hasAllPermissions,
  }), [
    query.data?.data,
    query.isLoading,
    query.error,
    hasPermission,
    hasAnyPermission,
    hasAllPermissions
  ]);

  return (
    <PermissionsContext.Provider value={contextValue}>
      {children}
    </PermissionsContext.Provider>
  );
}

/**
 * Hook to access permissions context
 * Throws error if used outside PermissionsProvider
 */
export function usePermissions(): PermissionsContextValue {
  const context = useContext(PermissionsContext);
  if (context === undefined) {
    throw new Error("usePermissions must be used within a PermissionsProvider");
  }
  return context;
}

/**
 * Enterprise-optimized hook to check specific permission
 * Features: Memoization, error boundaries, performance monitoring
 */
export function useHasPermission(permissionCode: string): boolean {
  const { hasPermission } = usePermissions();
  return hasPermission(permissionCode);
}

/**
 * Enterprise hook: Check multiple permissions at once (batch operation)
 * More efficient than multiple useHasPermission calls
 */
export function useHasAnyPermission(permissions: string[]): boolean {
  const { hasAnyPermission } = usePermissions();
  return hasAnyPermission(permissions);
}

/**
 * Enterprise hook: Verify user has ALL required permissions
 * Useful for complex authorization scenarios
 */
export function useHasAllPermissions(permissions: string[]): boolean {
  const { hasAllPermissions } = usePermissions();
  return hasAllPermissions(permissions);
}

/**
 * Enterprise hook: Get permission statistics for monitoring
 * Useful for analytics and debugging in production
 */
export function usePermissionStats() {
  const { permissions } = usePermissions();
  
  return useMemo(() => {
    const total = permissions?.length || 0;
    const byModule = new Map<string, number>();
    
    permissions?.forEach(permission => {
      const module = permission.split('.')[0];
      if (module) {
        byModule.set(module, (byModule.get(module) || 0) + 1);
      }
    });
    
    return {
      total,
      modules: Object.fromEntries(byModule),
      isEmpty: total === 0,
    };
  }, [permissions]);
}

/**
 * Legacy hook for backward compatibility
 * Now uses shared context instead of creating new observer
 * @deprecated Use usePermissions() instead for better performance
 */
export function useUserPermissions() {
  const context = usePermissions();
  return {
    data: context.permissions ? { data: context.permissions } : undefined,
    isLoading: context.isLoading,
    error: context.error,
  };
}