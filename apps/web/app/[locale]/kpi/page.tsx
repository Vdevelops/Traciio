import React from "react";
import { redirect } from "next/navigation";
import { getCurrentUser } from "@/features/auth/utils/getCurrentUser";

import SalesRepKPIPageClient from "./page.client";
import SalesManagerKPIPageClient from "./manager.client";

export default async function KPIPage({
  params,
  searchParams,
}: {
  params: Promise<{ locale: string }>;
  searchParams?: Promise<{
    startDate?: string;
    endDate?: string;
  }>;
}) {
  const { locale } = await params;
  const query = searchParams ? await searchParams : {};
  // server-side: get current user from auth util (session/JWT)
  const user = await getCurrentUser();
  if (!user) {
    // If not authenticated, redirect to login
    redirect(`/${locale}/login`);
  }

  const role = user.role;
  if (role === "sales_rep" || role === "sales") {
    return <SalesRepKPIPageClient initialStartDate={query.startDate} initialEndDate={query.endDate} />;
  }
  if (role === "sales_manager" || role === "admin") {
    return <SalesManagerKPIPageClient initialStartDate={query.startDate} initialEndDate={query.endDate} />;
  }

  // default: redirect to dashboard for other roles
  redirect(`/${locale}/dashboard`);
}
