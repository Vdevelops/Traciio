"use client";

import React from "react";
import { Users, Shield, Key } from "lucide-react";
import { useTranslations } from "next-intl";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { UserList } from "./user-list";
import { RoleList } from "./role-list";
import { PermissionList } from "./permission-list";
import { useHasPermission } from "../hooks/useHasPermission";

export function UserManagement() {
  const hasRolesPermission = useHasPermission("users.roles");
  const hasPermissionsPermission = useHasPermission("users.permissions");
  const t = useTranslations("userManagement.tabs");

  const router = useRouter();
  const searchParams = useSearchParams();
  const pathname = usePathname();

  // Build available tabs based on permissions
  const availableTabs = React.useMemo(() => {
    const tabs = ["users"];
    if (hasRolesPermission) tabs.push("roles");
    if (hasPermissionsPermission) tabs.push("permissions");
    return tabs;
  }, [hasRolesPermission, hasPermissionsPermission]);

  const sanitizeTab = (tab: string | null) => {
    if (!tab) return "users";
    return availableTabs.includes(tab) ? tab : "users";
  };

  // Initialize from ?tab=... if present and allowed
  const initialTab = sanitizeTab(searchParams.get("tab"));
  const [value, setValue] = React.useState<string>(initialTab);

  // Keep in sync when search params or permissions change (e.g., user navigates via link)
  React.useEffect(() => {
    const tab = sanitizeTab(searchParams.get("tab"));
    setValue(tab);
  }, [searchParams, hasRolesPermission, hasPermissionsPermission]);

  const handleChange = (val: string) => {
    if (val === value) return;
    setValue(val);
    const params = new URLSearchParams(searchParams.toString());
    params.set("tab", val);
    router.push(`${pathname}?${params.toString()}`);
  };

  return (
    <Tabs value={value} onValueChange={handleChange} className="w-full">
        <TabsList>
          <TabsTrigger value="users" className="gap-2">
            <Users className="h-4 w-4" />
            {t("users")}
          </TabsTrigger>
          {hasRolesPermission && (
          <TabsTrigger value="roles" className="gap-2">
            <Shield className="h-4 w-4" />
            {t("roles")}
          </TabsTrigger>
          )}
          {hasPermissionsPermission && (
          <TabsTrigger value="permissions" className="gap-2">
            <Key className="h-4 w-4" />
            {t("permissions")}
          </TabsTrigger>
          )}
        </TabsList>

        <TabsContent value="users" className="mt-6">
          <UserList />
        </TabsContent>

        {hasRolesPermission && (
        <TabsContent value="roles" className="mt-6">
          <RoleList />
        </TabsContent>
        )}

        {hasPermissionsPermission && (
        <TabsContent value="permissions" className="mt-6">
          <PermissionList />
        </TabsContent>
        )}
      </Tabs>
  );
}

