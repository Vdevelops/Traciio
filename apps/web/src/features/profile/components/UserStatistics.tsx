"use client";

import { useTranslations } from "next-intl";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { DollarSign, Briefcase, MapPin, CheckCircle2 } from "lucide-react";
import type { ProfileStats } from "../types";

interface UserStatisticsProps {
  readonly statistics: ProfileStats | undefined;
}

export function UserStatistics({ statistics }: UserStatisticsProps) {
  const t = useTranslations("profile");

  if (!statistics) {
    return null;
  }

  // Use extended stats with fallback to basic stats
  const totalRevenue = statistics.total_revenue_formatted ?? "Rp 0";
  const avgDealValue = statistics.average_deal_value_formatted ?? "Rp 0";
  const dealsCount = statistics.deals_won ?? 0;
  const conversionRate = statistics.conversion_rate ?? 0;
  const visitsCount = statistics.visits ?? 0;
  const tasksCount = statistics.tasks ?? 0;

  const stats = [
    {
      label: t("stats.total_revenue"),
      value: totalRevenue,
      detail: `${t("stats.avg_deal_value")}: ${avgDealValue}`,
      icon: DollarSign,
    },
    {
      label: t("stats.deals_closed"),
      value: dealsCount.toString(),
      detail: `${t("stats.conversion_rate")}: ${conversionRate.toFixed(1)}%`,
      icon: Briefcase,
    },
    {
      label: t("stats.visits_completed"),
      value: visitsCount.toString(),
      icon: MapPin,
    },
    {
      label: t("stats.tasks_completed"),
      value: tasksCount.toString(),
      icon: CheckCircle2,
    },
  ];

  return (
    <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
      {stats.map((stat, index) => {
        const Icon = stat.icon;

        return (
          <Card key={stat.label}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">{stat.label}</CardTitle>
              <Icon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-medium">{stat.value}</div>
              {stat.detail && (
                <div className="text-xs text-muted-foreground mt-1">{stat.detail}</div>
              )}
            </CardContent>
          </Card>
        );
      })}
    </div>
  );
}
