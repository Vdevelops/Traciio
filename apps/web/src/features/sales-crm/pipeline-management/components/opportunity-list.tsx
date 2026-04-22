"use client";

import { useState, useMemo } from "react";
import { useTranslations } from "next-intl";
import {
  Plus,
  Search,
  Filter,
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
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
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
import { useDeals, useCreateDeal, useUpdateDeal, useDeleteDeal } from "../hooks/useDeals";
import { usePipelines } from "../hooks/usePipelines";
import { DealForm } from "./deal-form";
import { DealDetailModal } from "./deal-detail-modal";
import type { Deal } from "../types";
import type { CreateDealFormData, UpdateDealFormData } from "../schemas/deal.schema";
import { useHasPermission } from "@/features/master-data/user-management/hooks/useHasPermission";
import { toast } from "sonner";
import { formatCurrency } from "@/lib/utils";
import { getProbabilityColor } from "../utils/color";

export function OpportunityList({ onDealClick }: { onDealClick?: (deal: { id: string }) => void }) {
  const t = useTranslations("pipelineManagement.opportunityList");
  const [searchQuery, setSearchQuery] = useState("");
  const [stageFilter, setStageFilter] = useState<string>("all");
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [editingDeal, setEditingDeal] = useState<Deal | null>(null);
  const [deletingDeal, setDeletingDeal] = useState<Deal | null>(null);
  const [viewingDealId, setViewingDealId] = useState<string | null>(null);

  const { data: dealsData, isLoading } = useDeals({});
  const { data: stagesData } = usePipelines({ is_active: true });
  const createDeal = useCreateDeal();
  const updateDeal = useUpdateDeal();
  const deleteDeal = useDeleteDeal();

  const stages = stagesData?.data ?? [];
  const deals = Array.isArray(dealsData) ? dealsData : [];

  const hasCreatePermission = useHasPermission("pipeline.opportunity-create");
  const hasEditPermission = useHasPermission("pipeline.opportunity-edit");
  const hasDeletePermission = useHasPermission("pipeline.opportunity-delete");

  const filteredDeals = useMemo(() => {
    return deals.filter((deal) => {
      const matchesSearch =
        deal.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
        deal.account?.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        deal.contact?.name?.toLowerCase().includes(searchQuery.toLowerCase());

      const matchesStage = stageFilter === "all" || deal.stage_id === stageFilter;

      return matchesSearch && matchesStage;
    });
  }, [deals, searchQuery, stageFilter]);

  const stageColors: Record<string, string> = useMemo(() => {
    const colors: Record<string, string> = {};
    stages.forEach((stage) => {
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
            <div className="font-medium">{row.title}</div>
            {row.description && (
              <div className="text-sm text-muted-foreground line-clamp-1">{row.description}</div>
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
          <span>{row.account?.name || "-"}</span>
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
          <span>{row.contact?.name || "-"}</span>
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
          {row.stage?.name || "-"}
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
              onClick={() => { setViewingDealId(row.id); onDealClick?.({ id: row.id }); }}
              className="h-8 w-8 p-0"
            >
              <Eye className="h-4 w-4" />
            </Button>

            {(canEdit || canDelete) && (
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
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
      <Card>
        <CardContent className="pt-6">
          <div className="flex items-center justify-between gap-4 flex-wrap">
            <div className="flex items-center gap-4 flex-1 min-w-[300px]">
              <div className="relative flex-1 max-w-md">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder={t("search.placeholder")}
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  className="pl-10"
                />
              </div>

              <Select value={stageFilter} onValueChange={setStageFilter}>
                <SelectTrigger className="w-[200px]">
                  <Filter className="h-4 w-4 mr-2" />
                  <SelectValue placeholder={t("filter.stage")} />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">{t("filter.allStages")}</SelectItem>
                  {stages.map((stage) => (
                    <SelectItem key={stage.id} value={stage.id}>
                      {stage.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            {hasCreatePermission && (
              <Button onClick={() => setIsCreateDialogOpen(true)}>
                <Plus className="h-4 w-4 mr-2" />
                {t("buttons.create")}
              </Button>
            )}
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardContent className="pt-6">
          <DataTable
            columns={columns}
            data={filteredDeals}
            isLoading={isLoading}
            emptyMessage={t("table.empty")}
          />
        </CardContent>
      </Card>

      <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
        <DialogContent className="max-w-3xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>{t("dialog.createTitle")}</DialogTitle>
            <DialogDescription>{t("dialog.createDescription")}</DialogDescription>
          </DialogHeader>
          <DealForm
            onSubmit={async (data: CreateDealFormData) => {
              try {
                await createDeal.mutateAsync(data as CreateDealFormData);
                toast.success(t("toast.createSuccess") || "Opportunity created successfully");
                setIsCreateDialogOpen(false);
              } catch {
                // Error handled by interceptor
              }
            }}
            onCancel={() => setIsCreateDialogOpen(false)}
            isLoading={createDeal.isPending}
          />
        </DialogContent>
      </Dialog>

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
                  toast.success(t("toast.updateSuccess") || "Opportunity updated successfully");
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

      <Dialog open={!!deletingDeal} onOpenChange={(open) => !open && setDeletingDeal(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("dialog.deleteTitle")}</DialogTitle>
            <DialogDescription>{t("dialog.deleteDescription")}</DialogDescription>
          </DialogHeader>
          <div className="py-4">
            <p className="text-sm text-muted-foreground">
              {t("dialog.deleteConfirm", { title: deletingDeal?.title || "" })}
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

      <DealDetailModal
        dealId={viewingDealId}
        open={!!viewingDealId}
        onOpenChange={(open) => !open && setViewingDealId(null)}
      />
    </div>
  );
}
