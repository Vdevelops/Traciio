"use client";

import { useTranslations } from "next-intl";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { useBrickPerformance } from "../hooks/useBrickAnalytics";
import { cn } from "@/lib/utils";
import { TrendingDown, TrendingUp } from "lucide-react";

interface BrickPipelineAnalyticsProps {
  readonly brickId: string;
  readonly periodStart?: string;
  readonly periodEnd?: string;
}

const DEAL_STAGES = [
  { key: "open",   label: "Open",        color: "bg-blue-500",   dot: "bg-blue-500" },
  { key: "won",    label: "Closed Won",  color: "bg-green-500",  dot: "bg-green-500" },
  { key: "lost",   label: "Closed Lost", color: "bg-red-500",    dot: "bg-red-500" },
];

export function BrickPipelineAnalytics({ brickId, periodStart, periodEnd }: BrickPipelineAnalyticsProps) {
  const t = useTranslations("brickAnalytics.pipeline");

  const { data, isLoading } = useBrickPerformance(brickId, periodStart, periodEnd);

  if (isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-28 w-full" />
        <Skeleton className="h-40 w-full" />
        <Skeleton className="h-28 w-full" />
      </div>
    );
  }

  const metrics = data?.data;

  if (!metrics) {
    return (
      <Card>
        <CardContent className="py-8 text-center text-muted-foreground">
          {t("noData")}
        </CardContent>
      </Card>
    );
  }

  const fmt = (value: number) =>
    new Intl.NumberFormat("id-ID", { style: "currency", currency: "IDR", minimumFractionDigits: 0 }).format(value);

  const winRate = metrics.total_deals > 0 ? (metrics.won_deals / metrics.total_deals) * 100 : 0;
  const lostRate = metrics.total_deals > 0 ? (metrics.lost_deals / metrics.total_deals) * 100 : 0;
  const openRate = metrics.total_deals > 0 ? (metrics.open_deals / metrics.total_deals) * 100 : 0;

  const stages = [
    { ...DEAL_STAGES[0], count: metrics.open_deals,  value: metrics.total_deal_value - metrics.won_deal_value, pct: openRate },
    { ...DEAL_STAGES[1], count: metrics.won_deals,   value: metrics.won_deal_value,  pct: winRate },
    { ...DEAL_STAGES[2], count: metrics.lost_deals,  value: 0,                       pct: lostRate },
  ];

  // Target / Achievement
  const tgt = metrics.monthly_target ?? 0;
  const achieved = metrics.target_achieved ?? 0;
  const achievementPct = tgt > 0 ? (achieved / tgt) * 100 : 0;
  const isAchievementGood = achievementPct >= 80;

  return (
    <div className="space-y-4">
      {/* Pipeline key metrics */}
      <Card>
        <CardContent className="pt-5">
          <div className="grid grid-cols-2 md:grid-cols-4 divide-x divide-border">
            <PipelineStat label={t("totalPipelineValue")} value={fmt(metrics.total_deal_value)} sub={t("totalPipelineValueDesc")} />
            <PipelineStat label={t("wonDealValue")}       value={fmt(metrics.won_deal_value)}   sub={t("wonDealValueDesc")} accent="green" />
            <PipelineStat label={t("averageDealSize")}   value={fmt(metrics.average_deal_size)} sub={t("averageDealSizeDesc")} />
            <PipelineStat label={t("winRate")}           value={`${winRate.toFixed(1)}%`}        sub={t("winRateDesc", { won: metrics.won_deals, total: metrics.total_deals })} />
          </div>
        </CardContent>
      </Card>

      {/* Deal Status — single segmented bar */}
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-sm font-medium">{t("dealStatusBreakdown")}</CardTitle>
          <p className="text-xs text-muted-foreground">{t("totalDealsDesc")}</p>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* Segmented bar */}
          <div className="flex h-3 w-full overflow-hidden rounded-full bg-muted">
            {metrics.total_deals > 0 ? (
              stages
                .filter((s) => s.pct > 0)
                .map((s) => (
                  <div
                    key={s.key}
                    className={cn("h-full transition-all", s.color)}
                    style={{ width: `${s.pct}%` }}
                    title={`${s.label}: ${s.pct.toFixed(1)}%`}
                  />
                ))
            ) : (
              <div className="h-full w-full bg-muted-foreground/20 rounded-full" />
            )}
          </div>

          {/* Legend */}
          <div className="space-y-2">
            {stages.map((s) => (
              <div key={s.key} className="flex items-center justify-between">
                <div className="flex items-center gap-2">
                  <span className={cn("h-2.5 w-2.5 rounded-full shrink-0", s.dot)} />
                  <span className="text-sm">{s.label}</span>
                  <span className="text-sm text-muted-foreground">
                    · {s.count} {s.count === 1 ? "deal" : "deals"}
                    {s.value > 0 && ` · ${fmt(s.value)}`}
                  </span>
                </div>
                <div className="flex items-center gap-3">
                  <div className="w-24 h-1.5 rounded-full bg-muted overflow-hidden">
                    <div
                      className={cn("h-full rounded-full", s.color)}
                      style={{ width: `${s.pct}%` }}
                    />
                  </div>
                  <span className="text-sm text-muted-foreground w-10 text-right">
                    {s.pct.toFixed(0)}%
                  </span>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

function PipelineStat({
  label,
  value,
  sub,
  accent,
}: Readonly<{ label: string; value: string; sub?: string; accent?: "green" }>) {
  return (
    <div className="flex flex-col gap-1 px-5 py-1 first:pl-0 last:pr-0">
      <span className="text-xs text-muted-foreground">{label}</span>
      <span className={cn("text-2xl font-semibold", accent === "green" && "text-green-600")}>{value}</span>
      {sub && <span className="text-xs text-muted-foreground">{sub}</span>}
    </div>
  );
}

