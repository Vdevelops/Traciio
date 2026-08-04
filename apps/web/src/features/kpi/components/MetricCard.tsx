"use client";

import React from "react";

interface Props {
  label: string;
  value: number | null | string;
  suffix?: string;
}

export default function MetricCard({ label, value, suffix }: Props) {
  const isNull = value === null || value === undefined;

  const renderValue = () => {
    if (isNull) {
      return (
        <span
          className="text-sm font-medium text-muted-foreground"
          role="status"
          aria-label="no-data"
        >
          Belum ada data
        </span>
      );
    }
    if (typeof value === "number") {
      return (
        <span
          className="text-2xl font-semibold tracking-tight text-foreground"
          aria-label={`metric-value-${label.replace(/\s+/g, "-").toLowerCase()}`}
        >
          {value.toLocaleString()}{suffix ? ` ${suffix}` : ""}
        </span>
      );
    }
    return (
      <span className="text-2xl font-semibold tracking-tight text-foreground">
        {value}
      </span>
    );
  };

  return (
    <div
      className="flex min-h-24 flex-col justify-between rounded-lg border bg-card p-4"
      role="group"
      aria-labelledby={`metric-${label.replace(/\s+/g, "-").toLowerCase()}`}
    >
      <span
        id={`metric-${label.replace(/\s+/g, "-").toLowerCase()}`}
        className="text-xs font-medium uppercase tracking-wide text-muted-foreground"
      >
        {label}
      </span>
      <div className="mt-3 flex items-baseline justify-between">
        {renderValue()}
      </div>
    </div>
  );
}
