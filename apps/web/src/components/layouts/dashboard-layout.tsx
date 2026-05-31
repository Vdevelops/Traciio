"use client";

import React, { memo, useMemo, useEffect } from "react";
import Image from "next/image";
import { Link, usePathname, useRouter } from "@/i18n/routing";
import { useSearchParams } from "next/navigation";
import { HelpCircle, Search } from "lucide-react";

import { usePermissions } from "@/features/auth/providers/permissions-provider";
import { useRoleValidation } from "@/features/auth/hooks/useRoleValidation";
import { useIsMobile } from "@/hooks/use-mobile";
import { NAVIGATION_CONFIG } from "@/lib/navigation-config";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { getMenuIcon } from "@/lib/menu-icons";
import { CommandPalette } from "@/components/ui/command-palette";
import { useDashboardCommandPalette } from "@/hooks/useDashboardCommandPalette";
import { NotificationDrawer } from "@/features/notifications/components/notification-drawer";
import { useNotificationStore } from "@/features/notifications/stores/useNotificationStore";
import { Breadcrumb } from "@/components/navigation/breadcrumb";
import { HeaderControls } from "@/components/ui/header-controls";
import {
  Sidebar,
  SidebarContent,
  SidebarHeader,
  SidebarInset,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
  SidebarTrigger,
  SidebarGroup,
  SidebarGroupLabel,
  SidebarGroupContent,
  useSidebar,
} from "@/components/ui/sidebar";

// Interface for what the Sidebar component actually consumes (runtime)
interface RuntimeNavigationItem {
  name: string;
  href: string;
  icon: React.ReactNode;
  group: string;
  permission?: string;
  children?: RuntimeNavigationItem[];
}

interface DashboardLayoutProps {
  readonly children: React.ReactNode;
}

const Header = memo(function Header() {
  return (
    <header className="surface-panel sticky top-0 z-50 mx-2 mt-2 flex h-16 shrink-0 items-center gap-3 rounded-[1.35rem] px-4">
      <SidebarTrigger className="-ml-1 size-8" />
      <Separator orientation="vertical" />

      <div className="lg:flex-1">
        {/* Desktop search input */}
        <div className="relative hidden max-w-sm flex-1 lg:block">
          <Search
            className="text-muted-foreground pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2"
            aria-hidden="true"
          />
          <input
            type="search"
            placeholder="Search..."
            className="file:text-foreground placeholder:text-muted-foreground selection:bg-primary selection:text-primary-foreground border-input h-10 w-full cursor-pointer rounded-2xl border bg-card/80 px-3.5 py-1 pr-4 pl-10 text-sm shadow-sm outline-none transition-[color,box-shadow] focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50"
          />
          <div className="surface-muted text-muted-foreground absolute right-2 top-1/2 hidden -translate-y-1/2 items-center gap-0.5 rounded-md px-1.5 py-0.5 font-mono text-[10px] font-medium sm:flex">
            <span>⌘</span>
            <span>K</span>
          </div>
        </div>

        {/* Mobile search button */}
        <div className="block lg:hidden">
          <Button
            variant="ghost"
            size="icon"
            className="size-9"
            type="button"
          >
            <Search className="h-4 w-4" aria-hidden="true" />
            <span className="sr-only">Open search</span>
          </Button>
        </div>
      </div>

      <div className="ml-auto flex items-center gap-1">
        <HeaderControls
          showNotifications
          showThemeToggle
          showLocaleToggle
          showProfile
          extraIcon={<HelpCircle className="h-4 w-4 text-muted-foreground" />}
        />
      </div>
    </header>
  );
});

const AppSidebar = memo(function AppSidebar({
  items,
}: {
  items: RuntimeNavigationItem[];
}) {

  const pathname = usePathname();
  const searchParams = useSearchParams();
  const { state } = useSidebar();

  // Group items by simplified category
  const grouped = useMemo(() => {
    const groups: Record<string, RuntimeNavigationItem[]> = {};
    for (const item of items) {
      const category = item.group || "Main";
      if (!groups[category]) {
        groups[category] = [];
      }
      groups[category].push(item);
    }
    return groups;
  }, [items]);

  // Enhanced active logic: only one active (users or roles)
  const isActive = (href: string) => {
    // Special handling for /master-data/users which may have ?tab=...
    if (pathname === "/master-data/users") {
      const tab = searchParams?.get("tab");

      // Attempt to parse the href query param (robust for relative hrefs)
      let hrefTab: string | null = null;
      try {
        const url = new URL(href, typeof window !== "undefined" ? window.location.origin : "http://localhost");
        if (url.pathname === "/master-data/users") {
          hrefTab = url.searchParams.get("tab");
        }
      } catch {
        hrefTab = null;
      }

      if (hrefTab) {
        return tab === hrefTab;
      }

      if (href === "/master-data/users") {
        return !tab || tab === "users";
      }
    }

    // Default: match path or subpath
    return pathname === href || pathname.startsWith(`${href}/`);
  };

  return (
    <Sidebar variant="inset" collapsible="icon">
      <SidebarHeader className="border-b border-sidebar-border/60 pb-3">
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild tooltip="Tracio">
              <Link href="/dashboard" className="flex w-full items-center gap-2">
                {state === "collapsed" ? (
                  <div className="flex w-full items-center justify-center rounded-xl bg-sidebar-accent/50 py-2">
                    <Image
                      src="/tracio-logo.svg"
                      alt="Tracio"
                      width={32}
                      height={32}
                      className="h-8 w-8 object-contain"
                    />
                  </div>
                ) : (
                  <div className="flex w-full items-center gap-2 rounded-2xl bg-sidebar-accent/45 px-2 py-2">
                    <Image
                      src="/tracio-logo.svg"
                      alt="Tracio"
                      width={36}
                      height={36}
                      className="h-9 w-9 object-contain"
                    />
                    <div className="flex min-w-0 flex-col leading-tight">
                      <span className="brand-text text-sm font-semibold tracking-tight">Tracio</span>
                      <span className="truncate text-[10px] text-muted-foreground">Track Better, Serve Smarter</span>
                    </div>
                  </div>
                )}
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>

      <SidebarContent className="gap-0">
        {Object.entries(grouped).map(([group, groupItems]) => (
          <SidebarGroup key={group} className="py-2">
            {group !== "Main" && (
              <SidebarGroupLabel className="mb-1 px-2 text-[10px] font-semibold uppercase tracking-[0.2em] text-muted-foreground/75">
                {group}
              </SidebarGroupLabel>
            )}
            <SidebarGroupContent>
              <SidebarMenu className="gap-1">
                {groupItems.map((item) => {
                  const active = isActive(item.href);
                  return (
                    <SidebarMenuItem key={item.href}>
                      <SidebarMenuButton
                        asChild
                        isActive={active}
                        tooltip={item.name}
                        className="h-10"
                      >
                        {/* prefetch=false: sidebar links must not eagerly prefetch all routes
                            on viewport entry — this was causing ~80 RSC requests per page load. */}
                        <Link href={item.href} prefetch={false}>
                          {item.icon}
                          <span>{item.name}</span>
                        </Link>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                  );
                })}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        ))}
      </SidebarContent>


    </Sidebar>
  );
});

const AutoCollapseSidebar = memo(function AutoCollapseSidebar() {
  const isMobile = useIsMobile();
  const { setOpen } = useSidebar();

  useEffect(() => {
    if (isMobile) {
      setOpen(false);
    }
  }, [isMobile, setOpen]);

  // Removed: Auto minimize sidebar when on AI chatbot page
  // The AI chatbot now has its own internal sidebar for chat history

  return null;
});

const FullScreenLayout = memo(function FullScreenLayout({
  isFullScreenPage,
  navigationItems,
  error,
  children,
}: {
  isFullScreenPage: boolean;
  navigationItems: RuntimeNavigationItem[];
  error: Error | null;
  children: React.ReactNode;
}) {
  return (
    <div className="page-frame flex min-h-screen w-full bg-background">
      <AppSidebar items={navigationItems} />
      <SidebarInset 
        className={`overflow-x-hidden transition-all duration-300 ${
          isFullScreenPage ? "w-full" : ""
        }`}
      >
        {!isFullScreenPage && (
          <>
            <Header />
            <Breadcrumb navigationItems={navigationItems} />
          </>
        )}
        <div 
          className={`flex flex-1 flex-col ${
            isFullScreenPage 
              ? "gap-0 p-0 h-screen" 
              : "gap-4 p-4"
          }`}
        >
          {!isFullScreenPage && error && (
            <div className="mb-2 rounded-2xl border border-destructive/30 bg-destructive/5 px-3 py-2 text-xs text-destructive">
              Failed to load menu permissions. Showing minimal navigation.
            </div>
          )}
          {children}
        </div>
      </SidebarInset>
    </div>
  );
});

export const DashboardLayout = memo(function DashboardLayout({
  children,
}: DashboardLayoutProps) {
  const { error, hasPermission } = usePermissions();
  // Validate role and auto logout if role is missing
  useRoleValidation();
  const { isDrawerOpen, closeDrawer } = useNotificationStore();

  const navigationItems: RuntimeNavigationItem[] = useMemo(() => {
    // Get permission list from context
    const items: RuntimeNavigationItem[] = [];

    // Iterate over static config and filter by permission
    NAVIGATION_CONFIG.forEach((group) => {
        group.items.forEach((item) => {
            // Check item permission
            let hasPermissionForItem = false;
            // Strict check: if permission defined, user MUST have it
            if (!item.permission) {
              hasPermissionForItem = true; 
            } else if (hasPermission(item.permission)) {
              hasPermissionForItem = true;
            }

            if (hasPermissionForItem) {
                // If item has children, filter them as well
                let validChildren: RuntimeNavigationItem[] | undefined = undefined;
                
                if (item.children && item.children.length > 0) {
                     // Filter children and map to runtime type
                     const filtered = item.children.filter(child => 
                  !child.permission || hasPermission(child.permission)
                    );
                    
                    if (filtered.length > 0) {
                        validChildren = filtered.map(child => ({
                            name: child.name,
                            href: child.href,
                            icon: getMenuIcon(child.icon),
                            group: group.label,
                            permission: child.permission,
                        }));
                    }
                }

                 // Add item if allowed
                items.push({
                    name: item.name,
                    href: item.href,
                    icon: getMenuIcon(item.icon), // Convert string icon key to ReactNode
                    group: group.label,         // Inject group label
                    permission: item.permission,
                    children: validChildren
                });
            }
        });
    });

    return items;
  }, [hasPermission]);

  const commandPalette = useDashboardCommandPalette();
  const router = useRouter();

  const handleSelectItem = (href: string) => {
    if (!href) return;
    router.push(href);
    commandPalette.close();
  };

  const pathname = usePathname();
  const searchParams = useSearchParams();
  // Full-screen pages: no header, no breadcrumb, no padding
  const isFullScreenPage = pathname?.includes("/ai-chatbot") || 
    pathname?.includes("/route-optimization") ||
    // Accounts map view (main page, no tab selected)
    (pathname?.endsWith("/accounts") && !searchParams?.get("tab")) ||
    (pathname?.includes("/master-data/bricks") && !pathname?.includes("/bricks/"));

  return (
    <SidebarProvider>
      <AutoCollapseSidebar />
      <FullScreenLayout 
        isFullScreenPage={!!isFullScreenPage}
        navigationItems={navigationItems}
        error={error}
      >
        {children}
      </FullScreenLayout>

      {/* Notification Drawer */}
      <NotificationDrawer open={isDrawerOpen} onOpenChange={closeDrawer} />

      {/* Command Palette */}
      <CommandPalette
        open={commandPalette.isOpen}
        onOpenChange={commandPalette.toggle}
        items={navigationItems}
        onSelectItem={handleSelectItem}
      />
    </SidebarProvider>
  );
});

