import { useStages } from "./useStages";
import type { PipelineStage } from "../types";

interface PipelineFilters {
  is_active?: boolean;
}

export const pipelineKeys = {
  all: ["pipelines"] as const,
  lists: () => [...pipelineKeys.all, "list"] as const,
  list: (filters?: PipelineFilters) => [...pipelineKeys.lists(), filters] as const,
};

// usePipelines is an alias for useStages with filtering support
export function usePipelines(filters?: PipelineFilters) {
  const { data, isLoading, error, refetch } = useStages();
  
  // Ensure data is an array
  const stagesArray = Array.isArray(data) ? data : [];
  
  // Filter stages based on filters
  const filteredData = stagesArray.filter((stage: PipelineStage) => {
    if (filters?.is_active !== undefined) {
      return stage.is_active === filters.is_active;
    }
    return true;
  });

  return {
    data: { data: filteredData },
    isLoading,
    error,
    refetch,
  };
}

// Re-export stage hooks for convenience
export { useStagesWithStats, useStage, useCreateStage, useUpdateStage, useDeleteStage, useReorderStages } from "./useStages";

// Re-export summary hook
export { usePipelineSummary } from "./useSummary";
