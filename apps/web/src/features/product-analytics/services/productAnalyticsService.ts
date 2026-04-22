import apiClient from "@/lib/api-client";
import type {
  GetProductPerformanceRequest,
  ProductPerformanceResponse,
  GetProductComparisonRequest,
  ProductComparisonResponse,
  GetProductTrendsRequest,
  ProductTrendsResponse,
} from "../types";
import type { GetMonthlySalesResponse } from "../types/monthly-sales";

export const productAnalyticsService = {
  /**
   * Get detailed performance metrics for a specific product
   */
  async getProductPerformance(
    productId: string,
    params?: Omit<GetProductPerformanceRequest, "product_id">
  ): Promise<ProductPerformanceResponse> {
    const response = await apiClient.get<ProductPerformanceResponse>(
      `/product-analytics/product/${productId}/performance`,
      {
        params,
      }
    );
    return response.data;
  },

  /**
   * Compare performance of multiple products
   */
  async getProductComparison(
    params: GetProductComparisonRequest
  ): Promise<ProductComparisonResponse> {
    const response = await apiClient.get<ProductComparisonResponse>(
      "/product-analytics/product-comparison",
      {
        params: {
          ...params,
          product_ids: params.product_ids.join(","),
        },
      }
    );
    return response.data;
  },

  /**
   * Get sales trends for a specific product over time
   */
  async getProductTrends(
    productId: string,
    params?: Omit<GetProductTrendsRequest, "product_id">
  ): Promise<ProductTrendsResponse> {
    const response = await apiClient.get<ProductTrendsResponse>(
      `/product-analytics/product-trends/${productId}`,
      {
        params,
      }
    );
    return response.data;
  },

  /**
   * Get products list with sales analytics and filtering
   */
  async getProductsList(params?: {
    period?: "all" | "day" | "week" | "month" | "year";
    start_date?: string;
    end_date?: string;
    search?: string;
    sort_by?: "total_sold" | "revenue" | "profit" | "name";
    order?: "asc" | "desc";
    limit?: number;
    page?: number;
    per_page?: number;
  }) {
    const response = await apiClient.get(
      "/product-analytics/products-list",
      {
        params,
      }
    );
    return response.data;
  },

  /**
   * Get monthly sales data for a specific period (all products aggregated)
   */
  async getMonthlySales(params?: { start_date?: string; end_date?: string }): Promise<GetMonthlySalesResponse> {
    const response = await apiClient.get<GetMonthlySalesResponse>(
      "/product-analytics/monthly-sales",
      {
        params,
      }
    );
    return response.data;
  },

  /**
   * Get monthly sales data for a specific product and period
   */
  async getProductMonthlySales(productId: string, params?: { start_date?: string; end_date?: string }): Promise<GetMonthlySalesResponse> {
    const response = await apiClient.get<GetMonthlySalesResponse>(
      `/product-analytics/product/${productId}/monthly-sales`,
      {
        params,
      }
    );
    return response.data;
  },

  /**
   * Get products sold by a specific user with sales analytics
   */
  async getUserProductSales(userId: string, params?: {
    period?: "all" | "day" | "week" | "month" | "year";
    start_date?: string; // YYYY-MM-DD format
    end_date?: string; // YYYY-MM-DD format
    sort_by?: "total_sold" | "revenue" | "profit" | "name";
    order?: "asc" | "desc";
    page?: number;
    per_page?: number;
  }) {
    const response = await apiClient.get(
      `/product-analytics/user/${userId}/products`,
      {
        params,
      }
    );
    return response.data;
  },
};
