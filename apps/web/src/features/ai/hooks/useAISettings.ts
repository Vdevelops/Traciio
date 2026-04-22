import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { aiService } from "../services/aiService";
import type { AISettingsResponse } from "../types";

export function useAISettings() {
  const queryClient = useQueryClient();

  // Fetch settings
  const {
    data: settingsResponse,
    isLoading,
    error,
    refetch: refetchSettings,
  } = useQuery({
    queryKey: ["ai-settings"],
    queryFn: async () => {
      const response = await aiService.getSettings();
      return response.data;
    },
  });

  // Update settings mutation
  const updateMutation = useMutation({
    mutationFn: async (data: Partial<AISettingsResponse>) => {
      return await aiService.updateSettings(data);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["ai-settings"] });
      toast.success("AI settings updated successfully");
    },
    onError: (error: unknown) => {
      const message =
        error instanceof Error ? error.message : "Failed to update settings";
      toast.error(message);
    },
  });

  const updateDataPrivacy = (
    key: string,
    value: boolean
  ) => {
    if (!settingsResponse) return;
    
    updateMutation.mutate({
      data_privacy: {
        ...settingsResponse.data_privacy,
        [key]: value,
      } as AISettingsResponse["data_privacy"],
    });
  };

  const toggleEnabled = () => {
    if (!settingsResponse) return;
    updateMutation.mutate({ enabled: !settingsResponse.enabled });
  };

  const updateProvider = (provider: string) => {
    updateMutation.mutate({ provider });
  };

  const updateModel = (model: string) => {
    updateMutation.mutate({ model });
  };

  const updateAPIKey = (apiKey: string) => {
    updateMutation.mutate({ api_key: apiKey });
  };

  const updateTimezone = (timezone: string) => {
    updateMutation.mutate({ timezone });
  };

  // Convert AISettingsResponse to AISettings format for backward compatibility
  const settings = settingsResponse
    ? {
        enabled: settingsResponse.enabled,
        data_privacy: settingsResponse.data_privacy,
        provider: settingsResponse.provider,
        model: settingsResponse.model,
        base_url: settingsResponse.base_url,
        timezone: settingsResponse.timezone || "Asia/Jakarta",
      }
    : {
        enabled: true,
        data_privacy: {
          // Sales domain
          allow_leads: true,
          allow_deals: true,
          allow_visit_reports: true,
          allow_activities: true,
          allow_tasks: true,
          allow_schedule: true,
          allow_pipelines: true,
          // Customer domain
          allow_accounts: true,
          allow_contacts: true,
          // Inventory domain
          allow_products: true,
          // Analytics domain
          allow_sales_performance: true,
          allow_product_analysis: true,
          allow_reports: true,
          // Management domain
          allow_users: true,
          allow_roles: true,
          allow_groups: true,
          allow_brick_management: true,
          allow_target: true,
          // Route Optimization domain
          allow_route_optimization: true,
        } as AISettingsResponse["data_privacy"],
        provider: "cerebras",
        model: "llama-3.1-8b",
        base_url: undefined,
        timezone: "Asia/Jakarta",
      };

  return {
    settings,
    settingsResponse,
    isLoading,
    error,
    updateDataPrivacy,
    toggleEnabled,
    updateProvider,
    updateModel,
    updateAPIKey,
    updateTimezone,
    isUpdating: updateMutation.isPending,
    refetch: refetchSettings,
  };
}

