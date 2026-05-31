"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { leadService } from "../services/leadService";
import type { CreateLeadFormData, UpdateLeadFormData, ConvertLeadFormData } from "../schemas/lead.schema";

export function useLeads(params?: {
  page?: number;
  per_page?: number;
  status?: string;
  source?: string;
  assigned_to?: string;
  search?: string;
  sort?: string;
  order?: "asc" | "desc";
}) {
  return useQuery({
    queryKey: ["leads", params],
    queryFn: () => leadService.list(params),
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

export function useLead(id: string) {
  return useQuery({
    queryKey: ["leads", id],
    queryFn: () => leadService.getById(id),
    enabled: !!id,
  });
}

export function useLeadFormData() {
  return useQuery({
    queryKey: ["leads", "form-data"],
    queryFn: () => leadService.getFormData(),
    // Form data (provinces, defaults, industries) is static master data
    staleTime: 10 * 60 * 1000,
    gcTime: 30 * 60 * 1000,
    refetchOnWindowFocus: false,
  });
}

export function useLeadAnalytics(params?: {
  start_date?: string;
  end_date?: string;
}) {
  return useQuery({
    queryKey: ["leads", "analytics", params],
    queryFn: () => leadService.getAnalytics(params),
  });
}

export function useCreateLead() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateLeadFormData) => leadService.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["leads"] });
      queryClient.invalidateQueries({ queryKey: ["leads", "analytics"] });
    },
  });
}

export function useUpdateLead() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateLeadFormData }) =>
      leadService.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["leads"] });
      queryClient.invalidateQueries({ queryKey: ["leads", variables.id] });
      queryClient.invalidateQueries({ queryKey: ["leads", "analytics"] });
    },
  });
}

export function useDeleteLead() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => leadService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["leads"] });
      queryClient.invalidateQueries({ queryKey: ["leads", "analytics"] });
    },
  });
}

export function useConvertLead() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: ConvertLeadFormData }) =>
      leadService.convert(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["leads"] });
      queryClient.invalidateQueries({ queryKey: ["leads", variables.id] });
      queryClient.invalidateQueries({ queryKey: ["leads", "analytics"] });
      queryClient.invalidateQueries({ queryKey: ["deals"] });
      queryClient.invalidateQueries({ queryKey: ["accounts"] });
      queryClient.invalidateQueries({ queryKey: ["contacts"] });
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
      queryClient.invalidateQueries({
        predicate: (query) => query.queryKey.includes("visit-reports") || query.queryKey.includes("activities"),
      });
    },
  });
}

export function useLeadVisitReports(leadId: string, params?: {
  page?: number;
  per_page?: number;
  search?: string;
  status?: string;
  account_id?: string;
  sales_rep_id?: string;
  start_date?: string;
  end_date?: string;
}) {
  return useQuery({
    queryKey: ["leads", leadId, "visit-reports", params],
    queryFn: () => leadService.getVisitReportsByLead(leadId, params),
    enabled: !!leadId,
  });
}

export function useLeadActivities(leadId: string, params?: {
  page?: number;
  per_page?: number;
  type?: string;
  account_id?: string;
  contact_id?: string;
  user_id?: string;
  start_date?: string;
  end_date?: string;
}) {
  return useQuery({
    queryKey: ["leads", leadId, "activities", params],
    queryFn: () => leadService.getActivitiesByLead(leadId, params),
    enabled: !!leadId,
  });
}

export function useCreateAccountFromLead() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data?: { category_id?: string; create_contact?: boolean } }) =>
      leadService.createAccountFromLead(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["leads"] });
      queryClient.invalidateQueries({ queryKey: ["leads", variables.id] });
      queryClient.invalidateQueries({ queryKey: ["accounts"] });
      queryClient.invalidateQueries({ queryKey: ["contacts"] });
    },
  });
}
