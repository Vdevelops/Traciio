"use client";

import { Edit, Trash2, Mail, MapPin, Phone, TrendingUp, UserPlus, Globe, FileText, Activity, Plus, CheckSquare, Package } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/avatar";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Drawer } from "@/components/ui/drawer";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useLead, useDeleteLead, useUpdateLead, useLeadVisitReports, useLeadActivities } from "../hooks/useLeads";
import { toast } from "sonner";
import { useState } from "react";
import { LeadForm } from "./lead-form";
import { LeadQualificationCard } from "./LeadQualificationCard";
import { ConvertLeadDialog } from "./convert-lead-dialog";
import { DealForm } from "@/features/sales-crm/pipeline-management/components/deal-form";
import { VisitReportForm } from "@/features/sales-crm/visit-report/components/visit-report-form";
import { CreateActivityDialog } from "@/features/sales-crm/visit-report/components/create-activity-dialog";
import { useCreateDeal } from "@/features/sales-crm/pipeline-management/hooks/useDeals";
import { useCreateVisitReport } from "@/features/sales-crm/visit-report/hooks/useVisitReports";
import { ProductInterestTab } from "@/features/sales-crm/visit-report/components/product-interest-tab";
import type { CreateDealFormData, UpdateDealFormData, DealFormData } from "@/features/sales-crm/pipeline-management/schemas/deal.schema";
import type { CreateVisitReportFormData, UpdateVisitReportFormData } from "@/features/sales-crm/visit-report/schemas/visit-report.schema";
import { useTranslations } from "next-intl";
import { useHasPermission } from "@/features/auth/providers/permissions-provider";
import { formatPhoneNumberToWA, formatEmailToMailto } from "@/lib/utils";
import { useTasks } from "@/features/sales-crm/task-management/hooks/useTasks";
import type { Task } from "@/features/sales-crm/task-management/types";
import type { Activity as CRMActivity } from "@/features/sales-crm/visit-report/types/activity";
import type { VisitReport } from "@/features/sales-crm/visit-report/types";

/** Builds a DiceBear `lorelei` URL to match other user avatars in the app */
function dicebearUrl(seed: string): string {
  return `https://api.dicebear.com/7.x/lorelei/svg?seed=${encodeURIComponent(seed)}`;
}

interface LeadDetailModalProps {
  readonly leadId: string | null;
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly onLeadUpdated?: () => void;
}

const statusColors: Record<string, "default" | "secondary" | "destructive" | "outline"> = {
  new: "outline",
  contacted: "secondary",
  interested: "secondary",
  qualified: "default",
  proposal_sent: "default",
  converted: "default",
  lost: "destructive",
};

export function LeadDetailModal({ leadId, open, onOpenChange, onLeadUpdated }: LeadDetailModalProps) {
  const { data, isLoading, error } = useLead(leadId || "");
  const deleteLead = useDeleteLead();
  const updateLead = useUpdateLead();
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [isConvertDialogOpen, setIsConvertDialogOpen] = useState(false);
  const [isCreateOpportunityDialogOpen, setIsCreateOpportunityDialogOpen] = useState(false);
  const t = useTranslations("leadManagement.leadDetail");
  const hasEditPermission = useHasPermission("leads.edit");
  const hasDeletePermission = useHasPermission("leads.delete");
  const hasConvertPermission = useHasPermission("leads.convert");
  const hasCreateOpportunityPermission = useHasPermission("pipeline.opportunity-create");
  const createDeal = useCreateDeal();

  const lead = data?.data;

  const handleUpdate = async (formData: Parameters<typeof updateLead.mutateAsync>[0]["data"]) => {
    if (!leadId) return;
    try {
      await updateLead.mutateAsync({ id: leadId, data: formData });
      toast.success(t("toast.updated"));
      setIsEditDialogOpen(false);
      onLeadUpdated?.();
    } catch {
      // Error handled by interceptor
    }
  };

  const handleDelete = async () => {
    if (!leadId) return;
    try {
      await deleteLead.mutateAsync(leadId);
      toast.success(t("toast.deleted"));
      setIsDeleteDialogOpen(false);
      onOpenChange(false);
      onLeadUpdated?.();
    } catch {
      // Error handled by interceptor
    }
  };

  const formatDate = (dateString?: string | null) => {
    if (!dateString) return t("fallbacks.noDate");
    const date = new Date(dateString);
    if (Number.isNaN(date.getTime())) return t("fallbacks.invalidDate");
    return date.toLocaleDateString("id-ID", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  };

  return (
    <>
      <Drawer
        open={open}
        onOpenChange={onOpenChange}
        title={t("drawerTitle")}
        side="right"
        defaultWidth={672}
      >
        {isLoading && (
          <div className="space-y-6">
            <div className="flex items-center gap-4">
              <Skeleton className="h-16 w-16 rounded-lg" />
              <div className="space-y-2">
                <Skeleton className="h-6 w-48" />
                <Skeleton className="h-4 w-32" />
              </div>
            </div>
            <Card className="surface-panel border-border/70">
              <CardHeader>
                <Skeleton className="h-6 w-32" />
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <Skeleton className="h-4 w-full" />
                  <Skeleton className="h-4 w-3/4" />
                  <Skeleton className="h-4 w-1/2" />
                </div>
              </CardContent>
            </Card>
          </div>
        )}

        {error && (
          <div className="text-center text-muted-foreground py-8">
            {t("loadError")}
          </div>
        )}

        {!isLoading && !error && lead && (
          <div className="space-y-6">
            {/* Lead Header */}
            <div className="crm-hero space-y-3 rounded-3xl border border-border/70 px-5 py-5">
              <div className="flex items-center gap-4">
                <Avatar className="h-12 w-12 sm:h-16 sm:w-16 rounded-xl shrink-0">
                  <AvatarImage
                    src={dicebearUrl(lead.email || `${lead.first_name} ${lead.last_name}`)}
                    alt={`${lead.first_name} ${lead.last_name}`}
                    className="rounded-xl"
                  />
                  <AvatarFallback className="text-xs bg-muted" />
                </Avatar>
                <div className="flex-1 min-w-0">
                  <h2 className="text-xl sm:text-2xl font-medium tracking-tight wrap-break-word">
                    {lead.first_name} {lead.last_name}
                  </h2>
                  <div className="flex flex-wrap items-center gap-2 mt-1">
                    <Badge variant={statusColors[lead.lead_status] || "outline"} className="capitalize text-xs">
                      {lead.lead_status}
                    </Badge>
                    {lead.company_name && (
                      <span className="text-xs sm:text-sm text-muted-foreground wrap-break-word">{lead.company_name}</span>
                    )}
                  </div>
                </div>
              </div>
              <div className="flex flex-wrap items-center gap-2">
                {hasEditPermission && (
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setIsEditDialogOpen(true)}
                    className="shrink-0"
                  >
                    <Edit className="h-4 w-4" />
                  </Button>
                )}
                {hasConvertPermission && lead.lead_status === "qualified" && (
                  <Button
                    size="sm"
                    onClick={() => setIsConvertDialogOpen(true)}
                    className="flex-1 sm:flex-initial"
                  >
                    <TrendingUp className="h-4 w-4 mr-2" />
                    {t("actions.convert")}
                  </Button>
                )}
                {hasCreateOpportunityPermission && (
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={() => setIsCreateOpportunityDialogOpen(true)}
                    className="flex-1 sm:flex-initial"
                  >
                    <Plus className="h-4 w-4 mr-2" />
                    <span className="truncate">{t("actions.createOpportunity") || "Create Opportunity"}</span>
                  </Button>
                )}
                {hasDeletePermission && (
                  <Button
                    variant="destructive"
                    size="sm"
                    onClick={() => setIsDeleteDialogOpen(true)}
                    disabled={deleteLead.isPending}
                    className="shrink-0"
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                )}
              </div>
            </div>

            {/* Basic Information */}
            <Card>
              <CardHeader>
                <CardTitle>{t("sections.basicInfo")}</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <div className="flex items-center gap-2 text-sm text-muted-foreground">
                      <Mail className="h-4 w-4" />
                      <span>{t("sections.email")}</span>
                    </div>
                    <div className="text-base font-medium">
                      <a
                        href={formatEmailToMailto(lead.email)}
                        className="text-primary hover:underline cursor-pointer truncate block max-w-full"
                        title={lead.email}
                      >
                        {lead.email}
                      </a>
                    </div>
                  </div>

                  {lead.phone && (
                    <div className="space-y-2">
                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <Phone className="h-4 w-4" />
                        <span>{t("sections.phone")}</span>
                      </div>
                      <a 
                        href={formatPhoneNumberToWA(lead.phone)} 
                        target="_blank" 
                        rel="noreferrer" 
                        className="text-base font-medium text-primary hover:underline cursor-pointer"
                      >
                        {lead.phone}
                      </a>
                    </div>
                  )}

                  {lead.job_title && (
                    <div className="space-y-2">
                      <div className="text-sm text-muted-foreground">
                        <span>{t("sections.jobTitle")}</span>
                      </div>
                      <div className="text-base">{lead.job_title}</div>
                    </div>
                  )}

                  {lead.industry && (
                    <div className="space-y-2">
                      <div className="text-sm text-muted-foreground">
                        <span>{t("sections.industry")}</span>
                      </div>
                      <div className="text-base">{lead.industry}</div>
                    </div>
                  )}

                  <div className="space-y-2">
                    <div className="flex items-center gap-2 text-sm text-muted-foreground">
                      <FileText className="h-4 w-4" />
                      <span>{t("sections.leadSource")}</span>
                    </div>
                    <div className="text-base">{lead.lead_source}</div>
                  </div>

                  {lead.assigned_user && (
                    <div className="space-y-2">
                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <UserPlus className="h-4 w-4" />
                        <span>{t("sections.assignedTo")}</span>
                      </div>
                      <div className="flex items-center gap-2">
                        <Avatar className="h-7 w-7 shrink-0">
                          <AvatarImage
                            src={dicebearUrl(lead.assigned_user.email || lead.assigned_user.name)}
                            alt={lead.assigned_user.name}
                          />
                          <AvatarFallback className="text-xs" />
                        </Avatar>
                        <span className="text-base">{lead.assigned_user.name}</span>
                      </div>
                    </div>
                  )}
                </div>

                {(lead.address || lead.city || lead.province) && (
                  <div className="space-y-2 pt-2 border-t">
                    <div className="flex items-center gap-2 text-sm text-muted-foreground">
                      <MapPin className="h-4 w-4" />
                      <span>{t("sections.address")}</span>
                    </div>
                    <div className="text-base">
                      {[lead.address, lead.city, lead.province, lead.postal_code, lead.country]
                        .filter(Boolean)
                        .join(", ")}
                    </div>
                  </div>
                )}

                {lead.website && (
                  <div className="space-y-2 pt-2 border-t">
                    <div className="flex items-center gap-2 text-sm text-muted-foreground">
                      <Globe className="h-4 w-4" />
                      <span>{t("sections.website")}</span>
                    </div>
                    <div className="text-base">
                      <a
                        href={lead.website}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-primary hover:underline"
                      >
                        {lead.website}
                      </a>
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>

            {/* Notes */}
            {lead.notes && (
              <Card className="surface-panel border-border/70">
                <CardHeader>
                  <CardTitle className="flex items-center gap-2">
                    <FileText className="h-5 w-5" />
                    {t("sections.notes")}
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-base whitespace-pre-wrap">{lead.notes}</p>
                </CardContent>
              </Card>
            )}

            {/* Visit Reports & Activities */}
            <Card className="surface-panel border-border/70">
              <CardHeader>
                <CardTitle>{t("sections.relatedActivities")}</CardTitle>
                <CardDescription>
                  {t("sections.relatedActivitiesDescription")}
                </CardDescription>
              </CardHeader>
              <CardContent>
                <Tabs defaultValue="visit-reports" className="w-full">
                  <TabsList className="grid w-full grid-cols-2 md:grid-cols-4">
                      <TabsTrigger value="visit-reports">
                        <MapPin className="h-4 w-4 mr-2" />
                        {t("sections.visitReports")}
                      </TabsTrigger>
                      <TabsTrigger value="activities">
                        <Activity className="h-4 w-4 mr-2" />
                        {t("sections.activities")}
                      </TabsTrigger>
                      <TabsTrigger value="tasks">
                        <CheckSquare className="h-4 w-4 mr-2" />
                        {t("sections.tasks")}
                      </TabsTrigger>
                      <TabsTrigger value="product-interest">
                        <Package className="h-4 w-4 mr-2" />
                        {t("sections.productInterest")}
                      </TabsTrigger>
                    </TabsList>
                    <TabsContent value="visit-reports" className="mt-4">
                      <LeadVisitReportsList leadId={leadId || ""} />
                    </TabsContent>
                    <TabsContent value="activities" className="mt-4">
                      <LeadActivitiesList leadId={leadId || ""} />
                    </TabsContent>
                    <TabsContent value="tasks" className="mt-4">
                      <LeadTasksList leadId={leadId || ""} />
                    </TabsContent>
                    <TabsContent value="product-interest" className="mt-4">
                      <LeadProductInterestList leadId={leadId || ""} />
                    </TabsContent>
                  </Tabs>
                </CardContent>
              </Card>

            {/* Metadata */}
            <Card className="surface-panel border-border/70">
              <CardHeader>
                <CardTitle>{t("sections.metadata")}</CardTitle>
              </CardHeader>
              <CardContent className="space-y-2 text-sm text-muted-foreground">
                <p>
                  {t("sections.createdAt")} {formatDate(lead.created_at)}
                </p>
                <p>
                  {t("sections.updatedAt")} {formatDate(lead.updated_at)}
                </p>
              </CardContent>
            </Card>
          </div>
        )}
      </Drawer>

      {/* Edit Dialog */}
      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent className="sm:max-w-[700px] max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>{t("editDialog.title")}</DialogTitle>
          </DialogHeader>
          {lead && (
            <Tabs defaultValue="profile" className="w-full">
              <TabsList className="grid w-full grid-cols-2">
                <TabsTrigger value="profile">{t("editDialog.tabs.profile")}</TabsTrigger>
                <TabsTrigger value="bant">{t("editDialog.tabs.bant")}</TabsTrigger>
              </TabsList>
              <TabsContent value="profile" className="mt-4">
                <LeadForm
                  lead={lead}
                  onSubmit={handleUpdate}
                  onCancel={() => setIsEditDialogOpen(false)}
                  isLoading={updateLead.isPending}
                />
              </TabsContent>
              <TabsContent value="bant" className="mt-4">
                <LeadQualificationCard leadId={lead.id} />
              </TabsContent>
            </Tabs>
          )}
        </DialogContent>
      </Dialog>

      {/* Delete Dialog */}
      <DeleteDialog
        open={isDeleteDialogOpen}
        onOpenChange={setIsDeleteDialogOpen}
        onConfirm={handleDelete}
        title={t("deleteDialog.title")}
        description={t("deleteDialog.description")}
        isLoading={deleteLead.isPending}
      />

      {/* Convert Dialog */}
      {lead && (
        <ConvertLeadDialog
          lead={lead}
          open={isConvertDialogOpen}
          onOpenChange={setIsConvertDialogOpen}
          onSuccess={() => {
            onLeadUpdated?.();
            onOpenChange(false);
          }}
        />
      )}

      {/* Create Opportunity Dialog */}
      {lead && (
        <Dialog open={isCreateOpportunityDialogOpen} onOpenChange={setIsCreateOpportunityDialogOpen}>
          <DialogContent className="sm:max-w-[600px] max-h-[90vh] overflow-y-auto mx-2 sm:mx-auto">
            <DialogHeader>
              <DialogTitle>{t("actions.createOpportunity")}</DialogTitle>
            </DialogHeader>
            <DealForm
              initialLeadId={lead.id}
              initialAccountId={lead.account_id || undefined}
              onSubmit={async (data: CreateDealFormData | UpdateDealFormData) => {
                try {
                  await createDeal.mutateAsync(data as DealFormData);
                  toast.success(t("toast.opportunityCreated"));
                  setIsCreateOpportunityDialogOpen(false);
                  onLeadUpdated?.();
                } catch {
                  // Error handled by interceptor
                }
              }}
              onCancel={() => setIsCreateOpportunityDialogOpen(false)}
              isLoading={createDeal.isPending}
            />
          </DialogContent>
        </Dialog>
      )}
    </>
  );
}

// Visit Reports List Component for Lead
function LeadVisitReportsList({ leadId }: { readonly leadId: string }) {
  const { data, isLoading } = useLeadVisitReports(leadId, { per_page: 10 });
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const createVisitReport = useCreateVisitReport();
  const t = useTranslations("leadManagement.leadDetail");
  const hasCreatePermission = useHasPermission("visit-reports.create");

  if (isLoading) {
    return (
      <div className="space-y-3">
        {Array.from({ length: 3 }).map((_, i) => (
          <Skeleton key={`visit-report-skeleton-${i}`} className="h-16 w-full" />
        ))}
      </div>
    );
  }

  const visitReports = data?.data ?? [];

  const handleCreate = async (formData: CreateVisitReportFormData | UpdateVisitReportFormData) => {
    try {
      await createVisitReport.mutateAsync(formData as CreateVisitReportFormData);
      toast.success("Visit report created successfully");
      setIsCreateOpen(false);
    } catch {
      // Error handled by interceptor
    }
  };

  return (
    <>
      {visitReports.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-12 px-4">
          <div className="flex h-14 w-14 items-center justify-center rounded-full bg-muted mb-4">
            <MapPin className="h-6 w-6 text-muted-foreground" />
          </div>
          <p className="text-sm font-medium text-foreground mb-1">{t("sections.noVisitReports")}</p>
          <p className="text-xs text-muted-foreground mb-4">{t("sections.relatedActivitiesDescription")}</p>
          {hasCreatePermission && (
            <Button size="sm" variant="outline" onClick={() => setIsCreateOpen(true)} className="gap-1.5">
              <Plus className="h-3.5 w-3.5" />
              {t("shortcuts.newVisitReport")}
            </Button>
          )}
        </div>
      ) : (
        <div className="space-y-3">
          {visitReports.map((vr) => (
            <div key={vr.id} className="flex items-start gap-3 p-3 border rounded-lg hover:bg-accent/50 transition-colors">
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 mb-1">
                  <Badge
                    variant={
                      vr.status === "approved"
                        ? "default"
                        : vr.status === "rejected"
                          ? "destructive"
                          : "secondary"
                    }
                  >
                    {vr.status}
                  </Badge>
                  <span className="text-sm text-muted-foreground">
                    {vr.visit_date ? new Date(vr.visit_date).toLocaleDateString("id-ID") : ""}
                  </span>
                </div>
                <p className="text-sm font-medium line-clamp-1">{vr.purpose}</p>
                {vr.account && (
                  <p className="text-xs text-muted-foreground mt-1">
                    {vr.account.name}
                  </p>
                )}
              </div>
            </div>
          ))}
          {data?.meta?.pagination && data.meta.pagination.total > 10 && (
            <div className="text-center pt-2">
              <p className="text-xs text-muted-foreground">
                {t("sections.showing")} {visitReports.length} {t("sections.of")} {data.meta.pagination.total} {t("sections.visitReports")}
              </p>
            </div>
          )}
          {hasCreatePermission && (
            <button
              type="button"
              onClick={() => setIsCreateOpen(true)}
              className="flex w-full items-center justify-center gap-2 p-3 border-2 border-dashed border-muted-foreground/25 rounded-lg text-muted-foreground hover:border-primary/50 hover:text-primary hover:bg-accent/30 transition-colors cursor-pointer"
            >
              <Plus className="h-4 w-4" />
              <span className="text-sm font-medium">{t("shortcuts.newVisitReport")}</span>
            </button>
          )}
        </div>
      )}

      {/* Create Visit Report Dialog - pre-fills the lead tab */}
      <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
        <DialogContent className="sm:max-w-[700px] max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>{t("shortcuts.newVisitReport")}</DialogTitle>
          </DialogHeader>
          <VisitReportForm
            open={isCreateOpen}
            initialLeadId={leadId}
            onSubmit={handleCreate}
            onCancel={() => setIsCreateOpen(false)}
            isLoading={createVisitReport.isPending}
          />
        </DialogContent>
      </Dialog>
    </>
  );
}

// Activities List Component for Lead
function LeadActivitiesList({ leadId }: { readonly leadId: string }) {
  const { data, isLoading } = useLeadActivities(leadId, { per_page: 10 });
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const t = useTranslations("leadManagement.leadDetail");
  const hasCreatePermission = useHasPermission("activities.create");

  if (isLoading) {
    return (
      <div className="space-y-3">
        {Array.from({ length: 3 }).map((_, i) => (
          <Skeleton key={`activity-skeleton-${i}`} className="h-16 w-full" />
        ))}
      </div>
    );
  }

  const activities = data?.data ?? [];

  return (
    <>
      {activities.length === 0 ? (
        <div className="flex flex-col items-center justify-center py-12 px-4">
          <div className="flex h-14 w-14 items-center justify-center rounded-full bg-muted mb-4">
            <Activity className="h-6 w-6 text-muted-foreground" />
          </div>
          <p className="text-sm font-medium text-foreground mb-1">{t("sections.noActivities")}</p>
          <p className="text-xs text-muted-foreground mb-4">{t("sections.relatedActivitiesDescription")}</p>
          {hasCreatePermission && (
            <Button size="sm" variant="outline" onClick={() => setIsCreateOpen(true)} className="gap-1.5">
              <Plus className="h-3.5 w-3.5" />
              {t("shortcuts.newActivity")}
            </Button>
          )}
        </div>
      ) : (
        <div className="space-y-3">
          {activities.map((activity) => (
            <div key={activity.id} className="flex items-start gap-3 p-3 border rounded-lg hover:bg-accent/50 transition-colors">
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 mb-1">
                  <Badge variant="outline" className="capitalize">
                    {activity.type}
                  </Badge>
                  <span className="text-sm text-muted-foreground">
                    {(() => {
                      let dateStr = activity.timestamp;
                      // For VISIT type activities, use visit_date from metadata
                      if (activity.type === "visit" && activity.metadata && typeof activity.metadata === "object") {
                        const meta = activity.metadata as Record<string, unknown>;
                        if (typeof meta.visit_date === "string") {
                          dateStr = meta.visit_date;
                        }
                      }
                      return new Date(dateStr).toLocaleDateString("id-ID", {
                        year: "numeric",
                        month: "short",
                        day: "numeric",
                        hour: "2-digit",
                        minute: "2-digit",
                      });
                    })()}
                  </span>
                </div>
                <p className="text-sm font-medium line-clamp-2">{activity.description}</p>
                {activity.account && (
                  <p className="text-xs text-muted-foreground mt-1">
                    {activity.account.name}
                  </p>
                )}
              </div>
            </div>
          ))}
          {data?.meta?.pagination && data.meta.pagination.total > 10 && (
            <div className="text-center pt-2">
              <p className="text-xs text-muted-foreground">
                {t("sections.showing")} {activities.length} {t("sections.of")} {data.meta.pagination.total} {t("sections.activities")}
              </p>
            </div>
          )}
          {hasCreatePermission && (
            <button
              type="button"
              onClick={() => setIsCreateOpen(true)}
              className="flex w-full items-center justify-center gap-2 p-3 border-2 border-dashed border-muted-foreground/25 rounded-lg text-muted-foreground hover:border-primary/50 hover:text-primary hover:bg-accent/30 transition-colors cursor-pointer"
            >
              <Plus className="h-4 w-4" />
              <span className="text-sm font-medium">{t("shortcuts.newActivity")}</span>
            </button>
          )}
        </div>
      )}

      {/* Create Activity Dialog - pre-filled with this lead */}
      <CreateActivityDialog
        open={isCreateOpen}
        onOpenChange={setIsCreateOpen}
        leadId={leadId}
        onSuccess={() => setIsCreateOpen(false)}
      />
    </>
  );
}

function LeadTasksList({ leadId }: { readonly leadId: string }) {
  const { data, isLoading } = useTasks({ lead_id: leadId, per_page: 10 });
  const t = useTranslations("leadManagement.leadDetail");

  if (isLoading) {
    return (
      <div className="space-y-3">
        {Array.from({ length: 3 }).map((_, i) => (
          <Skeleton key={`task-skeleton-${i}`} className="h-16 w-full" />
        ))}
      </div>
    );
  }

  const tasks = data?.data ?? [];

  if (tasks.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 px-4">
        <div className="mb-4 flex h-14 w-14 items-center justify-center rounded-full bg-muted">
          <CheckSquare className="h-6 w-6 text-muted-foreground" />
        </div>
        <p className="text-sm font-medium text-foreground mb-1">{t("sections.noTasks")}</p>
        <p className="text-xs text-muted-foreground">{t("sections.relatedActivitiesDescription")}</p>
      </div>
    );
  }

  return (
    <div className="space-y-3">
      {tasks.map((task) => (
        <TaskListCard key={task.id} task={task} />
      ))}
      {data?.meta?.pagination && data.meta.pagination.total > tasks.length && (
        <div className="text-center pt-2">
          <p className="text-xs text-muted-foreground">
            {t("sections.showing")} {tasks.length} {t("sections.of")} {data.meta.pagination.total} {t("sections.tasks")}
          </p>
        </div>
      )}
    </div>
  );
}

function LeadProductInterestList({ leadId }: { readonly leadId: string }) {
  const { data: visitsData, isLoading: isVisitsLoading } = useLeadVisitReports(leadId, { per_page: 50 });
  const { data: activitiesData, isLoading: isActivitiesLoading } = useLeadActivities(leadId, { per_page: 50 });

  const activities = buildProductInterestActivities(
    (visitsData?.data ?? []) as VisitReport[],
    (activitiesData?.data ?? []) as CRMActivity[],
  );

  return (
    <ProductInterestTab
      activities={activities}
      isLoading={isVisitsLoading || isActivitiesLoading}
    />
  );
}

function TaskListCard({ task }: { readonly task: Task }) {
  const dueDate = task.due_date ? new Date(task.due_date) : null;
  const dueDateLabel =
    dueDate && !Number.isNaN(dueDate.getTime())
      ? dueDate.toLocaleDateString("id-ID", {
          year: "numeric",
          month: "short",
          day: "numeric",
        })
      : null;

  return (
    <div className="flex items-start gap-3 rounded-lg border border-border/70 p-3 transition-colors hover:bg-accent/40">
      <div className="flex-1 min-w-0">
        <div className="mb-1 flex flex-wrap items-center gap-2">
          <Badge variant={task.status === "completed" ? "default" : "outline"} className="capitalize">
            {task.status}
          </Badge>
          <Badge variant="secondary" className="capitalize">
            {task.priority}
          </Badge>
          {dueDateLabel && <span className="text-sm text-muted-foreground">{dueDateLabel}</span>}
        </div>
        <p className="text-sm font-medium line-clamp-1">{task.title}</p>
        {task.description && <p className="mt-1 text-xs text-muted-foreground line-clamp-2">{task.description}</p>}
        {task.assigned_user?.name && (
          <p className="mt-2 text-xs text-muted-foreground">
            {task.assigned_user.name}
          </p>
        )}
      </div>
    </div>
  );
}

function buildProductInterestActivities(visits: VisitReport[], activities: CRMActivity[]): CRMActivity[] {
  const visitActivities: CRMActivity[] = visits.map((visit) => ({
    id: `visit-report-${visit.id}`,
    type: "visit",
    account_id: visit.account_id,
    contact_id: visit.contact_id,
    deal_id: visit.deal_id,
    lead_id: visit.lead_id,
    user_id: visit.sales_rep_id,
    description: visit.purpose,
    timestamp: visit.visit_date,
    metadata: {
      ...(visit.metadata ?? {}),
      visit_date: visit.visit_date,
    },
    created_at: visit.created_at,
    updated_at: visit.updated_at,
    account: visit.account,
    contact: visit.contact,
  }));

  return [...activities, ...visitActivities];
}
