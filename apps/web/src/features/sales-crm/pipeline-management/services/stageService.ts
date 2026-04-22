import apiClient from "@/lib/api-client";
import type { PipelineStage, StageWithStats } from "../types";
import type { StageFormData, StageUpdateData, StageReorderData } from "../schemas/stage.schema";

interface ApiResponse<T> {
  success: boolean;
  data: T;
  error?: {
    code: string;
    message: string;
    details?: unknown;
  };
}

export async function getStages(): Promise<ApiResponse<PipelineStage[]>> {
  const response = await apiClient.get<ApiResponse<PipelineStage[]>>(`/pipelines`);
  return response.data;
}

export async function getStagesWithStats(): Promise<ApiResponse<StageWithStats[]>> {
  // Note: Backend doesn't have a dedicated stats endpoint yet
  // For now, we'll use the regular stages endpoint
  // TODO: Add dedicated stats endpoint in backend or use summary endpoint
  const response = await apiClient.get<ApiResponse<StageWithStats[]>>(`/pipelines`);
  return response.data;
}

export async function getStage(id: string): Promise<ApiResponse<PipelineStage>> {
  const response = await apiClient.get<ApiResponse<PipelineStage>>(`/pipelines/${id}`);
  return response.data;
}

export async function createStage(data: StageFormData): Promise<ApiResponse<PipelineStage>> {
  const response = await apiClient.post<ApiResponse<PipelineStage>>(`/pipelines`, data);
  return response.data;
}

export async function updateStage(id: string, data: StageUpdateData): Promise<ApiResponse<PipelineStage>> {
  const response = await apiClient.put<ApiResponse<PipelineStage>>(`/pipelines/${id}`, data);
  return response.data;
}

export async function reorderStages(data: StageReorderData[]): Promise<ApiResponse<{ message: string }>> {
  const response = await apiClient.put<ApiResponse<{ message: string }>>(`/pipelines/order`, { stages: data });
  return response.data;
}

export async function deleteStage(id: string): Promise<ApiResponse<{ message: string }>> {
  const response = await apiClient.delete<ApiResponse<{ message: string }>>(`/pipelines/${id}`);
  return response.data;
}
