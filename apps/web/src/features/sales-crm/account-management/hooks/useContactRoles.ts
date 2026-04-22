"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { contactRoleService } from "../services/contactRoleService";
import type { CreateContactRoleFormData, UpdateContactRoleFormData } from "../schemas/contact-role.schema";

export function useContactRoles() {
  return useQuery({
    queryKey: ["contact-roles"],
    queryFn: () => contactRoleService.list(),
  });
}

export function useContactRole(id: string) {
  return useQuery({
    queryKey: ["contact-role", id],
    queryFn: () => contactRoleService.getById(id),
    enabled: !!id,
  });
}

export function useCreateContactRole() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateContactRoleFormData) => contactRoleService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["contact-roles"] });
      queryClient.invalidateQueries({ queryKey: ["contacts"] });
    },
  });
}

export function useUpdateContactRole() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateContactRoleFormData }) =>
      contactRoleService.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["contact-roles"] });
      queryClient.invalidateQueries({ queryKey: ["contact-role", variables.id] });
      queryClient.invalidateQueries({ queryKey: ["contacts"] });
    },
  });
}

export function useDeleteContactRole() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => contactRoleService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["contact-roles"] });
      queryClient.invalidateQueries({ queryKey: ["contacts"] });
    },
  });
}

