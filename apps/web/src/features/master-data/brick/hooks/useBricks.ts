import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { brickService } from "../services/brickService";
import type {
  CreateBrickFormData,
  UpdateBrickFormData,
  CreateBrickTargetDistributionFormData,
  UpdateBrickTargetDistributionFormData,
} from "../schemas/brick.schema";

export function useBricks(params?: {
  page?: number;
  per_page?: number;
  search?: string;
  province?: string;
  regency?: string;
  manager_id?: string;
  status?: string;
}) {
  return useQuery({
    queryKey: ["bricks", params],
    queryFn: () => brickService.list(params),
  });
}

export function useBrick(id: string) {
  return useQuery({
    queryKey: ["brick", id],
    queryFn: () => brickService.getById(id),
    enabled: !!id,
  });
}

export function useCreateBrick() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateBrickFormData) => brickService.create(data),
    // Optimistic update: immediately add a placeholder brick to all list caches
    onMutate: async (newBrick) => {
      await queryClient.cancelQueries({ queryKey: ["bricks"] });
      const previousSnapshots = queryClient.getQueriesData({ queryKey: ["bricks"] });
      queryClient.setQueriesData({ queryKey: ["bricks"] }, (old: any) => {
        if (!old?.data) return old;
        const placeholder = {
          id: `optimistic-${Date.now()}`,
          name: newBrick.name,
          code: "...",
          province: newBrick.province,
          regency: newBrick.regency,
          status: newBrick.status ?? "active",
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
        };
        return { ...old, data: [placeholder, ...old.data] };
      });
      return { previousSnapshots };
    },
    onError: (_err, _vars, context: any) => {
      // Roll back to previous state on error
      if (context?.previousSnapshots) {
        for (const [queryKey, snapshot] of context.previousSnapshots) {
          queryClient.setQueryData(queryKey, snapshot);
        }
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bricks"] });
    },
  });
}

export function useUpdateBrick() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateBrickFormData }) =>
      brickService.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["bricks"] });
      queryClient.invalidateQueries({ queryKey: ["brick", variables.id] });
    },
  });
}

export function useDeleteBrick() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => brickService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["bricks"] });
    },
  });
}

export function useBrickSales(brickId: string) {
  return useQuery({
    queryKey: ["brick-sales", brickId],
    queryFn: () => brickService.getSalesInBrick(brickId),
    enabled: !!brickId,
  });
}

export function useBrickTargetWithDistributions(
  brickId: string,
  year: number,
  month: number
) {
  return useQuery({
    queryKey: ["brick-target-distributions", brickId, year, month],
    queryFn: () => brickService.getBrickTargetWithDistributions(brickId, year, month),
    enabled: !!brickId && !!year && !!month,
  });
}

export function useDistributeBrickTarget() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      brickId,
      targetId,
      data,
    }: {
      brickId: string;
      targetId: string;
      data: CreateBrickTargetDistributionFormData;
    }) => brickService.distributeTarget(brickId, targetId, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: ["brick-target-distributions", variables.brickId],
      });
    },
  });
}

export function useUpdateBrickTargetDistribution() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      brickId,
      targetId,
      distributionId,
      data,
    }: {
      brickId: string;
      targetId: string;
      distributionId: string;
      data: UpdateBrickTargetDistributionFormData;
    }) =>
      brickService.updateDistribution(brickId, targetId, distributionId, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: ["brick-target-distributions", variables.brickId],
      });
    },
  });
}

export function useDeleteBrickTargetDistribution() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      brickId,
      targetId,
      distributionId,
    }: {
      brickId: string;
      targetId: string;
      distributionId: string;
    }) =>
      brickService.deleteDistribution(brickId, targetId, distributionId),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({
        queryKey: ["brick-target-distributions", variables.brickId],
      });
    },
  });
}

