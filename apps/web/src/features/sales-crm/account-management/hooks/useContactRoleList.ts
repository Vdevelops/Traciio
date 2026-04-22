"use client";

import { useState } from "react";
import { toast } from "sonner";
import {
  useContactRoles,
  useDeleteContactRole,
  useCreateContactRole,
  useUpdateContactRole,
  useContactRole,
} from "./useContactRoles";
import type { CreateContactRoleFormData, UpdateContactRoleFormData } from "../schemas/contact-role.schema";

export function useContactRoleList() {
  const [editingContactRole, setEditingContactRole] = useState<string | null>(null);
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [deletingContactRoleId, setDeletingContactRoleId] = useState<string | null>(null);

  const { data, isLoading } = useContactRoles();
  const deleteContactRole = useDeleteContactRole();
  const createContactRole = useCreateContactRole();
  const updateContactRole = useUpdateContactRole();

  const contactRoles = data?.data || [];
  const { data: editingContactRoleData } = useContactRole(editingContactRole || "");
  const contactRoleForEdit = editingContactRoleData;

  const handleCreate = async (formData: CreateContactRoleFormData) => {
    try {
      await createContactRole.mutateAsync(formData);
      setIsCreateDialogOpen(false);
      toast.success("Contact role created successfully");
    } catch (error) {
      // Error already handled in api-client interceptor
    }
  };

  const handleUpdate = async (formData: UpdateContactRoleFormData) => {
    if (editingContactRole) {
      try {
        await updateContactRole.mutateAsync({ id: editingContactRole, data: formData });
        setEditingContactRole(null);
        toast.success("Contact role updated successfully");
      } catch (error) {
        // Error already handled in api-client interceptor
      }
    }
  };

  const handleDeleteClick = (id: string) => {
    setDeletingContactRoleId(id);
  };

  const handleDeleteConfirm = async () => {
    if (deletingContactRoleId) {
      try {
        await deleteContactRole.mutateAsync(deletingContactRoleId);
        toast.success("Contact role deleted successfully");
        setDeletingContactRoleId(null);
      } catch (error) {
        // Error already handled in api-client interceptor
      }
    }
  };

  return {
    // State
    editingContactRole,
    setEditingContactRole,
    isCreateDialogOpen,
    setIsCreateDialogOpen,
    deletingContactRoleId,
    setDeletingContactRoleId,
    contactRoles,
    contactRoleForEdit,
    isLoading,
    // Handlers
    handleCreate,
    handleUpdate,
    handleDeleteClick,
    handleDeleteConfirm,
    // Mutations
    deleteContactRole,
    createContactRole,
    updateContactRole,
  };
}

