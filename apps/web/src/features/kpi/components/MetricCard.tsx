"use client";

import React from "react";

interface Props {
  label: string;
  value: number | null | string;
  suffix?: string;
}

export default function MetricCard({ label, value, suffix }: Props) {
  const renderValue = () => {
    if (value === null || value === undefined)
      return (
        <span className="text-sm text-gray-500" role="status" aria-label="no-data">Belum ada data</span>
      );
    if (typeof value === "number")
      return (
        <span className="text-lg font-semibold" aria-label={`metric-value-${label.replace(/\s+/g, "-").toLowerCase()}`}>{value.toLocaleString()}{suffix ? ` ${suffix}` : ""}</span>
      );
    return <span className="text-lg font-semibold">{value}</span>;
  };

  return (
    <figure className="p-3 bg-white rounded-md shadow-sm" role="group" aria-labelledby={`metric-${label.replace(/\s+/g, "-").toLowerCase()}`}>
      <figcaption id={`metric-${label.replace(/\s+/g, "-").toLowerCase()}`} className="text-xs text-gray-400">{label}</figcaption>
      <div className="mt-2">{renderValue()}</div>
    </figure>
  );
}
