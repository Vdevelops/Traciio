"use client";

import * as React from "react";
import { useState, Suspense, useEffect } from "react";
import dynamic from "next/dynamic";
import { useTranslations } from "next-intl";
import { PieChart, Search } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { DateRangePicker } from "@/components/ui/date-range-picker";
import { motion } from "framer-motion";
import type { DateRange } from "react-day-picker";
import { 
  useProductsList, 
  useMonthlySales 
} from "@/features/product-analytics/hooks/useProductAnalytics";
import type { MonthlySalesChartProps } from "@/features/product-analytics/components/MonthlySalesChart";

// Dynamic imports for heavy components
const ProductListTable = dynamic(
  () => import("@/features/product-analytics/components/ProductListTable").then(m => ({ default: m.ProductListTable })),
  { loading: () => <TableLoading /> }
);

const MonthlySalesChart = dynamic<MonthlySalesChartProps>(
  () => import("@/features/product-analytics/components/MonthlySalesChart").then(m => ({ default: m.MonthlySalesChart })),
  { loading: () => <ChartLoading /> }
);

// Loading Skeletons
function TableLoading() {
  return (
    <div className="pt-6">
      <Skeleton className="h-[400px] w-full" />
    </div>
  );
}

function ChartLoading() {
  return (
    <Card>
      <CardHeader>
        <Skeleton className="h-6 w-48" />
        <Skeleton className="h-4 w-64 mt-2" />
      </CardHeader>
      <CardContent>
        <Skeleton className="h-[400px] w-full" />
      </CardContent>
    </Card>
  );
}

// Section wrapper
function Section({ children, delay = 0 }: { children: React.ReactNode; delay?: number }) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3, ease: [0.4, 0, 0.2, 1], delay }}
    >
      {children}
    </motion.div>
  );
}

export function ProductAnalyticsPageClient() {
  const t = useTranslations("productAnalytics");

  // Metric State for Chart
  const [selectedMetric, setSelectedMetric] = useState<"total_sold" | "total_revenue" | "total_profit" | "sales_count">("total_sold");

  // Search State with Debounce
  const [searchInput, setSearchInput] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");
  
  // Date Filter Mode
  const [filterMode, setFilterMode] = useState<"year" | "range">("year");
  const [selectedYear, setSelectedYear] = useState<number>(() => new Date().getFullYear());

  // Date Range State for "range" mode
  const [dateRange, setDateRange] = useState<DateRange | undefined>(() => {
    const today = new Date();
    const startOfYear = new Date(today.getFullYear(), 0, 1);
    return { from: startOfYear, to: today };
  });

  const startDateStr = React.useMemo(() => {
    if (filterMode === "year") {
      return `${selectedYear}-01-01`;
    }
    if (!dateRange?.from) return undefined;
    const d = dateRange.from;
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
  }, [filterMode, selectedYear, dateRange]);

  const endDateStr = React.useMemo(() => {
    if (filterMode === "year") {
      return `${selectedYear}-12-31`;
    }
    if (!dateRange?.to) return undefined;
    const d = dateRange.to;
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
  }, [filterMode, selectedYear, dateRange]);

  const handleDateRangeChange = (range: DateRange | undefined) => {
    setDateRange(range);
    setFilterMode("range");
  };

  const handleYearChange = (year: number) => {
    setSelectedYear(year);
    setFilterMode("year");
  };

  // Debounce search input
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedSearch(searchInput);
    }, 500);
    
    return () => clearTimeout(timer);
  }, [searchInput]);

  // Data Hooks
  
  // 1. Monthly Sales Chart (Range-based)
  const { 
    monthlySales, 
    isLoading: isLoadingMonthlySales, 
  } = useMonthlySales(startDateStr, endDateStr);

  // 2. Product List (Table) - Paginated & Filtered by Date Range + Search
  const {
    products: productsList,
    isLoading: isLoadingProductsList,
    sortBy,
    setSortBy,
    orderBy,
    setOrderBy,
    pagination,
    setPage,
    perPage,
    setPerPage,
  } = useProductsList("all", startDateStr, endDateStr, debouncedSearch);

  // Reset page when filters change
  useEffect(() => {
    setPage(1);
  }, [startDateStr, endDateStr, debouncedSearch, setPage]);

  return (
    <div className="space-y-8">
      {/* HEADER */}
      <Section delay={0}>
        <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
          <div>
            <h1 className="text-3xl font-medium tracking-tight flex items-center gap-2">
              <PieChart className="h-8 w-8 text-primary" />
              {t("title")}
            </h1>
            <p className="text-muted-foreground mt-1 text-sm">{t("description")}</p>
          </div>
        </div>
      </Section>

      {/* MONTHLY SALES CHART */}
      <Section delay={0.1}>
        <div className="grid gap-4 grid-cols-1">
            <Suspense fallback={<ChartLoading />}>
                <MonthlySalesChart
                data={monthlySales}
                isLoading={isLoadingMonthlySales}
                filterMode={filterMode}
                onFilterModeChange={setFilterMode}
                selectedYear={selectedYear}
                onYearChange={handleYearChange}
                dateRange={dateRange}
                onDateChange={handleDateRangeChange}
                selectedMetric={selectedMetric}
                onMetricChange={setSelectedMetric}
                />
            </Suspense>
        </div>
      </Section>

      {/* PRODUCT LIST TABLE */}
      <Section delay={0.2}>
        <div className="space-y-4">
            <div className="flex items-center justify-between">
                <div>
                    <h2 className="text-xl font-semibold tracking-tight">{t("productList.title")}</h2>
                    <p className="text-sm text-muted-foreground">{t("productList.description")}</p>
                </div>
                {/* Search Input */}
                <div className="relative w-[300px]">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                    <Input
                    placeholder={t("productList.searchPlaceholder")}
                    value={searchInput}
                    onChange={(e) => setSearchInput(e.target.value)}
                    className="pl-10 h-9"
                    />
                </div>
            </div>
            <Suspense fallback={<TableLoading />}>
            <ProductListTable
                data={productsList}
                isLoading={isLoadingProductsList}
                sortBy={sortBy}
                onSortByChange={setSortBy}
                orderBy={orderBy}
                onOrderByChange={setOrderBy}
                pagination={pagination}
                onPageChange={setPage}
                onPerPageChange={setPerPage}
            />
            </Suspense>
        </div>
      </Section>
    </div>
  );
}
