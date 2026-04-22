"use client";

import { useTranslations } from "next-intl";
import { TerritoryList } from "./TerritoryList";

export function AreaMappingManagement() {
  const t = useTranslations("areaMapping");

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-medium tracking-tight">{t("title")}</h1>
        <p className="text-muted-foreground">{t("description")}</p>
      </div>

      <div className="space-y-4">
        <TerritoryList />
      </div>
    </div>
  );
}
