"use client";

import { useState, useEffect } from "react";
import { useTranslations } from "next-intl";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarImage } from "@/components/ui/avatar";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Skeleton } from "@/components/ui/skeleton";
import { Users, Search, BarChart3, Eye, Trophy, TrendingUp } from "lucide-react";
import { useBrickSales } from "../hooks/useBricks";
import { formatEmailToMailto } from "@/lib/utils";
import Link from "next/link";
import { useSalesPerformanceList } from "@/features/sales-overview/hooks/useSalesPerformanceList";
import type { SalesPerformanceListItem } from "@/features/sales-overview/types";

type BrickSalesUser = {
  id: string;
  name: string;
  email: string;
  status: string;
  monthly_target?: {
    target_amount: number;
    target_amount_formatted?: string;
  };
};

interface BrickSalesManagementProps {
  brickId: string;
  periodStart: string;
  periodEnd: string;
}

export function BrickSalesManagement({ brickId, periodStart, periodEnd }: BrickSalesManagementProps) {
  const t = useTranslations("brickSales");
  const [search, setSearch] = useState("");

  const { data, isLoading } = useBrickSales(brickId);

  const salesPerformanceHook = useSalesPerformanceList();

  // Sync sales performance range with the shared Brick period filter.
  useEffect(() => {
    salesPerformanceHook.setStartDate(periodStart);
    salesPerformanceHook.setEndDate(periodEnd);
  }, [periodStart, periodEnd, salesPerformanceHook]);

  const { performanceList } = salesPerformanceHook;

  // Create a map of user_id to performance data for quick lookup
  const performanceMap = new Map<string, SalesPerformanceListItem>(
    performanceList.map((perf: SalesPerformanceListItem) => [perf.user_id, perf])
  );

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Users className="h-5 w-5" />
            {t("title")}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-2">
            {Array.from({ length: 5 }).map((_, i) => (
              <Skeleton key={i} className="h-12 w-full" />
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  const allSales: BrickSalesUser[] = data?.data ?? [];

  // Filter by search
  const filteredSales = allSales.filter((sale) =>
    search === "" ||
    sale.name.toLowerCase().includes(search.toLowerCase()) ||
    sale.email.toLowerCase().includes(search.toLowerCase())
  );

  // Sort: active first by revenue, inactive (opacity-50) at the bottom
  const sortedSales = [...filteredSales].sort((a, b) => {
    const aActive = a.status === "active";
    const bActive = b.status === "active";
    if (aActive !== bActive) return aActive ? -1 : 1;

    const perfA = performanceMap.get(a.id);
    const perfB = performanceMap.get(b.id);
    return (perfB?.total_revenue ?? 0) - (perfA?.total_revenue ?? 0);
  });

  const getAvatarUrl = (user: BrickSalesUser) => {
    const perf = performanceMap.get(user.id);
    if (perf?.avatar_url) return perf.avatar_url;
    return `https://api.dicebear.com/7.x/lorelei/svg?seed=${encodeURIComponent(user.email)}`;
  };

  const getRankBadge = (index: number) => {
    if (index === 0) return <Trophy className="h-4 w-4 text-yellow-500" />;
    if (index === 1) return <Trophy className="h-4 w-4 text-gray-400" />;
    if (index === 2) return <Trophy className="h-4 w-4 text-amber-600" />;
    return null;
  };

  if (allSales.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Users className="h-5 w-5" />
            {t("title")}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground text-center py-8">
            {t("noSales")}
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Users className="h-5 w-5" />
            {t("title")} ({allSales.length})
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex gap-4 mb-6">
            <div className="flex-1 relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder={t("searchPlaceholder")}
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="pl-10"
              />
            </div>
          </div>

          {/* Sales Table */}
          <div className="border rounded-lg">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-[60px]">{t("rank")}</TableHead>
                  <TableHead>{t("name")}</TableHead>
                  <TableHead>{t("email")}</TableHead>
                  <TableHead>{t("monthlyTargetAchievement")}</TableHead>
                  <TableHead className="text-right">{t("revenue")}</TableHead>
                  <TableHead className="text-right">{t("deals")}</TableHead>
                  <TableHead className="text-right">{t("actions")}</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {sortedSales.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={7} className="text-center py-8 text-muted-foreground">
                      {t("noResults")}
                    </TableCell>
                  </TableRow>
                ) : (
                  sortedSales.map((sale, index: number) => {
                    const performance = performanceMap.get(sale.id);
                    const revenue = performance?.total_revenue ?? 0;
                    const dealsClosed = performance?.deals_closed ?? 0;
                    const conversionRate = performance?.conversion_rate ?? 0;
                    const isInactive = sale.status === "inactive";

                    // Determine monthly target: prefer BE monthly_target, fallback to performanceMap
                    const targetFormatted =
                      sale.monthly_target?.target_amount_formatted ??
                      performance?.target_amount_formatted ??
                      "-";
                    const revenueFormatted =
                      performance?.total_revenue_formatted ?? "-";
                    const achievementPct =
                      performance?.target_achievement_percentage ?? null;

                    return (
                      <TableRow
                        key={sale.id}
                        className={isInactive ? "opacity-50" : undefined}
                      >
                        <TableCell>
                          <div className="flex items-center gap-2">
                            {getRankBadge(index)}
                            <span className="text-sm font-medium text-muted-foreground">
                              #{index + 1}
                            </span>
                          </div>
                        </TableCell>
                        <TableCell>
                          <div className="flex items-center gap-3">
                            <Avatar className="h-8 w-8">
                              <AvatarImage src={getAvatarUrl(sale)} alt={sale.name} />
                            </Avatar>
                            <span className="font-medium">{sale.name}</span>
                          </div>
                        </TableCell>
                        <TableCell>
                          <a
                            href={formatEmailToMailto(sale.email)}
                            className="text-muted-foreground hover:text-primary hover:underline cursor-pointer min-w-0"
                          >
                            {sale.email}
                          </a>
                        </TableCell>
                        <TableCell>
                          <div className="flex flex-col gap-0.5">
                            <div className="text-sm">
                              <span className="text-muted-foreground">{targetFormatted}</span>
                              <span className="mx-1 text-muted-foreground">/</span>
                              <span className="font-medium">{revenueFormatted}</span>
                            </div>
                            {achievementPct !== null && (
                              <p
                                className={`text-[10px] font-medium ${
                                  achievementPct >= 100
                                    ? "text-green-600"
                                    : achievementPct >= 75
                                    ? "text-yellow-600"
                                    : "text-red-600"
                                }`}
                              >
                                {Math.round(achievementPct)}% {t("ofTarget")}
                              </p>
                            )}
                          </div>
                        </TableCell>
                        <TableCell className="text-right">
                          {performance ? (
                            <div className="flex flex-col items-end gap-1">
                              <div className="flex items-center gap-1">
                                <TrendingUp
                                  className={`h-3 w-3 ${
                                    conversionRate >= 50
                                      ? "text-green-500"
                                      : conversionRate >= 30
                                      ? "text-yellow-500"
                                      : "text-red-500"
                                  }`}
                                />
                                <span
                                  className={`text-sm font-medium ${
                                    conversionRate >= 50
                                      ? "text-green-600"
                                      : conversionRate >= 30
                                      ? "text-yellow-600"
                                      : "text-red-600"
                                  }`}
                                >
                                  {conversionRate.toFixed(1)}%
                                </span>
                              </div>
                              <span className="text-xs text-muted-foreground">
                                {performance.total_revenue_formatted ?? "Rp 0"}
                              </span>
                            </div>
                          ) : (
                            <span className="text-sm text-muted-foreground">-</span>
                          )}
                        </TableCell>
                        <TableCell className="text-right">
                          {performance ? (
                            <div className="flex flex-col items-end gap-1">
                              <span className="text-sm font-medium">
                                {dealsClosed} {t("dealsUnit")}
                              </span>
                              <span className="text-xs text-muted-foreground">
                                {conversionRate.toFixed(1)}% {t("conversion")}
                              </span>
                            </div>
                          ) : (
                            <span className="text-sm text-muted-foreground">-</span>
                          )}
                        </TableCell>
                        <TableCell className="text-right">
                          <div className="flex justify-end gap-2">
                            <Link href={`/sales-overview/sales-rep/${sale.id}`}>
                              <Button
                                variant="ghost"
                                size="icon"
                                title={t("viewPerformance")}
                                className="cursor-pointer"
                              >
                                <BarChart3 className="h-4 w-4" />
                              </Button>
                            </Link>
                            <Link href={`/sales-overview/sales-rep/${sale.id}`}>
                              <Button
                                variant="ghost"
                                size="icon"
                                title={t("viewDetails")}
                                className="cursor-pointer"
                              >
                                <Eye className="h-4 w-4" />
                              </Button>
                            </Link>
                          </div>
                        </TableCell>
                      </TableRow>
                    );
                  })
                )}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
