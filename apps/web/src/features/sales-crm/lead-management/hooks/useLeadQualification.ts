"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { leadQualificationService } from "../services/lead-qualification.service";
import type { UpdateLeadQualificationRequest } from "../types/qualification";

export function useLeadQualification(leadId: string) {
  const queryClient = useQueryClient();

  const { data, isLoading, error } = useQuery({
    queryKey: ["lead-qualification", leadId],
    queryFn: () => leadQualificationService.getQualification(leadId),
    enabled: !!leadId,
  });

  const updateMutation = useMutation({
    mutationFn: (req: UpdateLeadQualificationRequest) =>
      leadQualificationService.updateQualification(leadId, req),
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["lead-qualification", leadId],
      });
      queryClient.invalidateQueries({ queryKey: ["lead", leadId] });
    },
  });

  return {
    qualification: data,
    isLoading,
    error,
    updateQualification: updateMutation.mutate,
    isUpdating: updateMutation.isPending,
  };
}
