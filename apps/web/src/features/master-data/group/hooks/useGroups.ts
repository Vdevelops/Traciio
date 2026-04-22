"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { groupService } from "../services/groupService";
import type {
  CreateGroupFormData,
  UpdateGroupFormData,
} from "../schemas/group.schema";

export function useGroups(params?: {
  page?: number;
  per_page?: number;
  search?: string;
  status?: string;
}) {
  return useQuery({
    queryKey: ["groups", params],
    queryFn: () => groupService.list(params),
  });
}

export function useGroup(id: string) {
  return useQuery({
    queryKey: ["group", id],
    queryFn: () => groupService.getById(id),
    enabled: !!id,
  });
}

export function useCreateGroup() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateGroupFormData) =>
      groupService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["groups"] });
    },
  });
}

export function useUpdateGroup() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateGroupFormData }) =>
      groupService.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["groups"] });
      queryClient.invalidateQueries({ queryKey: ["group", variables.id] });
    },
  });
}

export function useDeleteGroup() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => groupService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["groups"] });
    },
  });
}

