"use client";

import dynamic from "next/dynamic";
import { Skeleton } from "@/components/ui/skeleton";

interface TerritoryMapDrawerProps {
  initialPolygon?: number[][][];
  color?: string;
  onPolygonChange: (coordinates: number[][][]) => void;
}

// Dynamically import the map component with SSR disabled
const TerritoryMapDrawer = dynamic(
  () => import("./TerritoryMapDrawer").then((mod) => mod.TerritoryMapDrawer),
  {
    ssr: false,
    loading: () => (
      <div className="w-full h-[400px] rounded-lg border">
        <Skeleton className="w-full h-full" />
        <div className="absolute inset-0 flex items-center justify-center">
          <p className="text-sm text-muted-foreground">Loading map...</p>
        </div>
      </div>
    ),
  }
);

export function TerritoryMapDrawerWrapper(props: TerritoryMapDrawerProps) {
  return <TerritoryMapDrawer {...props} />;
}
