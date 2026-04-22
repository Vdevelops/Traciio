"use client";

import { useTranslations } from "next-intl";
import { cn } from "@/lib/utils";
import { useProductCategories } from "../hooks/useProducts";
import { Skeleton } from "@/components/ui/skeleton";

interface CategorySidebarProps {
  selectedCategoryId: string | null;
  onCategorySelect: (categoryId: string | null) => void;
  totalProducts?: number;
}

export function CategorySidebar({
  selectedCategoryId,
  onCategorySelect,
  totalProducts = 0,
}: CategorySidebarProps) {
  const t = useTranslations("productManagement.categorySidebar");
  const { data: categoriesData, isLoading } = useProductCategories({ status: "active" });
  const categories = categoriesData?.data || [];

  if (isLoading) {
    return (
      <div className="w-full lg:w-64 lg:shrink-0 border-b lg:border-b-0 lg:border-r bg-muted/10 p-4 space-y-2">
        <Skeleton className="h-10 w-full" />
        <Skeleton className="h-8 w-full" />
        <Skeleton className="h-8 w-full" />
        <Skeleton className="h-8 w-full" />
        <Skeleton className="h-8 w-full" />
      </div>
    );
  }

  return (
    <div className="w-full lg:w-64 lg:shrink-0 border-b lg:border-b-0 lg:border-r bg-muted/10 p-4">
      <h3 className="font-semibold text-sm text-muted-foreground uppercase tracking-wider mb-4">
        {t("title")}
      </h3>
      <div className="flex gap-2 overflow-x-auto pb-1 lg:flex-col lg:overflow-visible lg:pb-0">
        {/* All Products */}
        <button
          type="button"
          onClick={() => onCategorySelect(null)}
          aria-pressed={selectedCategoryId === null}
          className={cn(
            "flex items-center justify-between px-3 py-2 rounded-md text-sm transition-colors shrink-0 lg:w-full",
            selectedCategoryId === null
              ? "bg-primary text-primary-foreground font-medium"
              : "hover:bg-muted text-foreground"
          )}
        >
          <span className="whitespace-nowrap lg:whitespace-normal">
            {t("allProducts")}
          </span>
          <span
            className={cn(
              "text-xs font-medium px-2 py-0.5 rounded-full shrink-0",
              selectedCategoryId === null
                ? "bg-primary-foreground/20 text-primary-foreground"
                : "bg-muted-foreground/10 text-muted-foreground"
            )}
          >
            {totalProducts}
          </span>
        </button>

        {/* Category List */}
        {categories.map((category) => (
          <button
            key={category.id}
            type="button"
            onClick={() => onCategorySelect(category.id)}
            aria-pressed={selectedCategoryId === category.id}
            className={cn(
              "flex items-center justify-between px-3 py-2 rounded-md text-sm transition-colors shrink-0 lg:w-full",
              selectedCategoryId === category.id
                ? "bg-primary text-primary-foreground font-medium"
                : "hover:bg-muted text-foreground"
            )}
          >
            <span className="whitespace-nowrap lg:truncate">
              {category.name}
            </span>
            <span
              className={cn(
                "text-xs font-medium px-2 py-0.5 rounded-full shrink-0 ml-2",
                selectedCategoryId === category.id
                  ? "bg-primary-foreground/20 text-primary-foreground"
                  : "bg-muted-foreground/10 text-muted-foreground"
              )}
            >
              {category.product_count ?? 0}
            </span>
          </button>
        ))}
      </div>
    </div>
  );
}
