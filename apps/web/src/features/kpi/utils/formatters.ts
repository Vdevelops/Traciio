import { ArrowUpRight, ArrowDownRight, Minus, AlertTriangle, AlertCircle, Info } from "lucide-react";
import type { KPIGrade, TrendDirection, DiagnosticSeverity } from "../types";

export const formatCurrency = (value: number | null): string => {
  if (value === null || value === undefined) return "Belum ada data";
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(value);
};

export const formatPercent = (value: number | null, decimals = 1): string => {
  if (value === null || value === undefined) return "N/A";
  return `${value.toFixed(decimals)}%`;
};

export const formatCompositeScore = (value: number | null): string => {
  if (value === null || value === undefined) return "N/A";
  return value.toFixed(1);
};

export const gradeColor = (grade: KPIGrade | null): { text: string; bg: string; border: string; raw: string } => {
  switch (grade) {
    case "Excellent":
      return {
        text: "text-emerald-700 dark:text-emerald-300",
        bg: "bg-emerald-50 dark:bg-emerald-950/30",
        border: "border-emerald-200 dark:border-emerald-800/55",
        raw: "#10b981",
      };
    case "Good":
      return {
        text: "text-sky-700 dark:text-sky-300",
        bg: "bg-sky-50 dark:bg-sky-950/30",
        border: "border-sky-200 dark:border-sky-800/55",
        raw: "#0ea5e9",
      };
    case "Needs Improvement":
      return {
        text: "text-amber-700 dark:text-amber-300",
        bg: "bg-amber-50 dark:bg-amber-950/30",
        border: "border-amber-200 dark:border-amber-800/55",
        raw: "#f59e0b",
      };
    case "Critical":
      return {
        text: "text-rose-700 dark:text-rose-300",
        bg: "bg-rose-50 dark:bg-rose-950/30",
        border: "border-rose-200 dark:border-rose-800/55",
        raw: "#f43f5e",
      };
    default:
      return {
        text: "text-slate-500",
        bg: "bg-slate-50",
        border: "border-slate-200",
        raw: "#64748b",
      };
  }
};

export const trendColor = (direction: TrendDirection | null): string => {
  switch (direction) {
    case "up":
      return "text-emerald-600 dark:text-emerald-400";
    case "down":
      return "text-rose-600 dark:text-rose-400";
    case "flat":
    default:
      return "text-slate-500";
  }
};

export const getTrendIcon = (direction: TrendDirection | null) => {
  switch (direction) {
    case "up":
      return ArrowUpRight;
    case "down":
      return ArrowDownRight;
    case "flat":
    default:
      return Minus;
  }
};

export const getDiagnosticSeverityDetails = (severity: DiagnosticSeverity) => {
  switch (severity) {
    case "critical":
      return {
        icon: AlertCircle,
        colorClass: "text-rose-600 dark:text-rose-400 bg-rose-50 dark:bg-rose-950/20 border-rose-100 dark:border-rose-900/30",
        rawColor: "#f43f5e"
      };
    case "warning":
      return {
        icon: AlertTriangle,
        colorClass: "text-amber-600 dark:text-amber-400 bg-amber-50 dark:bg-amber-950/20 border-amber-100 dark:border-amber-900/30",
        rawColor: "#f59e0b"
      };
    case "info":
    default:
      return {
        icon: Info,
        colorClass: "text-sky-600 dark:text-sky-400 bg-sky-50 dark:bg-sky-950/20 border-sky-100 dark:border-sky-900/30",
        rawColor: "#0ea5e9"
      };
  }
};
