"use client";

import { useState } from "react";
import { toast } from "sonner";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Checkbox } from "@/components/ui/checkbox";
import { useRoleMobilePermissions, useUpdateRoleMobilePermissions, useRole } from "../hooks/useRoles";
import type { MobileMenuPermission } from "../types";

interface MobilePermissionsDialogProps {
  readonly roleId: string;
  readonly onClose: () => void;
}

// Mobile menu definitions
const MOBILE_MENUS = [
  { code: "dashboard", label: "Dashboard" },
  { code: "task", label: "Task" },
  { code: "accounts", label: "Accounts" },
  { code: "contacts", label: "Contacts" },
  { code: "visit_reports", label: "Visit Reports" },
] as const;

const ACTIONS = [
  { code: "VIEW", label: "View" },
  { code: "CREATE", label: "Create" },
  { code: "EDIT", label: "Edit" },
  { code: "DELETE", label: "Delete" },
] as const;

// Inner component that only renders when data is loaded
function MobilePermissionsSelector({
  roleId,
  initialPermissions,
  onClose,
}: {
  readonly roleId: string;
  readonly initialPermissions: MobileMenuPermission[];
  readonly onClose: () => void;
}) {
  const updatePermissions = useUpdateRoleMobilePermissions();
  const t = useTranslations("userManagement.mobilePermissions");

  // Initialize state from initial permissions
  const [permissions, setPermissions] = useState<Record<string, string[]>>(() => {
    const state: Record<string, string[]> = {};
    initialPermissions.forEach((menu) => {
      state[menu.menu] = [...menu.actions];
    });
    // Initialize all menus if not present
    MOBILE_MENUS.forEach((menu) => {
      if (!state[menu.code]) {
        state[menu.code] = [];
      }
    });
    return state;
  });

  const toggleAction = (menuCode: string, actionCode: string) => {
    setPermissions((prev) => {
      const menuActions = prev[menuCode] || [];
      
      // Dashboard hanya bisa memiliki VIEW action
      if (menuCode === "dashboard") {
        if (actionCode === "VIEW") {
          // Toggle VIEW untuk dashboard
          const newActions = menuActions.includes("VIEW")
            ? []
            : ["VIEW"];
          return {
            ...prev,
            [menuCode]: newActions,
          };
        }
        // Untuk dashboard, ignore action selain VIEW
        return prev;
      }
      
      // Menu lain bisa memiliki semua actions
      const newActions = menuActions.includes(actionCode)
        ? menuActions.filter((a) => a !== actionCode)
        : [...menuActions, actionCode];
      return {
        ...prev,
        [menuCode]: newActions,
      };
    });
  };

  const handleSubmit = async () => {
    try {
      const menus: MobileMenuPermission[] = MOBILE_MENUS.map((menu) => {
        let actions = permissions[menu.code] || [];
        
        // Dashboard hanya bisa memiliki VIEW action
        if (menu.code === "dashboard") {
          actions = actions.filter((a) => a === "VIEW");
        }
        
        return {
          menu: menu.code,
          actions,
        };
      });

      await updatePermissions.mutateAsync({
        roleId,
        permissions: { menus },
      });
      toast.success(t("saved"));
      onClose();
    } catch {
      // Error already handled in api-client interceptor
    }
  };

  return (
    <>
      <div className="space-y-4">
        {MOBILE_MENUS.map((menu) => {
          const menuActions = permissions[menu.code] || [];
          return (
            <Card key={menu.code}>
              <CardHeader>
                <CardTitle className="text-sm">{menu.label}</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
                  {ACTIONS.map((action) => {
                    // Dashboard hanya menampilkan VIEW action
                    const isDashboard = menu.code === "dashboard";
                    const isViewAction = action.code === "VIEW";
                    const shouldShow = !isDashboard || isViewAction;
                    
                    if (!shouldShow) {
                      return null;
                    }
                    
                    return (
                      <div
                        key={`${menu.code}-${action.code}`}
                        className="flex items-center space-x-2 p-2 rounded-md hover:bg-muted"
                      >
                        <Checkbox
                          id={`${menu.code}-${action.code}`}
                          checked={menuActions.includes(action.code)}
                          onCheckedChange={() => toggleAction(menu.code, action.code)}
                        />
                        <label
                          htmlFor={`${menu.code}-${action.code}`}
                          className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer flex-1"
                        >
                          {action.label}
                        </label>
                      </div>
                    );
                  })}
                </div>
              </CardContent>
            </Card>
          );
        })}
      </div>

      <div className="flex justify-end gap-2 pt-4 border-t">
        <Button variant="outline" onClick={onClose} disabled={updatePermissions.isPending}>
          {t("cancel")}
        </Button>
        <Button onClick={handleSubmit} disabled={updatePermissions.isPending}>
          {updatePermissions.isPending ? t("saving") : t("save")}
        </Button>
      </div>
    </>
  );
}

export function MobilePermissionsDialog({ roleId, onClose }: MobilePermissionsDialogProps) {
  const { data: permissionsData, isLoading: isLoadingPermissions } = useRoleMobilePermissions(roleId);
  const { data: roleData, isLoading: isLoadingRole } = useRole(roleId);
  const t = useTranslations("userManagement.mobilePermissions");

  const permissions = permissionsData?.menus || [];
  const role = roleData;

  return (
    <Dialog open={!!roleId} onOpenChange={(open) => !open && onClose()} key={roleId}>
      <DialogContent className="max-w-4xl max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>
            {t("title")} - {role?.name ?? ""}
          </DialogTitle>
          <DialogDescription>
            {t("description")}
          </DialogDescription>
        </DialogHeader>

        {isLoadingPermissions || isLoadingRole || !role ? (
          <div className="space-y-4">
            {MOBILE_MENUS.map((menu) => (
              <Skeleton key={`mobile-permissions-skeleton-${menu.code}`} className="h-32 w-full" />
            ))}
          </div>
        ) : (
          <MobilePermissionsSelector
            roleId={roleId}
            initialPermissions={permissions}
            onClose={onClose}
          />
        )}
      </DialogContent>
    </Dialog>
  );
}
