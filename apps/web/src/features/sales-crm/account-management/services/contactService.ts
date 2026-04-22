import apiClient from "@/lib/api-client";
import type { Contact, ListContactsResponse, ContactResponse } from "../types";
import type { CreateContactFormData, UpdateContactFormData } from "../schemas/contact.schema";

export const contactService = {
  async list(params?: {
    page?: number;
    per_page?: number;
    search?: string;
    account_id?: string;
    role_id?: string;
  }): Promise<ListContactsResponse> {
    const response = await apiClient.get<ListContactsResponse>("/contacts", { params });
    return response.data;
  },

  async getById(id: string): Promise<ContactResponse> {
    const response = await apiClient.get<ContactResponse>(`/contacts/${id}`);
    return response.data;
  },

  async create(data: CreateContactFormData): Promise<ContactResponse> {
    const response = await apiClient.post<ContactResponse>("/contacts", data);
    return response.data;
  },

  async update(id: string, data: UpdateContactFormData): Promise<ContactResponse> {
    const response = await apiClient.put<ContactResponse>(`/contacts/${id}`, data);
    return response.data;
  },

  async delete(id: string): Promise<void> {
    await apiClient.delete(`/contacts/${id}`);
  },
};

