"use client";

import type React from "react";
import { usePathname } from "@/i18n/routing";
import { DashboardLayout } from "@/components/layouts/dashboard-layout";
import { useWebSocket } from "@/features/notifications/hooks/useWebSocket";

interface AppLayoutProps {
  readonly children: React.ReactNode;
}

// AppLayout wraps only authenticated pages with the main DashboardLayout (sidebar + header).
// Public routes (locale-scoped login page at "/[locale]/login") are rendered without the dashboard chrome.
export function AppLayout({ children }: AppLayoutProps) {
  // Locale-agnostic pathname from next-intl (e.g. "/login", "/dashboard")
  const pathname = usePathname();

  const publicRoutes: readonly string[] = ["/login"];
  const isPublicRoute = publicRoutes.includes(pathname);

  // Always call useWebSocket hook (Rules of Hooks) - it will handle connection internally
  useWebSocket();

  if (isPublicRoute) {
    return <>{children}</>;
  }

  return <DashboardLayout>{children}</DashboardLayout>;
}

