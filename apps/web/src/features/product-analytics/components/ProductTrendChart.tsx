"use client";

import * as React from "react";
import { useTranslations, useLocale } from "next-intl";
import { Line, LineChart, CartesianGrid, XAxis, YAxis } from "recharts";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import { Skeleton } from "@/components/ui/skeleton";
import { TrendingUp } from "lucide-react";
import type { ProductTrend } from "../types";

interface ProductTrendChartProps {
  readonly trend: ProductTrend | undefined;
  readonly isLoading?: boolean;
}

const chartConfig = {
  quantity: {
    label: "Quantity",
    color: "oklch(0.5234 0.1347 144.1672)", // Primary - Hijau
  },
  revenue: {
    label: "Revenue",
    color: "oklch(0.4851 0.1559 30.2175)", // Accent - Oren
  },
} satisfies ChartConfig;

const formatCurrency = (value: number): string => {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
  }).format(value);
};

const formatNumber = (value: number): string => {
  return value.toLocaleString("id-ID");
};

export function ProductTrendChart({ trend, isLoading }: ProductTrendChartProps) {
  const t = useTranslations("productAnalytics.trends");

  const chartData = React.useMemo(() => {
    if (!trend?.trends) return [];
    
    return trend.trends.map((item) => ({
      period: item.period,
      quantity: item.quantity,
      revenue: item.revenue,
    }));
  }, [trend]);

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>{t("title")}</CardTitle>
          <CardDescription>{t("description")}</CardDescription>
        </CardHeader>
        <CardContent>
          <Skeleton className="h-80 w-full" />
        </CardContent>
      </Card>
    );
  }

  if (!trend || chartData.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>{t("title")}</CardTitle>
          <CardDescription>{t("description")}</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-center h-80 text-muted-foreground">
            {t("noData")}
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <TrendingUp className="h-5 w-5" />
          {t("title")} - {trend.product_name}
        </CardTitle>
        <CardDescription>
          {t("sku")}: {trend.product_sku}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig} className="aspect-auto h-[350px] w-full">
          <LineChart data={chartData}>
            <CartesianGrid vertical={false} />
            <XAxis
              dataKey="period"
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              minTickGap={32}
            />
            <YAxis
              yAxisId="left"
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              tickFormatter={formatNumber}
            />
            <YAxis
              yAxisId="right"
              orientation="right"
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              tickFormatter={(value) => {
                if (value >= 1000000) {
                  return `Rp${(value / 1000000).toFixed(1)}M`;
                }
                if (value >= 1000) {
                  return `Rp${(value / 1000).toFixed(0)}K`;
                }
                return `Rp${value}`;
              }}
            />
            <ChartTooltip
              cursor={false}
              content={<ChartTooltipContent />}
              labelFormatter={(value) => {
                return value.toString();
              }}
              formatter={(value, name) => {
                if (name === "revenue") {
                  return formatCurrency(value as number);
                }
                return formatNumber(value as number);
              }}
            />
            <Line
              yAxisId="left"
              type="monotone"
              dataKey="quantity"
              stroke="var(--color-quantity)"
              strokeWidth={2}
              dot={{ r: 4 }}
              name={t("quantity")}
            />
            <Line
              yAxisId="right"
              type="monotone"
              dataKey="revenue"
              stroke="var(--color-revenue)"
              strokeWidth={2}
              dot={{ r: 4 }}
              name={t("revenue")}
            />
          </LineChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}
