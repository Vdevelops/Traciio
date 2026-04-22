"use client";

import { Settings, Calendar, LayoutGrid, List } from "lucide-react";
import { useTranslations } from "next-intl";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { VisitReportList } from "./visit-report-list";
import { VisitReportCalendar } from "./visit-report-calendar";
import { VisitReportTeamOverview } from "./visit-report-team-overview";
import { ActivityTypeList } from "./activity-type-list";
import { useHasPermission } from "@/features/master-data/user-management/hooks/useHasPermission";

export function VisitReportManagement() {
  const hasVisitReportsPermission = useHasPermission("visit-reports.view");
  const hasActivityPermission = useHasPermission("visit-reports.activity-type");
  const t = useTranslations("visitReportManagement.tabs");

  // Card view is the default for a better overview experience
  const defaultTab = hasVisitReportsPermission ? "cards" : "activity-types";

  return (
    <Tabs defaultValue={defaultTab} className="w-full">
      <TabsList>
        {hasVisitReportsPermission && (
          <>
            <TabsTrigger value="cards" className="gap-2">
              <LayoutGrid className="h-4 w-4" />
              {t("cards")}
            </TabsTrigger>
            <TabsTrigger value="list" className="gap-2">
              <List className="h-4 w-4" />
              {t("list")}
            </TabsTrigger>
            <TabsTrigger value="calendar" className="gap-2">
              <Calendar className="h-4 w-4" />
              {t("calendar")}
            </TabsTrigger>
          </>
        )}
        {hasActivityPermission && (
          <TabsTrigger value="activity-types" className="gap-2">
            <Settings className="h-4 w-4" />
            {t("activityTypes")}
          </TabsTrigger>
        )}
      </TabsList>

      {hasVisitReportsPermission && (
        <>
          <TabsContent value="cards" className="mt-6">
            <VisitReportTeamOverview />
          </TabsContent>
          <TabsContent value="list" className="mt-6">
            <VisitReportList />
          </TabsContent>
          <TabsContent value="calendar" className="mt-6">
            <VisitReportCalendar />
          </TabsContent>
        </>
      )}

      {hasActivityPermission && (
        <TabsContent value="activity-types" className="mt-6">
          <ActivityTypeList />
        </TabsContent>
      )}
    </Tabs>
  );
}

