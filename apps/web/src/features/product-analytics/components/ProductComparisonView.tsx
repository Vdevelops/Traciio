"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useTranslations } from "next-intl";
import { Package, DollarSign, Users, TrendingUp } from "lucide-react";
import { cn } from "@/lib/utils";
import { Skeleton } from "@/components/ui/skeleton";
import type { ProductPerformance } from "../types";

interface ProductComparisonViewProps {
  readonly products: readonly ProductPerformance[];
  readonly isLoading?: boolean;
}

const formatCurrency = (value: number): string => {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
  }).format(value);
};

const formatNumber = (value: number): string => {
  return value.toLocaleString("id-ID");
};

export function ProductComparisonView({ products, isLoading }: ProductComparisonViewProps) {
  const t = useTranslations("productAnalytics.comparison");

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>{t("title")}</CardTitle>
        </CardHeader>
        <CardContent>
          <Skeleton className="h-96 w-full" />
        </CardContent>
      </Card>
    );
  }

  if (!products || products.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>{t("title")}</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex items-center justify-center h-64 text-muted-foreground">
            {t("noData")}
          </div>
        </CardContent>
      </Card>
    );
  }

  const metrics = [
    {
      key: "total_quantity",
      label: t("quantity"),
      icon: Package,
      color: "text-blue-600",
      format: formatNumber,
    },
    {
      key: "total_revenue",
      label: t("revenue"),
      icon: DollarSign,
      color: "text-green-600",
      format: formatCurrency,
    },
    {
      key: "avg_price",
      label: t("avgPrice"),
      icon: DollarSign,
      color: "text-purple-600",
      format: formatCurrency,
    },
    {
      key: "unique_buyers",
      label: t("buyers"),
      icon: Users,
      color: "text-orange-600",
      format: formatNumber,
    },
  ];

  return (
    <Card>
      <CardHeader>
        <CardTitle>{t("title")}</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead>
              <tr className="border-b">
                <th className="text-left py-3 px-4 font-medium text-sm">
                  {t("product")}
                </th>
                {metrics.map((metric) => {
                  const Icon = metric.icon;
                  return (
                    <th key={metric.key} className="text-right py-3 px-4 font-medium text-sm">
                      <div className="flex items-center justify-end gap-2">
                        <Icon className={cn("h-4 w-4", metric.color)} />
                        {metric.label}
                      </div>
                    </th>
                  );
                })}
                <th className="text-right py-3 px-4 font-medium text-sm">
                  <div className="flex items-center justify-end gap-2">
                    <TrendingUp className="h-4 w-4 text-muted-foreground" />
                    {t("growth")}
                  </div>
                </th>
              </tr>
            </thead>
            <tbody>
              {products.map((product, index) => (
                <tr
                  key={product.product_id}
                  className={cn(
                    "border-b hover:bg-muted/50 transition-colors",
                    index === 0 && "bg-green-50/50 dark:bg-green-950/10"
                  )}
                >
                  <td className="py-3 px-4">
                    <div>
                      <div className="font-medium">{product.product_name}</div>
                      <div className="text-xs text-muted-foreground">
                        {product.product_sku}
                      </div>
                    </div>
                  </td>
                  {metrics.map((metric) => (
                    <td key={metric.key} className="py-3 px-4 text-right font-medium">
                      {metric.format((product as any)[metric.key])}
                    </td>
                  ))}
                  <td className="py-3 px-4">
                    <div
                      className={cn(
                        "flex items-center justify-end gap-1 font-medium",
                        product.growth_rate > 0 ? "text-green-600" : "text-red-600"
                      )}
                    >
                      <TrendingUp
                        className={cn(
                          "h-4 w-4",
                          product.growth_rate < 0 && "rotate-180"
                        )}
                      />
                      {Math.abs(product.growth_rate).toFixed(1)}%
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>
  );
}
