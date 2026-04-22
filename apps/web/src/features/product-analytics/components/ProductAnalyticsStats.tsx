"use client";

import { useTranslations } from "next-intl";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Package, TrendingUp, DollarSign, ShoppingCart } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import type { ProductListItem } from "../types";

interface ProductAnalyticsStatsProps {
  readonly data: readonly ProductListItem[];
  readonly isLoading?: boolean;
}

const formatCurrency = (value: number): string => {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(value);
};

const formatNumber = (value: number): string => {
  return value.toLocaleString("id-ID");
};

export function ProductAnalyticsStats({ data, isLoading }: ProductAnalyticsStatsProps) {
  const t = useTranslations("productAnalytics.stats");

  const stats = {
    totalProducts: data.length,
    totalSold: data.reduce((sum, item) => sum + item.total_sold, 0),
    totalRevenue: data.reduce((sum, item) => sum + item.total_revenue, 0),
    totalProfit: data.reduce((sum, item) => sum + item.total_profit, 0),
  };

  const cards = [
    {
      title: t("totalProducts"),
      value: formatNumber(stats.totalProducts),
      icon: Package,
    },
    {
      title: t("totalSold"),
      value: formatNumber(stats.totalSold),
      icon: ShoppingCart,
    },
    {
      title: t("totalRevenue"),
      value: formatCurrency(stats.totalRevenue),
      icon: DollarSign,
    },
    {
      title: t("totalProfit"),
      value: formatCurrency(stats.totalProfit),
      icon: TrendingUp,
    },
  ];

  if (isLoading) {
    return (
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {cards.map((_, index) => (
          <Card key={index}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <Skeleton className="h-4 w-24" />
              <Skeleton className="h-4 w-4 rounded" />
            </CardHeader>
            <CardContent>
              <Skeleton className="h-8 w-32" />
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
      {cards.map((card, index) => {
        const Icon = card.icon;
        return (
          <Card key={index}>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                {card.title}
              </CardTitle>
              <Icon className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-medium">{card.value}</div>
            </CardContent>
          </Card>
        );
      })}
    </div>
  );
}
