"use client";

import { useState } from "react";
import { Plus } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { TaskKanbanCard } from "./task-kanban-card";
import { TaskScrollLoader } from "./task-scroll-loader";
import { useProgressiveTaskKanban } from "../hooks/useProgressiveTaskKanban";
import { useTranslations } from "next-intl";
import type { Task } from "../types";

interface TaskKanbanBoardProps {
  readonly onTaskClick?: (task: Task) => void;
  readonly onCreateTask?: () => void;
}

export function TaskKanbanBoard({ onTaskClick, onCreateTask }: TaskKanbanBoardProps) {
  const t = useTranslations("taskManagement.kanban");

  const {
    statuses,
    tasksByStatus,
    loadingByStatus,
    hasNextPageByStatus,
    isFetchingNextPageByStatus,
    fetchNextPageForStatus,
    isLoading,
  } = useProgressiveTaskKanban({ pageSize: 20 });

  if (isLoading) {
    return (
      <div className="space-y-6">
        <div>
          <div className="h-8 bg-muted animate-pulse rounded w-64 mb-2" />
          <div className="h-4 bg-muted animate-pulse rounded w-96" />
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          {[...Array(4)].map((_, i) => (
            <div key={i} className="p-4 border rounded-lg">
              <div className="space-y-4">
                <div className="h-12 bg-muted animate-pulse rounded" />
                <div className="space-y-2">
                  {[...Array(2)].map((_, j) => (
                    <div key={j} className="h-32 bg-muted animate-pulse rounded" />
                  ))}
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header with Create Button */}
      {onCreateTask && (
        <div className="flex justify-end">
          <Button onClick={onCreateTask}>
            <Plus className="h-4 w-4 mr-2" />
            {t("addTask") || "Add Task"}
          </Button>
        </div>
      )}

      {/* Kanban Board - Exact same styling as Pipeline Kanban */}
      <div className="flex overflow-x-auto overflow-y-hidden pb-6 gap-6 -mx-6 px-6 scrollbar-thin scrollbar-thumb-muted-foreground/20">
        {statuses.map((status) => {
          const statusTasks = tasksByStatus[status.id] || [];

          return (
            <div
              key={status.id}
              className="shrink-0 w-80 h-full flex flex-col"
            >
              {/* Column Header - Exact same styling */}
              <div className="flex items-center gap-2.5 mb-4 shrink-0 pb-3 px-1 border-b border-border/40">
                <div
                  className="w-2.5 h-2.5 rounded-full shrink-0 ring-2 ring-offset-2 ring-offset-background/50 shadow-sm"
                  style={{
                    backgroundColor: status.color,
                  }}
                />
                <h3 className="font-semibold text-sm truncate flex-1 tracking-tight text-foreground/80">
                  {status.name}
                </h3>
                <Badge variant="secondary" className="shrink-0 text-[10px] font-bold h-5 px-1.5 bg-muted/40 border-none text-muted-foreground">
                  {statusTasks.length}
                </Badge>
              </div>

              {/* Column Content - Exact same styling */}
              <div className="space-y-3 min-h-[200px] flex-1 overflow-y-auto pr-1">
                {loadingByStatus[status.id] ? (
                  // Loading skeleton
                  <div className="space-y-3">
                    {[...Array(3)].map((_, i) => (
                      <div key={i} className="bg-muted/20 rounded-lg p-3 animate-pulse">
                        <div className="h-4 bg-muted/40 rounded w-3/4 mb-2" />
                        <div className="h-3 bg-muted/30 rounded w-1/2" />
                      </div>
                    ))}
                  </div>
                ) : statusTasks.length === 0 ? (
                  // Empty state
                  <div className="flex flex-col items-center justify-center py-12 text-center">
                    <div className="w-12 h-12 rounded-full bg-muted flex items-center justify-center mb-3">
                      <Plus className="h-5 w-5 text-muted-foreground" />
                    </div>
                    <p className="text-sm text-muted-foreground font-medium">
                      {t("noTasks") || "No tasks"}
                    </p>
                    <p className="text-xs text-muted-foreground mt-1">
                      {t("noTasksHint") || "Tasks will appear here"}
                    </p>
                  </div>
                ) : (
                  <>
                    {statusTasks.map((task) => (
                      <div key={task.id}>
                        <TaskKanbanCard
                          task={task}
                          onClick={() => onTaskClick?.(task)}
                        />
                      </div>
                    ))}

                    {/* Auto-load on scroll */}
                    <TaskScrollLoader
                      onLoadMore={() => fetchNextPageForStatus(status.id)}
                      hasMore={hasNextPageByStatus[status.id] ?? false}
                      isLoading={isFetchingNextPageByStatus[status.id] ?? false}
                    />

                    {/* Fallback "Load More" button */}
                    {hasNextPageByStatus[status.id] && !isFetchingNextPageByStatus[status.id] && (
                      <Button
                        variant="ghost"
                        size="sm"
                        className="w-full text-xs"
                        onClick={() => fetchNextPageForStatus(status.id)}
                      >
                        <Plus className="h-3 w-3 mr-1" />
                        Load More
                      </Button>
                    )}
                  </>
                )}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
