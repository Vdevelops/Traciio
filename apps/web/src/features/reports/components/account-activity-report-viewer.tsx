"use client";

import { useMemo } from "react";
import {
  Bar,
  BarChart,
  CartesianGrid,
  Pie,
  PieChart,
  Cell,
  XAxis,
  YAxis,
} from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
import { Building2, Activity, Users, CalendarDays } from "lucide-react";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  type ChartConfig,
} from "@/components/ui/chart";
import { useTranslations } from "next-intl";
import type { AccountActivityReport } from "../types";

interface AccountActivityReportViewerProps {
  data?: AccountActivityReport;
  isLoading: boolean;
}

const activityChartConfig = {
  count: {
    label: "Count",
    color: "var(--color-chart-1)",
  },
} satisfies ChartConfig;

const activityPieColors = [
  "var(--color-chart-1)",
  "var(--color-chart-2)",
  "var(--color-chart-3)",
  "var(--color-chart-4)",
  "var(--color-chart-5)",
];

export function AccountActivityReportViewer({
  data,
  isLoading,
}: AccountActivityReportViewerProps) {
  const t = useTranslations("reportsFeature.accountActivityReportViewer");
  const tCommon = useTranslations("reportsFeature.common");
  const activities = data?.activities ?? [];
  const visits = data?.visits ?? [];
  const accountName = data?.account_name ?? "";
  const summary = data?.summary ?? {
    total_visits: 0,
    total_activities: 0,
    total_contacts: 0,
  };

  const activityTypeChartData = useMemo(() => {
    const map = new Map<string, number>();

    activities.forEach((item) => {
      const current = map.get(item.type) ?? 0;
      map.set(item.type, current + 1);
    });

    return Array.from(map.entries())
      .map(([name, count]) => ({ name, count }))
      .sort((a, b) => b.count - a.count);
  }, [activities]);

  const visitStatusChartData = useMemo(() => {
    const map = new Map<string, number>();

    visits.forEach((item) => {
      const current = map.get(item.status) ?? 0;
      map.set(item.status, current + 1);
    });

    return Array.from(map.entries())
      .map(([name, value]) => ({ name, value }))
      .sort((a, b) => b.value - a.value);
  }, [visits]);

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
            <CardTitle>{accountName}</CardTitle>
          </div>
        </CardHeader>
        <CardContent>
          <div className="grid gap-4 md:grid-cols-3">
            <SummaryCard
              title={t("summary.totalVisits")}
              value={summary.total_visits}
            />
            <SummaryCard
              title={t("summary.totalActivities")}
              value={summary.total_activities}
            />
            <SummaryCard
              title={t("summary.totalContacts")}
              value={summary.total_contacts}
            />
          </div>
        </CardContent>
      </Card>

      <div className="grid gap-4 xl:grid-cols-5">
        <Card className="xl:col-span-3">
          <CardHeader>
            <CardTitle>{t("charts.activityTypeTitle")}</CardTitle>
          </CardHeader>
          <CardContent>
            <ChartContainer
              config={activityChartConfig}
              className="aspect-auto h-[320px] w-full"
            >
              <BarChart data={activityTypeChartData}>
                <CartesianGrid vertical={false} />
                <XAxis
                  dataKey="name"
                  tickLine={false}
                  axisLine={false}
                  tickMargin={8}
                />
                <YAxis tickLine={false} axisLine={false} tickMargin={8} />
                <ChartTooltip content={<ChartTooltipContent />} />
                <Bar dataKey="count" fill="var(--color-count)" radius={4} />
              </BarChart>
            </ChartContainer>
          </CardContent>
        </Card>

        <Card className="xl:col-span-2">
          <CardHeader>
            <CardTitle>{t("charts.visitStatusTitle")}</CardTitle>
          </CardHeader>
          <CardContent>
            <ChartContainer
              config={activityChartConfig}
              className="aspect-auto h-[320px] w-full"
            >
              <PieChart>
                <Pie
                  data={visitStatusChartData}
                  dataKey="value"
                  nameKey="name"
                  innerRadius={60}
                  outerRadius={96}
                  paddingAngle={3}
                >
                  {visitStatusChartData.map((entry, index) => (
                    <Cell
                      key={entry.name}
                      fill={activityPieColors[index % activityPieColors.length]}
                    />
                  ))}
                </Pie>
                <ChartTooltip content={<ChartTooltipContent hideLabel />} />
              </PieChart>
            </ChartContainer>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <Activity className="h-5 w-5" />
              <CardTitle>{t("activitiesTitle")}</CardTitle>
            </div>
          </CardHeader>
          <CardContent>
            {activities.length === 0 ? (
              <p className="text-sm text-muted-foreground">{t("emptyActivities")}</p>
            ) : (
              <div className="space-y-3">
                {activities.map((activity) => (
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
            {visits.length === 0 ? (
              <p className="text-sm text-muted-foreground">{t("emptyVisits")}</p>
            ) : (
              <div className="space-y-3">
                {visits.map((visit) => (
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
