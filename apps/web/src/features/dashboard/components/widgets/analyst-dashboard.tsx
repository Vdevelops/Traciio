"use client";

import { Suspense } from "react";
import { useTranslations } from "next-intl";
import { Line, LineChart, CartesianGrid, XAxis, YAxis } from "recharts";
import { Skeleton } from "@/components/ui/skeleton";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import {
  useAnalystAIInsights,
  useAnalystConversionRate,
  useAnalystRevenueTrend,
  useAnalystSalesVelocity,
} from "../../hooks/useDashboard";

const revenueChartConfig = {
  revenue: {
    label: "Revenue",
    color: "oklch(0.4851 0.1559 30.2175)",
  },
} satisfies ChartConfig;

const conversionChartConfig = {
  conversion_rate: {
    label: "Conversion Rate",
    color: "oklch(0.5234 0.1347 144.1672)",
  },
} satisfies ChartConfig;

export function AnalystDashboard({ startDate, endDate }: Readonly<{ startDate?: string; endDate?: string }>) {
  const t = useTranslations("dashboardAnalyst");

  return (
    <div className="space-y-4 sm:space-y-6 md:space-y-8">
      <section className="space-y-3 sm:space-y-4 md:space-y-6">
        <SectionLabel label={t("sections.revenue")} />
        <Suspense fallback={<Skeleton className="h-80 w-full" />}>
          <AnalystRevenueTrendCard startDate={startDate} endDate={endDate} />
        </Suspense>
      </section>

      <section className="space-y-3 sm:space-y-4 md:space-y-6">
        <SectionLabel label={t("sections.conversion")} />
        <div className="grid gap-4 lg:grid-cols-2">
          <Suspense fallback={<Skeleton className="h-80 w-full" />}>
            <AnalystConversionCard startDate={startDate} endDate={endDate} />
          </Suspense>
          <Suspense fallback={<Skeleton className="h-80 w-full" />}>
            <AnalystVelocityCard startDate={startDate} endDate={endDate} />
          </Suspense>
        </div>
      </section>

      <section className="space-y-3 sm:space-y-4 md:space-y-6">
        <SectionLabel label={t("sections.insights")} />
        <Suspense fallback={<Skeleton className="h-64 w-full" />}>
          <AnalystInsightsCard startDate={startDate} endDate={endDate} />
        </Suspense>
      </section>
    </div>
  );
}

function AnalystRevenueTrendCard({ startDate, endDate }: Readonly<{ startDate?: string; endDate?: string }>) {
  const t = useTranslations("dashboardAnalyst");
  const { data, isLoading } = useAnalystRevenueTrend({ start_date: startDate, end_date: endDate });

  if (isLoading) {
    return <Skeleton className="h-80 w-full" />;
  }

  const trend = data?.data;
  const chartData = trend?.trend ?? [];

  return (
    <Card className="border-0 shadow-sm">
      <CardHeader>
        <CardTitle>{t("revenue.title")}</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-4 md:grid-cols-3">
          <MetricCard label={t("revenue.totalRevenue")} value={formatCurrency(trend?.total_revenue ?? 0)} />
          <MetricCard label={t("revenue.growth")} value={`${(trend?.growth_percent ?? 0).toFixed(1)}%`} />
          <MetricCard label={t("revenue.averageDaily")} value={formatCurrency(trend?.average_daily ?? 0)} />
        </div>

        <ChartContainer config={revenueChartConfig} className="h-[280px] w-full">
          <LineChart data={chartData}>
            <CartesianGrid vertical={false} />
            <XAxis dataKey="date" tickLine={false} axisLine={false} tickMargin={8} />
            <YAxis tickLine={false} axisLine={false} tickMargin={8} />
            <ChartTooltip content={<ChartTooltipContent />} />
            <Line
              type="monotone"
              dataKey="revenue"
              stroke="var(--color-revenue)"
              strokeWidth={2}
              dot={{ r: 3 }}
            />
          </LineChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}

function AnalystConversionCard({ startDate, endDate }: Readonly<{ startDate?: string; endDate?: string }>) {
  const t = useTranslations("dashboardAnalyst");
  const { data, isLoading } = useAnalystConversionRate({ start_date: startDate, end_date: endDate });

  if (isLoading) {
    return <Skeleton className="h-80 w-full" />;
  }

  const conversion = data?.data;
  const chartData = conversion?.trend ?? [];

  return (
    <Card className="border-0 shadow-sm">
      <CardHeader>
        <CardTitle>{t("conversion.title")}</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-4 md:grid-cols-3">
          <MetricCard label={t("conversion.totalLeads")} value={(conversion?.total_leads ?? 0).toLocaleString("id-ID")} />
          <MetricCard label={t("conversion.convertedLeads")} value={(conversion?.converted_leads ?? 0).toLocaleString("id-ID")} />
          <MetricCard label={t("conversion.rate")} value={`${(conversion?.conversion_rate ?? 0).toFixed(1)}%`} />
        </div>

        <ChartContainer config={conversionChartConfig} className="h-[220px] w-full">
          <LineChart data={chartData}>
            <CartesianGrid vertical={false} />
            <XAxis dataKey="date" tickLine={false} axisLine={false} tickMargin={8} />
            <YAxis tickLine={false} axisLine={false} tickMargin={8} />
            <ChartTooltip content={<ChartTooltipContent />} />
            <Line
              type="monotone"
              dataKey="conversion_rate"
              stroke="var(--color-conversion_rate)"
              strokeWidth={2}
              dot={{ r: 3 }}
            />
          </LineChart>
        </ChartContainer>

        <div className="space-y-2">
          {conversion?.by_source?.map((source) => (
            <div key={source.source} className="flex items-center justify-between rounded-lg border p-3 text-sm">
              <span>{source.source}</span>
              <span className="text-muted-foreground">
                {source.converted}/{source.leads} · {source.conversion_rate.toFixed(1)}%
              </span>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}

function AnalystVelocityCard({ startDate, endDate }: Readonly<{ startDate?: string; endDate?: string }>) {
  const t = useTranslations("dashboardAnalyst");
  const { data, isLoading } = useAnalystSalesVelocity({ start_date: startDate, end_date: endDate });

  if (isLoading) {
    return <Skeleton className="h-80 w-full" />;
  }

  const velocity = data?.data;

  return (
    <Card className="border-0 shadow-sm">
      <CardHeader>
        <CardTitle>{t("velocity.title")}</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-4 md:grid-cols-3">
          <MetricCard label={t("velocity.salesCycle")} value={`${(velocity?.average_sales_cycle_days ?? 0).toFixed(1)} ${t("velocity.daysSuffix")}`} />
          <MetricCard label={t("velocity.averageDealValue")} value={formatCurrency(velocity?.average_deal_value ?? 0)} />
          <MetricCard label={t("velocity.salesVelocity")} value={(velocity?.sales_velocity ?? 0).toFixed(1)} />
        </div>

        <div className="space-y-2">
          {velocity?.by_stage?.map((stage) => (
            <div key={stage.stage} className="flex items-center justify-between rounded-lg border p-3 text-sm">
              <span>{stage.stage}</span>
              <span className="text-muted-foreground">
                {stage.average_days.toFixed(1)} {t("velocity.daysSuffix")}
              </span>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}

function AnalystInsightsCard({ startDate, endDate }: Readonly<{ startDate?: string; endDate?: string }>) {
  const t = useTranslations("dashboardAnalyst");
  const { data, isLoading } = useAnalystAIInsights({ start_date: startDate, end_date: endDate });

  if (isLoading) {
    return <Skeleton className="h-64 w-full" />;
  }

  const insights = data?.data?.insights ?? [];

  return (
    <Card className="border-0 shadow-sm">
      <CardHeader>
        <CardTitle>{t("insights.title")}</CardTitle>
      </CardHeader>
      <CardContent>
        {insights.length === 0 ? (
          <p className="text-sm text-muted-foreground">{t("insights.empty")}</p>
        ) : (
          <div className="space-y-3">
            {insights.map((insight) => (
              <div key={insight.id} className="rounded-lg border p-4">
                <div className="flex items-center justify-between gap-3">
                  <div className="font-medium">{insight.title}</div>
                  <span className="text-xs text-muted-foreground">{insight.impact}</span>
                </div>
                <p className="mt-2 text-sm text-muted-foreground">{insight.description}</p>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}

function MetricCard({ label, value }: Readonly<{ label: string; value: string }>) {
  return (
    <div className="rounded-lg border p-4">
      <div className="text-sm text-muted-foreground">{label}</div>
      <div className="mt-1 text-2xl font-medium">{value}</div>
    </div>
  );
}

function SectionLabel({ label }: Readonly<{ label: string }>) {
  return (
    <div className="flex items-center gap-1.5 sm:gap-2">
      <div className="h-1 w-1 rounded-full bg-primary" />
      <h2 className="text-sm sm:text-base md:text-lg font-medium text-muted-foreground uppercase tracking-wide">
        {label}
      </h2>
    </div>
  );
}

function formatCurrency(value: number) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(value);
}
