import apiClient from "@/lib/api-client";
import type { CreateGroupTargetWithUserAssignmentFormData } from "@/features/master-data/group/schemas/group-target.schema";

export const groupTargetService = {
  async createGroupTargetWithUserAssignment(data: CreateGroupTargetWithUserAssignmentFormData) {
    const response = await apiClient.post("/monthly-targets/group-with-users", data);
    return response.data;
  },
};
