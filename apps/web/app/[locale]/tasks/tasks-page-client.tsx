"use client";

import { Suspense } from "react";
import dynamic from "next/dynamic";
import { Skeleton } from "@/components/ui/skeleton";

// Dynamic import untuk TaskManagement dengan code splitting
const TaskManagement = dynamic(
  () =>
    import("@/features/sales-crm/task-management/components/task-management").then(
      (mod) => ({ default: mod.TaskManagement }),
    ),
  {
    loading: () => <Skeleton className="h-[600px] w-full" />,
    ssr: false, // Client component, no SSR needed
  },
);

export function TasksPageClient() {
  return (
    <Suspense fallback={<Skeleton className="h-[600px] w-full" />}>
      <TaskManagement />
    </Suspense>
  );
}


