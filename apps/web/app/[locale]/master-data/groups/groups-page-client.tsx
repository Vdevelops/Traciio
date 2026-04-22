"use client";

import { useTranslations } from "next-intl";
import dynamic from "next/dynamic";

// Dynamic import untuk code splitting
const GroupListComponent = dynamic(
  () =>
    import("@/features/master-data/group/components/group-list").then(
      (mod) => ({ default: mod.GroupList }),
    ),
  { loading: () => null }, // Use route-level loading.tsx
);

export function GroupsPageClient() {
  const t = useTranslations("groupManagement.page");

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-medium tracking-tight">{t("title")}</h1>
          <p className="text-muted-foreground mt-1">{t("description")}</p>
        </div>
      </div>

      <GroupListComponent />
    </div>
  );
}

