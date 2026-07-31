import React from "react";
import { redirect } from "next/navigation";
import { getCurrentUser } from "@/features/auth/utils/getCurrentUser";

import SalesRepKPIPageClient from "./page.client";
import SalesManagerKPIPageClient from "./manager.client";

export default async function KPIPage() {
  // server-side: get current user from auth util (session/JWT)
  const user = await getCurrentUser();
  if (!user) {
    // If not authenticated, redirect to login
    redirect("/auth/login");
  }

  const role = user.role;
  if (role === "sales_rep") return <SalesRepKPIPageClient />;
  if (role === "sales_manager") return <SalesManagerKPIPageClient />;

  // default: redirect to dashboard for other roles
  redirect("/dashboard");
}
