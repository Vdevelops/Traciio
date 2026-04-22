"use client";

import { useMemo } from "react";
import { useParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { BrickPerformanceOverview } from "@/features/master-data/brick/components/brick-performance-overview";
import { BrickAnalyticsTabs } from "@/features/master-data/brick/components/brick-analytics-tabs";
import { BrickSalesManagement } from "@/features/master-data/brick/components/brick-sales-management";
import { Skeleton } from "@/components/ui/skeleton";
import { useBrick } from "@/features/master-data/brick/hooks/useBricks";
import { ArrowLeft } from "lucide-react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { BrickPeriodFilter } from "@/features/master-data/brick/components/brick-period-filter";
import { useBrickPeriodQueryParams } from "@/features/master-data/brick/hooks/useBrickPeriodQueryParams";

export function BrickDashboardClient() {
  const t = useTranslations("brickDashboard");
  const params = useParams();
  const brickId = params?.id as string;

  const { mode, periodStart, periodEnd } = useBrickPeriodQueryParams();

  const periodQueryString = useMemo(() => {
    const p = new URLSearchParams();
    p.set("period_mode", mode);
    p.set("period_start", periodStart);
    p.set("period_end", periodEnd);
    return p.toString();
  }, [mode, periodStart, periodEnd]);

  const { data: brickData, isLoading: isLoadingBrick } = useBrick(brickId);

  if (isLoadingBrick) {
    return <BrickDashboardSkeleton />;
  }

  const brick = brickData?.data;

  if (!brick) {
    return (
      <div className="flex flex-col items-center justify-center py-12">
        <p className="text-muted-foreground">{t("brickNotFound")}</p>
        <Link href="/master-data/bricks">
          <Button variant="outline" className="mt-4 cursor-pointer">
            <ArrowLeft className="mr-2 h-4 w-4" />
            {t("backToBricks")}
          </Button>
        </Link>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <div className="flex items-center gap-2 mb-2">
            <Link href={`/master-data/bricks?${periodQueryString}`}>
              <Button variant="ghost" size="sm" className="cursor-pointer">
                <ArrowLeft className="h-4 w-4 mr-2" />
                {t("back")}
              </Button>
            </Link>
          </div>
          <h1 className="text-3xl font-medium tracking-tight">{brick.name}</h1>
          <p className="text-muted-foreground mt-1">
            {t("description")} • {brick.code} • {brick.province}, {brick.regency}
          </p>
          {brick.manager && (
            <p className="text-sm text-muted-foreground mt-1">
              {t("manager")}: {brick.manager.name}
            </p>
          )}
        </div>

        <BrickPeriodFilter />
      </div>

      {/* KPI Overview */}
      <div className="space-y-4">
        <SectionLabel label={t("performanceOverview")} />
        <BrickPerformanceOverview brickId={brickId} periodStart={periodStart} periodEnd={periodEnd} />
      </div>

      {/* Analytics — directly visible, no collapse */}
      <div className="space-y-4">
        <SectionLabel label={t("analytics")} />
        <BrickAnalyticsTabs brickId={brickId} periodStart={periodStart} periodEnd={periodEnd} />
      </div>

      {/* Sales Team */}
      <div className="space-y-4">
        <SectionLabel label={t("salesTeam")} />
        <BrickSalesManagement brickId={brickId} periodStart={periodStart} periodEnd={periodEnd} />
      </div>
    </div>
  );
}

function SectionLabel({ label }: Readonly<{ label: string }>) {
  return (
    <div className="flex items-center gap-1.5">
      <div className="h-1 w-1 rounded-full bg-primary" />
      <h2 className="text-sm font-medium text-muted-foreground uppercase tracking-wide">{label}</h2>
    </div>
  );
}

function BrickDashboardSkeleton() {
  return (
    <div className="space-y-8">
      <Skeleton className="h-8 w-64" />
      <div className="grid gap-4 grid-cols-2 md:grid-cols-4">
        {Array.from({ length: 4 }).map((_, i) => (
          <Skeleton key={i} className="h-32" />
        ))}
      </div>
      <Skeleton className="h-64 w-full" />
    </div>
  );
}

