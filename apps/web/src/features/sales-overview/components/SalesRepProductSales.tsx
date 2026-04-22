"use client";

import * as React from "react";
import { useState } from "react";
import { useTranslations } from "next-intl";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { formatCurrency } from "@/lib/utils";
import type { ProductListItem } from "@/features/product-analytics/types";
import { format } from "date-fns";
import { ProductDetailModal } from "@/features/sales-crm/product-management/components/product-detail-modal";
import { ProductImage } from "@/features/sales-crm/product-management/components/product-image";
import { useUserProductSales } from "@/features/product-analytics/hooks/useProductAnalytics";
import { DataTable, type Column } from "@/components/ui/data-table";

type SortBy = "total_sold" | "revenue" | "profit" | "name";
type OrderBy = "asc" | "desc";

interface SalesRepProductSalesProps {
  readonly userId: string;
  readonly startDate?: string;
  readonly endDate?: string;
}

export function SalesRepProductSales({ userId, startDate, endDate }: SalesRepProductSalesProps) {
  const t = useTranslations("salesOverview.productSales");
  const tFilters = useTranslations("productAnalytics.productList.filters");
  const tTable = useTranslations("productAnalytics.productList.table");
  
  const {
    products,
    pagination,
    isLoading,
    sortBy,
    setSortBy,
    orderBy,
    setOrderBy,
    page,
    setPage,
    perPage,
    setPerPage,
  } = useUserProductSales(userId, {
    startDate,
    endDate,
  });

  const [selectedProductId, setSelectedProductId] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  const handleProductClick = (productId: string) => {
    setSelectedProductId(productId);
    setIsModalOpen(true);
  };

  // Map products to include id property for DataTable
  const productsWithId = products.map((p: ProductListItem) => ({ ...p, id: p.product_id }));

  // Define columns for DataTable
  const columns: Column<ProductListItem & { id: string }>[] = [
    {
      id: "image",
      header: "",
      accessor: (item) => (
        <ProductImage
          src={item.image_url}
          alt={item.product_name}
          className="h-10 w-10 rounded-md"
        />
      ),
    },
    {
      id: "product",
      header: tTable("product"),
      accessor: (item) => (
        <div className="flex flex-col">
          <button
            onClick={() => handleProductClick(item.product_id)}
            className="font-medium text-primary hover:underline cursor-pointer text-left"
          >
            {item.product_name}
          </button>
          <span className="text-xs text-muted-foreground">
            SKU: {item.product_sku}
          </span>
        </div>
      ),
    },
    {
      id: "category",
      header: tTable("category"),
      accessor: (item) => (
        <Badge variant="outline">{item.category_name || "-"}</Badge>
      ),
    },
    {
      id: "total_sold",
      header: tTable("totalSold"),
      accessor: (item) => (
        <span className="font-medium">{item.total_sold.toLocaleString()}</span>
      ),
      className: "text-right",
    },
    {
      id: "revenue",
      header: tTable("revenue"),
      accessor: (item) => (
        <span className="font-medium">{formatCurrency(item.total_revenue)}</span>
      ),
      className: "text-right",
    },
    {
      id: "profit",
      header: tTable("profit"),
      accessor: (item) => (
        <span className="font-medium text-green-600">
          {formatCurrency(item.total_profit)}
        </span>
      ),
      className: "text-right",
    },
    {
      id: "avg_price",
      header: tTable("avgPrice"),
      accessor: (item) => formatCurrency(item.avg_unit_price),
      className: "text-right",
    },
    {
      id: "last_sold",
      header: tTable("lastSold"),
      accessor: (item) =>
        item.last_sold_at
          ? format(new Date(item.last_sold_at), "dd MMM yyyy")
          : "-",
      className: "text-muted-foreground text-sm",
    },
  ];

  return (
    <div className="space-y-4">
      {/* Filters */}
      <div className="flex flex-wrap gap-4 items-center">
        <div className="flex items-center gap-2">
          <label className="text-sm font-medium">{tFilters("sortBy")}</label>
          <Select value={sortBy} onValueChange={(v) => {
            setSortBy(v as SortBy);
            setPage(1);
          }}>
            <SelectTrigger className="w-[180px]">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="total_sold">{tFilters("sortOptions.totalSold")}</SelectItem>
              <SelectItem value="revenue">{tFilters("sortOptions.revenue")}</SelectItem>
              <SelectItem value="profit">{tFilters("sortOptions.profit")}</SelectItem>
              <SelectItem value="name">{tFilters("sortOptions.name")}</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div className="flex items-center gap-2">
          <label className="text-sm font-medium">{tFilters("order")}</label>
          <Select value={orderBy} onValueChange={(v) => {
            setOrderBy(v as OrderBy);
            setPage(1);
          }}>
            <SelectTrigger className="w-[140px]">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="desc">{tFilters("orderOptions.desc")}</SelectItem>
              <SelectItem value="asc">{tFilters("orderOptions.asc")}</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </div>

      {/* Data Table with Pagination */}
      <DataTable
        columns={columns}
        data={productsWithId}
        isLoading={isLoading}
        emptyMessage={t("noProducts")}
        pagination={pagination}
        onPageChange={setPage}
        onPerPageChange={setPerPage}
        itemName="product"
        perPageOptions={[10, 20, 50, 100]}
      />

      {/* Product Detail Modal */}
      <ProductDetailModal
        productId={selectedProductId}
        open={isModalOpen}
        onOpenChange={setIsModalOpen}
        onProductUpdated={() => {
          // Refresh handled by react-query
        }}
      />
    </div>
  );
}
