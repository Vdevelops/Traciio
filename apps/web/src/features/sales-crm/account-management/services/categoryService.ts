import { apiClient } from "@/lib/api-client";
import type { Category, ListCategoriesResponse, CategoryResponse } from "../types";
import type { CreateCategoryFormData, UpdateCategoryFormData } from "../schemas/category.schema";

export const categoryService = {
  async list(): Promise<ListCategoriesResponse> {
    const response = await apiClient.get<ListCategoriesResponse>("/categories");
    return response.data;
  },

  async getById(id: string): Promise<Category> {
    const response = await apiClient.get<CategoryResponse>(`/categories/${id}`);
    return response.data.data;
  },

  async create(data: CreateCategoryFormData): Promise<Category> {
    const response = await apiClient.post<CategoryResponse>("/categories", data);
    return response.data.data;
  },

  async update(id: string, data: UpdateCategoryFormData): Promise<Category> {
    const response = await apiClient.put<CategoryResponse>(`/categories/${id}`, data);
    return response.data.data;
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/categories/${id}`);
  },
};

