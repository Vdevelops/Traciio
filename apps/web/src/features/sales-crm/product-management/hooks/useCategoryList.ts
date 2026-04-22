"use client";

import { useState } from "react";
import { toast } from "sonner";
import {
  useProductCategories,
  useDeleteProductCategory,
  useProductCategory,
  useCreateProductCategory,
  useUpdateProductCategory,
} from "./useProducts";
import type { CreateCategoryFormData, UpdateCategoryFormData } from "../schemas/category.schema";

export function useCategoryList() {
  const [editingCategory, setEditingCategory] = useState<string | null>(null);
  const [deletingCategoryId, setDeletingCategoryId] = useState<string | null>(null);
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);

  const { data, isLoading } = useProductCategories();
  const { data: editingCategoryData } = useProductCategory(editingCategory || "");
  const deleteCategory = useDeleteProductCategory();
  const createCategory = useCreateProductCategory();
  const updateCategory = useUpdateProductCategory();

  const categories = data?.data || [];
  const categoryForEdit = editingCategoryData?.data;

  const handleCreate = async (formData: CreateCategoryFormData) => {
    try {
      await createCategory.mutateAsync(formData);
      setIsCreateDialogOpen(false);
      toast.success("Category created successfully");
    } catch {
      // Error already handled in api-client interceptor
    }
  };

  const handleUpdate = async (formData: UpdateCategoryFormData) => {
    if (editingCategory) {
      try {
        await updateCategory.mutateAsync({ id: editingCategory, data: formData });
        setEditingCategory(null);
        toast.success("Category updated successfully");
      } catch {
        // Error already handled in api-client interceptor
      }
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await deleteCategory.mutateAsync(id);
      toast.success("Category deleted successfully");
    } catch {
      // Error already handled in api-client interceptor
    }
  };

  return {
    // State
    editingCategory,
    setEditingCategory,
    deletingCategoryId,
    setDeletingCategoryId,
    isCreateDialogOpen,
    setIsCreateDialogOpen,
    // Data
    categories,
    categoryForEdit,
    isLoading,
    // Actions
    handleCreate,
    handleUpdate,
    handleDelete,
    // Mutations
    deleteCategory,
    createCategory,
    updateCategory,
  };
}

