import apiClient from "@/lib/api-client";
import type { ListProductsResponse, ProductResponse, ListProductCategoriesResponse, ProductCategory } from "../types";
import type { CreateProductFormData, UpdateProductFormData } from "../schemas/product.schema";

export const productService = {
  async list(params?: {
    page?: number;
    per_page?: number;
    search?: string;
    status?: string;
    category_id?: string;
  }): Promise<ListProductsResponse> {
    const response = await apiClient.get<ListProductsResponse>("/products", { params });
    return response.data;
  },

  async getById(id: string): Promise<ProductResponse> {
    const response = await apiClient.get<ProductResponse>(`/products/${id}`);
    return response.data;
  },

  async create(data: CreateProductFormData): Promise<ProductResponse> {
    const response = await apiClient.post<ProductResponse>("/products", data);
    return response.data;
  },

  async update(id: string, data: UpdateProductFormData): Promise<ProductResponse> {
    const response = await apiClient.put<ProductResponse>(`/products/${id}`, data);
    return response.data;
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/products/${id}`);
  },

  async uploadImage(file: File): Promise<{ url: string }> {
    const formData = new FormData();
    formData.append("image", file);
    
    const response = await apiClient.post<{ success: boolean; data: { url: string } }>(
      "/upload/image",
      formData,
      {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      }
    );
    
    return response.data.data;
  },
};

export const productCategoryService = {
  async list(params?: {
    status?: string;
  }): Promise<ListProductCategoriesResponse> {
    const response = await apiClient.get<ListProductCategoriesResponse>("/product-categories", { params });
    return response.data;
  },

  async getById(id: string): Promise<{ success: boolean; data: ProductCategory }> {
    const response = await apiClient.get<{ success: boolean; data: ProductCategory }>(`/product-categories/${id}`);
    return response.data;
  },

  async create(data: { name: string; slug?: string; description?: string; status?: string }): Promise<{ success: boolean; data: ProductCategory }> {
    const response = await apiClient.post<{ success: boolean; data: ProductCategory }>("/product-categories", data);
    return response.data;
  },

  async update(id: string, data: { name?: string; slug?: string; description?: string; status?: string }): Promise<{ success: boolean; data: ProductCategory }> {
    const response = await apiClient.put<{ success: boolean; data: ProductCategory }>(`/product-categories/${id}`, data);
    return response.data;
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/product-categories/${id}`);
  },
};

