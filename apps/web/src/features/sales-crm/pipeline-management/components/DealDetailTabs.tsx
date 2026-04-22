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
import { formatCurrency } from '../utils/format';
import {
  Info,
  Package,
  CheckSquare,
  MapPin,
  Activity as ActivityIcon,
  History,
} from 'lucide-react';

interface DealDetailTabsProps {
  readonly dealId: string;
}

export function DealDetailTabs({ dealId }: DealDetailTabsProps) {
  const [activeTab, setActiveTab] = useState('details');
  const { data: deal, isLoading } = useDeal(dealId);
  const { data: visitReports } = useDealVisitReports(dealId);
  const { data: activities } = useDealActivities(dealId);
  const { data: history } = useDealHistory(dealId);

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

  const productItems = deal.product_items ?? [];
  const visits = Array.isArray(visitReports) ? visitReports : [];
  const acts = Array.isArray(activities) ? activities : [];
  const historyItems = Array.isArray(history) ? history : [];

  return (
    <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
      <TabsList className="grid w-full grid-cols-6">
        <TabsTrigger value="details" className="gap-2 cursor-pointer">
          <Info size={16} />
          <span className="hidden sm:inline">Details</span>
        </TabsTrigger>
        <TabsTrigger value="products" className="gap-2 cursor-pointer">
          <Package size={16} />
          <span className="hidden sm:inline">Products</span>
        </TabsTrigger>
        <TabsTrigger value="tasks" className="gap-2 cursor-pointer">
          <CheckSquare size={16} />
          <span className="hidden sm:inline">Tasks</span>
        </TabsTrigger>
        <TabsTrigger value="visits" className="gap-2 cursor-pointer">
          <MapPin size={16} />
          <span className="hidden sm:inline">Visits</span>
        </TabsTrigger>
        <TabsTrigger value="activities" className="gap-2 cursor-pointer">
          <ActivityIcon size={16} />
          <span className="hidden sm:inline">Activities</span>
        </TabsTrigger>
        <TabsTrigger value="history" className="gap-2 cursor-pointer">
          <History size={16} />
          <span className="hidden sm:inline">History</span>
        </TabsTrigger>
      </TabsList>

      {/* Details Tab */}
      <TabsContent value="details" className="mt-4">
        <Card>
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
              <DetailItem label="Account" value={deal.account?.name ?? 'N/A'} />
              <DetailItem label="Contact" value={deal.contact?.name ?? '—'} />
              <DetailItem
                label="Assigned To"
                value={deal.assigned_user?.name ?? '—'}
              />
              <DetailItem
                label="Value"
                value={deal.value_formatted ?? formatCurrency(deal.value ?? 0)}
              />
              <DetailItem
                label="Probability"
                value={`${deal.probability ?? 0}%`}
              />
              <DetailItem label="Source" value={deal.source ?? '—'} />
              <DetailItem
                label="Expected Close"
                value={formatSafeDate(deal.expected_close_date)}
              />
              <DetailItem
                label="Actual Close"
                value={formatSafeDate(deal.actual_close_date)}
              />
            </div>

            {deal.notes && (
              <div>
                <h3 className="font-medium mb-2 text-sm">Notes</h3>
                <p className="text-sm text-muted-foreground whitespace-pre-wrap">
                  {deal.notes}
                </p>
              </div>
            )}

            <div className="text-xs text-muted-foreground pt-4 border-t space-y-1">
              <p>Created: {formatSafeDate(deal.created_at, true)}</p>
              <p>Updated: {formatSafeDate(deal.updated_at, true)}</p>
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
              <div className="text-center py-8">
                <Package className="mx-auto h-10 w-10 text-muted-foreground/50 mb-3" />
                <p className="text-sm text-muted-foreground">
                  No products added to this deal yet.
                </p>
              </div>
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

      {/* Tasks Tab */}
      <TabsContent value="tasks" className="mt-4">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Related Tasks</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-center py-8">
              <CheckSquare className="mx-auto h-10 w-10 text-muted-foreground/50 mb-3" />
              <p className="text-sm text-muted-foreground">
                Tasks related to this deal will appear here.
              </p>
            </div>
          </CardContent>
        </Card>
      </TabsContent>

      {/* Visits Tab */}
      <TabsContent value="visits" className="mt-4">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Visit Reports</CardTitle>
          </CardHeader>
          <CardContent>
            {visits.length === 0 ? (
              <div className="text-center py-8">
                <MapPin className="mx-auto h-10 w-10 text-muted-foreground/50 mb-3" />
                <p className="text-sm text-muted-foreground">
                  No visit reports linked to this deal yet.
                </p>
              </div>
            ) : (
              <div className="space-y-3">
                {visits.map((visit: VisitReport) => (
                  <div
                    key={visit.id}
                    className="p-4 border rounded-lg hover:bg-accent/50 transition-colors"
                  >
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="font-medium text-sm">
                          {visit.purpose ?? 'Visit Report'}
                        </p>
                        <p className="text-xs text-muted-foreground">
                          {formatSafeDate(visit.visit_date ?? visit.created_at)}
                        </p>
                      </div>
                      {visit.status && (
                        <Badge variant="outline" className="text-xs">
                          {visit.status}
                        </Badge>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </TabsContent>

      {/* Activities Tab */}
      <TabsContent value="activities" className="mt-4">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Activities</CardTitle>
          </CardHeader>
          <CardContent>
            {acts.length === 0 ? (
              <div className="text-center py-8">
                <ActivityIcon className="mx-auto h-10 w-10 text-muted-foreground/50 mb-3" />
                <p className="text-sm text-muted-foreground">
                  No activities logged for this deal yet.
                </p>
              </div>
            ) : (
              <div className="space-y-3">
                {acts.map((act: Activity) => (
                  <div
                    key={act.id ?? ''}
                    className="p-4 border rounded-lg hover:bg-accent/50 transition-colors"
                  >
                    <div className="flex items-center justify-between">
                      <div>
                        <p className="font-medium text-sm">
                          {act.type ?? 'Activity'}
                        </p>
                        <p className="text-xs text-muted-foreground">
                          {act.description ?? ''}
                        </p>
                      </div>
                      <p className="text-xs text-muted-foreground">
                        {formatSafeDate(act.created_at ?? '')}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </TabsContent>

      {/* History Tab */}
      <TabsContent value="history" className="mt-4">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Stage History</CardTitle>
          </CardHeader>
          <CardContent>
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
