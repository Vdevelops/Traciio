"use client";

import { useTranslations } from "next-intl";
import { cn } from "@/lib/utils";
import { useProductCategories } from "../hooks/useProducts";
import { Skeleton } from "@/components/ui/skeleton";

interface CategorySidebarProps {
  selectedCategoryId: string | null;
  onCategorySelect: (categoryId: string | null) => void;
}

export function CategorySidebar({
  selectedCategoryId,
  onCategorySelect,
}: CategorySidebarProps) {
  const t = useTranslations("productManagement.categorySidebar");
  const { data: categoriesData, isLoading } = useProductCategories({ status: "active" });
  const categories = categoriesData?.data || [];

  if (isLoading) {
    return (
      <div className="rounded-lg border bg-muted/10 p-4 space-y-2">
        <Skeleton className="h-10 w-full" />
        <Skeleton className="h-8 w-full" />
        <Skeleton className="h-8 w-full" />
        <Skeleton className="h-8 w-full" />
        <Skeleton className="h-8 w-full" />
      </div>
    );
  }

  return (
    <div className="rounded-lg border bg-muted/10 p-4">
      <h3 className="font-semibold text-sm text-muted-foreground uppercase tracking-wider mb-4">
        {t("title")}
      </h3>
      <div className="flex flex-wrap gap-2">
        {/* All Products */}
        <button
          type="button"
          onClick={() => onCategorySelect(null)}
          aria-pressed={selectedCategoryId === null}
          className={cn(
            "flex items-center px-3 py-2 rounded-md text-sm transition-colors",
            selectedCategoryId === null
              ? "bg-primary text-primary-foreground font-medium"
              : "hover:bg-muted text-foreground"
          )}
        >
          <span className="whitespace-nowrap">
            {t("allProducts")}
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
              "flex items-center px-3 py-2 rounded-md text-sm transition-colors",
              selectedCategoryId === category.id
                ? "bg-primary text-primary-foreground font-medium"
                : "hover:bg-muted text-foreground"
            )}
          >
            <span className="whitespace-nowrap">
              {category.name}
            </span>
          </button>
        ))}
      </div>
    </div>
  );
}
