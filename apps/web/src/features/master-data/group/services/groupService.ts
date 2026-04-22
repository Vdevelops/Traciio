import apiClient from "@/lib/api-client";
import type {
  Group,
  ListGroupsResponse,
  GroupResponse,
} from "../types";
import type {
  CreateGroupFormData,
  UpdateGroupFormData,
} from "../schemas/group.schema";

export const groupService = {
  async list(params?: {
    page?: number;
    per_page?: number;
    search?: string;
    status?: string;
  }): Promise<ListGroupsResponse> {
    const response = await apiClient.get<ListGroupsResponse>(
      "/groups",
      { params }
    );
    return response.data;
  },

  async getById(id: string): Promise<GroupResponse> {
    const response = await apiClient.get<GroupResponse>(`/groups/${id}`);
    return response.data;
  },

  async create(data: CreateGroupFormData): Promise<GroupResponse> {
    const response = await apiClient.post<GroupResponse>("/groups", data);
    return response.data;
  },

  async update(
    id: string,
    data: UpdateGroupFormData
  ): Promise<GroupResponse> {
    // Clean up the data - remove empty strings and undefined values
    const cleanData: Partial<UpdateGroupFormData> = {};

    if (data.name && data.name.trim() !== "") {
      cleanData.name = data.name;
    }
    if (data.code && data.code.trim() !== "") {
      cleanData.code = data.code;
    }
    if (data.description !== undefined) {
      cleanData.description = data.description;
    }
    if (data.status) {
      cleanData.status = data.status;
    }

    const response = await apiClient.put<GroupResponse>(
      `/groups/${id}`,
      cleanData
    );
    return response.data;
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/groups/${id}`);
  },
};

