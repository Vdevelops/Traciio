import apiClient from "@/lib/api-client";
import type {
  MonthlyTarget,
  ListMonthlyTargetsResponse,
  MonthlyTargetResponse,
} from "../types";
import type {
  CreateMonthlyTargetFormData,
  BulkCreateMonthlyTargetFormData,
  UpdateMonthlyTargetFormData,
  CreateGroupTargetWithUserAssignmentFormData,
  BulkSetTargetFormData,
} from "../schemas/monthly-target.schema";

export const monthlyTargetService = {
  async list(params?: {
    page?: number;
    per_page?: number;
    group_id?: string;
    user_id?: string;
    brick_id?: string;
    year?: number;
    month?: number;
    search?: string;
    manager_id?: string;
    scope?: "all" | "user" | "group" | "brick";
  }): Promise<ListMonthlyTargetsResponse> {
    const response = await apiClient.get<ListMonthlyTargetsResponse>(
      "/monthly-targets",
      { params }
    );
    return response.data;
  },

  async getById(id: string): Promise<MonthlyTargetResponse> {
    const response = await apiClient.get<MonthlyTargetResponse>(
      `/monthly-targets/${id}`
    );
    return response.data;
  },

  async getUserEffectiveTarget(params: {
    user_id: string;
    year: number;
    month: number;
  }): Promise<MonthlyTargetResponse> {
    const response = await apiClient.get<MonthlyTargetResponse>(
      "/monthly-targets/user-effective",
      { params }
    );
    return response.data;
  },

  async create(
    data: CreateMonthlyTargetFormData
  ): Promise<MonthlyTargetResponse> {
    const response = await apiClient.post<MonthlyTargetResponse>(
      "/monthly-targets",
      data
    );
    return response.data;
  },

  async update(
    id: string,
    data: UpdateMonthlyTargetFormData
  ): Promise<MonthlyTargetResponse> {
    // Clean up the data - remove undefined values
    const cleanData: Partial<UpdateMonthlyTargetFormData> = {};

    if (data.year !== undefined) {
      cleanData.year = data.year;
    }
    if (data.month !== undefined) {
      cleanData.month = data.month;
    }
    if (data.target_amount !== undefined) {
      cleanData.target_amount = data.target_amount;
    }

    const response = await apiClient.put<MonthlyTargetResponse>(
      `/monthly-targets/${id}`,
      cleanData
    );
    return response.data;
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/monthly-targets/${id}`);
  },

  async bulkCreate(
    data: BulkCreateMonthlyTargetFormData
  ): Promise<{ success: boolean; data: MonthlyTarget[]; timestamp: string; request_id: string }> {
    const response = await apiClient.post<{
      success: boolean;
      data: MonthlyTarget[];
      timestamp: string;
      request_id: string;
    }>("/monthly-targets/bulk", data);
    return response.data;
  },

  async createGroupTargetWithUserAssignment(
    data: CreateGroupTargetWithUserAssignmentFormData
  ): Promise<{
    success: boolean;
    data: {
      group_target: MonthlyTarget;
      user_targets: MonthlyTarget[];
      total_users: number;
    };
    timestamp: string;
    request_id: string;
  }> {
    const response = await apiClient.post<{
      success: boolean;
      data: {
        group_target: MonthlyTarget;
        user_targets: MonthlyTarget[];
        total_users: number;
      };
      timestamp: string;
      request_id: string;
    }>("/monthly-targets/group-with-users", data);
    return response.data;
  },

  async bulkSetTarget(
    data: BulkSetTargetFormData
  ): Promise<MonthlyTargetResponse[]> {
    const response = await apiClient.post<{
      success: boolean;
      data: MonthlyTargetResponse[];
      timestamp: string;
      request_id: string;
    }>("/monthly-targets/bulk-set", data);
    return response.data.data;
  },
};
