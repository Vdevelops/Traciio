"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { productService, productCategoryService } from "../services/productService";
import type { CreateProductFormData, UpdateProductFormData } from "../schemas/product.schema";

export function useProducts(params?: {
  page?: number;
  per_page?: number;
  search?: string;
  status?: string;
  category_id?: string;
}) {
  return useQuery({
    queryKey: ["products", params],
    queryFn: () => productService.list(params),
    retry: (failureCount, error) => {
      // Don't retry on 404 errors (might be expected)
      if (error && typeof error === "object" && "response" in error) {
        const axiosError = error as { response?: { status?: number } };
        if (axiosError.response?.status === 404) {
          return false;
        }
      }
      return failureCount < 1;
    },
  });
}

export function useProduct(id: string) {
  return useQuery({
    queryKey: ["products", id],
    queryFn: () => productService.getById(id),
    enabled: !!id,
  });
}

export function useCreateProduct() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateProductFormData) => productService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["products"] });
    },
  });
}

export function useUpdateProduct() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateProductFormData }) =>
      productService.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["products"] });
      queryClient.invalidateQueries({ queryKey: ["products", variables.id] });
    },
  });
}

export function useDeleteProduct() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => productService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["products"] });
    },
  });
}

export function useProductCategories(params?: {
  status?: string;
}) {
  return useQuery({
    queryKey: ["product-categories", params],
    queryFn: () => productCategoryService.list(params),
  });
}

export function useProductCategory(id: string) {
  return useQuery({
    queryKey: ["product-categories", id],
    queryFn: () => productCategoryService.getById(id),
    enabled: !!id,
  });
}

export function useCreateProductCategory() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: { name: string; slug?: string; description?: string; status?: string }) =>
      productCategoryService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["product-categories"] });
    },
  });
}

export function useUpdateProductCategory() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: { name?: string; slug?: string; description?: string; status?: string } }) =>
      productCategoryService.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["product-categories"] });
      queryClient.invalidateQueries({ queryKey: ["product-categories", variables.id] });
    },
  });
}

export function useDeleteProductCategory() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => productCategoryService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["product-categories"] });
      queryClient.invalidateQueries({ queryKey: ["products"] }); // Invalidate products too since they depend on categories
    },
  });
}

