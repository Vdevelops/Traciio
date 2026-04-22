"use client";

import * as React from "react";
import { useState } from "react";
import { useTranslations } from "next-intl";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
  PaginationFirst,
  PaginationLast,
} from "@/components/ui/pagination";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";
import { ArrowUpDown, ArrowDown, ArrowUp, Package } from "lucide-react";
import type { ProductListItem } from "../types";
import { format } from "date-fns";
import { ProductDetailModal } from "@/features/sales-crm/product-management/components/product-detail-modal";

type SortBy = "total_sold" | "revenue" | "profit" | "name";
type OrderBy = "asc" | "desc";

interface ProductListTableProps {
  readonly data: ProductListItem[];
  readonly isLoading?: boolean;
  readonly sortBy: SortBy;
  readonly orderBy: OrderBy;
  readonly onSortByChange: (value: SortBy) => void;
  readonly onOrderByChange: (value: OrderBy) => void;
  readonly pagination: {
    readonly page: number;
    readonly per_page: number;
    readonly total: number;
    readonly total_pages: number;
    readonly has_next: boolean;
    readonly has_prev: boolean;
  };
  readonly onPageChange: (page: number) => void;
  readonly onPerPageChange: (perPage: number) => void;
}

// Local format currency helper to match Leaderboard (raw integer)
const formatCurrencyValue = (value: number): string => {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(value);
};

export function ProductListTable({
  data,
  isLoading = false,
  sortBy,
  orderBy,
  onSortByChange,
  onOrderByChange,
  pagination,
  onPageChange,
  onPerPageChange,
}: ProductListTableProps) {
  const t = useTranslations("productAnalytics.productList");
  const [selectedProductId, setSelectedProductId] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  const handleProductClick = (productId: string) => {
    setSelectedProductId(productId);
    setIsModalOpen(true);
  };

  const handleSort = (column: SortBy) => {
    if (sortBy === column) {
      // Toggle order if same column
      onOrderByChange(orderBy === "asc" ? "desc" : "asc");
    } else {
      // Set new column with desc order by default
      onSortByChange(column);
      onOrderByChange("desc");
    }
  };

  const getSortIcon = (column: SortBy) => {
    if (sortBy !== column) {
      return <ArrowUpDown className="h-4 w-4 ml-1 opacity-50" />;
    }
    return orderBy === "asc" ? (
      <ArrowUp className="h-4 w-4 ml-1" />
    ) : (
      <ArrowDown className="h-4 w-4 ml-1" />
    );
  };

  const getPageNumbers = () => {
    const totalPages = pagination.total_pages;
    const currentPage = pagination.page;
    const pages: (number | string)[] = [];

    if (totalPages <= 7) {
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      pages.push(1);
      if (currentPage > 3) {
        pages.push("ellipsis-start");
      }
      const start = Math.max(2, currentPage - 1);
      const end = Math.min(totalPages - 1, currentPage + 1);
      for (let i = start; i <= end; i++) {
        pages.push(i);
      }
      if (currentPage < totalPages - 2) {
        pages.push("ellipsis-end");
      }
      pages.push(totalPages);
    }
    return pages;
  };

  const perPageOptions = [10, 20, 50, 100];

  return (
    <div className="space-y-4">
      {/* Table - Minimalist (No Card/Border wrapper) */}
      <div className="overflow-x-auto">
        <Table>
          <TableHeader>
            <TableRow className="border-b border-border/60 hover:bg-transparent">
              <TableHead className="w-[80px] font-medium text-foreground/90 uppercase tracking-wider text-xs px-4 py-3">
                Rank
              </TableHead>
              <TableHead className="w-[60px] font-medium text-foreground/90 uppercase tracking-wider text-xs px-4 py-3"></TableHead>
              <TableHead className="font-medium text-foreground/90 uppercase tracking-wider text-xs px-4 py-3">
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-8 px-2 -ml-2 hover:bg-transparent hover:text-primary"
                  onClick={() => handleSort("name")}
                >
                  {t("table.product")}
                  {getSortIcon("name")}
                </Button>
              </TableHead>
              <TableHead className="font-medium text-foreground/90 uppercase tracking-wider text-xs px-4 py-3">{t("table.category")}</TableHead>
              <TableHead className="text-right font-medium text-foreground/90 uppercase tracking-wider text-xs px-4 py-3">
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-8 px-2 ml-auto hover:bg-transparent hover:text-primary"
                  onClick={() => handleSort("total_sold")}
                >
                  {t("table.totalSold")}
                  {getSortIcon("total_sold")}
                </Button>
              </TableHead>
              <TableHead className="text-right font-medium text-foreground/90 uppercase tracking-wider text-xs px-4 py-3">
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-8 px-2 ml-auto hover:bg-transparent hover:text-primary"
                  onClick={() => handleSort("revenue")}
                >
                  {t("table.revenue")}
                  {getSortIcon("revenue")}
                </Button>
              </TableHead>
              <TableHead className="text-right font-medium text-foreground/90 uppercase tracking-wider text-xs px-4 py-3">
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-8 px-2 ml-auto hover:bg-transparent hover:text-primary"
                  onClick={() => handleSort("profit")}
                >
                  {t("table.profit")}
                  {getSortIcon("profit")}
                </Button>
              </TableHead>
              <TableHead className="text-right font-medium text-foreground/90 uppercase tracking-wider text-xs px-4 py-3">{t("table.avgPrice")}</TableHead>
              <TableHead className="text-right font-medium text-foreground/90 uppercase tracking-wider text-xs px-4 py-3">{t("table.salesCount")}</TableHead>
              <TableHead className="font-medium text-foreground/90 uppercase tracking-wider text-xs px-4 py-3">{t("table.lastSold")}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
          {isLoading ? (
            Array.from({ length: 5 }).map((_, index) => (
              <TableRow key={index} className="border-b border-border/30">
                <TableCell className="px-4 py-3"><Skeleton className="h-4 w-8" /></TableCell>
                <TableCell className="px-4 py-3"><Skeleton className="h-10 w-10 rounded-md" /></TableCell>
                <TableCell className="px-4 py-3"><Skeleton className="h-4 w-[200px]" /></TableCell>
                <TableCell className="px-4 py-3"><Skeleton className="h-4 w-[100px]" /></TableCell>
                <TableCell className="text-right px-4 py-3"><Skeleton className="h-4 w-16 ml-auto" /></TableCell>
                <TableCell className="text-right px-4 py-3"><Skeleton className="h-4 w-24 ml-auto" /></TableCell>
                <TableCell className="text-right px-4 py-3"><Skeleton className="h-4 w-24 ml-auto" /></TableCell>
                <TableCell className="text-right px-4 py-3"><Skeleton className="h-4 w-24 ml-auto" /></TableCell>
                <TableCell className="text-right px-4 py-3"><Skeleton className="h-4 w-16 ml-auto" /></TableCell>
                <TableCell className="px-4 py-3"><Skeleton className="h-4 w-24" /></TableCell>
              </TableRow>
            ))
          ) : data.length === 0 ? (
            <TableRow>
              <TableCell colSpan={10} className="h-24 text-center text-muted-foreground py-12 px-4">
                <div className="text-sm">{t("table.noData")}</div>
              </TableCell>
            </TableRow>
          ) : (
            data.map((item, index) => {
              const rank = (pagination.page - 1) * pagination.per_page + index + 1;
              let rankDisplay: React.ReactNode = <span className="text-muted-foreground font-medium">#{rank}</span>;
              
              if (rank === 1) {
                rankDisplay = (
                   <Badge variant="outline" className="bg-yellow-500/10 text-yellow-600 border-yellow-200 hover:bg-yellow-500/20 px-2 py-0.5 whitespace-nowrap">
                     🥇 1
                   </Badge>
                );
              } else if (rank === 2) {
                rankDisplay = (
                  <Badge variant="outline" className="bg-slate-300/10 text-slate-600 border-slate-300 hover:bg-slate-300/20 px-2 py-0.5 whitespace-nowrap">
                    🥈 2
                  </Badge>
                );
              } else if (rank === 3) {
                 rankDisplay = (
                  <Badge variant="outline" className="bg-orange-700/10 text-orange-700 border-orange-300 hover:bg-orange-700/20 px-2 py-0.5 whitespace-nowrap">
                    🥉 3
                  </Badge>
                );
              }

              return (
                <TableRow 
                  key={item.product_id}
                  className="hover:bg-muted/30 transition-colors border-b border-border/30"
                >
                  <TableCell className="px-4 py-3">
                    {rankDisplay}
                  </TableCell>
                  <TableCell className="px-4 py-3">
                    <div className="h-10 w-10 rounded-md bg-primary/10 flex items-center justify-center">
                      <Package className="h-5 w-5 text-primary" />
                    </div>
                  </TableCell>
                  <TableCell className="px-4 py-3">
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
                  </TableCell>
                  <TableCell className="px-4 py-3">
                    <Badge variant="outline" className="font-normal">{item.category_name || "-"}</Badge>
                  </TableCell>
                  <TableCell className="text-right font-medium px-4 py-3">
                    {item.total_sold.toLocaleString()}
                  </TableCell>
                  <TableCell className="text-right font-medium px-4 py-3">
                    {formatCurrencyValue(item.total_revenue)}
                  </TableCell>
                  <TableCell className="text-right font-medium text-green-600 px-4 py-3">
                    {formatCurrencyValue(item.total_profit)}
                  </TableCell>
                  <TableCell className="text-right px-4 py-3">
                    {formatCurrencyValue(item.avg_unit_price)}
                  </TableCell>
                  <TableCell className="text-right px-4 py-3">
                    {item.sales_count.toLocaleString()}
                  </TableCell>
                  <TableCell className="text-muted-foreground text-sm px-4 py-3">
                    {item.last_sold_at
                      ? format(new Date(item.last_sold_at), "dd MMM yyyy")
                      : "-"}
                  </TableCell>
                </TableRow>
              );
            })
          )}
          </TableBody>
        </Table>
      </div>
      
      {/* Pagination Controls - Minimalist */}
      <div className="py-2">
        <div className="flex flex-col lg:flex-row items-center justify-between gap-4">
          {/* Rows per page selector */}
          <div className="flex items-center gap-3 order-3 lg:order-1">
            <Label htmlFor="rows-per-page" className="text-sm whitespace-nowrap font-medium text-muted-foreground">
              Rows per page
            </Label>
            <Select
              value={String(pagination.per_page)}
              onValueChange={(value) => {
                onPerPageChange(Number(value));
                onPageChange(1);
              }}
            >
              <SelectTrigger
                id="rows-per-page"
                className="w-fit whitespace-nowrap h-8 border-border/60"
              >
                <SelectValue placeholder="Select rows" />
              </SelectTrigger>
              <SelectContent>
                {perPageOptions.map((option) => (
                  <SelectItem key={option} value={String(option)}>
                    {option}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          {/* Page number information */}
          <div className="flex grow justify-center lg:justify-end text-sm whitespace-nowrap text-muted-foreground order-2 lg:order-2">
            <p className="text-sm whitespace-nowrap">
              <span className="text-foreground font-medium">
                {(pagination.page - 1) * pagination.per_page + 1}-
                {Math.min(
                  pagination.page * pagination.per_page,
                  pagination.total
                )}
              </span>{" "}
              of{" "}
              <span className="text-foreground font-medium">
                {pagination.total}
              </span>
            </p>
          </div>

          {/* Pagination Controls */}
          {pagination.total_pages > 1 && (
            <div className="flex items-center gap-2 order-1 lg:order-3">
              <Pagination>
                <PaginationContent>
                  <PaginationItem>
                    <PaginationFirst
                      onClick={() => onPageChange(1)}
                      className={cn(
                        "cursor-pointer h-8 w-8",
                        (!pagination.has_prev || isLoading) && "pointer-events-none opacity-50"
                      )}
                    />
                  </PaginationItem>
                  <PaginationItem>
                    <PaginationPrevious
                      onClick={() => onPageChange(Math.max(1, pagination.page - 1))}
                      className={cn(
                        "cursor-pointer h-8 w-8",
                        (!pagination.has_prev || isLoading) && "pointer-events-none opacity-50"
                      )}
                    />
                  </PaginationItem>
                  
                  {getPageNumbers().map((pageNum, i) => {
                    if (pageNum === "ellipsis-start" || pageNum === "ellipsis-end") {
                      return (
                        <PaginationItem key={`ellipsis-${i}`}>
                          <PaginationEllipsis className="h-8 w-8" />
                        </PaginationItem>
                      );
                    }
                    
                    const p = pageNum as number;
                    const isActive = p === pagination.page;
                    
                    return (
                      <PaginationItem key={p}>
                        <PaginationLink
                          onClick={() => onPageChange(p)}
                          isActive={isActive}
                          className={cn(
                            "cursor-pointer h-8 w-8",
                            isLoading && "pointer-events-none opacity-50",
                            isActive && "hover:bg-primary/90 hover:text-primary-foreground"
                          )}
                        >
                          {p}
                        </PaginationLink>
                      </PaginationItem>
                    );
                  })}

                  <PaginationItem>
                    <PaginationNext
                      onClick={() => onPageChange(Math.min(pagination.total_pages, pagination.page + 1))}
                      className={cn(
                        "cursor-pointer h-8 w-8",
                        (!pagination.has_next || isLoading) && "pointer-events-none opacity-50"
                      )}
                    />
                  </PaginationItem>
                  <PaginationItem>
                    <PaginationLast
                      onClick={() => onPageChange(pagination.total_pages)}
                      className={cn(
                        "cursor-pointer h-8 w-8",
                        (!pagination.has_next || isLoading) && "pointer-events-none opacity-50"
                      )}
                    />
                  </PaginationItem>
                </PaginationContent>
              </Pagination>
            </div>
          )}
        </div>
      </div>

      {/* Product Detail Modal */}
      <ProductDetailModal
        productId={selectedProductId}
        open={isModalOpen}
        onOpenChange={setIsModalOpen}
        onProductUpdated={() => {
          // Optionally refresh the product list if needed
        }}
      />
    </div>
  );
}
