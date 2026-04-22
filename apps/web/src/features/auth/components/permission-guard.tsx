"use client";

import { useRef, useEffect } from "react";
import { useRouter, usePathname } from "@/i18n/routing";
import { usePermissions, useHasPermission } from "@/features/auth/providers/permissions-provider";
import { toast } from "sonner";
import { useTranslations } from "next-intl";

interface PermissionGuardProps {
  readonly children: React.ReactNode;
  readonly requiredPermission: string;
  readonly fallbackUrl?: string;
}

/**
 * PermissionGuard component that checks if user has required permission
 * If not, redirects to block page with toast notification
 * Real-time permission checking - redirects immediately when permission is revoked
 */
export function PermissionGuard({
  children,
  requiredPermission,
  fallbackUrl = "/block",
}: PermissionGuardProps) {
  const router = useRouter();
  const pathname = usePathname();
  const { permissions, isLoading, error } = usePermissions();
  const hasPermission = useHasPermission(requiredPermission);
  const hasRedirectedRef = useRef(false);
  const hasShownToastRef = useRef(false);
  const t = useTranslations("auth");

  useEffect(() => {
    // If loading finished and user doesn't have permission, redirect
    if (!isLoading && !hasPermission && !hasRedirectedRef.current) {
      hasRedirectedRef.current = true;
      
      // Show toast notification
      if (!hasShownToastRef.current && !error) {
        hasShownToastRef.current = true;
        toast.error(t("permissionRevoked.title") || "Access Denied", {
          description: t("permissionRevoked.description") || "You no longer have permission to access this page.",
          duration: 3000,
        });
      }

      // Ensure fallbackUrl is absolute (starts with /)
      const absoluteBlockUrl = fallbackUrl.startsWith("/") 
        ? fallbackUrl 
        : `/${fallbackUrl}`;
      
      // Use router.replace for client-side navigation
      router.replace(absoluteBlockUrl);
    } else if (hasPermission) {
      // Reset redirect flag if permission is restored
      hasRedirectedRef.current = false;
      hasShownToastRef.current = false;
    }
  }, [hasPermission, isLoading, permissions, router, pathname, fallbackUrl, t, error]);

  // Show nothing while checking permissions
  if (isLoading) {
    return null;
  }

  // If no permission, don't render children (redirect will happen)
  if (!hasPermission) {
    return null;
  }

  return <>{children}</>;
}

