"use client";

import { useState } from "react";
import { Edit, Trash2, Plus, Search, Eye, Package } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { DataTable, type Column } from "@/components/ui/data-table";
import { StatusSwitch } from "@/components/ui/status-switch";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useProductList } from "../hooks/useProductList";
import { ProductForm } from "./product-form";
import { ProductDetailModal } from "./product-detail-modal";
import { CategorySidebar } from "./category-sidebar";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import type { Product } from "../types";
import { useTranslations } from "next-intl";
import { useHasPermission } from "@/features/master-data/user-management/hooks/useHasPermission";
import type { CreateProductFormData, UpdateProductFormData } from "../schemas/product.schema";

export function ProductList() {
  const hasViewPermission = useHasPermission("products.view");
  const hasCreatePermission = useHasPermission("products.create");
  const hasEditPermission = useHasPermission("products.edit");
  const hasDeletePermission = useHasPermission("products.delete");

  const {
    page,
    setPage,
    setPerPage,
    search,
    setSearch,
    status,
    setStatus,
    categoryId,
    setCategoryId,
    isCreateDialogOpen,
    setIsCreateDialogOpen,
    editingProduct,
    setEditingProduct,
    products,
    pagination,
    categories,
    editingProductData,
    isLoading,
    isError,
    handleCreate,
    handleUpdate,
    handleDeleteClick,
    handleDeleteConfirm,
    deletingProductId,
    setDeletingProductId,
    deleteProduct,
    createProduct,
    updateProduct,
  } = useProductList();

  const [viewingProductId, setViewingProductId] = useState<string | null>(null);
  const [isDetailModalOpen, setIsDetailModalOpen] = useState(false);
  const t = useTranslations("productManagement.list");

  if (!hasViewPermission) {
    return (
      <div className="text-center text-muted-foreground py-8">
        {t("permissionDenied")}
      </div>
    );
  }

  const handleViewProduct = (productId: string) => {
    setViewingProductId(productId);
    setIsDetailModalOpen(true);
  };

  const formatCurrency = (amount: number) => {
    const rupiah = amount / 100;
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
    }).format(rupiah);
  };

  const columns: Column<Product>[] = [
    {
      id: "icon",
      header: "",
      accessor: () => (
        <div className="h-10 w-10 rounded-md bg-primary/10 flex items-center justify-center">
          <Package className="h-5 w-5 text-primary" />
        </div>
      ),
      className: "w-[60px]",
    },
    {
      id: "name",
      header: t("name"),
      accessor: (row) => (
        hasViewPermission ? (
          <button
            onClick={() => handleViewProduct(row.id)}
            className="flex items-center gap-3 font-medium text-primary hover:underline text-left"
          >
            <span>{row.name}</span>
          </button>
        ) : (
          <span className="font-medium">{row.name}</span>
        )
      ),
      className: "w-[200px]",
    },
    {
      id: "sku",
      header: t("sku"),
      accessor: (row) => (
        <span className="text-muted-foreground font-mono text-sm">{row.sku}</span>
      ),
    },
    {
      id: "category",
      header: t("category"),
      accessor: (row) => (
        <Badge variant="outline" className="font-normal">
          {row.category?.name || "-"}
        </Badge>
      ),
    },
    {
      id: "price",
      header: t("price"),
      accessor: (row) => (
        <span className="font-medium">{row.price_formatted || formatCurrency(row.price)}</span>
      ),
    },
    {
      id: "status",
      header: t("status"),
      accessor: (row) => (
        <StatusSwitch
          checked={row.status === "active"}
          onCheckedChange={(checked) => {
            updateProduct.mutate({
              id: row.id,
              data: { status: checked ? "active" : "inactive" },
            });
          }}
          disabled={!hasEditPermission}
        />
      ),
      className: "w-[100px]",
    },
    {
      id: "actions",
      header: t("actions"),
      accessor: (row) => (
        <div className="flex items-center justify-end gap-1">
          {hasViewPermission && (
            <Button
              variant="ghost"
              size="icon-sm"
              className="h-8 w-8"
              title={t("viewDetails")}
              onClick={() => handleViewProduct(row.id)}
            >
              <Eye className="h-3.5 w-3.5" />
            </Button>
          )}
          {hasEditPermission && (
            <Button
              variant="ghost"
              size="icon-sm"
              onClick={() => setEditingProduct(row.id)}
              className="h-8 w-8"
              title={t("edit")}
            >
              <Edit className="h-3.5 w-3.5" />
            </Button>
          )}
          {hasDeletePermission && (
            <Button
              variant="ghost"
              size="icon-sm"
              onClick={() => handleDeleteClick(row.id)}
              className="h-8 w-8 text-destructive hover:text-destructive"
              title={t("delete")}
            >
              <Trash2 className="h-3.5 w-3.5" />
            </Button>
          )}
        </div>
      ),
      className: "w-[140px] text-right",
    },
  ];

  return (
    <div className="space-y-4">
      {/* Category Filter */}
      <CategorySidebar
        selectedCategoryId={categoryId || null}
        onCategorySelect={(id) => {
          setCategoryId(id || "");
          setPage(1);
        }}
      />

      {/* Main Content */}
      <div className="space-y-4">
        {/* Header with Actions */}
        <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div className="flex min-w-0 flex-1 flex-col gap-3 sm:flex-row sm:items-center">
            <div className="relative w-full sm:max-w-sm">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                aria-label={t("searchPlaceholder")}
                placeholder={t("searchPlaceholder")}
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="pl-10 h-9"
              />
            </div>
            <Select
              value={status || "all"}
              onValueChange={(value) => setStatus(value === "all" ? "" : value)}
            >
              <SelectTrigger className="w-full sm:w-[140px] h-9" aria-label={t("allStatus")}>
                <SelectValue placeholder={t("allStatus")} />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">{t("allStatus")}</SelectItem>
                <SelectItem value="active">{t("statusActive")}</SelectItem>
                <SelectItem value="inactive">{t("statusInactive")}</SelectItem>
              </SelectContent>
            </Select>
          </div>
          {hasCreatePermission && (
            <Button
              onClick={() => setIsCreateDialogOpen(true)}
              size="sm"
              className="w-full sm:w-auto"
            >
              <Plus className="h-4 w-4 mr-2" />
              {t("addProduct")}
            </Button>
          )}
        </div>

        {/* Table */}
        {isError ? (
          <div className="rounded-lg border border-destructive/20 bg-destructive/5 px-4 py-6 text-sm text-destructive">
            {t("loadError")}
          </div>
        ) : (
          <div className="w-full">
            <DataTable
              columns={columns}
              data={products}
              isLoading={isLoading}
              emptyMessage={t("empty")}
              pagination={
                pagination
                  ? {
                      page: pagination.page,
                      per_page: pagination.per_page,
                      total: pagination.total,
                      total_pages: pagination.total_pages,
                      has_next: pagination.has_next,
                      has_prev: pagination.has_prev,
                    }
                  : undefined
              }
              onPageChange={setPage}
              onPerPageChange={setPerPage}
              itemName="product"
              perPageOptions={[10, 20, 50, 100]}
              onResetFilters={() => {
                setSearch("");
                setStatus("");
                setCategoryId("");
                setPage(1);
              }}
            />
          </div>
        )}

      {/* Create Dialog */}
      {hasCreatePermission && (
        <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
          <DialogContent className="sm:max-w-[500px]">
            <DialogHeader>
              <DialogTitle>{t("createTitle")}</DialogTitle>
            </DialogHeader>
            <ProductForm
              onSubmit={async (data) => {
                await handleCreate(data as CreateProductFormData);
              }}
              onCancel={() => setIsCreateDialogOpen(false)}
              isLoading={createProduct.isPending}
            />
          </DialogContent>
        </Dialog>
      )}

      {/* Edit Dialog */}
      {hasEditPermission && editingProduct && editingProductData?.data && (
        <Dialog open={!!editingProduct} onOpenChange={(open) => !open && setEditingProduct(null)}>
          <DialogContent className="sm:max-w-[500px]">
            <DialogHeader>
              <DialogTitle>{t("editTitle")}</DialogTitle>
            </DialogHeader>
            <ProductForm
              product={editingProductData.data}
              onSubmit={async (data) => {
                await handleUpdate(data as UpdateProductFormData);
              }}
              onCancel={() => setEditingProduct(null)}
              isLoading={updateProduct.isPending}
            />
          </DialogContent>
        </Dialog>
      )}

      {/* Product Detail Modal */}
      {hasViewPermission && (
        <ProductDetailModal
          productId={viewingProductId}
          open={isDetailModalOpen}
          onOpenChange={(open) => {
            setIsDetailModalOpen(open);
            if (!open) {
              setViewingProductId(null);
            }
          }}
          onProductUpdated={() => {
            // Refresh will be handled by query invalidation in hooks
          }}
        />
      )}

      {/* Delete Dialog */}
      {hasDeletePermission && (
        <DeleteDialog
          open={!!deletingProductId}
          onOpenChange={(open) => {
            if (!open) {
              setDeletingProductId(null);
            }
          }}
          onConfirm={handleDeleteConfirm}
          title={t("deleteTitle")}
          description={
            deletingProductId
              ? t("deleteDescriptionWithName", {
                  name:
                    products.find((p: Product) => p.id === deletingProductId)?.name ??
                    t("deleteDescription"),
                })
              : t("deleteDescription")
          }
          itemName="product"
          isLoading={deleteProduct.isPending}
        />
      )}
      </div>
    </div>
  );
}
