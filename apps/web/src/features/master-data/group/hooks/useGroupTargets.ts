"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { groupTargetService } from "@/features/master-data/group/services/groupTargetService";
import type { CreateGroupTargetWithUserAssignmentFormData } from "@/features/master-data/group/schemas/group-target.schema";

export function useCreateGroupTargetWithUserAssignment() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateGroupTargetWithUserAssignmentFormData) =>
      groupTargetService.createGroupTargetWithUserAssignment(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["monthly-targets"] });
    },
  });
}
