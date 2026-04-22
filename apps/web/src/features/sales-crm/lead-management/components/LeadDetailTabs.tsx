'use client';

import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { LeadQualificationCard } from './LeadQualificationCard';
import { useLead, useLeadVisitReports, useLeadActivities } from '../hooks/useLeads';
import type { VisitReport } from '@/features/sales-crm/visit-report/types';
import type { Activity } from '@/features/sales-crm/visit-report/types/activity';
import {
  ClipboardList,
  CheckSquare,
  MapPin,
  Activity as ActivityIcon,
  Info,
  Mail,
  Phone,
  Building2,
  Globe,
  User,
  Calendar,
  Plus
} from 'lucide-react';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { TaskForm } from '@/features/sales-crm/task-management/components/task-form';
import { useCreateTask } from '@/features/sales-crm/task-management/hooks/useTasks';
import type { CreateTaskFormData, UpdateTaskFormData } from '@/features/sales-crm/task-management/schemas/task.schema';
import { VisitReportForm } from '@/features/sales-crm/visit-report/components/visit-report-form';
import { useCreateVisitReport } from '@/features/sales-crm/visit-report/hooks/useVisitReports';
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

  const [isTaskModalOpen, setIsTaskModalOpen] = useState(false);
  const [isVisitModalOpen, setIsVisitModalOpen] = useState(false);
  const [isActivityModalOpen, setIsActivityModalOpen] = useState(false);

  const lead = leadData?.data;
  const visitReports = visitReportsData?.data ?? [];
  const visits = Array.isArray(visitReports) ? visitReports : [];
  const activitiesRaw = activitiesData?.data ?? [];
  const activities = Array.isArray(activitiesRaw) ? activitiesRaw : [];

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

      {/* Qualification (BANT) Section */}
      <h2 className="text-lg font-semibold flex items-center gap-2 pt-2 border-b pb-2">
        <ClipboardList size={18} /> Lead Qualification (BANT)
      </h2>
      <div>
        <LeadQualificationCard leadId={leadId} />
      </div>

      {/* Tasks Section */}
      <div className="flex items-center justify-between pt-2 border-b pb-2">
        <h2 className="text-lg font-semibold flex items-center gap-2">
          <CheckSquare size={18} /> Related Tasks
        </h2>
        <Button size="sm" onClick={() => setIsTaskModalOpen(true)}>
          <Plus className="h-4 w-4 mr-1" /> Add Task
        </Button>
      </div>
      <Card>
        <CardContent className="p-0 overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Title</TableHead>
                <TableHead>Type</TableHead>
                <TableHead>Priority</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Due Date</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow>
                <TableCell colSpan={5} className="text-center py-6 text-sm text-muted-foreground">
                  <div className="flex flex-col items-center justify-center">
                    <CheckSquare className="h-8 w-8 text-muted-foreground/30 mb-2" />
                    Tasks functionality for leads is being expanded. Currently no tasks synced.
                  </div>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      {/* Visits Section */}
      <div className="flex items-center justify-between pt-2 border-b pb-2">
        <h2 className="text-lg font-semibold flex items-center gap-2">
          <MapPin size={18} /> Visit Reports
        </h2>
        <Button size="sm" onClick={() => setIsVisitModalOpen(true)}>
          <Plus className="h-4 w-4 mr-1" /> Add Visit
        </Button>
      </div>
      <Card>
        <CardContent className="p-0 overflow-x-auto">
          {visits.length === 0 ? (
            <div className="text-center py-10">
              <MapPin className="mx-auto h-8 w-8 text-muted-foreground/30 mb-2" />
              <p className="text-sm text-muted-foreground">No visit reports linked to this lead yet.</p>
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Purpose</TableHead>
                  <TableHead>Date</TableHead>
                  <TableHead>Status</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {visits.map((visit: VisitReport) => (
                  <TableRow key={visit.id}>
                    <TableCell className="font-medium">{visit.purpose ?? 'Visit Report'}</TableCell>
                    <TableCell>{formatSafeDate(visit.visit_date ?? visit.created_at)}</TableCell>
                    <TableCell>
                      {visit.status && (
                        <Badge variant="outline" className="text-xs">
                          {visit.status}
                        </Badge>
                      )}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

      {/* Activities Section */}
      <div className="flex items-center justify-between pt-2 border-b pb-2">
        <h2 className="text-lg font-semibold flex items-center gap-2">
          <ActivityIcon size={18} /> Activities
        </h2>
        <Button size="sm" onClick={() => setIsActivityModalOpen(true)}>
          <Plus className="h-4 w-4 mr-1" /> Add Activity
        </Button>
      </div>
      <Card>
        <CardContent className="p-0 overflow-x-auto">
          {activities.length === 0 ? (
            <div className="text-center py-10">
              <ActivityIcon className="mx-auto h-8 w-8 text-muted-foreground/30 mb-2" />
              <p className="text-sm text-muted-foreground">No activities logged for this lead yet.</p>
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Type</TableHead>
                  <TableHead>Description</TableHead>
                  <TableHead>Date</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {activities.map((act: Activity) => (
                  <TableRow key={act.id ?? ''}>
                    <TableCell className="font-medium">{act.type ?? 'Activity'}</TableCell>
                    <TableCell className="max-w-md truncate">{act.description ?? '-'}</TableCell>
                    <TableCell>{formatSafeDate(act.created_at ?? '', true)}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

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

      {/* Activity Creation Modal */}
      <CreateActivityDialog
        open={isActivityModalOpen}
        onOpenChange={setIsActivityModalOpen}
        leadId={leadId}
        onSuccess={() => refetchActivities()}
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
