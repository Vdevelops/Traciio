"use client";

import type { ComponentType } from "react";
import { BarChart3, Target, TrendingUp, Users } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import { useLeadAnalytics } from "../hooks/useLeads";

export function LeadAnalytics() {
  const { data, isLoading } = useLeadAnalytics();
  const metrics = data?.data;

  if (isLoading) {
    return (
      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        {Array.from({ length: 4 }).map((_, index) => (
          <div key={index} className="rounded-lg border p-4">
            <Skeleton className="h-4 w-24" />
            <Skeleton className="mt-3 h-8 w-16" />
            <Skeleton className="mt-2 h-3 w-32" />
          </div>
        ))}
      </div>
    );
  }

  const topSource = metrics?.by_source?.[0];

  return (
    <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
      <MetricCard icon={Users} label="Total Leads" value={String(metrics?.total_leads ?? 0)} description="All captured lead records" />
      <MetricCard icon={Target} label="Conversion Rate" value={`${metrics?.conversion_rate ?? 0}%`} description="Converted lead ratio" />
      <MetricCard icon={TrendingUp} label="Average Score" value={String(metrics?.average_score ?? 0)} description="Mean lead score" />
      <MetricCard icon={BarChart3} label="Top Source" value={topSource?.source ?? "-"} description={topSource ? `${topSource.count} leads` : "No source data"} />
    </div>
  );
}

function MetricCard({
  icon: Icon,
  label,
  value,
  description,
}: {
  readonly icon: ComponentType<{ className?: string }>;
  readonly label: string;
  readonly value: string;
  readonly description: string;
}) {
  return (
    <div className="rounded-lg border bg-card p-4 shadow-sm">
      <div className="flex items-center justify-between gap-4">
        <div>
          <p className="text-sm text-muted-foreground">{label}</p>
          <p className="mt-2 text-2xl font-semibold tracking-tight">{value}</p>
          <p className="mt-1 text-xs text-muted-foreground">{description}</p>
        </div>
        <div className="rounded-lg bg-muted p-2 text-muted-foreground">
          <Icon className="h-5 w-5" />
        </div>
      </div>
    </div>
  );
}
