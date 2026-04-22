"use client";

import { useState } from "react";
import { toast } from "sonner";
import {
  useActivityTypes,
  useDeleteActivityType,
  useActivityType,
  useCreateActivityType,
  useUpdateActivityType,
} from "./useActivityTypes";
import type { CreateActivityTypeFormData, UpdateActivityTypeFormData } from "../schemas/activity-type.schema";

export function useActivityTypeList() {
  const [editingType, setEditingType] = useState<string | null>(null);
  const [deletingTypeId, setDeletingTypeId] = useState<string | null>(null);
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);

  const { data, isLoading } = useActivityTypes();
  const { data: editingTypeData } = useActivityType(editingType || "");
  const deleteType = useDeleteActivityType();
  const createType = useCreateActivityType();
  const updateType = useUpdateActivityType();

  const types = data?.data || [];
  const typeForEdit = editingTypeData?.data;

  const handleCreate = async (formData: CreateActivityTypeFormData | UpdateActivityTypeFormData) => {
    try {
      // Type guard to ensure required fields for create
      if ("name" in formData && "code" in formData && formData.name && formData.code) {
        await createType.mutateAsync({
          name: formData.name,
          code: formData.code,
          description: formData.description,
          icon: formData.icon,
          badge_color: formData.badge_color,
          status: formData.status,
          order: formData.order,
        });
        setIsCreateDialogOpen(false);
        toast.success("Activity type created successfully");
      }
    } catch {
      // Error already handled in api-client interceptor
    }
  };

  const handleUpdate = async (formData: CreateActivityTypeFormData | UpdateActivityTypeFormData) => {
    if (editingType) {
      try {
        await updateType.mutateAsync({ id: editingType, data: formData });
        setEditingType(null);
        toast.success("Activity type updated successfully");
      } catch {
        // Error already handled in api-client interceptor
      }
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await deleteType.mutateAsync(id);
      toast.success("Activity type deleted successfully");
    } catch {
      // Error already handled in api-client interceptor
    }
  };

  return {
    // State
    editingType,
    setEditingType,
    deletingTypeId,
    setDeletingTypeId,
    isCreateDialogOpen,
    setIsCreateDialogOpen,
    // Data
    types,
    typeForEdit,
    isLoading,
    // Handlers
    handleCreate,
    handleUpdate,
    handleDelete,
    // Mutations
    createType,
    updateType,
    deleteType,
  };
}

