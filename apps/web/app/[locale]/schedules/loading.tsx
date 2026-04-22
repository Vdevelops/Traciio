import { Skeleton } from "@/components/ui/skeleton";

export default function SchedulesLoading() {
  return (
    <div className="space-y-6 p-4">
      <div className="space-y-2">
        <Skeleton className="h-9 w-64" />
        <Skeleton className="h-5 w-96" />
      </div>
      <Skeleton className="h-[600px] w-full" />
    </div>
  );
}

