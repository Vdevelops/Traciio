"use client";

import { useState, useMemo } from "react";
import { toast } from "sonner";
import { useTranslations } from "next-intl";
import { Search } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { Drawer } from "@/components/ui/drawer";
import { Input } from "@/components/ui/input";
import { usePermissions } from "../hooks/usePermissions";
import { useRole, useAssignPermissionsToRole, useRolePermissions } from "../hooks/useRoles";
import { Checkbox } from "@/components/ui/checkbox";
import { cn } from "@/lib/utils";
import { getMenuIcon } from "@/lib/menu-icons";
import type { Permission } from "../types";

interface AssignPermissionsDrawerProps {
  readonly roleId: string;
  readonly onClose: () => void;
}

// Category mapping based on menu order ranges from seeder
const CATEGORY_MAP: Record<string, { name: string; order: number; iconName: string }> = {
  Dashboard: { name: "Dashboard", order: 0, iconName: "layout-dashboard" },
  Sales: { name: "Sales", order: 1, iconName: "shopping-cart" },
  Inventory: { name: "Inventory", order: 2, iconName: "package" },
  Customers: { name: "Customers", order: 3, iconName: "building" },
  Analytics: { name: "Analytics", order: 4, iconName: "bar-chart-3" },
  Management: { name: "Management", order: 5, iconName: "settings" },
  AI: { name: "AI", order: 6, iconName: "bot" },
};

// Map menu names to categories based on seeder structure
const MENU_TO_CATEGORY: Record<string, string> = {
  "Dashboard": "Dashboard",
  "Route Optimization": "Dashboard",
  "Leads": "Sales",
  "Pipeline": "Sales",
  "Schedules": "Sales",
  "Visits": "Sales",
  "Tasks": "Sales",
  "Products": "Inventory",
  "Accounts": "Customers",
  "Sales Performance": "Analytics",
  "Product Analytics": "Analytics",

  "Reports": "Analytics",
  "Users": "Management",
  "Groups": "Management",
  "Bricks": "Management",
  "Targets": "Management",
  "Chatbot": "AI",
  "AI Settings": "AI",
};

// Helper function to get action badge styling
function getActionBadgeClass(action: string): string {
  switch (action) {
    case "view":
    case "read":
      return "bg-muted text-muted-foreground";
    case "create":
      return "bg-success text-success-foreground";
    case "update":
    case "edit":
      return "bg-accent text-accent-foreground";
    case "delete":
      return "bg-destructive text-destructive-foreground";
    default:
      return "bg-muted text-muted-foreground";
  }
}

// Action badge component
function ActionBadge({ action }: { readonly action: string }) {
  return (
    <span className={cn("text-xs px-2 py-0.5 rounded-full capitalize", getActionBadgeClass(action))}>
      {action}
    </span>
  );
}

// Inner component that only renders when role data is loaded
function PermissionsSelector({
  permissions,
  rolePermissions,
  roleId,
  onClose,
}: {
  readonly permissions: Permission[];
  readonly rolePermissions: Permission[];
  readonly roleId: string;
  readonly onClose: () => void;
}) {
  const assignPermissions = useAssignPermissionsToRole();
  const t = useTranslations("userManagement.assignPermissions");

  // Initialize state from rolePermissions (fetched separately)
  const initialPermissions = useMemo(
    () => (rolePermissions?.length ? rolePermissions.map(p => p.id) : []),
    [rolePermissions]
  );

  const [selectedPermissions, setSelectedPermissions] = useState<string[]>(initialPermissions);
  const [searchQuery, setSearchQuery] = useState<string>("");
  const [activeCategory, setActiveCategory] = useState<string>("Dashboard");

  // Get category from menu name using mapping
  const getCategory = (perm: Permission): string => {
    // Handle standalone permissions without menu (e.g., ai-settings)
    if (!perm.menu) {
      if (perm.code.startsWith('ai-')) return 'AI';
      return 'Other';
    }
    const menuName = perm.menu?.name || "";
    return MENU_TO_CATEGORY[menuName] || "Other";
  };

  // Group permissions by Category -> Menu
  const groupedByCategory = useMemo(() => {
    const grouped: Record<string, Record<string, Permission[]>> = {};

    // Initialize categories
    Object.keys(CATEGORY_MAP).forEach(cat => {
      grouped[cat] = {};
    });
    grouped["Other"] = {};

    permissions.forEach(perm => {
      const category = getCategory(perm);
      // For standalone permissions, use a descriptive name from the permission
      let menuName = perm.menu?.name || "Other";
      
      // Special handling for AI Settings standalone permissions
      if (!perm.menu && perm.code.startsWith('ai-settings')) {
        menuName = "AI Settings";
      }
      
      if (!grouped[category]) grouped[category] = {};
      if (!grouped[category][menuName]) {
        grouped[category][menuName] = [];
      }
      
      grouped[category][menuName].push(perm);
    });

    return grouped;
  }, [permissions]);

  // Count permissions per category
  const categoryCounts = useMemo(() => {
    const counts: Record<string, { total: number; selected: number }> = {};
    Object.keys(CATEGORY_MAP).forEach(cat => {
      const categoryPermissions = Object.values(groupedByCategory[cat] || {}).flat();
      counts[cat] = {
        total: categoryPermissions.length,
        selected: categoryPermissions.filter(p => selectedPermissions.includes(p.id)).length,
      };
    });
    return counts;
  }, [groupedByCategory, selectedPermissions]);

  // Filter permissions based on search query
  const filteredPermissions = useMemo(() => {
    if (!searchQuery.trim()) {
      return groupedByCategory[activeCategory] || {};
    }

    const query = searchQuery.toLowerCase();
    const filtered: Record<string, Permission[]> = {};
    
    // When searching, look through ALL categories
    const allCategories = Object.keys(groupedByCategory);
    
    allCategories.forEach(category => {
      const categoryPerms = groupedByCategory[category] || {};
      
      Object.entries(categoryPerms).forEach(([menuName, perms]) => {
        const matchingPerms = perms.filter(perm => {
          const matchesName = perm.name.toLowerCase().includes(query);
          const matchesCode = perm.code.toLowerCase().includes(query);
          const matchesAction = perm.code.split('.')[1]?.toLowerCase().includes(query);
          const matchesMenu = menuName.toLowerCase().includes(query);
          return matchesName || matchesCode || matchesAction || matchesMenu;
        });

        if (matchingPerms.length > 0) {
          // If the menu already exists in filtered (from another category which shouldn't happen but good to be safe), merge it
          // Or just assign it. Since menu names might overlap across categories? 
          // Actually, our grouping logic assumes menu names map to categories. 
          // But just in case, let's keep it simple.
          if (!filtered[menuName]) {
             filtered[menuName] = [];
          }
          filtered[menuName] = [...filtered[menuName], ...matchingPerms];
        }
      });
    });

    return filtered;
  }, [groupedByCategory, activeCategory, searchQuery]);

  const togglePermission = (permissionId: string) => {
    const targetPermission = permissions.find((p) => p.id === permissionId);
    if (!targetPermission) return;

    const isCurrentlySelected = selectedPermissions.includes(permissionId);
    const [resource, action] = targetPermission.code.split(".");

    setSelectedPermissions((prev) => {
      let next = [...prev];

      if (isCurrentlySelected) {
        next = next.filter((id) => id !== permissionId);

        if (action === "view" || action === "read") {
          const relatedPermissions = permissions.filter(
            (p) => p.code.startsWith(`${resource}.`) && p.id !== permissionId
          );
          const relatedIds = new Set(relatedPermissions.map((p) => p.id));
          next = next.filter((id) => !relatedIds.has(id));
        }
      } else {
        next.push(permissionId);
        
        const viewPermission = permissions.find(
          (p) => p.code === `${resource}.view` || p.code === `${resource}.read`
        );
        if (viewPermission && !next.includes(viewPermission.id)) {
          next.push(viewPermission.id);
        }
      }

      return next;
    });
  };

  // Select all permissions in current category
  const handleSelectAll = () => {
    const categoryPermissions = Object.values(groupedByCategory[activeCategory] || {}).flat();
    const categoryIds = new Set(categoryPermissions.map(p => p.id));
    
    setSelectedPermissions(prev => {
      const withoutCategory = prev.filter(id => !categoryIds.has(id));
      return [...withoutCategory, ...categoryIds];
    });
  };

  // Unselect all permissions in current category
  const handleUnselectAll = () => {
    const categoryPermissions = Object.values(groupedByCategory[activeCategory] || {}).flat();
    const categoryIds = new Set(categoryPermissions.map(p => p.id));
    
    setSelectedPermissions(prev => prev.filter(id => !categoryIds.has(id)));
  };

  const handleSubmit = async () => {
    try {
      await assignPermissions.mutateAsync({
        roleId,
        permissionIds: selectedPermissions,
      });
      toast.success(t("save") || "Permissions saved successfully");
      onClose();
    } catch {
      // Error already handled in api-client interceptor
    }
  };

  const currentCount = categoryCounts[activeCategory] || { total: 0, selected: 0 };
  const sortedCategories = Object.keys(CATEGORY_MAP).sort((a, b) => 
    CATEGORY_MAP[a].order - CATEGORY_MAP[b].order
  );

  return (
    <div className="flex h-full">
      {/* Left Sidebar - Category Navigation */}
      <div className="w-56 border-r bg-muted/30 flex flex-col">
        <nav className="flex-1 py-2">
          {sortedCategories.map((category) => {
            const count = categoryCounts[category];
            const config = CATEGORY_MAP[category];
            const isActive = activeCategory === category;
            
            if (!count || count.total === 0) return null;
            
            return (
              <button
                key={category}
                type="button"
                onClick={() => {
                  setActiveCategory(category);
                  setSearchQuery("");
                }}
                className={cn(
                  "w-full flex items-center gap-3 px-4 py-3 text-left transition-colors",
                  "hover:bg-muted/80",
                  isActive && "bg-primary/10 border-l-2 border-primary text-primary"
                )}
              >
                <span className="shrink-0">{getMenuIcon(config.iconName)}</span>
                <div className="flex-1 min-w-0">
                  <div className="text-sm font-medium truncate">{config.name}</div>
                  {count && count.total > 0 && (
                    <div className="text-xs text-muted-foreground">
                      {count.selected}/{count.total} selected
                    </div>
                  )}
                </div>
              </button>
            );
          })}
        </nav>
      </div>

      {/* Right Content - Permissions List */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Search Bar */}
        <div className="px-6 py-3 border-b bg-muted/20">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              type="text"
              placeholder={t("searchPlaceholder") || "Search permissions, menus, or actions..."}
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="pl-9 h-9"
            />
          </div>
        </div>

        {/* Header with Select All / Unselect All */}
        <div className="flex items-center justify-between px-6 py-3 border-b bg-muted/20">
          <div className="text-sm text-muted-foreground">
            {currentCount.selected} of {currentCount.total} permissions selected
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={handleUnselectAll}
              className="text-xs h-7"
            >
              Unselect all
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={handleSelectAll}
              className="text-xs h-7"
            >
              Select all
            </Button>
          </div>
        </div>

        {/* Permissions List */}
        <div className="flex-1 overflow-y-auto">
          {Object.keys(filteredPermissions).length === 0 ? (
            <div className="flex items-center justify-center h-full text-muted-foreground text-sm">
              {searchQuery ? "No permissions found matching your search" : "No permissions in this category"}
            </div>
          ) : (
            <div className="divide-y">
              {Object.entries(filteredPermissions).map(([groupName, groupPermissions]) => (
                <div key={groupName}>
                  {/* Menu Group Header */}
                  <div className="px-6 py-2 bg-muted/30 border-b">
                    <span className="text-xs font-medium uppercase tracking-wider text-muted-foreground">
                      {groupName}
                    </span>
                  </div>
                  
                  {/* Permission Items */}
                  <div className="divide-y divide-border/50">
                    {groupPermissions.map((permission) => {
                      const isSelected = selectedPermissions.includes(permission.id);
                      const [, action] = permission.code.split(".");
                      
                      return (
                        <label
                          key={permission.id}
                          htmlFor={permission.id}
                          className={cn(
                            "flex items-center gap-4 px-6 py-3 cursor-pointer transition-colors",
                            "hover:bg-muted/50",
                            isSelected && "bg-primary/5"
                          )}
                        >
                          <Checkbox
                            id={permission.id}
                            checked={isSelected}
                            onCheckedChange={() => togglePermission(permission.id)}
                          />
                          <div className="flex-1 min-w-0">
                            <div className="text-sm font-medium">{permission.name}</div>
                          </div>
                          {action && (
                            <ActionBadge action={action} />
                          )}
                        </label>
                      );
                    })}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Footer Actions */}
        <div className="flex items-center justify-end gap-3 px-6 py-4 border-t bg-background">
          <Button variant="outline" onClick={onClose} disabled={assignPermissions.isPending}>
            {t("cancel")}
          </Button>
          <Button onClick={handleSubmit} disabled={assignPermissions.isPending}>
            {assignPermissions.isPending ? (t("saving") || "Saving...") : (t("save") || "Save")}
          </Button>
        </div>
      </div>
    </div>
  );
}

export function AssignPermissionsDrawer({ roleId, onClose }: AssignPermissionsDrawerProps) {
  const { data: permissionsData, isLoading: isLoadingPermissions, error: permissionsError } = usePermissions();
  const { data: roleData, isLoading: isLoadingRole, error: roleError } = useRole(roleId);
  const { data: rolePermissionsData, isLoading: isLoadingRolePermissions } = useRolePermissions(roleId);
  const t = useTranslations("userManagement.assignPermissions");

  const permissions = permissionsData?.data || [];
  const role = roleData;
  const rolePermissions = rolePermissionsData || [];

  // Debug logging
  if (typeof globalThis.window !== 'undefined') {
    console.log('AssignPermissionsDrawer Debug:', {
      roleId,
      isLoadingPermissions,
      isLoadingRole,
      hasPermissionsData: !!permissionsData,
      permissionsCount: permissions.length,
      hasRoleData: !!roleData,
      permissionsError,
      roleError
    });
  }

  const drawerContent = (() => {
    if (permissionsError || roleError) {
      return (
        <div className="flex items-center justify-center h-full p-6">
          <div className="text-center space-y-4">
            <div className="text-destructive font-medium text-lg">Error loading data</div>
            <div className="text-sm text-muted-foreground space-y-1">
              {permissionsError && <div>Failed to load permissions</div>}
              {roleError && (
                <div>
                  Failed to load role: {roleError.message}
                  {roleError.message?.includes('404') && (
                    <div className="mt-2 text-xs">
                      Role not found. The database may have been reseeded. Please refresh the page.
                    </div>
                  )}
                </div>
              )}
            </div>
            <Button onClick={onClose} variant="outline">
              Close
            </Button>
          </div>
        </div>
      );
    }

    if (isLoadingPermissions || isLoadingRole || isLoadingRolePermissions || !role) {
      return (
        <div className="flex h-full">
          {/* Sidebar Skeleton */}
          <div className="w-56 border-r bg-muted/30 p-4 space-y-3">
            <Skeleton className="h-12 w-full" />
            <Skeleton className="h-12 w-full" />
            <Skeleton className="h-12 w-full" />
            <Skeleton className="h-12 w-full" />
            <Skeleton className="h-12 w-full" />
          </div>
          {/* Content Skeleton */}
          <div className="flex-1 p-6 space-y-4">
            <Skeleton className="h-14 w-full" />
            <Skeleton className="h-14 w-full" />
            <Skeleton className="h-14 w-full" />
            <Skeleton className="h-14 w-full" />
            <Skeleton className="h-14 w-full" />
            <Skeleton className="h-14 w-full" />
          </div>
        </div>
      );
    }

    return (
      <PermissionsSelector
        permissions={permissions}
        rolePermissions={rolePermissions}
        roleId={roleId}
        onClose={onClose}
      />
    );
  })();

  return (
    <Drawer
      open={!!roleId}
      onOpenChange={(open) => !open && onClose()}
      title={t("title", { roleName: role?.name ?? "" })}
      description={t("description")}
      side="right"
      defaultWidth={800}
      minWidth={600}
      maxWidth={1200}
      resizable
    >
      {drawerContent}
    </Drawer>
  );
}

// Keep backward compatibility alias
export { AssignPermissionsDrawer as AssignPermissionsDialog };
