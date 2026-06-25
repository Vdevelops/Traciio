"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
import { Building2, Activity, Users, CalendarDays } from "lucide-react";
import { useTranslations } from "next-intl";
import type { AccountActivityReport } from "../types";

interface AccountActivityReportViewerProps {
  data?: AccountActivityReport;
  isLoading: boolean;
}

export function AccountActivityReportViewer({
  data,
  isLoading,
}: AccountActivityReportViewerProps) {
  const t = useTranslations("reportsFeature.accountActivityReportViewer");
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
      <Card>
        <CardHeader className="pb-4">
          <div className="flex items-center gap-2">
            <Building2 className="h-5 w-5" />
            <CardTitle>{data.account_name}</CardTitle>
          </div>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-3">
            <SummaryCard
              title={t("summary.totalVisits")}
              value={data.summary.total_visits}
            />
            <SummaryCard
              title={t("summary.totalActivities")}
              value={data.summary.total_activities}
            />
            <SummaryCard
              title={t("summary.totalContacts")}
              value={data.summary.total_contacts}
            />
          </div>
        </CardContent>
      </Card>

      <div className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <Activity className="h-5 w-5" />
              <CardTitle>{t("activitiesTitle")}</CardTitle>
            </div>
          </CardHeader>
          <CardContent>
            {data.activities.length === 0 ? (
              <p className="text-sm text-muted-foreground">{t("emptyActivities")}</p>
            ) : (
              <div className="space-y-3">
                {data.activities.map((activity) => (
                  <div key={activity.id} className="rounded-lg border p-3">
                    <div className="flex items-start justify-between gap-3">
                      <div className="space-y-1">
                        <div className="font-medium">{activity.type}</div>
                        <p className="text-sm text-muted-foreground">
                          {activity.description}
                        </p>
                      </div>
                      <Badge variant="outline">{activity.user.name}</Badge>
                    </div>
                    <div className="mt-2 flex items-center gap-2 text-xs text-muted-foreground">
                      <CalendarDays className="h-3.5 w-3.5" />
                      {new Date(activity.timestamp).toLocaleString("id-ID")}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <Users className="h-5 w-5" />
              <CardTitle>{t("visitsTitle")}</CardTitle>
            </div>
          </CardHeader>
          <CardContent>
            {data.visits.length === 0 ? (
              <p className="text-sm text-muted-foreground">{t("emptyVisits")}</p>
            ) : (
              <div className="space-y-3">
                {data.visits.map((visit) => (
                  <div key={visit.id} className="rounded-lg border p-3">
                    <div className="flex items-start justify-between gap-3">
                      <div className="space-y-1">
                        <div className="font-medium">{visit.purpose}</div>
                        <p className="text-sm text-muted-foreground">
                          {visit.sales_rep.name}
                        </p>
                      </div>
                      <Badge variant="outline">{visit.status}</Badge>
                    </div>
                    <div className="mt-2 flex items-center gap-2 text-xs text-muted-foreground">
                      <CalendarDays className="h-3.5 w-3.5" />
                      {new Date(visit.visit_date).toLocaleDateString("id-ID")}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

function SummaryCard({ title, value }: Readonly<{ title: string; value: number }>) {
  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-sm font-medium text-muted-foreground">
          {title}
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-medium">{value}</div>
      </CardContent>
    </Card>
  );
}
