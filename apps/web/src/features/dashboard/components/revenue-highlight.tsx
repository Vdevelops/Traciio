"use client";

import { useTranslations } from "next-intl";
import { Card, CardContent } from "@/components/ui/card";
import { useDashboardOverview } from "../hooks/useDashboard";
import { Skeleton } from "@/components/ui/skeleton";
import { DollarSign, TrendingUp, TrendingDown } from "lucide-react";
import { getCurrentMonthDateRange } from "../dashboard-date-util";

const formatCurrency = (value: number) =>
  new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
  }).format(value);

export function RevenueHighlight({ startDate, endDate }: Readonly<{ startDate?: string; endDate?: string }>) {
  const t = useTranslations("dashboardOverview");
  const monthParams = getCurrentMonthDateRange();
  
  const params = { ...monthParams, start_date: startDate, end_date: endDate };
  
  const { data, isLoading } = useDashboardOverview(params);

  console.log("RevenueHighlight Params:", params);
  console.log("RevenueHighlight Data:", data);

  if (isLoading) {
    return (
      <Card className="border shadow-sm">
        <CardContent className="px-4">
          <Skeleton className="h-4 sm:h-6 w-24 sm:w-32 mb-2 sm:mb-4" />
          <Skeleton className="h-8 sm:h-12 w-48 sm:w-64 mb-1 sm:mb-2" />
          <Skeleton className="h-3 sm:h-4 w-36 sm:w-48" />
        </CardContent>
      </Card>
    );
  }

  const overview = data?.data;
  const revenue = overview?.revenue ?? { 
    total_revenue: 0, 
    total_revenue_formatted: "Rp 0", 
    change_percent: 0 
  };

  const changePercent = revenue.change_percent ?? 0;
  const hasChange = changePercent !== 0;
  const isPositive = changePercent > 0;
  const revenueFormatted = revenue.total_revenue_formatted ?? formatCurrency(revenue.total_revenue ?? 0);

  return (
    <Card className="border shadow-sm h-full flex flex-col justify-center bg-linear-to-br from-card to-primary/5">
      <CardContent className="p-6">
        <div className="flex items-center justify-between">
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 mb-4">
              <DollarSign className="h-5 w-5 text-primary" />
              <h3 className="text-sm sm:text-base font-semibold text-muted-foreground truncate uppercase tracking-wider">{t("totalRevenue.title")}</h3>
            </div>
            <div className="flex flex-col gap-2 mb-4">
              <span className="text-3xl sm:text-4xl font-bold tracking-tight text-foreground">{revenueFormatted}</span>
              {hasChange && (
                <div className={`inline-flex items-center w-fit gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold ${
                  isPositive ? "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400" : "bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400"
                }`}>
                  {isPositive ? (
                    <TrendingUp className="h-3.5 w-3.5" />
                  ) : (
                    <TrendingDown className="h-3.5 w-3.5" />
                  )}
                  <span>{Math.abs(changePercent).toFixed(1)}%</span>
                </div>
              )}
            </div>
            <p className="text-xs text-muted-foreground leading-relaxed max-w-[240px] font-medium opacity-80">{t("totalRevenue.description")}</p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

