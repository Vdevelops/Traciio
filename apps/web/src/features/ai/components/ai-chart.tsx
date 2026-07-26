"use client";

import { memo } from "react";
import {
  Bar,
  BarChart,
  Cell,
  CartesianGrid,
  Line,
  LineChart,
  Pie,
  PieChart,
  XAxis,
  YAxis,
} from "recharts";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  type ChartConfig,
} from "@/components/ui/chart";

const CHART_PATTERN = /<!-- CHART:(.*?) -->/g;

type AIChartType = "bar" | "line" | "donut";

interface AIChartPoint {
  label: string;
  value: number;
}

export interface AIChartSpec {
  type: AIChartType;
  title: string;
  metric?: string;
  data: AIChartPoint[];
}

export function parseAICharts(message: string): {
  cleanMessage: string;
  charts: AIChartSpec[];
} {
  const charts: AIChartSpec[] = [];
  let match: RegExpExecArray | null;

  CHART_PATTERN.lastIndex = 0;
  while ((match = CHART_PATTERN.exec(message)) !== null) {
    try {
      const parsed = JSON.parse(match[1]) as Partial<AIChartSpec>;
      const type =
        parsed.type === "line" || parsed.type === "donut" ? parsed.type : "bar";
      const data = Array.isArray(parsed.data)
        ? parsed.data
            .map((item) => ({
              label: String((item as Partial<AIChartPoint>).label ?? "").trim(),
              value: Number((item as Partial<AIChartPoint>).value ?? 0),
            }))
            .filter((item) => item.label !== "" && Number.isFinite(item.value))
        : [];

      if (data.length > 0) {
        charts.push({
          type,
          title: String(parsed.title ?? "Grafik").trim() || "Grafik",
          metric: String(parsed.metric ?? "value").trim() || "value",
          data,
        });
      }
    } catch {
      // Ignore malformed chart markers so a bad AI response cannot crash chat.
    }
  }

  return {
    cleanMessage: message.replace(CHART_PATTERN, "").trim(),
    charts,
  };
}

interface AIChartsProps {
  charts: AIChartSpec[];
}

const chartConfig = {
  value: {
    label: "Value",
    color: "var(--color-chart-1)",
  },
} satisfies ChartConfig;

const COLORS = [
  "var(--color-chart-1)",
  "var(--color-chart-2)",
  "var(--color-chart-3)",
  "var(--color-chart-4)",
  "var(--color-chart-5)",
];

export const AICharts = memo(function AICharts({ charts }: AIChartsProps) {
  if (charts.length === 0) return null;

  return (
    <div className="my-4 space-y-4">
      {charts.map((chart, index) => {
        const chartData = chart.data.map((item) => ({
          label: item.label,
          value: item.value,
        }));

        return (
          <div
            key={`${chart.title}-${index}`}
            className="rounded-md border border-border bg-card p-4"
          >
            <div className="mb-3">
              <h4 className="text-sm font-medium text-foreground">
                {chart.title}
              </h4>
              {chart.metric && (
                <p className="text-xs text-muted-foreground">{chart.metric}</p>
              )}
            </div>
            <ChartContainer
              config={chartConfig}
              className={chart.type === "donut" ? "h-[300px] w-full" : "h-[260px] w-full"}
            >
              {chart.type === "donut" ? (
                <PieChart>
                  <Pie
                    data={chartData}
                    dataKey="value"
                    nameKey="label"
                    innerRadius={58}
                    outerRadius={96}
                    paddingAngle={3}
                    strokeWidth={0}
                  >
                    {chartData.map((_, sliceIndex) => (
                      <Cell
                        key={`slice-${sliceIndex}`}
                        fill={COLORS[sliceIndex % COLORS.length]}
                      />
                    ))}
                  </Pie>
                  <ChartTooltip
                    cursor={false}
                    content={<ChartTooltipContent nameKey="label" />}
                  />
                </PieChart>
              ) : chart.type === "line" ? (
                <LineChart data={chartData} margin={{ top: 12, right: 12, left: 12, bottom: 12 }}>
                  <CartesianGrid vertical={false} />
                  <XAxis dataKey="label" tickLine={false} axisLine={false} tickMargin={8} />
                  <YAxis tickLine={false} axisLine={false} tickMargin={8} />
                  <ChartTooltip content={<ChartTooltipContent />} />
                  <Line
                    type="monotone"
                    dataKey="value"
                    stroke="var(--color-value)"
                    strokeWidth={2}
                    dot={{ r: 3 }}
                  />
                </LineChart>
              ) : (
                <BarChart data={chartData} margin={{ top: 12, right: 12, left: 12, bottom: 12 }}>
                  <CartesianGrid vertical={false} />
                  <XAxis dataKey="label" tickLine={false} axisLine={false} tickMargin={8} />
                  <YAxis tickLine={false} axisLine={false} tickMargin={8} />
                  <ChartTooltip content={<ChartTooltipContent />} />
                  <Bar
                    dataKey="value"
                    fill="var(--color-value)"
                    radius={[4, 4, 0, 0]}
                  />
                </BarChart>
              )}
            </ChartContainer>
            {chart.type === "donut" && (
              <div className="mt-3 grid gap-2 text-xs sm:grid-cols-2">
                {chartData.map((item, itemIndex) => (
                  <div key={item.label} className="flex items-center justify-between gap-3">
                    <div className="flex min-w-0 items-center gap-2">
                      <span
                        className="h-2.5 w-2.5 shrink-0 rounded-full"
                        style={{ backgroundColor: COLORS[itemIndex % COLORS.length] }}
                      />
                      <span className="truncate text-muted-foreground">{item.label}</span>
                    </div>
                    <span className="shrink-0 font-medium text-foreground">
                      {item.value.toLocaleString("id-ID")}
                    </span>
                  </div>
                ))}
              </div>
            )}
          </div>
        );
      })}
    </div>
  );
});
