"use client";

import { Edit, Trash2, Package, DollarSign, Hash, Tag } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import { useProduct, useDeleteProduct, useUpdateProduct } from "../hooks/useProducts";
import { toast } from "sonner";
import { useState, Suspense, useMemo } from "react";
import { ProductForm } from "./product-form";
import { useTranslations } from "next-intl";
import { useHasPermission } from "@/features/master-data/user-management/hooks/useHasPermission";
import { useProductMonthlySales } from "@/features/product-analytics/hooks/useProductAnalytics";
import type { DateRange } from "react-day-picker";
import type { MonthlySalesChartProps } from "@/features/product-analytics/components/MonthlySalesChart";
import dynamic from "next/dynamic";

const MonthlySalesChart = dynamic<MonthlySalesChartProps>(
  () => import("@/features/product-analytics/components/MonthlySalesChart").then(m => ({ default: m.MonthlySalesChart })),
  { ssr: false }
);

interface ProductDetailModalProps {
  readonly productId: string | null;
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly onProductUpdated?: () => void;
}

export function ProductDetailModal({ productId, open, onOpenChange, onProductUpdated }: ProductDetailModalProps) {
  const { data, isLoading, error } = useProduct(productId || "");
  const deleteProduct = useDeleteProduct();
  const updateProduct = useUpdateProduct();
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const t = useTranslations("productManagement.detailModal");
  const hasEditPermission = useHasPermission("products.edit");
  const hasDeletePermission = useHasPermission("products.delete");

  // Date Filter Mode
  const [filterMode, setFilterMode] = useState<"year" | "range">("year");
  const [selectedYear, setSelectedYear] = useState<number>(() => new Date().getFullYear());

  // Date Range State for Analytics
  const [dateRange, setDateRange] = useState<DateRange | undefined>(() => {
    const today = new Date();
    const startOfYear = new Date(today.getFullYear(), 0, 1);
    return { from: startOfYear, to: today };
  });

  const startDateStr = useMemo(() => {
    if (filterMode === "year") {
      return `${selectedYear}-01-01`;
    }
    if (!dateRange?.from) return undefined;
    const d = dateRange.from;
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
  }, [filterMode, selectedYear, dateRange]);

  const endDateStr = useMemo(() => {
    if (filterMode === "year") {
      return `${selectedYear}-12-31`;
    }
    if (!dateRange?.to) return undefined;
    const d = dateRange.to;
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
  }, [filterMode, selectedYear, dateRange]);

  const {
    monthlySales,
    isLoading: isLoadingMonthlySales,
  } = useProductMonthlySales(productId || "", startDateStr, endDateStr);

  const product = data?.data;

  const handleDeleteConfirm = async () => {
    if (!product || !productId) return;
    try {
      await deleteProduct.mutateAsync(productId);
      toast.success(t("toastDeleted"));
      onOpenChange(false);
      onProductUpdated?.();
    } catch {
      // Error already handled in api-client interceptor
    }
  };

  const formatCurrency = (amount: number) => {
    const rupiah = amount / 100;
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
    }).format(rupiah);
  };

  return (
    <>
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent className="sm:max-w-[900px]">
          <DialogHeader>
            <DialogTitle>{t("title")}</DialogTitle>
          </DialogHeader>

          {isLoading && (
            <div className="space-y-6">
              <div className="flex items-center gap-4">
                <Skeleton className="h-20 w-20 rounded-full" />
                <div className="space-y-2">
                  <Skeleton className="h-6 w-48" />
                  <Skeleton className="h-4 w-32" />
                </div>
              </div>
              <Card>
                <CardHeader>
                  <Skeleton className="h-6 w-32" />
                  <Skeleton className="h-4 w-64 mt-2" />
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <Skeleton className="h-4 w-full" />
                    <Skeleton className="h-4 w-3/4" />
                    <Skeleton className="h-4 w-1/2" />
                  </div>
                </CardContent>
              </Card>
            </div>
          )}

          {error && (
            <div className="text-center text-muted-foreground py-8">
              {t("loadError")}
            </div>
          )}

          {!isLoading && !error && product && (
            <div className="space-y-6">
              {/* Product Header */}
              <div className="flex items-center gap-4 pb-4 border-b">
                <div className="h-20 w-20 rounded-lg bg-primary/10 flex items-center justify-center">
                  <Package className="h-10 w-10 text-primary" />
                </div>
                <div className="flex-1">
                  <h2 className="text-2xl font-medium tracking-tight">{product.name}</h2>
                  <p className="text-sm text-muted-foreground mt-1 font-mono">{product.sku}</p>
                </div>
                <div className="flex items-center gap-2">
                  {hasEditPermission && (
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setIsEditDialogOpen(true)}
                      className="shrink-0"
                    >
                      <Edit className="h-4 w-4" />
                    </Button>
                  )}
                  {hasDeletePermission && (
                    <Button
                      variant="destructive"
                      size="sm"
                      onClick={() => setIsDeleteDialogOpen(true)}
                      disabled={deleteProduct.isPending}
                      className="shrink-0"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  )}
                </div>
              </div>

              {/* Stats Cards */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {/* Price Card */}
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">
                      {t("stats.price")}
                    </CardTitle>
                    <DollarSign className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-medium">
                      {product.price_formatted || formatCurrency(product.price ?? 0)}
                    </div>
                    <p className="text-xs text-muted-foreground">
                      {t("stats.sellingPrice")}
                    </p>
                  </CardContent>
                </Card>

                {/* Cost Card */}
                <Card>
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">
                      {t("stats.cost")}
                    </CardTitle>
                    <DollarSign className="h-4 w-4 text-muted-foreground" />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-medium">
                      {formatCurrency(product.cost ?? 0)}
                    </div>
                    <p className="text-xs text-muted-foreground">
                      {t("stats.costPrice")}
                    </p>
                  </CardContent>
                </Card>
              </div>

              {/* Product Info Card */}
              <Card>
                <CardHeader>
                  <CardTitle>{t("productInfo.title")}</CardTitle>
                  <CardDescription>
                    {t("productInfo.sectionDescription")}
                  </CardDescription>
                </CardHeader>
                <CardContent className="space-y-6">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div className="space-y-2">
                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <Package className="h-4 w-4" />
                        <span>{t("productInfo.name")}</span>
                      </div>
                      <div className="text-base font-medium">{product.name}</div>
                    </div>

                    <div className="space-y-2">
                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <Hash className="h-4 w-4" />
                        <span>{t("productInfo.sku")}</span>
                      </div>
                      <div className="text-base font-medium font-mono">{product.sku}</div>
                    </div>

                    {product.barcode && (
                      <div className="space-y-2">
                        <div className="flex items-center gap-2 text-sm text-muted-foreground">
                          <Hash className="h-4 w-4" />
                          <span>{t("productInfo.barcode")}</span>
                        </div>
                        <div className="text-base font-medium font-mono">{product.barcode}</div>
                      </div>
                    )}

                    <div className="space-y-2">
                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <Tag className="h-4 w-4" />
                        <span>{t("productInfo.category")}</span>
                      </div>
                      <div>
                        <Badge variant="outline" className="font-normal">
                          {product.category?.name || "N/A"}
                        </Badge>
                      </div>
                    </div>
                  </div>

                  {product.description && (
                    <div className="space-y-2 pt-4 border-t">
                      <div className="text-sm text-muted-foreground">
                        {t("productInfo.description")}
                      </div>
                      <div className="text-base">{product.description}</div>
                    </div>
                  )}
                </CardContent>
              </Card>

              {/* Monthly Sales Chart */}
              {productId && (
                <Suspense fallback={
                  <Card>
                    <CardHeader>
                      <Skeleton className="h-6 w-48" />
                    </CardHeader>
                    <CardContent>
                      <Skeleton className="h-[400px] w-full" />
                    </CardContent>
                  </Card>
                }>
                  <MonthlySalesChart
                    data={monthlySales}
                    isLoading={isLoadingMonthlySales}
                    filterMode={filterMode}
                    onFilterModeChange={setFilterMode}
                    selectedYear={selectedYear}
                    onYearChange={setSelectedYear}
                    dateRange={dateRange}
                    onDateChange={setDateRange}
                  />
                </Suspense>
              )}
            </div>
          )}
        </DialogContent>
      </Dialog>

      {/* Edit Dialog */}
      {hasEditPermission && isEditDialogOpen && product && (
        <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
          <DialogContent className="sm:max-w-[500px]">
            <DialogHeader>
              <DialogTitle>{t("editDialogTitle")}</DialogTitle>
            </DialogHeader>
            <ProductForm
              product={product}
              onSubmit={async (formData) => {
                try {
                  await updateProduct.mutateAsync({ id: productId!, data: formData });
                  setIsEditDialogOpen(false);
                  toast.success(t("toastUpdated"));
                  onProductUpdated?.();
                } catch {
                  // Error already handled in api-client interceptor
                }
              }}
              onCancel={() => setIsEditDialogOpen(false)}
              isLoading={updateProduct.isPending}
            />
          </DialogContent>
        </Dialog>
      )}

      {/* Delete Dialog */}
      {hasDeletePermission && (
        <DeleteDialog
          open={isDeleteDialogOpen}
          onOpenChange={setIsDeleteDialogOpen}
          onConfirm={handleDeleteConfirm}
          title={t("deleteDialogTitle")}
          description={
            product
              ? t("deleteDialogDescriptionWithName", { name: product.name })
              : t("deleteDialogDescription")
          }
          itemName={t("deleteDialogItemName")}
          isLoading={deleteProduct.isPending}
        />
      )}
    </>
  );
}
