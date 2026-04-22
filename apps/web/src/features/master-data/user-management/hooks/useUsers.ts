"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { userService, roleService } from "../services/userService";
import { groupService } from "@/features/master-data/group/services/groupService";
import type { CreateUserFormData, UpdateUserFormData } from "../schemas/user.schema";

export function useUsers(params?: {
  page?: number;
  per_page?: number;
  search?: string;
  status?: string;
  role_id?: string;
  group_id?: string;
}) {
  return useQuery({
    queryKey: ["users", params],
    queryFn: () => userService.list(params),
    retry: (failureCount, error) => {
      // Don't retry on 404 errors (might be expected)
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

export function useUser(id: string) {
  return useQuery({
    queryKey: ["users", id],
    queryFn: () => userService.getById(id),
    enabled: !!id,
  });
}

export function useCreateUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateUserFormData) => userService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
    },
  });
}

export function useUpdateUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateUserFormData }) => {
      return userService.update(id, data);
    },
    onSuccess: (result, variables) => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
      queryClient.invalidateQueries({ queryKey: ["users", variables.id] });
    },
  });
}

export function useDeleteUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => userService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
    },
  });
}

export function useUserPermissions(userId: string) {
  return useQuery({
    queryKey: ["users", userId, "permissions"],
    queryFn: () => userService.getPermissions(userId),
    enabled: !!userId,
  });
}

export function useRoles() {
  return useQuery({
    queryKey: ["roles"],
    queryFn: () => roleService.list(),
  });
}

export function useRole(id: string) {
  return useQuery({
    queryKey: ["roles", id],
    queryFn: () => roleService.getById(id),
    enabled: !!id,
  });
}

export function useGroups() {
  return useQuery({
    queryKey: ["groups"],
    queryFn: () => groupService.list({ per_page: 100 }),
  });
}

