'use client';

import { type ReactNode, useMemo, useState } from 'react';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { DealHistoryTimeline } from './deal-history-timeline';
import { ProductInterestTab } from '@/features/sales-crm/visit-report/components/product-interest-tab';
import type { VisitReport } from '@/features/sales-crm/visit-report/types';
import type { Activity } from '@/features/sales-crm/visit-report/types/activity';
import { VisitReportDetailModal } from '@/features/sales-crm/visit-report/components/visit-report-detail-modal';
import {
  useDeal,
  useDealVisitReports,
  useDealActivities,
  useDealHistory,
} from '../hooks/useDeals';
import { useLead, useLeadActivities, useLeadVisitReports } from '@/features/sales-crm/lead-management/hooks/useLeads';
import { useLeadQualification } from '@/features/sales-crm/lead-management/hooks/useLeadQualification';
import { useTasks } from '@/features/sales-crm/task-management/hooks/useTasks';
import { TaskDetailModal } from '@/features/sales-crm/task-management/components/task-detail-modal';
import { formatCurrency } from '../utils/format';
import {
  Info,
  Package,
  CheckSquare,
  MapPin,
  Activity as ActivityIcon,
  History,
} from 'lucide-react';
import type { Task } from '@/features/sales-crm/task-management/types';

interface DealDetailTabsProps {
  readonly dealId: string;
}

export function DealDetailTabs({ dealId }: DealDetailTabsProps) {
  const [activeTab, setActiveTab] = useState('activities');
  const [viewingVisitReportId, setViewingVisitReportId] = useState<string | null>(null);
  const [viewingTaskId, setViewingTaskId] = useState<string | null>(null);
  const { data: deal, isLoading } = useDeal(dealId);
  const { data: leadData } = useLead(deal?.lead_id ?? '');
  const { qualification } = useLeadQualification(deal?.lead_id ?? '');
  const { data: visitReports } = useDealVisitReports(dealId);
  const { data: activities } = useDealActivities(dealId);
  const { data: leadVisitReportsData, isLoading: isLeadVisitReportsLoading } = useLeadVisitReports(deal?.lead_id ?? '', { per_page: 100 });
  const { data: leadActivitiesData, isLoading: isLeadActivitiesLoading } = useLeadActivities(deal?.lead_id ?? '', { per_page: 100 });
  const { data: history } = useDealHistory(dealId);
  const { data: dealTasksResponse } = useTasks({ deal_id: dealId, per_page: 20 });
  const { data: leadTasksResponse } = useTasks({ lead_id: deal?.lead_id, per_page: 20 });

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
  const visits = Array.isArray(visitReports) ? visitReports : [];
  const acts = Array.isArray(activities) ? activities : [];
  const leadVisits = Array.isArray(leadVisitReportsData?.data) ? leadVisitReportsData.data : [];
  const leadActs = Array.isArray(leadActivitiesData?.data) ? leadActivitiesData.data : [];
  const historyItems = Array.isArray(history) ? history : [];
  const dealTasks = Array.isArray(dealTasksResponse?.data) ? dealTasksResponse.data : [];
  const leadTasks = Array.isArray(leadTasksResponse?.data) ? leadTasksResponse.data : [];
  const tasks = dedupeTasks([...dealTasks, ...leadTasks]);
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
  const timelineItems = useMemo(() => buildTimelineItems(visits, acts), [visits, acts]);
  const productInterestActivities = useMemo(
    () => buildProductInterestActivities(
      dedupeVisitReports([...visits, ...leadVisits]),
      dedupeActivities([...acts, ...leadActs]),
    ),
    [visits, leadVisits, acts, leadActs],
  );
  const productInterestCount = getProductInterestCount(productInterestActivities);
  const isProductInterestLoading = isLeadVisitReportsLoading || isLeadActivitiesLoading;

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
        <TabsTrigger value="information" className="cursor-pointer gap-2 whitespace-nowrap rounded-none border-b-2 border-transparent px-4 py-3 text-sm text-muted-foreground data-[state=active]:border-primary data-[state=active]:bg-transparent data-[state=active]:text-foreground data-[state=active]:shadow-none">
          <Info size={16} />
          <span>Information</span>
        </TabsTrigger>
      </TabsList>

      <TabsContent value="activities" className="mt-4">
        <TabCard
          title="Activities"
          description="Track calls, visits, meetings, and follow ups for this deal."
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
          description="Products attached to this deal and products carried from lead qualification."
        >
          {productItems.length === 0 && productInterests.length === 0 && productInterestCount === 0 ? (
            isProductInterestLoading ? (
              <ProductInterestTab
                activities={productInterestActivities}
                isLoading={isProductInterestLoading}
              />
            ) : (
              <EmptyState
                icon={<Package className="h-10 w-10 text-muted-foreground/40" />}
                title="No product interest recorded"
                description="Products will appear here after they are attached to the deal or captured from qualification."
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
                      Products from the source lead qualification are shown separately from deal line items.
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

              {productItems.length > 0 && (
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
                              <p className="text-xs text-muted-foreground">
                                SKU: {item.product_sku}
                              </p>
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
                    <tfoot>
                      <tr className="border-t-2">
                        <td colSpan={4} className="py-3 text-right font-medium">
                          Total
                        </td>
                        <td className="py-3 text-right font-bold">
                          {deal.value_formatted ?? formatCurrency(deal.value ?? 0)}
                        </td>
                      </tr>
                    </tfoot>
                  </table>
                </div>
              )}
            </div>
          )}
        </TabCard>
      </TabsContent>

      <TabsContent value="tasks" className="mt-4">
        <TabCard
          title="Tasks"
          description="Follow-up tasks linked to this deal and its source lead."
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
      />
    </Tabs>
  );
}

function TabCard({
  title,
  description,
  children,
}: {
  readonly title: string;
  readonly description: string;
  readonly children: ReactNode;
}) {
  return (
    <Card className="border-slate-200 bg-white shadow-sm">
      <CardHeader className="flex flex-row items-center justify-between gap-3 pb-4">
        <div>
          <CardTitle className="text-base">{title}</CardTitle>
          <p className="text-sm text-slate-500">{description}</p>
        </div>
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
