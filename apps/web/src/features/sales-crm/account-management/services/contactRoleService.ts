import { apiClient } from "@/lib/api-client";
import type { ContactRole, ListContactRolesResponse, ContactRoleResponse } from "../types";
import type {
  CreateContactRoleFormData,
  UpdateContactRoleFormData,
} from "../schemas/contact-role.schema";

export const contactRoleService = {
  async list(): Promise<ListContactRolesResponse> {
    const response = await apiClient.get<ListContactRolesResponse>("/contact-roles");
    return response.data;
  },

  async getById(id: string): Promise<ContactRole> {
    const response = await apiClient.get<ContactRoleResponse>(`/contact-roles/${id}`);
    return response.data.data;
  },

  async create(data: CreateContactRoleFormData): Promise<ContactRole> {
    const response = await apiClient.post<ContactRoleResponse>("/contact-roles", data);
    return response.data.data;
  },

  async update(id: string, data: UpdateContactRoleFormData): Promise<ContactRole> {
    const response = await apiClient.put<ContactRoleResponse>(`/contact-roles/${id}`, data);
    return response.data.data;
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/contact-roles/${id}`);
  },
};

