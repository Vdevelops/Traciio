"use client";

import { Skeleton } from "@/components/ui/skeleton";
import type { PipelineReport } from "../types";
import { SalesFunnelTable } from "./sales-funnel-table";
import { useTranslations } from "next-intl";

interface SalesFunnelViewerProps {
  data?: PipelineReport;
  isLoading: boolean;
}

export function SalesFunnelViewer({ data, isLoading }: SalesFunnelViewerProps) {
  const tCommon = useTranslations("reportsFeature.common");
  if (isLoading) {
    return (
      <div className="space-y-4">
        <Skeleton className="h-32 w-full" />
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  if (!data) {
    return (
      <div className="text-center py-8 text-muted-foreground">
        <p className="text-sm">{tCommon("noData")}</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <SalesFunnelTable data={data} />
    </div>
  );
}
