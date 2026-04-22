import { getRequestConfig } from "next-intl/server";
import { routing } from "./routing";
import type { Locale } from "@/features/dashboard/types";

export default getRequestConfig(async ({ requestLocale }) => {
  let locale = await requestLocale;

  if (!locale || !routing.locales.includes(locale as Locale)) {
    locale = routing.defaultLocale;
  }

  const [
    dashboardMessages,
    userManagementMessages,
    groupManagementMessages,
    monthlyTargetMessages,
    brickManagementMessages,
    reportsMessages,
    accountManagementMessages,
    leadManagementMessages,
    pipelineManagementMessages,
    productManagementMessages,
    taskManagementMessages,
    areaMappingMessages,
    visitReportMessages,
    scheduleManagementMessages,
    aiMessages,
    profileMessages,
    notificationMessages,
    routeOptimizationMessages,
    salesOverviewMessages,
    productAnalyticsMessages,
  ] = await Promise.all([
    import(`@/features/dashboard/i18n/messages/${locale}.json`),
    import(`@/features/master-data/user-management/i18n/messages/${locale}.json`),
    import(`@/features/master-data/group/i18n/messages/${locale}.json`),
    import(`@/features/master-data/monthly-target/i18n/messages/${locale}.json`),
    import(`@/features/master-data/brick/i18n/messages/${locale}.json`),
    import(`@/features/reports/i18n/messages/${locale}.json`),
    import(`@/features/sales-crm/account-management/i18n/messages/${locale}.json`),
    import(`@/features/sales-crm/lead-management/i18n/messages/${locale}.json`),
    import(`@/features/sales-crm/pipeline-management/i18n/messages/${locale}.json`),
    import(`@/features/sales-crm/product-management/i18n/messages/${locale}.json`),
    import(`@/features/sales-crm/task-management/i18n/messages/${locale}.json`),
    import(`@/features/sales-crm/area-mapping/i18n/messages/${locale}.json`),
    import(`@/features/sales-crm/visit-report/i18n/messages/${locale}.json`),
    import(`@/features/sales-crm/schedule-management/i18n/messages/${locale}.json`),
    import(`@/features/ai/i18n/messages/${locale}.json`),
    import(`@/features/profile/i18n/messages/${locale}.json`),
    import(`@/features/notifications/i18n/messages/${locale}.json`),
    import(`@/features/route-optimization/i18n/messages/${locale}.json`),
    import(`@/features/sales-overview/i18n/messages/${locale}.json`),
    import(`@/features/product-analytics/i18n/messages/${locale}.json`),
  ]);

  const messages = {
    ...dashboardMessages.default,
    ...userManagementMessages.default,
    ...groupManagementMessages.default,
    ...monthlyTargetMessages.default,
    ...brickManagementMessages.default,
    ...reportsMessages.default,
    ...accountManagementMessages.default,
    ...leadManagementMessages.default,
    ...pipelineManagementMessages.default,
    ...productManagementMessages.default,
    ...taskManagementMessages.default,
    ...areaMappingMessages.default,
    ...visitReportMessages.default,
    ...scheduleManagementMessages.default,
    ...aiMessages.default,
    ...profileMessages.default,
    ...notificationMessages.default,
    ...routeOptimizationMessages.default,
    ...salesOverviewMessages.default,
    ...productAnalyticsMessages.default,
  };

  return {
    locale,
    messages,
  };
});


