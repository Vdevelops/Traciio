"use client";

import { useQuery, useMutation, useQueryClient, type QueryClient } from "@tanstack/react-query";
import { leadService } from "../services/leadService";
import type { CreateLeadFormData, UpdateLeadFormData, ConvertLeadFormData } from "../schemas/lead.schema";
import type { Lead, LeadResponse, ListLeadsResponse } from "../types";

function isLeadListQueryKey(queryKey: readonly unknown[]) {
  return (
    queryKey[0] === "leads" &&
    queryKey.length === 2 &&
    (queryKey[1] === undefined || queryKey[1] === null || typeof queryKey[1] === "object")
  );
}

function replaceLeadInList(current: ListLeadsResponse | undefined, updatedLead: Lead) {
  if (!current?.data) return current;

  return {
    ...current,
    data: current.data.map((lead) => (lead.id === updatedLead.id ? { ...lead, ...updatedLead } : lead)),
  };
}

function invalidateLeadDerivedQueries(queryClient: QueryClient) {
  return Promise.all([
    queryClient.invalidateQueries({ queryKey: ["sales-overview"] }),
    queryClient.invalidateQueries({ queryKey: ["reports"] }),
    queryClient.invalidateQueries({ queryKey: ["dashboard"] }),
  ]);
}

function invalidateLeadConversionDerivedQueries(queryClient: QueryClient) {
  return Promise.all([
    invalidateLeadDerivedQueries(queryClient),
    queryClient.invalidateQueries({ queryKey: ["product-analytics"] }),
  ]);
}

type DealLike = {
  id: string;
  stage_id?: string;
} & Record<string, unknown>;

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
      void invalidateLeadDerivedQueries(queryClient);
    },
  });
}

export function useUpdateLead() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateLeadFormData }) =>
      leadService.update(id, data),
    onSuccess: (response, variables) => {
      const updatedLead = response.data;

      queryClient.setQueryData<LeadResponse>(["leads", variables.id], response);
      queryClient.setQueriesData<ListLeadsResponse>(
        { predicate: (query) => isLeadListQueryKey(query.queryKey) },
        (current) => replaceLeadInList(current, updatedLead)
      );

      queryClient.invalidateQueries({ queryKey: ["leads"] });
      queryClient.invalidateQueries({ queryKey: ["leads", variables.id] });
      queryClient.invalidateQueries({ queryKey: ["leads", "analytics"] });
      void invalidateLeadDerivedQueries(queryClient);
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
      void invalidateLeadDerivedQueries(queryClient);
    },
  });
}

export function useConvertLead() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: ConvertLeadFormData }) =>
      leadService.convert(id, data),
    onSuccess: (response, variables) => {
      const convertedLead = response.data.lead;
      const convertedDeal = response.data.opportunity as DealLike | undefined;

      queryClient.setQueryData<LeadResponse>(["leads", variables.id], {
        success: response.success,
        data: convertedLead,
        timestamp: response.timestamp,
        request_id: response.request_id,
      });
      queryClient.setQueriesData<ListLeadsResponse>(
        { predicate: (query) => isLeadListQueryKey(query.queryKey) },
        (current) => replaceLeadInList(current, convertedLead)
      );

      if (convertedDeal?.id) {
        queryClient.setQueryData(["deals", "detail", convertedDeal.id], {
          success: response.success,
          data: convertedDeal,
          timestamp: response.timestamp,
          request_id: response.request_id,
        });
      }

      queryClient.invalidateQueries({ queryKey: ["leads"] });
      queryClient.invalidateQueries({ queryKey: ["leads", variables.id] });
      queryClient.invalidateQueries({ queryKey: ["leads", "analytics"] });
      queryClient.invalidateQueries({ queryKey: ["deals"] });
      if (convertedDeal?.id) {
        queryClient.invalidateQueries({ queryKey: ["deals", "detail", convertedDeal.id] });
      }
      queryClient.invalidateQueries({ queryKey: ["accounts"] });
      queryClient.invalidateQueries({ queryKey: ["contacts"] });
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
      void invalidateLeadConversionDerivedQueries(queryClient);
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
