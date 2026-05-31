"use client";

import { useMemo, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { useTranslations } from "next-intl";
import {
  ArrowLeft,
  ArrowRightLeft,
  Calendar,
  FileText,
  Mail,
  MapPinned,
  Phone,
  Trash2,
  User,
} from "lucide-react";
import { toast } from "sonner";

import { AuthGuard } from "@/features/auth/components/auth-guard";
import { PermissionGuard } from "@/features/auth/components/permission-guard";
import { PageMotion } from "@/components/motion";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Skeleton } from "@/components/ui/skeleton";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import { DealDetailTabs } from "@/features/sales-crm/pipeline-management/components/DealDetailTabs";
import { DealForm } from "@/features/sales-crm/pipeline-management/components/deal-form";
import { MoveStageModal } from "@/features/sales-crm/pipeline-management/components/move-stage-modal";
import { useDeal, useDeleteDeal, useUpdateDeal, useDealActivities, useDealVisitReports } from "@/features/sales-crm/pipeline-management/hooks/useDeals";
import { useStages } from "@/features/sales-crm/pipeline-management/hooks/useStages";
import { useLead } from "@/features/sales-crm/lead-management/hooks/useLeads";
import { useLeadQualification } from "@/features/sales-crm/lead-management/hooks/useLeadQualification";
import { useTasks } from "@/features/sales-crm/task-management/hooks/useTasks";
import type { Task } from "@/features/sales-crm/task-management/types";
import { formatCurrency, formatEmailToMailto, formatPhoneNumberToWA } from "@/lib/utils";

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
  const { data: activitiesData } = useDealActivities(dealId);
  const { data: visitReportsData } = useDealVisitReports(dealId);
  const { data: dealTasksResponse } = useTasks({ deal_id: dealId, per_page: 20 });
  const { data: leadTasksResponse } = useTasks({ lead_id: deal?.lead_id, per_page: 20 });
  const updateDeal = useUpdateDeal();
  const deleteDeal = useDeleteDeal();

  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [isMoveStageOpen, setIsMoveStageOpen] = useState(false);

  const lead = leadData?.data;
  const activities = Array.isArray(activitiesData) ? activitiesData : [];
  const visitReports = Array.isArray(visitReportsData) ? visitReportsData : [];
  const dealTasks = Array.isArray(dealTasksResponse?.data) ? dealTasksResponse.data : [];
  const leadTasks = Array.isArray(leadTasksResponse?.data) ? leadTasksResponse.data : [];
  const tasks = dedupeTasks([...dealTasks, ...leadTasks]);
  const productInterests = qualification?.need_target_products ?? deal?.qualification_snapshot?.need_target_products ?? [];
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
  const sourceText = lead?.lead_source ?? deal?.source ?? "—";
  const companyName = lead?.company_name ?? deal?.account?.name ?? "N/A";
  const contactName =
    lead?.contact?.name ??
    deal?.contact?.name ??
    ([lead?.first_name, lead?.last_name].filter(Boolean).join(" ") || "—");
  const contactEmail = lead?.contact?.email ?? deal?.contact?.email ?? lead?.email ?? "";
  const contactPhone = lead?.contact?.phone ?? deal?.contact?.phone ?? lead?.phone ?? "";
  const address = [lead?.address, lead?.city, lead?.province, lead?.postal_code, lead?.country]
    .filter(Boolean)
    .join(", ");

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

        <section className="grid grid-cols-1 items-start gap-4 md:grid-cols-[minmax(0,1.65fr)_minmax(300px,0.95fr)] xl:grid-cols-[minmax(0,1.78fr)_minmax(340px,0.88fr)] 2xl:grid-cols-[minmax(0,1.95fr)_minmax(360px,0.82fr)]">
          <div className="min-w-0">
            <DealDetailTabs dealId={dealId} />
          </div>

          <aside className="min-w-0 space-y-3 self-start md:sticky md:top-4">
            <SidebarCard title="CUSTOMER INFORMATION">
              <div className="flex items-start gap-3">
                <Avatar className="h-10 w-10 border border-slate-200">
                  <AvatarImage alt={companyName} />
                  <AvatarFallback className="bg-slate-100 text-sm font-semibold text-slate-700">
                    {getInitials(companyName)}
                  </AvatarFallback>
                </Avatar>
                <div className="min-w-0">
                  <div className="truncate font-medium">{companyName}</div>
                  <div className="mt-1 text-xs text-slate-500">{deal.description || "Deal workspace"}</div>
                </div>
              </div>
            </SidebarCard>

            <SidebarCard title="CONTACT">
              <div className="space-y-2">
                <div className="flex items-center gap-2 text-sm">
                  <User className="h-4 w-4 text-slate-400" />
                  <span className="truncate">{contactName}</span>
                </div>
                {contactEmail && (
                  <a href={formatEmailToMailto(contactEmail)} className="flex cursor-pointer items-center gap-2 text-sm text-blue-600 hover:underline">
                    <Mail className="h-4 w-4" />
                    <span className="truncate">{contactEmail}</span>
                  </a>
                )}
                {contactPhone && (
                  <a href={formatPhoneNumberToWA(contactPhone)} className="flex cursor-pointer items-center gap-2 text-sm text-blue-600 hover:underline">
                    <Phone className="h-4 w-4" />
                    <span>{contactPhone}</span>
                  </a>
                )}
              </div>
            </SidebarCard>

            <SidebarCard title="LOCATION">
              <div className="rounded-xl border border-dashed border-slate-200 bg-slate-50 px-4 py-6 text-center">
                <div className="mx-auto flex h-10 w-10 items-center justify-center rounded-full bg-white text-slate-400 shadow-sm ring-1 ring-slate-200">
                  <MapPinned className="h-5 w-5" />
                </div>
                <div className="mt-3 text-sm text-slate-500">{address ? "Location set" : "No location set"}</div>
              </div>
              <div className="mt-3 text-sm leading-relaxed text-slate-500">{address || "Address not available"}</div>
            </SidebarCard>

            <SidebarCard title="ASSIGNED TO">
              {deal.assigned_user ? (
                <div className="flex items-center gap-3">
                  <Avatar className="h-9 w-9">
                    <AvatarImage src={deal.assigned_user.avatar_url} alt={deal.assigned_user.name} />
                    <AvatarFallback className="bg-muted text-xs font-semibold">
                      {getInitials(deal.assigned_user.name)}
                    </AvatarFallback>
                  </Avatar>
                  <div>
                    <div className="font-medium">{deal.assigned_user.name}</div>
                    <div className="text-xs text-slate-500">{deal.assigned_user.email}</div>
                  </div>
                </div>
              ) : (
                <div className="text-sm text-slate-500">Unassigned</div>
              )}
            </SidebarCard>

            <SidebarCard title="DEAL BANT SUMMARY">
              <div className="space-y-2 text-sm text-slate-500">
                <BantSideRow label="Budget" value={budgetValueText(qualification?.budget_target_amount ?? deal.qualification_snapshot?.budget_target_amount)} />
                <BantSideRow label="Authority" value={qualification?.authority_target_person ?? deal.qualification_snapshot?.authority_target_person ?? "-"} />
                <BantSideRow label="Need" value={qualification?.need_notes ?? deal.qualification_snapshot?.need_notes ?? "-"} />
                <BantSideRow label="Timeline" value={formatSafeDate(qualification?.timeline_target_date ?? deal.qualification_snapshot?.timeline_target_date)} />
              </div>
            </SidebarCard>

            <SidebarCard title="SOURCE DATA">
              <div className="space-y-1">
                <div className="font-medium">{sourceText}</div>
                {lead?.id ? <div className="text-xs text-slate-500">Lead: {lead.id}</div> : null}
                <Badge variant="outline" className="mt-1 capitalize">
                  {deal.status}
                </Badge>
              </div>
            </SidebarCard>

            <SidebarCard title="DATES">
              <div className="space-y-2 text-sm">
                <DateRow label="Created" value={formatSafeDate(deal.created_at)} />
                <DateRow label="Updated" value={formatSafeDate(deal.updated_at)} />
                <DateRow label="Expected Close" value={formatSafeDate(deal.expected_close_date)} />
                <DateRow label="Actual Close" value={formatSafeDate(deal.actual_close_date)} />
              </div>
            </SidebarCard>

            <SidebarCard title="ACTIVITY SNAPSHOT">
              <div className="grid grid-cols-2 gap-3">
                <MiniMetric label="Visits" value={String(visitReports.length)} />
                <MiniMetric label="Activities" value={String(activities.length)} />
                <MiniMetric label="Tasks" value={String(tasks.length)} />
                <MiniMetric label="Products" value={String(productInterests.length)} />
              </div>
            </SidebarCard>
          </aside>
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

function SidebarCard({
  title,
  children,
}: {
  readonly title: string;
  readonly children: React.ReactNode;
}) {
  return (
    <Card className="border-slate-200 bg-white shadow-sm">
      <CardHeader className="space-y-0 pb-3">
        <CardTitle className="text-xs font-semibold tracking-[0.12em] text-slate-500">{title}</CardTitle>
      </CardHeader>
      <CardContent className="pt-0">{children}</CardContent>
    </Card>
  );
}

function DateRow({ label, value }: { readonly label: string; readonly value: string }) {
  return (
    <div className="flex items-center justify-between gap-3">
      <span className="text-slate-500">{label}</span>
      <span className="text-right text-slate-900">{value}</span>
    </div>
  );
}

function BantSideRow({ label, value }: { readonly label: string; readonly value: string }) {
  return (
    <div className="flex items-start gap-2 rounded-lg bg-slate-50 px-3 py-2">
      <span className="mt-1 h-2 w-2 rounded-full bg-emerald-500" />
      <div className="min-w-0">
        <div className="text-[11px] uppercase tracking-wide text-slate-500">{label}</div>
        <div className="mt-0.5 break-words text-sm text-slate-900">{value}</div>
      </div>
    </div>
  );
}

function MiniMetric({ label, value }: { readonly label: string; readonly value: string }) {
  return (
    <div className="rounded-lg border border-slate-200 bg-slate-50 p-3">
      <div className="text-[11px] uppercase tracking-wide text-slate-500">{label}</div>
      <div className="mt-1 text-lg font-semibold text-slate-900">{value}</div>
    </div>
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
        <div className="grid grid-cols-1 gap-4 md:grid-cols-[minmax(0,1.65fr)_minmax(300px,0.95fr)]">
          <Skeleton className="h-[720px] rounded-2xl bg-slate-200" />
          <div className="space-y-3">
            {Array.from({ length: 7 }).map((_, index) => (
              <Skeleton key={index} className="h-24 rounded-2xl bg-slate-200" />
            ))}
          </div>
        </div>
      </div>
    </PageMotion>
  );
}

function dedupeTasks(tasks: Task[]): Task[] {
  const map = new Map<string, Task>();
  for (const task of tasks) {
    if (!task?.id) continue;
    map.set(task.id, task);
  }
  return Array.from(map.values()).sort((a, b) => {
    const left = new Date(b.created_at).getTime();
    const right = new Date(a.created_at).getTime();
    return left - right;
  });
}

function getInitials(value: string) {
  return (
    value
      .split(/\s+/)
      .filter(Boolean)
      .slice(0, 2)
      .map((part) => part[0]?.toUpperCase() || "")
      .join("") || "D"
  );
}

function formatSafeDate(value?: string | null) {
  if (!value) return "-";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "-";
  return date.toLocaleDateString("id-ID", {
    year: "numeric",
    month: "long",
    day: "numeric",
  });
}

function budgetValueText(value?: number) {
  return value ? formatCurrency(value) : "-";
}
