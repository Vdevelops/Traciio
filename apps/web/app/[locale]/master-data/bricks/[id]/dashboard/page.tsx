import { Suspense } from "react";
import { AuthGuard } from "@/features/auth/components/auth-guard";
import { PermissionGuard } from "@/features/auth/components/permission-guard";
import { PageMotion } from "@/components/motion";
import { Skeleton } from "@/components/ui/skeleton";
import dynamic from "next/dynamic";

const BrickDashboardClient = dynamic(
  () => import("./brick-dashboard-client").then((mod) => ({ default: mod.BrickDashboardClient })),
  { loading: () => null }
);

export const metadata = {
  title: "Brick Dashboard | Salesview",
};

export default function BrickDashboardPage() {
  return (
    <AuthGuard>
      <PermissionGuard requiredPermission="bricks.view">
        <PageMotion>
          <Suspense fallback={<BrickDashboardSkeleton />}>
            <BrickDashboardClient />
          </Suspense>
        </PageMotion>
      </PermissionGuard>
    </AuthGuard>
  );
}

function BrickDashboardSkeleton() {
  return (
    <div className="space-y-6">
      <Skeleton className="h-8 w-64" />
      <div className="grid gagrid-cols-2 md:grid-cols-4">
        {Array.from({ length: 8 }).map((_, i) => (
          <Skeleton key={`skeleton-${i}`} className="h-32" />
        ))}
      </div>
    </div>
  );
}

