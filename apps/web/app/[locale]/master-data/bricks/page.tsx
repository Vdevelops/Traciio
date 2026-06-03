import { Suspense } from "react";
import { AuthGuard } from "@/features/auth/components/auth-guard";
import { PermissionGuard } from "@/features/auth/components/permission-guard";
import { BrickManagement } from "@/features/master-data/brick/components/brick-management";

export const metadata = {
  title: "Bricks | Tracio",
};

export default function BricksPage() {
  return (
    <AuthGuard>
      <PermissionGuard requiredPermission="bricks.view">
        <div className="absolute inset-0 overflow-hidden">
          <Suspense fallback={null}>
            <BrickManagement />
          </Suspense>
        </div>
      </PermissionGuard>
    </AuthGuard>
  );
}
