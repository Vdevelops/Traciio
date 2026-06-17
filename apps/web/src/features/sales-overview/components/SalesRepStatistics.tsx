"use client";

import { useTranslations } from "next-intl";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  DollarSign,
  Briefcase,
  MapPin,
  CheckCircle2,
  Target,
  TrendingDown,
  Trophy,
  Clock3,
} from "lucide-react";
import type { SalesRepStatistics } from "../types";

interface SalesRepStatisticsProps {
  readonly statistics: SalesRepStatistics | undefined;
}

export function SalesRepStatistics({ statistics }: SalesRepStatisticsProps) {
  const t = useTranslations("salesOverview");

  if (!statistics) {
    return null;
  }

  const prospectOutcome = statistics.prospect_outcome;

  const stats = [
    {
      label: t("total_revenue"),
      value: statistics.total_revenue_formatted,
      icon: DollarSign,
      subValue: `${t("avg_deal_value")}: ${statistics.average_deal_value_formatted}`,
    },
    {
      label: t("deals_closed"),
      value: statistics.deals_closed.toString(),
      icon: Briefcase,
      subValue: `${t("conversion_rate")}: ${statistics.conversion_rate.toFixed(1)}%`,
    },
    {
      label: t("prospects.total"),
      value: (prospectOutcome?.total_prospects ?? 0).toString(),
      icon: Target,
      subValue: `${t("prospects.conversion")}: ${(prospectOutcome?.prospect_conversion_rate ?? 0).toFixed(1)}%`,
    },
    {
      label: t("prospects.won"),
      value: (prospectOutcome?.won_prospects ?? 0).toString(),
      icon: Trophy,
      subValue: t("prospects.from_deals"),
    },
    {
      label: t("prospects.lost"),
      value: (prospectOutcome?.lost_prospects ?? 0).toString(),
      icon: TrendingDown,
      subValue: t("prospects.from_deals"),
    },
    {
      label: t("prospects.open"),
      value: (prospectOutcome?.open_prospects ?? 0).toString(),
      icon: Clock3,
      subValue: t("prospects.from_deals"),
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
    <div className="grid gap-6 md:grid-cols-2 xl:grid-cols-4">
      {stats.map((stat, index) => {
        const Icon = stat.icon;

        return (
          <Card key={index}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">
                {stat.label}
              </CardTitle>
              <Icon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-medium">{stat.value}</div>
              {stat.subValue ? (
                <div className="text-xs text-muted-foreground mt-1">
                  {stat.subValue}
                </div>
              ) : null}
            </CardContent>
          </Card>
        );
      })}
    </div>
  );
}
