"use client";

import { useTranslations } from "next-intl";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { usePipelineSummary } from "@/features/dashboard/hooks/useDashboard";
import { useHasPermission } from "@/features/auth/providers/permissions-provider";
import type { PipelineSummaryStage } from "@/features/dashboard/types";
import { getCurrentMonthDateRange } from "../dashboard-date-util";
import { Lock } from "lucide-react";

export function PipelineSummary({ startDate, endDate }: Readonly<{ startDate?: string; endDate?: string }>) {
  const t = useTranslations("pipelineSummary");
  const monthParams = getCurrentMonthDateRange();
  
  // Permission check
  const hasPermission = useHasPermission("pipeline.view");

  const params = { 
    ...monthParams, 
    ...(startDate ? { start_date: startDate } : {}),
    ...(endDate ? { end_date: endDate } : {})
  };

  const { data, isLoading } = usePipelineSummary(params, {
    enabled: hasPermission
  });

  console.log("PipelineSummary Params:", params);
  console.log("PipelineSummary Data:", data);
  
  const summary = data?.data;

  // Use stages from API response (already includes all stages with colors, sorted by order from backend)
  const stages: PipelineSummaryStage[] = summary?.by_stage || [];

  // Show warning if no permission
  if (!hasPermission) {
    return (
      <Card className="h-full flex flex-col border-0 shadow-sm bg-muted/10">
        <CardHeader className="px-3 sm:px-6 pb-2">
          <CardTitle className="text-xs sm:text-sm font-medium flex items-center gap-2 text-muted-foreground">
            <Lock className="h-3.5 w-3.5" />
            {t("title")}
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
            You don't have permission to view pipeline data.
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="h-full flex flex-col border-0 shadow-sm">
      <CardHeader className="px-3 sm:px-6">
        <CardTitle className="text-xs sm:text-sm font-medium">{t("title")}</CardTitle>
        <p className="text-[10px] sm:text-xs text-muted-foreground mt-0.5">{t("description")}</p>
      </CardHeader>
      <CardContent className="flex-1 px-3 sm:px-6 pb-3 sm:pb-6">
        {isLoading && stages.length === 0 ? (
          <>
            <Skeleton className="mb-3 sm:mb-4 h-4 sm:h-6 w-32 sm:w-40" />
            <Skeleton className="h-2 sm:h-3 w-full rounded-full" />
            <div className="mt-3 sm:mt-4 space-y-1.5 sm:space-y-2">
              {Array.from({ length: 4 }, (_, index) => (
                <Skeleton key={`pipeline-skeleton-${index}`} className="h-8 sm:h-10 w-full" />
              ))}
            </div>
          </>
        ) : null}

        {!isLoading && stages.length === 0 ? (
          <p className="text-[10px] sm:text-xs text-muted-foreground">
            {t("noData", { default: "No pipeline data" })}
          </p>
        ) : null}

        {stages.length > 0 && (
          <>
            {/* Top stacked bar */}
            <div className="mb-3 sm:mb-4 flex h-2 sm:h-3 overflow-hidden rounded-full bg-muted">
              {stages.map((stage) => (
                <div
                  key={stage.stage_id}
                  style={{
                    width: `${stage.percentage || 0}%`,
                    backgroundColor: stage.stage_color || "#CBD5F5",
                  }}
                />
              ))}
            </div>

            {/* Stage rows with max height and scroll */}
            <div className="max-h-full overflow-y-auto space-y-2 sm:space-y-3 pr-1 sm:pr-2">
              {stages.map((stage) => (
                <div
                  key={stage.stage_id}
                  className="flex items-center justify-between gap-2 sm:gap-4 text-[10px] sm:text-xs"
                >
                  <div className="flex flex-1 items-center gap-1.5 sm:gap-2 min-w-0">
                    <span
                      className="inline-flex h-1.5 w-1.5 sm:h-2 sm:w-2 rounded-full shrink-0"
                      style={{
                        backgroundColor: stage.stage_color || "#CBD5F5",
                      }}
                    />
                    <div className="min-w-0 flex-1">
                      <div className="truncate text-xs sm:text-sm font-medium">
                        {stage.stage_name}
                      </div>
                      <div className="text-[10px] sm:text-xs text-muted-foreground truncate">
                        {stage.deal_count.toLocaleString("id-ID")} deals ·{" "}
                        {stage.total_value_formatted}
                      </div>
                    </div>
                  </div>
                  <div className="flex w-16 sm:w-24 items-center gap-1 sm:gap-2 shrink-0">
                    <div className="h-1 sm:h-1.5 flex-1 overflow-hidden rounded-full bg-muted">
                      <div
                        className="h-1 sm:h-1.5 rounded-full"
                        style={{
                          width: `${stage.percentage || 0}%`,
                          backgroundColor: stage.stage_color || "#CBD5F5",
                        }}
                      />
                    </div>
                    <span className="w-8 sm:w-10 text-right text-[10px] sm:text-xs text-muted-foreground">
                      {Math.round(stage.percentage)}%
                    </span>
                  </div>
                </div>
              ))}
            </div>
          </>
        )}
      </CardContent>
    </Card>
  );
}

