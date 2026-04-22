import { Skeleton } from "@/components/ui/skeleton";

export default function ProductAnalyticsLoading() {
  return (
    <div className="space-y-8 p-6">
      {/* Header */}
      <div className="space-y-2">
        <Skeleton className="h-9 w-64" />
        <Skeleton className="h-4 w-96" />
      </div>

      {/* Stats Cards */}
      <div className="grid gap-4 md:grid-cols-4">
        <Skeleton className="h-32" />
        <Skeleton className="h-32" />
        <Skeleton className="h-32" />
        <Skeleton className="h-32" />
      </div>

      {/* Filters */}
      <Skeleton className="h-10 w-full" />

      {/* Trends Chart */}
      <Skeleton className="h-80 w-full" />

      {/* Podium */}
      <div className="grid grid-cols-3 gap-4">
        <Skeleton className="h-40" />
        <Skeleton className="h-40" />
        <Skeleton className="h-40" />
      </div>

      {/* Table */}
      <Skeleton className="h-96 w-full" />
    </div>
  );
}
