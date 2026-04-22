"use client";

import { Package, Tag } from "lucide-react";
import { useTranslations } from "next-intl";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { ProductList } from "./product-list";
import { CategoryList } from "./category-list";
import { useHasPermission } from "@/features/master-data/user-management/hooks/useHasPermission";

export function ProductManagement() {
  const hasProductsPermission = useHasPermission("products.view");
  const hasCategoriesPermission = useHasPermission("products.category-view");
  const t = useTranslations("productManagement.tabs");

  // Determine default tab - use first available tab
  const defaultTab = hasProductsPermission ? "products" : "categories";

  return (
    <Tabs defaultValue={defaultTab} className="w-full">
      <TabsList className="w-full flex-wrap h-auto">
        {hasProductsPermission && (
          <TabsTrigger value="products" className="gap-2">
            <Package className="h-4 w-4" />
            {t("products")}
          </TabsTrigger>
        )}
        {hasCategoriesPermission && (
          <TabsTrigger value="categories" className="gap-2">
            <Tag className="h-4 w-4" />
            {t("categories")}
          </TabsTrigger>
        )}
      </TabsList>

      {hasProductsPermission && (
        <TabsContent value="products" className="mt-6">
          <ProductList />
        </TabsContent>
      )}

      {hasCategoriesPermission && (
        <TabsContent value="categories" className="mt-6">
          <CategoryList />
        </TabsContent>
      )}
    </Tabs>
  );
}
