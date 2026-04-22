"use client";

import { useState } from "react";
import { Edit, Trash2, Plus, Settings, Smartphone } from "lucide-react";
import { useTranslations } from "next-intl";
import { StatusSwitch } from "@/components/ui/status-switch";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { DeleteDialog } from "@/components/ui/delete-dialog";
import { useRoleList } from "../hooks/useRoleList";
import { RoleForm } from "./role-form";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { AssignPermissionsDialog } from "./assign-permissions-dialog";
import { MobilePermissionsDialog } from "./mobile-permissions-dialog";
import type { CreateRoleFormData, UpdateRoleFormData } from "../schemas/role.schema";

export function RoleList() {
  const {
    editingRole,
    setEditingRole,
    assigningPermissions,
    setAssigningPermissions,
    isCreateDialogOpen,
    setIsCreateDialogOpen,
    deletingRoleId,
    setDeletingRoleId,
    roles,
    roleForEdit,
    isLoading,
    handleCreate,
    handleUpdate,
    handleDeleteClick,
    handleDeleteConfirm,
    deleteRole,
    createRole,
    updateRole,
  } = useRoleList();
  
  const [configuringMobilePermissions, setConfiguringMobilePermissions] = useState<string | null>(null);
  const t = useTranslations("userManagement.roleList");

  return (
    <div className="space-y-4">
      {/* Header with Actions */}
      <div className="flex items-center justify-end">
        <Button onClick={() => setIsCreateDialogOpen(true)} size="sm">
          <Plus className="h-4 w-4 mr-2" />
          {t("addRole")}
        </Button>
      </div>

      {/* Table */}
      <div className="border rounded-lg">
        {isLoading ? (
          <div className="p-4 space-y-3">
            {Array.from({ length: 5 }).map((_, i) => (
              <Skeleton key={`skeleton-${i}`} className="h-10 w-full" />
            ))}
          </div>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-[200px]">{t("name")}</TableHead>
                <TableHead>{t("code")}</TableHead>
                <TableHead>{t("description")}</TableHead>
                <TableHead className="w-[100px]">{t("status")}</TableHead>
                <TableHead className="w-[120px]">{t("permissions")}</TableHead>
                <TableHead className="w-[100px] text-center">{t("mobileAccess")}</TableHead>
                <TableHead className="w-[120px] text-right">
                  {t("actions")}
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {roles.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} className="text-center text-muted-foreground py-8">
                      {t("empty")}
                  </TableCell>
                </TableRow>
              ) : (
                roles.map((role) => (
                  <TableRow key={role.id} className="hover:bg-muted/50">
                    <TableCell className="font-medium">{role.name}</TableCell>
                    <TableCell>
                      <code className="text-xs bg-muted px-1.5 py-0.5 rounded">
                        {role.code}
                      </code>
                    </TableCell>
                    <TableCell className="text-muted-foreground max-w-xs truncate">
                      {role.description || "-"}
                    </TableCell>
                    <TableCell>
                      <StatusSwitch
                        checked={role.status === "active"}
                        onCheckedChange={(checked) => {
                          updateRole.mutate({
                            id: role.id,
                            data: { status: checked ? "active" : "inactive" },
                          });
                        }}
                        disabled={role.is_protected}
                      />
                    </TableCell>
                    <TableCell>
                      <Badge variant="outline" className="font-normal">
                        {role.permissions?.length || 0}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-center">
                      {role.mobile_access ? (
                        <button
                          onClick={() => setConfiguringMobilePermissions(role.id)}
                          className="mx-auto p-2 rounded-md hover:bg-primary/10 active:bg-primary/20 transition-all cursor-pointer group border border-transparent hover:border-primary/20"
                          aria-label={t("configureMobilePermissions") || "Configure Mobile Permissions"}
                          title={t("configureMobilePermissions") || "Click to configure mobile permissions"}
                        >
                          <Smartphone className="h-4 w-4 text-primary group-hover:scale-110 group-hover:text-primary/80 transition-all" />
                        </button>
                      ) : (
                        <span className="text-muted-foreground">-</span>
                      )}
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center justify-end gap-1">
                        <Button
                          variant="ghost"
                          size="icon-sm"
                          onClick={() => setEditingRole(role.id)}
                          className="h-8 w-8"
                        >
                          <Edit className="h-3.5 w-3.5" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon-sm"
                          onClick={() => setAssigningPermissions(role.id)}
                          className="h-8 w-8"
                          aria-label={t("assignPermissionsTitle")}
                          title="Assign Permissions"
                        >
                          <Settings className="h-3.5 w-3.5" />
                        </Button>
                        {(!role.user_count || role.user_count === 0) && !role.is_protected && (
                          <Button
                            variant="ghost"
                            size="icon-sm"
                            onClick={() => handleDeleteClick(role.id)}
                            className="h-8 w-8 text-destructive hover:text-destructive"
                            title="Delete role"
                          >
                            <Trash2 className="h-3.5 w-3.5" />
                          </Button>
                        )}
                      </div>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        )}
      </div>

      {/* Create Dialog */}
      <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
        <DialogContent className="sm:max-w-[500px]">
          <DialogHeader>
            <DialogTitle>Create Role</DialogTitle>
          </DialogHeader>
          <RoleForm
            onSubmit={async (data) => {
              await handleCreate(data as CreateRoleFormData);
            }}
            onCancel={() => setIsCreateDialogOpen(false)}
            isLoading={createRole.isPending}
          />
        </DialogContent>
      </Dialog>

      {/* Edit Dialog */}
      {editingRole && roleForEdit && (
        <Dialog open={!!editingRole} onOpenChange={(open) => !open && setEditingRole(null)}>
          <DialogContent className="sm:max-w-[500px]">
            <DialogHeader>
              <DialogTitle>Edit Role</DialogTitle>
            </DialogHeader>
            <RoleForm
              role={roleForEdit}
              onSubmit={async (data) => {
                await handleUpdate(data as UpdateRoleFormData);
              }}
              onCancel={() => setEditingRole(null)}
              isLoading={updateRole.isPending}
            />
          </DialogContent>
        </Dialog>
      )}

      {/* Assign Permissions Dialog */}
      {assigningPermissions && (
        <AssignPermissionsDialog
          roleId={assigningPermissions}
          onClose={() => setAssigningPermissions(null)}
        />
      )}

      {/* Mobile Permissions Dialog */}
      {configuringMobilePermissions && (
        <MobilePermissionsDialog
          roleId={configuringMobilePermissions}
          onClose={() => setConfiguringMobilePermissions(null)}
        />
      )}

      {/* Delete Dialog */}
      <DeleteDialog
        open={!!deletingRoleId}
        onOpenChange={(open) => {
          if (!open) {
            setDeletingRoleId(null);
          }
        }}
        onConfirm={handleDeleteConfirm}
        title={t("deleteTitle")}
        description={
          deletingRoleId
            ? t("deleteDescriptionWithName", {
                name:
                  roles.find((r) => r.id === deletingRoleId)?.name ||
                  "this role",
              })
            : t("deleteDescription")
        }
        itemName={t("deleteItemName")}
        isLoading={deleteRole.isPending}
      />
    </div>
  );
}
