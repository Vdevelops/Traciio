"use client";

import { Grid3X3 } from "lucide-react";

export default function BricksLoading() {
  return (
    <div className="absolute inset-0 flex items-center justify-center bg-muted/30">
      <div className="flex flex-col items-center gap-4">
        <div className="relative w-20 h-20">
          <div className="absolute inset-0 border-4 border-primary/20 rounded-full" />
          <div className="absolute inset-0 border-4 border-primary border-t-transparent rounded-full animate-spin" />
          <Grid3X3 className="absolute inset-0 m-auto w-8 h-8 text-primary/60" />
        </div>
        <div className="text-center">
          <p className="text-sm font-medium text-foreground">Loading Bricks</p>
          <p className="text-xs text-muted-foreground">Preparing territory map...</p>
        </div>
      </div>
    </div>
  );
}
