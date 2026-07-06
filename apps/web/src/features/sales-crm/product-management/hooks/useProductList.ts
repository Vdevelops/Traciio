"use client";

import { useState, useEffect } from "react";
import { toast } from "sonner";
import { useDebounce } from "@/hooks/use-debounce";
import {
  useDeleteProduct,
  useProduct,
  useProducts,
  useCreateProduct,
  useUpdateProduct,
  useProductCategories,
} from "./useProducts";
import type { CreateProductFormData, UpdateProductFormData } from "../schemas/product.schema";

export function useProductList() {
  const [page, setPage] = useState(1);
  const [perPage, setPerPage] = useState(20);
  const [search, setSearch] = useState("");
  const debouncedSearch = useDebounce(search, 500);

  useEffect(() => {
    setPage(1);
  }, [debouncedSearch]);
  const [status, setStatus] = useState<string>("");
  const [categoryId, setCategoryId] = useState<string>("");
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [editingProduct, setEditingProduct] = useState<string | null>(null);
  const [deletingProductId, setDeletingProductId] = useState<string | null>(null);

  // Use products hook from useProducts.ts
  const { 
    data: productsData,
    isLoading: isLoadingProducts,
    isError,
  } = useProducts({
    page,
    per_page: perPage,
    search: debouncedSearch,
    status: status === "all" ? undefined : status,
    category_id: categoryId === "all" ? undefined : categoryId,
  });

  const { data: categoriesData } = useProductCategories();
  const { data: editingProductData } = useProduct(editingProduct || "");
  const deleteProduct = useDeleteProduct();
  const createProduct = useCreateProduct();
  const updateProduct = useUpdateProduct();

  // Data extraction
  const products = productsData?.data || [];
  const pagination = productsData?.meta?.pagination;
  const categories = categoriesData?.data || [];
  const isLoading = isLoadingProducts;

  const handleCreate = async (formData: CreateProductFormData) => {
    try {
      await createProduct.mutateAsync(formData);
      setIsCreateDialogOpen(false);
      toast.success("Product created successfully");
    } catch {
      // Error already handled in api-client interceptor
    }
  };

  const handleUpdate = async (formData: UpdateProductFormData) => {
    if (editingProduct) {
      try {
        await updateProduct.mutateAsync({ id: editingProduct, data: formData });
        setEditingProduct(null);
        toast.success("Product updated successfully");
      } catch {
        // Error already handled in api-client interceptor
      }
    }
  };

  const handleDeleteClick = (id: string) => {
    setDeletingProductId(id);
  };

  const handleDeleteConfirm = async () => {
    if (deletingProductId) {
      try {
        await deleteProduct.mutateAsync(deletingProductId);
        toast.success("Product deleted successfully");
        setDeletingProductId(null);
      } catch {
        // Error already handled in api-client interceptor
      }
    }
  };

  const handlePerPageChange = (newPerPage: number) => {
    setPerPage(newPerPage);
    setPage(1); // Reset to first page when changing per page
  };

  return {
    // State
    page,
    setPage,
    perPage,
    setPerPage: handlePerPageChange,
    search,
    setSearch,
    status,
    setStatus,
    categoryId,
    setCategoryId,
    isCreateDialogOpen,
    setIsCreateDialogOpen,
    editingProduct,
    setEditingProduct,
    deletingProductId,
    setDeletingProductId,
    // Data
    products,
    pagination,
    categories,
    editingProductData,
    isLoading,
    isError,
    // Actions
    handleCreate,
    handleUpdate,
    handleDeleteClick,
    handleDeleteConfirm,
    // Mutations
    deleteProduct,
    createProduct,
    updateProduct,
  };
}
