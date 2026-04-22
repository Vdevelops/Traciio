"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import { taskService } from "../services/taskService";
import type { AddLeadFromTaskRequest } from "../types";
import { toast } from "sonner";

export function useTaskQuickActions(taskId: string) {
  const queryClient = useQueryClient();

  const addLeadMutation = useMutation({
    mutationFn: (data: AddLeadFromTaskRequest) =>
      taskService.addLeadFromTask(taskId, data),
    onSuccess: (result) => {
      toast.success(
        `Lead "${result.lead?.full_name ?? "New Lead"}" created and linked to task`,
      );

      // Invalidate related queries
      queryClient.invalidateQueries({ queryKey: ["task", taskId] });
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
      queryClient.invalidateQueries({ queryKey: ["leads"] });
    },
    onError: (error) => {
      toast.error(
        error instanceof Error ? error.message : "Failed to create lead",
      );
    },
  });

  return {
    addLeadFromTask: addLeadMutation.mutate,
    isAddingLead: addLeadMutation.isPending,
    addLeadResult: addLeadMutation.data,
  };
}
