"use client";

import { useMemo } from "react";
import { useTranslations, useLocale } from "next-intl";
import { Pie, PieChart, Cell } from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  type ChartConfig,
} from "@/components/ui/chart";
import { Skeleton } from "@/components/ui/skeleton";
import { useDashboardOverview } from "../hooks/useDashboard";
import { BarChart3, Lock } from "lucide-react";
import { useHasPermission } from "@/features/auth/providers/permissions-provider";
import { getCurrentMonthDateRange } from "../dashboard-date-util";

const COLORS = [
  "oklch(0.65 0.15 60)", // orange
  "oklch(0.55 0.15 240)", // blue
  "oklch(0.55 0.15 300)", // purple
  "oklch(0.52 0.13 144)", // green
];

const chartConfig = {
  value: {
    label: "Leads",
    color: "oklch(0.52 0.13 144)",
  },
} satisfies ChartConfig;

// Known source keys that have translations
const KNOWN_SOURCE_KEYS = ["social", "email", "call", "other"] as const;

// Helper function to normalize source name for translation key
function normalizeSourceKey(source: string): string {
  return source.toLowerCase().trim().replaceAll(/\s+/g, "_").replaceAll(/[^a-z0-9_]/g, "");
}

// Helper function to safely get translated source name
function getTranslatedSourceName(
  t: ReturnType<typeof useTranslations<"dashboardOverview">>,
  source: string
): string {
  const normalizedKey = normalizeSourceKey(source);
  
  // Only try to translate if it's a known source key
  if (KNOWN_SOURCE_KEYS.includes(normalizedKey as typeof KNOWN_SOURCE_KEYS[number])) {
    try {
      return t(`leadsBySource.source.${normalizedKey}` as `leadsBySource.source.${typeof KNOWN_SOURCE_KEYS[number]}`);
    } catch {
      // Fallback to original if translation fails
      return source;
    }
  }
  
  // For unknown sources, return the original name
  return source;
}

export function LeadsBySource({ startDate, endDate }: Readonly<{ startDate?: string; endDate?: string }>) {
  const t = useTranslations("dashboardOverview");
  const locale = useLocale();
  const monthParams = getCurrentMonthDateRange();

  // Permission check
  const hasPermission = useHasPermission("leads.view");

  const params = { ...monthParams, start_date: startDate, end_date: endDate };

  const { data, isLoading } = useDashboardOverview(params, {
    enabled: hasPermission,
  });

  console.log("LeadsBySource Params:", params);
  console.log("LeadsBySource Data:", data);

  const overview = data?.data;

  const chartData = useMemo(() => {
    if (!overview?.leads_by_source || overview.leads_by_source.by_source.length === 0) {
      return [];
    }

    return overview.leads_by_source.by_source.map((item) => ({
      name: item.source,
      value: item.count,
    }));
  }, [overview]);

  // Show warning if no permission
  if (!hasPermission) {
    return (
      <Card className="h-full flex flex-col border-0 shadow-sm bg-muted/10">
        <CardHeader className="px-3 sm:px-6 pb-2">
          <CardTitle className="text-xs sm:text-sm font-medium flex items-center gap-2 text-muted-foreground">
            <Lock className="h-3.5 w-3.5" />
            {t("leadsBySource.title")}
          </CardTitle>
        </CardHeader>
        <CardContent className="flex-1 flex flex-col items-center justify-center p-6 text-center">
          <div className="bg-muted/20 p-3 rounded-full mb-3">
            <Lock className="h-5 w-5 text-muted-foreground" />
          </div>
          <p className="text-sm font-medium text-muted-foreground mb-1">
            Access Restricted
          </p>
          <p className="text-xs text-muted-foreground/70 max-w-[200px]">
            You don't have permission to view leads data.
          </p>
        </CardContent>
      </Card>
    );
  }

  if (isLoading) {
    return (
      <Card className="h-full border-0 shadow-sm">
        <CardHeader className="px-3 sm:px-6">
          <CardTitle className="text-xs sm:text-sm">{t("leadsBySource.title")}</CardTitle>
        </CardHeader>
        <CardContent className="px-3 sm:px-6 pb-3 sm:pb-6">
          <Skeleton className="h-40 sm:h-48 w-full rounded-xl" />
        </CardContent>
      </Card>
    );
  }

  if (!overview?.leads_by_source) {
    return null;
  }

  if (chartData.length === 0) {
    return (
      <Card className="h-full border-0 shadow-sm">
        <CardHeader className="px-3 sm:px-6">
          <CardTitle className="text-xs sm:text-sm">{t("leadsBySource.title")}</CardTitle>
        </CardHeader>
        <CardContent className="relative min-h-[150px] sm:min-h-[200px] px-3 sm:px-6 pb-3 sm:pb-6">
          <div className="absolute inset-0 flex flex-col items-center justify-center">
            <BarChart3 className="h-8 w-8 sm:h-12 sm:w-12 text-muted-foreground/50 mb-2 sm:mb-3" />
            <p className="text-xs sm:text-sm text-muted-foreground">
              {t("leadsBySource.empty")}
            </p>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="h-full flex flex-col border-0 shadow-sm">
      <CardHeader className="px-3 sm:px-6">
        <CardTitle className="text-xs sm:text-sm font-medium">
          {t("leadsBySource.title")}
        </CardTitle>
      </CardHeader>
      <CardContent className="flex flex-col sm:flex-row items-center justify-between gap-3 sm:gap-4 flex-1 px-3 sm:px-6 pb-3 sm:pb-6">
        <ChartContainer
          config={chartConfig}
          className="mx-auto aspect-square h-32 w-32 sm:h-40 sm:w-40"
        >
          <PieChart>
            <Pie
              data={chartData}
              dataKey="value"
              nameKey="name"
              innerRadius={32}
              outerRadius={64}
              paddingAngle={4}
              strokeWidth={0}
            >
              {chartData.map((_, index) => (
                <Cell
                  key={`cell-${index}`}
                  fill={COLORS[index % COLORS.length]}
                />
              ))}
            </Pie>
            <ChartTooltip
              cursor={false}
              content={<ChartTooltipContent nameKey="name" />}
            />
          </PieChart>
        </ChartContainer>

        <div className="max-h-full overflow-y-auto space-y-2 sm:space-y-3 pr-1 sm:pr-2 text-xs sm:text-sm flex-1 w-full sm:w-auto">
          <div className="font-semibold text-xs sm:text-sm">
            {overview.leads_by_source.total.toLocaleString(locale)}{" "}
            {t("leadsBySource.totalLabel")}
          </div>
          {chartData.map((item, index) => (
            <div key={item.name} className="flex items-center justify-between gap-2 p-1 rounded-md">
              <div className="flex items-center gap-1.5 sm:gap-2 min-w-0 flex-1">
                <span
                  className="h-1.5 w-1.5 sm:h-2 sm:w-2 rounded-full shrink-0"
                  style={{ backgroundColor: COLORS[index % COLORS.length] }}
                />
                <span className="capitalize truncate text-xs sm:text-sm font-medium">
                  {getTranslatedSourceName(t, item.name)}
                </span>
              </div>
              <span className="text-muted-foreground text-[10px] sm:text-xs font-bold shrink-0">
                {item.value.toLocaleString(locale)}
              </span>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
