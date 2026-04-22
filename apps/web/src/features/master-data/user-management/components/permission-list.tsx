"use client";

import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
import { usePermissions } from "../hooks/usePermissions";
import type { Permission } from "../types";

export function PermissionList() {
  const { data, isLoading } = usePermissions();

  const permissions = data?.data || [];

  // Group permissions by menu
  const groupedPermissions = permissions.reduce((acc, perm) => {
    const key = perm.menu?.name || "Other";
    if (!acc[key]) {
      acc[key] = [];
    }
    acc[key].push(perm);
    return acc;
  }, {} as Record<string, Permission[]>);

  return (
    <div className="space-y-6">
      {isLoading ? (
        <div className="space-y-4">
          {Array.from({ length: 3 }).map((_, i) => (
            <Skeleton key={`skeleton-${i}`} className="h-32 w-full" />
          ))}
        </div>
      ) : (
        <div className="space-y-6">
          {Object.entries(groupedPermissions).map(([menuName, menuPermissions]) => (
            <div key={menuName} className="space-y-3">
              <h3 className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
                {menuName}
              </h3>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
                {menuPermissions.map((permission) => (
                  <div
                    key={permission.id}
                    className="p-3 border rounded-lg hover:bg-muted/30 transition-colors"
                  >
                    <div className="space-y-2">
                      <p className="font-medium text-sm">{permission.name}</p>
                      <div className="flex items-center gap-2 flex-wrap">
                        <code className="text-xs bg-muted px-1.5 py-0.5 rounded">
                          {permission.code}
                        </code>
                        <Badge variant="outline" className="text-xs">
                          {permission.action}
                        </Badge>
                      </div>
                      {permission.description && (
                        <p className="text-xs text-muted-foreground">
                          {permission.description}
                        </p>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
