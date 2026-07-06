import apiClient from "@/lib/api-client";
import type {
  ListUsersResponse,
  UserResponse,
  ListRolesResponse,
  Role,
  MobilePermissionsResponse,
  ListPermissionsResponse,
  Permission,
  UserPermissionsApiResponse,
  Menu,
} from "../types";
import type {
  CreateUserFormData,
  UpdateUserFormData,
} from "../schemas/user.schema";

export const userService = {
  async list(params?: {
    page?: number;
    per_page?: number;
    search?: string;
    status?: string;
    role_id?: string;
    group_id?: string;
  }): Promise<ListUsersResponse> {
    const response = await apiClient.get<ListUsersResponse>("/users", {
      params,
    });
    return response.data;
  },

  async getById(id: string): Promise<UserResponse> {
    const response = await apiClient.get<UserResponse>(`/users/${id}`);
    return response.data;
  },

  async create(data: CreateUserFormData): Promise<UserResponse> {
    const response = await apiClient.post<UserResponse>("/users", data);
    return response.data;
  },

  async update(id: string, data: UpdateUserFormData): Promise<UserResponse> {
    // Clean up the data - remove empty strings and undefined values
    const cleanData: Partial<UpdateUserFormData> = {};

    if (data.email && data.email.trim() !== "") {
      cleanData.email = data.email;
    }
    if (data.name && data.name.trim() !== "") {
      cleanData.name = data.name;
    }
    if (data.role_id && data.role_id.trim() !== "") {
      cleanData.role_id = data.role_id;
    }
    if ("group_id" in data) {
      cleanData.group_id = data.group_id ?? null;
    }
    if ("brick_id" in data) {
      cleanData.brick_id = data.brick_id ?? null;
    }
    if (data.status) {
      cleanData.status = data.status;
    }

    const response = await apiClient.patch<UserResponse>(
      `/users/${id}`,
      cleanData,
    );

    return response.data;
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/users/${id}`);
  },

  async getPermissions(userId: string): Promise<UserPermissionsApiResponse> {
    const response = await apiClient.get<UserPermissionsApiResponse>(
      `/users/${userId}/permissions`,
    );
    return response.data;
  },
};

export const roleService = {
  async list(): Promise<ListRolesResponse> {
    const response = await apiClient.get<ListRolesResponse>("/roles");
    return response.data;
  },

  async getById(id: string): Promise<Role> {
    const response = await apiClient.get<{ success: boolean; data: Role }>(
      `/roles/${id}`,
    );
    return response.data.data;
  },

  async create(data: {
    name: string;
    code: string;
    description?: string;
    status?: string;
    mobile_access?: boolean;
  }): Promise<Role> {
    const response = await apiClient.post<{ success: boolean; data: Role }>(
      "/roles",
      data,
    );
    return response.data.data;
  },

  async update(
    id: string,
    data: {
      name?: string;
      code?: string;
      description?: string;
      status?: string;
      mobile_access?: boolean;
    },
  ): Promise<Role> {
    const response = await apiClient.put<{ success: boolean; data: Role }>(
      `/roles/${id}`,
      data,
    );
    return response.data.data;
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/roles/${id}`);
  },

  async assignPermissions(
    roleId: string,
    permissionIds: string[],
  ): Promise<Role> {
    const response = await apiClient.put<{ success: boolean; data: Role }>(
      `/roles/${roleId}/permissions`,
      { permission_ids: permissionIds },
    );
    return response.data.data;
  },

  async getRolePermissions(roleId: string): Promise<Permission[]> {
    const response = await apiClient.get<{
      success: boolean;
      data: Permission[];
    }>(`/roles/${roleId}/permissions`);
    return response.data.data;
  },

  async getMobilePermissions(
    roleId: string,
  ): Promise<MobilePermissionsResponse> {
    const response = await apiClient.get<{
      success: boolean;
      data: MobilePermissionsResponse;
    }>(`/roles/${roleId}/mobile-permissions`);
    return response.data.data;
  },

  async updateMobilePermissions(
    roleId: string,
    permissions: MobilePermissionsResponse,
  ): Promise<MobilePermissionsResponse> {
    const response = await apiClient.put<{
      success: boolean;
      data: MobilePermissionsResponse;
    }>(`/roles/${roleId}/mobile-permissions`, permissions);
    return response.data.data;
  },
};

export const permissionService = {
  async list(): Promise<ListPermissionsResponse> {
    const response =
      await apiClient.get<ListPermissionsResponse>("/permissions");
    return response.data;
  },

  async getById(id: string): Promise<Permission> {
    const response = await apiClient.get<{
      success: boolean;
      data: Permission;
    }>(`/permissions/${id}`);
    return response.data.data;
  },
};

export const menuService = {
  async list(): Promise<{
    success: boolean;
    data: Menu[];
    timestamp: string;
    request_id: string;
  }> {
    const response = await apiClient.get<{
      success: boolean;
      data: Menu[];
      timestamp: string;
      request_id: string;
    }>("/menus");
    return response.data;
  },

  async getById(id: string): Promise<Menu> {
    const response = await apiClient.get<{ success: boolean; data: Menu }>(
      `/menus/${id}`,
    );
    return response.data.data;
  },
};
