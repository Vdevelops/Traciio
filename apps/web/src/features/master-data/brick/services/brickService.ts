import apiClient from "@/lib/api-client";
import type {
  Brick,
  ListBricksResponse,
  BrickResponse,
  BrickTargetWithDistributions,
} from "../types";
import type {
  CreateBrickFormData,
  UpdateBrickFormData,
  CreateBrickTargetDistributionFormData,
  UpdateBrickTargetDistributionFormData,
} from "../schemas/brick.schema";

export const brickService = {
  async list(params?: {
    page?: number;
    per_page?: number;
    search?: string;
    province?: string;
    regency?: string;
    manager_id?: string;
    status?: string;
  }): Promise<ListBricksResponse> {
    const response = await apiClient.get<ListBricksResponse>(
      "/bricks",
      { params }
    );
    return response.data;
  },

  async getById(id: string): Promise<BrickResponse> {
    const response = await apiClient.get<BrickResponse>(`/bricks/${id}`);
    return response.data;
  },

  async create(data: CreateBrickFormData): Promise<BrickResponse> {
    const response = await apiClient.post<BrickResponse>("/bricks", data);
    return response.data;
  },

  async update(
    id: string,
    data: UpdateBrickFormData
  ): Promise<BrickResponse> {
    // Clean up the data - remove empty strings and undefined values
    const cleanData: Partial<UpdateBrickFormData> = {};

    if (data.name && data.name.trim() !== "") {
      cleanData.name = data.name;
    }
    if (data.code && data.code.trim() !== "") {
      cleanData.code = data.code;
    }
    if (data.description !== undefined) {
      cleanData.description = data.description;
    }
    if (data.province && data.province.trim() !== "") {
      cleanData.province = data.province;
    }
    if (data.regency && data.regency.trim() !== "") {
      cleanData.regency = data.regency;
    }
    if (data.manager_id !== undefined) {
      cleanData.manager_id = data.manager_id;
    }
    if (data.status) {
      cleanData.status = data.status;
    }

    const response = await apiClient.patch<BrickResponse>(
      `/bricks/${id}`,
      cleanData
    );
    return response.data;
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/bricks/${id}`);
  },

  async getSalesInBrick(brickId: string) {
    const response = await apiClient.get<{
      success: boolean;
      data: Array<{
        id: string;
        name: string;
        email: string;
        status: string;
        monthly_target?: {
          target_amount: number;
          target_amount_formatted?: string;
        };
      }>;
    }>(`/bricks/${brickId}/sales`);
    return response.data;
  },

  async getBrickTargetWithDistributions(
    brickId: string,
    year: number,
    month: number
  ): Promise<{ success: boolean; data: BrickTargetWithDistributions }> {
    const response = await apiClient.get<{
      success: boolean;
      data: BrickTargetWithDistributions;
    }>(`/bricks/${brickId}/targets/${year}/${month}`);
    return response.data;
  },

  async distributeTarget(
    brickId: string,
    targetId: string,
    data: CreateBrickTargetDistributionFormData
  ): Promise<{
    success: boolean;
    data: {
      distributed_count: number;
      distributions: Array<{
        id: string;
        sales_user_id: string;
        distributed_amount: number;
      }>;
    };
  }> {
    const response = await apiClient.post<{
      success: boolean;
      data: {
        distributed_count: number;
        distributions: Array<{
          id: string;
          sales_user_id: string;
          distributed_amount: number;
        }>;
      };
    }>(`/bricks/${brickId}/targets/${targetId}/distribute`, data);
    return response.data;
  },

  async updateDistribution(
    brickId: string,
    targetId: string,
    distributionId: string,
    data: UpdateBrickTargetDistributionFormData
  ): Promise<{
    success: boolean;
    data: {
      id: string;
      distributed_amount: number;
    };
  }> {
    const response = await apiClient.patch<{
      success: boolean;
      data: {
        id: string;
        distributed_amount: number;
      };
    }>(`/bricks/${brickId}/targets/${targetId}/distributions/${distributionId}`, data);
    return response.data;
  },

  async deleteDistribution(
    brickId: string,
    targetId: string,
    distributionId: string
  ): Promise<void> {
    await apiClient.delete(
      `/bricks/${brickId}/targets/${targetId}/distributions/${distributionId}`
    );
  },
};

