"use client";

import { useTranslations } from "next-intl";
import { useDashboardOverview } from "../hooks/useDashboard";
import { Skeleton } from "@/components/ui/skeleton";
import { Card, CardContent } from "@/components/ui/card";
import { Target, Users, Briefcase, TrendingUp, TrendingDown } from "lucide-react";
import { getCurrentMonthDateRange } from "../dashboard-date-util";

import { useHasPermission } from "@/features/auth/providers/permissions-provider";
import { Lock } from "lucide-react";

export function DashboardStatsPlain({ startDate, endDate }: Readonly<{ startDate?: string; endDate?: string }>) {
  const t = useTranslations("dashboardOverview");
  const monthParams = getCurrentMonthDateRange();
  
  // Permission check
  const hasPermission = useHasPermission("dashboard.view");
  
  const params = { ...monthParams, start_date: startDate, end_date: endDate };
  
  const { data, isLoading } = useDashboardOverview(params, { enabled: hasPermission });

  console.log("DashboardStatsPlain Params:", params);
  console.log("DashboardStatsPlain Data:", data);

  // Show warning if no permission - showing 4 restricted cards to maintain layout
  if (!hasPermission) {
    return (
      <div className="grid grid-cols-2 gap-3 sm:gap-4 md:gap-6">
        {Array.from({ length: 4 }).map((_, i) => (
          <Card key={`restricted-${i}`} className="border-0 shadow-sm bg-muted/10 h-full">
            <CardContent className="flex flex-col items-center justify-center p-4 text-center h-[120px]">
              <Lock className="h-5 w-5 text-muted-foreground mb-2" />
              <p className="text-xs font-medium text-muted-foreground">Access Restricted</p>
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  if (isLoading) {
    return (
      <div className="grid grid-cols-2 gap-3 sm:gap-4 md:gap-6">
        {Array.from({ length: 4 }).map((_, i) => (
          <div key={i} className="flex-1 p-3 sm:p-4 rounded-lg border bg-card">
            <Skeleton className="h-3 sm:h-4 w-20 sm:w-24 mb-2" />
            <Skeleton className="h-6 sm:h-8 w-16 sm:w-20 mb-1" />
            <Skeleton className="h-2.5 sm:h-3 w-content" />
          </div>
        ))}
      </div>
    );
  }

  const overview = data?.data;

  if (!overview) {
    return null;
  }

  const targetProgress = overview.target?.progress_percent ?? 0;
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
  const leadStats = overview.lead_stats ?? {
    total: 0,
    new: 0,
    contacted: 0,
    qualified: 0,
    converted: 0,
    lost: 0,
    change_percent: 0,
  };

  const stats = [
    {
      label: "Target", 
      value: overview.target?.target_amount_formatted ?? "Rp 0",
      description: `${Math.round(targetProgress)}% Achieved`,
      icon: Target,
      changePercent: overview.target?.change_percent ?? 0,
    },
    {
      label: t("totalAccounts.title"),
      value: accountStats.total.toString(),
      description: t("totalAccounts.description", {
        active: accountStats.active,
        inactive: accountStats.inactive,
      }),
      icon: Users,
      changePercent: accountStats.change_percent,
    },
    {
      label: t("totalDeals.title"),
      value: deals.total_deals.toString(),
      description: t("totalDeals.description", {
        open: deals.open_deals,
        won: deals.won_deals,
      }),
      icon: Briefcase,
      changePercent: deals.change_percent,
    },
    {
      label: t("totalLeads.title"),
      value: leadStats.total.toString(),
      description: t("totalLeads.description", {
        qualified: leadStats.qualified,
        converted: leadStats.converted,
      }),
      icon: TrendingUp,
      changePercent: leadStats.change_percent,
    },
  ];

  return (
    <div className="grid grid-cols-2 gap-3 sm:gap-4 md:gap-6">
      {stats.map((stat) => {
        const Icon = stat.icon;
        const hasChange = stat.changePercent !== 0;
        const isPositive = stat.changePercent > 0;

        return (
          <Card key={stat.label} className="border-0 shadow-sm bg-card overflow-hidden">
            <CardContent className="px-4">
              <div className="flex items-center justify-between mb-3">
                <span className="text-xs sm:text-sm font-medium text-muted-foreground">{stat.label}</span>
                <Icon className="h-3.5 w-3.5 sm:h-4 sm:w-4 text-primary" />
              </div>
              <div className="flex items-baseline gap-2 mb-1">
                <span className="text-xl sm:text-2xl font-bold tracking-tight">{stat.value}</span>
                {hasChange && (
                  <div className={`flex items-center gap-0.5 text-[10px] font-medium ${isPositive ? "text-green-600" : "text-red-600"}`}>
                    {isPositive ? (
                      <TrendingUp className="h-2.5 w-2.5" />
                    ) : (
                      <TrendingDown className="h-2.5 w-2.5" />
                    )}
                    <span>{Math.abs(stat.changePercent).toFixed(1)}%</span>
                  </div>
                )}
              </div>
              <p className="text-[10px] sm:text-xs text-muted-foreground line-clamp-1">{stat.description}</p>
            </CardContent>
          </Card>
        );
      })}
    </div>
  );
}

