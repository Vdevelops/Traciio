"use client";

import { Tag, UserCircle } from "lucide-react";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { CategoryList } from "./category-list";
import { ContactRoleList } from "./contact-role-list";
import { useHasPermission } from "@/features/master-data/user-management/hooks/useHasPermission";
import { useTranslations } from "next-intl";

export function AccountManagement() {
  const hasCategoryPermission = useHasPermission("accounts.category");
  const hasRolePermission = useHasPermission("accounts.role");
  const tTabs = useTranslations("accountManagement.tabs");

  const defaultTab = hasCategoryPermission ? "categories" : "contact-roles";

  return (
    <div className="space-y-4 sm:space-y-6">
      <Tabs defaultValue={defaultTab} className="w-full">
        <TabsList className="w-full sm:w-auto overflow-x-auto">
          {hasCategoryPermission && (
            <TabsTrigger value="categories" className="gap-1.5 sm:gap-2 text-xs sm:text-sm">
              <Tag className="h-3.5 w-3.5 sm:h-4 sm:w-4" />
              <span className="whitespace-nowrap">{tTabs("categories")}</span>
            </TabsTrigger>
          )}
          {hasRolePermission && (
            <TabsTrigger value="contact-roles" className="gap-1.5 sm:gap-2 text-xs sm:text-sm">
              <UserCircle className="h-3.5 w-3.5 sm:h-4 sm:w-4" />
              <span className="whitespace-nowrap">{tTabs("contactRoles")}</span>
            </TabsTrigger>
          )}
        </TabsList>

        {hasCategoryPermission && (
          <TabsContent value="categories" className="mt-4 sm:mt-6">
            <CategoryList />
          </TabsContent>
        )}

        {hasRolePermission && (
          <TabsContent value="contact-roles" className="mt-4 sm:mt-6">
            <ContactRoleList />
          </TabsContent>
        )}
      </Tabs>
    </div>
  );
}


