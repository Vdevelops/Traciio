import { Suspense } from "react";
import { PageMotion } from "@/components/motion";
import { AuthGuard } from "@/features/auth/components/auth-guard";
import { PermissionGuard } from "@/features/auth/components/permission-guard";
import { GroupsPageClient } from "./groups-page-client";

export const metadata = {
  title: "Groups | Tracio",
};

export default function GroupsPage() {
  return (
    <AuthGuard>
      <PermissionGuard requiredPermission="groups.view">
        <PageMotion className=" ">
          <Suspense fallback={null}>
            <GroupsPageClient />
          </Suspense>
        </PageMotion>
      </PermissionGuard>
    </AuthGuard>
  );
}

