"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Progress } from "@/components/ui/progress";
import { Button } from "@/components/ui/button";
import { useBrickPerformance, useCurrentMonthPeriod } from "../hooks/useBrickAnalytics";
import { useBrickTargetWithDistributions } from "../hooks/useBricks";
import { Target, TrendingUp, Users, Calendar, Plus } from "lucide-react";
import { cn, formatCurrency } from "@/lib/utils";
import { useMemo } from "react";
import { BrickTargetDistributionDialog } from "./brick-target-distribution-dialog";

interface BrickTargetAnalyticsProps {
  brickId: string;
  periodStart?: string;
  periodEnd?: string;
}

export function BrickTargetAnalytics({ brickId, periodStart, periodEnd }: BrickTargetAnalyticsProps) {
  const t = useTranslations("brickAnalytics.target");
  const [isDistributionDialogOpen, setIsDistributionDialogOpen] = useState(false);
  const { periodStart: defaultStart, periodEnd: defaultEnd } = useCurrentMonthPeriod();
  const requestStart = periodStart || undefined;
  const requestEnd = periodEnd || undefined;

  const { data: performanceData, isLoading: isLoadingPerformance } = useBrickPerformance(
    brickId,
    requestStart,
    requestEnd
  );

  // Target distributions are monthly; use the month of the selected periodStart.
  const monthReferenceStart = periodStart || defaultStart;
  const periodStartDate = monthReferenceStart
    ? new Date(`${monthReferenceStart}T00:00:00`)
    : new Date();
  const currentYear = periodStartDate.getFullYear();
  const currentMonth = periodStartDate.getMonth() + 1;

  const { data: targetData, isLoading: isLoadingTarget } = useBrickTargetWithDistributions(
    brickId,
    currentYear,
    currentMonth
  );

  const isLoading = isLoadingPerformance || isLoadingTarget;

  if (isLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-64 w-full" />
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  const metrics = performanceData?.data;
  const targetDistributions = targetData?.data;

  if (!metrics) {
    return (
      <Card>
        <CardContent className="py-8 text-center text-muted-foreground">
          {t("noData")}
        </CardContent>
      </Card>
    );
  }

  const achievementPercentage = metrics.monthly_target > 0
    ? (metrics.target_achieved / metrics.monthly_target) * 100
    : 0;

  const totalDistributed = targetDistributions?.total_distributed ?? 0;
  const remainingAmount = targetDistributions?.remaining_amount ?? 0;
  const distributionPercentage = metrics.monthly_target > 0
    ? (totalDistributed / metrics.monthly_target) * 100
    : 0;

  const targetStats = [
    {
      label: t("monthlyTarget"),
      value: formatCurrency(metrics.monthly_target),
      icon: Target,
      description: t("monthlyTargetDesc"),
    },
    {
      label: t("targetAchieved"),
      value: formatCurrency(metrics.target_achieved),
      icon: TrendingUp,
      description: t("targetAchievedDesc"),
    },
    {
      label: t("achievementPercentage"),
      value: `${achievementPercentage.toFixed(1)}%`,
      icon: Calendar,
      description: t("achievementPercentageDesc"),
      isPositive: achievementPercentage >= 80,
    },
    {
      label: t("targetRemaining"),
      value: formatCurrency(metrics.target_remaining),
      icon: Users,
      description: t("targetRemainingDesc"),
    },
  ];

  const distributions = targetDistributions?.distributions ?? [];

  return (
    <div className="space-y-6">
      {/* Target Overview Cards */}
      <div className="grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-4">
        {targetStats.map((stat, index) => (
          <Card key={index}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">{stat.label}</CardTitle>
              <stat.icon 
                className={cn(
                  "h-4 w-4 text-muted-foreground",
                  stat.isPositive !== undefined && (
                    stat.isPositive ? "text-green-500" : "text-yellow-500"
                  )
                )} 
              />
            </CardHeader>
            <CardContent>
              <div className={cn(
                "text-2xl font-medium",
                stat.isPositive !== undefined && (
                  stat.isPositive ? "text-green-600" : "text-yellow-600"
                )
              )}>
                {stat.value}
              </div>
              <p className="text-xs text-muted-foreground mt-1">{stat.description}</p>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Achievement Progress */}
      <Card>
        <CardHeader>
          <CardTitle>{t("achievementProgress")}</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <span className="text-sm font-medium">{t("targetAchieved")}</span>
              <span className="text-sm text-muted-foreground">
                {formatCurrency(metrics.target_achieved)} / {formatCurrency(metrics.monthly_target)}
              </span>
            </div>
            <div className="relative h-3 w-full overflow-hidden rounded-full bg-primary/20">
              <div
                className={cn(
                  "h-full transition-all",
                  achievementPercentage >= 100 ? "bg-green-500" :
                  achievementPercentage >= 80 ? "bg-yellow-500" : "bg-red-500"
                )}
                style={{ width: `${Math.min(achievementPercentage, 100)}%` }}
              />
            </div>
            <div className="flex items-center justify-between text-xs text-muted-foreground">
              <span>{achievementPercentage.toFixed(1)}% {t("achieved")}</span>
              <span>{t("remaining")}: {formatCurrency(metrics.target_remaining)}</span>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Target Distribution */}
      {targetDistributions && (
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle>{t("targetDistribution")}</CardTitle>
              {remainingAmount > 0 && (
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setIsDistributionDialogOpen(true)}
                  className="cursor-pointer"
                >
                  <Plus className="h-4 w-4 mr-2" />
                  {t("distributeTarget")}
                </Button>
              )}
            </div>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-4 grid-cols-1 md:grid-cols-3">
              <div>
                <div className="text-2xl font-medium">{formatCurrency(totalDistributed)}</div>
                <p className="text-sm text-muted-foreground">{t("totalDistributed")}</p>
                <p className="text-xs text-muted-foreground mt-1">
                  {distributionPercentage.toFixed(1)}% {t("ofTarget")}
                </p>
              </div>
              <div>
                <div className="text-2xl font-medium text-yellow-600">
                  {formatCurrency(remainingAmount)}
                </div>
                <p className="text-sm text-muted-foreground">{t("remainingToDistribute")}</p>
              </div>
              <div>
                <div className="text-2xl font-medium text-green-600">
                  {distributions.length}
                </div>
                <p className="text-sm text-muted-foreground">{t("salesWithTarget")}</p>
              </div>
            </div>

            {distributions.length > 0 && (
              <div className="mt-6 space-y-2">
                <h4 className="text-sm font-medium">{t("distributionDetails")}</h4>
                <div className="space-y-2">
                  {distributions.map((dist) => {
                    const salesAchievement = dist.sales_user
                      ? (metrics.target_achieved / metrics.monthly_target) * (dist.distributed_amount / metrics.monthly_target) * 100
                      : 0;
                    
                    return (
                      <div key={dist.id} className="flex items-center justify-between p-3 border rounded-lg">
                        <div className="flex-1">
                          <p className="text-sm font-medium">
                            {dist.sales_user?.name ?? t("unknownSales")}
                          </p>
                          <p className="text-xs text-muted-foreground">
                            {formatCurrency(dist.distributed_amount)}
                          </p>
                        </div>
                        <div className="text-right">
                          <p className="text-sm font-medium">
                            {salesAchievement.toFixed(1)}%
                          </p>
                          <p className="text-xs text-muted-foreground">{t("achievement")}</p>
                        </div>
                      </div>
                    );
                  })}
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      )}

      {/* Target Distribution Dialog */}
      {targetDistributions && remainingAmount > 0 && (
        <BrickTargetDistributionDialog
          open={isDistributionDialogOpen}
          onOpenChange={setIsDistributionDialogOpen}
          brickId={brickId}
          targetData={targetDistributions}
          onSuccess={() => {
            // Refetch target data - hook will auto refetch due to query invalidation
          }}
        />
      )}
    </div>
  );
}
