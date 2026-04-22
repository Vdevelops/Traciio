/**
 * Brick Analytics Service
 * Service for interacting with brick analytics APIs
 */

import { apiClient } from "@/lib/api-client";
import type {
  BrickPerformanceResponse,
  BrickPerformanceListResponse,
} from "../types/analytics";

export const brickAnalyticsService = {
  /**
   * Get performance metrics for a single brick
   * @param brickId - Brick ID
   * @param periodStart - Start date (YYYY-MM-DD)
   * @param periodEnd - End date (YYYY-MM-DD)
   */
  async getBrickPerformance(
    brickId: string,
    periodStart?: string,
    periodEnd?: string
  ): Promise<BrickPerformanceResponse> {
    const params = new URLSearchParams();
    if (periodStart) params.append("period_start", periodStart);
    if (periodEnd) params.append("period_end", periodEnd);

    const queryString = params.toString();
    const url = `/bricks/${brickId}/performance${queryString ? `?${queryString}` : ""}`;

    const response = await apiClient.get<BrickPerformanceResponse>(url);
    return response.data;
  },

  /**
   * Get performance metrics for multiple bricks
   * @param brickIds - Array of brick IDs
   * @param periodStart - Start date (YYYY-MM-DD)
   * @param periodEnd - End date (YYYY-MM-DD)
   */
  async listBrickPerformance(
    brickIds: string[],
    periodStart?: string,
    periodEnd?: string
  ): Promise<BrickPerformanceListResponse> {
    const params = new URLSearchParams();
    params.append("brick_ids", brickIds.join(","));
    if (periodStart) params.append("period_start", periodStart);
    if (periodEnd) params.append("period_end", periodEnd);

    const url = `/bricks/performance?${params.toString()}`;

    const response = await apiClient.get<BrickPerformanceListResponse>(url);
    return response.data;
  },
};

