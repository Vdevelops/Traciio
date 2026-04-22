"use client";

import { useTranslations } from "next-intl";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { DollarSign, Briefcase, MapPin, CheckCircle2 } from "lucide-react";
import type { SalesRepStatistics } from "../types";

interface SalesRepStatisticsProps {
  readonly statistics: SalesRepStatistics | undefined;
}

export function SalesRepStatistics({ statistics }: SalesRepStatisticsProps) {
  const t = useTranslations("salesOverview");

  if (!statistics) {
    return null;
  }

  const stats = [
    {
      label: t("total_revenue"),
      value: statistics.total_revenue_formatted,
      icon: DollarSign,
    },
    {
      label: t("deals_closed"),
      value: statistics.deals_closed.toString(),
      icon: Briefcase,
    },
    {
      label: t("visits_completed"),
      value: statistics.visits_completed.toString(),
      icon: MapPin,
    },
    {
      label: t("tasks_completed"),
      value: statistics.tasks_completed.toString(),
      icon: CheckCircle2,
    },
  ];

  return (
    <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
      {stats.map((stat, index) => {
        const Icon = stat.icon;

        return (
          <Card key={index}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">{stat.label}</CardTitle>
              <Icon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-medium">{stat.value}</div>
              {stat.label === t("total_revenue") && (
                <div className="text-xs text-muted-foreground mt-1">
                  {t("avg_deal_value")}: {statistics.average_deal_value_formatted}
                </div>
              )}
              {stat.label === t("deals_closed") && (
                <div className="text-xs text-muted-foreground mt-1">
                  {t("conversion_rate")}: {statistics.conversion_rate.toFixed(1)}%
                </div>
              )}
            </CardContent>
          </Card>
        );
      })}
    </div>
  );
}

