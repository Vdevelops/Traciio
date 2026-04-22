"use client";

import { useState, useEffect, useMemo } from "react";
import Link from "next/link";
import { Edit, Trash2, Plus, Search, Eye, Target } from "lucide-react";
import { StatusSwitch } from "@/components/ui/status-switch";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarImage } from "@/components/ui/avatar";
import { DataTable, type Column } from "@/components/ui/data-table";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useUserList } from "../hooks/useUserList";
import { UserForm } from "./user-form";
import { UserDetailModal } from "./user-detail-modal";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { useHasPermission } from "@/features/auth/providers/permissions-provider";
import type { User } from "../types";
import { useTranslations } from "next-intl";
import { formatEmailToMailto } from "@/lib/utils";
import type { CreateUserFormData, UpdateUserFormData } from "../schemas/user.schema";
import { useSalesPerformanceList } from "@/features/sales-overview/hooks/useSalesPerformanceList";
import type { SalesPerformanceListItem } from "@/features/sales-overview/types";

export function UserList() {
  const {
    // page, // Unused
    setPage,
    setPerPage,
    search,
    setSearch,
    status,
    setStatus,
    roleId,
    setRoleId,
    isCreateDialogOpen,
    setIsCreateDialogOpen,
    editingUser,
    setEditingUser,
    users,
    pagination,
    roles,
    editingUserData,
    isLoading,
    handleCreate,
    handleUpdate,
    handleDeleteClick,
    handleDeleteConfirm,
    deletingUserId,
    setDeletingUserId,
    deleteUser,
    createUser,
    updateUser,
  } = useUserList();

  const [viewingUserId, setViewingUserId] = useState<string | null>(null);
  const [isDetailModalOpen, setIsDetailModalOpen] = useState(false);
  const t = useTranslations("userManagement.list");
  
  // Fetch sales performance data for achievement calculation
  // Use current month period for achievement
  const now = new Date();
  const firstDayOfMonth = new Date(now.getFullYear(), now.getMonth(), 1);
  const lastDayOfMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0);
  const startDate = firstDayOfMonth.toISOString().split("T")[0];
  const endDate = lastDayOfMonth.toISOString().split("T")[0];
  
  // Fetch sales performance data for achievement calculation
  // Use current month period for achievement
  const salesPerformanceHook = useSalesPerformanceList();
  
  // Set date range for current month (only once on mount)
  useEffect(() => {
    salesPerformanceHook.setStartDate(startDate);
    salesPerformanceHook.setEndDate(endDate);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []); // Only run once on mount
  
  const { performanceList } = salesPerformanceHook;
  
  // Create a map of user_id to performance data for quick lookup
  const performanceMap = new Map<string, SalesPerformanceListItem>(
    performanceList.map((perf: SalesPerformanceListItem) => [perf.user_id, perf])
  );
  
  // Permission checks
  const hasCreatePermission = useHasPermission("users.create");
  const hasEditPermission = useHasPermission("users.edit");
  const hasDeletePermission = useHasPermission("users.delete");



  const handleViewUser = (userId: string) => {
    setViewingUserId(userId);
    setIsDetailModalOpen(true);
  };

  const columns: Column<User>[] = useMemo(() => [
    {
      id: "name",
      header: t("name"),
      accessor: (row) => <NameCell row={row} onClick={() => handleViewUser(row.id)} />,
      className: "w-[200px]",
    },
    {
      id: "email",
      header: t("email"),
      accessor: (row) => (
        <a href={formatEmailToMailto(row.email)} className="text-muted-foreground hover:text-primary hover:underline cursor-pointer min-w-0">
          {row.email}
        </a>
      ),
    },
    {
      id: "role",
      header: t("role"),
      accessor: (row) => (
        <Badge variant="outline" className="font-normal">
          {row.role?.name || "N/A"}
        </Badge>
      ),
    },
    {
      id: "group",
      header: t("group"),
      accessor: (row) => (
        <Badge variant="outline" className="font-normal">
          {row.group?.name || "N/A"}
        </Badge>
      ),
    },
    {
      id: "brick",
      header: t("brick"),
      accessor: (row) => (
        row.brick ? (
          <Link href={`/master-data/bricks/${row.brick.id}/dashboard`} className="inline-block">
            <Badge variant="outline" className="font-normal cursor-pointer hover:bg-accent">
              {row.brick.name || row.brick.code || "N/A"}
            </Badge>
          </Link>
        ) : (
          <span className="text-muted-foreground">N/A</span>
        )
      ),
    },
    {
      id: "monthly_target_achievement",
      header: t("monthlyTargetAchievement"),
      accessor: (row) => <AchievementCell row={row} performanceMap={performanceMap} />,
      className: "w-[200px]",
    },
    {
      id: "status",
      header: t("status"),
      accessor: (row) => (
        <StatusSwitch
          checked={row.status === "active"}
          onCheckedChange={(checked) => {
            updateUser.mutate({
              id: row.id,
              data: { status: checked ? "active" : "inactive" },
            });
          }}
        />
      ),
      className: "w-[100px]",
    },
    {
      id: "actions",
      header: t("actions"),
      accessor: (row) => (
        <div className="flex items-center justify-end gap-1">
          <Button
            variant="ghost"
            size="icon-sm"
            className="h-8 w-8"
            title="View Details"
            onClick={() => handleViewUser(row.id)}
          >
            <Eye className="h-3.5 w-3.5" />
          </Button>
          {hasEditPermission && (
            <Button
              variant="ghost"
              size="icon-sm"
              onClick={() => setEditingUser(row.id)}
              className="h-8 w-8"
              title="Edit"
            >
              <Edit className="h-3.5 w-3.5" />
            </Button>
          )}
          {hasDeletePermission && (
            <Button
              variant="ghost"
              size="icon-sm"
              onClick={() => handleDeleteClick(row.id)}
              className="h-8 w-8 text-destructive hover:text-destructive"
              title="Delete"
            >
              <Trash2 className="h-3.5 w-3.5" />
            </Button>
          )}
        </div>
      ),
      className: "w-[140px] text-right",
    },
  ], [
    t, 
    viewingUserId,
    performanceMap,
    updateUser,
    handleDeleteClick,
    hasEditPermission,
    hasDeletePermission
  ]);

  return (
    <div className="space-y-4">
      {/* Header with Actions */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3 flex-1">
          {/* ... filters ... */}
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
              <SelectItem value="active">Active</SelectItem>
              <SelectItem value="inactive">Inactive</SelectItem>
            </SelectContent>
          </Select>
          <Select 
            value={roleId || "all"} 
            onValueChange={(value) => setRoleId(value === "all" ? "" : value)}
          >
            <SelectTrigger className="w-[140px] h-9">
              <SelectValue placeholder={t("allRoles")} />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">{t("allRoles")}</SelectItem>
            {roles.map((role) => (
                <SelectItem key={role.id} value={role.id}>
                {role.name}
                </SelectItem>
            ))}
            </SelectContent>
          </Select>
        </div>
        <div className="flex items-center gap-2">
            <Link href="/master-data/monthly-targets">
                <Button variant="outline" size="sm">
                    <Target className="h-4 w-4 mr-2" />
                    Manage Targets
                </Button>
            </Link>
            {hasCreatePermission && (
            <Button onClick={() => setIsCreateDialogOpen(true)} size="sm">
                <Plus className="h-4 w-4 mr-2" />
                {t("addUser")}
            </Button>
            )}
        </div>
      </div>

      {/* Table */}
      <DataTable
        columns={columns}
        data={users}
        isLoading={isLoading}
        emptyMessage={t("empty")}
        pagination={
          pagination
            ? {
                page: pagination.page,
                per_page: pagination.per_page,
                total: pagination.total,
                total_pages: pagination.total_pages,
                has_next: pagination.has_next,
                has_prev: pagination.has_prev,
              }
            : undefined
        }
        onPageChange={setPage}
        onPerPageChange={setPerPage}
        itemName="user"
        perPageOptions={[10, 20, 50, 100]}
        onResetFilters={() => {
          setSearch("");
          setStatus("");
          setRoleId("");
          setPage(1);
        }}
      />

      {/* Create Dialog */}
      <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle>Create User</DialogTitle>
          </DialogHeader>
          <UserForm
            onSubmit={async (data) => {
              await handleCreate(data as CreateUserFormData);
            }}
            onCancel={() => setIsCreateDialogOpen(false)}
            isLoading={createUser.isPending}
          />
        </DialogContent>
      </Dialog>

      {/* Edit Dialog */}
      {editingUser && editingUserData?.data && (
        <Dialog open={!!editingUser} onOpenChange={(open) => !open && setEditingUser(null)}>
          <DialogContent className="sm:max-w-[500px]">
            <DialogHeader>
              <DialogTitle>Edit User</DialogTitle>
            </DialogHeader>
            <UserForm
              user={editingUserData.data}
              onSubmit={async (data) => {
                await handleUpdate(data as UpdateUserFormData);
              }}
              onCancel={() => setEditingUser(null)}
              isLoading={updateUser.isPending}
            />
          </DialogContent>
        </Dialog>
      )}

      {/* User Detail Modal */}
      <UserDetailModal
        userId={viewingUserId}
        open={isDetailModalOpen}
        onOpenChange={(open) => {
          setIsDetailModalOpen(open);
          if (!open) {
            setViewingUserId(null);
          }
        }}
        onUserUpdated={() => {
          // Refresh will be handled by query invalidation in hooks
        }}
      />

      {/* Delete Dialog */}
      <DeleteDialog
        open={!!deletingUserId}
        onOpenChange={(open) => {
          if (!open) {
            setDeletingUserId(null);
          }
        }}
        onConfirm={handleDeleteConfirm}
        title={t("deleteTitle")}
        description={
          deletingUserId
            ? t("deleteDescriptionWithName", {
                name:
                  users.find((u) => u.id === deletingUserId)?.name ??
                  t("deleteDescription"),
              })
            : t("deleteDescription")
        }
        itemName="user"
        isLoading={deleteUser.isPending}
      />
    </div>
  );
}

// Extracted Components to fix lint warnings
const NameCell = ({ row, onClick }: { row: User; onClick: () => void }) => {
  const getAvatarUrl = (user: User) => {
    if (user.avatar_url) return user.avatar_url;
    return `https://api.dicebear.com/7.x/lorelei/svg?seed=${encodeURIComponent(user.email)}`;
  };

  return (
    <button
      onClick={onClick}
      className="flex items-center gap-3 font-medium text-primary hover:underline"
    >
      <Avatar className="h-8 w-8">
        <AvatarImage src={getAvatarUrl(row)} alt={row.name} />
      </Avatar>
      <span>{row.name}</span>
    </button>
  );
};

const AchievementCell = ({ row, performanceMap }: { row: User; performanceMap: Map<string, SalesPerformanceListItem> }) => {
  const performance = performanceMap.get(row.id);
  
  // Use performance data as primary source
  const targetFormatted = performance?.target_amount_formatted ?? "-";
  const formattedRevenue = performance?.total_revenue_formatted ?? "-";
  const percentage = performance?.target_achievement_percentage ?? null;
  
  // If no data at all
  if (targetFormatted === "-" && formattedRevenue === "-") {
    return <span className="text-muted-foreground">-</span>;
  }
  
  return (
    <div className="flex flex-col gap-0.5">
      <div className="text-sm">
        <span className="text-muted-foreground">{targetFormatted}</span>
        <span className="mx-1 text-muted-foreground">/</span>
        <span className="font-medium">{formattedRevenue}</span>
      </div>
      {percentage !== null && (
        <p className={`text-[10px] font-medium ${
          percentage >= 100 ? "text-green-600" : 
          percentage >= 75 ? "text-yellow-600" : 
          "text-red-600"
        }`}>
          {Math.round(percentage)}% dari target
        </p>
      )}
    </div>
  );
};
