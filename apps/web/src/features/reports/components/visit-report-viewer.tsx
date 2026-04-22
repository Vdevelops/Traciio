"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { BarChart3, Building2, Users, Calendar } from "lucide-react";
import type { VisitReportReport } from "../types";
import { useTranslations } from "next-intl";

interface VisitReportViewerProps {
  data?: VisitReportReport;
  isLoading: boolean;
}

export function VisitReportViewer({ data, isLoading }: VisitReportViewerProps) {
  const t = useTranslations("reportsFeature.visitReportViewer");
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

  const summary = data.summary ?? {
    total: 0,
    completed: 0,
    pending: 0,
    approved: 0,
    rejected: 0,
  };
  const byAccount = data.by_account ?? [];
  const bySalesRep = data.by_sales_rep ?? [];
  const byDate = data.by_date ?? [];

  return (
    <div className="space-y-6">
      {/* Summary Cards */}
      <div className="grid gap-4 md:grid-cols-5">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              {t("summary.total")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-medium">{summary.total}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              {t("summary.completed")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-medium">{summary.completed}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              {t("summary.pending")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-medium">{summary.pending}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              {t("summary.approved")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-medium">{summary.approved}</div>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              {t("summary.rejected")}
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-medium">{summary.rejected}</div>
          </CardContent>
        </Card>
      </div>

      {/* By Account */}
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
              {byAccount.map((item) => {
                const account = item.account ?? { id: "", name: "Unknown Account" };
                const visitCount = item.visit_count ?? 0;
                
                return (
                  <div
                    key={account.id || `account-${Math.random()}`}
                    className="flex items-center justify-between p-3 rounded-lg border"
                  >
                    <div className="font-medium">{account.name}</div>
                    <div className="text-sm text-muted-foreground">{visitCount} visits</div>
                  </div>
                );
              })}
            </div>
          </CardContent>
        </Card>
      )}

      {/* By Sales Rep */}
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
              {bySalesRep.map((item) => {
                const salesRep = item.sales_rep ?? { id: "", name: "Unknown Sales Rep" };
                const visitCount = item.visit_count ?? 0;
                
                return (
                  <div
                    key={salesRep.id || `sales-rep-${Math.random()}`}
                    className="flex items-center justify-between p-3 rounded-lg border"
                  >
                    <div className="font-medium">{salesRep.name}</div>
                    <div className="text-sm text-muted-foreground">{visitCount} visits</div>
                  </div>
                );
              })}
            </div>
          </CardContent>
        </Card>
      )}

      {/* By Date */}
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
              {byDate.slice(0, 20).map((item) => {
                const itemDate = item.date ?? "";
                const itemCount = item.count ?? 0;
                const maxCount = Math.max(...byDate.map((d) => d.count ?? 0), 1);
                
                return (
                  <div key={itemDate || `date-${Math.random()}`} className="flex items-center gap-3">
                    <div className="text-sm text-muted-foreground w-32">
                      {itemDate ? new Date(itemDate).toLocaleDateString("id-ID", {
                        month: "short",
                        day: "numeric",
                        year: "numeric",
                      }) : "-"}
                    </div>
                    <div className="flex-1">
                      <div className="h-6 bg-muted rounded-full relative overflow-hidden">
                        <div
                          className="h-full bg-primary rounded-full transition-all"
                          style={{
                            width: `${(itemCount / maxCount) * 100}%`,
                          }}
                        />
                      </div>
                    </div>
                    <div className="text-sm font-medium w-12 text-right">{itemCount}</div>
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

