"use client";

import React from "react";
import type { TargetGapItem } from "../types";
import { formatCurrency, formatPercent } from "../utils/formatters";
import { Target, TrendingUp, TrendingDown, CheckCircle2 } from "lucide-react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

interface Props {
  revenue: TargetGapItem | null;
  deals: TargetGapItem | null;
}

export default function TargetGapSummary({ revenue, deals }: Props) {
  const renderGapSentence = (type: "Revenue" | "Deals", item: TargetGapItem) => {
    if (item.status === "unknown") {
      return (
        <span className="text-sm font-semibold text-muted-foreground inline-flex items-center gap-1">
          <Target className="h-4 w-4" />
          Target {type.toLowerCase()} belum tersedia
        </span>
      );
    }

    const gapText = formatPercent(Math.abs(item.gapPercent ?? 0));
    
    if (item.status === "above") {
      return (
        <span className="text-sm font-semibold text-emerald-600 dark:text-emerald-400 inline-flex items-center gap-1">
          <TrendingUp className="h-4 w-4" />
          {type} {gapText} di atas target
        </span>
      );
    }
    if (item.status === "met") {
      return (
        <span className="text-sm font-semibold text-sky-600 dark:text-sky-400 inline-flex items-center gap-1">
          <CheckCircle2 className="h-4 w-4" />
          {type} mencapai target
        </span>
      );
    }
    return (
      <span className="text-sm font-semibold text-rose-600 dark:text-rose-400 inline-flex items-center gap-1">
        <TrendingDown className="h-4 w-4" />
        {type} {gapText} di bawah target
      </span>
    );
  };

  const getStatusColorClass = (status: TargetGapItem["status"]) => {
    if (status === "above") return "bg-emerald-500";
    if (status === "met") return "bg-sky-500";
    if (status === "unknown") return "bg-muted-foreground/40";
    return "bg-rose-500";
  };

  const renderCard = (type: "Revenue" | "Deals", item: TargetGapItem | null) => {
    if (!item) return null;

    const formattedActual =
      type === "Revenue" ? formatCurrency(item.actual) : item.actual.toString();
    const formattedTarget =
      type === "Revenue" ? formatCurrency(item.target) : item.target.toString();
    const rawProgress = item.status !== "unknown" && item.target > 0 ? (item.actual / item.target) * 100 : 0;
    const progressPercent = Math.min(100, Math.max(0, rawProgress));

    return (
      <div className="flex flex-col gap-3 rounded-lg border bg-card p-4">
        <div className="flex items-center justify-between">
          <span className="text-xs font-medium uppercase tracking-wide text-muted-foreground">
            Target {type}
          </span>
          <Target className="h-4 w-4 text-muted-foreground" />
        </div>

        <div className="flex flex-col gap-1">
          <div className="flex items-baseline justify-between">
            <span className="text-xl font-semibold text-foreground">
              {formattedActual}
            </span>
            <span className="text-xs font-medium text-muted-foreground">
              Target: {formattedTarget}
            </span>
          </div>

          <div className="mt-1">
            {renderGapSentence(type, item)}
          </div>
        </div>

        <div className="relative mt-1 h-2 w-full overflow-hidden rounded-full bg-muted">
          <div
            className={`h-full rounded-full transition-all duration-500 ${getStatusColorClass(item.status)}`}
            style={{ width: `${progressPercent}%` }}
          />
        </div>
      </div>
    );
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Pencapaian Target</CardTitle>
        <CardDescription>Perbandingan aktual dengan target pada periode yang dipilih.</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="grid gap-4 md:grid-cols-2">
          {revenue ? renderCard("Revenue", revenue) : (
            <div className="flex items-center justify-center rounded-lg border border-dashed p-6 text-sm text-muted-foreground">
              Target revenue belum ditentukan
            </div>
          )}
          {deals ? renderCard("Deals", deals) : (
            <div className="flex items-center justify-center rounded-lg border border-dashed p-6 text-sm text-muted-foreground">
              Target deals belum ditentukan
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
