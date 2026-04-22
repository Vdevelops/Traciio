"use client";

import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { Skeleton } from "@/components/ui/skeleton";
import type { PipelineReport } from "../types";
import { SalesFunnelTable } from "./sales-funnel-table";
import { SalesFunnelInsights } from "./sales-funnel-insights";
import { useTranslations } from "next-intl";

interface SalesFunnelViewerProps {
  data?: PipelineReport;
  isLoading: boolean;
}

export function SalesFunnelViewer({ data, isLoading }: SalesFunnelViewerProps) {
  const t = useTranslations("reportsFeature.salesFunnelViewer");
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
      {/* Main Tabs */}
      <Tabs defaultValue="table" className="w-full">
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="table">{t("tabTable")}</TabsTrigger>
          <TabsTrigger value="insights">{t("tabInsights")}</TabsTrigger>
        </TabsList>

        <TabsContent value="table" className="mt-6">
          <SalesFunnelTable data={data} />
        </TabsContent>

        <TabsContent value="insights" className="mt-6">
          <SalesFunnelInsights data={data} />
        </TabsContent>
      </Tabs>
    </div>
  );
}

