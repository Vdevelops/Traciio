"use client";

import { useState, useMemo } from "react";
import { useTranslations } from "next-intl";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { useSalesRepDetail } from "@/features/sales-overview/hooks/useSalesRepDetail";
import { useSalesRepCheckInLocations } from "@/features/sales-overview/hooks/useSalesRepCheckInLocations";
import { SalesRepStatistics } from "@/features/sales-overview/components/SalesRepStatistics";
import { SalesRepDetailTabs } from "@/features/sales-overview/components/sales-rep-detail-tabs";
import { Skeleton } from "@/components/ui/skeleton";
import { ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useRouter } from "next/navigation";
import { DateRangePicker } from "@/components/ui/date-range-picker";
import type { DateRange } from "react-day-picker";

interface SalesRepDetailPageClientProps {
  readonly userId: string;
}

// Helper function to get default monthly date range (month to date)
const getDefaultMonthlyRange = () => {
  const today = new Date();
  today.setHours(23, 59, 59, 999);
  
  const startOfMonth = new Date(today);
  startOfMonth.setDate(1);
  startOfMonth.setHours(0, 0, 0, 0);
  
  const startDateStr = `${startOfMonth.getFullYear()}-${String(startOfMonth.getMonth() + 1).padStart(2, "0")}-${String(startOfMonth.getDate()).padStart(2, "0")}`;
  const endDateStr = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, "0")}-${String(today.getDate()).padStart(2, "0")}`;
  
  return { startDate: startDateStr, endDate: endDateStr };
};

export function SalesRepDetailPageClient({ userId }: SalesRepDetailPageClientProps) {
  const t = useTranslations("salesOverview");
  const router = useRouter();
  
  // Set default to monthly (month to date)
  const defaultRange = useMemo(() => getDefaultMonthlyRange(), []);
  
  // State for date range filter (same as table list)
  const [startDate, setStartDate] = useState<string>(defaultRange.startDate);
  const [endDate, setEndDate] = useState<string>(defaultRange.endDate);
  
  // Convert date strings to DateRange for DateRangePicker
  const dateRange: DateRange | undefined = (() => {
    if (startDate && endDate) {
      const fromDate = new Date(startDate + "T00:00:00");
      fromDate.setHours(0, 0, 0, 0);
      const toDate = new Date(endDate + "T00:00:00");
      toDate.setHours(0, 0, 0, 0);
      return { from: fromDate, to: toDate };
    }
    if (startDate) {
      const fromDate = new Date(startDate + "T00:00:00");
      fromDate.setHours(0, 0, 0, 0);
      return { from: fromDate, to: undefined };
    }
    return undefined;
  })();
  
  // Handle date range change
  const handleDateRangeChange = (range: DateRange | undefined) => {
    if (range?.from) {
      const fromDate = new Date(range.from);
      fromDate.setHours(0, 0, 0, 0);
      const fromStr = `${fromDate.getFullYear()}-${String(fromDate.getMonth() + 1).padStart(2, "0")}-${String(fromDate.getDate()).padStart(2, "0")}`;
      setStartDate(fromStr);
      
      if (range.to) {
        const toDate = new Date(range.to);
        toDate.setHours(0, 0, 0, 0);
        const toStr = `${toDate.getFullYear()}-${String(toDate.getMonth() + 1).padStart(2, "0")}-${String(toDate.getDate()).padStart(2, "0")}`;
        setEndDate(toStr);
      } else {
        // If only from date, use default end date (today) for better UX
        setEndDate(defaultRange.endDate);
      }
    } else {
      // Clear date filter - reset to default monthly range for performance
      setStartDate(defaultRange.startDate);
      setEndDate(defaultRange.endDate);
    }
  };
  
  // Prepare params for API calls (same filtering as table list)
  // Always include date range for query optimization (default is monthly)
  const detailParams = (() => {
    const params: { start_date?: string; end_date?: string } = {};
    
    // Use provided dates or fallback to default monthly range for performance
    if (startDate && startDate.trim() !== "") {
      params.start_date = startDate.trim();
    } else {
      params.start_date = defaultRange.startDate;
    }
    
    if (endDate && endDate.trim() !== "") {
      params.end_date = endDate.trim();
    } else {
      params.end_date = defaultRange.endDate;
    }
    
    return params;
  })();
  
  const { detail, isLoading: detailLoading } = useSalesRepDetail(userId, detailParams);
  
  // Use same date range for check-in locations
  const { 
    locations, 
    isLoading: locationsLoading,
    totalVisits,
    page: locationsPage,
    setPage: setLocationsPage,
    perPage: locationsPerPage,
    setPerPage: setLocationsPerPage,
  } = useSalesRepCheckInLocations(userId, detailParams);

  if (detailLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-10 w-64" />
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  if (!detail?.user) {
    return (
      <div className="space-y-6">
        <Button variant="ghost" onClick={() => router.back()} className="mb-4">
          <ArrowLeft className="h-4 w-4 mr-2" />
          {t("back")}
        </Button>
        <Card>
          <CardContent className="py-8 text-center">
            <p className="text-muted-foreground">{t("sales_rep_not_found")}</p>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between gap-4 flex-wrap">
        <div>
          <Button variant="ghost" onClick={() => router.back()} className="mb-2">
            <ArrowLeft className="h-4 w-4 mr-2" />
            {t("back")}
          </Button>
          <h1 className="text-3xl font-medium">{detail.user.name}</h1>
          <p className="text-muted-foreground mt-1">{detail.user.email}</p>
        </div>
        <div className="flex items-center gap-3">
          <div>
            <DateRangePicker
              dateRange={dateRange}
              onDateChange={handleDateRangeChange}
              placeholder="All Time"
            />
          </div>
        </div>
      </div>

      {/* Statistics Cards */}
      <SalesRepStatistics statistics={detail.statistics} />

      {/* Check-in Locations, Products, and Customers Tabs */}
      <Card>
        <CardHeader>
          <CardTitle>{t("check_in_locations")}, {t("products_sold")} & {t("customers.title")}</CardTitle>
          <CardDescription>
            {t("check_in_locations_description")}, {t("products_sold_description")} and customer information
          </CardDescription>
        </CardHeader>
        <CardContent>
          <SalesRepDetailTabs
            userId={userId}
            startDate={startDate}
            endDate={endDate}
            checkInLocationsProps={{
              locations,
              isLoading: locationsLoading,
              totalVisits,
              page: locationsPage,
              perPage: locationsPerPage,
              onPageChange: setLocationsPage,
              onPerPageChange: setLocationsPerPage,
            }}
          />
        </CardContent>
      </Card>
    </div>
  );
}

