"use client";

import React, { useState, useMemo } from "react";
import { useTranslations } from "next-intl";
import { AuthGuard } from "@/features/auth/components/auth-guard";
import { PermissionGuard } from "@/features/auth/components/permission-guard";
import { useAuthStore } from "@/features/auth/stores/useAuthStore";
import { PageMotion } from "@/components/motion";
import { RoleBasedDashboard } from "@/features/dashboard/components/role-based-dashboard";
import { DateFilterControls } from "@/components/common/date-filter-controls";
import { startOfYear, endOfYear, format } from "date-fns";
import type { DateRange } from "react-day-picker";

function DashboardHeader() {
  const t = useTranslations("dashboard");
  const { user } = useAuthStore();

  return (
    <div className="flex items-center justify-between">
      <div>
        <h1 className="text-3xl font-medium tracking-tight">{t("title")}</h1>
        <p className="text-muted-foreground mt-1">
          {t("welcomeBack", { name: user?.name ?? "" })}
        </p>
      </div>
    </div>
  );
}

function DashboardContent() {
  return (
    <PageMotion className="space-y-6">
      <DashboardHeader />
      <RoleBasedDashboard />
    </PageMotion>
  );
}

export default function DashboardPage() {
  const t = useTranslations("dashboard");
  const { user } = useAuthStore();
  
  // Filter State
  const [filterMode, setFilterMode] = useState<"year" | "range">("year");
  const [selectedYear, setSelectedYear] = useState<number>(new Date().getFullYear());
  const [dateRange, setDateRange] = useState<DateRange | undefined>(() => {
    const start = startOfYear(new Date());
    const end = endOfYear(new Date());
    return { from: start, to: end };
  });

  // Calculate generic start/end dates
  const { startDate, endDate } = useMemo<{ startDate: string | undefined; endDate: string | undefined }>(() => {
    if (filterMode === "year") {
      const start = new Date(selectedYear, 0, 1);
      const end = new Date(selectedYear, 11, 31);
      return {
        startDate: format(start, "yyyy-MM-dd"),
        endDate: format(end, "yyyy-MM-dd"),
      };
    } else {
      // Range mode
      if (!dateRange?.from) return { startDate: undefined, endDate: undefined };
      
      const startStr = format(dateRange.from, "yyyy-MM-dd");
      let endStr = undefined;
      if (dateRange.to) {
        endStr = format(dateRange.to, "yyyy-MM-dd");
      }
      return { startDate: startStr, endDate: endStr };
    }
  }, [filterMode, selectedYear, dateRange]);

  return (
    <AuthGuard>
      <PermissionGuard requiredPermission="dashboard.view">
        <PageMotion className="space-y-6">
          <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
            <div>
              <h1 className="text-3xl font-medium tracking-tight">{t("title")}</h1>
              <p className="text-muted-foreground mt-1">
                {t("welcomeBack", { name: user?.name ?? "" })}
              </p>
            </div>
            <DateFilterControls
              filterMode={filterMode}
              onFilterModeChange={setFilterMode}
              selectedYear={selectedYear}
              onYearChange={setSelectedYear}
              dateRange={dateRange}
              onDateRangeChange={setDateRange}
            />
          </div>
          <RoleBasedDashboard startDate={startDate} endDate={endDate} />
        </PageMotion>
      </PermissionGuard>
    </AuthGuard>
  );
}
