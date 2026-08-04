"use client";

import React from "react";
import type { TeamSummary } from "../types";
import { formatCurrency, formatPercent } from "../utils/formatters";
import { Users, DollarSign, Target, ClipboardList, AlertCircle, CheckCircle } from "lucide-react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

interface Props {
  summary: TeamSummary | null;
}

export default function TeamSummaryCard({ summary }: Props) {
  if (!summary) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Ringkasan Performa Tim</CardTitle>
          <CardDescription>Agregasi performa sales pada brick yang berada dalam scope manager.</CardDescription>
        </CardHeader>
        <CardContent className="p-6">
          <div className="rounded-lg border border-dashed bg-muted/20 p-6 text-sm text-muted-foreground">
            Belum ada data ringkasan tim untuk periode ini.
          </div>
        </CardContent>
      </Card>
    );
  }

  const renderValue = (val: number | null, isCurrency = false, isPercent = false) => {
    if (val === null || val === undefined) {
      return (
        <span className="text-sm font-medium text-muted-foreground">
          Belum ada data
        </span>
      );
    }
    if (isCurrency) {
      return <span className="text-2xl font-semibold text-foreground">{formatCurrency(val)}</span>;
    }
    if (isPercent) {
      return <span className="text-2xl font-semibold text-foreground">{formatPercent(val)}</span>;
    }
    return <span className="text-2xl font-semibold text-foreground">{val.toLocaleString()}</span>;
  };

  const items = [
    {
      label: "Total Sales Reps",
      value: renderValue(summary.totalRepsCount),
      icon: Users,
      color: "text-blue-600 bg-blue-50 dark:bg-blue-950/20",
    },
    {
      label: "Total Revenue",
      value: renderValue(summary.totalRevenue, true),
      icon: DollarSign,
      color: "text-emerald-600 bg-emerald-50 dark:bg-emerald-950/20",
    },
    {
      label: "Total Deals Closed",
      value: renderValue(summary.totalDealsClosed),
      icon: Target,
      color: "text-violet-600 bg-violet-50 dark:bg-violet-950/20",
    },
    {
      label: "Rata-rata Konversi",
      value: renderValue(summary.teamConversionRate, false, true),
      icon: ClipboardList,
      color: "text-sky-600 bg-sky-50 dark:bg-sky-950/20",
    },
    {
      label: "Kepatuhan Kunjungan",
      value: renderValue(summary.teamVisitCompliance, false, true),
      icon: CheckCircle,
      color: "text-teal-600 bg-teal-50 dark:bg-teal-950/20",
    },
    {
      label: "Rasio Tugas Overdue",
      value: renderValue(summary.teamOverdueTaskRate, false, true),
      icon: AlertCircle,
      color: "text-rose-600 bg-rose-50 dark:bg-rose-950/20",
    },
  ];

  return (
    <Card>
      <CardHeader>
        <CardTitle>Ringkasan Performa Tim</CardTitle>
        <CardDescription>Agregasi performa sales pada brick yang berada dalam scope manager.</CardDescription>
      </CardHeader>
      <CardContent>
        <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-3">
        {items.map((item, idx) => {
          const Icon = item.icon;
          return (
            <div
              key={idx}
              className="flex items-center gap-4 rounded-lg border bg-card p-4"
            >
              <div className={`rounded-md p-2.5 ${item.color}`}>
                <Icon className="h-5 w-5" />
              </div>
              <div className="flex flex-col">
                <span className="text-xs font-medium uppercase tracking-wide text-muted-foreground">
                  {item.label}
                </span>
                <div className="mt-1">{item.value}</div>
              </div>
            </div>
          );
        })}
        </div>
      </CardContent>
    </Card>
  );
}
