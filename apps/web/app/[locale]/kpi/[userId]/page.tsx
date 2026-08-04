import React from "react";
import { redirect } from "next/navigation";
import { getCurrentUser } from "@/features/auth/utils/getCurrentUser";

import SalesRepKPIUserPageClient from "./page.client";

interface Props {
  params: Promise<{
    locale: string;
    userId: string;
  }>;
  searchParams?: Promise<{
    startDate?: string;
    endDate?: string;
    userName?: string;
  }>;
}

export default async function SalesRepKPIDrilldownPage({ params, searchParams }: Props) {
  const { locale, userId } = await params;
  const query = searchParams ? await searchParams : {};
  const user = await getCurrentUser();

  if (!user) {
    redirect(`/${locale}/login`);
  }

  if (user.role !== "sales_manager" && user.role !== "admin") {
    redirect(`/${locale}/dashboard`);
  }

  return (
    <SalesRepKPIUserPageClient
      userId={userId}
      userName={query.userName}
      initialStartDate={query.startDate}
      initialEndDate={query.endDate}
    />
  );
}