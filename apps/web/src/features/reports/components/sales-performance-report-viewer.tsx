"use client";

import { useMemo } from "react";
import {
  Bar,
  BarChart,
  CartesianGrid,
  Line,
  LineChart,
  XAxis,
  YAxis,
} from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Users, TrendingUp, Target, CircleDollarSign } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  type ChartConfig,
} from "@/components/ui/chart";
import type { SalesPerformanceReport } from "../types";
import { useTranslations } from "next-intl";
import { formatEmailToMailto } from "@/lib/utils";

interface SalesPerformanceReportViewerProps {
  data?: SalesPerformanceReport;
  isLoading: boolean;
}

const chartConfig = {
  leads: {
    label: "Leads",
    color: "var(--color-chart-1)",
  },
  deals: {
    label: "Deals",
    color: "var(--color-chart-2)",
  },
  revenue: {
    label: "Revenue",
    color: "var(--color-chart-3)",
  },
  conversion_rate: {
    label: "Conversion",
    color: "var(--color-chart-5)",
  },
} satisfies ChartConfig;

export function SalesPerformanceReportViewer({
  data,
  isLoading,
}: SalesPerformanceReportViewerProps) {
  const t = useTranslations("reportsFeature.salesPerformanceReportViewer");
  const tCommon = useTranslations("reportsFeature.common");

  const formatCurrency = (value: number) =>
    new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(value);

  const chartData = useMemo(
    () =>
      (data?.by_sales_rep ?? []).map((item) => ({
        name: item.sales_rep.name,
        leads: item.lead_count,
        deals: item.total_deals,
        revenue: item.total_revenue,
        conversion_rate: Number(item.conversion_rate.toFixed(1)),
      })),
    [data?.by_sales_rep],
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
      {/* Summary Cards */}
      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              {t("summary.totalLeads")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-medium">{data.summary.total_leads}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              {t("summary.totalDeals")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-medium">{data.summary.total_deals}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              {t("summary.totalRevenue")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-medium">
              {formatCurrency(data.summary.total_revenue)}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              {t("summary.conversionRate")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-medium">{data.summary.conversion_rate.toFixed(1)}%</div>
          </CardContent>
        </Card>
      </div>

      {chartData.length > 0 ? (
        <div className="grid gap-4 xl:grid-cols-5">
          <Card className="xl:col-span-3">
            <CardHeader>
              <CardTitle>{t("charts.volumeTitle")}</CardTitle>
            </CardHeader>
            <CardContent>
              <ChartContainer
                config={chartConfig}
                className="aspect-auto h-[320px] w-full"
              >
                <BarChart data={chartData}>
                  <CartesianGrid vertical={false} />
                  <XAxis
                    dataKey="name"
                    tickLine={false}
                    axisLine={false}
                    tickMargin={8}
                  />
                  <YAxis tickLine={false} axisLine={false} tickMargin={8} />
                  <ChartTooltip
                    content={<ChartTooltipContent indicator="line" />}
                  />
                  <Bar dataKey="leads" fill="var(--color-leads)" radius={4} />
                  <Bar dataKey="deals" fill="var(--color-deals)" radius={4} />
                </BarChart>
              </ChartContainer>
            </CardContent>
          </Card>

          <Card className="xl:col-span-2">
            <CardHeader>
              <CardTitle>{t("charts.revenueTitle")}</CardTitle>
            </CardHeader>
            <CardContent>
              <ChartContainer
                config={chartConfig}
                className="aspect-auto h-[320px] w-full"
              >
                <LineChart data={chartData}>
                  <CartesianGrid vertical={false} />
                  <XAxis
                    dataKey="name"
                    tickLine={false}
                    axisLine={false}
                    tickMargin={8}
                  />
                  <YAxis tickLine={false} axisLine={false} tickMargin={8} />
                  <ChartTooltip
                    content={<ChartTooltipContent indicator="line" />}
                  />
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
        </div>
      ) : null}

      {/* By Sales Rep */}
      {data.by_sales_rep && data.by_sales_rep.length > 0 && (
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <Users className="h-5 w-5" />
              <CardTitle>{t("bySalesRepTitle")}</CardTitle>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {data.by_sales_rep.map((item) => (
                <div
                  key={item.sales_rep.id}
                  className="flex items-center justify-between p-4 rounded-lg border hover:bg-muted/50 transition-colors"
                >
                  <div className="flex-1">
                    <div className="font-medium">{item.sales_rep.name}</div>
                    <a href={formatEmailToMailto(item.sales_rep.email)} className="text-sm text-muted-foreground hover:text-primary hover:underline cursor-pointer min-w-0">{item.sales_rep.email}</a>
                    <div className="flex items-center gap-4 mt-2 text-sm text-muted-foreground">
                      <div className="flex items-center gap-1">
                        <Target className="h-3 w-3" />
                        <span>
                          {t("leadsLabel", { count: item.lead_count })}
                        </span>
                      </div>
                      <div className="flex items-center gap-1">
                        <TrendingUp className="h-3 w-3" />
                        <span>
                          {t("dealsLabel", { count: item.total_deals })}
                        </span>
                      </div>
                      <div className="flex items-center gap-1">
                        <CircleDollarSign className="h-3 w-3" />
                        <span>
                          {t("revenueLabel", {
                            amount: formatCurrency(item.total_revenue),
                          })}
                        </span>
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center gap-3">
                    <Badge variant="outline">
                      {t("conversionBadge", {
                        rate: item.conversion_rate.toFixed(1),
                      })}
                    </Badge>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
