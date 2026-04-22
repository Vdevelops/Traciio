import apiClient from "@/lib/api-client";
import type {
  AnalyzeVisitReportRequest,
  AnalyzeVisitReportResponse,
  ChatRequest,
  ChatAPIResponse,
  AISettingsResponse,
} from "../types";

export const aiService = {
  async analyzeVisitReport(
    data: AnalyzeVisitReportRequest
  ): Promise<AnalyzeVisitReportResponse> {
    const response = await apiClient.post<AnalyzeVisitReportResponse>(
      "/ai/analyze/visit-report",
      data
    );
    return response.data;
  },

  async chat(data: ChatRequest): Promise<ChatAPIResponse> {
    const response = await apiClient.post<ChatAPIResponse>("/ai/chat", data);
    return response.data;
  },

  async getSettings(): Promise<{ success: boolean; data: AISettingsResponse }> {
    const response = await apiClient.get<{ success: boolean; data: AISettingsResponse }>(
      "/ai/settings"
    );
    return response.data;
  },

  async updateSettings(
    data: Partial<AISettingsResponse>
  ): Promise<{ success: boolean; data: AISettingsResponse }> {
    const response = await apiClient.put<{ success: boolean; data: AISettingsResponse }>(
      "/ai/settings",
      data
    );
    return response.data;
  },

};

