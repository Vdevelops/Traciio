"use client";

import { Users, Settings, Building2, Tag } from "lucide-react";
import { useTranslations } from "next-intl";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { LeadList } from "./lead-list";
import { LeadStatusList } from "./lead-status-list";
import { IndustryList } from "./industry-list";
import { LeadSourceList } from "./lead-source-list";
import { useHasPermission } from "@/features/master-data/user-management/hooks/useHasPermission";

export function LeadManagement() {
  const hasViewLeadsPermission = useHasPermission("leads.view");
  const hasViewLeadStatusPermission = useHasPermission("leads.status-view");
  const hasViewIndustriesPermission = useHasPermission("leads.industries-view");
  const hasViewLeadSourcesPermission = useHasPermission("leads.sources-view");
  const t = useTranslations("leadManagement.tabs");

  // Determine default tab - use first available tab
  let defaultTab = "leads";
  if (!hasViewLeadsPermission) {
    if (hasViewLeadStatusPermission) defaultTab = "lead-statuses";
    else if (hasViewIndustriesPermission) defaultTab = "industries";
    else if (hasViewLeadSourcesPermission) defaultTab = "lead-sources";
  }

  return (
    <div className="space-y-4 sm:space-y-6">
      <Tabs defaultValue={defaultTab} className="w-full">
        <TabsList className="w-full sm:w-auto overflow-x-auto">
          {hasViewLeadsPermission && (
            <TabsTrigger value="leads" className="gap-1.5 sm:gap-2 text-xs sm:text-sm">
              <Users className="h-3.5 w-3.5 sm:h-4 sm:w-4" />
              <span className="whitespace-nowrap">{t("leads")}</span>
            </TabsTrigger>
          )}
          {hasViewLeadStatusPermission && (
            <TabsTrigger value="lead-statuses" className="gap-1.5 sm:gap-2 text-xs sm:text-sm">
              <Settings className="h-3.5 w-3.5 sm:h-4 sm:w-4" />
              <span className="whitespace-nowrap">{t("leadStatuses")}</span>
            </TabsTrigger>
          )}
          {hasViewIndustriesPermission && (
            <TabsTrigger value="industries" className="gap-1.5 sm:gap-2 text-xs sm:text-sm">
              <Building2 className="h-3.5 w-3.5 sm:h-4 sm:w-4" />
              <span className="whitespace-nowrap">{t("industries")}</span>
            </TabsTrigger>
          )}
          {hasViewLeadSourcesPermission && (
            <TabsTrigger value="lead-sources" className="gap-1.5 sm:gap-2 text-xs sm:text-sm">
              <Tag className="h-3.5 w-3.5 sm:h-4 sm:w-4" />
              <span className="whitespace-nowrap">{t("leadSources")}</span>
            </TabsTrigger>
          )}
        </TabsList>

        {hasViewLeadsPermission && (
          <TabsContent value="leads" className="mt-4 sm:mt-6">
            <LeadList />
          </TabsContent>
        )}

        {hasViewLeadStatusPermission && (
          <TabsContent value="lead-statuses" className="mt-4 sm:mt-6">
            <LeadStatusList />
          </TabsContent>
        )}

        {hasViewIndustriesPermission && (
          <TabsContent value="industries" className="mt-4 sm:mt-6">
            <IndustryList />
          </TabsContent>
        )}

        {hasViewLeadSourcesPermission && (
          <TabsContent value="lead-sources" className="mt-4 sm:mt-6">
            <LeadSourceList />
          </TabsContent>
        )}
      </Tabs>
    </div>
  );
}

