"use client";

import { useEffect, useMemo, useState, type ReactNode } from "react";
import {
  Activity as ActivityIcon,
  ArrowLeft,
  ArrowRightLeft,
  ClipboardList,
  Clock3,
  Edit3,
  ListTodo,
  MapPinned,
  Package,
  Plus,
  Trash2,
} from "lucide-react";
import { toast } from "sonner";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { PageMotion } from "@/components/motion";
import { Skeleton } from "@/components/ui/skeleton";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { cn, formatCurrency } from "@/lib/utils";
import { useRouter } from "@/i18n/routing";
import { useHasPermission } from "@/features/auth/providers/permissions-provider";
import { CreateActivityDialog } from "@/features/sales-crm/visit-report/components/create-activity-dialog";
import { ProductInterestTab } from "@/features/sales-crm/visit-report/components/product-interest-tab";
import { useCreateVisitReport } from "@/features/sales-crm/visit-report/hooks/useVisitReports";
import { VisitReportDetailModal } from "@/features/sales-crm/visit-report/components/visit-report-detail-modal";
import { VisitReportForm } from "@/features/sales-crm/visit-report/components/visit-report-form";
import type { CreateVisitReportFormData, UpdateVisitReportFormData } from "@/features/sales-crm/visit-report/schemas/visit-report.schema";
import { useCreateTask, useTasks } from "@/features/sales-crm/task-management/hooks/useTasks";
import { TaskDetailModal } from "@/features/sales-crm/task-management/components/task-detail-modal";
import { TaskForm } from "@/features/sales-crm/task-management/components/task-form";
import type { CreateTaskFormData, UpdateTaskFormData } from "@/features/sales-crm/task-management/schemas/task.schema";
import type { Task } from "@/features/sales-crm/task-management/types";
import { useLead, useDeleteLead, useLeadActivities, useLeadVisitReports, useUpdateLead } from "../hooks/useLeads";
import { useLeadQualification } from "../hooks/useLeadQualification";
import { LeadForm } from "./lead-form";
import { LeadQualificationCard } from "./LeadQualificationCard";
import { ConvertLeadDialog } from "./convert-lead-dialog";
import type { Activity } from "@/features/sales-crm/visit-report/types/activity";
import type { VisitReport } from "@/features/sales-crm/visit-report/types";
import type { Lead } from "../types";

interface LeadDetailShellProps {
  readonly leadId: string;
}

const statusBadgeVariant: Record<string, "default" | "secondary" | "destructive" | "outline"> = {
  new: "outline",
  contacted: "secondary",
  interested: "secondary",
  qualified: "default",
  proposal_sent: "default",
  converted: "default",
  lost: "destructive",
};

const DETAIL_LIST_PAGE_SIZE = 10;
const PRODUCT_INTEREST_LIST_PAGE_SIZE = 100;

export function LeadDetailShell({ leadId }: LeadDetailShellProps) {
  const router = useRouter();
  const { data: leadData, isLoading, error } = useLead(leadId);
  const { qualification } = useLeadQualification(leadId);
  const { data: activitiesData, refetch: refetchActivities } = useLeadActivities(leadId, { per_page: DETAIL_LIST_PAGE_SIZE });
  const { data: visitReportsData, refetch: refetchVisitReports } = useLeadVisitReports(leadId, { per_page: DETAIL_LIST_PAGE_SIZE });
  const {
    data: productInterestActivitiesData,
    isLoading: isProductInterestActivitiesLoading,
    refetch: refetchProductInterestActivities,
  } = useLeadActivities(leadId, { per_page: PRODUCT_INTEREST_LIST_PAGE_SIZE });
  const {
    data: productInterestVisitReportsData,
    isLoading: isProductInterestVisitReportsLoading,
    refetch: refetchProductInterestVisitReports,
  } = useLeadVisitReports(leadId, { per_page: PRODUCT_INTEREST_LIST_PAGE_SIZE });
  const { data: tasksData, refetch: refetchTasks } = useTasks({ lead_id: leadId, per_page: DETAIL_LIST_PAGE_SIZE });
  const updateLead = useUpdateLead();
  const deleteLead = useDeleteLead();
  const createVisitReport = useCreateVisitReport();
  const createTask = useCreateTask();

  const [activeTab, setActiveTab] = useState("activities");
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [isConvertDialogOpen, setIsConvertDialogOpen] = useState(false);
  const [isTaskModalOpen, setIsTaskModalOpen] = useState(false);
  const [isVisitModalOpen, setIsVisitModalOpen] = useState(false);
  const [isActivityModalOpen, setIsActivityModalOpen] = useState(false);
  const [viewingVisitReportId, setViewingVisitReportId] = useState<string | null>(null);
  const [viewingTaskId, setViewingTaskId] = useState<string | null>(null);

  const hasEditPermission = useHasPermission("leads.edit");
  const hasDeletePermission = useHasPermission("leads.delete");
  const hasConvertPermission = useHasPermission("leads.convert");
  const hasCreateTaskPermission = useHasPermission("tasks.create");
  const hasVisitReportCreatePermission = useHasPermission("visit-reports.create");
  const hasCreateActivityPermission = useHasPermission("activities.create") || hasVisitReportCreatePermission;
  const hasCreateVisitPermission = hasVisitReportCreatePermission;

  const lead = leadData?.data as Lead | undefined;
  const activities = (activitiesData?.data ?? []) as Activity[];
  const visitReports = (visitReportsData?.data ?? []) as VisitReport[];
  const productInterestSourceActivities = (productInterestActivitiesData?.data ?? []) as Activity[];
  const productInterestSourceVisitReports = (productInterestVisitReportsData?.data ?? []) as VisitReport[];
  const tasks = (tasksData?.data ?? []) as Task[];
  const qualificationProductInterests = qualification?.need_target_products ?? [];
  const productInterestActivities = useMemo(
    () => buildProductInterestActivities(productInterestSourceVisitReports, productInterestSourceActivities),
    [productInterestSourceVisitReports, productInterestSourceActivities],
  );
  const isProductInterestLoading = isProductInterestActivitiesLoading || isProductInterestVisitReportsLoading;
  const timelineItems = useMemo(() => buildTimelineItems(visitReports, activities), [visitReports, activities]);
  const activityCount = timelineItems.length;
  const productInterestCount = getProductInterestCount(productInterestActivities);
  const hasActivityProductInterests = productInterestActivities.some((activity) =>
    Array.isArray(activity.metadata?.product_interests) && activity.metadata.product_interests.length > 0
  );
  const bantCount = [
    qualification?.budget_confirmed,
    qualification?.authority_confirmed,
    qualification?.need_confirmed,
  ].filter(Boolean).length;

  const leadName = useMemo(() => {
    if (!lead) return "Lead";
    return [lead.first_name, lead.last_name].filter(Boolean).join(" ") || lead.company_name || lead.email || "Lead";
  }, [lead]);

  const isConvertedLead = Boolean(
    lead && ((lead.converted_at || lead.lead_status === "converted") && lead.opportunity?.id)
  );
  const leadStatusValue = isConvertedLead ? "converted" : lead?.lead_status;
  const leadStatus = formatLabel(leadStatusValue);
  const statusVariant = leadStatusValue ? statusBadgeVariant[leadStatusValue] || "outline" : "outline";
  const canConvertLead = hasConvertPermission && leadStatusValue === "qualified";
  const valueText = lead?.estimated_value ? formatCurrency(lead.estimated_value) : "Rp 0";
  const leadScoreText = `${lead?.lead_score ?? 0}`;
  const activityCountText = `${activityCount}`;
  const leadMetaItems = [lead?.company_name, lead?.lead_source].filter(Boolean);

  useEffect(() => {
    if (!isConvertedLead || !lead?.opportunity?.id) return;
    router.replace(`/deals/${lead.opportunity.id}`);
  }, [isConvertedLead, lead?.opportunity?.id, router]);

  const handleDelete = async () => {
    try {
      await deleteLead.mutateAsync(leadId);
      toast.success("Lead deleted successfully");
      setIsDeleteDialogOpen(false);
    } catch {
      // handled upstream
    }
  };

  const handleUpdate = async (formData: Parameters<typeof updateLead.mutateAsync>[0]["data"]) => {
    try {
      await updateLead.mutateAsync({ id: leadId, data: formData });
      toast.success("Lead updated successfully");
      setIsEditModalOpen(false);
    } catch {
      // handled upstream
    }
  };

  const handleCreateTask = async (formData: CreateTaskFormData | UpdateTaskFormData) => {
    try {
      await createTask.mutateAsync({
        ...(formData as Record<string, unknown>),
        lead_id: leadId,
      } as CreateTaskFormData);
      toast.success("Task created successfully");
      setIsTaskModalOpen(false);
      refetchTasks();
    } catch {
      // handled upstream
    }
  };

  const handleCreateVisit = async (formData: CreateVisitReportFormData | UpdateVisitReportFormData) => {
    try {
      const createdVisit = await createVisitReport.mutateAsync({
        ...(formData as CreateVisitReportFormData),
        lead_id: leadId,
      });
      toast.success("Visit report created successfully");
      setIsVisitModalOpen(false);
      refetchVisitReports();
      refetchActivities();
      refetchProductInterestVisitReports();
      refetchProductInterestActivities();
      setViewingVisitReportId(createdVisit.data.id);
    } catch {
      // handled upstream
    }
  };

  const formatDate = (value?: string | null) => {
    if (!value) return "-";
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return "-";
    return date.toLocaleDateString("id-ID", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  };

  const formatDateTime = (value?: string | null) => {
    if (!value) return "-";
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return "-";
    return date.toLocaleString("id-ID", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  if (isLoading) {
    return <LeadDetailSkeleton />;
  }

  if (error || !lead) {
    return (
      <PageMotion className="p-2 sm:p-4 lg:p-6">
        <Card className="border-border/70 bg-card/80">
          <CardContent className="py-10 text-center text-muted-foreground">
            Failed to load lead details.
          </CardContent>
        </Card>
      </PageMotion>
    );
  }

  if (isConvertedLead) {
    return (
      <PageMotion className="p-2 sm:p-4 lg:p-6">
        <Card className="border-border/70 bg-card/80">
          <CardContent className="py-10 text-center text-muted-foreground">
            Redirecting to converted deal...
          </CardContent>
        </Card>
      </PageMotion>
    );
  }

  return (
    <PageMotion className="min-h-screen p-2 sm:p-4 lg:p-6">
      <div className="mx-auto max-w-[1680px] space-y-5 lg:space-y-6">
        <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
          <div className="flex items-start gap-3 sm:gap-4 min-w-0">
            <Button
              variant="ghost"
              size="icon"
              className="mt-0.5 shrink-0 cursor-pointer rounded-xl"
              onClick={() => window.history.back()}
            >
              <ArrowLeft className="h-4 w-4" />
            </Button>
            <div className="min-w-0">
              <div className="flex flex-wrap items-center gap-2">
                <h1 className="page-title text-2xl sm:text-3xl font-semibold tracking-tight truncate max-w-[280px] sm:max-w-[620px]">
                  {leadName}
                </h1>
                <Badge variant={statusVariant} className="rounded-full capitalize text-[0.68rem]">
                  {leadStatus}
                </Badge>
              </div>
              {leadMetaItems.length > 0 && (
                <div className="mt-1 flex flex-wrap items-center gap-2 text-sm text-muted-foreground">
                  {leadMetaItems.map((item, index) => (
                    <span key={`${item}-${index}`}>
                      {index > 0 && <span className="mr-2">•</span>}
                      {item}
                    </span>
                  ))}
                </div>
              )}
            </div>
          </div>

          <div className="flex flex-wrap items-center gap-2 lg:justify-end">
              {canConvertLead && (
              <Button
                variant="outline"
                size="sm"
                className="cursor-pointer"
                onClick={() => setIsConvertDialogOpen(true)}
              >
                <ArrowRightLeft className="mr-2 h-4 w-4" />
                Convert
              </Button>
            )}
            {hasEditPermission && (
              <Button
                variant="outline"
                size="sm"
                className="cursor-pointer"
                onClick={() => setIsEditModalOpen(true)}
              >
                <Edit3 className="mr-2 h-4 w-4" />
                Edit
              </Button>
            )}
            {hasDeletePermission && (
              <Button
                variant="ghost"
                size="icon"
                className="cursor-pointer rounded-xl text-destructive hover:text-destructive"
                onClick={() => setIsDeleteDialogOpen(true)}
              >
                <Trash2 className="h-4 w-4" />
              </Button>
            )}
          </div>
        </div>

        <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 xl:grid-cols-4">
          <MetricCard label="Value" value={valueText} />
          <MetricCard label="Activities" value={activityCountText} />
          <MetricCard label="Lead Score" value={leadScoreText} />
          <MetricCard label="BANT Qualification" value={`${bantCount}/3`} accent="success" />
        </div>

        <section className="grid grid-cols-1 items-start gap-4">
          <div className="min-w-0 space-y-4">
            <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
              <TabsList className="h-auto w-full justify-start gap-4 overflow-x-auto rounded-none border-b border-border bg-transparent p-0">
                <TabsTrigger value="activities" className="cursor-pointer gap-2 whitespace-nowrap rounded-none border-b-2 border-transparent px-4 py-3 text-sm text-muted-foreground data-[state=active]:border-primary data-[state=active]:bg-transparent data-[state=active]:text-foreground data-[state=active]:shadow-none">
                  <ActivityIcon className="h-4 w-4" />
                  Activities
                  <span className="ml-1 rounded-full bg-muted px-1.5 py-0.5 text-xs text-muted-foreground">
                    {activityCount}
                  </span>
                </TabsTrigger>
                <TabsTrigger value="tasks" className="cursor-pointer gap-2 whitespace-nowrap rounded-none border-b-2 border-transparent px-4 py-3 text-sm text-muted-foreground data-[state=active]:border-primary data-[state=active]:bg-transparent data-[state=active]:text-foreground data-[state=active]:shadow-none">
                  <ListTodo className="h-4 w-4" />
                  Tasks
                  <span className="ml-1 rounded-full bg-muted px-1.5 py-0.5 text-xs text-muted-foreground">
                    {tasks.length}
                  </span>
                </TabsTrigger>
                <TabsTrigger value="product-interest" className="cursor-pointer gap-2 whitespace-nowrap rounded-none border-b-2 border-transparent px-4 py-3 text-sm text-muted-foreground data-[state=active]:border-primary data-[state=active]:bg-transparent data-[state=active]:text-foreground data-[state=active]:shadow-none">
                  <Package className="h-4 w-4" />
                  Product Interest
                  <span className="ml-1 rounded-full bg-muted px-1.5 py-0.5 text-xs text-muted-foreground">
                    {productInterestCount}
                  </span>
                </TabsTrigger>
                <TabsTrigger value="bant" className="cursor-pointer gap-2 whitespace-nowrap rounded-none border-b-2 border-transparent px-4 py-3 text-sm text-muted-foreground data-[state=active]:border-primary data-[state=active]:bg-transparent data-[state=active]:text-foreground data-[state=active]:shadow-none">
                  <ClipboardList className="h-4 w-4" />
                  BANT
                </TabsTrigger>
              </TabsList>

              <TabsContent value="activities" className="mt-4">
                <TabCard
                  title="Activities"
                  description="Track calls, visits, meetings, and follow ups for this lead."
                  actions={
                    <>
                      {hasCreateVisitPermission && (
                        <Button size="sm" variant="outline" onClick={() => setIsVisitModalOpen(true)} className="cursor-pointer">
                          <MapPinned className="mr-1 h-4 w-4" />
                          Log Visit
                        </Button>
                      )}
                      {hasCreateActivityPermission && (
                        <Button size="sm" variant="outline" onClick={() => setIsActivityModalOpen(true)} className="cursor-pointer">
                          <Plus className="mr-1 h-4 w-4" />
                          Log Activity
                        </Button>
                      )}
                    </>
                  }
                >
                  {timelineItems.length === 0 ? (
                    <EmptyState
                      icon={<Clock3 className="h-10 w-10 text-muted-foreground/40" />}
                      title="No activities yet"
                      description="Create the first activity to start the lead timeline."
                    />
                  ) : (
                    <div className="space-y-4">
                      {timelineItems.map((item, index) => (
                        <div key={item.id} className="flex gap-3">
                          <div className="relative flex flex-col items-center">
                            <div className="flex h-8 w-8 items-center justify-center rounded-full border border-primary/30 bg-primary/10 text-primary">
                              {item.kind === "visit" ? <MapPinned className="h-4 w-4" /> : <ActivityIcon className="h-4 w-4" />}
                            </div>
                            {index !== timelineItems.length - 1 && <div className="mt-1 h-full w-px flex-1 bg-border/70" />}
                          </div>
                          <div
                            className={`min-w-0 flex-1 rounded-xl border border-border bg-muted/30 p-4 ${item.kind === "visit" ? "cursor-pointer transition-colors hover:bg-accent/40" : ""}`}
                            onClick={item.kind === "visit" ? () => setViewingVisitReportId(item.visit.id) : undefined}
                            role={item.kind === "visit" ? "button" : undefined}
                            tabIndex={item.kind === "visit" ? 0 : undefined}
                            onKeyDown={item.kind === "visit" ? (event) => {
                              if (event.key === "Enter" || event.key === " ") {
                                event.preventDefault();
                                setViewingVisitReportId(item.visit.id);
                              }
                            } : undefined}
                          >
                            <div className="flex flex-wrap items-center gap-2">
                              <Badge variant="outline" className="capitalize">
                                {item.kind === "visit" ? "Visit" : item.activity.activity_type?.name || item.activity.type}
                              </Badge>
                              <span className="text-xs text-muted-foreground">{formatDateTime(item.dateValue)}</span>
                              {item.kind === "activity" && item.activity.user?.name && <span className="text-xs text-muted-foreground">• {item.activity.user.name}</span>}
                              {item.kind === "visit" && item.visit.sales_rep?.name && <span className="text-xs text-muted-foreground">• {item.visit.sales_rep.name}</span>}
                            </div>
                            <p className="mt-2 whitespace-pre-wrap text-sm leading-relaxed text-foreground">
                              {item.description}
                            </p>
                            {item.accountName && (
                              <p className="mt-2 text-xs text-muted-foreground">{item.accountName}</p>
                            )}
                            {item.kind === "visit" && (
                              <div className="mt-3">
                                <Button
                                  size="sm"
                                  variant="outline"
                                  className="cursor-pointer"
                                  onClick={(event) => {
                                    event.stopPropagation();
                                    setViewingVisitReportId(item.visit.id);
                                  }}
                                >
                                  Open Visit Detail
                                </Button>
                              </div>
                            )}
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </TabCard>
              </TabsContent>

              <TabsContent value="tasks" className="mt-4">
                <TabCard
                  title="Tasks"
                  description="Follow-up tasks linked to this lead."
                  actions={
                    hasCreateTaskPermission ? (
                      <Button size="sm" variant="outline" onClick={() => setIsTaskModalOpen(true)} className="cursor-pointer">
                        <Plus className="mr-1 h-4 w-4" />
                        Add Task
                      </Button>
                    ) : null
                  }
                >
                  {tasks.length === 0 ? (
                    <EmptyState
                      icon={<ListTodo className="h-10 w-10 text-muted-foreground/40" />}
                      title="No tasks yet"
                      description="Create a task to keep the lead follow-up in view."
                    />
                  ) : (
                    <div className="space-y-3">
                      {tasks.map((task) => (
                        <button
                          key={task.id}
                          type="button"
                          onClick={() => setViewingTaskId(task.id)}
                          className="w-full rounded-xl border border-border bg-muted/30 p-4 text-left transition-colors hover:bg-accent/40"
                        >
                          <div className="flex items-start justify-between gap-3">
                            <div className="min-w-0">
                              <div className="font-medium">{task.title}</div>
                              <div className="mt-1 text-sm text-muted-foreground line-clamp-2">
                                {task.description || "No description"}
                              </div>
                            </div>
                            <Badge
                              variant={task.status === "completed" ? "default" : "outline"}
                              className="capitalize"
                            >
                              {task.status.replace(/_/g, " ")}
                            </Badge>
                          </div>
                          <div className="mt-3 flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
                            <span className="rounded-full bg-card px-2 py-1 capitalize ring-1 ring-border">{task.type.replace(/_/g, " ")}</span>
                            <span className="rounded-full bg-card px-2 py-1 capitalize ring-1 ring-border">{task.priority}</span>
                            {task.due_date && <span>Due {formatDateTime(task.due_date)}</span>}
                          </div>
                        </button>
                      ))}
                    </div>
                  )}
                </TabCard>
              </TabsContent>

              <TabsContent value="product-interest" className="mt-4">
                <TabCard
                  title="Product Interest"
                  description="Products captured from activity and visit logs with their actual interest values."
                >
                  {isProductInterestLoading ? (
                    <ProductInterestTab
                      activities={productInterestActivities}
                      isLoading={isProductInterestLoading}
                    />
                  ) : productInterestCount === 0 ? (
                    <EmptyState
                      icon={<Package className="h-10 w-10 text-muted-foreground/40" />}
                      title="No product interest recorded"
                      description="Product interests will appear here after activity or visit input."
                    />
                  ) : (
                    <div className="space-y-4">
                      {qualificationProductInterests.length > 0 && (
                        <div className="space-y-3">
                          <div className="rounded-xl border border-dashed border-border bg-background/70 p-4">
                            <p className="text-sm font-medium">BANT Need Products</p>
                            <p className="mt-1 text-xs text-muted-foreground">
                              Produk dari BANT ditampilkan terpisah karena tidak menyimpan nilai interest level, quantity, dan price.
                            </p>
                          </div>
                          {qualificationProductInterests.map((item) => (
                            <div key={`${item.product_id ?? "qualification"}-${item.product_name}`} className="rounded-xl border border-border bg-muted/30 p-4">
                              <div className="flex items-start justify-between gap-3">
                                <div>
                                  <p className="font-medium">{item.product_name}</p>
                                  <p className="mt-1 text-xs text-muted-foreground">
                                    {item.category_name || "BANT qualification"}
                                  </p>
                                </div>
                                <Badge variant="outline" className="capitalize">
                                  {item.category_name || "BANT"}
                                </Badge>
                              </div>
                            </div>
                          ))}
                        </div>
                      )}
                      {hasActivityProductInterests && (
                        <ProductInterestTab
                          activities={productInterestActivities}
                          isLoading={isProductInterestLoading}
                        />
                      )}
                    </div>
                  )}
                </TabCard>
              </TabsContent>

              <TabsContent value="bant" className="mt-4">
                <TabCard
                  title="BANT"
                  description="Update budget, authority, dan need lead dari satu tempat."
                >
                  <LeadQualificationCard leadId={leadId} />
                </TabCard>
              </TabsContent>
            </Tabs>
          </div>
        </section>
      </div>

      {lead && (
        <ConvertLeadDialog
          lead={lead}
          open={isConvertDialogOpen}
          onOpenChange={setIsConvertDialogOpen}
          onSuccess={(response) => {
            setIsConvertDialogOpen(false);
            const opportunityId = response.data.opportunity && typeof response.data.opportunity === "object" && "id" in response.data.opportunity
              ? String(response.data.opportunity.id)
              : "";
            if (opportunityId) {
              router.replace(`/deals/${opportunityId}`);
            }
          }}
        />
      )}

      <Dialog open={isEditModalOpen} onOpenChange={setIsEditModalOpen}>
        <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-[780px]">
          <DialogHeader>
            <DialogTitle>Edit Lead</DialogTitle>
          </DialogHeader>
          {lead && (
            <LeadForm
              lead={lead}
              onSubmit={handleUpdate}
              onCancel={() => setIsEditModalOpen(false)}
              isLoading={updateLead.isPending}
            />
          )}
        </DialogContent>
      </Dialog>

      <DeleteDialog
        open={isDeleteDialogOpen}
        onOpenChange={setIsDeleteDialogOpen}
        onConfirm={handleDelete}
        title="Delete Lead"
        description="Are you sure you want to delete this lead? This action cannot be undone."
        isLoading={deleteLead.isPending}
      />

      <Dialog open={isTaskModalOpen} onOpenChange={setIsTaskModalOpen}>
        <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-[700px] mx-2 sm:mx-auto">
          <DialogHeader>
            <DialogTitle>Create Task</DialogTitle>
          </DialogHeader>
          <TaskForm
            onSubmit={handleCreateTask}
            onCancel={() => setIsTaskModalOpen(false)}
            isLoading={createTask.isPending}
          />
        </DialogContent>
      </Dialog>

      <Dialog open={isVisitModalOpen} onOpenChange={setIsVisitModalOpen}>
        <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-[700px] mx-2 sm:mx-auto">
          <DialogHeader>
            <DialogTitle>Log Visit</DialogTitle>
          </DialogHeader>
          <VisitReportForm
            open={isVisitModalOpen}
            initialLeadId={lead.id}
            onSubmit={handleCreateVisit}
            onCancel={() => setIsVisitModalOpen(false)}
            isLoading={createVisitReport.isPending}
          />
        </DialogContent>
      </Dialog>

      <CreateActivityDialog
        open={isActivityModalOpen}
        onOpenChange={setIsActivityModalOpen}
        leadId={leadId}
        onSuccess={() => {
          refetchActivities();
          refetchProductInterestActivities();
        }}
      />

      <VisitReportDetailModal
        visitReportId={viewingVisitReportId}
        open={!!viewingVisitReportId}
        onOpenChange={(open) => {
          if (!open) setViewingVisitReportId(null);
        }}
        onVisitReportUpdated={() => {
          refetchVisitReports();
          refetchActivities();
          refetchProductInterestVisitReports();
          refetchProductInterestActivities();
        }}
      />

      <TaskDetailModal
        taskId={viewingTaskId}
        open={!!viewingTaskId}
        onOpenChange={(open) => {
          if (!open) setViewingTaskId(null);
        }}
        onTaskUpdated={() => {
          refetchTasks();
        }}
      />
    </PageMotion>
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
        <div className={cn(
          "mt-1 text-2xl font-semibold tracking-tight",
          accent === "success" ? "text-emerald-500" : "text-slate-900",
        )}>
          {value}
        </div>
      </CardContent>
    </Card>
  );
}

function TabCard({
  title,
  description,
  actions,
  children,
}: {
  readonly title: string;
  readonly description: string;
  readonly actions?: ReactNode;
  readonly children: ReactNode;
}) {
  return (
    <Card className="border-slate-200 bg-white shadow-sm">
      <CardHeader className="flex flex-row items-center justify-between gap-3 pb-4">
        <div>
          <CardTitle className="text-base">{title}</CardTitle>
          <p className="text-sm text-slate-500">{description}</p>
        </div>
        <div className="flex items-center gap-2">{actions}</div>
      </CardHeader>
      <CardContent>{children}</CardContent>
    </Card>
  );
}

function EmptyState({
  icon,
  title,
  description,
}: {
  readonly icon: ReactNode;
  readonly title: string;
  readonly description: string;
}) {
  return (
    <div className="flex flex-col items-center justify-center rounded-lg border border-dashed border-slate-200 bg-slate-50 px-4 py-10 text-center">
      <div className="mb-3 flex h-12 w-12 items-center justify-center rounded-full bg-white text-slate-400 shadow-sm ring-1 ring-slate-200">
        {icon}
      </div>
      <p className="text-sm font-medium text-slate-900">{title}</p>
      <p className="mt-1 max-w-sm text-xs text-slate-500">{description}</p>
    </div>
  );
}

type LeadTimelineItem =
  | {
      kind: "visit";
      id: string;
      dateValue: string;
      description: string;
      accountName?: string;
      visit: VisitReport;
      activity?: never;
    }
  | {
      kind: "activity";
      id: string;
      dateValue: string;
      description: string;
      accountName?: string;
      activity: Activity;
      visit?: never;
    };

function buildTimelineItems(visits: VisitReport[], activities: Activity[]): LeadTimelineItem[] {
  const visitIds = new Set(visits.map((visit) => visit.id));
  const visitItems: LeadTimelineItem[] = visits.map((visit) => ({
    kind: "visit",
    id: `visit-${visit.id}`,
    dateValue: visit.visit_date || visit.created_at || visit.updated_at,
    description: visit.notes?.trim() || visit.purpose || "Visit recorded",
    accountName: visit.account?.name,
    visit,
  }));

  const activityItems: LeadTimelineItem[] = activities
    .filter((activity) => {
      const visitReportId = getVisitReportIdFromActivity(activity);
      return !visitReportId || !visitIds.has(visitReportId);
    })
    .map((activity) => {
    const dateValue = getActivityDateValue(activity);
    return {
      kind: "activity",
      id: activity.id,
      dateValue,
      description: activity.description,
      accountName: activity.account?.name,
      activity,
    };
  });

  return [...visitItems, ...activityItems].sort((a, b) => {
    const timeA = new Date(a.dateValue).getTime();
    const timeB = new Date(b.dateValue).getTime();
    if (Number.isNaN(timeA) && Number.isNaN(timeB)) return 0;
    if (Number.isNaN(timeA)) return 1;
    if (Number.isNaN(timeB)) return -1;
    return timeB - timeA;
  });
}

function getActivityDateValue(activity: Activity): string {
  if (activity.metadata && typeof activity.metadata === "object") {
    const meta = activity.metadata as Record<string, unknown>;
    if (typeof meta.visit_date === "string") {
      return meta.visit_date;
    }
  }

  return activity.timestamp || activity.created_at || activity.updated_at;
}

function buildProductInterestActivities(visits: VisitReport[], activities: Activity[]): Activity[] {
  const visitIds = new Set(visits.map((visit) => visit.id));
  const visitActivities: Activity[] = visits.map((visit) => ({
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

  const standaloneActivities = activities.filter((activity) => {
    const visitReportId = getVisitReportIdFromActivity(activity);
    return !visitReportId || !visitIds.has(visitReportId);
  });

  return [...standaloneActivities, ...visitActivities];
}

function getProductInterestCount(activities: Activity[]) {
  return activities.reduce((count, activity) => {
    const items = Array.isArray(activity.metadata?.product_interests)
      ? activity.metadata.product_interests
      : [];
    return count + items.length;
  }, 0);
}

function getVisitReportIdFromActivity(activity: Activity): string | null {
  if (!activity.metadata || typeof activity.metadata !== "object") return null;
  const meta = activity.metadata as Record<string, unknown>;
  return typeof meta.visit_report_id === "string" ? meta.visit_report_id : null;
}

function LeadDetailSkeleton() {
  return (
    <PageMotion className="min-h-screen bg-slate-50/80 p-2 sm:p-4 lg:p-6">
      <div className="space-y-4 lg:space-y-6">
        <Skeleton className="h-4 w-80 bg-slate-200" />
        <div className="flex items-start gap-3">
          <Skeleton className="h-9 w-9 rounded-full bg-slate-200" />
          <div className="space-y-2">
            <Skeleton className="h-8 w-64 bg-slate-200" />
            <Skeleton className="h-4 w-48 bg-slate-200" />
          </div>
        </div>
        <div className="grid grid-cols-2 gap-3 sm:gap-4 xl:grid-cols-4">
          {Array.from({ length: 4 }).map((_, index) => (
            <Skeleton key={index} className="h-20 rounded-xl bg-slate-200" />
          ))}
        </div>
        <div className="grid grid-cols-1 gap-4">
          <div className="space-y-4">
            <Skeleton className="h-12 w-full rounded-xl bg-slate-200" />
            <Skeleton className="h-[520px] w-full rounded-2xl bg-slate-200" />
          </div>
        </div>
      </div>
    </PageMotion>
  );
}

function formatLabel(value?: string | null) {
  if (!value) return "Unknown";
  return value
    .replace(/_/g, " ")
    .replace(/\b\w/g, (char) => char.toUpperCase());
}
