"use client";

import { useState } from "react";
import { toast } from "sonner";
import {
  useCategories,
  useDeleteCategory,
  useCreateCategory,
  useUpdateCategory,
  useCategory,
} from "./useCategories";
import type { CreateCategoryFormData, UpdateCategoryFormData } from "../schemas/category.schema";

export function useCategoryList() {
  const [editingCategory, setEditingCategory] = useState<string | null>(null);
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [deletingCategoryId, setDeletingCategoryId] = useState<string | null>(null);

  const { data, isLoading } = useCategories();
  const deleteCategory = useDeleteCategory();
  const createCategory = useCreateCategory();
  const updateCategory = useUpdateCategory();

  const categories = data?.data || [];
  const { data: editingCategoryData } = useCategory(editingCategory || "");
  const categoryForEdit = editingCategoryData;

  const handleCreate = async (formData: CreateCategoryFormData) => {
    try {
      await createCategory.mutateAsync(formData);
      setIsCreateDialogOpen(false);
      toast.success("Category created successfully");
    } catch (error) {
      // Error already handled in api-client interceptor
    }
  };

  const handleUpdate = async (formData: UpdateCategoryFormData) => {
    if (editingCategory) {
      try {
        await updateCategory.mutateAsync({ id: editingCategory, data: formData });
        setEditingCategory(null);
        toast.success("Category updated successfully");
      } catch (error) {
        // Error already handled in api-client interceptor
      }
    }
  };

  const handleDeleteClick = (id: string) => {
    setDeletingCategoryId(id);
  };

  const handleDeleteConfirm = async () => {
    if (deletingCategoryId) {
      try {
        await deleteCategory.mutateAsync(deletingCategoryId);
        toast.success("Category deleted successfully");
        setDeletingCategoryId(null);
      } catch (error) {
        // Error already handled in api-client interceptor
      }
    }
  };

  return {
    // State
    editingCategory,
    setEditingCategory,
    isCreateDialogOpen,
    setIsCreateDialogOpen,
    deletingCategoryId,
    setDeletingCategoryId,
    categories,
    categoryForEdit,
    isLoading,
    // Handlers
    handleCreate,
    handleUpdate,
    handleDeleteClick,
    handleDeleteConfirm,
    // Mutations
    deleteCategory,
    createCategory,
    updateCategory,
  };
}

