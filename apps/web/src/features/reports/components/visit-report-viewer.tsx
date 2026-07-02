"use client";

import { useMemo } from "react";
import {
  CartesianGrid,
  Cell,
  Line,
  LineChart,
  Pie,
  PieChart,
  XAxis,
  YAxis,
} from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Building2, Users, Calendar } from "lucide-react";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  type ChartConfig,
} from "@/components/ui/chart";
import type { VisitReportReport } from "../types";
import { useTranslations } from "next-intl";

interface VisitReportViewerProps {
  data?: VisitReportReport;
  isLoading: boolean;
}

const visitChartConfig = {
  visits: {
    label: "Visits",
    color: "var(--color-chart-1)",
  },
} satisfies ChartConfig;

const pieColors = [
  "var(--color-chart-1)",
  "var(--color-chart-2)",
  "var(--color-chart-3)",
  "var(--color-chart-4)",
  "var(--color-chart-5)",
];

export function VisitReportViewer({ data, isLoading }: VisitReportViewerProps) {
  const t = useTranslations("reportsFeature.visitReportViewer");
  const tCommon = useTranslations("reportsFeature.common");
  const summary = data?.summary ?? {
    total: 0,
    completed: 0,
    pending: 0,
    approved: 0,
    rejected: 0,
  };
  const byAccount = data?.by_account ?? [];
  const bySalesRep = data?.by_sales_rep ?? [];
  const byDate = data?.by_date ?? [];
  const byStatus = data?.by_status ?? {};

  const byDateChartData = useMemo(
    () =>
      [...byDate]
        .sort((a, b) => a.date.localeCompare(b.date))
        .map((item) => ({
          date: new Date(item.date).toLocaleDateString("id-ID", {
            day: "2-digit",
            month: "short",
          }),
          rawDate: item.date,
          visits: item.count ?? 0,
        })),
    [byDate],
  );

  const statusChartData = useMemo(
    () =>
      Object.entries(byStatus)
        .map(([status, count]) => ({
          name: status,
          value: count,
        }))
        .sort((a, b) => b.value - a.value),
    [byStatus],
  );

  if (isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-32 w-full" />
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  if (!data) {
    return (
      <div className="text-center py-8 text-muted-foreground">
        <p className="text-sm">{tCommon("noData")}</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="grid gap-4 md:grid-cols-5">
        <SummaryCard title={t("summary.total")} value={summary.total} />
        <SummaryCard title={t("summary.completed")} value={summary.completed} />
        <SummaryCard title={t("summary.pending")} value={summary.pending} />
        <SummaryCard title={t("summary.approved")} value={summary.approved} />
        <SummaryCard title={t("summary.rejected")} value={summary.rejected} />
      </div>

      <div className="grid gap-4 xl:grid-cols-5">
        <Card className="xl:col-span-3">
          <CardHeader>
            <CardTitle>{t("charts.trendTitle")}</CardTitle>
          </CardHeader>
          <CardContent>
            <ChartContainer
              config={visitChartConfig}
              className="aspect-auto h-[320px] w-full"
            >
              <LineChart data={byDateChartData}>
                <CartesianGrid vertical={false} />
                <XAxis
                  dataKey="date"
                  tickLine={false}
                  axisLine={false}
                  tickMargin={8}
                />
                <YAxis tickLine={false} axisLine={false} tickMargin={8} />
                <ChartTooltip content={<ChartTooltipContent />} />
                <Line
                  type="monotone"
                  dataKey="visits"
                  stroke="var(--color-visits)"
                  strokeWidth={2}
                  dot={{ r: 3 }}
                />
              </LineChart>
            </ChartContainer>
          </CardContent>
        </Card>

        <Card className="xl:col-span-2">
          <CardHeader>
            <CardTitle>{t("charts.statusTitle")}</CardTitle>
          </CardHeader>
          <CardContent>
            <ChartContainer
              config={visitChartConfig}
              className="aspect-auto h-[320px] w-full"
            >
              <PieChart>
                <Pie
                  data={statusChartData}
                  dataKey="value"
                  nameKey="name"
                  innerRadius={60}
                  outerRadius={96}
                  paddingAngle={3}
                >
                  {statusChartData.map((entry, index) => (
                    <Cell
                      key={entry.name}
                      fill={pieColors[index % pieColors.length]}
                    />
                  ))}
                </Pie>
                <ChartTooltip content={<ChartTooltipContent hideLabel />} />
              </PieChart>
            </ChartContainer>
          </CardContent>
        </Card>
      </div>

      {byAccount.length > 0 && (
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <Building2 className="h-5 w-5" />
              <CardTitle>{t("byAccountTitle")}</CardTitle>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {[...byAccount]
                .sort((a, b) => (b.visit_count ?? 0) - (a.visit_count ?? 0))
                .map((item) => {
                  const account = item.account ?? {
                    id: "",
                    name: t("unknownAccount"),
                  };
                  const visitCount = item.visit_count ?? 0;

                  return (
                    <div
                      key={account.id || account.name}
                      className="flex items-center justify-between rounded-lg border p-3"
                    >
                      <div className="font-medium">{account.name}</div>
                      <div className="text-sm text-muted-foreground">
                        {t("visitCount", { count: visitCount })}
                      </div>
                    </div>
                  );
                })}
            </div>
          </CardContent>
        </Card>
      )}

      {bySalesRep.length > 0 && (
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <Users className="h-5 w-5" />
              <CardTitle>{t("bySalesRepTitle")}</CardTitle>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {[...bySalesRep]
                .sort((a, b) => (b.visit_count ?? 0) - (a.visit_count ?? 0))
                .map((item) => {
                  const salesRep = item.sales_rep ?? {
                    id: "",
                    name: t("unknownSalesRep"),
                  };
                  const visitCount = item.visit_count ?? 0;

                  return (
                    <div
                      key={salesRep.id || salesRep.name}
                      className="flex items-center justify-between rounded-lg border p-3"
                    >
                      <div className="font-medium">{salesRep.name}</div>
                      <div className="text-sm text-muted-foreground">
                        {t("visitCount", { count: visitCount })}
                      </div>
                    </div>
                  );
                })}
            </div>
          </CardContent>
        </Card>
      )}

      {byDate.length > 0 && (
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <Calendar className="h-5 w-5" />
              <CardTitle>{t("byDateTitle")}</CardTitle>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              {byDateChartData.slice(0, 20).map((item) => {
                const maxCount = Math.max(
                  ...byDateChartData.map((d) => d.visits),
                  1,
                );

                return (
                  <div
                    key={item.rawDate}
                    className="flex items-center gap-3"
                  >
                    <div className="w-32 text-sm text-muted-foreground">
                      {new Date(item.rawDate).toLocaleDateString("id-ID", {
                        month: "short",
                        day: "numeric",
                        year: "numeric",
                      })}
                    </div>
                    <div className="flex-1">
                      <div className="relative h-6 overflow-hidden rounded-full bg-muted">
                        <div
                          className="h-full rounded-full bg-primary transition-all"
                          style={{
                            width: `${(item.visits / maxCount) * 100}%`,
                          }}
                        />
                      </div>
                    </div>
                    <div className="w-12 text-right text-sm font-medium">
                      {item.visits}
                    </div>
                  </div>
                );
              })}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}

function SummaryCard({ title, value }: Readonly<{ title: string; value: number }>) {
  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">
          {title}
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-medium">{value}</div>
      </CardContent>
    </Card>
  );
}
