import { apiClient } from "@/lib/api-client";
import type {
  CustomerPurchaseHistory,
  CustomerProductAnalytics,
  CustomerPurchaseSummary,
} from "../types/purchase-history";
import type { ApiResponse, PaginatedResponse } from "@/types/api";

export const customerPurchaseService = {
  // Get purchase history for account
  async getAccountPurchaseHistory(
    accountId: string,
    page = 1,
    perPage = 20,
  ): Promise<PaginatedResponse<CustomerPurchaseHistory>> {
    const response = await apiClient.get<
      ApiResponse<PaginatedResponse<CustomerPurchaseHistory>>
    >(`/accounts/${accountId}/purchases`, {
      params: { page, per_page: perPage },
    });
    return response.data.data;
  },

  // Get product analytics for account
  async getAccountProductAnalytics(
    accountId: string,
  ): Promise<CustomerProductAnalytics[]> {
    const response = await apiClient.get<
      ApiResponse<CustomerProductAnalytics[]>
    >(`/accounts/${accountId}/purchases/analytics`);
    return response.data.data;
  },

  // Get purchase summary for account
  async getAccountPurchaseSummary(
    accountId: string,
  ): Promise<CustomerPurchaseSummary> {
    const response = await apiClient.get<ApiResponse<CustomerPurchaseSummary>>(
      `/accounts/${accountId}/purchases/summary`,
    );
    return response.data.data;
  },

  // Get purchase history by deal
  async getPurchaseHistoryByDeal(
    dealId: string,
  ): Promise<CustomerPurchaseHistory> {
    const response = await apiClient.get<ApiResponse<CustomerPurchaseHistory>>(
      `/pipeline/deals/${dealId}/purchase-history`,
    );
    return response.data.data;
  },
};
