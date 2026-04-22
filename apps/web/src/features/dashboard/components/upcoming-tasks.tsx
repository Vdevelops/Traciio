"use client";

import { useTranslations, useLocale } from "next-intl";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useDashboardOverview } from "../hooks/useDashboard";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
import { Clock, CheckSquare } from "lucide-react";

import { useHasPermission } from "@/features/auth/providers/permissions-provider";
import { Lock } from "lucide-react";

const priorityVariant: Record<string, "default" | "secondary" | "destructive" | "outline"> = {
  low: "secondary",
  medium: "default",
  high: "destructive",
  urgent: "destructive",
};

export function UpcomingTasks({ startDate, endDate }: Readonly<{ startDate?: string; endDate?: string }>) {
  const t = useTranslations("dashboardOverview");
  const locale = useLocale();
  
  // Permission check
  const hasPermission = useHasPermission("tasks.view");

  const params = { period: "month" as const, start_date: startDate, end_date: endDate };

  const { data, isLoading } = useDashboardOverview(params, { enabled: hasPermission });

  console.log("UpcomingTasks Params:", params);
  console.log("UpcomingTasks Data:", data);

  // Show warning if no permission
  if (!hasPermission) {
    return (
      <Card className="border-0 shadow-sm bg-muted/10 h-full">
        <CardHeader className="px-3 sm:px-6 pb-2">
          <CardTitle className="text-xs sm:text-sm font-medium flex items-center gap-2 text-muted-foreground">
            <Lock className="h-3.5 w-3.5" />
            {t("upcomingTasks.title")}
          </CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col items-center justify-center p-6 text-center h-[200px]">
          <div className="bg-muted/20 p-3 rounded-full mb-3">
            <Lock className="h-5 w-5 text-muted-foreground" />
          </div>
          <p className="text-xs font-medium text-muted-foreground mb-1">
            Access Restricted
          </p>
          <p className="text-[10px] text-muted-foreground/70 max-w-[150px]">
            You don't have permission to view tasks.
          </p>
        </CardContent>
      </Card>
    );
  }

  if (isLoading) {
    return (
      <Card className="border-0 shadow-sm h-full">
        <CardHeader className="px-3 sm:px-6">
          <CardTitle className="text-xs sm:text-sm">{t("upcomingTasks.title")}</CardTitle>
        </CardHeader>
        <CardContent className="px-3 sm:px-6 pb-3 sm:pb-6">
          <div className="space-y-2 sm:space-y-3">
            {[...Array(3)].map((_, i) => (
              <Skeleton key={i} className="h-10 sm:h-12 w-full" />
            ))}
          </div>
        </CardContent>
      </Card>
    );
  }

  const overview = data?.data;
  const tasks = Array.isArray(overview?.upcoming_tasks) ? overview.upcoming_tasks : [];

  if (!tasks || tasks.length === 0) {
    return (
      <Card className="border-0 shadow-sm">
        <CardHeader className="px-3 sm:px-6">
          <CardTitle className="text-xs sm:text-sm">{t("upcomingTasks.title")}</CardTitle>
        </CardHeader>
        <CardContent className="relative min-h-[150px] sm:min-h-[200px] px-3 sm:px-6 pb-3 sm:pb-6">
          <div className="absolute inset-0 flex flex-col items-center justify-center">
            <CheckSquare className="h-8 w-8 sm:h-12 sm:w-12 text-muted-foreground/50 mb-2 sm:mb-3" />
            <p className="text-xs sm:text-sm text-muted-foreground">
              {t("upcomingTasks.empty")}
            </p>
          </div>
        </CardContent>
      </Card>
    );
  }

  const formatDueDate = (value?: string | null) => {
    if (!value) return t("upcomingTasks.noDueDate");
    const date = new Date(value);
    if (isNaN(date.getTime())) return t("upcomingTasks.invalidDate");
    return date.toLocaleDateString(locale, {
      day: "2-digit",
      month: "short",
    });
  };

  return (
    <Card className="border-0 shadow-sm">
      <CardHeader className="px-3 sm:px-6">
        <CardTitle className="text-xs sm:text-sm font-medium">
          {t("upcomingTasks.title")}
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-1.5 sm:space-y-2 px-3 sm:px-6 pb-3 sm:pb-6">
        {tasks.slice(0, 4).map((task) => {
          if (!task) return null;
          return (
            <div
              key={task.id}
              className="flex items-center justify-between gap-2 sm:gap-3 rounded-md border px-2 sm:px-3 py-1.5 sm:py-2 text-xs sm:text-sm cursor-pointer hover:bg-muted/50 transition-colors"
            >
              <div className="min-w-0 flex-1">
                <div className="truncate font-medium text-xs sm:text-sm">{task.title ?? t("upcomingTasks.noTitle")}</div>
                <div className="mt-0.5 flex items-center gap-1 sm:gap-2 text-[10px] sm:text-xs text-muted-foreground">
                  <Clock className="h-2.5 w-2.5 sm:h-3 sm:w-3" />
                  <span>{formatDueDate(task.due_date)}</span>
                </div>
              </div>
              <Badge
                variant={priorityVariant[task.priority ?? ""] ?? "outline"}
                className="text-[10px] sm:text-xs capitalize shrink-0"
              >
                {task.priority ?? "unknown"}
              </Badge>
            </div>
          );
        })}
      </CardContent>
    </Card>
  );
}


