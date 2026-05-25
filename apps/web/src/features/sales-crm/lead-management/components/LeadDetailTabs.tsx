'use client';

import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { LeadQualificationCard } from './LeadQualificationCard';
import { useLead, useLeadVisitReports, useLeadActivities } from '../hooks/useLeads';
import type { VisitReport } from '@/features/sales-crm/visit-report/types';
import type { Activity } from '@/features/sales-crm/visit-report/types/activity';
import {
  ClipboardList,
  CheckSquare,
  MapPin,
  Activity as ActivityIcon,
  GitCommitHorizontal,
  Eye,
  Info,
  Mail,
  Pencil,
  Phone,
  Building2,
  Globe,
  User,
  Calendar,
  Plus
} from 'lucide-react';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { TaskForm } from '@/features/sales-crm/task-management/components/task-form';
import { useCreateTask } from '@/features/sales-crm/task-management/hooks/useTasks';
import type { CreateTaskFormData, UpdateTaskFormData } from '@/features/sales-crm/task-management/schemas/task.schema';
import { VisitReportDetailModal } from '@/features/sales-crm/visit-report/components/visit-report-detail-modal';
import { VisitReportForm } from '@/features/sales-crm/visit-report/components/visit-report-form';
import { useCreateVisitReport, useUpdateVisitReport } from '@/features/sales-crm/visit-report/hooks/useVisitReports';
import type { CreateVisitReportFormData, UpdateVisitReportFormData } from '@/features/sales-crm/visit-report/schemas/visit-report.schema';
import { CreateActivityDialog } from '@/features/sales-crm/visit-report/components/create-activity-dialog';
import { toast } from 'sonner';

interface LeadDetailTabsProps {
  readonly leadId: string;
}

export function LeadDetailTabs({ leadId }: LeadDetailTabsProps) {
  const { data: leadData } = useLead(leadId);
  const { data: visitReportsData, refetch: refetchVisits } = useLeadVisitReports(leadId);
  const { data: activitiesData, refetch: refetchActivities } = useLeadActivities(leadId);

  const createTask = useCreateTask();
  const createVisitReport = useCreateVisitReport();
  const updateVisitReport = useUpdateVisitReport();

  const [isTaskModalOpen, setIsTaskModalOpen] = useState(false);
  const [isVisitModalOpen, setIsVisitModalOpen] = useState(false);
  const [isActivityModalOpen, setIsActivityModalOpen] = useState(false);
  const [editingVisit, setEditingVisit] = useState<VisitReport | null>(null);
  const [viewingVisitReportId, setViewingVisitReportId] = useState<string | null>(null);

  const lead = leadData?.data;
  const visitReports = visitReportsData?.data ?? [];
  const visits = Array.isArray(visitReports) ? visitReports : [];
  const activitiesRaw = activitiesData?.data ?? [];
  const activities = Array.isArray(activitiesRaw) ? activitiesRaw : [];
  const timelineItems = buildTimelineItems(visits, activities);

  const handleCreateTask = async (data: CreateTaskFormData | UpdateTaskFormData) => {
    try {
      // Cast the payload because we add lead_id mapping manually to the request
      await createTask.mutateAsync({ ...data, lead_id: leadId } as any);
      toast.success("Task created successfully");
      setIsTaskModalOpen(false);
    } catch {
      // Extractor interceptor handles it
    }
  };

  const handleCreateVisit = async (data: CreateVisitReportFormData | UpdateVisitReportFormData) => {
    try {
      await createVisitReport.mutateAsync(data as CreateVisitReportFormData);
      toast.success("Visit report created successfully");
      setIsVisitModalOpen(false);
      refetchVisits();
      refetchActivities();
    } catch {
      // Extractor interceptor handles it
    }
  };

  const handleUpdateVisit = async (data: CreateVisitReportFormData | UpdateVisitReportFormData) => {
    if (!editingVisit) return;

    try {
      await updateVisitReport.mutateAsync({
        id: editingVisit.id,
        data: data as UpdateVisitReportFormData,
      });
      toast.success("Visit report updated successfully");
      setEditingVisit(null);
      refetchVisits();
      refetchActivities();
    } catch {
      // Extractor interceptor handles it
    }
  };

  if (!lead) {
    return (
      <Card>
        <CardContent className="text-center py-8">
          <p className="text-sm text-muted-foreground">Loading lead details...</p>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {/* Contact Information */}
        <Card className="h-full">
          <CardHeader className="py-4">
            <CardTitle className="text-base flex items-center gap-2">
              <User size={16} /> Contact Information
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <DetailItem icon={<User size={14} />} label="Name" value={`${lead.first_name ?? ''} ${lead.last_name ?? ''}`} />
              <DetailItem icon={<Mail size={14} />} label="Email" value={lead.email ?? '—'} />
              <DetailItem icon={<Phone size={14} />} label="Phone" value={lead.phone ?? '—'} />
              <DetailItem icon={<Building2 size={14} />} label="Company" value={lead.company_name ?? '—'} />
              <DetailItem icon={<User size={14} />} label="Job Title" value={lead.job_title ?? '—'} />
              <DetailItem icon={<Globe size={14} />} label="Website" value={lead.website ?? '—'} />
            </div>
          </CardContent>
        </Card>

        {/* Location & SummaryInfo */}
        <div className="flex flex-col gap-4 h-full">
          <Card>
            <CardHeader className="py-4">
              <CardTitle className="text-base flex items-center gap-2">
                <MapPin size={16} /> Location
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <DetailItem label="Address" value={lead.address ?? '—'} />
                <DetailItem label="City" value={lead.city ?? '—'} />
                <DetailItem label="Province" value={lead.province ?? '—'} />
                <DetailItem label="Country" value={lead.country ?? '—'} />
                <DetailItem label="Postal Code" value={lead.postal_code ?? '—'} />
              </div>
            </CardContent>
          </Card>

          <Card className="flex-1">
            <CardHeader className="py-4">
              <CardTitle className="text-base flex items-center gap-2">
                <Info size={16} /> Lead Information
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <DetailItem label="Source" value={lead.lead_source ?? '—'} />
                <DetailItem label="Industry" value={lead.industry ?? '—'} />
                <DetailItem label="Score" value={`${lead.lead_score ?? 0}/100`} />
                <DetailItem label="Assigned To" value={lead.assigned_user?.name ?? '—'} />
              </div>
              {lead.notes && (
                <div className="mt-4 pt-4 border-t">
                  <h3 className="font-medium mb-1 text-xs text-muted-foreground">Notes</h3>
                  <p className="text-sm whitespace-pre-wrap">{lead.notes}</p>
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>

      {/* Qualification (BANT) Section (content hidden) */}

      {/* Tabbed: BANT / Log / Task */}
      <Tabs defaultValue="bant">
        <TabsList className="grid grid-cols-3 gap-2">
          <TabsTrigger value="bant" className="text-sm">
            <ClipboardList size={16} /> BANT
          </TabsTrigger>
          <TabsTrigger value="log" className="text-sm">
            <GitCommitHorizontal size={16} /> Activity Log
          </TabsTrigger>
          <TabsTrigger value="task" className="text-sm">
            <CheckSquare size={16} /> Task
          </TabsTrigger>
        </TabsList>

        <TabsContent value="bant" className="mt-4">
          <LeadQualificationCard leadId={leadId} />
        </TabsContent>

        <TabsContent value="log" className="mt-4">
          <div className="mb-5 flex flex-wrap items-center justify-end gap-2">
            <Button size="sm" variant="outline" className="rounded-full px-4" onClick={() => setIsVisitModalOpen(true)}>
              <Plus className="h-4 w-4 mr-1" /> Add Visit
            </Button>
            <Button size="sm" className="rounded-full px-4" onClick={() => setIsActivityModalOpen(true)}>
              <Plus className="h-4 w-4 mr-1" /> Add Activity
            </Button>
          </div>
          {timelineItems.length === 0 ? (
            <div className="rounded-2xl border border-dashed border-border/70 bg-card/40 py-12 text-center">
              <ActivityIcon className="mx-auto h-8 w-8 text-muted-foreground/30 mb-2" />
              <p className="text-sm text-muted-foreground">No visits or activities linked to this lead yet.</p>
            </div>
          ) : (
            <div className="rounded-3xl border border-border/70 bg-card/40 p-5 shadow-[0_18px_40px_-32px_rgba(15,23,42,0.45)]">
              <div className="space-y-1">
                {timelineItems.map((item, index) => (
                  <div key={`${item.kind}-${item.id}`} className="relative pl-12">
                    {index !== timelineItems.length - 1 && (
                      <div className="absolute left-[19px] top-10 bottom-[-12px] w-px bg-border/80" />
                    )}
                    <div className="absolute left-0 top-1 flex h-10 w-10 items-center justify-center rounded-full border border-border/70 bg-background shadow-sm">
                      {item.kind === 'visit' ? (
                        <MapPin className="h-4 w-4 text-sky-600" />
                      ) : (
                        <ActivityIcon className="h-4 w-4 text-amber-600" />
                      )}
                    </div>
                    <div className="rounded-2xl border border-border/70 bg-card/80 p-5">
                      <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
                        <div className="min-w-0 space-y-3">
                          <div className="flex flex-wrap items-center gap-2">
                            <span className="rounded-full border border-border/70 bg-background/80 px-2.5 py-1 text-[11px] font-medium uppercase tracking-[0.16em] text-muted-foreground">
                              {item.kind === 'visit' ? 'Visit' : formatActivityType(item.activity?.type)}
                            </span>
                            <span className="text-xs text-muted-foreground">
                              {formatSafeDate(item.dateValue, true)}
                            </span>
                          </div>
                          <div className="space-y-2">
                            <h3 className="text-base font-semibold tracking-tight text-foreground">
                              {item.title}
                            </h3>
                            <p className="max-w-3xl text-sm leading-6 text-muted-foreground whitespace-pre-wrap">
                              {item.description}
                            </p>
                          </div>
                        </div>
                        {item.kind === 'visit' && item.visit ? (
                          <div className="flex items-center gap-2 self-start">
                            <Button
                              size="sm"
                              variant="outline"
                              className="rounded-full px-4"
                              onClick={() => setViewingVisitReportId(item.visit?.id ?? null)}
                            >
                              <Eye className="mr-1 h-4 w-4" />
                              Detail
                            </Button>
                            <Button
                              size="sm"
                              variant="ghost"
                              className="rounded-full px-4"
                              onClick={() => setEditingVisit(item.visit ?? null)}
                              disabled={item.visit.status !== 'draft'}
                            >
                              <Pencil className="mr-1 h-4 w-4" />
                              Edit
                            </Button>
                          </div>
                        ) : null}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </TabsContent>

        <TabsContent value="task" className="mt-4">
          <div className="mb-4 flex items-center justify-end">
            <Button size="sm" className="rounded-full px-4" onClick={() => setIsTaskModalOpen(true)}>
              <Plus className="h-4 w-4 mr-1" /> Add Task
            </Button>
          </div>
          <div className="rounded-2xl border border-dashed border-border/70 bg-card/40 py-12 text-center">
            <CheckSquare className="mx-auto mb-3 h-8 w-8 text-muted-foreground/30" />
            <p className="text-sm font-medium text-foreground">Task timeline is being prepared.</p>
            <p className="mt-1 text-sm text-muted-foreground">
              Tasks functionality for leads is being expanded. Currently no tasks synced.
            </p>
          </div>
        </TabsContent>
      </Tabs>

      {/* Task Creation Modal */}
      <Dialog open={isTaskModalOpen} onOpenChange={setIsTaskModalOpen}>
        <DialogContent className="sm:max-w-[600px] max-h-[90vh] overflow-y-auto">
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

      {/* Visit Report Creation Modal */}
      <Dialog open={isVisitModalOpen} onOpenChange={setIsVisitModalOpen}>
        <DialogContent className="sm:max-w-xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Create Visit Report</DialogTitle>
          </DialogHeader>
          <VisitReportForm
            onSubmit={handleCreateVisit}
            onCancel={() => setIsVisitModalOpen(false)}
            isLoading={createVisitReport.isPending}
            initialLeadId={leadId}
          />
        </DialogContent>
      </Dialog>

      <Dialog open={!!editingVisit} onOpenChange={(open) => {
        if (!open) setEditingVisit(null);
      }}>
        <DialogContent className="sm:max-w-xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Edit Visit Report</DialogTitle>
          </DialogHeader>
          {editingVisit && (
            <VisitReportForm
              visitReport={editingVisit}
              onSubmit={handleUpdateVisit}
              onCancel={() => setEditingVisit(null)}
              isLoading={updateVisitReport.isPending}
            />
          )}
        </DialogContent>
      </Dialog>

      {/* Activity Creation Modal */}
      <CreateActivityDialog
        open={isActivityModalOpen}
        onOpenChange={setIsActivityModalOpen}
        leadId={leadId}
        onSuccess={() => refetchActivities()}
      />

      <VisitReportDetailModal
        visitReportId={viewingVisitReportId}
        open={!!viewingVisitReportId}
        onOpenChange={(open) => {
          if (!open) setViewingVisitReportId(null);
        }}
        onVisitReportUpdated={() => {
          refetchVisits();
          refetchActivities();
        }}
      />
    </div>
  );
}

// Helper components
function DetailItem({
  icon,
  label,
  value,
}: {
  readonly icon?: React.ReactNode;
  readonly label: string;
  readonly value: string;
}) {
  return (
    <div className="flex items-start gap-2">
      {icon && (
        <span className="text-muted-foreground mt-0.5 shrink-0">{icon}</span>
      )}
      <div>
        <p className="text-xs text-muted-foreground">{label}</p>
        <p className="text-sm">{value}</p>
      </div>
    </div>
  );
}

function formatSafeDate(
  value?: string | null,
  includeTime = false,
): string {
  if (!value) return '—';
  const date = new Date(value);
  if (isNaN(date.getTime())) return '—';
  if (includeTime) {
    return date.toLocaleString('id-ID', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  }
  return date.toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
}

function formatActivityType(value?: string | null): string {
  if (!value) return 'Activity';
  return value
    .split('_')
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join(' ');
}

type LeadTimelineItem =
  | {
      kind: 'visit';
      id: string;
      dateValue: string;
      title: string;
      description: string;
      visit: VisitReport;
      activity?: undefined;
    }
  | {
      kind: 'activity';
      id: string;
      dateValue: string;
      title: string;
      description: string;
      visit?: undefined;
      activity: Activity;
    };

function buildTimelineItems(visits: VisitReport[], activities: Activity[]): LeadTimelineItem[] {
  const visitItems: LeadTimelineItem[] = visits.map((visit) => ({
    kind: 'visit',
    id: visit.id,
    dateValue: visit.visit_date || visit.created_at || visit.updated_at,
    title: visit.purpose || 'Visit',
    description: visit.notes?.trim() || 'No additional notes for this visit.',
    visit,
  }));

  const activityItems: LeadTimelineItem[] = activities.map((activity) => {
    let dateValue = activity.metadata && typeof activity.metadata === "object"
      ? (() => {
          const meta = activity.metadata as Record<string, unknown>;
          if (typeof meta.visit_date === "string") {
            return meta.visit_date;
          }
          return activity.timestamp || activity.created_at || activity.updated_at;
        })()
      : activity.timestamp || activity.created_at || activity.updated_at;
    if (activity.type === "visit" && activity.metadata && typeof activity.metadata === "object") {
      const meta = activity.metadata as Record<string, unknown>;
      if (typeof meta.visit_date === "string") {
        dateValue = meta.visit_date;
      }
    }
    return {
      kind: 'activity',
      id: activity.id ?? `${activity.type}-${activity.created_at}`,
      dateValue,
      title: formatActivityType(activity.type),
      description: activity.description || 'No activity description.',
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
