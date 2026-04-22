"use client";

import { useState, Suspense } from "react";
import dynamic from "next/dynamic";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Skeleton } from "@/components/ui/skeleton";
import { useTranslations } from "next-intl";
import { LayoutDashboard, Table } from "lucide-react";
import { TaskDetailModal } from "./task-detail-modal";

// Dynamic import untuk TaskKanbanBoard (lazy load - hanya load saat tab "board" dipilih)
const TaskKanbanBoard = dynamic<{ onTaskClick?: (task: any) => void }>(
  () => import("./task-kanban-board").then((mod) => ({ default: mod.TaskKanbanBoard })),
  {
    loading: () => <Skeleton className="h-[600px] w-full" />,
    ssr: false,
  },
);

// Dynamic import untuk TaskList (lazy load - hanya load saat tab "list" dipilih)
const TaskList = dynamic<{ onTaskClick?: (task: { id: string }) => void }>(
  () => import("./task-list").then((mod) => ({ default: mod.TaskList })),
  {
    loading: () => <Skeleton className="h-[600px] w-full" />,
    ssr: false,
  },
);

export function TaskManagement() {
  const t = useTranslations("taskManagement.page");
  const [activeTab, setActiveTab] = useState("board");
  const [viewingTaskId, setViewingTaskId] = useState<string | null>(null);

  const handleTaskClick = (task: { id: string }) => {
    setViewingTaskId(task.id);
  };

  return (
    <div className="space-y-6">
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList>
          <TabsTrigger value="board" className="gap-2">
            <LayoutDashboard className="h-4 w-4" />
            {t("tabKanban")}
          </TabsTrigger>
          <TabsTrigger value="list" className="gap-2">
            <Table className="h-4 w-4" />
            {t("tabList")}
          </TabsTrigger>
        </TabsList>
        
        <TabsContent value="board" className="mt-6">
          <Suspense fallback={<Skeleton className="h-[600px] w-full" />}>
            <TaskKanbanBoard onTaskClick={handleTaskClick} />
          </Suspense>
        </TabsContent>
        
        <TabsContent value="list" className="mt-6">
          <Suspense fallback={<Skeleton className="h-[600px] w-full" />}>
            <TaskList onTaskClick={handleTaskClick} />
          </Suspense>
        </TabsContent>
      </Tabs>

      {/* Task Detail Modal */}
      <TaskDetailModal
        taskId={viewingTaskId}
        open={!!viewingTaskId}
        onOpenChange={(open) => !open && setViewingTaskId(null)}
      />
    </div>
  );
}


