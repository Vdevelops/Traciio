"use client";

// DEPRECATED: This file is kept for backward compatibility only
// Use useHasPermission from @/features/auth/providers/permissions-provider instead

import { useHasPermission as useOptimizedHasPermission } from "@/features/auth/providers/permissions-provider";

/**
 * @deprecated Use useHasPermission from @/features/auth/providers/permissions-provider for better performance
 */
export function useHasPermission(permissionCode: string): boolean {
  return useOptimizedHasPermission(permissionCode);
}

