"use client";

import { Suspense, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { useTranslations } from "next-intl";
import { AuthGuard } from "@/features/auth/components/auth-guard";
import { PermissionGuard } from "@/features/auth/components/permission-guard";
import { PageMotion } from "@/components/motion";
import { PageDetailLayout } from "@/components/layouts/page-detail-layout";
import { DealDetailTabs } from "@/features/sales-crm/pipeline-management/components/DealDetailTabs";
import { useDeal, useUpdateDeal, useDeleteDeal } from "@/features/sales-crm/pipeline-management/hooks/useDeals";
import { DealForm } from "@/features/sales-crm/pipeline-management/components/deal-form";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Trash2, Pencil, FileText } from "lucide-react";
import { formatCurrency } from "@/features/sales-crm/pipeline-management/utils/format";
import { toast } from "sonner";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog";

function DealDetailPageContent() {
  const params = useParams();
  const router = useRouter();
  const dealId = params.id as string;
  const t = useTranslations("deals.detail");
  const tCommon = useTranslations("common");

  const { data: deal, isLoading } = useDeal(dealId);
  const updateDeal = useUpdateDeal();
  const deleteDeal = useDeleteDeal();
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);

  const handleUpdate = async (formData: Parameters<typeof updateDeal.mutateAsync>[0]["data"]) => {
    try {
      await updateDeal.mutateAsync({ id: dealId, data: formData });
      setIsEditDialogOpen(false);
      toast.success(t("toast.updated"));
    } catch {
      // Error already handled in api-client interceptor
    }
  };

  const handleDelete = async () => {
    try {
      await deleteDeal.mutateAsync(dealId);
      setIsDeleteDialogOpen(false);
      toast.success(t("toast.deleted"));
      router.push("/pipeline");
    } catch {
      // Error already handled in api-client interceptor
    }
  };

  const valueFormatted = deal?.value_formatted ?? formatCurrency(deal?.value ?? 0);

  return (
    <PageMotion className="p-2 sm:p-4">
      <PageDetailLayout
        title={
          isLoading ? (
            <Skeleton className="w-64 h-8" />
          ) : deal ? (
            deal.title
          ) : (
            "Deal Not Found"
          )
        }
        subtitle={
          deal ? (
            <div className="flex flex-wrap items-center gap-2 mt-1">
              <span className="font-medium">{valueFormatted}</span>
              {deal.stage && (
                <>
                  <div
                    className="w-3 h-3 rounded-full shrink-0"
                    style={{ backgroundColor: deal.stage.color ?? '#6366f1' }}
                  />
                  <span>{deal.stage.name}</span>
                </>
              )}
              <Badge
                variant={
                  deal.status === "won"
                    ? "default"
                    : deal.status === "lost"
                      ? "destructive"
                      : "secondary"
                }
                className="text-xs"
              >
                {(deal.status ?? 'open').toUpperCase()}
              </Badge>
              {deal.probability > 0 && (
                <span className="text-xs text-muted-foreground">
                  {deal.probability}% probability
                </span>
              )}
            </div>
          ) : undefined
        }
        backHref="/pipeline"
        actions={
          deal ? (
            <>
              <Button
                type="button"
                variant="default"
                size="sm"
                onClick={() => toast.info("Convert to Quotation / Sales Order flow will be implemented in subsequent phases", { icon: "📝" })}
                className="cursor-pointer bg-blue-600 hover:bg-blue-700 text-white"
              >
                <FileText className="h-4 w-4 mr-2" />
                Convert to Quotation
              </Button>
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => setIsEditDialogOpen(true)}
                className="cursor-pointer"
              >
                <Pencil className="h-4 w-4 mr-2" />
                {tCommon("edit")}
              </Button>
              <Button
                type="button"
                variant="destructive"
                size="sm"
                onClick={() => setIsDeleteDialogOpen(true)}
                className="cursor-pointer"
              >
                <Trash2 className="h-4 w-4 mr-2" />
                {tCommon("delete")}
              </Button>
            </>
          ) : undefined
        }
      >
        <Suspense fallback={<Skeleton className="h-[600px] w-full" />}>
          <DealDetailTabs dealId={dealId} />
        </Suspense>
      </PageDetailLayout>

      {/* Edit Dialog */}
      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent className="sm:max-w-[600px] max-h-[90vh] overflow-y-auto">
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
      <Dialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("deleteDialog.title")}</DialogTitle>
            <DialogDescription>
              {t("deleteDialog.description")}
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => setIsDeleteDialogOpen(false)}
              className="cursor-pointer"
            >
              {t("deleteDialog.cancel")}
            </Button>
            <Button
              type="button"
              variant="destructive"
              onClick={handleDelete}
              disabled={deleteDeal.isPending}
              className="cursor-pointer"
            >
              {deleteDeal.isPending
                ? t("deleteDialog.confirmLoading")
                : t("deleteDialog.confirm")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </PageMotion>
  );
}

export default function DealDetailPage() {
  return (
    <AuthGuard>
      <PermissionGuard requiredPermission="pipeline.view">
        <DealDetailPageContent />
      </PermissionGuard>
    </AuthGuard>
  );
}
