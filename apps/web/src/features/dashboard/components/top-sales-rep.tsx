"use client";

import { useTranslations } from "next-intl";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useTopSalesRep } from "../hooks/useDashboard";
import { Skeleton } from "@/components/ui/skeleton";
import { Users, MapPin, Building2, Activity, Target, TrendingUp } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarImage } from "@/components/ui/avatar";
import { cn } from "@/lib/utils";

import { useHasPermission } from "@/features/auth/providers/permissions-provider";
import { Lock } from "lucide-react";

export function TopSalesRep({ startDate, endDate }: Readonly<{ startDate?: string; endDate?: string }>) {
  const t = useTranslations("topSalesReps");
  
  // Permission check
  const hasPermission = useHasPermission("dashboard.view");

  const params = { limit: 5, start_date: startDate, end_date: endDate };

  const { data, isLoading } = useTopSalesRep(params, { enabled: hasPermission });

  // Show warning if no permission
  if (!hasPermission) {
    return (
      <Card className="border-0 shadow-sm bg-muted/10 h-full flex flex-col">
        <CardHeader className="px-4 sm:px-6 pb-2">
          <CardTitle className="text-sm sm:text-base flex items-center gap-2 text-muted-foreground">
            <Lock className="h-4 w-4" />
            {t("title")}
          </CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col items-center justify-center p-6 text-center h-[250px]">
          <div className="bg-muted/20 p-3 rounded-full mb-3">
            <Lock className="h-6 w-6 text-muted-foreground" />
          </div>
          <p className="text-sm font-medium text-muted-foreground mb-1">
            Access Restricted
          </p>
          <p className="text-xs text-muted-foreground/70 max-w-[200px]">
            You don&apos;t have permission to view top sales reps.
          </p>
        </CardContent>
      </Card>
    );
  }

  if (isLoading) {
    return (
      <Card className="border-0 shadow-sm h-full flex flex-col">
        <CardHeader className="px-4 sm:px-6">
          <CardTitle className="text-sm sm:text-base">{t("title")}</CardTitle>
        </CardHeader>
        <CardContent className="px-4 sm:px-6 pb-4 sm:pb-6 flex-1 overflow-y-auto max-h-[520px]">
          <div className="space-y-2 sm:space-y-3">
            {[...Array(5)].map((_, i) => (
              <Skeleton key={i} className="h-20 sm:h-24 w-full" />
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  const topSalesReps = data?.data || [];

  if (topSalesReps.length === 0) {
    return (
      <Card className="border-0 shadow-sm h-full flex flex-col">
        <CardHeader className="px-4 sm:px-6">
          <div className="flex items-center gap-1.5 sm:gap-2">
            <Users className="h-4 w-4 sm:h-5 sm:w-5" />
            <CardTitle className="text-sm sm:text-base">{t("title")}</CardTitle>
          </div>
        </CardHeader>
        <CardContent className="relative min-h-[150px] sm:min-h-[200px] px-4 sm:px-6 pb-4 sm:pb-6 flex-1 overflow-y-auto max-h-[520px]">
          <div className="absolute inset-0 flex flex-col items-center justify-center">
            <Users className="h-8 w-8 sm:h-12 sm:w-12 text-muted-foreground/50 mb-2 sm:mb-3" />
            <p className="text-xs sm:text-sm text-muted-foreground">{t("empty")}</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  const getAvatarUrl = (email: string) => {
    return `https://api.dicebear.com/7.x/lorelei/svg?seed=${encodeURIComponent(email)}`;
  };

  const getAchievementColor = (pct: number) => {
    if (pct >= 100) return "bg-green-500";
    if (pct >= 75) return "bg-yellow-500";
    return "bg-red-500";
  };

  const getAchievementBadgeVariant = (pct: number): "default" | "secondary" | "destructive" | "outline" => {
    if (pct >= 100) return "default";
    if (pct >= 75) return "secondary";
    return "destructive";
  };

  return (
    <Card className="border-0 shadow-sm h-full flex flex-col">
      <CardHeader className="px-4 sm:px-6">
        <div className="flex items-center gap-1.5 sm:gap-2">
          <Users className="h-4 w-4 sm:h-5 sm:w-5" />
          <CardTitle className="text-sm sm:text-base">{t("title")}</CardTitle>
        </div>
      </CardHeader>
      <CardContent className="px-4 sm:px-6 pb-4 sm:pb-6 flex-1 overflow-y-auto max-h-[520px]">
        <div className="space-y-3 sm:space-y-4">
          {topSalesReps.map((salesRep, index) => {
            const salesRepData = salesRep.sales_rep ?? { id: "", name: "Unknown User", email: "" };
            const visitCount = salesRep.visit_count ?? 0;
            const accountCount = salesRep.account_count ?? 0;
            const activityCount = salesRep.activity_count ?? 0;
            const dealsClosed = salesRep.deals_closed ?? 0;
            const actualRevenue = salesRep.actual_revenue_formatted ?? "Rp 0";
            const targetAmount = salesRep.target_amount_formatted ?? "Rp 0";
            const achievementPct = salesRep.target_achievement_percent ?? 0;
            const hasTarget = (salesRep.target_amount ?? 0) > 0;
            
            return (
                <div
                  key={salesRepData.id || `sales-rep-${index}`}
                  className="p-2.5 sm:p-3 rounded-md sm:rounded-lg border hover:bg-muted/50 transition-colors cursor-pointer space-y-2.5"
                >
                {/* Header: Avatar + Name + Achievement badge */}
                <div className="flex items-center justify-between gap-2">
                  <div className="flex items-center gap-2 sm:gap-3 flex-1 min-w-0">
                    <Avatar className="h-8 w-8 sm:h-9 sm:w-9 shrink-0">
                      <AvatarImage
                        src={getAvatarUrl(salesRepData.email || "unknown")}
                        alt={salesRepData.name || "Unknown User"}
                      />
                    </Avatar>
                    <div className="min-w-0">
                      <div className="font-medium text-xs sm:text-sm truncate">{salesRepData.name || "Unknown User"}</div>
                      {/* Activity stats row */}
                      <div className="flex flex-wrap items-center gap-2 sm:gap-3 mt-0.5 text-[10px] sm:text-xs text-muted-foreground">
                        <div className="flex items-center gap-0.5 sm:gap-1">
                          <MapPin className="h-2.5 w-2.5 sm:h-3 sm:w-3" />
                          <span>{t("visits", { count: visitCount })}</span>
                        </div>
                        <div className="flex items-center gap-0.5 sm:gap-1">
                          <Building2 className="h-2.5 w-2.5 sm:h-3 sm:w-3" />
                          <span>{t("accounts", { count: accountCount })}</span>
                        </div>
                        <div className="flex items-center gap-0.5 sm:gap-1">
                          <Activity className="h-2.5 w-2.5 sm:h-3 sm:w-3" />
                          <span>{t("activities", { count: activityCount })}</span>
                        </div>
                      </div>
                    </div>
                  </div>
                  {hasTarget && (
                    <Badge
                      variant={getAchievementBadgeVariant(achievementPct)}
                      className="text-[10px] sm:text-xs shrink-0 font-medium"
                    >
                      {Math.round(achievementPct)}%
                    </Badge>
                  )}
                </div>

                {/* Achievement section */}
                <div className="space-y-1.5">
                  {/* Revenue + Deals row */}
                  <div className="flex items-center justify-between text-[10px] sm:text-xs">
                    <div className="flex items-center gap-1 text-muted-foreground">
                      <TrendingUp className="h-2.5 w-2.5 sm:h-3 sm:w-3" />
                      <span className="font-medium text-foreground">{actualRevenue}</span>
                      {hasTarget && (
                        <span className="text-muted-foreground">/ {targetAmount}</span>
                      )}
                    </div>
                    <div className="flex items-center gap-1 text-muted-foreground">
                      <Target className="h-2.5 w-2.5 sm:h-3 sm:w-3" />
                      <span>{dealsClosed} deals</span>
                    </div>
                  </div>

                  {/* Progress bar — shown only if target is set */}
                  {hasTarget && (
                    <div className="relative h-1.5 w-full overflow-hidden rounded-full bg-muted">
                      <div
                        className={cn(
                          "h-full transition-all rounded-full",
                          getAchievementColor(achievementPct)
                        )}
                        style={{ width: `${Math.min(achievementPct, 100)}%` }}
                      />
                    </div>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      </CardContent>
    </Card>
  );
}
