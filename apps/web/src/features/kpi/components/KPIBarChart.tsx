"use client";

import React from "react";
import { Bar, BarChart, CartesianGrid, XAxis, YAxis } from "recharts";

import {
  ChartConfig,
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

export interface KPIBarChartItem {
  label: string;
  value: number;
}

interface Props {
  title: string;
  data: KPIBarChartItem[];
  description?: string;
  valueSuffix?: string;
  maxValue?: number;
}

const chartConfig = {
  value: {
    label: "Nilai",
    color: "#2563EB",
  },
} satisfies ChartConfig;

export default function KPIBarChart({ title, data, description, valueSuffix, maxValue }: Props) {
  const chartData = data.map((item) => ({
    ...item,
    value: Number.isFinite(item.value) ? item.value : 0,
  }));

  return (
    <Card>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
        {description && <CardDescription>{description}</CardDescription>}
      </CardHeader>
      <CardContent>
        {chartData.length === 0 ? (
          <div className="flex h-[260px] items-center justify-center rounded-lg border border-dashed bg-muted/20 text-sm text-muted-foreground">
            Belum ada data chart untuk periode ini.
          </div>
        ) : (
          <ChartContainer config={chartConfig} className="h-[260px] w-full">
            <BarChart data={chartData} margin={{ left: -16, right: 12, top: 12, bottom: 8 }}>
              <CartesianGrid vertical={false} strokeDasharray="3 3" />
              <XAxis
                dataKey="label"
                tickLine={false}
                axisLine={false}
                interval={0}
                tickMargin={8}
                tick={{ fontSize: 11 }}
              />
              <YAxis
                tickLine={false}
                axisLine={false}
                tickMargin={8}
                tickFormatter={(value) => `${value}${valueSuffix ?? ""}`}
                domain={[0, maxValue ?? "auto"]}
              />
              <ChartTooltip
                content={
                  <ChartTooltipContent
                    formatter={(value) => (
                      <span className="font-medium">
                        {Number(value).toLocaleString("id-ID")}
                        {valueSuffix ?? ""}
                      </span>
                    )}
                  />
                }
              />
              <Bar dataKey="value" fill="var(--color-value)" radius={[6, 6, 0, 0]} />
            </BarChart>
          </ChartContainer>
        )}
      </CardContent>
    </Card>
  );
}
