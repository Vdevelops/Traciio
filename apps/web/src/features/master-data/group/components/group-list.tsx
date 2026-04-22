"use client";

import { Edit, Trash2, Plus, Search, Target, ChevronDown, ChevronRight } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import { GroupForm } from "./group-form";
import { GroupTargetForm } from "./group-target-form";
import { GroupUsersDropdown } from "./group-users-dropdown";
import { useGroupList } from "../hooks/useGroupList";
import { useState } from "react";
import type { Group } from "../types";
import { useTranslations } from "next-intl";
import type { CreateGroupFormData } from "../schemas/group.schema";
import { useIsMobile } from "@/hooks/use-mobile";
import { cn } from "@/lib/utils";

export function GroupList() {
  const {
    setPage,
    setPerPage,
    search,
    setSearch,
    status,
    setStatus,
    isCreateDialogOpen,
    setIsCreateDialogOpen,
    editingGroup,
    setEditingGroup,
    groups,
    pagination,
    editingGroupData,
    isLoading,
    handleCreate,
    handleUpdate,
    handleDeleteClick,
    handleDeleteConfirm,
    deletingGroupId,
    setDeletingGroupId,
    targetGroupId,
    setTargetGroupId,
    targetGroupData,
    handleSetTarget,
    handleCreateTarget,
    createGroup,
    updateGroup,
    createGroupTarget,
  } = useGroupList();

  const t = useTranslations("groupManagement.list");
  const tForm = useTranslations("groupManagement.form");
  const isMobile = useIsMobile();
  
  const [expandedRows, setExpandedRows] = useState<Set<string>>(new Set());

  const toggleRow = (id: string) => {
    const newExpanded = new Set(expandedRows);
    if (newExpanded.has(id)) {
      newExpanded.delete(id);
    } else {
      newExpanded.add(id);
    }
    setExpandedRows(newExpanded);
  };

  return (
    <div className="space-y-4">
      {/* Header with Actions */}
      <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-3">
        <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-2 sm:gap-3 flex-1">
          <div className="relative flex-1 max-w-sm">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder={t("searchPlaceholder")}
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              className="pl-10 h-9"
            />
          </div>
          <Select
            value={status || "all"}
            onValueChange={(value) => setStatus(value === "all" ? "" : value)}
          >
            <SelectTrigger className="w-[140px] h-9">
              <SelectValue placeholder={t("allStatus")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">{t("allStatus")}</SelectItem>
              <SelectItem value="active">{tForm("statusActive")}</SelectItem>
              <SelectItem value="inactive">{tForm("statusInactive")}</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <Button
          onClick={() => setIsCreateDialogOpen(true)}
          size="sm"
          className="cursor-pointer"
        >
          <Plus className="h-4 w-4 mr-2" />
          {t("addGroup")}
        </Button>
      </div>

      {/* Table */}
      <div className="border rounded-lg">
        {isLoading ? (
          <div className="p-4 space-y-3">
            {Array.from({ length: 5 }, (_, i) => (
              <Skeleton key={`skeleton-row-${i}`} className={cn(isMobile ? "h-48" : "h-10", "w-full")} />
            ))}
          </div>
        ) : (
          <>
            {isMobile ? (
              // Mobile Card View
              <div className="divide-y">
                {!Array.isArray(groups) || groups.length === 0 ? (
                  <div className="text-center text-muted-foreground py-12 px-4">
                    {t("empty")}
                  </div>
                ) : (
                  groups.map((group) => {
                    if (!group) return null;
                    return (
                      <GroupMobileCard
                        key={group.id}
                        group={group}
                        isExpanded={expandedRows.has(group.id)}
                        onToggle={() => toggleRow(group.id)}
                        onEdit={() => setEditingGroup(group.id)}
                        onDelete={() => handleDeleteClick(group.id)}
                        onSetTarget={() => handleSetTarget(group.id)}
                      />
                    );
                  })
                )}
              </div>
            ) : (
              // Desktop Table View
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead className="w-[50px]"></TableHead>
                    <TableHead className="w-[250px]">{t("name")}</TableHead>
                    <TableHead>{t("description")}</TableHead>
                    <TableHead className="w-[100px]">{t("status")}</TableHead>
                    <TableHead className="w-[140px] text-right">
                      {t("actions")}
                    </TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {!Array.isArray(groups) || groups.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={5} className="text-center text-muted-foreground py-8">
                        {t("empty")}
                      </TableCell>
                    </TableRow>
                  ) : (
                    groups.map((group) => {
                      if (!group) return null;
                      return (
                        <GroupRow
                          key={group.id}
                          group={group}
                          isExpanded={expandedRows.has(group.id)}
                          onToggle={() => toggleRow(group.id)}
                          onEdit={() => setEditingGroup(group.id)}
                          onDelete={() => handleDeleteClick(group.id)}
                          onSetTarget={() => handleSetTarget(group.id)}
                        />
                      );
                    })
                  )}
                </TableBody>
              </Table>
            )}

            {/* Pagination */}
            {pagination && (
              <div className="border-t bg-muted/30 px-3 sm:px-6 py-3 sm:py-4">
                <div className="flex flex-col lg:flex-row items-center justify-between gap-4">
                  <div className="flex items-center gap-3 order-3 lg:order-1">
                    <label htmlFor="rows-per-page" className="text-sm whitespace-nowrap">
                      Rows per page
                    </label>
                    <Select
                      value={String(pagination.per_page)}
                      onValueChange={(value) => {
                        setPerPage(Number(value));
                        setPage(1);
                      }}
                    >
                      <SelectTrigger
                        id="rows-per-page"
                        className="w-fit whitespace-nowrap h-9"
                      >
                        <SelectValue placeholder="Select rows" />
                      </SelectTrigger>
                      <SelectContent>
                        {[10, 20, 50, 100].map((option) => (
                          <SelectItem key={option} value={String(option)}>
                            {option}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="flex grow justify-center lg:justify-end text-sm whitespace-nowrap text-muted-foreground order-2">
                    <p>
                      <span className="text-foreground font-medium">
                        {(pagination.page - 1) * pagination.per_page + 1}-
                        {Math.min(pagination.page * pagination.per_page, pagination.total)}
                      </span>{" "}
                      of {pagination.total}
                    </p>
                  </div>
                  {pagination.total_pages > 1 && (
                    <div className="order-1 lg:order-3 flex items-center gap-2">
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setPage(1)}
                        disabled={!pagination.has_prev}
                        className="cursor-pointer"
                      >
                        First
                      </Button>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setPage(Math.max(1, pagination.page - 1))}
                        disabled={!pagination.has_prev}
                        className="cursor-pointer"
                      >
                        Prev
                      </Button>
                      <span className="text-sm">
                        Page {pagination.page} of {pagination.total_pages}
                      </span>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setPage(pagination.page + 1)}
                        disabled={!pagination.has_next}
                        className="cursor-pointer"
                      >
                        Next
                      </Button>
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setPage(pagination.total_pages)}
                        disabled={!pagination.has_next}
                        className="cursor-pointer"
                      >
                        Last
                      </Button>
                    </div>
                  )}
                </div>
              </div>
            )}
          </>
        )}
      </div>

      {/* Create Dialog */}
      <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("createTitle")}</DialogTitle>
          </DialogHeader>
          <GroupForm
            onSubmit={async (data) => {
              await handleCreate(data as CreateGroupFormData);
            }}
            onCancel={() => setIsCreateDialogOpen(false)}
            isLoading={createGroup.isPending}
          />
        </DialogContent>
      </Dialog>

      {/* Edit Dialog */}
      <Dialog open={!!editingGroup} onOpenChange={(open) => !open && setEditingGroup(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{t("editTitle")}</DialogTitle>
          </DialogHeader>
          {editingGroupData && (
            <GroupForm
              group={editingGroupData}
              onSubmit={handleUpdate}
              onCancel={() => setEditingGroup(null)}
              isLoading={updateGroup.isPending}
            />
          )}
        </DialogContent>
      </Dialog>

      {/* Delete Dialog */}
      <DeleteDialog
        open={!!deletingGroupId}
        onOpenChange={(open) => !open && setDeletingGroupId(null)}
        onConfirm={handleDeleteConfirm}
        title={t("deleteTitle")}
        description={t("deleteDescription")}
      />

      {/* Set Target Dialog */}
      <Dialog
        open={!!targetGroupId}
        onOpenChange={(open) => !open && setTargetGroupId(null)}
      >
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>{t("setTargetTitle")}</DialogTitle>
            <DialogDescription>
              Set monthly target for this group. The target will be automatically assigned to all users in the group.
            </DialogDescription>
          </DialogHeader>
          {targetGroupData && (
            <GroupTargetForm
              group={targetGroupData}
              onSubmit={handleCreateTarget}
              onCancel={() => setTargetGroupId(null)}
              isLoading={createGroupTarget.isPending}
            />
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
}

interface GroupRowProps {
  readonly group: Group;
  readonly isExpanded: boolean;
  readonly onToggle: () => void;
  readonly onEdit: () => void;
  readonly onDelete: () => void;
  readonly onSetTarget: () => void;
}

function GroupRow({
  group,
  isExpanded,
  onToggle,
  onEdit,
  onDelete,
  onSetTarget,
}: GroupRowProps) {
  return (
    <>
      <TableRow className="hover:bg-muted/50">
        <TableCell>
          <Button
            variant="ghost"
            size="icon-sm"
            className="h-6 w-6 cursor-pointer"
            onClick={onToggle}
          >
            {isExpanded ? (
              <ChevronDown className="h-4 w-4" />
            ) : (
              <ChevronRight className="h-4 w-4" />
            )}
          </Button>
        </TableCell>
        <TableCell>
          <div className="flex flex-col">
            <span className="font-medium">{group.name}</span>
            <span className="text-xs text-muted-foreground">{group.code}</span>
          </div>
        </TableCell>
        <TableCell>
          <span className="text-muted-foreground">
            {group.description || "-"}
          </span>
        </TableCell>
        <TableCell>
          <Badge variant={group.status === "active" ? "active" : "inactive"}>
            {group.status}
          </Badge>
        </TableCell>
        <TableCell className="text-right">
          <div className="flex items-center justify-end gap-1">
            <Button
              variant="ghost"
              size="icon-sm"
              onClick={onSetTarget}
              className="h-8 w-8 cursor-pointer"
              title="Set Target"
            >
              <Target className="h-3.5 w-3.5" />
            </Button>
            <Button
              variant="ghost"
              size="icon-sm"
              onClick={onEdit}
              className="h-8 w-8 cursor-pointer"
              title="Edit"
            >
              <Edit className="h-3.5 w-3.5" />
            </Button>
            <Button
              variant="ghost"
              size="icon-sm"
              onClick={onDelete}
              className="h-8 w-8 text-destructive hover:text-destructive cursor-pointer"
              title="Delete"
            >
              <Trash2 className="h-3.5 w-3.5" />
            </Button>
          </div>
        </TableCell>
      </TableRow>
      {isExpanded && (
        <TableRow>
          <TableCell colSpan={5} className="bg-muted/20 p-0">
            <div className="p-6">
              <GroupUsersDropdown groupId={group.id} groupName={group.name} />
            </div>
          </TableCell>
        </TableRow>
      )}
    </>
  );
}

interface GroupMobileCardProps {
  readonly group: Group;
  readonly isExpanded: boolean;
  readonly onToggle: () => void;
  readonly onEdit: () => void;
  readonly onDelete: () => void;
  readonly onSetTarget: () => void;
}

function GroupMobileCard({
  group,
  isExpanded,
  onToggle,
  onEdit,
  onDelete,
  onSetTarget,
}: GroupMobileCardProps) {
  return (
    <div className="p-4 space-y-3">
      <div className="flex items-start gap-3">
        <button
          onClick={onToggle}
          className="mt-0.5 shrink-0 cursor-pointer"
        >
          {isExpanded ? (
            <ChevronDown className="h-4 w-4 text-muted-foreground" />
          ) : (
            <ChevronRight className="h-4 w-4 text-muted-foreground" />
          )}
        </button>
        <div className="flex-1 min-w-0">
          <p className="font-medium text-base">{group.name}</p>
          <p className="text-xs text-muted-foreground">{group.code}</p>
          <div className="flex items-center gap-2 mt-1.5">
            <Badge variant={group.status === "active" ? "active" : "inactive"} className="text-xs">
              {group.status}
            </Badge>
          </div>
        </div>
        <div className="flex items-center gap-1 shrink-0">
          <Button
            variant="ghost"
            size="icon-sm"
            className="h-8 w-8 cursor-pointer"
            onClick={onSetTarget}
            title="Set Target"
          >
            <Target className="h-4 w-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon-sm"
            onClick={onEdit}
            className="h-8 w-8 cursor-pointer"
            title="Edit"
          >
            <Edit className="h-4 w-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon-sm"
            onClick={onDelete}
            className="h-8 w-8 text-destructive hover:text-destructive cursor-pointer"
            title="Delete"
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </div>
      </div>

      {group.description && (
        <div className="pl-7 text-sm text-muted-foreground">
          {group.description}
        </div>
      )}

      {isExpanded && (
        <div className="pl-7 pt-2">
          <GroupUsersDropdown groupId={group.id} groupName={group.name} />
        </div>
      )}
    </div>
  );
}

