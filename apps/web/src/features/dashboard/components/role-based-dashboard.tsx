"use client";

import { Suspense } from "react";
import { useAuthStore } from "@/features/auth/stores/useAuthStore";
import { DashboardOverview } from "./dashboard-overview";
import { PipelineSummary } from "./pipeline-summary";
import { LeadsBySource } from "./leads-by-source";
import { UpcomingTasks } from "./upcoming-tasks";
import { StaggerContainer } from "@/components/motion";
import { Skeleton } from "@/components/ui/skeleton";
import { AdminDashboard } from "./widgets/admin-dashboard";
import { SalesDashboard } from "./widgets/sales-dashboard";
import { SalesManagerDashboard } from "./widgets/sales-manager-dashboard";
import { AnalystDashboard } from "./widgets/analyst-dashboard";
interface RoleBasedDashboardProps {
  startDate?: string;
  endDate?: string;
}

export function RoleBasedDashboard({ startDate, endDate }: RoleBasedDashboardProps) {
  const { user } = useAuthStore();
  const role = user?.role ?? "";

  // Admin Dashboard
  if (role === "admin") {
    return <AdminDashboard startDate={startDate} endDate={endDate} />;
  }
  
  // Sales Manager Dashboard
  if (role === "sales_manager") {
    return <SalesManagerDashboard startDate={startDate} endDate={endDate} />;
  }

  // Sales/Field Staff Dashboard
  if (role === "sales") {
    return <SalesDashboard startDate={startDate} endDate={endDate} />;
  }

  if (role === "analyst") {
    return <AnalystDashboard startDate={startDate} endDate={endDate} />;
  }

  // Default dashboard (fallback for other roles)
  return (
    <div className="space-y-4 sm:space-y-6">
      <Suspense
        fallback={
          <div className="grid gap-3 sm:gap-4 md:gap-6 grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
            {Array.from({ length: 4 }).map((_, i) => (
              <Skeleton key={i} className="h-24 sm:h-32" />
            ))}
          </div>
        }
      >
        <DashboardOverview startDate={startDate} endDate={endDate} />
      </Suspense>

      <StaggerContainer className="grid gap-3 sm:gap-4 md:gap-6 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
        <Suspense fallback={<Skeleton className="h-48 sm:h-64" />}>
          <LeadsBySource />
        </Suspense>
        <Suspense fallback={<Skeleton className="h-48 sm:h-64" />}>
          <UpcomingTasks />
        </Suspense>
        <Suspense fallback={<Skeleton className="h-48 sm:h-64" />}>
          <PipelineSummary />
        </Suspense>
      </StaggerContainer>
    </div>
  );
}
