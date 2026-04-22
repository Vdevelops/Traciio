"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Users, TrendingUp, Building2, Activity } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import type { SalesPerformanceReport } from "../types";
import { useTranslations } from "next-intl";
import { formatEmailToMailto } from "@/lib/utils";

interface SalesPerformanceReportViewerProps {
  data?: SalesPerformanceReport;
  isLoading: boolean;
}

export function SalesPerformanceReportViewer({
  data,
  isLoading,
}: SalesPerformanceReportViewerProps) {
  const t = useTranslations("reportsFeature.salesPerformanceReportViewer");
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
      <div className="grid gap-4 md:grid-cols-3">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              {t("summary.totalVisits")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-medium">{data.summary.total_visits}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              {t("summary.totalAccounts")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-medium">{data.summary.total_accounts}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              {t("summary.avgVisitsPerAccount")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-medium">
              {data.summary.average_visits_per_account.toFixed(2)}
            </div>
          </CardContent>
        </Card>
      </div>

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
                        <Activity className="h-3 w-3" />
                        <span>
                          {t("visitsLabel", { count: item.visit_count })}
                        </span>
                      </div>
                      <div className="flex items-center gap-1">
                        <Building2 className="h-3 w-3" />
                        <span>
                          {t("accountsLabel", { count: item.account_count })}
                        </span>
                      </div>
                      <div className="flex items-center gap-1">
                        <TrendingUp className="h-3 w-3" />
                        <span>
                          {t("activitiesLabel", { count: item.activity_count })}
                        </span>
                      </div>
                    </div>
                  </div>
                  <div className="flex items-center gap-3">
                    <Badge variant="outline">
                      {t("completionBadge", {
                        rate: item.completion_rate.toFixed(1),
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

