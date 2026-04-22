"use client";

import React, { useMemo } from "react";
import { usePathname, Link } from "@/i18n/routing";
import { ChevronRight } from "lucide-react";
import { getMenuIcon } from "@/lib/menu-icons";
import type { Menu } from "@/features/master-data/user-management/types";
import { useMenus } from "@/features/master-data/user-management/hooks/useMenus";

interface BreadcrumbItem {
  readonly label: string;
  readonly href: string;
  readonly icon: React.ReactNode;
}

interface BreadcrumbProps {
  readonly navigationItems?: Array<{
    readonly name: string;
    readonly href: string;
    readonly icon: React.ReactNode;
  }>;
}

function findMenuByUrl(
  menus: Menu[],
  url: string
): Menu | null {
  for (const menu of menus) {
    if (menu.url === url) {
      return menu;
    }
    if (menu.children && menu.children.length > 0) {
      const found = findMenuByUrl(menu.children, url);
      if (found) {
        return found;
      }
    }
  }
  return null;
}

function getRouteLabel(path: string): string {
  // Convert path segments to readable labels
  const labelMap: Record<string, string> = {
    dashboard: "Dashboard",
    accounts: "Accounts",
    leads: "Leads",
    deals: "Deals",
    pipeline: "Pipeline",
    products: "Products",
    tasks: "Tasks",
    "visit-reports": "Visit Reports",
    reports: "Reports",
    "sales-overview": "Sales Overview",
    "sales-rep": "Sales Rep",
    leaderboard: "Leaderboard",
    "master-data": "Master Data",
    users: "Users",
    profile: "Profile",
    "ai-chatbot": "AI Chatbot",
    block: "Block",
  };

  // Remove leading slash and split
  const segments = path.split("/").filter(Boolean);
  const lastSegment = segments.at(-1) ?? "";

  // Check if it's a UUID or ID (dynamic route)
  const isDynamicRoute = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(
    lastSegment
  ) || /^\d+$/.test(lastSegment);

  if (isDynamicRoute && segments.length > 1) {
    // For dynamic routes, use the parent segment label
    const parentSegment = segments[segments.length - 2] ?? "";
    return labelMap[parentSegment] ?? parentSegment;
  }

  return labelMap[lastSegment] ?? lastSegment;
}

export function Breadcrumb({ navigationItems }: BreadcrumbProps) {
  const pathname = usePathname();
  const { data: menusData } = useMenus();

  const breadcrumbItems = useMemo<BreadcrumbItem[]>(() => {
    if (!pathname) {
      return [];
    }

    // Skip if at root or dashboard (no breadcrumb needed)
    if (pathname === "/dashboard" || pathname === "/") {
      return [];
    }

    // Always start with Home/Dashboard
    const items: BreadcrumbItem[] = [
      {
        label: "Dashboard",
        href: "/dashboard",
        icon: getMenuIcon("home"),
      },
    ];

    // Split pathname into segments
    const segments = pathname.split("/").filter(Boolean);

    // Build breadcrumb items from segments
    let currentPath = "";
    const menus = menusData?.data ?? [];

    for (let i = 0; i < segments.length; i++) {
      const segment = segments[i] ?? "";
      currentPath += `/${segment}`;

      // Skip dashboard segment (already added)
      if (currentPath === "/dashboard") {
        continue;
      }

      // Try to find menu from permissions data
      const menu = findMenuByUrl(menus, currentPath);
      let label = menu?.name ?? getRouteLabel(currentPath);
      let icon: React.ReactNode = menu
        ? getMenuIcon(menu.icon)
        : getMenuIcon(segment);

      // If not found in menus, try navigationItems
      if (!menu && navigationItems) {
        const navItem = navigationItems.find((item) => item.href === currentPath);
        if (navItem) {
          label = navItem.name;
          icon = navItem.icon;
        }
      }

      // For dynamic routes (UUIDs or numbers), use parent label + "Detail"
      const isDynamicSegment =
        /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(
          segment
        ) || /^\d+$/.test(segment);

      if (isDynamicSegment && i > 0) {
        const parentSegment = segments.at(-2) ?? "";
        const parentPath = `/${parentSegment}`;
        const parentMenu = findMenuByUrl(menus, parentPath);
        label = `${parentMenu?.name ?? getRouteLabel(parentPath)} Detail`;
        icon = parentMenu ? getMenuIcon(parentMenu.icon) : getMenuIcon(parentSegment);
      }

      items.push({
        label,
        href: currentPath,
        icon,
      });
    }

    return items;
  }, [pathname, menusData, navigationItems]);

  // Don't show breadcrumb if empty or only dashboard
  if (breadcrumbItems.length === 0 || breadcrumbItems.length <= 1) {
    return null;
  }

  return (
    <nav
      className="flex items-center gap-2 border-b bg-background/50 px-4 py-2 text-sm text-muted-foreground backdrop-blur-sm"
      aria-label="Breadcrumb"
    >
      <ol className="flex items-center gap-2">
        {breadcrumbItems.map((item, index) => {
          const isLast = index === breadcrumbItems.length - 1;

          return (
            <li key={item.href} className="flex items-center gap-2">
              {index > 0 && (
                <ChevronRight className="h-4 w-4 text-muted-foreground/50" />
              )}
              {isLast ? (
                <div className="flex items-center gap-2 font-medium text-foreground">
                  <span className="flex items-center justify-center text-primary">
                    {item.icon}
                  </span>
                  <span>{item.label}</span>
                </div>
              ) : (
                <Link
                  href={item.href}
                  className="flex items-center gap-2 hover:text-foreground transition-colors cursor-pointer"
                >
                  <span className="flex items-center justify-center text-muted-foreground">
                    {item.icon}
                  </span>
                  <span>{item.label}</span>
                </Link>
              )}
            </li>
          );
        })}
      </ol>
    </nav>
  );
}
