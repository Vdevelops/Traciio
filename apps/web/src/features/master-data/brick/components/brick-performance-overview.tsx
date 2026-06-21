"use client";

import { useTranslations } from "next-intl";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Target,
  Users,
  Briefcase,
  DollarSign,
  TrendingUp,
  TrendingDown,
} from "lucide-react";
import { useBrickPerformance } from "../hooks/useBrickAnalytics";
import { cn, formatCurrency } from "@/lib/utils";

interface BrickPerformanceOverviewProps {
  brickId: string;
  periodStart?: string;
  periodEnd?: string;
}

export function BrickPerformanceOverview({ brickId, periodStart, periodEnd }: Readonly<BrickPerformanceOverviewProps>) {
  const t = useTranslations("brickDashboard");
  const requestStart = periodStart || undefined;
  const requestEnd = periodEnd || undefined;
  const { data, isLoading } = useBrickPerformance(brickId, requestStart, requestEnd);

  if (isLoading) {
    return (
      <div className="grid gap-3 sm:gap-4 md:gap-6 grid-cols-2 md:grid-cols-4">
        {[0, 1, 2, 3].map((i) => (
          <Card key={`brick-performance-skeleton-${i}`} className="border-0 shadow-sm">
            <CardHeader className="pb-1.5 sm:pb-2 px-3 sm:px-6 pt-3 sm:pt-6">
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

  const metrics = data?.data;
  if (!metrics) return null;

  const stats = [
    {
      label: t("totalSales.title"),
      value: metrics.total_sales.toString(),
      description: t("totalSales.description", { active: metrics.active_sales }),
      icon: Users,
    },
    {
      label: t("totalDeals.title"),
      value: metrics.total_deals.toString(),
      description: t("totalDeals.description", { open: metrics.open_deals, won: metrics.won_deals }),
      icon: Briefcase,
      changePercent: metrics.win_rate,
    },
    {
      label: t("totalRevenue.title"),
      value: formatCurrency(metrics.total_revenue),
      description: t("totalRevenue.description", { thisMonth: formatCurrency(metrics.revenue_this_month) }),
      icon: DollarSign,
      changePercent: metrics.revenue_growth_percentage,
    },
    {
      label: t("targetAchievement.title"),
      value: `${Math.round(metrics.achievement_percentage)}%`,
      description: t("targetAchievement.description", {
        achieved: formatCurrency(metrics.target_achieved),
        target: formatCurrency(metrics.monthly_target),
      }),
      icon: Target,
      changePercent: metrics.achievement_percentage,
    },
  ];

  return (
    <div className="grid gap-3 sm:gap-4 md:gap-6 grid-cols-2 md:grid-cols-4">
      {stats.map((stat) => {
        const Icon = stat.icon;
        const hasChange = stat.changePercent !== undefined && stat.changePercent !== null;
        const isPositive = hasChange && stat.changePercent! > 0;

        return (
          <Card key={stat.label} className="border-0 shadow-sm">
            <CardHeader className="pb-1.5 sm:pb-2 flex flex-row items-center justify-between space-y-0 px-3 sm:px-6 pt-3 sm:pt-6">
              <CardTitle className="text-xs sm:text-sm font-medium text-muted-foreground">
                {stat.label}
              </CardTitle>
              <Icon className="h-3 w-3 sm:h-4 sm:w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent className="px-3 sm:px-6 pb-3 sm:pb-6">
              <div className="flex items-baseline gap-1.5 sm:gap-2 flex-wrap">
                <div className="text-xl sm:text-2xl font-medium">{stat.value}</div>
                {hasChange && (
                  <div className={cn(
                    "flex items-center gap-0.5 sm:gap-1",
                    isPositive ? "text-green-600" : "text-destructive"
                  )}>
                    {isPositive ? (
                      <TrendingUp className="h-2.5 w-2.5 sm:h-3 sm:w-3" />
                    ) : (
                      <TrendingDown className="h-2.5 w-2.5 sm:h-3 sm:w-3" />
                    )}
                    <span className="text-[10px] sm:text-xs font-medium">
                      {Math.abs(stat.changePercent!).toFixed(1)}%
                    </span>
                  </div>
                )}
              </div>
              <p className="text-[10px] sm:text-xs text-muted-foreground mt-0.5 sm:mt-1">
                {stat.description}
              </p>
            </CardContent>
          </Card>
        );
      })}
    </div>
  );
}
