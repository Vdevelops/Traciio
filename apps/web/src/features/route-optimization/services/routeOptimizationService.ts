import apiClient from "@/lib/api-client";
import type {
  OptimizeRouteRequest,
  OptimizedRouteResponse,
  CalculateDistanceRequest,
  CalculateDistanceResponseWrapper,
  ListRoutesResponse,
  OptimizedRoute,
} from "../types";

export const routeOptimizationService = {
  async optimize(data: OptimizeRouteRequest): Promise<OptimizedRouteResponse> {
    const response = await apiClient.post<OptimizedRouteResponse>(
      "/route-optimization/optimize",
      data,
    );
    return response.data;
  },

  async getById(id: string): Promise<OptimizedRoute> {
    const response = await apiClient.get<OptimizedRouteResponse>(
      `/route-optimization/route/${id}`,
    );
    return response.data.data;
  },

  async list(params?: {
    page?: number;
    per_page?: number;
    user_id?: string;
  }): Promise<ListRoutesResponse> {
    const response = await apiClient.get<ListRoutesResponse>(
      "/route-optimization/history",
      { params },
    );
    return response.data;
  },

  async calculateDistance(
    data: CalculateDistanceRequest,
  ): Promise<CalculateDistanceResponseWrapper> {
    const response = await apiClient.post<CalculateDistanceResponseWrapper>(
      "/route-optimization/calculate-distance",
      data,
    );
    return response.data;
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/route-optimization/route/${id}`);
  },
};


