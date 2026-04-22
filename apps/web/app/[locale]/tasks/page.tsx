import { Suspense } from "react";
import { getTranslations } from "next-intl/server";
import { AuthGuard } from "@/features/auth/components/auth-guard";
import { PermissionGuard } from "@/features/auth/components/permission-guard";
import { PageMotion } from "@/components/motion";
import dynamic from "next/dynamic";
import { Skeleton } from "@/components/ui/skeleton";

export const metadata = {
  title: "Tasks | Tracio",
};

// Dynamic import untuk TasksPageClient (client wrapper dengan code splitting)
const TasksPageClient = dynamic(
  () => import("./tasks-page-client").then((mod) => ({ default: mod.TasksPageClient })),
  {
    loading: () => null, // Use route-level loading.tsx
  },
);

async function TasksHeader() {
  const t = await getTranslations("taskManagement.page");

  return (
    <div className="flex items-center justify-between">
      <div>
        <h1 className="text-3xl font-medium tracking-tight">{t("title")}</h1>
        <p className="text-muted-foreground mt-1">{t("description")}</p>
      </div>
    </div>
  );
}

export default function TasksPage() {
  return (
    <AuthGuard>
      <PermissionGuard requiredPermission="tasks.view">
        <PageMotion className="space-y-6">
          <Suspense fallback={<Skeleton className="h-9 w-48" />}>
            <TasksHeader />
          </Suspense>

          <Suspense fallback={<Skeleton className="h-[600px] w-full" />}>
            <TasksPageClient />
          </Suspense>
        </PageMotion>
      </PermissionGuard>
    </AuthGuard>
  );
}


