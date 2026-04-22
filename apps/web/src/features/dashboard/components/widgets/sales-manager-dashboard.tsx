"use client";

import { Suspense } from "react";
import { Skeleton } from "@/components/ui/skeleton";
import { VisitStatistics } from "../visit-statistics";
import { TopSalesRep } from "../top-sales-rep";
import { TeamDraftApprovals } from "../team-draft-approvals";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { TrendingUp, Target, AlertTriangle, CheckCircle2 } from "lucide-react";
import { cn } from "@/lib/utils";
import {
  useSalesManagerPipelineFunnel,
  useSalesManagerTargetVsActual,
  useSalesManagerVisitCompletion,
  useSalesManagerDealsAtRisk,
} from "../../hooks/useDashboard";

export function SalesManagerDashboard({ startDate, endDate }: Readonly<{ startDate?: string; endDate?: string }>) {
  return (
    <div className="space-y-4 sm:space-y-6 md:space-y-8">
      {/* ========================================================================
          SECTION 1: PERFORMANCE MONITORING
          Pipeline Funnel, Target vs Actual, Visit Completion, and Deals at Risk
      ========================================================================= */}
      <section className="space-y-3 sm:space-y-4 md:space-y-6">
        <div className="flex items-center gap-1.5 sm:gap-2">
          <div className="h-1 w-1 rounded-full bg-primary" />
          <h2 className="text-sm sm:text-base md:text-lg font-medium text-muted-foreground uppercase tracking-wide">
            Performance Monitoring
          </h2>
        </div>

        <div className="grid gap-3 sm:gap-4 md:gap-6 grid-cols-2 md:grid-cols-2 lg:grid-cols-4">
        <Suspense fallback={<Skeleton className="h-24 sm:h-32" />}>
          <SalesManagerPipelineFunnelCard startDate={startDate} endDate={endDate} />
        </Suspense>

        <Suspense fallback={<Skeleton className="h-24 sm:h-32" />}>
          <SalesManagerTargetVsActualCard startDate={startDate} endDate={endDate} />
        </Suspense>

        <Suspense fallback={<Skeleton className="h-24 sm:h-32" />}>
          <SalesManagerVisitCompletionCard startDate={startDate} endDate={endDate} />
        </Suspense>

        <Suspense fallback={<Skeleton className="h-24 sm:h-32" />}>
          <SalesManagerDealsAtRiskCard startDate={startDate} endDate={endDate} />
        </Suspense>
        </div>
      </section>

      {/* ========================================================================
          SECTION 2: TEAM OVERVIEW
          Team Performance and Team Pending Approvals side-by-side
      ========================================================================= */}
      <section className="space-y-3 sm:space-y-4 md:space-y-6">
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 items-stretch h-full">
          <div className="space-y-3 sm:space-y-4 md:space-y-6 h-full flex flex-col">
            <div className="flex items-center gap-1.5 sm:gap-2">
              <div className="h-1 w-1 rounded-full bg-primary" />
              <h2 className="text-sm sm:text-base md:text-lg font-medium text-muted-foreground uppercase tracking-wide">
                Team Performance
              </h2>
            </div>

            <Suspense fallback={<Skeleton className="h-64 sm:h-80 w-full" />}>
              <TopSalesRep startDate={startDate} endDate={endDate} />
            </Suspense>
          </div>

          <div className="space-y-3 sm:space-y-4 md:space-y-6 h-full flex flex-col">
            <div className="flex items-center gap-1.5 sm:gap-2">
              <div className="h-1 w-1 rounded-full bg-primary" />
              <h2 className="text-sm sm:text-base md:text-lg font-medium text-muted-foreground uppercase tracking-wide">
                Team Pending Approvals
              </h2>
            </div>

            <Suspense fallback={<Skeleton className="h-64 sm:h-80 w-full" />}>
              <TeamDraftApprovals />
            </Suspense>
          </div>
        </div>
      </section>

      {/* ========================================================================
          SECTION 4: VISIT STATISTICS
          Visit Reports and Statistics
      ========================================================================= */}
      <section className="space-y-3 sm:space-y-4 md:space-y-6">
        <div className="flex items-center gap-1.5 sm:gap-2">
          <div className="h-1 w-1 rounded-full bg-primary" />
          <h2 className="text-sm sm:text-base md:text-lg font-medium text-muted-foreground uppercase tracking-wide">
            Visit Statistics
          </h2>
        </div>

        {/* Visit Statistics */}
        <Suspense fallback={<Skeleton className="h-64 sm:h-80 w-full" />}>
          <VisitStatistics startDate={startDate} endDate={endDate} />
        </Suspense>
      </section>
    </div>
  );
}

// ============================================================================
// Sales Manager Dashboard Widget Components
// ============================================================================

function SalesManagerPipelineFunnelCard({ startDate, endDate }: Readonly<{ startDate?: string; endDate?: string }>) {
  const params = { start_date: startDate, end_date: endDate };
  const { data, isLoading, isError } = useSalesManagerPipelineFunnel(params);

  if (isLoading) {
    return (
      <Card>
        <CardHeader className="pb-2 flex flex-row items-center justify-between space-y-0">
          <Skeleton className="h-3 sm:h-4 w-24 sm:w-28" />
          <TrendingUp className="h-3 w-3 sm:h-4 sm:w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent className="px-3 sm:px-6 pb-3 sm:pb-6">
          <Skeleton className="h-6 sm:h-8 w-12 sm:w-16 mb-1 sm:mb-2" />
          <Skeleton className="h-2.5 sm:h-3 w-24 sm:w-32" />
        </CardContent>
      </Card>
    );
  }

  if (isError) {
    return (
      <Card>
        <CardHeader className="pb-2 flex flex-row items-center justify-between space-y-0">
          <CardTitle className="text-sm font-medium text-muted-foreground">
            Pipeline Funnel
          </CardTitle>
          <TrendingUp className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-sm text-destructive">Failed to load</div>
        </CardContent>
      </Card>
    );
  }

  const funnel = data?.data;
  const conversionRate = funnel?.conversion_rate ?? 0;

  return (
    <Card className="border-0 shadow-sm">
      <CardHeader className="pb-1.5 sm:pb-2 flex flex-row items-center justify-between space-y-0 px-3 sm:px-6 pt-3 sm:pt-6">
        <CardTitle className="text-xs sm:text-sm font-medium text-muted-foreground">
          Pipeline Funnel
        </CardTitle>
        <TrendingUp className="h-3 w-3 sm:h-4 sm:w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent className="px-3 sm:px-6 pb-3 sm:pb-6">
        <div className="text-xl sm:text-2xl font-medium">{conversionRate.toFixed(1)}%</div>
        <p className="text-[10px] sm:text-xs text-muted-foreground mt-0.5 sm:mt-1">
          Conversion rate | Stages: {funnel?.funnel?.length ?? 0}
        </p>
      </CardContent>
    </Card>
  );
}

function SalesManagerTargetVsActualCard({ startDate, endDate }: Readonly<{ startDate?: string; endDate?: string }>) {
  const params = { period: "month" as const, start_date: startDate, end_date: endDate };
  const { data, isLoading, isError } = useSalesManagerTargetVsActual(params);

  if (isLoading) {
    return (
      <Card>
        <CardHeader className="pb-2 flex flex-row items-center justify-between space-y-0">
          <Skeleton className="h-4 w-32" />
          <Target className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent className="px-3 sm:px-6 pb-3 sm:pb-6">
          <Skeleton className="h-6 sm:h-8 w-12 sm:w-16 mb-1 sm:mb-2" />
          <Skeleton className="h-1.5 w-full my-2" />
          <Skeleton className="h-2.5 sm:h-3 w-24 sm:w-32" />
        </CardContent>
      </Card>
    );
  }

  if (isError) {
    return (
      <Card>
        <CardHeader className="pb-2 flex flex-row items-center justify-between space-y-0">
          <CardTitle className="text-sm font-medium text-muted-foreground">
            Target vs Actual
          </CardTitle>
          <Target className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-sm text-destructive">Failed to load</div>
        </CardContent>
      </Card>
    );
  }

  const targetVsActual = data?.data;
  const achievement = targetVsActual?.achievement;
  const target = targetVsActual?.target;
  const actual = targetVsActual?.actual;
  const revenuePercent = achievement?.revenue_percent ?? 0;
  const targetRevenue = target?.revenue ?? 0;
  const actualRevenue = actual?.revenue ?? 0;

  // Derive palette from CSS vars based on attainment thresholds
  const getAchievementTextColor = (pct: number) =>
    pct >= 100 ? "text-success" : pct >= 75 ? "text-primary" : "text-destructive";

  const getBarBg = (pct: number) =>
    pct >= 100 ? "bg-success" : pct >= 75 ? "bg-primary" : "bg-destructive";

  const formatRevenue = (cents: number) =>
    new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
      notation: "compact",
    }).format(cents / 100);

  // No target set — minimal state
  if (targetRevenue === 0) {
    return (
      <Card className="border-0 shadow-sm">
        <CardHeader className="pb-1.5 sm:pb-2 flex flex-row items-center justify-between space-y-0 px-3 sm:px-6 pt-3 sm:pt-6">
          <CardTitle className="text-xs sm:text-sm font-medium text-muted-foreground">
            Target vs Actual
          </CardTitle>
          <Target className="h-3 w-3 sm:h-4 sm:w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent className="px-3 sm:px-6 pb-3 sm:pb-6">
          <div className="text-base sm:text-lg font-medium text-muted-foreground">—</div>
          <p className="text-[10px] sm:text-xs text-muted-foreground mt-1">
            No target this month
          </p>
          {actualRevenue > 0 && (
            <p className="text-[10px] sm:text-xs text-muted-foreground mt-0.5">
              Actual: {formatRevenue(actualRevenue)}
            </p>
          )}
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="border-0 shadow-sm">
      <CardHeader className="pb-1.5 sm:pb-2 flex flex-row items-center justify-between space-y-0 px-3 sm:px-6 pt-3 sm:pt-6">
        <CardTitle className="text-xs sm:text-sm font-medium text-muted-foreground">
          Target vs Actual
        </CardTitle>
        <Target className={cn("h-3 w-3 sm:h-4 sm:w-4", getAchievementTextColor(revenuePercent))} />
      </CardHeader>
      <CardContent className="px-3 sm:px-6 pb-3 sm:pb-6">
        <div className={cn("text-xl sm:text-2xl font-medium", getAchievementTextColor(revenuePercent))}>
          {revenuePercent.toFixed(1)}%
        </div>
        {/* Attainment progress bar */}
        <div className="mt-1.5 h-1.5 w-full rounded-full bg-muted overflow-hidden">
          <div
            className={cn("h-full rounded-full transition-all duration-500", getBarBg(revenuePercent))}
            style={{ width: `${Math.min(revenuePercent, 100)}%` }}
          />
        </div>
        <p className="text-[10px] sm:text-xs text-muted-foreground mt-1.5">
          {formatRevenue(actualRevenue)}
          <span className="mx-1 opacity-50">/</span>
          {formatRevenue(targetRevenue)}
        </p>
      </CardContent>
    </Card>
  );
}

function SalesManagerVisitCompletionCard({ startDate, endDate }: Readonly<{ startDate?: string; endDate?: string }>) {
  const params = { period: "month" as const, start_date: startDate, end_date: endDate };
  const { data, isLoading, isError } = useSalesManagerVisitCompletion(params);

  if (isLoading) {
    return (
      <Card>
        <CardHeader className="pb-2 flex flex-row items-center justify-between space-y-0">
          <Skeleton className="h-4 w-32" />
          <CheckCircle2 className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <Skeleton className="h-8 w-16 mb-2" />
          <Skeleton className="h-3 w-32" />
        </CardContent>
      </Card>
    );
  }

  if (isError) {
    return (
      <Card>
        <CardHeader className="pb-2 flex flex-row items-center justify-between space-y-0">
          <CardTitle className="text-sm font-medium text-muted-foreground">
            Visit Completion
          </CardTitle>
          <CheckCircle2 className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-sm text-destructive">Failed to load</div>
        </CardContent>
      </Card>
    );
  }

  const completion = data?.data;
  const completionRate = completion?.completion_rate ?? 0;
  const completed = completion?.completed ?? 0;
  const totalScheduled = completion?.total_scheduled ?? 0;

  return (
    <Card className="border-0 shadow-sm">
      <CardHeader className="pb-1.5 sm:pb-2 flex flex-row items-center justify-between space-y-0 px-3 sm:px-6 pt-3 sm:pt-6">
        <CardTitle className="text-xs sm:text-sm font-medium text-muted-foreground">
          Visit Completion
        </CardTitle>
        <CheckCircle2 className="h-4 w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-medium">{completionRate.toFixed(1)}%</div>
        <p className="text-[10px] sm:text-xs text-muted-foreground mt-0.5 sm:mt-1">
          {completed}/{totalScheduled} completed | Pending: {completion?.pending ?? 0}
        </p>
      </CardContent>
    </Card>
  );
}

function SalesManagerDealsAtRiskCard({ startDate, endDate }: Readonly<{ startDate?: string; endDate?: string }>) {
  const params = { limit: 10, start_date: startDate, end_date: endDate };
  const { data, isLoading, isError } = useSalesManagerDealsAtRisk(params);

  if (isLoading) {
    return (
      <Card>
        <CardHeader className="pb-2 flex flex-row items-center justify-between space-y-0">
          <Skeleton className="h-4 w-28" />
          <AlertTriangle className="h-4 w-4 text-destructive" />
        </CardHeader>
        <CardContent>
          <Skeleton className="h-8 w-16 mb-2" />
          <Skeleton className="h-3 w-32" />
        </CardContent>
      </Card>
    );
  }

  if (isError) {
    return (
      <Card>
        <CardHeader className="pb-2 flex flex-row items-center justify-between space-y-0">
          <CardTitle className="text-sm font-medium text-muted-foreground">
            Deals at Risk
          </CardTitle>
          <AlertTriangle className="h-4 w-4 text-destructive" />
        </CardHeader>
        <CardContent>
          <div className="text-sm text-destructive">Failed to load</div>
        </CardContent>
      </Card>
    );
  }

  const dealsAtRisk = data?.data;
  const totalAtRisk = dealsAtRisk?.total_at_risk ?? 0;

  return (
    <Card className="border-0 shadow-sm">
      <CardHeader className="pb-1.5 sm:pb-2 flex flex-row items-center justify-between space-y-0 px-3 sm:px-6 pt-3 sm:pt-6">
        <CardTitle className="text-xs sm:text-sm font-medium text-muted-foreground">
          Deals at Risk
        </CardTitle>
        <AlertTriangle className="h-4 w-4 text-destructive" />
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-medium">{totalAtRisk}</div>
        <p className="text-[10px] sm:text-xs text-muted-foreground mt-0.5 sm:mt-1">
          Deals requiring attention
        </p>
      </CardContent>
    </Card>
  );
}

