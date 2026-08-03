"use client";

import React, { useMemo } from "react";
import type { DiagnosticFlag } from "../types";
import { getDiagnosticSeverityDetails } from "../utils/formatters";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

interface Props {
  items: DiagnosticFlag[];
  onBrickClick?: (brickId: string) => void;
}

export default function DiagnosticList({ items, onBrickClick }: Props) {
  // Group or sort by severity: critical first, then warning, then info
  const sortedItems = useMemo(() => {
    if (!items) return [];
    const severityWeight = {
      critical: 3,
      warning: 2,
      info: 1,
    };
    return [...items].sort((a, b) => {
      const weightA = severityWeight[a.severity] || 0;
      const weightB = severityWeight[b.severity] || 0;
      return weightB - weightA;
    });
  }, [items]);

  if (sortedItems.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle id="kpi-diagnostics">Diagnostic & Actionable Insights</CardTitle>
          <CardDescription>Catatan evaluasi otomatis berdasarkan rule KPI.</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="rounded-lg border border-dashed bg-muted/20 p-6 text-sm text-muted-foreground">
            Tidak ada diagnostic untuk periode ini.
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle id="kpi-diagnostics">Diagnostic & Actionable Insights</CardTitle>
        <CardDescription>Catatan evaluasi otomatis berdasarkan rule KPI.</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="space-y-3" aria-live="polite">
          {sortedItems.map((flag, idx) => {
            const details = getDiagnosticSeverityDetails(flag.severity);
            const Icon = details.icon;

            return (
              <div
                key={`${flag.code}-${flag.brickId || ""}-${idx}`}
                className={`flex items-start gap-3 rounded-lg border p-4 ${details.colorClass}`}
              >
                <div className="mt-0.5 flex-shrink-0">
                  <Icon className="h-5 w-5" />
                </div>
                <div className="flex-1">
                  <div className="flex flex-wrap items-center gap-2">
                    <span className="text-xs font-semibold uppercase tracking-wide">
                      {flag.code.replace(/_/g, " ")}
                    </span>
                    {flag.brickId && (
                      <button
                        onClick={() => {
                          if (flag.brickId) {
                            onBrickClick?.(flag.brickId);
                          }
                        }}
                        className="inline-flex items-center rounded-md bg-muted px-2 py-0.5 text-[10px] font-medium text-muted-foreground transition-colors hover:bg-muted/80"
                      >
                        Brick: {flag.brickId}
                      </button>
                    )}
                  </div>
                  <p className="mt-1 text-sm leading-relaxed opacity-90">
                    {flag.message}
                  </p>
                </div>
              </div>
            );
          })}
        </div>
      </CardContent>
    </Card>
  );
}
