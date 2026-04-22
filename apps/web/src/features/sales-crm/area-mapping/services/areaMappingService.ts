import apiClient from "@/lib/api-client";
import type {
  Territory,
  AreaCapture,
  CoverageAnalysis,
  TerritoriesResponse,
  AreaCapturesResponse,
  HeatmapResponse,
  CreateTerritoryRequest,
  UpdateTerritoryRequest,
  CreateAreaCaptureRequest,
  CreateTerritoryRequest as CreateTerritoryRequestType,
} from "../types";

export const areaMappingService = {
  // Territory APIs
  async listTerritories(params?: {
    page?: number;
    page_size?: number;
    search?: string;
    assigned_to?: string;
  }): Promise<TerritoriesResponse> {
    const response = await apiClient.get<{
      success: boolean;
      data: Territory[];
      meta: {
        pagination: {
          page: number;
          per_page: number;
          total: number;
          total_pages: number;
          has_next: boolean;
          has_prev: boolean;
        };
      };
    }>("/area-mapping/territories", {
      params,
    });
    
    // Transform API response to match TerritoriesResponse interface
    return {
      territories: response.data.data,
      total: response.data.meta.pagination.total,
      page: response.data.meta.pagination.page,
      page_size: response.data.meta.pagination.per_page,
    };
  },

  async getTerritory(id: string): Promise<{ territory: Territory }> {
    const response = await apiClient.get<{
      success: boolean;
      data: Territory;
    }>(`/area-mapping/territories/${id}`);
    return { territory: response.data.data };
  },

  async createTerritory(data: CreateTerritoryRequestType): Promise<{ territory: Territory }> {
    
    // Simple validation
    if (!data.polygon || data.polygon.type !== "Polygon") {
      throw new Error("Invalid GeoJSON Polygon");
    }
    
    if (!data.polygon.coordinates || data.polygon.coordinates.length === 0) {
      throw new Error("Polygon coordinates are required");
    }
    
    const firstRing = data.polygon.coordinates[0];
    if (!firstRing || firstRing.length < 4) {
      throw new Error("Polygon must have at least 4 points (including closing point)");
    }
    
    // Transform GeoJSON to backend format
    // Backend expects coordinates as flat array of [lng, lat] pairs (first ring only)
    const payload = {
      name: data.name,
      description: data.description,
      assigned_to: data.assigned_to,
      color: data.color,
      coordinates: firstRing, // Send first ring as flat coordinate array
    };
    
    
    const response = await apiClient.post<{
      success: boolean;
      data: Territory;
    }>("/area-mapping/territories", payload);
    
    return { territory: response.data.data };
  },

  async updateTerritory(
    id: string,
    data: UpdateTerritoryRequest
  ): Promise<{ territory: Territory }> {
    
    const payload: any = {
      name: data.name,
      description: data.description,
      assigned_to: data.assigned_to,
      color: data.color,
    };
    
    if (data.polygon) {
      // Simple validation
      if (data.polygon.type !== "Polygon") {
        throw new Error("Invalid GeoJSON Polygon");
      }
      
      if (!data.polygon.coordinates || data.polygon.coordinates.length === 0) {
        throw new Error("Polygon coordinates are required");
      }
      
      const firstRing = data.polygon.coordinates[0];
      if (!firstRing || firstRing.length < 4) {
        throw new Error("Polygon must have at least 4 points (including closing point)");
      }
      
      // Transform GeoJSON to backend format
      // Backend expects coordinates as flat array of [lng, lat] pairs (first ring only)
      payload.coordinates = firstRing;
    }
    
    
    const response = await apiClient.put<{
      success: boolean;
      data: Territory;
    }>(`/area-mapping/territories/${id}`, payload);
    
    return { territory: response.data.data };
  },

  async deleteTerritory(id: string): Promise<void> {
    await apiClient.delete(`/area-mapping/territories/${id}`);
  },

  // Area Capture APIs
  async listCaptures(params?: {
    page?: number;
    per_page?: number;
    visit_report_id?: string;
    capture_type?: string;
    captured_after?: string;
    captured_before?: string;
  }): Promise<AreaCapturesResponse> {
    const response = await apiClient.get<{
      success: boolean;
      data: AreaCapture[];
      meta: {
        pagination: {
          page: number;
          per_page: number;
          total: number;
          total_pages: number;
          has_next: boolean;
          has_prev: boolean;
        };
      };
    }>("/area-mapping/captures", {
      params,
    });
    
    return {
      captures: response.data.data,
      total: response.data.meta.pagination.total,
      page: response.data.meta.pagination.page,
      page_size: response.data.meta.pagination.per_page,
    };
  },


  async createCapture(data: CreateAreaCaptureRequest): Promise<{ capture: AreaCapture }> {
    const response = await apiClient.post<{
      success: boolean;
      data: AreaCapture;
    }>("/area-mapping/capture", data);
    return { capture: response.data.data };
  },

  // Coverage Analysis APIs
  async getCoverageAnalysis(params: {
    territory_id: string;
    start_date: string;
    end_date: string;
  }): Promise<{ analysis: CoverageAnalysis }> {
    const response = await apiClient.get<{
      success: boolean;
      data: CoverageAnalysis;
    }>("/area-mapping/coverage", { params });
    return { analysis: response.data.data };
  },

  // Heatmap API
  async getHeatmapData(params?: {
    territory_id?: string;
    capture_type?: string;
    start_date?: string;
    end_date?: string;
  }): Promise<HeatmapResponse> {
    const response = await apiClient.get<{
      success: boolean;
      data: Array<{ lat: number; lng: number; intensity: number }>;
    }>("/area-mapping/heatmap", { params });
    
    // Calculate bounds from data points
    const points = response.data.data;
    let bounds = {
      north: -90,
      south: 90,
      east: -180,
      west: 180,
    };
    
    if (points.length > 0) {
      bounds = points.reduce(
        (acc, point) => ({
          north: Math.max(acc.north, point.lat),
          south: Math.min(acc.south, point.lat),
          east: Math.max(acc.east, point.lng),
          west: Math.min(acc.west, point.lng),
        }),
        bounds
      );
    }
    
    return {
      data: points,
      bounds,
    };
  },
};
