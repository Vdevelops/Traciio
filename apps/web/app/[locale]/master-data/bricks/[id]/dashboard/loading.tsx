import { Skeleton } from "@/components/ui/skeleton";

export default function BrickDashboardLoading() {
  return (
    <div className="space-y-6 p-4 sm:p-6">
      <Skeleton className="h-8 w-64 mb-4" />
      <div className="grid gap-4 grid-cols-2 md:grid-cols-4">
        {Array.from({ length: 8 }).map((_, i) => (
          <Skeleton key={i} className="h-32" />
        ))}
      </div>
    </div>
  );
}

