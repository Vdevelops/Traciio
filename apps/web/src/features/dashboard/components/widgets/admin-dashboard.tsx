"use client";

import { Suspense } from "react";
import { useTranslations } from "next-intl";
import { StaggerContainer } from "@/components/motion";
import { Skeleton } from "@/components/ui/skeleton";
import { DashboardStatsPlain } from "../dashboard-stats-plain";
import { RevenueHighlight } from "../revenue-highlight";
import { LeadsBySource } from "../leads-by-source";
import { PipelineSummary } from "../pipeline-summary";
import { VisitStatistics } from "../visit-statistics";

export function AdminDashboard({ startDate, endDate }: Readonly<{ startDate?: string; endDate?: string }>) {
  const tOverview = useTranslations("dashboardOverview");

  return (
    <div className="space-y-4 sm:space-y-6">
      {/* ========================================================================
          SECTION 1: OVERVIEW METRICS
          Top-level KPIs and revenue summary
      ========================================================================= */}
      <section className="space-y-3">
        <div className="flex items-center gap-1.5 sm:gap-2">
          <div className="h-1.5 w-1.5 rounded-full bg-primary shadow-[0_0_8px_rgba(var(--primary),0.5)]" />
          <h2 className="text-sm sm:text-base md:text-lg font-bold text-foreground/80 lowercase tracking-tight first-letter:uppercase">
            {tOverview("sections.overview")}
          </h2>
        </div>

        {/* Main Metrics Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 sm:gap-6">
          {/* Revenue Highlight - Prominent with primary background */}
          <Suspense fallback={<Skeleton className="h-[200px] sm:h-[300px] w-full rounded-xl shadow-xs" />}>
            <RevenueHighlight startDate={startDate} endDate={endDate} />
          </Suspense>

          {/* Stats Plain - 4 sub-stats in a 2x2 grid inside */}
          <Suspense
            fallback={
              <div className="grid grid-cols-2 gap-4 sm:gap-6">
                {Array.from({ length: 4 }).map((_, i) => (
                  <Skeleton key={`admin-stat-skeleton-${i}`} className="h-32 w-full rounded-xl shadow-xs" />
                ))}
              </div>
            }
          >
            <DashboardStatsPlain startDate={startDate} endDate={endDate} />
          </Suspense>
        </div>
      </section>

      {/* ========================================================================
          SECTION 2: ANALYTICS & INSIGHTS
          Charts and detailed breakdowns
      ========================================================================= */}
      <section className="space-y-3">
        <div className="flex items-center gap-1.5 sm:gap-2">
          <div className="h-1.5 w-1.5 rounded-full bg-primary shadow-[0_0_8px_rgba(var(--primary),0.5)]" />
          <h2 className="text-sm sm:text-base md:text-lg font-bold text-foreground/80 lowercase tracking-tight first-letter:uppercase">
            {tOverview("sections.analyticsInsights")}
          </h2>
        </div>

        <StaggerContainer className="grid gap-4 sm:gap-6 grid-cols-1 lg:grid-cols-2">
          <Suspense fallback={<Skeleton className="h-64 sm:h-80 w-full rounded-xl shadow-xs" />}>
            <LeadsBySource startDate={startDate} endDate={endDate} />
          </Suspense>
          <Suspense fallback={<Skeleton className="h-64 sm:h-80 w-full rounded-xl shadow-xs" />}>
            <PipelineSummary startDate={startDate} endDate={endDate} />
          </Suspense>
        </StaggerContainer>

        {/* Visit Reports Summary */}
        <Suspense fallback={<Skeleton className="h-64 sm:h-80 w-full rounded-xl shadow-xs" />}>
          <VisitStatistics startDate={startDate} endDate={endDate} />
        </Suspense>
      </section>

    </div>
  );
}
