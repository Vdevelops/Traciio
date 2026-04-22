"use client";

import { Suspense } from "react";
import { useTranslations } from "next-intl";
import { AuthGuard } from "@/features/auth/components/auth-guard";
import { PermissionGuard } from "@/features/auth/components/permission-guard";
import { PageMotion } from "@/components/motion";
import { ScheduleManagement } from "@/features/sales-crm/schedule-management/components/schedule-management";
import { Skeleton } from "@/components/ui/skeleton";

function SchedulesHeader() {
  const t = useTranslations("scheduleManagement.page");

  return (
    <div className="flex items-center justify-between">
      <div>
        <h1 className="text-3xl font-medium tracking-tight">{t("title")}</h1>
        <p className="text-muted-foreground mt-1">{t("description")}</p>
      </div>
    </div>
  );
}

function SchedulesPageContent() {
  return (
    <PageMotion className="space-y-6">
      <SchedulesHeader />

      <Suspense fallback={<Skeleton className="h-[600px] w-full" />}>
        <ScheduleManagement />
      </Suspense>
    </PageMotion>
  );
}

export default function SchedulesPage() {
  return (
    <AuthGuard>
      <PermissionGuard requiredPermission="schedules.view">
        <SchedulesPageContent />
      </PermissionGuard>
    </AuthGuard>
  );
}
