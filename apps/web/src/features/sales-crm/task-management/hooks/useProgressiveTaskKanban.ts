"use client";

import { useState, useCallback, useMemo } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { taskService } from "../services/taskService";
import type { Task, TaskStatus } from "../types";

interface UseProgressiveTaskKanbanParams {
  readonly pageSize?: number;
}

interface ApiResponse<T> {
  success: boolean;
  data: T;
  meta?: {
    pagination?: {
      page: number;
      per_page: number;
      total: number;
      total_pages: number;
      has_next: boolean;
      has_prev: boolean;
    };
  };
}

// Cached data structure returned by the initial query
interface KanbanInitialData {
  tasks: Record<string, Task[]>;
  pages: Record<string, number>;
  hasMore: Record<string, boolean>;
}

// Define task statuses as "stages"
const TASK_STATUSES: { id: TaskStatus; name: string; color: string }[] = [
  { id: "pending", name: "Pending", color: "#f59e0b" },
  { id: "completed", name: "Completed", color: "#10b981" },
];

export function useProgressiveTaskKanban(
  params: UseProgressiveTaskKanbanParams = {}
) {
  const { pageSize = 20 } = params;
  const queryClient = useQueryClient();

  // Local state only for pagination beyond the initial load
  const [extraTasks, setExtraTasks] = useState<Record<string, Task[]>>({});
  const [currentPages, setCurrentPages] = useState<Record<string, number>>({});
  const [paginationHasMore, setPaginationHasMore] = useState<Record<string, boolean>>({});
  const [loadingMore, setLoadingMore] = useState<Record<string, boolean>>({});

  // Initial load - return actual data so TanStack Query caches it properly
  const { data: initialData, isLoading: initialLoading } = useQuery<KanbanInitialData>({
    queryKey: ["tasks", "kanban", "initial"],
    queryFn: async () => {
      const promises = TASK_STATUSES.map(async (status) => {
        const response = await taskService.list({
          status: status.id,
          page: 1,
          per_page: pageSize,
        }) as unknown as ApiResponse<Task[]>;

        const currentPage = response.meta?.pagination?.page ?? 0;
        const totalPages = response.meta?.pagination?.total_pages ?? 0;
        const hasMore = currentPage < totalPages;

        return {
          statusId: status.id,
          tasks: response.data || [],
          hasMore,
        };
      });

      const results = await Promise.all(promises);

      const tasks: Record<string, Task[]> = {};
      const pages: Record<string, number> = {};
      const hasMore: Record<string, boolean> = {};

      results.forEach((result) => {
        tasks[result.statusId] = result.tasks;
        pages[result.statusId] = 1;
        hasMore[result.statusId] = result.hasMore;
      });

      // Reset pagination state when initial data is refetched
      setExtraTasks({});
      setCurrentPages(pages);
      setPaginationHasMore(hasMore);

      return { tasks, pages, hasMore };
    },
    staleTime: 30000,
  });

  // Merge initial cached data with extra paginated tasks
  const mergedHasMore = useMemo(() => {
    if (!initialData) return paginationHasMore;
    return { ...initialData.hasMore, ...paginationHasMore };
  }, [initialData, paginationHasMore]);

  const mergedPages = useMemo(() => {
    if (!initialData) return currentPages;
    return { ...initialData.pages, ...currentPages };
  }, [initialData, currentPages]);

  // Function to load more tasks for a specific status
  const fetchNextPageForStatus = useCallback(
    async (statusId: TaskStatus) => {
      if (loadingMore[statusId] || !mergedHasMore[statusId]) {
        return;
      }

      const page = mergedPages[statusId] || 1;
      const nextPage = page + 1;

      setLoadingMore((prev) => ({ ...prev, [statusId]: true }));

      try {
        const response = await taskService.list({
          status: statusId,
          page: nextPage,
          per_page: pageSize,
        }) as unknown as ApiResponse<Task[]>;

        const newTasks = response.data || [];
        const currentPageNum = response.meta?.pagination?.page ?? 0;
        const totalPages = response.meta?.pagination?.total_pages ?? 0;
        const hasMore = currentPageNum < totalPages;

        // Append new tasks, filtering out duplicates
        setExtraTasks((prev) => {
          const existing = prev[statusId] || [];
          const existingIds = new Set([
            ...existing.map((t) => t.id),
            ...(initialData?.tasks[statusId] || []).map((t) => t.id),
          ]);
          const uniqueNewTasks = newTasks.filter((t) => !existingIds.has(t.id));

          return {
            ...prev,
            [statusId]: [...existing, ...uniqueNewTasks],
          };
        });

        setCurrentPages((prev) => ({ ...prev, [statusId]: nextPage }));
        setPaginationHasMore((prev) => ({ ...prev, [statusId]: hasMore }));
      } catch (error) {
        console.error(`[Pagination] Error loading more tasks for status ${statusId}:`, error);
      } finally {
        setLoadingMore((prev) => ({ ...prev, [statusId]: false }));
      }
    },
    [loadingMore, mergedHasMore, mergedPages, pageSize, initialData]
  );

  // Derive tasks from cached query data + extra paginated tasks
  const tasksByStatus = useMemo(() => {
    const result: Record<string, Task[]> = {};
    TASK_STATUSES.forEach((status) => {
      const base = initialData?.tasks[status.id] || [];
      const extra = extraTasks[status.id] || [];
      result[status.id] = [...base, ...extra];
    });
    return result;
  }, [initialData, extraTasks]);

  const loadingByStatus = useMemo(() => {
    const result: Record<string, boolean> = {};
    TASK_STATUSES.forEach((status) => {
      result[status.id] = initialLoading;
    });
    return result;
  }, [initialLoading]);

  const hasNextPageByStatus = useMemo(() => {
    const result: Record<string, boolean> = {};
    TASK_STATUSES.forEach((status) => {
      result[status.id] = mergedHasMore[status.id] ?? false;
    });
    return result;
  }, [mergedHasMore]);

  const isFetchingNextPageByStatus = useMemo(() => {
    const result: Record<string, boolean> = {};
    TASK_STATUSES.forEach((status) => {
      result[status.id] = loadingMore[status.id] ?? false;
    });
    return result;
  }, [loadingMore]);

  const refreshAll = useCallback(() => {
    setExtraTasks({});
    queryClient.invalidateQueries({ queryKey: ["tasks", "kanban"] });
  }, [queryClient]);

  return {
    statuses: TASK_STATUSES,
    tasksByStatus,
    loadingByStatus,
    hasNextPageByStatus,
    isFetchingNextPageByStatus,
    fetchNextPageForStatus,
    isLoading: initialLoading,
    refreshAll,
  };
}
