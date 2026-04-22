"use client";

import { Suspense } from "react";
import { StaggerContainer } from "@/components/motion";
import { Skeleton } from "@/components/ui/skeleton";
import { UpcomingTasks } from "../upcoming-tasks";
import { PipelineSummary } from "../pipeline-summary";
import { LeadsBySource } from "../leads-by-source";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Calendar, UserPlus, Bell } from "lucide-react";
import {
  useSalesTodayTasks,
  useSalesAssignedLeads,
  useSalesUpcomingVisits,
  useSalesReminders,
} from "../../hooks/useDashboard";

interface SalesDashboardProps {
  startDate?: string;
  endDate?: string;
}

export function SalesDashboard({ startDate, endDate }: Readonly<SalesDashboardProps>) {
  return (
    <div className="space-y-4 sm:space-y-6 md:space-y-8">
      {/* ========================================================================
          SECTION 1: DAILY WORK
          Tasks, Leads, Visits, and Reminders
      ========================================================================= */}
      <section className="space-y-3 sm:space-y-4 md:space-y-6">
        <div className="flex items-center gap-1.5 sm:gap-2">
          <div className="h-1 w-1 rounded-full bg-primary" />
          <h2 className="text-sm sm:text-base md:text-lg font-medium text-muted-foreground uppercase tracking-wide">
            Daily Work
          </h2>
        </div>

        <div className="grid gap-3 sm:gap-4 md:gap-6 grid-cols-2 md:grid-cols-2 lg:grid-cols-4">
        <Suspense fallback={<Skeleton className="h-24 sm:h-32" />}>
          <SalesTodayTasksCard startDate={startDate} endDate={endDate} />
        </Suspense>

        <Suspense fallback={<Skeleton className="h-24 sm:h-32" />}>
          <SalesAssignedLeadsCard startDate={startDate} endDate={endDate} />
        </Suspense>

        <Suspense fallback={<Skeleton className="h-24 sm:h-32" />}>
          <SalesUpcomingVisitsCard startDate={startDate} endDate={endDate} />
        </Suspense>

        <Suspense fallback={<Skeleton className="h-24 sm:h-32" />}>
          <SalesRemindersCard startDate={startDate} endDate={endDate} />
        </Suspense>
        </div>
      </section>

      {/* ========================================================================
          SECTION 2: PIPELINE & LEADS
          Pipeline Summary and Leads by Source
      ========================================================================= */}
      <section className="space-y-3 sm:space-y-4 md:space-y-6">
        <div className="flex items-center gap-1.5 sm:gap-2">
          <div className="h-1 w-1 rounded-full bg-primary" />
          <h2 className="text-sm sm:text-base md:text-lg font-medium text-muted-foreground uppercase tracking-wide">
            Pipeline & Leads
          </h2>
        </div>

        <StaggerContainer className="grid gap-3 sm:gap-4 md:gap-6 grid-cols-1 lg:grid-cols-2">
          <Suspense fallback={<Skeleton className="h-48 sm:h-64" />}>
            <PipelineSummary />
          </Suspense>

          <Suspense fallback={<Skeleton className="h-48 sm:h-64" />}>
            <LeadsBySource />
          </Suspense>
        </StaggerContainer>
      </section>

      {/* ========================================================================
          SECTION 3: TASKS
          Upcoming Tasks
      ========================================================================= */}
      <section className="space-y-3 sm:space-y-4 md:space-y-6">
        <div className="flex items-center gap-1.5 sm:gap-2">
          <div className="h-1 w-1 rounded-full bg-primary" />
          <h2 className="text-sm sm:text-base md:text-lg font-medium text-muted-foreground uppercase tracking-wide">
            Tasks
          </h2>
        </div>

        {/* Today Tasks */}
        <Suspense fallback={<Skeleton className="h-48 sm:h-64 w-full" />}>
          <UpcomingTasks />
        </Suspense>
      </section>
    </div>
  );
}

// ============================================================================
// Sales Dashboard Widget Components
// ============================================================================

function SalesTodayTasksCard({ startDate, endDate }: Readonly<{ startDate?: string; endDate?: string }>) {
  const params = { start_date: startDate, end_date: endDate };
  const { data, isLoading, isError } = useSalesTodayTasks(params);
  console.log("SalesTodayTasksCard Params:", params);
  console.log("SalesTodayTasksCard Data:", data);

  if (isLoading) {
    return (
      <Card>
        <CardHeader className="pb-2 flex flex-row items-center justify-between space-y-0">
          <Skeleton className="h-4 w-24" />
          <Calendar className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <Skeleton className="h-8 w-16 mb-2" />
          <Skeleton className="h-3 w-32" />
        </CardContent>
      </Card>
    );
  }

  if (isError) {
    return (
      <Card>
        <CardHeader className="pb-2 flex flex-row items-center justify-between space-y-0">
          <CardTitle className="text-sm font-medium text-muted-foreground">
            Today Tasks
          </CardTitle>
          <Calendar className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-sm text-destructive">Failed to load</div>
        </CardContent>
      </Card>
    );
  }

  const tasks = data?.data;
  const total = tasks?.total ?? 0;
  const completed = tasks?.completed ?? 0;
  const overdue = tasks?.overdue ?? 0;

  return (
    <Card className="border-0 shadow-sm">
      <CardHeader className="pb-1.5 sm:pb-2 flex flex-row items-center justify-between space-y-0 px-3 sm:px-6 pt-3 sm:pt-6">
        <CardTitle className="text-xs sm:text-sm font-medium text-muted-foreground">
          Today Tasks
        </CardTitle>
        <Calendar className="h-3 w-3 sm:h-4 sm:w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent className="px-3 sm:px-6 pb-3 sm:pb-6">
        <div className="text-xl sm:text-2xl font-medium">{total}</div>
        <p className="text-[10px] sm:text-xs text-muted-foreground mt-0.5 sm:mt-1">
          Completed: {completed} | Overdue: {overdue}
        </p>
      </CardContent>
    </Card>
  );
}

function SalesAssignedLeadsCard({ startDate, endDate }: Readonly<{ startDate?: string; endDate?: string }>) {
  const params = { limit: 5, start_date: startDate, end_date: endDate };
  const { data, isLoading, isError } = useSalesAssignedLeads(params);
  console.log("SalesAssignedLeadsCard Params:", params);
  console.log("SalesAssignedLeadsCard Data:", data);

  if (isLoading) {
    return (
      <Card>
        <CardHeader className="pb-2 flex flex-row items-center justify-between space-y-0">
          <Skeleton className="h-4 w-28" />
          <UserPlus className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <Skeleton className="h-8 w-16 mb-2" />
          <Skeleton className="h-3 w-32" />
        </CardContent>
      </Card>
    );
  }

  if (isError) {
    return (
      <Card>
        <CardHeader className="pb-2 flex flex-row items-center justify-between space-y-0">
          <CardTitle className="text-sm font-medium text-muted-foreground">
            Assigned Leads
          </CardTitle>
          <UserPlus className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-sm text-destructive">Failed to load</div>
        </CardContent>
      </Card>
    );
  }

  const leads = data?.data;
  const total = leads?.total ?? 0;
  const newCount = leads?.new ?? 0;
  const qualified = leads?.qualified ?? 0;

  return (
    <Card className="border-0 shadow-sm">
      <CardHeader className="pb-1.5 sm:pb-2 flex flex-row items-center justify-between space-y-0 px-3 sm:px-6 pt-3 sm:pt-6">
        <CardTitle className="text-xs sm:text-sm font-medium text-muted-foreground">
          Assigned Leads
        </CardTitle>
        <UserPlus className="h-3 w-3 sm:h-4 sm:w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent className="px-3 sm:px-6 pb-3 sm:pb-6">
        <div className="text-xl sm:text-2xl font-medium">{total}</div>
        <p className="text-[10px] sm:text-xs text-muted-foreground mt-0.5 sm:mt-1">
          New: {newCount} | Qualified: {qualified}
        </p>
      </CardContent>
    </Card>
  );
}

function SalesUpcomingVisitsCard({ startDate, endDate }: Readonly<{ startDate?: string; endDate?: string }>) {
  const params = { limit: 5, start_date: startDate, end_date: endDate };
  const { data, isLoading, isError } = useSalesUpcomingVisits(params);
  console.log("SalesUpcomingVisitsCard Params:", params);
  console.log("SalesUpcomingVisitsCard Data:", data);

  if (isLoading) {
    return (
      <Card>
        <CardHeader className="pb-2 flex flex-row items-center justify-between space-y-0">
          <Skeleton className="h-4 w-32" />
          <Calendar className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <Skeleton className="h-8 w-16 mb-2" />
          <Skeleton className="h-3 w-32" />
        </CardContent>
      </Card>
    );
  }

  if (isError) {
    return (
      <Card>
        <CardHeader className="pb-2 flex flex-row items-center justify-between space-y-0">
          <CardTitle className="text-sm font-medium text-muted-foreground">
            Upcoming Visits
          </CardTitle>
          <Calendar className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-sm text-destructive">Failed to load</div>
        </CardContent>
      </Card>
    );
  }

  const visits = data?.data;
  const total = visits?.total ?? 0;
  const today = visits?.today ?? 0;
  const thisWeek = visits?.this_week ?? 0;

  return (
    <Card className="border-0 shadow-sm">
      <CardHeader className="pb-1.5 sm:pb-2 flex flex-row items-center justify-between space-y-0 px-3 sm:px-6 pt-3 sm:pt-6">
        <CardTitle className="text-xs sm:text-sm font-medium text-muted-foreground">
          Upcoming Visits
        </CardTitle>
        <Calendar className="h-3 w-3 sm:h-4 sm:w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent className="px-3 sm:px-6 pb-3 sm:pb-6">
        <div className="text-xl sm:text-2xl font-medium">{total}</div>
        <p className="text-[10px] sm:text-xs text-muted-foreground mt-0.5 sm:mt-1">
          Today: {today} | This Week: {thisWeek}
        </p>
      </CardContent>
    </Card>
  );
}

function SalesRemindersCard({ startDate, endDate }: Readonly<{ startDate?: string; endDate?: string }>) {
  const params = { limit: 5, start_date: startDate, end_date: endDate };
  const { data, isLoading, isError } = useSalesReminders(params);
  console.log("SalesRemindersCard Params:", params);
  console.log("SalesRemindersCard Data:", data);

  if (isLoading) {
    return (
      <Card>
        <CardHeader className="pb-2 flex flex-row items-center justify-between space-y-0">
          <Skeleton className="h-4 w-20" />
          <Bell className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <Skeleton className="h-8 w-16 mb-2" />
          <Skeleton className="h-3 w-32" />
        </CardContent>
      </Card>
    );
  }

  if (isError) {
    return (
      <Card>
        <CardHeader className="pb-2 flex flex-row items-center justify-between space-y-0">
          <CardTitle className="text-sm font-medium text-muted-foreground">
            Reminders
          </CardTitle>
          <Bell className="h-4 w-4 text-muted-foreground" />
        </CardHeader>
        <CardContent>
          <div className="text-sm text-destructive">Failed to load</div>
        </CardContent>
      </Card>
    );
  }

  const reminders = data?.data;
  const total = reminders?.total ?? 0;
  const unread = reminders?.unread ?? 0;

  return (
    <Card className="border-0 shadow-sm">
      <CardHeader className="pb-1.5 sm:pb-2 flex flex-row items-center justify-between space-y-0 px-3 sm:px-6 pt-3 sm:pt-6">
        <CardTitle className="text-xs sm:text-sm font-medium text-muted-foreground">
          Reminders
        </CardTitle>
        <Bell className="h-3 w-3 sm:h-4 sm:w-4 text-muted-foreground" />
      </CardHeader>
      <CardContent className="px-3 sm:px-6 pb-3 sm:pb-6">
        <div className="text-xl sm:text-2xl font-medium">{total}</div>
        <p className="text-[10px] sm:text-xs text-muted-foreground mt-0.5 sm:mt-1">
          Unread: {unread}
        </p>
      </CardContent>
    </Card>
  );
}

