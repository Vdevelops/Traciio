"use client";

import { Suspense } from "react";
import { LayoutDashboard, Settings, Table, ClipboardList } from "lucide-react";
import { AuthGuard } from "@/features/auth/components/auth-guard";
import { PermissionGuard } from "@/features/auth/components/permission-guard";
import { PageMotion } from "@/components/motion";
import { KanbanBoard } from "@/features/sales-crm/pipeline-management/components/kanban-board";
import { PipelineTableView } from "@/features/sales-crm/pipeline-management/components/pipeline-table-view";
import { StagesManagement } from "@/features/sales-crm/pipeline-management/components/stages-management";
import { ActivityTypeList } from "@/features/sales-crm/visit-report/components/activity-type-list";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useTranslations } from "next-intl";
import { useHasPermission } from "@/features/master-data/user-management/hooks/useHasPermission";
import { Skeleton } from "@/components/ui/skeleton";
import { useRouter } from "@/i18n/routing";

function PipelineHeader() {
  const t = useTranslations("pipelineManagement.page");

  return (
    <div>
      <h1 className="text-3xl font-bold tracking-tight">{t("title")}</h1>
      <p className="text-muted-foreground mt-1">{t("description")}</p>
    </div>
  );
}

function PipelinePageContent() {
  const t = useTranslations("pipelineManagement.page");
  const hasStagesPermission = useHasPermission("pipeline.stages-view");
  const hasActivityTypesPermission = useHasPermission("visit-reports.activity-type");
  const router = useRouter();

  const handleDealClick = (deal: { id: string }) => {
    router.push(`/deals/${deal.id}`);
  };

  return (
    <PageMotion className="space-y-6 p-4">
      <PipelineHeader />

      <Tabs defaultValue="kanban" className="w-full">
        <TabsList>
          <TabsTrigger value="kanban" className="gap-2">
            <LayoutDashboard className="h-4 w-4" />
            {t("tabKanban")}
          </TabsTrigger>
          <TabsTrigger value="table" className="gap-2">
            <Table className="h-4 w-4" />
            {t("tabTable")}
          </TabsTrigger>
          {hasStagesPermission && (
            <TabsTrigger value="stages" className="gap-2">
              <Settings className="h-4 w-4" />
              {t("tabStages")}
            </TabsTrigger>
          )}
          {hasActivityTypesPermission && (
            <TabsTrigger value="activity-types" className="gap-2">
              <ClipboardList className="h-4 w-4" />
              {t("tabActivityTypes")}
            </TabsTrigger>
          )}
        </TabsList>

        <TabsContent value="kanban" className="mt-6">
          <Suspense fallback={<Skeleton className="h-[600px] w-full" />}>
            <KanbanBoard onDealClick={handleDealClick} />
          </Suspense>
        </TabsContent>

        <TabsContent value="table" className="mt-6">
          <Suspense fallback={<Skeleton className="h-[600px] w-full" />}>
            <PipelineTableView onDealClick={handleDealClick} />
          </Suspense>
        </TabsContent>

        {hasStagesPermission && (
          <TabsContent value="stages" className="mt-6">
            <Suspense fallback={<Skeleton className="h-[600px] w-full" />}>
              <StagesManagement />
            </Suspense>
          </TabsContent>
        )}

        {hasActivityTypesPermission && (
          <TabsContent value="activity-types" className="mt-6">
            <Suspense fallback={<Skeleton className="h-[600px] w-full" />}>
              <ActivityTypeList />
            </Suspense>
          </TabsContent>
        )}
      </Tabs>
    </PageMotion>
  );
}

export default function PipelinePageClient() {
  return (
    <AuthGuard>
      <PermissionGuard requiredPermission="pipeline.view">
        <PipelinePageContent />
      </PermissionGuard>
    </AuthGuard>
  );
}
