"use client";

import { Suspense } from "react";
import { useTranslations } from "next-intl";
import { AuthGuard } from "@/features/auth/components/auth-guard";
import { PermissionGuard } from "@/features/auth/components/permission-guard";
import { PageMotion } from "@/components/motion";
import { VisitReportManagement } from "@/features/sales-crm/visit-report/components/visit-report-management";
import { Skeleton } from "@/components/ui/skeleton";

function VisitReportsHeader() {
  const t = useTranslations("visitReportManagement.page");

  return (
    <div className="flex items-center justify-between">
      <div>
        <h1 className="text-3xl font-medium tracking-tight">{t("title")}</h1>
        <p className="text-muted-foreground mt-1">{t("description")}</p>
      </div>
    </div>
  );
}

function VisitReportsPageContent() {
  return (
    <PageMotion className="space-y-6">
      <VisitReportsHeader />

      <Suspense fallback={<Skeleton className="h-[600px] w-full" />}>
        <VisitReportManagement />
      </Suspense>
    </PageMotion>
  );
}

export default function VisitReportsPage() {
  return (
    <AuthGuard>
      <PermissionGuard requiredPermission="visit-reports.view">
        <VisitReportsPageContent />
      </PermissionGuard>
    </AuthGuard>
  );
}
