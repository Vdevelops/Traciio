"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { TargetMatrix } from "./target-matrix";
import { Input } from "@/components/ui/input";
import { formatCurrency } from "@/lib/utils";
import { useInfiniteMonthlyTargets } from "../hooks/useMonthlyTargets";
import { Calendar } from "lucide-react";
import { useDebounce } from "@/hooks/use-debounce";
import { useAuthStore } from "@/features/auth/stores/useAuthStore";

export function TeamTargetPlanner() {
  const t = useTranslations("monthlyTargetManagement.planner");
  const [year, setYear] = useState<number>(new Date().getFullYear());
  
  const currentYear = new Date().getFullYear();
  const years = Array.from({ length: 5 }, (_, i) => currentYear - 1 + i);
  const rowsPerLoad = 20;
  const monthsPerYear = 12;
  const perPage = rowsPerLoad * monthsPerYear;
  const [search, setSearch] = useState("");
  const debouncedSearch = useDebounce(search, 500);
  const { user } = useAuthStore();

  // If user is a sales manager, limit results to their scope on server
  const managerIdParam = user?.role === "sales_manager" ? user.id : undefined;

  const {
    data,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
  } = useInfiniteMonthlyTargets({
    year,
    per_page: perPage,
    search: debouncedSearch || undefined,
    manager_id: managerIdParam,
    scope: "user",
  });

  const targets = data?.pages.flatMap((page) => page.data) || [];
  
  // Calculate summaries (based on loaded data)
  const serverTotalTarget = data?.pages?.[0]?.meta?.additional?.total_target_amount ?? 0;
  const totalTarget = Number.isFinite(serverTotalTarget)
    ? serverTotalTarget
    : targets.reduce((sum, t) => sum + (t.target_amount || 0), 0);
  const uniqueUsers = new Set(targets.map(t => t.user_id).filter(Boolean)).size;

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div className="flex flex-col md:flex-row md:items-center gap-4">

          <div className="flex items-center gap-3">
            <Input
              placeholder={t("searchPlaceholder")}
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="w-[220px] h-9"
            />

            <Select
              value={year.toString()}
              onValueChange={(value) => setYear(Number.parseInt(value, 10))}
            >
              <SelectTrigger className="w-[120px] h-9">
                  <Calendar className="mr-2 h-4 w-4 text-muted-foreground" />
                  <SelectValue placeholder={t("selectYear")} />
                </SelectTrigger>
              <SelectContent>
                {years.map((y) => (
                  <SelectItem key={y} value={y.toString()}>
                    {y}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>

        <div className="flex flex-wrap gap-4">
          <div className="flex flex-col items-end">
            <span className="text-sm text-muted-foreground">{t("totalTargetLoaded")}</span>
            <span className="text-xl font-bold">{formatCurrency(totalTarget)}</span>
          </div>
          <div className="flex flex-col items-end border-l pl-4">
            <span className="text-sm text-muted-foreground">{t("usersWithTargets")}</span>
            <span className="text-xl font-bold">{uniqueUsers}</span>
          </div>
        </div>
      </div>

      <TargetMatrix
        initialYear={year}
        showHeader={true}
        data={targets}
        isLoading={isLoading}
        onLoadMore={() => fetchNextPage()}
        hasMore={hasNextPage}
        isLoadingMore={isFetchingNextPage}
        searchQuery={debouncedSearch}
      />
    </div>
  );
}
