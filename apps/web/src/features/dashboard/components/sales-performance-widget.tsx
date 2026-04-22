"use client";

import { useTranslations } from "next-intl";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { useSalesPerformanceList } from "@/features/sales-overview/hooks/useSalesPerformanceList";
import { Avatar, AvatarImage } from "@/components/ui/avatar";
import { TrendingUp, Users } from "lucide-react";
import Link from "next/link";
import type { SalesPerformanceListItem } from "@/features/sales-overview/types";

import { useHasPermission } from "@/features/auth/providers/permissions-provider";
import { Lock } from "lucide-react";

interface SalesPerformanceWidgetProps {
  startDate?: string;
  endDate?: string;
}

export function SalesPerformanceWidget({ startDate, endDate }: Readonly<SalesPerformanceWidgetProps>) {
  const t = useTranslations("dashboard");
  
  // Permission check
  const hasPermission = useHasPermission("sales-overview.view");

  // Pass date filters to the hook
  const { performanceList, isLoading } = useSalesPerformanceList(5, { startDate, endDate });

  console.log("SalesPerformanceWidget Params:", { startDate, endDate });
  console.log("SalesPerformanceWidget Data:", performanceList);

  // Show warning if no permission
  if (!hasPermission) {
    return (
      <Card className="h-full border-0 shadow-sm bg-muted/10">
        <CardHeader className="px-3 sm:px-6">
          <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-2 sm:gap-0">
            <div className="flex-1 min-w-0">
              <CardTitle className="flex items-center gap-1.5 sm:gap-2 text-sm sm:text-base text-muted-foreground">
                <Lock className="h-4 w-4 sm:h-5 sm:w-5" />
                {t("salesPerformance.title")}
              </CardTitle>
            </div>
          </div>
        </CardHeader>
        <CardContent className="flex flex-col items-center justify-center p-6 text-center h-[200px]">
          <div className="bg-muted/20 p-3 rounded-full mb-3">
            <Lock className="h-6 w-6 text-muted-foreground" />
          </div>
          <p className="text-sm font-medium text-muted-foreground mb-1">
            Access Restricted
          </p>
          <p className="text-xs text-muted-foreground/70 max-w-[200px]">
            You don't have permission to view sales performance.
          </p>
        </CardContent>
      </Card>
    );
  }

  if (isLoading) {
    return (
      <Card className="h-full border-0 shadow-sm">
        <CardHeader className="px-3 sm:px-6">
          <Skeleton className="h-4 sm:h-6 w-36 sm:w-48" />
          <Skeleton className="h-3 sm:h-4 w-48 sm:w-64 mt-1 sm:mt-2" />
        </CardHeader>
        <CardContent className="px-3 sm:px-6 pb-3 sm:pb-6">
          {Array.from({ length: 5 }).map((_, i) => (
            <Skeleton key={i} className="h-12 sm:h-16 w-full mb-1.5 sm:mb-2" />
          ))}
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="h-full flex flex-col border-0 shadow-sm">
      <CardHeader className="px-3 sm:px-6">
        <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-2 sm:gap-0">
          <div className="flex-1 min-w-0">
            <CardTitle className="flex items-center gap-1.5 sm:gap-2 text-sm sm:text-base">
              <TrendingUp className="h-4 w-4 sm:h-5 sm:w-5" />
              {t("salesPerformance.title")}
            </CardTitle>
            <CardDescription className="mt-0.5 sm:mt-1 text-xs sm:text-sm">
              {t("salesPerformance.description")}
            </CardDescription>
          </div>
          <Link
            href="/sales-overview"
            className="text-xs sm:text-sm text-primary hover:underline cursor-pointer shrink-0"
          >
            {t("viewAll")}
          </Link>
        </div>
      </CardHeader>
      <CardContent className="flex-1 px-3 sm:px-6 pb-3 sm:pb-6">
        {performanceList.length === 0 ? (
          <div className="text-center py-6 sm:py-8 text-muted-foreground">
            <Users className="h-8 w-8 sm:h-12 sm:w-12 mx-auto mb-2 sm:mb-3 opacity-50" />
            <p className="text-xs sm:text-sm">{t("salesPerformance.noData")}</p>
          </div>
        ) : (
          <div className="max-h-full overflow-y-auto space-y-2 sm:space-y-3 pr-1 sm:pr-2">
            {performanceList.map((perf: SalesPerformanceListItem) => (
              <Link
                key={perf.user_id}
                href={`/sales-overview/sales-rep/${perf.user_id}`}
                className="flex items-center gap-2 sm:gap-3 font-medium text-primary hover:underline cursor-pointer p-1.5 sm:p-2 rounded-md hover:bg-muted/50 transition-colors"
              >
                <Avatar className="h-7 w-7 sm:h-8 sm:w-8 shrink-0">
                  <AvatarImage
                    src={perf.avatar_url}
                    alt={perf.user_name ?? "User"}
                  />
                </Avatar>
                <div className="flex-1 min-w-0">
                  <div className="truncate text-xs sm:text-sm">{perf.user_name ?? "-"}</div>
                  <div className="text-[10px] sm:text-xs text-muted-foreground truncate">
                    {perf.user_email}
                  </div>
                </div>
                <div className="text-right shrink-0 ml-2">
                  <div className="text-xs sm:text-sm font-medium">{perf.total_revenue_formatted}</div>
                  <div className="text-[10px] sm:text-xs text-muted-foreground">
                    {perf.deals_closed} {t("salesPerformance.deals")}
                  </div>
                </div>
              </Link>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}

