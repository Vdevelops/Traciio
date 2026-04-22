import { Skeleton } from "@/components/ui/skeleton";

export default function DashboardLoading() {
  return (
    <div className="space-y-6 p-4">
      {/* Header skeleton */}
      <div className="space-y-2">
        <Skeleton className="h-9 w-48" />
        <Skeleton className="h-5 w-64" />
      </div>

      {/* Overview cards skeleton */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {Array.from({ length: 4 }).map((_, i) => (
          <Skeleton key={i} className="h-32" />
        ))}
      </div>

      {/* Grid section skeleton */}
      <div className="grid gap-6 lg:grid-cols-3">
        {Array.from({ length: 3 }).map((_, i) => (
          <Skeleton key={i} className="h-64" />
        ))}
      </div>

      {/* Table skeleton */}
      <Skeleton className="h-96 w-full" />

      {/* Statistics skeleton */}
      <Skeleton className="h-80 w-full" />

      {/* Trends skeleton */}
      <Skeleton className="h-80 w-full" />

      {/* Grid section skeleton */}
      <div className="grid gap-6 md:grid-cols-2">
        {Array.from({ length: 2 }).map((_, i) => (
          <Skeleton key={i} className="h-64" />
        ))}
      </div>

      {/* Activities skeleton */}
      <Skeleton className="h-64 w-full" />
    </div>
  );
}

