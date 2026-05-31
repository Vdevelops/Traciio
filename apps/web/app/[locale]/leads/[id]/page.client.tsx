"use client";

import { AuthGuard } from "@/features/auth/components/auth-guard";
import { PermissionGuard } from "@/features/auth/components/permission-guard";
import { LeadDetailShell } from "@/features/sales-crm/lead-management/components/lead-detail-shell";

export default function LeadDetailPageClient({ id }: { id: string }) {
  return (
    <AuthGuard>
      <PermissionGuard requiredPermission="leads.view">
        <LeadDetailShell leadId={id} />
      </PermissionGuard>
    </AuthGuard>
  );
}
