'use client';

import { type ReactNode, useMemo, useState } from 'react';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Switch } from '@/components/ui/switch';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Textarea } from '@/components/ui/textarea';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { DealHistoryTimeline } from './deal-history-timeline';
import { ProductInterestTab } from '@/features/sales-crm/visit-report/components/product-interest-tab';
import { LeadQualificationCard } from '@/features/sales-crm/lead-management/components/LeadQualificationCard';
import { CreateActivityDialog } from '@/features/sales-crm/visit-report/components/create-activity-dialog';
import { VisitReportForm } from '@/features/sales-crm/visit-report/components/visit-report-form';
import type { VisitReport } from '@/features/sales-crm/visit-report/types';
import type { Activity } from '@/features/sales-crm/visit-report/types/activity';
import { VisitReportDetailModal } from '@/features/sales-crm/visit-report/components/visit-report-detail-modal';
import { useCreateVisitReport } from '@/features/sales-crm/visit-report/hooks/useVisitReports';
import {
  useDeal,
  useDealVisitReports,
  useDealActivities,
  useDealHistory,
  useUpdateDeal,
} from '../hooks/useDeals';
import { useLead } from '@/features/sales-crm/lead-management/hooks/useLeads';
import { useLeadQualification } from '@/features/sales-crm/lead-management/hooks/useLeadQualification';
import { useCreateTask, useTasksByDeal } from '@/features/sales-crm/task-management/hooks/useTasks';
import { TaskDetailModal } from '@/features/sales-crm/task-management/components/task-detail-modal';
import { TaskForm } from '@/features/sales-crm/task-management/components/task-form';
import { useHasPermission } from '@/features/auth/providers/permissions-provider';
import { formatCurrency } from '../utils/format';
import type { Deal } from '../types';
import type { CreateVisitReportFormData, UpdateVisitReportFormData } from '@/features/sales-crm/visit-report/schemas/visit-report.schema';
import type { CreateTaskFormData, UpdateTaskFormData } from '@/features/sales-crm/task-management/schemas/task.schema';
import type { DealUpdateData } from '../schemas/deal.schema';
import {
  Info,
  Package,
  CheckSquare,
  MapPin,
  Activity as ActivityIcon,
  History,
  ClipboardList,
  MapPinned,
  Plus,
} from 'lucide-react';

interface DealDetailTabsProps {
  readonly dealId: string;
}

export function DealDetailTabs({ dealId }: DealDetailTabsProps) {
  const [activeTab, setActiveTab] = useState('activities');
  const [isTaskModalOpen, setIsTaskModalOpen] = useState(false);
  const [isVisitModalOpen, setIsVisitModalOpen] = useState(false);
  const [isActivityModalOpen, setIsActivityModalOpen] = useState(false);
  const [viewingVisitReportId, setViewingVisitReportId] = useState<string | null>(null);
  const [viewingTaskId, setViewingTaskId] = useState<string | null>(null);
  const { data: deal, isLoading } = useDeal(dealId);
  const { data: leadData } = useLead(deal?.lead_id ?? '');
  const { qualification } = useLeadQualification(deal?.lead_id ?? '');
  const { data: visitReports, refetch: refetchVisitReports, isLoading: isVisitReportsLoading } = useDealVisitReports(dealId, 1, 100);
  const { data: activities, refetch: refetchActivities, isLoading: isActivitiesLoading } = useDealActivities(dealId, 1, 100);
  const { data: history } = useDealHistory(dealId);
  const { data: dealTasksResponse, refetch: refetchTasks } = useTasksByDeal(dealId, 1, 20);
  const createVisitReport = useCreateVisitReport();
  const createTask = useCreateTask();
  const hasCreateTaskPermission = useHasPermission('tasks.create');
  const hasVisitReportCreatePermission = useHasPermission('visit-reports.create');
  const hasCreateActivityPermission = useHasPermission('activities.create') || hasVisitReportCreatePermission;
  const visits = useMemo(() => (Array.isArray(visitReports) ? visitReports : []), [visitReports]);
  const acts = useMemo(() => (Array.isArray(activities) ? activities : []), [activities]);
  const timelineItems = useMemo(() => buildTimelineItems(visits, acts), [visits, acts]);
  const productInterestActivities = useMemo(
    () => buildProductInterestActivities(
      dedupeVisitReports(visits),
      dedupeActivities(acts),
    ),
    [visits, acts],
  );

  if (isLoading) {
    return <Skeleton className="h-[500px] w-full rounded-lg" />;
  }

  if (!deal) {
    return (
      <Card>
        <CardContent className="text-center py-8">
          <p className="text-sm text-muted-foreground">Deal not found</p>
        </CardContent>
      </Card>
    );
  }

  const lead = leadData?.data;
  const productItems = deal.product_items ?? [];
  const productInterests = qualification?.need_target_products ?? deal.qualification_snapshot?.need_target_products ?? [];
  const historyItems = Array.isArray(history) ? history : [];
  const dealTasks = Array.isArray(dealTasksResponse?.data) ? dealTasksResponse.data : [];
  const tasks = dealTasks.filter((task) => task.deal_id === dealId);
  const contactName = lead?.contact?.name ?? deal.contact?.name ?? '—';
  const contactEmail = lead?.contact?.email ?? deal.contact?.email ?? lead?.email ?? '—';
  const contactPhone = lead?.contact?.phone ?? deal.contact?.phone ?? lead?.phone ?? '—';
  const companyName = lead?.company_name ?? deal.account?.name ?? 'N/A';
  const address = [lead?.address, lead?.city, lead?.province, lead?.postal_code, lead?.country]
    .filter(Boolean)
    .join(', ');
  const budgetValue = qualification?.budget_target_amount ?? deal.qualification_snapshot?.budget_target_amount;
  const authorityValue = qualification?.authority_target_person ?? deal.qualification_snapshot?.authority_target_person;
  const needValue = qualification?.need_notes ?? deal.qualification_snapshot?.need_notes;
  const timelineValue = qualification?.timeline_target_date ?? deal.qualification_snapshot?.timeline_target_date;
  const productInterestCount = getProductInterestCount(productInterestActivities);
  const isProductInterestLoading = isVisitReportsLoading || isActivitiesLoading;

  const handleCreateVisit = async (formData: CreateVisitReportFormData | UpdateVisitReportFormData) => {
    const createdVisit = await createVisitReport.mutateAsync({
      ...(formData as CreateVisitReportFormData),
      deal_id: dealId,
      account_id: deal.account_id,
      contact_id: deal.contact_id || undefined,
    });
    setIsVisitModalOpen(false);
    await Promise.all([refetchVisitReports(), refetchActivities()]);
    setViewingVisitReportId(createdVisit.data.id);
  };

  const handleCreateTask = async (formData: CreateTaskFormData | UpdateTaskFormData) => {
    await createTask.mutateAsync({
      ...(formData as Record<string, unknown>),
      deal_id: dealId,
      account_id: deal.account_id,
      contact_id: deal.contact_id || undefined,
      assigned_to: (formData as CreateTaskFormData).assigned_to || deal.assigned_to || undefined,
    } as CreateTaskFormData);
    setIsTaskModalOpen(false);
    await refetchTasks();
  };

  return (
    <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
      <TabsList className="h-auto w-full justify-start gap-4 overflow-x-auto rounded-none border-b border-border bg-transparent p-0">
        <TabsTrigger value="activities" className="cursor-pointer gap-2 whitespace-nowrap rounded-none border-b-2 border-transparent px-4 py-3 text-sm text-muted-foreground data-[state=active]:border-primary data-[state=active]:bg-transparent data-[state=active]:text-foreground data-[state=active]:shadow-none">
          <ActivityIcon size={16} />
          <span>Activities</span>
        </TabsTrigger>
        <TabsTrigger value="tasks" className="cursor-pointer gap-2 whitespace-nowrap rounded-none border-b-2 border-transparent px-4 py-3 text-sm text-muted-foreground data-[state=active]:border-primary data-[state=active]:bg-transparent data-[state=active]:text-foreground data-[state=active]:shadow-none">
          <CheckSquare size={16} />
          <span>Tasks</span>
        </TabsTrigger>
        <TabsTrigger value="products" className="cursor-pointer gap-2 whitespace-nowrap rounded-none border-b-2 border-transparent px-4 py-3 text-sm text-muted-foreground data-[state=active]:border-primary data-[state=active]:bg-transparent data-[state=active]:text-foreground data-[state=active]:shadow-none">
          <Package size={16} />
          <span>Product Interest</span>
        </TabsTrigger>
        <TabsTrigger value="bant" className="cursor-pointer gap-2 whitespace-nowrap rounded-none border-b-2 border-transparent px-4 py-3 text-sm text-muted-foreground data-[state=active]:border-primary data-[state=active]:bg-transparent data-[state=active]:text-foreground data-[state=active]:shadow-none">
          <ClipboardList size={16} />
          <span>BANT</span>
        </TabsTrigger>
        <TabsTrigger value="information" className="cursor-pointer gap-2 whitespace-nowrap rounded-none border-b-2 border-transparent px-4 py-3 text-sm text-muted-foreground data-[state=active]:border-primary data-[state=active]:bg-transparent data-[state=active]:text-foreground data-[state=active]:shadow-none">
          <Info size={16} />
          <span>Information</span>
        </TabsTrigger>
      </TabsList>

      <TabsContent value="activities" className="mt-4">
        <TabCard
          title="Activities"
          description="Track calls, visits, meetings, and follow ups for this deal."
          actions={
            <>
              {hasVisitReportCreatePermission && (
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
              icon={<ActivityIcon className="h-10 w-10 text-muted-foreground/40" />}
              title="No activities yet"
              description="Activities and visit reports linked to this deal will appear here."
            />
          ) : (
            <div className="space-y-4">
              {timelineItems.map((item, index) => (
                <div key={item.id} className="flex gap-3">
                  <div className="relative flex flex-col items-center">
                    <div className="flex h-8 w-8 items-center justify-center rounded-full border border-primary/30 bg-primary/10 text-primary">
                      {item.kind === 'visit' ? <MapPin className="h-4 w-4" /> : <ActivityIcon className="h-4 w-4" />}
                    </div>
                    {index !== timelineItems.length - 1 && <div className="mt-1 h-full w-px flex-1 bg-border/70" />}
                  </div>
                  <div
                    className={`min-w-0 flex-1 rounded-xl border border-border bg-muted/30 p-4 ${
                      item.kind === 'visit' ? 'cursor-pointer transition-colors hover:bg-accent/40' : ''
                    }`}
                    role={item.kind === 'visit' ? 'button' : undefined}
                    tabIndex={item.kind === 'visit' ? 0 : undefined}
                    onClick={item.kind === 'visit' && item.visitId ? () => setViewingVisitReportId(item.visitId ?? null) : undefined}
                    onKeyDown={
                      item.kind === 'visit' && item.visitId
                        ? (event) => {
                          if (event.key === 'Enter' || event.key === ' ') {
                            event.preventDefault();
                            setViewingVisitReportId(item.visitId ?? null);
                          }
                        }
                        : undefined
                    }
                  >
                    <div className="flex flex-wrap items-center gap-2">
                      <Badge variant="outline" className="capitalize">
                        {item.kind === 'visit' ? 'Visit' : item.type}
                      </Badge>
                      <span className="text-xs text-muted-foreground">{formatSafeDate(item.dateValue, true)}</span>
                      {item.ownerName && <span className="text-xs text-muted-foreground">• {item.ownerName}</span>}
                    </div>
                    <p className="mt-2 whitespace-pre-wrap text-sm leading-relaxed text-foreground">
                      {item.description}
                    </p>
                    {item.accountName && (
                      <p className="mt-2 text-xs text-muted-foreground">{item.accountName}</p>
                    )}
                    {item.kind === 'visit' && item.visitId && (
                      <div className="mt-3">
                        <Button
                          size="sm"
                          variant="outline"
                          className="cursor-pointer"
                          onClick={(event) => {
                            event.stopPropagation();
                            setViewingVisitReportId(item.visitId ?? null);
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

      {/* Products Tab */}
      <TabsContent value="products" className="mt-4">
        <TabCard
          title="Product Interest"
          description="Products captured as customer interest from activities, visits, and lead qualification."
          actions={
            <>
              {hasVisitReportCreatePermission && (
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
          {productInterests.length === 0 && productInterestCount === 0 ? (
            isProductInterestLoading ? (
              <ProductInterestTab
                activities={productInterestActivities}
                isLoading={isProductInterestLoading}
              />
            ) : (
              <EmptyState
                icon={<Package className="h-10 w-10 text-muted-foreground/40" />}
                title="No product interest recorded"
                description="Product interest captured from qualification, activities, or visits will appear here."
              />
            )
          ) : (
            <div className="space-y-4">
              {productInterestCount > 0 && (
                <ProductInterestTab
                  activities={productInterestActivities}
                  isLoading={isProductInterestLoading}
                />
              )}

              {productInterests.length > 0 && (
                <div className="space-y-3">
                  <div className="rounded-xl border border-dashed border-border bg-background/70 p-4">
                      <p className="text-sm font-medium">BANT Need Products</p>
                      <p className="mt-1 text-xs text-muted-foreground">
                        Products from the source lead qualification are interest data only.
                      </p>
                  </div>
                  {productInterests.map((product, idx) => (
                    <div key={`${product.product_id ?? idx}-${product.product_name ?? idx}`} className="rounded-xl border border-border bg-muted/30 p-4">
                      <div className="flex items-start justify-between gap-3">
                        <div>
                          <p className="font-medium">{product.product_name ?? '—'}</p>
                          <p className="mt-1 text-xs text-muted-foreground">
                            {product.product_id || 'BANT qualification'}
                          </p>
                        </div>
                        <Badge variant="outline" className="capitalize">
                          BANT
                        </Badge>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}
        </TabCard>
      </TabsContent>

      <TabsContent value="tasks" className="mt-4">
        <TabCard
          title="Tasks"
          description="Follow-up tasks linked to this deal."
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
              icon={<CheckSquare className="h-10 w-10 text-muted-foreground/40" />}
              title="No tasks yet"
              description="Tasks related to this deal will appear here."
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
                        {task.description || 'No description'}
                      </div>
                    </div>
                    <Badge variant={task.status === 'completed' ? 'default' : 'outline'} className="capitalize">
                      {task.status.replace(/_/g, ' ')}
                    </Badge>
                  </div>
                  <div className="mt-3 flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
                    <span className="rounded-full bg-card px-2 py-1 capitalize ring-1 ring-border">{task.type.replace(/_/g, ' ')}</span>
                    <span className="rounded-full bg-card px-2 py-1 capitalize ring-1 ring-border">{task.priority}</span>
                    {task.due_date && <span>Due {formatSafeDate(task.due_date, true)}</span>}
                  </div>
                </button>
              ))}
            </div>
          )}
        </TabCard>
      </TabsContent>

      <TabsContent value="bant" className="mt-4">
        <TabCard
          title="BANT"
          description="Budget, authority, need, and timeline qualification for this opportunity."
        >
          {deal.lead_id ? (
            <LeadQualificationCard leadId={deal.lead_id} />
          ) : (
            <DealBantSnapshot deal={deal} />
          )}
        </TabCard>
      </TabsContent>

      <TabsContent value="information" className="mt-4">
        <TabCard
          title="Information"
          description="Deal, customer, BANT, and stage history details."
        >
          <div className="space-y-6">
            {deal.description && (
              <div>
                <h3 className="font-medium mb-2 text-sm">Description</h3>
                <p className="text-sm text-muted-foreground whitespace-pre-wrap">
                  {deal.description}
                </p>
              </div>
            )}

            {productItems.length > 0 && (
              <div>
                <h3 className="font-medium mb-2 text-sm">Products Sold</h3>
                <div className="overflow-x-auto rounded-xl border border-border bg-muted/20 p-4">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b text-left">
                        <th className="pb-2 pr-4 font-medium">Product</th>
                        <th className="pb-2 pr-4 font-medium text-right">Qty</th>
                        <th className="pb-2 pr-4 font-medium text-right">Unit Price</th>
                        <th className="pb-2 pr-4 font-medium text-right">Discount</th>
                        <th className="pb-2 font-medium text-right">Subtotal</th>
                      </tr>
                    </thead>
                    <tbody>
                      {productItems.map((item, idx) => (
                        <tr key={item.id ?? idx} className="border-b last:border-0">
                          <td className="py-3 pr-4">
                            <p className="font-medium">{item.product_name ?? '—'}</p>
                            {item.product_sku && (
                              <p className="text-xs text-muted-foreground">SKU: {item.product_sku}</p>
                            )}
                          </td>
                          <td className="py-3 pr-4 text-right">{item.quantity}</td>
                          <td className="py-3 pr-4 text-right">
                            {item.unit_price_formatted ?? formatCurrency(item.unit_price)}
                          </td>
                          <td className="py-3 pr-4 text-right">
                            {item.discount_amount_formatted ?? formatCurrency(item.discount_amount ?? 0)}
                          </td>
                          <td className="py-3 text-right font-medium">
                            {item.subtotal_formatted ?? formatCurrency(item.subtotal ?? 0)}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </div>
            )}

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <DetailItem label="Account" value={companyName} />
              <DetailItem label="Contact" value={contactName} />
              <DetailItem label="Contact Email" value={contactEmail} />
              <DetailItem label="Contact Phone" value={contactPhone} />
              <DetailItem label="Assigned To" value={deal.assigned_user?.name ?? '—'} />
              <DetailItem label="Value" value={deal.value_formatted ?? formatCurrency(deal.value ?? 0)} />
              <DetailItem label="Probability" value={`${deal.probability ?? 0}%`} />
              <DetailItem label="Source" value={deal.source ?? '—'} />
              <DetailItem label="Expected Close" value={formatSafeDate(deal.expected_close_date)} />
              <DetailItem label="Actual Close" value={formatSafeDate(deal.actual_close_date)} />
              <DetailItem label="Lead Source" value={lead?.lead_source ?? deal.source ?? '—'} />
              <DetailItem label="Location" value={address || '—'} />
            </div>

            {(budgetValue || authorityValue || needValue || timelineValue) && (
              <div>
                <h3 className="font-medium mb-2 text-sm">BANT Snapshot</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <DetailItem label="Budget" value={budgetValue ? formatCurrency(budgetValue) : '—'} />
                  <DetailItem label="Authority" value={authorityValue ?? '—'} />
                  <DetailItem label="Need" value={needValue ?? '—'} />
                  <DetailItem label="Timeline" value={formatSafeDate(timelineValue)} />
                </div>
              </div>
            )}

            {deal.notes && (
              <div>
                <h3 className="font-medium mb-2 text-sm">Notes</h3>
                <p className="text-sm text-muted-foreground whitespace-pre-wrap">
                  {deal.notes}
                </p>
              </div>
            )}

            <div>
              <h3 className="font-medium mb-3 text-sm">Stage History</h3>
              {historyItems.length === 0 ? (
                <div className="text-center py-8">
                  <History className="mx-auto h-10 w-10 text-muted-foreground/50 mb-3" />
                  <p className="text-sm text-muted-foreground">
                    No stage transitions recorded yet.
                  </p>
                </div>
              ) : (
                <DealHistoryTimeline dealId={dealId} />
              )}
            </div>

            <div className="text-xs text-muted-foreground pt-4 border-t space-y-1">
              <p>Created: {formatSafeDate(deal.created_at, true)}</p>
              <p>Updated: {formatSafeDate(deal.updated_at, true)}</p>
            </div>
          </div>
        </TabCard>
      </TabsContent>

      <VisitReportDetailModal
        visitReportId={viewingVisitReportId}
        open={!!viewingVisitReportId}
        onOpenChange={(open) => {
          if (!open) {
            setViewingVisitReportId(null);
          }
        }}
      />

      <TaskDetailModal
        taskId={viewingTaskId}
        open={!!viewingTaskId}
        onOpenChange={(open) => {
          if (!open) {
            setViewingTaskId(null);
          }
        }}
        onTaskUpdated={() => {
          void refetchTasks();
        }}
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
            initialDealId={dealId}
            initialAccountId={deal.account_id}
            initialContactId={deal.contact_id}
            onSubmit={handleCreateVisit}
            onCancel={() => setIsVisitModalOpen(false)}
            isLoading={createVisitReport.isPending}
          />
        </DialogContent>
      </Dialog>

      <CreateActivityDialog
        open={isActivityModalOpen}
        onOpenChange={setIsActivityModalOpen}
        accountId={deal.account_id}
        contactId={deal.contact_id}
        dealId={dealId}
        onSuccess={() => {
          void refetchActivities();
          void refetchVisitReports();
        }}
      />
    </Tabs>
  );
}

function DealBantSnapshot({ deal }: { readonly deal: Deal }) {
  const updateDeal = useUpdateDeal();
  const [isEditing, setIsEditing] = useState(false);
  const snapshot = deal.qualification_snapshot;
  const budgetAmount = snapshot?.budget_target_amount ?? (deal.value > 0 ? deal.value : undefined);
  const [formData, setFormData] = useState({
    budget_target_amount: budgetAmount ? String(Math.round(budgetAmount / 100)) : '',
    budget_notes: snapshot?.budget_notes ?? '',
    budget_confirmed: Boolean(deal.budget_confirmed || snapshot?.budget_confirmed),
    authority_target_person: snapshot?.authority_target_person ?? '',
    authority_target_role: snapshot?.authority_target_role ?? '',
    authority_notes: snapshot?.authority_notes ?? '',
    authority_confirmed: Boolean(deal.authority_confirmed || snapshot?.authority_confirmed),
    need_priority_level: snapshot?.need_priority_level ?? '',
    need_notes: snapshot?.need_notes ?? '',
    need_confirmed: Boolean(deal.need_confirmed || snapshot?.need_confirmed),
    timeline_target_date: snapshot?.timeline_target_date ? snapshot.timeline_target_date.split('T')[0] : '',
    timeline_notes: snapshot?.timeline_notes ?? '',
    timeline_confirmed: Boolean(deal.timeline_confirmed || snapshot?.timeline_confirmed),
  });

  const handleEdit = () => {
    setFormData({
      budget_target_amount: budgetAmount ? String(Math.round(budgetAmount / 100)) : '',
      budget_notes: snapshot?.budget_notes ?? '',
      budget_confirmed: Boolean(deal.budget_confirmed || snapshot?.budget_confirmed),
      authority_target_person: snapshot?.authority_target_person ?? '',
      authority_target_role: snapshot?.authority_target_role ?? '',
      authority_notes: snapshot?.authority_notes ?? '',
      authority_confirmed: Boolean(deal.authority_confirmed || snapshot?.authority_confirmed),
      need_priority_level: snapshot?.need_priority_level ?? '',
      need_notes: snapshot?.need_notes ?? '',
      need_confirmed: Boolean(deal.need_confirmed || snapshot?.need_confirmed),
      timeline_target_date: snapshot?.timeline_target_date ? snapshot.timeline_target_date.split('T')[0] : '',
      timeline_notes: snapshot?.timeline_notes ?? '',
      timeline_confirmed: Boolean(deal.timeline_confirmed || snapshot?.timeline_confirmed),
    });
    setIsEditing(true);
  };

  const handleSave = async () => {
    const parsedBudget = Number(formData.budget_target_amount);
    const qualificationSnapshot: Record<string, unknown> = {
      ...(snapshot ?? {}),
      budget_target_amount: Number.isFinite(parsedBudget) && parsedBudget > 0 ? Math.round(parsedBudget * 100) : undefined,
      budget_target_currency: snapshot?.budget_target_currency ?? 'IDR',
      budget_confirmed: formData.budget_confirmed,
      budget_notes: formData.budget_notes || undefined,
      authority_target_person: formData.authority_target_person || undefined,
      authority_target_role: formData.authority_target_role || undefined,
      authority_confirmed: formData.authority_confirmed,
      authority_notes: formData.authority_notes || undefined,
      need_priority_level: formData.need_priority_level || undefined,
      need_confirmed: formData.need_confirmed,
      need_notes: formData.need_notes || undefined,
      timeline_target_date: formData.timeline_target_date ? new Date(`${formData.timeline_target_date}T00:00:00`).toISOString() : undefined,
      timeline_confirmed: formData.timeline_confirmed,
      timeline_notes: formData.timeline_notes || undefined,
    };

    Object.keys(qualificationSnapshot).forEach((key) => {
      if (qualificationSnapshot[key] === undefined || qualificationSnapshot[key] === '') {
        delete qualificationSnapshot[key];
      }
    });

    await updateDeal.mutateAsync({
      id: deal.id,
      data: {
        budget_confirmed: formData.budget_confirmed,
        authority_confirmed: formData.authority_confirmed,
        need_confirmed: formData.need_confirmed,
        timeline_confirmed: formData.timeline_confirmed,
        qualification_snapshot: qualificationSnapshot,
      } as DealUpdateData,
    });
    setIsEditing(false);
  };

  const bantItems = [
    {
      key: 'budget',
      label: 'Budget',
      confirmed: deal.budget_confirmed || snapshot?.budget_confirmed,
      primary: budgetAmount ? formatCurrency(budgetAmount) : 'Not specified',
      secondary: snapshot?.budget_notes,
    },
    {
      key: 'authority',
      label: 'Authority',
      confirmed: deal.authority_confirmed || snapshot?.authority_confirmed,
      primary: snapshot?.authority_target_person || 'Decision maker not identified',
      secondary: [snapshot?.authority_target_role, snapshot?.authority_notes].filter(Boolean).join(' • '),
    },
    {
      key: 'need',
      label: 'Need',
      confirmed: deal.need_confirmed || snapshot?.need_confirmed,
      primary: snapshot?.need_priority_level ? `${snapshot.need_priority_level} priority` : 'Need priority not specified',
      secondary: snapshot?.need_notes,
    },
    {
      key: 'timeline',
      label: 'Timeline',
      confirmed: deal.timeline_confirmed || snapshot?.timeline_confirmed,
      primary: formatSafeDate(snapshot?.timeline_target_date),
      secondary: snapshot?.timeline_notes,
    },
  ];
  const completedCount = bantItems.filter((item) => item.confirmed).length;
  const products = snapshot?.need_target_products ?? [];

  return (
    <div className="space-y-4">
      <div className="rounded-xl border border-border bg-muted/30 p-4">
        <div className="flex flex-wrap items-center justify-between gap-3">
          <div>
            <h3 className="text-base font-medium">Qualification Checklist (BANT)</h3>
            <p className="mt-1 text-sm text-muted-foreground">
              {completedCount}/4 qualification items confirmed
            </p>
          </div>
          <Badge variant="outline" className="rounded-full px-3">
            {snapshot?.qualification_status?.toUpperCase() || (completedCount === 4 ? 'QUALIFIED' : 'PENDING')}
          </Badge>
        </div>
        {!isEditing && (
          <div className="mt-4">
            <Button size="sm" variant="outline" onClick={handleEdit} className="cursor-pointer">
              Edit
            </Button>
          </div>
        )}
        {snapshot?.qualification_score !== undefined && (
          <div className="mt-4">
            <div className="mb-1 flex justify-between text-sm">
              <span>Qualification Score</span>
              <span className="font-medium">{snapshot.qualification_score}/100</span>
            </div>
            <div className="h-2 overflow-hidden rounded-full bg-muted">
              <div
                className="h-full rounded-full bg-primary"
                style={{ width: `${Math.min(100, Math.max(0, snapshot.qualification_score))}%` }}
              />
            </div>
          </div>
        )}
      </div>

      {isEditing ? (
        <div className="space-y-4">
          <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
            <div className="rounded-xl border border-border bg-muted/30 p-4">
              <div className="mb-3 flex items-center justify-between gap-3">
                <h3 className="font-medium">Budget</h3>
                <Switch
                  checked={formData.budget_confirmed}
                  onCheckedChange={(checked) => setFormData((prev) => ({ ...prev, budget_confirmed: checked }))}
                />
              </div>
              <Input
                type="number"
                min={0}
                value={formData.budget_target_amount}
                onChange={(event) => setFormData((prev) => ({ ...prev, budget_target_amount: event.target.value }))}
                placeholder="Budget amount"
              />
              <Textarea
                className="mt-3"
                value={formData.budget_notes}
                onChange={(event) => setFormData((prev) => ({ ...prev, budget_notes: event.target.value }))}
                placeholder="Budget notes"
              />
            </div>

            <div className="rounded-xl border border-border bg-muted/30 p-4">
              <div className="mb-3 flex items-center justify-between gap-3">
                <h3 className="font-medium">Authority</h3>
                <Switch
                  checked={formData.authority_confirmed}
                  onCheckedChange={(checked) => setFormData((prev) => ({ ...prev, authority_confirmed: checked }))}
                />
              </div>
              <Input
                value={formData.authority_target_person}
                onChange={(event) => setFormData((prev) => ({ ...prev, authority_target_person: event.target.value }))}
                placeholder="Decision maker"
              />
              <Input
                className="mt-3"
                value={formData.authority_target_role}
                onChange={(event) => setFormData((prev) => ({ ...prev, authority_target_role: event.target.value }))}
                placeholder="Role"
              />
              <Textarea
                className="mt-3"
                value={formData.authority_notes}
                onChange={(event) => setFormData((prev) => ({ ...prev, authority_notes: event.target.value }))}
                placeholder="Authority notes"
              />
            </div>

            <div className="rounded-xl border border-border bg-muted/30 p-4">
              <div className="mb-3 flex items-center justify-between gap-3">
                <h3 className="font-medium">Need</h3>
                <Switch
                  checked={formData.need_confirmed}
                  onCheckedChange={(checked) => setFormData((prev) => ({ ...prev, need_confirmed: checked }))}
                />
              </div>
              <Input
                value={formData.need_priority_level}
                onChange={(event) => setFormData((prev) => ({ ...prev, need_priority_level: event.target.value }))}
                placeholder="Priority level"
              />
              <Textarea
                className="mt-3"
                value={formData.need_notes}
                onChange={(event) => setFormData((prev) => ({ ...prev, need_notes: event.target.value }))}
                placeholder="Need notes"
              />
            </div>

            <div className="rounded-xl border border-border bg-muted/30 p-4">
              <div className="mb-3 flex items-center justify-between gap-3">
                <h3 className="font-medium">Timeline</h3>
                <Switch
                  checked={formData.timeline_confirmed}
                  onCheckedChange={(checked) => setFormData((prev) => ({ ...prev, timeline_confirmed: checked }))}
                />
              </div>
              <Input
                type="date"
                value={formData.timeline_target_date}
                onChange={(event) => setFormData((prev) => ({ ...prev, timeline_target_date: event.target.value }))}
              />
              <Textarea
                className="mt-3"
                value={formData.timeline_notes}
                onChange={(event) => setFormData((prev) => ({ ...prev, timeline_notes: event.target.value }))}
                placeholder="Timeline notes"
              />
            </div>
          </div>
          <div className="flex justify-end gap-2">
            <Button variant="outline" onClick={() => setIsEditing(false)} disabled={updateDeal.isPending}>
              Cancel
            </Button>
            <Button onClick={handleSave} disabled={updateDeal.isPending}>
              {updateDeal.isPending ? 'Saving...' : 'Save Changes'}
            </Button>
          </div>
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
          {bantItems.map((item) => (
            <div key={item.key} className="rounded-xl border border-border bg-muted/30 p-4">
              <div className="mb-2 flex items-center justify-between gap-3">
                <h3 className="font-medium">{item.label}</h3>
                <Badge variant={item.confirmed ? 'default' : 'outline'}>
                  {item.confirmed ? 'Confirmed' : 'Open'}
                </Badge>
              </div>
              <p className="text-sm text-foreground">{item.primary}</p>
              {item.secondary && (
                <p className="mt-2 whitespace-pre-wrap text-sm text-muted-foreground">{item.secondary}</p>
              )}
            </div>
          ))}
        </div>
      )}

      {products.length > 0 && (
        <div className="space-y-3">
          <h3 className="text-sm font-medium">Need Products</h3>
          {products.map((product, index) => (
            <div key={`${product.product_id ?? index}-${product.product_name ?? index}`} className="rounded-xl border border-border bg-muted/30 p-4">
              <p className="font-medium">{product.product_name ?? '—'}</p>
              {product.product_id && (
                <p className="mt-1 text-xs text-muted-foreground">{product.product_id}</p>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
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
        <div className="flex shrink-0 flex-wrap items-center gap-2">{actions}</div>
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

function DetailItem({
  label,
  value,
}: {
  readonly label: string;
  readonly value: string;
}) {
  return (
    <div className="rounded-xl border border-border bg-muted/30 p-4">
      <h3 className="mb-1 text-xs font-medium uppercase tracking-wide text-muted-foreground">{label}</h3>
      <p className="break-words text-sm text-foreground">{value}</p>
    </div>
  );
}

type DealTimelineItem = {
  id: string;
  kind: 'visit' | 'activity';
  type: string;
  dateValue: string;
  description: string;
  visitId?: string;
  accountName?: string;
  ownerName?: string;
};

function buildTimelineItems(visits: VisitReport[], activities: Activity[]): DealTimelineItem[] {
  const visitItems = visits.map((visit) => ({
    id: `visit-${visit.id}`,
    kind: 'visit' as const,
    type: 'Visit',
    dateValue: visit.visit_date ?? visit.created_at ?? '',
    description: visit.purpose || visit.notes || 'Visit Report',
    visitId: visit.id,
    accountName: visit.account?.name,
    ownerName: visit.sales_rep?.name,
  }));

  const activityItems = activities.map((activity, index) => ({
    id: `activity-${activity.id ?? activity.created_at ?? index}`,
    kind: 'activity' as const,
    type: activity.activity_type?.name ?? activity.type ?? 'Activity',
    dateValue: activity.timestamp ?? activity.created_at ?? '',
    description: activity.description || 'Activity',
    accountName: activity.account?.name,
    ownerName: activity.user?.name,
  }));

  return [...visitItems, ...activityItems].sort((a, b) => {
    const left = new Date(b.dateValue).getTime();
    const right = new Date(a.dateValue).getTime();
    return (Number.isNaN(left) ? 0 : left) - (Number.isNaN(right) ? 0 : right);
  });
}

function buildProductInterestActivities(visits: VisitReport[], activities: Activity[]): Activity[] {
  const visitIds = new Set(visits.map((visit) => visit.id));
  const visitActivities: Activity[] = visits.map((visit) => ({
    id: `deal-visit-report-${visit.id}`,
    type: 'visit',
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
  if (!activity.metadata || typeof activity.metadata !== 'object') {
    return null;
  }

  const value = activity.metadata.visit_report_id;
  return typeof value === 'string' ? value : null;
}

function dedupeVisitReports(visits: VisitReport[]) {
  const map = new Map<string, VisitReport>();
  for (const visit of visits) {
    if (!visit?.id) continue;
    map.set(visit.id, visit);
  }
  return Array.from(map.values());
}

function dedupeActivities(activities: Activity[]) {
  const map = new Map<string, Activity>();
  for (const activity of activities) {
    if (!activity?.id) continue;
    map.set(activity.id, activity);
  }
  return Array.from(map.values());
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
