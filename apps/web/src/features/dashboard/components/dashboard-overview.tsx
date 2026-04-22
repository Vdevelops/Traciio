"use client";

import { useTranslations } from "next-intl";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useDashboardOverview } from "../hooks/useDashboard";
import { Skeleton } from "@/components/ui/skeleton";
import { Target, Users, Briefcase, DollarSign, TrendingUp } from "lucide-react";
import { getCurrentMonthDateRange } from "../dashboard-date-util";

interface DashboardOverviewProps {
  startDate?: string;
  endDate?: string;
}

export function DashboardOverview({ startDate, endDate }: Readonly<DashboardOverviewProps>) {
  const t = useTranslations("dashboardOverview");

  // Use provided date range or fallback (though parent should provide generic defaults)
  const queryParams = {
    start_date: startDate,
    end_date: endDate,
  };

  const { data, isLoading } = useDashboardOverview(queryParams);

  if (isLoading) {
    return (
      <div className="grid gap-3 sm:gap-4 md:gap-6 grid-cols-2 md:grid-cols-3 lg:grid-cols-5">
        {["a", "b", "c", "d", "e"].map((key) => (
          <Card key={`dashboard-overview-skeleton-${key}`} className="border-0 shadow-sm">
            <CardHeader className="px-3 sm:px-6">
              <Skeleton className="h-3 sm:h-4 w-20 sm:w-24" />
            </CardHeader>
            <CardContent className="px-3 sm:px-6 pb-3 sm:pb-6">
              <Skeleton className="h-6 sm:h-8 w-12 sm:w-16 mb-1 sm:mb-2" />
              <Skeleton className="h-2.5 sm:h-3 w-24 sm:w-32" />
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  const overview = data?.data;

  if (!overview) {
    return null;
  }

  const target = overview.target ?? {
    target_amount: 0,
    target_amount_formatted: "Rp 0",
    achieved_amount: 0,
    achieved_amount_formatted: "Rp 0",
    progress_percent: 0,
    change_percent: 0,
  };
  const targetProgress = target.progress_percent ?? 0;
  const accountStats = overview.account_stats ?? { total: 0, active: 0, inactive: 0, change_percent: 0 };
  const deals = overview.deals ?? { 
    total_deals: 0, 
    open_deals: 0, 
    won_deals: 0, 
    lost_deals: 0, 
    total_value: 0, 
    total_value_formatted: "Rp 0", 
    change_percent: 0 
  };
  const revenue = overview.revenue ?? { 
    total_revenue: 0, 
    total_revenue_formatted: "Rp 0", 
    change_percent: 0 
  };
  const leadStats = overview.lead_stats ?? {
    total: 0,
    new: 0,
    contacted: 0,
    qualified: 0,
    converted: 0,
    lost: 0,
    change_percent: 0,
  };

  const formatCurrency = (value: number) =>
    new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
    }).format(value);

  return (
    <div className="grid gap-3 sm:gap-4 md:gap-6 grid-cols-2 md:grid-cols-3 lg:grid-cols-5">
      {/* Target progress */}
      <Card className="border-0 shadow-sm">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 px-3 sm:px-6">
          <CardTitle className="text-xs sm:text-sm font-medium text-muted-foreground">
            {t("targetProgress.title")}
          </CardTitle>
          <Target className="h-3 w-3 sm:h-4 sm:w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent className="px-3 sm:px-6 pb-3 sm:pb-6">
          <div className="text-xl sm:text-2xl font-medium">
            {Math.round(targetProgress)}%
          </div>
          <p className="text-[10px] sm:text-xs text-muted-foreground mt-0.5 sm:mt-1">
            {t("targetProgress.description", {
              progress: Math.round(targetProgress),
            })}
          </p>
        </CardContent>
      </Card>

      {/* Total customers/accounts */}
      <Card className="border-0 shadow-sm">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 px-3 sm:px-6">
          <CardTitle className="text-xs sm:text-sm font-medium text-muted-foreground">
            {t("totalAccounts.title")}
          </CardTitle>
          <Users className="h-3 w-3 sm:h-4 sm:w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent className="px-3 sm:px-6 pb-3 sm:pb-6">
          <div className="text-xl sm:text-2xl font-medium">
            {accountStats.total}
          </div>
          <p className="text-[10px] sm:text-xs text-muted-foreground mt-0.5 sm:mt-1">
            {t("totalAccounts.description", {
              active: accountStats.active,
              inactive: accountStats.inactive,
            })}
          </p>
        </CardContent>
      </Card>

      {/* Total deals */}
      <Card className="border-0 shadow-sm">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 px-3 sm:px-6">
          <CardTitle className="text-xs sm:text-sm font-medium text-muted-foreground">
            {t("totalDeals.title")}
          </CardTitle>
          <Briefcase className="h-3 w-3 sm:h-4 sm:w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent className="px-3 sm:px-6 pb-3 sm:pb-6">
          <div className="text-xl sm:text-2xl font-medium">
            {deals.total_deals}
          </div>
          <p className="text-[10px] sm:text-xs text-muted-foreground mt-0.5 sm:mt-1">
            {t("totalDeals.description", {
              open: deals.open_deals,
              won: deals.won_deals,
            })}
          </p>
        </CardContent>
      </Card>

      {/* Total revenue */}
      <Card className="border-0 shadow-sm">
        <CardHeader className="flex flex-row items-center justify-between space-y-0 px-3 sm:px-6">
          <CardTitle className="text-xs sm:text-sm font-medium text-muted-foreground">
            {t("totalRevenue.title")}
          </CardTitle>
          <DollarSign className="h-3 w-3 sm:h-4 sm:w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent className="px-3 sm:px-6 pb-3 sm:pb-6">
          <div className="text-lg sm:text-2xl font-medium wrap-break-word">
            {revenue.total_revenue_formatted ?? formatCurrency(revenue.total_revenue)}
          </div>
          {target.target_amount > 0 ? (
            <div className="mt-1 space-y-0.5">
              <p className="text-[10px] sm:text-xs text-muted-foreground">
                {revenue.total_revenue_formatted ?? formatCurrency(revenue.total_revenue)} / {target.target_amount_formatted}
              </p>
              <p className={`text-[10px] sm:text-xs font-medium ${
                targetProgress >= 100 ? "text-green-600" : 
                targetProgress >= 75 ? "text-yellow-600" : 
                "text-red-600"
              }`}>
                {Math.round(targetProgress)}% dari target
              </p>
            </div>
          ) : (
            <p className="text-[10px] sm:text-xs text-muted-foreground mt-0.5 sm:mt-1">
              {t("totalRevenue.description")}
            </p>
          )}
        </CardContent>
      </Card>

      {/* Total leads */}
      <Card className="border-0 shadow-sm">
        <CardHeader className="pb-1.5 sm:pb-2 flex flex-row items-center justify-between space-y-0 px-3 sm:px-6">
          <CardTitle className="text-xs sm:text-sm font-medium text-muted-foreground">
            {t("totalLeads.title")}
          </CardTitle>
          <TrendingUp className="h-3 w-3 sm:h-4 sm:w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent className="px-3 sm:px-6 pb-3 sm:pb-6">
          <div className="text-xl sm:text-2xl font-medium">
            {leadStats.total}
          </div>
          <p className="text-[10px] sm:text-xs text-muted-foreground mt-0.5 sm:mt-1">
            {t("totalLeads.description", {
              qualified: leadStats.qualified,
              converted: leadStats.converted,
            })}
          </p>
        </CardContent>
      </Card>
    </div>
  );
}

