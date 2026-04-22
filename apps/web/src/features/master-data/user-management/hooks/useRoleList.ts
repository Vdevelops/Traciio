"use client";

import { useState } from "react";
import { toast } from "sonner";
import {
  useRoles,
  useDeleteRole,
  useCreateRole,
  useUpdateRole,
  useRole,
} from "./useRoles";
import type { CreateRoleFormData, UpdateRoleFormData } from "../schemas/role.schema";

export function useRoleList() {
  const [editingRole, setEditingRole] = useState<string | null>(null);
  const [assigningPermissions, setAssigningPermissions] = useState<string | null>(null);
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [deletingRoleId, setDeletingRoleId] = useState<string | null>(null);

  const { data, isLoading } = useRoles();
  const deleteRole = useDeleteRole();
  const createRole = useCreateRole();
  const updateRole = useUpdateRole();

  const roles = data?.data || [];
  const { data: editingRoleData } = useRole(editingRole || "");
  const roleForEdit = editingRoleData;

  const handleCreate = async (formData: CreateRoleFormData) => {
    try {
      await createRole.mutateAsync(formData);
      setIsCreateDialogOpen(false);
      toast.success("Role created successfully");
    } catch (error) {
      // Error already handled in api-client interceptor
    }
  };

  const handleUpdate = async (formData: UpdateRoleFormData) => {
    if (editingRole) {
      try {
        await updateRole.mutateAsync({ id: editingRole, data: formData });
        setEditingRole(null);
        toast.success("Role updated successfully");
      } catch (error) {
        // Error already handled in api-client interceptor
      }
    }
  };

  const handleDeleteClick = (id: string) => {
    setDeletingRoleId(id);
  };

  const handleDeleteConfirm = async () => {
    if (deletingRoleId) {
      try {
        await deleteRole.mutateAsync(deletingRoleId);
        setDeletingRoleId(null);
        toast.success("Role deleted successfully");
      } catch (error) {
        // Error already handled in api-client interceptor
      }
    }
  };

  return {
    // State
    editingRole,
    setEditingRole,
    assigningPermissions,
    setAssigningPermissions,
    isCreateDialogOpen,
    setIsCreateDialogOpen,
    deletingRoleId,
    setDeletingRoleId,
    // Data
    roles,
    roleForEdit,
    isLoading,
    // Actions
    handleCreate,
    handleUpdate,
    handleDeleteClick,
    handleDeleteConfirm,
    // Mutations
    deleteRole,
    createRole,
    updateRole,
  };
}
