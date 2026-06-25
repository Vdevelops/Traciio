"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { BarChart3, Target, TrendingUp } from "lucide-react";
import { BrickPipelineAnalytics } from "./brick-pipeline-analytics";
import { BrickVisitAnalytics } from "./brick-visit-analytics";
import { BrickTargetAnalytics } from "./brick-target-analytics";

interface BrickAnalyticsTabsProps {
  readonly brickId: string;
  readonly periodStart?: string;
  readonly periodEnd?: string;
}

export function BrickAnalyticsTabs({ brickId, periodStart, periodEnd }: BrickAnalyticsTabsProps) {
  const t = useTranslations("brickAnalytics");
  const [activeTab, setActiveTab] = useState("pipeline");

  return (
    <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
      <TabsList className="grid w-full grid-cols-3">
        <TabsTrigger value="pipeline" className="gap-2">
          <BarChart3 className="h-4 w-4" />
          {t("pipelineTab")}
        </TabsTrigger>
        <TabsTrigger value="visits" className="gap-2">
          <TrendingUp className="h-4 w-4" />
          {t("visitsTab")}
        </TabsTrigger>
        <TabsTrigger value="target" className="gap-2">
          <Target className="h-4 w-4" />
          {t("targetTab")}
        </TabsTrigger>
      </TabsList>

      <TabsContent value="pipeline" className="mt-6">
        <BrickPipelineAnalytics brickId={brickId} periodStart={periodStart} periodEnd={periodEnd} />
      </TabsContent>

      <TabsContent value="visits" className="mt-6">
        <BrickVisitAnalytics brickId={brickId} periodStart={periodStart} periodEnd={periodEnd} />
      </TabsContent>

      <TabsContent value="target" className="mt-6">
        <BrickTargetAnalytics brickId={brickId} periodStart={periodStart} periodEnd={periodEnd} />
      </TabsContent>
    </Tabs>
  );
}
