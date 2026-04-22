"use client";

import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import {
  CheckCircle2,
  XCircle,
  ClipboardList,
  FileText,
  CalendarDays,
  UserPlus,
  TrendingUp,
  ExternalLink,
  AlertCircle,
  Loader2,
  Plus,
} from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { useSalesManagerTeamDraftApprovals } from "../hooks/useDashboard";
import {
  useApproveVisitReport,
  useRejectVisitReport,
} from "@/features/sales-crm/visit-report/hooks/useVisitReports";
import { useCreateTask } from "@/features/sales-crm/task-management/hooks/useTasks";
import { useHasPermission } from "@/features/master-data/user-management/hooks/useHasPermission";
import { TaskForm } from "@/features/sales-crm/task-management/components/task-form";
import type { CreateTaskFormData } from "@/features/sales-crm/task-management/schemas/task.schema";
import type {
  DraftVisitItem,
  DraftTaskItem,
  DraftScheduleItem,
  DraftLeadItem,
  DraftPipelineItem,
} from "../types";
import { formatCurrency } from "@/lib/utils";
import { format } from "date-fns";
import Link from "next/link";

// ============================================================================
// Types
// ============================================================================

interface TeamDraftApprovalsProps {
  readonly startDate?: string;
  readonly endDate?: string;
}

// ============================================================================
// Top-level Component
// ============================================================================

export function TeamDraftApprovals(_props: TeamDraftApprovalsProps) {
  const { data, isLoading, isError } = useSalesManagerTeamDraftApprovals();

  if (isLoading) {
    return <TeamDraftApprovalsSkeleton />;
  }

  if (isError) {
    return <TeamDraftApprovalsError />;
  }

  const approvals = data?.data;

  const tabs = [
    {
      key: "visits",
      label: "Visits",
      icon: FileText,
      total: approvals?.visits?.total ?? 0,
    },
    {
      key: "tasks",
      label: "Tasks",
      icon: ClipboardList,
      total: approvals?.tasks?.total ?? 0,
    },
    {
      key: "schedules",
      label: "Schedules",
      icon: CalendarDays,
      total: approvals?.schedules?.total ?? 0,
    },
    {
      key: "leads",
      label: "Leads",
      icon: UserPlus,
      total: approvals?.leads?.total ?? 0,
    },
    {
      key: "pipeline",
      label: "Pipeline",
      icon: TrendingUp,
      total: approvals?.pipeline?.total ?? 0,
    },
  ] as const;

  const totalPending = approvals?.total ?? 0;

  return (
    <Card className="border-0 shadow-sm h-full flex flex-col">
      <CardHeader className="px-4 sm:px-6 pb-2 sm:pb-3">
        <div className="flex items-center justify-between gap-2">
          <div className="flex items-center gap-1.5 sm:gap-2">
            <AlertCircle className="h-4 w-4 sm:h-5 sm:w-5 text-warning shrink-0" />
            <CardTitle className="text-sm sm:text-base">Team Pending Approvals</CardTitle>
          </div>
          {totalPending > 0 ? (
            <Badge
              variant="destructive"
              className="text-[10px] sm:text-xs font-medium px-2 py-0.5 shrink-0"
            >
              {totalPending} pending
            </Badge>
          ) : (
            <Badge
              variant="secondary"
              className="text-[10px] sm:text-xs font-medium px-2 py-0.5 shrink-0 text-green-700 bg-green-100"
            >
              All clear
            </Badge>
          )}
        </div>
        <p className="text-[11px] sm:text-xs text-muted-foreground mt-1">
          Items submitted by your team awaiting your review
        </p>
      </CardHeader>

      <CardContent className="px-4 sm:px-6 pb-4 sm:pb-6 flex-1 overflow-y-auto max-h-[520px]">
        <Tabs defaultValue="visits" className="w-full">
          {/* Tab list ─ scrollable on small screens */}
          <TabsList className="flex flex-wrap h-auto gap-1 bg-muted/50 rounded-lg p-1 mb-4">
            {tabs.map(({ key, label, icon: Icon, total }) => (
              <TabsTrigger
                key={key}
                value={key}
                className="flex items-center gap-1.5 text-[11px] sm:text-xs px-2.5 py-1.5 rounded-md data-[state=active]:bg-background data-[state=active]:shadow-sm cursor-pointer"
              >
                <Icon className="h-3.5 w-3.5" />
                <span>{label}</span>
                {total > 0 && (
                  <span className="ml-0.5 inline-flex items-center justify-center rounded-full bg-destructive text-destructive-foreground text-[10px] font-semibold w-4 h-4 leading-none">
                    {total}
                  </span>
                )}
              </TabsTrigger>
            ))}
          </TabsList>

          {/* ── Visits ── */}
          <TabsContent value="visits">
            <VisitApprovalList
              items={approvals?.visits?.items ?? []}
              total={approvals?.visits?.total ?? 0}
            />
          </TabsContent>

          {/* ── Tasks ── */}
          <TabsContent value="tasks">
            <TaskApprovalList
              items={approvals?.tasks?.items ?? []}
              total={approvals?.tasks?.total ?? 0}
            />
          </TabsContent>

          {/* ── Schedules ── */}
          <TabsContent value="schedules">
            <ScheduleApprovalList
              items={approvals?.schedules?.items ?? []}
              total={approvals?.schedules?.total ?? 0}
            />
          </TabsContent>

          {/* ── Leads ── */}
          <TabsContent value="leads">
            <LeadApprovalList
              items={approvals?.leads?.items ?? []}
              total={approvals?.leads?.total ?? 0}
            />
          </TabsContent>

          {/* ── Pipeline ── */}
          <TabsContent value="pipeline">
            <PipelineApprovalList
              items={approvals?.pipeline?.items ?? []}
              total={approvals?.pipeline?.total ?? 0}
            />
          </TabsContent>
        </Tabs>
      </CardContent>
    </Card>
  );
}

// ============================================================================
// Visit Approval List
// ============================================================================

function VisitApprovalList({
  items,
  total,
}: {
  readonly items: DraftVisitItem[];
  readonly total: number;
}) {
  const queryClient = useQueryClient();
  const approveMutation = useApproveVisitReport();
  const rejectMutation = useRejectVisitReport();
  const [processingId, setProcessingId] = useState<string | null>(null);
  const [rejectTarget, setRejectTarget] = useState<string | null>(null);
  const [rejectReason, setRejectReason] = useState("");

  const handleApprove = async (id: string) => {
    setProcessingId(id);
    try {
      await approveMutation.mutateAsync(id);
      queryClient.invalidateQueries({
        queryKey: ["dashboard", "sales-manager", "team-draft-approvals"],
      });
      toast.success("Visit report approved successfully");
    } catch {
      // Error handled by api-client interceptor
    } finally {
      setProcessingId(null);
    }
  };

  const handleRejectConfirm = async () => {
    if (!rejectTarget) return;
    setProcessingId(rejectTarget);
    try {
      await rejectMutation.mutateAsync({
        id: rejectTarget,
        data: { reason: rejectReason.trim() || "Rejected by manager" },
      });
      queryClient.invalidateQueries({
        queryKey: ["dashboard", "sales-manager", "team-draft-approvals"],
      });
      toast.success("Visit report rejected");
    } catch {
      // Error handled by api-client interceptor
    } finally {
      setProcessingId(null);
      setRejectTarget(null);
      setRejectReason("");
    }
  };

  if (items.length === 0) {
    return <EmptyApprovalState label="visit reports" />;
  }

  return (
    <>
      {/* Reject reason dialog */}
      <Dialog
        open={!!rejectTarget}
        onOpenChange={(open) => {
          if (!open) {
            setRejectTarget(null);
            setRejectReason("");
          }
        }}
      >
        <DialogContent className="max-w-sm">
          <DialogHeader>
            <DialogTitle>Reject Visit Report</DialogTitle>
            <DialogDescription>
              Please provide a reason for rejecting this visit report.
            </DialogDescription>
          </DialogHeader>
          <div className="space-y-2 py-2">
            <Label htmlFor="reject-reason" className="text-sm">
              Reason
            </Label>
            <Textarea
              id="reject-reason"
              placeholder="Explain why this visit report is being rejected..."
              value={rejectReason}
              onChange={(e) => setRejectReason(e.target.value)}
              rows={3}
              className="resize-none text-sm"
            />
          </div>
          <DialogFooter className="gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => {
                setRejectTarget(null);
                setRejectReason("");
              }}
              className="cursor-pointer"
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              size="sm"
              onClick={handleRejectConfirm}
              disabled={rejectMutation.isPending}
              className="cursor-pointer"
            >
              {rejectMutation.isPending ? (
                <Loader2 className="h-3.5 w-3.5 animate-spin" />
              ) : (
                "Reject"
              )}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <div className="space-y-2">
        {total > items.length && (
          <MoreItemsBanner total={total} shown={items.length} href="/visit-reports" />
        )}
        {items.map((item) => {
          const isProcessing = processingId === item.id;
          return (
            <div
              key={item.id}
              className="flex items-start justify-between gap-3 rounded-lg border bg-card px-3 py-2.5 hover:bg-accent/30 transition-colors"
            >
              <div className="min-w-0 flex-1">
                <div className="flex items-center gap-1.5 flex-wrap">
                  <span className="text-xs sm:text-sm font-medium truncate">
                    {item.purpose}
                  </span>
                  <StatusBadge status={item.status} />
                </div>
                <div className="flex items-center gap-2 mt-0.5 flex-wrap text-[11px] text-muted-foreground">
                  <span>by {item.assigned_to}</span>
                  <span>·</span>
                  <span>{safeFormat(item.visit_date)}</span>
                </div>
              </div>
              <div className="flex items-center gap-1 shrink-0">
                <TooltipProvider>
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Button
                        size="icon"
                        variant="ghost"
                        className="h-7 w-7 text-green-600 hover:text-green-700 hover:bg-green-50 cursor-pointer"
                        onClick={() => handleApprove(item.id)}
                        disabled={isProcessing}
                      >
                        {isProcessing ? (
                          <Loader2 className="h-3.5 w-3.5 animate-spin" />
                        ) : (
                          <CheckCircle2 className="h-3.5 w-3.5" />
                        )}
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>Approve</TooltipContent>
                  </Tooltip>

                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Button
                        size="icon"
                        variant="ghost"
                        className="h-7 w-7 text-destructive hover:text-destructive hover:bg-destructive/10 cursor-pointer"
                        onClick={() => setRejectTarget(item.id)}
                        disabled={isProcessing}
                      >
                        <XCircle className="h-3.5 w-3.5" />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>Reject</TooltipContent>
                  </Tooltip>
                </TooltipProvider>

                <Link href={`/visit-reports/${item.id}`}>
                  <Button
                    size="icon"
                    variant="ghost"
                    className="h-7 w-7 text-muted-foreground hover:text-foreground cursor-pointer"
                  >
                    <ExternalLink className="h-3.5 w-3.5" />
                  </Button>
                </Link>
              </div>
            </div>
          );
        })}
      </div>
    </>
  );
}

// ============================================================================
// Task Approval List
// ============================================================================

function TaskApprovalList({
  items,
  total,
}: {
  readonly items: DraftTaskItem[];
  readonly total: number;
}) {
  const queryClient = useQueryClient();
  const hasCreatePermission = useHasPermission("tasks.create");
  const createTask = useCreateTask();
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);

  const handleCreateSubmit = async (data: CreateTaskFormData) => {
    try {
      await createTask.mutateAsync(data);
      setIsCreateDialogOpen(false);
      queryClient.invalidateQueries({
        queryKey: ["dashboard", "sales-manager", "team-draft-approvals"],
      });
      toast.success("Task created successfully");
    } catch {
      // Error handled by api-client interceptor
    }
  };

  return (
    <>
      {/* Create task dialog */}
      <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
        <DialogContent className="sm:max-w-[600px] max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Create Task</DialogTitle>
          </DialogHeader>
          <TaskForm
            onSubmit={(data) => handleCreateSubmit(data as CreateTaskFormData)}
            onCancel={() => setIsCreateDialogOpen(false)}
            isLoading={createTask.isPending}
          />
        </DialogContent>
      </Dialog>

      <div className="space-y-2">
        {hasCreatePermission && (
          <div className="flex justify-end pb-1">
            <Button
              size="sm"
              className="h-7 gap-1.5 text-xs cursor-pointer"
              onClick={() => setIsCreateDialogOpen(true)}
            >
              <Plus className="h-3.5 w-3.5" />
              Create Task
            </Button>
          </div>
        )}

        {items.length === 0 ? (
          <EmptyApprovalState label="tasks" />
        ) : (
          <>
            {total > items.length && (
              <MoreItemsBanner total={total} shown={items.length} href="/tasks" />
            )}
            {items.map((item) => (
              <div
                key={item.id}
                className="flex items-start justify-between gap-3 rounded-lg border bg-card px-3 py-2.5 hover:bg-accent/30 transition-colors"
              >
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-1.5 flex-wrap">
                    <span className="text-xs sm:text-sm font-medium truncate">{item.title}</span>
                    <PriorityBadge priority={item.priority} />
                  </div>
                  <div className="flex items-center gap-2 mt-0.5 flex-wrap text-[11px] text-muted-foreground">
                    <span>by {item.assigned_to}</span>
                    {item.due_date && (
                      <>
                        <span>·</span>
                        <span>due {safeFormat(item.due_date)}</span>
                      </>
                    )}
                  </div>
                </div>
                <Link href="/tasks">
                  <Button
                    size="icon"
                    variant="ghost"
                    className="h-7 w-7 text-muted-foreground hover:text-foreground shrink-0 cursor-pointer"
                  >
                    <ExternalLink className="h-3.5 w-3.5" />
                  </Button>
                </Link>
              </div>
            ))}
          </>
        )}
      </div>
    </>
  );
}

// ============================================================================
// Schedule Approval List
// ============================================================================

function ScheduleApprovalList({
  items,
  total,
}: {
  readonly items: DraftScheduleItem[];
  readonly total: number;
}) {
  if (items.length === 0) {
    return <EmptyApprovalState label="schedules" />;
  }

  return (
    <div className="space-y-2">
      {total > items.length && (
        <MoreItemsBanner total={total} shown={items.length} href="/schedules" />
      )}
      {items.map((item) => (
        <div
          key={item.id}
          className="flex items-start justify-between gap-3 rounded-lg border bg-card px-3 py-2.5 hover:bg-accent/30 transition-colors"
        >
          <div className="min-w-0 flex-1">
            <span className="text-xs sm:text-sm font-medium truncate block">{item.title}</span>
            <div className="flex items-center gap-2 mt-0.5 flex-wrap text-[11px] text-muted-foreground">
              <span>by {item.assigned_to}</span>
              <span>·</span>
              <span>{safeFormat(item.scheduled_at)}</span>
            </div>
          </div>
          <Link href="/schedules">
            <Button
              size="icon"
              variant="ghost"
              className="h-7 w-7 text-muted-foreground hover:text-foreground shrink-0 cursor-pointer"
            >
              <ExternalLink className="h-3.5 w-3.5" />
            </Button>
          </Link>
        </div>
      ))}
    </div>
  );
}

// ============================================================================
// Lead Approval List
// ============================================================================

function LeadApprovalList({
  items,
  total,
}: {
  readonly items: DraftLeadItem[];
  readonly total: number;
}) {
  if (items.length === 0) {
    return <EmptyApprovalState label="leads" />;
  }

  return (
    <div className="space-y-2">
      {total > items.length && (
        <MoreItemsBanner total={total} shown={items.length} href="/leads" />
      )}
      {items.map((item) => (
        <div
          key={item.id}
          className="flex items-start justify-between gap-3 rounded-lg border bg-card px-3 py-2.5 hover:bg-accent/30 transition-colors"
        >
          <div className="min-w-0 flex-1">
            <div className="flex items-center gap-1.5 flex-wrap">
              <span className="text-xs sm:text-sm font-medium truncate">{item.name}</span>
              <StatusBadge status={item.status} />
            </div>
            <div className="flex items-center gap-2 mt-0.5 flex-wrap text-[11px] text-muted-foreground">
              <span>{item.company}</span>
              <span>·</span>
              <span>by {item.assigned_to}</span>
            </div>
          </div>
          <Link href="/leads">
            <Button
              size="icon"
              variant="ghost"
              className="h-7 w-7 text-muted-foreground hover:text-foreground shrink-0 cursor-pointer"
            >
              <ExternalLink className="h-3.5 w-3.5" />
            </Button>
          </Link>
        </div>
      ))}
    </div>
  );
}

// ============================================================================
// Pipeline Approval List
// ============================================================================

function PipelineApprovalList({
  items,
  total,
}: {
  readonly items: DraftPipelineItem[];
  readonly total: number;
}) {
  if (items.length === 0) {
    return <EmptyApprovalState label="pipeline deals" />;
  }

  return (
    <div className="space-y-2">
      {total > items.length && (
        <MoreItemsBanner total={total} shown={items.length} href="/pipeline" />
      )}
      {items.map((item) => (
        <div
          key={item.id}
          className="flex items-start justify-between gap-3 rounded-lg border bg-card px-3 py-2.5 hover:bg-accent/30 transition-colors"
        >
          <div className="min-w-0 flex-1">
            <div className="flex items-center gap-1.5 flex-wrap">
              <span className="text-xs sm:text-sm font-medium truncate">{item.name}</span>
              {item.stage && (
                <Badge variant="outline" className="text-[10px] px-1.5 py-0">
                  {item.stage}
                </Badge>
              )}
            </div>
            <div className="flex items-center gap-2 mt-0.5 flex-wrap text-[11px] text-muted-foreground">
              <span>{formatCurrency(item.value / 100)}</span>
              <span>·</span>
              <span>by {item.assigned_to}</span>
            </div>
          </div>
          <Link href="/pipeline">
            <Button
              size="icon"
              variant="ghost"
              className="h-7 w-7 text-muted-foreground hover:text-foreground shrink-0 cursor-pointer"
            >
              <ExternalLink className="h-3.5 w-3.5" />
            </Button>
          </Link>
        </div>
      ))}
    </div>
  );
}

// ============================================================================
// Shared Sub-components
// ============================================================================

function EmptyApprovalState({ label }: { readonly label: string }) {
  return (
    <div className="flex flex-col items-center justify-center py-8 gap-2 text-muted-foreground">
      <CheckCircle2 className="h-8 w-8 text-green-500 opacity-70" />
      <p className="text-xs sm:text-sm font-medium">No pending {label}</p>
      <p className="text-[11px] text-center">
        All {label} have been reviewed or none submitted yet.
      </p>
    </div>
  );
}

function MoreItemsBanner({
  total,
  shown,
  href,
}: {
  readonly total: number;
  readonly shown: number;
  readonly href: string;
}) {
  return (
    <div className="flex items-center justify-between text-[11px] text-muted-foreground mb-2 px-1">
      <span>
        Showing {shown} of {total} items
      </span>
      <Link
        href={href}
        className="text-primary hover:underline font-medium cursor-pointer"
      >
        View all →
      </Link>
    </div>
  );
}

function StatusBadge({ status }: { readonly status: string }) {
  const map: Record<string, { label: string; className: string }> = {
    submitted: {
      label: "Submitted",
      className: "bg-blue-100 text-blue-700 border-blue-200",
    },
    draft: {
      label: "Draft",
      className: "bg-gray-100 text-gray-600 border-gray-200",
    },
    new: {
      label: "New",
      className: "bg-purple-100 text-purple-700 border-purple-200",
    },
    pending: {
      label: "Pending",
      className: "bg-yellow-100 text-yellow-700 border-yellow-200",
    },
    open: {
      label: "Open",
      className: "bg-blue-100 text-blue-700 border-blue-200",
    },
  };
  const cfg = map[status] ?? {
    label: status,
    className: "bg-muted text-muted-foreground border-border",
  };
  return (
    <Badge
      variant="outline"
      className={`text-[10px] font-medium px-1.5 py-0 ${cfg.className}`}
    >
      {cfg.label}
    </Badge>
  );
}

function PriorityBadge({ priority }: { readonly priority: string }) {
  const map: Record<string, { label: string; className: string }> = {
    urgent: {
      label: "Urgent",
      className: "bg-red-100 text-red-700 border-red-200",
    },
    high: {
      label: "High",
      className: "bg-orange-100 text-orange-700 border-orange-200",
    },
    medium: {
      label: "Medium",
      className: "bg-yellow-100 text-yellow-700 border-yellow-200",
    },
    low: {
      label: "Low",
      className: "bg-green-100 text-green-700 border-green-200",
    },
  };
  const cfg = map[priority];
  if (!cfg) return null;
  return (
    <Badge
      variant="outline"
      className={`text-[10px] font-medium px-1.5 py-0 ${cfg.className}`}
    >
      {cfg.label}
    </Badge>
  );
}

function safeFormat(dateStr: string | undefined | null): string {
  if (!dateStr) return "—";
  try {
    return format(new Date(dateStr), "dd MMM yyyy");
  } catch {
    return "—";
  }
}

// ============================================================================
// Loading & Error States
// ============================================================================

function TeamDraftApprovalsSkeleton() {
  return (
    <Card className="border-0 shadow-sm h-full flex flex-col">
      <CardHeader className="pb-3 px-4 sm:px-6 pt-4 sm:pt-6">
        <Skeleton className="h-5 w-48" />
        <Skeleton className="h-3 w-64 mt-2" />
      </CardHeader>
      <CardContent className="px-4 sm:px-6 pb-4 sm:pb-6 flex-1 overflow-y-auto max-h-[520px]">
        <Skeleton className="h-8 w-full mb-4" />
        <div className="space-y-2">
          {Array.from({ length: 3 }).map((_, i) => (
            <Skeleton key={i} className="h-14 w-full rounded-lg" />
          ))}
        </div>
      </CardContent>
    </Card>
  );
}

function TeamDraftApprovalsError() {
  return (
    <Card className="border-0 shadow-sm h-full flex flex-col">
      <CardContent className="px-4 sm:px-6 py-8 flex flex-col items-center gap-2 text-muted-foreground flex-1 overflow-y-auto max-h-[520px]">
        <AlertCircle className="h-6 w-6 text-destructive" />
        <p className="text-xs">Failed to load team approvals</p>
      </CardContent>
    </Card>
  );
}
