"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { TrendingUp, Package, DollarSign, Users, ShoppingCart } from "lucide-react";
import { cn } from "@/lib/utils";
import { useTranslations } from "next-intl";
import type { ProductPerformance, BuyerData } from "../types";
import { Skeleton } from "@/components/ui/skeleton";

interface ProductPerformanceCardProps {
  readonly performance: ProductPerformance;
  readonly isLoading?: boolean;
}

const formatCurrency = (value: number): string => {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
  }).format(value);
};

export function ProductPerformanceCard({ performance, isLoading }: ProductPerformanceCardProps) {
  const t = useTranslations("productAnalytics.performance");

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <Skeleton className="h-6 w-48" />
          <Skeleton className="h-4 w-32 mt-2" />
        </CardHeader>
        <CardContent className="space-y-4">
          <Skeleton className="h-20 w-full" />
          <Skeleton className="h-20 w-full" />
          <Skeleton className="h-20 w-full" />
        </CardContent>
      </Card>
    );
  }

  if (!performance) {
    return null;
  }

  const metrics = [
    {
      label: t("totalQuantity"),
      value: performance.total_quantity.toLocaleString("id-ID"),
      icon: Package,
      color: "text-blue-600",
    },
    {
      label: t("totalRevenue"),
      value: formatCurrency(performance.total_revenue),
      icon: DollarSign,
      color: "text-green-600",
    },
    {
      label: t("avgPrice"),
      value: formatCurrency(performance.avg_price),
      icon: ShoppingCart,
      color: "text-purple-600",
    },
    {
      label: t("uniqueBuyers"),
      value: performance.unique_buyers.toLocaleString("id-ID"),
      icon: Users,
      color: "text-orange-600",
    },
  ];

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center justify-between">
          <span>{performance.product_name}</span>
          {performance.growth_rate !== 0 && (
            <span
              className={cn(
                "flex items-center gap-1 text-sm font-medium",
                performance.growth_rate > 0 ? "text-green-600" : "text-red-600"
              )}
            >
              <TrendingUp
                className={cn(
                  "h-4 w-4",
                  performance.growth_rate < 0 && "rotate-180"
                )}
              />
              {Math.abs(performance.growth_rate).toFixed(1)}%
            </span>
          )}
        </CardTitle>
        <CardDescription>{t("sku")}: {performance.product_sku}</CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="grid grid-cols-2 gap-4">
          {metrics.map((metric) => {
            const Icon = metric.icon;
            return (
              <div key={metric.label} className="flex items-center gap-3 p-3 rounded-lg border">
                <div className={cn("p-2 rounded-md bg-muted", metric.color)}>
                  <Icon className="h-5 w-5 text-white" />
                </div>
                <div>
                  <p className="text-sm text-muted-foreground">{metric.label}</p>
                  <p className="text-lg font-medium">{metric.value}</p>
                </div>
              </div>
            );
          })}
        </div>

        {performance.top_buyers && performance.top_buyers.length > 0 && (
          <div>
            <h4 className="text-sm font-medium mb-3">{t("topBuyers")}</h4>
            <div className="space-y-2">
              {performance.top_buyers.slice(0, 5).map((buyer: BuyerData) => (
                <div
                  key={buyer.user_id}
                  className="flex items-center justify-between p-2 rounded-md hover:bg-muted/50 transition-colors"
                >
                  <div className="flex items-center gap-2">
                    <Users className="h-4 w-4 text-muted-foreground" />
                    <span className="text-sm font-medium">{buyer.user_name}</span>
                  </div>
                  <div className="text-right">
                    <div className="text-sm font-medium">
                      {formatCurrency(buyer.revenue)}
                    </div>
                    <div className="text-xs text-muted-foreground">
                      {buyer.quantity} {t("units")}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
