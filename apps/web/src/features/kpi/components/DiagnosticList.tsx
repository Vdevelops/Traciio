"use client";

import React from "react";
import type { DiagnosticFlag } from "../types";

interface Props {
  items: DiagnosticFlag[];
}

export default function DiagnosticList({ items }: Props) {
  if (!items || items.length === 0) return null;
  return (
    <section aria-labelledby="kpi-diagnostics" className="p-3 bg-white rounded-md shadow-sm" role="region">
      <h3 id="kpi-diagnostics" className="text-sm font-semibold mb-2">Diagnostics</h3>
      <ul className="space-y-2" aria-live="polite">
        {items.map((d) => (
          <li key={d.code + (d.brickId || "")} className="flex items-start gap-3">
            <div className="w-3 h-3 rounded-full" aria-hidden style={{ background: d.severity === "critical" ? "#ef4444" : d.severity === "warning" ? "#f59e0b" : "#3b82f6" }} />
            <div>
              <div className="text-sm font-medium">{d.code}{d.brickId ? ` • ${d.brickId}` : ""}</div>
              <div className="text-sm text-gray-600">{d.message}</div>
            </div>
          </li>
        ))}
      </ul>
    </section>
  );
}
