"use client";

import * as React from "react";
import { Bar, BarChart, CartesianGrid, XAxis, YAxis } from "recharts";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import {
  ChartConfig,
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import { TrendingUp, DollarSign, Target, AlertCircle } from "lucide-react";
import { useTranslations } from "next-intl";
import type { PipelineReport } from "../types";

interface SalesFunnelInsightsProps {
  data: PipelineReport;
}

const chartConfig = {
  deals: {
    label: "Deals",
    color: "oklch(0.5234 0.1347 144.1672)", // Primary - Hijau
  },
  value: {
    label: "Value",
    color: "oklch(0.55 0.15 240)", // Biru
  },
} satisfies ChartConfig;

export function SalesFunnelInsights({ data }: SalesFunnelInsightsProps) {
  const t = useTranslations("reportsFeature.salesFunnelInsights");

  // Extract with null safety
  const summary = data?.summary ?? {
    total_deals: 0,
    total_value: 0,
    won_deals: 0,
    lost_deals: 0,
    open_deals: 0,
    expected_revenue: 0,
    won_value: 0,
    open_value: 0,
  };
  const byStage = data?.by_stage ?? {};
  const deals = data?.deals ?? [];

  // Prepare chart data from by_stage and deals
  const stageChartData = React.useMemo(() => {
    if (!byStage || Object.keys(byStage).length === 0) {
      return [];
    }

    // Calculate value per stage from deals
    const stageValues: Record<string, number> = {};
    if (deals.length > 0) {
      deals.forEach((deal) => {
        const stageKey = deal.stage_code || deal.stage;
        if (stageKey) {
          const dealValue = deal.value ?? 0;
          stageValues[stageKey] = (stageValues[stageKey] || 0) + dealValue;
        }
      });
    }

    return Object.entries(byStage).map(([stage, count]) => ({
      stage: stage.replace(/_/g, " ").replace(/\b\w/g, (l) => l.toUpperCase()),
      deals: count ?? 0,
      value: stageValues[stage] || 0,
    }));
  }, [byStage, deals]);

  // Calculate metrics
  const winRate = summary.total_deals > 0
    ? ((summary.won_deals / summary.total_deals) * 100).toFixed(1)
    : "0.0";

  const averageDealValue = summary.total_deals > 0
    ? summary.total_value / summary.total_deals
    : 0;

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
    }).format(value);
  };

  return (
    <div className="space-y-6">
      {/* Key Metrics */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              {t("metricWinRate")}
            </CardTitle>
            <Target className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-medium">{winRate}%</div>
            <p className="text-xs text-muted-foreground">
              {t("metricWinRateDetail", {
                won: summary.won_deals,
                total: summary.total_deals,
              })}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              {t("metricAvgDealValue")}
            </CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-medium">{formatCurrency(averageDealValue)}</div>
            <p className="text-xs text-muted-foreground">
              {t("metricAvgDealValueDetail")}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              {t("metricTotalPipelineValue")}
            </CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-medium">{formatCurrency(summary.total_value)}</div>
            <p className="text-xs text-muted-foreground">
              {t("metricTotalPipelineValueDetail")}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              {t("metricLostDeals")}
            </CardTitle>
            <AlertCircle className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-medium text-red-600">{summary.lost_deals}</div>
            <p className="text-xs text-muted-foreground">
              {t("metricLostDealsDetail")}
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Additional Metrics */}
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Won Value
            </CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-medium">{formatCurrency(summary.won_value ?? 0)}</div>
            <p className="text-xs text-muted-foreground">
              Total value of won deals
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Expected Revenue
            </CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-medium">{formatCurrency(summary.expected_revenue ?? 0)}</div>
            <p className="text-xs text-muted-foreground">
              Total expected revenue
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Open Deals
            </CardTitle>
            <Target className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-medium">{summary.open_deals ?? 0}</div>
            <p className="text-xs text-muted-foreground">
              Currently open opportunities
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">
              Open Value
            </CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-medium">{formatCurrency(summary.open_value ?? 0)}</div>
            <p className="text-xs text-muted-foreground">
              Total value of open deals
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Deals by Stage Chart */}
      {stageChartData.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>{t("chartDealsByStageTitle")}</CardTitle>
            <CardDescription>
              {t("chartDealsByStageDescription")}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <ChartContainer config={chartConfig} className="aspect-auto h-[300px] w-full">
              <BarChart data={stageChartData}>
                <CartesianGrid vertical={false} />
                <XAxis
                  dataKey="stage"
                  tickLine={false}
                  axisLine={false}
                  tickMargin={8}
                  angle={-45}
                  textAnchor="end"
                  height={80}
                />
                <YAxis
                  tickLine={false}
                  axisLine={false}
                  tickMargin={8}
                  tickFormatter={(value) => value.toString()}
                />
                <ChartTooltip
                  cursor={false}
                  content={<ChartTooltipContent indicator="line" />}
                />
                <Bar dataKey="deals" fill="var(--color-deals)" radius={[4, 4, 0, 0]} />
                <ChartLegend content={<ChartLegendContent />} />
              </BarChart>
            </ChartContainer>
          </CardContent>
        </Card>
      )}

      {/* Stage Breakdown */}
      {byStage && Object.keys(byStage).length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle>{t("stageBreakdownTitle")}</CardTitle>
            <CardDescription>
              {t("stageBreakdownDescription")}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {Object.entries(byStage)
                .sort(([, a], [, b]) => (b ?? 0) - (a ?? 0))
                .map(([stage, count]) => {
                  const stageCount = count ?? 0;
                  const percentage = summary.total_deals > 0
                    ? ((stageCount / summary.total_deals) * 100).toFixed(1)
                    : "0.0";

                  return (
                    <div key={stage} className="space-y-2">
                      <div className="flex items-center justify-between text-sm">
                        <span className="font-medium capitalize">
                          {stage.replace(/_/g, " ")}
                        </span>
                        <span className="text-muted-foreground">
                          {t("stageBreakdownItem", { count: stageCount, percentage })}
                        </span>
                      </div>
                      <div className="h-2 w-full rounded-full bg-muted">
                        <div
                          className="h-2 rounded-full bg-primary transition-all"
                          style={{ width: `${percentage}%` }}
                        />
                      </div>
                    </div>
                  );
                })}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Empty State */}
      {(!byStage || Object.keys(byStage).length === 0) && (
        <Card>
          <CardContent className="py-8 text-center text-muted-foreground">
            <p className="text-sm">
              {t("emptyStageTitle")}
            </p>
          </CardContent>
        </Card>
      )}
    </div>
  );
}

