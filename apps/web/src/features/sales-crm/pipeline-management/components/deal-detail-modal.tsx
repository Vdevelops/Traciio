"use client";

import { User, Building2, FileText, Edit, Trash2, DollarSign, TrendingUp, Calendar, Circle, Activity, MapPin, Package, Contact, CheckSquare } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { Avatar, AvatarImage } from "@/components/ui/avatar";
import { useDeal, useUpdateDeal, useDeleteDeal, useDealVisitReports, useDealActivities } from "../hooks/useDeals";
import { DealForm } from "./deal-form";
import { formatCurrency } from "@/lib/utils";
import { toast } from "sonner";
import { useState } from "react";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { renderIcon } from "@/features/sales-crm/visit-report/lib/icon-utils";
import { useTranslations } from "next-intl";
import { AccountDetailModal } from "@/features/sales-crm/account-management/components/account-detail-modal";
import { ContactDetailModal } from "@/features/sales-crm/account-management/components/contact-detail-modal";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { ProductInterestTab } from "@/features/sales-crm/visit-report/components/product-interest-tab";
import { useTasks } from "@/features/sales-crm/task-management/hooks/useTasks";
import type { Task } from "@/features/sales-crm/task-management/types";
import type { Activity as CRMActivity } from "@/features/sales-crm/visit-report/types/activity";
import type { VisitReport } from "@/features/sales-crm/visit-report/types";
import { VisitReportDetailModal } from "@/features/sales-crm/visit-report/components/visit-report-detail-modal";
import { TaskDetailModal } from "@/features/sales-crm/task-management/components/task-detail-modal";

const statusVariantMap: Record<string, "default" | "secondary" | "outline" | "destructive"> = {
  open: "secondary",
  won: "default",
  lost: "destructive",
};

interface DealDetailModalProps {
  readonly dealId: string | null;
  readonly open: boolean;
  readonly onOpenChange: (open: boolean) => void;
  readonly onDealUpdated?: () => void;
}

export function DealDetailModal({
  dealId,
  open,
  onOpenChange,
  onDealUpdated,
}: DealDetailModalProps) {
  const t = useTranslations("deals.detail");
  const tCommon = useTranslations("common");
  
  const { data, isLoading, error } = useDeal(dealId || "");
  const updateDeal = useUpdateDeal();
  const deleteDeal = useDeleteDeal();
  
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
  const [viewingAccountId, setViewingAccountId] = useState<string | null>(null);
  const [viewingContactId, setViewingContactId] = useState<string | null>(null);
  const [viewingVisitReportId, setViewingVisitReportId] = useState<string | null>(null);
  const [viewingTaskId, setViewingTaskId] = useState<string | null>(null);

  const deal = data;

  const handleUpdate = async (formData: Parameters<typeof updateDeal.mutateAsync>[0]["data"]) => {
    if (!dealId) return;
    try {
      await updateDeal.mutateAsync({ id: dealId, data: formData });
      toast.success(t("toast.updated"));
      setIsEditDialogOpen(false);
      onDealUpdated?.();
    } catch {
      // Error handled by interceptor
    }
  };

  const handleDelete = async () => {
    if (!dealId) return;
    try {
      await deleteDeal.mutateAsync(dealId);
      toast.success(t("toast.deleted"));
      setIsDeleteDialogOpen(false);
      onOpenChange(false);
      onDealUpdated?.();
    } catch {
      // Error handled by interceptor
    }
  };

  return (
    <>
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent className="sm:max-w-7xl max-h-[90vh]">
          <DialogHeader>
            <DialogTitle>{t("drawerTitle")}</DialogTitle>
          </DialogHeader>

          {isLoading && (
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              {/* Left Column Skeleton */}
              <div className="space-y-6">
                <div className="flex items-center gap-4">
                  <Skeleton className="h-20 w-20 rounded-full" />
                  <div className="space-y-2">
                    <Skeleton className="h-6 w-48" />
                    <Skeleton className="h-4 w-32" />
                  </div>
                </div>
                <Card className="surface-panel border-border/70">
                  <CardHeader>
                    <Skeleton className="h-6 w-32" />
                    <Skeleton className="h-4 w-64 mt-2" />
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
              {/* Right Column Skeleton */}
              <div className="space-y-3">
                <Skeleton className="h-8 w-48" />
                {[...Array(5)].map((_, i) => (
                  <Skeleton key={i} className="h-20 w-full" />
                ))}
              </div>
            </div>
          )}

          {error && (
            <div className="text-center text-muted-foreground py-8">
              {t("notFound")}
            </div>
          )}

          {!isLoading && !error && deal && (
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 overflow-hidden">
              {/* Left Column: Deal Information */}
              <div className="space-y-6 overflow-y-auto max-h-[calc(90vh-120px)] pr-2">
                {/* Deal Header */}
                <div className="crm-hero flex items-center gap-4 rounded-3xl border border-border/70 px-5 py-5">
                  <div className="flex-1">
                    <h2 className="text-2xl font-medium tracking-tight">{deal.title}</h2>
                    <div className="flex items-center gap-2 mt-2">
                      <Badge variant={statusVariantMap[deal.status] || "outline"} className="capitalize">
                        {deal.status}
                      </Badge>
                      {deal.stage && (
                        <>
                          <Circle
                            className="h-3 w-3"
                            style={{ color: deal.stage.color || "#CBD5F5", fill: deal.stage.color || "#CBD5F5" }}
                          />
                          <span className="text-sm text-muted-foreground">{deal.stage.name || t("fallbacks.unknownStage")}</span>
                        </>
                      )}
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setIsEditDialogOpen(true)}
                      className="shrink-0"
                    >
                      <Edit className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="destructive"
                      size="sm"
                      onClick={() => setIsDeleteDialogOpen(true)}
                      disabled={deleteDeal.isPending}
                      className="shrink-0"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </div>

                {/* Deal Info Card */}
                <Card>
                  <CardHeader>
                    <CardTitle>{t("sections.basicInfo")}</CardTitle>
                    <CardDescription>
                      {t("sections.description")}
                    </CardDescription>
                  </CardHeader>
                  <CardContent className="space-y-6">
                    {deal.description && (
                      <div className="space-y-2">
                        <div className="text-sm text-muted-foreground">
                          {t("sections.description")}
                        </div>
                        <div className="text-base whitespace-pre-wrap">{deal.description}</div>
                      </div>
                    )}

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                      <div className="space-y-2">
                        <div className="flex items-center gap-2 text-sm text-muted-foreground">
                          <DollarSign className="h-4 w-4" />
                          <span>{t("sections.value")}</span>
                        </div>
                        <div className="text-base font-medium">
                          {deal.value_formatted || formatCurrency(deal.value ?? 0)}
                        </div>
                      </div>

                      {(deal.probability ?? 0) > 0 && (
                        <div className="space-y-2">
                          <div className="flex items-center gap-2 text-sm text-muted-foreground">
                            <TrendingUp className="h-4 w-4" />
                            <span>Probability</span>
                          </div>
                          <div className="text-base font-medium">{deal.probability ?? 0}%</div>
                        </div>
                      )}

                      {deal.account && (
                        <div className="space-y-2">
                          <div className="flex items-center gap-2 text-sm text-muted-foreground">
                            <Building2 className="h-4 w-4" />
                            <span>{t("sections.account")}</span>
                          </div>
                          <div>
                            <Button
                              variant="link"
                              className="h-auto p-0 text-base font-medium"
                              onClick={() => setViewingAccountId(deal.account?.id || null)}
                            >
                              {deal.account?.name || t("fallbacks.unknownAccount")}
                            </Button>
                          </div>
                        </div>
                      )}

                      {deal.contact && (
                        <div className="space-y-2">
                          <div className="flex items-center gap-2 text-sm text-muted-foreground">
                            <User className="h-4 w-4" />
                            <span>{t("sections.contact")}</span>
                          </div>
                          <div>
                            <Button
                              variant="link"
                              className="h-auto p-0 text-base font-medium"
                              onClick={() => setViewingContactId(deal.contact?.id || null)}
                            >
                              {deal.contact?.name || t("fallbacks.unknownContact")}
                            </Button>
                          </div>
                        </div>
                      )}

                      {deal.assigned_user && (
                        <div className="space-y-2">
                          <div className="flex items-center gap-2 text-sm text-muted-foreground">
                            <User className="h-4 w-4" />
                            <span>{t("sections.assignedTo")}</span>
                          </div>
                          <div className="flex items-center gap-2">
                            {deal.assigned_user?.avatar_url && (
                              <Avatar className="h-6 w-6">
                                <AvatarImage src={deal.assigned_user.avatar_url} alt={deal.assigned_user?.name ?? t("fallbacks.unknownUser")} />
                              </Avatar>
                            )}
                            <div className="text-base font-medium">{deal.assigned_user?.name ?? t("fallbacks.unknownUser")}</div>
                          </div>
                        </div>
                      )}

                      {deal.expected_close_date && (
                        <div className="space-y-2">
                          <div className="flex items-center gap-2 text-sm text-muted-foreground">
                            <Calendar className="h-4 w-4" />
                            <span>{t("sections.expectedCloseDate")}</span>
                          </div>
                          <div className="text-base font-medium">
                            {(() => {
                              if (!deal.expected_close_date) return t("fallbacks.noDate");
                              const date = new Date(deal.expected_close_date);
                              if (isNaN(date.getTime())) return t("fallbacks.invalidDate");
                              return date.toLocaleDateString("id-ID", {
                                year: "numeric",
                                month: "long",
                                day: "numeric",
                              });
                            })()}
                          </div>
                        </div>
                      )}

                      {deal.actual_close_date && (
                        <div className="space-y-2">
                          <div className="flex items-center gap-2 text-sm text-muted-foreground">
                            <Calendar className="h-4 w-4" />
                            <span>{t("sections.actualCloseDate")}</span>
                          </div>
                          <div className="text-base font-medium">
                            {(() => {
                              if (!deal.actual_close_date) return t("fallbacks.noDate");
                              const date = new Date(deal.actual_close_date);
                              if (isNaN(date.getTime())) return t("fallbacks.invalidDate");
                              return date.toLocaleDateString("id-ID", {
                                year: "numeric",
                                month: "long",
                                day: "numeric",
                              });
                            })()}
                          </div>
                        </div>
                      )}

                      {deal.source && (
                        <div className="space-y-2">
                          <div className="text-sm text-muted-foreground">
                            <span>{t("sections.source")}</span>
                          </div>
                          <div className="text-base font-medium">{deal.source}</div>
                        </div>
                      )}
                    </div>
                  </CardContent>
                </Card>

                {/* Product Items */}
                {Array.isArray(deal.product_items) && deal.product_items.length > 0 && (
                  <Card className="surface-panel border-border/70">
                    <CardHeader>
                      <CardTitle className="flex items-center gap-2">
                        <Package className="h-5 w-5" />
                        {t("sections.productItems") || "Product Items"}
                      </CardTitle>
                      <CardDescription>
                        {t("sections.productItemsDescription") || "Products attached to this opportunity"}
                      </CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-4">
                      <div className="space-y-3">
                        {deal.product_items.map((item) => {
                          const qty = item.quantity ?? 0;
                          const unitPrice = item.unit_price ?? 0;
                          const discount = item.discount_amount ?? 0;
                          const subtotal = item.subtotal ?? (unitPrice * qty - discount);

                          return (
                            <div
                              key={item.id || `${item.product_id}-${item.product_name}`}
                              className="flex flex-col gap-2 rounded-lg border p-3"
                            >
                              <div className="flex items-start justify-between gap-3">
                                <div className="min-w-0">
                                  <div className="font-medium truncate">
                                    {item.product_name || t("fallbacks.unknownProduct") || "Unknown product"}
                                  </div>
                                  {item.product_sku && (
                                    <div className="text-xs text-muted-foreground">SKU: {item.product_sku}</div>
                                  )}
                                </div>
                                <div className="text-right">
                                  <div className="text-sm text-muted-foreground">
                                    {t("labels.subtotal") || "Subtotal"}
                                  </div>
                                  <div className="font-medium">
                                    {item.subtotal_formatted || formatCurrency(subtotal)}
                                  </div>
                                </div>
                              </div>

                              <div className="grid grid-cols-2 md:grid-cols-4 gap-2 text-sm">
                                <div>
                                  <div className="text-muted-foreground">{t("labels.quantity") || "Qty"}</div>
                                  <div className="font-medium">{qty}</div>
                                </div>
                                <div>
                                  <div className="text-muted-foreground">{t("labels.unitPrice") || "Unit price"}</div>
                                  <div className="font-medium">{item.unit_price_formatted || formatCurrency(unitPrice)}</div>
                                </div>
                                <div>
                                  <div className="text-muted-foreground">{t("labels.discount") || "Discount"}</div>
                                  <div className="font-medium">{item.discount_amount_formatted || formatCurrency(discount)}</div>
                                </div>
                                <div>
                                  <div className="text-muted-foreground">{t("labels.lineTotal") || "Line total"}</div>
                                  <div className="font-medium">{item.subtotal_formatted || formatCurrency(subtotal)}</div>
                                </div>
                              </div>

                              {item.notes && (
                                <div className="text-sm">
                                  <div className="text-muted-foreground">{t("labels.notes") || "Notes"}</div>
                                  <div className="whitespace-pre-wrap">{item.notes}</div>
                                </div>
                              )}
                            </div>
                          );
                        })}
                      </div>

                      <div className="flex items-center justify-between rounded-lg bg-muted/50 p-3">
                        <div className="text-sm text-muted-foreground">
                          {t("labels.itemsTotal") || "Items total"}
                        </div>
                        <div className="font-medium">
                          {deal.value_formatted || formatCurrency(deal.value ?? 0)}
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                )}

                {(deal.qualification_snapshot || deal.budget_confirmed || deal.authority_confirmed || deal.need_confirmed || deal.timeline_confirmed) && (
                  <Card className="surface-panel border-border/70">
                    <CardHeader>
                      <CardTitle>{t("sections.bant") || "BANT"}</CardTitle>
                      <CardDescription>
                        {t("sections.bantDescription") || "Qualification data carried from the source lead"}
                      </CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-4">
                      <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
                        {[
                          { label: "Budget", value: deal.budget_confirmed },
                          { label: "Authority", value: deal.authority_confirmed },
                          { label: "Need", value: deal.need_confirmed },
                          { label: "Timeline", value: deal.timeline_confirmed },
                        ].map((item) => (
                          <div key={item.label} className="rounded-lg border border-border/70 bg-muted/30 p-3">
                            <div className="text-xs text-muted-foreground">{item.label}</div>
                            <Badge variant={item.value ? "default" : "outline"} className="mt-2">
                              {item.value ? (t("labels.confirmed") || "Confirmed") : (t("labels.pending") || "Pending")}
                            </Badge>
                          </div>
                        ))}
                      </div>
                      {Array.isArray(deal.qualification_snapshot?.need_target_products) && deal.qualification_snapshot.need_target_products.length > 0 && (
                        <div className="space-y-2">
                          <div className="text-sm font-medium">{t("sections.productInterest") || "Product Interest"}</div>
                          <div className="flex flex-wrap gap-2">
                            {deal.qualification_snapshot.need_target_products.map((product, index) => (
                              <Badge key={product.product_id || `${product.product_name}-${index}`} variant="secondary">
                                {product.product_name || t("fallbacks.unknownProduct") || "Unknown product"}
                              </Badge>
                            ))}
                          </div>
                        </div>
                      )}
                    </CardContent>
                  </Card>
                )}

                {/* Notes */}
                {deal.notes && (
                  <Card className="surface-panel border-border/70">
                    <CardHeader>
                      <CardTitle className="flex items-center gap-2">
                        <FileText className="h-5 w-5" />
                        {t("sections.notes")}
                      </CardTitle>
                    </CardHeader>
                    <CardContent>
                      <p className="text-base whitespace-pre-wrap">{deal.notes}</p>
                    </CardContent>
                  </Card>
                )}

                {/* Metadata */}
                <Card className="surface-panel border-border/70">
                  <CardHeader>
                    <CardTitle>{t("sections.metadata")}</CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-2 text-sm text-muted-foreground">
                    <p>
                      {t("sections.createdAt")}{" "}
                      {new Date(deal.created_at).toLocaleString("id-ID")}
                    </p>
                    <p>
                      {t("sections.updatedAt")}{" "}
                      {new Date(deal.updated_at).toLocaleString("id-ID")}
                    </p>
                  </CardContent>
                </Card>
              </div>

              {/* Right Column: Related Activities (Scrollable) */}
              <div className="flex flex-col h-full">
                <div className="mb-4">
                  <h3 className="text-lg font-medium">{t("sections.relatedActivities")}</h3>
                  <p className="text-sm text-muted-foreground">
                    {t("sections.relatedActivitiesDescription")}
                  </p>
                </div>
                
                <Card className="surface-panel flex-1 flex flex-col overflow-hidden border-border/70">
                  <CardContent className="p-4 flex-1 overflow-hidden">
                    <Tabs defaultValue="visit-reports" className="w-full h-full flex flex-col">
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
                      <TabsContent value="visit-reports" className="mt-4 flex-1 overflow-y-auto">
                        <DealVisitReportsList
                          dealId={dealId || ""}
                          onOpenVisit={(visitId) => setViewingVisitReportId(visitId)}
                        />
                      </TabsContent>
                      <TabsContent value="activities" className="mt-4 flex-1 overflow-y-auto">
                        <DealActivitiesList dealId={dealId || ""} />
                      </TabsContent>
                      <TabsContent value="tasks" className="mt-4 flex-1 overflow-y-auto">
                        <DealTasksList
                          dealId={dealId || ""}
                          onOpenTask={(taskId) => setViewingTaskId(taskId)}
                        />
                      </TabsContent>
                      <TabsContent value="product-interest" className="mt-4 flex-1 overflow-y-auto">
                        <DealProductInterestList dealId={dealId || ""} />
                      </TabsContent>
                    </Tabs>
                  </CardContent>
                </Card>
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>

      {/* Edit Dialog */}
      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent className="sm:max-w-[600px]">
          <DialogHeader>
            <DialogTitle>{tCommon("edit")}</DialogTitle>
          </DialogHeader>
          {deal && (
            <DealForm
              deal={deal}
              onSubmit={handleUpdate}
              onCancel={() => setIsEditDialogOpen(false)}
              isLoading={updateDeal.isPending}
            />
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
        isLoading={deleteDeal.isPending}
      />

      {/* Account Detail Modal */}
      <AccountDetailModal
        accountId={viewingAccountId}
        open={!!viewingAccountId}
        onOpenChange={(open) => !open && setViewingAccountId(null)}
      />

      {/* Contact Detail Modal */}
      <ContactDetailModal
        contactId={viewingContactId}
        open={!!viewingContactId}
        onOpenChange={(open) => !open && setViewingContactId(null)}
      />

      <VisitReportDetailModal
        visitReportId={viewingVisitReportId}
        open={!!viewingVisitReportId}
        onOpenChange={(open) => {
          if (!open) setViewingVisitReportId(null);
        }}
      />

      <TaskDetailModal
        taskId={viewingTaskId}
        open={!!viewingTaskId}
        onOpenChange={(open) => {
          if (!open) setViewingTaskId(null);
        }}
      />
    </>
  );
}

// Visit Reports List Component for Deal
function DealVisitReportsList({
  dealId,
  onOpenVisit,
}: {
  readonly dealId: string;
  readonly onOpenVisit: (visitId: string) => void;
}) {
  const { data, isLoading } = useDealVisitReports(dealId);
  const t = useTranslations("deals.detail");

  if (isLoading) {
    return (
      <div className="space-y-3">
        {[...Array(3)].map((_, i) => (
          <Skeleton key={i} className="h-16 w-full" />
        ))}
      </div>
    );
  }

  const visitReports = (data ?? []) as Array<{
    id: string;
    status?: string;
    visit_date?: string;
    purpose?: string;
    account?: { name: string };
    photos?: string[];
  }>;

  if (visitReports.length === 0) {
    return (
      <div className="text-center py-8 text-muted-foreground">
        <p className="text-sm">{t("sections.noVisitReports")}</p>
      </div>
    );
  }

  // Helper to convert relative photo URL to absolute
  const getPhotoUrl = (photoUrl: string): string => {
    if (photoUrl.startsWith("http://") || photoUrl.startsWith("https://")) {
      return photoUrl;
    }
    const cleanUrl = photoUrl.startsWith("/") ? photoUrl : `/${photoUrl}`;
    if (typeof window !== "undefined") {
      const baseUrl = process.env.NEXT_PUBLIC_API_URL || window.location.origin;
      return new URL(cleanUrl, baseUrl).toString();
    }
    return cleanUrl;
  };

  return (
    <div className="space-y-3">
      {visitReports.map((vr) => (
        <button
          key={vr.id}
          type="button"
          onClick={() => onOpenVisit(vr.id)}
          className="flex w-full items-start gap-3 rounded-lg border p-3 text-left transition-colors hover:bg-accent/50"
        >
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 mb-1">
              <Badge variant={vr.status === "completed" ? "default" : "secondary"}>
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

            {/* Photo thumbnails */}
            {vr.photos && vr.photos.length > 0 && (
              <div className="flex items-center gap-2 mt-2">
                <MapPin className="h-3 w-3 text-muted-foreground" />
                <div className="flex gap-1">
                  {vr.photos.slice(0, 3).map((photo, idx) => (
                    <div key={idx} className="relative w-12 h-12 border rounded overflow-hidden">
                      <img
                        src={getPhotoUrl(photo)}
                        alt={`Photo ${idx + 1}`}
                        className="w-full h-full object-cover"
                      />
                    </div>
                  ))}
                  {vr.photos.length > 3 && (
                    <div className="w-12 h-12 border rounded flex items-center justify-center bg-muted text-xs text-muted-foreground">
                      +{vr.photos.length - 3}
                    </div>
                  )}
                </div>
              </div>
            )}
          </div>
        </button>
      ))}
    </div>
  );
}

// Activities List Component for Deal
function DealActivitiesList({ dealId }: { readonly dealId: string }) {
  const { data, isLoading } = useDealActivities(dealId);
  const t = useTranslations("deals.detail");

  if (isLoading) {
    return (
      <div className="space-y-3">
        {[...Array(3)].map((_, i) => (
          <Skeleton key={i} className="h-16 w-full" />
        ))}
      </div>
    );
  }

  const activities = (data ?? []) as Array<{
    id: string;
    type?: string;
    timestamp?: string;
    description?: string;
    metadata?: Record<string, unknown>;
    account?: { id: string; name: string };
    contact?: { id: string; name: string };
    user?: { id: string; name: string };
    activity_type?: {
      id: string;
      name: string;
      icon?: string;
      badge_color?: string;
    };
  }>;

  if (activities.length === 0) {
    return (
      <div className="text-center py-8 text-muted-foreground">
        <p className="text-sm">{t("sections.noActivities")}</p>
      </div>
    );
  }

  return (
    <div className="space-y-3">
      {activities.map((activity) => {
        const activityType = activity.activity_type;
        const iconName = activityType?.icon;
        const badgeColor = activityType?.badge_color as "default" | "secondary" | "destructive" | "outline" | undefined;
        const typeName = activityType?.name ?? activity.type;

        return (
          <div key={activity.id} className="flex items-start gap-3 p-3 border rounded-lg hover:bg-accent/50 transition-colors">
            <div className="flex-1 min-w-0 space-y-2">
              {/* Type and timestamp */}
              <div className="flex items-center gap-2 flex-wrap">
                <div className="flex items-center gap-2">
                  <div className="text-muted-foreground">
                    {iconName ? renderIcon(iconName, "h-4 w-4") : <Activity className="h-4 w-4" />}
                  </div>
                  <Badge variant={badgeColor ?? "outline"} className="text-xs">
                    {typeName}
                  </Badge>
                </div>
                <span className="text-xs text-muted-foreground">
                  {(() => {
                    let dateStr = activity.timestamp;
                    // For VISIT type activities, use visit_date from metadata
                    if (activity.type === "visit" && activity.metadata && typeof activity.metadata === "object") {
                      const meta = activity.metadata as Record<string, unknown>;
                      if (typeof meta.visit_date === "string") {
                        dateStr = meta.visit_date;
                      }
                    }
                    return dateStr ? new Date(dateStr).toLocaleDateString("id-ID", {
                      year: "numeric",
                      month: "short",
                      day: "numeric",
                      hour: "2-digit",
                      minute: "2-digit",
                    }) : "";
                  })()}
                </span>
              </div>

              {/* Description */}
              <p className="text-sm font-medium line-clamp-2">{activity.description}</p>

              {/* Related entities */}
              <div className="flex items-center gap-3 flex-wrap text-xs text-muted-foreground">
                {activity.account && (
                  <div className="flex items-center gap-1">
                    <Building2 className="h-3 w-3" />
                    <span>{activity.account.name}</span>
                  </div>
                )}
                {activity.contact && (
                  <div className="flex items-center gap-1">
                    <Contact className="h-3 w-3" />
                    <span>{activity.contact.name}</span>
                  </div>
                )}
                {activity.user && (
                  <div className="flex items-center gap-1">
                    <User className="h-3 w-3" />
                    <span>{activity.user.name}</span>
                  </div>
                )}
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
}

function DealTasksList({
  dealId,
  onOpenTask,
}: {
  readonly dealId: string;
  readonly onOpenTask: (taskId: string) => void;
}) {
  const { data, isLoading } = useTasks({ deal_id: dealId, per_page: 10 });
  const t = useTranslations("deals.detail");

  if (isLoading) {
    return (
      <div className="space-y-3">
        {[...Array(3)].map((_, i) => (
          <Skeleton key={`deal-task-skeleton-${i}`} className="h-16 w-full" />
        ))}
      </div>
    );
  }

  const tasks = data?.data ?? [];

  if (tasks.length === 0) {
    return (
      <div className="text-center py-8 text-muted-foreground">
        <p className="text-sm">{t("sections.noTasks")}</p>
      </div>
    );
  }

  return (
    <div className="space-y-3">
      {tasks.map((task) => (
        <TaskCard key={task.id} task={task} onOpenTask={onOpenTask} />
      ))}
      {data?.meta?.pagination && data.meta.pagination.total > tasks.length && (
        <div className="text-center pt-2">
          <p className="text-xs text-muted-foreground">
            {tasks.length} / {data.meta.pagination.total} {t("sections.tasks")}
          </p>
        </div>
      )}
    </div>
  );
}

function DealProductInterestList({ dealId }: { readonly dealId: string }) {
  const { data: visits, isLoading: isVisitsLoading } = useDealVisitReports(dealId);
  const { data: activities, isLoading: isActivitiesLoading } = useDealActivities(dealId);

  return (
    <ProductInterestTab
      activities={buildProductInterestActivities(
        (visits ?? []) as VisitReport[],
        (activities ?? []) as CRMActivity[],
      )}
      isLoading={isVisitsLoading || isActivitiesLoading}
    />
  );
}

function TaskCard({
  task,
  onOpenTask,
}: {
  readonly task: Task;
  readonly onOpenTask: (taskId: string) => void;
}) {
  const dueDate = task.due_date ? new Date(task.due_date) : null;
  const dueLabel =
    dueDate && !Number.isNaN(dueDate.getTime())
      ? dueDate.toLocaleDateString("id-ID", {
          year: "numeric",
          month: "short",
          day: "numeric",
        })
      : null;

  return (
    <button
      type="button"
      onClick={() => onOpenTask(task.id)}
      className="flex w-full items-start gap-3 rounded-lg border border-border/70 p-3 text-left transition-colors hover:bg-accent/40"
    >
      <div className="flex-1 min-w-0">
        <div className="mb-1 flex flex-wrap items-center gap-2">
          <Badge variant={task.status === "completed" ? "default" : "outline"} className="capitalize">
            {task.status}
          </Badge>
          <Badge variant="secondary" className="capitalize">
            {task.priority}
          </Badge>
          {dueLabel && <span className="text-sm text-muted-foreground">{dueLabel}</span>}
        </div>
        <p className="text-sm font-medium line-clamp-1">{task.title}</p>
        {task.description && <p className="mt-1 text-xs text-muted-foreground line-clamp-2">{task.description}</p>}
        {task.assigned_user?.name && (
          <p className="mt-2 text-xs text-muted-foreground">
            {task.assigned_user.name}
          </p>
        )}
      </div>
    </button>
  );
}

function buildProductInterestActivities(visits: VisitReport[], activities: CRMActivity[]): CRMActivity[] {
  const visitActivities: CRMActivity[] = visits.map((visit) => ({
    id: `deal-visit-report-${visit.id}`,
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
