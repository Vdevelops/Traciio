"use client";

import { lazy, Suspense } from "react";
import type { PermissionsProviderProps } from "./permissions-provider";

// Lazy load the permissions provider for better code splitting
const LazyPermissionsProviderComponent = lazy(() => 
  import("./permissions-provider").then(module => ({ 
    default: module.PermissionsProvider 
  }))
);

/**
 * Permissions loader with code splitting and preloading
 * Provides better performance for large scale deployment
 */
export function LazyPermissionsProvider(props: PermissionsProviderProps) {
  return (
    <Suspense 
      fallback={
        <div className="min-h-screen bg-background" aria-label="Loading permissions...">
          {/* Minimal loading state to prevent layout shift */}
          <div className="sr-only">Loading user permissions...</div>
        </div>
      }
    >
      <LazyPermissionsProviderComponent {...props} />
    </Suspense>
  );
}

export default LazyPermissionsProvider;