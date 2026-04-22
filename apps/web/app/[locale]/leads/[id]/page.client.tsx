"use client";

import { Suspense } from "react";
import { useTranslations } from "next-intl";
import { AuthGuard } from "@/features/auth/components/auth-guard";
import { PermissionGuard } from "@/features/auth/components/permission-guard";
import { PageMotion } from "@/components/motion";
import { LeadDetailTabs } from "@/features/sales-crm/lead-management/components/LeadDetailTabs";
import { Skeleton } from "@/components/ui/skeleton";
import { PageDetailLayout } from "@/components/layouts/page-detail-layout";
import { useLead } from "@/features/sales-crm/lead-management/hooks/useLeads";
import { Badge } from "@/components/ui/badge";

export default function LeadDetailPageClient({ id }: { id: string }) {
  const t = useTranslations("leadManagement.page");
  const { data: leadData, isLoading } = useLead(id);

  const lead = leadData?.data;

  return (
    <AuthGuard>
      <PermissionGuard requiredPermission="leads.view">
        <PageMotion className="p-2 sm:p-4">
          <PageDetailLayout
            title={isLoading ? <Skeleton className="w-48 h-8" /> : (
              lead ? `${lead.first_name} ${lead.last_name || ""}` : "Lead Not Found"
            )}
            subtitle={
              <div className="flex items-center gap-2 mt-1">
                <span>{lead?.company_name || ""}</span>
                {lead?.lead_status && (
                  <Badge variant="outline" className="text-xs">
                    {lead.lead_status.toUpperCase()}
                  </Badge>
                )}
              </div>
            }
            backHref="/leads"
          >
            <Suspense fallback={<Skeleton className="h-[600px] w-full" />}>
              <LeadDetailTabs leadId={id} />
            </Suspense>
          </PageDetailLayout>
        </PageMotion>
      </PermissionGuard>
    </AuthGuard>
  );
}
