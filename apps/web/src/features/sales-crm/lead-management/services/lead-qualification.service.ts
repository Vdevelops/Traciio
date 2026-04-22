import { apiClient } from "@/lib/api-client";
import type {
  LeadQualificationChecklist,
  UpdateLeadQualificationRequest,
} from "../types/qualification";
import type { ApiResponse, PaginatedResponse } from "@/types/api";

// Assume we import Task, VisitReport, Activity from somewhere common if needed,
// for now we'll do an any or minimal type definition just to get it functional.
// This is the structure expected by the backend and UI.

export const leadQualificationService = {
  // Get qualification checklist
  async getQualification(leadId: string): Promise<LeadQualificationChecklist> {
    const response = await apiClient.get<
      ApiResponse<LeadQualificationChecklist>
    >(`/leads/${leadId}/qualification`);
    return response.data.data;
  },

  // Update qualification
  async updateQualification(
    leadId: string,
    data: UpdateLeadQualificationRequest,
  ): Promise<LeadQualificationChecklist> {
    const response = await apiClient.post<
      ApiResponse<LeadQualificationChecklist>
    >(`/leads/${leadId}/qualification`, data);
    return response.data.data;
  },

  // Get lead tasks
  async getLeadTasks(
    leadId: string,
    page = 1,
    perPage = 20,
  ): Promise<PaginatedResponse<any>> {
    const response = await apiClient.get<ApiResponse<PaginatedResponse<any>>>(
      `/leads/${leadId}/tasks`,
      { params: { page, per_page: perPage } },
    );
    return response.data.data;
  },

  // Create task for lead
  async createLeadTask(leadId: string, data: any): Promise<any> {
    const response = await apiClient.post<ApiResponse<any>>(
      `/leads/${leadId}/tasks`,
      data,
    );
    return response.data.data;
  },

  // Get lead visit reports
  async getLeadVisitReports(leadId: string): Promise<any[]> {
    const response = await apiClient.get<ApiResponse<any[]>>(
      `/leads/${leadId}/visit-reports`,
    );
    return response.data.data;
  },

  // Get lead activities
  async getLeadActivities(leadId: string): Promise<any[]> {
    const response = await apiClient.get<ApiResponse<any[]>>(
      `/leads/${leadId}/activities`,
    );
    return response.data.data;
  },
};
