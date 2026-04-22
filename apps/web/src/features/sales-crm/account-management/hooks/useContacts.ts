"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { contactService } from "../services/contactService";
import type { CreateContactFormData, UpdateContactFormData } from "../schemas/contact.schema";

export function useContacts(
  params?: {
    page?: number;
    per_page?: number;
    search?: string;
    account_id?: string;
    role_id?: string;
  },
  options?: { enabled?: boolean }
) {
  return useQuery({
    queryKey: ["contacts", params],
    queryFn: () => contactService.list(params),
    enabled: options?.enabled ?? true,
    retry: (failureCount, error) => {
      if (error && typeof error === "object" && "response" in error) {
        const axiosError = error as { response?: { status?: number } };
        if (axiosError.response?.status === 404) {
          return false;
        }
      }
      return failureCount < 1;
    },
  });
}

export function useContact(id: string) {
  return useQuery({
    queryKey: ["contacts", id],
    queryFn: () => contactService.getById(id),
    enabled: !!id,
  });
}

export function useCreateContact() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateContactFormData) => contactService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["contacts"] });
    },
  });
}

export function useUpdateContact() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateContactFormData }) =>
      contactService.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["contacts"] });
      queryClient.invalidateQueries({ queryKey: ["contacts", variables.id] });
    },
  });
}

export function useDeleteContact() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => contactService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["contacts"] });
    },
  });
}

