"use client";

import { useTranslations } from "next-intl";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { useBrickPerformance } from "../hooks/useBrickAnalytics";

interface BrickVisitAnalyticsProps {
  readonly brickId: string;
  readonly periodStart?: string;
  readonly periodEnd?: string;
}

export function BrickVisitAnalytics({ brickId, periodStart, periodEnd }: BrickVisitAnalyticsProps) {
  const t = useTranslations("brickAnalytics.visits");
  const { data, isLoading } = useBrickPerformance(brickId, periodStart, periodEnd);

  if (isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-28 w-full" />
        <Skeleton className="h-20 w-full" />
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

  return (
    <div className="space-y-4">
      {/* Minimalist visit stats row */}
      <Card>
        <CardContent className="pt-5">
          <div className="grid grid-cols-2 md:grid-cols-4 divide-x divide-border">
            <StatItem
              label={t("totalVisits")}
              value={metrics.total_visits}
              sub={t("totalVisitsDesc")}
            />
            <StatItem
              label={t("visitsThisMonth")}
              value={metrics.visits_this_month}
              sub={t("visitsThisMonthDesc")}
              highlight
            />
            <StatItem
              label={t("averageVisitsPerSales")}
              value={metrics.average_visits_per_sales.toFixed(1)}
              sub={t("averageVisitsPerSalesDesc", { sales: metrics.active_sales })}
            />
            <StatItem
              label={t("totalAccounts")}
              value={metrics.total_accounts}
              sub={t("totalAccountsDesc", { active: metrics.active_accounts })}
            />
          </div>
        </CardContent>
      </Card>

      {/* Account activity */}
      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-sm font-medium text-muted-foreground uppercase tracking-wide">
            {t("accountActivity")}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-3 divide-x divide-border">
            <StatItem label={t("totalAccounts")} value={metrics.total_accounts} />
            <StatItem label={t("activeAccounts")} value={metrics.active_accounts} accent="green" />
            <StatItem label={t("newAccountsThisMonth")} value={metrics.new_accounts_this_month} accent="primary" />
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

function resolveValueClass(highlight?: boolean, accent?: "green" | "primary") {
  if (accent === "green") return "text-3xl font-semibold text-green-600";
  if (accent === "primary") return "text-3xl font-semibold text-primary";
  if (highlight) return "text-3xl font-semibold text-primary";
  return "text-3xl font-semibold";
}

function StatItem({
  label,
  value,
  sub,
  highlight,
  accent,
}: Readonly<{
  label: string;
  value: string | number;
  sub?: string;
  highlight?: boolean;
  accent?: "green" | "primary";
}>) {
  const valueClass = resolveValueClass(highlight, accent);

  return (
    <div className="flex flex-col gap-1 px-5 py-1 first:pl-0 last:pr-0">
      <span className="text-xs text-muted-foreground">{label}</span>
      <span className={valueClass}>{value}</span>
      {sub && <span className="text-xs text-muted-foreground">{sub}</span>}
    </div>
  );
}

