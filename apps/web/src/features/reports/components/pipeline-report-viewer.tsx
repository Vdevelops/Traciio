"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { TrendingUp } from "lucide-react";
import { useTranslations } from "next-intl";
import type { PipelineReport } from "../types";

interface PipelineReportViewerProps {
  data?: PipelineReport;
  isLoading: boolean;
}

export function PipelineReportViewer({ data, isLoading }: PipelineReportViewerProps) {
  const t = useTranslations("reportsFeature.pipelineReportViewer");
  const tCommon = useTranslations("reportsFeature.common");

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
      <div className="grid gap-4 md:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              {t("totalDeals")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-medium">{data.summary.total_deals}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              {t("totalValue")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-medium">
              {new Intl.NumberFormat("id-ID", {
                style: "currency",
                currency: "IDR",
                minimumFractionDigits: 0,
              }).format(data.summary.total_value)}
            </div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              {t("wonDeals")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-medium text-green-600">{data.summary.won_deals}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              {t("lostDeals")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-medium text-red-600">{data.summary.lost_deals}</div>
          </CardContent>
        </Card>
      </div>

      {/* By Stage */}
      {data.by_stage && Object.keys(data.by_stage).length > 0 && (
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <TrendingUp className="h-5 w-5" />
              <CardTitle>{t("byStageTitle")}</CardTitle>
            </div>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {Object.entries(data.by_stage).map(([stage, count]) => (
                <div key={stage} className="flex items-center justify-between p-3 rounded-lg border">
                  <div className="font-medium capitalize">{stage.replace(/_/g, " ")}</div>
                  <div className="text-sm text-muted-foreground">
                    {t("byStageItem", { count })}
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

