"use client";

import { useState, useMemo } from "react";
import { useTranslations } from "next-intl";
import {
  Plus,
  MoreVertical,
  Eye,
  Edit,
  Trash2,
  TrendingUp,
  Calendar,
  DollarSign,
  Building2,
  User,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { CreateDealDialog } from "./create-deal-dialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import { DataTable, type Column } from "@/components/ui/data-table";
import { useDeals, useUpdateDeal, useDeleteDeal } from "../hooks/useDeals";
import { useStages } from "../hooks/useStages";
import { DealForm } from "./deal-form";
import { DealDetailModal } from "./deal-detail-modal";
import { PipelineFilters } from "./pipeline-filters";
import type { Deal, DealFilters } from "../types";
import type { UpdateDealFormData } from "../schemas/deal.schema";
import { useHasPermission } from "@/features/master-data/user-management/hooks/useHasPermission";
import { toast } from "sonner";
import { formatCurrency } from "@/lib/utils";
import { getProbabilityColor } from "../utils/color";

interface PipelineTableViewProps {
  readonly onDealClick?: (deal: { id: string }) => void;
}

export function PipelineTableView({ onDealClick }: PipelineTableViewProps) {
  const t = useTranslations("pipelineManagement.tableView");
  const [filters, setFilters] = useState<DealFilters>({});
  
  // Pagination state
  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(10);
  
  const [editingDeal, setEditingDeal] = useState<Deal | null>(null);
  const [deletingDeal, setDeletingDeal] = useState<Deal | null>(null);
  const [viewingDealId, setViewingDealId] = useState<string | null>(null);
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);

  // Pass pagination and sort params. Sort by created_at desc (newest first)
  const { data: dealsResponse, isLoading } = useDeals(
    filters, 
    page, 
    perPage,
    "created_at",
    "desc"
  );
  
  const { data: stagesData } = useStages();
  const updateDeal = useUpdateDeal();
  const deleteDeal = useDeleteDeal();

  const stages = Array.isArray(stagesData) ? stagesData : [];
  const deals = dealsResponse?.data ?? [];
  const pagination = dealsResponse?.meta?.pagination;

  const hasCreatePermission = useHasPermission("pipeline.opportunity-create");
  const hasEditPermission = useHasPermission("pipeline.opportunity-edit");
  const hasDeletePermission = useHasPermission("pipeline.opportunity-delete");

  const stageColors: Record<string, string> = useMemo(() => {
    const colors: Record<string, string> = {};
    stages.forEach((stage: { id: string; color?: string }) => {
      colors[stage.id] = stage.color ?? "";
    });
    return colors;
  }, [stages]);

  const handleDeleteDeal = async () => {
    if (!deletingDeal) return;

    try {
      await deleteDeal.mutateAsync(deletingDeal.id);
      toast.success(t("toast.deleteSuccess"));
      setDeletingDeal(null);
    } catch (error) {
      toast.error(t("toast.deleteError"));
    }
  };

  const handleResetFilters = () => {
    setFilters({});
    setPage(1); // Reset to first page
  };

  const handleCreateSuccess = () => {
    setIsCreateDialogOpen(false);
  };

  const columns: Column<Deal>[] = [
    {
      id: "title",
      header: t("table.title"),
      accessor: (row) => (
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
            <TrendingUp className="h-5 w-5 text-primary" />
          </div>
          <div>
            <div className="font-medium cursor-pointer hover:underline" onClick={() => {
              setViewingDealId(row.id);
              onDealClick?.({ id: row.id });
            }}>
              {row.title}
            </div>
            {row.description && (
              <div className="text-sm text-muted-foreground line-clamp-1">
                {row.description}
              </div>
            )}
          </div>
        </div>
      ),
      className: "min-w-[250px]",
    },
    {
      id: "account",
      header: t("table.account"),
      accessor: (row) => (
        <div className="flex items-center gap-2">
          <Building2 className="h-4 w-4 text-muted-foreground" />
          <span>{row.account?.name ?? "-"}</span>
        </div>
      ),
      className: "w-[200px]",
    },
    {
      id: "contact",
      header: t("table.contact"),
      accessor: (row) => (
        <div className="flex items-center gap-2">
          <User className="h-4 w-4 text-muted-foreground" />
          <span>{row.contact?.name ?? "-"}</span>
        </div>
      ),
      className: "w-[180px]",
    },
    {
      id: "value",
      header: t("table.value"),
      accessor: (row) => (
        <div className="flex items-center gap-2 font-medium">
          <DollarSign className="h-4 w-4 text-muted-foreground" />
          <span>{formatCurrency(row.value)}</span>
        </div>
      ),
      className: "w-[150px]",
    },
    {
      id: "stage",
      header: t("table.stage"),
      accessor: (row) => (
        <Badge
          variant="outline"
          style={{
            backgroundColor: `${stageColors[row.stage_id]}20`,
            borderColor: stageColors[row.stage_id],
            color: stageColors[row.stage_id],
          }}
        >
          {row.stage?.name ?? "-"}
        </Badge>
      ),
      className: "w-[150px]",
    },
    {
      id: "expected_close_date",
      header: t("table.expectedCloseDate"),
      accessor: (row) => (
        <div className="flex items-center gap-2 text-sm">
          <Calendar className="h-4 w-4 text-muted-foreground" />
          <span>
            {row.expected_close_date
              ? new Date(row.expected_close_date).toLocaleDateString()
              : "-"}
          </span>
        </div>
      ),
      className: "w-[150px]",
    },
    {
      id: "probability",
      header: t("table.probability"),
      accessor: (row) => {
        const probColor = getProbabilityColor(row.probability ?? 0);
        return (
          <div className="flex items-center gap-2">
            <div className="flex-1 h-2 bg-muted rounded-full overflow-hidden">
              <div
                className="h-full transition-all"
                style={{ 
                  width: `${row.probability || 0}%`,
                  backgroundColor: probColor
                }}
              />
            </div>
            <span 
              className="text-sm font-medium w-12 text-right transition-colors"
              style={{ color: probColor }}
            >
              {row.probability || 0}%
            </span>
          </div>
        );
      },
      className: "w-[150px]",
    },
    {
      id: "assigned_user",
      header: t("table.assignedTo"),
      accessor: (row) => (
        <div className="flex items-center gap-2">
          <User className="h-4 w-4 text-muted-foreground" />
          <span>{row.assigned_user?.name ?? "-"}</span>
        </div>
      ),
      className: "w-[150px]",
    },
    {
      id: "actions",
      header: t("table.actions"),
      accessor: (row) => {
        const canEdit = hasEditPermission;
        const canDelete = hasDeletePermission;

        return (
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => {
                setViewingDealId(row.id);
                onDealClick?.({ id: row.id });
              }}
              className="h-8 w-8 p-0 cursor-pointer"
            >
              <Eye className="h-4 w-4" />
            </Button>

            {(canEdit || canDelete) && (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="h-8 w-8 p-0 cursor-pointer"
                  >
                    <MoreVertical className="h-4 w-4" />
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  {canEdit && (
                    <DropdownMenuItem onClick={() => setEditingDeal(row)}>
                      <Edit className="h-4 w-4 mr-2" />
                      {t("buttons.edit")}
                    </DropdownMenuItem>
                  )}
                  {canDelete && (
                    <DropdownMenuItem
                      onClick={() => setDeletingDeal(row)}
                      className="text-destructive"
                    >
                      <Trash2 className="h-4 w-4 mr-2" />
                      {t("buttons.delete")}
                    </DropdownMenuItem>
                  )}
                </DropdownMenuContent>
              </DropdownMenu>
            )}
          </div>
        );
      },
      className: "w-[100px]",
    },
  ];

  return (
    <div className="space-y-6">
      {/* Filters with Add Button */}
      <PipelineFilters
        filters={filters}
        onFiltersChange={setFilters}
        onReset={handleResetFilters}
        onAdd={hasCreatePermission ? () => setIsCreateDialogOpen(true) : undefined}
        addLabel={t("table.add") || "Add Opportunity"}
      />

      {/* Table */}
      <div className="border border-border/50 rounded-xl bg-card/30 overflow-hidden">
        <DataTable
          columns={columns}
          data={deals}
          isLoading={isLoading}
          pagination={pagination}
          onPageChange={setPage}
          onPerPageChange={setPerPage}
          emptyMessage={t("table.empty")}
        />
      </div>

      {/* Create Deal Dialog */}
      <CreateDealDialog
        open={isCreateDialogOpen}
        onOpenChange={setIsCreateDialogOpen}
        onSuccess={handleCreateSuccess}
      />

      {/* Edit Deal Dialog */}
      <Dialog open={!!editingDeal} onOpenChange={(open) => !open && setEditingDeal(null)}>
        <DialogContent className="max-w-3xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>{t("dialog.editTitle")}</DialogTitle>
            <DialogDescription>{t("dialog.editDescription")}</DialogDescription>
          </DialogHeader>
          {editingDeal && (
            <DealForm
              deal={editingDeal}
              onSubmit={async (data: UpdateDealFormData) => {
                try {
                  await updateDeal.mutateAsync({
                    id: editingDeal.id,
                    data: data as UpdateDealFormData,
                  });
                  toast.success(t("toast.updateSuccess"));
                  setEditingDeal(null);
                } catch {
                  // Error handled by interceptor
                }
              }}
              onCancel={() => setEditingDeal(null)}
              isLoading={updateDeal.isPending}
            />
          )}
        </DialogContent>
      </Dialog>

      {/* Delete Deal Dialog */}
      <Dialog open={!!deletingDeal} onOpenChange={(open) => !open && setDeletingDeal(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("dialog.deleteTitle")}</DialogTitle>
            <DialogDescription>{t("dialog.deleteDescription")}</DialogDescription>
          </DialogHeader>
          <div className="py-4">
            <p className="text-sm text-muted-foreground">
              {t("dialog.deleteConfirm", { title: deletingDeal?.title ?? "" })}
            </p>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeletingDeal(null)}>
              {t("buttons.cancel")}
            </Button>
            <Button
              variant="destructive"
              onClick={handleDeleteDeal}
              disabled={deleteDeal.isPending}
            >
              {deleteDeal.isPending ? t("buttons.deleting") : t("buttons.delete")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Deal Detail Modal */}
      <DealDetailModal
        dealId={viewingDealId}
        open={!!viewingDealId}
        onOpenChange={(open) => !open && setViewingDealId(null)}
      />
    </div>
  );
}

