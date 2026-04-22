"use client";

import { Eye, ArrowUpDown, ArrowUp, ArrowDown, Search } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/avatar";
import { useRouter } from "next/navigation";
import { useTranslations } from "next-intl";
import type { SalesPerformanceListItem } from "../types";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
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
import { cn, formatEmailToMailto } from "@/lib/utils";

type SortBy = "revenue" | "deals" | "visits" | "tasks" | "name" | "target" | "achievement";

interface SalesPerformanceListProps {
  page: number;
  setPage: (page: number) => void;
  perPage: number;
  setPerPage: (perPage: number) => void;
  search: string;
  setSearch: (search: string) => void;
  startDate: string;
  endDate: string;
  setStartDate: (date: string) => void;
  setEndDate: (date: string) => void;
  performanceList: SalesPerformanceListItem[];
  pagination?: {
    page: number;
    per_page: number;
    total: number;
    total_pages: number;
    has_next: boolean;
    has_prev: boolean;
  };
  isLoading: boolean;
  sortBy: SortBy;
  setSortBy: (sortBy: SortBy) => void;
  order: "asc" | "desc";
  setOrder: (order: "asc" | "desc") => void;
}

export function SalesPerformanceList({
  setPage,
  setPerPage,
  search,
  setSearch,
  startDate,
  endDate,
  setStartDate,
  setEndDate,
  performanceList,
  pagination,
  isLoading,
  sortBy,
  setSortBy,
  order,
  setOrder,
}: Readonly<SalesPerformanceListProps>) {

  const router = useRouter();
  const t = useTranslations("salesOverview");

  const handleViewUser = (userId: string) => {
    router.push(`/sales-overview/sales-rep/${userId}`);
  };

  // Convert date strings to DateRange for DateRangePicker


  const handleSort = (column: SortBy) => {
    if (sortBy === column) {
      setOrder(order === "asc" ? "desc" : "asc");
    } else {
      setSortBy(column);
      setOrder("desc");
    }
  };

  const getSortIcon = (column: SortBy) => {
    if (sortBy !== column) {
      return <ArrowUpDown className="h-4 w-4 ml-1 opacity-50" />;
    }
    return order === "asc" ? (
      <ArrowUp className="h-4 w-4 ml-1" />
    ) : (
      <ArrowDown className="h-4 w-4 ml-1" />
    );
  };

  const perPageOptions = [10, 20, 50, 100];

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

  return (
    <div className="space-y-4">
      {/* Header with Search (Date Filter follows Chart) */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-semibold tracking-tight">{t("performance_list")}</h2>
          <p className="text-sm text-muted-foreground">{t("performance_list_desc")}</p>
        </div>
        <div className="relative w-[300px]">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" aria-hidden="true" />
          <Input 
            className="flex w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground placeholder:opacity-50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 pl-10 h-9" 
            placeholder={t("table.searchPlaceholder")} 
            value={search}
            onChange={(e) => {
              setSearch(e.target.value);
              setPage(1);
            }}
          />
        </div>
      </div>
      
      {/* Table - Minimalist (No Card/Border wrapper) */}
      <div className="overflow-x-auto">
        <Table>
          <TableHeader>
            <TableRow className="border-b border-border/60 hover:bg-transparent">
              <TableHead className="w-[80px] font-medium text-foreground/90 uppercase tracking-wider text-xs px-4 py-3">
                {t("table.rank")}
              </TableHead>
              <TableHead className="font-medium text-foreground/90 uppercase tracking-wider text-xs px-4 py-3">
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-8 px-2 -ml-2 hover:bg-transparent hover:text-primary"
                  onClick={() => handleSort("name")}
                >
                  {t("table.user")}
                  {getSortIcon("name")}
                </Button>
              </TableHead>
              <TableHead className="font-medium text-foreground/90 uppercase tracking-wider text-xs px-4 py-3">
                {t("table.email")}
              </TableHead>
              <TableHead className="text-right font-medium text-foreground/90 uppercase tracking-wider text-xs px-4 py-3">
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-8 px-2 ml-auto hover:bg-transparent hover:text-primary"
                  onClick={() => handleSort("revenue")}
                >
                  Monthly Target / Achievement
                  {getSortIcon("revenue")}
                </Button>
              </TableHead>
              <TableHead className="text-right font-medium text-foreground/90 uppercase tracking-wider text-xs px-4 py-3">
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-8 px-2 ml-auto hover:bg-transparent hover:text-primary"
                  onClick={() => handleSort("deals")}
                >
                  {t("table.deals")}
                  {getSortIcon("deals")}
                </Button>
              </TableHead>
              <TableHead className="text-right font-medium text-foreground/90 uppercase tracking-wider text-xs px-4 py-3">
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-8 px-2 ml-auto hover:bg-transparent hover:text-primary"
                  onClick={() => handleSort("visits")}
                >
                  {t("table.visits")}
                  {getSortIcon("visits")}
                </Button>
              </TableHead>
              <TableHead className="text-right font-medium text-foreground/90 uppercase tracking-wider text-xs px-4 py-3">
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-8 px-2 ml-auto hover:bg-transparent hover:text-primary"
                  onClick={() => handleSort("tasks")}
                >
                  {t("table.tasks")}
                  {getSortIcon("tasks")}
                </Button>
              </TableHead>
              <TableHead className="text-right font-medium text-foreground/90 uppercase tracking-wider text-xs px-4 py-3">
                {t("table.actions")}
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              Array.from({ length: 5 }).map((_, index) => (
                <TableRow key={index} className="border-b border-border/30">
                  <TableCell className="px-4 py-3"><Skeleton className="h-4 w-8" /></TableCell>
                  <TableCell className="px-4 py-3">
                    <div className="flex items-center gap-3">
                      <Skeleton className="h-8 w-8 rounded-full" />
                      <Skeleton className="h-4 w-24" />
                    </div>
                  </TableCell>
                  <TableCell className="px-4 py-3"><Skeleton className="h-4 w-32" /></TableCell>
                  <TableCell className="text-right px-4 py-3">
                    <div className="flex flex-col items-end gap-1">
                      <Skeleton className="h-4 w-28" />
                      <Skeleton className="h-3 w-16" />
                    </div>
                  </TableCell>
                  <TableCell className="text-right px-4 py-3"><Skeleton className="h-4 w-12 ml-auto" /></TableCell>
                  <TableCell className="text-right px-4 py-3"><Skeleton className="h-4 w-12 ml-auto" /></TableCell>
                  <TableCell className="text-right px-4 py-3"><Skeleton className="h-4 w-8 ml-auto" /></TableCell>
                  <TableCell className="text-right px-4 py-3"><Skeleton className="h-8 w-8 ml-auto rounded" /></TableCell>
                </TableRow>
              ))
            ) : performanceList.length === 0 ? (
              <TableRow>
                <TableCell colSpan={8} className="h-24 text-center text-muted-foreground py-12 px-4">
                  <div className="text-sm">{t("table.noData")}</div>
                </TableCell>
              </TableRow>
            ) : (
              performanceList.map((item, index) => {
                const rank = ((pagination?.page || 1) - 1) * (pagination?.per_page || 10) + index + 1;
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
                    key={item.user_id}
                    className="hover:bg-muted/30 transition-colors border-b border-border/30"
                  >
                    <TableCell className="px-4 py-3">
                      {rankDisplay}
                    </TableCell>
                    <TableCell className="px-4 py-3">
                      <button
                        onClick={() => handleViewUser(item.user_id)}
                        className="flex items-center gap-3 font-medium text-primary hover:underline cursor-pointer text-left"
                      >
                         <Avatar className="h-8 w-8">
                            <AvatarImage 
                              src={item.avatar_url || `https://api.dicebear.com/7.x/avataaars/svg?seed=${encodeURIComponent(item.user_name)}`} 
                              alt={item.user_name} 
                            />
                            <AvatarFallback className="bg-transparent text-transparent">
                              {/* Fallback hidden as we force dicebear */}
                            </AvatarFallback>
                          </Avatar>
                        <span>{item.user_name}</span>
                      </button>
                    </TableCell>
                    <TableCell className="px-4 py-3 text-muted-foreground">
                      <a href={formatEmailToMailto(item.user_email)} className="hover:text-primary hover:underline cursor-pointer min-w-0">
                        {item.user_email}
                      </a>
                    </TableCell>
                    <TableCell className="text-right px-4 py-3">
                      {(item.target_amount_formatted || item.total_revenue_formatted) ? (
                        <div className="flex flex-col items-end gap-0.5">
                          <div className="text-sm">
                            <span className="text-muted-foreground">{item.target_amount_formatted || "-"}</span>
                            <span className="mx-1 text-muted-foreground">/</span>
                            <span className="font-medium">{item.total_revenue_formatted}</span>
                          </div>
                          {item.target_achievement_percentage !== undefined && (() => {
                            const pct = item.target_achievement_percentage;
                            let colorClass = "text-red-600";
                            if (pct >= 100) colorClass = "text-green-600";
                            else if (pct >= 75) colorClass = "text-yellow-600";
                            return (
                              <p className={cn("text-[10px] font-medium", colorClass)}>
                                {Math.round(pct)}% dari target
                              </p>
                            );
                          })()}
                        </div>
                      ) : (
                        <span className="text-muted-foreground">-</span>
                      )}
                    </TableCell>
                    <TableCell className="text-right font-medium px-4 py-3">
                      {item.deals_closed}
                    </TableCell>
                    <TableCell className="text-right font-medium px-4 py-3">
                      {item.visits_completed}
                    </TableCell>
                    <TableCell className="text-right font-medium px-4 py-3">
                      {item.tasks_completed}
                    </TableCell>
                    <TableCell className="text-right px-4 py-3">
                      <div className="flex items-center justify-end gap-1">
                        <Button
                          variant="ghost"
                          size="icon-sm"
                          className="h-8 w-8 cursor-pointer text-muted-foreground hover:text-primary"
                          title="View Details"
                          onClick={() => handleViewUser(item.user_id)}
                        >
                          <Eye className="h-4 w-4" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                );
              })
            )}
          </TableBody>
        </Table>
      </div>

       {/* Pagination Controls - Minimalist */}
      {pagination && (
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
                  setPerPage(Number(value));
                  setPage(1);
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
                        onClick={() => setPage(1)}
                        className={cn(
                          "cursor-pointer h-8 w-8",
                          (!pagination.has_prev || isLoading) && "pointer-events-none opacity-50"
                        )}
                      />
                    </PaginationItem>
                    <PaginationItem>
                      <PaginationPrevious
                        onClick={() => setPage(Math.max(1, pagination.page - 1))}
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
                            onClick={() => setPage(p)}
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
                         onClick={() => setPage(Math.min(pagination.total_pages, pagination.page + 1))}
                        className={cn(
                          "cursor-pointer h-8 w-8",
                          (!pagination.has_next || isLoading) && "pointer-events-none opacity-50"
                        )}
                      />
                    </PaginationItem>
                    <PaginationItem>
                      <PaginationLast
                         onClick={() => setPage(pagination.total_pages)}
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
      )}
    </div>
  );
}
