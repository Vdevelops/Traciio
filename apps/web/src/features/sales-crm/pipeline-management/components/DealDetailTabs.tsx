'use client';

import { useState } from 'react';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { DealHistoryTimeline } from './deal-history-timeline';
import type { VisitReport } from '@/features/sales-crm/visit-report/types';
import type { Activity } from '@/features/sales-crm/visit-report/types/activity';
import {
  useDeal,
  useDealVisitReports,
  useDealActivities,
  useDealHistory,
} from '../hooks/useDeals';
import { useLead } from '@/features/sales-crm/lead-management/hooks/useLeads';
import { useLeadQualification } from '@/features/sales-crm/lead-management/hooks/useLeadQualification';
import { useTasks } from '@/features/sales-crm/task-management/hooks/useTasks';
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
  const { data: deal, isLoading } = useDeal(dealId);
  const { data: leadData } = useLead(deal?.lead_id ?? '');
  const { qualification } = useLeadQualification(deal?.lead_id ?? '');
  const { data: visitReports } = useDealVisitReports(dealId);
  const { data: activities } = useDealActivities(dealId);
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

  return (
    <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
      <TabsList className="h-auto w-full justify-start gap-1 overflow-x-auto rounded-none border-b border-border bg-transparent p-0">
        <TabsTrigger value="activities" className="cursor-pointer gap-1.5 rounded-none border-b-2 border-transparent px-1 py-3 text-sm text-muted-foreground data-[state=active]:border-primary data-[state=active]:bg-transparent data-[state=active]:text-foreground data-[state=active]:shadow-none">
          <ActivityIcon size={16} />
          <span>Activities</span>
        </TabsTrigger>
        <TabsTrigger value="tasks" className="cursor-pointer gap-1.5 rounded-none border-b-2 border-transparent px-1 py-3 text-sm text-muted-foreground data-[state=active]:border-primary data-[state=active]:bg-transparent data-[state=active]:text-foreground data-[state=active]:shadow-none">
          <CheckSquare size={16} />
          <span>Tasks</span>
        </TabsTrigger>
        <TabsTrigger value="products" className="cursor-pointer gap-1.5 rounded-none border-b-2 border-transparent px-1 py-3 text-sm text-muted-foreground data-[state=active]:border-primary data-[state=active]:bg-transparent data-[state=active]:text-foreground data-[state=active]:shadow-none">
          <Package size={16} />
          <span>Product Interest</span>
        </TabsTrigger>
        <TabsTrigger value="information" className="cursor-pointer gap-1.5 rounded-none border-b-2 border-transparent px-1 py-3 text-sm text-muted-foreground data-[state=active]:border-primary data-[state=active]:bg-transparent data-[state=active]:text-foreground data-[state=active]:shadow-none">
          <Info size={16} />
          <span>Information</span>
        </TabsTrigger>
      </TabsList>

      <TabsContent value="activities" className="mt-4">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Activities & Visits</CardTitle>
          </CardHeader>
          <CardContent className="space-y-6">
            {visits.length > 0 && (
              <div className="space-y-3">
                <h3 className="text-sm font-medium text-muted-foreground">Visit Reports</h3>
                {visits.map((visit: VisitReport) => (
                  <div key={visit.id} className="rounded-lg border p-4">
                    <div className="flex items-center justify-between gap-3">
                      <div>
                        <p className="font-medium text-sm">{visit.purpose ?? 'Visit Report'}</p>
                        <p className="text-xs text-muted-foreground">
                          {formatSafeDate(visit.visit_date ?? visit.created_at, true)}
                        </p>
                      </div>
                      {visit.status && (
                        <Badge variant="outline" className="text-xs capitalize">
                          {visit.status}
                        </Badge>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            )}

            <div className="space-y-3">
              <h3 className="text-sm font-medium text-muted-foreground">Activities</h3>
              {acts.length === 0 ? (
                <div className="text-center py-8">
                  <ActivityIcon className="mx-auto h-10 w-10 text-muted-foreground/50 mb-3" />
                  <p className="text-sm text-muted-foreground">
                    No activities logged for this deal yet.
                  </p>
                </div>
              ) : (
                acts.map((act: Activity) => (
                  <div key={act.id ?? ''} className="rounded-lg border p-4">
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="font-medium text-sm">{act.type ?? 'Activity'}</p>
                        <p className="text-xs text-muted-foreground">{act.description ?? ''}</p>
                      </div>
                      <p className="text-xs text-muted-foreground">
                        {formatSafeDate(act.created_at ?? '', true)}
                      </p>
                    </div>
                  </div>
                ))
              )}
            </div>
          </CardContent>
        </Card>
      </TabsContent>

      {/* Products Tab */}
      <TabsContent value="products" className="mt-4">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Product Items</CardTitle>
          </CardHeader>
          <CardContent>
            {productItems.length === 0 ? (
              productInterests.length === 0 ? (
                <div className="text-center py-8">
                  <Package className="mx-auto h-10 w-10 text-muted-foreground/50 mb-3" />
                  <p className="text-sm text-muted-foreground">
                    No products added to this deal yet.
                  </p>
                </div>
              ) : (
                <div className="space-y-3">
                  {productInterests.map((product, idx) => (
                    <div key={`${product.product_id ?? idx}-${product.product_name ?? idx}`} className="rounded-lg border p-4">
                      <p className="font-medium">{product.product_name ?? '—'}</p>
                      {product.product_id && (
                        <p className="text-xs text-muted-foreground mt-1">{product.product_id}</p>
                      )}
                    </div>
                  ))}
                </div>
              )
            ) : (
              <div className="overflow-x-auto">
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
          </CardContent>
        </Card>
      </TabsContent>

      <TabsContent value="tasks" className="mt-4">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Related Tasks</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-center py-8">
              {tasks.length === 0 ? (
                <>
                  <CheckSquare className="mx-auto h-10 w-10 text-muted-foreground/50 mb-3" />
                  <p className="text-sm text-muted-foreground">
                    Tasks related to this deal will appear here.
                  </p>
                </>
              ) : (
                <div className="space-y-3 text-left">
                  {tasks.map((task) => (
                    <div key={task.id} className="rounded-lg border p-4">
                      <div className="flex items-start justify-between gap-3">
                        <div className="min-w-0">
                          <div className="font-medium">{task.title}</div>
                          <div className="mt-1 text-sm text-muted-foreground">
                            {task.description || 'No description'}
                          </div>
                        </div>
                        <Badge variant={task.status === 'completed' ? 'default' : 'outline'} className="capitalize">
                          {task.status.replace(/_/g, ' ')}
                        </Badge>
                      </div>
                      <div className="mt-2 text-xs text-muted-foreground">
                        {task.due_date ? `Due ${formatSafeDate(task.due_date, true)}` : 'No due date'}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      </TabsContent>

      <TabsContent value="information" className="mt-4">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Information</CardTitle>
          </CardHeader>
          <CardContent className="pt-6 space-y-6">
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
          </CardContent>
        </Card>
      </TabsContent>
    </Tabs>
  );
}

// Helper components
function DetailItem({
  label,
  value,
}: {
  readonly label: string;
  readonly value: string;
}) {
  return (
    <div>
      <h3 className="font-medium mb-1 text-sm text-muted-foreground">{label}</h3>
      <p className="text-sm">{value}</p>
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
