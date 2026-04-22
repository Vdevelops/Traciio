"use client";

import { Suspense } from "react";
import { useTranslations } from "next-intl";
import { AuthGuard } from "@/features/auth/components/auth-guard";
import { PermissionGuard } from "@/features/auth/components/permission-guard";
import { PageMotion } from "@/components/motion";
import { LeadManagement } from "@/features/sales-crm/lead-management/components/lead-management";
import { Skeleton } from "@/components/ui/skeleton";

function LeadsHeader() {
  const t = useTranslations("leadManagement.page");

  return (
    <div className="flex items-center justify-between">
      <div>
        <h1 className="text-3xl font-medium tracking-tight">{t("title")}</h1>
        <p className="text-muted-foreground mt-1">{t("description")}</p>
      </div>
    </div>
  );
}

function LeadsPageContent() {
  return (
    <PageMotion className="space-y-6">
      <LeadsHeader />

      <Suspense fallback={<Skeleton className="h-[600px] w-full" />}>
        <LeadManagement />
      </Suspense>
    </PageMotion>
  );
}

export default function LeadsPage() {
  return (
    <AuthGuard>
      <PermissionGuard requiredPermission="leads.view">
        <LeadsPageContent />
      </PermissionGuard>
    </AuthGuard>
  );
}
