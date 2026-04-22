"use client";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";
import { useIsMobile } from "@/hooks/use-mobile";

export interface Column<T> {
  id: string;
  header: string;
  accessor: (row: T) => React.ReactNode;
  className?: string;
}

interface DataTableProps<T> {
  readonly columns: readonly Column<T>[];
  readonly data: readonly T[];
  readonly isLoading?: boolean;
  readonly emptyMessage?: string;
  readonly pagination?: {
    readonly page: number;
    readonly per_page: number;
    readonly total: number;
    readonly total_pages: number;
    readonly has_next: boolean;
    readonly has_prev: boolean;
  };
  readonly onPageChange?: (page: number) => void;
  readonly onPerPageChange?: (perPage: number) => void;
  readonly itemName?: string; // e.g., "diagnosis", "procedure"
  readonly perPageOptions?: readonly number[]; // e.g., [10, 20, 50, 100]
  readonly onResetFilters?: () => void;
}

export function DataTable<T extends { id: string }>({
  columns,
  data,
  isLoading = false,
  emptyMessage = "No data found",
  pagination,
  onPageChange,
  onPerPageChange,
  itemName = "item",
  perPageOptions = [10, 20, 50, 100],
  onResetFilters,
}: DataTableProps<T>) {
  const isMobile = useIsMobile();

  const getPageNumbers = () => {
    if (!pagination) return [];

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

  // Mobile card view
  if (isMobile) {
    return (
      <div className="border border-border/60 rounded-lg bg-card shadow-sm overflow-hidden">
        {isLoading ? (
          <div className="p-4 space-y-3">
            {Array.from({ length: 5 }, (_, i) => (
              <Skeleton key={`skeleton-card-${i}`} className="h-32 w-full rounded-lg" />
            ))}
          </div>
        ) : (
          <>
            {data.length === 0 ? (
              <div className="text-center text-muted-foreground py-12 px-4">
                <div className="text-sm">{emptyMessage}</div>
              </div>
            ) : (
              <div className="divide-y divide-border/50">
                {data.map((row) => (
                  <div
                    key={row.id}
                    className="p-4 hover:bg-accent/40 transition-all duration-200 border-b border-border/30 last:border-b-0 active:bg-accent/50"
                  >
                    <div className="space-y-2.5">
                      {columns.map((column, index) => {
                        const isFirstColumn = index === 0;
                        return (
                          <div
                            key={column.id}
                            className={cn(
                              "flex flex-col gap-1",
                              isFirstColumn && "pb-2 border-b border-primary/20"
                            )}
                          >
                            <div className="text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                              {column.header}
                            </div>
                            <div className={cn(
                              "text-sm transition-colors",
                              isFirstColumn && "font-semibold text-foreground"
                            )}>
                              {column.accessor(row)}
                            </div>
                          </div>
                        );
                      })}
                    </div>
                  </div>
                ))}
              </div>
            )}

            {/* Mobile Pagination */}
            {pagination && (
              <div className="border-t border-border/60 bg-gradient-to-b from-muted/40 to-muted/20 px-4 py-3 space-y-3">
                {/* Rows per page selector */}
                {onPerPageChange && (
                  <div className="flex items-center justify-between">
                    <Label htmlFor="rows-per-page-mobile" className="text-sm font-medium text-foreground">
                      Rows per page
                    </Label>
                    <Select
                      value={String(pagination.per_page)}
                      onValueChange={(value) => {
                        onPerPageChange?.(Number(value));
                        onPageChange?.(1);
                      }}
                    >
                      <SelectTrigger
                        id="rows-per-page-mobile"
                        className="w-20 h-8 text-sm border-border/60 hover:border-primary/40 transition-colors"
                      >
                        <SelectValue />
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
                )}

                {/* Page info */}
                <div className="text-center text-sm text-muted-foreground">
                  <span className="text-foreground font-semibold">
                    {(pagination.page - 1) * pagination.per_page + 1}-
                    {Math.min(
                      pagination.page * pagination.per_page,
                      pagination.total
                    )}
                  </span>{" "}
                  of{" "}
                  <span className="text-foreground font-semibold">
                    {pagination.total}
                  </span>
                </div>

                {/* Pagination controls */}
                {pagination.total_pages > 1 && (
                  <div className="flex items-center justify-center gap-1">
                    <Pagination>
                      <PaginationContent>
                        <PaginationItem>
                          <PaginationPrevious
                            onClick={() =>
                              onPageChange?.(Math.max(1, pagination.page - 1))
                            }
                            disabled={!pagination.has_prev || isLoading}
                            className={cn(
                              "transition-all duration-200",
                              (!pagination.has_prev || isLoading) &&
                                "pointer-events-none opacity-50 cursor-not-allowed"
                            )}
                            aria-disabled={!pagination.has_prev || isLoading}
                          />
                        </PaginationItem>

                        <PaginationItem>
                          <span className="px-3 py-1 text-sm font-medium text-foreground">
                            Page {pagination.page} of {pagination.total_pages}
                          </span>
                        </PaginationItem>

                        <PaginationItem>
                          <PaginationNext
                            onClick={() => onPageChange?.(pagination.page + 1)}
                            disabled={!pagination.has_next || isLoading}
                            className={cn(
                              "transition-all duration-200",
                              (!pagination.has_next || isLoading) &&
                                "pointer-events-none opacity-50 cursor-not-allowed"
                            )}
                            aria-disabled={!pagination.has_next || isLoading}
                          />
                        </PaginationItem>
                      </PaginationContent>
                    </Pagination>
                  </div>
                )}
              </div>
            )}
          </>
        )}
      </div>
    );
  }

  // Desktop table view
  return (
    <div className="border border-border/60 rounded-lg bg-card shadow-sm overflow-hidden">
      {isLoading ? (
        <div className="p-4 space-y-3">
          {Array.from({ length: 5 }, (_, i) => (
            <Skeleton key={`skeleton-row-${i}`} className="h-10 w-full rounded-md" />
          ))}
        </div>
      ) : (
        <>
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow className="bg-gradient-to-b from-muted/50 to-muted/30 border-b border-border/60 hover:bg-muted/50">
                  {columns.map((column) => (
                    <TableHead
                      key={column.id}
                      className={cn(
                        "font-semibold text-foreground/90 uppercase tracking-wider text-xs px-6 py-4",
                        column.className
                      )}
                    >
                      {column.header}
                    </TableHead>
                  ))}
                </TableRow>
              </TableHeader>
              <TableBody>
                {data.length === 0 ? (
                  <TableRow>
                    <TableCell
                      colSpan={columns.length}
                      className="text-center text-muted-foreground py-12 px-6"
                    >
                      <div className="text-sm">{emptyMessage}</div>
                    </TableCell>
                  </TableRow>
                ) : (
                  data.map((row, rowIndex) => (
                    <TableRow 
                      key={row.id} 
                      className={cn(
                        "hover:bg-accent/40 transition-all duration-200 border-b border-border/30",
                        "hover:shadow-sm",
                        rowIndex % 2 === 0 && "bg-card",
                        rowIndex % 2 === 1 && "bg-muted/20"
                      )}
                    >
                      {columns.map((column) => (
                        <TableCell
                          key={`${row.id}-${column.id}`}
                          className={cn(
                            "transition-colors duration-150 px-6 py-4",
                            column.className
                          )}
                        >
                          {column.accessor(row)}
                        </TableCell>
                      ))}
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>

          {pagination && (
            <div className="border-t border-border/60 bg-gradient-to-b from-muted/40 to-muted/20 px-6 py-4">
              <div className="flex flex-col lg:flex-row items-center justify-between gap-4">
                {/* Rows per page selector */}
                {onPerPageChange && (
                  <div className="flex items-center gap-3 order-3 lg:order-1">
                    <Label htmlFor="rows-per-page" className="text-sm whitespace-nowrap font-medium text-foreground">
                      Rows per page
                    </Label>
                    <Select
                      value={String(pagination.per_page)}
                      onValueChange={(value) => {
                        onPerPageChange?.(Number(value));
                        // Reset to page 1 when changing per page
                        onPageChange?.(1);
                      }}
                    >
                      <SelectTrigger
                        id="rows-per-page"
                        className="w-fit whitespace-nowrap h-9 border-border/60 hover:border-primary/40 transition-colors"
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
                )}

                {/* Page number information */}
                <div className="flex grow justify-center lg:justify-end text-sm whitespace-nowrap text-muted-foreground order-2 lg:order-2">
                  <p className="text-sm whitespace-nowrap" aria-live="polite">
                    <span className="text-foreground font-semibold">
                      {(pagination.page - 1) * pagination.per_page + 1}-
                      {Math.min(
                        pagination.page * pagination.per_page,
                        pagination.total
                      )}
                    </span>{" "}
                    of{" "}
                    <span className="text-foreground font-semibold">
                      {pagination.total}
                    </span>
                  </p>
                </div>

                {/* Pagination controls + optional reset filters */}
                {pagination.total_pages > 1 && (
                  <div className="flex items-center gap-3 order-1 lg:order-3">
                    {onResetFilters && (
                      <button
                        type="button"
                        onClick={() => onResetFilters()}
                        className="text-xs font-medium text-muted-foreground hover:text-primary transition-colors duration-200 underline-offset-4 hover:underline cursor-pointer"
                      >
                        Reset filters
                      </button>
                    )}
                    <Pagination>
                      <PaginationContent>
                        <PaginationItem>
                          <PaginationFirst
                            onClick={() => onPageChange?.(1)}
                            disabled={!pagination.has_prev || isLoading}
                            className={cn(
                              "transition-all duration-200",
                              (!pagination.has_prev || isLoading) &&
                                "pointer-events-none opacity-50 cursor-not-allowed",
                              pagination.has_prev &&
                                !isLoading &&
                                "hover:bg-accent hover:text-accent-foreground"
                            )}
                            aria-disabled={!pagination.has_prev || isLoading}
                          />
                        </PaginationItem>

                        <PaginationItem>
                          <PaginationPrevious
                            onClick={() =>
                              onPageChange?.(Math.max(1, pagination.page - 1))
                            }
                            disabled={!pagination.has_prev || isLoading}
                            className={cn(
                              "transition-all duration-200",
                              (!pagination.has_prev || isLoading) &&
                                "pointer-events-none opacity-50 cursor-not-allowed",
                              pagination.has_prev &&
                                !isLoading &&
                                "hover:bg-accent hover:text-accent-foreground"
                            )}
                            aria-disabled={!pagination.has_prev || isLoading}
                          />
                        </PaginationItem>

                        {getPageNumbers().map((pageNum) => {
                          if (
                            pageNum === "ellipsis-start" ||
                            pageNum === "ellipsis-end"
                          ) {
                            return (
                              <PaginationItem key={`ellipsis-${pageNum}`}>
                                <PaginationEllipsis />
                              </PaginationItem>
                            );
                          }

                          const pageNumber = pageNum as number;
                          const isActive = pageNumber === pagination.page;

                          return (
                            <PaginationItem key={pageNumber}>
                              <PaginationLink
                                onClick={() => onPageChange?.(pageNumber)}
                                disabled={isLoading}
                                isActive={isActive}
                                className={cn(
                                  "transition-all duration-200",
                                  isLoading &&
                                    "pointer-events-none opacity-50 cursor-not-allowed",
                                  isActive &&
                                    "bg-primary text-primary-foreground shadow-sm hover:bg-primary/90 font-semibold",
                                  !isActive &&
                                    "hover:bg-accent hover:text-accent-foreground"
                                )}
                              >
                                {pageNumber}
                              </PaginationLink>
                            </PaginationItem>
                          );
                        })}

                        <PaginationItem>
                          <PaginationNext
                            onClick={() => onPageChange?.(pagination.page + 1)}
                            disabled={!pagination.has_next || isLoading}
                            className={cn(
                              "transition-all duration-200",
                              (!pagination.has_next || isLoading) &&
                                "pointer-events-none opacity-50 cursor-not-allowed",
                              pagination.has_next &&
                                !isLoading &&
                                "hover:bg-accent hover:text-accent-foreground"
                            )}
                            aria-disabled={!pagination.has_next || isLoading}
                          />
                        </PaginationItem>

                        <PaginationItem>
                          <PaginationLast
                            onClick={() =>
                              onPageChange?.(pagination.total_pages)
                            }
                            disabled={!pagination.has_next || isLoading}
                            className={cn(
                              "transition-all duration-200",
                              (!pagination.has_next || isLoading) &&
                                "pointer-events-none opacity-50 cursor-not-allowed",
                              pagination.has_next &&
                                !isLoading &&
                                "hover:bg-accent hover:text-accent-foreground"
                            )}
                            aria-disabled={!pagination.has_next || isLoading}
                          />
                        </PaginationItem>
                      </PaginationContent>
                    </Pagination>
                  </div>
                )}
              </div>
            </div>
          )}
        </>
      )}
    </div>
  );
}

