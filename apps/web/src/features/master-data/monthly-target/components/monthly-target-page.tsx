"use client";

import { PageMotion } from "@/components/motion";
import { useTranslations } from "next-intl";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { TargetMatrix } from "./target-matrix";

export function MonthlyTargetPage() {
  const t = useTranslations("monthlyTargetManagement.page");
  const tp = useTranslations("monthlyTargetManagement.planner");
  return (
    <PageMotion className="p-6">
      <div className="flex flex-col gap-6">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">{t("title")}</h2>
          <p className="text-muted-foreground">{t("description")}</p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>{tp("teamMatrixTitle", { year: new Date().getFullYear() })}</CardTitle>
            <CardDescription>{tp("teamMatrixDescription")}</CardDescription>
          </CardHeader>
          <CardContent>
            <TargetMatrix />
          </CardContent>
        </Card>
      </div>
    </PageMotion>
  );
}
