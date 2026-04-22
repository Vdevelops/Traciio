"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { accountService } from "../services/accountService";
import type { CreateAccountFormData, UpdateAccountFormData } from "../schemas/account.schema";

export function useAccounts(params?: {
  page?: number;
  per_page?: number;
  search?: string;
  status?: string;
  category_id?: string;
  assigned_to?: string;
}) {
  return useQuery({
    queryKey: ["accounts", params],
    queryFn: () => accountService.list(params),
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

export function useAccount(id: string) {
  return useQuery({
    queryKey: ["accounts", id],
    queryFn: () => accountService.getById(id),
    enabled: !!id,
  });
}

export function useCreateAccount() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateAccountFormData) => accountService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["accounts"] });
    },
  });
}

export function useUpdateAccount() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateAccountFormData }) =>
      accountService.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["accounts"] });
      queryClient.invalidateQueries({ queryKey: ["accounts", variables.id] });
    },
  });
}

export function useDeleteAccount() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => accountService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["accounts"] });
    },
  });
}

export function useAccountsForMap(params?: {
  status?: string;
}) {
  return useQuery({
    queryKey: ["accounts", "map", params],
    queryFn: () => accountService.listForMap(params),
    retry: 2,
    staleTime: 5 * 60 * 1000, // Cache for 5 minutes (geocoding takes time)
  });
}

export function useAccountsByBBox(params: {
  north: number;
  south: number;
  east: number;
  west: number;
  search?: string;
  status?: string;
  category_id?: string;
  limit?: number;
} | null) {
  return useQuery({
    queryKey: ["accounts", "bbox", params],
    queryFn: () => accountService.listByBBox(params!),
    enabled: !!params,
    retry: 1,
    staleTime: 10 * 1000, // Cache viewport data for 10 seconds
    placeholderData: (prev) => prev, // Keep previous data while loading new viewport
  });
}

