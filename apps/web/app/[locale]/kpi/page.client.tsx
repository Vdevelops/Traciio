"use client";

import React from "react";
import SalesRepKPIView from "@/features/kpi/components/sales-rep-kpi-view";

interface Props {
  initialStartDate?: string;
  initialEndDate?: string;
}

export default function SalesRepKPIPageClient({ initialStartDate, initialEndDate }: Props) {
  return <SalesRepKPIView initialStartDate={initialStartDate} initialEndDate={initialEndDate} />;
}
