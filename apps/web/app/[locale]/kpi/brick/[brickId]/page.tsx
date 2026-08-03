import React from "react";
import { redirect } from "next/navigation";
import { getCurrentUser } from "@/features/auth/utils/getCurrentUser";
import SalesManagerKPIPageClient from "../../manager.client";

interface Props {
  params: Promise<{
    locale: string;
    brickId: string;
  }>;
  searchParams?: Promise<{
    startDate?: string;
    endDate?: string;
  }>;
}

export default async function BrickKPIPage({ params, searchParams }: Props) {
  const { brickId, locale } = await params;
  const query = searchParams ? await searchParams : {};
  const user = await getCurrentUser();
  if (!user) {
    redirect(`/${locale}/login`);
  }

  // Only manager and admin roles are allowed to view brick-level KPI reports
  if (user.role !== "sales_manager" && user.role !== "admin") {
    redirect(`/${locale}/dashboard`);
  }

  return (
    <SalesManagerKPIPageClient
      lockedBrickId={brickId}
      initialStartDate={query.startDate}
      initialEndDate={query.endDate}
    />
  );
}
