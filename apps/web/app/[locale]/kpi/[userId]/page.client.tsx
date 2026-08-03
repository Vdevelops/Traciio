"use client";

import React from "react";
import { useRouter } from "@/i18n/routing";

import SalesRepKPIView from "@/features/kpi/components/sales-rep-kpi-view";

interface Props {
  userId: string;
  userName?: string;
  initialStartDate?: string;
  initialEndDate?: string;
}

export default function SalesRepKPIUserPageClient({
  userId,
  userName,
  initialStartDate,
  initialEndDate,
}: Props) {
  const router = useRouter();

  return (
    <SalesRepKPIView
      userId={userId}
      userName={userName}
      initialStartDate={initialStartDate}
      initialEndDate={initialEndDate}
      onBack={() => {
        const params = new URLSearchParams();

        if (initialStartDate) {
          params.set("startDate", initialStartDate);
        }

        if (initialEndDate) {
          params.set("endDate", initialEndDate);
        }

        const query = params.toString();
        router.push(query ? `/kpi?${query}` : "/kpi");
      }}
    />
  );
}