"use client";

import { useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { useTranslations } from "next-intl";
import {
  ArrowLeft,
  ArrowRightLeft,
  FileText,
  Trash2,
} from "lucide-react";
import { toast } from "sonner";

import { AuthGuard } from "@/features/auth/components/auth-guard";
import { PermissionGuard } from "@/features/auth/components/permission-guard";
import { PageMotion } from "@/components/motion";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Skeleton } from "@/components/ui/skeleton";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import { DealDetailTabs } from "@/features/sales-crm/pipeline-management/components/DealDetailTabs";
import { DealForm } from "@/features/sales-crm/pipeline-management/components/deal-form";
import { MoveStageModal } from "@/features/sales-crm/pipeline-management/components/move-stage-modal";
import { useDeal, useDeleteDeal, useUpdateDeal } from "@/features/sales-crm/pipeline-management/hooks/useDeals";
import { useStages } from "@/features/sales-crm/pipeline-management/hooks/useStages";
import { useLead } from "@/features/sales-crm/lead-management/hooks/useLeads";
import { useLeadQualification } from "@/features/sales-crm/lead-management/hooks/useLeadQualification";
import { formatCurrency } from "@/lib/utils";

function DealDetailPageContent() {
  const params = useParams();
  const router = useRouter();
  const dealId = params.id as string;
  const t = useTranslations("deals.detail");
  const tCommon = useTranslations("common");

  const { data: deal, isLoading, error } = useDeal(dealId);
  const { data: stages } = useStages();
  const { data: leadData } = useLead(deal?.lead_id ?? "");
  const { qualification } = useLeadQualification(deal?.lead_id ?? "");
  const updateDeal = useUpdateDeal();
  const deleteDeal = useDeleteDeal();

  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [isMoveStageOpen, setIsMoveStageOpen] = useState(false);

  const lead = leadData?.data;
  const bantCount = [
    qualification?.budget_confirmed ?? deal?.budget_confirmed,
    qualification?.authority_confirmed ?? deal?.authority_confirmed,
    qualification?.need_confirmed ?? deal?.need_confirmed,
    qualification?.timeline_confirmed ?? deal?.timeline_confirmed,
  ].filter(Boolean).length;

  const title = deal?.title || "Deal";
  const statusVariant = deal?.status === "won" ? "default" : deal?.status === "lost" ? "destructive" : "secondary";
  const stageName = deal?.stage?.name || "No Stage";
  const valueText = deal?.value_formatted ?? formatCurrency(deal?.value ?? 0);
  const probabilityText = `${deal?.probability ?? 0}%`;
  const companyName = lead?.company_name ?? deal?.account?.name ?? "N/A";

  const handleUpdate = async (formData: Parameters<typeof updateDeal.mutateAsync>[0]["data"]) => {
    try {
      await updateDeal.mutateAsync({ id: dealId, data: formData });
      setIsEditDialogOpen(false);
      toast.success(t("toast.updated"));
    } catch {
      // handled by interceptor
    }
  };

  const handleDelete = async () => {
    try {
      await deleteDeal.mutateAsync(dealId);
      setIsDeleteDialogOpen(false);
      toast.success(t("toast.deleted"));
      router.push("/pipeline");
    } catch {
      // handled by interceptor
    }
  };

  if (isLoading) {
    return <DealDetailSkeleton />;
  }

  if (error || !deal) {
    return (
      <PageMotion className="p-2 sm:p-4 lg:p-6">
        <Card className="border-border/70 bg-card/80">
          <CardContent className="py-10 text-center text-muted-foreground">
            Deal not found.
          </CardContent>
        </Card>
      </PageMotion>
    );
  }

  return (
    <PageMotion className="min-h-screen p-2 sm:p-4 lg:p-6">
      <div className="mx-auto max-w-[1680px] space-y-5 lg:space-y-6">
        <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
          <div className="flex min-w-0 items-start gap-3 sm:gap-4">
            <Button
              variant="ghost"
              size="icon"
              className="mt-0.5 shrink-0 rounded-xl"
              onClick={() => router.back()}
            >
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <div className="min-w-0">
              <div className="flex flex-wrap items-center gap-2">
                <h1 className="page-title max-w-[620px] truncate text-2xl font-semibold tracking-tight sm:text-3xl">
                  {title}
                </h1>
                <Badge variant={statusVariant} className="rounded-full capitalize text-[0.68rem]">
                  {deal.status}
                </Badge>
              </div>
              <div className="mt-1 flex flex-wrap items-center gap-2 text-sm text-muted-foreground">
                <span>{deal.id}</span>
                <span>• {companyName}</span>
                <span>• {stageName}</span>
              </div>
            </div>
          </div>

          <div className="flex flex-wrap items-center gap-2 lg:justify-end">
            {deal.stage_id && Array.isArray(stages) && stages.length > 0 && (
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => setIsMoveStageOpen(true)}
              >
                <ArrowRightLeft className="mr-2 h-4 w-4" />
                Move Stage
              </Button>
            )}
            <Button
              type="button"
              variant="default"
              size="sm"
              onClick={() => toast.info("Convert to Quotation / Sales Order flow will be implemented in subsequent phases")}
              className="bg-blue-600 text-white hover:bg-blue-700"
            >
              <FileText className="mr-2 h-4 w-4" />
              Convert to Quotation
            </Button>
            <Button type="button" variant="outline" size="sm" onClick={() => setIsEditDialogOpen(true)}>
              {tCommon("edit")}
            </Button>
            <Button type="button" variant="destructive" size="sm" onClick={() => setIsDeleteDialogOpen(true)}>
              <Trash2 className="mr-2 h-4 w-4" />
              {tCommon("delete")}
            </Button>
          </div>
        </div>

        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 xl:grid-cols-4">
          <MetricCard label="Value" value={valueText} />
          <MetricCard label="Probability (%)" value={probabilityText} />
          <MetricCard label="Pipeline Stage" value={stageName} />
          <MetricCard label="BANT Qualification" value={`${bantCount}/4`} accent="success" />
        </div>

        <section className="grid grid-cols-1 items-start gap-4">
          <div className="min-w-0">
            <DealDetailTabs dealId={dealId} />
          </div>
        </section>
      </div>

      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-[700px]">
          <DialogHeader>
            <DialogTitle>{tCommon("edit")}</DialogTitle>
          </DialogHeader>
          <DealForm
            deal={deal}
            onSubmit={handleUpdate}
            onCancel={() => setIsEditDialogOpen(false)}
            isLoading={updateDeal.isPending}
          />
        </DialogContent>
      </Dialog>

      <DeleteDialog
        open={isDeleteDialogOpen}
        onOpenChange={setIsDeleteDialogOpen}
        onConfirm={handleDelete}
        title={t("deleteDialog.title")}
        description={t("deleteDialog.description")}
        isLoading={deleteDeal.isPending}
      />

      {Array.isArray(stages) && stages.length > 0 && (
        <MoveStageModal
          dealId={deal.id}
          currentStageId={deal.stage_id}
          availableStages={stages}
          isOpen={isMoveStageOpen}
          onClose={() => setIsMoveStageOpen(false)}
        />
      )}
    </PageMotion>
  );
}

export default function DealDetailPage() {
  return (
    <AuthGuard>
      <PermissionGuard requiredPermission="pipeline.view">
        <DealDetailPageContent />
      </PermissionGuard>
    </AuthGuard>
  );
}

function MetricCard({
  label,
  value,
  accent,
}: {
  readonly label: string;
  readonly value: string;
  readonly accent?: "default" | "success";
}) {
  return (
    <Card className="border-slate-200 bg-white shadow-sm">
      <CardContent className="p-4 sm:p-5">
        <div className="text-xs text-slate-500">{label}</div>
        <div className={accent === "success" ? "mt-1 text-2xl font-semibold tracking-tight text-emerald-500" : "mt-1 text-2xl font-semibold tracking-tight text-slate-900"}>
          {value}
        </div>
      </CardContent>
    </Card>
  );
}

function DealDetailSkeleton() {
  return (
    <PageMotion className="min-h-screen p-2 sm:p-4 lg:p-6">
      <div className="space-y-4 lg:space-y-6">
        <div className="flex items-start gap-3">
          <Skeleton className="h-10 w-10 rounded-xl bg-slate-200" />
          <div className="space-y-2">
            <Skeleton className="h-8 w-64 bg-slate-200" />
            <Skeleton className="h-4 w-48 bg-slate-200" />
          </div>
        </div>
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 xl:grid-cols-4">
          {Array.from({ length: 4 }).map((_, index) => (
            <Skeleton key={index} className="h-20 rounded-xl bg-slate-200" />
          ))}
        </div>
        <div className="grid grid-cols-1 gap-4">
          <Skeleton className="h-[720px] rounded-2xl bg-slate-200" />
        </div>
      </div>
    </PageMotion>
  );
}
